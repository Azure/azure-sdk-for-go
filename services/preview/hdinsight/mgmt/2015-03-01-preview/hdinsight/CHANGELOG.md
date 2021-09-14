# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. OperationResource

#### Removed Struct Fields

1. CapabilitiesResult.VMSizeFilters
1. CapabilitiesResult.VMSizes
1. Extension.autorest.Response
1. VMSizeCompatibilityFilter.Vmsizes

### Signature Changes

#### Funcs

1. ExtensionClient.Create
	- Returns
		- From: autorest.Response, error
		- To: ExtensionCreateFuture, error
1. ExtensionClient.CreateSender
	- Returns
		- From: *http.Response, error
		- To: ExtensionCreateFuture, error
1. ExtensionClient.Delete
	- Returns
		- From: autorest.Response, error
		- To: ExtensionDeleteFuture, error
1. ExtensionClient.DeleteSender
	- Returns
		- From: *http.Response, error
		- To: ExtensionDeleteFuture, error
1. ExtensionClient.Get
	- Returns
		- From: Extension, error
		- To: ClusterMonitoringResponse, error
1. ExtensionClient.GetResponder
	- Returns
		- From: Extension, error
		- To: ClusterMonitoringResponse, error

#### Struct Fields

1. Usage.CurrentValue changed type from *int32 to *int64
1. Usage.Limit changed type from *int32 to *int64
1. VersionSpec.IsDefault changed type from *string to *bool

## Additive Changes

### New Constants

1. FilterMode.Default
1. FilterMode.Recommend

### New Funcs

1. *ClustersUpdateIdentityCertificateFuture.UnmarshalJSON([]byte) error
1. *ExtensionCreateFuture.UnmarshalJSON([]byte) error
1. *ExtensionDeleteFuture.UnmarshalJSON([]byte) error
1. *ExtensionsDisableAzureMonitorFuture.UnmarshalJSON([]byte) error
1. *ExtensionsEnableAzureMonitorFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesRestartHostsFuture.UnmarshalJSON([]byte) error
1. ApplicationGetHTTPSEndpoint.MarshalJSON() ([]byte, error)
1. ApplicationsClient.GetAzureAsyncOperationStatus(context.Context, string, string, string, string) (AsyncOperationResult, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. AzureMonitorSelectedConfigurations.MarshalJSON() ([]byte, error)
1. BillingResponseListResult.MarshalJSON() ([]byte, error)
1. ClusterCreateRequestValidationParameters.MarshalJSON() ([]byte, error)
1. ClustersClient.GetAzureAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. ClustersClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. ClustersClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ClustersClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. ClustersClient.UpdateIdentityCertificate(context.Context, string, string, UpdateClusterIdentityCertificateParameters) (ClustersUpdateIdentityCertificateFuture, error)
1. ClustersClient.UpdateIdentityCertificatePreparer(context.Context, string, string, UpdateClusterIdentityCertificateParameters) (*http.Request, error)
1. ClustersClient.UpdateIdentityCertificateResponder(*http.Response) (autorest.Response, error)
1. ClustersClient.UpdateIdentityCertificateSender(*http.Request) (ClustersUpdateIdentityCertificateFuture, error)
1. ExtensionsClient.DisableAzureMonitor(context.Context, string, string) (ExtensionsDisableAzureMonitorFuture, error)
1. ExtensionsClient.DisableAzureMonitorPreparer(context.Context, string, string) (*http.Request, error)
1. ExtensionsClient.DisableAzureMonitorResponder(*http.Response) (autorest.Response, error)
1. ExtensionsClient.DisableAzureMonitorSender(*http.Request) (ExtensionsDisableAzureMonitorFuture, error)
1. ExtensionsClient.EnableAzureMonitor(context.Context, string, string, AzureMonitorRequest) (ExtensionsEnableAzureMonitorFuture, error)
1. ExtensionsClient.EnableAzureMonitorPreparer(context.Context, string, string, AzureMonitorRequest) (*http.Request, error)
1. ExtensionsClient.EnableAzureMonitorResponder(*http.Response) (autorest.Response, error)
1. ExtensionsClient.EnableAzureMonitorSender(*http.Request) (ExtensionsEnableAzureMonitorFuture, error)
1. ExtensionsClient.GetAzureAsyncOperationStatus(context.Context, string, string, string, string) (AsyncOperationResult, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. ExtensionsClient.GetAzureMonitorStatus(context.Context, string, string) (AzureMonitorResponse, error)
1. ExtensionsClient.GetAzureMonitorStatusPreparer(context.Context, string, string) (*http.Request, error)
1. ExtensionsClient.GetAzureMonitorStatusResponder(*http.Response) (AzureMonitorResponse, error)
1. ExtensionsClient.GetAzureMonitorStatusSender(*http.Request) (*http.Response, error)
1. KafkaRestProperties.MarshalJSON() ([]byte, error)
1. LocationsClient.CheckNameAvailability(context.Context, string, NameAvailabilityCheckRequestParameters) (NameAvailabilityCheckResult, error)
1. LocationsClient.CheckNameAvailabilityPreparer(context.Context, string, NameAvailabilityCheckRequestParameters) (*http.Request, error)
1. LocationsClient.CheckNameAvailabilityResponder(*http.Response) (NameAvailabilityCheckResult, error)
1. LocationsClient.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)
1. LocationsClient.GetAzureAsyncOperationStatus(context.Context, string, string) (AsyncOperationResult, error)
1. LocationsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string) (*http.Request, error)
1. LocationsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. LocationsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. LocationsClient.ValidateClusterCreateRequest(context.Context, string, ClusterCreateRequestValidationParameters) (ClusterCreateValidationResult, error)
1. LocationsClient.ValidateClusterCreateRequestPreparer(context.Context, string, ClusterCreateRequestValidationParameters) (*http.Request, error)
1. LocationsClient.ValidateClusterCreateRequestResponder(*http.Response) (ClusterCreateValidationResult, error)
1. LocationsClient.ValidateClusterCreateRequestSender(*http.Request) (*http.Response, error)
1. NameAvailabilityCheckResult.MarshalJSON() ([]byte, error)
1. NewExtensionsClient(string) ExtensionsClient
1. NewExtensionsClientWithBaseURI(string, string) ExtensionsClient
1. NewVirtualMachinesClient(string) VirtualMachinesClient
1. NewVirtualMachinesClientWithBaseURI(string, string) VirtualMachinesClient
1. ScriptActionsClient.GetExecutionAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. VirtualMachinesClient.GetAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. VirtualMachinesClient.GetAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. VirtualMachinesClient.GetAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. VirtualMachinesClient.GetAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. VirtualMachinesClient.ListHosts(context.Context, string, string) (ListHostInfo, error)
1. VirtualMachinesClient.ListHostsPreparer(context.Context, string, string) (*http.Request, error)
1. VirtualMachinesClient.ListHostsResponder(*http.Response) (ListHostInfo, error)
1. VirtualMachinesClient.ListHostsSender(*http.Request) (*http.Response, error)
1. VirtualMachinesClient.RestartHosts(context.Context, string, string, []string) (VirtualMachinesRestartHostsFuture, error)
1. VirtualMachinesClient.RestartHostsPreparer(context.Context, string, string, []string) (*http.Request, error)
1. VirtualMachinesClient.RestartHostsResponder(*http.Response) (autorest.Response, error)
1. VirtualMachinesClient.RestartHostsSender(*http.Request) (VirtualMachinesRestartHostsFuture, error)

