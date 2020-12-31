Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

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

### New Constants

1. DNSServiceLogLevelEnum.DEBUG
1. DNSServiceLogLevelEnum.ERROR
1. DNSServiceLogLevelEnum.FATAL
1. DNSServiceLogLevelEnum.INFO
1. DNSServiceLogLevelEnum.WARNING
1. DNSServiceStatusEnum.FAILURE
1. DNSServiceStatusEnum.SUCCESS
1. WorkloadNetworkDNSServiceProvisioningState.WorkloadNetworkDNSServiceProvisioningStateBuilding
1. WorkloadNetworkDNSServiceProvisioningState.WorkloadNetworkDNSServiceProvisioningStateDeleting
1. WorkloadNetworkDNSServiceProvisioningState.WorkloadNetworkDNSServiceProvisioningStateFailed
1. WorkloadNetworkDNSServiceProvisioningState.WorkloadNetworkDNSServiceProvisioningStateSucceeded
1. WorkloadNetworkDNSServiceProvisioningState.WorkloadNetworkDNSServiceProvisioningStateUpdating
1. WorkloadNetworkDNSZoneProvisioningState.WorkloadNetworkDNSZoneProvisioningStateBuilding
1. WorkloadNetworkDNSZoneProvisioningState.WorkloadNetworkDNSZoneProvisioningStateDeleting
1. WorkloadNetworkDNSZoneProvisioningState.WorkloadNetworkDNSZoneProvisioningStateFailed
1. WorkloadNetworkDNSZoneProvisioningState.WorkloadNetworkDNSZoneProvisioningStateSucceeded
1. WorkloadNetworkDNSZoneProvisioningState.WorkloadNetworkDNSZoneProvisioningStateUpdating

### New Funcs

