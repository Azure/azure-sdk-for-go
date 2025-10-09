# AGENTS.md

This file provides guidance for AI agents (e.g., GitHub Copilot, MCP, or LLM-based assistants) interacting with the Azure SDK for Go repository.

## Repository Overview

### Purpose
This repository contains the Azure SDK for Go, providing Go packages for interacting with Azure services. The SDK includes:

- **Client modules** (`sdk/`) - For consuming Azure services (e.g., uploading blobs, querying databases)
- **Management modules** (`sdk/resourcemanager/`) - For configuring and managing Azure resources
- **Historical packages** (`services/`, `profiles/`) - Deprecated track 1 SDKs (no longer actively maintained)

### Key Documentation
- [Main README](README.md) - Getting started and package information
- [Contributing Guide](CONTRIBUTING.md) - Contribution guidelines and PR requirements
- [Developer Setup](documentation/developer_setup.md) - Environment setup for SDK development
- [Release Documentation](documentation/release.md) - Package release process
- [Copilot Instructions](.github/copilot-instructions.md) - Copilot-specific guidance

### Go Version Support
The SDK is compatible with the two most recent major Go releases, following Go's official [support policy](https://go.dev/doc/devel/release#policy).

## Agent Capabilities and Boundaries

### Supported Actions

AI agents can assist with the following activities:

