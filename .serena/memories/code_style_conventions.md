# Code Style and Conventions for f5xcctl

## Go Code Style
- **Formatting**: Use `gofmt` (enforced by pre-commit hooks)
- **Indentation**: Tabs for Go files
- **Line Length**: Standard Go conventions apply
- **Naming**: Follow Go conventions (camelCase for unexported, PascalCase for exported)

## Linting (golangci-lint)
Enabled linters:
- `errcheck` - Check for unchecked errors
- `gosimple` - Simplify code suggestions
- `govet` - Report suspicious constructs
- `ineffassign` - Detect ineffective assignments
- `staticcheck` - Static analysis checks
- `unused` - Check for unused code

Disabled:
- `depguard` - Dependency guard disabled

## File Conventions
| File Type | Indent | Size |
|-----------|--------|------|
| *.go | tabs | - |
| *.yml, *.yaml | spaces | 2 |
| *.json | spaces | 2 |
| *.sh | spaces | 2 |
| Makefile | tabs | - |
| *.md | spaces | - (no trailing whitespace trim) |

## Markdown Style
- Line length: 120 chars max (code blocks and tables exempt)
- Inline HTML allowed
- Duplicate headings allowed in different sections
- Bare URLs allowed

## Project Structure Patterns
- **cmd/**: CLI commands (Cobra-based)
- **pkg/**: Reusable packages
  - `config/`: Configuration management
  - `client/`: HTTP client and auth
  - `output/`: Output formatting
  - `types/`: Resource type definitions
- **docs/**: MkDocs documentation
- **tests/**: Integration tests
- **scripts/**: Utility scripts

## Pre-commit Hooks
- Trailing whitespace trimming
- Unix line endings (LF)
- YAML/JSON validation
- No commits to main/master
- Go formatting and mod tidy
- golangci-lint
- Shell script linting (shellcheck)
- Markdown linting (auto-fix enabled)
