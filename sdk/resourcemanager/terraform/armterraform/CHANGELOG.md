# Release History

## 0.2.0 (2025-06-16)
### Breaking Changes

- Struct `ErrorAdditionalInfoInfo` has been removed

### Features Added

- Type of `ErrorAdditionalInfo.Info` has been changed from `*ErrorAdditionalInfoInfo` to `any`
- New enum type `AuthorizationScopeFilter` with values `AuthorizationScopeFilterAtScopeAboveAndBelow`, `AuthorizationScopeFilterAtScopeAndAbove`, `AuthorizationScopeFilterAtScopeAndBelow`, `AuthorizationScopeFilterAtScopeExact`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource` in struct `BaseExportModel`
- New field `AuthorizationScopeFilter`, `ExcludeAzureResource`, `ExcludeTerraformResource`, `Table` in struct `ExportQuery`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource` in struct `ExportResource`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource` in struct `ExportResourceGroup`
- New field `Import` in struct `ExportResult`


## 0.1.0 (2024-11-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/terraform/armterraform` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).