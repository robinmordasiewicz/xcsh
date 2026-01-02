/**
 * GenAI Chat Command
 *
 * Interactive multi-turn conversation with the AI assistant
 */

import * as readline from "node:readline";
import type { CommandDefinition, DomainCommandResult } from "../registry.js";
import { successResult, errorResult } from "../registry.js";
import type { REPLSession } from "../../repl/session.js";
import { getCommandSpec, formatSpec } from "../../output/index.js";
import { parseDomainOutputFlags } from "../../output/domain-formatter.js";
import { getGenAIClient } from "./client.js";
import { renderResponse } from "./response-renderer.js";
import {
	updateLastQueryState,
	clearLastQueryState,
	getLastQueryState,
} from "./query.js";
import { FEEDBACK_TYPE_MAP } from "./types.js";

/**
 * Parse chat args for namespace and spec flag
 */
function parseChatArgs(
	args: string[],
	session: REPLSession,
): {
	spec: boolean;
	namespace: string;
} {
	const { remainingArgs } = parseDomainOutputFlags(
		args,
		session.getOutputFormat(),
	);

	let spec = false;
	let namespace = session.getNamespace();

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
		}
		i++;
	}

	return { spec, namespace };
}

/**
 * Display chat help
 */
function showChatHelp(): string[] {
	return [
		"",
		"=== AI Chat Commands ===",
		"",
		"  /exit, /quit, /q    - Exit chat mode",
		"  /help, /h           - Show this help",
		"  /clear, /c          - Clear conversation context",
		"  /feedback <type>    - Submit feedback for last response",
		"                        Types: positive, negative",
		"  1, 2, 3...          - Select a follow-up question by number",
		"",
		"Just type your question to query the AI assistant.",
		"",
	];
}

/**
 * Handle feedback submission within chat
 */
async function handleFeedback(
	input: string,
	session: REPLSession,
): Promise<string[]> {
	const state = getLastQueryState();
	if (!state.lastQueryId || !state.lastQuery) {
		return ["No previous query to provide feedback for."];
	}

	const parts = input.split(/\s+/);
	const feedbackType = parts[1]?.toLowerCase();

	if (!feedbackType) {
		return [
			"Usage: /feedback <positive|negative>",
			"  Optional: /feedback negative <type> [comment]",
			"  Types: other, inaccurate, irrelevant, poor_format, slow",
		];
	}

	const apiClient = session.getAPIClient();
	if (!apiClient) {
		return ["Not connected to API."];
	}

	try {
		const client = getGenAIClient(apiClient);

		if (feedbackType === "positive" || feedbackType === "+") {
			await client.feedback({
				query: state.lastQuery,
				query_id: state.lastQueryId,
				namespace: state.namespace,
				positive_feedback: {},
			});
			return ["Positive feedback submitted. Thank you!"];
		}

		if (feedbackType === "negative" || feedbackType === "-") {
			const negType = parts[2]?.toLowerCase();
			const mappedType = negType ? FEEDBACK_TYPE_MAP[negType] : undefined;
			const comment = parts.slice(3).join(" ") || undefined;

			await client.feedback({
				query: state.lastQuery,
				query_id: state.lastQueryId,
				namespace: state.namespace,
				negative_feedback: {
					remarks: mappedType ? [mappedType] : ["OTHER"],
				},
				comment,
			});
			return [
				"Negative feedback submitted. Thank you for helping improve the AI.",
			];
		}

		return [
			`Unknown feedback type: ${feedbackType}`,
			"Use 'positive' or 'negative'.",
		];
	} catch (error) {
		const message = error instanceof Error ? error.message : String(error);
		return [`Feedback failed: ${message}`];
	}
}

/**
 * Run interactive chat loop
 */
