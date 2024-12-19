# Release History

## 1.7.0 (2024-12-27)
### Features Added

- New value `ModelLifecycleStatusDeprecated`, `ModelLifecycleStatusDeprecating`, `ModelLifecycleStatusStable` added to enum type `ModelLifecycleStatus`
- New enum type `ByPassSelection` with values `ByPassSelectionAzureServices`, `ByPassSelectionNone`
- New enum type `ContentLevel` with values `ContentLevelHigh`, `ContentLevelLow`, `ContentLevelMedium`
- New enum type `DefenderForAISettingState` with values `DefenderForAISettingStateDisabled`, `DefenderForAISettingStateEnabled`
- New enum type `EncryptionScopeProvisioningState` with values `EncryptionScopeProvisioningStateAccepted`, `EncryptionScopeProvisioningStateCanceled`, `EncryptionScopeProvisioningStateCreating`, `EncryptionScopeProvisioningStateDeleting`, `EncryptionScopeProvisioningStateFailed`, `EncryptionScopeProvisioningStateMoving`, `EncryptionScopeProvisioningStateSucceeded`
- New enum type `EncryptionScopeState` with values `EncryptionScopeStateDisabled`, `EncryptionScopeStateEnabled`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `RaiPolicyContentSource` with values `RaiPolicyContentSourceCompletion`, `RaiPolicyContentSourcePrompt`
- New enum type `RaiPolicyMode` with values `RaiPolicyModeAsynchronousFilter`, `RaiPolicyModeBlocking`, `RaiPolicyModeDefault`, `RaiPolicyModeDeferred`
- New enum type `RaiPolicyType` with values `RaiPolicyTypeSystemManaged`, `RaiPolicyTypeUserManaged`
- New function `*ClientFactory.NewDefenderForAISettingsClient() *DefenderForAISettingsClient`
- New function `*ClientFactory.NewEncryptionScopesClient() *EncryptionScopesClient`
- New function `*ClientFactory.NewLocationBasedModelCapacitiesClient() *LocationBasedModelCapacitiesClient`
- New function `*ClientFactory.NewModelCapacitiesClient() *ModelCapacitiesClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ClientFactory.NewRaiBlocklistItemsClient() *RaiBlocklistItemsClient`
- New function `*ClientFactory.NewRaiBlocklistsClient() *RaiBlocklistsClient`
- New function `*ClientFactory.NewRaiContentFiltersClient() *RaiContentFiltersClient`
- New function `*ClientFactory.NewRaiPoliciesClient() *RaiPoliciesClient`
- New function `NewDefenderForAISettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DefenderForAISettingsClient, error)`
- New function `*DefenderForAISettingsClient.CreateOrUpdate(context.Context, string, string, string, DefenderForAISetting, *DefenderForAISettingsClientCreateOrUpdateOptions) (DefenderForAISettingsClientCreateOrUpdateResponse, error)`
- New function `*DefenderForAISettingsClient.Get(context.Context, string, string, string, *DefenderForAISettingsClientGetOptions) (DefenderForAISettingsClientGetResponse, error)`
- New function `*DefenderForAISettingsClient.NewListPager(string, string, *DefenderForAISettingsClientListOptions) *runtime.Pager[DefenderForAISettingsClientListResponse]`
- New function `*DefenderForAISettingsClient.Update(context.Context, string, string, string, DefenderForAISetting, *DefenderForAISettingsClientUpdateOptions) (DefenderForAISettingsClientUpdateResponse, error)`
- New function `*DeploymentsClient.NewListSKUsPager(string, string, string, *DeploymentsClientListSKUsOptions) *runtime.Pager[DeploymentsClientListSKUsResponse]`
- New function `*DeploymentsClient.BeginUpdate(context.Context, string, string, string, PatchResourceTagsAndSKU, *DeploymentsClientBeginUpdateOptions) (*runtime.Poller[DeploymentsClientUpdateResponse], error)`
- New function `NewEncryptionScopesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EncryptionScopesClient, error)`
- New function `*EncryptionScopesClient.CreateOrUpdate(context.Context, string, string, string, EncryptionScope, *EncryptionScopesClientCreateOrUpdateOptions) (EncryptionScopesClientCreateOrUpdateResponse, error)`
- New function `*EncryptionScopesClient.BeginDelete(context.Context, string, string, string, *EncryptionScopesClientBeginDeleteOptions) (*runtime.Poller[EncryptionScopesClientDeleteResponse], error)`
- New function `*EncryptionScopesClient.Get(context.Context, string, string, string, *EncryptionScopesClientGetOptions) (EncryptionScopesClientGetResponse, error)`
- New function `*EncryptionScopesClient.NewListPager(string, string, *EncryptionScopesClientListOptions) *runtime.Pager[EncryptionScopesClientListResponse]`
- New function `NewLocationBasedModelCapacitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LocationBasedModelCapacitiesClient, error)`
- New function `*LocationBasedModelCapacitiesClient.NewListPager(string, string, string, string, *LocationBasedModelCapacitiesClientListOptions) *runtime.Pager[LocationBasedModelCapacitiesClientListResponse]`
- New function `*ManagementClient.CalculateModelCapacity(context.Context, CalculateModelCapacityParameter, *ManagementClientCalculateModelCapacityOptions) (ManagementClientCalculateModelCapacityResponse, error)`
- New function `NewModelCapacitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ModelCapacitiesClient, error)`
- New function `*ModelCapacitiesClient.NewListPager(string, string, string, *ModelCapacitiesClientListOptions) *runtime.Pager[ModelCapacitiesClientListResponse]`
- New function `NewRaiBlocklistItemsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiBlocklistItemsClient, error)`
- New function `*RaiBlocklistItemsClient.BatchAdd(context.Context, string, string, string, []*RaiBlocklistItemBulkRequest, *RaiBlocklistItemsClientBatchAddOptions) (RaiBlocklistItemsClientBatchAddResponse, error)`
- New function `*RaiBlocklistItemsClient.BatchDelete(context.Context, string, string, string, any, *RaiBlocklistItemsClientBatchDeleteOptions) (RaiBlocklistItemsClientBatchDeleteResponse, error)`
- New function `*RaiBlocklistItemsClient.CreateOrUpdate(context.Context, string, string, string, string, RaiBlocklistItem, *RaiBlocklistItemsClientCreateOrUpdateOptions) (RaiBlocklistItemsClientCreateOrUpdateResponse, error)`
- New function `*RaiBlocklistItemsClient.BeginDelete(context.Context, string, string, string, string, *RaiBlocklistItemsClientBeginDeleteOptions) (*runtime.Poller[RaiBlocklistItemsClientDeleteResponse], error)`
- New function `*RaiBlocklistItemsClient.Get(context.Context, string, string, string, string, *RaiBlocklistItemsClientGetOptions) (RaiBlocklistItemsClientGetResponse, error)`
- New function `*RaiBlocklistItemsClient.NewListPager(string, string, string, *RaiBlocklistItemsClientListOptions) *runtime.Pager[RaiBlocklistItemsClientListResponse]`
- New function `NewRaiBlocklistsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiBlocklistsClient, error)`
- New function `*RaiBlocklistsClient.CreateOrUpdate(context.Context, string, string, string, RaiBlocklist, *RaiBlocklistsClientCreateOrUpdateOptions) (RaiBlocklistsClientCreateOrUpdateResponse, error)`
- New function `*RaiBlocklistsClient.BeginDelete(context.Context, string, string, string, *RaiBlocklistsClientBeginDeleteOptions) (*runtime.Poller[RaiBlocklistsClientDeleteResponse], error)`
- New function `*RaiBlocklistsClient.Get(context.Context, string, string, string, *RaiBlocklistsClientGetOptions) (RaiBlocklistsClientGetResponse, error)`
- New function `*RaiBlocklistsClient.NewListPager(string, string, *RaiBlocklistsClientListOptions) *runtime.Pager[RaiBlocklistsClientListResponse]`
- New function `NewRaiContentFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiContentFiltersClient, error)`
- New function `*RaiContentFiltersClient.Get(context.Context, string, string, *RaiContentFiltersClientGetOptions) (RaiContentFiltersClientGetResponse, error)`
- New function `*RaiContentFiltersClient.NewListPager(string, *RaiContentFiltersClientListOptions) *runtime.Pager[RaiContentFiltersClientListResponse]`
- New function `NewRaiPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiPoliciesClient, error)`
- New function `*RaiPoliciesClient.CreateOrUpdate(context.Context, string, string, string, RaiPolicy, *RaiPoliciesClientCreateOrUpdateOptions) (RaiPoliciesClientCreateOrUpdateResponse, error)`
- New function `*RaiPoliciesClient.BeginDelete(context.Context, string, string, string, *RaiPoliciesClientBeginDeleteOptions) (*runtime.Poller[RaiPoliciesClientDeleteResponse], error)`
- New function `*RaiPoliciesClient.Get(context.Context, string, string, string, *RaiPoliciesClientGetOptions) (RaiPoliciesClientGetResponse, error)`
- New function `*RaiPoliciesClient.NewListPager(string, string, *RaiPoliciesClientListOptions) *runtime.Pager[RaiPoliciesClientListResponse]`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Get(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetOptions) (NetworkSecurityPerimeterConfigurationsClientGetResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListPager(string, string, *NetworkSecurityPerimeterConfigurationsClientListOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileResponse], error)`
- New struct `BillingMeterInfo`
- New struct `CalculateModelCapacityParameter`
- New struct `CalculateModelCapacityResult`
- New struct `CalculateModelCapacityResultEstimatedCapacity`
- New struct `CustomBlocklistConfig`
- New struct `DefenderForAISetting`
- New struct `DefenderForAISettingProperties`
- New struct `DefenderForAISettingResult`
- New struct `DeploymentCapacitySettings`
- New struct `DeploymentSKUListResult`
- New struct `EncryptionScope`
- New struct `EncryptionScopeListResult`
- New struct `EncryptionScopeProperties`
- New struct `ModelCapacityCalculatorWorkload`
- New struct `ModelCapacityCalculatorWorkloadRequestParam`
- New struct `ModelCapacityListResult`
- New struct `ModelCapacityListResultValueItem`
- New struct `ModelSKUCapacityProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterAccessRule`
- New struct `NetworkSecurityPerimeterAccessRuleProperties`
- New struct `NetworkSecurityPerimeterAccessRulePropertiesSubscriptionsItem`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationAssociationInfo`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterProfileInfo`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `RaiBlockListItemsResult`
- New struct `RaiBlockListResult`
- New struct `RaiBlocklist`
- New struct `RaiBlocklistConfig`
- New struct `RaiBlocklistItem`
- New struct `RaiBlocklistItemBulkRequest`
- New struct `RaiBlocklistItemProperties`
- New struct `RaiBlocklistProperties`
- New struct `RaiContentFilter`
- New struct `RaiContentFilterListResult`
- New struct `RaiContentFilterProperties`
- New struct `RaiMonitorConfig`
- New struct `RaiPolicy`
- New struct `RaiPolicyContentFilter`
- New struct `RaiPolicyListResult`
- New struct `RaiPolicyProperties`
- New struct `SKUResource`
- New struct `UserOwnedAmlWorkspace`
- New field `Publisher`, `SourceAccount` in struct `AccountModel`
- New field `AmlWorkspace`, `RaiMonitorConfig` in struct `AccountProperties`
- New field `AllowedValues` in struct `CapacityConfig`
- New field `Tags` in struct `CommitmentPlanAccountAssociation`
- New field `Tags` in struct `Deployment`
- New field `Publisher`, `SourceAccount` in struct `DeploymentModel`
- New field `CapacitySettings`, `CurrentCapacity`, `DynamicThrottlingEnabled`, `ParentDeploymentName` in struct `DeploymentProperties`
- New field `Description` in struct `Model`
- New field `Cost` in struct `ModelSKU`
- New field `Bypass` in struct `NetworkRuleSet`


## 1.6.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.5.0 (2023-07-28)
### Features Added

- New value `DeploymentProvisioningStateCanceled`, `DeploymentProvisioningStateDisabled` added to enum type `DeploymentProvisioningState`
- New value `HostingModelProvisionedWeb` added to enum type `HostingModel`
- New enum type `AbusePenaltyAction` with values `AbusePenaltyActionBlock`, `AbusePenaltyActionThrottle`
- New enum type `DeploymentModelVersionUpgradeOption` with values `DeploymentModelVersionUpgradeOptionNoAutoUpgrade`, `DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired`, `DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable`
- New function `*ClientFactory.NewModelsClient() *ModelsClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New function `NewModelsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ModelsClient, error)`
- New function `*ModelsClient.NewListPager(string, *ModelsClientListOptions) *runtime.Pager[ModelsClientListResponse]`
- New function `NewUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsagesClient, error)`
- New function `*UsagesClient.NewListPager(string, *UsagesClientListOptions) *runtime.Pager[UsagesClientListResponse]`
- New struct `AbusePenalty`
- New struct `CapacityConfig`
- New struct `Model`
- New struct `ModelListResult`
- New struct `ModelSKU`
- New field `IsDefaultVersion`, `SKUs`, `Source` in struct `AccountModel`
- New field `AbusePenalty` in struct `AccountProperties`
- New field `ProvisioningIssues` in struct `CommitmentPlanProperties`
- New field `SKU` in struct `Deployment`
- New field `Source` in struct `DeploymentModel`
- New field `RateLimits`, `VersionUpgradeOption` in struct `DeploymentProperties`
- New field `NextLink` in struct `UsageListResult`


