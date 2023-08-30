# Release History

## 2.0.0 (2023-08-30)
### Breaking Changes

- Operation `*ApplicationDefinitionsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ApplicationDefinitionsClient.CreateOrUpdate` instead.
- Operation `*ApplicationDefinitionsClient.BeginDelete` has been changed to non-LRO, use `*ApplicationDefinitionsClient.Delete` instead.
- Operation `*ApplicationsClient.Update` has been changed to LRO, use `*ApplicationsClient.BeginUpdate` instead.
- Struct `ApplicationPropertiesPatchable` has been removed
- Field `ProvisioningState` of struct `ApplicationDefinitionProperties` has been removed

### Features Added

- New enum type `Status` with values `StatusElevate`, `StatusNotSpecified`, `StatusRemove`
- New enum type `Substatus` with values `SubstatusApproved`, `SubstatusDenied`, `SubstatusExpired`, `SubstatusFailed`, `SubstatusNotSpecified`, `SubstatusTimeout`
- New function `*ApplicationDefinitionsClient.CreateOrUpdateByID(context.Context, string, string, ApplicationDefinition, *ApplicationDefinitionsClientCreateOrUpdateByIDOptions) (ApplicationDefinitionsClientCreateOrUpdateByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.DeleteByID(context.Context, string, string, *ApplicationDefinitionsClientDeleteByIDOptions) (ApplicationDefinitionsClientDeleteByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.GetByID(context.Context, string, string, *ApplicationDefinitionsClientGetByIDOptions) (ApplicationDefinitionsClientGetByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.UpdateByID(context.Context, string, string, ApplicationDefinitionPatchable, *ApplicationDefinitionsClientUpdateByIDOptions) (ApplicationDefinitionsClientUpdateByIDResponse, error)`
- New function `*ApplicationsClient.BeginCreateOrUpdateByID(context.Context, string, Application, *ApplicationsClientBeginCreateOrUpdateByIDOptions) (*runtime.Poller[ApplicationsClientCreateOrUpdateByIDResponse], error)`
- New function `*ApplicationsClient.BeginDeleteByID(context.Context, string, *ApplicationsClientBeginDeleteByIDOptions) (*runtime.Poller[ApplicationsClientDeleteByIDResponse], error)`
- New function `*ApplicationsClient.GetByID(context.Context, string, *ApplicationsClientGetByIDOptions) (ApplicationsClientGetByIDResponse, error)`
- New function `*ApplicationsClient.ListTokens(context.Context, string, string, ListTokenRequest, *ApplicationsClientListTokensOptions) (ApplicationsClientListTokensResponse, error)`
- New function `*ApplicationsClient.BeginUpdateAccess(context.Context, string, string, UpdateAccessDefinition, *ApplicationsClientBeginUpdateAccessOptions) (*runtime.Poller[ApplicationsClientUpdateAccessResponse], error)`
- New function `*ApplicationsClient.BeginUpdateByID(context.Context, string, *ApplicationsClientBeginUpdateByIDOptions) (*runtime.Poller[ApplicationsClientUpdateByIDResponse], error)`
- New struct `AllowedUpgradePlansResult`
- New struct `JitRequestMetadata`
- New struct `ListTokenRequest`
- New struct `ManagedIdentityToken`
- New struct `ManagedIdentityTokenResult`
- New struct `UpdateAccessDefinition`
- New anonymous field `AllowedUpgradePlansResult` in struct `ApplicationsClientListAllowedUpgradePlansResponse`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/solutions/armmanagedapplications` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).