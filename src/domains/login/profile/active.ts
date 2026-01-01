/**
 * login profile (default) - Show active profile details
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";
import {
	parseDomainOutputFlags,
	formatKeyValueOutput,
} from "../../../output/domain-formatter.js";

export const activeCommand: CommandDefinition = {
	name: "active",
	description:
		"Display the currently active profile configuration including tenant URL, authentication method, and default namespace. Shows masked credentials for security while confirming connection settings.",
	descriptionShort: "Display active profile configuration",
	descriptionMedium:
		"Show current active profile details including tenant URL, auth type, and namespace settings.",

	async execute(args, session) {
		const { options } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);
		const manager = getProfileManager();

		try {
			const activeName = await manager.getActive();

			if (!activeName) {
				// Handle none format
				if (options.format === "none") {
					return successResult([]);
				}
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

			// Handle none format
			if (options.format === "none") {
				return successResult([]);
			}

			// Build data object for unified formatter
			const data: Record<string, string | undefined> = {
				name: profile.name,
				apiUrl: masked.apiUrl,
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
						title: "Active Profile",
					}),
				);
			}

			// Table format (default) - custom text output
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
