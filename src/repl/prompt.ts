/**
 * Prompt building utilities for the REPL.
 * Constructs the prompt string based on current context.
 */

import type { REPLSession } from "./session.js";
import { CLI_NAME } from "../branding/index.js";

/**
 * Build a plain text prompt string.
 * Shows context-aware prompt based on current navigation state:
 * - Root: "xcsh> "
 * - Domain: "xcsh:virtual> "
 * Note: Actions don't create sub-contexts - they execute immediately
 */
export function buildPlainPrompt(session: REPLSession): string {
	const ctx = session.getContextPath();

	if (ctx.isDomain()) {
		return `${CLI_NAME}:${ctx.domain}> `;
	}

	return `${CLI_NAME}> `;
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
	separator: string; // ":"
	wrapper: { open: string; close: string }; // "<" and ">"
}

export function getPromptParts(session: REPLSession): PromptParts {
	const ctx = session.getContextPath();
	const profile = session.getActiveProfileName();

	return {
		profile: profile,
		prefix: CLI_NAME,
		domain: ctx.domain || null,
		separator: ctx.isDomain() ? ":" : "",
		wrapper: { open: "", close: ">" },
	};
}
