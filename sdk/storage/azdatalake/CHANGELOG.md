# Release History

## 1.1.2 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 1.1.2-beta.1 (2024-04-10)

### Features Added
* Append API Bundled with Flush functionality
* HNS Encryption Scope support
* Append API with acquire lease, release lease and renewal of lease support.
* Flush API bundled with release lease option.
* HNS Encryption Context support
* Pagination Support for recursive directory deletion
* Bundle ability to set permission, owner, group, acl, lease, expiry time and umask along with FileSystem.CreateFile and FileSystem.CreateDirectory APIs.
* Added support for AAD Audience when OAuth is used.
* Updated service version to `2023-11-03`
* Integrate `InsecureAllowCredentialWithHTTP` client options.

### Bugs Fixed
* Fixed an issue where GetSASURL() was providing HTTPS SAS, instead of the default http+https SAS. Fixes [#22448](https://github.com/Azure/azure-sdk-for-go/issues/22448)

### Other Changes
* Updated azcore version to `1.11.1`

## 1.1.1 (2024-02-29)

### Bugs Fixed
* Exposing x-ms-resource-type response header in GetProperties API for file and directory.

* Re-enabled `SharedKeyCredential` authentication mode for non TLS protected endpoints.

### Other Changes
* Updated version of azblob to `1.3.1`

## 1.1.0 (2024-02-14)

### Bugs Fixed
* Escape paths for NewDirectoryClient and NewFileClient in a file system. Fixes [#22281](https://github.com/Azure/azure-sdk-for-go/issues/22281).

### Other Changes
* Updated version of azblob to `1.3.0`
* Updated azcore version to `1.9.2` and azidentity version to `1.5.1`.

## 1.1.0-beta.1 (2024-01-10)

### Features Added
* Encryption Scope For SAS
* CPK for Datalake
* Create SubDirectory Client
* Service Version upgrade to 2021-06-08

### Bugs Fixed

* Block `SharedKeyCredential` authentication mode for non TLS protected endpoints. Fixes [#21841](https://github.com/Azure/azure-sdk-for-go/issues/21841).

### Other Changes
* Updated version of azblob to `1.3.0-beta.1`
* Updated azcore version to `1.9.1` and azidentity version to `1.4.0`.

## 1.0.0 (2023-10-18)

### Bugs Fixed
* Fixed an issue where customers could not capture the raw HTTP response of directory and file GetProperties operations.
* Fixed an issue where file/directory renames with source/destination SAS tokens fail with authorization failures.

## 0.1.0-beta.1 (2023-08-16)

### Features Added

* This is the initial preview release of the `azdatalake` library
