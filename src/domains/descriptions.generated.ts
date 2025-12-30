/**
 * Generated Description Types
 * Auto-generated from config/custom-domain-descriptions.yaml
 */

export interface DescriptionTiers {
	short: string;
	medium: string;
	long: string;
}

export type CommandDescriptions = DescriptionTiers;

export interface SubcommandDescriptions extends DescriptionTiers {
	commands?: Record<string, CommandDescriptions>;
}

export interface DomainDescriptions extends DescriptionTiers {
	subcommands?: Record<string, SubcommandDescriptions>;
	commands?: Record<string, CommandDescriptions>;
}

export interface GeneratedDescriptionsData {
	version: string;
	generatedAt: string;
	domains: Record<string, DomainDescriptions>;
}

/**
 * Generated Descriptions Data
 * Auto-generated from config/custom-domain-descriptions.yaml
 * Generated at: 2025-12-30T00:05:36.416Z
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "2025-12-30T00:05:36.416Z",
	domains: {
		login: {
			short: "Configure authentication and session profiles",
			medium: "Store named credential sets, switch between saved configurations, and display current session context including authenticated identity.",
			long: "Control authentication lifecycle for CLI sessions. Create and organize named profiles containing credentials and target configurations. Switch between saved profiles for different environments or accounts. Display current session state including active profile, authenticated identity, and connection details. Subcommands: 'show' displays auth status and connection info, 'profile' manages saved credential sets, 'context' controls active settings, 'banner' configures session messages.",
			subcommands: {
				profile: {
					short: "Manage saved connection profiles",
					medium: "Store named configurations containing tenant URLs and authentication credentials for quick switching between environments.",
					long: "Organize multiple tenant connections as reusable entries that persist locally between sessions. Each profile captures a URL and authentication token, removing the need for repeated credential input when moving across development, staging, and production. Available operations include listing all saved entries, inspecting individual details, adding new configurations, removing obsolete ones, viewing the currently active selection, and activating a different context.",
				},
				context: {
					short: "Manage default namespace for scoping operations",
					medium: "Configure the active namespace determining resource visibility. Show current settings, set new defaults, or list available options.",
					long: "Control the default namespace that scopes CLI operations to resource boundaries. The active namespace determines which resources are visible and accessible when running commands without explicit namespace flags. Use 'show' to display the currently selected namespace, 'set' to change defaults for subsequent operations, or 'list' to enumerate all namespaces you have access to. Setting a default avoids repetitive --namespace flags and reduces command verbosity for focused work.",
				},
			},
		},
		cloudstatus: {
			short: "Check cloud service health and active incidents",
			medium: "Display operational status, component availability, active incidents with severity levels, and scheduled maintenance windows.",
			long: "Query real-time infrastructure health across all service regions. View component-level breakdowns showing availability percentages and degradation states. List ongoing incidents filtered by severity or affected systems. Review upcoming maintenance schedules with expected duration and impact scope. Output formats include JSON for automation pipelines and human-readable tables for dashboards. Commands: 'status' for quick checks, 'summary' for overviews, 'components' for granular details, 'incidents' for active issues, 'maintenance' for planned outages.",
		},
		completion: {
			short: "Create tab-completion scripts for your terminal",
			medium: "Output shell-specific scripts enabling tab-completion of commands, arguments, and flags for bash, zsh, and fish terminals.",
			long: "Create scripts that integrate with your terminal environment to suggest valid options as you type. Supports bash, zsh, and fish with automatic command discovery, argument hints, and flag suggestions. Generated output can be sourced directly in your shell configuration file or saved to the appropriate directory for persistent activation. Tab-based suggestions reduce typing errors and accelerate command entry by presenting available choices contextually.",
		},
		subscription: {
			short: "Manage subscription tiers, addons, and quota limits",
			medium: "Inspect tier details, list addon services, check quota usage, and validate feature availability before deployment.",
			long: "Monitor tenant-level subscription configuration including service tier, enabled addons, and resource allocations. Run pre-deployment checks to verify feature access and quota headroom. Available subcommands: 'overview' for summary display, 'addons' for filtering by state or access level, 'quota' for usage against limits, 'validate' for readiness assessment, and 'activation-status' for enablement verification. Supports JSON output for automation workflows.",
		},
	},
};

/**
 * Get descriptions for a domain
 */
export function getDomainDescriptions(
	domainName: string,
): DomainDescriptions | undefined {
	return generatedDescriptions.domains[domainName];
}

/**
 * Get descriptions for a subcommand within a domain
 */
export function getSubcommandDescriptions(
	domainName: string,
	subcommandName: string,
): SubcommandDescriptions | undefined {
	const domain = generatedDescriptions.domains[domainName];
	return domain?.subcommands?.[subcommandName];
}

/**
 * Get descriptions for a command within a domain or subcommand
 */
export function getCommandDescriptions(
	domainName: string,
	commandName: string,
	subcommandName?: string,
): CommandDescriptions | undefined {
	const domain = generatedDescriptions.domains[domainName];
	if (!domain) return undefined;

	if (subcommandName) {
		const subcommand = domain.subcommands?.[subcommandName];
		return subcommand?.commands?.[commandName];
	}

	return domain.commands?.[commandName];
}
