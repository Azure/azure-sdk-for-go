# Release History

## 0.3.0 (Unreleased)

### Features Added
* Added PossibleValues functions for `CertificateKeyUsage`, `CertificateKeyType`, `CertificateKeyCurveName`, and `CertificatePolicyAction` constants.
* Added the `ResumeToken` method on pollers for resuming operations later
* Added the `ResumeToken` field to the options structs of `Begin` methods for resuming operations

### Breaking Changes
* Fixed a misspelling of `CerificateKeyUsage`, changed to `CertificateKeyUsage`
* Removed all `ToPtr` methods from constants

### Bugs Fixed

### Other Changes

## 0.2.0 (2022-03-08)

### Breaking Changes
* Changed pager APIs for `ListCertificatesPager`, `ListDeletedCertificatesPager`, `ListPropertiesOfIssuersPager`, and `ListCertificateVersionsPager`
    * Use the `More()` method to determine if there are more pages to fetch
    * Use the `NextPage(context.Context)` to fetch the next page of results
* Removed all `RawResponse *http.Response` fields from response structs.

## 0.1.0 (2022-02-08)
* This is the initial release of the `azcertificates` library
