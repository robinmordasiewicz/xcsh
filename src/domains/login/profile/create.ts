/**
 * login profile create - Create a new profile
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";
import type { Profile } from "../../../profile/index.js";

export const createCommand: CommandDefinition = {
	name: "create",
	description:
		"Create a new connection profile with tenant URL and authentication credentials. Profiles store all settings needed to connect to a tenant, including optional default namespace for scoped operations.",
	descriptionShort: "Create a new connection profile",
	descriptionMedium:
		"Create a profile with tenant URL, authentication token, and optional default namespace.",
	usage: "<name> --url <api-url> --token <api-token> [--namespace <ns>]",
	aliases: ["add", "new"],

	async execute(args, _session) {
		const manager = getProfileManager();

		// Parse arguments
		const name = args[0];
		if (!name || name.startsWith("-")) {
			return errorResult(
				[
					"Usage: login profile create <name> --url <api-url> --token <api-token>",
					"",
					"Options:",
					"  --url       F5 XC API URL (e.g., https://tenant.console.ves.volterra.io)",
					"  --token     API token for authentication",
					"  --namespace Default namespace (optional)",
					"",
					"Example:",
					"  login profile create myprofile --url https://myco.console.ves.volterra.io --token abc123",
				].join("\n"),
			);
		}

		// Check if profile already exists
		const existing = await manager.get(name);
		if (existing) {
			return errorResult(
				`Profile '${name}' already exists. Use 'login profile delete ${name}' first, or choose a different name.`,
			);
		}

		// Parse flags
		let apiUrl = "";
		let apiToken = "";
		let defaultNamespace = "";

		for (let i = 1; i < args.length; i++) {
			const arg = args[i];
			const next = args[i + 1];

			if ((arg === "--url" || arg === "-u") && next) {
				apiUrl = next;
				i++;
			} else if ((arg === "--token" || arg === "-t") && next) {
				apiToken = next;
				i++;
			} else if ((arg === "--namespace" || arg === "-n") && next) {
				defaultNamespace = next;
				i++;
			}
		}

		// Validate required fields
		if (!apiUrl) {
			return errorResult("Missing required --url option");
		}

		if (!apiToken) {
			return errorResult("Missing required --token option");
		}

		// Create profile
		const profile: Profile = {
			name,
			apiUrl,
			apiToken,
		};

		if (defaultNamespace) {
			profile.defaultNamespace = defaultNamespace;
		}

		// Save profile
		const result = await manager.save(profile);

		if (!result.success) {
			return errorResult(result.message);
		}

		return successResult([
			`Profile '${name}' created successfully.`,
			``,
			`Use 'login profile use ${name}' to activate this profile.`,
		]);
	},
};
