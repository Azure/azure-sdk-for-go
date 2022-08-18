# Release History

## 0.6.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
