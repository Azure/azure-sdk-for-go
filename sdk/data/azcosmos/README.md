# Microsoft Azure Cosmos DB SDK for Go, Golang

## Introduction

This client library enables client applications to connect to Azure Cosmos via the SQL API. Azure Cosmos is a globally distributed, multi-model database service. 

Go: [cosmodb package](github.com/Azure/azure-sdk-for-go/tree/feature/cosmos/sdk/data/azcosmos)

## Getting Started

### Prerequisites

* Go versions 1.16 or higher
* An Azure subscription or free Azure Cosmos DB trial account

Note: If you don't have an Azure subscription, create a free account before you begin.
You can Try Azure Cosmos DB for free without an Azure subscription, free of charge and commitments, or create an Azure Cosmos DB free tier account, with the first 400 RU/s and 5 GB of storage for free. You can also use the Azure Cosmos DB Emulator with a URI of https://localhost:8081. For the key to use with the emulator, see Authenticating requests.

#### Create an Azure Cosmos DB account

* From the Azure portal menu or the Home page, select Create a resource.
* On the New page, search for and select Azure Cosmos DB.
* On the Azure Cosmos DB page, select [Create.](github.com/Azure/azure-sdk-for-go/tree/feature/cosmos/sdk/data/azcosmos)
* In the Create Azure Cosmos DB Account page, enter the basic settings for the new Azure Cosmos account.

#### Install the package
* Install the Azure Cosmos DB SDK for Go with `go get`:
  ```bash
  go get -u github.com/Azure/azure-sdk-for-go/tree/feature/cosmos/sdk/data/azcosmos
  ```
  
#### Authenticate the client

In order to interact with the Azure CosmosDB service you'll need to create an instance of the Cosmos Client class. To make this possible you will need an URL and key of the Azure CosmosDB service.

#### Create CosmosClient
```go
endpoint := "someEndpoint"
key := "someKey"
cred, _ := NewSharedKeyCredential(key)
	client, err := NewCosmosClient(endpoint, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
```

## Examples

The following section provides several code snippets covering some of the most common CosmosDB SQL API tasks, including:
* [Create Cosmos Client](#create-cosmos-client "Create Cosmos Client")
* [Create Database](#create-database "Create Database")
* [Create Container](#create-container "Create Container")
* [CRUD operation on Items](#crud-operation-on-items "CRUD operation on Items")

### Create Cosmos Client
```go
// Create a new CosmosClient via the CosmosClientBuilder
// It only requires endpoint and key, but other useful settings are available

const (
	CosmosDbEndpoint = "someEndpoint"
	CosmoDbKey = "someKey"
)

cred, _ := NewSharedKeyCredential(CosmosDbEndpoint)
	client := newCosmosClientConnection(CosmoDbKey, cred, &CosmosClientOptions{})
	if connection == nil {
		t.Error("Expected connection to be not nil")
	}
```

### Create Database
Using the client created in previous example, you can create a database like this:

```go
// Get a reference to the container
// This will create (or read) a database and its container.
    database := CosmosDatabaseProperties{Id: dbName}
	resp, err := client.CreateDatabase(ctx, database, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
```

### Create Container
Using the above created database for creating a container, like this:

```go
resp, err := database.CreateContainer(context, properties, throughput, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
container := resp.ContainerProperties.Container
```
### CRUD operation on Items

```go

// Create an item
database := client.createDatabase(t, context, client, "itemCRUD")
	properties := CosmosContainerProperties{
		Id: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}

	resp, err := database.CreateContainer(context, properties, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	item := map[string]string{
		"id":    "1",
		"value": "2",
	}

	container := resp.ContainerProperties.Container
	pk, err := NewPartitionKey("1")
	if err != nil {
		t.Fatalf("Failed to create pk: %v", err)
	}

	itemResponse, err := container.CreateItem(context, pk, item, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

```

For many more scenarios and examples see
[Azure-Samples/azure-sdk-for-go-samples][samples_repo].

## Next steps

- [Get Started APP](https://docs.microsoft.com/azure/cosmos-db/create-sql-api-dotnet-v4)
- [Github samples](https://github.com/Azure/azure-cosmos-dotnet-v3/tree/master/Microsoft.Azure.Cosmos.Samples/CodeSamples)
- [Resource Model of Azure Cosmos DB Service](https://docs.microsoft.com/azure/cosmos-db/sql-api-resources)
- [Cosmos DB Resource URI](https://docs.microsoft.com/rest/api/documentdb/documentdb-resource-uri-syntax-for-rest)
- [Partitioning](https://docs.microsoft.com/azure/cosmos-db/partition-data)
- [Introduction to SQL API of Azure Cosmos DB Service](https://docs.microsoft.com/azure/cosmos-db/sql-api-sql-query)
- [SDK API](https://docs.microsoft.com/dotnet/api/azure.cosmos?view=azure-dotnet)
- [Using emulator](https://github.com/Azure/azure-documentdb-dotnet/blob/master/docs/documentdb-nosql-local-emulator.md)
- [Capture traces](https://github.com/Azure/azure-documentdb-dotnet/blob/master/docs/documentdb-sdk_capture_etl.md)
- [Release notes](https://github.com/Azure/azure-cosmos-dotnet-v3/blob/v4/changelog.md)
- [Diagnose and troubleshooting](https://docs.microsoft.com/azure/cosmos-db/troubleshoot-dot-net-sdk)


## License

This project is licensed under MIT.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Azure.data.azcosmos` label.

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

<!-- LINKS -->
[nuget_package]: https://www.nuget.org/packages/Azure.Cosmos
[cosmos_emulator]: https://docs.microsoft.com/azure/cosmos-db/local-emulator
[cosmos_resource_portal]: https://docs.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[cosmos_resource_cli]: https://docs.microsoft.com/azure/cosmos-db/scripts/cli/sql/create
[cosmos_resource_arm]: https://docs.microsoft.com/azure/cosmos-db/quick-create-template
[cosmos_throughput]: https://docs.microsoft.com/azure/cosmos-db/set-throughput
[cosmos_partition]: https://docs.microsoft.com/azure/cosmos-db/partitioning-overview#choose-partitionkey
[cosmos_optimistic]: https://docs.microsoft.com/azure/cosmos-db/database-transactions-optimistic-concurrency#optimistic-concurrency-control
[cosmos_scripts]: https://docs.microsoft.com/azure/cosmos-db/how-to-write-stored-procedures-triggers-udfs
[cosmos_resourcemodel]: https://docs.microsoft.com/azure/cosmos-db/databases-containers-items

# Resources

- SDK docs are at [godoc.org](https://godoc.org/github.com/Azure/azure-sdk-for-go/).
- SDK samples are at [Azure-Samples/azure-sdk-for-go-samples](https://github.com/Azure-Samples/azure-sdk-for-go-samples).
- SDK notifications are published via the [Azure update feed](https://azure.microsoft.com/updates/).
- Azure API docs are at [docs.microsoft.com/rest/api](https://docs.microsoft.com/rest/api/).
- General Azure docs are at [docs.microsoft.com/azure](https://docs.microsoft.com/azure).

## License

Apache 2.0, see [LICENSE](./LICENSE).

## Contribute

See [CONTRIBUTING.md](./CONTRIBUTING.md).

[samples_repo]: https://github.com/Azure-Samples/azure-sdk-for-go-samples