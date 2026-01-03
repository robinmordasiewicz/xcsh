# F5 XC AI Assistant User Guide

The AI Assistant is an intelligent chat-based interface that helps you interact with F5 Distributed Cloud. Ask questions about your configurations, security events, troubleshooting, and platform operations using natural language.

## Quick Start

### Your First Query

Ask a simple question to get started:

```bash
xcsh ai query "What load balancers do I have?"
```

**Example Output:**

```text
Found 3 HTTP load balancers:
- lb-prod-frontend (production) - ACTIVE
- lb-staging-api (staging) - ACTIVE
- lb-dev-test (development) - INACTIVE

Follow-up Questions:
  1. Show details for lb-prod-frontend
  2. Which load balancers are inactive?
  3. Create a new load balancer
```

### Starting an Interactive Chat

For multi-turn conversations, use chat mode:

```bash
xcsh ai chat
```

```text
=== F5 XC AI Assistant Chat ===
Namespace: default
Type /help for commands, /exit to quit.

ai> What security events occurred today?
...

ai> 1
Following up: Show me the most severe events...
```

### Submitting Feedback

Help improve the AI by providing feedback:

```bash
# Positive feedback on last response
xcsh ai feedback --positive

# Negative feedback with reason
xcsh ai feedback --negative inaccurate --comment "Data was outdated"
```

---

## Command Reference

### ai query - Single Questions

Ask one-off questions without entering interactive mode.

**Syntax:**

```bash
xcsh ai query "<question>" [options]
```

**Options:**

| Flag | Short | Description |
|------|-------|-------------|
| `--namespace` | `-ns` | Target namespace (default: session namespace) |
| `--output` | `-o` | Output format: json, yaml, table, text, tsv, none |
| `--spec` | - | Output command specification as JSON |

**Aliases:** `ask`, `q`

**Examples:**

```bash
# Basic query
xcsh ai query "How many sites are active?"

# Query with JSON output
xcsh ai query "List all origin pools" -o json

# Query in specific namespace
xcsh ai query "What WAF policies exist?" -ns production

# Short alias
xcsh ai q "Show my tenant configuration"
```

---

### ai chat - Interactive Conversations

Enter a multi-turn conversation with context preservation.

**Syntax:**

```bash
xcsh ai chat [options]
```

**Options:**

| Flag | Short | Description |
|------|-------|-------------|
| `--namespace` | `-ns` | Default namespace for the session |
| `--spec` | - | Output command specification as JSON |

**Aliases:** `interactive`, `i`

**Built-in Chat Commands:**

| Command | Aliases | Description |
|---------|---------|-------------|
| `/help` | `/h` | Show available commands |
| `/clear` | `/c` | Clear conversation context |
| `/exit` | `/quit`, `/q` | Exit chat mode |
| `/feedback <type>` | - | Submit feedback on last response |
| `1`, `2`, `3`... | - | Select follow-up question by number |

**Examples:**

```bash
# Start chat in default namespace
xcsh ai chat

# Start chat in production namespace
xcsh ai chat -ns production
```

**Chat Session Example:**

```text
=== F5 XC AI Assistant Chat ===
Namespace: production
Type /help for commands, /exit to quit.

ai> What's the status of my load balancers?

Load Balancer Status Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Active: 5
Degraded: 1
Inactive: 2

Follow-up Questions:
  1. Which load balancer is degraded?
  2. Show details for inactive load balancers
  3. How do I troubleshoot degraded status?

ai> 1
Following up: Which load balancer is degraded?

lb-api-gateway is showing degraded status:
- Health check failures on 2 of 5 origin servers
- Origin pool: api-backend-pool
- Last healthy: 2026-01-02 09:45:00 UTC

Follow-up Questions:
  1. Show health check details
  2. How do I fix origin server issues?
  3. Remove unhealthy origins temporarily

ai> /feedback positive
Positive feedback submitted. Thank you!

ai> /clear
Conversation context cleared.

ai> /exit
Exiting chat mode.
```

---

### ai feedback - Improving Responses

Submit feedback to help improve AI response quality.

**Syntax:**

```bash
xcsh ai feedback <--positive | --negative <type>> [options]
```

**Options:**

| Flag | Short | Description |
|------|-------|-------------|
| `--positive` | `-p` | Submit positive feedback |
| `--negative` | `-n` | Submit negative feedback with reason type |
| `--comment` | `-c` | Add optional comment text |
| `--query-id` | `-q` | Target specific query (default: last query) |
| `--namespace` | `-ns` | Namespace context |
| `--output` | `-o` | Output format |

**Negative Feedback Types:**

