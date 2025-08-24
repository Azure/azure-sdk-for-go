# Azure Go SDK Automation Troubleshooting Guide

## Overview

The Azure Go SDK automation tool:

- Generates Go SDK code from TypeSpec/OpenAPI specifications
- Creates changelogs for API updates
- Detects breaking changes automatically

## Log Structure Analysis

### Key Log Keywords

1. **Start Marker**: `Reading generate input file from`
2. **Code Generation For TypeSpec**: `Start to process typespec project`
3. **Code Generation For Swagger**: `Start to process swagger project`

### Success Indicators

- **For Each TypeSpec Specification's Generation**: `Finish processing typespec project`
- **For Each Swagger Specification's Generation**: `Finish processing swagger project`
- **Output File**: contains both changelog details and breaking changes summary

### Failure Indicators

- Any `[ERROR]` messages in the log
- Non-zero exit codes
- Exception stack traces

## Error Classification and Resolution

### 1. Internal Errors

**Log Keywords**: `The emitter encountered an internal error during preprocessing.`

**Context**: This indicates a failure in the TypeSpec emitter, which should be a bug in the `@azure-tools/typespec-go` package.

**Analysis Actions**:

- Extract package information
- Capture complete error message and stack trace
- Identify the failing emitter component

**Resolution**: File issue at https://github.com/Azure/autorest.go/issues with:

- Complete error details
- Package/service context
- Stack trace information

### 2. TypeSpec Configuration Errors

#### Invalid Emitter Arguments

**Log Keywords**: `Invalid arguments were passed to the emitter.`

**Context**: This happens when config is not set correctly.

**Analysis Actions**:

- Check the error log to see which config is wrong.
- Refer the [doc](https://azure.github.io/typespec-azure/docs/emitters/clients/typespec-go/reference/emitter/#emitter-options) for the usage of emitter config
- Compare against standard templates:
  - [Management Plane Template](https://github.com/Azure/azure-rest-api-specs/blob/a8f97793186c7680c62519da238c6d07a20f2023/specification/contosowidgetmanager/Contoso.Management/tspconfig.yaml#L35)
  - [Data Plane Template](https://github.com/Azure/azure-rest-api-specs/blob/a8f97793186c7680c62519da238c6d07a20f2023/specification/contosowidgetmanager/Contoso.WidgetManager/tspconfig.yaml#L41)

**Resolution**: Fix `@azure-tools/typespec-go` configuration in `tspconfig.yaml`

#### Module Path Errors

**Log Keywords**: `module not found, package path:`

**Context**: This happens when emitter could not resolve the package you want to generate through configs.

**Analysis Actions**:

- Validate module path configuration
- Check service directory structure
- Verify package directory alignment
- Refer the [doc](https://azure.github.io/typespec-azure/docs/emitters/clients/typespec-go/reference/emitter/#emitter-options) for the usage of emitter config

**Resolution**: Correct these `tspconfig.yaml` properties:

- `module`: Go module path
- `service-dir`: Service directory path
- `package-dir`: Package directory path
- `module-version`: Module version

### 3. Naming Collision Errors

**Log Keywords**: `The emitter automatically renamed one or more types which resulted in a type name collision.`

**Context**: Go SDK automatically removes service name prefixes to prevent stuttering (e.g., `armcompute.ComputeDisk` → `armcompute.Disk`)

**Analysis Actions**:

- Identify conflicting type names
- Locate the original model definitions
- Determine rename conflicts

**Resolution**: Rename conflicting models in `client.tsp` using `@clientName` decorator

### 4. Unsupported TypeSpec Features

**Log Keywords**: `UnsupportedTsp`

**Common Unsupported Features**:

- Paging with re-injected parameters
- Cookie parameters

**Analysis Actions**:

- Identify specific unsupported feature
- Review error message for feature details

**Resolution**: Modify TypeSpec authoring to use supported patterns if possible, or wait for emitter support.
