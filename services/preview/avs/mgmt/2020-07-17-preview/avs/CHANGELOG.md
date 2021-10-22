# Unreleased

## Breaking Changes

### Removed Constants

1. NsxtAdminRotateEnum.OnetimeRotate
1. VcsaAdminRotateEnum.VcsaAdminRotateEnumOnetimeRotate

### Removed Funcs

1. PossibleNsxtAdminRotateEnumValues() []NsxtAdminRotateEnum
1. PossibleVcsaAdminRotateEnumValues() []VcsaAdminRotateEnum

### Struct Changes

#### Removed Struct Fields

1. PrivateCloudUpdateProperties.NsxtPassword
1. PrivateCloudUpdateProperties.VcenterPassword

### Signature Changes

#### Const Types

1. Cancelled changed type from ClusterProvisioningState to AddonProvisioningState
1. Deleting changed type from ClusterProvisioningState to AddonProvisioningState
1. Failed changed type from ClusterProvisioningState to AddonProvisioningState
1. Succeeded changed type from ClusterProvisioningState to AddonProvisioningState
1. Updating changed type from ClusterProvisioningState to AddonProvisioningState

#### Struct Fields

1. PrivateCloudProperties.NsxtPassword changed type from NsxtAdminRotateEnum to *string
1. PrivateCloudProperties.VcenterPassword changed type from VcsaAdminRotateEnum to *string

## Additive Changes

### New Constants

1. AddonType.SRM
1. AddonType.VR
1. ClusterProvisioningState.ClusterProvisioningStateCancelled
1. ClusterProvisioningState.ClusterProvisioningStateDeleting
1. ClusterProvisioningState.ClusterProvisioningStateFailed
1. ClusterProvisioningState.ClusterProvisioningStateSucceeded
1. ClusterProvisioningState.ClusterProvisioningStateUpdating
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