async function runChatLoop(
	session: REPLSession,
	namespace: string,
): Promise<string[]> {
	const apiClient = session.getAPIClient();
	if (!apiClient) {
		return ["Not connected to API. Please configure connection first."];
	}

	if (!session.isTokenValidated()) {
		return ["Not authenticated. Please check your API token."];
	}

	const client = getGenAIClient(apiClient);
	const output: string[] = [];

	output.push("");
	output.push("=== F5 XC AI Assistant Chat ===");
	output.push(`Namespace: ${namespace}`);
	output.push("Type /help for commands, /exit to quit.");
	output.push("");

	// Create readline interface
	const rl = readline.createInterface({
		input: process.stdin,
		output: process.stdout,
		terminal: process.stdin.isTTY ?? false,
	});

	// Handle Ctrl+C
	let interrupted = false;
	rl.on("SIGINT", () => {
		interrupted = true;
		console.log("\n(Use /exit to leave chat mode)");
		rl.prompt();
	});

	// Promise-based question helper
	const askQuestion = (prompt: string): Promise<string> => {
		return new Promise((resolve) => {
			rl.question(prompt, (answer) => {
				resolve(answer);
			});
		});
	};

	// Print initial output
	for (const line of output) {
		console.log(line);
	}

	// Main chat loop
	let running = true;
	while (running && !interrupted) {
		const input = await askQuestion("ai> ");

		if (interrupted) {
			break;
		}

		const trimmed = input.trim();

		// Empty input
		if (!trimmed) {
			continue;
		}

		// Exit commands
		if (trimmed === "/exit" || trimmed === "/quit" || trimmed === "/q") {
			console.log("Exiting chat mode.");
			running = false;
			break;
		}

		// Help command
		if (trimmed === "/help" || trimmed === "/h") {
			for (const line of showChatHelp()) {
				console.log(line);
			}
			continue;
		}

		// Clear command
		if (trimmed === "/clear" || trimmed === "/c") {
			clearLastQueryState();
			console.log("Conversation context cleared.");
			continue;
		}

		// Feedback command
		if (trimmed.startsWith("/feedback")) {
			const feedbackLines = await handleFeedback(trimmed, session);
			for (const line of feedbackLines) {
				console.log(line);
			}
			continue;
		}

		// Follow-up selection by number
		if (/^\d+$/.test(trimmed)) {
			const num = parseInt(trimmed, 10);
			const state = getLastQueryState();
			if (
				state.followUpQueries.length > 0 &&
				num >= 1 &&
				num <= state.followUpQueries.length
			) {
				const followUp = state.followUpQueries[num - 1];
				if (followUp) {
					console.log(`\nFollowing up: ${followUp}\n`);
					try {
						const response = await client.query(
							namespace,
							followUp,
						);

						updateLastQueryState({
							namespace,
							lastQueryId: response.query_id,
							lastQuery: followUp,
							followUpQueries: response.follow_up_queries ?? [],
						});

						const lines = renderResponse(response);
						for (const line of lines) {
							console.log(line);
						}
						console.log("");
					} catch (error) {
						const message =
							error instanceof Error
								? error.message
								: String(error);
						console.log(`Query failed: ${message}`);
					}
					continue;
				}
			}
			console.log(
				`Invalid selection. Choose 1-${state.followUpQueries.length} from suggested follow-ups.`,
			);
			continue;
		}

		// Unknown command
		if (trimmed.startsWith("/")) {
			console.log(
				`Unknown command: ${trimmed}. Type /help for commands.`,
			);
			continue;
		}

		// Regular query
		try {
			const response = await client.query(namespace, trimmed);

			updateLastQueryState({
				namespace,
				lastQueryId: response.query_id,
				lastQuery: trimmed,
				followUpQueries: response.follow_up_queries ?? [],
			});

			console.log("");
			const lines = renderResponse(response);
			for (const line of lines) {
				console.log(line);
			}
			console.log("");
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			console.log(`Query failed: ${message}`);
		}
	}

	rl.close();

	return ["Chat session ended."];
}

/**
 * Chat command - Interactive conversation with AI assistant
 */
export const chatCommand: CommandDefinition = {
	name: "chat",
	description:
		"Start an interactive conversation with the F5 XC AI assistant. Enter a multi-turn dialog where you can ask questions, receive responses with follow-up suggestions, and navigate through topics naturally. Use numbered responses to quickly select suggested follow-up questions. Supports in-chat feedback submission. Type /exit to return to the main CLI.",
	descriptionShort: "Interactive AI chat mode",
	descriptionMedium:
		"Start an interactive multi-turn conversation with the AI assistant. Supports follow-up suggestions and in-chat commands.",
	usage: "[--namespace <ns>]",
	aliases: ["interactive", "i"],

	async execute(args, session): Promise<DomainCommandResult> {
		const { spec, namespace } = parseChatArgs(args, session);

		// Handle --spec flag
		if (spec) {
			const cmdSpec = getCommandSpec("generative_ai chat");
			if (cmdSpec) {
				return successResult([formatSpec(cmdSpec)]);
			}
		}

		// Check if running in a TTY
		if (!process.stdin.isTTY) {
			return errorResult(
				"Chat mode requires an interactive terminal. Use 'ai query' for non-interactive queries.",
			);
		}

		try {
			const result = await runChatLoop(session, namespace);
			return successResult(result);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Chat session failed: ${message}`);
		}
	},
};
