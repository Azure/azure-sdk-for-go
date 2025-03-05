# Overview

This guide describes the step-by-step process for [deprecating](https://go.dev/wiki/Deprecated) a Go module.

Deprecated modules are still available to use by the customer. When a module is deprecated, the Go toolchain will display a deprecation warning. This helps customers learn the module is deprecated and how to migrate their code.

Deprecation is different from [retraction](https://go.dev/ref/mod#go-mod-file-retract). Retractions are for removing problematic versions of a module. Use deprecation when you want to retire a module and transition customers to a new one.

# Step 1: Create Migration Guide

Lay out a migration plan for customers. Users should have a clear path to migrate their code from the old, deprecated module to the new, recommended module.

There should be a [`MIGRATION.md`](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/MIGRATION.md) file in the directory of the new module. This guide should include key differences between the modules and samples of the old vs new code.


# Step 2: Release the Module with a Deprecation Message

There should be a [final patch release](https://github.com/Azure/azure-sdk-for-go/pull/22578/files) of the module that includes the deprecation message. The deprecation message should also include migration instructions.

The deprecation message should be included in the following files:

### go.mod

You MUST include the deprecation message in the `go.mod` file. The [deprecation message](https://go.dev/ref/mod#go-mod-file-module-deprecation) starts with `// Deprecated: ` followed by migration instructions.

```
// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs instead
module github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest

go 1.18
```

### CHANGELOG.md

```md
# Release History

## 0.1.2 (2024-03-13)

### Other Changes
* This module is now DEPRECATED. The latest supported version of this module is at [github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs)

## 0.1.1 (2023-10-11)
```

### README.md

State the deprecation message at the beginning of the file, under the title.
```md
# Azure Monitor Ingestion client module for Go
> Deprecated: use [github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs) instead
The Azure Monitor Ingestion client module is used to send custom logs to [Azure Monitor][azure_monitor_overview] using the [Logs Ingestion API][ingestion_overview].
```

# Step 3: Request Removal from [pkg.go.dev](https://pkg.go.dev/)

After the final release with the deprecation message, it's time to remove the module from [pkg.go.dev](https://pkg.go.dev/about). Users will still be able to access the module via `go get` or `go install` but removing should prevent new users from discovering the deprecated module.

To do this, [file a request](https://github.com/golang/go/issues/66302) to the pkgsite team by filing an issue at https://github.com/golang/go/issues.

# Step 4: Delete Deprecated Module Code

Lastly, [remove the old module](https://github.com/Azure/azure-sdk-for-go/pull/22587/files) and all references to it from the `azure-sdk-for-go` repo. Double check all documentation (including MS Learn docs) for accuracy.