| Type | When to Use |
|------|-------------|
| `other` | General feedback not fitting other categories |
| `inaccurate` | Response contained incorrect information |
| `irrelevant` | Response didn't address the question |
| `poor_format` | Response was hard to read or understand |
| `slow` | Response took too long |

**Aliases:** `fb`, `rate`

**Examples:**

```bash
# Quick positive feedback
xcsh ai feedback -p

# Negative with type
xcsh ai feedback --negative inaccurate

# Detailed negative feedback
xcsh ai feedback -n irrelevant -c "I asked about WAF but got DNS info"

# Feedback for specific query
xcsh ai feedback -p --query-id "qid-abc123"
```

---

### ai eval - RBAC Testing Mode

Test AI queries under different access permissions for RBAC validation.

**Subcommands:**

#### eval query

```bash
xcsh ai eval query "<question>" [options]
```

Sends queries through the eval endpoint with `[EVAL MODE]` context. Used for testing what responses would be returned for different permission levels.

#### eval feedback

```bash
xcsh ai eval feedback <--positive | --negative <type>> [options]
```

Submit feedback for eval queries on a separate analytics track.

**Examples:**

```bash
# Test query as different user context
xcsh ai eval query "List all namespaces"

# Feedback on eval response
xcsh ai eval feedback --positive
```

---

## Response Types Explained

The AI assistant returns different response types based on your question:

### Generic Text Responses

Most common response type with text explanations and optional links.

```bash
xcsh ai query "What is F5 Distributed Cloud?"
```

```text
F5 Distributed Cloud (F5 XC) is a SaaS-based platform that provides
application security, multi-cloud networking, and edge computing services.

Links:
  - F5 XC Documentation: https://docs.cloud.f5.com
  - Getting Started: https://docs.cloud.f5.com/getting-started
```

### Security Event Explanations

Detailed analysis of WAF, bot, and security events.

```bash
xcsh ai query "Why was request xyz blocked?"
```

```text
Security Event Analysis
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Summary: WAF blocked request due to SQL injection attempt

Request Details:
  Method: GET
  Path: /api/users?id=1 OR 1=1
  Source IP: 192.168.1.100
  Time: 2026-01-02 10:30:00 UTC

Violations:
  - SQL_INJECTION: UNION SELECT pattern detected in query parameter
  - ANOMALY_SCORE: Score 15 exceeded threshold of 5

Action: BLOCK
Confidence: HIGH

Threat Campaigns:
  - SQLi-Campaign-2024-001
```

### List Responses

Resource listings with counts and summaries.

```bash
xcsh ai query "List my origin pools"
```

```text
Found 4 origin pools:

  NAME              NAMESPACE    SERVERS    STATUS
  api-backend       production   3          HEALTHY
  web-frontend      production   5          HEALTHY
  staging-pool      staging      2          DEGRADED
  dev-pool          development  1          INACTIVE

Total: 4 pools (2 healthy, 1 degraded, 1 inactive)
```

### Site Analysis Reports

Health analysis with metrics and recommendations.

```bash
xcsh ai query "Analyze site-1 health"
```

```text
Site Analysis: site-1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Health Status: DEGRADED

Metrics:
  Availability:     98.5%
  Latency (P99):    245ms
  Error Rate:       1.2%
  Throughput:       1.5K rps

Recommendations:
  1. Consider scaling up to handle current traffic levels
  2. Investigate high latency - check origin server response times
  3. Review error logs for the 1.2% error rate
```

### Dashboard Filters

Filter expressions for use in dashboards.

```bash
xcsh ai query "Create a filter for blocked requests from today"
```

```text
Dashboard Filter Generated:

  Expression: waf_action:BLOCK AND timestamp:[now-24h TO now]
  Context: Security Events Dashboard

You can apply this filter in the F5 XC Console under:
  Dashboards > Security > Events
```

### Widget Data Tables

Tabular data for dashboard visualization.

```bash
xcsh ai query "Show traffic summary by load balancer"
```

```text
Traffic Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  LOAD BALANCER     REQUESTS    AVAILABILITY
  lb-prod           1.2M        99.9%
  lb-staging        45K         98.5%
  lb-dev            2.3K        100.0%
```

---

## Working with Context

### Multi-Turn Conversations

In chat mode, the AI maintains context between messages:

```text
ai> Show me production load balancers
[Lists 3 production LBs]

ai> Which one has the most traffic?
[Knows you're asking about the 3 production LBs just listed]

ai> Show its origin pool configuration
[Knows "its" refers to the highest-traffic LB]
```

### Follow-Up Suggestions

After each response, the AI suggests relevant follow-up questions:

```text
Follow-up Questions:
  1. Show details for this load balancer
  2. View health check configuration
  3. How do I update the origin pool?
```

Select by typing the number:

```text
ai> 2
Following up: View health check configuration...
```

### Clearing Context

