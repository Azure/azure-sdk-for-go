---
name: azsdk-common-generate-sdk-locally
license: MIT
metadata:
  version: "1.0.0"
description: "**UTILITY SKILL** Generate, build, and test Azure SDKs locally from TypeSpec. USE FOR: \"generate SDK locally\", \"build SDK\", \"run SDK tests\", \"update changelog\". INVOKES: azsdk_verify_setup, azsdk_package_generate_code, azsdk_package_build_code, azsdk_package_run_check, azsdk_package_run_tests, azsdk_package_update_metadata, azsdk_package_update_changelog_content, azsdk_package_update_version."
compatibility:
  requires: "azure-sdk-mcp server, local azure-sdk-for-{language} clone, language build tools"
---

# Generate SDK Locally

## MCP Prerequisites

Requires `azure-sdk-mcp` server. Run `azsdk_verify_setup` to confirm.

## MCP Tools

| Tool                          | Purpose                    |
| ----------------------------- | -------------------------- |
| `azsdk_package_generate_code` | Generate SDK from TypeSpec |
| `azsdk_package_build_code`    | Build package              |
| `azsdk_package_run_check`     | Validate package           |
| `azsdk_package_run_tests`     | Run tests                  |

## Steps

1. **Select Language** — .NET, Java, JavaScript, Python, or Go.
2. **Verify** — Run `azsdk_verify_setup` to confirm environment.
3. **Generate** — Run `azsdk_package_generate_code` with the path to `tspconfig.yaml` (spec repo) or `tsp-location.yaml` (SDK repo). Accepts local paths or HTTPS URLs.
4. **Build** — Run `azsdk_package_build_code`. On failure, use typespec-customization.
5. **Validate** — Run `azsdk_package_run_check` and `azsdk_package_run_tests`.
6. **Metadata** — Update metadata, changelog, and version.
7. **Next Steps** — Push API spec and SDK changes to create PRs.

[SDK repos](references/sdk-repos.md)

## CLI Fallback

Without MCP, use `npx tsp-client` directly:
- Spec repo: `npx tsp-client init --update-if-exists`
- SDK repo: `npx tsp-client update`

Then use language build tools manually.
