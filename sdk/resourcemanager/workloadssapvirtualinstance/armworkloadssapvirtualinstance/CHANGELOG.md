# Release History

## 1.0.0 (2025-05-21)
### Breaking Changes

- Type of `CreateAndMountFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `ErrorAdditionalInfo.Info` has been changed from `any` to `*ErrorAdditionalInfoInfo`
- Type of `FileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `MountFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `OperationStatusResult.PercentComplete` has been changed from `*float32` to `*float64`
- Type of `SAPVirtualInstance.Identity` has been changed from `*UserAssignedServiceIdentity` to `*SAPVirtualInstanceIdentity`
- Type of `SkipFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `UpdateSAPVirtualInstanceRequest.Identity` has been changed from `*UserAssignedServiceIdentity` to `*SAPVirtualInstanceIdentity`
- Enum `ConfigurationType` has been removed
- Enum `ManagedServiceIdentityType` has been removed
- Function `*ClientFactory.NewSAPCentralInstancesClient` has been removed
- Function `*ClientFactory.NewWorkloadsClient` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginStartInstance` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginStopInstance` has been removed
- Function `NewSAPCentralInstancesClient` has been removed
- Function `*SAPCentralInstancesClient.BeginCreate` has been removed
- Function `*SAPCentralInstancesClient.BeginDelete` has been removed
- Function `*SAPCentralInstancesClient.Get` has been removed
- Function `*SAPCentralInstancesClient.NewListPager` has been removed
- Function `*SAPCentralInstancesClient.BeginStartInstance` has been removed
- Function `*SAPCentralInstancesClient.BeginStopInstance` has been removed
- Function `*SAPCentralInstancesClient.Update` has been removed
- Function `*SAPDatabaseInstancesClient.BeginStartInstance` has been removed
- Function `*SAPDatabaseInstancesClient.BeginStopInstance` has been removed
- Function `NewWorkloadsClient` has been removed
- Function `*WorkloadsClient.SAPAvailabilityZoneDetails` has been removed
- Function `*WorkloadsClient.SAPDiskConfigurations` has been removed
- Function `*WorkloadsClient.SAPSizingRecommendations` has been removed
- Function `*WorkloadsClient.SAPSupportedSKU` has been removed
- Struct `SAPApplicationServerInstanceList` has been removed
- Struct `SAPCentralInstanceList` has been removed
- Struct `SAPDatabaseInstanceList` has been removed
- Struct `SAPVirtualInstanceList` has been removed
- Struct `UserAssignedServiceIdentity` has been removed
- Field `SAPApplicationServerInstanceList` of struct `SAPApplicationServerInstancesClientListResponse` has been removed
- Field `SAPDatabaseInstanceList` of struct `SAPDatabaseInstancesClientListResponse` has been removed
- Field `SAPVirtualInstanceList` of struct `SAPVirtualInstancesClientListByResourceGroupResponse` has been removed
- Field `SAPVirtualInstanceList` of struct `SAPVirtualInstancesClientListBySubscriptionResponse` has been removed

### Features Added

