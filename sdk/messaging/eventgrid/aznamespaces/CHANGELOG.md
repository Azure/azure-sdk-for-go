# Release History

## 1.0.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 1.0.0 (2024-06-11)

### Features Added

- First stable release of the aznamespaces package targeted at API version `2024-06-01`.

### Breaking Changes

- Sending and receiving operations have been moved to separate clients (SenderClient and ReceiverClient).
- Method names have been shortened from <Verb>CloudEvent(s) to <Verb>Event(s)
- LockTokens for AcknowledgeEvents, RejectEvents and ReleaseEvents are now a positional argument, instead of optional.
- Topic and subscription name are now set at the Client level, as part of `NewSenderClient` or `NewReceiverClient`.

## 0.4.1 (2024-03-05)

### Breaking Changes

- This module has been moved from its previous location in `azeventgrid` to this location (`github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces`).

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
