/**
 * Show Command
 * Connection and identity information
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getWhoamiInfo } from "./service.js";
import { formatWhoami } from "./formatter.js";
import {
	parseDomainOutputFlags,
	formatKeyValueOutput,
} from "../../../output/domain-formatter.js";
import type { WhoamiInfo } from "./types.js";

// Re-export types and utilities for external use (e.g., Banner)
export type { WhoamiInfo, WhoamiOptions } from "./types.js";
export { getWhoamiInfo, getWhoamiInfoBasic } from "./service.js";
export { formatWhoami, formatWhoamiCompact } from "./formatter.js";

/**
 * Convert WhoamiInfo to key-value data for unified formatter
 */
function whoamiToKeyValue(
	info: WhoamiInfo,
): Record<string, string | boolean | undefined> {
	const data: Record<string, string | boolean | undefined> = {};

	if (info.tenant) data.tenant = info.tenant;
	if (info.email) {
		data.user = info.email;
	} else if (info.username) {
		data.user = info.username;
	}
	data.namespace = info.namespace;
	data.serverUrl = info.serverUrl;
	data.isAuthenticated = info.isAuthenticated;

	return data;
}

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

	async execute(args, session) {
		// Parse output format from args
		const { options } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);

		try {
			const info = await getWhoamiInfo(session);

			// Handle none format - return empty
			if (options.format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (
				options.format === "json" ||
				options.format === "yaml" ||
				options.format === "tsv"
			) {
				const keyValueData = whoamiToKeyValue(info);
				return successResult(
					formatKeyValueOutput(keyValueData, {
						...options,
						title: "Connection Info",
					}),
				);
			}

			// Table format (default) - use the original beautiful box formatter
			const output = formatWhoami(info);
			return successResult(output);
		} catch (error) {
			return errorResult(
				`Failed to get connection info: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};
