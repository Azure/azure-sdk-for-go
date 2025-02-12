# Azure Cosmos DB SDK for Go

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Getting Started

### Prerequisites

* Go versions 1.21 or higher
* An Azure subscription or free Azure Cosmos DB trial account

Note: If you don't have an Azure subscription, create a free account before you begin.
You can Try Azure Cosmos DB for free without an Azure subscription, free of charge and commitments, or create an Azure Cosmos DB free tier account, with the first 400 RU/s and 5 GB of storage for free. You can also use the Azure Cosmos DB Emulator with a URI of https://localhost:8081. For the key to use with the emulator, see [how to develop with the emulator](https://learn.microsoft.com/azure/cosmos-db/how-to-develop-emulator).

### Create an Azure Cosmos DB account

You can create an Azure Cosmos DB account using:

* [Azure Portal](https://portal.azure.com).
* [Azure CLI](https://learn.microsoft.com/cli/azure).
* [Azure ARM](https://learn.microsoft.com/azure/cosmos-db/quick-create-template).

#### Install the package

* Install the Azure Cosmos DB SDK for Go with `go get`:

  ```bash
  go get -u github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos
  ```

#### Authenticate the client

In order to interact with the Azure Cosmos DB service you'll need to create an instance of the `Client` struct. To make this possible you will need a URL and key of the Azure Cosmos DB service.

#### Logging

The SDK can make use of `azcore`'s logging implementation to collect useful information for debugging your application. In order to make use of logs, one must set the environment variable `"AZURE_SDK_GO_LOGGING"` to `"all"` like outlined in this [public document](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#hdr-Built_in_Logging).

Once that is done, the SDK will begin to collect diagnostics. By default, it will output the logs to `stdout` - printing directly to your console - and will record all types of events (requests, responses, retries). If you'd like to configure a listener that acts differently, the small snippet below shows how you could do so.

```go
import (
	"os"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
)

f, err := os.Create("cosmos-log-file.txt")
handle(err)
defer f.Close()

// Configure the listener to write to a file rather than to the console
azlog.SetListener(func(event azlog.Event, s string) {
	f.WriteString(s + "\n")
})

// Filter the types of events you'd like to log by removing the ones you're not interested in (if any)
// We recommend using the default logging with no filters - but if filtering we recommend *always* including 
// `azlog.EventResponseError` since this is the event type that will help with debugging errors
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse, azlog.EventRetryPolicy, azlog.EventResponseError) 
```

## Examples

The following section provides several code snippets covering some of the most common Azure Cosmos DB NoSQL API tasks, including:
* [Create Client](#create-cosmos-db-client "Create Cosmos DB client")
* [Create Database](#create-database "Create Database")
* [Create Container](#create-container "Create Container")
* [CRUD operation on Items](#crud-operation-on-items "CRUD operation on Items")

### Create Cosmos DB Client

The clients support different forms of authentication. The azcosmos library supports authorization via Microsoft Entra identities or an account key.

**Using Microsoft Entra identities**

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

cred, err := azidentity.NewDefaultAzureCredential(nil)
handle(err)
client, err := azcosmos.NewClient("myAccountEndpointURL", cred, nil)
handle(err)
```

**Using account keys**

```go
const (
    cosmosDbEndpoint = "someEndpoint"
    cosmosDbKey = "someKey"
)

cred, err := azcosmos.NewKeyCredential(cosmosDbKey)
handle(err)
client, err := azcosmos.NewClientWithKey(cosmosDbEndpoint, cred, nil)
handle(err)
```

### Create Database

Using the client created in previous example, you can create a database like this:

```go
databaseProperties := azcosmos.DatabaseProperties{ID: dbName}
response, err := client.CreateDatabase(context, databaseProperties, nil)
handle(err)
database, err := client.NewDatabase(dbName)
handle(err)
```

### Create Container

Using the above created database for creating a container, like this:

```go
properties := azcosmos.ContainerProperties{
    ID: "aContainer",
    PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
        Paths: []string{"/id"},
    },
}

throughput := azcosmos.NewManualThroughputProperties(400)
response, err := database.CreateContainer(context, properties, &azcosmos.CreateContainerOptions{ThroughputProperties: &throughput})
handle(err)
```

### CRUD operation on Items

```go
item := map[string]string{
    "id":    "1",
    "value": "2",
}

marshalled, err := json.Marshal(item)
if err != nil {
    log.Fatal(err)
}

container, err := client.NewContainer(dbName, containerName)
handle(err)

pk := azcosmos.NewPartitionKeyString("1")
id := "1"

// Create an item
itemResponse, err := container.CreateItem(context, pk, marshalled, nil)
handle(err)

// Read an item
itemResponse, err = container.ReadItem(context, pk, id, nil)
handle(err)

var itemResponseBody map[string]string
err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
if err != nil {
    log.Fatal(err)
}

itemResponseBody["value"] = "3"
marshalledReplace, err := json.Marshal(itemResponseBody)
if err != nil {
    log.Fatal(err)
}

// Replace an item
itemResponse, err = container.ReplaceItem(context, pk, id, marshalledReplace, nil)
handle(err)

// Patch an item
patch := PatchOperations{}
patch.AppendAdd("/newField", "newValue")
patch.AppendRemove("/oldFieldToRemove")

itemResponse, err := container.PatchItem(context.Background(), pk, id, patch, nil)
handle(err)

// Delete an item
itemResponse, err = container.DeleteItem(context, pk, id, nil)
handle(err)
```

## Next steps

- [Resource Model of Azure Cosmos DB Service](https://learn.microsoft.com/azure/cosmos-db/sql-api-resources)
- [Azure Cosmos DB Resource URI](https://learn.microsoft.com/rest/api/documentdb/documentdb-resource-uri-syntax-for-rest)
- [Partitioning](https://learn.microsoft.com/azure/cosmos-db/partition-data)
- [Using emulator](https://github.com/Azure/azure-documentdb-dotnet/blob/master/docs/documentdb-nosql-local-emulator.md)


## License

This project is licensed under MIT.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Cosmos` label.

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License
Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For
details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.


