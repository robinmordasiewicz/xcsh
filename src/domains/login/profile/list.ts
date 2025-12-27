/**
 * login profile list - List all saved profiles
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const listCommand: CommandDefinition = {
	name: "list",
	description: "List all saved profiles",
	aliases: ["ls"],

	async execute(_args, _session) {
		const manager = getProfileManager();

		try {
			const profiles = await manager.list();
			const activeProfile = await manager.getActive();

			if (profiles.length === 0) {
				return successResult([
					"No profiles configured.",
					"",
					"Create a profile with: login profile create <name>",
				]);
			}

			const output: string[] = ["Saved profiles:"];

			for (const profile of profiles) {
				const isActive = profile.name === activeProfile;
				const marker = isActive ? " [active]" : "";
				const authType = profile.apiToken
					? "token"
					: profile.cert
						? "cert"
						: profile.p12Bundle
							? "p12"
							: "none";
				output.push(
					`  ${profile.name}${marker} - ${profile.apiUrl} (${authType})`,
				);
			}

			return successResult(output);
		} catch (error) {
			return errorResult(
				`Failed to list profiles: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};
