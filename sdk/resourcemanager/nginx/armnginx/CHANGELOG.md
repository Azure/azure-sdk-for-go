# Release History

## 2.0.0 (2022-10-13)
### Breaking Changes

- Function `*CertificatesClient.BeginCreate` has been removed
- Function `*DeploymentsClient.BeginCreate` has been removed
- Struct `CertificatesClientBeginCreateOptions` has been removed
- Struct `CertificatesClientCreateResponse` has been removed
- Struct `DeploymentsClientBeginCreateOptions` has been removed
- Struct `DeploymentsClientCreateResponse` has been removed

### Features Added

- New function `*CertificatesClient.BeginCreateOrUpdate(context.Context, string, string, string, *CertificatesClientBeginCreateOrUpdateOptions) (*runtime.Poller[CertificatesClientCreateOrUpdateResponse], error)`
- New function `*DeploymentsClient.BeginCreateOrUpdate(context.Context, string, string, *DeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DeploymentsClientCreateOrUpdateResponse], error)`
- New struct `CertificatesClientBeginCreateOrUpdateOptions`
- New struct `CertificatesClientCreateOrUpdateResponse`
- New struct `DeploymentsClientBeginCreateOrUpdateOptions`
- New struct `DeploymentsClientCreateOrUpdateResponse`
- New field `ProtectedFiles` in struct `ConfigurationProperties`


## 1.0.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 1.0.0 (2022-08-19)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/nginx/armnginx` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).