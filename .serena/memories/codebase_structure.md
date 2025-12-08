# vesctl Codebase Structure

## Root Directory
```
vesctl/
├── main.go              # Entry point - calls cmd.Execute()
├── go.mod               # Go module definition (github.com/robinmordasiewicz/vesctl)
├── go.sum               # Dependency checksums
├── Makefile             # Build, test, lint commands
├── .goreleaser.yaml     # GoReleaser configuration
├── .golangci.yml        # Linter configuration
├── .pre-commit-config.yaml # Pre-commit hooks
├── .editorconfig        # Editor settings
├── .markdownlint.yaml   # Markdown linting rules
├── mkdocs.yml           # Documentation config
├── install.sh           # Installation script
└── README.md            # Project documentation
```

## cmd/ - CLI Commands (Cobra)
```
cmd/
├── root.go              # Root command and global flags
├── version.go           # Version command
├── completion.go        # Shell completion command
├── configure.go         # Configuration command
├── login.go             # Login command
├── resource.go          # Generic resource operations
├── request.go           # API request handling
├── request_rpc.go       # RPC-style requests
├── request_rpc_generated.go # Auto-generated RPC handlers
├── request_cmdseq.go    # Command sequencing
├── request_secrets.go   # Secret handling
├── configuration.go     # Configuration management
├── api_endpoint.go      # API endpoint command
├── site.go              # Site management
├── site_aws_vpc.go      # AWS VPC site operations
├── site_azure_vnet.go   # Azure VNet site operations
├── configuration/       # Configuration subcommands
├── gen_rpc/             # RPC generation tooling
├── site/                # Site subcommands
├── request/             # Request subcommands
└── api_endpoint/        # API endpoint subcommands
```

## pkg/ - Reusable Packages
```
pkg/
├── config/
│   ├── config.go        # Configuration loading and management
│   └── config_test.go   # Configuration tests
├── client/
│   ├── client.go        # HTTP client for F5 XC API
│   ├── client_test.go   # Client tests
│   └── auth.go          # Authentication (P12, cert/key)
├── output/
│   ├── formatter.go     # Output formatting (JSON, YAML, table, TSV)
│   └── formatter_test.go # Formatter tests
├── types/
│   ├── resource.go      # Resource type definitions
│   ├── resources_generated.go # Auto-generated resource types
│   └── registry.go      # Resource registry
├── terraform/           # Terraform-related utilities
└── openapi/             # OpenAPI schema handling
```

## docs/ - MkDocs Documentation
```
docs/
├── index.md             # Documentation home
├── getting-started/     # Installation, authentication, quickstart
├── commands/            # Command reference
├── examples/            # Usage examples
├── reference/           # Global flags, environment variables
├── specifications/api/  # API specifications (JSON)
├── overrides/           # MkDocs theme overrides
└── stylesheets/         # Custom CSS
```

## Other Directories
```
├── tests/               # Integration tests
├── scripts/             # Utility scripts
├── dist/                # Build output directory
├── .github/             # GitHub Actions workflows
└── claudedocs/          # Claude-specific documentation
```

## Key Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `gopkg.in/yaml.v3` - YAML parsing
- `software.sslmate.com/src/go-pkcs12` - P12 certificate handling
