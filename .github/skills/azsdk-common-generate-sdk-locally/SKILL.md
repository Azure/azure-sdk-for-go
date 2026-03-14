---
name: azsdk-common-generate-sdk-locally
license: MIT
metadata:
  version: "1.0.0"
description: "Generate, build, and test Azure SDKs locally from TypeSpec. **UTILITY SKILL**. USE FOR: \"generate SDK locally\", \"build SDK\", \"run SDK tests\", \"update changelog\". DO NOT USE FOR: publishing to package registries, CI pipeline configuration, API design review. INVOKES: azsdk_verify_setup, azsdk_package_generate_code, azsdk_package_build_code, azsdk_package_run_check, azsdk_package_run_tests, azsdk_package_update_metadata, azsdk_package_update_changelog_content, azsdk_package_update_version."
compatibility:
  requires: "azure-sdk-mcp server, local azure-sdk-for-{language} clone, language build tools"
---

# Generate SDK Locally

## MCP Tools

| Tool | Purpose |
|------|---------|
| `azsdk_package_generate_code` | Generate SDK from TypeSpec |
| `azsdk_package_build_code` | Build package |
| `azsdk_package_run_check` | Validate package |
| `azsdk_package_run_tests` | Run tests |

**Prerequisites:** azure-sdk-mcp server must be running. Without MCP, use `npx tsp-client` CLI.

## Steps

1. **Verify** — Run `azsdk_verify_setup` to confirm environment.
2. **Generate** — Run `azsdk_package_generate_code` with `tspconfig.yaml` or `tsp-location.yaml` path (local or HTTPS).
3. **Build** — Run `azsdk_package_build_code`. On failure, use typespec-customization.
4. **Validate** — Run `azsdk_package_run_check` and `azsdk_package_run_tests`.
5. **Metadata** — Update metadata, changelog, and version.

[SDK repos](references/sdk-repos.md)

## Examples

- "Generate the SDK locally for my TypeSpec service"
- "Build and test the Python SDK package"

## Troubleshooting

- Run `azsdk_verify_setup` first to confirm MCP and tools are ready.
- Without MCP, use `npx tsp-client init` (spec repo) or `npx tsp-client update` (SDK repo).
