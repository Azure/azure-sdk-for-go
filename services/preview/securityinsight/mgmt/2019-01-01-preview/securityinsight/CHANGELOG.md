Generated from https://github.com/Azure/azure-rest-api-specs/tree/138759b8a5987252fd66658078907e1d93969c85/specification/securityinsights/resource-manager/readme.md tag: `package-2019-01-preview-only`

Code generator @microsoft.azure/autorest.go@2.1.169


## Breaking Changes

### Removed Constants

1. AlertRuleKind.Fusion
1. AlertRuleKind.MicrosoftSecurityIncidentCreation
1. AlertRuleKind.Scheduled
1. FileHashAlgorithm.MD5
1. FileHashAlgorithm.SHA1
1. FileHashAlgorithm.SHA256
1. FileHashAlgorithm.SHA256AC
1. OSFamily.Android
1. OSFamily.IOS
1. OSFamily.Linux
1. OSFamily.Windows

### Removed Funcs

1. *EntityQuery.UnmarshalJSON([]byte) error
1. AlertRulesClient.CreateOrUpdateAction(context.Context, string, string, string, string, string, ActionRequest) (ActionResponse, error)
1. AlertRulesClient.CreateOrUpdateActionPreparer(context.Context, string, string, string, string, string, ActionRequest) (*http.Request, error)
1. AlertRulesClient.CreateOrUpdateActionResponder(*http.Response) (ActionResponse, error)
1. AlertRulesClient.CreateOrUpdateActionSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.DeleteAction(context.Context, string, string, string, string, string) (autorest.Response, error)
1. AlertRulesClient.DeleteActionPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. AlertRulesClient.DeleteActionResponder(*http.Response) (autorest.Response, error)
1. AlertRulesClient.DeleteActionSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.GetAction(context.Context, string, string, string, string, string) (ActionResponse, error)
1. AlertRulesClient.GetActionPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. AlertRulesClient.GetActionResponder(*http.Response) (ActionResponse, error)
1. AlertRulesClient.GetActionSender(*http.Request) (*http.Response, error)

## Struct Changes

### Removed Structs

1. EntityQueryProperties

### Removed Struct Fields

1. EntityQuery.*EntityQueryProperties

## Signature Changes

### Const Types

1. Unknown changed type from FileHashAlgorithm to AntispamMailDirection

### Funcs

1. EntityQueriesClient.Get
	- Returns
		- From: EntityQuery, error
		- To: EntityQueryModel, error
1. EntityQueriesClient.GetResponder
	- Returns
		- From: EntityQuery, error
		- To: EntityQueryModel, error
1. EntityQueryListIterator.Value
	- Returns
		- From: EntityQuery
		- To: BasicEntityQuery
1. EntityQueryListPage.Values
	- Returns
		- From: []EntityQuery
		- To: []BasicEntityQuery

### Struct Fields

1. EntityQueryList.Value changed type from *[]EntityQuery to *[]BasicEntityQuery

### New Constants