## 1.4.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.4.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.3.0 (2023-02-24)
### Features Added

- New type alias `CommitmentPlanProvisioningState` with values `CommitmentPlanProvisioningStateAccepted`, `CommitmentPlanProvisioningStateCanceled`, `CommitmentPlanProvisioningStateCreating`, `CommitmentPlanProvisioningStateDeleting`, `CommitmentPlanProvisioningStateFailed`, `CommitmentPlanProvisioningStateMoving`, `CommitmentPlanProvisioningStateSucceeded`
- New type alias `ModelLifecycleStatus` with values `ModelLifecycleStatusGenerallyAvailable`, `ModelLifecycleStatusPreview`
- New type alias `RoutingMethods` with values `RoutingMethodsPerformance`, `RoutingMethodsPriority`, `RoutingMethodsWeighted`
- New function `*CommitmentPlansClient.BeginCreateOrUpdateAssociation(context.Context, string, string, string, CommitmentPlanAccountAssociation, *CommitmentPlansClientBeginCreateOrUpdateAssociationOptions) (*runtime.Poller[CommitmentPlansClientCreateOrUpdateAssociationResponse], error)`
- New function `*CommitmentPlansClient.BeginCreateOrUpdatePlan(context.Context, string, string, CommitmentPlan, *CommitmentPlansClientBeginCreateOrUpdatePlanOptions) (*runtime.Poller[CommitmentPlansClientCreateOrUpdatePlanResponse], error)`
- New function `*CommitmentPlansClient.BeginDeleteAssociation(context.Context, string, string, string, *CommitmentPlansClientBeginDeleteAssociationOptions) (*runtime.Poller[CommitmentPlansClientDeleteAssociationResponse], error)`
- New function `*CommitmentPlansClient.BeginDeletePlan(context.Context, string, string, *CommitmentPlansClientBeginDeletePlanOptions) (*runtime.Poller[CommitmentPlansClientDeletePlanResponse], error)`
- New function `*CommitmentPlansClient.GetAssociation(context.Context, string, string, string, *CommitmentPlansClientGetAssociationOptions) (CommitmentPlansClientGetAssociationResponse, error)`
- New function `*CommitmentPlansClient.GetPlan(context.Context, string, string, *CommitmentPlansClientGetPlanOptions) (CommitmentPlansClientGetPlanResponse, error)`
- New function `*CommitmentPlansClient.NewListAssociationsPager(string, string, *CommitmentPlansClientListAssociationsOptions) *runtime.Pager[CommitmentPlansClientListAssociationsResponse]`
- New function `*CommitmentPlansClient.NewListPlansByResourceGroupPager(string, *CommitmentPlansClientListPlansByResourceGroupOptions) *runtime.Pager[CommitmentPlansClientListPlansByResourceGroupResponse]`
- New function `*CommitmentPlansClient.NewListPlansBySubscriptionPager(*CommitmentPlansClientListPlansBySubscriptionOptions) *runtime.Pager[CommitmentPlansClientListPlansBySubscriptionResponse]`
- New function `*CommitmentPlansClient.BeginUpdatePlan(context.Context, string, string, PatchResourceTagsAndSKU, *CommitmentPlansClientBeginUpdatePlanOptions) (*runtime.Poller[CommitmentPlansClientUpdatePlanResponse], error)`
- New struct `CommitmentPlanAccountAssociation`
- New struct `CommitmentPlanAccountAssociationListResult`
- New struct `CommitmentPlanAccountAssociationProperties`
- New struct `CommitmentPlanAssociation`
- New struct `CommitmentPlansClientCreateOrUpdateAssociationResponse`
- New struct `CommitmentPlansClientCreateOrUpdatePlanResponse`
- New struct `CommitmentPlansClientDeleteAssociationResponse`
- New struct `CommitmentPlansClientDeletePlanResponse`
- New struct `CommitmentPlansClientListAssociationsResponse`
- New struct `CommitmentPlansClientListPlansByResourceGroupResponse`
- New struct `CommitmentPlansClientListPlansBySubscriptionResponse`
- New struct `CommitmentPlansClientUpdatePlanResponse`
- New struct `MultiRegionSettings`
- New struct `PatchResourceTags`
- New struct `PatchResourceTagsAndSKU`
- New struct `RegionSetting`
- New field `FinetuneCapabilities` in struct `AccountModel`
- New field `LifecycleStatus` in struct `AccountModel`
- New field `CommitmentPlanAssociations` in struct `AccountProperties`
- New field `Locations` in struct `AccountProperties`
- New field `Kind` in struct `CommitmentPlan`
- New field `Location` in struct `CommitmentPlan`
- New field `SKU` in struct `CommitmentPlan`
- New field `Tags` in struct `CommitmentPlan`
- New field `CommitmentPlanGUID` in struct `CommitmentPlanProperties`
- New field `ProvisioningState` in struct `CommitmentPlanProperties`


## 1.2.0 (2022-10-20)
### Features Added

- New field `CallRateLimit` in struct `DeploymentProperties`
- New field `Capabilities` in struct `DeploymentProperties`
- New field `RaiPolicyName` in struct `DeploymentProperties`
- New field `CallRateLimit` in struct `AccountModel`
- New field `CallRateLimit` in struct `DeploymentModel`


## 1.1.0 (2022-06-09)
### Features Added

- New const `DeploymentScaleTypeStandard`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).