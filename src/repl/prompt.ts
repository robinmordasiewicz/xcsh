/**
 * Prompt building utilities for the REPL.
 * Constructs the prompt string based on current context.
 */

import type { REPLSession } from "./session.js";

/**
 * Build a plain text prompt string.
 * Format: <profile@xc.domain.action> or <xc.domain.action>
 */
export function buildPlainPrompt(session: REPLSession): string {
	const parts: string[] = [];

	// Add profile prefix if active
	const profileName = session.getActiveProfileName();
	if (profileName) {
		parts.push(`${profileName}@xc`);
	} else {
		parts.push("xc");
	}

	const ctx = session.getContextPath();
	if (ctx.domain !== "") {
		parts.push(ctx.domain);
		if (ctx.action !== "") {
			parts.push(ctx.action);
		}
	}

	return `<${parts.join(".")}> `;
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

export function getPromptParts(session: REPLSession): PromptParts {
	const ctx = session.getContextPath();

	return {
		profile: session.getActiveProfileName(),
		prefix: "xc",
		domain: ctx.domain || null,
		action: ctx.action || null,
		separator: ".",
		wrapper: { open: "<", close: ">" },
	};
}
