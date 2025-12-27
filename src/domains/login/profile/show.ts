/**
 * login profile show - Show profile details
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const showCommand: CommandDefinition = {
	name: "show",
	description: "Show profile details (sensitive data masked)",
	usage: "<name>",
	aliases: ["get", "view"],

	async execute(args, _session) {
		const manager = getProfileManager();
		const name = args[0];

		if (!name) {
			return errorResult("Usage: login profile show <name>");
		}

		try {
			const profile = await manager.get(name);

			if (!profile) {
				return errorResult(`Profile '${name}' not found.`);
			}

			const masked = manager.maskProfile(profile);
			const activeProfile = await manager.getActive();
			const isActive = profile.name === activeProfile;

			const output: string[] = [
				`Profile: ${profile.name}${isActive ? " [active]" : ""}`,
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
				`Failed to show profile: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},

	async completion(partial, _args, _session) {
		const manager = getProfileManager();
		const profiles = await manager.list();
		return profiles
			.map((p) => p.name)
			.filter((name) =>
				name.toLowerCase().startsWith(partial.toLowerCase()),
			);
	},
};
