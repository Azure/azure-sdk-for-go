# Release History


## 1.1.2 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes


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
