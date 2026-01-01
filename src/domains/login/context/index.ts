/**
 * Context Subcommand Group - Manage default namespace context for API operations
 *
 * Similar to kubectl config use-context, this manages the default namespace
 * that will be used for API operations when no namespace is explicitly specified.
 */

import type { CommandDefinition, SubcommandGroup } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { ENV_PREFIX } from "../../../branding/index.js";
import type { REPLSession } from "../../../repl/session.js";
import {
	parseDomainOutputFlags,
	formatKeyValueOutput,
	formatListOutput,
} from "../../../output/domain-formatter.js";

/**
 * Show command - Display current default namespace
 */
const showCommand: CommandDefinition = {
	name: "show",
	description:
		"Display the currently active namespace context used for scoping operations. Shows both the namespace value and its source (environment variable, profile default, or session configuration).",
	descriptionShort: "Show current default namespace",
	descriptionMedium:
		"Display active namespace context and its configuration source.",
	usage: "",
	aliases: ["current", "get"],

	async execute(args, session) {
		const { options } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);
		const namespace = session.getNamespace();
		const source = determineNamespaceSource(namespace);

		// Handle none format
		if (options.format === "none") {
			return successResult([]);
		}

		// Build data object for unified formatter
		const data = {
			namespace,
			source,
		};

		// Use unified formatter for json/yaml/tsv
		if (
			options.format === "json" ||
			options.format === "yaml" ||
			options.format === "tsv"
		) {
			return successResult(
				formatKeyValueOutput(data, {
					...options,
					title: "Context",
				}),
			);
		}

		// Table format (default) - custom text output
		const lines: string[] = [
			`Current namespace: ${namespace}`,
			`Source: ${source}`,
		];

		return successResult(lines);
	},
};

/**
 * Determine where the namespace value came from
 */
function determineNamespaceSource(namespace: string): string {
	const envNamespace = process.env[`${ENV_PREFIX}_NAMESPACE`];

	if (envNamespace && envNamespace === namespace) {
		return `environment variable (${ENV_PREFIX}_NAMESPACE)`;
	}

	if (namespace === "default") {
		return "default value";
	}

	return "session configuration";
}

/**
 * Set command - Change default namespace
 */
const setCommand: CommandDefinition = {
	name: "set",
	description:
		"Change the default namespace context for all subsequent operations. Updates the session scope so operations target the specified namespace unless explicitly overridden.",
	descriptionShort: "Set default namespace context",
	descriptionMedium:
		"Switch namespace context for scoped operations without specifying namespace each time.",
	usage: "<namespace>",
	aliases: ["use", "switch"],

	async execute(args, session) {
		const namespace = args[0];

		if (!namespace) {
			return errorResult("Usage: login context set <namespace>");
		}

		// Validate namespace format (alphanumeric, dash, underscore)
		if (!/^[a-zA-Z0-9_-]+$/.test(namespace)) {
			return errorResult(
				"Invalid namespace. Use alphanumeric characters, dashes, and underscores only.",
			);
		}

		const previousNamespace = session.getNamespace();
		session.setNamespace(namespace);

		const lines: string[] = [`Namespace context changed.`];
		if (previousNamespace !== namespace) {
			lines.push(`  Previous: ${previousNamespace}`);
			lines.push(`  Current:  ${namespace}`);
		} else {
			lines.push(`  Namespace remains: ${namespace}`);
		}

		return successResult(lines, true); // contextChanged = true to update prompt
	},

	async completion(partial: string, _args: string[], session: REPLSession) {
		// Try to use cached namespaces first
		const cached = session.getNamespaceCache();
		if (cached.length > 0) {
			return cached.filter((ns) =>
				ns.toLowerCase().startsWith(partial.toLowerCase()),
			);
		}

		// Fetch namespaces from API if available
		const client = session.getAPIClient();
		if (client?.isAuthenticated()) {
			try {
				const response = await client.get<{
					items?: Array<{ name?: string }>;
				}>("/api/web/namespaces");
				if (response.ok && response.data?.items) {
					const namespaces = response.data.items
						.map((item) => item.name)
						.filter((name): name is string => !!name)
						.sort();
					// Cache the fetched namespaces
					session.setNamespaceCache(namespaces);
					return namespaces.filter((ns) =>
						ns.toLowerCase().startsWith(partial.toLowerCase()),
					);
				}
			} catch {
				// Fall back to defaults on error
			}
		}
		// Return common namespace patterns as fallback
		const suggestions = ["default", "system", "shared"];
		return suggestions.filter((ns) => ns.startsWith(partial));
	},
};

/**
 * List command - Show available namespaces
 */
const listCommand: CommandDefinition = {
	name: "list",
	description:
		"Fetch and display all available namespaces from the tenant. Requires authenticated connection. Shows current namespace indicator and provides switch command guidance.",
	descriptionShort: "List available namespaces",
	descriptionMedium:
		"Query tenant for available namespaces with current context indicator.",
	usage: "",
	aliases: ["ls"],

	async execute(args, session) {
		const { options } = parseDomainOutputFlags(
			args,
			session.getOutputFormat(),
		);
		const currentNamespace = session.getNamespace();

		// Check if we have API connection
		if (!session.isConnected()) {
			return errorResult(
				"Not connected to F5 XC API. Use 'login profile use <profile>' to connect.",
			);
		}

		// Check if authenticated
		const client = session.getAPIClient();
		if (!client?.isAuthenticated()) {
			return errorResult(
				"Not authenticated. Configure a profile with API token.",
			);
		}

		// Fetch namespaces from API
		try {
			const response = await client.get<{
				items?: Array<{ name?: string }>;
			}>("/api/web/namespaces");

			if (!response.ok || !response.data?.items) {
				return errorResult("Failed to fetch namespaces from API.");
			}

			const namespaces = response.data.items
				.map((item) => item.name)
				.filter((name): name is string => !!name)
				.sort();

			// Cache the namespaces for tab completion
			session.setNamespaceCache(namespaces);

			// Handle none format
			if (options.format === "none") {
				return successResult([]);
			}

			// Build data array for unified formatter
			const data = namespaces.map((ns) => ({
				name: ns,
				current: ns === currentNamespace,
			}));

			// Use unified formatter for json/yaml/tsv
			if (
				options.format === "json" ||
				options.format === "yaml" ||
				options.format === "tsv"
			) {
				return successResult(formatListOutput(data, options));
			}

			// Table format (default) - custom text output
			const lines: string[] = ["Available namespaces:", ""];
			for (const ns of namespaces) {
				if (ns === currentNamespace) {
					lines.push(`  ${ns} (current)`);
				} else {
					lines.push(`  ${ns}`);
				}
			}
			lines.push("");
			lines.push("Use 'login context set <namespace>' to switch.");

			return successResult(lines);
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			return errorResult(`Failed to list namespaces: ${message}`);
		}
	},
};

/**
 * Context subcommand group
 */
export const contextSubcommands: SubcommandGroup = {
	name: "context",
	description:
		"Manage default namespace context for scoping operations. Set, display, and list namespaces to control which namespace is used when no explicit namespace is specified in commands.",
	descriptionShort: "Manage default namespace context",
	descriptionMedium:
		"Set, display, and list namespaces for scoping operations without explicit namespace flags.",
	commands: new Map([
		["show", showCommand],
		["set", setCommand],
		["list", listCommand],
	]),
};
