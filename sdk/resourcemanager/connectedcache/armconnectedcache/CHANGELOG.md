# Release History

## 0.2.0 (2025-12-24)
### Breaking Changes

- Enum `CycleType` has been removed
- Function `NewCacheNodesOperationsClient` has been removed
- Function `*CacheNodesOperationsClient.BeginCreateorUpdate` has been removed
- Function `*CacheNodesOperationsClient.Delete` has been removed
- Function `*CacheNodesOperationsClient.Get` has been removed
- Function `*CacheNodesOperationsClient.NewListByResourceGroupPager` has been removed
- Function `*CacheNodesOperationsClient.NewListBySubscriptionPager` has been removed
- Function `*CacheNodesOperationsClient.Update` has been removed
- Function `*ClientFactory.NewCacheNodesOperationsClient` has been removed
- Function `*ClientFactory.NewEnterpriseCustomerOperationsClient` has been removed
- Function `NewEnterpriseCustomerOperationsClient` has been removed
- Function `*EnterpriseCustomerOperationsClient.BeginCreateOrUpdate` has been removed
- Function `*EnterpriseCustomerOperationsClient.Delete` has been removed
- Function `*EnterpriseCustomerOperationsClient.Get` has been removed
- Function `*EnterpriseCustomerOperationsClient.NewListByResourceGroupPager` has been removed
- Function `*EnterpriseCustomerOperationsClient.NewListBySubscriptionPager` has been removed
- Function `*EnterpriseCustomerOperationsClient.Update` has been removed
- Struct `CacheNodeOldResponse` has been removed
- Struct `CacheNodePreviewResource` has been removed
- Struct `CacheNodePreviewResourceListResult` has been removed
- Struct `EnterprisePreviewResource` has been removed
- Struct `EnterprisePreviewResourceListResult` has been removed
- Struct `ErrorAdditionalInfoInfo` has been removed
- Field `ProxyURL`, `UpdateCycleType` of struct `AdditionalCacheNodeProperties` has been removed
- Field `PeeringDbLastUpdateTime` of struct `AdditionalCustomerProperties` has been removed

### Features Added

- Type of `ErrorAdditionalInfo.Info` has been changed from `*ErrorAdditionalInfoInfo` to `any`
- New function `*EnterpriseMccCacheNodesOperationsClient.GetCacheNodeAutoUpdateHistory(ctx context.Context, resourceGroupName string, customerResourceName string, cacheNodeResourceName string, options *EnterpriseMccCacheNodesOperationsClientGetCacheNodeAutoUpdateHistoryOptions) (EnterpriseMccCacheNodesOperationsClientGetCacheNodeAutoUpdateHistoryResponse, error)`
- New function `*EnterpriseMccCacheNodesOperationsClient.GetCacheNodeMccIssueDetailsHistory(ctx context.Context, resourceGroupName string, customerResourceName string, cacheNodeResourceName string, options *EnterpriseMccCacheNodesOperationsClientGetCacheNodeMccIssueDetailsHistoryOptions) (EnterpriseMccCacheNodesOperationsClientGetCacheNodeMccIssueDetailsHistoryResponse, error)`
- New function `*EnterpriseMccCacheNodesOperationsClient.GetCacheNodeTLSCertificateHistory(ctx context.Context, resourceGroupName string, customerResourceName string, cacheNodeResourceName string, options *EnterpriseMccCacheNodesOperationsClientGetCacheNodeTLSCertificateHistoryOptions) (EnterpriseMccCacheNodesOperationsClientGetCacheNodeTLSCertificateHistoryResponse, error)`
- New function `*IspCacheNodesOperationsClient.GetCacheNodeAutoUpdateHistory(ctx context.Context, resourceGroupName string, customerResourceName string, cacheNodeResourceName string, options *IspCacheNodesOperationsClientGetCacheNodeAutoUpdateHistoryOptions) (IspCacheNodesOperationsClientGetCacheNodeAutoUpdateHistoryResponse, error)`
- New function `*IspCacheNodesOperationsClient.GetCacheNodeMccIssueDetailsHistory(ctx context.Context, resourceGroupName string, customerResourceName string, cacheNodeResourceName string, options *IspCacheNodesOperationsClientGetCacheNodeMccIssueDetailsHistoryOptions) (IspCacheNodesOperationsClientGetCacheNodeMccIssueDetailsHistoryResponse, error)`
- New struct `MccCacheNodeAutoUpdateHistory`
- New struct `MccCacheNodeAutoUpdateHistoryProperties`
- New struct `MccCacheNodeAutoUpdateInfo`
- New struct `MccCacheNodeIssueHistory`
- New struct `MccCacheNodeIssueHistoryProperties`
- New struct `MccCacheNodeTLSCertificate`
- New struct `MccCacheNodeTLSCertificateHistory`
- New struct `MccCacheNodeTLSCertificateProperties`
- New struct `MccIssue`
- New field `CreationMethod`, `CurrentTLSCertificate`, `IssuesCount`, `IssuesList`, `LastAutoUpdateInfo`, `TLSStatus` in struct `AdditionalCacheNodeProperties`
- New field `DriveConfiguration`, `ProxyURLConfiguration`, `TLSCertificateProvisioningKey` in struct `CacheNodeInstallProperties`


## 0.1.0 (2024-11-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedcache/armconnectedcache` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).