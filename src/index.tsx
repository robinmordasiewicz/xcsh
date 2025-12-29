#!/usr/bin/env node
/**
 * F5 Distributed Cloud Shell - Interactive CLI for F5 XC
 * Main entry point for the xcsh CLI application.
 *
 * Usage:
 *   xcsh                    # Enter interactive REPL mode
 *   xcsh <command> [args]   # Execute command non-interactively
 *   xcsh --help             # Show help
 *   xcsh --version          # Show version
 */

import { render } from "ink";
import { Command } from "commander";
import { App } from "./repl/index.js";
import {
	CLI_NAME,
	CLI_VERSION,
	CLI_FULL_NAME,
	F5_LOGO,
	colorRed,
	colorBoldWhite,
	colors,
} from "./branding/index.js";
import { executeCommand } from "./repl/executor.js";
import { REPLSession } from "./repl/session.js";
import { formatRootHelp } from "./repl/help.js";

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
 * Print banner to scrollback BEFORE Ink takes over.
 * This ensures banner remains in scrollback buffer and isn't
 * cleared when Ink restructures the component tree.
 */
function printBannerToScrollback(): void {
	const BOX = {
		topLeft: "\u256D",
		topRight: "\u256E",
		bottomLeft: "\u2570",
		bottomRight: "\u256F",
		horizontal: "\u2500",
		vertical: "\u2502",
		leftT: "\u251C",
		rightT: "\u2524",
	};

	// Calculate logo dimensions
	const logoLines = F5_LOGO.split("\n");
	const logoWidth = Math.max(...logoLines.map((l) => [...l].length));
	const TOTAL_WIDTH = Math.max(80, logoWidth + 4); // At least 80 or logo width + padding
	const INNER_WIDTH = TOTAL_WIDTH - 2;

	// Help text positioned on specific logo rows (vertically centered)
	const HELP_LINES = [
		"Type 'help' for commands",
		"Run 'namespace <ns>' to set",
		"Press Ctrl+C twice to exit",
	];
	const HELP_START_ROW = 8; // Start at row 8 of the logo

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

	// Empty line for spacing before prompt
	output.push("");

	process.stdout.write(output.join("\n") + "\n");
}

const program = new Command();

// Custom help: override default help to show comprehensive root help
program.configureHelp({
	formatHelp: () => formatRootHelp().join("\n"),
});

program
	.name(CLI_NAME)
	.description("F5 Distributed Cloud Shell - Interactive CLI for F5 XC")
	.version(CLI_VERSION, "-v, --version", "Show version number")
	.option("-i, --interactive", "Force interactive mode")
	.option("--no-color", "Disable color output")
	.option("-h, --help", "Show help") // Manual help option to prevent auto-exit
	.argument("[command...]", "Command to execute non-interactively")
	.allowUnknownOption(true) // Pass through unknown options to commands
	.helpOption(false) // Disable auto-help so domain --help works
	.action(
		async (
			commandArgs: string[],
			options: { interactive?: boolean; help?: boolean },
		) => {
			// Handle root-level help (xcsh --help or xcsh -h with no domain)
			if (options.help && commandArgs.length === 0) {
				formatRootHelp().forEach((line) => console.log(line));
				process.exit(0);
			}

			// If --help with a domain, re-inject --help into args for domain help
			// Commander consumes --help as an option, so we add it back
			if (options.help && commandArgs.length > 0) {
				commandArgs.push("--help");
			}

			// If no command args or --interactive flag, enter REPL mode
			if (commandArgs.length === 0 || options.interactive) {
				// Check if stdin is a TTY (interactive terminal)
				if (!process.stdin.isTTY && !options.interactive) {
					console.error(
						"Error: Interactive mode requires a terminal (TTY).",
					);
					console.error(
						"Use: xcsh <command> for non-interactive execution.",
					);
					process.exit(1);
				}

				// Print banner to scrollback BEFORE Ink takes over
				printBannerToScrollback();

				// Enter interactive REPL mode
				// WORKAROUND: Bun doesn't call process.stdin.resume() automatically,
				// which breaks Ink's useInput hook. This is a known Bun bug:
				// https://github.com/oven-sh/bun/issues/6862
				process.stdin.resume();
				render(<App />);
				return;
			}

			// Non-interactive mode: execute command and exit
			await executeNonInteractive(commandArgs);
		},
	);

/**
 * Execute a command non-interactively
 */
async function executeNonInteractive(args: string[]): Promise<void> {
	const session = new REPLSession();
	await session.initialize();

	// Join args into a command string
	const command = args.join(" ");

	// Execute the command
	const result = await executeCommand(command, session);

	// Output results
	result.output.forEach((line) => {
		// eslint-disable-next-line no-console
		console.log(line);
	});

	// Exit with appropriate code
	if (result.error) {
		process.exit(1);
	}

	process.exit(0);
}

// Parse command line arguments
program.parse();
