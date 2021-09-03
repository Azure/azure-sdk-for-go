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

### Signature Changes

#### Const Types

1. KindAAD changed type from KindEnum to KindEnum1
1. KindATA changed type from KindEnum to KindEnum1
1. KindCEF changed type from KindEnum to KindEnum1
1. KindExternalSecuritySolution changed type from KindEnum to KindEnum1
1. None changed type from EnforcementMode to EndOfSupportStatus

#### Funcs

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

#### Struct Fields

1. AadExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. AtaExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. CefExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. DataExportSetting.Kind changed type from SettingKind to KindEnum
1. ExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. Setting.Kind changed type from SettingKind to KindEnum
1. SettingsList.Value changed type from *[]Setting to *[]BasicSetting

## Additive Changes

### New Constants

1. EndOfSupportStatus.NoLongerSupported
1. EndOfSupportStatus.UpcomingNoLongerSupported
1. EndOfSupportStatus.UpcomingVersionNoLongerSupported
1. EndOfSupportStatus.VersionNoLongerSupported
1. EnforcementMode.EnforcementModeAudit
1. EnforcementMode.EnforcementModeEnforce
1. EnforcementMode.EnforcementModeNone
1. EventSource.EventSourceAlerts
1. EventSource.EventSourceAssessments
1. EventSource.EventSourceRegulatoryComplianceAssessment
1. EventSource.EventSourceRegulatoryComplianceAssessmentSnapshot
1. EventSource.EventSourceSecureScoreControls
1. EventSource.EventSourceSecureScoreControlsSnapshot
1. EventSource.EventSourceSecureScores
1. EventSource.EventSourceSecureScoresSnapshot
1. EventSource.EventSourceSubAssessments
1. KindEnum.KindDataExportSetting
1. KindEnum.KindSetting

### New Funcs

1. *IngestionSettingListIterator.Next() error
1. *IngestionSettingListIterator.NextWithContext(context.Context) error
1. *IngestionSettingListPage.Next() error
1. *IngestionSettingListPage.NextWithContext(context.Context) error
1. *SettingModel.UnmarshalJSON([]byte) error
1. *SettingsList.UnmarshalJSON([]byte) error
1. *Software.UnmarshalJSON([]byte) error
1. *SoftwaresListIterator.Next() error
1. *SoftwaresListIterator.NextWithContext(context.Context) error
1. *SoftwaresListPage.Next() error
1. *SoftwaresListPage.NextWithContext(context.Context) error
1. DataExportSetting.AsBasicSetting() (BasicSetting, bool)
1. DataExportSetting.AsDataExportSetting() (*DataExportSetting, bool)
1. DataExportSetting.AsSetting() (*Setting, bool)
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
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
1. NewIngestionSettingListIterator(IngestionSettingListPage) IngestionSettingListIterator
1. NewIngestionSettingListPage(IngestionSettingList, func(context.Context, IngestionSettingList) (IngestionSettingList, error)) IngestionSettingListPage
1. NewIngestionSettingsClient(string, string) IngestionSettingsClient
1. NewIngestionSettingsClientWithBaseURI(string, string, string) IngestionSettingsClient
1. NewSoftwareInventoriesClient(string, string) SoftwareInventoriesClient
1. NewSoftwareInventoriesClientWithBaseURI(string, string, string) SoftwareInventoriesClient
1. NewSoftwaresListIterator(SoftwaresListPage) SoftwaresListIterator
1. NewSoftwaresListPage(SoftwaresList, func(context.Context, SoftwaresList) (SoftwaresList, error)) SoftwaresListPage
1. PossibleEndOfSupportStatusValues() []EndOfSupportStatus
1. PossibleKindEnum1Values() []KindEnum1
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

1. ConnectionStrings
1. ErrorAdditionalInfo
1. IngestionConnectionString
1. IngestionSetting
1. IngestionSettingList
1. IngestionSettingListIterator
1. IngestionSettingListPage
1. IngestionSettingToken
1. IngestionSettingsClient
1. SettingModel
1. Software
1. SoftwareInventoriesClient
1. SoftwareProperties
1. SoftwaresList
1. SoftwaresListIterator
1. SoftwaresListPage

#### New Struct Fields

1. CloudErrorBody.AdditionalInfo
1. CloudErrorBody.Details
1. CloudErrorBody.Target
