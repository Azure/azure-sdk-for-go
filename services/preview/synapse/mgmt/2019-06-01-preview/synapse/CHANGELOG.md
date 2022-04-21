# Unreleased

## Breaking Changes

### Removed Constants

1. VulnerabilityAssessmentPolicyBaselineName.Master

### Struct Changes

#### Removed Struct Fields

1. BigDataPoolResourceProperties.HaveLibraryRequirementsChanged

### Signature Changes

#### Const Types

1. Default changed type from VulnerabilityAssessmentPolicyBaselineName to CreateMode

#### Struct Fields

1. SQLPoolResourceProperties.CreateMode changed type from *string to CreateMode

## Additive Changes

### New Constants

1. CreateMode.PointInTimeRestore
1. CreateMode.Recovery
1. CreateMode.Restore
1. DayOfWeek.Friday
1. DayOfWeek.Monday
1. DayOfWeek.Saturday
1. DayOfWeek.Sunday
1. DayOfWeek.Thursday
1. DayOfWeek.Tuesday
1. DayOfWeek.Wednesday
1. RecommendedSensitivityLabelUpdateKind.Disable
1. RecommendedSensitivityLabelUpdateKind.Enable
1. SensitivityLabelRank.SensitivityLabelRankCritical
1. SensitivityLabelRank.SensitivityLabelRankHigh
1. SensitivityLabelRank.SensitivityLabelRankLow
1. SensitivityLabelRank.SensitivityLabelRankMedium
1. SensitivityLabelRank.SensitivityLabelRankNone
1. SensitivityLabelUpdateKind.Remove
1. SensitivityLabelUpdateKind.Set
1. VulnerabilityAssessmentPolicyBaselineName.VulnerabilityAssessmentPolicyBaselineNameDefault
1. VulnerabilityAssessmentPolicyBaselineName.VulnerabilityAssessmentPolicyBaselineNameMaster

### New Funcs

