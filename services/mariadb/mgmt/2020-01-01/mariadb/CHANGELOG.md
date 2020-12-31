Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *ConfigurationsCreateOrUpdateFuture.Result(ConfigurationsClient) (Configuration, error)
1. *CreateRecommendedActionSessionFuture.Result(BaseClient) (autorest.Response, error)
1. *DatabasesCreateOrUpdateFuture.Result(DatabasesClient) (Database, error)
1. *DatabasesDeleteFuture.Result(DatabasesClient) (autorest.Response, error)
1. *FirewallRulesCreateOrUpdateFuture.Result(FirewallRulesClient) (FirewallRule, error)
1. *FirewallRulesDeleteFuture.Result(FirewallRulesClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsCreateOrUpdateFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsUpdateTagsFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *ServerSecurityAlertPoliciesCreateOrUpdateFuture.Result(ServerSecurityAlertPoliciesClient) (ServerSecurityAlertPolicy, error)
1. *ServersCreateFuture.Result(ServersClient) (Server, error)
1. *ServersDeleteFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersRestartFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersStartFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersStopFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersUpdateFuture.Result(ServersClient) (Server, error)
1. *VirtualNetworkRulesCreateOrUpdateFuture.Result(VirtualNetworkRulesClient) (VirtualNetworkRule, error)
1. *VirtualNetworkRulesDeleteFuture.Result(VirtualNetworkRulesClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. ConfigurationsCreateOrUpdateFuture.azure.Future
1. CreateRecommendedActionSessionFuture.azure.Future
1. DatabasesCreateOrUpdateFuture.azure.Future
1. DatabasesDeleteFuture.azure.Future
1. FirewallRulesCreateOrUpdateFuture.azure.Future
1. FirewallRulesDeleteFuture.azure.Future
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future
1. PrivateEndpointConnectionsUpdateTagsFuture.azure.Future
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.Future
1. ServersCreateFuture.azure.Future
1. ServersDeleteFuture.azure.Future
1. ServersRestartFuture.azure.Future
1. ServersStartFuture.azure.Future
1. ServersStopFuture.azure.Future
1. ServersUpdateFuture.azure.Future
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.Future
1. VirtualNetworkRulesDeleteFuture.azure.Future

### New Constants

1. QueryPerformanceInsightResetDataResultState.QueryPerformanceInsightResetDataResultStateFailed
1. QueryPerformanceInsightResetDataResultState.QueryPerformanceInsightResetDataResultStateSucceeded

### New Funcs

1. *RecoverableServerResource.UnmarshalJSON([]byte) error
1. BaseClient.ResetQueryPerformanceInsightData(context.Context, string, string) (QueryPerformanceInsightResetDataResult, error)
1. BaseClient.ResetQueryPerformanceInsightDataPreparer(context.Context, string, string) (*http.Request, error)
1. BaseClient.ResetQueryPerformanceInsightDataResponder(*http.Response) (QueryPerformanceInsightResetDataResult, error)
1. BaseClient.ResetQueryPerformanceInsightDataSender(*http.Request) (*http.Response, error)
1. NewRecoverableServersClient(string) RecoverableServersClient
1. NewRecoverableServersClientWithBaseURI(string, string) RecoverableServersClient
1. NewServerBasedPerformanceTierClient(string) ServerBasedPerformanceTierClient
1. NewServerBasedPerformanceTierClientWithBaseURI(string, string) ServerBasedPerformanceTierClient
1. NewServerParametersClient(string) ServerParametersClient
1. NewServerParametersClientWithBaseURI(string, string) ServerParametersClient
1. PossibleQueryPerformanceInsightResetDataResultStateValues() []QueryPerformanceInsightResetDataResultState
1. RecoverableServerResource.MarshalJSON() ([]byte, error)
1. RecoverableServersClient.Get(context.Context, string, string) (RecoverableServerResource, error)
1. RecoverableServersClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. RecoverableServersClient.GetResponder(*http.Response) (RecoverableServerResource, error)
1. RecoverableServersClient.GetSender(*http.Request) (*http.Response, error)
1. ServerBasedPerformanceTierClient.List(context.Context, string, string) (PerformanceTierListResult, error)
1. ServerBasedPerformanceTierClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ServerBasedPerformanceTierClient.ListResponder(*http.Response) (PerformanceTierListResult, error)
1. ServerBasedPerformanceTierClient.ListSender(*http.Request) (*http.Response, error)
1. ServerParametersClient.ListUpdateConfigurations(context.Context, string, string, ConfigurationListResult) (ServerParametersListUpdateConfigurationsFuture, error)
1. ServerParametersClient.ListUpdateConfigurationsPreparer(context.Context, string, string, ConfigurationListResult) (*http.Request, error)
1. ServerParametersClient.ListUpdateConfigurationsResponder(*http.Response) (ConfigurationListResult, error)
1. ServerParametersClient.ListUpdateConfigurationsSender(*http.Request) (ServerParametersListUpdateConfigurationsFuture, error)

## Struct Changes

### New Structs

1. QueryPerformanceInsightResetDataResult
1. RecoverableServerProperties
1. RecoverableServerResource
1. RecoverableServersClient
1. ServerBasedPerformanceTierClient
1. ServerParametersClient
1. ServerParametersListUpdateConfigurationsFuture

### New Struct Fields

1. ConfigurationsCreateOrUpdateFuture.Result
1. ConfigurationsCreateOrUpdateFuture.azure.FutureAPI
1. CreateRecommendedActionSessionFuture.Result
1. CreateRecommendedActionSessionFuture.azure.FutureAPI
1. DatabasesCreateOrUpdateFuture.Result
1. DatabasesCreateOrUpdateFuture.azure.FutureAPI
1. DatabasesDeleteFuture.Result
1. DatabasesDeleteFuture.azure.FutureAPI
1. FirewallRulesCreateOrUpdateFuture.Result
1. FirewallRulesCreateOrUpdateFuture.azure.FutureAPI
1. FirewallRulesDeleteFuture.Result
1. FirewallRulesDeleteFuture.azure.FutureAPI
1. PerformanceTierProperties.MaxBackupRetentionDays
1. PerformanceTierProperties.MaxLargeStorageMB
1. PerformanceTierProperties.MaxStorageMB
1. PerformanceTierProperties.MinBackupRetentionDays
1. PerformanceTierProperties.MinLargeStorageMB
1. PerformanceTierProperties.MinStorageMB
1. PrivateEndpointConnectionsCreateOrUpdateFuture.Result
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
1. PrivateEndpointConnectionsUpdateTagsFuture.Result
1. PrivateEndpointConnectionsUpdateTagsFuture.azure.FutureAPI
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.Result
1. ServerSecurityAlertPoliciesCreateOrUpdateFuture.azure.FutureAPI
1. ServersCreateFuture.Result
1. ServersCreateFuture.azure.FutureAPI
1. ServersDeleteFuture.Result
1. ServersDeleteFuture.azure.FutureAPI
1. ServersRestartFuture.Result
1. ServersRestartFuture.azure.FutureAPI
1. ServersStartFuture.Result
1. ServersStartFuture.azure.FutureAPI
1. ServersStopFuture.Result
1. ServersStopFuture.azure.FutureAPI
1. ServersUpdateFuture.Result
1. ServersUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesCreateOrUpdateFuture.Result
1. VirtualNetworkRulesCreateOrUpdateFuture.azure.FutureAPI
1. VirtualNetworkRulesDeleteFuture.Result
1. VirtualNetworkRulesDeleteFuture.azure.FutureAPI
