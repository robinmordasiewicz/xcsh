/**
 * GenAI Eval Subcommand Group
 *
 * Eval mode commands for RBAC testing and validation
 */

import type {
	SubcommandGroup,
	CommandDefinition,
	DomainCommandResult,
} from "../registry.js";
import { successResult, errorResult } from "../registry.js";
import type { REPLSession } from "../../repl/session.js";
import type { OutputFormat } from "../../output/index.js";
import { getCommandSpec, formatSpec } from "../../output/index.js";
import {
	parseDomainOutputFlags,
	formatDomainOutput,
} from "../../output/domain-formatter.js";
import { getGenAIClient } from "./client.js";
import { renderResponse } from "./response-renderer.js";
import { updateLastQueryState, getLastQueryState } from "./query.js";
import { FEEDBACK_TYPE_MAP, getValidFeedbackTypes } from "./types.js";
import type { NegativeFeedbackType } from "./types.js";

/**
 * Parse eval query args
 */
function parseEvalQueryArgs(
	args: string[],
	session: REPLSession,
): {
	format: OutputFormat;
	spec: boolean;
	noColor: boolean;
	namespace: string;
	question: string;
} {
	const { options, remainingArgs } = parseDomainOutputFlags(
		args,
		session.getOutputFormat(),
	);

	let spec = false;
	let namespace = session.getNamespace();
	const questionParts: string[] = [];

	let i = 0;
	while (i < remainingArgs.length) {
		const arg = remainingArgs[i];

		if (arg === "--spec") {
			spec = true;
		} else if (arg === "--namespace" || arg === "-ns") {
			if (i + 1 < remainingArgs.length) {
				namespace = remainingArgs[i + 1] ?? namespace;
				i++;
			}
		} else {
			questionParts.push(arg ?? "");
		}
		i++;
	}

	return {
		format: options.format,
		noColor: options.noColor,
		spec,
		namespace,
		question: questionParts.join(" "),
	};
}

/**
 * Eval Query command - Query using eval endpoint (for RBAC testing)
 */
