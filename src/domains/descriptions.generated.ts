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
 * CLI Title from upstream OpenAPI spec (short description)
 * Extracted at build time from .specs/openapi.json info.title
 * This is the single source of truth - no local enrichment
 */
export const CLI_TITLE_FROM_SPEC: string | null = "F5 Distributed Cloud API";

/**
 * CLI Summary from upstream OpenAPI spec (medium description)
 * Extracted at build time from .specs/openapi.json info.summary
 * This is the single source of truth - no local enrichment
 */
export const CLI_SUMMARY_FROM_SPEC: string | null = "Multi-cloud application services with load balancing, WAF, DNS, and edge infrastructure. Unified platform for security and connectivity.";

/**
 * CLI Description from upstream OpenAPI spec (long description)
 * Extracted at build time from .specs/openapi.json info.description
 * This is the single source of truth - no local enrichment
 */
export const CLI_DESCRIPTION_FROM_SPEC: string | null = "Unified application services across multi-cloud, edge, and hybrid environments. Load balancers with origin pools and health checks for traffic distribution. Web application firewall and bot defense for application protection. DNS zones with geographic routing for name resolution. Cloud sites on AWS, Azure, and GCP for infrastructure deployment. Service policies, network security, and observability dashboards from a single control plane.";


/**
 * Generated Descriptions Data
 * Auto-generated from config/custom-domain-descriptions.yaml
 * Generated at: 2026-01-02T08:07:08.466Z
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "2026-01-02T08:07:08.466Z",

	domains: {
		"login": {
			short: "Manage authentication and session profiles",
			medium: "Configure authentication credentials, manage named profiles for different environments, and control connection context for CLI sessions.",
			long: "Handle identity verification and persistent session state for CLI operations. Use 'show' to display current token status and credential validity. Create named profiles for distinct tenants or environments with the 'profile' command, enabling quick switching without re-entering passwords. Set active endpoint URLs and default namespaces using 'context'. View system announcements with 'banner'. Stored configurations persist securely across sessions.",
			subcommands: {
				"profile": {
					short: "Manage saved credential sets for tenant connections",
					medium: "Store named authentication configurations with URLs and tokens. Switch between environments without re-entering credentials each time.",
					long: "Configure named credential sets that persist tenant URLs and access tokens for quick environment switching. Create entries for development, staging, and production tenants, then activate them as needed. List all saved entries, view their configuration details, set one as active for subsequent commands, or remove entries no longer needed. Credentials are stored securely in your local configuration directory and persist across sessions.",
				},
				"context": {
					short: "Manage default namespace for scoping operations",
					medium: "Configure, display, and switch between namespace settings that scope subsequent commands to organizational boundaries.",
					long: "Manage namespace settings that determine the default scope for resource operations. The active namespace applies when commands do not explicitly provide one, reducing repetitive flags across multiple operations. Use 'show' to display the current active namespace, 'set' to change the default for future commands, and 'list' to view all available namespaces. Settings persist across sessions until explicitly changed, streamlining workflows within a single boundary.",
				},
			},
		},
		"cloudstatus": {
			short: "Check platform status and service incidents",
			medium: "View operational health, active incidents, component availability, and scheduled maintenance windows for infrastructure services.",
			long: "Monitor real-time service availability across infrastructure components. Query active incidents with severity classification and affected systems, inspect individual component states, and review upcoming maintenance schedules. Filter by component type, date range, or impact level. Generate human-readable summaries or structured JSON output for integration with alerting tools and automation pipelines. Designed for operations teams tracking system reliability and coordinating around planned outages.",
		},
		"completion": {
			short: "Generate tab-completion scripts for supported shells",
			medium: "Create shell integration scripts for bash, zsh, and fish that provide intelligent tab-triggered suggestions for commands, flags, and arguments.",
			long: "Output scripts enabling your terminal to suggest valid options when pressing Tab. Supports bash, zsh, and fish with automatic command discovery, argument hints, and flag recommendations. Install the generated output in your shell's designated directory or source it directly in your configuration file (e.g., .bashrc, .zshrc). Once active, Tab displays available subcommands, recognized flags, and expected parameter values as you type, reducing errors and accelerating CLI usage through context-aware assistance.",
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
