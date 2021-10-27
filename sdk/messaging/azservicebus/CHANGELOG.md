# Release History

## 0.2.0 (Unreleased)

### Features Added

- Scheduling messages to be delivered at a later date, via the `Sender.ScheduleMessage(s)` function or 
  setting `Message.ScheduledEnqueueTime`.
- Added in the `Sender.SendMessages([slice of sendable messages])` function, which batches messages 
  automatically. Useful when you're sending multiple messages that you are already sure will be small
  enough to fit into a single batch.
- Receiving from sessions using a SessionReceiver, created using Client.AcceptSessionFor(Queue|Subscription)
  or Client.AcceptNextSessionFor(Queue|Subscription).
- Can fully create, update, delete and list queues (and queue runtime properties) using the `AdministrationClient`.
- Can now renew a message lock for a ReceivedMessage using Receiver.RenewMessageLock()
- Can now renew a session lock for a SessionReceiver using SessionReceiver.RenewSessionLock()

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-05)

- Initial preview for the new version of the Azure Service Bus Go SDK. 
