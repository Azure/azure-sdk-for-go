# Service Bus package reliability tests

These are the stress/reliability tests for the `azservicebus` package.

The entrypoint for the tests is [`stress.go`](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azservicebus/internal/stress/stress.go). All of the individual tests are in the `tests` sub-folder. These tests should run fine on your local machine - you'll need to create an `.env` file with the following values:

```bash
SERVICEBUS_CONNECTION_STRING=<connection string to the Service Bus namespace>
APPINSIGHTS_INSTRUMENTATIONKEY=<instrumentation key for appinsights>
```

To run one of the more basic tests, where we just send and receive messages for a few days:

```bash
go run . tests infiniteSendAndReceive
```

To see all the tests that are available:

```bash
go run . tests
```

For convenience there's a deploy.ps1 file that'll launch the deployment - by default it'll go to the `pg` cluster, which we use for adhoc workloads. For more information about prerequisites look at the official stress test docs here: [stress test readme.md](https://github.com/Azure/azure-sdk-tools/tree/main/tools/stress-cluster/chaos).
