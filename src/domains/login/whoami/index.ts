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
	description: "Show connection and identity information",

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
