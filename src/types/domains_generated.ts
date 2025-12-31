/**
 * AUTO-GENERATED FILE - DO NOT EDIT
 * Generated from .specs/index.json v1.0.82
 * Run: npx tsx scripts/generate-domains.ts
 */

import type { DomainInfo } from "./domains.js";

/**
 * Spec version used for generation
 */
export const SPEC_VERSION = "1.0.82";

/**
 * Generated domain data from upstream API specifications
 */
export const generatedDomains: Map<string, DomainInfo> = new Map([
	[
		"admin_console_and_ui",
		{
			name: "admin_console_and_ui",
			displayName: "Admin Console And Ui",
			description:
				"Create administrative dashboard building blocks with tailored setup data and view bindings. Organize presentational materials by namespace and fetch them by name or list all available items. Define display parameters, track system object relationships, and maintain consistent portal appearance through centralized resource management workflows.",
			descriptionShort: "Manage static UI components for admin console.",
			descriptionMedium:
				"Deploy and retrieve graphical elements within namespaces. Configure custom startup parameters and view references for display composition.",
			aliases: ["console-ui", "ui-assets", "static-components"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Manage static UI components for admin console",
				"Deploy and retrieve UI assets within namespaces",
				"Configure console interface elements",
				"Manage custom UI component metadata",
			],
			relatedDomains: ["admin", "system"],
		},
	],
	[
		"api",
		{
			name: "api",
			displayName: "Api",
			description:
				"Catalog services automatically to maintain an inventory of operations and their characteristics. Organize related resources by function or ownership through logical groupings. Establish verification procedures that confirm authentication requirements and expected response structures. Link definitions with load balancers for traffic routing decisions. Flag non-standard paths for exclusion from automated scanning. Monitor resource status and metadata throughout deployment zones.",
			descriptionShort: "Discover, catalog, and test service interfaces.",
			descriptionMedium:
				"Define interface groups and discovery policies. Set up verification rules to check security posture and expected patterns across environments.",
			aliases: ["apisec", "api-discovery"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Discover and catalog APIs",
				"Test API security and behavior",
				"Manage API credentials",
				"Define API groups and testing policies",
			],
			relatedDomains: ["waf", "network_security"],
			cliMetadata: {
				quick_start: {
					command:
						"curl $F5XC_API_URL/api/config/namespaces/default/api_catalogs -H 'Authorization: APIToken $F5XC_API_TOKEN'",
					description: "List all API catalogs in default namespace",
					expected_output: "JSON array of API catalog objects",
				},
				common_workflows: [
					{
						name: "Protect API with Security Policy",
						description:
							"Discover and protect APIs with WAF security policies",
						steps: [
							{
								step: 1,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/api_catalogs -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...catalog_config...}'",
								description:
									"Create API catalog for API discovery and documentation",
							},
							{
								step: 2,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/api_definitions -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...api_config...}'",
								description:
									"Create API definition with security enforcement",
							},
						],
						prerequisites: [
							"API endpoints documented",
							"Security policies defined",
							"WAF rules configured",
						],
						expected_outcome:
							"APIs protected, violations logged and blocked",
					},
				],
				troubleshooting: [
					{
						problem: "API traffic blocked by security policy",
						symptoms: [
							"HTTP 403 Forbidden",
							"Requests rejected at edge",
						],
						diagnosis_commands: [
							"curl $F5XC_API_URL/api/config/namespaces/default/api_definitions/{api} -H 'Authorization: APIToken $F5XC_API_TOKEN'",
							"Check security policy enforcement rules",
						],
						solutions: [
							"Review API definition and security policy rules",
							"Adjust rule sensitivity to reduce false positives",
							"Add exception rules for legitimate traffic patterns",
						],
					},
				],
				icon: "🔐",
			},
		},
	],
	[
		"authentication",
		{
			name: "authentication",
			displayName: "Authentication",
			description:
				"F5 Distributed Cloud Authentication API specifications",
			descriptionShort: "Authentication API",
			descriptionMedium:
				"F5 Distributed Cloud Authentication API specifications",
			aliases: ["authn", "oidc", "sso"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Configure authentication mechanisms",
				"Manage OIDC and OAuth providers",
				"Configure SCIM user provisioning",
				"Manage API credentials and access",
				"Configure account signup policies",
			],
			relatedDomains: ["system", "users"],
		},
	],
	[
		"bigip",
		{
			name: "bigip",
			displayName: "Bigip",
			description:
				"Define custom rule-based policies governing routing decisions and request handling. Build organized collections for network ranges, string patterns, and key-value entries. Map cloud services to physical appliances through connector setups. Link identity workflows using access modules. Track performance metrics and coordinate synchronization between components.",
			descriptionShort:
				"Manage iRules, data groups, and virtual servers.",
			descriptionMedium:
				"Configure traffic logic scripts and structured list entries. Establish appliance bindings and access module integrations.",
			aliases: ["f5-bigip", "irule", "ltm"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Platform",
			useCases: [
				"Manage BigIP F5 appliances",
				"Configure iRule scripts",
				"Manage data groups",
				"Integrate BigIP CNE",
			],
			relatedDomains: ["marketplace"],
		},
	],
	[
		"billing_and_usage",
		{
			name: "billing_and_usage",
			displayName: "Billing And Usage",
			description:
				"Set up payment methods with primary and secondary designations for redundancy. Initiate plan transitions between subscription tiers with state tracking. Download invoice PDFs and query custom invoice lists by date range or status. Define quota limits per namespace and monitor current usage against allocated capacity. Swap payment method roles without service interruption.",
			descriptionShort: "Manage subscription plans and payment methods.",
			descriptionMedium:
				"Configure billing transitions and payment processing. Track invoices and monitor resource quota consumption across namespaces.",
			aliases: ["billing-usage", "quotas", "usage-tracking"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Manage subscription plans and billing transitions",
				"Configure payment methods and invoices",
				"Track resource quota usage across namespaces",
				"Monitor usage limits and capacity",
			],
			relatedDomains: ["system", "users"],
		},
	],
	[
		"blindfold",
		{
			name: "blindfold",
			displayName: "Blindfold",
			description:
				"Define policy rules with label matching and combining algorithms. Set up transformers and matchers to control data safeguarding. Track access patterns through timestamped records with scroll queries and date groupings. Retrieve public keys for cryptographic operations and process policy information for decryption workflows.",
			descriptionShort: "Manage secret encryption and policy rules.",
			descriptionMedium:
				"Configure protection policies and access controls for sensitive data. Monitor usage through detailed logs and date-based rollups.",
			aliases: ["bf", "encrypt", "secrets"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure secret policies for encryption",
				"Manage sensitive data encryption",
				"Enforce data protection policies",
			],
			relatedDomains: ["client_side_defense", "certificates"],
		},
	],
	[
		"bot_and_threat_defense",
		{
			name: "bot_and_threat_defense",
			displayName: "Bot And Threat Defense",
			description:
				"Deploy namespace-scoped protection using behavioral analysis and machine learning. Provision dedicated keys for system automation and real-time intelligence feeds. Coordinate detection across protected applications through centralized managers. Configure pre-authentication checks to identify suspicious patterns before they reach backends. Enable adaptive blocking decisions based on risk scoring and historical activity profiles.",
			descriptionShort: "Detect and block automated attacks.",
			descriptionMedium:
				"Create bot defense instances with Shape integration. Set up traffic classification rules and automated response policies for malicious actors.",
			aliases: ["threat-defense", "tpm", "shape-bot"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure bot defense instances per namespace",
				"Manage TPM threat categories for classification",
				"Provision API keys for automated defense systems",
				"Integrate threat intelligence services",
			],
			relatedDomains: ["bot_defense", "shape", "waf"],
		},
	],
	[
		"cdn",
		{
			name: "cdn",
			displayName: "Cdn",
			description:
				"Set up cache eligibility based on headers, cookies, and query parameters. Create expression-based rules with custom TTL settings and path matchers. Deploy load balancers that handle content distribution across origin pools. Monitor access logs and metrics, aggregate performance data, and execute cache purge operations when content updates require immediate invalidation.",
			descriptionShort: "Configure caching rules and load balancing.",
			descriptionMedium:
				"Define cache rules, TTLs, and path matching. Manage load balancers with origin pools and purge operations.",
			aliases: ["cache", "content"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Networking",
			useCases: [
				"Configure CDN load balancing",
				"Manage content delivery network services",
				"Configure caching policies",
				"Manage data delivery and distribution",
			],
			relatedDomains: ["virtual"],
		},
	],
	[
		"ce_management",
		{
			name: "ce_management",
			displayName: "Ce Management",
			description:
				"Define network connectivity parameters including address allocation ranges, dual-stack protocol support, and isolated administrative ports for out-of-band access. Group physical locations under common policy templates for streamlined oversight. Onboard new deployments through secure credential workflows with expiration policies. Execute controlled software transitions featuring pre-flight validation, rollback capabilities, and progress tracking to maintain service continuity.",
			descriptionShort:
				"Manage Customer Edge sites and network interfaces.",
			descriptionMedium:
				"Configure DHCP pools, IPv6 addressing, and dedicated management ports. Handle site tokens with lifecycle controls and software version transitions.",
			aliases: ["ce-mgmt", "edge-management", "ce-lifecycle"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Infrastructure",
			useCases: [
				"Manage Customer Edge site lifecycle",
				"Configure network interfaces and fleet settings",
				"Handle site registration and token workflows",
				"Execute site upgrades with pre-upgrade checks",
			],
			relatedDomains: ["customer_edge", "sites"],
		},
	],
	[
		"certificates",
		{
			name: "certificates",
			displayName: "Certificates",
			description:
				"Create PKI artifacts organizing cryptographic identity materials by namespace for multi-tenant isolation. Deploy keypair bundles with issuer hierarchies for TLS termination. Establish verification anchor collections governing which external parties can authenticate. Maintain deny-lists blocking compromised identities from initiating sessions. Organize resources within independent security boundaries supporting granular access control.",
			descriptionShort:
				"Manage SSL/TLS certificate chains and trusted CAs.",
			descriptionMedium:
				"Configure certificate manifests linking keys to credential bundles. Define trust anchors for validating client authenticity during mutual TLS.",
			aliases: ["cert", "certs", "ssl", "tls"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Security",
			useCases: [
				"Manage SSL/TLS certificates",
				"Configure trusted CAs",
				"Manage certificate revocation lists (CRL)",
				"Configure certificate manifests",
			],
			relatedDomains: ["blindfold", "system"],
		},
	],
	[
		"cloud_infrastructure",
		{
			name: "cloud_infrastructure",
			displayName: "Cloud Infrastructure",
			description:
				"Establish connections to AWS, Azure, and GCP environments with secure authentication and network discovery. Define gateway links, edge site peering, and elastic provisioning workflows. Monitor segment performance and connection health across geographic regions. Create automated VPC attachment policies with intelligent path selection between customer locations and cloud workloads.",
			descriptionShort: "Connect and manage multi-cloud providers.",
			descriptionMedium:
				"Configure cloud provider credentials and VPC attachments. Manage AWS transit gateways, Azure route tables, and cross-cloud connectivity.",
			aliases: ["cloud", "infra", "provider"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Infrastructure",
			useCases: [
				"Connect to cloud providers (AWS, Azure, GCP)",
				"Manage cloud credentials and authentication",
				"Configure cloud connectivity and elastic provisioning",
				"Link and manage cloud regions",
			],
			relatedDomains: ["sites", "customer_edge"],
		},
	],
	[
		"container_services",
		{
			name: "container_services",
			displayName: "Container Services",
			description:
				"Create definitions for applications running on distributed infrastructure. Establish standardized templates controlling resource consumption and disk limits. Set up partitioned execution contexts supporting namespace separation and multi-tenant isolation. Track persistent volume claims and usage metrics. Connect with mesh networking for traffic routing.",
			descriptionShort: "Deploy containerized workloads across sites.",
			descriptionMedium:
				"Run services with simplified orchestration. Define blueprints governing processor and storage allocation.",
			aliases: ["vk8s", "containers", "workloads"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Infrastructure",
			useCases: [
				"Deploy XCCS (Container Services) namespaces for multi-tenant workloads",
				"Manage container workloads with simplified orchestration",
				"Configure distributed edge container deployments",
				"Run containerized applications without full K8s complexity",
			],
			relatedDomains: ["managed_kubernetes", "sites", "service_mesh"],
		},
	],
	[
		"data_and_privacy_security",
		{
			name: "data_and_privacy_security",
			displayName: "Data And Privacy Security",
			description:
				"Set up sensitive data policies that identify and protect personally identifiable information across traffic flows. Create custom data type definitions matching organizational privacy standards and industry regulations. Configure LMA region parameters including Clickhouse, Elastic, and Kafka integrations. Deploy geo-configurations enforcing data residency rules and regional compliance mandates. Monitor detection status through condition tracking and secret management with blindfold encryption.",
			descriptionShort:
				"Configure sensitive data detection and privacy policies.",
			descriptionMedium:
				"Define custom data types for PII classification. Manage LMA regions and geo-configurations to meet regulatory compliance requirements.",
			aliases: ["data-privacy", "pii", "sensitive-data", "lma"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure sensitive data detection policies",
				"Define custom data types for PII classification",
				"Manage LMA region configurations",
				"Integrate geo-configurations for compliance",
			],
			relatedDomains: ["blindfold", "client_side_defense"],
		},
	],
	[
		"data_intelligence",
		{
			name: "data_intelligence",
			displayName: "Data Intelligence",
			description:
				"F5 Distributed Cloud Data Intelligence API specifications",
			descriptionShort: "Data Intelligence API",
			descriptionMedium:
				"F5 Distributed Cloud Data Intelligence API specifications",
			aliases: ["di", "intelligence", "insights"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Operations",
			useCases: [
				"Analyze security and traffic data",
				"Generate intelligent insights from logs",
				"Configure data analytics policies",
			],
			relatedDomains: ["statistics", "observability"],
		},
	],
	[
		"ddos",
		{
			name: "ddos",
			displayName: "Ddos",
			description:
				"Deploy definitions that block IP addresses and network segments from accessing protected resources. Organize by threat type or source classification. Manage secure channels routing suspicious packets for analysis before reaching origin servers. Update status for real-time visibility into active defenses. Add items during attacks and monitor health metrics.",
			descriptionShort:
				"Configure blocking policies and tunnel protection.",
			descriptionMedium:
				"Set up firewall configurations with deny list rules. Filter malicious traffic through inspection points.",
			aliases: ["dos", "ddos-protect"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure DDoS protection policies",
				"Monitor and analyze DDoS threats",
				"Configure infrastructure protection",
			],
			relatedDomains: ["network_security", "virtual"],
		},
	],
	[
		"dns",
		{
			name: "dns",
			displayName: "Dns",
			description:
				"Set up primary and secondary zones with support for A, AAAA, CNAME, CAA, CERT, and AFSDB record types. Define health checks to monitor target availability and enable automatic failover between record destinations. Clone existing domains, import zone configurations from external servers, or export zone files for backup. Track query metrics and request logs to analyze resolution patterns across namespaces.",
			descriptionShort: "Manage zones, records, and load balancing.",
			descriptionMedium:
				"Configure authoritative name services with record sets and health checks. Import zones from BIND files or transfer via AXFR protocol.",
			aliases: ["dns-zone", "zones"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Networking",
			useCases: [
				"Configure DNS load balancing",
				"Manage DNS zones and domains",
				"Configure DNS compliance policies",
				"Manage resource record sets (RRSets)",
			],
			relatedDomains: ["virtual", "network"],
			cliMetadata: {
				quick_start: {
					command:
						"curl $F5XC_API_URL/api/config/namespaces/default/dns_domains -H 'Authorization: APIToken $F5XC_API_TOKEN'",
					description:
						"List all DNS domains configured in default namespace",
					expected_output: "JSON array of DNS domain objects",
				},
				common_workflows: [
					{
						name: "Create DNS Domain",
						description:
							"Configure DNS domain with load balancer backend",
						steps: [
							{
								step: 1,
								command:
									"Create load balancer endpoint first (virtual domain)",
								description:
									"Ensure target load balancer exists",
							},
							{
								step: 2,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/dns_domains -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...dns_config...}'",
								description:
									"Create DNS domain pointing to load balancer",
							},
						],
						prerequisites: [
							"DNS domain registered",
							"Load balancer configured",
							"SOA and NS records prepared",
						],
						expected_outcome:
							"DNS domain in Active status, queries resolving to load balancer",
					},
				],
				troubleshooting: [
					{
						problem: "DNS queries not resolving",
						symptoms: [
							"NXDOMAIN responses",
							"Timeout on DNS queries",
						],
						diagnosis_commands: [
							"curl $F5XC_API_URL/api/config/namespaces/default/dns_domains/{domain} -H 'Authorization: APIToken $F5XC_API_TOKEN'",
							"nslookup {domain} @ns-server",
						],
						solutions: [
							"Verify domain delegation to F5 XC nameservers",
							"Check DNS domain configuration and backend load balancer status",
							"Validate zone file and record configuration",
						],
					},
				],
				icon: "🌐",
			},
		},
	],
	[
		"generative_ai",
		{
			name: "generative_ai",
			displayName: "Generative Ai",
			description:
				"Set up query evaluation and response handling for intelligent assistant workflows. Manage rating collection with positive and negative outcome tracking. Subscribe to data streams for traffic pattern detection and behavioral analysis. Allocate and deallocate IP resources for ML infrastructure. Control feature enablement and token management for telemetry collection paths.",
			descriptionShort: "Access AI assistant queries and feedback.",
			descriptionMedium:
				"Configure machine learning interactions and collect response ratings. Enable flow pattern monitoring through data subscription channels.",
			aliases: ["ai", "genai", "assistant"],
			complexity: "simple" as const,
			isPreview: true,
			requiresTier: "Advanced",
			category: "AI",
			useCases: [
				"Access AI-powered features",
				"Configure AI assistant policies",
				"Enable flow anomaly detection",
				"Manage AI data collection",
			],
			relatedDomains: [],
		},
	],
	[
		"managed_kubernetes",
		{
			name: "managed_kubernetes",
			displayName: "Managed Kubernetes",
			description:
				"Create granular access controls for namespace resources and non-resource URLs. Map permissions to users, groups, or service accounts through binding configurations. Deploy security admission enforcement using baseline, restricted, or privileged profiles. Register private image sources with credential management for secure pulls. Integrate with external managed solutions including EKS, AKS, and GKE infrastructure.",
			descriptionShort:
				"Configure Kubernetes RBAC and pod security policies.",
			descriptionMedium:
				"Define permission boundaries for workload access. Set up private image repositories with authentication for enterprise deployments.",
			aliases: ["mk8s", "appstack", "k8s-mgmt"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Infrastructure",
			useCases: [
				"Manage XCKS (Managed Kubernetes) cluster RBAC and security",
				"Configure pod security policies and admission controllers",
				"Manage container registries for enterprise deployments",
				"Integrate with external Kubernetes clusters (EKS, AKS, GKE)",
			],
			relatedDomains: ["container_services", "sites", "service_mesh"],
		},
	],
	[
		"marketplace",
		{
			name: "marketplace",
			displayName: "Marketplace",
			description:
				"Set up secure tunnel connections using IKEv1/IKEv2 parameters, GRE encapsulation with source and destination addressing, or dedicated link types. Manage DPD keep-alive timers and tunnel termination points for reliable connectivity. Activate purchasable services with namespace-scoped status tracking. Create custom portal widgets for interface integration and configure Cloud Manager instances for Terraform and infrastructure automation workflows.",
			descriptionShort: "Manage third-party integrations and add-ons.",
			descriptionMedium:
				"Configure connector tunnels with IPSec, GRE, or direct links. Deploy purchasable services and portal customizations across namespaces.",
			aliases: ["market", "addons", "extensions"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Platform",
			useCases: [
				"Access third-party integrations and add-ons",
				"Manage marketplace extensions",
				"Configure Terraform and external integrations",
				"Manage TPM policies",
			],
			relatedDomains: ["bigip", "admin"],
		},
	],
	[
		"network",
		{
			name: "network",
			displayName: "Network",
			description:
				"Deploy secure site connectivity using IPsec tunnels with customizable IKE phase settings, encryption algorithms, and DH groups. Configure BGP routing with peer state monitoring, ASN management, and traffic policies. Set up SRv6 segments, IP prefix sets, and subnet definitions. Manage DC cluster groups for data center integration and define routes for traffic steering across distributed infrastructure.",
			descriptionShort:
				"Configure BGP routing, tunnels, and connectivity.",
			descriptionMedium:
				"Manage IPsec tunnels and IKE configurations. Define BGP peers, ASN assignments, and routing policies for site-to-site connections.",
			aliases: ["net", "routing", "bgp"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Networking",
			useCases: [
				"Configure BGP routing and ASN management",
				"Manage IPsec tunnels and IKE phases",
				"Configure network connectors and routes",
				"Manage SRv6 and subnetting",
				"Define segment connections and policies",
				"Configure IP prefix sets",
			],
			relatedDomains: ["virtual", "network_security", "dns"],
		},
	],
	[
		"network_security",
		{
			name: "network_security",
			displayName: "Network Security",
			description:
				"Manage firewall configurations with match criteria and action rules. Create NAT policies using dynamic pools and port configurations for address translation. Define segment connections to isolate traffic between network zones. Configure policy-based routing to direct packets based on source, destination, or protocol attributes. Set up forward proxy policies and access control lists to govern outbound connections.",
			descriptionShort: "Configure firewalls, NAT, and routing policies.",
			descriptionMedium:
				"Define network firewall rules and NAT policies. Set up policy-based routing with segment connections for traffic control.",
			aliases: ["netsec", "nfw"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure network firewall and ACL policies",
				"Manage NAT policies and port forwarding",
				"Configure policy-based routing",
				"Define network segments and policies",
				"Configure forward proxy policies",
			],
			relatedDomains: ["waf", "api", "network"],
		},
	],
	[
		"nginx_one",
		{
			name: "nginx_one",
			displayName: "Nginx One",
			description:
				"Set up load balancing configurations with backend server definitions and routing logic. Create monitoring schedules for availability tracking across distributed nodes. Build request handling pipelines with rate controls and authentication layers. Track instance performance metrics and traffic patterns. Coordinate failover mechanisms using weighted distribution and priority-based selection.",
			descriptionShort:
				"Configure NGINX proxy instances and deployments.",
			descriptionMedium:
				"Manage upstream server pools and health monitors. Define SSL termination rules and connection parameters for gateway endpoints.",
			aliases: ["nginx", "nms", "nginx-plus"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Platform",
			useCases: [
				"Manage NGINX One platform integrations",
				"Configure NGINX Plus instances",
				"Integrate NGINX configuration management",
			],
			relatedDomains: ["marketplace"],
		},
	],
	[
		"object_storage",
		{
			name: "object_storage",
			displayName: "Object Storage",
			description:
				"Deploy binary artifacts and configuration bundles with automatic version tracking and lifecycle policies. Organize content by category including protection signatures, SDK packages, and third-party connector files. Enable time-limited download links for secure distribution without credential exposure. Track revision history for audit trails and support rollback to previous artifact versions when needed.",
			descriptionShort: "Manage stored objects and bucket versioning.",
			descriptionMedium:
				"Create versioned content within tenant buckets. Generate secure access URLs for SDK distributions and application protection resources.",
			aliases: ["storage", "s3", "buckets"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Platform",
			useCases: [
				"Manage object storage services",
				"Configure stored objects and buckets",
				"Manage storage policies",
			],
			relatedDomains: ["marketplace"],
		},
	],
	[
		"observability",
		{
			name: "observability",
			displayName: "Observability",
			description:
				"Set up synthetic monitoring for DNS resolution and HTTP services across AWS regions. Generate health reports with historical trends and summary dashboards. Monitor certificate validity, track response times, and aggregate results by namespace for capacity planning.",
			descriptionShort: "Configure synthetic monitors and health checks.",
			descriptionMedium:
				"Define DNS and HTTP monitors with regional testing. Track certificate expiration and service availability across zones.",
			aliases: ["obs", "monitoring", "synth"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Operations",
			useCases: [
				"Configure synthetic monitoring",
				"Define monitoring and testing policies",
				"Manage observability dashboards",
			],
			relatedDomains: ["statistics", "support"],
		},
	],
	[
		"rate_limiting",
		{
			name: "rate_limiting",
			displayName: "Rate Limiting",
			description:
				"Create rate limiter policies with configurable time periods using seconds, minutes, or hours granularity. Deploy policers and protocol policers to enforce bandwidth constraints across namespaces. Define limit values, burst allowances, and blocking behaviors when thresholds trigger. Integrate with load balancers and security policies for layered traffic management and abuse prevention.",
			descriptionShort: "Configure traffic throttling and policer rules.",
			descriptionMedium:
				"Define request limits and burst thresholds for traffic control. Set up leaky bucket algorithms and block actions for exceeded quotas.",
			aliases: ["ratelimit", "throttle", "policer"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Networking",
			useCases: [
				"Configure rate limiter policies",
				"Manage policer configurations",
				"Control traffic flow and queuing",
			],
			relatedDomains: ["virtual", "network_security"],
		},
	],
	[
		"secops_and_incident_response",
		{
			name: "secops_and_incident_response",
			displayName: "Secops And Incident Response",
			description:
				"Manage incident response workflows that detect and mitigate malicious users automatically. Create rules matching threat levels to actions like blocking, rate limiting, or alerting. Set up mitigation policies per namespace to isolate security responses. Define thresholds for user behavior analysis and configure graduated responses based on severity. Integrate with bot defense and WAF systems for coordinated protection across application layers.",
			descriptionShort: "Configure automated threat mitigation rules.",
			descriptionMedium:
				"Define malicious user detection policies and response actions. Apply blocking or rate limiting based on threat levels.",
			aliases: ["secops", "incident-response", "mitigation"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure automated threat mitigation policies",
				"Define rules for malicious user detection",
				"Manage incident response workflows",
				"Apply blocking or rate limiting to threats",
			],
			relatedDomains: ["bot_defense", "waf", "network_security"],
		},
	],
	[
		"service_mesh",
		{
			name: "service_mesh",
			displayName: "Service Mesh",
			description:
				"Create classifications to organize services and support automatic identification of interconnected components. Set up analysis pipelines to understand patterns and build intelligent routing rules. Define network function virtualization for regional architectures. Configure authentication settings including location, state, and type recognition.",
			descriptionShort: "Configure application types and discovery.",
			descriptionMedium:
				"Manage NFV integrations and workload categories. Enable traffic learning across distributed deployments.",
			aliases: ["mesh", "svc-mesh"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Infrastructure",
			useCases: [
				"Configure service mesh connectivity",
				"Manage endpoint discovery and routing",
				"Configure NFV services",
				"Define application settings and types",
			],
			relatedDomains: [
				"managed_kubernetes",
				"container_services",
				"virtual",
			],
		},
	],
	[
		"shape",
		{
			name: "shape",
			displayName: "Shape",
			description:
				"Set up bot defense infrastructure across namespaces with deployment tracking and status monitoring. Integrate mobile SDK attributes for app shielding and device recognition. Subscribe to threat intelligence services for real-time protection updates. Define cluster states and location-based policies for distributed bot mitigation. Track deployment history and manage policy configurations through centralized infrastructure objects.",
			descriptionShort: "Configure bot defense and threat prevention.",
			descriptionMedium:
				"Deploy bot infrastructure with mobile SDK integration. Manage subscription services and policy enforcement for automated threat protection.",
			aliases: ["shape-sec", "safeap"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure Shape Security policies",
				"Manage bot and threat prevention",
				"Configure SafeAP policies",
				"Enable threat recognition",
			],
			relatedDomains: ["bot_defense", "waf"],
		},
	],
	[
		"sites",
		{
			name: "sites",
			displayName: "Sites",
			description:
				"Deploy edge nodes across AWS, Azure, and GCP with automated provisioning. Configure VPC peering, transit gateway attachments, and VPN tunnel settings. Define virtual groupings with label selectors for policy targeting. Manage Kubernetes cluster integrations and secure mesh deployments. Monitor node health, validate configurations, and set IP prefix allocations.",
			descriptionShort: "Deploy edge nodes across cloud providers.",
			descriptionMedium:
				"Configure AWS, Azure, GCP deployments with VPC integration. Manage transit gateways and VPN tunnels.",
			aliases: ["site", "deployment"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Infrastructure",
			useCases: [
				"Deploy F5 XC across cloud providers (AWS, Azure, GCP)",
				"Manage XCKS (Managed Kubernetes) site deployments (formerly AppStack)",
				"Deploy Secure Mesh sites for networking-focused edge deployments",
				"Integrate external Kubernetes clusters as Customer Edge",
				"Configure AWS VPC, Azure VNet, and GCP VPC sites",
				"Manage virtual sites and site policies",
			],
			relatedDomains: [
				"cloud_infrastructure",
				"customer_edge",
				"managed_kubernetes",
			],
			cliMetadata: {
				quick_start: {
					command:
						"curl $F5XC_API_URL/api/config/namespaces/default/sites -H 'Authorization: APIToken $F5XC_API_TOKEN'",
					description:
						"List all configured sites in default namespace",
					expected_output:
						"JSON array of site objects with deployment status",
				},
				common_workflows: [
					{
						name: "Deploy AWS Cloud Site",
						description:
							"Deploy F5 XC in AWS for traffic management",
						steps: [
							{
								step: 1,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/cloud_credentials -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...aws_credentials...}'",
								description:
									"Create cloud credentials for AWS access",
							},
							{
								step: 2,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/sites -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...site_config...}'",
								description:
									"Create site definition for AWS deployment",
							},
						],
						prerequisites: [
							"AWS account configured",
							"Cloud credentials created",
							"VPC and security groups prepared",
						],
						expected_outcome:
							"Site deployed in AWS, nodes connected and healthy",
					},
				],
				troubleshooting: [
					{
						problem: "Site deployment fails",
						symptoms: [
							"Status: Error",
							"Nodes not coming online",
							"Connectivity issues",
						],
						diagnosis_commands: [
							"curl $F5XC_API_URL/api/config/namespaces/default/sites/{site} -H 'Authorization: APIToken $F5XC_API_TOKEN'",
							"Check site events and node status",
						],
						solutions: [
							"Verify cloud credentials have required permissions",
							"Check VPC and security group configuration",
							"Review site logs for deployment errors",
							"Ensure sufficient cloud resources available",
						],
					},
				],
				icon: "🌍",
			},
		},
	],
	[
		"statistics",
		{
			name: "statistics",
			displayName: "Statistics",
			description:
				"Set up alert policies with custom matchers, label filters, and group-by rules for targeted notifications. Define routing channels via email, webhook, or integration receivers with confirmation and verification workflows. Access flow analytics, historical alert data, and namespace-scoped metrics. Build capacity planning graphs and operational summaries. Observe deployment health and service discovery mapping across distributed environments.",
			descriptionShort: "Monitor alerts, logs, and flow analytics.",
			descriptionMedium:
				"Configure alerting policies and notification receivers. Track service topology, build dashboards, and view site health summaries.",
			aliases: ["stats", "metrics", "logs"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Operations",
			useCases: [
				"Access flow statistics and analytics",
				"Manage alerts and alerting policies",
				"View logs and log receivers",
				"Generate reports and graphs",
				"Track topology and service discovery",
				"Monitor status at sites",
			],
			relatedDomains: ["observability", "support"],
		},
	],
	[
		"support",
		{
			name: "support",
			displayName: "Support",
			description:
				"Open new cases and assign severity ratings based on business impact. Append notes throughout resolution workflows. Mark items as closed or reinstate them if symptoms recur. Execute diagnostic packet captures on deployed sites for network troubleshooting. Handle tax exemption verification through certificate submission.",
			descriptionShort: "Create and track customer tickets.",
			descriptionMedium:
				"Submit requests with file uploads and priority levels. Add comments and escalate critical incidents to engineering teams.",
			aliases: ["tickets", "help-desk"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Operations",
			useCases: [
				"Submit and manage support tickets",
				"Track customer support requests",
				"Access operational support documentation",
			],
			relatedDomains: ["statistics", "observability"],
		},
	],
	[
		"telemetry_and_insights",
		{
			name: "telemetry_and_insights",
			displayName: "Telemetry And Insights",
			description:
				"F5 Distributed Cloud Telemetry And Insights API specifications",
			descriptionShort: "Telemetry And Insights API",
			descriptionMedium:
				"F5 Distributed Cloud Telemetry And Insights API specifications",
			aliases: ["telemetry", "ti"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Operations",
			useCases: [
				"Collect and analyze telemetry data",
				"Generate actionable insights from metrics",
				"Configure telemetry collection policies",
			],
			relatedDomains: ["observability", "statistics"],
		},
	],
	[
		"tenant_and_identity",
		{
			name: "tenant_and_identity",
			displayName: "Tenant And Identity",
			description:
				"Set up granular alert routing for administrative and combined channels with personalized delivery options. Control active login sessions and enforce one-time password resets for security compliance. Define display layouts and avatar images for customized user experiences. Process onboarding access submissions and toggle account management features. Coordinate support ticket attachments and client relationship interactions across managed tenant hierarchies.",
			descriptionShort: "Manage user profiles and session controls.",
			descriptionMedium:
				"Configure OTP resets and admin alert channels. Handle view settings and profile customization for platform participants.",
			aliases: ["tenant-identity", "idm", "user-settings"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Manage user profiles and notification preferences",
				"Configure session controls and OTP settings",
				"Handle identity management operations",
				"Process initial user access requests",
			],
			relatedDomains: ["users", "authentication", "system"],
		},
	],
	[
		"threat_campaign",
		{
			name: "threat_campaign",
			displayName: "Threat Campaign",
			description:
				"F5 Distributed Cloud Threat Campaign API specifications",
			descriptionShort: "Threat Campaign API",
			descriptionMedium:
				"F5 Distributed Cloud Threat Campaign API specifications",
			aliases: ["threats", "campaigns", "threat-intel"],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Security",
			useCases: [
				"Track and analyze threat campaigns",
				"Monitor active threats and attack patterns",
				"Configure threat intelligence integration",
			],
			relatedDomains: ["bot_defense", "ddos"],
		},
	],
	[
		"users",
		{
			name: "users",
			displayName: "Users",
			description:
				"Deploy namespace-scoped access credentials with lifecycle state tracking for secure machine enrollment. Build hierarchical tagging frameworks that enable systematic organization of infrastructure elements. Retrieve automated provisioning payloads for streamlined node initialization. Enable system-level automatic tagging that applies predefined metadata to newly created objects without operator action.",
			descriptionShort: "Manage account tokens and label settings.",
			descriptionMedium:
				"Configure credential issuance and cloud-init provisioning. Establish key-value taxonomies for consistent resource categorization across deployments.",
			aliases: ["user", "accounts", "iam"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Manage user accounts and tokens",
				"Configure user identification",
				"Manage user settings and preferences",
				"Configure implicit and known labels",
			],
			relatedDomains: ["system", "admin"],
		},
	],
	[
		"virtual",
		{
			name: "virtual",
			displayName: "Virtual",
			description:
				"Deploy load balancers across protocols with origin pool management and service discovery. Set up geo-location routing to direct traffic based on client location. Define rate limiter policies to control request volume and protect services from abuse. Configure health checks for origin monitoring and automatic failover. Manage service policies for access control and traffic filtering. Enable malware protection and threat campaign blocking for security enforcement.",
			descriptionShort: "Configure load balancers and origin pools.",
			descriptionMedium:
				"Create HTTP, TCP, and UDP load balancers with origin pools. Define routing rules, health checks, and rate limiting policies.",
			aliases: ["lb", "loadbalancer", "vhost"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Networking",
			useCases: [
				"Configure HTTP/TCP/UDP load balancers",
				"Manage origin pools and services",
				"Configure virtual hosts and routing",
				"Define rate limiter and service policies",
				"Manage geo-location-based routing",
				"Configure proxy and forwarding policies",
				"Manage malware protection and threat campaigns",
				"Configure health checks and endpoint monitoring",
			],
			relatedDomains: ["dns", "service_policy", "network"],
			cliMetadata: {
				quick_start: {
					command:
						"curl $F5XC_API_URL/api/config/namespaces/default/http_loadbalancers -H 'Authorization: APIToken $F5XC_API_TOKEN'",
					description:
						"List all HTTP load balancers in default namespace",
					expected_output:
						"JSON array of load balancer objects with status",
				},
				common_workflows: [
					{
						name: "Create HTTP Load Balancer",
						description:
							"Deploy basic HTTP load balancer with origin pool backend",
						steps: [
							{
								step: 1,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/origin_pools -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...pool_config...}'",
								description:
									"Create backend origin pool with target endpoints",
							},
							{
								step: 2,
								command:
									"curl -X POST $F5XC_API_URL/api/config/namespaces/default/http_loadbalancers -H 'Authorization: APIToken $F5XC_API_TOKEN' -H 'Content-Type: application/json' -d '{...lb_config...}'",
								description:
									"Create HTTP load balancer pointing to origin pool",
							},
						],
						prerequisites: [
							"Active namespace",
							"Origin pool targets reachable",
							"DNS domain configured",
						],
						expected_outcome:
							"Load balancer in Active status, traffic routed to origins",
					},
				],
				troubleshooting: [
					{
						problem:
							"Load balancer shows Configuration Error status",
						symptoms: [
							"Status: Configuration Error",
							"No traffic routing",
							"Requests timeout",
						],
						diagnosis_commands: [
							"curl $F5XC_API_URL/api/config/namespaces/default/http_loadbalancers/{name} -H 'Authorization: APIToken $F5XC_API_TOKEN'",
							"Check origin_pool status and endpoint connectivity",
						],
						solutions: [
							"Verify origin pool targets are reachable from edge",
							"Check DNS configuration and domain propagation",
							"Validate certificate configuration if using HTTPS",
							"Review security policies not blocking traffic",
						],
					},
				],
				icon: "⚖️",
			},
		},
	],
	[
		"vpm_and_node_management",
		{
			name: "vpm_and_node_management",
			displayName: "Vpm And Node Management",
			description:
				"F5 Distributed Cloud Vpm And Node Management API specifications",
			descriptionShort: "Vpm And Node Management API",
			descriptionMedium:
				"F5 Distributed Cloud Vpm And Node Management API specifications",
			aliases: ["vpm", "nodes", "node-mgmt"],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Platform",
			useCases: [
				"Manage Virtual Private Mesh (VPM) configuration",
				"Configure node lifecycle and management",
				"Monitor VPM and node status",
			],
			relatedDomains: ["sites", "system"],
		},
	],
	[
		"waf",
		{
			name: "waf",
			displayName: "Waf",
			description:
				"Set up firewall configurations with attack type settings and violation detection. Create exclusion policies to tune false positives and customize blocking responses. Deploy staged signatures before production release and monitor rule hits through security event metrics. Integrate with virtual hosts for layered protection using AI-based risk blocking and anonymization settings for sensitive data handling.",
			descriptionShort:
				"Configure application firewall rules and bot protection.",
			descriptionMedium:
				"Define security policies for web applications. Manage attack signatures, exclusion rules, and threat detection settings.",
			aliases: ["firewall", "appfw"],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Advanced",
			category: "Security",
			useCases: [
				"Configure web application firewall rules",
				"Manage application security policies",
				"Enable enhanced firewall capabilities",
				"Configure protocol inspection",
			],
			relatedDomains: ["api", "network_security", "virtual"],
		},
	],
]);

/**
 * Total domain count
 */
export const DOMAIN_COUNT = 38;
