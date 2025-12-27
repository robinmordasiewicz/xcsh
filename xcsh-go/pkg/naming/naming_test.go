package naming

import (
	"testing"
)

func TestToHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic conversions
		{"simple word", "certificate", "Certificate"},
		{"underscore separated", "origin_pool", "Origin Pool"},
		{"kebab separated", "origin-pool", "Origin Pool"},

		// Uppercase acronyms
		{"dns zone", "dns_zone", "DNS Zone"},
		{"http loadbalancer", "http_loadbalancer", "HTTP Load Balancer"},
		{"api endpoint", "api_endpoint", "API Endpoint"},
		{"tcp proxy", "tcp_proxy", "TCP Proxy"},
		{"udp loadbalancer", "udp_loadbalancer", "UDP Load Balancer"},
		{"tls certificate", "tls_certificate", "TLS Certificate"},
		{"ssl certificate", "ssl_certificate", "SSL Certificate"},
		{"jwt token", "jwt_token", "JWT Token"},
		{"vpn tunnel", "vpn_tunnel", "VPN Tunnel"},
		{"bgp config", "bgp_config", "BGP Config"},
		{"vlan config", "vlan_config", "VLAN Config"},
		{"waf policy", "waf_policy", "WAF Policy"},
		{"cdn config", "cdn_config", "CDN Config"},
		{"aws vpc site", "aws_vpc_site", "AWS VPC Site"},
		{"gcp vpc site", "gcp_vpc_site", "GCP VPC Site"},
		{"k8s cluster", "k8s_cluster", "K8S Cluster"},
		{"oidc config", "oidc_config", "OIDC Config"},
		{"saml config", "saml_config", "SAML Config"},
		{"ldap config", "ldap_config", "LDAP Config"},
		{"rbac policy", "rbac_policy", "RBAC Policy"},
		{"iam policy", "iam_policy", "IAM Policy"},
		{"rest api", "rest_api", "REST API"},
		{"json schema", "json_schema", "JSON Schema"},
		{"xml config", "xml_config", "XML Config"},
		{"cors policy", "cors_policy", "CORS Policy"},

		// F5-specific acronyms
		{"bigip apm", "bigip_apm", "BIG-IP APM"},
		{"bigip ltm", "bigip_ltm", "BIG-IP LTM"},
		{"bigip gtm", "bigip_gtm", "BIG-IP GTM"},
		{"bigip asm", "bigip_asm", "BIG-IP ASM"},
		{"api sec", "api_sec", "API SEC"},
		{"xc config", "xc_config", "XC Config"},

		// Mixed-case acronyms
		{"mtls config", "mtls_config", "mTLS Config"},
		{"oauth provider", "oauth_provider", "OAuth Provider"},
		{"graphql endpoint", "graphql_endpoint", "GraphQL Endpoint"},
		{"websocket config", "websocket_config", "WebSocket Config"},
		{"ipv4 address", "ipv4_address", "IPv4 Address"},
		{"ipv6 address", "ipv6_address", "IPv6 Address"},
		{"irule config", "irule_config", "iRule Config"},

		// Compound words
		{"http loadbalancer", "http_loadbalancer", "HTTP Load Balancer"},
		{"origin pool", "origin_pool", "Origin Pool"},
		{"health check", "health_check", "Health Check"},
		{"service policy", "service_policy", "Service Policy"},

		// Kebab-case inputs
		{"kebab dns zone", "dns-zone", "DNS Zone"},
		{"kebab http lb", "http-loadbalancer", "HTTP Load Balancer"},
		{"kebab api endpoint", "api-endpoint", "API Endpoint"},

		// Complex combinations
		{"full resource name", "http_loadbalancer_service_policy", "HTTP Load Balancer Service Policy"},
		{"multiple acronyms", "dns_tcp_udp_config", "DNS TCP UDP Config"},
		{"api sec api", "api_sec_api_definition", "API SEC API Definition"},

		// Edge cases
		{"empty string", "", ""},
		{"single char", "a", "A"},
		{"already uppercase", "DNS", "DNS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("ToHumanReadable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http_load_balancer", "HTTP Load Balancer"},
		{"dns_zone", "DNS Zone"},
		{"api_endpoint", "API Endpoint"},
		{"simple_name", "Simple Name"},
		{"mtls_config", "mTLS Config"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToTitleCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToTitleCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToTitleCaseFromAnchor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http-load-balancer", "HTTP Load Balancer"},
		{"dns-zone", "DNS Zone"},
		{"api-endpoint", "API Endpoint"},
		{"simple-name", "Simple Name"},
		{"mtls-config", "mTLS Config"},
		{"bigip-apm", "BIG-IP APM"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToTitleCaseFromAnchor(tt.input)
			if result != tt.expected {
				t.Errorf("ToTitleCaseFromAnchor(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HTTPLoadBalancer", "http_load_balancer"},
		{"DNSZone", "dns_zone"},
		{"APIEndpoint", "api_endpoint"},
		{"SimpleName", "simple_name"},
		{"mTLSConfig", "m_tls_config"}, // Note: mTLS becomes m_tls in snake_case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http_loadbalancer", "http-loadbalancer"},
		{"dns_zone", "dns-zone"},
		{"api_endpoint", "api-endpoint"},
		{"simple_name", "simple-name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToKebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeAcronyms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic acronym correction
		{"dns lowercase", "Configure dns settings", "Configure DNS settings"},
		{"api lowercase", "Call the api endpoint", "Call the API endpoint"},
		{"http lowercase", "Use http protocol", "Use HTTP protocol"},
		{"multiple acronyms", "Configure dns and api settings for http", "Configure DNS and API settings for HTTP"},

		// Mixed-case acronyms
		{"mtls", "Enable mtls authentication", "Enable mTLS authentication"},
		{"oauth", "Use oauth for login", "Use OAuth for login"},
		{"graphql", "Query graphql endpoint", "Query GraphQL endpoint"},
		{"ipv4 and ipv6", "Support ipv4 and ipv6", "Support IPv4 and IPv6"},

		// Preserve existing correct casing
		{"already correct", "DNS is configured", "DNS is configured"},
		{"mixed correct", "mTLS and HTTP are enabled", "mTLS and HTTP are enabled"},

		// No acronyms
		{"no acronyms", "This is a simple sentence", "This is a simple sentence"},

		// Idempotency test
		{"idempotent", "DNS and API", "DNS and API"},

		// Edge cases
		{"empty string", "", ""},
		{"single acronym", "dns", "DNS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAcronyms(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeAcronyms(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeAcronymsIdempotent(t *testing.T) {
	// Test that running NormalizeAcronyms multiple times produces the same result
	inputs := []string{
		"Configure dns settings for the api endpoint",
		"DNS is already correct",
		"Enable mtls authentication",
		"HTTP Load Balancer",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			first := NormalizeAcronyms(input)
			second := NormalizeAcronyms(first)
			third := NormalizeAcronyms(second)

			if first != second || second != third {
				t.Errorf("NormalizeAcronyms is not idempotent for %q: first=%q, second=%q, third=%q",
					input, first, second, third)
			}
		})
	}
}

func TestToResourceTypeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http_loadbalancer", "HTTPLoadBalancer"},
		{"dns_zone", "DNSZone"},
		{"api_endpoint", "APIEndpoint"},
		{"simple_name", "SimpleName"},
		{"origin_pool", "OriginPool"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToResourceTypeName(tt.input)
			if result != tt.expected {
				t.Errorf("ToResourceTypeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStartsWithVowel(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"API", true},
		{"HTTP", false},
		{"endpoint", true},
		{"origin", true},
		{"DNS", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := StartsWithVowel(tt.input)
			if result != tt.expected {
				t.Errorf("StartsWithVowel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetArticle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"API", "an"},
		{"HTTP", "a"},
		{"endpoint", "an"},
		{"origin", "an"},
		{"DNS", "a"},
		{"certificate", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetArticle(tt.input)
			if result != tt.expected {
				t.Errorf("GetArticle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsUppercaseAcronym(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"dns", true},
		{"DNS", true},
		{"Dns", true},
		{"http", true},
		{"api", true},
		{"hello", false},
		{"world", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsUppercaseAcronym(tt.input)
			if result != tt.expected {
				t.Errorf("IsUppercaseAcronym(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetMixedCaseAcronym(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mtls", "mTLS"},
		{"MTLS", "mTLS"},
		{"oauth", "OAuth"},
		{"OAUTH", "OAuth"},
		{"bigip", "BIG-IP"},
		{"hello", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetMixedCaseAcronym(tt.input)
			if result != tt.expected {
				t.Errorf("GetMixedCaseAcronym(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCompoundWordHumanReadable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"loadbalancer", "Load Balancer"},
		{"LOADBALANCER", "Load Balancer"},
		{"originpool", "Origin Pool"},
		{"healthcheck", "Health Check"},
		{"hello", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetCompoundWordHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("GetCompoundWordHumanReadable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkToHumanReadable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToHumanReadable("http_loadbalancer")
	}
}

func BenchmarkNormalizeAcronyms(b *testing.B) {
	text := "Configure dns settings for the api endpoint using http protocol"
	for i := 0; i < b.N; i++ {
		NormalizeAcronyms(text)
	}
}
