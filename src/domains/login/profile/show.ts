/**
 * login profile show - Show profile details
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";
import {
	parseDomainOutputFlags,
	formatKeyValueOutput,
} from "../../../output/domain-formatter.js";

export const showCommand: CommandDefinition = {
	name: "show",
	description:
		"Display detailed configuration for a specific profile. Shows tenant URL, authentication method, and namespace settings with sensitive credentials securely masked for safe viewing.",
	descriptionShort: "Show profile details (masked credentials)",
	descriptionMedium:
		"Display profile configuration with tenant URL, auth type, and masked credentials.",
	usage: "<name>",
	aliases: ["get", "view"],

	async execute(args, session) {
		const { options, remainingArgs } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);
		const manager = getProfileManager();
		const name = remainingArgs[0];

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

			// Handle none format
			if (options.format === "none") {
				return successResult([]);
			}

			// Build data object for unified formatter
			const data: Record<string, string | boolean | undefined> = {
				name: profile.name,
				apiUrl: masked.apiUrl,
				active: isActive,
			};
			if (masked.apiToken) data.apiToken = masked.apiToken;
			if (masked.p12Bundle) data.p12Bundle = masked.p12Bundle;
			if (masked.cert) data.cert = masked.cert;
			if (masked.key) data.key = masked.key;
			if (masked.defaultNamespace)
				data.namespace = masked.defaultNamespace;

			// Use unified formatter for json/yaml/tsv
			if (
				options.format === "json" ||
				options.format === "yaml" ||
				options.format === "tsv"
			) {
				return successResult(
					formatKeyValueOutput(data, {
						...options,
						title: "Profile",
					}),
				);
			}

			// Table format (default) - custom text output
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
