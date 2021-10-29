# Release History

## 0.2.0 (2021-10-29)
### Breaking Changes

- Function `NewIotSecuritySolutionAnalyticsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSubAssessmentsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewJitNetworkAccessPoliciesClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAllowedConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRegulatoryComplianceStandardsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAutoProvisioningSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPricingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewLocationsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*AdvancedThreatProtectionClient.Create` parameter(s) have been changed from `(context.Context, string, Enum6, AdvancedThreatProtectionSetting, *AdvancedThreatProtectionCreateOptions)` to `(context.Context, string, Enum7, AdvancedThreatProtectionSetting, *AdvancedThreatProtectionCreateOptions)`
- Function `*AdvancedThreatProtectionClient.Get` parameter(s) have been changed from `(context.Context, string, Enum6, *AdvancedThreatProtectionGetOptions)` to `(context.Context, string, Enum7, *AdvancedThreatProtectionGetOptions)`
- Function `NewSecureScoresClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSQLVulnerabilityAssessmentScanResultsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAdvancedThreatProtectionClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSQLVulnerabilityAssessmentBaselineRulesClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSecurityContactsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAssessmentsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCustomEntityStoreAssignmentsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewWorkspaceSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewExternalSecuritySolutionsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAdaptiveNetworkHardeningsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSQLVulnerabilityAssessmentScansClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCompliancesClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotSecuritySolutionsAnalyticsAggregatedAlertClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSoftwareInventoriesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotSecuritySolutionsAnalyticsRecommendationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAdaptiveApplicationControlsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAssessmentsMetadataClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSecureScoreControlsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRegulatoryComplianceControlsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConnectorsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAlertsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDiscoveredSecuritySolutionsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewComplianceResultsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotSecuritySolutionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewInformationProtectionPoliciesClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAutomationsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAlertsSuppressionRulesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCustomAssessmentAutomationsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDeviceSecurityGroupsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSecuritySolutionsClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSecuritySolutionsReferenceDataClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSecureScoreControlDefinitionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewServerVulnerabilityAssessmentClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*CustomAssessmentAutomationsClient.Create` parameter(s) have been changed from `(context.Context, string, string, CustomAssessmentAutomation, *CustomAssessmentAutomationsCreateOptions)` to `(context.Context, string, string, CustomAssessmentAutomationRequest, *CustomAssessmentAutomationsCreateOptions)`
- Function `NewIngestionSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRegulatoryComplianceAssessmentsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTopologyClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTasksClient` parameter(s) have been changed from `(*arm.Connection, string, string)` to `(string, string, azcore.TokenCredential, *arm.ClientOptions)`
- Const `Enum6Current` has been removed
- Function `Enum6.ToPtr` has been removed
- Function `PossibleEnum6Values` has been removed

### New Content

