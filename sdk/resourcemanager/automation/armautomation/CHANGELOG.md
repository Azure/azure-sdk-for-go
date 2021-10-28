# Release History

## 0.2.0 (2021-10-28)
### Breaking Changes

- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAgentRegistrationInformationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewLinkedWorkspaceClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewObjectDataTypesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAutomationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUsagesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConnectionTypeClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDscCompilationJobClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewKeysClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPython2PackageClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSoftwareUpdateConfigurationRunsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewModuleClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAutomationAccountClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewWatcherClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewJobClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNodeCountInformationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewActivityClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewHybridRunbookWorkerGroupClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSourceControlSyncJobStreamsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSourceControlSyncJobClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewJobScheduleClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewHybridRunbookWorkersClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewWebhookClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewStatisticsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSoftwareUpdateConfigurationMachineRunsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSoftwareUpdateConfigurationsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNodeReportsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSourceControlClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTestJobClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDscConfigurationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewFieldsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewScheduleClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRunbookClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRunbookDraftClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDscNodeClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConnectionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDscCompilationJobStreamClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCertificateClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDscNodeConfigurationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewJobStreamClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewVariableClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCredentialClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTestJobStreamsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 45 breaking change(s), 0 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-15)

- Initial preview release.
