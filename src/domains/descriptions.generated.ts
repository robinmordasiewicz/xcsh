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
 * Generated at: 2026-01-02T03:29:16.852Z
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "2026-01-02T03:29:16.852Z",
	cli: {
		xcsh: {
			short: "Configure multi-cloud application delivery services",
			medium: "Manage load balancers, security policies, DNS routing, and cloud site provisioning across hybrid environments.",
			long: "Deploy distributed application infrastructure spanning AWS, Azure, GCP, and edge locations. Create origin pools with health monitoring and traffic management. Apply WAF rules, bot mitigation, and DDoS protection. Set up DNS zones supporting geographic and weighted routing strategies. Establish network segmentation and service mesh configurations. Monitor performance through observability dashboards. Validate subscription quotas before resource creation. Generate JSON output for CI/CD pipeline integration.",
		},
	},
	domains: {
		login: {
			short: "Configure credentials, profiles, and session context",
			medium: "Set up authentication credentials, organize named profiles for different environments, and control which connection applies to subsequent operations.",
			long: "Handle authentication and session configuration for CLI operations. Create named profiles to maintain separate credentials for development, staging, and production environments. View current authentication status with 'show', organize credentials using 'profile', switch active connections with 'context', and display session information via 'banner'. Profiles store tokens and target URLs, enabling quick environment switching without re-entering credentials.",
			subcommands: {
				profile: {
					short: "Manage saved authentication credentials for tenants",
					medium: "Store named configurations with URLs and tokens for quick switching between staging, production, and development environments.",
					long: "Organize credentials and tenant configurations in named entries for streamlined environment management. Create entries to store connection URLs and tokens, list available configurations, switch between active contexts, and remove outdated items. Settings persist across sessions, eliminating repeated credential entry when working with multiple tenants. Use the active command to display the current selection, or switch contexts with the use command for immediate environment changes.",
				},
				context: {
					short: "Manage default namespace scope for CLI operations",
					medium: "Configure and switch between namespace settings that control the default boundary for all commands.",
					long: "Organize resources into logical groupings through namespace scope management. Setting a default namespace eliminates the need to pass --namespace on every command, streamlining workflows when operating within a single boundary for extended periods. Use 'show' to display current settings, 'set' to change the active namespace, and 'list' to view available options. These preferences persist across sessions until explicitly changed.",
				},
			},
		},
		cloudstatus: {
			short: "Track service health and operational incidents",
			medium: "Display real-time platform availability, active incident reports, scheduled maintenance windows, and component status across global regions.",
			long: "Query operational health for infrastructure components and services. Retrieve active incident details with severity levels and resolution timelines. List upcoming maintenance windows and generate status summaries. Filter results by component type, geographic region, or time range. Output formats support automation workflows and alerting integrations. Commands include status overview, component breakdown, incident history, and maintenance schedules.",
		},
		completion: {
			short: "Generate tab-completion scripts for popular shells",
			medium: "Create autocompletion support for bash, zsh, and fish terminals enabling tab-triggered command and flag discovery.",
			long: "Build shell integration scripts that activate tab-triggered suggestions in your terminal environment. Supports bash, zsh, and fish with automatic command discovery, flag hints, and argument recommendations. Install the generated output in your shell's configuration directory to enable this functionality. Each supported shell includes tailored setup instructions and behaviors matching its native conventions for seamless integration.",
		},
	},
};

/**
 * Get CLI descriptions
 */
export function getCliDescriptions(
	cliName: string = "xcsh",
): CliDescriptions | undefined {
	return generatedDescriptions.cli?.[cliName];
}

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
