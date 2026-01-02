/**
 * GenAI Feedback Command
 *
 * Submit feedback for AI assistant queries
 */

import type { CommandDefinition, DomainCommandResult } from "../registry.js";
import { successResult, errorResult } from "../registry.js";
import type { REPLSession } from "../../repl/session.js";
import { getCommandSpec, formatSpec } from "../../output/index.js";
import { parseDomainOutputFlags } from "../../output/domain-formatter.js";
import { getGenAIClient } from "./client.js";
import { getLastQueryState } from "./query.js";
import { FEEDBACK_TYPE_MAP, getValidFeedbackTypes } from "./types.js";
import type { NegativeFeedbackType } from "./types.js";

/**
 * Parse feedback command args
 */
function parseFeedbackArgs(
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
					// Default to OTHER if type not recognized
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
 * Feedback command - Submit feedback for AI queries
 */
export const feedbackCommand: CommandDefinition = {
	name: "feedback",
	description:
		"Submit feedback on AI assistant responses to help improve future answers. Provide positive feedback when responses are helpful, or negative feedback with a reason type when improvements are needed. Feedback is associated with the most recent query unless a specific query ID is provided. Negative feedback types include: inaccurate (wrong information), irrelevant (off-topic), poor_format (hard to read), slow (response time), or other.",
	descriptionShort: "Submit feedback for AI responses",
	descriptionMedium:
		"Provide positive or negative feedback for AI assistant responses. Use --negative with a type: inaccurate, irrelevant, poor_format, slow, or other.",
	usage: "--positive | --negative <type> [--comment <text>] [--query-id <id>]",
	aliases: ["fb", "rate"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { spec, namespace, positive, negativeType, comment, queryId } =
			parseFeedbackArgs(args, session);

		// Handle --spec flag
		if (spec) {
			const cmdSpec = getCommandSpec("generative_ai feedback");
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
				"No query to provide feedback for. Make a query first, or use --query-id to specify one.",
			);
		}

		if (!targetQuery && !queryId) {
			return errorResult(
				"Could not find the original query. Use --query-id to specify the query ID.",
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
				await client.feedback({
					query: targetQuery ?? "",
					query_id: targetQueryId,
					namespace: state.namespace || namespace,
					positive_feedback: {},
					comment: comment ?? undefined,
				});
				return successResult([
					"Positive feedback submitted successfully.",
					`Query ID: ${targetQueryId}`,
				]);
			}

			// Negative feedback
			await client.feedback({
				query: targetQuery ?? "",
				query_id: targetQueryId,
				namespace: state.namespace || namespace,
				negative_feedback: {
					remarks: negativeType ? [negativeType] : ["OTHER"],
				},
				comment: comment ?? undefined,
			});

			return successResult([
				"Negative feedback submitted successfully.",
				`Query ID: ${targetQueryId}`,
				`Reason: ${negativeType ?? "OTHER"}`,
				...(comment ? [`Comment: ${comment}`] : []),
			]);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Feedback submission failed: ${message}`);
		}
	},
};
