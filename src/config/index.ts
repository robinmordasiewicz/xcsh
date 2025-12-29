/**
 * Configuration module exports.
 */

export {
	EnvVarRegistry,
	formatEnvVarsSection,
	formatConfigSection,
	type EnvVar,
} from "./envvars.js";

export {
	loadSettings,
	loadSettingsSync,
	getConfigPath,
	isValidLogoMode,
	getLogoModeDescription,
	DEFAULT_SETTINGS,
	LOGO_MODES,
	type AppSettings,
	type LogoDisplayMode,
	type LogoModeDefinition,
} from "./settings.js";
