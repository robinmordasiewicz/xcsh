/**
 * Application settings from config.yaml file.
 * Handles loading, parsing, and validation of user configuration.
 *
 * Location: ~/.config/xcsh/config.yaml (XDG Base Directory compliant)
 */

import { promises as fs } from "fs";
import YAML from "yaml";
import { paths } from "./paths.js";

/**
 * Logo display mode definition with descriptions.
 * Single source of truth for mode names and their descriptions.
 * Used for both validation and help text generation.
 */
export interface LogoModeDefinition {
	/** Mode identifier */
	mode: string;
	/** Human-readable description */
	description: string;
}

/**
 * All supported logo display modes with descriptions.
 * This is the canonical definition - used for validation and help generation.
 */
export const LOGO_MODES: readonly LogoModeDefinition[] = [
	{
		mode: "image",
		description: "Image if terminal supports, else ASCII (default)",
	},
	{ mode: "ascii", description: "ASCII art only" },
	{ mode: "none", description: "No logo" },
] as const;

/**
 * Logo display mode type - derived from LOGO_MODES for type safety.
 */
export type LogoDisplayMode = (typeof LOGO_MODES)[number]["mode"];

/**
 * Helper string for help text generation
 * Automatically derived from LOGO_MODES constant
 */
export const LOGO_MODE_HELP = LOGO_MODES.map((m) => m.mode).join(", ");

/**
 * Application settings from .xcshconfig file.
 */
export interface AppSettings {
	/** Logo display mode */
	logo: LogoDisplayMode;
}

/**
 * Default settings used when config file doesn't exist or values are missing.
 */
export const DEFAULT_SETTINGS: AppSettings = {
	logo: "image",
};

/**
 * Validate if a string is a valid logo display mode.
 * Uses LOGO_MODES as single source of truth.
 */
export function isValidLogoMode(mode: string): mode is LogoDisplayMode {
	return LOGO_MODES.some((m) => m.mode === mode);
}

/**
 * Get logo mode description for help text.
 */
export function getLogoModeDescription(mode: string): string | undefined {
	return LOGO_MODES.find((m) => m.mode === mode)?.description;
}

/**
 * Validate and sanitize settings from config file.
 * Invalid values are ignored and defaults are used.
 */
function validateSettings(
	settings: Partial<AppSettings>,
): Partial<AppSettings> {
	const validated: Partial<AppSettings> = {};

	if (settings.logo && isValidLogoMode(settings.logo)) {
		validated.logo = settings.logo;
	}

	return validated;
}

/**
 * Load settings from config.yaml file.
 *
 * File format: YAML
 * Location: ~/.config/xcsh/config.yaml (XDG Base Directory compliant)
 *
 * Example:
 * ```yaml
 * # F5 Distributed Cloud Shell Configuration
 * logo: image  # image | ascii | none
 * ```
 *
 * @returns Merged settings with defaults for missing values
 */
export async function loadSettings(): Promise<AppSettings> {
	const configPath = paths.settings;

	try {
		const content = await fs.readFile(configPath, "utf-8");
		const parsed = YAML.parse(content) as Partial<AppSettings>;

		return {
			...DEFAULT_SETTINGS,
			...validateSettings(parsed),
		};
	} catch {
		// File doesn't exist or is invalid - use defaults
		return DEFAULT_SETTINGS;
	}
}

/**
 * Load settings synchronously (for non-async contexts).
 * Uses defaults if file doesn't exist or can't be read.
 */
export function loadSettingsSync(): AppSettings {
	const configPath = paths.settings;

	try {
		// eslint-disable-next-line @typescript-eslint/no-require-imports
		const content = require("fs").readFileSync(
			configPath,
			"utf-8",
		) as string;
		const parsed = YAML.parse(content) as Partial<AppSettings>;

		return {
			...DEFAULT_SETTINGS,
			...validateSettings(parsed),
		};
	} catch {
		return DEFAULT_SETTINGS;
	}
}

/**
 * Get the path to the config file.
 */
export function getConfigPath(): string {
	return paths.settings;
}
