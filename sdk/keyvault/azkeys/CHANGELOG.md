# Release History

## 0.2.0 (Unreleased)

### Features Added
* Adds the `CreateOKPKey` method for creating OKP public keys
* Adds the `Ed25519` `JSONWebKeyCurveName`
*Adds the `OKP` and `OKPHSM` `JSONWebKeyType`s

### Breaking Changes
* `NewClient` returns an instance of a `Client`, instead of a `*Client`

### Bugs Fixed
* Fixes a bug in `crypto.NewClient` where the key version was required in the path, it is no longer required but is recommended.

### Other Changes

## 0.1.0 (2021-11-09)
* This is the initial release of the `azkeys` library
