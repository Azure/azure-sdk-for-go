# Release History

## 1.1.0-beta.1 (2026-10-09)
### Features Added

- New enum type `MetricConfigurationType` with values `MetricConfigurationTypeAggregated`, `MetricConfigurationTypeRaw`
- New enum type `TraceMetricsState` with values `TraceMetricsStateDisabled`, `TraceMetricsStateEnabled`
- New function `*ClientFactory.NewMetricConfigurationsClient() *MetricConfigurationsClient`
- New function `*ClientFactory.NewMetricNamespacesClient() *MetricNamespacesClient`
- New function `*ClientFactory.NewTraceAssociationsAtResourceGroupClient() *TraceAssociationsAtResourceGroupClient`
- New function `*ClientFactory.NewTraceAssociationsAtSubscriptionClient() *TraceAssociationsAtSubscriptionClient`
- New function `*ClientFactory.NewTraceAssociationsClient() *TraceAssociationsClient`
- New function `*ClientFactory.NewTraceContainersClient() *TraceContainersClient`
- New function `NewMetricConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MetricConfigurationsClient, error)`
- New function `*MetricConfigurationsClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, encodedMetricNamespace string, encodedMetricName string, resource MetricConfigurationResource, options *MetricConfigurationsClientCreateOrUpdateOptions) (MetricConfigurationsClientCreateOrUpdateResponse, error)`
- New function `*MetricConfigurationsClient.Delete(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, encodedMetricNamespace string, encodedMetricName string, options *MetricConfigurationsClientDeleteOptions) (MetricConfigurationsClientDeleteResponse, error)`
- New function `*MetricConfigurationsClient.Get(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, encodedMetricNamespace string, encodedMetricName string, options *MetricConfigurationsClientGetOptions) (MetricConfigurationsClientGetResponse, error)`
- New function `*MetricConfigurationsClient.NewListByMetricNamespacePager(resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, encodedMetricNamespace string, options *MetricConfigurationsClientListByMetricNamespaceOptions) *runtime.Pager[MetricConfigurationsClientListByMetricNamespaceResponse]`
- New function `*MetricConfigurationsClient.NewListByMetricsContainerPager(resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, options *MetricConfigurationsClientListByMetricsContainerOptions) *runtime.Pager[MetricConfigurationsClientListByMetricsContainerResponse]`
- New function `NewMetricNamespacesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MetricNamespacesClient, error)`
- New function `*MetricNamespacesClient.Get(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, encodedMetricNamespace string, options *MetricNamespacesClientGetOptions) (MetricNamespacesClientGetResponse, error)`
- New function `*MetricNamespacesClient.NewListByMetricsContainerPager(resourceGroupName string, azureMonitorWorkspaceName string, metricsContainerName string, options *MetricNamespacesClientListByMetricsContainerOptions) *runtime.Pager[MetricNamespacesClientListByMetricsContainerResponse]`
- New function `NewTraceAssociationsAtResourceGroupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TraceAssociationsAtResourceGroupClient, error)`
- New function `*TraceAssociationsAtResourceGroupClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, resource TraceAssociationResource, options *TraceAssociationsAtResourceGroupClientCreateOrUpdateOptions) (TraceAssociationsAtResourceGroupClientCreateOrUpdateResponse, error)`
- New function `*TraceAssociationsAtResourceGroupClient.Delete(ctx context.Context, resourceGroupName string, options *TraceAssociationsAtResourceGroupClientDeleteOptions) (TraceAssociationsAtResourceGroupClientDeleteResponse, error)`
- New function `*TraceAssociationsAtResourceGroupClient.Get(ctx context.Context, resourceGroupName string, options *TraceAssociationsAtResourceGroupClientGetOptions) (TraceAssociationsAtResourceGroupClientGetResponse, error)`
- New function `*TraceAssociationsAtResourceGroupClient.NewListPager(resourceGroupName string, options *TraceAssociationsAtResourceGroupClientListOptions) *runtime.Pager[TraceAssociationsAtResourceGroupClientListResponse]`
- New function `NewTraceAssociationsAtSubscriptionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TraceAssociationsAtSubscriptionClient, error)`
- New function `*TraceAssociationsAtSubscriptionClient.CreateOrUpdate(ctx context.Context, resource TraceAssociationResource, options *TraceAssociationsAtSubscriptionClientCreateOrUpdateOptions) (TraceAssociationsAtSubscriptionClientCreateOrUpdateResponse, error)`
- New function `*TraceAssociationsAtSubscriptionClient.Delete(ctx context.Context, options *TraceAssociationsAtSubscriptionClientDeleteOptions) (TraceAssociationsAtSubscriptionClientDeleteResponse, error)`
- New function `*TraceAssociationsAtSubscriptionClient.Get(ctx context.Context, options *TraceAssociationsAtSubscriptionClientGetOptions) (TraceAssociationsAtSubscriptionClientGetResponse, error)`
- New function `*TraceAssociationsAtSubscriptionClient.NewListPager(options *TraceAssociationsAtSubscriptionClientListOptions) *runtime.Pager[TraceAssociationsAtSubscriptionClientListResponse]`
- New function `NewTraceAssociationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TraceAssociationsClient, error)`
- New function `*TraceAssociationsClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, providerName string, providerType string, resourceName string, resource TraceAssociationResource, options *TraceAssociationsClientCreateOrUpdateOptions) (TraceAssociationsClientCreateOrUpdateResponse, error)`
- New function `*TraceAssociationsClient.Delete(ctx context.Context, resourceGroupName string, providerName string, providerType string, resourceName string, options *TraceAssociationsClientDeleteOptions) (TraceAssociationsClientDeleteResponse, error)`
- New function `*TraceAssociationsClient.Get(ctx context.Context, resourceGroupName string, providerName string, providerType string, resourceName string, options *TraceAssociationsClientGetOptions) (TraceAssociationsClientGetResponse, error)`
- New function `*TraceAssociationsClient.NewListPager(resourceGroupName string, providerName string, providerType string, resourceName string, options *TraceAssociationsClientListOptions) *runtime.Pager[TraceAssociationsClientListResponse]`
- New function `NewTraceContainersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TraceContainersClient, error)`
- New function `*TraceContainersClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, resource TraceContainerResource, options *TraceContainersClientBeginCreateOrUpdateOptions) (*runtime.Poller[TraceContainersClientCreateOrUpdateResponse], error)`
- New function `*TraceContainersClient.Delete(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, options *TraceContainersClientDeleteOptions) (TraceContainersClientDeleteResponse, error)`
- New function `*TraceContainersClient.Get(ctx context.Context, resourceGroupName string, azureMonitorWorkspaceName string, options *TraceContainersClientGetOptions) (TraceContainersClientGetResponse, error)`
- New function `*TraceContainersClient.NewListByAzureMonitorWorkspacePager(resourceGroupName string, azureMonitorWorkspaceName string, options *TraceContainersClientListByAzureMonitorWorkspaceOptions) *runtime.Pager[TraceContainersClientListByAzureMonitorWorkspaceResponse]`
- New struct `AzureMonitorWorkspaceActions`
- New struct `AzureMonitorWorkspaceEndpoints`
- New struct `DefaultActionGroupResource`
- New struct `MetricAggregationConfiguration`
- New struct `MetricAggregationFunctions`
- New struct `MetricConfigurationProperties`
- New struct `MetricConfigurationResource`
- New struct `MetricConfigurationResourceListResult`
- New struct `MetricNamespaceProperties`
- New struct `MetricNamespaceResource`
- New struct `MetricNamespaceResourceListResult`
- New struct `MetricsLimits`
- New struct `PagedMetricConfigurationResource`
- New struct `TraceAssociation`
- New struct `TraceAssociationResource`
- New struct `TraceAssociationResourceListResult`
- New struct `TraceContainer`
- New struct `TraceContainerResource`
- New struct `TraceContainerResourceListResult`
- New field `Actions`, `Endpoints` in struct `AzureMonitorWorkspace`
- New field `Limits` in struct `MetricsContainer`


## 1.0.0 (2026-06-12)

### Other Changes

- General availability of the `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitorworkspaces` package.

## 0.1.0 (2026-05-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitorworkspaces` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).