// Package naming provides consistent case conversion and acronym handling
// for the f5xcctl CLI tool, ensuring industry-standard acronyms are displayed
// correctly in documentation, help text, and user-facing output.
package naming

// UppercaseAcronyms defines acronyms that should always be uppercase.
// Based on RFC 4949, IEEE standards, and industry style guides (Google, Microsoft, Apple).
var UppercaseAcronyms = map[string]bool{
	// Networking protocols
	"DNS": true, "HTTP": true, "HTTPS": true, "TCP": true, "UDP": true,
	"TLS": true, "SSL": true, "SSH": true, "FTP": true, "SFTP": true,
	"SMTP": true, "IMAP": true, "POP": true, "LDAP": true, "DHCP": true,
	"ARP": true, "ICMP": true, "SNMP": true, "NTP": true, "SIP": true,
	"RTP": true, "RTSP": true, "QUIC": true, "IP": true, "GRPC": true,
	// Web/API
	"API": true, "URL": true, "URI": true, "REST": true, "SOAP": true,
	"JSON": true, "XML": true, "HTML": true, "CSS": true, "CORS": true,
	"CDN": true, "WAF": true, "JWT": true, "SAML": true,
	// Network infrastructure
	"VPN": true, "NAT": true, "VLAN": true, "BGP": true, "OSPF": true,
	"QOS": true, "MTU": true, "TTL": true, "ACL": true, "CIDR": true,
	"VIP": true, "LB": true, "HA": true, "DR": true,
	// Security
	"PKI": true, "CA": true, "CSR": true, "CRL": true, "OCSP": true,
	"PEM": true, "AES": true, "RSA": true, "SHA": true, "MD5": true,
	"HMAC": true, "MFA": true, "SSO": true, "RBAC": true, "IAM": true,
	"DDOS": true, "DOS": true, "XSS": true, "CSRF": true, "SQL": true,
	// Cloud/Infrastructure
	"AWS": true, "GCP": true, "CPU": true, "RAM": true, "SSD": true,
	"HDD": true, "GPU": true, "RAID": true, "VM": true, "OS": true,
	"SLA": true, "RPO": true, "RTO": true, "VPC": true, "VNET": true,
	"TGW": true, "IKE": true, "ID": true, "SLI": true, "S2S": true,
	"RE": true, "CE": true, "SPO": true, "SMG": true,
	"APM": true, "PII": true, "OIDC": true, "K8S": true,
	// F5-specific
	"ASM": true, "LTM": true, "GTM": true, "CNE": true, "XC": true,
	"SSLO": true, "AFM": true, "AVR": true, "ASN": true, "SEC": true,
	"RPC": true,
}

// MixedCaseAcronyms defines acronyms with specific mixed-case conventions.
var MixedCaseAcronyms = map[string]string{
	"mtls":      "mTLS",
	"oauth":     "OAuth",
	"graphql":   "GraphQL",
	"websocket": "WebSocket",
	"iscsi":     "iSCSI",
	"ipv4":      "IPv4",
	"ipv6":      "IPv6",
	"macos":     "macOS",
	"ios":       "iOS",
	"nosql":     "NoSQL",
	"bigip":     "BIG-IP",
	"irule":     "iRule",
}

// CompoundWords defines compound words for Go type name formatting.
var CompoundWords = map[string]string{
	"loadbalancer":  "LoadBalancer",
	"bigip":         "BigIP",
	"websocket":     "WebSocket",
	"fastcgi":       "FastCGI",
	"originpool":    "OriginPool",
	"healthcheck":   "HealthCheck",
	"servicepolicy": "ServicePolicy",
}

// CompoundWordsHumanReadable defines compound words for documentation purposes.
var CompoundWordsHumanReadable = map[string]string{
	"loadbalancer":  "Load Balancer",
	"bigip":         "BIG-IP",
	"websocket":     "WebSocket",
	"fastcgi":       "FastCGI",
	"originpool":    "Origin Pool",
	"healthcheck":   "Health Check",
	"servicepolicy": "Service Policy",
	"apiendpoint":   "API Endpoint",
	"apidefinition": "API Definition",
	"apisecurity":   "API Security",
}

// IsUppercaseAcronym returns true if the given string is a known uppercase acronym.
func IsUppercaseAcronym(s string) bool {
	return UppercaseAcronyms[toUpperASCII(s)]
}

// GetMixedCaseAcronym returns the correct mixed-case form for known acronyms.
// Returns empty string if not a mixed-case acronym.
func GetMixedCaseAcronym(s string) string {
	return MixedCaseAcronyms[toLowerASCII(s)]
}

// GetCompoundWordHumanReadable returns the human-readable form of a compound word.
// Returns empty string if not a known compound word.
func GetCompoundWordHumanReadable(s string) string {
	return CompoundWordsHumanReadable[toLowerASCII(s)]
}

// toUpperASCII converts ASCII string to uppercase without importing strings package.
func toUpperASCII(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// toLowerASCII converts ASCII string to lowercase without importing strings package.
func toLowerASCII(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}
