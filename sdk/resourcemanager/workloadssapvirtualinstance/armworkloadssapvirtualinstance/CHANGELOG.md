# Release History

## 1.0.0 (2025-01-14)
### Breaking Changes

- Type of `CreateAndMountFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `FileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `MountFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `SAPVirtualInstance.Identity` has been changed from `*UserAssignedServiceIdentity` to `*SAPVirtualInstanceIdentity`
- Type of `SkipFileShareConfiguration.ConfigurationType` has been changed from `*ConfigurationType` to `*FileShareConfigurationType`
- Type of `UpdateSAPVirtualInstanceRequest.Identity` has been changed from `*UserAssignedServiceIdentity` to `*SAPVirtualInstanceIdentity`
- Enum `ConfigurationType` has been removed
- Enum `ManagedServiceIdentityType` has been removed
- Function `*ClientFactory.NewSAPApplicationServerInstancesClient` has been removed
- Function `*ClientFactory.NewSAPCentralInstancesClient` has been removed
- Function `*ClientFactory.NewSAPDatabaseInstancesClient` has been removed
- Function `*ClientFactory.NewSAPVirtualInstancesClient` has been removed
- Function `*ClientFactory.NewWorkloadsClient` has been removed
- Function `NewSAPApplicationServerInstancesClient` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginCreate` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginDelete` has been removed
- Function `*SAPApplicationServerInstancesClient.Get` has been removed
- Function `*SAPApplicationServerInstancesClient.NewListPager` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginStartInstance` has been removed
- Function `*SAPApplicationServerInstancesClient.BeginStopInstance` has been removed
- Function `*SAPApplicationServerInstancesClient.Update` has been removed
- Function `NewSAPCentralInstancesClient` has been removed
- Function `*SAPCentralInstancesClient.BeginCreate` has been removed
- Function `*SAPCentralInstancesClient.BeginDelete` has been removed
- Function `*SAPCentralInstancesClient.Get` has been removed
- Function `*SAPCentralInstancesClient.NewListPager` has been removed
- Function `*SAPCentralInstancesClient.BeginStartInstance` has been removed
- Function `*SAPCentralInstancesClient.BeginStopInstance` has been removed
- Function `*SAPCentralInstancesClient.Update` has been removed
- Function `NewSAPDatabaseInstancesClient` has been removed
- Function `*SAPDatabaseInstancesClient.BeginCreate` has been removed
- Function `*SAPDatabaseInstancesClient.BeginDelete` has been removed
- Function `*SAPDatabaseInstancesClient.Get` has been removed
- Function `*SAPDatabaseInstancesClient.NewListPager` has been removed
- Function `*SAPDatabaseInstancesClient.BeginStartInstance` has been removed
- Function `*SAPDatabaseInstancesClient.BeginStopInstance` has been removed
- Function `*SAPDatabaseInstancesClient.Update` has been removed
- Function `NewSAPVirtualInstancesClient` has been removed
- Function `*SAPVirtualInstancesClient.BeginCreate` has been removed
- Function `*SAPVirtualInstancesClient.BeginDelete` has been removed
- Function `*SAPVirtualInstancesClient.Get` has been removed
- Function `*SAPVirtualInstancesClient.NewListByResourceGroupPager` has been removed
- Function `*SAPVirtualInstancesClient.NewListBySubscriptionPager` has been removed
- Function `*SAPVirtualInstancesClient.BeginStart` has been removed
- Function `*SAPVirtualInstancesClient.BeginStop` has been removed
- Function `*SAPVirtualInstancesClient.BeginUpdate` has been removed
- Function `NewWorkloadsClient` has been removed
- Function `*WorkloadsClient.SAPAvailabilityZoneDetails` has been removed
- Function `*WorkloadsClient.SAPDiskConfigurations` has been removed
- Function `*WorkloadsClient.SAPSizingRecommendations` has been removed
- Function `*WorkloadsClient.SAPSupportedSKU` has been removed
- Struct `SAPApplicationServerInstanceList` has been removed
- Struct `SAPCentralInstanceList` has been removed
- Struct `SAPDatabaseInstanceList` has been removed
- Struct `SAPVirtualInstanceList` has been removed
- Struct `UserAssignedIdentity` has been removed
- Struct `UserAssignedServiceIdentity` has been removed

### Features Added

