# Release History

## 0.4.0 (2022-04-22)

### Breaking Changes
* Updated `ExpiringResource` and its dependent types to use generics.

### Other Changes
* Remove reference to `TokenRequestOptions.TenantID` as it's been removed and wasn't working anyways.

## 0.3.0 (2022-04-04)

### Features Added
* Adds the `ParseKeyvaultID` function to parse an ID into the Key Vault URL, item name, and item version

### Breaking Changes
* Updates to azcore v0.23.0

## 0.2.1 (2022-01-31)

### Bugs Fixed
* Avoid retries on terminal failures (#16932)

## 0.2.0 (2022-01-12)

### Bugs Fixed
* Fixes a bug with Managed HSMs that prevented correctly authorizing requests.

## 0.1.0 (2021-11-09)
* This is the initial release of the `internal` library for KeyVault
