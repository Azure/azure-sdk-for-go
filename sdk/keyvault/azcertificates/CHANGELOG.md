# Release History

## 0.3.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
