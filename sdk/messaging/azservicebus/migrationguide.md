# Guide to migrate from `azure-service-bus-go` to `azservicebus` 0.1.0

This guide is intended to assist in the migration from the pre-release `azure-service-bus-go` package to the latest beta releases (and eventual GA) of the `github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus`.

# Migration benefits

The redesign of the Service Bus SDK offers better integration with Azure Identity, a simpler API surface that allows you to uniformly work with queues, topics, subscriptions and subqueues (for instance: dead letter queues).

## Simplified API surface

The redesign for the API surface of Service Bus involves changing the way that clients are created. We wanted to simplify the number of types needed to get started, while also providing clarity on how, as a user of the SDK, to manage the resources the SDK creates (connections, links, etc...)

- [`Namespace` to `Client` migration](#namespace-to-client-migration)
- [Sending messages](#sending-messages)
- [Sending messages in batches](#sending-messages-in-batches)
- [Processing and receiving messages](#processing-and-receiving-messages)
- [Using dead letter queues](#using-dead-letter-queues)

### Namespace to Client migration

One big change is that the top level "client" is now [Client](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#Client), not `Namespace`:

Previous code:

```go
// previous code

ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString())
```

New (using `azservicebus`):

```go
// new code

client, err = azservicebus.NewClientWithConnectionString(connectionString, nil)
```

### Sending messages

Sending is done from a [Sender](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#Sender), which
works the same for queues or topics:

```go
sender, err := client.NewSender(queueOrTopicName)

sender.SendMessage(&azservicebus.Message{
  Body: []byte("hello world"),
})
```

### Sending messages in batches

Sending messages in batches is similar, except that the focus has been moved more
towards giving the user full control using the `MessageBatch` type.

```go
batch, err := sender.NewMessageBatch(ctx, nil)

// can be called multiple times
added, err := batch.Add(&azservicebus.Message{
  Body: []byte("hello world")
})

sender.SendMessage(ctx, batch)
```

### Processing and receiving messages

Receiving has split into two types:
- the [Processor](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#Processor), for continuously streaming messages to a user provided callback (similar to `Listen`).
- the [Receiver](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#Receiver), for receiving of messages in batches.

The `Processor` replaces the `Listen` functions on the previous Receiver type that could be created from a `Queue` or `Subscription`. It also adds in more robust error handling, which means that previous errors (like a link detaching) do not cause the Processor to exit. 

The Processor will only exit when you call `Close`.

```go
// NOTE: there is also NewProcessorForSubscription for subscriptions.
processor, err = client.NewProcessorForQueue(
  queueName,
  nil)

handleMessage := func(message *azservicebus.ReceivedMessage) error {
  log.Printf("Message arrived: %s", message.)
  processor.CompleteMessage(ctx, message)
}

handleError := func(err error) {
  // called whenever errors occur. Note, that unlike the
  // Listen, errors are automatically recovered.
}

// blocks until the Processor is closed.
processor.Start(ctx, handleMessage, handleError, nil)

// close at a time of your choosing
processor.Close(ctx)
```

### Receivers

Receivers allow you to request messages in batches, or easily receive a single message. This can be useful
for programs that need more control over when messages are received and when they are processed.

```go
receiver, err := client.NewReceiverForQueue(queue)
// or for a subscription
receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName)
```

`ReceiveOne` has been split into two functions to allow for receiving
multiple messages at a time (`ReceiveMessages`) or a single message (`ReceiveMessage`).

```go
// new code

// receiving multiple messages at a time, with a configurable timeout.
var messages []*azservicebus.ReceivedMessage
messages, err = receiver.ReceiveMessages(ctx, numMessages, nil)

// receiving a single message time, with a configurable timeout.
var message *azservicebus.ReceivedMessage

// this is similar to the `ReceiveOne`
message, err = receiver.ReceiveMessage(ctx, nil)
```

### Using dead letter queues

Previously, you created a receiver through an entity struct, like Queue or Subscription:

```go
// previous code

queue, err := ns.NewQueue()
deadLetterReceiver, err := queue.NewDeadLetterReceiver()

// or

topic, err := ns.NewTopic("topic")
subscription, err := topic.NewSubscription("subscription")
deadLetterReceiver, err := subscription.NewDeadLetterReceiver()

// the resulting receiver was a `ReceiveOner` which had different
// functions than some of the more full-fledged receiving types.
```

Now, in `azservicebus`:

```go
// new code

receiver, err = client.NewReceiverForQueue(
  queueName,
  &azservicebus.ReceiverOptions{
    ReceiveMode: azservicebus.PeekLock,
    SubQueue:    azservicebus.SubQueueDeadLetter,
  })

//or

receiver, err = client.NewReceiverForSubscription(
  topicName,
  subscriptionName,
  &azservicebus.ReceiverOptions{
    ReceiveMode: azservicebus.PeekLock,
    SubQueue:    azservicebus.SubQueueDeadLetter,
  })
```

The `Receiver` type for a dead letter queue is the same as the receiver for a 
queue or subscription, making things more consistent.

### Message settlement

Message settlement functions have moved to the `Receiver`, rather than being on the `Message`. 

Previously:

```go
// previous code

receiver.Listen(ctx, servicebus.HandlerFunc(func(c context.Context, m *servicebus.Message) error {
  m.Complete(ctx)
  return nil
}))
```

Now, using `azservicebus`:

```go
// new code

// with the Processor
processor.Start(ctx, func(msg *azservicebus.ReceivedMessage) {
  processor.CompleteMessage(ctx, msg)
})

// or with a Receiver
message, err := receiver.ReceiveMessage(ctx)   // or ReceiveMessages()
receiver.CompleteMessage(ctx, message)
```

# Azure Identity integration

Azure Identity has been directly integrated into the `Client` via the `NewClient()` function. This allows you to take advantage of conveniences like [DefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#section-readme) or any of the supported types within the package.

In `azservicebus`:

```go
credential, err := azidentity.NewDefaultAzureCredential(nil)
client, err = azservicebus.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)
```

# Upcoming features

Some features that are coming in the next beta:
- Management of entities (AdministrationClient) within Service Bus.
- Scheduling and cancellation of messages.
- Sending and receiving to Service Bus session enabled entities.

