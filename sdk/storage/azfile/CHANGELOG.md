# Release History

## 0.2.0 (Unreleased)

### Features Added

* Updated service version to STG 87(2022-11-02).
* Added OAuth support.
* Added [Rename Directory API](https://learn.microsoft.com/rest/api/storageservices/rename-directory).
* Added [Rename File API](https://learn.microsoft.com/rest/api/storageservices/rename-file).
* Updated the SAS Version to `2022-11-02` and added `Encryption Scope` to Account SAS.

### Breaking Changes

### Bugs Fixed

* Fixed service SAS creation where expiry time or permissions can be omitted when stored access policy is used.

### Other Changes

* Updated version of azcore to 1.7.0 and azidentity to 1.3.0.

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

* Added `dragonfly` and `aix` to build constraints in `mmf_unix.go`.

## 1.0.0 (2023-07-12)

### Features Added

* Added `ParseNTFSFileAttributes` method for parsing the file attributes to `file.NTFSFileAttributes` type.

### Bugs Fixed

* Fixed the issue where trailing slash is encoded when passed in directory or subdirectory name while creating the directory client.

## 0.1.0 (2023-05-09)

### Features Added

* This is the initial preview release of the `azfile` library
