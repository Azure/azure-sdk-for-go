# Microsoft Azure Cosmos DB SDK for Go, Golang

## Introduction

This client library enables client applications to connect to Azure Cosmos via the SQL API. Azure Cosmos is a globally distributed, multi-model database service.

## Getting Started

### Prerequisites

* Go versions 1.16 or higher
* An Azure subscription or free Azure Cosmos DB trial account

Note: If you don't have an Azure subscription, create a free account before you begin.
You can Try Azure Cosmos DB for free without an Azure subscription, free of charge and commitments, or create an Azure Cosmos DB free tier account, with the first 400 RU/s and 5 GB of storage for free. You can also use the Azure Cosmos DB Emulator with a URI of https://localhost:8081. For the key to use with the emulator, see Authenticating requests.

### Create an Azure Cosmos account

You can create an Azure Cosmos account using:

* [Azure Portal](https://portal.azure.com).
* [Azure CLI](https://docs.microsoft.com/cli/azure).
* [Azure ARM](https://docs.microsoft.com/azure/cosmos-db/quick-create-template).

#### Install the package

* Install the Azure Cosmos DB SDK for Go with `go get`:

  ```bash
  go get -u github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos
  ```

#### Authenticate the client

In order to interact with the Azure CosmosDB service you'll need to create an instance of the Cosmos client class. To make this possible you will need an URL and key of the Azure CosmosDB service.

## Examples

The following section provides several code snippets covering some of the most common CosmosDB SQL API tasks, including:
* [Create Client](#create-cosmos-client "Create Cosmos client")
* [Create Database](#create-database "Create Database")
* [Create Container](#create-container "Create Container")
* [CRUD operation on Items](#crud-operation-on-items "CRUD operation on Items")

### Create Cosmos Client

```go
const (
    cosmosDbEndpoint = "someEndpoint"
    cosmosDbKey = "someKey"
)

cred, _ := azcosmos.NewKeyCredential(cosmosDbKey)
client, err := azcosmos.NewClientWithKey(cosmosDbEndpoint, cred, nil)
handle(err)
```

### Create Database

Using the client created in previous example, you can create a database like this:

```go
database := azcosmos.DatabaseProperties{Id: dbName}
response, err := client.CreateDatabase(context, database, nil)
handle(err)
database, err := azcosmos.NewDatabase(dbName)
handle(err)
```

### Create Container

Using the above created database for creating a container, like this:

```go
properties := azcosmos.ContainerProperties{
    Id: "aContainer",
    PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
        Paths: []string{"/id"},
    },
}

throughput := azcosmos.NewManualThroughputProperties(400)
response, err := database.CreateContainer(context, properties, &CreateContainerOptions{ThroughputProperties: &throughput})
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

// Delete an item
itemResponse, err = container.DeleteItem(context, pk, id, nil)
handle(err)
```

## Next steps

- [Resource Model of Azure Cosmos DB Service](https://docs.microsoft.com/azure/cosmos-db/sql-api-resources)
- [Cosmos DB Resource URI](https://docs.microsoft.com/rest/api/documentdb/documentdb-resource-uri-syntax-for-rest)
- [Partitioning](https://docs.microsoft.com/azure/cosmos-db/partition-data)
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

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go/sdk/data/azcosmos/README.png)
