/**
 * GenAI Query Command
 *
 * Send a single query to the AI assistant and display the response
 */

import type { CommandDefinition, DomainCommandResult } from "../registry.js";
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
import type { GenAIChatSession } from "./types.js";

/**
 * Session state for tracking last query (for feedback)
 */
let lastQueryState: GenAIChatSession = {
	namespace: "",
	lastQueryId: null,
	lastQuery: null,
	followUpQueries: [],
	isActive: false,
};

/**
 * Get the last query state (used by feedback command)
 */
export function getLastQueryState(): GenAIChatSession {
	return lastQueryState;
}

/**
 * Update the last query state
 */
export function updateLastQueryState(partial: Partial<GenAIChatSession>): void {
	lastQueryState = { ...lastQueryState, ...partial };
}

/**
 * Clear the last query state
 */
export function clearLastQueryState(): void {
	lastQueryState = {
		namespace: "",
		lastQueryId: null,
		lastQuery: null,
		followUpQueries: [],
		isActive: false,
	};
}

/**
 * Parse output format, spec flag, and namespace from args
 */
function parseQueryArgs(
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
				i++; // Skip the value
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
 * Query command - Send a single query to the AI assistant
 */
export const queryCommand: CommandDefinition = {
	name: "query",
	description:
		"Send a natural language query to the F5 Distributed Cloud AI assistant. Ask about load balancers, WAF configurations, site status, security events, or any platform topic. Returns AI-generated responses with optional follow-up suggestions. Use --namespace to specify the context for namespace-scoped resources.",
	descriptionShort: "Query the AI assistant",
	descriptionMedium:
		"Send natural language queries to the AI assistant for help with F5 XC platform operations, configurations, and troubleshooting.",
	usage: "<question> [--namespace <ns>]",
	aliases: ["ask", "q"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { format, noColor, spec, namespace, question } = parseQueryArgs(
			args,
			session,
		);

		// Handle --spec flag
		if (spec) {
			const cmdSpec = getCommandSpec("generative_ai query");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		// Validate we have a question
		if (!question.trim()) {
			return errorResult(
				"Please provide a question. Usage: ai query <question>",
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
			const response = await client.query(namespace, question);

			// Store query state for feedback
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
			const lines = renderResponse(response);
			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Query failed: ${message}`);
		}
	},
};
