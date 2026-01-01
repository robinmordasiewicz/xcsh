/**
 * login profile list - List all saved profiles
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";
import {
	parseDomainOutputFlags,
	formatListOutput,
} from "../../../output/domain-formatter.js";

export const listCommand: CommandDefinition = {
	name: "list",
	description:
		"Display all saved connection profiles with their tenant URLs and authentication types. Highlights the currently active profile for easy identification when managing multiple tenants.",
	descriptionShort: "List all saved profiles",
	descriptionMedium:
		"Show all profiles with tenant URLs, auth types, and active status indicator.",
	aliases: ["ls"],

	async execute(args, session) {
		const { options } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);
		const manager = getProfileManager();

		try {
			const profiles = await manager.list();
			const activeProfile = await manager.getActive();

			if (profiles.length === 0) {
				// Handle none format
				if (options.format === "none") {
					return successResult([]);
				}
				return successResult([
					"No profiles configured.",
					"",
					"Create a profile with: login profile create <name>",
				]);
			}

			// Build data array for unified formatter
			const data = profiles.map((profile) => {
				const isActive = profile.name === activeProfile;
				const authType = profile.apiToken
					? "token"
					: profile.cert
						? "cert"
						: profile.p12Bundle
							? "p12"
							: "none";
				return {
					name: profile.name,
					apiUrl: profile.apiUrl,
					authType,
					active: isActive,
				};
			});

			// Handle none format
			if (options.format === "none") {
				return successResult([]);
			}

			// Use unified formatter for json/yaml/tsv
			if (
				options.format === "json" ||
				options.format === "yaml" ||
				options.format === "tsv"
			) {
				return successResult(formatListOutput(data, options));
			}

			// Table format (default) - custom text output
			const output: string[] = ["Saved profiles:"];

			for (const item of data) {
				const marker = item.active ? " [active]" : "";
				output.push(
					`  ${item.name}${marker} - ${item.apiUrl} (${item.authType})`,
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