1. AlertRuleKind.AlertRuleKindAnomaly
1. AlertRuleKind.AlertRuleKindFusion
1. AlertRuleKind.AlertRuleKindMLBehaviorAnalytics
1. AlertRuleKind.AlertRuleKindMicrosoftSecurityIncidentCreation
1. AlertRuleKind.AlertRuleKindScheduled
1. AlertRuleKind.AlertRuleKindThreatIntelligence
1. AntispamMailDirection.Inbound
1. AntispamMailDirection.Intraorg
1. AntispamMailDirection.Outbound
1. DataConnectorKind.DataConnectorKindDynamics365
1. DeliveryAction.DeliveryActionBlocked
1. DeliveryAction.DeliveryActionDelivered
1. DeliveryAction.DeliveryActionDeliveredAsSpam
1. DeliveryAction.DeliveryActionReplaced
1. DeliveryAction.DeliveryActionUnknown
1. DeliveryLocation.DeliveryLocationDeletedFolder
1. DeliveryLocation.DeliveryLocationDropped
1. DeliveryLocation.DeliveryLocationExternal
1. DeliveryLocation.DeliveryLocationFailed
1. DeliveryLocation.DeliveryLocationForwarded
1. DeliveryLocation.DeliveryLocationInbox
1. DeliveryLocation.DeliveryLocationJunkFolder
1. DeliveryLocation.DeliveryLocationQuarantine
1. DeliveryLocation.DeliveryLocationUnknown
1. EntityKind.EntityKindMailCluster
1. EntityKind.EntityKindMailMessage
1. EntityKind.EntityKindMailbox
1. EntityKind.EntityKindSubmissionMail
1. EntityQueryKind.Expansion
1. EntityQueryKind.Insight
1. EntityType.EntityTypeMailCluster
1. EntityType.EntityTypeMailMessage
1. EntityType.EntityTypeMailbox
1. EntityType.EntityTypeSubmissionMail
1. FileHashAlgorithm.FileHashAlgorithmMD5
1. FileHashAlgorithm.FileHashAlgorithmSHA1
1. FileHashAlgorithm.FileHashAlgorithmSHA256
1. FileHashAlgorithm.FileHashAlgorithmSHA256AC
1. FileHashAlgorithm.FileHashAlgorithmUnknown
1. GroupingEntityType.FileHash
1. KindBasicAlertRule.KindMLBehaviorAnalytics
1. KindBasicAlertRuleTemplate.KindBasicAlertRuleTemplateKindMLBehaviorAnalytics
1. KindBasicDataConnector.KindDynamics365
1. KindBasicDataConnectorsCheckRequirements.KindBasicDataConnectorsCheckRequirementsKindDynamics365
1. KindBasicEntity.KindMailCluster
1. KindBasicEntity.KindMailMessage
1. KindBasicEntity.KindMailbox
1. KindBasicEntity.KindSubmissionMail
1. KindBasicEntityQuery.KindEntityQuery
1. KindBasicEntityQuery.KindExpansion
1. KindBasicEntityQueryItem.KindEntityQueryItem
1. KindBasicEntityQueryItem.KindInsight
1. OSFamily.OSFamilyAndroid
1. OSFamily.OSFamilyIOS
1. OSFamily.OSFamilyLinux
1. OSFamily.OSFamilyUnknown
1. OSFamily.OSFamilyWindows
1. OutputType.OutputTypeDate
1. OutputType.OutputTypeEntity
1. OutputType.OutputTypeNumber
1. OutputType.OutputTypeString

### New Funcs

