# Release History

## 2.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.1 (2023-04-28)
### Breaking Changes

- Type of `AppResourceProperties.AddonConfigs` has been changed from `map[string]map[string]any` to `map[string]any`
- Type of `BindingResourceProperties.BindingParameters` has been changed from `map[string]any` to `map[string]*string`
- Type of `DeploymentSettings.AddonConfigs` has been changed from `map[string]map[string]any` to `map[string]any`

### Features Added

- New value `BindingTypeCACertificates` added to enum type `BindingType`
- New enum type `APIPortalProvisioningState` with values `APIPortalProvisioningStateCreating`, `APIPortalProvisioningStateDeleting`, `APIPortalProvisioningStateFailed`, `APIPortalProvisioningStateSucceeded`, `APIPortalProvisioningStateUpdating`
- New enum type `ApmType` with values `ApmTypeAppDynamics`, `ApmTypeApplicationInsights`, `ApmTypeDynatrace`, `ApmTypeElasticAPM`, `ApmTypeNewRelic`
- New enum type `ApplicationAcceleratorProvisioningState` with values `ApplicationAcceleratorProvisioningStateCreating`, `ApplicationAcceleratorProvisioningStateDeleting`, `ApplicationAcceleratorProvisioningStateFailed`, `ApplicationAcceleratorProvisioningStateSucceeded`, `ApplicationAcceleratorProvisioningStateUpdating`
- New enum type `ApplicationLiveViewProvisioningState` with values `ApplicationLiveViewProvisioningStateCanceled`, `ApplicationLiveViewProvisioningStateCreating`, `ApplicationLiveViewProvisioningStateDeleting`, `ApplicationLiveViewProvisioningStateFailed`, `ApplicationLiveViewProvisioningStateSucceeded`, `ApplicationLiveViewProvisioningStateUpdating`
- New enum type `BackendProtocol` with values `BackendProtocolDefault`, `BackendProtocolGRPC`
- New enum type `CertificateResourceProvisioningState` with values `CertificateResourceProvisioningStateCreating`, `CertificateResourceProvisioningStateDeleting`, `CertificateResourceProvisioningStateFailed`, `CertificateResourceProvisioningStateSucceeded`, `CertificateResourceProvisioningStateUpdating`
- New enum type `CustomDomainResourceProvisioningState` with values `CustomDomainResourceProvisioningStateCreating`, `CustomDomainResourceProvisioningStateDeleting`, `CustomDomainResourceProvisioningStateFailed`, `CustomDomainResourceProvisioningStateSucceeded`, `CustomDomainResourceProvisioningStateUpdating`
- New enum type `CustomizedAcceleratorProvisioningState` with values `CustomizedAcceleratorProvisioningStateCreating`, `CustomizedAcceleratorProvisioningStateDeleting`, `CustomizedAcceleratorProvisioningStateFailed`, `CustomizedAcceleratorProvisioningStateSucceeded`, `CustomizedAcceleratorProvisioningStateUpdating`
- New enum type `CustomizedAcceleratorValidateResultState` with values `CustomizedAcceleratorValidateResultStateInvalid`, `CustomizedAcceleratorValidateResultStateValid`
- New enum type `DevToolPortalFeatureState` with values `DevToolPortalFeatureStateDisabled`, `DevToolPortalFeatureStateEnabled`
- New enum type `DevToolPortalProvisioningState` with values `DevToolPortalProvisioningStateCanceled`, `DevToolPortalProvisioningStateCreating`, `DevToolPortalProvisioningStateDeleting`, `DevToolPortalProvisioningStateFailed`, `DevToolPortalProvisioningStateSucceeded`, `DevToolPortalProvisioningStateUpdating`
- New enum type `GatewayProvisioningState` with values `GatewayProvisioningStateCreating`, `GatewayProvisioningStateDeleting`, `GatewayProvisioningStateFailed`, `GatewayProvisioningStateSucceeded`, `GatewayProvisioningStateUpdating`
- New enum type `GatewayRouteConfigProtocol` with values `GatewayRouteConfigProtocolHTTP`, `GatewayRouteConfigProtocolHTTPS`
- New enum type `HTTPSchemeType` with values `HTTPSchemeTypeHTTP`, `HTTPSchemeTypeHTTPS`
- New enum type `PowerState` with values `PowerStateRunning`, `PowerStateStopped`
- New enum type `PredefinedAcceleratorProvisioningState` with values `PredefinedAcceleratorProvisioningStateCreating`, `PredefinedAcceleratorProvisioningStateFailed`, `PredefinedAcceleratorProvisioningStateSucceeded`, `PredefinedAcceleratorProvisioningStateUpdating`
- New enum type `PredefinedAcceleratorState` with values `PredefinedAcceleratorStateDisabled`, `PredefinedAcceleratorStateEnabled`
- New enum type `ProbeActionType` with values `ProbeActionTypeExecAction`, `ProbeActionTypeHTTPGetAction`, `ProbeActionTypeTCPSocketAction`
- New enum type `SessionAffinity` with values `SessionAffinityCookie`, `SessionAffinityNone`
- New enum type `StorageType` with values `StorageTypeStorageAccount`
- New enum type `Type` with values `TypeAzureFileVolume`
- New function `NewAPIPortalCustomDomainsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIPortalCustomDomainsClient, error)`
- New function `*APIPortalCustomDomainsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, APIPortalCustomDomainResource, *APIPortalCustomDomainsClientBeginCreateOrUpdateOptions) (*runtime.Poller[APIPortalCustomDomainsClientCreateOrUpdateResponse], error)`
- New function `*APIPortalCustomDomainsClient.BeginDelete(context.Context, string, string, string, string, *APIPortalCustomDomainsClientBeginDeleteOptions) (*runtime.Poller[APIPortalCustomDomainsClientDeleteResponse], error)`
- New function `*APIPortalCustomDomainsClient.Get(context.Context, string, string, string, string, *APIPortalCustomDomainsClientGetOptions) (APIPortalCustomDomainsClientGetResponse, error)`
- New function `*APIPortalCustomDomainsClient.NewListPager(string, string, string, *APIPortalCustomDomainsClientListOptions) *runtime.Pager[APIPortalCustomDomainsClientListResponse]`
- New function `NewAPIPortalsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIPortalsClient, error)`
- New function `*APIPortalsClient.BeginCreateOrUpdate(context.Context, string, string, string, APIPortalResource, *APIPortalsClientBeginCreateOrUpdateOptions) (*runtime.Poller[APIPortalsClientCreateOrUpdateResponse], error)`
- New function `*APIPortalsClient.BeginDelete(context.Context, string, string, string, *APIPortalsClientBeginDeleteOptions) (*runtime.Poller[APIPortalsClientDeleteResponse], error)`
- New function `*APIPortalsClient.Get(context.Context, string, string, string, *APIPortalsClientGetOptions) (APIPortalsClientGetResponse, error)`
- New function `*APIPortalsClient.NewListPager(string, string, *APIPortalsClientListOptions) *runtime.Pager[APIPortalsClientListResponse]`
- New function `*APIPortalsClient.ValidateDomain(context.Context, string, string, string, CustomDomainValidatePayload, *APIPortalsClientValidateDomainOptions) (APIPortalsClientValidateDomainResponse, error)`
- New function `*AcceleratorAuthSetting.GetAcceleratorAuthSetting() *AcceleratorAuthSetting`
- New function `*AcceleratorBasicAuthSetting.GetAcceleratorAuthSetting() *AcceleratorAuthSetting`
- New function `*AcceleratorPublicSetting.GetAcceleratorAuthSetting() *AcceleratorAuthSetting`
- New function `*AcceleratorSSHSetting.GetAcceleratorAuthSetting() *AcceleratorAuthSetting`
- New function `NewApplicationAcceleratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApplicationAcceleratorsClient, error)`
- New function `*ApplicationAcceleratorsClient.BeginCreateOrUpdate(context.Context, string, string, string, ApplicationAcceleratorResource, *ApplicationAcceleratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ApplicationAcceleratorsClientCreateOrUpdateResponse], error)`
- New function `*ApplicationAcceleratorsClient.BeginDelete(context.Context, string, string, string, *ApplicationAcceleratorsClientBeginDeleteOptions) (*runtime.Poller[ApplicationAcceleratorsClientDeleteResponse], error)`
- New function `*ApplicationAcceleratorsClient.Get(context.Context, string, string, string, *ApplicationAcceleratorsClientGetOptions) (ApplicationAcceleratorsClientGetResponse, error)`
- New function `*ApplicationAcceleratorsClient.NewListPager(string, string, *ApplicationAcceleratorsClientListOptions) *runtime.Pager[ApplicationAcceleratorsClientListResponse]`
- New function `NewApplicationLiveViewsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ApplicationLiveViewsClient, error)`
- New function `*ApplicationLiveViewsClient.BeginCreateOrUpdate(context.Context, string, string, string, ApplicationLiveViewResource, *ApplicationLiveViewsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ApplicationLiveViewsClientCreateOrUpdateResponse], error)`
- New function `*ApplicationLiveViewsClient.BeginDelete(context.Context, string, string, string, *ApplicationLiveViewsClientBeginDeleteOptions) (*runtime.Poller[ApplicationLiveViewsClientDeleteResponse], error)`
- New function `*ApplicationLiveViewsClient.Get(context.Context, string, string, string, *ApplicationLiveViewsClientGetOptions) (ApplicationLiveViewsClientGetResponse, error)`
- New function `*ApplicationLiveViewsClient.NewListPager(string, string, *ApplicationLiveViewsClientListOptions) *runtime.Pager[ApplicationLiveViewsClientListResponse]`
- New function `*AzureFileVolume.GetCustomPersistentDiskProperties() *CustomPersistentDiskProperties`
- New function `*BuildServiceBuilderClient.ListDeployments(context.Context, string, string, string, string, *BuildServiceBuilderClientListDeploymentsOptions) (BuildServiceBuilderClientListDeploymentsResponse, error)`
- New function `*BuildpackBindingClient.NewListForClusterPager(string, string, *BuildpackBindingClientListForClusterOptions) *runtime.Pager[BuildpackBindingClientListForClusterResponse]`
- New function `*ClientFactory.NewAPIPortalCustomDomainsClient() *APIPortalCustomDomainsClient`
- New function `*ClientFactory.NewAPIPortalsClient() *APIPortalsClient`
- New function `*ClientFactory.NewApplicationAcceleratorsClient() *ApplicationAcceleratorsClient`
- New function `*ClientFactory.NewApplicationLiveViewsClient() *ApplicationLiveViewsClient`
- New function `*ClientFactory.NewCustomizedAcceleratorsClient() *CustomizedAcceleratorsClient`
- New function `*ClientFactory.NewDevToolPortalsClient() *DevToolPortalsClient`
- New function `*ClientFactory.NewGatewayCustomDomainsClient() *GatewayCustomDomainsClient`
- New function `*ClientFactory.NewGatewayRouteConfigsClient() *GatewayRouteConfigsClient`
- New function `*ClientFactory.NewGatewaysClient() *GatewaysClient`
- New function `*ClientFactory.NewPredefinedAcceleratorsClient() *PredefinedAcceleratorsClient`
- New function `*ClientFactory.NewStoragesClient() *StoragesClient`
- New function `*CustomContainerUserSourceInfo.GetUserSourceInfo() *UserSourceInfo`
- New function `*CustomPersistentDiskProperties.GetCustomPersistentDiskProperties() *CustomPersistentDiskProperties`
- New function `NewCustomizedAcceleratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomizedAcceleratorsClient, error)`
- New function `*CustomizedAcceleratorsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, CustomizedAcceleratorResource, *CustomizedAcceleratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomizedAcceleratorsClientCreateOrUpdateResponse], error)`
- New function `*CustomizedAcceleratorsClient.BeginDelete(context.Context, string, string, string, string, *CustomizedAcceleratorsClientBeginDeleteOptions) (*runtime.Poller[CustomizedAcceleratorsClientDeleteResponse], error)`
- New function `*CustomizedAcceleratorsClient.Get(context.Context, string, string, string, string, *CustomizedAcceleratorsClientGetOptions) (CustomizedAcceleratorsClientGetResponse, error)`
- New function `*CustomizedAcceleratorsClient.NewListPager(string, string, string, *CustomizedAcceleratorsClientListOptions) *runtime.Pager[CustomizedAcceleratorsClientListResponse]`
- New function `*CustomizedAcceleratorsClient.Validate(context.Context, string, string, string, string, CustomizedAcceleratorProperties, *CustomizedAcceleratorsClientValidateOptions) (CustomizedAcceleratorsClientValidateResponse, error)`
- New function `*DeploymentsClient.BeginDisableRemoteDebugging(context.Context, string, string, string, string, *DeploymentsClientBeginDisableRemoteDebuggingOptions) (*runtime.Poller[DeploymentsClientDisableRemoteDebuggingResponse], error)`
- New function `*DeploymentsClient.BeginEnableRemoteDebugging(context.Context, string, string, string, string, *DeploymentsClientBeginEnableRemoteDebuggingOptions) (*runtime.Poller[DeploymentsClientEnableRemoteDebuggingResponse], error)`
- New function `*DeploymentsClient.GetRemoteDebuggingConfig(context.Context, string, string, string, string, *DeploymentsClientGetRemoteDebuggingConfigOptions) (DeploymentsClientGetRemoteDebuggingConfigResponse, error)`
- New function `NewDevToolPortalsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DevToolPortalsClient, error)`
- New function `*DevToolPortalsClient.BeginCreateOrUpdate(context.Context, string, string, string, DevToolPortalResource, *DevToolPortalsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DevToolPortalsClientCreateOrUpdateResponse], error)`
- New function `*DevToolPortalsClient.BeginDelete(context.Context, string, string, string, *DevToolPortalsClientBeginDeleteOptions) (*runtime.Poller[DevToolPortalsClientDeleteResponse], error)`
- New function `*DevToolPortalsClient.Get(context.Context, string, string, string, *DevToolPortalsClientGetOptions) (DevToolPortalsClientGetResponse, error)`
- New function `*DevToolPortalsClient.NewListPager(string, string, *DevToolPortalsClientListOptions) *runtime.Pager[DevToolPortalsClientListResponse]`
- New function `*ExecAction.GetProbeAction() *ProbeAction`
- New function `NewGatewayCustomDomainsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatewayCustomDomainsClient, error)`
- New function `*GatewayCustomDomainsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, GatewayCustomDomainResource, *GatewayCustomDomainsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewayCustomDomainsClientCreateOrUpdateResponse], error)`
- New function `*GatewayCustomDomainsClient.BeginDelete(context.Context, string, string, string, string, *GatewayCustomDomainsClientBeginDeleteOptions) (*runtime.Poller[GatewayCustomDomainsClientDeleteResponse], error)`
- New function `*GatewayCustomDomainsClient.Get(context.Context, string, string, string, string, *GatewayCustomDomainsClientGetOptions) (GatewayCustomDomainsClientGetResponse, error)`
- New function `*GatewayCustomDomainsClient.NewListPager(string, string, string, *GatewayCustomDomainsClientListOptions) *runtime.Pager[GatewayCustomDomainsClientListResponse]`
- New function `NewGatewayRouteConfigsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatewayRouteConfigsClient, error)`
- New function `*GatewayRouteConfigsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, GatewayRouteConfigResource, *GatewayRouteConfigsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewayRouteConfigsClientCreateOrUpdateResponse], error)`
- New function `*GatewayRouteConfigsClient.BeginDelete(context.Context, string, string, string, string, *GatewayRouteConfigsClientBeginDeleteOptions) (*runtime.Poller[GatewayRouteConfigsClientDeleteResponse], error)`
- New function `*GatewayRouteConfigsClient.Get(context.Context, string, string, string, string, *GatewayRouteConfigsClientGetOptions) (GatewayRouteConfigsClientGetResponse, error)`
- New function `*GatewayRouteConfigsClient.NewListPager(string, string, string, *GatewayRouteConfigsClientListOptions) *runtime.Pager[GatewayRouteConfigsClientListResponse]`
- New function `NewGatewaysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatewaysClient, error)`
- New function `*GatewaysClient.BeginCreateOrUpdate(context.Context, string, string, string, GatewayResource, *GatewaysClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewaysClientCreateOrUpdateResponse], error)`
- New function `*GatewaysClient.BeginDelete(context.Context, string, string, string, *GatewaysClientBeginDeleteOptions) (*runtime.Poller[GatewaysClientDeleteResponse], error)`
- New function `*GatewaysClient.Get(context.Context, string, string, string, *GatewaysClientGetOptions) (GatewaysClientGetResponse, error)`
- New function `*GatewaysClient.ListEnvSecrets(context.Context, string, string, string, *GatewaysClientListEnvSecretsOptions) (GatewaysClientListEnvSecretsResponse, error)`
- New function `*GatewaysClient.NewListPager(string, string, *GatewaysClientListOptions) *runtime.Pager[GatewaysClientListResponse]`
- New function `*GatewaysClient.BeginUpdateCapacity(context.Context, string, string, string, SKUObject, *GatewaysClientBeginUpdateCapacityOptions) (*runtime.Poller[GatewaysClientUpdateCapacityResponse], error)`
- New function `*GatewaysClient.ValidateDomain(context.Context, string, string, string, CustomDomainValidatePayload, *GatewaysClientValidateDomainOptions) (GatewaysClientValidateDomainResponse, error)`
- New function `*HTTPGetAction.GetProbeAction() *ProbeAction`
- New function `NewPredefinedAcceleratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PredefinedAcceleratorsClient, error)`
- New function `*PredefinedAcceleratorsClient.BeginDisable(context.Context, string, string, string, string, *PredefinedAcceleratorsClientBeginDisableOptions) (*runtime.Poller[PredefinedAcceleratorsClientDisableResponse], error)`
- New function `*PredefinedAcceleratorsClient.BeginEnable(context.Context, string, string, string, string, *PredefinedAcceleratorsClientBeginEnableOptions) (*runtime.Poller[PredefinedAcceleratorsClientEnableResponse], error)`
- New function `*PredefinedAcceleratorsClient.Get(context.Context, string, string, string, string, *PredefinedAcceleratorsClientGetOptions) (PredefinedAcceleratorsClientGetResponse, error)`
- New function `*PredefinedAcceleratorsClient.NewListPager(string, string, string, *PredefinedAcceleratorsClientListOptions) *runtime.Pager[PredefinedAcceleratorsClientListResponse]`
- New function `*ProbeAction.GetProbeAction() *ProbeAction`
- New function `*ServicesClient.BeginStart(context.Context, string, string, *ServicesClientBeginStartOptions) (*runtime.Poller[ServicesClientStartResponse], error)`
- New function `*ServicesClient.BeginStop(context.Context, string, string, *ServicesClientBeginStopOptions) (*runtime.Poller[ServicesClientStopResponse], error)`
- New function `*StorageAccount.GetStorageProperties() *StorageProperties`
- New function `*StorageProperties.GetStorageProperties() *StorageProperties`
- New function `NewStoragesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*StoragesClient, error)`
- New function `*StoragesClient.BeginCreateOrUpdate(context.Context, string, string, string, StorageResource, *StoragesClientBeginCreateOrUpdateOptions) (*runtime.Poller[StoragesClientCreateOrUpdateResponse], error)`
- New function `*StoragesClient.BeginDelete(context.Context, string, string, string, *StoragesClientBeginDeleteOptions) (*runtime.Poller[StoragesClientDeleteResponse], error)`
- New function `*StoragesClient.Get(context.Context, string, string, string, *StoragesClientGetOptions) (StoragesClientGetResponse, error)`
- New function `*StoragesClient.NewListPager(string, string, *StoragesClientListOptions) *runtime.Pager[StoragesClientListResponse]`
- New function `*TCPSocketAction.GetProbeAction() *ProbeAction`
- New struct `APIPortalCustomDomainProperties`
- New struct `APIPortalCustomDomainResource`
- New struct `APIPortalCustomDomainResourceCollection`
- New struct `APIPortalInstance`
- New struct `APIPortalProperties`
- New struct `APIPortalResource`
- New struct `APIPortalResourceCollection`
- New struct `APIPortalResourceRequests`
- New struct `AcceleratorBasicAuthSetting`
- New struct `AcceleratorGitRepository`
- New struct `AcceleratorPublicSetting`
- New struct `AcceleratorSSHSetting`
- New struct `AppVNetAddons`
- New struct `ApplicationAcceleratorComponent`
- New struct `ApplicationAcceleratorInstance`
- New struct `ApplicationAcceleratorProperties`
- New struct `ApplicationAcceleratorResource`
- New struct `ApplicationAcceleratorResourceCollection`
- New struct `ApplicationAcceleratorResourceRequests`
- New struct `ApplicationLiveViewComponent`
- New struct `ApplicationLiveViewInstance`
- New struct `ApplicationLiveViewProperties`
- New struct `ApplicationLiveViewResource`
- New struct `ApplicationLiveViewResourceCollection`
- New struct `ApplicationLiveViewResourceRequests`
- New struct `AzureFileVolume`
- New struct `BuildResourceRequests`
- New struct `ContainerProbeSettings`
- New struct `CustomContainer`
- New struct `CustomContainerUserSourceInfo`
- New struct `CustomPersistentDiskResource`
- New struct `CustomScaleRule`
- New struct `CustomizedAcceleratorProperties`
- New struct `CustomizedAcceleratorResource`
- New struct `CustomizedAcceleratorResourceCollection`
- New struct `CustomizedAcceleratorValidateResult`
- New struct `DeploymentList`
- New struct `DevToolPortalFeatureDetail`
- New struct `DevToolPortalFeatureSettings`
- New struct `DevToolPortalInstance`
- New struct `DevToolPortalProperties`
- New struct `DevToolPortalResource`
- New struct `DevToolPortalResourceCollection`
- New struct `DevToolPortalResourceRequests`
- New struct `DevToolPortalSsoProperties`
- New struct `ExecAction`
- New struct `GatewayAPIMetadataProperties`
- New struct `GatewayAPIRoute`
- New struct `GatewayCorsProperties`
- New struct `GatewayCustomDomainProperties`
- New struct `GatewayCustomDomainResource`
- New struct `GatewayCustomDomainResourceCollection`
- New struct `GatewayInstance`
- New struct `GatewayOperatorProperties`
- New struct `GatewayOperatorResourceRequests`
- New struct `GatewayProperties`
- New struct `GatewayPropertiesEnvironmentVariables`
- New struct `GatewayResource`
- New struct `GatewayResourceCollection`
- New struct `GatewayResourceRequests`
- New struct `GatewayRouteConfigOpenAPIProperties`
- New struct `GatewayRouteConfigProperties`
- New struct `GatewayRouteConfigResource`
- New struct `GatewayRouteConfigResourceCollection`
- New struct `HTTPGetAction`
- New struct `HTTPScaleRule`
- New struct `ImageRegistryCredential`
- New struct `IngressConfig`
- New struct `IngressSettings`
- New struct `IngressSettingsClientAuth`
- New struct `MarketplaceResource`
- New struct `PredefinedAcceleratorProperties`
- New struct `PredefinedAcceleratorResource`
- New struct `PredefinedAcceleratorResourceCollection`
- New struct `Probe`
- New struct `QueueScaleRule`
- New struct `RemoteDebugging`
- New struct `RemoteDebuggingPayload`
- New struct `SKUObject`
- New struct `Scale`
- New struct `ScaleRule`
- New struct `ScaleRuleAuth`
- New struct `Secret`
- New struct `ServiceVNetAddons`
- New struct `SsoProperties`
- New struct `StorageAccount`
- New struct `StorageResource`
- New struct `StorageResourceCollection`
- New struct `TCPScaleRule`
- New struct `TCPSocketAction`
- New struct `UserAssignedManagedIdentity`
- New field `CustomPersistentDisks` in struct `AppResourceProperties`
- New field `IngressSettings` in struct `AppResourceProperties`
- New field `Secrets` in struct `AppResourceProperties`
- New field `VnetAddons` in struct `AppResourceProperties`
- New field `ResourceRequests` in struct `BuildProperties`
- New field `Error` in struct `BuildResultProperties`
- New field `ExitCode` in struct `BuildStageProperties`
- New field `Reason` in struct `BuildStageProperties`
- New field `ProvisioningState` in struct `CertificateProperties`
- New field `InfraResourceGroup` in struct `ClusterResourceProperties`
- New field `ManagedEnvironmentID` in struct `ClusterResourceProperties`
- New field `MarketplaceResource` in struct `ClusterResourceProperties`
- New field `PowerState` in struct `ClusterResourceProperties`
- New field `VnetAddons` in struct `ClusterResourceProperties`
- New field `ProvisioningState` in struct `ContentCertificateProperties`
- New field `ProvisioningState` in struct `CustomDomainProperties`
- New field `ContainerProbeSettings` in struct `DeploymentSettings`
- New field `LivenessProbe` in struct `DeploymentSettings`
- New field `ReadinessProbe` in struct `DeploymentSettings`
- New field `Scale` in struct `DeploymentSettings`
- New field `StartupProbe` in struct `DeploymentSettings`
- New field `TerminationGracePeriodSeconds` in struct `DeploymentSettings`
- New field `ProvisioningState` in struct `KeyVaultCertificateProperties`
- New field `UserAssignedIdentities` in struct `ManagedIdentityProperties`
- New field `IngressConfig` in struct `NetworkProfile`
- New field `OutboundType` in struct `NetworkProfile`


## 1.1.0 (2023-04-06)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).