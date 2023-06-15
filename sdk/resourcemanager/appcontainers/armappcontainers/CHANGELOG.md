# Release History

## 2.0.0-beta.3 (2023-05-26)
### Breaking Changes

- Type of `BillingMeterProperties.Category` has been changed from `*Category` to `*string`
- Type of `CustomDomainConfiguration.CertificatePassword` has been changed from `[]byte` to `*string`
- Type of `DaprSecretsCollection.Value` has been changed from `[]*Secret` to `[]*DaprSecret`
- Enum `Category` has been removed
- Enum `ManagedEnvironmentOutBoundType` has been removed
- Enum `SKUName` has been removed
- Struct `EnvironmentSKUProperties` has been removed
- Struct `ManagedEnvironmentOutboundSettings` has been removed
- Field `BillingMeterCategory` of struct `AvailableWorkloadProfileProperties` has been removed
- Field `WorkloadProfileType` of struct `ContainerAppProperties` has been removed
- Field `SKU` of struct `ManagedEnvironment` has been removed
- Field `OutboundSettings`, `RuntimeSubnetID` of struct `VnetConfiguration` has been removed

### Features Added

- New value `StorageTypeSecret` added to enum type `StorageType`
- New enum type `Affinity` with values `AffinityNone`, `AffinitySticky`
- New enum type `IngressClientCertificateMode` with values `IngressClientCertificateModeAccept`, `IngressClientCertificateModeIgnore`, `IngressClientCertificateModeRequire`
- New enum type `JobExecutionRunningState` with values `JobExecutionRunningStateDegraded`, `JobExecutionRunningStateFailed`, `JobExecutionRunningStateProcessing`, `JobExecutionRunningStateRunning`, `JobExecutionRunningStateStopped`, `JobExecutionRunningStateSucceeded`, `JobExecutionRunningStateUnknown`
- New enum type `JobProvisioningState` with values `JobProvisioningStateCanceled`, `JobProvisioningStateDeleting`, `JobProvisioningStateFailed`, `JobProvisioningStateInProgress`, `JobProvisioningStateSucceeded`
- New enum type `ManagedCertificateDomainControlValidation` with values `ManagedCertificateDomainControlValidationCNAME`, `ManagedCertificateDomainControlValidationHTTP`, `ManagedCertificateDomainControlValidationTXT`
- New enum type `TriggerType` with values `TriggerTypeEvent`, `TriggerTypeManual`, `TriggerTypeScheduled`
- New function `*ClientFactory.NewJobsClient() *JobsClient`
- New function `*ClientFactory.NewJobsExecutionsClient() *JobsExecutionsClient`
- New function `*ClientFactory.NewManagedCertificatesClient() *ManagedCertificatesClient`
- New function `NewJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*JobsClient, error)`
- New function `*JobsClient.BeginCreateOrUpdate(context.Context, string, string, Job, *JobsClientBeginCreateOrUpdateOptions) (*runtime.Poller[JobsClientCreateOrUpdateResponse], error)`
- New function `*JobsClient.BeginDelete(context.Context, string, string, *JobsClientBeginDeleteOptions) (*runtime.Poller[JobsClientDeleteResponse], error)`
- New function `*JobsClient.Get(context.Context, string, string, *JobsClientGetOptions) (JobsClientGetResponse, error)`
- New function `*JobsClient.NewListByResourceGroupPager(string, *JobsClientListByResourceGroupOptions) *runtime.Pager[JobsClientListByResourceGroupResponse]`
- New function `*JobsClient.NewListBySubscriptionPager(*JobsClientListBySubscriptionOptions) *runtime.Pager[JobsClientListBySubscriptionResponse]`
- New function `*JobsClient.ListSecrets(context.Context, string, string, *JobsClientListSecretsOptions) (JobsClientListSecretsResponse, error)`
- New function `*JobsClient.BeginStart(context.Context, string, string, JobExecutionTemplate, *JobsClientBeginStartOptions) (*runtime.Poller[JobsClientStartResponse], error)`
- New function `*JobsClient.BeginStopExecution(context.Context, string, string, string, *JobsClientBeginStopExecutionOptions) (*runtime.Poller[JobsClientStopExecutionResponse], error)`
- New function `*JobsClient.BeginStopMultipleExecutions(context.Context, string, string, JobExecutionNamesCollection, *JobsClientBeginStopMultipleExecutionsOptions) (*runtime.Poller[JobsClientStopMultipleExecutionsResponse], error)`
- New function `*JobsClient.BeginUpdate(context.Context, string, string, JobPatchProperties, *JobsClientBeginUpdateOptions) (*runtime.Poller[JobsClientUpdateResponse], error)`
- New function `NewJobsExecutionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*JobsExecutionsClient, error)`
- New function `*JobsExecutionsClient.NewListPager(string, string, *JobsExecutionsClientListOptions) *runtime.Pager[JobsExecutionsClientListResponse]`
- New function `NewManagedCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagedCertificatesClient, error)`
- New function `*ManagedCertificatesClient.BeginCreateOrUpdate(context.Context, string, string, string, *ManagedCertificatesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedCertificatesClientCreateOrUpdateResponse], error)`
- New function `*ManagedCertificatesClient.Delete(context.Context, string, string, string, *ManagedCertificatesClientDeleteOptions) (ManagedCertificatesClientDeleteResponse, error)`
- New function `*ManagedCertificatesClient.Get(context.Context, string, string, string, *ManagedCertificatesClientGetOptions) (ManagedCertificatesClientGetResponse, error)`
- New function `*ManagedCertificatesClient.NewListPager(string, string, *ManagedCertificatesClientListOptions) *runtime.Pager[ManagedCertificatesClientListResponse]`
- New function `*ManagedCertificatesClient.Update(context.Context, string, string, string, ManagedCertificatePatch, *ManagedCertificatesClientUpdateOptions) (ManagedCertificatesClientUpdateResponse, error)`
- New struct `ContainerAppJobExecutions`
- New struct `CorsPolicy`
- New struct `DaprConfiguration`
- New struct `IngressStickySessions`
- New struct `Job`
- New struct `JobConfiguration`
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
- New struct `JobSecretsCollection`
- New struct `JobTemplate`
- New struct `JobsCollection`
- New struct `KedaConfiguration`
- New struct `ManagedCertificate`
- New struct `ManagedCertificateCollection`
- New struct `ManagedCertificatePatch`
- New struct `ManagedCertificateProperties`
- New struct `SecretVolumeItem`
- New field `Category` in struct `AvailableWorkloadProfileProperties`
- New field `ManagedBy` in struct `ContainerApp`
- New field `LatestReadyRevisionName`, `WorkloadProfileName` in struct `ContainerAppProperties`
- New field `Identity`, `KeyVaultURL` in struct `ContainerAppSecret`
- New anonymous field `ContainerApp` in struct `ContainerAppsClientUpdateResponse`
- New field `ClientCertificateMode`, `CorsPolicy`, `StickySessions` in struct `Ingress`
- New field `Kind` in struct `ManagedEnvironment`
- New field `DaprConfiguration`, `InfrastructureResourceGroup`, `KedaConfiguration` in struct `ManagedEnvironmentProperties`
- New anonymous field `ManagedEnvironment` in struct `ManagedEnvironmentsClientUpdateResponse`
- New field `Identity`, `KeyVaultURL` in struct `Secret`
- New field `Secrets` in struct `Volume`
- New field `Name` in struct `WorkloadProfile`


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).