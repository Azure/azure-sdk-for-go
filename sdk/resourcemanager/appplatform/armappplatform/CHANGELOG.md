# Release History

## 0.2.0 (2021-10-28)
### Breaking Changes

- Function `NewDeploymentsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewMonitoringSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewBindingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAppsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewServicesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCertificatesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCustomDomainsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRuntimeVersionsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSKUsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConfigServersClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Type of `CertificateResource.Properties` has been changed from `*CertificateProperties` to `CertificatePropertiesClassification`
- Function `CertificateProperties.MarshalJSON` has been removed
- Field `VaultURI` of struct `CertificateProperties` has been removed
- Field `CertVersion` of struct `CertificateProperties` has been removed
- Field `KeyVaultCertName` of struct `CertificateProperties` has been removed

### New Content

- New const `CustomPersistentDiskPropertiesTypeAzureFileVolume`
- New const `PowerStateRunning`
- New const `CreatedByTypeApplication`
- New const `CreatedByTypeUser`
- New const `CreatedByTypeManagedIdentity`
- New const `StoragePropertiesStorageTypeStorageAccount`
- New const `CreatedByTypeKey`
- New const `PowerStateStopped`
- New function `CustomPersistentDiskResource.MarshalJSON() ([]byte, error)`
- New function `*DeploymentsGenerateHeapDumpPoller.ResumeToken() (string, error)`
- New function `*CertificateProperties.UnmarshalJSON([]byte) error`
- New function `*ServicesStopPoller.Poll(context.Context) (*http.Response, error)`
- New function `*DeploymentsClient.BeginGenerateHeapDump(context.Context, string, string, string, string, DiagnosticParameters, *DeploymentsBeginGenerateHeapDumpOptions) (DeploymentsGenerateHeapDumpPollerResponse, error)`
- New function `*StoragesListPager.PageResponse() StoragesListResponse`
- New function `*ServicesStopPoller.FinalResponse(context.Context) (ServicesStopResponse, error)`
- New function `*DeploymentsStartJFRPoller.Poll(context.Context) (*http.Response, error)`
- New function `*StoragesDeletePollerResponse.Resume(context.Context, *StoragesClient, string) error`
- New function `*ServicesStopPollerResponse.Resume(context.Context, *ServicesClient, string) error`
- New function `*StoragesCreateOrUpdatePoller.Done() bool`
- New function `DeploymentsGenerateHeapDumpPollerResponse.PollUntilDone(context.Context, time.Duration) (DeploymentsGenerateHeapDumpResponse, error)`
- New function `*StoragesDeletePoller.Done() bool`
- New function `*StoragesCreateOrUpdatePollerResponse.Resume(context.Context, *StoragesClient, string) error`
- New function `*KeyVaultCertificateProperties.UnmarshalJSON([]byte) error`
- New function `*DeploymentsGenerateHeapDumpPoller.Poll(context.Context) (*http.Response, error)`
- New function `ServicesStopPollerResponse.PollUntilDone(context.Context, time.Duration) (ServicesStopResponse, error)`
- New function `*DeploymentsGenerateHeapDumpPoller.Done() bool`
- New function `*ServicesClient.BeginStart(context.Context, string, string, *ServicesBeginStartOptions) (ServicesStartPollerResponse, error)`
- New function `StoragesDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (StoragesDeleteResponse, error)`
- New function `*DeploymentsGenerateThreadDumpPoller.ResumeToken() (string, error)`
- New function `*TrackedResource.UnmarshalJSON([]byte) error`
- New function `*BindingResource.UnmarshalJSON([]byte) error`
- New function `ContentCertificateProperties.MarshalJSON() ([]byte, error)`
- New function `*DeploymentsClient.BeginStartJFR(context.Context, string, string, string, string, DiagnosticParameters, *DeploymentsBeginStartJFROptions) (DeploymentsStartJFRPollerResponse, error)`
- New function `ServicesStartPollerResponse.PollUntilDone(context.Context, time.Duration) (ServicesStartResponse, error)`
- New function `*ServicesClient.BeginStop(context.Context, string, string, *ServicesBeginStopOptions) (ServicesStopPollerResponse, error)`
- New function `*Resource.UnmarshalJSON([]byte) error`
- New function `*CustomPersistentDiskResource.UnmarshalJSON([]byte) error`
- New function `CreatedByType.ToPtr() *CreatedByType`
- New function `PowerState.ToPtr() *PowerState`
- New function `*StoragesClient.List(string, string, *StoragesListOptions) *StoragesListPager`
- New function `*CustomPersistentDiskProperties.GetCustomPersistentDiskProperties() *CustomPersistentDiskProperties`
- New function `DeploymentsStartJFRPollerResponse.PollUntilDone(context.Context, time.Duration) (DeploymentsStartJFRResponse, error)`
- New function `KeyVaultCertificateProperties.MarshalJSON() ([]byte, error)`
- New function `CustomPersistentDiskPropertiesType.ToPtr() *CustomPersistentDiskPropertiesType`
- New function `*CustomPersistentDiskProperties.UnmarshalJSON([]byte) error`
- New function `*MonitoringSettingResource.UnmarshalJSON([]byte) error`
- New function `*StoragesDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New function `*ServicesStartPoller.FinalResponse(context.Context) (ServicesStartResponse, error)`
- New function `*DeploymentsGenerateHeapDumpPoller.FinalResponse(context.Context) (DeploymentsGenerateHeapDumpResponse, error)`
- New function `NewStoragesClient(string, azcore.TokenCredential, *arm.ClientOptions) *StoragesClient`
- New function `*DeploymentsGenerateThreadDumpPoller.Done() bool`
- New function `StoragePropertiesStorageType.ToPtr() *StoragePropertiesStorageType`
- New function `*StoragesCreateOrUpdatePoller.FinalResponse(context.Context) (StoragesCreateOrUpdateResponse, error)`
- New function `*StoragesListPager.Err() error`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `*CertificateResource.UnmarshalJSON([]byte) error`
- New function `*StorageResource.UnmarshalJSON([]byte) error`
- New function `*DeploymentsGenerateThreadDumpPollerResponse.Resume(context.Context, *DeploymentsClient, string) error`
- New function `StorageAccount.MarshalJSON() ([]byte, error)`
- New function `*ConfigServerResource.UnmarshalJSON([]byte) error`
- New function `*ServiceResource.UnmarshalJSON([]byte) error`
- New function `*CertificateProperties.GetCertificateProperties() *CertificateProperties`
- New function `*StoragesClient.BeginDelete(context.Context, string, string, string, *StoragesBeginDeleteOptions) (StoragesDeletePollerResponse, error)`
- New function `*CustomDomainResource.UnmarshalJSON([]byte) error`
- New function `*DeploymentsClient.BeginGenerateThreadDump(context.Context, string, string, string, string, DiagnosticParameters, *DeploymentsBeginGenerateThreadDumpOptions) (DeploymentsGenerateThreadDumpPollerResponse, error)`
- New function `*AzureFileVolume.UnmarshalJSON([]byte) error`
- New function `*DeploymentResource.UnmarshalJSON([]byte) error`
- New function `StorageResourceCollection.MarshalJSON() ([]byte, error)`
- New function `*StoragesDeletePoller.FinalResponse(context.Context) (StoragesDeleteResponse, error)`
- New function `*StorageProperties.GetStorageProperties() *StorageProperties`
- New function `*StoragesClient.Get(context.Context, string, string, string, *StoragesGetOptions) (StoragesGetResponse, error)`
- New function `*StoragesDeletePoller.ResumeToken() (string, error)`
- New function `*DeploymentsGenerateThreadDumpPoller.FinalResponse(context.Context) (DeploymentsGenerateThreadDumpResponse, error)`
- New function `*ServicesStartPoller.ResumeToken() (string, error)`
- New function `StorageResource.MarshalJSON() ([]byte, error)`
- New function `*ServicesStartPollerResponse.Resume(context.Context, *ServicesClient, string) error`
- New function `AzureFileVolume.MarshalJSON() ([]byte, error)`
- New function `*ServicesStartPoller.Done() bool`
- New function `*DeploymentsStartJFRPoller.FinalResponse(context.Context) (DeploymentsStartJFRResponse, error)`
- New function `*DeploymentsStartJFRPoller.ResumeToken() (string, error)`
- New function `*DeploymentsStartJFRPoller.Done() bool`
- New function `*StoragesListPager.NextPage(context.Context) bool`
- New function `*StoragesCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*DeploymentsGenerateHeapDumpPollerResponse.Resume(context.Context, *DeploymentsClient, string) error`
- New function `*DeploymentsGenerateThreadDumpPoller.Poll(context.Context) (*http.Response, error)`
- New function `*ServicesStopPoller.Done() bool`
- New function `PossibleStoragePropertiesStorageTypeValues() []StoragePropertiesStorageType`
- New function `*StorageAccount.UnmarshalJSON([]byte) error`
- New function `*AppResource.UnmarshalJSON([]byte) error`
- New function `DeploymentsGenerateThreadDumpPollerResponse.PollUntilDone(context.Context, time.Duration) (DeploymentsGenerateThreadDumpResponse, error)`
- New function `*ServicesStartPoller.Poll(context.Context) (*http.Response, error)`
- New function `*ServicesStopPoller.ResumeToken() (string, error)`
- New function `*ContentCertificateProperties.UnmarshalJSON([]byte) error`
- New function `PossiblePowerStateValues() []PowerState`
- New function `StoragesCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (StoragesCreateOrUpdateResponse, error)`
- New function `*StoragesCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*StorageProperties.UnmarshalJSON([]byte) error`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `*StoragesClient.BeginCreateOrUpdate(context.Context, string, string, string, StorageResource, *StoragesBeginCreateOrUpdateOptions) (StoragesCreateOrUpdatePollerResponse, error)`
- New function `*DeploymentsStartJFRPollerResponse.Resume(context.Context, *DeploymentsClient, string) error`
- New function `PossibleCustomPersistentDiskPropertiesTypeValues() []CustomPersistentDiskPropertiesType`
- New struct `AzureFileVolume`
- New struct `ContentCertificateProperties`
- New struct `CustomPersistentDiskProperties`
- New struct `CustomPersistentDiskResource`
- New struct `DeploymentSettingsContainerProbeSettings`
- New struct `DeploymentsBeginGenerateHeapDumpOptions`
- New struct `DeploymentsBeginGenerateThreadDumpOptions`
- New struct `DeploymentsBeginStartJFROptions`
- New struct `DeploymentsGenerateHeapDumpPoller`
- New struct `DeploymentsGenerateHeapDumpPollerResponse`
- New struct `DeploymentsGenerateHeapDumpResponse`
- New struct `DeploymentsGenerateThreadDumpPoller`
- New struct `DeploymentsGenerateThreadDumpPollerResponse`
- New struct `DeploymentsGenerateThreadDumpResponse`
- New struct `DeploymentsStartJFRPoller`
- New struct `DeploymentsStartJFRPollerResponse`
- New struct `DeploymentsStartJFRResponse`
- New struct `DiagnosticParameters`
- New struct `KeyVaultCertificateProperties`
- New struct `LoadedCertificate`
- New struct `ServicesBeginStartOptions`
- New struct `ServicesBeginStopOptions`
- New struct `ServicesStartPoller`
- New struct `ServicesStartPollerResponse`
- New struct `ServicesStartResponse`
- New struct `ServicesStopPoller`
- New struct `ServicesStopPollerResponse`
- New struct `ServicesStopResponse`
- New struct `StorageAccount`
- New struct `StorageProperties`
- New struct `StorageResource`
- New struct `StorageResourceCollection`
- New struct `StoragesBeginCreateOrUpdateOptions`
- New struct `StoragesBeginDeleteOptions`
- New struct `StoragesClient`
- New struct `StoragesCreateOrUpdatePoller`
- New struct `StoragesCreateOrUpdatePollerResponse`
- New struct `StoragesCreateOrUpdateResponse`
- New struct `StoragesCreateOrUpdateResult`
- New struct `StoragesDeletePoller`
- New struct `StoragesDeletePollerResponse`
- New struct `StoragesDeleteResponse`
- New struct `StoragesGetOptions`
- New struct `StoragesGetResponse`
- New struct `StoragesGetResult`
- New struct `StoragesListOptions`
- New struct `StoragesListPager`
- New struct `StoragesListResponse`
- New struct `StoragesListResult`
- New struct `SystemData`
- New field `ContainerProbeSettings` in struct `DeploymentSettings`
- New field `LoadedCertificates` in struct `AppResourceProperties`
- New field `CustomPersistentDisks` in struct `AppResourceProperties`
- New field `Type` in struct `CertificateProperties`
- New field `PowerState` in struct `ClusterResourceProperties`

Total 14 breaking change(s), 209 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-20)

- Initial preview release.
