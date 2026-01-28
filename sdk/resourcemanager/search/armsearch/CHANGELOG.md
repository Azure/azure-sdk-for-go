# Release History

## 2.0.0 (2026-01-28)
### Breaking Changes

- Function `*AdminKeysClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *AdminKeysClientGetOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *AdminKeysClientGetOptions)`
- Function `*AdminKeysClient.Regenerate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, keyKind AdminKeyKind, searchManagementRequestOptions *SearchManagementRequestOptions, options *AdminKeysClientRegenerateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, keyKind AdminKeyKind, searchManagementRequestOptions ManagementRequestOptions, options *AdminKeysClientRegenerateOptions)`
- Function `*ManagementClient.UsageBySubscriptionSKU` parameter(s) have been changed from `(ctx context.Context, location string, skuName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *ManagementClientUsageBySubscriptionSKUOptions)` to `(ctx context.Context, location string, skuName string, searchManagementRequestOptions ManagementRequestOptions, options *ManagementClientUsageBySubscriptionSKUOptions)`
- Function `*PrivateEndpointConnectionsClient.Delete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *PrivateEndpointConnectionsClientDeleteOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, searchManagementRequestOptions ManagementRequestOptions, options *PrivateEndpointConnectionsClientDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *PrivateEndpointConnectionsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, searchManagementRequestOptions ManagementRequestOptions, options *PrivateEndpointConnectionsClientGetOptions)`
- Function `*PrivateEndpointConnectionsClient.NewListByServicePager` parameter(s) have been changed from `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *PrivateEndpointConnectionsClientListByServiceOptions)` to `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *PrivateEndpointConnectionsClientListByServiceOptions)`
- Function `*PrivateEndpointConnectionsClient.Update` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, privateEndpointConnection PrivateEndpointConnection, searchManagementRequestOptions *SearchManagementRequestOptions, options *PrivateEndpointConnectionsClientUpdateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, privateEndpointConnectionName string, privateEndpointConnection PrivateEndpointConnection, searchManagementRequestOptions ManagementRequestOptions, options *PrivateEndpointConnectionsClientUpdateOptions)`
- Function `*PrivateLinkResourcesClient.NewListSupportedPager` parameter(s) have been changed from `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *PrivateLinkResourcesClientListSupportedOptions)` to `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *PrivateLinkResourcesClientListSupportedOptions)`
- Function `*QueryKeysClient.Create` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, name string, searchManagementRequestOptions *SearchManagementRequestOptions, options *QueryKeysClientCreateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, name string, searchManagementRequestOptions ManagementRequestOptions, options *QueryKeysClientCreateOptions)`
- Function `*QueryKeysClient.Delete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, key string, searchManagementRequestOptions *SearchManagementRequestOptions, options *QueryKeysClientDeleteOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, key string, searchManagementRequestOptions ManagementRequestOptions, options *QueryKeysClientDeleteOptions)`
- Function `*QueryKeysClient.NewListBySearchServicePager` parameter(s) have been changed from `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *QueryKeysClientListBySearchServiceOptions)` to `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *QueryKeysClientListBySearchServiceOptions)`
- Function `*ServicesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, service Service, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, resource Service, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientBeginCreateOrUpdateOptions)`
- Function `*ServicesClient.CheckNameAvailability` parameter(s) have been changed from `(ctx context.Context, checkNameAvailabilityInput CheckNameAvailabilityInput, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientCheckNameAvailabilityOptions)` to `(ctx context.Context, checkNameAvailabilityInput CheckNameAvailabilityInput, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientCheckNameAvailabilityOptions)`
- Function `*ServicesClient.Delete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientDeleteOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientDeleteOptions)`
- Function `*ServicesClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientGetOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientGetOptions)`
- Function `*ServicesClient.NewListByResourceGroupPager` parameter(s) have been changed from `(resourceGroupName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientListByResourceGroupOptions)` to `(resourceGroupName string, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientListByResourceGroupOptions)`
- Function `*ServicesClient.NewListBySubscriptionPager` parameter(s) have been changed from `(searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientListBySubscriptionOptions)` to `(searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientListBySubscriptionOptions)`
- Function `*ServicesClient.Update` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, service ServiceUpdate, searchManagementRequestOptions *SearchManagementRequestOptions, options *ServicesClientUpdateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, properties ServiceUpdate, searchManagementRequestOptions ManagementRequestOptions, options *ServicesClientUpdateOptions)`
- Function `*SharedPrivateLinkResourcesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, sharedPrivateLinkResource SharedPrivateLinkResource, searchManagementRequestOptions *SearchManagementRequestOptions, options *SharedPrivateLinkResourcesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, sharedPrivateLinkResource SharedPrivateLinkResource, searchManagementRequestOptions ManagementRequestOptions, options *SharedPrivateLinkResourcesClientBeginCreateOrUpdateOptions)`
- Function `*SharedPrivateLinkResourcesClient.BeginDelete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *SharedPrivateLinkResourcesClientBeginDeleteOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, searchManagementRequestOptions ManagementRequestOptions, options *SharedPrivateLinkResourcesClientBeginDeleteOptions)`
- Function `*SharedPrivateLinkResourcesClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *SharedPrivateLinkResourcesClientGetOptions)` to `(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, searchManagementRequestOptions ManagementRequestOptions, options *SharedPrivateLinkResourcesClientGetOptions)`
- Function `*SharedPrivateLinkResourcesClient.NewListByServicePager` parameter(s) have been changed from `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions *SearchManagementRequestOptions, options *SharedPrivateLinkResourcesClientListByServiceOptions)` to `(resourceGroupName string, searchServiceName string, searchManagementRequestOptions ManagementRequestOptions, options *SharedPrivateLinkResourcesClientListByServiceOptions)`
- Function `*UsagesClient.NewListBySubscriptionPager` parameter(s) have been changed from `(location string, searchManagementRequestOptions *SearchManagementRequestOptions, options *UsagesClientListBySubscriptionOptions)` to `(location string, searchManagementRequestOptions ManagementRequestOptions, options *UsagesClientListBySubscriptionOptions)`
- Struct `SearchManagementRequestOptions` has been removed

### Features Added

- New struct `ManagementRequestOptions`
- New field `NextLink` in struct `PrivateLinkResourcesResult`


## 1.4.0 (2025-07-21)
### Features Added

- New value `IdentityTypeSystemAssignedUserAssigned`, `IdentityTypeUserAssigned` added to enum type `IdentityType`
- New value `PublicNetworkAccessSecuredByPerimeter` added to enum type `PublicNetworkAccess`
- New value `SearchServiceStatusStopped` added to enum type `SearchServiceStatus`
- New enum type `AccessRuleDirection` with values `AccessRuleDirectionInbound`, `AccessRuleDirectionOutbound`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `ComputeType` with values `ComputeTypeConfidential`, `ComputeTypeDefault`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `IssueType` with values `IssueTypeConfigurationPropagationFailure`, `IssueTypeMissingIdentityConfiguration`, `IssueTypeMissingPerimeterConfiguration`, `IssueTypeUnknown`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAudit`, `ResourceAssociationAccessModeEnforced`, `ResourceAssociationAccessModeLearning`
- New enum type `SearchBypass` with values `SearchBypassAzureServices`, `SearchBypassNone`
- New enum type `SearchDataExfiltrationProtection` with values `SearchDataExfiltrationProtectionBlockAll`
- New enum type `Severity` with values `SeverityError`, `SeverityWarning`
- New enum type `UpgradeAvailable` with values `UpgradeAvailableAvailable`, `UpgradeAvailableNotAvailable`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ServicesClient.BeginUpgrade(context.Context, string, string, *ServicesClientBeginUpgradeOptions) (*runtime.Poller[ServicesClientUpgradeResponse], error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Get(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetOptions) (NetworkSecurityPerimeterConfigurationsClientGetResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListByServicePager(string, string, *NetworkSecurityPerimeterConfigurationsClientListByServiceOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListByServiceResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileResponse], error)`
- New struct `AccessRule`
- New struct `AccessRuleProperties`
- New struct `AccessRulePropertiesSubscriptionsItem`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityProfile`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `SystemData`
- New struct `UserAssignedIdentity`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `Bypass` in struct `NetworkRuleSet`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `SystemData` in struct `Service`
- New field `ComputeType`, `DataExfiltrationProtections`, `ETag`, `Endpoint`, `ServiceUpgradedAt`, `UpgradeAvailable` in struct `ServiceProperties`
- New field `SystemData` in struct `ServiceUpdate`
- New field `SystemData` in struct `SharedPrivateLinkResource`


## 1.4.0-beta.3 (2025-04-07)
### Features Added

- New enum type `ComputeType` with values `ComputeTypeConfidential`, `ComputeTypeDefault`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `FeatureName` with values `FeatureNameAvailabilityZones`, `FeatureNameDocumentIntelligence`, `FeatureNameGrok`, `FeatureNameImageVectorization`, `FeatureNameMegaStore`, `FeatureNameQueryRewrite`, `FeatureNameS3`, `FeatureNameSemanticSearch`, `FeatureNameStorageOptimized`
- New function `*ClientFactory.NewOfferingsClient() *OfferingsClient`
- New function `NewOfferingsClient(azcore.TokenCredential, *arm.ClientOptions) (*OfferingsClient, error)`
- New function `*OfferingsClient.NewListPager(*OfferingsClientListOptions) *runtime.Pager[OfferingsClientListResponse]`
- New function `*ServicesClient.BeginUpgrade(context.Context, string, string, *ServicesClientBeginUpgradeOptions) (*runtime.Poller[ServicesClientUpgradeResponse], error)`
- New struct `FeatureOffering`
- New struct `OfferingsByRegion`
- New struct `OfferingsListResult`
- New struct `SKUOffering`
- New struct `SKUOfferingLimits`
- New struct `SystemData`
- New field `SystemData` in struct `Service`
- New field `ComputeType`, `Endpoint`, `ServiceUpgradeDate`, `UpgradeAvailable` in struct `ServiceProperties`
- New field `SystemData` in struct `ServiceUpdate`


## 1.4.0-beta.2 (2024-06-21)
### Features Added

- New value `SearchBypassAzureServices` added to enum type `SearchBypass`


## 1.4.0-beta.1 (2024-03-22)
### Features Added

- New value `IdentityTypeSystemAssignedUserAssigned`, `IdentityTypeUserAssigned` added to enum type `IdentityType`
- New value `SearchServiceStatusStopped` added to enum type `SearchServiceStatus`
- New enum type `SearchBypass` with values `SearchBypassAzurePortal`, `SearchBypassNone`
- New enum type `SearchDisabledDataExfiltrationOption` with values `SearchDisabledDataExfiltrationOptionAll`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Get(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetOptions) (NetworkSecurityPerimeterConfigurationsClientGetResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListByServicePager(string, string, *NetworkSecurityPerimeterConfigurationsClientListByServiceOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListByServiceResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileResponse], error)`
- New struct `NSPConfigAccessRule`
- New struct `NSPConfigAccessRuleProperties`
- New struct `NSPConfigAssociation`
- New struct `NSPConfigNetworkSecurityPerimeterRule`
- New struct `NSPConfigPerimeter`
- New struct `NSPConfigProfile`
- New struct `NSPProvisioningIssue`
- New struct `NSPProvisioningIssueProperties`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `OperationAvailability`
- New struct `OperationLogsSpecification`
- New struct `OperationMetricDimension`
- New struct `OperationMetricsSpecification`
- New struct `OperationProperties`
- New struct `OperationServiceSpecification`
- New struct `UserAssignedManagedIdentity`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `Bypass` in struct `NetworkRuleSet`
- New field `IsDataAction`, `Origin`, `Properties` in struct `Operation`
- New field `DisabledDataExfiltrationOptions`, `ETag` in struct `ServiceProperties`


## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-10-27)
### Features Added

- New enum type `AADAuthFailureMode` with values `AADAuthFailureModeHttp401WithBearerChallenge`, `AADAuthFailureModeHttp403`
- New enum type `PrivateLinkServiceConnectionProvisioningState` with values `PrivateLinkServiceConnectionProvisioningStateCanceled`, `PrivateLinkServiceConnectionProvisioningStateDeleting`, `PrivateLinkServiceConnectionProvisioningStateFailed`, `PrivateLinkServiceConnectionProvisioningStateIncomplete`, `PrivateLinkServiceConnectionProvisioningStateSucceeded`, `PrivateLinkServiceConnectionProvisioningStateUpdating`
- New enum type `SearchEncryptionComplianceStatus` with values `SearchEncryptionComplianceStatusCompliant`, `SearchEncryptionComplianceStatusNonCompliant`
- New enum type `SearchEncryptionWithCmk` with values `SearchEncryptionWithCmkDisabled`, `SearchEncryptionWithCmkEnabled`, `SearchEncryptionWithCmkUnspecified`
- New enum type `SearchSemanticSearch` with values `SearchSemanticSearchDisabled`, `SearchSemanticSearchFree`, `SearchSemanticSearchStandard`
- New function `*ClientFactory.NewManagementClient() *ManagementClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New function `NewManagementClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagementClient, error)`
- New function `*ManagementClient.UsageBySubscriptionSKU(context.Context, string, string, *SearchManagementRequestOptions, *ManagementClientUsageBySubscriptionSKUOptions) (ManagementClientUsageBySubscriptionSKUResponse, error)`
- New function `NewUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsagesClient, error)`
- New function `*UsagesClient.NewListBySubscriptionPager(string, *SearchManagementRequestOptions, *UsagesClientListBySubscriptionOptions) *runtime.Pager[UsagesClientListBySubscriptionResponse]`
- New struct `DataPlaneAADOrAPIKeyAuthOption`
- New struct `DataPlaneAuthOptions`
- New struct `EncryptionWithCmk`
- New struct `QuotaUsageResult`
- New struct `QuotaUsageResultName`
- New struct `QuotaUsagesListResult`
- New field `GroupID`, `ProvisioningState` in struct `PrivateEndpointConnectionProperties`
- New field `AuthOptions`, `DisableLocalAuth`, `EncryptionWithCmk`, `SemanticSearch` in struct `ServiceProperties`


## 1.1.0 (2023-04-03)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).