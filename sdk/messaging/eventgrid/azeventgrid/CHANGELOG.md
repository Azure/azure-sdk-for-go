# Release History

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed
* The publisher client no longer follows HTTP redirects. Event Grid topic endpoints do not issue redirects; not following them prevents the publishing credential (`aeg-sas-key`, `aeg-sas-token`, or `Authorization`) and the event payload from being sent to a redirect target.

### Other Changes
* Regenerated code with the latest emitter.
* Updated dependencies.

## 1.0.0 (2024-04-09)

- GA for the Event Grid basic module.

## 0.1.0 (2024-03-05)

### Features Added

- Initial preview for the Event Grid Basic module. This module is the new home the `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher` package.
