# Release History

## 0.1.2 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
