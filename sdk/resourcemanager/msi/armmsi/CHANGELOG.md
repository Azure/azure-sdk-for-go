# Release History

## 0.7.0 (2022-06-27)
### Features Added

- New function `*FederatedIdentityCredentialsClient.Delete(context.Context, string, string, string, *FederatedIdentityCredentialsClientDeleteOptions) (FederatedIdentityCredentialsClientDeleteResponse, error)`
- New function `*FederatedIdentityCredentialsClient.CreateOrUpdate(context.Context, string, string, string, FederatedIdentityCredential, *FederatedIdentityCredentialsClientCreateOrUpdateOptions) (FederatedIdentityCredentialsClientCreateOrUpdateResponse, error)`
- New function `*FederatedIdentityCredentialsClient.NewListPager(string, string, *FederatedIdentityCredentialsClientListOptions) *runtime.Pager[FederatedIdentityCredentialsClientListResponse]`
- New function `NewFederatedIdentityCredentialsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FederatedIdentityCredentialsClient, error)`
- New function `*FederatedIdentityCredentialsClient.Get(context.Context, string, string, string, *FederatedIdentityCredentialsClientGetOptions) (FederatedIdentityCredentialsClientGetResponse, error)`
- New struct `FederatedIdentityCredential`
- New struct `FederatedIdentityCredentialProperties`
- New struct `FederatedIdentityCredentialsClient`
- New struct `FederatedIdentityCredentialsClientCreateOrUpdateOptions`
- New struct `FederatedIdentityCredentialsClientCreateOrUpdateResponse`
- New struct `FederatedIdentityCredentialsClientDeleteOptions`
- New struct `FederatedIdentityCredentialsClientDeleteResponse`
- New struct `FederatedIdentityCredentialsClientGetOptions`
- New struct `FederatedIdentityCredentialsClientGetResponse`
- New struct `FederatedIdentityCredentialsClientListOptions`
- New struct `FederatedIdentityCredentialsClientListResponse`
- New struct `FederatedIdentityCredentialsListResult`


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).