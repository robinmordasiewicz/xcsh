/**
 * Logo rendering utility with iTerm2 inline image support.
 * Handles mode resolution and generates appropriate output.
 */

import type { LogoDisplayMode } from "../config/settings.js";
import {
	detectTerminalCapabilities,
	getTerminalImageSequence,
	type TerminalCapabilities,
} from "./terminal.js";
import { F5_LOGO } from "./index.js";
import {
	F5_LOGO_PNG_BASE64,
	F5_LOGO_DISPLAY_WIDTH,
	F5_LOGO_DISPLAY_HEIGHT,
} from "./logo-image.js";

/**
 * Options for resolving logo display mode.
 */
export interface LogoRenderOptions {
	/** CLI flag value (highest priority) */
	cliMode?: LogoDisplayMode | undefined;
	/** Environment variable value */
	envMode?: LogoDisplayMode | undefined;
	/** Config file value */
	configMode?: LogoDisplayMode | undefined;
}

/**
 * Result of logo rendering.
 */
export interface RenderedLogo {
	/** The logo content to display (escape sequences and/or ASCII) */
	content: string;
	/** Number of lines the logo occupies */
	lineCount: number;
	/** Whether inline image was used */
	usedImage: boolean;
	/** Terminal capabilities detected */
	capabilities: TerminalCapabilities;
	/** Effective mode that was used */
	effectiveMode: "image" | "ascii" | "none";
}

/**
 * Resolve the effective logo mode using priority:
 * CLI flag > Environment variable > Config file > image (default)
 *
 * @param options - Logo mode sources
 * @returns Effective logo mode to use
 */
export function resolveLogoMode(options: LogoRenderOptions): LogoDisplayMode {
	// CLI flag has highest priority
	if (options.cliMode) {
		return options.cliMode;
	}

	// Environment variable is next
	if (options.envMode) {
		return options.envMode;
	}

	// Config file is last explicit option
	if (options.configMode) {
		return options.configMode;
	}

	// Default: image mode (tries image first, falls back to ASCII)
	return "image";
}

/**
 * Check if image logo data is available.
 */
export function hasImageData(): boolean {
	return F5_LOGO_PNG_BASE64.length > 0;
}

/**
 * Render the F5 logo based on mode and terminal capabilities.
 *
 * @param mode - Display mode to use
 * @returns Rendered logo with content and metadata
 */
export function renderLogo(mode: LogoDisplayMode): RenderedLogo {
	const capabilities = detectTerminalCapabilities();
	const logoLines = F5_LOGO.split("\n");
	const asciiLineCount = logoLines.length;

	// Determine effective mode based on requested mode and capabilities
	let effectiveMode: "image" | "ascii" | "none";

	if (mode === "none") {
		effectiveMode = "none";
	} else if (mode === "image") {
		// Image mode: use image if supported and available, otherwise fall back to ASCII
		if (capabilities.supportsInlineImages && hasImageData()) {
			effectiveMode = "image";
		} else {
			effectiveMode = "ascii";
		}
	} else {
		// ASCII mode
		effectiveMode = "ascii";
	}

	// Generate output based on effective mode
	switch (effectiveMode) {
		case "image": {
			const imageSeq = getTerminalImageSequence(
				F5_LOGO_PNG_BASE64,
				capabilities,
				{
					width: F5_LOGO_DISPLAY_WIDTH,
					height: "auto",
					preserveAspectRatio: true,
				},
			);

			return {
				content: imageSeq ?? "",
				lineCount: F5_LOGO_DISPLAY_HEIGHT,
				usedImage: true,
				capabilities,
				effectiveMode,
			};
		}

		case "none":
			return {
				content: "",
				lineCount: 0,
				usedImage: false,
				capabilities,
				effectiveMode,
			};

		case "ascii":
		default:
			return {
				content: F5_LOGO,
				lineCount: asciiLineCount,
				usedImage: false,
				capabilities,
				effectiveMode,
			};
	}
}

/**
 * Get logo mode from environment variable.
 *
 * @param envPrefix - Environment variable prefix (e.g., "F5XC")
 * @returns Logo mode or undefined if not set or invalid
 */
export function getLogoModeFromEnv(
	envPrefix: string,
): LogoDisplayMode | undefined {
	const value = process.env[`${envPrefix}_LOGO`];
	if (!value) return undefined;

	const normalized = value.toLowerCase().trim();
	const validModes: LogoDisplayMode[] = ["image", "ascii", "none"];

	if (validModes.includes(normalized as LogoDisplayMode)) {
		return normalized as LogoDisplayMode;
	}

	return undefined;
}
