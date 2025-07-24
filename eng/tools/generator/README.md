# Generator

This is a command line tool for generating new releases and managing automation for `github.com/Azure/azure-sdk-for-go`.

## Overview

The generator tool provides several commands to support the Azure SDK for Go development lifecycle:

- **Environment**: Check and validate development environment prerequisites
- **Generate**: Generate individual SDK packages from TypeSpec specifications
- **Issue Management**: Parse GitHub release request issues into configuration
- **Release Generation**: Generate new SDK releases from TypeSpec or Swagger specifications  
- **Automation**: Process batch SDK generation for CI/CD pipelines
- **Refresh**: Regenerate all existing SDK packages
- **Templates**: Scaffold the package for onboard services

## Commands

This CLI tool provides the following commands:

### The `environment` command

The `environment` command checks and validates environment prerequisites for Azure Go SDK generation. It verifies the installation and versions of required tools and can automatically install missing TypeSpec tools.

**Usage:**
```bash
generator environment [flags]
```

**Flags:**
- `--auto-install`: Automatically install missing TypeSpec tools (default: true)
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**What it checks:**
- **Go**: Minimum version 1.23
- **Node.js**: Minimum version 20.0.0
- **TypeSpec compiler**: `@typespec/compiler` package
- **TypeSpec client generator CLI**: `@azure-tools/typespec-client-generator-cli` package
- **GitHub CLI**: Installation and authentication status

**Examples:**
```bash
# Check environment with auto-install (default)
generator environment

# Check environment without auto-installing missing tools
generator environment --auto-install=false

# Get results in JSON format
generator environment --output json

# Get help
generator environment --help
```

**Sample Output:**
```
All environment checks are satisfied! ✓

✓ Go: Go version 1.24 is installed ✓
✓ Node.js: Node.js version 22.17.1 is installed ✓
✓ TypeSpec Compiler: TypeSpec compiler is installed ✓
✓ TypeSpec Client Generator CLI: TypeSpec client generator CLI is installed ✓
✓ GitHub CLI: GitHub CLI 2.40.1 is installed ✓
✓ GitHub CLI Authentication: GitHub CLI is authenticated ✓

✓ Automatically installed: TypeSpec compiler, TypeSpec client generator CLI
```

### The `generate` command

The `generate` command generates Azure Go SDK packages from TypeSpec specifications. It can work with either a direct path to a TypeSpec configuration file or a GitHub PR link.

**Usage:**
```bash
generator generate <sdk-repo-path> <spec-repo-path> [flags]
```

**Arguments:**
- `sdk-repo-path`: Path to the local Azure SDK for Go repository
- `spec-repo-path`: Path to the local Azure REST API Specs repository

**Flags:**
- `--tsp-config`: Direct path to tspconfig.yaml file (relative to spec repo root)
- `--github-pr`: GitHub PR link to extract TypeSpec configuration from
- `--debug`: Enable debug output
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**Note:** You must provide exactly one of `--tsp-config` or `--github-pr`.

**Examples:**
```bash
# Generate from direct TypeSpec config path
generator generate /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --tsp-config specification/cognitiveservices/OpenAI.Inference/tspconfig.yaml

# Generate from GitHub PR
generator generate /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --github-pr https://github.com/Azure/azure-rest-api-specs/pull/12345

# Generate with debug output and JSON format
generator generate /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --tsp-config specification/service/namespace/tspconfig.yaml \
  --debug --output json
```

**What it does:**
1. Validates the provided repository paths
2. Resolves the TypeSpec configuration (checking out PR branch if needed using GitHub CLI)
3. Generates the Go SDK using the TypeSpec-Go emitter
4. Reports generation results including package info, version, and breaking changes

**Sample Output:**
```
✓ SDK generation completed successfully!

Package Name: armcognitiveservices
Package Path: /path/to/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices
Spec Folder: /path/to/azure-rest-api-specs/specification/cognitiveservices/OpenAI.Inference
Version: 1.0.0
Generation Type: mgmt
✓ Has Breaking Changes: No

Changelog:
### Features Added
- New client `armcognitiveservices.ClientFactory` which is a client factory used to create any client in this module
```

### The `issue` command

The `issue` command fetches release request issues from `github.com/Azure/sdk-release-request/issues` and parses them into configuration that other commands consume. The configuration outputs to stdout.

