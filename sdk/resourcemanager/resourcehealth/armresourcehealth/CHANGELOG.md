# Release History

## 2.0.0 (2022-09-15)
### Breaking Changes

- Struct `ErrorResponseError` has been removed
- Field `UnavailableOccurredTime` of struct `AvailabilityStatusPropertiesRecentlyResolved` has been removed
- Field `UnavailabilitySummary` of struct `AvailabilityStatusPropertiesRecentlyResolved` has been removed
- Field `OccurredTime` of struct `AvailabilityStatusProperties` has been removed
- Field `Error` of struct `ErrorResponse` has been removed

### Features Added

- New field `OccuredTime` in struct `AvailabilityStatusProperties`
- New field `UnavailableOccuredTime` in struct `AvailabilityStatusPropertiesRecentlyResolved`
- New field `UnavailableSummary` in struct `AvailabilityStatusPropertiesRecentlyResolved`
- New field `Message` in struct `ErrorResponse`
- New field `Code` in struct `ErrorResponse`
- New field `Details` in struct `ErrorResponse`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).