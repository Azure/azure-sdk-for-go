# Release History

## 0.4.1 (Unreleased)

### Features Added

### Breaking Changes

- This module has been moved from it's previous location in `azeventgrid` to this location (`github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces`).

### Bugs Fixed

### Other Changes

## 0.4.0 (2023-11-27)

### Features Added

- New functionality for Event Grid namespaces: 
  - Client.PublishCloudEvent can be used to publish a single `messaging.CloudEvent`.
  - Client.RenewCloudEventLocks can extend the lock time for a set of events.
  - Client.ReleaseCloudEvents (via ReleaseCloudEventsOptions.ReleaseDelayInSeconds) can release an event with a 
    server-side delay, allowing the message to remain unavailable for a configured period of time.

### Breaking Changes

- FailedLockToken, included in the response for settlement functions, has an `Error` field, which contains the data previously
  in `ErrorDescription` and `ErrorCode`.
- Settlement functions (AcknowledgeCloudEvents, ReleaseCloudEvents, RejectCloudEvents) take lock tokens as a parameter.

## 0.3.0 (2023-10-17)

### Breaking Changes

- Client constructors that take a `key string` parameter for a credential now require an `*azcore.KeyCredential` or `*azcore.SASCredential`.

## 0.2.0 (2023-09-12)

### Features Added

- The publisher client for Event Grid topics has been added as a sub-package under `publisher`.

### Other Changes

- Documentation and examples added for Event Grid namespace client.

## 0.1.0 (2023-07-11)

### Features Added

- Initial preview for the Event Grid package for Event Grid Namespaces
