# Release History

## 0.1.1 (Unreleased)

### Features Added

- Adding in the new Processor type, which can be used to do distributed (and load balanced) consumption of events, using a 
  CheckpointStore. The built-in CheckpointStore (in [github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints](github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints))
  uses Azure Blob storage for persistence. A full fledged example is in [example_processor_test.go](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go).

### Breaking Changes

- In the first beta, ConsumerClient took constructor parameter that required a partition ID, which meant you had to create
  multiple ConsumerClients if you wanted to consume multiple partitions. 
  
  This has been changed: ConsumerClient has a new method (NewPartitionClient) which allows you to create a PartitionClient,
  which targets a specific partition instead. Each PartitionClient also shares an AMQP connection.

## 0.1.0 (2022-08-11)

- Initial preview for the new version of the Azure Event Hubs Go SDK. 
