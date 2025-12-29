/**
 * Display Command - Render the xcsh banner with configurable logo modes
 *
 * Supports logo display modes:
 * - image: Image if terminal supports it, otherwise ASCII (default)
 * - ascii: ASCII art only
 * - none: No logo
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult, rawStdoutResult } from "../../registry.js";
import {
	CLI_FULL_NAME,
	CLI_VERSION,
	F5_LOGO,
	ENV_PREFIX,
	colorRed,
	colorBoldWhite,
	colors,
	renderLogo,
	resolveLogoMode,
	getLogoModeFromEnv,
	detectTerminalCapabilities,
	getTerminalImageSequence,
} from "../../../branding/index.js";
import {
	F5_LOGO_PNG_BASE64,
	F5_LOGO_DISPLAY_WIDTH,
	F5_LOGO_DISPLAY_HEIGHT,
} from "../../../branding/logo-image.js";
import {
	loadSettingsSync,
	isValidLogoMode,
	LOGO_MODES,
	type LogoDisplayMode,
} from "../../../config/index.js";

/**
 * Colorize a logo line for raw ANSI output.
 * Applies red to circle chars and white to F5 block chars.
 */
function colorizeLogoLine(line: string): string {
	let result = "";
	let currentColor: "red" | "white" | "none" = "none";

	for (const char of line) {
		let newColor: "red" | "white" | "none";

		switch (char) {
			case "\u2593": // Dark shade - red (rendered as solid block)
			case "\u2592": // Medium shade - red
			case "(":
			case ")":
			case "|":
			case "_":
				newColor = "red";
				break;
			case "\u2588": // Full block - white F5 text
				newColor = "white";
				break;
			default:
				newColor = "none";
		}

		if (newColor !== currentColor) {
			// Close previous color if active
			if (currentColor !== "none") {
				result += colors.reset;
			}
			// Apply new color
			if (newColor === "red") {
				result += colors.red;
			} else if (newColor === "white") {
				result += colors.boldWhite;
			}
			currentColor = newColor;
		}

		// Render dark shade as solid block for consistency
		result += char === "\u2593" ? "\u2588" : char;
	}

	// Close any active color
	if (currentColor !== "none") {
		result += colors.reset;
	}

	return result;
}

/**
 * Print banner with inline image logo inside the box.
 * Image on left, help text on right - uses cursor positioning.
 */
function printImageBanner(
	imageSeq: string,
	imageHeight: number,
	imageWidth: number,
): void {
	const BOX = {
		topLeft: "\u256D",
		topRight: "\u256E",
		bottomLeft: "\u2570",
		bottomRight: "\u256F",
		horizontal: "\u2500",
		vertical: "\u2502",
	};

	const TOTAL_WIDTH = 80;
	const INNER_WIDTH = TOTAL_WIDTH - 2;

	// Image column: padding (1) + image + padding (1) - borders are separate
	const IMAGE_COL_WIDTH = 1 + imageWidth + 1;
	// Text column: remaining space
	const TEXT_COL_WIDTH = INNER_WIDTH - IMAGE_COL_WIDTH;

	// Help text - will be vertically centered
	const HELP_LINES = [
		"Type 'help' for commands",
		"Run 'namespace <ns>' to set",
		"Press Ctrl+C twice to exit",
	];

	// Calculate vertical centering for help text
	const helpStartRow = Math.floor((imageHeight - HELP_LINES.length) / 2);

	// Build title line
	const title = ` ${CLI_FULL_NAME} v${CLI_VERSION} `;
	const leftDashes = 3;
	const rightDashes = TOTAL_WIDTH - 1 - leftDashes - title.length - 1;

	// Blank line after command for visual separation
	process.stdout.write("\n");

	// STEP 1: Draw the complete frame first (all borders + text content)
	// This ensures clean lines before we overlay the image

	// Top border
	process.stdout.write(
		colorRed(BOX.topLeft + BOX.horizontal.repeat(leftDashes)) +
			colorBoldWhite(title) +
			colorRed(
				BOX.horizontal.repeat(Math.max(0, rightDashes)) + BOX.topRight,
			) +
			"\n",
	);

	// One row of padding above image (matches ASCII banner)
	process.stdout.write(
		colorRed(BOX.vertical) +
			" ".repeat(INNER_WIDTH) +
			colorRed(BOX.vertical) +
			"\n",
	);

	// Content rows: left border + spaces for image area + text + right border
	for (let row = 0; row < imageHeight; row++) {
		const helpIndex = row - helpStartRow;
		const helpText =
			helpIndex >= 0 && helpIndex < HELP_LINES.length
				? (HELP_LINES[helpIndex] ?? "")
				: "";

		// Build the full row: │ + image_space + text + │
		const imageSpace = " ".repeat(IMAGE_COL_WIDTH);
		const paddedHelp = helpText.padEnd(TEXT_COL_WIDTH);

		process.stdout.write(
			colorRed(BOX.vertical) +
				imageSpace +
				colorBoldWhite(paddedHelp) +
				colorRed(BOX.vertical) +
				"\n",
		);
	}

	// Bottom border (no extra padding below - matches ASCII banner)
	process.stdout.write(
		colorRed(
			BOX.bottomLeft +
				BOX.horizontal.repeat(INNER_WIDTH) +
				BOX.bottomRight,
		) + "\n",
	);

	// STEP 2: Now overlay the image on top of the frame
	// Move cursor back up to the first content row (after padding), then to image position
	// Total rows to go up: image rows + bottom border (1)
	const rowsUp = imageHeight + 1;
	process.stdout.write(`\x1b[${rowsUp}A`); // Move up to first image row (after padding)
	process.stdout.write(`\x1b[2C`); // Move right: past left border (1) + padding (1)

	// Output the image - it will overwrite the spaces we left
	process.stdout.write(imageSeq);

	// STEP 3: Restore cursor to a known good position
	// After iTerm2 renders the image, cursor position is unpredictable
	// Move down past the banner content, then to column 1
	const rowsDown = imageHeight + 1; // Match what we moved up
	process.stdout.write(`\x1b[${rowsDown}B`); // Move down past banner
	process.stdout.write(`\x1b[1G`); // Move to column 1 (absolute)
	process.stdout.write("\n"); // Extra newline to ensure clean separation
}

