/**
 * Prompt building utilities for the REPL.
 * Constructs the prompt string based on current context.
 */

import type { REPLSession } from "./session.js";

/**
 * Build a plain text prompt string.
 * Simple ">" prompt.
 */
export function buildPlainPrompt(_session: REPLSession): string {
	return "> ";
}

/**
 * Build a prompt with color formatting.
 * Uses F5 brand colors for visual enhancement.
 */
export function buildColoredPrompt(session: REPLSession): string {
	// For now, return plain prompt - we'll add colors when we integrate with Ink
	return buildPlainPrompt(session);
}

/**
 * Build prompt parts for use in Ink components.
 * Returns structured data for flexible rendering.
 */
export interface PromptParts {
	profile: string | null; // Active profile name
	prefix: string; // "xc"
	domain: string | null; // e.g., "http_loadbalancer"
	action: string | null; // e.g., "list"
	separator: string; // "."
	wrapper: { open: string; close: string }; // "<" and ">"
}

export function getPromptParts(_session: REPLSession): PromptParts {
	return {
		profile: null,
		prefix: "",
		domain: null,
		action: null,
		separator: "",
		wrapper: { open: "", close: ">" },
	};
}
