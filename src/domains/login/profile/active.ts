/**
 * login profile (default) - Show active profile details
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const activeCommand: CommandDefinition = {
	name: "active",
	description:
		"Display the currently active profile configuration including tenant URL, authentication method, and default namespace. Shows masked credentials for security while confirming connection settings.",
	descriptionShort: "Display active profile configuration",
	descriptionMedium:
		"Show current active profile details including tenant URL, auth type, and namespace settings.",

	async execute(_args, _session) {
		const manager = getProfileManager();

		try {
			const activeName = await manager.getActive();

			if (!activeName) {
				return successResult([
					"No active profile set.",
					"",
					"Run 'login profile list' to see available profiles.",
					"Run 'login profile use <name>' to activate a profile.",
				]);
			}

			const profile = await manager.get(activeName);

			if (!profile) {
				return errorResult(
					`Active profile '${activeName}' not found. Run 'login profile list' to see available profiles.`,
				);
			}

			const masked = manager.maskProfile(profile);

			const output: string[] = [
				`Active Profile: ${profile.name}`,
				``,
				`  API URL:    ${masked.apiUrl}`,
			];

			if (masked.apiToken) {
				output.push(`  API Token:  ${masked.apiToken}`);
			}

			if (masked.p12Bundle) {
				output.push(`  P12 Bundle: ${masked.p12Bundle}`);
			}

			if (masked.cert) {
				output.push(`  Certificate: ${masked.cert}`);
			}

			if (masked.key) {
				output.push(`  Private Key: ${masked.key}`);
			}

			if (masked.defaultNamespace) {
				output.push(`  Namespace:  ${masked.defaultNamespace}`);
			}

			return successResult(output);
		} catch (error) {
			return errorResult(
				`Failed to get active profile: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};
