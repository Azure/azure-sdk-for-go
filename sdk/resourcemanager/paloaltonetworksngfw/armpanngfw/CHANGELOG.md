# Release History

## 2.0.0 (2023-11-24)
### Breaking Changes

- Function `*LocalRulestacksClient.ListAppIDs` has been removed
- Function `*LocalRulestacksClient.ListCountries` has been removed
- Function `*LocalRulestacksClient.ListPredefinedURLCategories` has been removed
- Function `timeRFC3339.MarshalText` has been removed
- Function `*timeRFC3339.Parse` has been removed
- Function `*timeRFC3339.UnmarshalText` has been removed

### Features Added

- New function `dateTimeRFC3339.MarshalText() ([]byte, error)`
- New function `*dateTimeRFC3339.Parse(string) error`
- New function `*dateTimeRFC3339.UnmarshalText([]byte) error`
- New field `TrustedRanges` in struct `NetworkProfile`


## 1.0.0 (2023-07-14)
### Other Changes

- Release stable version.

## 0.1.0 (2023-04-28)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).