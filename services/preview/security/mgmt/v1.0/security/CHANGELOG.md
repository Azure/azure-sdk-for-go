# Unreleased

## Breaking Changes

### Removed Constants

1. EnforcementMode.Audit
1. EnforcementMode.Enforce
1. EventSource.Alerts
1. EventSource.Assessments
1. EventSource.SecureScoreControls
1. EventSource.SecureScores
1. EventSource.SubAssessments
1. SettingKind.SettingKindAlertSuppressionSetting
1. SettingKind.SettingKindDataExportSetting

### Removed Funcs

1. PossibleSettingKindValues() []SettingKind
1. SettingResource.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. SettingResource

#### Removed Struct Fields

1. BaseClient.AscLocation

### Signature Changes

#### Const Types

1. BinarySignature changed type from Type to Type1
1. File changed type from Type to Type1
1. FileHash changed type from Type to Type1
1. KindAAD changed type from KindEnum to KindEnum1
1. KindATA changed type from KindEnum to KindEnum1
1. KindCEF changed type from KindEnum to KindEnum1
1. KindExternalSecuritySolution changed type from KindEnum to KindEnum1
1. None changed type from EnforcementMode to EndOfSupportStatus
1. ProductSignature changed type from Type to Type1
1. PublisherSignature changed type from Type to Type1
1. VersionAndAboveSignature changed type from Type to Type1

#### Funcs

1. AdaptiveApplicationControlsClient.Delete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.DeletePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.Get
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.GetPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.Put
	- Params
		- From: context.Context, string, AppWhitelistingPutGroupData
		- To: context.Context, string, string, AppWhitelistingPutGroupData
1. AdaptiveApplicationControlsClient.PutPreparer
	- Params
		- From: context.Context, string, AppWhitelistingPutGroupData
		- To: context.Context, string, string, AppWhitelistingPutGroupData
1. AlertsClient.GetResourceGroupLevelAlerts
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.GetResourceGroupLevelAlertsPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.GetSubscriptionLevelAlert
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.GetSubscriptionLevelAlertPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.ListResourceGroupLevelAlertsByRegion
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string, string, string
1. AlertsClient.ListResourceGroupLevelAlertsByRegionComplete
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string, string, string
1. AlertsClient.ListResourceGroupLevelAlertsByRegionPreparer
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, string, string, string, string
1. AlertsClient.ListSubscriptionLevelAlertsByRegion
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, string
1. AlertsClient.ListSubscriptionLevelAlertsByRegionComplete
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, string
1. AlertsClient.ListSubscriptionLevelAlertsByRegionPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, string
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismiss
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismissPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivate
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivatePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismiss
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismissPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivate
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivatePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AllowedConnectionsClient.Get
	- Params
		- From: context.Context, string, ConnectionType
		- To: context.Context, string, string, ConnectionType
1. AllowedConnectionsClient.GetPreparer
	- Params
		- From: context.Context, string, ConnectionType
		- To: context.Context, string, string, ConnectionType
1. AllowedConnectionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. AllowedConnectionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. AllowedConnectionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. ConnectorsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, ConnectorSetting
		- To: context.Context, string, string, Connector
	- Returns
		- From: ConnectorSetting, error
		- To: Connector, error
1. ConnectorsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, ConnectorSetting
		- To: context.Context, string, string, Connector
1. ConnectorsClient.CreateOrUpdateResponder
	- Returns
		- From: ConnectorSetting, error
		- To: Connector, error
1. ConnectorsClient.Delete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. ConnectorsClient.DeletePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. ConnectorsClient.Get
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
	- Returns
		- From: ConnectorSetting, error
		- To: Connector, error
1. ConnectorsClient.GetPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. ConnectorsClient.GetResponder
	- Returns
		- From: ConnectorSetting, error
		- To: Connector, error
1. ConnectorsClient.List
	- Returns
		- From: ConnectorSettingListPage, error
		- To: ConnectorsListPage, error
1. ConnectorsClient.ListComplete
	- Returns
		- From: ConnectorSettingListIterator, error
		- To: ConnectorsListIterator, error
1. ConnectorsClient.ListResponder
	- Returns
		- From: ConnectorSettingList, error
		- To: ConnectorsList, error
1. DiscoveredSecuritySolutionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. DiscoveredSecuritySolutionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. ExternalSecuritySolutionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. ExternalSecuritySolutionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. InformationProtectionPoliciesClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, InformationProtectionPolicy
		- To: context.Context, string, InformationProtectionPolicyName, InformationProtectionPolicy
1. InformationProtectionPoliciesClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, InformationProtectionPolicy
		- To: context.Context, string, InformationProtectionPolicyName, InformationProtectionPolicy
1. InformationProtectionPoliciesClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, InformationProtectionPolicyName
1. InformationProtectionPoliciesClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, InformationProtectionPolicyName
1. JitNetworkAccessPoliciesClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicy
		- To: context.Context, string, string, string, JitNetworkAccessPolicy
1. JitNetworkAccessPoliciesClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicy
		- To: context.Context, string, string, string, JitNetworkAccessPolicy
1. JitNetworkAccessPoliciesClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.Initiate
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicyInitiateRequest
		- To: context.Context, string, string, string, JitNetworkAccessPolicyInitiateRequest
1. JitNetworkAccessPoliciesClient.InitiatePreparer
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicyInitiateRequest
		- To: context.Context, string, string, string, JitNetworkAccessPolicyInitiateRequest
1. JitNetworkAccessPoliciesClient.ListByRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegion
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegionComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. LocationsClient.Get
	- Params
		- From: context.Context
		- To: context.Context, string
1. LocationsClient.GetPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. New
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveApplicationControlsClient
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveApplicationControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAdaptiveNetworkHardeningsClient
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveNetworkHardeningsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAdvancedThreatProtectionClient
	- Params
		- From: string, string
		- To: string
1. NewAdvancedThreatProtectionClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAlertsClient
	- Params
		- From: string, string
		- To: string
1. NewAlertsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAlertsSuppressionRulesClient
	- Params
		- From: string, string
		- To: string
1. NewAlertsSuppressionRulesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAllowedConnectionsClient
	- Params
		- From: string, string
		- To: string
1. NewAllowedConnectionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAutoProvisioningSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewAutoProvisioningSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAutomationsClient
	- Params
		- From: string, string
		- To: string
1. NewAutomationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewCompliancesClient
	- Params
		- From: string, string
		- To: string
1. NewCompliancesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewConnectorsClient
	- Params
		- From: string, string
		- To: string
1. NewConnectorsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewContactsClient
	- Params
		- From: string, string
		- To: string
1. NewContactsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewDeviceSecurityGroupsClient
	- Params
		- From: string, string
		- To: string
1. NewDeviceSecurityGroupsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewDiscoveredSecuritySolutionsClient
	- Params
		- From: string, string
		- To: string
1. NewDiscoveredSecuritySolutionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewExternalSecuritySolutionsClient
	- Params
		- From: string, string
		- To: string
1. NewExternalSecuritySolutionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewInformationProtectionPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewInformationProtectionPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewJitNetworkAccessPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewJitNetworkAccessPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewLocationsClient
	- Params
		- From: string, string
		- To: string
1. NewLocationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewOperationsClient
	- Params
		- From: string, string
		- To: string
1. NewOperationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewPricingsClient
	- Params
		- From: string, string
		- To: string
1. NewPricingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceAssessmentsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceAssessmentsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceControlsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceStandardsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceStandardsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentBaselineRulesClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentBaselineRulesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentScanResultsClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentScanResultsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentScansClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentScansClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoreControlDefinitionsClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoreControlDefinitionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoreControlsClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoreControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoresClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoresClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSubAssessmentsClient
	- Params
		- From: string, string
		- To: string
1. NewSubAssessmentsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewTasksClient
	- Params
		- From: string, string
		- To: string
1. NewTasksClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewTopologyClient
	- Params
		- From: string, string
		- To: string
1. NewTopologyClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWorkspaceSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewWorkspaceSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. SettingsClient.Get
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.GetResponder
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.Update
	- Params
		- From: context.Context, string, Setting
		- To: context.Context, string, BasicSetting
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.UpdatePreparer
	- Params
		- From: context.Context, string, Setting
		- To: context.Context, string, BasicSetting
1. SettingsClient.UpdateResponder
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsListIterator.Value
	- Returns
		- From: Setting
		- To: BasicSetting
1. SettingsListPage.Values
	- Returns
		- From: []Setting
		- To: []BasicSetting
1. TasksClient.GetResourceGroupLevelTask
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.GetResourceGroupLevelTaskPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.GetSubscriptionLevelTask
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.GetSubscriptionLevelTaskPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegion
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegionComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByResourceGroup
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.ListByResourceGroupComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.ListByResourceGroupPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.UpdateResourceGroupLevelTaskState
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, TaskUpdateActionType
1. TasksClient.UpdateResourceGroupLevelTaskStatePreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, TaskUpdateActionType
1. TasksClient.UpdateSubscriptionLevelTaskState
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, TaskUpdateActionType
1. TasksClient.UpdateSubscriptionLevelTaskStatePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, TaskUpdateActionType
1. TopologyClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TopologyClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TopologyClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. TopologyClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. TopologyClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string