#### Code Development
- **Reading and understanding code**: Browse SDK packages, understand APIs, and explain functionality
- **Code suggestions**: Propose improvements, bug fixes, or new features following [Azure Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
- **Testing**: Write or update unit tests using `github.com/stretchr/testify/require`
- **Examples**: Create example code in `example*_test.go` files following the [Go examples guidelines](.github/instructions/go-examples.instructions.md)

#### Documentation
- **README updates**: Improve module READMEs and documentation
- **Code comments**: Add or improve GoDoc comments following [documentation style](https://azure.github.io/azure-sdk/golang_introduction.html#documentation-style)
- **CHANGELOG updates**: Document changes in CHANGELOG.md files

#### Issue and PR Management
- **Issue triage**: Review issues, suggest labels, identify duplicates
- **PR review assistance**: Analyze PRs, suggest improvements, check for guideline compliance
- **Question answering**: Help developers with SDK usage questions

### Automation Boundaries

AI agents should **NOT** perform the following actions without human approval:

#### Build and Release
- **Triggering releases**: Only humans should use `CheckPackageReleaseReadiness` and `ReleasePackage` MCP tools
- **Modifying CI/CD pipelines**: Changes to `ci.yml`, Azure Pipelines configurations, or workflow files require careful review
- **Approving releases**: Release stage approvals in pipelines must be done by authorized personnel

#### Code Generation
- **Regenerating SDK code**: Packages under `services/` are generated from [Azure API specs](https://github.com/Azure/azure-rest-api-specs) and should not be manually modified
- **AutoRest/TypeSpec changes**: SDK generation from specifications requires specific tools and workflows (see [code generation docs](documentation/code-generation.md))

#### Security and Compliance
- **CODEOWNERS modifications**: Changes require following [CODEOWNERS validation workflow](eng/common/instructions/azsdk-tools/validate-codeowners.instructions.md)
- **Security issues**: Must be reported privately to <secure@microsoft.com>, not in public issues
- **License changes**: No modifications to licensing without explicit approval

## Key Workflows

### Development Workflow

```bash
# Navigate to the SDK module you want to work with
cd sdk/azcore  # or any other module

# Build the module
go build ./...

# Run tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.txt ./...
```

### TypeSpec/Code Generation Workflow

For modules with `tsp-location.yaml`:

```bash
# Install prerequisites
npm install -g @typespec/compiler
npm install -g @azure-tools/typespec-client-generator-cli

# Navigate to the module directory
cd sdk/<service>/<module>

# Regenerate the SDK
go generate
```

See [TypeSpec location instructions](.github/instructions/tsp-location.instructions.md) for details.

### Contributing Workflow

1. **Fork and clone** the repository
2. **Create a feature branch** from `main`
3. **Make changes** following the [Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
4. **Add tests** to ensure CI catches future regressions
5. **Update documentation** (README, CHANGELOG, examples)
6. **Run tests locally** to verify changes
7. **Submit a PR** with a descriptive title and reference to related issues
8. **Address review feedback** in additional commits

### PR Review Checklist

Agents can help verify PRs meet these requirements:

- [ ] Code follows [Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
- [ ] Tests are added/updated with proper assertions (`require.Equal`, `require.NoError`)
- [ ] CHANGELOG.md is updated with changes
- [ ] Examples are provided for new APIs
- [ ] GoDoc comments are added for public APIs
- [ ] No unrelated changes are mixed in the PR
- [ ] go.mod only has direct dependencies to Azure SDK modules or standard library
- [ ] All CI checks pass

## SDK-Specific Automation

### Module Structure

Each SDK module follows this structure:

```
sdk/<service>/<module>/
├── ci.yml                    # CI/CD pipeline configuration
├── CHANGELOG.md              # Version history and changes
├── README.md                 # Module documentation
├── go.mod                    # Go module dependencies
├── *_client.go              # Client implementation
├── models.go                # Data models
├── options.go               # Client options
├── responses.go             # Response types
├── *_test.go                # Unit tests
├── example*_test.go         # Example code
└── testdata/                # Test fixtures
```

### Test Conventions

- Use `github.com/stretchr/testify/require` for assertions
- Environment variables for live testing go in `.env` files at module root
- Look for `recording.Getenv()` or `os.Getenv()` calls to find required environment variables
- See [Go tests guidelines](.github/instructions/go-tests.instructions.md)

### Go Module Standards

- go.mod should only have direct references to:
  - Azure SDK modules (`github.com/Azure/azure-sdk-for-go/sdk/...`)
  - Standard library modules
  - `golang.org/x/...` modules
- Exception: `github.com/stretchr/testify` can be an indirect dependency
- See [Go mod standards](.github/instructions/go-mod-standards.instructions.md)

### Code Standards

- Acronyms in exported names should be uppercased (e.g., `UserID`, not `UserId`)
- All Go files should have a copyright header
- Error handling should use descriptive messages
- See [Go code guidelines](.github/instructions/go-code.instructions.md)

## Communication Channels

### Getting Help
- **Issues**: File issues via [GitHub Issues](https://github.com/Azure/azure-sdk-for-go/issues)
- **Stack Overflow**: Ask questions with tags `azure` and `go`
- **Slack**: Chat in [#Azure SDK channel](https://gophers.slack.com/messages/CA7HK8EEP) on Gophers Slack
- **Teams**: Contact the Go team via "Language - Go" channel under "Azure SDK" team

### Code Ownership
- Use [CODEOWNERS file](.github/CODEOWNERS) to find module owners
- Owners may be listed for parent folders due to wildcards
- Service owners handle issues; source owners handle PRs

## Safety and Best Practices

### For AI Agents

1. **Read before writing**: Always review existing code and documentation before suggesting changes
2. **Follow patterns**: Match the style and patterns used in the repository
3. **Test assertions**: Use the same testing patterns as existing tests
4. **Small changes**: Prefer small, focused changes over large refactorings
5. **Explain decisions**: Provide context for why changes are being suggested
6. **Reference guidelines**: Link to relevant guidelines when suggesting changes
7. **Verify compatibility**: Ensure changes don't break backward compatibility without explicit approval

### For Code Reviews

When assisting with code reviews:

1. Check compliance with [Azure Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
2. Verify tests are comprehensive and use proper assertions
3. Ensure documentation is complete and accurate
4. Check for potential breaking changes
5. Verify CHANGELOG entries describe the changes clearly
6. Look for security concerns or potential issues

### For Testing

When working with tests:

1. Run tests locally before suggesting they're complete
2. Check for required environment variables in `.env` files
3. Use `require` package for assertions, not `assert`
4. Ensure tests are repeatable and don't depend on external state
5. Add error messages to assertions to help debugging

## Additional Resources

### Azure SDK Guidelines
- [Azure Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html) - Primary reference for SDK development
- [API Design Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html) - API design principles

### Documentation
- [Developer Setup](documentation/developer_setup.md) - Machine setup for development
- [Release Guidelines](documentation/release.md) - Package release process
- [Migration Guide](documentation/MIGRATION_GUIDE.md) - Migrating from track 1 to track 2
- [Breaking Changes Guide](documentation/sdk-breaking-changes-guide.md) - Handling breaking changes

### Tools and Automation
- [Code Generation](documentation/code-generation.md) - SDK generation from specs
- [Engineering System](documentation/eng_sys.md) - Build and CI/CD systems
- [SDK Automation](documentation/sdk-automation-tsg.md) - Automation troubleshooting

## Version History

This AGENTS.md file follows the emerging [AGENTS.md standard](https://agents.md) to provide consistent AI agent guidance across Azure SDK repositories.

**Last Updated**: 2025-01-09
**Maintainers**: Azure SDK Go Team
