
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewAlertRulesListPage` signature has been changed from `(func(context.Context, AlertRulesList) (AlertRulesList, error))` to `(AlertRulesList,func(context.Context, AlertRulesList) (AlertRulesList, error))`
- Function `NewOperationsListPage` signature has been changed from `(func(context.Context, OperationsList) (OperationsList, error))` to `(OperationsList,func(context.Context, OperationsList) (OperationsList, error))`
- Function `NewAlertsListPage` signature has been changed from `(func(context.Context, AlertsList) (AlertsList, error))` to `(AlertsList,func(context.Context, AlertsList) (AlertsList, error))`
- Function `NewSmartGroupsClientWithBaseURI` has been removed
- Function `SmartGroupsClient.ChangeStateResponder` has been removed
- Function `SmartGroupsClient.GetHistoryResponder` has been removed
- Function `SmartGroupsClient.GetHistoryPreparer` has been removed
- Function `SmartGroupsClient.GetByIDSender` has been removed
- Function `SmartGroupsClient.ChangeStateSender` has been removed
- Function `SmartGroupsClient.GetAll` has been removed
- Function `SmartGroupsClient.GetByIDPreparer` has been removed
- Function `SmartGroupsClient.GetByIDResponder` has been removed
- Function `SmartGroupsClient.GetAllResponder` has been removed
- Function `SmartGroupsClient.ChangeStatePreparer` has been removed
- Function `SmartGroupsClient.GetHistory` has been removed
- Function `SmartGroupsClient.GetAllSender` has been removed
- Function `SmartGroupsClient.GetHistorySender` has been removed
- Function `NewSmartGroupsClient` has been removed
- Function `SmartGroupsClient.GetByID` has been removed
- Function `SmartGroupsClient.ChangeState` has been removed
- Function `SmartGroupsClient.GetAllPreparer` has been removed
- Struct `SmartGroupsClient` has been removed
- Field `autorest.Response` of struct `SmartGroupModification` has been removed
- Field `autorest.Response` of struct `SmartGroupsList` has been removed
- Field `autorest.Response` of struct `SmartGroup` has been removed

## New Content

- Const `MetadataIdentifierMonitorServiceList` is added
- Const `MetadataIdentifierAlertsMetaDataProperties` is added
- Function `AlertsClient.MetaDataResponder(*http.Response) (AlertsMetaData,error)` is added
- Function `AlertsMetaDataProperties.AsAlertsMetaDataProperties() (*AlertsMetaDataProperties,bool)` is added
- Function `MonitorServiceList.AsAlertsMetaDataProperties() (*AlertsMetaDataProperties,bool)` is added
- Function `AlertsMetaDataProperties.AsMonitorServiceList() (*MonitorServiceList,bool)` is added
- Function `PossibleMetadataIdentifierValues() []MetadataIdentifier` is added
- Function `MonitorServiceList.AsMonitorServiceList() (*MonitorServiceList,bool)` is added
- Function `MonitorServiceList.MarshalJSON() ([]byte,error)` is added
- Function `AlertsClient.MetaDataSender(*http.Request) (*http.Response,error)` is added
- Function `AlertsMetaDataProperties.MarshalJSON() ([]byte,error)` is added
- Function `AlertsClient.MetaDataPreparer(context.Context) (*http.Request,error)` is added
- Function `*AlertsMetaData.UnmarshalJSON([]byte) error` is added
- Function `AlertsClient.MetaData(context.Context) (AlertsMetaData,error)` is added
- Function `MonitorServiceList.AsBasicAlertsMetaDataProperties() (BasicAlertsMetaDataProperties,bool)` is added
- Function `AlertsMetaDataProperties.AsBasicAlertsMetaDataProperties() (BasicAlertsMetaDataProperties,bool)` is added
- Struct `AlertsMetaData` is added
- Struct `AlertsMetaDataProperties` is added
- Struct `MonitorServiceDetails` is added
- Struct `MonitorServiceList` is added

