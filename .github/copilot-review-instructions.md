# Copilot Code Review Instructions

## Release Process Files

- Do not comment on CHANGELOG.md formatting or structure. The changelog follows a specific automated release process where the "Unreleased" header is intentionally removed when a version is released. A new "Unreleased" section is added as part of the next development cycle.
- Do not comment on files that only change a version string constant (e.g., `version.go`). These changes are part of the automated release process.

## Review Focus

- Focus on substantive code issues: bugs, security vulnerabilities, correctness, and performance problems.
- Do not comment on style, formatting, or whitespace unless it causes a functional issue.
