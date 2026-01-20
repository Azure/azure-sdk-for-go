# Generator

This is a command line tool for generating new releases and managing automation for `github.com/Azure/azure-sdk-for-go`.

## Overview

The generator tool provides several commands to support the Azure SDK for Go development lifecycle:

- **Build**: Compile and validate Go packages using `go build` and `go vet`
- **Changelog**: Generate and update changelog content for SDK packages
- **Version**: Calculate and update version numbers across package files
- **Environment**: Check and validate development environment prerequisites
- **Generate**: Generate individual SDK packages from TypeSpec specifications
- **Issue Management**: Parse GitHub release request issues into configuration
- **Release Generation**: Generate new SDK releases from TypeSpec or Swagger specifications
- **Automation**: Process batch SDK generation for CI/CD pipelines
- **Refresh**: Regenerate all existing SDK packages
- **Templates**: Scaffold the package for onboard services
- **Generate Go Readme**: Generate or update readme.go.md files for Swagger specifications

## Commands

This CLI tool provides the following commands:

### Inner Loop Development Commands

The following commands support inner loop development workflows for working with SDK packages:

#### The `build` command

The `build` command compiles and validates Go packages in a specified folder using `go build` and `go vet`.

**Usage:**

```bash
generator build <folder-path>
```

**Arguments:**

- `folder-path`: Path to the folder containing Go packages to build and vet

**Flags:**

- `-v, --verbose`: Enable verbose output showing build details
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**What it does:**

1. Runs `go build` to compile the Go packages
2. Runs `go vet` to check for common Go programming errors
3. Reports any issues found during build or vet process

**Examples:**

```bash
# Build and vet packages in a specific folder
generator build /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

# Build with verbose output
generator build /path/to/package --verbose

# Build with JSON output
generator build /path/to/package --output json
```

#### The `changelog` command

The `changelog` command generates and updates changelog content for SDK packages based on code changes.

**Usage:**

```bash
generator changelog <package-path>
```

**Arguments:**

- `package-path`: Absolute path to a Go module (containing go.mod file)

**Flags:**

- `-v, --verbose`: Enable verbose output
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**What it does:**

1. Determines the package status (new package vs. existing package)
2. For new packages: generates changelog according to the template
3. For existing packages: compares current package exports with previous released version and calculates the changelog
4. Updates the CHANGELOG.md file, replacing existing version entry if it exists

**Examples:**

```bash
# Update changelog for an existing package
generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

# Update changelog with verbose output
generator changelog /path/to/package --verbose

# Generate changelog with JSON output
generator changelog /path/to/package --output json
```

#### The `version` command

The `version` command calculates and updates version numbers across all version-related files in an SDK package.

**Usage:**

```bash
generator version <package-path>
```

**Arguments:**

- `package-path`: Absolute path to a Go module (containing go.mod and version.go files)

**Flags:**

- `--sdkversion`: Specific SDK version to set (e.g., "1.2.0" or "1.2.0-beta.1")
- `--sdkreleasetype`: SDK release type ("beta" or "stable"), only used when --sdkversion is not specified
- `-v, --verbose`: Enable verbose output
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**What it does:**

1. If `--sdkversion` is specified: updates all version files with the provided version
2. If `--sdkversion` is not specified: calculates new version based on package changes and release type, then updates all version files
3. Updates version in autorest.md (if exists), version.go, go.mod, README.md, and import paths

**Examples:**

```bash
# Update version files with a specific version
generator version /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --sdkversion 1.2.0

# Calculate and update version based on changes
generator version /path/to/package

# Calculate and update version as stable release
generator version /path/to/package --sdkreleasetype stable

# Calculate and update version as beta release
generator version /path/to/package --sdkreleasetype beta

# Calculate and update version with JSON output
generator version /path/to/package --output json
```

**Sample Output:**

```
✓ Version updated successfully!

Package: /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute
Previous Version: 1.1.0
New Version: 1.2.0
```

#### The `environment` command

The `environment` command checks and validates environment prerequisites for Azure Go SDK generation. It verifies the installation and versions of required tools and can automatically install missing TypeSpec tools.