const evalQueryCommand: CommandDefinition = {
	name: "query",
	description:
		"Send a query to the AI assistant using the eval endpoint. This endpoint is used for RBAC testing and validation purposes, allowing administrators to test permission scenarios without affecting production query analytics.",
	descriptionShort: "Eval mode AI query",
	descriptionMedium:
		"Query the AI assistant in eval mode for RBAC testing and permission validation.",
	usage: "<question> [--namespace <ns>]",
	aliases: ["q"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, namespace, question } =
			parseEvalQueryArgs(args, session);

		// Handle --spec flag
		if (spec) {
			const cmdSpec = getCommandSpec("generative_ai eval query");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		// Validate we have a question
		if (!question.trim()) {
			return errorResult(
				"Please provide a question. Usage: ai eval query <question>",
			);
		}

		// Check for API client
		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult(
				"Not connected to API. Please configure connection first.",
			);
		}

		if (!session.isTokenValidated()) {
			return errorResult(
				"Not authenticated. Please check your API token.",
			);
		}

		try {
			const client = getGenAIClient(apiClient);
			const response = await client.evalQuery(namespace, question);

			// Store query state for feedback (mark as eval)
			updateLastQueryState({
				namespace,
				lastQueryId: response.query_id,
				lastQuery: question,
				followUpQueries: response.follow_up_queries ?? [],
				isActive: true,
			});

			// Handle none format
			if (format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (format === "json" || format === "yaml" || format === "tsv") {
				return successResult(
					formatDomainOutput(response, { format, noColor }),
				);
			}

			// Table/text format - use custom rendering
			const lines = ["[EVAL MODE]", "", ...renderResponse(response)];
			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Eval query failed: ${message}`);
		}
	},
};

/**
 * Parse eval feedback args
 */
function parseEvalFeedbackArgs(
	args: string[],
	session: REPLSession,
): {
	spec: boolean;
	namespace: string;
	positive: boolean;
	negativeType: NegativeFeedbackType | null;
	comment: string | null;
	queryId: string | null;
} {
	const { remainingArgs } = parseDomainOutputFlags(
		args,
		session.getOutputFormat(),
	);

	let spec = false;
	let namespace = session.getNamespace();
	let positive = false;
	let negativeType: NegativeFeedbackType | null = null;
	let comment: string | null = null;
	let queryId: string | null = null;

	let i = 0;
	while (i < remainingArgs.length) {
		const arg = remainingArgs[i];

		if (arg === "--spec") {
			spec = true;
		} else if (arg === "--namespace" || arg === "-ns") {
			if (i + 1 < remainingArgs.length) {
				namespace = remainingArgs[i + 1] ?? namespace;
				i++;
			}
		} else if (arg === "--positive" || arg === "-p") {
			positive = true;
		} else if (arg === "--negative" || arg === "-n") {
			if (i + 1 < remainingArgs.length) {
				const typeArg = remainingArgs[i + 1]?.toLowerCase();
				if (typeArg && FEEDBACK_TYPE_MAP[typeArg]) {
					negativeType = FEEDBACK_TYPE_MAP[typeArg] ?? null;
					i++;
				} else {
					negativeType = "OTHER";
				}
			} else {
				negativeType = "OTHER";
			}
		} else if (arg === "--comment" || arg === "-c") {
			if (i + 1 < remainingArgs.length) {
				comment = remainingArgs[i + 1] ?? null;
				i++;
			}
		} else if (arg === "--query-id" || arg === "-q") {
			if (i + 1 < remainingArgs.length) {
				queryId = remainingArgs[i + 1] ?? null;
				i++;
			}
		}
		i++;
	}

	return {
		spec,
		namespace,
		positive,
		negativeType,
		comment,
		queryId,
	};
}

/**
 * Eval Feedback command - Submit feedback using eval endpoint
 */
const evalFeedbackCommand: CommandDefinition = {
	name: "feedback",
	description:
		"Submit feedback for an eval mode query. Use this endpoint when providing feedback for queries made through the eval endpoint, ensuring proper RBAC testing analytics separation.",
	descriptionShort: "Eval mode feedback submission",
	descriptionMedium:
		"Submit feedback for eval mode AI queries, keeping RBAC testing analytics separate.",
	usage: "--positive | --negative <type> [--comment <text>] [--query-id <id>]",
	aliases: ["fb"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { spec, namespace, positive, negativeType, comment, queryId } =
			parseEvalFeedbackArgs(args, session);

		// Handle --spec flag
		if (spec) {
			const cmdSpec = getCommandSpec("generative_ai eval feedback");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		// Validate we have feedback type
		if (!positive && !negativeType) {
			const validTypes = getValidFeedbackTypes().join(", ");
			return errorResult(
				`Please specify feedback type:\n` +
					`  --positive (-p)       Rate the response positively\n` +
					`  --negative (-n) <type> Rate negatively with reason\n` +
					`\nNegative types: ${validTypes}`,
			);
		}

		// Get query state
		const state = getLastQueryState();
		const targetQueryId = queryId ?? state.lastQueryId;
		const targetQuery = state.lastQuery;

		if (!targetQueryId) {
			return errorResult(
				"No query to provide feedback for. Make an eval query first, or use --query-id to specify one.",
			);
		}

		// Check API client
		const apiClient = session.getAPIClient();
		if (!apiClient) {
			return errorResult(
				"Not connected to API. Please configure connection first.",
			);
		}

		if (!session.isTokenValidated()) {
			return errorResult(
				"Not authenticated. Please check your API token.",
			);
		}

		try {
			const client = getGenAIClient(apiClient);

			if (positive) {
				await client.evalFeedback({
					query: targetQuery ?? "",
					query_id: targetQueryId,
					namespace: state.namespace || namespace,
					positive_feedback: {},
					comment: comment ?? undefined,
				});
				return successResult([
					"[EVAL MODE] Positive feedback submitted successfully.",
					`Query ID: ${targetQueryId}`,
				]);
			}

			// Negative feedback
			await client.evalFeedback({
				query: targetQuery ?? "",
				query_id: targetQueryId,
				namespace: state.namespace || namespace,
				negative_feedback: {
					remarks: negativeType ? [negativeType] : ["OTHER"],
				},
				comment: comment ?? undefined,
			});

			return successResult([
				"[EVAL MODE] Negative feedback submitted successfully.",
				`Query ID: ${targetQueryId}`,
				`Reason: ${negativeType ?? "OTHER"}`,
				...(comment ? [`Comment: ${comment}`] : []),
			]);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Eval feedback submission failed: ${message}`);
		}
	},
};

/**
 * Eval subcommand group
 */
export const evalSubcommands: SubcommandGroup = {
	name: "eval",
	description:
		"Eval mode commands for RBAC testing and permission validation. Use these endpoints to test AI assistant queries and feedback without affecting production analytics. Useful for administrators validating access controls.",
	descriptionShort: "RBAC testing mode commands",
	descriptionMedium:
		"Query and provide feedback in eval mode for RBAC testing and permission validation.",
	commands: new Map([
		["query", evalQueryCommand],
		["feedback", evalFeedbackCommand],
	]),
	defaultCommand: evalQueryCommand,
};
