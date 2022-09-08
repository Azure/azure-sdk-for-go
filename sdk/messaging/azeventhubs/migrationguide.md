# Guide to migrate from `azure-event-hubs-go` to `azeventhubs`

This guide is intended to assist in the migration from the `azure-event-hubs-go` package to the latest beta releases (and eventual GA) of the `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs`.

## Clients

The `Hub` type has been replaced by two types:

* Consuming events, using the `azeventhubs.ConsumerClient`: [docs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#ConsumerClient) | [example](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_consuming_events_test.go)
* Sending events, use the `azeventhubs.ProducerClient`: [docs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#ProducerClient) | [example](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_producing_events_test.go)

`EventProcessorHost` has been replaced by the `azeventhubs.Processor` type: [docs](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#Processor) | [example](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go)

## Authentication

The older Event Hubs package provided some authentication methods like hub.NewHubFromEnvironment. These have been replaced by by using Azure Identity credentials from [azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#section-readme). 

You can also still authenticate using connection strings.

* `azeventhubs.ConsumerClient`: [using azidentity](https://github.com/Azure/azure-sdk-for-go/blob/a46bd74e113d6a045541b82a0f3f6497011d8417/sdk/messaging/azeventhubs/example_consumerclient_test.go#L16) | [using a connection string](https://github.com/Azure/azure-sdk-for-go/blob/a46bd74e113d6a045541b82a0f3f6497011d8417/sdk/messaging/azeventhubs/example_consumerclient_test.go#L30)

* `azeventhubs.ProducerClient`: [using azidentity](https://github.com/Azure/azure-sdk-for-go/blob/a46bd74e113d6a045541b82a0f3f6497011d8417/sdk/messaging/azeventhubs/example_producerclient_test.go#L16) | [using a connection string](https://github.com/Azure/azure-sdk-for-go/blob/a46bd74e113d6a045541b82a0f3f6497011d8417/sdk/messaging/azeventhubs/example_producerclient_test.go#L30)

## EventBatchIterator

Sending events has changed to be more explicit about when batches are formed and sent.

The older module had a type (EventBatchIterator). This type has been removed and replaced
with explicit batching, using `azeventhubs.EventDataBatch`. See here for an example: [link](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_producing_events_test.go).

## Getting hub/partition information

In the older module functions to get the partition IDs, as well as runtime properties
like the last enqueued sequence number were on the `Hub` type. These are now on both
of the client types instead (`ProducerClient`, `ConsumerClient`).

```go
// old
hub.GetPartitionInformation(context.TODO(), "0")
hub.GetRuntimeInformation(context.TODO())
```

```go
// new

// equivalent to: hub.GetRuntimeInformation(context.TODO())
consumerClient.GetEventHubProperties(context.TODO(), nil)   

// equivalent to: hub.GetPartitionInformation
consumerClient.GetPartitionProperties(context.TODO(), "partition-id", nil)  

//
// or, using the ProducerClient
//

producerClient.GetEventHubProperties(context.TODO(), nil)
producerClient.GetPartitionProperties(context.TODO(), "partition-id", nil)
```

