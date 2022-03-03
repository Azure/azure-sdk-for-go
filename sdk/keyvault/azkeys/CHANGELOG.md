# Release History

## 0.4.0 (2022-03-08)

### Features Added
* Adds the `ReleasePolicy` parameter to the `UpdateKeyPropertiesOptions` struct.
* Adds the `Immutable` boolean to the `KeyReleasePolicy` model.
* Added a `ToPtr` method on `KeyType` constant

### Breaking Changes
* Changed the `Data` to `EncodedPolicy` on the `KeyReleasePolicy` struct.
* Changed the `Updated`, `Created`, and `Expires` properties to `UpdatedOn`, `CreatedOn`, and `ExpiresOn`.
* Renamed `JSONWebKeyOperation` to `Operation`.
* Renamed `JSONWebKeyCurveName` to `CurveName`
* Prefixed all KeyType constants with `KeyType`
* Changed `KeyBundle` to `KeyVaultKey` and `DeletedKeyBundle` to `DeletedKey`
* Renamed `KeyAttributes` to `KeyProperties`
* Renamed `ListKeyVersions` to `ListPropertiesOfKeyVersions`
* Removed `Attributes` struct
* Changed `CreateOCTKey`/`Response`/`Options` to `CreateOctKey`/`Response`/`Options`

## 0.3.0 (2022-02-08)

### Breaking Changes
* Changed the `Tags` properties from `map[string]*string` to `map[string]string`

### Bugs Fixed
* Fixed a bug in `UpdateKeyProperties` where the `KeyOps` would be deleted if the `UpdateKeyProperties.KeyOps` value was left empty.

## 0.2.0 (2022-01-12)

### Bugs Fixed
* Fixes a bug in `crypto.NewClient` where the key version was required in the path, it is no longer required but is recommended.

### Other Changes
* Updates `azcore` dependency from `v0.20.0` to `v0.21.0`

## 0.1.0 (2021-11-09)
* This is the initial release of the `azkeys` library