#### Struct Fields

1. AadExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. AtaExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. CefExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. DataExportSetting.Kind changed type from SettingKind to KindEnum
1. ExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. PathRecommendation.Type changed type from Type to Type1
1. Setting.Kind changed type from SettingKind to KindEnum
1. SettingsList.Value changed type from *[]Setting to *[]BasicSetting

## Additive Changes

### New Constants

1. CloudName.AWS
1. CloudName.Azure
1. CloudName.GCP
1. CloudName.Github
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. EndOfSupportStatus.NoLongerSupported
1. EndOfSupportStatus.UpcomingNoLongerSupported
1. EndOfSupportStatus.UpcomingVersionNoLongerSupported
1. EndOfSupportStatus.VersionNoLongerSupported
1. EnforcementMode.EnforcementModeAudit
1. EnforcementMode.EnforcementModeEnforce
1. EnforcementMode.EnforcementModeNone
1. EnvironmentType.EnvironmentTypeAwsAccount
1. EnvironmentType.EnvironmentTypeEnvironmentData
1. EnvironmentType.EnvironmentTypeGcpProject
1. EnvironmentType.EnvironmentTypeGithubScope
1. EventSource.EventSourceAlerts
1. EventSource.EventSourceAssessments
1. EventSource.EventSourceAssessmentsSnapshot
1. EventSource.EventSourceRegulatoryComplianceAssessment
1. EventSource.EventSourceRegulatoryComplianceAssessmentSnapshot
1. EventSource.EventSourceSecureScoreControls
1. EventSource.EventSourceSecureScoreControlsSnapshot
1. EventSource.EventSourceSecureScores
1. EventSource.EventSourceSecureScoresSnapshot
1. EventSource.EventSourceSubAssessments
1. EventSource.EventSourceSubAssessmentsSnapshot
1. InformationProtectionPolicyName.InformationProtectionPolicyNameCustom
1. InformationProtectionPolicyName.InformationProtectionPolicyNameEffective
1. KindEnum.KindDataExportSetting
1. KindEnum.KindSetting
1. OfferingType.OfferingTypeCloudOffering
1. OfferingType.OfferingTypeCspmMonitorAws
1. OfferingType.OfferingTypeCspmMonitorGcp
1. OfferingType.OfferingTypeCspmMonitorGithub
1. OfferingType.OfferingTypeDefenderForContainersAws
1. OfferingType.OfferingTypeDefenderForContainersGcp
1. OfferingType.OfferingTypeDefenderForServersAws
1. OfferingType.OfferingTypeDefenderForServersGcp
1. OfferingType.OfferingTypeInformationProtectionAws
1. OrganizationMembershipType.OrganizationMembershipTypeAwsOrganizationalData
1. OrganizationMembershipType.OrganizationMembershipTypeMember
1. OrganizationMembershipType.OrganizationMembershipTypeOrganization
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeGcpOrganizationalData
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeMember
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeOrganization
1. SeverityEnum.SeverityEnumHigh
1. SeverityEnum.SeverityEnumLow
1. SeverityEnum.SeverityEnumMedium
1. SubPlan.P1
1. SubPlan.P2
1. SupportedCloudEnum.SupportedCloudEnumAWS
1. SupportedCloudEnum.SupportedCloudEnumGCP
1. TaskUpdateActionType.Activate
1. TaskUpdateActionType.Close
1. TaskUpdateActionType.Dismiss
1. TaskUpdateActionType.Resolve
1. TaskUpdateActionType.Start
1. Type.Qualys
1. Type.TVM

### New Funcs

