# Release History

## 0.6.0 (Unreleased)

### Features Added

- Added the `ConsumerClientOptions.InstanceID` field. This optional field can enhance error messages from 
  Event Hubs. For example, error messages related to ownership changes for a partition will contain the 
  name of the link that has taken ownership, which can help with traceability.

### Breaking Changes

- `ConsumerClient.ID()` renamed to `ConsumerClient.InstanceID()`.

### Bugs Fixed

### Other Changes

## 0.5.0 (2023-02-07)

### Features Added

- Adds ProcessorOptions.Prefetch field, allowing configuration of Prefetch values for PartitionClients created using the Processor. (PR#19786)
- Added new function to parse connection string into values using `ParseConnectionString` and `ConnectionStringProperties`. (PR#19855)

### Breaking Changes

- ProcessorOptions.OwnerLevel has been removed. The Processor uses 0 as the owner level.
- Uses the public release of `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob` package rather than using an internal copy. 
  For an example, see [example_processor_test.go](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go).

## 0.4.0 (2023-01-10)

### Bugs Fixed

- User-Agent was incorrectly formatted in our AMQP-based clients. (PR#19712)
- Connection recovery has been improved, removing some unnecessasry retries as well as adding a bound around 
  some operations (Close) that could potentially block recovery for a long time. (PR#19683)

## 0.3.0 (2022-11-10)

### Bugs Fixed

- $cbs link is properly closed, even on cancellation (#19492)

### Breaking Changes

- ProducerClient.SendEventBatch renamed to ProducerClient.SendEventDataBatch, to align with
  the name of the type.

## 0.2.0 (2022-10-17)

### Features Added

- Raw AMQP message support, including full support for encoding Body (Value, Sequence and also multiple byte slices for Data). See ExampleEventDataBatch_AddEventData_rawAMQPMessages for some concrete examples. (PR#19156)
- Prefetch is now enabled by default. Prefetch allows the Event Hubs client to maintain a continuously full cache of events, controlled by PartitionClientOptions.Prefetch. (PR#19281)
- ConsumerClient.ID() returns a unique ID representing each instance of ConsumerClient.

### Breaking Changes

- EventDataBatch.NumMessages() renamed to EventDataBatch.NumEvents()
- Prefetch is now enabled by default. To disable it set PartitionClientOptions.Prefetch to -1.
- NewWebSocketConnArgs renamed to WebSocketConnParams
- Code renamed to ErrorCode, including associated constants like `ErrorCodeOwnershipLost`.
- OwnershipData, CheckpointData, and CheckpointStoreAddress have been folded into their individual structs: Ownership and Checkpoint.
- StartPosition and OwnerLevel were erroneously included in the ConsumerClientOptions struct - they've been removed. These can be 
  configured in the PartitionClientOptions.

### Bugs Fixed

- Retries now respect cancellation when they're in the "delay before next try" phase. (PR#19295)
- Fixed a potential leak which could cause us to open and leak a $cbs link connection, resulting in errors. (PR#19326)

## 0.1.1 (2022-09-08)

### Features Added

- Adding in the new Processor type, which can be used to do distributed (and load balanced) consumption of events, using a 
  CheckpointStore. The built-in checkpoints.BlobStore uses Azure Blob Storage for persistence. A full example is 
  in [example_processor_test.go](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go).

### Breaking Changes

- In the first beta, ConsumerClient took constructor parameter that required a partition ID, which meant you had to create
  multiple ConsumerClients if you wanted to consume multiple partitions. ConsumerClient can now create multiple PartitionClient
  instances (using ConsumerClient.NewPartitionClient), which allows you to share the same AMQP connection and receive from multiple
  partitions simultaneously.
- Changes to EventData/ReceivedEventData:
  - ReceivedEventData now embeds EventData for fields common between the two, making it easier to change and resend.
  - `ApplicationProperties` renamed to `Properties`.
  - `PartitionKey` removed from `EventData`. To send events using a PartitionKey you must set it in the options
    when creating the EventDataBatch:

    ```go
    batch, err := producerClient.NewEventDataBatch(context.TODO(), &azeventhubs.NewEventDataBatchOptions{
		  PartitionKey: to.Ptr("partition key"),
	  })
    ```

### Bugs Fixed

- ReceivedEventData.Offset was incorrectly parsed, resulting in it always being 0.
- Added missing fields to ReceivedEventData and EventData (CorrelationID)
- PartitionKey property was not populated for messages sent via batch.

## 0.1.0 (2022-08-11)

- Initial preview for the new version of the Azure Event Hubs Go SDK. 
