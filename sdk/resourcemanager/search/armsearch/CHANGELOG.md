# Release History

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