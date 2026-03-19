# Release History

## 2.0.0 (2026-03-19)
### Breaking Changes

- `SKUNamePremium` from enum `SKUName` has been removed
- `SKUTierPremium` from enum `SKUTier` has been removed
- Enum `ApplicationGroupPolicyType` has been removed
- Enum `CaptureIdentityType` has been removed
- Enum `CleanupPolicyRetentionDescription` has been removed
- Enum `ClusterSKUName` has been removed
- Enum `CreatedByType` has been removed
- Enum `EndPointProvisioningState` has been removed
- Enum `ManagedServiceIdentityType` has been removed
- Enum `MetricID` has been removed
- Enum `NetworkSecurityPerimeterConfigurationProvisioningState` has been removed
- Enum `NspAccessRuleDirection` has been removed
- Enum `PrivateLinkConnectionStatus` has been removed
- Enum `ProvisioningState` has been removed
- Enum `PublicNetworkAccess` has been removed
- Enum `PublicNetworkAccessFlag` has been removed
- Enum `ResourceAssociationAccessMode` has been removed
- Enum `SchemaCompatibility` has been removed
- Enum `SchemaType` has been removed
- Enum `TLSVersion` has been removed
- Function `NewApplicationGroupClient` has been removed
- Function `*ApplicationGroupClient.CreateOrUpdateApplicationGroup` has been removed
- Function `*ApplicationGroupClient.Delete` has been removed
- Function `*ApplicationGroupClient.Get` has been removed
- Function `*ApplicationGroupClient.NewListByNamespacePager` has been removed
- Function `*ApplicationGroupPolicy.GetApplicationGroupPolicy` has been removed
- Function `*ClientFactory.NewApplicationGroupClient` has been removed
- Function `*ClientFactory.NewClustersClient` has been removed
- Function `*ClientFactory.NewConfigurationClient` has been removed
- Function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationClient` has been removed
- Function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient` has been removed
- Function `*ClientFactory.NewPrivateEndpointConnectionsClient` has been removed
- Function `*ClientFactory.NewPrivateLinkResourcesClient` has been removed
- Function `*ClientFactory.NewSchemaRegistryClient` has been removed
- Function `NewClustersClient` has been removed
- Function `*ClustersClient.BeginCreateOrUpdate` has been removed
- Function `*ClustersClient.BeginDelete` has been removed
- Function `*ClustersClient.Get` has been removed
- Function `*ClustersClient.ListAvailableClusterRegion` has been removed
- Function `*ClustersClient.NewListByResourceGroupPager` has been removed
- Function `*ClustersClient.NewListBySubscriptionPager` has been removed
- Function `*ClustersClient.ListNamespaces` has been removed
- Function `*ClustersClient.BeginUpdate` has been removed
- Function `NewConfigurationClient` has been removed
- Function `*ConfigurationClient.Get` has been removed
- Function `*ConfigurationClient.Patch` has been removed
- Function `NewPrivateEndpointConnectionsClient` has been removed
- Function `*PrivateEndpointConnectionsClient.CreateOrUpdate` has been removed
- Function `*PrivateEndpointConnectionsClient.BeginDelete` has been removed
- Function `*PrivateEndpointConnectionsClient.Get` has been removed
- Function `*PrivateEndpointConnectionsClient.NewListPager` has been removed
- Function `NewPrivateLinkResourcesClient` has been removed
- Function `*PrivateLinkResourcesClient.Get` has been removed
- Function `NewSchemaRegistryClient` has been removed
- Function `*SchemaRegistryClient.CreateOrUpdate` has been removed
- Function `*SchemaRegistryClient.Delete` has been removed
- Function `*SchemaRegistryClient.Get` has been removed
- Function `*SchemaRegistryClient.NewListByNamespacePager` has been removed
- Function `*ThrottlingPolicy.GetApplicationGroupPolicy` has been removed
- Function `*NamespacesClient.ListNetworkRuleSet` has been removed
- Function `NewNetworkSecurityPerimeterConfigurationClient` has been removed
- Function `*NetworkSecurityPerimeterConfigurationClient.List` has been removed
- Function `NewNetworkSecurityPerimeterConfigurationsClient` has been removed
- Function `*NetworkSecurityPerimeterConfigurationsClient.BeginCreateOrUpdate` has been removed
- Struct `ApplicationGroup` has been removed
- Struct `ApplicationGroupListResult` has been removed
- Struct `ApplicationGroupProperties` has been removed
- Struct `AvailableCluster` has been removed
- Struct `AvailableClustersList` has been removed
- Struct `CaptureIdentity` has been removed
- Struct `Cluster` has been removed
- Struct `ClusterListResult` has been removed
- Struct `ClusterProperties` has been removed
- Struct `ClusterQuotaConfigurationProperties` has been removed
- Struct `ClusterSKU` has been removed
- Struct `ConnectionState` has been removed
- Struct `EHNamespaceIDContainer` has been removed
- Struct `EHNamespaceIDListResult` has been removed
- Struct `Encryption` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `Identity` has been removed
- Struct `KeyVaultProperties` has been removed
- Struct `NetworkSecurityPerimeter` has been removed
- Struct `NetworkSecurityPerimeterConfiguration` has been removed
- Struct `NetworkSecurityPerimeterConfigurationList` has been removed
- Struct `NetworkSecurityPerimeterConfigurationProperties` has been removed
- Struct `NetworkSecurityPerimeterConfigurationPropertiesProfile` has been removed
- Struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation` has been removed
- Struct `NspAccessRule` has been removed
- Struct `NspAccessRuleProperties` has been removed
- Struct `NspAccessRulePropertiesSubscriptionsItem` has been removed
- Struct `PrivateEndpoint` has been removed
- Struct `PrivateEndpointConnection` has been removed
- Struct `PrivateEndpointConnectionListResult` has been removed
- Struct `PrivateEndpointConnectionProperties` has been removed
- Struct `PrivateLinkResource` has been removed
- Struct `PrivateLinkResourceProperties` has been removed
- Struct `PrivateLinkResourcesListResult` has been removed
- Struct `ProvisioningIssue` has been removed
- Struct `ProvisioningIssueProperties` has been removed
- Struct `ProxyResource` has been removed
- Struct `RetentionDescription` has been removed
- Struct `SchemaGroup` has been removed
- Struct `SchemaGroupListResult` has been removed
- Struct `SchemaGroupProperties` has been removed
- Struct `SystemData` has been removed
- Struct `ThrottlingPolicy` has been removed
- Struct `UserAssignedIdentity` has been removed
- Struct `UserAssignedIdentityProperties` has been removed
- Field `Location`, `SystemData` of struct `ArmDisasterRecovery` has been removed
- Field `Location`, `SystemData` of struct `AuthorizationRule` has been removed
- Field `Location`, `SystemData` of struct `ConsumerGroup` has been removed
- Field `Identity` of struct `Destination` has been removed
- Field `DataLakeAccountName`, `DataLakeFolderPath`, `DataLakeSubscriptionID` of struct `DestinationProperties` has been removed
- Field `Identity`, `SystemData` of struct `EHNamespace` has been removed
- Field `AlternateName`, `ClusterArmID`, `DisableLocalAuth`, `Encryption`, `MinimumTLSVersion`, `PrivateEndpointConnections`, `PublicNetworkAccess`, `Status`, `ZoneRedundant` of struct `EHNamespaceProperties` has been removed
- Field `Error` of struct `ErrorResponse` has been removed
- Field `Location`, `SystemData` of struct `Eventhub` has been removed
- Field `Location`, `SystemData` of struct `NetworkRuleSet` has been removed
- Field `PublicNetworkAccess`, `TrustedServiceAccessEnabled` of struct `NetworkRuleSetProperties` has been removed
- Field `IsDataAction`, `Origin`, `Properties` of struct `Operation` has been removed
- Field `Description` of struct `OperationDisplay` has been removed
- Field `RetentionDescription`, `UserMetadata` of struct `Properties` has been removed

### Features Added

- New function `*ClientFactory.NewRegionsClient() *RegionsClient`
- New function `NewRegionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RegionsClient, error)`
- New function `*RegionsClient.NewListBySKUPager(sku string, options *RegionsClientListBySKUOptions) *runtime.Pager[RegionsClientListBySKUResponse]`
- New function `*NamespacesClient.GetMessagingPlan(ctx context.Context, resourceGroupName string, namespaceName string, options *NamespacesClientGetMessagingPlanOptions) (NamespacesClientGetMessagingPlanResponse, error)`
- New function `*NamespacesClient.NewListNetworkRuleSetsPager(resourceGroupName string, namespaceName string, options *NamespacesClientListNetworkRuleSetsOptions) *runtime.Pager[NamespacesClientListNetworkRuleSetsResponse]`
- New struct `MessagingPlan`
- New struct `MessagingPlanProperties`
- New struct `MessagingRegions`
- New struct `MessagingRegionsListResult`
- New struct `MessagingRegionsProperties`
- New field `Code`, `Message` in struct `ErrorResponse`


## 1.4.0-beta.1 (2025-01-23)
### Features Added

- New value `CleanupPolicyRetentionDescriptionDeleteOrCompact` added to enum type `CleanupPolicyRetentionDescription`
- New value `SchemaTypeJSON`, `SchemaTypeProtoBuf` added to enum type `SchemaType`
- New enum type `GeoDRRoleType` with values `GeoDRRoleTypePrimary`, `GeoDRRoleTypeSecondary`
- New enum type `TimestampType` with values `TimestampTypeCreate`, `TimestampTypeLogAppend`
- New function `*NamespacesClient.BeginFailover(context.Context, string, string, FailOver, *NamespacesClientBeginFailoverOptions) (*runtime.Poller[NamespacesClientFailoverResponse], error)`
- New struct `ErrorDetailAutoGenerated`
- New struct `ErrorResponseAutoGenerated`
- New struct `FailOver`
- New struct `FailOverProperties`
- New struct `GeoDataReplicationProperties`
- New struct `MessageTimestampDescription`
- New struct `NamespaceReplicaLocation`
- New field `GeoDataReplication` in struct `EHNamespaceProperties`
- New field `Identifier`, `MessageTimestampDescription` in struct `Properties`
- New field `MinCompactionLagInMins` in struct `RetentionDescription`


## 1.3.0 (2024-07-24)
### Features Added

- New value `PublicNetworkAccessFlagSecuredByPerimeter` added to enum type `PublicNetworkAccessFlag`
- New enum type `ApplicationGroupPolicyType` with values `ApplicationGroupPolicyTypeThrottlingPolicy`
- New enum type `CaptureIdentityType` with values `CaptureIdentityTypeSystemAssigned`, `CaptureIdentityTypeUserAssigned`
- New enum type `CleanupPolicyRetentionDescription` with values `CleanupPolicyRetentionDescriptionCompact`, `CleanupPolicyRetentionDescriptionDelete`
- New enum type `MetricID` with values `MetricIDIncomingBytes`, `MetricIDIncomingMessages`, `MetricIDOutgoingBytes`, `MetricIDOutgoingMessages`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleted`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateInvalidResponse`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceededWithIssues`, `NetworkSecurityPerimeterConfigurationProvisioningStateUnknown`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `ProvisioningState` with values `ProvisioningStateActive`, `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateScaling`, `ProvisioningStateSucceeded`, `ProvisioningStateUnknown`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`, `PublicNetworkAccessSecuredByPerimeter`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAuditMode`, `ResourceAssociationAccessModeEnforcedMode`, `ResourceAssociationAccessModeLearningMode`, `ResourceAssociationAccessModeNoAssociationMode`, `ResourceAssociationAccessModeUnspecifiedMode`
- New enum type `TLSVersion` with values `TLSVersionOne0`, `TLSVersionOne1`, `TLSVersionOne2`
- New function `NewApplicationGroupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApplicationGroupClient, error)`
- New function `*ApplicationGroupClient.CreateOrUpdateApplicationGroup(context.Context, string, string, string, ApplicationGroup, *ApplicationGroupClientCreateOrUpdateApplicationGroupOptions) (ApplicationGroupClientCreateOrUpdateApplicationGroupResponse, error)`
- New function `*ApplicationGroupClient.Delete(context.Context, string, string, string, *ApplicationGroupClientDeleteOptions) (ApplicationGroupClientDeleteResponse, error)`
- New function `*ApplicationGroupClient.Get(context.Context, string, string, string, *ApplicationGroupClientGetOptions) (ApplicationGroupClientGetResponse, error)`
- New function `*ApplicationGroupClient.NewListByNamespacePager(string, string, *ApplicationGroupClientListByNamespaceOptions) *runtime.Pager[ApplicationGroupClientListByNamespaceResponse]`
- New function `*ApplicationGroupPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `*ClientFactory.NewApplicationGroupClient() *ApplicationGroupClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationClient() *NetworkSecurityPerimeterConfigurationClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ThrottlingPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `NewNetworkSecurityPerimeterConfigurationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationClient.List(context.Context, string, string, *NetworkSecurityPerimeterConfigurationClientListOptions) (NetworkSecurityPerimeterConfigurationClientListResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginCreateOrUpdate(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse], error)`
- New struct `ApplicationGroup`
- New struct `ApplicationGroupListResult`
- New struct `ApplicationGroupProperties`
- New struct `CaptureIdentity`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesProfile`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `RetentionDescription`
- New struct `ThrottlingPolicy`
- New field `ProvisioningState`, `SupportsScaling` in struct `ClusterProperties`
- New field `Identity` in struct `Destination`
- New field `MinimumTLSVersion`, `PublicNetworkAccess` in struct `EHNamespaceProperties`
- New field `RetentionDescription`, `UserMetadata` in struct `Properties`


## 1.3.0-beta.1 (2023-11-30)
### Features Added

- New value `PublicNetworkAccessFlagSecuredByPerimeter` added to enum type `PublicNetworkAccessFlag`
- New enum type `ApplicationGroupPolicyType` with values `ApplicationGroupPolicyTypeThrottlingPolicy`
- New enum type `CleanupPolicyRetentionDescription` with values `CleanupPolicyRetentionDescriptionCompaction`, `CleanupPolicyRetentionDescriptionDelete`
- New enum type `MetricID` with values `MetricIDIncomingBytes`, `MetricIDIncomingMessages`, `MetricIDOutgoingBytes`, `MetricIDOutgoingMessages`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleted`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateInvalidResponse`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceededWithIssues`, `NetworkSecurityPerimeterConfigurationProvisioningStateUnknown`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`, `PublicNetworkAccessSecuredByPerimeter`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAuditMode`, `ResourceAssociationAccessModeEnforcedMode`, `ResourceAssociationAccessModeLearningMode`, `ResourceAssociationAccessModeNoAssociationMode`, `ResourceAssociationAccessModeUnspecifiedMode`
- New enum type `TLSVersion` with values `TLSVersionOne0`, `TLSVersionOne1`, `TLSVersionOne2`
- New function `NewApplicationGroupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApplicationGroupClient, error)`
- New function `*ApplicationGroupClient.CreateOrUpdateApplicationGroup(context.Context, string, string, string, ApplicationGroup, *ApplicationGroupClientCreateOrUpdateApplicationGroupOptions) (ApplicationGroupClientCreateOrUpdateApplicationGroupResponse, error)`
- New function `*ApplicationGroupClient.Delete(context.Context, string, string, string, *ApplicationGroupClientDeleteOptions) (ApplicationGroupClientDeleteResponse, error)`
- New function `*ApplicationGroupClient.Get(context.Context, string, string, string, *ApplicationGroupClientGetOptions) (ApplicationGroupClientGetResponse, error)`
- New function `*ApplicationGroupClient.NewListByNamespacePager(string, string, *ApplicationGroupClientListByNamespaceOptions) *runtime.Pager[ApplicationGroupClientListByNamespaceResponse]`
- New function `*ApplicationGroupPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `*ClientFactory.NewApplicationGroupClient() *ApplicationGroupClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationClient() *NetworkSecurityPerimeterConfigurationClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ThrottlingPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `NewNetworkSecurityPerimeterConfigurationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationClient.List(context.Context, string, string, *NetworkSecurityPerimeterConfigurationClientListOptions) (NetworkSecurityPerimeterConfigurationClientListResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginCreateOrUpdate(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse], error)`
- New struct `ApplicationGroup`
- New struct `ApplicationGroupListResult`
- New struct `ApplicationGroupProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesProfile`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `RetentionDescription`
- New struct `ThrottlingPolicy`
- New field `SupportsScaling` in struct `ClusterProperties`
- New field `MinimumTLSVersion`, `PublicNetworkAccess` in struct `EHNamespaceProperties`
- New field `RetentionDescription` in struct `Properties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.1 (2023-04-28)
### Features Added

- New value `PublicNetworkAccessFlagSecuredByPerimeter` added to enum type `PublicNetworkAccessFlag`
- New enum type `ApplicationGroupPolicyType` with values `ApplicationGroupPolicyTypeThrottlingPolicy`
- New enum type `CleanupPolicyRetentionDescription` with values `CleanupPolicyRetentionDescriptionCompaction`, `CleanupPolicyRetentionDescriptionDelete`
- New enum type `MetricID` with values `MetricIDIncomingBytes`, `MetricIDIncomingMessages`, `MetricIDOutgoingBytes`, `MetricIDOutgoingMessages`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleted`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateInvalidResponse`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceededWithIssues`, `NetworkSecurityPerimeterConfigurationProvisioningStateUnknown`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`, `PublicNetworkAccessSecuredByPerimeter`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAuditMode`, `ResourceAssociationAccessModeEnforcedMode`, `ResourceAssociationAccessModeLearningMode`, `ResourceAssociationAccessModeNoAssociationMode`, `ResourceAssociationAccessModeUnspecifiedMode`
- New enum type `TLSVersion` with values `TLSVersionOne0`, `TLSVersionOne1`, `TLSVersionOne2`
- New function `NewApplicationGroupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApplicationGroupClient, error)`
- New function `*ApplicationGroupClient.CreateOrUpdateApplicationGroup(context.Context, string, string, string, ApplicationGroup, *ApplicationGroupClientCreateOrUpdateApplicationGroupOptions) (ApplicationGroupClientCreateOrUpdateApplicationGroupResponse, error)`
- New function `*ApplicationGroupClient.Delete(context.Context, string, string, string, *ApplicationGroupClientDeleteOptions) (ApplicationGroupClientDeleteResponse, error)`
- New function `*ApplicationGroupClient.Get(context.Context, string, string, string, *ApplicationGroupClientGetOptions) (ApplicationGroupClientGetResponse, error)`
- New function `*ApplicationGroupClient.NewListByNamespacePager(string, string, *ApplicationGroupClientListByNamespaceOptions) *runtime.Pager[ApplicationGroupClientListByNamespaceResponse]`
- New function `*ApplicationGroupPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `*ClientFactory.NewApplicationGroupClient() *ApplicationGroupClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationClient() *NetworkSecurityPerimeterConfigurationClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ThrottlingPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `NewNetworkSecurityPerimeterConfigurationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationClient.List(context.Context, string, string, *NetworkSecurityPerimeterConfigurationClientListOptions) (NetworkSecurityPerimeterConfigurationClientListResponse, error)`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginCreateOrUpdate(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse], error)`
- New struct `ApplicationGroup`
- New struct `ApplicationGroupListResult`
- New struct `ApplicationGroupProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesProfile`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `RetentionDescription`
- New struct `ThrottlingPolicy`
- New field `SupportsScaling` in struct `ClusterProperties`
- New field `MinimumTLSVersion` in struct `EHNamespaceProperties`
- New field `PublicNetworkAccess` in struct `EHNamespaceProperties`
- New field `RetentionDescription` in struct `Properties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-06)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).