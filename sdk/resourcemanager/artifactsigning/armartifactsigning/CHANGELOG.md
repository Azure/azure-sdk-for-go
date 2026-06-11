# Release History

## 0.2.0 (2026-06-09)
### Breaking Changes

- Function `*CertificateProfilesClient.RevokeCertificate` has been removed

### Features Added

- New function `*CertificateProfilesClient.RevokeCertificates(ctx context.Context, resourceGroupName string, accountName string, profileName string, body RevokeCertificateList, options *CertificateProfilesClientRevokeCertificatesOptions) (CertificateProfilesClientRevokeCertificatesResponse, error)`
- New struct `RevokeCertificateList`
- New field `ProgramType` in struct `CertificateProfileProperties`


## 0.1.0 (2026-02-11)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/artifactsigning/armartifactsigning` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).