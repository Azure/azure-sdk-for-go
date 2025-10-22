# Deprecating Azure SDK for Go Modules

This document provides comprehensive instructions for deprecating Go modules in the Azure SDK for Go repository. For general Azure SDK deprecation policies and guidelines, please refer to the [Azure SDK Deprecation Policies](https://azure.github.io/azure-sdk/policies_releases.html#deprecation).

## Overview

This guide describes the step-by-step process for [deprecating](https://go.dev/wiki/Deprecated) a Go module in the Azure SDK for Go.

Deprecated modules are still available to use by customers. When a module is deprecated, the Go toolchain will display a deprecation warning. This helps customers learn the module is deprecated and how to migrate their code.

Deprecation is different from [retraction](https://go.dev/ref/mod#go-mod-file-retract). Retractions are for removing problematic versions of a module. Use deprecation when you want to retire a module and transition customers to a new one.

## When to Deprecate

Consider deprecating a module when:
- The module's functionality is being replaced by a new module
- The module's underlying Azure service has been deprecated
- The module architecture needs significant changes that break compatibility
- The module is being consolidated with other modules

Before deprecating, ensure that:
- There is a clear migration path for customers
- The replacement module is stable and production-ready
- Adequate notice is given to customers (refer to [Azure SDK support policies](https://azure.github.io/azure-sdk/policies_support.html))

## Deprecation Process

### Step 1: Create Migration Guide

Lay out a migration plan for customers. Users should have a clear path to migrate their code from the old, deprecated module to the new, recommended module.

Create a `MIGRATION.md` file in the directory of the new module. Here is an example of a [`MIGRATION.md`](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/MIGRATION.md) file in the `azidentity` module. This guide should include:

- Key differences between the modules
- Code samples showing old vs new patterns
- Breaking changes and their resolutions
- Timeline for deprecation and support end-of-life

### Step 2: Release the Module with a Deprecation Message

Create a PR containing the final patch release of the module that includes the deprecation message. The deprecation message should also include migration instructions. Here is an example of the [final patch release](https://github.com/Azure/azure-sdk-for-go/pull/22578/files) for the `azingest` module.

The deprecation message should be included in the following files:

#### go.mod

You MUST include the deprecation message in the `go.mod` file. The [deprecation message](https://go.dev/ref/mod#go-mod-file-module-deprecation) starts with `// Deprecated: ` followed by migration instructions.

```go
// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs instead
module github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest

go 1.18
```

#### CHANGELOG.md

Add a deprecation notice to the changelog for the final release:

```md
# Release History

## 0.1.2 (2024-03-13)

### Other Changes
* This module is now DEPRECATED. The latest supported version of this module is at [github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs)

## 0.1.1 (2023-10-11)
```

#### README.md

State the deprecation message at the beginning of the file, under the title:

```md
# Azure Monitor Ingestion client module for Go
> **DEPRECATED**: use [github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs) instead

The Azure Monitor Ingestion client module is used to send custom logs to [Azure Monitor][azure_monitor_overview] using the [Logs Ingestion API][ingestion_overview].
```

### Step 3: Request Removal from pkg.go.dev

After the final release with the deprecation message, request removal of the module from [pkg.go.dev](https://pkg.go.dev/about). Users will still be able to access the module via `go get` or `go install`, but removal prevents new users from discovering the deprecated module.

To do this, file a request to the pkgsite team by filing an issue at https://github.com/golang/go/issues. Here is an example [request](https://github.com/golang/go/issues/66302) to remove `azingest`.

### Step 4: Update Documentation and References

- Update any documentation that references the deprecated module
- Remove the module from "getting started" guides and samples
- Update Azure documentation (MS Learn docs) for accuracy
- Ensure SDK overview pages reflect the current module recommendations

### Step 5: Delete Deprecated Module Code

After an appropriate deprecation period (minimum 12 months as per [Azure SDK support policies](https://azure.github.io/azure-sdk/policies_support.html#package-lifecycle)), remove the old module and all references to it from the `azure-sdk-for-go` repository. 

Before deletion:
- Ensure customers have had sufficient time to migrate
- Verify that usage metrics show minimal adoption of the deprecated module
- Confirm that replacement modules are stable and well-adopted

Here is an example [PR](https://github.com/Azure/azure-sdk-for-go/pull/22587/files) that removes the old `azingest` module code.

## Best Practices

1. **Communication**: Announce deprecations through appropriate channels (release notes, blog posts, Azure updates)
2. **Timeline**: Provide a clear timeline for deprecation milestones
3. **Support**: Continue to provide security fixes for deprecated modules during the deprecation period
4. **Documentation**: Keep migration guides up-to-date and easily discoverable
5. **Monitoring**: Track usage metrics to understand migration progress

## Related Documentation

- [Azure SDK Deprecation Policies](https://azure.github.io/azure-sdk/policies_releases.html#deprecation)
- [Azure SDK Support Policies](https://azure.github.io/azure-sdk/policies_support.html)
- [Go Module Deprecation](https://go.dev/wiki/Deprecated)
- [Azure SDK Release Guidelines](https://azure.github.io/azure-sdk/policies_releases.html)
- [Azure SDK for Go Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)

## Questions or Issues?

If you have questions about deprecating a module or need assistance with the process:
- File an issue in this repository
- Contact the Azure SDK team via the established support channels
- Refer to the [Azure SDK communication guidelines](https://azure.github.io/azure-sdk/policies_opensource.html)
