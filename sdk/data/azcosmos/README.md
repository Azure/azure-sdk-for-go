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

## Connection Modes

The SDK supports two connection modes for communicating with Azure Cosmos DB:

### Gateway Mode (Default)

Gateway mode routes all requests through the Azure Cosmos DB gateway using HTTPS. This is the default mode and works with any network configuration.

```go
client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
```

### Direct Mode

Direct mode connects directly to Azure Cosmos DB backend nodes for document operations using the RNTBD (binary) protocol over TCP. This provides lower latency and higher throughput compared to Gateway mode.

```go
options := &azcosmos.ClientOptions{
    ConnectionMode: azcosmos.ConnectionModeDirect,
}
client, err := azcosmos.NewClientWithKey(endpoint, cred, options)
```

#### Direct Mode Configuration Options

You can tune Direct mode connection behavior with these options:

```go
options := &azcosmos.ClientOptions{
    ConnectionMode: azcosmos.ConnectionModeDirect,
    DirectModeOptions: &azcosmos.DirectModeOptions{
        MaxRequestsPerConnection:  50,              // Max concurrent requests per TCP connection (default: 30)
        IdleConnectionTimeout:     10 * time.Minute, // Idle connection timeout (default: server value)
        MaxConnectionsPerEndpoint: 20,              // Max TCP connections per backend (default: 10)
        ConnectTimeout:            10 * time.Second, // TCP connect timeout (default: 5s)
    },
}
client, err := azcosmos.NewClientWithKey(endpoint, cred, options)
```

| Option | Default | Description |
|--------|---------|-------------|
| `MaxRequestsPerConnection` | 30 | Maximum concurrent requests per TCP connection. Higher values improve throughput but increase memory usage. |
| `IdleConnectionTimeout` | 0 (server) | Duration after which idle connections are closed. 0 uses server-provided value. |
| `MaxConnectionsPerEndpoint` | 10 | Maximum TCP connections per backend endpoint. |
| `ConnectTimeout` | 5s | Timeout for establishing new TCP connections. |

**When to use Direct mode:**
- Low-latency requirements for document operations
- High-throughput scenarios
- Applications running in Azure or with direct network connectivity to Cosmos DB

**Requirements:**
- TCP connectivity to Cosmos DB backend ports (typically 10255)
- Not behind restrictive firewalls that block non-HTTPS traffic

**Note:** Control plane operations (database and container management) always use HTTPS through the gateway, regardless of connection mode. Only document operations (CRUD, queries) use the RNTBD protocol in Direct mode.

## Query Optimization

The SDK includes several query optimization features to reduce latency and improve throughput.

### Query Plan Caching

By default, query plans are cached to reduce Gateway roundtrips for repeated queries. When executing a query that requires a query plan, the SDK:

1. Checks the local cache for a matching plan
2. On cache hit, uses the cached plan (saving a Gateway roundtrip)
3. On cache miss, fetches the plan from Gateway and caches it

**Cache Configuration:**
- **Max Size:** 5,000 entries (LRU eviction when full)
- **TTL:** 5 minutes per entry
- **Cache Key:** Container link + query text (SHA256 hashed)

#### When to Disable Query Plan Caching

You can disable caching on a per-query basis using `DisableQueryPlanCache`:

```go
pk := azcosmos.NewPartitionKeyString("myPartition")
pager := container.NewQueryItemsPager("SELECT * FROM c", pk, &azcosmos.QueryOptions{
    DisableQueryPlanCache: true,
})
```

**Disable query plan caching when:**

| Scenario | Reason |
|----------|--------|
| **Container schema changes in progress** | Avoid stale plans from old indexing policy |
| **Dynamic/ad-hoc queries** | Prevent cache pollution from one-off unique queries |
| **Memory-constrained environment** | Reduce memory footprint |
| **Debugging query performance** | Force fresh plan fetch to see current execution strategy |
| **Very few repeated queries** | Cache overhead may exceed benefit |

#### Important: Use Parameterized Queries

The cache key is based on the **query text only**, not parameter values. To benefit from caching, use parameterized queries:

```go
// ✅ GOOD - Same cache entry reused for different values
query := "SELECT * FROM c WHERE c.city = @city"
pager := container.NewQueryItemsPager(query, pk, &azcosmos.QueryOptions{
    QueryParameters: []azcosmos.QueryParameter{
        {Name: "@city", Value: "Seattle"},
    },
})

// ❌ BAD - Each unique query creates a new cache entry (cache pollution)
query := fmt.Sprintf("SELECT * FROM c WHERE c.city = '%s'", userCity)
pager := container.NewQueryItemsPager(query, pk, nil)
```

Non-parameterized queries with embedded values will cause cache pollution, as each unique query string creates a separate cache entry.

### Optimistic Direct Execution (ODE)

When Direct mode is enabled, simple single-partition queries can bypass the query plan fetch entirely and execute directly against the target partition via RNTBD. This is called **Optimistic Direct Execution (ODE)**.

**ODE is enabled by default** when:
- Direct mode connection is configured
- A partition key is provided with the query

To disable ODE for a specific query:

```go
enabled := false
pager := container.NewQueryItemsPager("SELECT * FROM c", pk, &azcosmos.QueryOptions{
    EnableOptimisticDirectExecution: &enabled,
})
```

**ODE Benefits:**
- Eliminates Gateway query plan fetch latency (~350ms savings observed)
- Reduces RU consumption (no query plan request)
- Lower overall query latency for simple queries

**When ODE may not be suitable:**
- Complex queries with ORDER BY, GROUP BY, aggregates, DISTINCT, TOP, or OFFSET
- Cross-partition queries
- When you need the query plan for debugging

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