/**
 * Print ASCII banner to stdout
 */
function printAsciiBanner(): void {
	const BOX = {
		topLeft: "\u256D",
		topRight: "\u256E",
		bottomLeft: "\u2570",
		bottomRight: "\u256F",
		horizontal: "\u2500",
		vertical: "\u2502",
	};

	// Calculate logo dimensions
	const logoLines = F5_LOGO.split("\n");
	const logoWidth = Math.max(...logoLines.map((l) => [...l].length));
	const TOTAL_WIDTH = Math.max(80, logoWidth + 4);
	const INNER_WIDTH = TOTAL_WIDTH - 2;

	// Help text positioned on specific logo rows (vertically centered)
	const HELP_LINES = [
		"Type 'help' for commands",
		"Run 'namespace <ns>' to set",
		"Press Ctrl+C twice to exit",
	];
	const HELP_START_ROW = 8;

	// Build title line
	const title = ` ${CLI_FULL_NAME} v${CLI_VERSION} `;
	const leftDashes = 3;
	const rightDashes = TOTAL_WIDTH - 1 - leftDashes - title.length - 1;

	const output: string[] = [];

	// Top border with title
	output.push(
		colorRed(BOX.topLeft + BOX.horizontal.repeat(leftDashes)) +
			colorBoldWhite(title) +
			colorRed(
				BOX.horizontal.repeat(Math.max(0, rightDashes)) + BOX.topRight,
			),
	);

	// Logo lines with help text overlay
	for (let i = 0; i < logoLines.length; i++) {
		const logoLine = logoLines[i] ?? "";
		const helpIndex = i - HELP_START_ROW;
		const helpText: string =
			helpIndex >= 0 && helpIndex < HELP_LINES.length
				? (HELP_LINES[helpIndex] ?? "")
				: "";

		// Pad logo to consistent width, then add help text
		const paddedLogo = logoLine.padEnd(logoWidth);
		const coloredLogo = colorizeLogoLine(paddedLogo);

		// Calculate remaining space for help column
		const helpColumnWidth = INNER_WIDTH - logoWidth - 1;
		const paddedHelp = helpText.padEnd(helpColumnWidth);

		output.push(
			colorRed(BOX.vertical) +
				coloredLogo +
				" " +
				colorBoldWhite(paddedHelp) +
				colorRed(BOX.vertical),
		);
	}

	// Bottom border
	output.push(
		colorRed(
			BOX.bottomLeft +
				BOX.horizontal.repeat(INNER_WIDTH) +
				BOX.bottomRight,
		),
	);

	// Empty line for spacing
	output.push("");

	// Blank line before banner for visual separation
	process.stdout.write("\n");
	process.stdout.write(output.join("\n") + "\n");
}