1. *Dynamics365CheckRequirements.UnmarshalJSON([]byte) error
1. *Dynamics365DataConnector.UnmarshalJSON([]byte) error
1. *EntityQueryList.UnmarshalJSON([]byte) error
1. *EntityQueryModel.UnmarshalJSON([]byte) error
1. *ExpansionEntityQuery.UnmarshalJSON([]byte) error
1. *GetQueriesResponse.UnmarshalJSON([]byte) error
1. *MLBehaviorAnalyticsAlertRule.UnmarshalJSON([]byte) error
1. *MLBehaviorAnalyticsAlertRuleTemplate.UnmarshalJSON([]byte) error
1. *MailClusterEntity.UnmarshalJSON([]byte) error
1. *MailMessageEntity.UnmarshalJSON([]byte) error
1. *MailboxEntity.UnmarshalJSON([]byte) error
1. *SubmissionMailEntity.UnmarshalJSON([]byte) error
1. *WatchlistItem.UnmarshalJSON([]byte) error
1. AADCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AADDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. AATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. ASCCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. ASCDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. AccountEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. AccountEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. AccountEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. AccountEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. ActionsClient.CreateOrUpdate(context.Context, string, string, string, string, string, ActionRequest) (ActionResponse, error)
1. ActionsClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, string, ActionRequest) (*http.Request, error)
1. ActionsClient.CreateOrUpdateResponder(*http.Response) (ActionResponse, error)
1. ActionsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ActionsClient.Delete(context.Context, string, string, string, string, string) (autorest.Response, error)
1. ActionsClient.DeletePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ActionsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ActionsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ActionsClient.Get(context.Context, string, string, string, string, string) (ActionResponse, error)
1. ActionsClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. ActionsClient.GetResponder(*http.Response) (ActionResponse, error)
1. ActionsClient.GetSender(*http.Request) (*http.Response, error)
1. AlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. AlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. AwsCloudTrailCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. AwsCloudTrailDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. AzureResourceEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. AzureResourceEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. AzureResourceEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. AzureResourceEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. CloudApplicationEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. CloudApplicationEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. CloudApplicationEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. CloudApplicationEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. DNSEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. DNSEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. DNSEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. DNSEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. DataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. DataConnectorsCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAADCheckRequirements() (*AADCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAATPCheckRequirements() (*AATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsASCCheckRequirements() (*ASCCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsAwsCloudTrailCheckRequirements() (*AwsCloudTrailCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsBasicDataConnectorsCheckRequirements() (BasicDataConnectorsCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsDataConnectorsCheckRequirements() (*DataConnectorsCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMCASCheckRequirements() (*MCASCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMDATPCheckRequirements() (*MDATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsOfficeATPCheckRequirements() (*OfficeATPCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsTICheckRequirements() (*TICheckRequirements, bool)
1. Dynamics365CheckRequirements.AsTiTaxiiCheckRequirements() (*TiTaxiiCheckRequirements, bool)
1. Dynamics365CheckRequirements.MarshalJSON() ([]byte, error)
1. Dynamics365DataConnector.AsAADDataConnector() (*AADDataConnector, bool)
1. Dynamics365DataConnector.AsAATPDataConnector() (*AATPDataConnector, bool)
1. Dynamics365DataConnector.AsASCDataConnector() (*ASCDataConnector, bool)
1. Dynamics365DataConnector.AsAwsCloudTrailDataConnector() (*AwsCloudTrailDataConnector, bool)
1. Dynamics365DataConnector.AsBasicDataConnector() (BasicDataConnector, bool)
1. Dynamics365DataConnector.AsDataConnector() (*DataConnector, bool)
1. Dynamics365DataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. Dynamics365DataConnector.AsMCASDataConnector() (*MCASDataConnector, bool)
1. Dynamics365DataConnector.AsMDATPDataConnector() (*MDATPDataConnector, bool)
1. Dynamics365DataConnector.AsOfficeATPDataConnector() (*OfficeATPDataConnector, bool)
1. Dynamics365DataConnector.AsOfficeDataConnector() (*OfficeDataConnector, bool)
1. Dynamics365DataConnector.AsTIDataConnector() (*TIDataConnector, bool)
1. Dynamics365DataConnector.AsTiTaxiiDataConnector() (*TiTaxiiDataConnector, bool)
1. Dynamics365DataConnector.MarshalJSON() ([]byte, error)
1. EntitiesClient.GetInsights(context.Context, string, string, string, string, EntityGetInsightsParameters) (EntityGetInsightsResponse, error)
1. EntitiesClient.GetInsightsPreparer(context.Context, string, string, string, string, EntityGetInsightsParameters) (*http.Request, error)
1. EntitiesClient.GetInsightsResponder(*http.Response) (EntityGetInsightsResponse, error)
1. EntitiesClient.GetInsightsSender(*http.Request) (*http.Response, error)
1. EntitiesClient.Queries(context.Context, string, string, string, string) (GetQueriesResponse, error)
1. EntitiesClient.QueriesPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. EntitiesClient.QueriesResponder(*http.Response) (GetQueriesResponse, error)
1. EntitiesClient.QueriesSender(*http.Request) (*http.Response, error)
1. Entity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. Entity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. Entity.AsMailboxEntity() (*MailboxEntity, bool)
1. Entity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. EntityQuery.AsBasicEntityQuery() (BasicEntityQuery, bool)
1. EntityQuery.AsEntityQuery() (*EntityQuery, bool)
1. EntityQuery.AsExpansionEntityQuery() (*ExpansionEntityQuery, bool)
1. EntityQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. EntityQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. EntityQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. EntityQueryItem.MarshalJSON() ([]byte, error)
1. ExpansionEntityQuery.AsBasicEntityQuery() (BasicEntityQuery, bool)
1. ExpansionEntityQuery.AsEntityQuery() (*EntityQuery, bool)
1. ExpansionEntityQuery.AsExpansionEntityQuery() (*ExpansionEntityQuery, bool)
1. ExpansionEntityQuery.MarshalJSON() ([]byte, error)
1. FileEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. FileEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. FileEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. FileEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. FileHashEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. FileHashEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. FileHashEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. FileHashEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. FusionAlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. FusionAlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. HostEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. HostEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. HostEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. HostEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. HuntingBookmark.AsMailClusterEntity() (*MailClusterEntity, bool)
1. HuntingBookmark.AsMailMessageEntity() (*MailMessageEntity, bool)
1. HuntingBookmark.AsMailboxEntity() (*MailboxEntity, bool)
1. HuntingBookmark.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. IPEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. IPEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. IPEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. IPEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. InsightQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. InsightQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. InsightQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. InsightQueryItem.MarshalJSON() ([]byte, error)
1. IoTDeviceEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. IoTDeviceEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. IoTDeviceEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. IoTDeviceEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. MCASCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MCASDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. MDATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MDATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. MLBehaviorAnalyticsAlertRule.AsAlertRule() (*AlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.AsBasicAlertRule() (BasicAlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.AsFusionAlertRule() (*FusionAlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.AsMicrosoftSecurityIncidentCreationAlertRule() (*MicrosoftSecurityIncidentCreationAlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.AsScheduledAlertRule() (*ScheduledAlertRule, bool)
1. MLBehaviorAnalyticsAlertRule.MarshalJSON() ([]byte, error)
1. MLBehaviorAnalyticsAlertRuleProperties.MarshalJSON() ([]byte, error)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsAlertRuleTemplate() (*AlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsBasicAlertRuleTemplate() (BasicAlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsFusionAlertRuleTemplate() (*FusionAlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsMicrosoftSecurityIncidentCreationAlertRuleTemplate() (*MicrosoftSecurityIncidentCreationAlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsScheduledAlertRuleTemplate() (*ScheduledAlertRuleTemplate, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.MarshalJSON() ([]byte, error)
1. MLBehaviorAnalyticsAlertRuleTemplateProperties.MarshalJSON() ([]byte, error)
1. MailClusterEntity.AsAccountEntity() (*AccountEntity, bool)
1. MailClusterEntity.AsAzureResourceEntity() (*AzureResourceEntity, bool)
1. MailClusterEntity.AsBasicEntity() (BasicEntity, bool)
1. MailClusterEntity.AsCloudApplicationEntity() (*CloudApplicationEntity, bool)
1. MailClusterEntity.AsDNSEntity() (*DNSEntity, bool)
1. MailClusterEntity.AsEntity() (*Entity, bool)
1. MailClusterEntity.AsFileEntity() (*FileEntity, bool)
1. MailClusterEntity.AsFileHashEntity() (*FileHashEntity, bool)
1. MailClusterEntity.AsHostEntity() (*HostEntity, bool)
1. MailClusterEntity.AsHuntingBookmark() (*HuntingBookmark, bool)
1. MailClusterEntity.AsIPEntity() (*IPEntity, bool)
1. MailClusterEntity.AsIoTDeviceEntity() (*IoTDeviceEntity, bool)
1. MailClusterEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. MailClusterEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. MailClusterEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. MailClusterEntity.AsMalwareEntity() (*MalwareEntity, bool)
1. MailClusterEntity.AsProcessEntity() (*ProcessEntity, bool)
1. MailClusterEntity.AsRegistryKeyEntity() (*RegistryKeyEntity, bool)
1. MailClusterEntity.AsRegistryValueEntity() (*RegistryValueEntity, bool)
1. MailClusterEntity.AsSecurityAlert() (*SecurityAlert, bool)
1. MailClusterEntity.AsSecurityGroupEntity() (*SecurityGroupEntity, bool)
1. MailClusterEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. MailClusterEntity.AsURLEntity() (*URLEntity, bool)
1. MailClusterEntity.MarshalJSON() ([]byte, error)
1. MailClusterEntityProperties.MarshalJSON() ([]byte, error)
1. MailMessageEntity.AsAccountEntity() (*AccountEntity, bool)
1. MailMessageEntity.AsAzureResourceEntity() (*AzureResourceEntity, bool)
1. MailMessageEntity.AsBasicEntity() (BasicEntity, bool)
1. MailMessageEntity.AsCloudApplicationEntity() (*CloudApplicationEntity, bool)
1. MailMessageEntity.AsDNSEntity() (*DNSEntity, bool)
1. MailMessageEntity.AsEntity() (*Entity, bool)
1. MailMessageEntity.AsFileEntity() (*FileEntity, bool)
1. MailMessageEntity.AsFileHashEntity() (*FileHashEntity, bool)
1. MailMessageEntity.AsHostEntity() (*HostEntity, bool)
1. MailMessageEntity.AsHuntingBookmark() (*HuntingBookmark, bool)
1. MailMessageEntity.AsIPEntity() (*IPEntity, bool)
1. MailMessageEntity.AsIoTDeviceEntity() (*IoTDeviceEntity, bool)
1. MailMessageEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. MailMessageEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. MailMessageEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. MailMessageEntity.AsMalwareEntity() (*MalwareEntity, bool)
1. MailMessageEntity.AsProcessEntity() (*ProcessEntity, bool)
1. MailMessageEntity.AsRegistryKeyEntity() (*RegistryKeyEntity, bool)
1. MailMessageEntity.AsRegistryValueEntity() (*RegistryValueEntity, bool)
1. MailMessageEntity.AsSecurityAlert() (*SecurityAlert, bool)
1. MailMessageEntity.AsSecurityGroupEntity() (*SecurityGroupEntity, bool)
1. MailMessageEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. MailMessageEntity.AsURLEntity() (*URLEntity, bool)
1. MailMessageEntity.MarshalJSON() ([]byte, error)
1. MailMessageEntityProperties.MarshalJSON() ([]byte, error)
1. MailboxEntity.AsAccountEntity() (*AccountEntity, bool)
1. MailboxEntity.AsAzureResourceEntity() (*AzureResourceEntity, bool)
1. MailboxEntity.AsBasicEntity() (BasicEntity, bool)
1. MailboxEntity.AsCloudApplicationEntity() (*CloudApplicationEntity, bool)
1. MailboxEntity.AsDNSEntity() (*DNSEntity, bool)
1. MailboxEntity.AsEntity() (*Entity, bool)
1. MailboxEntity.AsFileEntity() (*FileEntity, bool)
1. MailboxEntity.AsFileHashEntity() (*FileHashEntity, bool)
1. MailboxEntity.AsHostEntity() (*HostEntity, bool)
1. MailboxEntity.AsHuntingBookmark() (*HuntingBookmark, bool)
1. MailboxEntity.AsIPEntity() (*IPEntity, bool)
1. MailboxEntity.AsIoTDeviceEntity() (*IoTDeviceEntity, bool)
1. MailboxEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. MailboxEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. MailboxEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. MailboxEntity.AsMalwareEntity() (*MalwareEntity, bool)
1. MailboxEntity.AsProcessEntity() (*ProcessEntity, bool)
1. MailboxEntity.AsRegistryKeyEntity() (*RegistryKeyEntity, bool)
1. MailboxEntity.AsRegistryValueEntity() (*RegistryValueEntity, bool)
1. MailboxEntity.AsSecurityAlert() (*SecurityAlert, bool)
1. MailboxEntity.AsSecurityGroupEntity() (*SecurityGroupEntity, bool)
1. MailboxEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. MailboxEntity.AsURLEntity() (*URLEntity, bool)
1. MailboxEntity.MarshalJSON() ([]byte, error)
1. MailboxEntityProperties.MarshalJSON() ([]byte, error)
1. MalwareEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. MalwareEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. MalwareEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. MalwareEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. MicrosoftSecurityIncidentCreationAlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. MicrosoftSecurityIncidentCreationAlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. NewWatchlistItemClient(string) WatchlistItemClient
1. NewWatchlistItemClientWithBaseURI(string, string) WatchlistItemClient
1. OfficeATPCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. OfficeATPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. OfficeDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. PossibleAntispamMailDirectionValues() []AntispamMailDirection
1. PossibleDeliveryActionValues() []DeliveryAction
1. PossibleDeliveryLocationValues() []DeliveryLocation
1. PossibleEntityQueryKindValues() []EntityQueryKind
1. PossibleKindBasicEntityQueryItemValues() []KindBasicEntityQueryItem
1. PossibleKindBasicEntityQueryValues() []KindBasicEntityQuery
1. PossibleOutputTypeValues() []OutputType
1. ProcessEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. ProcessEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. ProcessEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. ProcessEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. RegistryKeyEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. RegistryKeyEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. RegistryKeyEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. RegistryKeyEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. RegistryValueEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. RegistryValueEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. RegistryValueEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. RegistryValueEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. ScheduledAlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. ScheduledAlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. SecurityAlert.AsMailClusterEntity() (*MailClusterEntity, bool)
1. SecurityAlert.AsMailMessageEntity() (*MailMessageEntity, bool)
1. SecurityAlert.AsMailboxEntity() (*MailboxEntity, bool)
1. SecurityAlert.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. SecurityGroupEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. SecurityGroupEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. SecurityGroupEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. SecurityGroupEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. SubmissionMailEntity.AsAccountEntity() (*AccountEntity, bool)
1. SubmissionMailEntity.AsAzureResourceEntity() (*AzureResourceEntity, bool)
1. SubmissionMailEntity.AsBasicEntity() (BasicEntity, bool)
1. SubmissionMailEntity.AsCloudApplicationEntity() (*CloudApplicationEntity, bool)
1. SubmissionMailEntity.AsDNSEntity() (*DNSEntity, bool)
1. SubmissionMailEntity.AsEntity() (*Entity, bool)
1. SubmissionMailEntity.AsFileEntity() (*FileEntity, bool)
1. SubmissionMailEntity.AsFileHashEntity() (*FileHashEntity, bool)
1. SubmissionMailEntity.AsHostEntity() (*HostEntity, bool)
1. SubmissionMailEntity.AsHuntingBookmark() (*HuntingBookmark, bool)
1. SubmissionMailEntity.AsIPEntity() (*IPEntity, bool)
1. SubmissionMailEntity.AsIoTDeviceEntity() (*IoTDeviceEntity, bool)
1. SubmissionMailEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. SubmissionMailEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. SubmissionMailEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. SubmissionMailEntity.AsMalwareEntity() (*MalwareEntity, bool)
1. SubmissionMailEntity.AsProcessEntity() (*ProcessEntity, bool)
1. SubmissionMailEntity.AsRegistryKeyEntity() (*RegistryKeyEntity, bool)
1. SubmissionMailEntity.AsRegistryValueEntity() (*RegistryValueEntity, bool)
1. SubmissionMailEntity.AsSecurityAlert() (*SecurityAlert, bool)
1. SubmissionMailEntity.AsSecurityGroupEntity() (*SecurityGroupEntity, bool)
1. SubmissionMailEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. SubmissionMailEntity.AsURLEntity() (*URLEntity, bool)
1. SubmissionMailEntity.MarshalJSON() ([]byte, error)
1. SubmissionMailEntityProperties.MarshalJSON() ([]byte, error)
1. TICheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. TIDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. TiTaxiiCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. TiTaxiiDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. URLEntity.AsMailClusterEntity() (*MailClusterEntity, bool)
1. URLEntity.AsMailMessageEntity() (*MailMessageEntity, bool)
1. URLEntity.AsMailboxEntity() (*MailboxEntity, bool)
1. URLEntity.AsSubmissionMailEntity() (*SubmissionMailEntity, bool)
1. WatchlistItem.MarshalJSON() ([]byte, error)
1. WatchlistItemClient.CreateOrUpdate(context.Context, string, string, string, string, string, WatchlistItem) (WatchlistItem, error)
1. WatchlistItemClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, string, WatchlistItem) (*http.Request, error)
1. WatchlistItemClient.CreateOrUpdateResponder(*http.Response) (WatchlistItem, error)
1. WatchlistItemClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. WatchlistItemClient.Delete(context.Context, string, string, string, string, string) (autorest.Response, error)
1. WatchlistItemClient.DeletePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. WatchlistItemClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. WatchlistItemClient.DeleteSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. Dynamics365CheckRequirements
1. Dynamics365CheckRequirementsProperties
1. Dynamics365DataConnector
1. Dynamics365DataConnectorDataTypes
1. Dynamics365DataConnectorDataTypesDynamics365CdsActivities
1. Dynamics365DataConnectorProperties
1. EntityGetInsightsParameters
1. EntityGetInsightsResponse
1. EntityInsightItem
1. EntityInsightItemQueryTimeInterval
1. EntityQueryItem
1. EntityQueryItemProperties
1. EntityQueryItemPropertiesDataTypesItem
1. EntityQueryKind1
1. EntityQueryModel
1. ExpansionEntityQueriesProperties
1. ExpansionEntityQuery
1. GetInsightsError
1. GetInsightsResultsMetadata
1. GetQueriesResponse
1. InsightQueryItem
1. InsightQueryItemProperties
1. InsightQueryItemPropertiesAdditionalQuery
1. InsightQueryItemPropertiesDefaultTimeRange
1. InsightQueryItemPropertiesReferenceTimeRange
1. InsightQueryItemPropertiesTableQuery
1. InsightQueryItemPropertiesTableQueryColumnsDefinitionsItem
1. InsightQueryItemPropertiesTableQueryQueriesDefinitionsItem
1. InsightQueryItemPropertiesTableQueryQueriesDefinitionsItemLinkColumnsDefinitionsItem
1. InsightsTableResult
1. InsightsTableResultColumnsItem
1. MLBehaviorAnalyticsAlertRule
1. MLBehaviorAnalyticsAlertRuleProperties
1. MLBehaviorAnalyticsAlertRuleTemplate
1. MLBehaviorAnalyticsAlertRuleTemplateProperties
1. MailClusterEntity
1. MailClusterEntityProperties
1. MailMessageEntity
1. MailMessageEntityProperties
1. MailboxEntity
1. MailboxEntityProperties
1. SubmissionMailEntity
1. SubmissionMailEntityProperties
1. WatchlistItem
1. WatchlistItemClient
1. WatchlistItemProperties

### New Struct Fields

1. AlertRuleTemplatePropertiesBase.LastUpdatedDateUTC
1. CaseProperties.Metrics
1. CaseProperties.RelatedAlertProductNames
1. CasesAggregationByStatusProperties.TotalFalsePositiveStatus
1. CasesAggregationByStatusProperties.TotalTruePositiveStatus
1. EntityQuery.Etag
1. EntityQuery.Kind
1. FusionAlertRuleTemplateProperties.LastUpdatedDateUTC
1. IncidentProperties.ProviderIncidentID
1. IncidentProperties.ProviderName
1. IoTDeviceEntityProperties.DeviceName
1. IoTDeviceEntityProperties.FirmwareVersion
1. IoTDeviceEntityProperties.IPAddressEntityID
1. IoTDeviceEntityProperties.MacAddress
1. IoTDeviceEntityProperties.Model
1. IoTDeviceEntityProperties.OperatingSystem
1. IoTDeviceEntityProperties.Protocols
1. IoTDeviceEntityProperties.SerialNumber
1. IoTDeviceEntityProperties.Source
1. MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties.LastUpdatedDateUTC
1. ScheduledAlertRuleTemplateProperties.LastUpdatedDateUTC
1. WatchlistProperties.UploadStatus
1. WatchlistProperties.WatchlistItemsCount
