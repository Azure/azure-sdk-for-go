# Release History

## 1.1.0 (2022-10-31)
### Features Added

- New const `ConnectionStateTypeChangeInProgress`
- New const `ConnectionStateTypeChangeRequested`
- New function `*RpUnbilledPrefixesClient.NewListPager(string, string, *RpUnbilledPrefixesClientListOptions) *runtime.Pager[RpUnbilledPrefixesClientListResponse]`
- New function `*RegisteredPrefixesClient.Validate(context.Context, string, string, string, *RegisteredPrefixesClientValidateOptions) (RegisteredPrefixesClientValidateResponse, error)`
- New function `NewRpUnbilledPrefixesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RpUnbilledPrefixesClient, error)`
- New struct `RegisteredPrefixesClientValidateOptions`
- New struct `RegisteredPrefixesClientValidateResponse`
- New struct `RpUnbilledPrefix`
- New struct `RpUnbilledPrefixListResult`
- New struct `RpUnbilledPrefixesClient`
- New struct `RpUnbilledPrefixesClientListOptions`
- New struct `RpUnbilledPrefixesClientListResponse`
- New field `DirectPeeringType` in struct `LegacyPeeringsClientListOptions`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/peering/armpeering` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).