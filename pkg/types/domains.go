package types

// DomainInfo contains metadata about a resource domain
type DomainInfo struct {
	Name        string   // Canonical: "load_balancer"
	DisplayName string   // Human: "Load Balancer"
	Description string   // Functional description
	Aliases     []string // Short forms: ["lb"]
}

// DomainRegistry maps canonical domain names to metadata
var DomainRegistry = map[string]*DomainInfo{
	"load_balancer": {
		Name:        "load_balancer",
		DisplayName: "Load Balancer",
		Description: "HTTP, TCP, UDP load balancing and origin pool management",
		Aliases:     []string{"lb"},
	},
	"security": {
		Name:        "security",
		DisplayName: "Security",
		Description: "WAF policies, bot defense, service policies, and threat protection",
		Aliases:     []string{"sec"},
	},
	"networking": {
		Name:        "networking",
		DisplayName: "Networking",
		Description: "Networks, DNS zones, routing, BGP, and virtual connectivity",
		Aliases:     []string{"net"},
	},
	"infrastructure": {
		Name:        "infrastructure",
		DisplayName: "Infrastructure",
		Description: "Sites, fleets, cloud provisioning, and cluster management",
		Aliases:     []string{"infra"},
	},
	"observability": {
		Name:        "observability",
		DisplayName: "Observability",
		Description: "Monitoring, logging, alerts, metrics, and dashboards",
		Aliases:     []string{"obs", "o11y"},
	},
	"identity": {
		Name:        "identity",
		DisplayName: "Identity",
		Description: "Namespaces, users, roles, authentication, and access control",
		Aliases:     []string{"iam"},
	},
	"api_security": {
		Name:        "api_security",
		DisplayName: "API Security",
		Description: "API discovery, testing, crawling, and endpoint protection",
		Aliases:     []string{"apisec"},
	},
	"service_mesh": {
		Name:        "service_mesh",
		DisplayName: "Service Mesh",
		Description: "Service mesh configuration, discovery, and orchestration",
		Aliases:     []string{"mesh"},
	},
	"shape_security": {
		Name:        "shape_security",
		DisplayName: "Shape Security",
		Description: "Advanced bot protection, device identification, and client defense",
		Aliases:     []string{"shape"},
	},
	"infrastructure_protection": {
		Name:        "infrastructure_protection",
		DisplayName: "Infrastructure Protection",
		Description: "DDoS protection, alerts, events, and mitigation policies",
		Aliases:     []string{"infraprot"},
	},
	"applications": {
		Name:        "applications",
		DisplayName: "Applications",
		Description: "Application deployment, management, and configuration",
		Aliases:     []string{"app", "apps"},
	},
	"integrations": {
		Name:        "integrations",
		DisplayName: "Integrations",
		Description: "Third-party integrations and external connectors",
		Aliases:     []string{"int", "connect"},
	},
	"cdn": {
		Name:        "cdn",
		DisplayName: "CDN",
		Description: "Content delivery network, caching, and distribution",
		Aliases:     []string{},
	},
	"bigip": {
		Name:        "bigip",
		DisplayName: "BIG-IP",
		Description: "BIG-IP integration and management",
		Aliases:     []string{},
	},
	"nginx": {
		Name:        "nginx",
		DisplayName: "NGINX",
		Description: "NGINX configuration and management",
		Aliases:     []string{},
	},
	"operations": {
		Name:        "operations",
		DisplayName: "Operations",
		Description: "Operational tasks, workflows, and system operations",
		Aliases:     []string{"ops"},
	},
	"subscriptions": {
		Name:        "subscriptions",
		DisplayName: "Subscriptions",
		Description: "Subscription management and plan administration",
		Aliases:     []string{"sub"},
	},
	"tenant_management": {
		Name:        "tenant_management",
		DisplayName: "Tenant Management",
		Description: "Tenant administration, organization management, and governance",
		Aliases:     []string{"tenant"},
	},
	"billing": {
		Name:        "billing",
		DisplayName: "Billing",
		Description: "Billing, usage tracking, and payment management",
		Aliases:     []string{},
	},
	"vpn": {
		Name:        "vpn",
		DisplayName: "VPN",
		Description: "VPN configuration and remote access management",
		Aliases:     []string{},
	},
	"ai_intelligence": {
		Name:        "ai_intelligence",
		DisplayName: "AI Intelligence",
		Description: "AI and machine learning features and capabilities",
		Aliases:     []string{"ai"},
	},
	"config": {
		Name:        "config",
		DisplayName: "Configuration",
		Description: "System configuration and settings management",
		Aliases:     []string{},
	},
}

// AliasRegistry maps aliases and canonical names to canonical names
var AliasRegistry = map[string]string{}

func init() {
	// Initialize AliasRegistry on startup
	for canonical, info := range DomainRegistry {
		// Map canonical name to itself
		AliasRegistry[canonical] = canonical
		// Map all aliases to canonical name
		for _, alias := range info.Aliases {
			AliasRegistry[alias] = canonical
		}
	}
}

// ResolveDomain converts an alias or canonical name to the canonical domain name
func ResolveDomain(nameOrAlias string) (string, bool) {
	canonical, ok := AliasRegistry[nameOrAlias]
	return canonical, ok
}

// GetDomainInfo retrieves domain metadata by canonical name or alias
func GetDomainInfo(nameOrAlias string) (*DomainInfo, bool) {
	canonical, ok := ResolveDomain(nameOrAlias)
	if !ok {
		return nil, false
	}
	info, found := DomainRegistry[canonical]
	return info, found
}

// AllDomains returns all canonical domain names in sorted order
func AllDomains() []string {
	domains := make([]string, 0, len(DomainRegistry))
	for domain := range DomainRegistry {
		domains = append(domains, domain)
	}
	return domains
}
