/**
 * Show Command
 * Connection and identity information
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getWhoamiInfo } from "./service.js";
import { formatWhoami } from "./formatter.js";

// Re-export types and utilities for external use (e.g., Banner)
export type { WhoamiInfo, WhoamiOptions, DisplayTier } from "./types.js";
export { toDisplayTier } from "./types.js";
export { getWhoamiInfo, getWhoamiInfoBasic } from "./service.js";
export { formatWhoami, formatWhoamiCompact } from "./formatter.js";

/**
 * Show command - displays connection and identity information
 */
export const whoamiCommand: CommandDefinition = {
	name: "show",
	description:
		"Display current connection status and authenticated identity information. Shows active profile, tenant URL, username, tenant details, and session context including namespace targeting.",
	descriptionShort: "Show connection and identity info",
	descriptionMedium:
		"Display active profile, tenant URL, user identity, and current namespace context.",

	async execute(_args, session) {
		try {
			const info = await getWhoamiInfo(session);
			const output = formatWhoami(info);
			return successResult(output);
		} catch (error) {
			return errorResult(
				`Failed to get connection info: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};
