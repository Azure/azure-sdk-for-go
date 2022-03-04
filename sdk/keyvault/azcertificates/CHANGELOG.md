# Release History

## 0.2.0 (Unreleased)

### Features Added

### Breaking Changes
* Changed pager APIs for `ListCertificatesPager`, `ListDeletedCertificatesPager`, `ListPropertiesOfIssuersPager`, and `ListCertificateVersionsPager`
    * Use the `More()` method to determine if there are more pages to fetch
    * Use the `NextPage(context.Context)` to fetch the next page of results

### Bugs Fixed

### Other Changes

## 0.1.0 (2022-02-08)
* This is the initial release of the `azcertificates` library
