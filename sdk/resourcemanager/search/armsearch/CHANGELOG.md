# Release History

## 1.4.0-beta.3 (2025-02-28)
### Features Added

- New enum type `ComputeType` with values `ComputeTypeConfidential`, `ComputeTypeDefault`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `FeatureName` with values `FeatureNameAvailabilityZones`, `FeatureNameDocumentIntelligence`, `FeatureNameGrok`, `FeatureNameImageVectorization`, `FeatureNameMegaStore`, `FeatureNameQueryRewrite`, `FeatureNameS3`, `FeatureNameSemanticSearch`, `FeatureNameStorageOptimized`
- New function `*ClientFactory.NewOfferingsClient() *OfferingsClient`
- New function `*ClientFactory.NewServiceClient() *ServiceClient`
- New function `NewOfferingsClient(azcore.TokenCredential, *arm.ClientOptions) (*OfferingsClient, error)`
- New function `*OfferingsClient.NewListPager(*OfferingsClientListOptions) *runtime.Pager[OfferingsClientListResponse]`
- New function `NewServiceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServiceClient, error)`
- New function `*ServiceClient.BeginUpgrade(context.Context, string, string, *ServiceClientBeginUpgradeOptions) (*runtime.Poller[ServiceClientUpgradeResponse], error)`
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