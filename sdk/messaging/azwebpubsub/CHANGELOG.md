# Release History

## 0.1.2 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed
* Fixed `GenerateClientAccessURL` panicking when `options` is `nil` and the client was created with `NewClient`.
* Fixed `GenerateClientAccessURL` failing with `MinutesToExpire must be greater than 0` when `ExpirationTimeInMinutes` was not set and the client was created with `NewClient`. It now defaults to 60 minutes, matching `NewClientFromConnectionString`.
* Fixed `GenerateClientAccessURL` producing a malformed audience and client URL (for example `wss://<host>client/hubs/<hub>`) when the endpoint passed to `NewClient` has no trailing slash.
* `GenerateClientAccessURL` now validates a negative `ExpirationTimeInMinutes` on both credential types.

### Other Changes
* Regenerated code with the latest emitter.
* Updated dependencies.

## 0.1.1 (2026-03-26)

### Other Changes

- Remove duplicate `github.com/golang-jwt/jwt` dependency

## 0.1.0 (2024-02-27)

### Features Added

- Initial preview for the Web PubSub Service
