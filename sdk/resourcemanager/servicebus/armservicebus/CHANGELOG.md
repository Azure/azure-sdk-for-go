# Release History

## 2.0.0-beta.4 (2025-10-31)
### Breaking Changes

- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ResourceNamespacePatch` has been removed
- Struct `SQLRuleAction` has been removed
- Struct `TrackedResource` has been removed

### Features Added

- New enum type `GeoDRRoleType` with values `GeoDRRoleTypePrimary`, `GeoDRRoleTypeSecondary`
- New enum type `Mode` with values `ModeDisabled`, `ModeEnabled`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleted`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateInvalidResponse`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceededWithIssues`, `NetworkSecurityPerimeterConfigurationProvisioningStateUnknown`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAuditMode`, `ResourceAssociationAccessModeEnforcedMode`, `ResourceAssociationAccessModeLearningMode`, `ResourceAssociationAccessModeNoAssociationMode`, `ResourceAssociationAccessModeUnspecifiedMode`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationClient() *NetworkSecurityPerimeterConfigurationClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*NamespacesClient.BeginFailover(context.Context, string, string, FailOver, *NamespacesClientBeginFailoverOptions) (*runtime.Poller[NamespacesClientFailoverResponse], error)`
- New function `NewNetworkSecurityPerimeterConfigurationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationClient.NewListPager(string, string, *NetworkSecurityPerimeterConfigurationClientListOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationClientListResponse]`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.GetResourceAssociationName(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetResourceAssociationNameOptions) (NetworkSecurityPerimeterConfigurationsClientGetResourceAssociationNameResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Reconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientReconcileOptions) (NetworkSecurityPerimeterConfigurationsClientReconcileResponse, error)`
- New struct `ConfidentialCompute`
- New struct `FailOver`
- New struct `FailOverProperties`
- New struct `GeoDataReplicationProperties`
- New struct `NamespaceReplicaLocation`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesProfile`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `PlatformCapabilities`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New field `GeoDataReplication`, `PlatformCapabilities` in struct `SBNamespaceProperties`
- New field `SystemData` in struct `SBNamespaceUpdateParameters`
- New field `UserMetadata` in struct `SBQueueProperties`
- New field `UserMetadata` in struct `SBSubscriptionProperties`
- New field `UserMetadata` in struct `SBTopicProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).