/**
 * Get banner as array of lines (for REPL mode)
 * Returns banner lines without writing to stdout - suitable for Ink's output system
 * Supports image mode (inline image without cursor positioning) and ASCII mode
 */
function getBannerLines(
	logoMode?: LogoDisplayMode,
	useImage?: boolean,
): string[] {
	// Handle "none" mode
	if (logoMode === "none") {
		return [];
	}

	const BOX = {
		topLeft: "\u256D",
		topRight: "\u256E",
		bottomLeft: "\u2570",
		bottomRight: "\u256F",
		horizontal: "\u2500",
		vertical: "\u2502",
	};

	// Calculate logo dimensions
	const logoLines = F5_LOGO.split("\n");
	const logoWidth = Math.max(...logoLines.map((l) => [...l].length));
	const TOTAL_WIDTH = Math.max(80, logoWidth + 4);
	const INNER_WIDTH = TOTAL_WIDTH - 2;

	// Help text
	const HELP_LINES = [
		"Type 'help' for commands",
		"Run 'namespace <ns>' to set",
		"Press Ctrl+C twice to exit",
	];

	// Build title line
	const title = ` ${CLI_FULL_NAME} v${CLI_VERSION} `;
	const leftDashes = 3;
	const rightDashes = TOTAL_WIDTH - 1 - leftDashes - title.length - 1;

	const output: string[] = [];

	// IMAGE MODE: Output inline image inside the banner box (matching startup layout)
	if (useImage) {
		const capabilities = detectTerminalCapabilities();
		const imageSeq = getTerminalImageSequence(
			F5_LOGO_PNG_BASE64,
			capabilities,
			{
				width: F5_LOGO_DISPLAY_WIDTH,
				height: F5_LOGO_DISPLAY_HEIGHT,
				preserveAspectRatio: true,
			},
		);

		if (imageSeq) {
			// Build the entire banner as a single string with cursor positioning
			// This matches the startup banner layout exactly
			const imageHeight = F5_LOGO_DISPLAY_HEIGHT;
			const imageWidth = F5_LOGO_DISPLAY_WIDTH;

			// Image column: padding (1) + image + padding (1)
			const IMAGE_COL_WIDTH = 1 + imageWidth + 1;
			// Text column: remaining space
			const TEXT_COL_WIDTH = INNER_WIDTH - IMAGE_COL_WIDTH;

			// Calculate vertical centering for help text
			const helpStartRow = Math.floor(
				(imageHeight - HELP_LINES.length) / 2,
			);

			let bannerStr = "";

			// CRITICAL: Single blank line for visual separation from command prompt
			// DO NOT ADD MORE - causes excessive deadspace above banner
			// DO NOT REMOVE - banner title will touch command line
			bannerStr += "\n";

			// Top border
			bannerStr +=
				colorRed(BOX.topLeft + BOX.horizontal.repeat(leftDashes)) +
				colorBoldWhite(title) +
				colorRed(
					BOX.horizontal.repeat(Math.max(0, rightDashes)) +
						BOX.topRight,
				) +
				"\n";

			// One row of padding above image
			bannerStr +=
				colorRed(BOX.vertical) +
				" ".repeat(INNER_WIDTH) +
				colorRed(BOX.vertical) +
				"\n";

			// Content rows: left border + spaces for image area + text + right border
			for (let row = 0; row < imageHeight; row++) {
				const helpIndex = row - helpStartRow;
				const helpText =
					helpIndex >= 0 && helpIndex < HELP_LINES.length
						? (HELP_LINES[helpIndex] ?? "")
						: "";

				const imageSpace = " ".repeat(IMAGE_COL_WIDTH);
				const paddedHelp = helpText.padEnd(TEXT_COL_WIDTH);

				bannerStr +=
					colorRed(BOX.vertical) +
					imageSpace +
					colorBoldWhite(paddedHelp) +
					colorRed(BOX.vertical) +
					"\n";
			}

			// Bottom border
			bannerStr +=
				colorRed(
					BOX.bottomLeft +
						BOX.horizontal.repeat(INNER_WIDTH) +
						BOX.bottomRight,
				) + "\n";

			// Now overlay the image using cursor positioning
			// Match printImageBanner exactly: move up (imageHeight + 1) rows
			// This skips: bottom border (1) + content rows (imageHeight) to land on first content row
			const rowsUp = imageHeight + 1;
			bannerStr += `\x1b[${rowsUp}A`; // Move up to first content row
			bannerStr += `\x1b[2C`; // Move right: past left border (1) + padding (1)
			bannerStr += imageSeq; // Output the image
			// Restore cursor position after image overlay
			const rowsDown = imageHeight + 1;
			bannerStr += `\x1b[${rowsDown}B`; // Move down past banner
			bannerStr += `\x1b[1G`; // Move to column 1
			// CRITICAL: No trailing newline here - App.tsx useEffect adds exactly 3
			// Adding newlines here causes excessive deadspace below banner

			// Return as single item - Ink will output it all at once
			return [bannerStr];
		}
		// Fall through to ASCII if image sequence failed
	}

	// ASCII MODE: Traditional banner with ASCII logo
	const HELP_START_ROW = 8;

	// Blank line for visual separation
	output.push("");

	// Top border with title
	output.push(
		colorRed(BOX.topLeft + BOX.horizontal.repeat(leftDashes)) +
			colorBoldWhite(title) +
			colorRed(
				BOX.horizontal.repeat(Math.max(0, rightDashes)) + BOX.topRight,
			),
	);

	// Logo lines with help text overlay
	for (let i = 0; i < logoLines.length; i++) {
		const logoLine = logoLines[i] ?? "";
		const helpIndex = i - HELP_START_ROW;
		const helpText: string =
			helpIndex >= 0 && helpIndex < HELP_LINES.length
				? (HELP_LINES[helpIndex] ?? "")
				: "";

		// Pad logo to consistent width, then add help text
		const paddedLogo = logoLine.padEnd(logoWidth);
		const coloredLogo = colorizeLogoLine(paddedLogo);

		// Calculate remaining space for help column
		const helpColumnWidth = INNER_WIDTH - logoWidth - 1;
		const paddedHelp = helpText.padEnd(helpColumnWidth);

		output.push(
			colorRed(BOX.vertical) +
				coloredLogo +
				" " +
				colorBoldWhite(paddedHelp) +
				colorRed(BOX.vertical),
		);
	}

	// Bottom border
	output.push(
		colorRed(
			BOX.bottomLeft +
				BOX.horizontal.repeat(INNER_WIDTH) +
				BOX.bottomRight,
		),
	);

	// Empty line for spacing
	output.push("");

	return output;
}

