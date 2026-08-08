# Release History

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed
* Credential headers (`aeg-sas-key`, `aeg-sas-token`, `Authorization`) are now stripped when an HTTP redirect crosses to a different host, preventing the publishing credential from being disclosed to a redirect target on another host.

### Other Changes
* Regenerated code with the latest emitter.
* Updated dependencies.

## 1.0.0 (2024-04-09)

- GA for the Event Grid basic module.

## 0.1.0 (2024-03-05)

### Features Added

- Initial preview for the Event Grid Basic module. This module is the new home the `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher` package.
