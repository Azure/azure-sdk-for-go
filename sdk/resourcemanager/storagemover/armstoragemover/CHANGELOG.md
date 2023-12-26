# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-10-27)
### Breaking Changes

- Type of `EndpointBaseUpdateParameters.Properties` has been changed from `*EndpointBaseUpdateProperties` to `EndpointBaseUpdatePropertiesClassification`

### Features Added

- New value `EndpointTypeAzureStorageSmbFileShare`, `EndpointTypeSmbMount` added to enum type `EndpointType`
- New enum type `CredentialType` with values `CredentialTypeAzureKeyVaultSmb`
- New function `*AzureKeyVaultSmbCredentials.GetCredentials() *Credentials`
- New function `*AzureStorageBlobContainerEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*AzureStorageSmbFileShareEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureStorageSmbFileShareEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*Credentials.GetCredentials() *Credentials`
- New function `*EndpointBaseUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*SmbMountEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*SmbMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*NfsMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureKeyVaultSmbCredentials`
- New struct `AzureStorageBlobContainerEndpointUpdateProperties`
- New struct `AzureStorageSmbFileShareEndpointProperties`
- New struct `AzureStorageSmbFileShareEndpointUpdateProperties`
- New struct `NfsMountEndpointUpdateProperties`
- New struct `SmbMountEndpointProperties`
- New struct `SmbMountEndpointUpdateProperties`


## 2.0.0-beta.1 (2023-07-28)
### Breaking Changes

- Type of `EndpointBaseUpdateParameters.Properties` has been changed from `*EndpointBaseUpdateProperties` to `EndpointBaseUpdatePropertiesClassification`
### Features Added
- New value `EndpointTypeAzureStorageSmbFileShare`, `EndpointTypeSmbMount` added to enum type `EndpointType`
- New enum type `CredentialType` with values `CredentialTypeAzureKeyVaultSmb`
- New function `*AzureKeyVaultSmbCredentials.GetCredentials() *Credentials`
- New function `*AzureStorageBlobContainerEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*AzureStorageSmbFileShareEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureStorageSmbFileShareEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*Credentials.GetCredentials() *Credentials`
- New function `*EndpointBaseUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*SmbMountEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*SmbMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*NfsMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureKeyVaultSmbCredentials`
- New struct `AzureStorageBlobContainerEndpointUpdateProperties`
- New struct `AzureStorageSmbFileShareEndpointProperties`
- New struct `AzureStorageSmbFileShareEndpointUpdateProperties`
- New struct `NfsMountEndpointUpdateProperties`
- New struct `SmbMountEndpointProperties`
- New struct `SmbMountEndpointUpdateProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2023-03-07)
### Other Changes

- Release stable version.

## 0.1.0 (2023-02-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).