### Struct Changes

#### New Structs

1. AaddsResourceDetails
1. AsyncOperationResult
1. AzureMonitorRequest
1. AzureMonitorResponse
1. AzureMonitorSelectedConfigurations
1. AzureMonitorTableConfiguration
1. ClusterCreateRequestValidationParameters
1. ClusterCreateValidationResult
1. ClustersUpdateIdentityCertificateFuture
1. ComputeIsolationProperties
1. Dimension
1. ExcludedServicesConfig
1. ExtensionCreateFuture
1. ExtensionDeleteFuture
1. ExtensionsClient
1. ExtensionsDisableAzureMonitorFuture
1. ExtensionsEnableAzureMonitorFuture
1. HostInfo
1. ListHostInfo
1. MetricSpecifications
1. NameAvailabilityCheckRequestParameters
1. NameAvailabilityCheckResult
1. OperationProperties
1. ServiceSpecification
1. UpdateClusterIdentityCertificateParameters
1. VMSizeProperty
1. ValidationErrorInfo
1. VirtualMachinesClient
1. VirtualMachinesRestartHostsFuture

#### New Struct Fields

1. ApplicationGetEndpoint.PrivateIPAddress
1. ApplicationGetHTTPSEndpoint.PrivateIPAddress
1. BillingResponseListResult.VMSizeProperties
1. BillingResponseListResult.VMSizesWithEncryptionAtHost
1. CapabilitiesResult.VmsizeFilters
1. CapabilitiesResult.Vmsizes
1. ClusterCreateProperties.ComputeIsolationProperties
1. ClusterGetProperties.ClusterHdpVersion
1. ClusterGetProperties.ComputeIsolationProperties
1. ClusterGetProperties.ExcludedServicesConfig
1. ClusterGetProperties.StorageProfile
1. ClusterIdentityUserAssignedIdentitiesValue.TenantID
1. ConnectivityEndpoint.PrivateIPAddress
1. KafkaRestProperties.ConfigurationOverride
1. Operation.Properties
1. OperationDisplay.Description
1. Role.EncryptDataDisks
1. Role.VMGroupName
1. StorageAccount.Fileshare
1. StorageAccount.Saskey
1. VMSizeCompatibilityFilter.ComputeIsolationSupported
1. VMSizeCompatibilityFilter.ESPApplied
1. VMSizeCompatibilityFilter.OsType
1. VMSizeCompatibilityFilter.VMSizes
