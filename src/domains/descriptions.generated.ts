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

export type CliDescriptions = DescriptionTiers;

export interface GeneratedDescriptionsData {
	version: string;
	generatedAt: string;
	cli?: Record<string, CliDescriptions>;
	domains: Record<string, DomainDescriptions>;
}


/**
 * Generated Descriptions Data
 * Auto-generated from config/custom-domain-descriptions.yaml
 * Generated at: 2025-12-30T03:58:36.781Z
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "2025-12-30T03:58:36.781Z",
	cli: {
		"xcsh": {
			short: "Navigate cloud services via interactive shell",
			medium: "Manage multi-tenant connections, execute domain operations across 100+ services, and output results in JSON, YAML, or table format.",
			long: "Interact with cloud services through an intelligent shell environment. Navigate over 100 domain operations using tab completion for commands, flags, and values. Organize multiple tenant connections with named profiles and switch contexts without re-authenticating. Execute commands directly from scripts or explore interactively with history and suggestions. Configure output format as JSON, YAML, or formatted tables. Set behavior through environment variables or persistent profile settings. Generate shell completions for bash, zsh, and fish terminals.",
		},
	},
	domains: {
		"login": {
			short: "Configure session credentials and environment profiles",
			medium: "Set up authentication tokens, organize named profiles for multiple tenants, and switch between target environments for CLI operations.",
			long: "Control authentication and session state across multiple environments. Run 'show' to display current connection details and token status. Organize credentials with 'profile' to maintain separate configurations for development, staging, and production tenants. Switch active targets via 'context' without re-authenticating. Run 'banner' for visual confirmation of the active environment. Profiles persist locally with token-based and certificate authentication support.",
			subcommands: {
				"profile": {
					short: "Manage saved connection configurations for authentication",
					medium: "Store and switch between multiple tenant connection settings. Create, list, and activate named configurations to avoid repeated credential entry.",
					long: "Organize tenant connections as reusable named entries for rapid environment switching. Each entry persists URLs and authentication tokens, removing manual reconfiguration overhead. Available operations: 'list' enumerates saved entries, 'show' reveals configuration details, 'create' registers new connections, 'delete' purges obsolete ones, 'active' identifies the current selection, and 'use' changes context. Ideal for workflows spanning development, staging, and production tiers.",
				},
				"context": {
					short: "Manage default namespace for scoping operations",
					medium: "Configure and display the active namespace used to scope commands. Set, view, or list available namespaces for your session.",
					long: "Control which namespace subsequent operations target by default. Namespaces partition resources and configurations, ensuring commands affect only the intended area. Use 'show' to display the current selection, 'set' to switch contexts, and 'list' to enumerate accessible options. Once configured, the choice persists across commands until explicitly changed, eliminating repeated namespace flag usage.",
				},
			},
		},
		"cloudstatus": {
			short: "Check infrastructure health and active incidents",
			medium: "Query service availability, component health, ongoing incidents, and scheduled maintenance windows across infrastructure regions.",
			long: "Display real-time operational state for infrastructure components and services. Track ongoing incidents with severity levels and resolution timelines. View upcoming maintenance windows affecting particular regions or deployments. Filter results by component type, severity, time range, or operational condition. Retrieve incident history and outage notifications. Commands support monitoring degraded performance indicators and uptime metrics useful for operations teams.",
		},
		"completion": {
			short: "Generate tab-assist scripts for supported shells",
			medium: "Create shell scripts enabling tab-triggered suggestions for commands, subcommands, flags, and option values in bash, zsh, and fish.",
			long: "Output autocomplete functionality for your terminal environment. Bash, zsh, and fish are fully supported with context-aware prompts covering command names, nested subcommands, available flags, and valid argument values. Installation requires sourcing the generated content in your shell's configuration file (.bashrc, .zshrc, or config.fish). Execute the appropriate subcommand to produce the script, then follow shell-specific setup instructions to activate intelligent tab behavior.",
		},
	},
};

/**
 * Get CLI descriptions
 */
export function getCliDescriptions(cliName: string = "xcsh"): CliDescriptions | undefined {
	return generatedDescriptions.cli?.[cliName];
}

/**
 * Get descriptions for a domain
 */
export function getDomainDescriptions(domainName: string): DomainDescriptions | undefined {
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
