/**
 * Centralized branding information for the CLI.
 * This is the single source of truth for CLI names, descriptions, and branding.
 * Update this file to rebrand the entire application.
 */

import {
	getCliDescriptions,
	CLI_TITLE_FROM_SPEC,
	CLI_DESCRIPTION_FROM_SPEC,
} from "../domains/descriptions.generated.js";

// CLI identification
export const CLI_NAME = "xcsh";
export const CLI_FULL_NAME = "F5 Distributed Cloud Shell";

/**
 * CLI Version
 *
 * Priority:
 * 1. XCSH_VERSION environment variable (set by CI/CD or build scripts)
 * 2. npm package version from package.json
 * 3. Fallback to "dev" for local development
 *
 * The version is determined at build time via tsup's define option,
 * or at runtime via environment variable.
 */

// BUILD_VERSION is replaced at build time by tsup, or undefined if not set
declare const BUILD_VERSION: string | undefined;

function getVersion(): string {
	// Check for build-time injected version
	if (typeof BUILD_VERSION !== "undefined" && BUILD_VERSION) {
		return BUILD_VERSION;
	}

	// Check for runtime environment variable
	if (process.env.XCSH_VERSION) {
		return process.env.XCSH_VERSION;
	}

	// Development fallback
	return "dev";
}

export const CLI_VERSION = getVersion();

// Load CLI descriptions from generated data
const cliDescs = getCliDescriptions();

// CLI descriptions - upstream spec is single source of truth (compiled at build time)
// SHORT/MEDIUM: Use spec title (shorter, suitable for banner)
// LONG: Use spec description (full paragraph)
// Fallback to generated descriptions if spec not available
export const CLI_DESCRIPTION_SHORT =
	CLI_TITLE_FROM_SPEC ?? cliDescs?.short ?? "F5 Distributed Cloud Shell";
export const CLI_DESCRIPTION_MEDIUM =
	CLI_TITLE_FROM_SPEC ?? cliDescs?.medium ?? CLI_DESCRIPTION_SHORT;
export const CLI_DESCRIPTION_LONG =
	CLI_DESCRIPTION_FROM_SPEC ?? cliDescs?.long ?? CLI_DESCRIPTION_MEDIUM;

// Configuration
export const ENV_PREFIX = "F5XC";

// F5 Logo - compact circular logo with F5 text
// Character encoding (edit in VIM - human-readable):
// - ▓ (U+2593) = red circle background (rendered as solid red block)
// - ▒ (U+2592) = red outline elements (lighter red shade)
// - █ (U+2588) = white F5 text
// - (, ), |, _ = red circle outline
export const F5_LOGO = `\
                   ________
              (▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒)
         (▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒)
      (▒▒▓▓▓▓██████████▓▓▓▓█████████████)
    (▒▓▓▓▓██████▒▒▒▒▒███▓▓██████████████▒)
   (▒▓▓▓▓██████▒▓▓▓▓▓▒▒▒▓██▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒)
  (▒▓▓▓▓▓██████▓▓▓▓▓▓▓▓▓██▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒)
 (▒▓▓███████████████▓▓▓▓█████████████▓▓▓▓▓▓▒)
(▒▓▓▓▒▒▒███████▒▒▒▒▒▓▓▓████████████████▓▓▓▓▓▒)
|▒▓▓▓▓▓▓▒██████▓▓▓▓▓▓▓████████████████████▓▓▒|
|▒▓▓▓▓▓▓▓██████▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒██████████▓▒|
(▒▓▓▓▓▓▓▓██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒████████▒▒)
 (▒▓▓▓▓▓▓██████▓▓▓▓▓▓▓███▓▓▓▓▓▓▓▓▓▓▒▒▒████▒▒)
  (▒▓▓▓▓▓██████▓▓▓▓▓▓█████▓▓▓▓▓▓▓▓▓▓▓▓███▒▒)
   (▒▒██████████▓▓▓▓▓▒██████▓▓▓▓▓▓▓▓███▒▒▒)
    (▒▒▒▒▒██████████▓▓▒▒█████████████▒▒▓▒)
      (▒▓▓▒▒▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒)
         (▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒)
              (▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒)`;

// ANSI color codes for terminal output
export const colors = {
	// F5 Brand colors - brightness-adjusted for terminal rendering
	// PNG uses #E4002B but ANSI needs higher values to match visually
	red: "\x1b[38;2;202;38;10m", // F5 Brand Red (adjusted for ANSI)
	boldWhite: "\x1b[1;97m", // Bold bright white
	reset: "\x1b[0m", // Reset to default

	// Status colors
	green: "\x1b[38;2;0;200;83m", // Git clean status
	yellow: "\x1b[38;2;255;193;7m", // Git dirty status
	blue: "\x1b[38;2;33;150;243m", // Git ahead/behind
	dim: "\x1b[2m", // Dimmed text
} as const;

// Repository information
export const REPO_OWNER = "robinmordasiewicz";
export const REPO_NAME = CLI_NAME;
export const REPO_URL = `https://github.com/${REPO_OWNER}/${REPO_NAME}`;
export const DOCS_URL = `https://${REPO_OWNER}.github.io/${REPO_NAME}/`;

/**
 * Colorize text with F5 brand red
 */
export function colorRed(text: string): string {
	return `${colors.red}${text}${colors.reset}`;
}

/**
 * Colorize text as bold white
 */
export function colorBoldWhite(text: string): string {
	return `${colors.boldWhite}${text}${colors.reset}`;
}

/**
 * Colorize text for git clean status
 */
export function colorGreen(text: string): string {
	return `${colors.green}${text}${colors.reset}`;
}

/**
 * Colorize text for git dirty status
 */
export function colorYellow(text: string): string {
	return `${colors.yellow}${text}${colors.reset}`;
}

/**
 * Colorize text for git ahead/behind status
 */
export function colorBlue(text: string): string {
	return `${colors.blue}${text}${colors.reset}`;
}

/**
 * Dim text
 */
export function colorDim(text: string): string {
	return `${colors.dim}${text}${colors.reset}`;
}

// Re-export terminal detection and logo rendering utilities
export {
	detectTerminalCapabilities,
	generateITerm2ImageSequence,
	getTerminalImageSequence,
	type TerminalCapabilities,
	type ITerm2ImageOptions,
} from "./terminal.js";

export {
	renderLogo,
	resolveLogoMode,
	getLogoModeFromEnv,
	hasImageData,
	type LogoRenderOptions,
	type RenderedLogo,
} from "./logo-renderer.js";
