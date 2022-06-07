# Release History

## 0.5.0 (2022-05-16)

### Breaking Changes
* Removed types `CreateCertificatePoller`, `DeleteCertificatePoller`, and `RecoverDeletedCertificatePoller`.
* Methods `BeginCreateCertificate`, `BeginDeleteCertificate`, and `BeginRecoverDeletedCertificate` now return a `*runtime.Poller[T]` with their respective response types.
* Options types with a `ResumeToken` field now take the token by value.
* The poller for `BeginCreateCertificate` now returns the created certificate from its `PollUntilDone` method.
* Changed type of certificate `Tags` to `map[string]*string`
* Deleted `UpdateCertificatePropertiesOptions` fields
* Renamed types
  * `ListIssuersPropertiesOfIssuersResponse` to `ListPropertiesOfIssuersResponse`
  * `ListCertificatesOptions` to `ListPropertiesOfCertificatesOptions`
  * `ListCertificateVersionsOptions` to `ListPropertiesOfCertificateVersionsOptions`
* Renamed `ListDeletedCertificatesResponse.Certificates` to `.DeletedCertificates`
* `UpdateCertificateProperties()` has a `Properties` parameter instead of a `string` parameter
* Removed JSON tags from models

### Bugs Fixed
* LROs now correctly exit the polling loop in `PollUntilDone()` when the operations reach a terminal state.

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