1. *MaintenanceWindowOptions.UnmarshalJSON([]byte) error
1. *MaintenanceWindows.UnmarshalJSON([]byte) error
1. *RecommendedSensitivityLabelUpdate.UnmarshalJSON([]byte) error
1. *SensitivityLabelUpdate.UnmarshalJSON([]byte) error
1. BigDataPoolResourceProperties.MarshalJSON() ([]byte, error)
1. DataMaskingRulesClient.Get(context.Context, string, string, string, string) (DataMaskingRule, error)
1. DataMaskingRulesClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DataMaskingRulesClient.GetResponder(*http.Response) (DataMaskingRule, error)
1. DataMaskingRulesClient.GetSender(*http.Request) (*http.Response, error)
1. LibraryInfo.MarshalJSON() ([]byte, error)
1. MaintenanceWindowOptions.MarshalJSON() ([]byte, error)
1. MaintenanceWindows.MarshalJSON() ([]byte, error)
1. NewPrivateLinkHubPrivateLinkResourcesClient(string) PrivateLinkHubPrivateLinkResourcesClient
1. NewPrivateLinkHubPrivateLinkResourcesClientWithBaseURI(string, string) PrivateLinkHubPrivateLinkResourcesClient
1. NewSQLPoolMaintenanceWindowOptionsClient(string) SQLPoolMaintenanceWindowOptionsClient
1. NewSQLPoolMaintenanceWindowOptionsClientWithBaseURI(string, string) SQLPoolMaintenanceWindowOptionsClient
1. NewSQLPoolMaintenanceWindowsClient(string) SQLPoolMaintenanceWindowsClient
1. NewSQLPoolMaintenanceWindowsClientWithBaseURI(string, string) SQLPoolMaintenanceWindowsClient
1. NewSQLPoolRecommendedSensitivityLabelsClient(string) SQLPoolRecommendedSensitivityLabelsClient
1. NewSQLPoolRecommendedSensitivityLabelsClientWithBaseURI(string, string) SQLPoolRecommendedSensitivityLabelsClient
1. PossibleCreateModeValues() []CreateMode
1. PossibleDayOfWeekValues() []DayOfWeek
1. PossibleRecommendedSensitivityLabelUpdateKindValues() []RecommendedSensitivityLabelUpdateKind
1. PossibleSensitivityLabelRankValues() []SensitivityLabelRank
1. PossibleSensitivityLabelUpdateKindValues() []SensitivityLabelUpdateKind
1. PrivateLinkHubPrivateLinkResourcesClient.Get(context.Context, string, string, string) (PrivateLinkResource, error)
1. PrivateLinkHubPrivateLinkResourcesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateLinkHubPrivateLinkResourcesClient.GetResponder(*http.Response) (PrivateLinkResource, error)
1. PrivateLinkHubPrivateLinkResourcesClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkHubPrivateLinkResourcesClient.List(context.Context, string, string) (PrivateLinkResourceListResultPage, error)
1. PrivateLinkHubPrivateLinkResourcesClient.ListComplete(context.Context, string, string) (PrivateLinkResourceListResultIterator, error)
1. PrivateLinkHubPrivateLinkResourcesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkHubPrivateLinkResourcesClient.ListResponder(*http.Response) (PrivateLinkResourceListResult, error)
1. PrivateLinkHubPrivateLinkResourcesClient.ListSender(*http.Request) (*http.Response, error)
1. RecommendedSensitivityLabelUpdate.MarshalJSON() ([]byte, error)
1. SQLPoolColumnProperties.MarshalJSON() ([]byte, error)
1. SQLPoolGeoBackupPoliciesClient.CreateOrUpdate(context.Context, string, string, string, GeoBackupPolicy) (GeoBackupPolicy, error)
1. SQLPoolGeoBackupPoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, string, GeoBackupPolicy) (*http.Request, error)
1. SQLPoolGeoBackupPoliciesClient.CreateOrUpdateResponder(*http.Response) (GeoBackupPolicy, error)
1. SQLPoolGeoBackupPoliciesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. SQLPoolMaintenanceWindowOptionsClient.Get(context.Context, string, string, string, string) (MaintenanceWindowOptions, error)
1. SQLPoolMaintenanceWindowOptionsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. SQLPoolMaintenanceWindowOptionsClient.GetResponder(*http.Response) (MaintenanceWindowOptions, error)
1. SQLPoolMaintenanceWindowOptionsClient.GetSender(*http.Request) (*http.Response, error)
1. SQLPoolMaintenanceWindowsClient.CreateOrUpdate(context.Context, string, string, string, string, MaintenanceWindows) (autorest.Response, error)
1. SQLPoolMaintenanceWindowsClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, MaintenanceWindows) (*http.Request, error)
1. SQLPoolMaintenanceWindowsClient.CreateOrUpdateResponder(*http.Response) (autorest.Response, error)
1. SQLPoolMaintenanceWindowsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. SQLPoolMaintenanceWindowsClient.Get(context.Context, string, string, string, string) (MaintenanceWindows, error)
1. SQLPoolMaintenanceWindowsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. SQLPoolMaintenanceWindowsClient.GetResponder(*http.Response) (MaintenanceWindows, error)
1. SQLPoolMaintenanceWindowsClient.GetSender(*http.Request) (*http.Response, error)
1. SQLPoolRecommendedSensitivityLabelsClient.Update(context.Context, string, string, string, RecommendedSensitivityLabelUpdateList) (autorest.Response, error)
1. SQLPoolRecommendedSensitivityLabelsClient.UpdatePreparer(context.Context, string, string, string, RecommendedSensitivityLabelUpdateList) (*http.Request, error)
1. SQLPoolRecommendedSensitivityLabelsClient.UpdateResponder(*http.Response) (autorest.Response, error)
1. SQLPoolRecommendedSensitivityLabelsClient.UpdateSender(*http.Request) (*http.Response, error)
1. SQLPoolSensitivityLabelsClient.Update(context.Context, string, string, string, SensitivityLabelUpdateList) (autorest.Response, error)
1. SQLPoolSensitivityLabelsClient.UpdatePreparer(context.Context, string, string, string, SensitivityLabelUpdateList) (*http.Request, error)
1. SQLPoolSensitivityLabelsClient.UpdateResponder(*http.Response) (autorest.Response, error)
1. SQLPoolSensitivityLabelsClient.UpdateSender(*http.Request) (*http.Response, error)
1. SensitivityLabelUpdate.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. DynamicExecutorAllocation
1. GitHubClientSecret
1. LibraryInfo
1. MaintenanceWindowOptions
1. MaintenanceWindowOptionsProperties
1. MaintenanceWindowTimeRange
1. MaintenanceWindows
1. MaintenanceWindowsProperties
1. ManagedVirtualNetworkReference
1. PrivateLinkHubPrivateLinkResourcesClient
1. RecommendedSensitivityLabelUpdate
1. RecommendedSensitivityLabelUpdateList
1. RecommendedSensitivityLabelUpdateProperties
1. SQLPoolMaintenanceWindowOptionsClient
1. SQLPoolMaintenanceWindowsClient
1. SQLPoolRecommendedSensitivityLabelsClient
1. SensitivityLabelUpdate
1. SensitivityLabelUpdateList
1. SensitivityLabelUpdateProperties

#### New Struct Fields

1. BigDataPoolResourceProperties.CacheSize
1. BigDataPoolResourceProperties.CustomLibraries
1. BigDataPoolResourceProperties.DynamicExecutorAllocation
1. BigDataPoolResourceProperties.LastSucceededTimestamp
1. IntegrationRuntimeDataFlowProperties.Cleanup
1. ManagedIntegrationRuntime.ManagedVirtualNetwork
1. SQLPoolColumnProperties.IsComputed
1. SensitivityLabel.ManagedBy
1. SensitivityLabelProperties.ColumnName
1. SensitivityLabelProperties.Rank
1. SensitivityLabelProperties.SchemaName
1. SensitivityLabelProperties.TableName
1. WorkspaceProperties.AdlaResourceID
1. WorkspaceRepositoryConfiguration.ClientID
1. WorkspaceRepositoryConfiguration.ClientSecret
1. WorkspaceRepositoryConfiguration.LastCommitID
1. WorkspaceRepositoryConfiguration.TenantID
