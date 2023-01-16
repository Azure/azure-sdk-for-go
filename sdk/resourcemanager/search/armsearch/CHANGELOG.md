# Release History

## 2.0.0 (2023-01-16)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Struct `ListQueryKeysResult` has been removed
- Struct `PrivateEndpointConnectionListResult` has been removed
- Struct `PrivateLinkResourcesResult` has been removed
- Struct `ServiceListResult` has been removed
- Struct `SharedPrivateLinkResourceListResult` has been removed

### Features Added

- New value `SearchServiceStatusStopped` added to type alias `SearchServiceStatus`
- New type alias `AADAuthFailureMode` with values `AADAuthFailureModeHttp401WithBearerChallenge`, `AADAuthFailureModeHttp403`
- New type alias `PrivateLinkServiceConnectionProvisioningState` with values `PrivateLinkServiceConnectionProvisioningStateCanceled`, `PrivateLinkServiceConnectionProvisioningStateDeleting`, `PrivateLinkServiceConnectionProvisioningStateFailed`, `PrivateLinkServiceConnectionProvisioningStateIncomplete`, `PrivateLinkServiceConnectionProvisioningStateSucceeded`, `PrivateLinkServiceConnectionProvisioningStateUpdating`
- New type alias `SearchEncryptionComplianceStatus` with values `SearchEncryptionComplianceStatusCompliant`, `SearchEncryptionComplianceStatusNonCompliant`
- New type alias `SearchEncryptionWithCmk` with values `SearchEncryptionWithCmkDisabled`, `SearchEncryptionWithCmkEnabled`, `SearchEncryptionWithCmkUnspecified`
- New struct `DataPlaneAADOrAPIKeyAuthOption`
- New struct `DataPlaneAuthOptions`
- New struct `EncryptionWithCmk`
- New field `GroupID` in struct `PrivateEndpointConnectionProperties`
- New field `ProvisioningState` in struct `PrivateEndpointConnectionProperties`
- New field `AuthOptions` in struct `ServiceProperties`
- New field `DisableLocalAuth` in struct `ServiceProperties`
- New field `EncryptionWithCmk` in struct `ServiceProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).