1. *WorkloadNetworkDNSService.UnmarshalJSON([]byte) error
1. *WorkloadNetworkDNSServicesListIterator.Next() error
1. *WorkloadNetworkDNSServicesListIterator.NextWithContext(context.Context) error
1. *WorkloadNetworkDNSServicesListPage.Next() error
1. *WorkloadNetworkDNSServicesListPage.NextWithContext(context.Context) error
1. *WorkloadNetworkDNSZone.UnmarshalJSON([]byte) error
1. *WorkloadNetworkDNSZonesListIterator.Next() error
1. *WorkloadNetworkDNSZonesListIterator.NextWithContext(context.Context) error
1. *WorkloadNetworkDNSZonesListPage.Next() error
1. *WorkloadNetworkDNSZonesListPage.NextWithContext(context.Context) error
1. NewWorkloadNetworkDNSServicesListIterator(WorkloadNetworkDNSServicesListPage) WorkloadNetworkDNSServicesListIterator
1. NewWorkloadNetworkDNSServicesListPage(WorkloadNetworkDNSServicesList, func(context.Context, WorkloadNetworkDNSServicesList) (WorkloadNetworkDNSServicesList, error)) WorkloadNetworkDNSServicesListPage
1. NewWorkloadNetworkDNSZonesListIterator(WorkloadNetworkDNSZonesListPage) WorkloadNetworkDNSZonesListIterator
1. NewWorkloadNetworkDNSZonesListPage(WorkloadNetworkDNSZonesList, func(context.Context, WorkloadNetworkDNSZonesList) (WorkloadNetworkDNSZonesList, error)) WorkloadNetworkDNSZonesListPage
1. PossibleDNSServiceLogLevelEnumValues() []DNSServiceLogLevelEnum
1. PossibleDNSServiceStatusEnumValues() []DNSServiceStatusEnum
1. PossibleWorkloadNetworkDNSServiceProvisioningStateValues() []WorkloadNetworkDNSServiceProvisioningState
1. PossibleWorkloadNetworkDNSZoneProvisioningStateValues() []WorkloadNetworkDNSZoneProvisioningState
1. WorkloadNetworkDNSService.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSServiceProperties.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSServicesList.IsEmpty() bool
1. WorkloadNetworkDNSServicesListIterator.NotDone() bool
1. WorkloadNetworkDNSServicesListIterator.Response() WorkloadNetworkDNSServicesList
1. WorkloadNetworkDNSServicesListIterator.Value() WorkloadNetworkDNSService
1. WorkloadNetworkDNSServicesListPage.NotDone() bool
1. WorkloadNetworkDNSServicesListPage.Response() WorkloadNetworkDNSServicesList
1. WorkloadNetworkDNSServicesListPage.Values() []WorkloadNetworkDNSService
1. WorkloadNetworkDNSZone.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSZoneProperties.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSZonesList.IsEmpty() bool
1. WorkloadNetworkDNSZonesListIterator.NotDone() bool
1. WorkloadNetworkDNSZonesListIterator.Response() WorkloadNetworkDNSZonesList
1. WorkloadNetworkDNSZonesListIterator.Value() WorkloadNetworkDNSZone
1. WorkloadNetworkDNSZonesListPage.NotDone() bool
1. WorkloadNetworkDNSZonesListPage.Response() WorkloadNetworkDNSZonesList
1. WorkloadNetworkDNSZonesListPage.Values() []WorkloadNetworkDNSZone
1. WorkloadNetworksClient.CreateDNSService(context.Context, string, string, string, WorkloadNetworkDNSService) (WorkloadNetworksCreateDNSServiceFuture, error)
1. WorkloadNetworksClient.CreateDNSServicePreparer(context.Context, string, string, string, WorkloadNetworkDNSService) (*http.Request, error)
1. WorkloadNetworksClient.CreateDNSServiceResponder(*http.Response) (WorkloadNetworkDNSService, error)
1. WorkloadNetworksClient.CreateDNSServiceSender(*http.Request) (WorkloadNetworksCreateDNSServiceFuture, error)
1. WorkloadNetworksClient.CreateDNSZone(context.Context, string, string, string, WorkloadNetworkDNSZone) (WorkloadNetworksCreateDNSZoneFuture, error)
1. WorkloadNetworksClient.CreateDNSZonePreparer(context.Context, string, string, string, WorkloadNetworkDNSZone) (*http.Request, error)
1. WorkloadNetworksClient.CreateDNSZoneResponder(*http.Response) (WorkloadNetworkDNSZone, error)
1. WorkloadNetworksClient.CreateDNSZoneSender(*http.Request) (WorkloadNetworksCreateDNSZoneFuture, error)
1. WorkloadNetworksClient.DeleteDNSService(context.Context, string, string, string) (WorkloadNetworksDeleteDNSServiceFuture, error)
1. WorkloadNetworksClient.DeleteDNSServicePreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkloadNetworksClient.DeleteDNSServiceResponder(*http.Response) (autorest.Response, error)
1. WorkloadNetworksClient.DeleteDNSServiceSender(*http.Request) (WorkloadNetworksDeleteDNSServiceFuture, error)
1. WorkloadNetworksClient.DeleteDNSZone(context.Context, string, string, string) (WorkloadNetworksDeleteDNSZoneFuture, error)
1. WorkloadNetworksClient.DeleteDNSZonePreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkloadNetworksClient.DeleteDNSZoneResponder(*http.Response) (autorest.Response, error)
1. WorkloadNetworksClient.DeleteDNSZoneSender(*http.Request) (WorkloadNetworksDeleteDNSZoneFuture, error)
1. WorkloadNetworksClient.GetDNSService(context.Context, string, string, string) (WorkloadNetworkDNSService, error)
1. WorkloadNetworksClient.GetDNSServicePreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkloadNetworksClient.GetDNSServiceResponder(*http.Response) (WorkloadNetworkDNSService, error)
1. WorkloadNetworksClient.GetDNSServiceSender(*http.Request) (*http.Response, error)
1. WorkloadNetworksClient.GetDNSZone(context.Context, string, string, string) (WorkloadNetworkDNSZone, error)
1. WorkloadNetworksClient.GetDNSZonePreparer(context.Context, string, string, string) (*http.Request, error)
1. WorkloadNetworksClient.GetDNSZoneResponder(*http.Response) (WorkloadNetworkDNSZone, error)
1. WorkloadNetworksClient.GetDNSZoneSender(*http.Request) (*http.Response, error)
1. WorkloadNetworksClient.ListDNSServices(context.Context, string, string) (WorkloadNetworkDNSServicesListPage, error)
1. WorkloadNetworksClient.ListDNSServicesComplete(context.Context, string, string) (WorkloadNetworkDNSServicesListIterator, error)
1. WorkloadNetworksClient.ListDNSServicesPreparer(context.Context, string, string) (*http.Request, error)
1. WorkloadNetworksClient.ListDNSServicesResponder(*http.Response) (WorkloadNetworkDNSServicesList, error)
1. WorkloadNetworksClient.ListDNSServicesSender(*http.Request) (*http.Response, error)
1. WorkloadNetworksClient.ListDNSZones(context.Context, string, string) (WorkloadNetworkDNSZonesListPage, error)
1. WorkloadNetworksClient.ListDNSZonesComplete(context.Context, string, string) (WorkloadNetworkDNSZonesListIterator, error)
1. WorkloadNetworksClient.ListDNSZonesPreparer(context.Context, string, string) (*http.Request, error)
1. WorkloadNetworksClient.ListDNSZonesResponder(*http.Response) (WorkloadNetworkDNSZonesList, error)
1. WorkloadNetworksClient.ListDNSZonesSender(*http.Request) (*http.Response, error)
1. WorkloadNetworksClient.UpdateDNSService(context.Context, string, string, string, WorkloadNetworkDNSService) (WorkloadNetworksUpdateDNSServiceFuture, error)
1. WorkloadNetworksClient.UpdateDNSServicePreparer(context.Context, string, string, string, WorkloadNetworkDNSService) (*http.Request, error)
1. WorkloadNetworksClient.UpdateDNSServiceResponder(*http.Response) (WorkloadNetworkDNSService, error)
1. WorkloadNetworksClient.UpdateDNSServiceSender(*http.Request) (WorkloadNetworksUpdateDNSServiceFuture, error)
1. WorkloadNetworksClient.UpdateDNSZone(context.Context, string, string, string, WorkloadNetworkDNSZone) (WorkloadNetworksUpdateDNSZoneFuture, error)
1. WorkloadNetworksClient.UpdateDNSZonePreparer(context.Context, string, string, string, WorkloadNetworkDNSZone) (*http.Request, error)
1. WorkloadNetworksClient.UpdateDNSZoneResponder(*http.Response) (WorkloadNetworkDNSZone, error)
1. WorkloadNetworksClient.UpdateDNSZoneSender(*http.Request) (WorkloadNetworksUpdateDNSZoneFuture, error)

## Struct Changes

### New Structs

1. WorkloadNetworkDNSService
1. WorkloadNetworkDNSServiceProperties
1. WorkloadNetworkDNSServicesList
1. WorkloadNetworkDNSServicesListIterator
1. WorkloadNetworkDNSServicesListPage
1. WorkloadNetworkDNSZone
1. WorkloadNetworkDNSZoneProperties
1. WorkloadNetworkDNSZonesList
1. WorkloadNetworkDNSZonesListIterator
1. WorkloadNetworkDNSZonesListPage
1. WorkloadNetworksCreateDNSServiceFuture
1. WorkloadNetworksCreateDNSZoneFuture
1. WorkloadNetworksDeleteDNSServiceFuture
1. WorkloadNetworksDeleteDNSZoneFuture
1. WorkloadNetworksUpdateDNSServiceFuture
1. WorkloadNetworksUpdateDNSZoneFuture

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
