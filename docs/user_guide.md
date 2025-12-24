# xcsh User Guide - Complete Feature Reference

**Version**: 1.0
**Last Updated**: 2025-12-24
**Status**: Complete

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Understanding Domain Categories](#understanding-domain-categories)
3. [Finding Domains](#finding-domains)
4. [Subscription Tier System](#subscription-tier-system)
5. [Preview Features](#preview-features)
6. [Use Cases and Workflows](#use-cases-and-workflows)
7. [Command Reference](#command-reference)
8. [Common Tasks and Examples](#common-tasks-and-examples)
9. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Installation

#### Using Homebrew

```bash
brew install robinmordasiewicz/tap/xcsh
```

#### Using Install Script

```bash
sh -c "$(curl -sSL https://raw.githubusercontent.com/robinmordasiewicz/xcsh/main/install.sh)"
```

#### Manual Build

```bash
git clone https://github.com/robinmordasiewicz/xcsh.git
cd xcsh
go build -o xcsh
./xcsh version
```

### First Steps

1. **View help**:

   ```bash
   xcsh help
   ```

2. **List all domains**:

   ```bash
   xcsh domains list
   ```

3. **Get help for a specific domain**:

   ```bash
   xcsh api --help
   ```

4. **View your configuration**:

   ```bash
   xcsh configure show
   ```

---

## Understanding Domain Categories

xcsh organizes all 42 API domains into 7 categories to help you find and understand what's available.

### Categories Overview

| Category | Description | Domains | Examples |
|----------|-------------|---------|----------|
| **Security** | API security, firewalls, threat protection | 9 | api, application_firewall, ddos, threat_campaign |
| **Platform** | User management, authentication, identity | 7 | authentication, users, tenant_and_identity, marketplace |
| **Networking** | DNS, CDN, load balancing, network config | 5 | dns, virtual, cdn, rate_limiting, network |
| **Infrastructure** | Kubernetes, cloud, service mesh, sites | 4 | kubernetes, service_mesh, cloud_infrastructure, site |
| **Operations** | Monitoring, analytics, observability, telemetry | 5 | observability, statistics, telemetry_and_insights, data_intelligence |
| **Other** | General utilities and management tools | 11 | marketplace, bigip, nginx_one, shape, blindfold, etc. |
| **AI** | AI and machine learning features | 1 | generative_ai |

### Why Categories Matter

Categories help you:

1. **Discover related features** - Domains in the same category often work together
2. **Understand scope** - Quickly grasp what a domain controls
3. **Plan workflows** - Build multi-domain solutions logically
4. **Learn progressively** - Start with one category, expand as needed

---

## Finding Domains

### Method 1: List by Category

List all domains in a specific category:

```bash
xcsh domains list --category Security
xcsh domains list --category Platform
xcsh domains list --category Infrastructure
```

Output shows all domains with their descriptions and tier requirements.

### Method 2: Search by Use Case

Find domains that help with specific tasks:

```bash
xcsh api --help
# Shows use cases: "Discover and catalog APIs", "Test API security", etc.

xcsh dns --help
# Shows use cases: "Configure DNS load balancing", "Manage DNS zones", etc.
```

### Method 3: Explore Related Domains

See which domains work well together:

```bash
xcsh api --help
# Under "RELATED DOMAINS:" shows:
# - application_firewall (same security category)
# - authentication (authentication is related to API security)
# - network_security (complementary security functionality)
# etc.
```

### Method 4: Direct Domain Help

Get complete information about a domain:

```bash
xcsh kubernetes --help
```

This displays:

- Domain description
- Category and complexity level
- Use cases (practical examples)
- Related domains (work well with this domain)
- Suggested workflows (common task sequences)
- Available operations

---

## Subscription Tier System

xcsh respects your F5 Distributed Cloud subscription tier. Different tiers unlock different domains and features.

### Tier Levels

**Standard Tier** (Base subscription)

- Access to 25 core domains
- Essential API management and security features
- No advanced infrastructure features

#### Professional Tier

- Access to all 42 domains (full feature set)
- Includes infrastructure management
- Enables advanced workflows
- Recommended for most teams

#### Enterprise Tier

- All Professional features
- Advanced analytics and integrations
- Premium support features

### Checking Your Tier

Your tier is determined by your F5 account. Set it via environment variable:

```bash
# See your current tier
echo $F5XC_SUBSCRIPTION_TIER

# Set your tier (for testing with Standard-only features)
export F5XC_SUBSCRIPTION_TIER=Standard
xcsh domains list

# Check which domains you have access to
xcsh domains list --tier Standard
```

### Understanding Tier Requirements

Every domain has a tier requirement shown in help text:

```bash
$ xcsh kubernetes --help
Category: Infrastructure
Tier: Professional
# This domain requires Professional or Enterprise tier
```

Domains show their tier in:

- Domain help text
- Domain list output
- Completion suggestions

---

## Preview Features

Some domains are marked as "preview" or "beta" because they're still under development. Use preview features to test upcoming capabilities.

### Identifying Preview Domains

Preview domains are marked with a badge in help text:

```bash
$ xcsh [preview-domain] --help
[PREVIEW] This is a preview feature

# Additional information about preview status
```

### Using Preview Features

1. **Preview features are experimental** - API may change
2. **Not recommended for production** - Features are still being tested
3. **Feedback welcome** - Help us improve by reporting issues
4. **May be limited** - Some functionality may not be complete

### Common Preview Patterns

Preview domains work like standard domains:

- Use normal commands and syntax
- Include use cases and related domains
- Subject to tier requirements
- Appear in help and completion

---

## Use Cases and Workflows

### What Are Use Cases?

Use cases are **practical examples** of what you can do with a domain:

```bash
$ xcsh api --help

USE CASES:
  • Discover and catalog APIs
  • Test API security and behavior
  • Manage API credentials
  • Define API groups and testing policies
```

Each use case describes a real-world task you can accomplish.

### What Are Workflows?

Workflows are **recommended sequences** of domains that work together to accomplish larger goals:

```bash
$ xcsh api --help

SUGGESTED WORKFLOWS:
  • API Security Workflow
    Secure APIs with firewall and threat detection
    Involves: api, application_firewall, threat_campaign

  • Network Protection Workflow
    Protect network infrastructure and applications
    Involves: network_security, ddos, infrastructure_protection
```

### Workflow Benefits

1. **See the big picture** - Understand how domains fit together
2. **Avoid manual discovery** - Pre-built solutions for common tasks
3. **Learn best practices** - Workflows follow recommended patterns
4. **Save time** - Know exactly which domains you need

### Workflow Categories

| Workflow | Domains | Purpose |
|----------|---------|---------|
| **API Security** | api, application_firewall, threat_campaign | Protect your APIs |
| **Network Protection** | network_security, ddos, infrastructure_protection | Protect infrastructure |
| **Load Balancing** | dns, virtual, cdn | Distribute traffic |
| **Access Management** | authentication, users, tenant_and_identity | Manage access |
| **Kubernetes Mgmt** | kubernetes, service_mesh, observability | Manage Kubernetes |
| **Cloud Connectivity** | cloud_infrastructure, site, network | Connect to cloud |
| **Monitoring** | observability, statistics, telemetry_and_insights | Monitor systems |

---

## Command Reference

### Core Commands

#### `xcsh help`

Show general help and list available domains.

```bash
xcsh help
```

#### `xcsh [domain] --help`

Get detailed help for a specific domain.

```bash
xcsh api --help
xcsh kubernetes --help
xcsh dns --help
```

**Help text includes**:

- Domain description
- Tier requirement
- Category and complexity
- Use cases (practical examples)
- Related domains (work well together)
- Suggested workflows (common task sequences)
- Available operations

#### `xcsh domains list`

List all 42 domains with brief descriptions.

```bash
# List all domains
xcsh domains list

# Filter by category
xcsh domains list --category Security

# Filter by tier
xcsh domains list --tier Professional

# Combine filters
xcsh domains list --category Infrastructure --tier Professional
```

#### `xcsh domains completion`

Generate shell completion scripts.

```bash
# For Bash
xcsh completion bash | sudo tee /usr/local/etc/bash_completion.d/xcsh

# For Zsh
xcsh completion zsh | sudo tee /usr/share/zsh/site-functions/_xcsh

# Then reload shell
source ~/.bashrc  # or ~/.zshrc
```

#### `xcsh configure show`

View current configuration.

```bash
xcsh configure show
```

Shows:

- API endpoint
- Authentication status
- Current subscription tier
- Namespace (if set)

### Domain Operations

Once you've selected a domain, available operations include:

```bash
xcsh [domain] list              # List resources
xcsh [domain] get [name]        # Get specific resource
xcsh [domain] create [file]     # Create from YAML/JSON
xcsh [domain] delete [name]     # Delete resource
xcsh [domain] apply [file]      # Create or update
xcsh [domain] replace [name] [file]   # Replace entirely
xcsh [domain] patch [name]      # Partial update
xcsh [domain] status [name]     # Check status
xcsh [domain] add-labels [name] [labels]  # Add labels
```

---

## Common Tasks and Examples

### Task 1: Set Up API Security

**Goal**: Secure your APIs with firewall and threat detection

**Domains involved**:

- `api` - Define and manage your APIs
- `application_firewall` - Add security rules
- `threat_campaign` - Monitor threats

**Steps**:

1. **View the recommended workflow**:

   ```bash
   xcsh api --help | grep -A 5 "API Security Workflow"
   ```

2. **Explore each domain**:

   ```bash
   xcsh api --help          # Understand API management
   xcsh application_firewall --help  # Learn firewall options
   xcsh threat_campaign --help       # Understand threat detection
   ```

3. **Start with API management**:

   ```bash
   xcsh api list            # See existing APIs
   xcsh api create api.yaml # Add your API definition
   ```

4. **Add firewall rules**:

   ```bash
   xcsh application_firewall create rules.yaml
   ```

5. **Enable threat detection**:

   ```bash
   xcsh threat_campaign create campaign.yaml
   ```

### Task 2: Deploy to Kubernetes

**Goal**: Manage Kubernetes clusters and services

**Domains involved**:

- `kubernetes` - Core Kubernetes management
- `service_mesh` - Advanced traffic control
- `observability` - Monitor performance

**Steps**:

1. **Check prerequisites**:

   ```bash
   # Verify you have Professional tier (Kubernetes requires it)
   echo $F5XC_SUBSCRIPTION_TIER
   ```

2. **View the workflow**:

   ```bash
   xcsh kubernetes --help | grep -A 5 "Kubernetes Management"
   ```

3. **List existing clusters**:

   ```bash
   xcsh kubernetes list
   ```

4. **Deploy your cluster config**:

   ```bash
   xcsh kubernetes create cluster.yaml
   ```

5. **Set up service mesh (advanced)**:

   ```bash
   xcsh service_mesh create service-mesh.yaml
   ```

6. **Monitor health**:

   ```bash
   xcsh observability list
   xcsh observability get [name]
   ```

### Task 3: Configure DNS Load Balancing

**Goal**: Distribute traffic across regions

**Domains involved**:

- `dns` - DNS configuration
- `virtual` - Virtual hosts
- `cdn` - Content delivery

**Steps**:

1. **View load balancing workflow**:

   ```bash
   xcsh dns --help | grep -A 5 "Load Balancing"
   ```

2. **Check available DNS configurations**:

   ```bash
   xcsh dns list
   ```

3. **Create DNS zone**:

   ```bash
   xcsh dns create zone.yaml
   ```

4. **Set up virtual hosts**:

   ```bash
   xcsh virtual create virtual-host.yaml
   ```

5. **Enable CDN caching**:

   ```bash
   xcsh cdn create cdn-config.yaml
   ```

### Task 4: Manage User Access

**Goal**: Control who can access your resources

**Domains involved**:

- `authentication` - Login and auth setup
- `users` - User management
- `tenant_and_identity` - Identity governance

**Steps**:

1. **View access management workflow**:

   ```bash
   xcsh authentication --help | grep -A 5 "Access Management"
   ```

2. **Check authentication methods**:

   ```bash
   xcsh authentication list
   ```

3. **Configure authentication**:

   ```bash
   xcsh authentication create auth.yaml
   ```

4. **Create users**:

   ```bash
   xcsh users create users.yaml
   ```

5. **Set up tenant identity**:

   ```bash
   xcsh tenant_and_identity apply identity.yaml
   ```

---

## Troubleshooting

### Issue: "Domain not found"

**Problem**: When running `xcsh [domain] --help`, you get "domain not found"

**Solution**:

1. Check domain name spelling - use `xcsh domains list` to see exact names
2. Verify tier requirement - some domains need Professional/Enterprise tier
3. Check if domain is preview - preview domains may need additional setup

```bash
# Find the right domain
xcsh domains list | grep -i keyword

# Check tier requirement
xcsh domains list | grep [domain-name]
```

### Issue: "Permission denied" for authentication

**Problem**: Getting authentication errors

**Solution**:

1. Verify credentials are configured correctly
2. Check API token validity
3. Ensure you're using correct API endpoint

```bash
# Check configuration
xcsh configure show

# Verify authentication works
xcsh api list
```

### Issue: "Tier requirement not met"

**Problem**: Domain requires Professional tier but you have Standard

**Solution**:

1. Check your actual subscription tier
2. Upgrade your subscription (contact F5 sales)
3. Or test with domains available at your tier

```bash
# See your tier
echo $F5XC_SUBSCRIPTION_TIER

# List only domains at your tier
xcsh domains list --tier Standard
```

### Issue: Slow command response

**Problem**: Commands are taking longer than expected

**Solution**:

1. Check network connectivity
2. Verify API endpoint is responsive
3. Try with a single domain instead of listing all

```bash
# Test API connectivity
xcsh api list

# If slow, check network:
ping api.volterra.us
```

### Issue: Completion not working in terminal

**Problem**: Tab completion isn't working for domain names

**Solution**:

1. Ensure completion script is installed
2. Reload shell configuration

```bash
# For Bash
source ~/.bashrc

# For Zsh
source ~/.zshrc

# Or restart terminal
```

---

## Tips and Best Practices

### 1. Use Help Frequently

Help text is comprehensive and always current:

```bash
xcsh [domain] --help
```

### 2. Explore Related Domains

Every domain shows 5 related domains that work well together:

```bash
xcsh api --help
# Look at "RELATED DOMAINS:" section
```

### 3. Follow Suggested Workflows

Workflows are pre-planned sequences that follow best practices:

```bash
xcsh [domain] --help
# Look at "SUGGESTED WORKFLOWS:" section
```

### 4. Check Use Cases

Use cases show practical examples of what you can accomplish:

```bash
xcsh [domain] --help
# Look at "USE CASES:" section
```

### 5. Organize by Category

Group related work by domain category:

- Security domains for API/network protection
- Platform domains for user/identity management
- Infrastructure domains for Kubernetes/cloud
- Operations domains for monitoring/analytics

### 6. Start Simple

Begin with one domain, then expand to related domains:

```bash
# Start with basic API management
xcsh api list

# Then explore related domains
xcsh application_firewall --help
```

### 7. Use Tier Appropriately

Know what your subscription includes and use tier-appropriate domains:

```bash
# Standard tier - 25 domains
# Professional tier - all 42 domains
```

---

## Getting Help and Reporting Issues

### Documentation

- **User Guide** (this document)
- **API Reference** - `docs/api_reference.md`
- **Examples** - `docs/examples/`
- **Architecture** - `claudedocs/architecture_overview.md`

### Support

- **Report Issues**: https://github.com/robinmordasiewicz/xcsh/issues
- **Discussions**: https://github.com/robinmordasiewicz/xcsh/discussions
- **Documentation**: https://github.com/robinmordasiewicz/xcsh/tree/main/docs

### Contributing

- **Code contributions**: See CONTRIBUTING.md
- **Documentation improvements**: Pull requests welcome
- **Feature requests**: Use GitHub Issues

---

## Summary

**xcsh** is a comprehensive CLI for managing F5 Distributed Cloud resources across 42 domains organized into 7 categories. It features:

✅ **Smart Discovery** - Find domains by category, use case, or workflow
✅ **Tier-Based Access** - Respects your subscription level
✅ **Practical Guidance** - Use cases and workflows for common tasks
✅ **Related Domain Discovery** - See what works together
✅ **Complete Help** - Everything in one place

**Get started today**: `xcsh help`

---

*Last updated: 2025-12-24*
*Version: 1.0*
*Part of the xcsh data-driven CLI architecture*