**Usage:**

```bash
generator environment [flags]
```

**Flags:**

- `--auto-install`: Automatically install missing TypeSpec tools (default: true)
- `-o, --output`: Output format, either "text" or "json" (default: "text")

**What it checks:**

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

### SDK Generation Commands

The following commands support SDK generation from TypeSpec and Swagger specifications:

#### The `generate` command

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

#### The `issue` command

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

#### The `automation-v2` command

The `automation-v2` command processes batch SDK generation for automation pipelines. This command is designed to run in the root directory of azure-sdk-for-go and handles multiple SDK generations in a single execution.

**Usage:**

```bash
generator automation-v2 <generate input filepath> <generate output filepath>
```

**Arguments:**

- `generate input filepath`: Path to the generation input JSON file
- `generate output filepath`: Path where generation output JSON will be written

**What it does:**

1. Reads the input configuration file containing specification details
2. Processes multiple TypeSpec projects and README files
3. Generates SDKs for all specified services
4. Writes comprehensive output including generation results and errors
5. Handles batch operations for CI/CD automation pipelines

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

**Output Format:**
The command generates a JSON output file with generation results:

```json
{
  "packages": [
    {
      "packageName": "armservice",
      "path": "sdk/resourcemanager/service/armservice",
      "readmeMd": "specification/service/resource-manager/readme.md",
      "changelog": "### Features Added\n- Initial release",
      "breaking": false
    }
  ],
  "errors": []
}
```

**Examples:**

```bash
# Process automation input file
generator automation-v2 ./input.json ./output.json

# Typical CI/CD usage
generator automation-v2 /tmp/generation-input.json /tmp/generation-output.json
```

#### The `release-v2` command

The `release-v2` command generates individual SDK releases for specific resource providers. It creates new SDK packages or updates existing ones with new API versions.

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

