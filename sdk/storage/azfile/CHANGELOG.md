# Release History

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed
* Fixed issue where some requests fail with mismatch in string to sign.

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
