# Event Hubs package reliability tests

These are the stress/reliability tests for the `azeventhubs` package.

The entrypoint for the tests is [`stress.go`](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/internal/eh/stress/stress.go). All of the individual tests are in the `tests` sub-folder. These tests should run fine on your local machine - you'll need to create an `.env` file, in the root of the `azeventhubs` module, with the following values:

```bash
EVENTHUB_CONNECTION_STRING=<connection string to the Event Hubs namespace>
EVENTHUB_NAME=<already created Event Hub - should have at least 4 partitions>
CHECKPOINTSTORE_STORAGE_CONNECTION_STRING=<connection string to an Azure Storage account>
APPINSIGHTS_INSTRUMENTATIONKEY=<instrumentation key for appinsights>
```

There are two types of tests - batch and processor. Each test takes a variety of flags to control the duration, number of events, etc..

For instance, to run a `Processor` test to receive events:

```bash
go run . processor
```

To see more options just run:

```bash
go run . processor --help
```

For convenience there's a deploy.ps1 file that'll launch the deployment - by default it'll go to the `pg` cluster, which we use for adhoc workloads. For more information about prerequisites look at the official stress test docs here: [stress test readme.md](https://github.com/Azure/azure-sdk-tools/tree/main/tools/stress-cluster/chaos).
