# Release History

## 1.1.0-beta.3 (2023-04-07)
### Other Changes


## 1.1.0-beta.2 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0-beta.1 (2022-05-19)
### Features Added

- New const `ResourceAssociationAccessModeLearningMode`
- New const `MetricIDIncomingMessages`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`
- New const `PublicNetworkAccessSecuredByPerimeter`
- New const `ResourceAssociationAccessModeNoAssociationMode`
- New const `ApplicationGroupPolicyTypeThrottlingPolicy`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`
- New const `PublicNetworkAccessFlagSecuredByPerimeter`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateDeleted`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateInvalidResponse`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateUnknown`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`
- New const `MetricIDOutgoingBytes`
- New const `NspAccessRuleDirectionInbound`
- New const `MetricIDIncomingBytes`
- New const `PublicNetworkAccessEnabled`
- New const `NspAccessRuleDirectionOutbound`
- New const `TLSVersionOne2`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New const `PublicNetworkAccessDisabled`
- New const `TLSVersionOne1`
- New const `ResourceAssociationAccessModeEnforcedMode`
- New const `TLSVersionOne0`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateSucceededWithIssues`
- New const `MetricIDOutgoingMessages`
- New const `ResourceAssociationAccessModeUnspecifiedMode`
- New const `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`
- New const `ResourceAssociationAccessModeAuditMode`
- New function `PossibleTLSVersionValues() []TLSVersion`
- New function `NetworkSecurityPerimeterConfigurationProperties.MarshalJSON() ([]byte, error)`
- New function `*ThrottlingPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `*ThrottlingPolicy.UnmarshalJSON([]byte) error`
- New function `PossibleMetricIDValues() []MetricID`
- New function `NetworkSecurityPerimeterConfigurationPropertiesProfile.MarshalJSON() ([]byte, error)`
- New function `ThrottlingPolicy.MarshalJSON() ([]byte, error)`
- New function `PossibleApplicationGroupPolicyTypeValues() []ApplicationGroupPolicyType`
- New function `*ApplicationGroupPolicy.GetApplicationGroupPolicy() *ApplicationGroupPolicy`
- New function `NspAccessRuleProperties.MarshalJSON() ([]byte, error)`
- New function `*ApplicationGroupProperties.UnmarshalJSON([]byte) error`
- New function `ApplicationGroupProperties.MarshalJSON() ([]byte, error)`
- New function `PossibleResourceAssociationAccessModeValues() []ResourceAssociationAccessMode`
- New function `PossibleNetworkSecurityPerimeterConfigurationProvisioningStateValues() []NetworkSecurityPerimeterConfigurationProvisioningState`
- New function `PossibleNspAccessRuleDirectionValues() []NspAccessRuleDirection`
- New function `PossiblePublicNetworkAccessValues() []PublicNetworkAccess`
- New function `NetworkSecurityPerimeterConfiguration.MarshalJSON() ([]byte, error)`
- New struct `ApplicationGroup`
- New struct `ApplicationGroupClientCreateOrUpdateApplicationGroupOptions`
- New struct `ApplicationGroupClientCreateOrUpdateApplicationGroupResponse`
- New struct `ApplicationGroupClientDeleteOptions`
- New struct `ApplicationGroupClientDeleteResponse`
- New struct `ApplicationGroupClientGetOptions`
- New struct `ApplicationGroupClientGetResponse`
- New struct `ApplicationGroupClientListByNamespaceOptions`
- New struct `ApplicationGroupClientListByNamespaceResponse`
- New struct `ApplicationGroupListResult`
- New struct `ApplicationGroupPolicy`
- New struct `ApplicationGroupProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationClientListOptions`
- New struct `NetworkSecurityPerimeterConfigurationClientListResponse`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesProfile`
- New struct `NetworkSecurityPerimeterConfigurationPropertiesResourceAssociation`
- New struct `NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions`
- New struct `NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ThrottlingPolicy`
- New field `SupportsScaling` in struct `ClusterProperties`
- New field `PublicNetworkAccess` in struct `EHNamespaceProperties`
- New field `MinimumTLSVersion` in struct `EHNamespaceProperties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).