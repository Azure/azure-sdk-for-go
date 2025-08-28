# Release History

## 0.3.0 (2025-08-28)
### Breaking Changes

- Function `*SolutionTemplatesClient.Update` parameter(s) have been changed from `(context.Context, string, string, SolutionTemplate, *SolutionTemplatesClientUpdateOptions)` to `(context.Context, string, string, SolutionTemplateUpdate, *SolutionTemplatesClientUpdateOptions)`

### Features Added

- New struct `SolutionTemplateUpdate`
- New struct `SolutionTemplateUpdateProperties`


## 0.2.0 (2025-08-27)
### Breaking Changes

- Function `*ConfigTemplatesClient.Update` parameter(s) have been changed from `(context.Context, string, string, ConfigTemplate, *ConfigTemplatesClientUpdateOptions)` to `(context.Context, string, string, ConfigTemplateUpdate, *ConfigTemplatesClientUpdateOptions)`
- Function `*ContextsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Context, *ContextsClientBeginUpdateOptions)` to `(context.Context, string, string, ContextUpdate, *ContextsClientBeginUpdateOptions)`
- Function `*DiagnosticsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Diagnostic, *DiagnosticsClientBeginUpdateOptions)` to `(context.Context, string, string, DiagnosticUpdate, *DiagnosticsClientBeginUpdateOptions)`
- Function `*SchemasClient.Update` parameter(s) have been changed from `(context.Context, string, string, Schema, *SchemasClientUpdateOptions)` to `(context.Context, string, string, SchemaUpdate, *SchemasClientUpdateOptions)`
- Function `*SolutionsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, Solution, *SolutionsClientBeginUpdateOptions)` to `(context.Context, string, string, string, SolutionUpdate, *SolutionsClientBeginUpdateOptions)`
- Function `*TargetsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Target, *TargetsClientBeginUpdateOptions)` to `(context.Context, string, string, TargetUpdate, *TargetsClientBeginUpdateOptions)`

### Features Added

- New struct `ConfigTemplateUpdate`
- New struct `ConfigTemplateUpdateProperties`
- New struct `ContextUpdate`
- New struct `ContextUpdateProperties`
- New struct `DiagnosticUpdate`
- New struct `DiagnosticUpdateProperties`
- New struct `SchemaUpdate`
- New struct `SchemaUpdateProperties`
- New struct `SolutionUpdate`
- New struct `SolutionUpdateProperties`
- New struct `TargetUpdate`
- New struct `TargetUpdateProperties`


## 0.1.0 (2025-08-13)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloadorchestration/armworkloadorchestration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).