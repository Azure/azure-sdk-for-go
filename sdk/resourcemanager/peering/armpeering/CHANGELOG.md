# Release History

## 2.0.0 (2025-10-31)
### Breaking Changes

- Function `*LocationsClient.NewListPager` parameter(s) have been changed from `(PeeringLocationsKind, *LocationsClientListOptions)` to `(LocationsKind, *LocationsClientListOptions)`
- Type of `LocationsClientListOptions.DirectPeeringType` has been changed from `*PeeringLocationsDirectPeeringType` to `*LocationsDirectPeeringType`
- Enum `Enum0` has been removed
- Enum `PeeringLocationsDirectPeeringType` has been removed
- Enum `PeeringLocationsKind` has been removed
- Function `*ClientFactory.NewManagementClient` has been removed
- Function `NewManagementClient` has been removed
- Function `*ManagementClient.CheckServiceProviderAvailability` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Resource` has been removed

### Features Added

- New value `ConnectionStateExternalBlocker`, `ConnectionStateTypeChangeInProgress`, `ConnectionStateTypeChangeRequested` added to enum type `ConnectionState`
- New value `DirectPeeringTypePeerProp` added to enum type `DirectPeeringType`
- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New enum type `CheckServiceProviderAvailabilityResponse` with values `CheckServiceProviderAvailabilityResponseAvailable`, `CheckServiceProviderAvailabilityResponseUnavailable`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `LocationsDirectPeeringType` with values `LocationsDirectPeeringTypeCdn`, `LocationsDirectPeeringTypeEdge`, `LocationsDirectPeeringTypeEdgeZoneForOperators`, `LocationsDirectPeeringTypeInternal`, `LocationsDirectPeeringTypeIx`, `LocationsDirectPeeringTypeIxRs`, `LocationsDirectPeeringTypePeerProp`, `LocationsDirectPeeringTypeTransit`, `LocationsDirectPeeringTypeVoice`
- New enum type `LocationsKind` with values `LocationsKindDirect`, `LocationsKindExchange`
- New enum type `Protocol` with values `ProtocolICMP`, `ProtocolNone`, `ProtocolTCP`
- New function `NewClient(string, azcore.TokenCredential, *arm.ClientOptions) (*Client, error)`
- New function `*Client.CheckServiceProviderAvailability(context.Context, CheckServiceProviderAvailabilityInput, *ClientCheckServiceProviderAvailabilityOptions) (ClientCheckServiceProviderAvailabilityResponse, error)`
- New function `*Client.NewCdnPeeringPrefixesClient() *CdnPeeringPrefixesClient`
- New function `*Client.NewConnectionMonitorTestsClient() *ConnectionMonitorTestsClient`
- New function `*Client.NewLegacyPeeringsClient() *LegacyPeeringsClient`
- New function `*Client.NewLocationsClient() *LocationsClient`
- New function `*Client.NewLookingGlassClient() *LookingGlassClient`
- New function `*Client.NewOperationsClient() *OperationsClient`
- New function `*Client.NewPeerAsnsClient() *PeerAsnsClient`
- New function `*Client.NewPeeringsClient() *PeeringsClient`
- New function `*Client.NewPrefixesClient() *PrefixesClient`
- New function `*Client.NewReceivedRoutesClient() *ReceivedRoutesClient`
- New function `*Client.NewRegisteredAsnsClient() *RegisteredAsnsClient`
- New function `*Client.NewRegisteredPrefixesClient() *RegisteredPrefixesClient`
- New function `*Client.NewRpUnbilledPrefixesClient() *RpUnbilledPrefixesClient`
- New function `*Client.NewServiceCountriesClient() *ServiceCountriesClient`
- New function `*Client.NewServiceLocationsClient() *ServiceLocationsClient`
- New function `*Client.NewServiceProvidersClient() *ServiceProvidersClient`
- New function `*Client.NewServicesClient() *ServicesClient`
- New function `*ClientFactory.NewClient() *Client`
- New function `*ClientFactory.NewRpUnbilledPrefixesClient() *RpUnbilledPrefixesClient`
- New function `*RegisteredPrefixesClient.Validate(context.Context, string, string, string, *RegisteredPrefixesClientValidateOptions) (RegisteredPrefixesClientValidateResponse, error)`
- New function `NewRpUnbilledPrefixesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RpUnbilledPrefixesClient, error)`
- New function `*RpUnbilledPrefixesClient.NewListPager(string, string, *RpUnbilledPrefixesClientListOptions) *runtime.Pager[RpUnbilledPrefixesClientListResponse]`
- New struct `ConnectivityProbe`
- New struct `RpUnbilledPrefix`
- New struct `RpUnbilledPrefixListResult`
- New struct `SystemData`
- New field `SystemData` in struct `ConnectionMonitorTest`
- New field `DirectPeeringType` in struct `LegacyPeeringsClientListOptions`
- New field `SystemData` in struct `PeerAsn`
- New field `SystemData` in struct `Peering`
- New field `ConnectivityProbes` in struct `Properties`
- New field `SystemData` in struct `RegisteredAsn`
- New field `SystemData` in struct `RegisteredPrefix`
- New field `SystemData` in struct `Service`
- New field `SystemData` in struct `ServicePrefix`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/peering/armpeering` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).