- `--package-title`: Package title for the release
- `--sdk-repo`: SDK repository URL (default: https://github.com/Azure/azure-sdk-for-go)
- `--spec-repo`: Spec repository URL (default: https://github.com/Azure/azure-rest-api-specs)
- `--spec-rp-name`: Swagger spec RP name (default: same as rp-name)
- `--release-date`: Release date for changelog
- `--skip-create-branch`: Skip creating release branch
- `--skip-generate-example`: Skip generating examples
- `--package-config`: Additional package configuration
- `-t, --token`: Personal access token for GitHub operations
- `--force-stable-version`: Force generation of stable SDK versions even when input files contain preview API versions. The tag must not contain preview when using this flag
- `--tsp-config`: Path to TypeSpec tspconfig.yaml
- `--tsp-option`: TypeSpec-go emit options (format: option1=value1;option2=value2)
- `--tsp-client-option`: tsp-client options (e.g., --save-inputs, --debug)

**What it does:**

1. Creates or updates SDK packages for specified resource providers
2. Generates Go client code from TypeSpec or Swagger specifications
3. Creates appropriate changelogs and documentation
4. Handles versioning and breaking change detection
5. Optionally creates release branches for the changes

**Examples:**

```bash
# Generate release for a specific RP
generator release-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs network

# Generate with custom TypeSpec config
generator release-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs network \
  --tsp-config specification/network/tspconfig.yaml

# Generate from release request config file
generator release-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs config.json

# Skip branch creation and examples
generator release-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs compute \
  --skip-create-branch --skip-generate-example

# Force stable version generation even with preview input files
generator release-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs network \
  --force-stable-version
```

#### The `refresh-v2` command

The `refresh-v2` command regenerates all existing SDK packages using the latest specifications. This is useful for bulk updates across multiple packages.

**Usage:**

```bash
generator refresh-v2 <azure-sdk-for-go directory> <azure-rest-api-specs directory> [flags]
```

**Arguments:**

- `azure-sdk-for-go directory`: Path to azure-sdk-for-go repository
- `azure-rest-api-specs directory`: Path to azure-rest-api-specs repository

**Flags:**

- `--sdk-repo`: SDK repository URL
- `--spec-repo`: Spec repository URL
- `--release-date`: Release date for changelog
- `--skip-create-branch`: Skip creating release branch
- `--skip-generate-example`: Skip generating examples
- `--rps`: Comma-separated list of RPs to refresh (default: all)
- `--update-spec-version`: Whether to update commit ID (default: true)

**What it does:**

1. Discovers all existing SDK packages in the repository
2. Regenerates each package using current specifications
3. Updates version numbers and changelogs as needed
4. Handles batch processing of multiple resource providers
5. Optionally creates release branches for all changes

**Examples:**

```bash
# Refresh all existing packages
generator refresh-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs

# Refresh specific resource providers only
generator refresh-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --rps network,compute,storage

# Refresh without creating branches or examples
generator refresh-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --skip-create-branch --skip-generate-example

# Refresh with custom release date
generator refresh-v2 /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs \
  --release-date 2024-01-15
```

#### The `template` command

The `template` command generates package templates and scaffolding for new SDK packages. It creates the necessary directory structure and boilerplate code to onboard new services to the Azure SDK for Go.

**Usage:**

```bash
generator template (<rpName> <packageName>) | <packagePath>
```

**Arguments:**

- `rpName`: Name of the resource provider (e.g., "compute", "network")
- `packageName`: Name of the package (e.g., "armcompute")
- `packagePath`: Alternative format using `rpName/packageName` (e.g., "compute/armcompute")

**Flags:**

- `--go-sdk-folder`: Specifies the path of root of azure-sdk-for-go (default: ".")
- `--template-path`: Specifies the path of the template (default: "eng/tools/generator/template/rpName/packageName")
- `--package-title`: Specifies the title of this package (required)
- `--commit`: Specifies the commit hash of azure-rest-api-specs (required)
- `--release-date`: Specifies the release date in changelog
- `--package-config`: Additional config for package
- `--package-version`: Specify the version number of this release

**What it does:**

1. Creates the standard Azure SDK for Go package directory structure under `sdk/resourcemanager/<rpName>/<packageName>`
2. Reads template files from the template directory
3. Replaces placeholder values (rpName, packageName, packageTitle, commitID, releaseDate, etc.)
4. Generates the package files with proper content

**Examples:**

```bash
# Generate a template for a new service using two arguments
generator template compute armcompute --package-title "Compute" --commit abc123

# Generate a template using package path format
generator template compute/armcompute --package-title "Compute" --commit abc123

# Generate with custom release date and version
generator template network armnetwork --package-title "Network" --commit abc123 \
  --release-date 2024-01-15 --package-version 1.0.0

# Generate with custom SDK folder path
generator template storage armstorage --package-title "Storage" --commit abc123 \
  --go-sdk-folder /path/to/azure-sdk-for-go
```

#### The `generate-go-readme` command

The `generate-go-readme` command generates or updates a `readme.go.md` file for a resource provider based on the base swagger readme file.

**Usage:**

```bash
generator generate-go-readme <rp readme filepath>
```

**Arguments:**

- `rp readme filepath`: Path to the resource provider's readme.md file in the swagger spec repository

**What it does:**

1. Checks if the readme.md file exists at the specified path
2. Looks for an existing `readme.go.md` file in the same directory
3. If the Go track2 configuration section doesn't exist, appends the standard Go SDK configuration
4. Creates or updates the `readme.go.md` file with the appropriate module configuration

**Examples:**

```bash
# Generate readme.go.md for network RP
generator generate-go-readme /path/to/azure-rest-api-specs/specification/network/resource-manager/readme.md

# Generate readme.go.md for compute RP
generator generate-go-readme /path/to/azure-rest-api-specs/specification/compute/resource-manager/readme.md
```

**Generated Content:**
The command adds a Go track2 configuration section like:

```yaml
``` yaml $(go) && $(track2)
license-header: MICROSOFT_MIT_NO_VERSION
module-name: sdk/resourcemanager/<rpName>/arm<rpName>
module: github.com/Azure/azure-sdk-for-go/$(module-name)
output-folder: $(go-sdk-folder)/$(module-name)
azure-arm: true
```
