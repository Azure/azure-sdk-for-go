# Release History

## 0.2.0 (2026-08-27)
### Breaking Changes

- Function `*ClientFactory.NewTerraformClient` has been removed
- Function `NewTerraformClient` has been removed
- Function `*TerraformClient.BeginExportTerraform` has been removed
- Struct `ErrorAdditionalInfoInfo` has been removed

### Features Added

- Type of `ErrorAdditionalInfo.Info` has been changed from `*ErrorAdditionalInfoInfo` to `any`
- New enum type `AuthorizationScopeFilter` with values `AuthorizationScopeFilterAtScopeAboveAndBelow`, `AuthorizationScopeFilterAtScopeAndAbove`, `AuthorizationScopeFilterAtScopeAndBelow`, `AuthorizationScopeFilterAtScopeExact`
- New enum type `AzureExtensionResourceType` with values `AzureExtensionResourceTypeDiagnosticSettings`, `AzureExtensionResourceTypeLocks`, `AzureExtensionResourceTypeRoleAssignments`
- New function `NewClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.BeginExportTerraform(ctx context.Context, body BaseExportModelClassification, options *ClientBeginExportTerraformOptions) (*runtime.Poller[ClientExportTerraformResponse], error)`
- New function `*ClientFactory.NewClient() *Client`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource`, `IncludeExtensions`, `IncludeManagedResource`, `IncludeRoleAssignment` in struct `BaseExportModel`
- New field `AuthorizationScopeFilter`, `ExcludeAzureResource`, `ExcludeTerraformResource`, `IncludeExtensions`, `IncludeManagedResource`, `IncludeResourceGroup`, `IncludeRoleAssignment`, `Table` in struct `ExportQuery`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource`, `IncludeExtensions`, `IncludeManagedResource`, `IncludeResourceGroup`, `IncludeRoleAssignment`, `Recursive` in struct `ExportResource`
- New field `ExcludeAzureResource`, `ExcludeTerraformResource`, `IncludeExtensions`, `IncludeManagedResource`, `IncludeRoleAssignment` in struct `ExportResourceGroup`
- New field `Import` in struct `ExportResult`


## 0.1.0 (2024-11-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/terraform/armterraform` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).