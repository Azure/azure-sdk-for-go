# Release History

## 0.2.0 (2021-10-29)
### Breaking Changes

- Function `NewVirtualNetworkLinksClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateZonesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRecordSetsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 3 breaking change(s), 0 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/privatedns/armprivatedns") to avoid confusion.