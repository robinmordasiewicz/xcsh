/**
 * Generative AI Domain - AI Assistant for F5 Distributed Cloud
 *
 * Query the AI assistant for help with platform operations, configurations,
 * and troubleshooting. Supports single queries, interactive chat, and feedback.
 */

import type { DomainDefinition } from "../registry.js";
import { queryCommand } from "./query.js";
import { chatCommand } from "./chat.js";
import { feedbackCommand } from "./feedback.js";
import { evalSubcommands } from "./eval.js";

/**
 * Generative AI domain definition
 */
export const generativeAiDomain: DomainDefinition = {
	name: "generative_ai",
	description:
		"Interact with the F5 Distributed Cloud AI assistant for natural language queries about platform operations. Ask questions about load balancers, WAF configurations, site status, security events, or any platform topic. Supports single queries with follow-up suggestions, interactive multi-turn chat sessions, and feedback submission to improve AI responses.",
	descriptionShort: "AI assistant queries and feedback",
	descriptionMedium:
		"Query the AI assistant for help with F5 XC platform operations, configurations, security analysis, and troubleshooting.",
	defaultCommand: queryCommand,
	commands: new Map([
		["query", queryCommand],
		["chat", chatCommand],
		["feedback", feedbackCommand],
	]),
	subcommands: new Map([["eval", evalSubcommands]]),
};

/**
 * Domain aliases
 */
export const generativeAiAliases = ["ai", "genai", "assistant"];

// Re-export types and utilities for external use
export type {
	GenAIQueryRequest,
	GenAIQueryResponse,
	GenAIFeedbackRequest,
	GenAIChatSession,
	NegativeFeedbackType,
} from "./types.js";

export { getResponseType, FEEDBACK_TYPE_MAP } from "./types.js";
export { GenAIClient, getGenAIClient } from "./client.js";
export { renderResponse } from "./response-renderer.js";