/** Rendering context - startup uses direct stdout, repl returns lines */
export type RenderContext = "startup" | "repl";

/**
 * Render banner with specified logo mode
 * Exported for use by main entry point during REPL startup
 *
 * @param logoMode - Logo display mode (auto, image, ascii, both, none)
 * @param context - Rendering context:
 *   - "startup": Write directly to stdout (before Ink takes over)
 *   - "repl": Return lines as array (for Ink's output system)
 * @returns Lines array when context is "repl", void when "startup"
 */
export function renderBanner(
	logoMode?: LogoDisplayMode,
	context: RenderContext = "startup",
): string[] | void {
	// Resolve effective logo mode using priority: CLI > env > config > auto
	const settings = loadSettingsSync();
	const envMode = getLogoModeFromEnv(ENV_PREFIX);
	const effectiveMode = resolveLogoMode({
		cliMode: logoMode,
		envMode,
		configMode: settings.logo,
	});

	// Handle "none" mode - skip banner entirely
	if (effectiveMode === "none") {
		return context === "repl" ? [] : undefined;
	}

	// REPL mode: return lines as array for Ink's output system
	// Uses inline image (without cursor positioning) when terminal supports it
	if (context === "repl") {
		const rendered = renderLogo(effectiveMode);
		const useImage = rendered.effectiveMode === "image";
		return getBannerLines(rendered.effectiveMode, useImage);
	}

	// === STARTUP MODE: Write directly to stdout ===

	// Render the logo based on mode and terminal capabilities
	const rendered = renderLogo(effectiveMode);

	// If image mode succeeded, render image inside the banner box
	if (rendered.usedImage && rendered.effectiveMode === "image") {
		// Generate image sequence with proper dimensions for the banner
		const capabilities = detectTerminalCapabilities();
		const imageSeq = getTerminalImageSequence(
			F5_LOGO_PNG_BASE64,
			capabilities,
			{
				width: F5_LOGO_DISPLAY_WIDTH,
				height: F5_LOGO_DISPLAY_HEIGHT,
				preserveAspectRatio: true,
			},
		);
		if (imageSeq) {
			printImageBanner(
				imageSeq,
				F5_LOGO_DISPLAY_HEIGHT,
				F5_LOGO_DISPLAY_WIDTH,
			);
		}
		return;
	}

	// ASCII mode or fallback - render traditional banner
	printAsciiBanner();
}

