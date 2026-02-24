# Release History

## 0.4.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.4.0 (2026-01-13)

### Other Changes

* Linux caches are stored in $XDG_CACHE_HOME by default, falling back to $HOME/.cache

## 0.3.2 (2025-01-15)

### Other Changes
* Upgraded dependencies

## 0.3.1 (2024-11-05)

### Bugs Fixed
* `New` no longer returns an error when called simultaneously in two processes

## 0.3.0 (2024-07-18)

### Features Added
* Added `New`, a constructor for persistent caches. See `azidentity` docs,
  in particular the `PersistentUserAuthentication` example, for usage details.

### Breaking Changes
* Removed optional fallback to plaintext storage. `azidentity/cache` now
  always returns an error when it can't encrypt a persistent cache.

## 0.2.2 (2024-05-07)

### Bugs Fixed
* On Linux, prevent "permission denied" errors by linking the session keyring
  to the user keyring so the process possesses any keys it adds

### Other Changes
* Upgraded dependencies

## 0.2.1 (2023-11-07)

### Other Changes
* Upgraded dependencies and documentation

## 0.2.0 (2023-10-10)

### Bugs Fixed
* Correct dependency versions

## 0.1.0 (2023-10-10)

### Features Added
* Initial release
