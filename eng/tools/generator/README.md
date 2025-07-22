# Generator

This is a command line tool for generating new releases and managing automation for `github.com/Azure/azure-sdk-for-go`.

## Overview

The generator tool provides several commands to support the Azure SDK for Go development lifecycle:

- **Issue Management**: Parse GitHub release request issues into configuration
- **Release Generation**: Generate new SDK releases from TypeSpec or Swagger specifications  
- **Automation**: Process batch SDK generation for CI/CD pipelines
- **Refresh**: Regenerate all existing SDK packages
- **Templates**: Scaffold the package for onboard services

## Commands

This CLI tool provides the following commands:

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
generator automation-v2 <generate input filepath> <generate output filepath> [goVersion]
```

**Arguments:**
- `generate input filepath`: Path to the generation input JSON file
- `generate output filepath`: Path where generation output JSON will be written
- `goVersion`: Go version to use (default: "1.18")

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
- `--go-version`: Go version (default: "1.18")
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
- `--go-version`: Go version (default: "1.18")
- `--rps`: Comma-separated list of RPs to refresh (default: all)
- `--update-spec-version`: Whether to update commit ID (default: true)

### The `template` command

The `template` command generates package templates and scaffolding for new SDK packages.

## TypeSpec Support

The generator supports both traditional Swagger/OpenAPI specifications and modern TypeSpec definitions:

### TypeSpec Generation
- Uses `@azure-tools/typespec-client-generator-cli` (tsp-client) v0.21.0
- Requires `tspconfig.yaml` with `@azure-tools/typespec-go` emitter configured
- Supports both data-plane and management-plane TypeSpec projects

### Swagger Generation  
- Uses AutoRest with Go extensions
- Processes `readme.go.md` configuration files
- Supports management-plane SDK generation

## Prerequisites

For full functionality, ensure you have:

1. **Go 1.23 or later**
2. **Node.js 20 or later** 
3. **Generator tool**: Install with:
   ```bash
   go install github.com/Azure/azure-sdk-for-go/eng/tools/generator@latest
   ```
4. **tsp-client**: Install with:
   ```bash
   npm install -g @azure-tools/typespec-client-generator-cli@v0.21.0
   ```
