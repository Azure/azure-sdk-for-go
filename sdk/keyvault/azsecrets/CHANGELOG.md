# Release History

## 0.5.0 (Unreleased)

### Features Added

### Breaking Changes
* Fixes a bug where `UpdateSecretProperties` will delete properties that are not explicitly set each time. This is only a breaking change at runtime, where the request body will change.

### Bugs Fixed

### Other Changes

## 0.4.0 (2022-01-11)

### Other Changes
* Bumps `azcore` dependency from `v0.20.0` to `v0.21.0`

## 0.3.0 (2021-11-09)

### Features Added
* Clients can now connect to Key Vaults in any cloud

## 0.2.0 (2021-11-02)

### Other Changes
* Bumps `azcore` dependency to `v0.20.0` and `azidentity` to `v0.12.0`

## 0.1.1 (2021-10-06)
* Adds the MIT License for redistribution

## 0.1.0 (2021-10-05)
* This is the initial release of the `azsecrets` library
