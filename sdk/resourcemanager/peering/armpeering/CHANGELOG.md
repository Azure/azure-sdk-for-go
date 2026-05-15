# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Resource` has been removed

### Features Added

- New value `ConnectionStateExternalBlocker`, `ConnectionStateTypeChangeInProgress`, `ConnectionStateTypeChangeRequested` added to enum type `ConnectionState`
- New value `DirectPeeringTypePeerProp` added to enum type `DirectPeeringType`
- New value `PeeringLocationsDirectPeeringTypePeerProp` added to enum type `PeeringLocationsDirectPeeringType`
- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Protocol` with values `ProtocolICMP`, `ProtocolNone`, `ProtocolTCP`
- New function `*ClientFactory.NewRpUnbilledPrefixesClient() *RpUnbilledPrefixesClient`
- New function `*RegisteredPrefixesClient.Validate(ctx context.Context, resourceGroupName string, peeringName string, registeredPrefixName string, options *RegisteredPrefixesClientValidateOptions) (RegisteredPrefixesClientValidateResponse, error)`
- New function `NewRpUnbilledPrefixesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RpUnbilledPrefixesClient, error)`
- New function `*RpUnbilledPrefixesClient.NewListPager(resourceGroupName string, peeringName string, options *RpUnbilledPrefixesClientListOptions) *runtime.Pager[RpUnbilledPrefixesClientListResponse]`
- New struct `ConnectivityProbe`
- New struct `RpUnbilledPrefix`
- New struct `RpUnbilledPrefixListResult`
- New struct `SystemData`
- New field `SystemData` in struct `CdnPeeringPrefix`
- New field `SystemData` in struct `ConnectionMonitorTest`
- New field `DirectPeeringType` in struct `LegacyPeeringsClientListOptions`
- New field `SystemData` in struct `Location`
- New field `ContentType` in struct `ManagementClientCheckServiceProviderAvailabilityResponse`
- New field `SystemData` in struct `PeerAsn`
- New field `SystemData` in struct `Peering`
- New field `ConnectivityProbes` in struct `Properties`
- New field `SystemData` in struct `RegisteredAsn`
- New field `SystemData` in struct `RegisteredPrefix`
- New field `SystemData` in struct `Service`
- New field `SystemData` in struct `ServiceCountry`
- New field `SystemData` in struct `ServiceLocation`
- New field `SystemData` in struct `ServicePrefix`
- New field `SystemData` in struct `ServiceProvider`


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