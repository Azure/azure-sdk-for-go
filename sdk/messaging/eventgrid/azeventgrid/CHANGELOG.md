# Release History

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed
* HTTP redirects that leave the configured origin are no longer followed by the publisher client, preventing the publishing credential (`aeg-sas-key`, `aeg-sas-token`, or `Authorization`) and the event payload from being sent to a different host or port, or over an `https`-to-`http` downgrade. Same-origin redirects (and `http`-to-`https` upgrades of the same host) continue to be followed.

### Other Changes
* Regenerated code with the latest emitter.
* Updated dependencies.

## 1.0.0 (2024-04-09)

- GA for the Event Grid basic module.

## 0.1.0 (2024-03-05)

### Features Added

- Initial preview for the Event Grid Basic module. This module is the new home the `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher` package.