- New const `OfferingTypeDefenderForContainersAws`
- New const `CloudNameAWS`
- New const `Enum7Current`
- New const `OrganizationMembershipTypeMember`
- New const `OfferingTypeDefenderForServersAws`
- New const `OfferingTypeCspmMonitorAws`
- New const `EventSourceSubAssessmentsSnapshot`
- New const `OrganizationMembershipTypeOrganization`
- New const `CloudNameAzure`
- New const `CloudNameGCP`
- New const `EventSourceAssessmentsSnapshot`
- New function `PossibleOfferingTypeValues() []OfferingType`
- New function `*SecurityConnectorsClient.CreateOrUpdate(context.Context, string, string, SecurityConnector, *SecurityConnectorsCreateOrUpdateOptions) (SecurityConnectorsCreateOrUpdateResponse, error)`
- New function `*DefenderForServersAwsOffering.UnmarshalJSON([]byte) error`
- New function `*SecurityConnectorsClient.Update(context.Context, string, string, SecurityConnector, *SecurityConnectorsUpdateOptions) (SecurityConnectorsUpdateResponse, error)`
- New function `*SecurityConnectorsClient.List(*SecurityConnectorsListOptions) *SecurityConnectorsListPager`
- New function `*MdeOnboardingData.UnmarshalJSON([]byte) error`
- New function `PossibleOrganizationMembershipTypeValues() []OrganizationMembershipType`
- New function `PossibleEnum7Values() []Enum7`
- New function `*SecurityConnector.UnmarshalJSON([]byte) error`
- New function `SecurityConnectorProperties.MarshalJSON() ([]byte, error)`
- New function `*CspmMonitorAwsOffering.UnmarshalJSON([]byte) error`
- New function `*SecurityConnectorsListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*MdeOnboardingsClient.Get(context.Context, *MdeOnboardingsGetOptions) (MdeOnboardingsGetResponse, error)`
- New function `*CloudOffering.UnmarshalJSON([]byte) error`
- New function `MdeOnboardingDataProperties.MarshalJSON() ([]byte, error)`
- New function `*CloudOffering.GetCloudOffering() *CloudOffering`
- New function `Enum7.ToPtr() *Enum7`
- New function `*SecurityConnectorsClient.Delete(context.Context, string, string, *SecurityConnectorsDeleteOptions) (SecurityConnectorsDeleteResponse, error)`
- New function `*SecurityConnectorsClient.Get(context.Context, string, string, *SecurityConnectorsGetOptions) (SecurityConnectorsGetResponse, error)`
- New function `CspmMonitorAwsOffering.MarshalJSON() ([]byte, error)`
- New function `*SecurityConnectorsClient.ListByResourceGroup(string, *SecurityConnectorsListByResourceGroupOptions) *SecurityConnectorsListByResourceGroupPager`
- New function `PossibleCloudNameValues() []CloudName`
- New function `*SecurityConnectorsListPager.NextPage(context.Context) bool`
- New function `*SecurityConnectorsListByResourceGroupPager.Err() error`
- New function `SecurityConnectorsList.MarshalJSON() ([]byte, error)`
- New function `*SecurityConnectorsListPager.PageResponse() SecurityConnectorsListResponse`
- New function `*CustomAssessmentAutomationRequest.UnmarshalJSON([]byte) error`
- New function `*SecurityConnectorsListByResourceGroupPager.PageResponse() SecurityConnectorsListByResourceGroupResponse`
- New function `CustomAssessmentAutomationRequest.MarshalJSON() ([]byte, error)`
- New function `MdeOnboardingDataList.MarshalJSON() ([]byte, error)`
- New function `*MdeOnboardingsClient.List(context.Context, *MdeOnboardingsListOptions) (MdeOnboardingsListResponse, error)`
- New function `*SecurityConnectorProperties.UnmarshalJSON([]byte) error`
- New function `SecurityConnectorPropertiesOrganizationalData.MarshalJSON() ([]byte, error)`
- New function `OrganizationMembershipType.ToPtr() *OrganizationMembershipType`
- New function `MdeOnboardingData.MarshalJSON() ([]byte, error)`
- New function `DefenderForServersAwsOffering.MarshalJSON() ([]byte, error)`
- New function `*DefenderForContainersAwsOffering.UnmarshalJSON([]byte) error`
- New function `*MdeOnboardingDataProperties.UnmarshalJSON([]byte) error`
- New function `*SecurityConnectorsListPager.Err() error`
- New function `OfferingType.ToPtr() *OfferingType`
- New function `CloudName.ToPtr() *CloudName`
- New function `DefenderForContainersAwsOffering.MarshalJSON() ([]byte, error)`
- New function `NewSecurityConnectorsClient(string, azcore.TokenCredential, *arm.ClientOptions) *SecurityConnectorsClient`
- New function `NewMdeOnboardingsClient(string, azcore.TokenCredential, *arm.ClientOptions) *MdeOnboardingsClient`
- New function `SecurityConnector.MarshalJSON() ([]byte, error)`
- New struct `CloudOffering`
- New struct `CspmMonitorAwsOffering`
- New struct `CspmMonitorAwsOfferingNativeCloudConnection`
- New struct `CustomAssessmentAutomationRequest`
- New struct `CustomAssessmentAutomationRequestProperties`
- New struct `DefenderForContainersAwsOffering`
- New struct `DefenderForContainersAwsOfferingCloudWatchToKinesis`
- New struct `DefenderForContainersAwsOfferingKinesisToS3`
- New struct `DefenderForContainersAwsOfferingKubernetesScubaReader`
- New struct `DefenderForContainersAwsOfferingKubernetesService`
- New struct `DefenderForServersAwsOffering`
- New struct `DefenderForServersAwsOfferingArcAutoProvisioning`
- New struct `DefenderForServersAwsOfferingArcAutoProvisioningServicePrincipalSecretMetadata`
- New struct `DefenderForServersAwsOfferingDefenderForServers`
- New struct `MdeOnboardingData`
- New struct `MdeOnboardingDataList`
- New struct `MdeOnboardingDataProperties`
- New struct `MdeOnboardingsClient`
- New struct `MdeOnboardingsGetOptions`
- New struct `MdeOnboardingsGetResponse`
- New struct `MdeOnboardingsGetResult`
- New struct `MdeOnboardingsListOptions`
- New struct `MdeOnboardingsListResponse`
- New struct `MdeOnboardingsListResult`
- New struct `SecurityConnector`
- New struct `SecurityConnectorProperties`
- New struct `SecurityConnectorPropertiesOrganizationalData`
- New struct `SecurityConnectorsClient`
- New struct `SecurityConnectorsCreateOrUpdateOptions`
- New struct `SecurityConnectorsCreateOrUpdateResponse`
- New struct `SecurityConnectorsCreateOrUpdateResult`
- New struct `SecurityConnectorsDeleteOptions`
- New struct `SecurityConnectorsDeleteResponse`
- New struct `SecurityConnectorsGetOptions`
- New struct `SecurityConnectorsGetResponse`
- New struct `SecurityConnectorsGetResult`
- New struct `SecurityConnectorsList`
- New struct `SecurityConnectorsListByResourceGroupOptions`
- New struct `SecurityConnectorsListByResourceGroupPager`
- New struct `SecurityConnectorsListByResourceGroupResponse`
- New struct `SecurityConnectorsListByResourceGroupResult`
- New struct `SecurityConnectorsListOptions`
- New struct `SecurityConnectorsListPager`
- New struct `SecurityConnectorsListResponse`
- New struct `SecurityConnectorsListResult`
- New struct `SecurityConnectorsUpdateOptions`
- New struct `SecurityConnectorsUpdateResponse`
- New struct `SecurityConnectorsUpdateResult`
- New field `SystemData` in struct `CustomEntityStoreAssignment`
- New field `AssessmentKey` in struct `CustomAssessmentAutomationProperties`
- New field `SystemData` in struct `CustomAssessmentAutomation`

Total 53 breaking change(s), 156 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-15)

- Initial preview release.
