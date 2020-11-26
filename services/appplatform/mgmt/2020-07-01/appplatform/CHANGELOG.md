
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewServiceResourceListPage` signature has been changed from `(func(context.Context, ServiceResourceList) (ServiceResourceList, error))` to `(ServiceResourceList,func(context.Context, ServiceResourceList) (ServiceResourceList, error))`
- Function `NewAvailableOperationsPage` signature has been changed from `(func(context.Context, AvailableOperations) (AvailableOperations, error))` to `(AvailableOperations,func(context.Context, AvailableOperations) (AvailableOperations, error))`
- Function `NewDeploymentResourceCollectionPage` signature has been changed from `(func(context.Context, DeploymentResourceCollection) (DeploymentResourceCollection, error))` to `(DeploymentResourceCollection,func(context.Context, DeploymentResourceCollection) (DeploymentResourceCollection, error))`
- Function `NewCustomDomainResourceCollectionPage` signature has been changed from `(func(context.Context, CustomDomainResourceCollection) (CustomDomainResourceCollection, error))` to `(CustomDomainResourceCollection,func(context.Context, CustomDomainResourceCollection) (CustomDomainResourceCollection, error))`
- Function `NewResourceSkuCollectionPage` signature has been changed from `(func(context.Context, ResourceSkuCollection) (ResourceSkuCollection, error))` to `(ResourceSkuCollection,func(context.Context, ResourceSkuCollection) (ResourceSkuCollection, error))`
- Function `NewBindingResourceCollectionPage` signature has been changed from `(func(context.Context, BindingResourceCollection) (BindingResourceCollection, error))` to `(BindingResourceCollection,func(context.Context, BindingResourceCollection) (BindingResourceCollection, error))`
- Function `NewCertificateResourceCollectionPage` signature has been changed from `(func(context.Context, CertificateResourceCollection) (CertificateResourceCollection, error))` to `(CertificateResourceCollection,func(context.Context, CertificateResourceCollection) (CertificateResourceCollection, error))`
- Function `NewAppResourceCollectionPage` signature has been changed from `(func(context.Context, AppResourceCollection) (AppResourceCollection, error))` to `(AppResourceCollection,func(context.Context, AppResourceCollection) (AppResourceCollection, error))`

## New Content

- Const `SupportedRuntimeValueNetCore31` is added
- Const `NetCoreZip` is added
- Const `NetCore31` is added
- Const `SupportedRuntimeValueJava11` is added
- Const `SupportedRuntimeValueJava8` is added
- Const `Java` is added
- Const `NETCore` is added
- Function `RuntimeVersionsClient.ListRuntimeVersionsSender(*http.Request) (*http.Response,error)` is added
- Function `PossibleSupportedRuntimePlatformValues() []SupportedRuntimePlatform` is added
- Function `ConfigServersClient.ValidatePreparer(context.Context,string,string,ConfigServerSettings) (*http.Request,error)` is added
- Function `RuntimeVersionsClient.ListRuntimeVersionsPreparer(context.Context) (*http.Request,error)` is added
- Function `NewRuntimeVersionsClientWithBaseURI(string,string) RuntimeVersionsClient` is added
- Function `*ConfigServersValidateFuture.Result(ConfigServersClient) (ConfigServerSettingsValidateResult,error)` is added
- Function `ConfigServersClient.Validate(context.Context,string,string,ConfigServerSettings) (ConfigServersValidateFuture,error)` is added
- Function `NetworkProfile.MarshalJSON() ([]byte,error)` is added
- Function `PossibleSupportedRuntimeValueValues() []SupportedRuntimeValue` is added
- Function `ConfigServersClient.ValidateResponder(*http.Response) (ConfigServerSettingsValidateResult,error)` is added
- Function `RuntimeVersionsClient.ListRuntimeVersionsResponder(*http.Response) (AvailableRuntimeVersions,error)` is added
- Function `RuntimeVersionsClient.ListRuntimeVersions(context.Context) (AvailableRuntimeVersions,error)` is added
- Function `NewRuntimeVersionsClient(string) RuntimeVersionsClient` is added
- Function `ConfigServersClient.ValidateSender(*http.Request) (ConfigServersValidateFuture,error)` is added
- Struct `AvailableRuntimeVersions` is added
- Struct `ConfigServerSettingsErrorRecord` is added
- Struct `ConfigServerSettingsValidateResult` is added
- Struct `ConfigServersValidateFuture` is added
- Struct `NetworkProfileOutboundIPs` is added
- Struct `RuntimeVersionsClient` is added
- Struct `SupportedRuntimeVersion` is added
- Field `OutboundIPs` is added to struct `NetworkProfile`
- Field `StartTime` is added to struct `DeploymentInstance`
- Field `NetCoreMainEntryPath` is added to struct `DeploymentSettings`

