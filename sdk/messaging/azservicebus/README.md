# Azure Service Bus Client Module for Go

[Azure Service Bus](https://azure.microsoft.com/services/service-bus/) is a highly-reliable cloud messaging service from Microsoft.

Use the client library `github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus` in your application to:

- Send messages to an Azure Service Bus Queue or Topic
- Receive messages from an Azure Service Bus Queue or Subscription

**NOTE**: This library is currently a preview. There may be breaking interface changes until it reaches semantic version `v1.0.0`.

Key links:
- [Source code][source]
- [API Reference Documentation](https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus)
- [Product documentation](https://azure.microsoft.com/services/service-bus/)
- [Samples][godoc_samples]

## Getting started

### Install the package

Install the Azure Service Bus client module for Go with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus
```

### Prerequisites
- Go, version 1.16 or higher
- An [Azure subscription](https://azure.microsoft.com/free/)
- A [Service Bus Namespace](https://docs.microsoft.com/azure/service-bus-messaging/) 

### Authenticate the client

The Service Bus [Client][gopkg_client] can be created using a Service Bus connection string or a credential from the [Azure Identity package][azure_identity_pkg], like [DefaultAzureCredential][default_azure_credential].

#### Using a connection string

```go
import (
  "log"
  "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
  client, err := azservicebus.NewClient("<Service Bus connection string>")
 
  if err != nil {
    log.Fatalf("Failed to create Service Bus Client: %s", err.Error())
  }
}
```

#### Using an Azure Active Directory Credential

```go
import (
  "log"
  "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
  "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
  // For more information about the DefaultAzureCredential:
  // https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
  cred, err := azidentity.NewDefaultAzureCredential(nil)

  if err != nil {
    log.Fatalf("Failed creating DefaultAzureCredential: %s", err.Error())
  }

  client, err := azservicebus.NewClient("<ex: my-service-bus.servicebus.windows.net>", cred)

  if err != nil {
    log.Fatalf("Failed to create Service Bus Client: %s", err.Error())
  }
}
```

## Key concepts

Once you've created a [Client][godoc_client], you can interact with resources within a Service Bus Namespace:

- [Queues][queue_concept]: Allows for sending and receiving messages. Often used for point-to-point communication.
- [Topics][topic_concept]: As opposed to Queues, Topics are better suited to publish/subscribe scenarios. A topic can be sent to, but requires a subscription, of which there can be multiple in parallel, to consume from.
- [Subscriptions][subscription_concept]: The mechanism to consume from a Topic. Each subscription is independent, and receives a copy of each message sent to the topic. Rules and Filters can be used to tailor which messages are received by a specific subscription.

For more information about these resources, see [What is Azure Service Bus?][service_bus_overview].

Using a `Client` you can do the following:

- Send messages, to a queue or topic, using a [Sender][godoc_sender] created using [Client.NewSender()][godoc_newsender]. [sample link]()
- Receive messages, from either a queue or a subscription, using a [Receiver][godoc_receiver] created using [client.NewReceiverForQueue()][godoc_newreceiver_queue] or [client.NewReceiverForSubscription()][godoc_newreceiver_subscription]

Please note that the Queues, Topics and Subscriptions should be created prior to using this library.

### Samples

Find up-to-date examples and documentation on [godoc.org](https://godoc.org/github.com/Azure/azure-service-bus-go#pkg-examples).

## Examples

The following sections provide code snippets that cover some of the common tasks using Azure Service Bus

- [Send messages](#send-messages)

### Send messages

Once you've created a [Client][godoc_client] you can create a [Sender][godoc_sender], which will allow you to send messages.

NOTE: Creating a `client` is covered in the ["Authenticate the client"](#authenticate-the-client) section of the readme.

```go
sender, err := client.NewSender("<queue or topic>")

if err != nil {
  log.Fatalf("Failed to create Sender: %s", err.Error())
}

// send a single message
err = sender.SendMessage(context.TODO(), &azservicebus.Message{
  Body: []byte("hello world!"),
})
```

You can also send messages in batches, which can be more efficient than sending them individually

```go
// Create a message batch. It will automatically be sized for the Service Bus
// Namespace's maximum message size.
messageBatch, err := sender.NewMessageBatch(context.TODO())

if err != nil {
  log.Fatalf("Failed to create a message batch: %s", err.Error())
}

// Add a message using TryAdd.
// This can be called multiple times, and will return (false, nil)
// if the message cannot be added because the batch is full.
added, err := messageBatch.TryAdd(&azservicebus.Message{
    Body: []byte(fmt.Sprintf("hello world")),
})

if err != nil {
  log.Fatalf("Failed to add message to batch because of an error: %s", err.Error())
}

if !added {
  log.Printf("Message batch is full. We should send it and create a new one.")
  err := sender.SendMessageBatch(context.TODO(), messageBatch)

  if err != nil {
    log.Fatalf("Failed to send message batch: %s", err.Error())
  }

  // add the next message to a new batch and start again.
}
```

## Next steps

Please take a look at the [samples][godoc_samples] for detailed examples on how to use this library to send and receive messages to/from [Service Bus Queues, Topics and Subscriptions](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-messaging-overview).

## Contributing

If you'd like to contribute to this library, please read the [contributing guide](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md) to learn more about how to build and test the code.

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fmessaging%2Fazservicebus%2FREADME.png)

[source]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/messaging/azservicebus
[new_issue]: https://github.com/Azure/azure-sdk-for-go/issues/new
[azure_identity_pkg]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[default_azure_credential]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
[queue_concept]: https://docs.microsoft.com/azure/service-bus-messaging/service-bus-messaging-overview#queues
[topic_concept]: https://docs.microsoft.com/azure/service-bus-messaging/service-bus-messaging-overview#topics
[subscription_concept]: https://docs.microsoft.com/azure/service-bus-messaging/service-bus-queues-topics-subscriptions#topics-and-subscriptions
[service_bus_overview]: https://docs.microsoft.com/azure/service-bus-messaging/service-bus-messaging-overview
[godoc_samples]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#pkg-examples
[godoc_receiver]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Receiver
[godoc_client]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Client
[godoc_sender]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Sender
[godoc_newsender]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Client.NewSender
[godoc_newreceiver_queue]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Client.NewReceiverWithQueue
[godoc_newreceiver_subscription]: https://godoc.org/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/#Client.NewReceiverWithQueue