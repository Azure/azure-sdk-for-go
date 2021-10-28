# Release History

## 0.2.0 (2021-10-28)
### Breaking Changes

- Function `NewDisasterRecoveryConfigsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRulesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSubscriptionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNamespacesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewMigrationConfigsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTopicsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewQueuesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 10 breaking change(s), 0 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/servicebus/armservicebus") to avoid confusion.