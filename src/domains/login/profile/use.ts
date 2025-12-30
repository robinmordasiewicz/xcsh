/**
 * login profile use - Switch active profile
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const useCommand: CommandDefinition = {
	name: "use",
	description:
		"Switch the active connection to a different saved profile. Updates the session context to use the new profile's tenant URL, credentials, and default namespace for all subsequent operations.",
	descriptionShort: "Switch to a different profile",
	descriptionMedium:
		"Activate a profile and update session context with its tenant and namespace settings.",
	usage: "<name>",
	aliases: ["switch", "activate"],

	async execute(args, session) {
		const name = args[0];

		if (!name) {
			return errorResult("Usage: login profile use <name>");
		}

		try {
			// Use session.switchProfile() which properly updates the API client
			const success = await session.switchProfile(name);

			if (!success) {
				return errorResult(`Profile '${name}' not found.`);
			}

			// Get the profile for display
			const profile = session.getActiveProfile();

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
