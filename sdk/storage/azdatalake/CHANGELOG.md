# Release History

## 1.0.1 (Unreleased)

### Features Added
* Encryption Scope For SAS
* CPK for Datalake
* Create SubDirectory Client
* Service Version upgrade to 2021-06-08

### Breaking Changes

### Bugs Fixed

* Block `SharedKeyCredential` authentication mode for non TLS protected endpoints. Fixes [#21841](https://github.com/Azure/azure-sdk-for-go/issues/21841).

### Other Changes
* Updated version of azblob to 1.2.1

* Updated azcore version to `1.9.1` and azidentity version to `1.4.0`.

## 1.0.0 (2023-10-18)

### Bugs Fixed
* Fixed an issue where customers could not capture the raw HTTP response of directory and file GetProperties operations.
* Fixed an issue where file/directory renames with source/destination SAS tokens fail with authorization failures.

## 0.1.0-beta.1 (2023-08-16)

### Features Added

* This is the initial preview release of the `azdatalake` library