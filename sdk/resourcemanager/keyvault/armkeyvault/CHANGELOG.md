# Release History

## 0.2.1 (2021-11-26)

### Other Changes

- Now use `github.com/Azure/azure-sdk-for-go/sdk/azidentity@v0.12.0` explicitly.

## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-09-29)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/keyvault/armkeyvault") to avoid confusion. 
