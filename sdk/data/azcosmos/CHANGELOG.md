# Release History

## 0.3.1 (Unreleased)

### Features Added

* Added Transactional Batch support

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.3.0 (2022-05-10)

### Features Added
* Added single partition query support.
* Added Azure AD authentication support through `azcosmos.NewClient`

### Breaking Changes
* This module now requires Go 1.18

## 0.2.0 (2022-01-13)

### Features Added
* Failed API calls will now return an `*azcore.ResponseError` type.

### Breaking Changes
* Updated to latest `azcore`. Public surface area is unchanged.  However, the `azcore.HTTPResponse` interface has been removed.

## 0.1.0 (2021-11-09)
* This is the initial preview release of the `azcosmos` library
