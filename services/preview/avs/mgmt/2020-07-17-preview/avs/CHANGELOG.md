Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/vmware/resource-manager/readme.md tag: `package-2020-07-17-preview`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *AuthorizationsCreateOrUpdateFuture.Result(AuthorizationsClient) (ExpressRouteAuthorization, error)
1. *AuthorizationsDeleteFuture.Result(AuthorizationsClient) (autorest.Response, error)
1. *ClustersCreateOrUpdateFuture.Result(ClustersClient) (Cluster, error)
1. *ClustersDeleteFuture.Result(ClustersClient) (autorest.Response, error)
1. *ClustersUpdateFuture.Result(ClustersClient) (Cluster, error)
1. *GlobalReachConnectionsCreateOrUpdateFuture.Result(GlobalReachConnectionsClient) (GlobalReachConnection, error)
1. *GlobalReachConnectionsDeleteFuture.Result(GlobalReachConnectionsClient) (autorest.Response, error)
1. *PrivateCloudsCreateOrUpdateFuture.Result(PrivateCloudsClient) (PrivateCloud, error)
1. *PrivateCloudsDeleteFuture.Result(PrivateCloudsClient) (autorest.Response, error)
1. *PrivateCloudsUpdateFuture.Result(PrivateCloudsClient) (PrivateCloud, error)
1. *WorkloadNetworksCreateDhcpFuture.Result(WorkloadNetworksClient) (WorkloadNetworkDhcp, error)
1. *WorkloadNetworksCreatePortMirroringFuture.Result(WorkloadNetworksClient) (WorkloadNetworkPortMirroring, error)
1. *WorkloadNetworksCreateSegmentsFuture.Result(WorkloadNetworksClient) (WorkloadNetworkSegment, error)
1. *WorkloadNetworksCreateVMGroupFuture.Result(WorkloadNetworksClient) (WorkloadNetworkVMGroup, error)
1. *WorkloadNetworksDeleteDhcpFuture.Result(WorkloadNetworksClient) (autorest.Response, error)
1. *WorkloadNetworksDeletePortMirroringFuture.Result(WorkloadNetworksClient) (autorest.Response, error)
1. *WorkloadNetworksDeleteSegmentFuture.Result(WorkloadNetworksClient) (autorest.Response, error)
1. *WorkloadNetworksDeleteVMGroupFuture.Result(WorkloadNetworksClient) (autorest.Response, error)
1. *WorkloadNetworksUpdateDhcpFuture.Result(WorkloadNetworksClient) (WorkloadNetworkDhcp, error)
1. *WorkloadNetworksUpdatePortMirroringFuture.Result(WorkloadNetworksClient) (WorkloadNetworkPortMirroring, error)
1. *WorkloadNetworksUpdateSegmentsFuture.Result(WorkloadNetworksClient) (WorkloadNetworkSegment, error)
1. *WorkloadNetworksUpdateVMGroupFuture.Result(WorkloadNetworksClient) (WorkloadNetworkVMGroup, error)

## Struct Changes

### Removed Struct Fields

1. AuthorizationsCreateOrUpdateFuture.azure.Future
1. AuthorizationsDeleteFuture.azure.Future
1. ClustersCreateOrUpdateFuture.azure.Future
1. ClustersDeleteFuture.azure.Future
1. ClustersUpdateFuture.azure.Future
1. GlobalReachConnectionsCreateOrUpdateFuture.azure.Future
1. GlobalReachConnectionsDeleteFuture.azure.Future
1. PrivateCloudsCreateOrUpdateFuture.azure.Future
1. PrivateCloudsDeleteFuture.azure.Future
1. PrivateCloudsUpdateFuture.azure.Future
1. WorkloadNetworksCreateDhcpFuture.azure.Future
1. WorkloadNetworksCreatePortMirroringFuture.azure.Future
1. WorkloadNetworksCreateSegmentsFuture.azure.Future
1. WorkloadNetworksCreateVMGroupFuture.azure.Future
1. WorkloadNetworksDeleteDhcpFuture.azure.Future
1. WorkloadNetworksDeletePortMirroringFuture.azure.Future
1. WorkloadNetworksDeleteSegmentFuture.azure.Future
1. WorkloadNetworksDeleteVMGroupFuture.azure.Future
1. WorkloadNetworksUpdateDhcpFuture.azure.Future
1. WorkloadNetworksUpdatePortMirroringFuture.azure.Future
1. WorkloadNetworksUpdateSegmentsFuture.azure.Future
1. WorkloadNetworksUpdateVMGroupFuture.azure.Future

## Struct Changes

### New Struct Fields

1. AuthorizationsCreateOrUpdateFuture.Result
1. AuthorizationsCreateOrUpdateFuture.azure.FutureAPI
1. AuthorizationsDeleteFuture.Result
1. AuthorizationsDeleteFuture.azure.FutureAPI
1. ClustersCreateOrUpdateFuture.Result
1. ClustersCreateOrUpdateFuture.azure.FutureAPI
1. ClustersDeleteFuture.Result
1. ClustersDeleteFuture.azure.FutureAPI
1. ClustersUpdateFuture.Result
1. ClustersUpdateFuture.azure.FutureAPI
1. GlobalReachConnectionsCreateOrUpdateFuture.Result
1. GlobalReachConnectionsCreateOrUpdateFuture.azure.FutureAPI
1. GlobalReachConnectionsDeleteFuture.Result
1. GlobalReachConnectionsDeleteFuture.azure.FutureAPI
1. PrivateCloudsCreateOrUpdateFuture.Result
1. PrivateCloudsCreateOrUpdateFuture.azure.FutureAPI
1. PrivateCloudsDeleteFuture.Result
1. PrivateCloudsDeleteFuture.azure.FutureAPI
1. PrivateCloudsUpdateFuture.Result
1. PrivateCloudsUpdateFuture.azure.FutureAPI
1. WorkloadNetworksCreateDhcpFuture.Result
1. WorkloadNetworksCreateDhcpFuture.azure.FutureAPI
1. WorkloadNetworksCreatePortMirroringFuture.Result
1. WorkloadNetworksCreatePortMirroringFuture.azure.FutureAPI
1. WorkloadNetworksCreateSegmentsFuture.Result
1. WorkloadNetworksCreateSegmentsFuture.azure.FutureAPI
1. WorkloadNetworksCreateVMGroupFuture.Result
1. WorkloadNetworksCreateVMGroupFuture.azure.FutureAPI
1. WorkloadNetworksDeleteDhcpFuture.Result
1. WorkloadNetworksDeleteDhcpFuture.azure.FutureAPI
1. WorkloadNetworksDeletePortMirroringFuture.Result
1. WorkloadNetworksDeletePortMirroringFuture.azure.FutureAPI
1. WorkloadNetworksDeleteSegmentFuture.Result
1. WorkloadNetworksDeleteSegmentFuture.azure.FutureAPI
1. WorkloadNetworksDeleteVMGroupFuture.Result
1. WorkloadNetworksDeleteVMGroupFuture.azure.FutureAPI
1. WorkloadNetworksUpdateDhcpFuture.Result
1. WorkloadNetworksUpdateDhcpFuture.azure.FutureAPI
1. WorkloadNetworksUpdatePortMirroringFuture.Result
1. WorkloadNetworksUpdatePortMirroringFuture.azure.FutureAPI
1. WorkloadNetworksUpdateSegmentsFuture.Result
1. WorkloadNetworksUpdateSegmentsFuture.azure.FutureAPI
1. WorkloadNetworksUpdateVMGroupFuture.Result
1. WorkloadNetworksUpdateVMGroupFuture.azure.FutureAPI