/**
 * Parse --logo flag from args
 */
function parseLogoArg(args: string[]): LogoDisplayMode | undefined {
	for (let i = 0; i < args.length; i++) {
		const arg = args[i];
		if (arg === "--logo" && args[i + 1]) {
			const mode = args[i + 1];
			if (mode && isValidLogoMode(mode)) {
				return mode as LogoDisplayMode;
			}
		} else if (arg?.startsWith("--logo=")) {
			const mode = arg.slice("--logo=".length);
			if (isValidLogoMode(mode)) {
				return mode as LogoDisplayMode;
			}
		}
	}
	return undefined;
}

/**
 * Banner command - renders the xcsh banner
 */
export const bannerCommand: CommandDefinition = {
	name: "banner",
	description: "Display the xcsh banner with optional logo mode",
	usage: "[--logo <mode>]",

	async execute(args, _session) {
		// Check for help flags
		if (args.includes("--help") || args.includes("-h")) {
			// Generate help dynamically from LOGO_MODES (single source of truth)
			const helpLines = [
				"login banner - Display the xcsh banner",
				"",
				"Usage: login banner [--logo <mode>]",
				"",
				"Options:",
				"  --logo <mode>    Logo display mode:",
				// Generate mode list dynamically from LOGO_MODES
				...LOGO_MODES.map(
					(m) =>
						`                   - ${m.mode.padEnd(6)} ${m.description}`,
				),
				"  --help, -h       Show this help message",
				"",
				"Priority Resolution:",
				"  1. CLI flag --logo (highest)",
				"  2. Environment variable F5XC_LOGO",
				"  3. Config file ~/.xcshconfig → logo: <mode>",
				"  4. Auto-detect based on terminal capabilities (fallback)",
				"",
				"Examples:",
				"  login banner              # Use default/auto mode",
				"  login banner --logo image # Force image mode",
				"  login banner --logo ascii # Force ASCII mode",
			];
			return successResult(helpLines);
		}

		// Parse logo mode from args
		const logoMode = parseLogoArg(args);

		// Validate logo mode if provided
		if (
			args.some((a) => a === "--logo" || a.startsWith("--logo=")) &&
			!logoMode
		) {
			return errorResult("Invalid logo mode. Use: image, ascii, or none");
		}

		try {
			// Resolve effective logo mode using priority: CLI > env > config > auto
			// This follows the SAME logic as startup
			const settings = loadSettingsSync();
			const envMode = getLogoModeFromEnv(ENV_PREFIX);
			const resolvedMode = resolveLogoMode({
				cliMode: logoMode,
				envMode,
				configMode: settings.logo,
			});

			// Handle "none" mode - skip banner entirely
			if (resolvedMode === "none") {
				return successResult([]);
			}

			// Use renderLogo to determine effective mode based on terminal capabilities
			// This follows the SAME logic as startup:
			// - "image" mode with image-capable terminal → effectiveMode = "image"
			// - "image" mode without image support → effectiveMode = "ascii" (auto fallback)
			// - "ascii" mode → effectiveMode = "ascii"
			const rendered = renderLogo(resolvedMode);

			// Image mode: Generate banner content for raw stdout
			// This bypasses Ink's rendering to avoid status bar in scrollback
			if (rendered.effectiveMode === "image") {
				const bannerLines = getBannerLines("image", true);
				if (bannerLines.length > 0 && bannerLines[0]) {
					// getBannerLines returns the image banner as a single string
					// with embedded cursor positioning sequences
					return rawStdoutResult(bannerLines[0]);
				}
				// Fall through to ASCII if image generation failed
			}

			// ASCII mode: Return lines for Ink's output system
			const lines = getBannerLines("ascii", false);
			return successResult(lines);
		} catch (error) {
			return errorResult(
				`Failed to display banner: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},

	async completion(partial: string, args: string[], _session) {
		// If previous arg was --logo, suggest modes (derived from LOGO_MODES)
		const lastArg = args[args.length - 1];
		if (lastArg === "--logo") {
			const modes = LOGO_MODES.map((m) => m.mode);
			return modes.filter((m) => m.startsWith(partial.toLowerCase()));
		}

		// Suggest --logo flag if not already present
		if (!args.includes("--logo") && "--logo".startsWith(partial)) {
			return ["--logo"];
		}

		return [];
	},
};
