# Release History

## 1.5.1-beta.1 (Unreleased)

### Features Added
* Upgraded service version to `2025-07-05`.
* NFS Over Rest API Request/Response changes in the following APIs - 
  CreateFile
  SetFileProperties
  GetFileProperties
  StartCopyFromUrl
  DownloadFile
  DeleteFile
  CreateDirectory
  GetDirectoryProperties
  SetDirectoryProperties
* NFS Over Rest, new API File CreateHardLink() Added 
  

### Breaking Changes
* The following APIs no longer send the x-ms-file-permission-key, x-ms-file-attributes, x-ms-file-creation-time, and x-ms-file-last-write-time request headers by default. These headers have been optional in the REST API since x-ms-version 2021-06-08:
  FileClient.Create()
  FileClient.SetHttpHeaders()
  DirectoryClient.Create()
  DirectoryClient.SetHttpHeaders()

### Bugs Fixed

### Other Changes

## 1.5.0 (2025-01-23)

### Features Added
* Upgraded service version to `2025-01-05`.

## 1.5.0-beta.1 (2025-01-13)

### Features Added
* Add Premium access tier to Files Shares
* Provisioned v2 billing model for Azure files

## 1.4.0 (2024-11-13)

### Features Added
* Fix compareHeaders custom sorting algorithm for String To Sign.

## 1.4.0-beta.1 (2024-10-22)

### Features Added
* Add Paid Burst IOPS/Bandwidth support for Premium Files.
* REST API for binary ACE in Azure Files.

### Other Changes
* Updated `azidentity` version to `1.8.0`
* Updated `azcore` version to `1.16.0`

## 1.3.1 (2024-09-18)

### Features Added
* Upgraded service version to `2024-08-04`.

## 1.3.1-beta.1 (2024-08-27)

### Features Added
* Snapshot management support via REST for NFS shares.
* Upgraded service version to `2024-08-04`.

### Other Changes
* Updated `azcore` version to `1.14.0`

## 1.3.0 (2024-07-18)

### Other Changes
* GetProperties() was called twice in DownloadFile method. Enhanced to call it only once, reducing latency.
* Updated `azcore` version to `1.13.0`

## 1.3.0-beta.1 (2024-06-14)

### Features Added
* Add Rename Support to List Ranges API
* Updated service version to `2024-05-04`

### Other Changes
* Updated `azidentity` version to `1.6.0`
* Updated `azcore` version to `1.12.0`

## 1.2.2 (2024-04-09)

### Bugs Fixed
* Fixed an issue where GetSASURL() was providing HTTPS SAS, instead of the default http+https SAS. Fixes [#22448](https://github.com/Azure/azure-sdk-for-go/issues/22448)

### Other Changes
* Integrate `InsecureAllowCredentialWithHTTP` client options.
* Update dependencies.

## 1.2.1 (2024-02-29)

### Bugs Fixed

* Re-enabled `SharedKeyCredential` authentication mode for non TLS protected endpoints.

### Other Changes

* Updated `azidentity` version to `1.5.1`.

## 1.2.0 (2024-02-12)

### Other Changes

* Updated `azcore` version to `1.9.2`.

## 1.2.0-beta.1 (2024-01-09)

### Features Added

* Updated service version to `2023-11-03`.
* Added support for Audience when OAuth is used.

### Bugs Fixed

* Block `SharedKeyCredential` authentication mode for non TLS protected endpoints. Fixes [#21841](https://github.com/Azure/azure-sdk-for-go/issues/21841).
* Fixed a bug where `UploadRangeFromURL` using OAuth was returning error.

### Other Changes

* Updated azcore version to `1.9.1`.

## 1.1.1 (2023-11-15)

### Bugs Fixed

* Fixed a bug where Optional fields which were mandatory earlier create a failure when passed an older service version
* Made SourceContentCRC64 header as optional. Changed the type from uint64 to a generic interface implementation. 
  These changes impact: `file.UploadRangeFromURL()`

### Other Changes

* Updated azcore version to `1.9.0` and azidentity version to `1.4.0`.

## 1.1.0 (2023-10-11)

### Features Added

* Updated service version to `2022-11-02`.

### Bugs Fixed

* Fixed a bug where the `x-ms-file-attributes` header could be set to contain invalid trailing or leading | characters.

## 1.1.0-beta.1 (2023-09-12)

### Features Added

* Updated service version to `2022-11-02`.
* Added OAuth support.
* Added [Rename Directory API](https://learn.microsoft.com/rest/api/storageservices/rename-directory).
* Added [Rename File API](https://learn.microsoft.com/rest/api/storageservices/rename-file).
* Added `x-ms-file-change-time` request header in
  * Create File/Directory
  * Set File/Directory Properties
  * Copy File
* Added `x-ms-file-last-write-time` request header in Put Range and Put Range from URL.
* Updated the SAS Version to `2022-11-02` and added `Encryption Scope` to Account SAS.
* Trailing dot support for files and directories.

### Bugs Fixed

* Fixed service SAS creation where expiry time or permissions can be omitted when stored access policy is used.
* Fixed issue where some requests fail with mismatch in string to sign.

### Other Changes

* Updated version of azcore to 1.7.2 and azidentity to 1.3.1.
* Added `dragonfly` and `aix` to build constraints in `mmf_unix.go`.

## 1.0.0 (2023-07-12)

### Features Added

* Added `ParseNTFSFileAttributes` method for parsing the file attributes to `file.NTFSFileAttributes` type.

### Bugs Fixed

* Fixed the issue where trailing slash is encoded when passed in directory or subdirectory name while creating the directory client.

## 0.1.0 (2023-05-09)

### Features Added

* This is the initial preview release of the `azfile` library
