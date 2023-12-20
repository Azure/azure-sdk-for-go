# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-08-25)
### Breaking Changes

- Type of `DaprSecretsCollection.Value` has been changed from `[]*Secret` to `[]*DaprSecret`
- Struct `CustomHostnameAnalysisResultProperties` has been removed
- Field `ID`, `Name`, `Properties`, `SystemData`, `Type` of struct `CustomHostnameAnalysisResult` has been removed
- Field `RuntimeSubnetID` of struct `VnetConfiguration` has been removed

### Features Added

- New value `ContainerAppProvisioningStateDeleting` added to enum type `ContainerAppProvisioningState`
- New value `IngressTransportMethodTCP` added to enum type `IngressTransportMethod`
- New value `StorageTypeSecret` added to enum type `StorageType`
- New enum type `Action` with values `ActionAllow`, `ActionDeny`
- New enum type `Affinity` with values `AffinityNone`, `AffinitySticky`
- New enum type `Applicability` with values `ApplicabilityCustom`, `ApplicabilityLocationDefault`
- New enum type `ConnectedEnvironmentProvisioningState` with values `ConnectedEnvironmentProvisioningStateCanceled`, `ConnectedEnvironmentProvisioningStateFailed`, `ConnectedEnvironmentProvisioningStateInfrastructureSetupComplete`, `ConnectedEnvironmentProvisioningStateInfrastructureSetupInProgress`, `ConnectedEnvironmentProvisioningStateInitializationInProgress`, `ConnectedEnvironmentProvisioningStateScheduledForDelete`, `ConnectedEnvironmentProvisioningStateSucceeded`, `ConnectedEnvironmentProvisioningStateWaiting`
- New enum type `ContainerAppContainerRunningState` with values `ContainerAppContainerRunningStateRunning`, `ContainerAppContainerRunningStateTerminated`, `ContainerAppContainerRunningStateWaiting`
- New enum type `ContainerAppReplicaRunningState` with values `ContainerAppReplicaRunningStateNotRunning`, `ContainerAppReplicaRunningStateRunning`, `ContainerAppReplicaRunningStateUnknown`
- New enum type `ExtendedLocationTypes` with values `ExtendedLocationTypesCustomLocation`
- New enum type `IngressClientCertificateMode` with values `IngressClientCertificateModeAccept`, `IngressClientCertificateModeIgnore`, `IngressClientCertificateModeRequire`
- New enum type `JobExecutionRunningState` with values `JobExecutionRunningStateDegraded`, `JobExecutionRunningStateFailed`, `JobExecutionRunningStateProcessing`, `JobExecutionRunningStateRunning`, `JobExecutionRunningStateStopped`, `JobExecutionRunningStateSucceeded`, `JobExecutionRunningStateUnknown`
- New enum type `JobProvisioningState` with values `JobProvisioningStateCanceled`, `JobProvisioningStateDeleting`, `JobProvisioningStateFailed`, `JobProvisioningStateInProgress`, `JobProvisioningStateSucceeded`
- New enum type `LogLevel` with values `LogLevelDebug`, `LogLevelError`, `LogLevelInfo`, `LogLevelWarn`
- New enum type `ManagedCertificateDomainControlValidation` with values `ManagedCertificateDomainControlValidationCNAME`, `ManagedCertificateDomainControlValidationHTTP`, `ManagedCertificateDomainControlValidationTXT`
- New enum type `RevisionRunningState` with values `RevisionRunningStateDegraded`, `RevisionRunningStateFailed`, `RevisionRunningStateProcessing`, `RevisionRunningStateRunning`, `RevisionRunningStateStopped`, `RevisionRunningStateUnknown`
- New enum type `TriggerType` with values `TriggerTypeEvent`, `TriggerTypeManual`, `TriggerTypeSchedule`
- New function `NewAvailableWorkloadProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AvailableWorkloadProfilesClient, error)`
- New function `*AvailableWorkloadProfilesClient.NewGetPager(string, *AvailableWorkloadProfilesClientGetOptions) *runtime.Pager[AvailableWorkloadProfilesClientGetResponse]`
- New function `NewBillingMetersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BillingMetersClient, error)`
- New function `*BillingMetersClient.Get(context.Context, string, *BillingMetersClientGetOptions) (BillingMetersClientGetResponse, error)`
- New function `*ClientFactory.NewAvailableWorkloadProfilesClient() *AvailableWorkloadProfilesClient`
- New function `*ClientFactory.NewBillingMetersClient() *BillingMetersClient`
- New function `*ClientFactory.NewConnectedEnvironmentsCertificatesClient() *ConnectedEnvironmentsCertificatesClient`
- New function `*ClientFactory.NewConnectedEnvironmentsClient() *ConnectedEnvironmentsClient`
- New function `*ClientFactory.NewConnectedEnvironmentsDaprComponentsClient() *ConnectedEnvironmentsDaprComponentsClient`
- New function `*ClientFactory.NewConnectedEnvironmentsStoragesClient() *ConnectedEnvironmentsStoragesClient`
- New function `*ClientFactory.NewContainerAppsAPIClient() *ContainerAppsAPIClient`
- New function `*ClientFactory.NewContainerAppsDiagnosticsClient() *ContainerAppsDiagnosticsClient`
- New function `*ClientFactory.NewJobsClient() *JobsClient`
- New function `*ClientFactory.NewJobsExecutionsClient() *JobsExecutionsClient`
- New function `*ClientFactory.NewManagedCertificatesClient() *ManagedCertificatesClient`
- New function `*ClientFactory.NewManagedEnvironmentDiagnosticsClient() *ManagedEnvironmentDiagnosticsClient`
- New function `*ClientFactory.NewManagedEnvironmentsDiagnosticsClient() *ManagedEnvironmentsDiagnosticsClient`
- New function `NewConnectedEnvironmentsCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedEnvironmentsCertificatesClient, error)`
- New function `*ConnectedEnvironmentsCertificatesClient.CreateOrUpdate(context.Context, string, string, string, *ConnectedEnvironmentsCertificatesClientCreateOrUpdateOptions) (ConnectedEnvironmentsCertificatesClientCreateOrUpdateResponse, error)`
- New function `*ConnectedEnvironmentsCertificatesClient.Delete(context.Context, string, string, string, *ConnectedEnvironmentsCertificatesClientDeleteOptions) (ConnectedEnvironmentsCertificatesClientDeleteResponse, error)`
- New function `*ConnectedEnvironmentsCertificatesClient.Get(context.Context, string, string, string, *ConnectedEnvironmentsCertificatesClientGetOptions) (ConnectedEnvironmentsCertificatesClientGetResponse, error)`
- New function `*ConnectedEnvironmentsCertificatesClient.NewListPager(string, string, *ConnectedEnvironmentsCertificatesClientListOptions) *runtime.Pager[ConnectedEnvironmentsCertificatesClientListResponse]`
- New function `*ConnectedEnvironmentsCertificatesClient.Update(context.Context, string, string, string, CertificatePatch, *ConnectedEnvironmentsCertificatesClientUpdateOptions) (ConnectedEnvironmentsCertificatesClientUpdateResponse, error)`
- New function `NewConnectedEnvironmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedEnvironmentsClient, error)`
- New function `*ConnectedEnvironmentsClient.CheckNameAvailability(context.Context, string, string, CheckNameAvailabilityRequest, *ConnectedEnvironmentsClientCheckNameAvailabilityOptions) (ConnectedEnvironmentsClientCheckNameAvailabilityResponse, error)`
- New function `*ConnectedEnvironmentsClient.BeginCreateOrUpdate(context.Context, string, string, ConnectedEnvironment, *ConnectedEnvironmentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConnectedEnvironmentsClientCreateOrUpdateResponse], error)`
- New function `*ConnectedEnvironmentsClient.BeginDelete(context.Context, string, string, *ConnectedEnvironmentsClientBeginDeleteOptions) (*runtime.Poller[ConnectedEnvironmentsClientDeleteResponse], error)`
- New function `*ConnectedEnvironmentsClient.Get(context.Context, string, string, *ConnectedEnvironmentsClientGetOptions) (ConnectedEnvironmentsClientGetResponse, error)`
- New function `*ConnectedEnvironmentsClient.NewListByResourceGroupPager(string, *ConnectedEnvironmentsClientListByResourceGroupOptions) *runtime.Pager[ConnectedEnvironmentsClientListByResourceGroupResponse]`
- New function `*ConnectedEnvironmentsClient.NewListBySubscriptionPager(*ConnectedEnvironmentsClientListBySubscriptionOptions) *runtime.Pager[ConnectedEnvironmentsClientListBySubscriptionResponse]`
- New function `*ConnectedEnvironmentsClient.Update(context.Context, string, string, *ConnectedEnvironmentsClientUpdateOptions) (ConnectedEnvironmentsClientUpdateResponse, error)`
- New function `NewConnectedEnvironmentsDaprComponentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedEnvironmentsDaprComponentsClient, error)`
- New function `*ConnectedEnvironmentsDaprComponentsClient.CreateOrUpdate(context.Context, string, string, string, DaprComponent, *ConnectedEnvironmentsDaprComponentsClientCreateOrUpdateOptions) (ConnectedEnvironmentsDaprComponentsClientCreateOrUpdateResponse, error)`
- New function `*ConnectedEnvironmentsDaprComponentsClient.Delete(context.Context, string, string, string, *ConnectedEnvironmentsDaprComponentsClientDeleteOptions) (ConnectedEnvironmentsDaprComponentsClientDeleteResponse, error)`
- New function `*ConnectedEnvironmentsDaprComponentsClient.Get(context.Context, string, string, string, *ConnectedEnvironmentsDaprComponentsClientGetOptions) (ConnectedEnvironmentsDaprComponentsClientGetResponse, error)`
- New function `*ConnectedEnvironmentsDaprComponentsClient.NewListPager(string, string, *ConnectedEnvironmentsDaprComponentsClientListOptions) *runtime.Pager[ConnectedEnvironmentsDaprComponentsClientListResponse]`
- New function `*ConnectedEnvironmentsDaprComponentsClient.ListSecrets(context.Context, string, string, string, *ConnectedEnvironmentsDaprComponentsClientListSecretsOptions) (ConnectedEnvironmentsDaprComponentsClientListSecretsResponse, error)`
- New function `NewConnectedEnvironmentsStoragesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedEnvironmentsStoragesClient, error)`
- New function `*ConnectedEnvironmentsStoragesClient.CreateOrUpdate(context.Context, string, string, string, ConnectedEnvironmentStorage, *ConnectedEnvironmentsStoragesClientCreateOrUpdateOptions) (ConnectedEnvironmentsStoragesClientCreateOrUpdateResponse, error)`
- New function `*ConnectedEnvironmentsStoragesClient.Delete(context.Context, string, string, string, *ConnectedEnvironmentsStoragesClientDeleteOptions) (ConnectedEnvironmentsStoragesClientDeleteResponse, error)`
- New function `*ConnectedEnvironmentsStoragesClient.Get(context.Context, string, string, string, *ConnectedEnvironmentsStoragesClientGetOptions) (ConnectedEnvironmentsStoragesClientGetResponse, error)`
- New function `*ConnectedEnvironmentsStoragesClient.List(context.Context, string, string, *ConnectedEnvironmentsStoragesClientListOptions) (ConnectedEnvironmentsStoragesClientListResponse, error)`
- New function `NewContainerAppsAPIClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ContainerAppsAPIClient, error)`
- New function `*ContainerAppsAPIClient.JobExecution(context.Context, string, string, string, *ContainerAppsAPIClientJobExecutionOptions) (ContainerAppsAPIClientJobExecutionResponse, error)`
- New function `*ContainerAppsClient.GetAuthToken(context.Context, string, string, *ContainerAppsClientGetAuthTokenOptions) (ContainerAppsClientGetAuthTokenResponse, error)`
- New function `*ContainerAppsClient.BeginStart(context.Context, string, string, *ContainerAppsClientBeginStartOptions) (*runtime.Poller[ContainerAppsClientStartResponse], error)`
- New function `*ContainerAppsClient.BeginStop(context.Context, string, string, *ContainerAppsClientBeginStopOptions) (*runtime.Poller[ContainerAppsClientStopResponse], error)`
- New function `NewContainerAppsDiagnosticsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ContainerAppsDiagnosticsClient, error)`
- New function `*ContainerAppsDiagnosticsClient.GetDetector(context.Context, string, string, string, *ContainerAppsDiagnosticsClientGetDetectorOptions) (ContainerAppsDiagnosticsClientGetDetectorResponse, error)`
- New function `*ContainerAppsDiagnosticsClient.GetRevision(context.Context, string, string, string, *ContainerAppsDiagnosticsClientGetRevisionOptions) (ContainerAppsDiagnosticsClientGetRevisionResponse, error)`
- New function `*ContainerAppsDiagnosticsClient.GetRoot(context.Context, string, string, *ContainerAppsDiagnosticsClientGetRootOptions) (ContainerAppsDiagnosticsClientGetRootResponse, error)`
- New function `*ContainerAppsDiagnosticsClient.NewListDetectorsPager(string, string, *ContainerAppsDiagnosticsClientListDetectorsOptions) *runtime.Pager[ContainerAppsDiagnosticsClientListDetectorsResponse]`
- New function `*ContainerAppsDiagnosticsClient.NewListRevisionsPager(string, string, *ContainerAppsDiagnosticsClientListRevisionsOptions) *runtime.Pager[ContainerAppsDiagnosticsClientListRevisionsResponse]`
- New function `NewJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*JobsClient, error)`
- New function `*JobsClient.BeginCreateOrUpdate(context.Context, string, string, Job, *JobsClientBeginCreateOrUpdateOptions) (*runtime.Poller[JobsClientCreateOrUpdateResponse], error)`
- New function `*JobsClient.BeginDelete(context.Context, string, string, *JobsClientBeginDeleteOptions) (*runtime.Poller[JobsClientDeleteResponse], error)`
- New function `*JobsClient.Get(context.Context, string, string, *JobsClientGetOptions) (JobsClientGetResponse, error)`
- New function `*JobsClient.NewListByResourceGroupPager(string, *JobsClientListByResourceGroupOptions) *runtime.Pager[JobsClientListByResourceGroupResponse]`
- New function `*JobsClient.NewListBySubscriptionPager(*JobsClientListBySubscriptionOptions) *runtime.Pager[JobsClientListBySubscriptionResponse]`
- New function `*JobsClient.ListSecrets(context.Context, string, string, *JobsClientListSecretsOptions) (JobsClientListSecretsResponse, error)`
- New function `*JobsClient.BeginStart(context.Context, string, string, *JobsClientBeginStartOptions) (*runtime.Poller[JobsClientStartResponse], error)`
- New function `*JobsClient.BeginStopExecution(context.Context, string, string, string, *JobsClientBeginStopExecutionOptions) (*runtime.Poller[JobsClientStopExecutionResponse], error)`
- New function `*JobsClient.BeginStopMultipleExecutions(context.Context, string, string, *JobsClientBeginStopMultipleExecutionsOptions) (*runtime.Poller[JobsClientStopMultipleExecutionsResponse], error)`
- New function `*JobsClient.BeginUpdate(context.Context, string, string, JobPatchProperties, *JobsClientBeginUpdateOptions) (*runtime.Poller[JobsClientUpdateResponse], error)`
- New function `NewJobsExecutionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*JobsExecutionsClient, error)`
- New function `*JobsExecutionsClient.NewListPager(string, string, *JobsExecutionsClientListOptions) *runtime.Pager[JobsExecutionsClientListResponse]`
- New function `NewManagedCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedCertificatesClient, error)`
- New function `*ManagedCertificatesClient.BeginCreateOrUpdate(context.Context, string, string, string, *ManagedCertificatesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedCertificatesClientCreateOrUpdateResponse], error)`
- New function `*ManagedCertificatesClient.Delete(context.Context, string, string, string, *ManagedCertificatesClientDeleteOptions) (ManagedCertificatesClientDeleteResponse, error)`
- New function `*ManagedCertificatesClient.Get(context.Context, string, string, string, *ManagedCertificatesClientGetOptions) (ManagedCertificatesClientGetResponse, error)`
- New function `*ManagedCertificatesClient.NewListPager(string, string, *ManagedCertificatesClientListOptions) *runtime.Pager[ManagedCertificatesClientListResponse]`
- New function `*ManagedCertificatesClient.Update(context.Context, string, string, string, ManagedCertificatePatch, *ManagedCertificatesClientUpdateOptions) (ManagedCertificatesClientUpdateResponse, error)`
- New function `NewManagedEnvironmentDiagnosticsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedEnvironmentDiagnosticsClient, error)`
- New function `*ManagedEnvironmentDiagnosticsClient.GetDetector(context.Context, string, string, string, *ManagedEnvironmentDiagnosticsClientGetDetectorOptions) (ManagedEnvironmentDiagnosticsClientGetDetectorResponse, error)`
- New function `*ManagedEnvironmentDiagnosticsClient.ListDetectors(context.Context, string, string, *ManagedEnvironmentDiagnosticsClientListDetectorsOptions) (ManagedEnvironmentDiagnosticsClientListDetectorsResponse, error)`
- New function `*ManagedEnvironmentsClient.GetAuthToken(context.Context, string, string, *ManagedEnvironmentsClientGetAuthTokenOptions) (ManagedEnvironmentsClientGetAuthTokenResponse, error)`
- New function `*ManagedEnvironmentsClient.NewListWorkloadProfileStatesPager(string, string, *ManagedEnvironmentsClientListWorkloadProfileStatesOptions) *runtime.Pager[ManagedEnvironmentsClientListWorkloadProfileStatesResponse]`
- New function `NewManagedEnvironmentsDiagnosticsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedEnvironmentsDiagnosticsClient, error)`
- New function `*ManagedEnvironmentsDiagnosticsClient.GetRoot(context.Context, string, string, *ManagedEnvironmentsDiagnosticsClientGetRootOptions) (ManagedEnvironmentsDiagnosticsClientGetRootResponse, error)`
- New struct `AvailableWorkloadProfile`
- New struct `AvailableWorkloadProfileProperties`
- New struct `AvailableWorkloadProfilesCollection`
- New struct `BaseContainer`
- New struct `BillingMeter`
- New struct `BillingMeterCollection`
- New struct `BillingMeterProperties`
- New struct `ConnectedEnvironment`
- New struct `ConnectedEnvironmentCollection`
- New struct `ConnectedEnvironmentProperties`
- New struct `ConnectedEnvironmentStorage`
- New struct `ConnectedEnvironmentStorageProperties`
- New struct `ConnectedEnvironmentStoragesCollection`
- New struct `ContainerAppAuthToken`
- New struct `ContainerAppAuthTokenProperties`
- New struct `ContainerAppJobExecutions`
- New struct `CorsPolicy`
- New struct `CustomDomainConfiguration`
- New struct `CustomHostnameAnalysisResultCustomDomainVerificationFailureInfo`
- New struct `CustomHostnameAnalysisResultCustomDomainVerificationFailureInfoDetailsItem`
- New struct `DaprConfiguration`
- New struct `DaprSecret`
- New struct `DiagnosticDataProviderMetadata`
- New struct `DiagnosticDataProviderMetadataPropertyBagItem`
- New struct `DiagnosticDataTableResponseColumn`
- New struct `DiagnosticDataTableResponseObject`
- New struct `DiagnosticRendering`
- New struct `DiagnosticSupportTopic`
- New struct `Diagnostics`
- New struct `DiagnosticsCollection`
- New struct `DiagnosticsDataAPIResponse`
- New struct `DiagnosticsDefinition`
- New struct `DiagnosticsProperties`
- New struct `DiagnosticsStatus`
- New struct `EnvironmentAuthToken`
- New struct `EnvironmentAuthTokenProperties`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `ExtendedLocation`
- New struct `IPSecurityRestrictionRule`
- New struct `IngressStickySessions`
- New struct `InitContainer`
- New struct `Job`
- New struct `JobConfiguration`
- New struct `JobConfigurationEventTriggerConfig`
- New struct `JobConfigurationManualTriggerConfig`
- New struct `JobConfigurationScheduleTriggerConfig`
- New struct `JobExecution`
- New struct `JobExecutionBase`
- New struct `JobExecutionContainer`
- New struct `JobExecutionNamesCollection`
- New struct `JobExecutionTemplate`
- New struct `JobPatchProperties`
- New struct `JobPatchPropertiesProperties`
- New struct `JobProperties`
- New struct `JobScale`
- New struct `JobScaleRule`
- New struct `JobSecretsCollection`
- New struct `JobTemplate`
- New struct `JobsCollection`
- New struct `KedaConfiguration`
- New struct `ManagedCertificate`
- New struct `ManagedCertificateCollection`
- New struct `ManagedCertificatePatch`
- New struct `ManagedCertificateProperties`
- New struct `ManagedEnvironmentPropertiesPeerAuthentication`
- New struct `Mtls`
- New struct `SecretVolumeItem`
- New struct `Service`
- New struct `ServiceBind`
- New struct `TCPScaleRule`
- New struct `WorkloadProfile`
- New struct `WorkloadProfileStates`
- New struct `WorkloadProfileStatesCollection`
- New struct `WorkloadProfileStatesProperties`
- New field `Kind` in struct `AzureCredentials`
- New field `SubjectAlternativeNames` in struct `CertificateProperties`
- New field `MaxInactiveRevisions`, `Service` in struct `Configuration`
- New field `ExtendedLocation`, `ManagedBy` in struct `ContainerApp`
- New field `EnvironmentID`, `EventStreamEndpoint`, `LatestReadyRevisionName`, `WorkloadProfileName` in struct `ContainerAppProperties`
- New field `Identity`, `KeyVaultURL` in struct `ContainerAppSecret`
- New anonymous field `ContainerApp` in struct `ContainerAppsClientUpdateResponse`
- New field `ARecords`, `AlternateCNameRecords`, `AlternateTxtRecords`, `CNameRecords`, `ConflictWithEnvironmentCustomDomain`, `ConflictingContainerAppResourceID`, `CustomDomainVerificationFailureInfo`, `CustomDomainVerificationTest`, `HasConflictOnManagedEnvironment`, `HostName`, `IsHostnameAlreadyVerified`, `TxtRecords` in struct `CustomHostnameAnalysisResult`
- New field `EnableAPILogging`, `HTTPMaxRequestSize`, `HTTPReadBufferSize`, `LogLevel` in struct `Dapr`
- New field `SecretStoreComponent` in struct `DaprComponentProperties`
- New field `GithubPersonalAccessToken` in struct `GithubActionConfiguration`
- New field `ClientCertificateMode`, `CorsPolicy`, `ExposedPort`, `IPSecurityRestrictions`, `StickySessions` in struct `Ingress`
- New field `Kind` in struct `ManagedEnvironment`
- New field `CustomDomainConfiguration`, `DaprConfiguration`, `EventStreamEndpoint`, `InfrastructureResourceGroup`, `KedaConfiguration`, `PeerAuthentication`, `WorkloadProfiles` in struct `ManagedEnvironmentProperties`
- New anonymous field `ManagedEnvironment` in struct `ManagedEnvironmentsClientUpdateResponse`
- New field `ExecEndpoint`, `LogStreamEndpoint`, `RunningState`, `RunningStateDetails` in struct `ReplicaContainer`
- New field `InitContainers`, `RunningState`, `RunningStateDetails` in struct `ReplicaProperties`
- New field `LastActiveTime`, `RunningState` in struct `RevisionProperties`
- New field `TCP` in struct `ScaleRule`
- New field `Identity`, `KeyVaultURL` in struct `Secret`
- New field `InitContainers`, `ServiceBinds`, `TerminationGracePeriodSeconds` in struct `Template`
- New field `MountOptions`, `Secrets` in struct `Volume`
- New field `SubPath` in struct `VolumeMount`


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).