**Usage:**
```bash
generator issue [flags]
```

**Flags:**
- `-t, --token`: Personal access token for GitHub authentication
- `-u, --username`: GitHub username 
- `-p, --password`: GitHub password
- `--otp`: Two-factor authentication code
- `--include-data-plane`: Include data plane RP requests
- `-l, --skip-validate`: Skip validation for readme files and tags
- `--request-issues`: Specify release request IDs to parse

**Authentication:**
You need to provide authentication to query GitHub issues. Either:
1. Use a personal access token: `-t $YOUR_PERSONAL_ACCESS_TOKEN`
2. Use username/password (and OTP if needed): `-u username -p password --otp code`

**Example:**
```bash
generator issue -t $YOUR_PERSONAL_ACCESS_TOKEN > sdk-release.json
```

**Output Format:**
The command outputs a JSON configuration:
```json
{
  "track2Requests": {
    "specification/network/resource-manager/readme.md": {
      "package-2020-12-01": [
        {
          "targetDate": "2021-02-11T00:00:00Z",
          "requestLink": "https://github.com/Azure/sdk-release-request/issues/1212"
        }
      ]
    }
  },
  "typespecRequests": {}
}
```

### The `automation-v2` command

The `automation-v2` command processes batch SDK generation for automation pipelines. This command is designed to run in the root directory of azure-sdk-for-go.

**Usage:**
```bash
generator automation-v2 <generate input filepath> <generate output filepath
```

**Arguments:**
- `generate input filepath`: Path to the generation input JSON file
- `generate output filepath`: Path where generation output JSON will be written

**Input Format:**
The input file should contain:
```json
{
  "specFolder": "/path/to/azure-rest-api-specs",
  "headSha": "commit-hash",
  "repoHttpsUrl": "https://github.com/Azure/azure-rest-api-specs",
  "relatedTypeSpecProjectFolder": ["specification/service/data-plane"],
  "relatedReadmeMdFiles": ["specification/service/resource-manager/readme.md"],
  "apiVersion": "2023-01-01",
  "sdkReleaseType": "stable"
}
```

### The `release-v2` command

The `release-v2` command generates individual SDK releases for specific resource providers.

**Usage:**
```bash
generator release-v2 <azure-sdk-for-go directory> <azure-rest-api-specs directory> <rp-name/config-file> [namespaceName]
```

**Arguments:**
- `azure-sdk-for-go directory`: Path to azure-sdk-for-go repository or commit ID
- `azure-rest-api-specs directory`: Path to azure-rest-api-specs repository or commit ID  
- `rp-name/config-file`: Resource provider name or JSON config file from `issue` command
- `namespaceName`: Namespace name (default: "arm" + rp-name)

**Flags:**
- `--version-number`: Specify the version number for release
- `--package-title`: Package title for the release
- `--sdk-repo`: SDK repository URL (default: https://github.com/Azure/azure-sdk-for-go)
- `--spec-repo`: Spec repository URL (default: https://github.com/Azure/azure-rest-api-specs)
- `--spec-rp-name`: Swagger spec RP name (default: same as rp-name)
- `--release-date`: Release date for changelog
- `--skip-create-branch`: Skip creating release branch
- `--skip-generate-example`: Skip generating examples
- `--package-config`: Additional package configuration
- `-t, --token`: Personal access token for GitHub operations
- `--tsp-config`: Path to TypeSpec tspconfig.yaml
- `--tsp-option`: TypeSpec-go emit options (format: option1=value1;option2=value2)
- `--tsp-client-option`: tsp-client options (e.g., --save-inputs, --debug)

### The `refresh-v2` command

The `refresh-v2` command regenerates all existing SDK packages.

**Usage:**
```bash
generator refresh-v2 <azure-sdk-for-go directory> <azure-rest-api-specs directory>
```

**Flags:**
- `--version-number`: Specify version number for refresh
- `--sdk-repo`: SDK repository URL
- `--spec-repo`: Spec repository URL  
- `--release-date`: Release date for changelog
- `--skip-create-branch`: Skip creating release branch
- `--skip-generate-example`: Skip generating examples
- `--rps`: Comma-separated list of RPs to refresh (default: all)
- `--update-spec-version`: Whether to update commit ID (default: true)

### The `template` command

The `template` command generates package templates and scaffolding for new SDK packages.