1. *AWSEnvironmentData.UnmarshalJSON([]byte) error
1. *Connector.UnmarshalJSON([]byte) error
1. *ConnectorProperties.UnmarshalJSON([]byte) error
1. *ConnectorsListIterator.Next() error
1. *ConnectorsListIterator.NextWithContext(context.Context) error
1. *ConnectorsListPage.Next() error
1. *ConnectorsListPage.NextWithContext(context.Context) error
1. *CustomAssessmentAutomation.UnmarshalJSON([]byte) error
1. *CustomAssessmentAutomationRequest.UnmarshalJSON([]byte) error
1. *CustomAssessmentAutomationsListResultIterator.Next() error
1. *CustomAssessmentAutomationsListResultIterator.NextWithContext(context.Context) error
1. *CustomAssessmentAutomationsListResultPage.Next() error
1. *CustomAssessmentAutomationsListResultPage.NextWithContext(context.Context) error
1. *CustomEntityStoreAssignment.UnmarshalJSON([]byte) error
1. *CustomEntityStoreAssignmentRequest.UnmarshalJSON([]byte) error
1. *CustomEntityStoreAssignmentsListResultIterator.Next() error
1. *CustomEntityStoreAssignmentsListResultIterator.NextWithContext(context.Context) error
1. *CustomEntityStoreAssignmentsListResultPage.Next() error
1. *CustomEntityStoreAssignmentsListResultPage.NextWithContext(context.Context) error
1. *GcpProjectEnvironmentData.UnmarshalJSON([]byte) error
1. *IngestionSettingListIterator.Next() error
1. *IngestionSettingListIterator.NextWithContext(context.Context) error
1. *IngestionSettingListPage.Next() error
1. *IngestionSettingListPage.NextWithContext(context.Context) error
1. *MdeOnboardingData.UnmarshalJSON([]byte) error
1. *SettingModel.UnmarshalJSON([]byte) error
1. *SettingsList.UnmarshalJSON([]byte) error
1. *Software.UnmarshalJSON([]byte) error
1. *SoftwaresListIterator.Next() error
1. *SoftwaresListIterator.NextWithContext(context.Context) error
1. *SoftwaresListPage.Next() error
1. *SoftwaresListPage.NextWithContext(context.Context) error
1. AWSEnvironmentData.AsAWSEnvironmentData() (*AWSEnvironmentData, bool)
1. AWSEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. AWSEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. AWSEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. AWSEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. AWSEnvironmentData.MarshalJSON() ([]byte, error)
1. AwsOrganizationalData.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalData.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalData.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalData.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalData.MarshalJSON() ([]byte, error)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalDataMaster.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalDataMaster.MarshalJSON() ([]byte, error)
1. AwsOrganizationalDataMember.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalDataMember.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalDataMember.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalDataMember.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalDataMember.MarshalJSON() ([]byte, error)
1. CloudOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CloudOffering.AsCloudOffering() (*CloudOffering, bool)
1. CloudOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CloudOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CloudOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CloudOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CloudOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CloudOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CloudOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CloudOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CloudOffering.MarshalJSON() ([]byte, error)
1. Connector.MarshalJSON() ([]byte, error)
1. ConnectorsClient.ListByResourceGroup(context.Context, string) (ConnectorsListPage, error)
1. ConnectorsClient.ListByResourceGroupComplete(context.Context, string) (ConnectorsListIterator, error)
1. ConnectorsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. ConnectorsClient.ListByResourceGroupResponder(*http.Response) (ConnectorsList, error)
1. ConnectorsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. ConnectorsClient.Update(context.Context, string, string, Connector) (Connector, error)
1. ConnectorsClient.UpdatePreparer(context.Context, string, string, Connector) (*http.Request, error)
1. ConnectorsClient.UpdateResponder(*http.Response) (Connector, error)
1. ConnectorsClient.UpdateSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.CreateOrUpdate(context.Context, string, ConnectorSetting) (ConnectorSetting, error)
1. ConnectorsGroupClient.CreateOrUpdatePreparer(context.Context, string, ConnectorSetting) (*http.Request, error)
1. ConnectorsGroupClient.CreateOrUpdateResponder(*http.Response) (ConnectorSetting, error)
1. ConnectorsGroupClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.Delete(context.Context, string) (autorest.Response, error)
1. ConnectorsGroupClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. ConnectorsGroupClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConnectorsGroupClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.Get(context.Context, string) (ConnectorSetting, error)
1. ConnectorsGroupClient.GetPreparer(context.Context, string) (*http.Request, error)
1. ConnectorsGroupClient.GetResponder(*http.Response) (ConnectorSetting, error)
1. ConnectorsGroupClient.GetSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.List(context.Context) (ConnectorSettingListPage, error)
1. ConnectorsGroupClient.ListComplete(context.Context) (ConnectorSettingListIterator, error)
1. ConnectorsGroupClient.ListPreparer(context.Context) (*http.Request, error)
1. ConnectorsGroupClient.ListResponder(*http.Response) (ConnectorSettingList, error)
1. ConnectorsGroupClient.ListSender(*http.Request) (*http.Response, error)
1. ConnectorsList.IsEmpty() bool
1. ConnectorsList.MarshalJSON() ([]byte, error)
1. ConnectorsListIterator.NotDone() bool
1. ConnectorsListIterator.Response() ConnectorsList
1. ConnectorsListIterator.Value() Connector
1. ConnectorsListPage.NotDone() bool
1. ConnectorsListPage.Response() ConnectorsList
1. ConnectorsListPage.Values() []Connector
1. CspmMonitorAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorAwsOffering.MarshalJSON() ([]byte, error)
1. CspmMonitorGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorGcpOffering.MarshalJSON() ([]byte, error)
1. CspmMonitorGithubOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorGithubOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorGithubOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorGithubOffering.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomation.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationRequest.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationsClient.Create(context.Context, string, string, CustomAssessmentAutomationRequest) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.CreatePreparer(context.Context, string, string, CustomAssessmentAutomationRequest) (*http.Request, error)
1. CustomAssessmentAutomationsClient.CreateResponder(*http.Response) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.CreateSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. CustomAssessmentAutomationsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. CustomAssessmentAutomationsClient.DeleteSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.Get(context.Context, string, string) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.GetResponder(*http.Response) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.GetSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroup(context.Context, string) (CustomAssessmentAutomationsListResultPage, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupComplete(context.Context, string) (CustomAssessmentAutomationsListResultIterator, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupResponder(*http.Response) (CustomAssessmentAutomationsListResult, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.ListBySubscription(context.Context) (CustomAssessmentAutomationsListResultPage, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionComplete(context.Context) (CustomAssessmentAutomationsListResultIterator, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionResponder(*http.Response) (CustomAssessmentAutomationsListResult, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsListResult.IsEmpty() bool
1. CustomAssessmentAutomationsListResult.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationsListResultIterator.NotDone() bool
1. CustomAssessmentAutomationsListResultIterator.Response() CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultIterator.Value() CustomAssessmentAutomation
1. CustomAssessmentAutomationsListResultPage.NotDone() bool
1. CustomAssessmentAutomationsListResultPage.Response() CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultPage.Values() []CustomAssessmentAutomation
1. CustomEntityStoreAssignment.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentRequest.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentsClient.Create(context.Context, string, string, CustomEntityStoreAssignmentRequest) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.CreatePreparer(context.Context, string, string, CustomEntityStoreAssignmentRequest) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.CreateResponder(*http.Response) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.CreateSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. CustomEntityStoreAssignmentsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. CustomEntityStoreAssignmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.Get(context.Context, string, string) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.GetResponder(*http.Response) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.GetSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroup(context.Context, string) (CustomEntityStoreAssignmentsListResultPage, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupComplete(context.Context, string) (CustomEntityStoreAssignmentsListResultIterator, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupResponder(*http.Response) (CustomEntityStoreAssignmentsListResult, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscription(context.Context) (CustomEntityStoreAssignmentsListResultPage, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionComplete(context.Context) (CustomEntityStoreAssignmentsListResultIterator, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionResponder(*http.Response) (CustomEntityStoreAssignmentsListResult, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsListResult.IsEmpty() bool
1. CustomEntityStoreAssignmentsListResult.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentsListResultIterator.NotDone() bool
1. CustomEntityStoreAssignmentsListResultIterator.Response() CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultIterator.Value() CustomEntityStoreAssignment
1. CustomEntityStoreAssignmentsListResultPage.NotDone() bool
1. CustomEntityStoreAssignmentsListResultPage.Response() CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultPage.Values() []CustomEntityStoreAssignment
1. DataExportSetting.AsBasicSetting() (BasicSetting, bool)
1. DataExportSetting.AsDataExportSetting() (*DataExportSetting, bool)
1. DataExportSetting.AsSetting() (*Setting, bool)
1. DefenderForContainersAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForContainersAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForContainersAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderForContainersGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForContainersGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForContainersGcpOffering.MarshalJSON() ([]byte, error)
1. DefenderForServersAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForServersAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForServersAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForServersAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderForServersGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForServersGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForServersGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForServersGcpOffering.MarshalJSON() ([]byte, error)
1. EnvironmentData.AsAWSEnvironmentData() (*AWSEnvironmentData, bool)
1. EnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. EnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. EnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. EnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. EnvironmentData.MarshalJSON() ([]byte, error)
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. GcpOrganizationalData.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalData.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalData.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalData.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalData.MarshalJSON() ([]byte, error)
1. GcpOrganizationalDataMember.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalDataMember.MarshalJSON() ([]byte, error)
1. GcpOrganizationalDataOrganization.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalDataOrganization.MarshalJSON() ([]byte, error)
1. GcpProjectDetails.MarshalJSON() ([]byte, error)
1. GcpProjectEnvironmentData.AsAWSEnvironmentData() (*AWSEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. GcpProjectEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. GcpProjectEnvironmentData.MarshalJSON() ([]byte, error)
1. GithubScopeEnvironmentData.AsAWSEnvironmentData() (*AWSEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. GithubScopeEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. GithubScopeEnvironmentData.MarshalJSON() ([]byte, error)
1. InformationProtectionAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. InformationProtectionAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. InformationProtectionAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. InformationProtectionAwsOffering.MarshalJSON() ([]byte, error)
1. IngestionConnectionString.MarshalJSON() ([]byte, error)
1. IngestionSetting.MarshalJSON() ([]byte, error)
1. IngestionSettingList.IsEmpty() bool
1. IngestionSettingList.MarshalJSON() ([]byte, error)
1. IngestionSettingListIterator.NotDone() bool
1. IngestionSettingListIterator.Response() IngestionSettingList
1. IngestionSettingListIterator.Value() IngestionSetting
1. IngestionSettingListPage.NotDone() bool
1. IngestionSettingListPage.Response() IngestionSettingList
1. IngestionSettingListPage.Values() []IngestionSetting
1. IngestionSettingToken.MarshalJSON() ([]byte, error)
1. IngestionSettingsClient.Create(context.Context, string, IngestionSetting) (IngestionSetting, error)
1. IngestionSettingsClient.CreatePreparer(context.Context, string, IngestionSetting) (*http.Request, error)
1. IngestionSettingsClient.CreateResponder(*http.Response) (IngestionSetting, error)
1. IngestionSettingsClient.CreateSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.Delete(context.Context, string) (autorest.Response, error)
1. IngestionSettingsClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IngestionSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.Get(context.Context, string) (IngestionSetting, error)
1. IngestionSettingsClient.GetPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.GetResponder(*http.Response) (IngestionSetting, error)
1. IngestionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.List(context.Context) (IngestionSettingListPage, error)
1. IngestionSettingsClient.ListComplete(context.Context) (IngestionSettingListIterator, error)
1. IngestionSettingsClient.ListConnectionStrings(context.Context, string) (ConnectionStrings, error)
1. IngestionSettingsClient.ListConnectionStringsPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.ListConnectionStringsResponder(*http.Response) (ConnectionStrings, error)
1. IngestionSettingsClient.ListConnectionStringsSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.ListPreparer(context.Context) (*http.Request, error)
1. IngestionSettingsClient.ListResponder(*http.Response) (IngestionSettingList, error)
1. IngestionSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.ListTokens(context.Context, string) (IngestionSettingToken, error)
1. IngestionSettingsClient.ListTokensPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.ListTokensResponder(*http.Response) (IngestionSettingToken, error)
1. IngestionSettingsClient.ListTokensSender(*http.Request) (*http.Response, error)
1. MdeOnboardingData.MarshalJSON() ([]byte, error)
1. MdeOnboardingsClient.Get(context.Context) (MdeOnboardingData, error)
1. MdeOnboardingsClient.GetPreparer(context.Context) (*http.Request, error)
1. MdeOnboardingsClient.GetResponder(*http.Response) (MdeOnboardingData, error)
1. MdeOnboardingsClient.GetSender(*http.Request) (*http.Response, error)
1. MdeOnboardingsClient.List(context.Context) (MdeOnboardingDataList, error)
1. MdeOnboardingsClient.ListPreparer(context.Context) (*http.Request, error)
1. MdeOnboardingsClient.ListResponder(*http.Response) (MdeOnboardingDataList, error)
1. MdeOnboardingsClient.ListSender(*http.Request) (*http.Response, error)
1. NewConnectorsGroupClient(string) ConnectorsGroupClient
1. NewConnectorsGroupClientWithBaseURI(string, string) ConnectorsGroupClient
1. NewConnectorsListIterator(ConnectorsListPage) ConnectorsListIterator
1. NewConnectorsListPage(ConnectorsList, func(context.Context, ConnectorsList) (ConnectorsList, error)) ConnectorsListPage
1. NewCustomAssessmentAutomationsClient(string) CustomAssessmentAutomationsClient
1. NewCustomAssessmentAutomationsClientWithBaseURI(string, string) CustomAssessmentAutomationsClient
1. NewCustomAssessmentAutomationsListResultIterator(CustomAssessmentAutomationsListResultPage) CustomAssessmentAutomationsListResultIterator
1. NewCustomAssessmentAutomationsListResultPage(CustomAssessmentAutomationsListResult, func(context.Context, CustomAssessmentAutomationsListResult) (CustomAssessmentAutomationsListResult, error)) CustomAssessmentAutomationsListResultPage
1. NewCustomEntityStoreAssignmentsClient(string) CustomEntityStoreAssignmentsClient
1. NewCustomEntityStoreAssignmentsClientWithBaseURI(string, string) CustomEntityStoreAssignmentsClient
1. NewCustomEntityStoreAssignmentsListResultIterator(CustomEntityStoreAssignmentsListResultPage) CustomEntityStoreAssignmentsListResultIterator
1. NewCustomEntityStoreAssignmentsListResultPage(CustomEntityStoreAssignmentsListResult, func(context.Context, CustomEntityStoreAssignmentsListResult) (CustomEntityStoreAssignmentsListResult, error)) CustomEntityStoreAssignmentsListResultPage
1. NewIngestionSettingListIterator(IngestionSettingListPage) IngestionSettingListIterator
1. NewIngestionSettingListPage(IngestionSettingList, func(context.Context, IngestionSettingList) (IngestionSettingList, error)) IngestionSettingListPage
1. NewIngestionSettingsClient(string) IngestionSettingsClient
1. NewIngestionSettingsClientWithBaseURI(string, string) IngestionSettingsClient
1. NewMdeOnboardingsClient(string) MdeOnboardingsClient
1. NewMdeOnboardingsClientWithBaseURI(string, string) MdeOnboardingsClient
1. NewSoftwareInventoriesClient(string) SoftwareInventoriesClient
1. NewSoftwareInventoriesClientWithBaseURI(string, string) SoftwareInventoriesClient
1. NewSoftwaresListIterator(SoftwaresListPage) SoftwaresListIterator
1. NewSoftwaresListPage(SoftwaresList, func(context.Context, SoftwaresList) (SoftwaresList, error)) SoftwaresListPage
1. PossibleCloudNameValues() []CloudName
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleEndOfSupportStatusValues() []EndOfSupportStatus
1. PossibleEnvironmentTypeValues() []EnvironmentType
1. PossibleInformationProtectionPolicyNameValues() []InformationProtectionPolicyName
1. PossibleKindEnum1Values() []KindEnum1
1. PossibleOfferingTypeValues() []OfferingType
1. PossibleOrganizationMembershipTypeBasicGcpOrganizationalDataValues() []OrganizationMembershipTypeBasicGcpOrganizationalData
1. PossibleOrganizationMembershipTypeValues() []OrganizationMembershipType
1. PossibleSeverityEnumValues() []SeverityEnum
1. PossibleSubPlanValues() []SubPlan
1. PossibleSupportedCloudEnumValues() []SupportedCloudEnum
1. PossibleTaskUpdateActionTypeValues() []TaskUpdateActionType
1. PossibleType1Values() []Type1
1. Setting.AsBasicSetting() (BasicSetting, bool)
1. Setting.AsDataExportSetting() (*DataExportSetting, bool)
1. Setting.AsSetting() (*Setting, bool)
1. Software.MarshalJSON() ([]byte, error)
1. SoftwareInventoriesClient.Get(context.Context, string, string, string, string, string) (Software, error)
1. SoftwareInventoriesClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. SoftwareInventoriesClient.GetResponder(*http.Response) (Software, error)
1. SoftwareInventoriesClient.GetSender(*http.Request) (*http.Response, error)
1. SoftwareInventoriesClient.ListByExtendedResource(context.Context, string, string, string, string) (SoftwaresListPage, error)
1. SoftwareInventoriesClient.ListByExtendedResourceComplete(context.Context, string, string, string, string) (SoftwaresListIterator, error)
1. SoftwareInventoriesClient.ListByExtendedResourcePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. SoftwareInventoriesClient.ListByExtendedResourceResponder(*http.Response) (SoftwaresList, error)
1. SoftwareInventoriesClient.ListByExtendedResourceSender(*http.Request) (*http.Response, error)
1. SoftwareInventoriesClient.ListBySubscription(context.Context) (SoftwaresListPage, error)
1. SoftwareInventoriesClient.ListBySubscriptionComplete(context.Context) (SoftwaresListIterator, error)
1. SoftwareInventoriesClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. SoftwareInventoriesClient.ListBySubscriptionResponder(*http.Response) (SoftwaresList, error)
1. SoftwareInventoriesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. SoftwaresList.IsEmpty() bool
1. SoftwaresList.MarshalJSON() ([]byte, error)
1. SoftwaresListIterator.NotDone() bool
1. SoftwaresListIterator.Response() SoftwaresList
1. SoftwaresListIterator.Value() Software
1. SoftwaresListPage.NotDone() bool
1. SoftwaresListPage.Response() SoftwaresList
1. SoftwaresListPage.Values() []Software

### Struct Changes

#### New Structs

1. AWSEnvironmentData
1. AwsOrganizationalData
1. AwsOrganizationalDataMaster
1. AwsOrganizationalDataMember
1. CloudOffering
1. ConnectionStrings
1. Connector
1. ConnectorProperties
1. ConnectorsGroupClient
1. ConnectorsList
1. ConnectorsListIterator
1. ConnectorsListPage
1. CspmMonitorAwsOffering
1. CspmMonitorAwsOfferingNativeCloudConnection
1. CspmMonitorGcpOffering
1. CspmMonitorGcpOfferingNativeCloudConnection
1. CspmMonitorGithubOffering
1. CustomAssessmentAutomation
1. CustomAssessmentAutomationProperties
1. CustomAssessmentAutomationRequest
1. CustomAssessmentAutomationRequestProperties
1. CustomAssessmentAutomationsClient
1. CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultIterator
1. CustomAssessmentAutomationsListResultPage
1. CustomEntityStoreAssignment
1. CustomEntityStoreAssignmentProperties
1. CustomEntityStoreAssignmentRequest
1. CustomEntityStoreAssignmentRequestProperties
1. CustomEntityStoreAssignmentsClient
1. CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultIterator
1. CustomEntityStoreAssignmentsListResultPage
1. DefenderForContainersAwsOffering
1. DefenderForContainersAwsOfferingCloudWatchToKinesis
1. DefenderForContainersAwsOfferingKinesisToS3
1. DefenderForContainersAwsOfferingKubernetesScubaReader
1. DefenderForContainersAwsOfferingKubernetesService
1. DefenderForContainersGcpOffering
1. DefenderForContainersGcpOfferingDataPipelineNativeCloudConnection
1. DefenderForContainersGcpOfferingNativeCloudConnection
1. DefenderForServersAwsOffering
1. DefenderForServersAwsOfferingArcAutoProvisioning
1. DefenderForServersAwsOfferingArcAutoProvisioningServicePrincipalSecretMetadata
1. DefenderForServersAwsOfferingDefenderForServers
1. DefenderForServersAwsOfferingMdeAutoProvisioning
1. DefenderForServersAwsOfferingSubPlan
1. DefenderForServersAwsOfferingVaAutoProvisioning
1. DefenderForServersAwsOfferingVaAutoProvisioningConfiguration
1. DefenderForServersGcpOffering
1. DefenderForServersGcpOfferingArcAutoProvisioning
1. DefenderForServersGcpOfferingArcAutoProvisioningConfiguration
1. DefenderForServersGcpOfferingDefenderForServers
1. DefenderForServersGcpOfferingMdeAutoProvisioning
1. DefenderForServersGcpOfferingSubPlan
1. DefenderForServersGcpOfferingVaAutoProvisioning
1. DefenderForServersGcpOfferingVaAutoProvisioningConfiguration
1. EnvironmentData
1. ErrorAdditionalInfo
1. GcpOrganizationalData
1. GcpOrganizationalDataMember
1. GcpOrganizationalDataOrganization
1. GcpProjectDetails
1. GcpProjectEnvironmentData
1. GithubScopeEnvironmentData
1. InformationProtectionAwsOffering
1. InformationProtectionAwsOfferingInformationProtection
1. IngestionConnectionString
1. IngestionSetting
1. IngestionSettingList
1. IngestionSettingListIterator
1. IngestionSettingListPage
1. IngestionSettingToken
1. IngestionSettingsClient
1. MdeOnboardingData
1. MdeOnboardingDataList
1. MdeOnboardingDataProperties
1. MdeOnboardingsClient
1. SettingModel
1. Software
1. SoftwareInventoriesClient
1. SoftwareProperties
1. SoftwaresList
1. SoftwaresListIterator
1. SoftwaresListPage
1. SystemData

#### New Struct Fields

1. CloudErrorBody.AdditionalInfo
1. CloudErrorBody.Details
1. CloudErrorBody.Target
