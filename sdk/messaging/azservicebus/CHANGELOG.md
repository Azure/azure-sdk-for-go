# Release History

## 0.3.4 (Unreleased)

### Bugs Fixed

- Fix unaligned 64-bit atomic operation on mips.  Thanks to @jackesdavid for contributing this fix. (#16847)
- Multiple fixes to address connection/link recovery (#16831)
- Fixing panic() when the links haven't been initialized (early cancellation) (#16941)
- Handle 500 as a retryable code (no recovery needed) (#16925)

## 0.3.3 (2022-01-12)

### Features Added

- Support the pass-through of an Application ID when constructing an Azure Service Bus Client. PR#16558 (thanks halspang!)

### Bugs Fixed 

- Fixing connection/link recovery in Sender.SendMessages() and Sender.SendMessageBatch(). PR#16790
- Fixing bug in the management link which could cause it to panic during recovery. PR#16790

## 0.3.2 (2021-12-08)

### Features Added

- Enabling websocket support via `ClientOptions.NewWebSocketConn`. For an example, see the `ExampleNewClient_usingWebsockets` 
  function in `example_client_test.go`.

### Breaking Changes

- Message properties that come from the standard AMQP message have been made into pointers, to allow them to be 
  properly omitted (or indicate that they've been omitted) when sending and receiving.  

### Bugs Fixed

- Session IDs can now be blank - prior to this release it would cause an error. PR#16530
- Drain will no longer hang if there is a link failure. Thanks to @flexarts for reporting this issue: PR#16530
- Attempting to settle messages received in ReceiveAndDelete mode would cause a panic. PR#16255

### Other Changes
- Removed legacy dependencies, resulting in a much smaller package.

## 0.3.1 (2021-11-16)

### Bugs Fixed

- Updating go-amqp to v0.16.4 to fix a race condition found when running `go test -race`.  Thanks to @peterzeller for reporting this issue. PR: #16168

## 0.3.0 (2021-11-12)

### Features Added

- AbandonMessage and DeferMessage now take an additional `PropertiesToModify` option, allowing
  the message properties to be modified when they are settled.
- Missing fields for entities in the admin.Client have been added (UserMetadata, etc..)

### Breaking Changes

- AdminClient has been moved into the `admin` subpackage.
- ReceivedMessage.Body is now a function that returns a ([]byte, error), rather than being a field.
  This protects against a potential data-loss scenario where a message is received with a payload 
  encoded in the sequence or value sections of an AMQP message, which cannot be prpoerly represented
  in the .Body. This will now return an error.
- Functions that have options or might have options in the future have an additional *options parameter.
  As usual, passing 'nil' ignores the options, and will cause the function to use defaults.
- MessageBatch.Add() has been renamed to MessageBatch.AddMessage(). AddMessage() now returns only an `error`, 
  with a sentinel error (ErrMessageTooLarge) signaling that the batch cannot fit a new message.
- Sender.SendMessages() has been removed in favor of simplifications made in MessageBatch.

### Bugs Fixed

- ReceiveMessages has been tuned to match the .NET limits (which has worked well in practice). This partly addresses #15963, 
  as our default limit was far higher than needed.

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
