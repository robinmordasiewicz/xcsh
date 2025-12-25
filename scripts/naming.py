"""
Naming utilities for consistent case conversion and acronym handling.

This module provides functions to convert resource names to human-readable format
with proper acronym capitalization, matching the Go pkg/naming package.

Usage:
    from naming import to_human_readable, normalize_acronyms

    to_human_readable("http_loadbalancer")  # "HTTP Load Balancer"
    normalize_acronyms("Configure dns settings")  # "Configure DNS settings"
"""

import re

# UppercaseAcronyms defines acronyms that should always be uppercase.
# Based on RFC 4949, IEEE standards, and industry style guides.
UPPERCASE_ACRONYMS: set[str] = {
    # Networking protocols
    "DNS",
    "HTTP",
    "HTTPS",
    "TCP",
    "UDP",
    "TLS",
    "SSL",
    "SSH",
    "FTP",
    "SFTP",
    "SMTP",
    "IMAP",
    "POP",
    "LDAP",
    "DHCP",
    "ARP",
    "ICMP",
    "SNMP",
    "NTP",
    "SIP",
    "RTP",
    "RTSP",
    "QUIC",
    "IP",
    "GRPC",
    # Web/API
    "API",
    "URL",
    "URI",
    "REST",
    "SOAP",
    "JSON",
    "XML",
    "HTML",
    "CSS",
    "CORS",
    "CDN",
    "WAF",
    "JWT",
    "SAML",
    # Network infrastructure
    "VPN",
    "NAT",
    "VLAN",
    "BGP",
    "OSPF",
    "QOS",
    "MTU",
    "TTL",
    "ACL",
    "CIDR",
    "VIP",
    "LB",
    "HA",
    "DR",
    # Security
    "PKI",
    "CA",
    "CSR",
    "CRL",
    "OCSP",
    "PEM",
    "AES",
    "RSA",
    "SHA",
    "MD5",
    "HMAC",
    "MFA",
    "SSO",
    "RBAC",
    "IAM",
    "DDOS",
    "DOS",
    "XSS",
    "CSRF",
    "SQL",
    # Cloud/Infrastructure
    "AWS",
    "GCP",
    "CPU",
    "RAM",
    "SSD",
    "HDD",
    "GPU",
    "RAID",
    "VM",
    "OS",
    "SLA",
    "RPO",
    "RTO",
    "VPC",
    "VNET",
    "TGW",
    "IKE",
    "ID",
    "SLI",
    "S2S",
    "RE",
    "CE",
    "SPO",
    "SMG",
    "APM",
    "PII",
    "OIDC",
    "K8S",
    # F5-specific
    "ASM",
    "LTM",
    "GTM",
    "CNE",
    "XC",
    "SSLO",
    "AFM",
    "AVR",
    "ASN",
    "SEC",
    "RPC",
}

# MixedCaseAcronyms defines acronyms with specific mixed-case conventions.
MIXED_CASE_ACRONYMS: dict[str, str] = {
    "mtls": "mTLS",
    "oauth": "OAuth",
    "graphql": "GraphQL",
    "websocket": "WebSocket",
    "iscsi": "iSCSI",
    "ipv4": "IPv4",
    "ipv6": "IPv6",
    "macos": "macOS",
    "ios": "iOS",
    "nosql": "NoSQL",
    "bigip": "BIG-IP",
    "irule": "iRule",
}

# CompoundWordsHumanReadable defines compound words for documentation purposes.
COMPOUND_WORDS_HUMAN_READABLE: dict[str, str] = {
    "loadbalancer": "Load Balancer",
    "bigip": "BIG-IP",
    "websocket": "WebSocket",
    "fastcgi": "FastCGI",
    "originpool": "Origin Pool",
    "healthcheck": "Health Check",
    "servicepolicy": "Service Policy",
    "apiendpoint": "API Endpoint",
    "apidefinition": "API Definition",
    "apisecurity": "API Security",
}

# Pre-compiled regex for word boundary matching
_WORD_REGEX = re.compile(r"\b([A-Za-z0-9]+)\b")


def to_human_readable(s: str) -> str:
    """
    Convert a snake_case or kebab-case name to human-readable format
    with proper acronym capitalization and compound word spacing.

    Examples:
        "http_loadbalancer" -> "HTTP Load Balancer"
        "dns-zone" -> "DNS Zone"
        "bigip_apm" -> "BIG-IP APM"
        "mtls_config" -> "mTLS Config"

    Args:
        s: Input string in snake_case or kebab-case format

    Returns:
        Human-readable string with proper capitalization
    """
    if not s:
        return ""

    # Normalize separators: replace both underscores and hyphens with spaces
    s = s.replace("_", " ").replace("-", " ")

    parts = s.split()
    result = []

    for part in parts:
        lower = part.lower()
        upper = part.upper()

        # Check for uppercase acronyms first (e.g., DNS, HTTP, API)
        if upper in UPPERCASE_ACRONYMS:
            result.append(upper)
        # Handle compound words with spaces (e.g., "loadbalancer" -> "Load Balancer")
        elif lower in COMPOUND_WORDS_HUMAN_READABLE:
            result.append(COMPOUND_WORDS_HUMAN_READABLE[lower])
        # Handle mixed-case acronyms (e.g., "mtls" -> "mTLS")
        elif lower in MIXED_CASE_ACRONYMS:
            result.append(MIXED_CASE_ACRONYMS[lower])
        # Standard title case: capitalize first letter
        elif part:
            result.append(part[0].upper() + part[1:].lower())

    return " ".join(result)


