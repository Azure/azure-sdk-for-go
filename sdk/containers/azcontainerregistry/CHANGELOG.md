# Release History

## 0.2.1 (Unreleased)

### Features Added
* Add `ConfigMediaType` and `MediaType` properties to `ManifestAttributes`

### Other Changes
* Refine some logics and comments

## 0.2.0 (2023-06-06)

### Features Added
* Add `DigestValidationReader` to help to do digest validation when read manifest or blob

### Breaking Changes
* Remove `MarshalJSON` for some of the types that are not used in the request.

### Bugs Fixed
* Add state restore for hash calculator when upload fails
* Do not re-calculate digest when retry

### Other Changes
* Change default audience to https://containerregistry.azure.net
* Refine examples of image upload and download

## 0.1.1 (2023-03-07)

### Bugs Fixed
* Fix possible failure when request retry

### Other Changes
* Rewrite auth policy to promote efficiency of auth process

## 0.1.0 (2023-02-07)

* This is the initial release of the `azcontainerregistry` library
