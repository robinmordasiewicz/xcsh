/**
 * AUTO-GENERATED FILE - DO NOT EDIT
 * Generated from .specs/index.json v1.0.62
 * Run: npx tsx scripts/generate-domains.ts
 */

import type { DomainInfo } from "./domains.js";

/**
 * Spec version used for generation
 */
export const SPEC_VERSION = "1.0.62";

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
				"F5 Distributed Cloud Admin Console And Ui API specifications",
			aliases: [],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"api",
		{
			name: "api",
			displayName: "Api",
			description: "F5 Distributed Cloud Api API specifications",
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
			description: "F5 Distributed Cloud Bigip API specifications",
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
				"F5 Distributed Cloud Billing And Usage API specifications",
			aliases: [],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"blindfold",
		{
			name: "blindfold",
			displayName: "Blindfold",
			description: "F5 Distributed Cloud Blindfold API specifications",
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
				"F5 Distributed Cloud Bot And Threat Defense API specifications",
			aliases: [],
			complexity: "moderate" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"cdn",
		{
			name: "cdn",
			displayName: "Cdn",
			description: "F5 Distributed Cloud Cdn API specifications",
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
				"F5 Distributed Cloud Ce Management API specifications",
			aliases: [],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"certificates",
		{
			name: "certificates",
			displayName: "Certificates",
			description: "F5 Distributed Cloud Certificates API specifications",
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
				"F5 Distributed Cloud Cloud Infrastructure API specifications",
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
				"F5 Distributed Cloud Container Services API specifications",
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
				"F5 Distributed Cloud Data And Privacy Security API specifications",
			aliases: [],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"data_intelligence",
		{
			name: "data_intelligence",
			displayName: "Data Intelligence",
			description:
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
			description: "F5 Distributed Cloud Ddos API specifications",
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
			description: "F5 Distributed Cloud Dns API specifications",
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
				"F5 Distributed Cloud Generative Ai API specifications",
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
				"F5 Distributed Cloud Managed Kubernetes API specifications",
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
			description: "F5 Distributed Cloud Marketplace API specifications",
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
			description: "F5 Distributed Cloud Network API specifications",
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
				"F5 Distributed Cloud Network Security API specifications",
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
			description: "F5 Distributed Cloud Nginx One API specifications",
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
				"F5 Distributed Cloud Object Storage API specifications",
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
				"F5 Distributed Cloud Observability API specifications",
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
				"F5 Distributed Cloud Rate Limiting API specifications",
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
				"F5 Distributed Cloud Secops And Incident Response API specifications",
			aliases: [],
			complexity: "simple" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"service_mesh",
		{
			name: "service_mesh",
			displayName: "Service Mesh",
			description: "F5 Distributed Cloud Service Mesh API specifications",
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
			description: "F5 Distributed Cloud Shape API specifications",
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
			description: "F5 Distributed Cloud Sites API specifications",
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
			description: "F5 Distributed Cloud Statistics API specifications",
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
			description: "F5 Distributed Cloud Support API specifications",
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
				"F5 Distributed Cloud Tenant And Identity API specifications",
			aliases: [],
			complexity: "advanced" as const,
			isPreview: false,
			requiresTier: "Standard",
			category: "Other",
			useCases: [],
			relatedDomains: [],
		},
	],
	[
		"threat_campaign",
		{
			name: "threat_campaign",
			displayName: "Threat Campaign",
			description:
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
			description: "F5 Distributed Cloud Users API specifications",
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
			description: "F5 Distributed Cloud Virtual API specifications",
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
			description: "F5 Distributed Cloud Waf API specifications",
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