1. *Addon.UnmarshalJSON([]byte) error
1. *AddonListIterator.Next() error
1. *AddonListIterator.NextWithContext(context.Context) error
1. *AddonListPage.Next() error
1. *AddonListPage.NextWithContext(context.Context) error
1. *AddonUpdate.UnmarshalJSON([]byte) error
1. *AddonUpdateProperties.UnmarshalJSON([]byte) error
1. *AddonsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *AddonsDeleteFuture.UnmarshalJSON([]byte) error
1. *PrivateCloudsRotateNsxtPasswordFuture.UnmarshalJSON([]byte) error
1. *PrivateCloudsRotateVcenterPasswordFuture.UnmarshalJSON([]byte) error
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
1. *WorkloadNetworksCreateDNSServiceFuture.UnmarshalJSON([]byte) error
1. *WorkloadNetworksCreateDNSZoneFuture.UnmarshalJSON([]byte) error
1. *WorkloadNetworksDeleteDNSServiceFuture.UnmarshalJSON([]byte) error
1. *WorkloadNetworksDeleteDNSZoneFuture.UnmarshalJSON([]byte) error
1. *WorkloadNetworksUpdateDNSServiceFuture.UnmarshalJSON([]byte) error
1. *WorkloadNetworksUpdateDNSZoneFuture.UnmarshalJSON([]byte) error
1. Addon.MarshalJSON() ([]byte, error)
1. AddonList.IsEmpty() bool
1. AddonList.MarshalJSON() ([]byte, error)
1. AddonListIterator.NotDone() bool
1. AddonListIterator.Response() AddonList
1. AddonListIterator.Value() Addon
1. AddonListPage.NotDone() bool
1. AddonListPage.Response() AddonList
1. AddonListPage.Values() []Addon
1. AddonProperties.MarshalJSON() ([]byte, error)
1. AddonUpdate.MarshalJSON() ([]byte, error)
1. AddonUpdateProperties.MarshalJSON() ([]byte, error)
1. AddonsClient.CreateOrUpdate(context.Context, string, string, string, Addon) (AddonsCreateOrUpdateFuture, error)
1. AddonsClient.CreateOrUpdatePreparer(context.Context, string, string, string, Addon) (*http.Request, error)
1. AddonsClient.CreateOrUpdateResponder(*http.Response) (Addon, error)
1. AddonsClient.CreateOrUpdateSender(*http.Request) (AddonsCreateOrUpdateFuture, error)
1. AddonsClient.Delete(context.Context, string, string, string) (AddonsDeleteFuture, error)
1. AddonsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. AddonsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. AddonsClient.DeleteSender(*http.Request) (AddonsDeleteFuture, error)
1. AddonsClient.Get(context.Context, string, string, string) (Addon, error)
1. AddonsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. AddonsClient.GetResponder(*http.Response) (Addon, error)
1. AddonsClient.GetSender(*http.Request) (*http.Response, error)
1. AddonsClient.List(context.Context, string, string) (AddonListPage, error)
1. AddonsClient.ListComplete(context.Context, string, string) (AddonListIterator, error)
1. AddonsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. AddonsClient.ListResponder(*http.Response) (AddonList, error)
1. AddonsClient.ListSender(*http.Request) (*http.Response, error)
1. NewAddonListIterator(AddonListPage) AddonListIterator
1. NewAddonListPage(AddonList, func(context.Context, AddonList) (AddonList, error)) AddonListPage
1. NewAddonsClient(string) AddonsClient
1. NewAddonsClientWithBaseURI(string, string) AddonsClient
1. NewWorkloadNetworkDNSServicesListIterator(WorkloadNetworkDNSServicesListPage) WorkloadNetworkDNSServicesListIterator
1. NewWorkloadNetworkDNSServicesListPage(WorkloadNetworkDNSServicesList, func(context.Context, WorkloadNetworkDNSServicesList) (WorkloadNetworkDNSServicesList, error)) WorkloadNetworkDNSServicesListPage
1. NewWorkloadNetworkDNSZonesListIterator(WorkloadNetworkDNSZonesListPage) WorkloadNetworkDNSZonesListIterator
1. NewWorkloadNetworkDNSZonesListPage(WorkloadNetworkDNSZonesList, func(context.Context, WorkloadNetworkDNSZonesList) (WorkloadNetworkDNSZonesList, error)) WorkloadNetworkDNSZonesListPage
1. PossibleAddonProvisioningStateValues() []AddonProvisioningState
1. PossibleAddonTypeValues() []AddonType
1. PossibleDNSServiceLogLevelEnumValues() []DNSServiceLogLevelEnum
1. PossibleDNSServiceStatusEnumValues() []DNSServiceStatusEnum
1. PossibleWorkloadNetworkDNSServiceProvisioningStateValues() []WorkloadNetworkDNSServiceProvisioningState
1. PossibleWorkloadNetworkDNSZoneProvisioningStateValues() []WorkloadNetworkDNSZoneProvisioningState
1. PrivateCloudsClient.RotateNsxtPassword(context.Context, string, string) (PrivateCloudsRotateNsxtPasswordFuture, error)
1. PrivateCloudsClient.RotateNsxtPasswordPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateCloudsClient.RotateNsxtPasswordResponder(*http.Response) (autorest.Response, error)
1. PrivateCloudsClient.RotateNsxtPasswordSender(*http.Request) (PrivateCloudsRotateNsxtPasswordFuture, error)
1. PrivateCloudsClient.RotateVcenterPassword(context.Context, string, string) (PrivateCloudsRotateVcenterPasswordFuture, error)
1. PrivateCloudsClient.RotateVcenterPasswordPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateCloudsClient.RotateVcenterPasswordResponder(*http.Response) (autorest.Response, error)
1. PrivateCloudsClient.RotateVcenterPasswordSender(*http.Request) (PrivateCloudsRotateVcenterPasswordFuture, error)
1. WorkloadNetworkDNSService.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSServiceProperties.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSServicesList.IsEmpty() bool
1. WorkloadNetworkDNSServicesList.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSServicesListIterator.NotDone() bool
1. WorkloadNetworkDNSServicesListIterator.Response() WorkloadNetworkDNSServicesList
1. WorkloadNetworkDNSServicesListIterator.Value() WorkloadNetworkDNSService
1. WorkloadNetworkDNSServicesListPage.NotDone() bool
1. WorkloadNetworkDNSServicesListPage.Response() WorkloadNetworkDNSServicesList
1. WorkloadNetworkDNSServicesListPage.Values() []WorkloadNetworkDNSService
1. WorkloadNetworkDNSZone.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSZoneProperties.MarshalJSON() ([]byte, error)
1. WorkloadNetworkDNSZonesList.IsEmpty() bool
1. WorkloadNetworkDNSZonesList.MarshalJSON() ([]byte, error)
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

### Struct Changes

#### New Structs

1. Addon
1. AddonList
1. AddonListIterator
1. AddonListPage
1. AddonProperties
1. AddonSrmProperties
1. AddonUpdate
1. AddonUpdateProperties
1. AddonsClient
1. AddonsCreateOrUpdateFuture
1. AddonsDeleteFuture
1. PrivateCloudsRotateNsxtPasswordFuture
1. PrivateCloudsRotateVcenterPasswordFuture
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

#### New Struct Fields

1. MetricDimension.InternalName
1. MetricDimension.ToBeExportedForShoebox
