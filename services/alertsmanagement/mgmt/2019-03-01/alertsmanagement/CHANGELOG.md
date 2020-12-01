Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewAlertsListPage` parameter(s) have been changed from `(func(context.Context, AlertsList) (AlertsList, error))` to `(AlertsList, func(context.Context, AlertsList) (AlertsList, error))`
- Function `NewAlertRulesListPage` parameter(s) have been changed from `(func(context.Context, AlertRulesList) (AlertRulesList, error))` to `(AlertRulesList, func(context.Context, AlertRulesList) (AlertRulesList, error))`
- Function `NewOperationsListPage` parameter(s) have been changed from `(func(context.Context, OperationsList) (OperationsList, error))` to `(OperationsList, func(context.Context, OperationsList) (OperationsList, error))`
- Function `SmartGroupsClient.ChangeState` has been removed
- Function `SmartGroupsClient.GetByIDResponder` has been removed
- Function `SmartGroupsClient.GetHistoryPreparer` has been removed
- Function `SmartGroupsClient.GetAllSender` has been removed
- Function `SmartGroupsClient.ChangeStateSender` has been removed
- Function `SmartGroupsClient.GetAll` has been removed
- Function `SmartGroupsClient.GetHistorySender` has been removed
- Function `SmartGroupsClient.GetHistory` has been removed
- Function `SmartGroupsClient.ChangeStateResponder` has been removed
- Function `SmartGroupsClient.GetByIDSender` has been removed
- Function `SmartGroupsClient.GetHistoryResponder` has been removed
- Function `SmartGroupsClient.GetByID` has been removed
- Function `NewSmartGroupsClientWithBaseURI` has been removed
- Function `SmartGroupsClient.GetAllResponder` has been removed
- Function `SmartGroupsClient.ChangeStatePreparer` has been removed
- Function `SmartGroupsClient.GetAllPreparer` has been removed
- Function `SmartGroupsClient.GetByIDPreparer` has been removed
- Function `NewSmartGroupsClient` has been removed
- Struct `SmartGroupsClient` has been removed
- Field `autorest.Response` of struct `SmartGroupsList` has been removed
- Field `autorest.Response` of struct `SmartGroup` has been removed
- Field `autorest.Response` of struct `SmartGroupModification` has been removed

## New Content

- New const `MetadataIdentifierMonitorServiceList`
- New const `MetadataIdentifierAlertsMetaDataProperties`
- New function `AlertsMetaDataProperties.MarshalJSON() ([]byte, error)`
- New function `AlertsClient.MetaDataResponder(*http.Response) (AlertsMetaData, error)`
- New function `*AlertsMetaData.UnmarshalJSON([]byte) error`
- New function `MonitorServiceList.AsMonitorServiceList() (*MonitorServiceList, bool)`
- New function `AlertsMetaDataProperties.AsMonitorServiceList() (*MonitorServiceList, bool)`
- New function `MonitorServiceList.MarshalJSON() ([]byte, error)`
- New function `AlertsClient.MetaDataPreparer(context.Context) (*http.Request, error)`
- New function `AlertsMetaDataProperties.AsAlertsMetaDataProperties() (*AlertsMetaDataProperties, bool)`
- New function `AlertsClient.MetaData(context.Context) (AlertsMetaData, error)`
- New function `PossibleMetadataIdentifierValues() []MetadataIdentifier`
- New function `MonitorServiceList.AsBasicAlertsMetaDataProperties() (BasicAlertsMetaDataProperties, bool)`
- New function `AlertsClient.MetaDataSender(*http.Request) (*http.Response, error)`
- New function `MonitorServiceList.AsAlertsMetaDataProperties() (*AlertsMetaDataProperties, bool)`
- New function `AlertsMetaDataProperties.AsBasicAlertsMetaDataProperties() (BasicAlertsMetaDataProperties, bool)`
- New struct `AlertsMetaData`
- New struct `AlertsMetaDataProperties`
- New struct `MonitorServiceDetails`
- New struct `MonitorServiceList`
