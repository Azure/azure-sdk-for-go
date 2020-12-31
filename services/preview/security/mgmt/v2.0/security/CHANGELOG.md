Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Constants

1. EventSource.Alerts
1. EventSource.Assessments
1. EventSource.SecureScoreControls
1. EventSource.SecureScores
1. EventSource.SubAssessments
1. SettingKind.SettingKindAlertSuppressionSetting
1. SettingKind.SettingKindDataExportSetting

### Removed Funcs

1. *AdaptiveNetworkHardeningsEnforceFuture.Result(AdaptiveNetworkHardeningsClient) (autorest.Response, error)
1. PossibleSettingKindValues() []SettingKind
1. SettingResource.MarshalJSON() ([]byte, error)

## Struct Changes

### Removed Structs

1. SettingResource

### Removed Struct Fields

1. AdaptiveNetworkHardeningsEnforceFuture.azure.Future

## Signature Changes

### Const Types

1. KindAAD changed type from KindEnum to KindEnum1
1. KindATA changed type from KindEnum to KindEnum1
1. KindCEF changed type from KindEnum to KindEnum1
1. KindExternalSecuritySolution changed type from KindEnum to KindEnum1

### Funcs

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

### Struct Fields

1. AadExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. AtaExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. CefExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. DataExportSetting.Kind changed type from SettingKind to KindEnum
1. ExternalSecuritySolution.Kind changed type from KindEnum to KindEnum1
1. Setting.Kind changed type from SettingKind to KindEnum
1. SettingsList.Value changed type from *[]Setting to *[]BasicSetting

### New Constants

1. EventSource.EventSourceAlerts
1. EventSource.EventSourceAssessments
1. EventSource.EventSourceRegulatoryComplianceAssessment
1. EventSource.EventSourceSecureScoreControls
1. EventSource.EventSourceSecureScores
1. EventSource.EventSourceSubAssessments
1. KindEnum.KindDataExportSetting
1. KindEnum.KindSetting

### New Funcs

1. *SettingModel.UnmarshalJSON([]byte) error
1. *SettingsList.UnmarshalJSON([]byte) error
1. DataExportSetting.AsBasicSetting() (BasicSetting, bool)
1. DataExportSetting.AsDataExportSetting() (*DataExportSetting, bool)
1. DataExportSetting.AsSetting() (*Setting, bool)
1. PossibleKindEnum1Values() []KindEnum1
1. Setting.AsBasicSetting() (BasicSetting, bool)
1. Setting.AsDataExportSetting() (*DataExportSetting, bool)
1. Setting.AsSetting() (*Setting, bool)

## Struct Changes

### New Structs

1. SettingModel

### New Struct Fields

1. AdaptiveNetworkHardeningsEnforceFuture.Result
1. AdaptiveNetworkHardeningsEnforceFuture.azure.FutureAPI