If the conversation becomes confusing, clear and start fresh:

```text
ai> /clear
Conversation context cleared.
```

### Namespace Scoping

Queries are scoped to a namespace for relevant results:

```bash
# Set namespace at start
xcsh ai chat -ns production

# Or use flag per query
xcsh ai query "List resources" -ns staging
```

---

## Best Practices

### Crafting Effective Queries

**Be Specific:**

```bash
# Good - specific and actionable
xcsh ai query "Why is lb-prod-api returning 502 errors since 10am?"

# Less effective - too vague
xcsh ai query "Why errors?"
```

**Provide Context:**

```bash
# Good - includes relevant details
xcsh ai query "Explain WAF event with request ID req-abc123 from today"

# Less effective - missing specifics
xcsh ai query "Explain security event"
```

**Ask One Thing at a Time:**

```bash
# Good - single focused question
xcsh ai query "How do I add a new origin server to my-pool?"

# Less effective - multiple questions
xcsh ai query "How do I add origins and configure health checks and set up routing?"
```

### When to Use Chat vs Query

**Use `ai query` when:**

- Asking a single question
- Scripting or automation
- Piping output to other commands
- Quick lookups

**Use `ai chat` when:**

- Exploring a topic
- Troubleshooting (needs back-and-forth)
- Learning about features
- Working through a complex task

### Providing Useful Feedback

Good feedback helps improve the AI:

```bash
# Specific negative feedback
xcsh ai feedback -n inaccurate -c "The API endpoint mentioned doesn't exist in v2"

# Always provide feedback for helpful responses
xcsh ai feedback -p -c "Clear step-by-step instructions"
```

---

## Troubleshooting

### "Not connected to API"

**Cause:** No active connection to F5 XC API.

**Solution:**

```bash
# Set environment variables
export F5XC_API_URL="https://your-tenant.console.ves.volterra.io"
export F5XC_API_TOKEN="your-api-token"

# Or start xcsh and configure
xcsh
> login
```

### "Not authenticated"

**Cause:** API token is invalid or expired.

**Solution:**

1. Generate a new API token in F5 XC Console
2. Update your credentials:

   ```bash
   export F5XC_API_TOKEN="new-token"
   ```

### "Chat mode requires an interactive terminal"

**Cause:** Running chat mode in a non-interactive environment (pipe, script, etc.)

**Solution:** Use `ai query` instead for non-interactive use:

```bash
# In scripts, use query mode
xcsh ai query "Your question" -o json
```

### Empty or Unexpected Responses

**Possible Causes:**

1. **Wrong namespace:** Check you're querying the right namespace
2. **No matching resources:** The query might be correct but there's nothing to show
3. **Permission issues:** Your token might lack access to certain resources

**Solution:**

```bash
# Verify namespace
xcsh ai query "What is my current namespace?"

# Try with explicit namespace
xcsh ai query "List resources" -ns system

# Check for errors in JSON output
xcsh ai query "List resources" -o json
```

### No Follow-Up Suggestions

**Cause:** Some queries don't generate follow-ups if they're self-contained.

**This is normal** - not all responses will have follow-up suggestions.

---

## Output Format Examples

### JSON Output

```bash
xcsh ai query "What namespace am I in?" -o json
```

```json
{
  "query_id": "qid-abc123",
  "current_query": "What namespace am I in?",
  "generic_response": {
    "text": "You are currently in the 'default' namespace."
  },
  "follow_up_queries": [
    "List all namespaces",
    "Switch to a different namespace"
  ]
}
```

### YAML Output

```bash
xcsh ai query "What namespace am I in?" -o yaml
```

```yaml
query_id: qid-abc123
current_query: What namespace am I in?
generic_response:
  text: You are currently in the 'default' namespace.
follow_up_queries:
  - List all namespaces
  - Switch to a different namespace
```

### TSV Output (for parsing)

```bash
xcsh ai query "List load balancers" -o tsv
```

```text
NAME NAMESPACE STATUS
lb-prod production ACTIVE
lb-staging staging ACTIVE
```

---

## Command Aliases Summary

| Full Command | Aliases |
|--------------|---------|
| `xcsh ai_services` | `ai`, `genai`, `assistant` |
| `ai query` | `q`, `ask` |
| `ai chat` | `i`, `interactive` |
| `ai feedback` | `fb`, `rate` |
| `ai eval query` | `ai eval q` |
| `ai eval feedback` | `ai eval fb` |

---

## See Also

- [AI Services Command Reference](../commands/ai_services/index.md)
- [Query Command Details](../commands/ai_services/query/index.md)
- [Chat Command Details](../commands/ai_services/chat/index.md)
- [Feedback Command Details](../commands/ai_services/feedback/index.md)
