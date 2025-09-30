# Release History

## 1.4.1-beta.1 (Unreleased)

### Features Added
* Added support for IP addresses and URIs in `SubjectAlternativeNames` through new `IPAddresses` and `Uris` fields

### Breaking Changes

### Bugs Fixed

### Other Changes
* Upgraded to API service version `2025-07-01`

## 1.4.0 (2025-06-12)

### Features Added
* Add fakes support (https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/samples/fakes)

### Other Changes
* Upgraded to API service version `7.6`

## 1.4.0-beta.1 (2025-04-10)

### Features Added
* Added `PreserveCertOrder`

### Other Changes
* Upgraded to API service version `7.6-preview.2`

## 1.3.1 (2025-02-13)

### Other Changes
* Upgraded dependencies

## 1.3.0 (2024-11-06)

### Features Added
* Added API Version support. Users can now change the default API Version by setting ClientOptions.APIVersion

## 1.2.0 (2024-10-21)

### Features Added
* Added CAE support
* Client requests tokens from the Vault's tenant, overriding any credential default
  (thanks @francescomari)

## 1.1.0 (2024-02-13)

### Other Changes
* Upgraded to API service version `7.5`
* Upgraded dependencies

## 1.1.0-beta.1 (2023-11-08)

### Other Changes
* Upgraded service version to `7.5-preview.1`
* Updated to latest version of `azcore`.
* Enabled spans for distributed tracing.

## 1.0.0 (2023-09-12)

### Features Added
* First stable release of the azcertificates module

### Other Changes
* Upgraded dependencies

## 0.11.0 (2023-07-17)

### Breaking Changes
* Rename `ListCertificates` to `ListCertificateProperties`
* `ListCertificateIssuers` to `ListIssuerProperties`
* `ListCertificateVersions` to `ListCertificatePropertiesVersions`
* `ListDeletedCertificates` to `ListDeletedCertificateProperties`
* `CertificateListResult` to `CertificatePropertiesListResult`
* `DeletedCertificateListResult` to `DeletedCertificatePropertiesListResult`
* `SetCertificateContacts` to `SetContacts`
* `GetCertificateContacts` to `GetContacts`
* `DeleteCertificateContacts` to `DeleteContacts`
* `SetCertificateIssuer` to `SetIssuer`
* `UpdateCertificateIssuer` to `UpdateIssuer`
* `GetCertificateIssuer` to `GetIssuer`
* `DeleteCertificateIssuer` to `DeleteIssuer`
* `CertificateIssuerListResult` to `IssuerPropertiesListResult`
* `UpdateCertificateIssuerParameters` to `UpdateIssuerParameters`
* `SetCertificateIssuerParameters` to `SetIssuerParameters`
* `CertificateBundle` to `Certificate`
* `CertificateItem` to `CertificateProperties`
* `DeletedCertificateBundle` to `DeletedCertificate`
* `DeletedCertificateItem` to `DeletedCertificateProperties`
* `IssuerBundle` to `Issuer`
* `CertificateIssuerItem` to `IssuerProperties`
* `RestoreCertificateParameters.CertificateBundleBackup` to `RestoreCertificateParameters.CertificateBackup`
* `JSONWebKeyCurveName` to `CurveName`
* `JSONWebKeyType` to `KeyType`
* `Trigger` to `LifetimeActionTrigger`
* `Action` to `LifetimeActionType`
* `AdministratorDetails` to ``AdministratorContact`
* `OrganizationDetails.AdminDetails` to `OrganizationDetails.AdminContacts`
* `EmailAddress` to `Email`
* `UPNs` to `UserPrincipalNames`
* `EKUs` to `EnhancedKeyUsage`
* remove `MaxResults` parameter
* remove `DeletionRecoveryLevel` type

### Other Changes
* Updated dependencies

## 0.10.0 (2023-04-13)

### Breaking Changes
* Moved module from `sdk/keyvault/azadmin` to `sdk/security/keyvault/azadmin`

## 0.9.0 (2023-04-13)

### Features Added
* Upgraded to api version 7.4

### Breaking Changes
* Renamed `ActionType` to `CertificatePolicyAction`
* Changed `Error` struct to `ErrorInfo`

## 0.8.0 (2022-11-08)

### Breaking Changes
* `NewClient` returns an `error`

## 0.7.1 (2022-09-20)

### Features Added
* Added `ClientOptions.DisableChallengeResourceVerification`.
  See https://aka.ms/azsdk/blog/vault-uri for more information.

## 0.7.0 (2022-09-12)

### Breaking Changes
* Verify the challenge resource matches the vault domain.

## 0.6.0 (2022-08-09)

### Breaking Changes
* Changed type of `NewClient` options parameter to `azcertificates.ClientOptions`, which embeds
  the former type, `azcore.ClientOptions`

## 0.5.0 (2022-07-07)

### Breaking Changes
* The `Client` API now corresponds more directly to the Key Vault REST API.
  Most method signatures and types have changed. See the
  [module documentation](https://aka.ms/azsdk/go/keyvault-certificates/docs)
  for updated code examples and more details.

### Other Changes
* Upgrade to latest `azcore`

## 0.4.1 (2022-05-12)

### Other Changes
* Updated to latest `azcore` and `internal` modules.

## 0.4.0 (2022-04-21)

### Breaking Changes
* Renamed the following methods
  * `Client.ListPropertiesOfCertificates` to `Client.NewListPropertiesOfCertificatesPager`
  * `Client.ListPropertiesOfCertificateVersions` to `Client.NewListPropertiesOfCertificateVersionsPager`
  * `Client.ListPropertiesOfIssuers` to `Client.NewListPropertiesOfIssuersPager`
  * `Client.ListDeletedCertificates` to `Client.NewListDeletedCertificatesPager`

### Other Changes
* Regenerated code with latest code generator (no visible changes).

## 0.3.0 (2022-04-06)

### Features Added
* Added PossibleValues functions for `CertificateKeyUsage`, `CertificateKeyType`, `CertificateKeyCurveName`, and `CertificatePolicyAction` constants.
* Added the `ResumeToken` method on pollers for resuming operations later
* Added the `ResumeToken` field to the options structs of `Begin` methods for resuming operations
* Added the `Name *string` field to `Certificate`, `CertificateItem`, `DeletedCertificate`, `DeletedCertificateItem`

### Breaking Changes
* Requires Go 1.18
* Fixed a misspelling of `CerificateKeyUsage`, changed to `CertificateKeyUsage`
* Removed all `ToPtr` methods from constants
* Renamed `CertificateOperation` to `Operation`
* Renamed `Operation.Csr` to `Operation.CSR`
* Renamed `KeyVaultCertificateWithPolicy` to `CertificateWithPolicy`
* Abbreviated `EmailAddress` to `Email`
* Changed `Upns` to `UserPrincipalNames`
* Removed the `Trigger` struct and elevated it to the `LifetimeAction`
* Renamed `DeletedDate` to `DeletedOn` and `Expires` to `ExpiresOn`

## 0.2.0 (2022-03-08)

### Breaking Changes
* Changed pager APIs for `ListCertificatesPager`, `ListDeletedCertificatesPager`, `ListPropertiesOfIssuersPager`, and `ListCertificateVersionsPager`
    * Use the `More()` method to determine if there are more pages to fetch
    * Use the `NextPage(context.Context)` to fetch the next page of results
* Removed all `RawResponse *http.Response` fields from response structs.

## 0.1.0 (2022-02-08)
* This is the initial release of the `azcertificates` library
