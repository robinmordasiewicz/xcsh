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
export const CLI_SUMMARY_FROM_SPEC: string | null =
	"Multi-cloud application services with load balancing, WAF, DNS, and edge infrastructure. Unified platform for security and connectivity.";

/**
 * CLI Description from upstream OpenAPI spec (long description)
 * Extracted at build time from .specs/openapi.json info.description
 * This is the single source of truth - no local enrichment
 */
export const CLI_DESCRIPTION_FROM_SPEC: string | null =
	"Unified application services across multi-cloud, edge, and hybrid environments. Load balancers with origin pools and health checks for traffic distribution. Web application firewall and bot defense for application protection. DNS zones with geographic routing for name resolution. Cloud sites on AWS, Azure, and GCP for infrastructure deployment. Service policies, network security, and observability dashboards from a single control plane.";

/**
 * Generated Descriptions Data
 * Auto-generated from config/custom-domain-descriptions.yaml
 * Generated at: 2026-01-02T23:47:33.236Z
 *
 * DO NOT EDIT MANUALLY - Regenerate with: npm run generate:descriptions
 */
export const generatedDescriptions: GeneratedDescriptionsData = {
	version: "1.0.0",
	generatedAt: "2026-01-02T23:47:33.236Z",

	domains: {
		login: {
			short: "Configure authentication profiles and connection context",
			medium: "Store named credential profiles for multiple tenants, switch active connection contexts, verify session state, and display environment banners.",
			long: "Manage authentication lifecycle and identity configuration for CLI operations. Define reusable credential sets targeting different tenants without repeated login prompts. Toggle between saved identities to execute commands against distinct environments. Inspect current session details including endpoint URL and token expiration status. Show informational banners reflecting active configuration. Supports environment variable injection for automated pipelines and scripted workflows requiring non-interactive authentication.",
			subcommands: {
				profile: {
					short: "Manage saved authentication credentials for tenants",
					medium: "Store named credential sets for different environments, switch between active connections, and remove outdated entries.",
					long: "Organize authentication credentials across multiple tenants and environments. Create named entries for production, staging, or development contexts containing URLs and tokens. List available configurations, display stored details, designate which connection to use for subsequent commands, and delete obsolete records. Credentials persist locally, eliminating repeated entry when working across different tenants throughout the day.",
				},
				context: {
					short: "Set default namespace for scoped command execution",
					medium: "Configure, display, and switch the active namespace scope. Commands automatically target this namespace when no explicit --namespace flag is provided.",
					long: "Manage persistent namespace selection that determines where CLI operations execute by default. Subcommands include 'show' to reveal the current setting, 'set' to designate a new active namespace, and 'list' to enumerate all available choices. After configuration, subsequent commands operate within the selected namespace without requiring repetitive flag usage. Settings persist across terminal sessions, reducing overhead during extended work within a single namespace. Remove the stored value to restore explicit per-command namespace specification.",
				},
			},
		},
		cloudstatus: {
			short: "Check platform status and active incidents",
			medium: "Query operational health, view active incidents, track scheduled maintenance windows, and verify component availability across deployment regions.",
			long: "Display real-time service health and infrastructure component status. Retrieve incident details including severity levels, affected systems, and resolution progress. List upcoming maintenance windows for service interruption planning. Verify regional component availability before deployments. Access historical incident records for post-mortem analysis and SLA tracking. Output supports both human-readable summaries and JSON format for monitoring integration.",
		},
		completion: {
			short: "Generate tab-completion scripts for supported shells",
			medium: "Create shell-specific scripts enabling tab-assisted input for commands, subcommands, and flags in bash, zsh, or fish environments.",
			long: "Output completion scripts that integrate with native shell mechanisms for context-aware suggestions. Supports bash, zsh, and fish with format-specific output suitable for sourcing directly or saving to configuration files. After installation, pressing Tab auto-suggests partial command names, displays available options, and recommends valid flag values. Select the subcommand matching your environment to generate properly formatted output.",
		},
		subscription: {
			short: "Subscription and billing management",
			medium: "Manage subscription tier, addon services, quota limits, usage metrics, and billing information.",
			long: "Manage F5 Distributed Cloud subscription, billing, quotas, and usage. View plan details, addon services, resource limits, usage metrics, payment methods, invoices, and generate comprehensive reports. Monitor tenant-level quota utilization and current billing period costs. Access historical usage data and download invoices for accounting purposes.",
			subcommands: {
				plan: {
					short: "Manage subscription plans",
					medium: "View current plan, list available options, and manage plan transitions.",
					long: "View and manage subscription plan information. Show current plan details including tier, features, and limits. List all available plans for comparison. Initiate plan transitions to upgrade or modify subscription level.",
				},
				addon: {
					short: "Manage addon services",
					medium: "List, view, and manage addon service subscriptions and activation status.",
					long: "Manage addon services for your subscription. List available addons with their status and access level. View detailed addon information including features, pricing, and requirements. Check activation status for all subscribed addons. Subscribe to new addons or manage existing subscriptions.",
				},
				quota: {
					short: "View quota limits and usage",
					medium: "Display tenant-level quota limits and current usage with utilization metrics.",
					long: "View tenant-level quota limits and current usage. Monitor resource utilization to avoid quota exhaustion. Display all quota limits defined by your subscription plan. Track current usage against limits with utilization percentages. Identify critical quotas approaching their limits.",
				},
				usage: {
					short: "View usage metrics",
					medium: "Display usage metrics, cost breakdowns, and historical billing data.",
					long: "View usage metrics and cost data for current and historical billing periods. Display current billing period usage with itemized costs and projected totals. Access monthly usage summaries with cost breakdowns. Analyze usage trends and plan for future capacity needs.",
				},
				billing: {
					short: "Manage billing",
					medium: "View and manage payment methods, invoices, and billing details.",
					long: "Manage billing information including payment methods and invoices. List configured payment methods with their status. View invoice history and download invoice PDFs. Manage primary and secondary payment method designations.",
				},
				report: {
					short: "Generate reports",
					medium: "Create detailed subscription reports combining multiple data sources.",
					long: "Generate comprehensive subscription reports for analysis and planning. Create summary reports combining plan details, addon status, quota utilization, usage metrics, and billing information. Export reports in various formats for stakeholder communication and compliance documentation.",
				},
			},
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