def to_title_case(s: str) -> str:
    """
    Convert a snake_case or dot.separated string to Title Case,
    preserving acronym capitalization.

    Example: "http_load_balancer" -> "HTTP Load Balancer"

    Args:
        s: Input string

    Returns:
        Title case string with proper acronym handling
    """
    # Replace underscores, dots, and hyphens with spaces
    s = s.replace("_", " ").replace(".", " ").replace("-", " ")

    words = s.split()
    result = []
    for word in words:
        if word:
            # Apply standard title case first
            result.append(word[0].upper() + word[1:].lower())

    # Apply acronym normalization
    return normalize_acronyms(" ".join(result))


def to_title_case_from_anchor(anchor: str) -> str:
    """
    Convert an anchor name (kebab-case) to Title Case,
    preserving acronym capitalization.

    Example: "http-load-balancer" -> "HTTP Load Balancer"

    Args:
        anchor: Input string in kebab-case format

    Returns:
        Title case string with proper acronym handling
    """
    words = anchor.split("-")
    result = []

    for word in words:
        upper = word.upper()
        lower = word.lower()

        if upper in UPPERCASE_ACRONYMS:
            result.append(upper)
        elif lower in MIXED_CASE_ACRONYMS:
            result.append(MIXED_CASE_ACRONYMS[lower])
        elif word:
            result.append(word[0].upper() + word[1:].lower())

    return " ".join(result)


def normalize_acronyms(text: str) -> str:
    """
    Correct acronym capitalization in free text.
    This function is idempotent - running it multiple times produces the same result.

    Example: "Configure dns settings for the api endpoint"
             -> "Configure DNS settings for the API endpoint"

    Args:
        text: Input text with potentially incorrect acronym casing

    Returns:
        Text with corrected acronym capitalization
    """
    if not text:
        return ""

    def replace_word(match):
        word = match.group(1)
        upper_word = word.upper()
        lower_word = word.lower()

        # Check for mixed-case acronyms first (e.g., mTLS, OAuth)
        if lower_word in MIXED_CASE_ACRONYMS:
            return MIXED_CASE_ACRONYMS[lower_word]

        # Check for uppercase acronyms (e.g., DNS, HTTP, TCP)
        if upper_word in UPPERCASE_ACRONYMS:
            return upper_word

        # Return original word unchanged
        return word

    return _WORD_REGEX.sub(replace_word, text)


def to_kebab_case(s: str) -> str:
    """
    Convert a snake_case string to kebab-case.

    Example: "http_loadbalancer" -> "http-loadbalancer"

    Args:
        s: Input string in snake_case format

    Returns:
        String in kebab-case format
    """
    return s.lower().replace("_", "-")


def to_anchor_name(name: str) -> str:
    """
    Convert a name to an anchor-friendly format (kebab-case).

    Example: "http_load_balancer" -> "http-load-balancer"

    Args:
        name: Input name

    Returns:
        Anchor-friendly string in kebab-case
    """
    return name.lower().replace("_", "-")


def is_uppercase_acronym(s: str) -> bool:
    """
    Check if a string is a known uppercase acronym.

    Args:
        s: String to check

    Returns:
        True if the string is a known uppercase acronym
    """
    return s.upper() in UPPERCASE_ACRONYMS


def get_mixed_case_acronym(s: str) -> str:
    """
    Get the correct mixed-case form for known acronyms.

    Args:
        s: String to check

    Returns:
        The correct mixed-case form, or empty string if not a mixed-case acronym
    """
    return MIXED_CASE_ACRONYMS.get(s.lower(), "")


def get_compound_word_human_readable(s: str) -> str:
    """
    Get the human-readable form of a compound word.

    Args:
        s: String to check

    Returns:
        The human-readable form, or empty string if not a known compound word
    """
    return COMPOUND_WORDS_HUMAN_READABLE.get(s.lower(), "")


def get_article(s: str) -> str:
    """
    Get the appropriate article ("a" or "an") for a string.

    Args:
        s: String to check

    Returns:
        "an" if the string starts with a vowel, otherwise "a"
    """
    if not s:
        return "a"
    first_char = s[0].lower()
    if first_char in "aeiou":
        return "an"
    return "a"


# Jinja2 filter aliases for backward compatibility
def underscore_to_space(s: str) -> str:
    """Convert underscores to spaces."""
    return s.replace("_", " ") if s else ""


title_case = to_title_case


if __name__ == "__main__":
    # Test examples
    test_cases = [
        "http_loadbalancer",
        "dns_zone",
        "api_endpoint",
        "bigip_apm",
        "mtls_config",
        "aws_vpc_site",
        "origin_pool",
        "health_check",
    ]

    print("Testing to_human_readable:")
    for tc in test_cases:
        print(f"  {tc!r} -> {to_human_readable(tc)!r}")

    print("\nTesting normalize_acronyms:")
    texts = [
        "Configure dns settings for the api endpoint",
        "Enable mtls authentication",
        "Use oauth for login",
    ]
    for text in texts:
        print(f"  {text!r} -> {normalize_acronyms(text)!r}")
