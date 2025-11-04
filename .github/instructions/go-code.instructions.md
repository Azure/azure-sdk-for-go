---
applyTo: '**/*.go'
---

All code should follow the guidelines from the [Azure Go SDK Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html). This document is a summary of the most important guidelines to follow when contributing to the Azure Go SDK.

Some additional rules:
- Acronyms in exported type/struct/function names (for instance, 'id' in UserID, or 'acs' in ACSRecordedEvent) should be uppercased. This rule does not apply to string constants.
- All Go files should have a copyright header. The header might be after go build directives, like "//go:build go1.18"
