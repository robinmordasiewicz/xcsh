/**
 * Prompt building utilities for the REPL.
 * Constructs the prompt string based on current context.
 */

import type { REPLSession } from "./session.js";
import {
	detectTerminalCapabilities,
	generateITerm2ImageSequence,
	generateKittyImageSequence,
} from "../branding/terminal.js";
import { F5_LOGO_PNG_BASE64 } from "../branding/logo-image.js";

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

/**
 * Build a prompt with F5 logo for terminals that support inline images.
 * Falls back to plain "> " prompt for non-image terminals.
 */
export function buildPromptWithLogo(_session: REPLSession): string {
	const capabilities = detectTerminalCapabilities();

	if (!capabilities.supportsInlineImages) {
		return "> "; // Plain fallback for non-image terminals
	}

	// Generate inline image (1 cell height)
	let imageSequence: string;

	if (capabilities.isKitty) {
		// Kitty uses different protocol - fall back to plain for now
		// Kitty sizing with c/r parameters is complex and may not work consistently
		imageSequence = generateKittyImageSequence(F5_LOGO_PNG_BASE64);
	} else {
		// iTerm2, WezTerm, Mintty - use height=1 for single character cell height
		imageSequence = generateITerm2ImageSequence(F5_LOGO_PNG_BASE64, {
			height: 1, // 1 character cell height - terminal scales image automatically
			preserveAspectRatio: true,
			inline: true,
		});
	}

	return `${imageSequence} `; // Logo + space for input
}