- New enum type `FileShareConfigurationType` with values `FileShareConfigurationTypeCreateAndMount`, `FileShareConfigurationTypeMount`, `FileShareConfigurationTypeSkip`
- New enum type `SAPVirtualInstanceIdentityType` with values `SAPVirtualInstanceIdentityTypeNone`, `SAPVirtualInstanceIdentityTypeUserAssigned`
- New function `*ClientFactory.NewSapApplicationServerInstancesClient() *SapApplicationServerInstancesClient`
- New function `*ClientFactory.NewSapCentralServerInstancesClient() *SapCentralServerInstancesClient`
- New function `*ClientFactory.NewSapDatabaseInstancesClient() *SapDatabaseInstancesClient`
- New function `*ClientFactory.NewSapVirtualInstancesClient() *SapVirtualInstancesClient`
- New function `NewSapApplicationServerInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SapApplicationServerInstancesClient, error)`
- New function `*SapApplicationServerInstancesClient.BeginCreate(context.Context, string, string, string, SAPApplicationServerInstance, *SapApplicationServerInstancesClientBeginCreateOptions) (*runtime.Poller[SapApplicationServerInstancesClientCreateResponse], error)`
- New function `*SapApplicationServerInstancesClient.BeginDelete(context.Context, string, string, string, *SapApplicationServerInstancesClientBeginDeleteOptions) (*runtime.Poller[SapApplicationServerInstancesClientDeleteResponse], error)`
- New function `*SapApplicationServerInstancesClient.Get(context.Context, string, string, string, *SapApplicationServerInstancesClientGetOptions) (SapApplicationServerInstancesClientGetResponse, error)`
- New function `*SapApplicationServerInstancesClient.NewListPager(string, string, *SapApplicationServerInstancesClientListOptions) *runtime.Pager[SapApplicationServerInstancesClientListResponse]`
- New function `*SapApplicationServerInstancesClient.BeginStart(context.Context, string, string, string, *SapApplicationServerInstancesClientBeginStartOptions) (*runtime.Poller[SapApplicationServerInstancesClientStartResponse], error)`
- New function `*SapApplicationServerInstancesClient.BeginStop(context.Context, string, string, string, *SapApplicationServerInstancesClientBeginStopOptions) (*runtime.Poller[SapApplicationServerInstancesClientStopResponse], error)`
- New function `*SapApplicationServerInstancesClient.Update(context.Context, string, string, string, UpdateSAPApplicationInstanceRequest, *SapApplicationServerInstancesClientUpdateOptions) (SapApplicationServerInstancesClientUpdateResponse, error)`
- New function `NewSapCentralServerInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SapCentralServerInstancesClient, error)`
- New function `*SapCentralServerInstancesClient.BeginCreate(context.Context, string, string, string, SAPCentralServerInstance, *SapCentralServerInstancesClientBeginCreateOptions) (*runtime.Poller[SapCentralServerInstancesClientCreateResponse], error)`
- New function `*SapCentralServerInstancesClient.BeginDelete(context.Context, string, string, string, *SapCentralServerInstancesClientBeginDeleteOptions) (*runtime.Poller[SapCentralServerInstancesClientDeleteResponse], error)`
- New function `*SapCentralServerInstancesClient.Get(context.Context, string, string, string, *SapCentralServerInstancesClientGetOptions) (SapCentralServerInstancesClientGetResponse, error)`
- New function `*SapCentralServerInstancesClient.NewListPager(string, string, *SapCentralServerInstancesClientListOptions) *runtime.Pager[SapCentralServerInstancesClientListResponse]`
- New function `*SapCentralServerInstancesClient.BeginStart(context.Context, string, string, string, *SapCentralServerInstancesClientBeginStartOptions) (*runtime.Poller[SapCentralServerInstancesClientStartResponse], error)`
- New function `*SapCentralServerInstancesClient.BeginStop(context.Context, string, string, string, *SapCentralServerInstancesClientBeginStopOptions) (*runtime.Poller[SapCentralServerInstancesClientStopResponse], error)`
- New function `*SapCentralServerInstancesClient.Update(context.Context, string, string, string, UpdateSAPCentralInstanceRequest, *SapCentralServerInstancesClientUpdateOptions) (SapCentralServerInstancesClientUpdateResponse, error)`
- New function `NewSapDatabaseInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SapDatabaseInstancesClient, error)`
- New function `*SapDatabaseInstancesClient.BeginCreate(context.Context, string, string, string, SAPDatabaseInstance, *SapDatabaseInstancesClientBeginCreateOptions) (*runtime.Poller[SapDatabaseInstancesClientCreateResponse], error)`
- New function `*SapDatabaseInstancesClient.BeginDelete(context.Context, string, string, string, *SapDatabaseInstancesClientBeginDeleteOptions) (*runtime.Poller[SapDatabaseInstancesClientDeleteResponse], error)`
- New function `*SapDatabaseInstancesClient.Get(context.Context, string, string, string, *SapDatabaseInstancesClientGetOptions) (SapDatabaseInstancesClientGetResponse, error)`
- New function `*SapDatabaseInstancesClient.NewListPager(string, string, *SapDatabaseInstancesClientListOptions) *runtime.Pager[SapDatabaseInstancesClientListResponse]`
- New function `*SapDatabaseInstancesClient.BeginStart(context.Context, string, string, string, *SapDatabaseInstancesClientBeginStartOptions) (*runtime.Poller[SapDatabaseInstancesClientStartResponse], error)`
- New function `*SapDatabaseInstancesClient.BeginStop(context.Context, string, string, string, *SapDatabaseInstancesClientBeginStopOptions) (*runtime.Poller[SapDatabaseInstancesClientStopResponse], error)`
- New function `*SapDatabaseInstancesClient.Update(context.Context, string, string, string, UpdateSAPDatabaseInstanceRequest, *SapDatabaseInstancesClientUpdateOptions) (SapDatabaseInstancesClientUpdateResponse, error)`
- New function `NewSapVirtualInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SapVirtualInstancesClient, error)`
- New function `*SapVirtualInstancesClient.BeginCreate(context.Context, string, string, SAPVirtualInstance, *SapVirtualInstancesClientBeginCreateOptions) (*runtime.Poller[SapVirtualInstancesClientCreateResponse], error)`
- New function `*SapVirtualInstancesClient.BeginDelete(context.Context, string, string, *SapVirtualInstancesClientBeginDeleteOptions) (*runtime.Poller[SapVirtualInstancesClientDeleteResponse], error)`
- New function `*SapVirtualInstancesClient.Get(context.Context, string, string, *SapVirtualInstancesClientGetOptions) (SapVirtualInstancesClientGetResponse, error)`
- New function `*SapVirtualInstancesClient.InvokeAvailabilityZoneDetails(context.Context, string, SAPAvailabilityZoneDetailsRequest, *SapVirtualInstancesClientInvokeAvailabilityZoneDetailsOptions) (SapVirtualInstancesClientInvokeAvailabilityZoneDetailsResponse, error)`
- New function `*SapVirtualInstancesClient.InvokeDiskConfigurations(context.Context, string, SAPDiskConfigurationsRequest, *SapVirtualInstancesClientInvokeDiskConfigurationsOptions) (SapVirtualInstancesClientInvokeDiskConfigurationsResponse, error)`
- New function `*SapVirtualInstancesClient.InvokeSapSupportedSKU(context.Context, string, SAPSupportedSKUsRequest, *SapVirtualInstancesClientInvokeSapSupportedSKUOptions) (SapVirtualInstancesClientInvokeSapSupportedSKUResponse, error)`
- New function `*SapVirtualInstancesClient.InvokeSizingRecommendations(context.Context, string, SAPSizingRecommendationRequest, *SapVirtualInstancesClientInvokeSizingRecommendationsOptions) (SapVirtualInstancesClientInvokeSizingRecommendationsResponse, error)`
- New function `*SapVirtualInstancesClient.NewListByResourceGroupPager(string, *SapVirtualInstancesClientListByResourceGroupOptions) *runtime.Pager[SapVirtualInstancesClientListByResourceGroupResponse]`
- New function `*SapVirtualInstancesClient.NewListBySubscriptionPager(*SapVirtualInstancesClientListBySubscriptionOptions) *runtime.Pager[SapVirtualInstancesClientListBySubscriptionResponse]`
- New function `*SapVirtualInstancesClient.BeginStart(context.Context, string, string, *SapVirtualInstancesClientBeginStartOptions) (*runtime.Poller[SapVirtualInstancesClientStartResponse], error)`
- New function `*SapVirtualInstancesClient.BeginStop(context.Context, string, string, *SapVirtualInstancesClientBeginStopOptions) (*runtime.Poller[SapVirtualInstancesClientStopResponse], error)`
- New function `*SapVirtualInstancesClient.BeginUpdate(context.Context, string, string, UpdateSAPVirtualInstanceRequest, *SapVirtualInstancesClientBeginUpdateOptions) (*runtime.Poller[SapVirtualInstancesClientUpdateResponse], error)`
- New struct `Components1IrwhnvSchemasSapvirtualinstanceidentityPropertiesUserassignedidentitiesAdditionalproperties`
- New struct `SAPApplicationServerInstanceListResult`
- New struct `SAPCentralServerInstanceListResult`
- New struct `SAPDatabaseInstanceListResult`
- New struct `SAPVirtualInstanceIdentity`
- New struct `SAPVirtualInstanceListResult`
- New field `ResourceID` in struct `OperationStatusResult`


## 0.1.0 (2024-02-23)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloadssapvirtualinstance/armworkloadssapvirtualinstance` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).