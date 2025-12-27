/**
 * login profile use - Switch active profile
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const useCommand: CommandDefinition = {
	name: "use",
	description: "Switch to a different profile",
	usage: "<name>",
	aliases: ["switch", "activate"],

	async execute(args, session) {
		const manager = getProfileManager();
		const name = args[0];

		if (!name) {
			return errorResult("Usage: login profile use <name>");
		}

		try {
			// Set active profile
			const result = await manager.setActive(name);

			if (!result.success) {
				return errorResult(result.message);
			}

			// Get the profile to update session
			const profile = await manager.get(name);
			if (profile) {
				// Update session with profile settings
				if (profile.defaultNamespace) {
					session.setNamespace(profile.defaultNamespace);
				}
				// Note: API URL and token will be used when making API calls
				// The session will need to be extended to support this
			}

			return successResult(
				[
					`Switched to profile '${name}'.`,
					profile?.apiUrl ? `  API URL: ${profile.apiUrl}` : "",
					profile?.defaultNamespace
						? `  Namespace: ${profile.defaultNamespace}`
						: "",
				].filter(Boolean),
				true, // contextChanged - prompt should update
			);
		} catch (error) {
			return errorResult(
				`Failed to switch profile: ${error instanceof Error ? error.message : "Unknown error"}`,
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
