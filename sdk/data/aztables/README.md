# Azure Tables client library for Go

Azure Table storage is a service that stores large amounts of structured NoSQL data in the cloud, providing
a key/attribute store with a schema-less design.

Azure Cosmos DB provides a Table API for applications that are written for Azure Table storage that need premium capabilities like:

- Turnkey global distribution.
- Dedicated throughput worldwide.
- Single-digit millisecond latencies at the 99th percentile.
- Guaranteed high availability.
- Automatic secondary indexing.

The Azure Tables client library can seamlessly target either Azure Table storage or Azure Cosmos DB table service endpoints with no code changes.

## Getting started

### Install the package
Install the Azure Tables client library for Go :

### Prerequisites
* An [Azure subscription][azure_sub].
* An existing Azure storage account or Azure Cosmos DB database with Azure Table API specified.

If you need to create either of these, you can use the [Azure CLI][azure_cli].

#### Creating a storage account

Create a storage account `mystorageaccount` in resource group `MyResourceGroup`
in the subscription `MySubscription` in the West US region.


#### Creating a Cosmos DB

Create a Cosmos DB account `MyCosmosDBDatabaseAccount` in resource group `MyResourceGroup`
in the subscription `MySubscription` and a table named `MyTableName` in the account.


### Authenticate the Client

Learn more about options for authentication _(including Connection Strings, Shared Key, and Shared Key Signatures)_ in our samples.

## Key concepts

- `TableServiceClient` - Client that provides methods to interact at the Table Service level such as creating, listing, and deleting tables
- `TableClient` - Client that provides methods to interact at an table entity level such as creating, querying, and deleting entities within a table.
- `Table` - Tables store data as collections of entities.
- `Entity` - Entities are similar to rows. An entity has a primary key and a set of properties. A property is a name value pair, similar to a column.

Common uses of the Table service include:

- Storing TBs of structured data capable of serving web scale applications
- Storing datasets that don't require complex joins, foreign keys, or stored procedures and can be de-normalized for fast access
- Quickly querying data using a clustered index
### Create the Table service client

First, we need to construct a `TableServiceClient`.

### Create an Azure table
Next, we can create a new table.


### Get an Azure table
The set of existing Azure tables can be queries using an OData filter.


### Delete an Azure table

Individual tables can be deleted from the service.


### Create the Table client

To interact with table entities, we must first construct a `TableClient`.


### Add table entities

Let's define a new `TableEntity` so that we can add it to the table.

Using the `TableClient` we can now add our new entity to the table.


### Query table entities

To inspect the set of existing table entities, we can query the table using an OData filter.


If you prefer LINQ style query expressions, we can query the table using that syntax as well.


### Delete table entities

If we no longer need our new table entity, it can be deleted.


## Troubleshooting

When you interact with the Azure table library using the .NET SDK, errors returned by the service correspond to the same HTTP
status codes returned for [REST API][tables_rest] requests.

For example, if you try to create a table that already exists, a `409` error is returned, indicating "Conflict".


### Setting up console logging

The simplest way to see the logs is to enable the console logging.
To create an Azure SDK log listener that outputs messages to console use AzureEventSourceListener.CreateConsoleLogger method.


To learn more about other logging mechanisms see [here][logging].

## Next steps

Get started with our [Table samples][table_client_samples].

## Known Issues

A list of currently known issues relating to Cosmos DB table endpoints can be found [here](https://aka.ms/tablesknownissues).

## Contributing

This project welcomes contributions and suggestions.  Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution. For
details, visit [cla.microsoft.com][cla].

This project has adopted the [Microsoft Open Source Code of Conduct][coc].
For more information see the [Code of Conduct FAQ][coc_faq] or contact
[opencode@microsoft.com][coc_contact] with any additional questions or comments.

## Generating the client

From the tables dir:

autorest --use=@autorest/go@4.0.0-preview.20 https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification/cosmos-db/data-plane/readme.md --tag=package-2019-02 --file-prefix="zz_generated_" --modelerfour.lenient-model-deduplication --license-header=MICROSOFT_MIT_NO_VERSION --output-folder=aztables --module=aztables --openapi-type="data-plane" --credential-scope=none

<!-- LINKS -->
[tables_rest]: https://docs.microsoft.com/rest/api/storageservices/table-service-rest-api
[azure_cli]: https://docs.microsoft.com/cli/azure
[azure_sub]: https://azure.microsoft.com/free/
[table_client_nuget_package]: https://www.nuget.org/packages?q=Azure.Data.Tables
[table_client_samples]: https://github\.com/Azure/azure-sdk-for-go
[table_client_src]: https://github\.com/Azure/azure-sdk-for-go
[logging]: https://github\.com/Azure/azure-sdk-for-go
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-net%2Fsdk%2Ftables%2FAzure.Data.Tables%2FREADME.png)