- New enum type `FileShareConfigurationType` with values `FileShareConfigurationTypeCreateAndMount`, `FileShareConfigurationTypeMount`, `FileShareConfigurationTypeSkip`
- New enum type `SAPVirtualInstanceIdentityType` with values `SAPVirtualInstanceIdentityTypeNone`, `SAPVirtualInstanceIdentityTypeUserAssigned`
- New function `*ClientFactory.NewSAPCentralServerInstancesClient() *SAPCentralServerInstancesClient`
- New function `*SAPApplicationServerInstancesClient.BeginStart(context.Context, string, string, string, *SAPApplicationServerInstancesClientBeginStartOptions) (*runtime.Poller[SAPApplicationServerInstancesClientStartResponse], error)`
- New function `*SAPApplicationServerInstancesClient.BeginStop(context.Context, string, string, string, *SAPApplicationServerInstancesClientBeginStopOptions) (*runtime.Poller[SAPApplicationServerInstancesClientStopResponse], error)`
- New function `NewSAPCentralServerInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SAPCentralServerInstancesClient, error)`
- New function `*SAPCentralServerInstancesClient.BeginCreate(context.Context, string, string, string, SAPCentralServerInstance, *SAPCentralServerInstancesClientBeginCreateOptions) (*runtime.Poller[SAPCentralServerInstancesClientCreateResponse], error)`
- New function `*SAPCentralServerInstancesClient.BeginDelete(context.Context, string, string, string, *SAPCentralServerInstancesClientBeginDeleteOptions) (*runtime.Poller[SAPCentralServerInstancesClientDeleteResponse], error)`
- New function `*SAPCentralServerInstancesClient.Get(context.Context, string, string, string, *SAPCentralServerInstancesClientGetOptions) (SAPCentralServerInstancesClientGetResponse, error)`
- New function `*SAPCentralServerInstancesClient.NewListPager(string, string, *SAPCentralServerInstancesClientListOptions) *runtime.Pager[SAPCentralServerInstancesClientListResponse]`
- New function `*SAPCentralServerInstancesClient.BeginStart(context.Context, string, string, string, *SAPCentralServerInstancesClientBeginStartOptions) (*runtime.Poller[SAPCentralServerInstancesClientStartResponse], error)`
- New function `*SAPCentralServerInstancesClient.BeginStop(context.Context, string, string, string, *SAPCentralServerInstancesClientBeginStopOptions) (*runtime.Poller[SAPCentralServerInstancesClientStopResponse], error)`
- New function `*SAPCentralServerInstancesClient.Update(context.Context, string, string, string, UpdateSAPCentralInstanceRequest, *SAPCentralServerInstancesClientUpdateOptions) (SAPCentralServerInstancesClientUpdateResponse, error)`
- New function `*SAPDatabaseInstancesClient.BeginStart(context.Context, string, string, string, *SAPDatabaseInstancesClientBeginStartOptions) (*runtime.Poller[SAPDatabaseInstancesClientStartResponse], error)`
- New function `*SAPDatabaseInstancesClient.BeginStop(context.Context, string, string, string, *SAPDatabaseInstancesClientBeginStopOptions) (*runtime.Poller[SAPDatabaseInstancesClientStopResponse], error)`
- New function `*SAPVirtualInstancesClient.GetAvailabilityZoneDetails(context.Context, string, SAPAvailabilityZoneDetailsRequest, *SAPVirtualInstancesClientGetAvailabilityZoneDetailsOptions) (SAPVirtualInstancesClientGetAvailabilityZoneDetailsResponse, error)`
- New function `*SAPVirtualInstancesClient.GetDiskConfigurations(context.Context, string, SAPDiskConfigurationsRequest, *SAPVirtualInstancesClientGetDiskConfigurationsOptions) (SAPVirtualInstancesClientGetDiskConfigurationsResponse, error)`
- New function `*SAPVirtualInstancesClient.GetSapSupportedSKU(context.Context, string, SAPSupportedSKUsRequest, *SAPVirtualInstancesClientGetSapSupportedSKUOptions) (SAPVirtualInstancesClientGetSapSupportedSKUResponse, error)`
- New function `*SAPVirtualInstancesClient.GetSizingRecommendations(context.Context, string, SAPSizingRecommendationRequest, *SAPVirtualInstancesClientGetSizingRecommendationsOptions) (SAPVirtualInstancesClientGetSizingRecommendationsResponse, error)`
- New struct `ErrorAdditionalInfoInfo`
- New struct `SAPApplicationServerInstanceListResult`
- New struct `SAPCentralServerInstanceListResult`
- New struct `SAPDatabaseInstanceListResult`
- New struct `SAPVirtualInstanceIdentity`
- New struct `SAPVirtualInstanceListResult`
- New field `ResourceID` in struct `OperationStatusResult`
- New anonymous field `SAPApplicationServerInstanceListResult` in struct `SAPApplicationServerInstancesClientListResponse`
- New anonymous field `SAPDatabaseInstanceListResult` in struct `SAPDatabaseInstancesClientListResponse`
- New anonymous field `SAPVirtualInstanceListResult` in struct `SAPVirtualInstancesClientListByResourceGroupResponse`
- New anonymous field `SAPVirtualInstanceListResult` in struct `SAPVirtualInstancesClientListBySubscriptionResponse`


## 0.1.0 (2024-02-23)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloadssapvirtualinstance/armworkloadssapvirtualinstance` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
