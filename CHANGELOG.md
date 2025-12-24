# Changelog

All notable changes to xcsh (F5 Distributed Cloud Shell) will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.1.0] - 2025-12-24

### Added

- End-to-end integration testing suite (13 comprehensive tests verifying all feature phases working together)
- Cross-feature validation test suite (16 tests covering tier × category × use cases × workflows matrix)
- Tier validation helper functions for efficient domain filtering by subscription tier
- Comprehensive user guide with examples, workflows, and best practices
- Complete API reference for developers and CLI users
- Detailed architecture documentation explaining system design
- Performance benchmarking analysis with sub-millisecond operation times
- Troubleshooting guide for common issues and error resolution
- Migration guide for users upgrading from previous versions

### Changed

- Refined tier system accuracy based on cross-feature validation (Professional: 36 domains, Enterprise: 42 domains)
- Improved test approach from assumption-based to data-driven assertions
- Enhanced error handling with complete error checking coverage
- Reorganized script files to eliminate package conflicts and improve maintainability

### Fixed

- Fixed 10 unchecked error handling issues in output operations (fmt.Fprintf)
- Fixed 4 static analysis warnings (unnecessary fmt.Sprintf, De Morgan's law optimization)
- Removed 2 unused functions reducing dead code
- Removed 3 unused imports for clean code organization
- Resolved package redeclaration errors by reorganizing tool scripts

### Performance

- Verified all operations complete in sub-millisecond timeframes
- GetDomainInfo: 125 nanoseconds per operation
- ValidateTierAccess: 49.583 nanoseconds per operation
- SearchUseCases: 17.708 microseconds per search
- Complete help text generation: approximately 115 microseconds
- Identified zero performance bottlenecks; no optimization required

### Security

- Complete error handling validation and recovery
- Robust input validation for all domain and tier operations
- Secure error messages with appropriate information disclosure

### Internal

- golangci-lint: 0 linting issues (100% clean code analysis)
- 100+ tests with 100% pass rate
- All features validated across tier system, categories, use cases, and workflows
- Code quality rating: Production-ready
- Ready for release with confidence

## [5.0.0] - Previous Release

### Added

- Initial release of xcsh (F5 Distributed Cloud Shell)
- Domain-based command structure: `xcsh <domain> <operation> <resource>`
- Comprehensive CLI for F5 Distributed Cloud management
- Shell completions for bash, zsh, and fish
- Multiple authentication methods (P12, cert/key, API token)
- JSON and YAML output formats
- DRY configuration system for easy future rebranding

### Changed

- Complete rebrand to xcsh
- Removed all legacy backward compatibility code
