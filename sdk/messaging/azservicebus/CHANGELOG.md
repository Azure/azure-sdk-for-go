# Release History

## 0.2.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.2.0 (2021-11-02)

### Features Added

- Scheduling messages to be delivered at a later date, via the `Sender.ScheduleMessage(s)` function or 
  setting `Message.ScheduledEnqueueTime`.
- Added in the `Sender.SendMessages([slice of sendable messages])` function, which batches messages 
  automatically. Useful when you're sending multiple messages that you are already sure will be small
  enough to fit into a single batch.
- Receiving from sessions using a SessionReceiver, created using Client.AcceptSessionFor(Queue|Subscription)
  or Client.AcceptNextSessionFor(Queue|Subscription).
- Can fully create, update, delete and list queues, topics and subscriptions using the `AdministrationClient`.
- Can renew message and session locks, using Receiver.RenewMessageLock() and SessionReceiver.RenewSessionLock(), respectively.

### Bugs Fixed

- Receiver.ReceiveMessages() had a bug where multiple calls could result in the link no longer receiving messages.
  This was fixed with an update in go-amqp.

## 0.1.0 (2021-10-05)

- Initial preview for the new version of the Azure Service Bus Go SDK. 
