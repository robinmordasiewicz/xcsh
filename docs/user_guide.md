# xcsh User Guide - Complete Feature Reference

**Version**: 6.0.0
**Last Updated**: 2025-12-27
**Status**: Complete
**Platform**: TypeScript / Ink (React for CLI)

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Interactive REPL](#interactive-repl)
3. [Understanding Domain Categories](#understanding-domain-categories)
4. [Finding Domains](#finding-domains)
5. [Subscription Tier System](#subscription-tier-system)
6. [Preview Features](#preview-features)
7. [Use Cases and Workflows](#use-cases-and-workflows)
8. [Command Reference](#command-reference)
9. [Common Tasks and Examples](#common-tasks-and-examples)
10. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Installation

#### Using Homebrew

```bash
brew install robinmordasiewicz/tap/xcsh
```

#### Using Install Script

```bash
sh -c "$(curl -sSL https://raw.githubusercontent.com/robinmordasiewicz/f5xc-xcsh/main/install.sh)"
```

#### Using npm (for development)

```bash
git clone https://github.com/robinmordasiewicz/f5xc-xcsh.git
cd xcsh
npm install
npm run build
npm start
```

### First Steps

1. **Start the interactive shell**:

   ```bash
   xcsh
   ```

2. **View help**:

   ```bash
   # From within REPL:
   /help
   ```

3. **List all domains**:

   ```bash
   # From within REPL:
   /domains
   ```

4. **Navigate to a domain**:

   ```bash
   # From within REPL:
   /http_loadbalancers
   ```

5. **View your profile configuration**:

   ```bash
   # From within REPL:
   /profile show
   ```

---

## Interactive REPL

xcsh v6.0 features a modern interactive REPL (Read-Eval-Print Loop) built with Ink.

### Starting the REPL

```bash
xcsh
```

### REPL Features

- **Command History**: Use up/down arrows to navigate history
- **Tab Completion**: Press Tab to see available commands and domains
- **Slash Commands**: Quick actions starting with `/`
- **Status Bar**: Shows connection status and git information
- **Real-time Completion**: Suggestions appear as you type

### Key Bindings

| Key | Action |
|-----|--------|
| Tab | Trigger/cycle completion |
| Up/Down | Navigate history or suggestions |
| Enter | Execute command or select suggestion |
| Escape | Cancel suggestions |
| Ctrl+C (2x) | Exit REPL |
| Ctrl+D | Exit REPL immediately |

### Slash Commands

| Command | Description |
|---------|-------------|
| `/help` | Show help information |
| `/domains` | List all available domains |
| `/login` | Configure authentication |
| `/profile` | Manage connection profiles |
| `/namespace [ns]` | Switch namespace |
| `/clear` | Clear screen |
| `/exit` | Exit REPL |

---

## Understanding Domain Categories

xcsh organizes all 42 API domains into 7 categories to help you find and understand what's available.

### Categories Overview

| Category | Description | Domains | Examples |
|----------|-------------|---------|----------|
| **Security** | API security, firewalls, threat protection | 9 | api, waf, ddos, threat_campaign |
| **Platform** | User management, authentication, identity | 7 | authentication, users, tenant_and_identity, marketplace |
| **Networking** | DNS, CDN, load balancing, network config | 5 | dns, virtual, cdn, rate_limiting, network |
| **Infrastructure** | Kubernetes, cloud, service mesh, sites | 4 | kubernetes, service_mesh, cloud_infrastructure, site |
| **Operations** | Monitoring, analytics, observability, telemetry | 5 | observability, statistics, telemetry_and_insights, data_intelligence |
| **Other** | General utilities and management tools | 9 | marketplace, bigip, nginx_one, shape, blindfold, etc. |
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

List all domains using the `/domains` command:

```bash
# From within REPL:
/domains
```

Output shows all 42 domains organized by category with descriptions and tier requirements.

### Method 2: Navigate to a Domain

Navigate directly to any domain to see its resources:

```bash
# From within REPL:
/http_loadbalancers    # Navigate to HTTP load balancers
/origin_pools          # Navigate to origin pools
/app_firewalls         # Navigate to application firewalls
```

### Method 3: Use Tab Completion

The REPL provides intelligent tab completion:

```bash
# Type partial domain name and press Tab:
/http<Tab>             # Shows http_loadbalancers, etc.
/app<Tab>              # Shows app_firewalls, etc.
```

### Method 4: Explore Domain Resources

Once in a domain, list available resources:

```bash
# Navigate to domain, then:
list                   # List all resources in namespace
get [name]             # Get specific resource details
```

The REPL displays:

- Available resources in the domain
- Resource metadata and configuration
- Namespace context
- Related domains for navigation

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
xcsh

# Within REPL, view domains available at your tier:
/domains
```

### Understanding Tier Requirements

The REPL shows tier information in the status bar and domain navigation. Higher-tier domains will indicate their requirements.

Tier information appears in:

- Status bar display
- Domain navigation responses
- Completion suggestions

---

## Preview Features

Some domains are marked as "preview" or "beta" because they're still under development. Use preview features to test upcoming capabilities.

### Identifying Preview Domains

Preview domains are marked with a badge in the REPL:

```bash
# Within REPL, navigate to a preview domain:
/[preview-domain]

# The REPL displays: [PREVIEW] This is a preview feature
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

Use cases are **practical examples** of what you can do with a domain. Each domain supports specific operations:

- **Discover and catalog** - List and explore resources
- **Configure and manage** - Create, update, delete resources
- **Monitor and observe** - View status and metrics
- **Secure and protect** - Apply security policies

Each use case describes a real-world task you can accomplish within the REPL.

### What Are Workflows?

Workflows are **recommended sequences** of domains that work together to accomplish larger goals:

#### API Security Workflow

- Navigate to `/http_loadbalancers` for load balancer configuration
- Use `/app_firewalls` for security policies
- Configure `/origin_pools` for backend services

#### Network Protection Workflow

- Start with `/http_loadbalancers` for entry points
- Add `/service_policys` for traffic rules
- Configure `/rate_limiting` for protection

### Workflow Benefits

1. **See the big picture** - Understand how domains fit together
2. **Avoid manual discovery** - Pre-built solutions for common tasks
3. **Learn best practices** - Workflows follow recommended patterns
4. **Save time** - Know exactly which domains you need

### Workflow Categories

| Workflow | Domains | Purpose |
|----------|---------|---------|
| **API Security** | api, waf, threat_campaign | Protect your APIs |
| **Network Protection** | network_security, ddos, infrastructure_protection | Protect infrastructure |
| **Load Balancing** | dns, virtual, cdn | Distribute traffic |
| **Access Management** | authentication, users, tenant_and_identity | Manage access |
| **Kubernetes Mgmt** | kubernetes, service_mesh, observability | Manage Kubernetes |
| **Cloud Connectivity** | cloud_infrastructure, site, network | Connect to cloud |
| **Monitoring** | observability, statistics, telemetry_and_insights | Monitor systems |

---

## Command Reference

### REPL Commands

All commands are executed within the interactive REPL. Start with `xcsh` to enter the shell.

#### `/help`

Show general help and list available commands.

```bash
/help
```

#### `/domains`

List all available API domains with brief descriptions.

```bash
/domains
```

#### `/[domain]`

Navigate to a domain to see available actions and resources.

```bash
/http_loadbalancers
/origin_pools
/app_firewalls
```

#### `/login`

Configure authentication with the F5 Distributed Cloud API.

```bash
/login
```

#### `/profile`

Manage connection profiles for different environments.

```bash
/profile show          # Show current profile
/profile list          # List all profiles
/profile use [name]    # Switch to a profile
```

#### `/namespace [name]`

Switch the current namespace context.

```bash
/namespace shared
/namespace system
```

#### `/cloudstatus`

Check F5 Distributed Cloud service status.

```bash
/cloudstatus           # Show overall status
/cloudstatus incidents # Show current incidents
```

### Domain Operations

Once you've navigated to a domain, available operations include:

```bash
# List resources in the domain
list

# Get a specific resource by name
get [name]

# Create a new resource from JSON
create [json-payload]

# Delete a resource by name
delete [name]

# Escape back to show domains (prefix with /)
/domains
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `F5XC_API_URL` | API endpoint URL |
| `F5XC_API_TOKEN` | API authentication token |
| `F5XC_NAMESPACE` | Default namespace |
| `F5XC_DEBUG` | Enable debug output |

---

## Common Tasks and Examples

### Task 1: Set Up HTTP Load Balancer

**Goal**: Configure an HTTP load balancer with origin pools

**Domains involved**:

- `http_loadbalancers` - Frontend load balancer configuration
- `origin_pools` - Backend server pools
- `app_firewalls` - Security policies

**Steps**:

1. **Start the REPL and navigate**:

   ```bash
   xcsh
   # Within REPL:
   /http_loadbalancers
   ```

2. **List existing load balancers**:

   ```bash
   list
   ```

3. **Get details of a specific load balancer**:

   ```bash
   get my-load-balancer
   ```

4. **Navigate to origin pools**:

   ```bash
   /origin_pools
   list
   ```

5. **Configure application firewall**:

   ```bash
   /app_firewalls
   list
   ```

### Task 2: Manage Namespaces

**Goal**: Work with different namespaces

**Domains involved**:

- All domains support namespace context
- `/namespace` command for switching

**Steps**:

1. **Start the REPL**:

   ```bash
   xcsh
   ```

2. **View current namespace**:

   ```bash
   /profile show
   ```

3. **Switch namespace**:

   ```bash
   /namespace shared
   /namespace system
   ```

4. **List resources in current namespace**:

   ```bash
   /http_loadbalancers
   list
   ```

5. **Resources are scoped to namespace**:

   ```bash
   # Each domain shows resources in current namespace
   /origin_pools
   list
   ```

### Task 3: Check Cloud Status

**Goal**: Monitor F5 Distributed Cloud service status

**Steps**:

1. **Start the REPL**:

   ```bash
   xcsh
   ```

2. **Check overall status**:

   ```bash
   /cloudstatus
   ```

3. **View current incidents**:

   ```bash
   /cloudstatus incidents
   ```

4. **Monitor maintenance windows**:

   ```bash
   /cloudstatus maintenance
   ```

### Task 4: Configure Authentication

**Goal**: Set up and manage API authentication

**Steps**:

1. **Start the REPL**:

   ```bash
   xcsh
   ```

2. **Configure login**:

   ```bash
   /login
   ```

3. **View current profile**:

   ```bash
   /profile show
   ```

4. **List all profiles**:

   ```bash
   /profile list
   ```

5. **Switch to different profile**:

   ```bash
   /profile use production
   ```

---

## Troubleshooting

### Issue: "Domain not found"

**Problem**: Navigation to a domain doesn't work

**Solution**:

1. Check domain name spelling - use `/domains` to see available domains
2. Verify tier requirement - some domains need Professional/Enterprise tier
3. Use tab completion to find correct domain names

```bash
# Within REPL:
/domains            # List all available domains
/http<Tab>          # Use tab completion
```

### Issue: "Permission denied" for authentication

**Problem**: Getting authentication errors

**Solution**:

1. Verify credentials are configured correctly
2. Check API token validity
3. Ensure you're using correct API endpoint

```bash
# Within REPL:
/profile show       # Check current configuration
/login              # Reconfigure credentials
```

### Issue: "Tier requirement not met"

**Problem**: Domain requires Professional tier but you have Standard

**Solution**:

1. Check your actual subscription tier
2. Upgrade your subscription (contact F5 sales)
3. Or test with domains available at your tier

```bash
# Check environment variable
echo $F5XC_SUBSCRIPTION_TIER

# Within REPL:
/domains            # Shows available domains at your tier
```

### Issue: Slow command response

**Problem**: Commands are taking longer than expected

**Solution**:

1. Check network connectivity
2. Verify API endpoint is responsive
3. Check status bar in REPL for connection status

```bash
# Within REPL:
/cloudstatus        # Check service status

# From terminal:
ping api.volterra.us
```

### Issue: REPL not starting

**Problem**: The interactive shell fails to start

**Solution**:

1. Ensure Node.js 18+ is installed
2. Check environment variables are set
3. Try running with debug mode

```bash
# Check Node version
node --version

# Run with debug
F5XC_DEBUG=true xcsh
```

---

## Tips and Best Practices

### 1. Use Tab Completion

Tab completion is your best friend for navigating domains:

```bash
# Type partial names and press Tab:
/http<Tab>          # Shows http_loadbalancers
/app<Tab>           # Shows app_firewalls
```

### 2. Navigate Efficiently

Use slash commands for quick domain switching:

```bash
/domains            # See all domains
/http_loadbalancers # Navigate directly
/origin_pools       # Switch domains
```

### 3. Use History Navigation

The REPL remembers your command history:

- **Up/Down arrows**: Navigate previous commands
- **Tab**: Cycle through suggestions
- **Escape**: Cancel current suggestion

### 4. Monitor Connection Status

The status bar shows important information:

- Connection status to API
- Current namespace
- Git branch information

### 5. Organize by Category

Group related work by domain category:

- Security domains for API/network protection
- Platform domains for user/identity management
- Infrastructure domains for Kubernetes/cloud
- Operations domains for monitoring/analytics

### 6. Start Simple

Begin with listing resources, then explore details:

```bash
# Within REPL:
/http_loadbalancers
list                # List all resources
get my-lb           # Get specific resource
```

### 7. Use Profiles for Environments

Manage multiple environments with profiles:

```bash
# Within REPL:
/profile list       # See available profiles
/profile use prod   # Switch to production
```

---

## Getting Help and Reporting Issues

### Documentation

- **User Guide** (this document)
- **API Reference** - `docs/api_reference.md`
- **Examples** - `docs/examples/`
- **Architecture** - `claudedocs/architecture_overview.md`

### Support

- **Report Issues**: https://github.com/robinmordasiewicz/f5xc-xcsh/issues
- **Discussions**: https://github.com/robinmordasiewicz/f5xc-xcsh/discussions
- **Documentation**: https://github.com/robinmordasiewicz/f5xc-xcsh/tree/main/docs

### Contributing

- **Code contributions**: See CONTRIBUTING.md
- **Documentation improvements**: Pull requests welcome
- **Feature requests**: Use GitHub Issues

---

## Summary

**xcsh** is an interactive CLI for managing F5 Distributed Cloud resources with a modern REPL interface built on TypeScript and Ink. It features:

✅ **Interactive REPL** - Real-time command completion and navigation
✅ **42 API Domains** - Comprehensive coverage of F5 XC resources
✅ **Smart Completion** - Tab completion for domains and commands
✅ **Profile Management** - Multiple environment configurations
✅ **Status Monitoring** - Cloud status and connection information
✅ **Namespace Support** - Seamless namespace switching

**Get started today**: `xcsh`

---

*Last updated: 2025-12-27*
*Version: 6.0.0*
*Built with TypeScript and Ink (React for CLI)*
