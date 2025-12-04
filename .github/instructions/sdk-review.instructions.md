---
description: 'SDK Review Instructions for Go Packages'
applyTo: 'sdk/**'
---

When reviewing Go SDK packages under `sdk/` directory, follow this systematic process to ensure compliance with Azure SDK guidelines.

## Scope
Apply to all packages under `sdk/` excluding `fake/`, `testdata/`, and `*_test.go` files.

## Review Process

### Phase 1: Inventory & Discovery
Use `file_search` and `grep_search` to:
- Find all `*.go` files in the target module
- Locate key files: `client_factory.go`, `*_client.go`, `models.go`, `constants.go`, `options.go`, `responses.go`, `interfaces.go`
- Locate documentation: `README.md`, `CHANGELOG.md`
- Locate configuration: `go.mod`, `ci.yml`, `tsp-location.yaml`
- Map API surface

### Phase 2: API Design Review
Use `read_file` to examine code and `azure-sdk-code-review` tool to validate client factory, client methods, models, constants, options, and responses against guidelines.

**Review azure-sdk-code-review Results:**
- When the tool reports issues, VERIFY if the issue actually exists in the code
- Use `read_file` or `grep_search` to check if the flagged code/pattern is truly present
- If the issue is a FALSE POSITIVE (code is correct but tool misidentified it), do NOT include it in the report
- If the issue is REAL and present in the code, include it with severity assessment:
  - 🚫 Blocker: Violates required SDK guidelines, must fix before release
  - ⚠️ Warning: Violates recommended guidelines or best practices, should fix
  - ℹ️ Informational: Suggestions for improvement, optional
- For generated code, note that some violations may be acceptable if from code generator

### Phase 3: Documentation Review
Use `read_file` and `azure-sdk-code-review` tool to validate `README.md` and `CHANGELOG.md` against documentation standards.

**CHANGELOG.md Review:**
- Only review the **latest changelog entry** (the most recent version section)

**Review azure-sdk-code-review Results:**
- Verify that flagged documentation issues actually exist
- Check if required sections are truly missing or if tool missed them
- Only report confirmed documentation gaps or issues

### Phase 4: Configuration Review
Use `read_file` to check `go.mod`, `ci.yml`, and `tsp-location.yaml` against configuration requirements.

### Phase 5: Release Readiness
Verify final checklist items and produce assessment: READY / NOT READY / CONDITIONAL

## Tools
- `file_search` - Discover files by pattern
- `grep_search` - Find patterns across multiple files
- `read_file` - Read specific files or sections (use `limit`/`offset` for large files)
- `azure-sdk-code-review` - Validate code snippets against SDK guidelines (results must be verified before reporting)

## Important: Validating Code Review Results
The `azure-sdk-code-review` tool may occasionally produce false positives. Always:
1. **Verify the issue exists**: Use `read_file` or `grep_search` to confirm the flagged code pattern is actually present
2. **Understand the context**: Some patterns may be acceptable in generated code or specific scenarios
3. **Filter false positives**: Only include issues that you can confirm exist in the actual code
4. **Check for fixes**: If the tool reports something is missing (e.g., ResumeToken field), verify it's not present elsewhere in the file

## Output Format
Provide structured findings with:
- ✅ Compliant areas
- ⚠️ Warnings (non-blocking issues)
- 🚫 Blockers (must fix before release)
- ℹ️ Informational notes
- 📋 Release readiness assessment
- 🎯 Prioritized next steps

## Reference Documentation

### Internal References
- `documentation/development/README.md` - Development overview
- `documentation/development/ARM/new-version-guideline.md` - New version guidelines
- `documentation/development/ARM/new-version-quickstart.md` - Quickstart guide
- `documentation/development/ARM/go-mgmt-sdk-release-guideline.md` - Release guidelines
- `documentation/development/breaking-changes/sdk-breaking-changes-guide.md` - Breaking changes guide (for CHANGELOG.md review)
- `documentation/development/breaking-changes/deprecation-guide.md` - Deprecation guidelines
- `documentation/development/release.md` - Release process and checklist
- `documentation/development/generate.md` - Code generation guide
- `documentation/development/testing.md` - Testing guidelines

### External References
- [Azure Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
- [Release Notes Policy](https://azure.github.io/azure-sdk/policies_releasenotes.html)
- [Changelog Guidance](https://azure.github.io/azure-sdk/policies_releases.html#changelog-guidance)
