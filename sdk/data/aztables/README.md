# Azure Data Tables client library for Go

Azure Data Tables is a NoSQL data storage service that can be accessed from anywhere in the world via authenticated calls using HTTP or HTTPS.
Tables scales as needed to support the amount of data inserted, and allow for the storing of data with non-complex accessing.
The Azure Data Tables client can be used to access Azure Storage or Cosmos accounts.

[Source code][source_code] | [Package][Tables_gopackage] | [API reference documentation][Tables_ref_docs] | [Samples][Tables_samples]

## Getting started
The Azure Data Tables SDK can access an Azure Storage or CosmosDB account.

### Prerequisites
* Go versions 1.16 or higher
* You must have an [Azure subscription][azure_subscription] and either
    * an [Azure Storage account][azure_storage_account] or
    * an [Azure Cosmos Account][azure_cosmos_account].

#### Create account
* To create a new storage account, you can use [Azure Portal][azure_portal_create_account], [Azure PowerShell][azure_powershell_create_account], or [Azure CLI][azure_cli_create_account]:
* To create a new cosmos storage account, you can use the [Azure CLI][azure_cli_create_cosmos] or [Azure Portal][azure_portal_create_cosmos].

### Install the package
Install the Azure Data Tables client library for Go with `go get`:
```bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/aztables
```

#### Create the client
The Azure Data Tables library allows you to interact with two types of resources:
* the tables in your account
* the entities within those tables.
Interaction with these resources starts with an instance of a [client](#clients). To create a client object, you will need the account's table service endpoint URL and a credential that allows you to access the account. The `endpoint` can be found on the page for your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or by running the following Azure CLI command:

```bash
# Get the table service URL for the account
az storage account show -n mystorageaccount -g MyResourceGroup --query "primaryEndpoints.table"
```

Once you have the account URL, it can be used to create the service client:
```golang
cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
serviceClient, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net/", cred, nil)
handle(err)
```

For more information about table service URL's and how to configure custom domain names for Azure Storage check out the [official documentation][azure_portal_account_url]

#### Types of credentials
The clients support different forms of authentication. Cosmos accounts can use a Shared Key Credential, Connection String, or an Shared Access Signature Token for authentication. Storage account can use the same credentials as a Cosmos account and can use the credentials in [`azidentity`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) like `azidentity.NewDefaultAzureCredential()`.

The aztables library supports any of the `azcore.TokenCredential` interfaces, authorization via a Connection String, or authorization with a Shared Access Signature Token.

##### Creating the client from a shared key
To use an account [shared key][azure_shared_key] (aka account key or access key), provide the key as a string. This can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or by running the following Azure CLI command:

```bash
az storage account keys list -g MyResourceGroup -n MyStorageAccount
```

Use AAD authentication as the credential parameter to authenticate the client:
```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
handle(err)
serviceClient, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net/", cred, nil)
handle(err)
```

##### Creating the client from a connection string
Depending on your use case and authorization method, you may prefer to initialize a client instance with a connection string instead of providing the account URL and credential separately. To do this, pass the
connection string to the client's `from_connection_string` class method. The connection string can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or with the following Azure CLI command:

```bash
az storage account show-connection-string -g MyResourceGroup -n MyStorageAccount
```

```golang
connStr := "DefaultEndpointsProtocol=https;AccountName=<my_account_name>;AccountKey=<my_account_key>;EndpointSuffix=core.windows.net"
serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
```

##### Creating the client from a SAS token
To use a [shared access signature (SAS) token][azure_sas_token], provide the token as a string. If your account URL includes the SAS token, omit the credential parameter. You can generate a SAS token from the Azure Portal under [Shared access signature](https://docs.microsoft.com/rest/api/storageservices/create-service-sas) or use the `ServiceClient.GetAccountSASToken` or `Client.GetTableSASToken()` functions.

```golang
cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)

resources := aztables.AccountSASResourceTypes{Service: true}
permission := aztables.AccountSASPermissions{Read: true}
start := time.Date(2021, time.August, 21, 1, 1, 0, 0, time.UTC)
expiry := time.Date(2022, time.August, 21, 1, 1, 0, 0, time.UTC)
sasUrl, err := service.GetAccountSASToken(resources, permission, start, expiry)
handle(err)

sasService, err := aztables.NewServiceClient(sasUrl, azcore.AnonymousCredential(), nil)
handle(err)
```


## Key concepts
Common uses of the Table service included:
* Storing TBs of structured data capable of serving web scale applications
* Storing datasets that do not require complex joins, foreign keys, or stored procedures and can be de-normalized for fast access
* Quickly querying data using a clustered index
* Accessing data using the OData protocol and LINQ filter expressions

The following components make up the Azure Data Tables Service:
* The account
* A table within the account, which contains a set of entities
* An entity within a table, as a dictionary

The Azure Data Tables client library for Go allows you to interact with each of these components through the
use of a dedicated client object.

### Clients
Two different clients are provided to interact with the various components of the Table Service:
1. **`ServiceClient`** -
    * Get and set account setting
    * Query, create, and delete tables within the account.
    * Get a `Client` to access a specific table using the `NewClient` method.
2. **`Client`** -
    * Interacts with a specific table (which need not exist yet).
    * Create, delete, query, and upsert entities within the specified table.
    * Create or delete the specified table itself.

### Entities
Entities are similar to rows. An entity has a **`PartitionKey`**, a **`RowKey`**, and a set of properties. A property is a name value pair, similar to a column. Every entity in a table does not need to have the same properties. Entities are returned as JSON, allowing developers to use JSON marshalling and unmarshalling techniques. Additionally, you can use the `aztables.EDMEntity` to ensure proper round-trip serialization of all properties.
```golang
aztables.EDMEntity{
    Entity: aztables.Entity{
        PartitionKey: "pencils",
        RowKey: "id-003",
    },
    Properties: map[string]interface{}{
        "Product": "Ticonderoga Pencils",
        "Price": 5.00,
        "Count": aztables.EDMInt64(12345678901234),
        "ProductGUID": aztables.EDMGUID("some-guid-value"),
        "DateReceived": aztables.EDMDateTime(time.Date{....}),
        "ProductCode": aztables.EDMBinary([]byte{"somebinaryvalue"})
    }
}
```
<!-- * **[create_entity][create_entity]** - Add an entity to the table.
* **[delete_entity][delete_entity]** - Delete an entity from the table.
* **[update_entity][update_entity]** - Update an entity's information by either merging or replacing the existing entity.
    * `UpdateMode.MERGE` will add new properties to an existing entity it will not delete an existing properties
    * `UpdateMode.REPLACE` will replace the existing entity with the given one, deleting any existing properties not included in the submitted entity
* **[query_entities][query_entities]** - Query existing entities in a table using OData filters
* **[get_entity][get_entity]** - Get a specific entity from a table by partition and row key.
* **[upsert_entity][upsert_entity]** - Merge or replace an entity in a table, or if the entity does not exist, inserts the entity.
    * `UpdateMode.MERGE` will add new properties to an existing entity it will not delete an existing properties
    * `UpdateMode.REPLACE` will replace the existing entity with the given one, deleting any existing properties not included in the submitted entity -->

## Examples

The following sections provide several code snippets covering some of the most common Table tasks, including:

* [Creating a table](#creating-a-table "Creating a table")
* [Creating entities](#creating-entities "Creating entities")
* [Querying entities](#querying-entities "Querying entities")


### Creating a table
Create a table in your account and get a `Client` to perform operations on the newly created table:

```golang
cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)
handle(err)
resp, err := service.CreateTable("myTable")
```

### Creating entities
Create entities in the table:

```golang
cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
service, err := aztables.NewServiceClient("https://<my_account_name>.table.core.windows.net", cred, nil)
handle(err)

myEntity := aztables.EDMEntity{
    Entity: aztables.Entity{
        PartitionKey: "001234",
        RowKey: "RedMarker",
    },
    Properties: map[string]interface{}{
        "Stock": 15,
        "Price": 9.99,
        "Comments": "great product",
        "OnSale": true,
        "ReducedPrice": 7.99,
        "PurchaseDate": aztables.EDMDateTime(time.Date(2021, time.August, 21, 1, 1, 0, 0, time.UTC)),
        "BinaryRepresentation": aztables.EDMBinary([]byte{"Bytesliceinfo"})
    }
}
marshalled, err := json.Marshal(myEntity)
handle(err)

client, err := service.NewClient("myTable")
handle(err)

resp, err := client.AddEntity(context.Background(), marshalled, nil)
handle(err)
```

### Querying entities
Querying entities in the table:

```golang
cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
client, err := aztables.NewClient("https://myAccountName.table.core.windows.net/myTableName", cred, nil)
handle(err)

filter := "PartitionKey eq 'markers' or RowKey eq 'id-001'"
options := &ListEntitiesOptions{
    Filter: &filter,
    Select: to.StringPtr("RowKey,Value,Product,Available"),
    Top: to.Int32Ptr(15),
}

pager := client.List(options)
for pager.NextPage(context.Background()) {
    resp := pager.PageResponse()
    fmt.Printf("Received: %v entities\n", len(resp.Entities))

    for _, entity := range resp.Entities {
        var myEntity aztables.EDMEntity
        err = json.Unmarshal(entity, &myEntity)
        handle(err)

        fmt.Printf("Received: %v, %v, %v, %v\n", myEntity.Properties["RowKey"], myEntity.Properties["Value"], myEntity.Properties["Product"], myEntity.Properties["Available"])
    }
}

err := pager.Err()
handle(err)
```

## Troubleshooting

### Error Handling

All I/O operations will return an `error` that can be investigated to discover more information about the error. In addition, you can investigate the raw response of any response object:
```golang
resp, err := client.CreateTable(context.Background(), nil)
if err != nil {
    err = errors.As(err, azcore.HTTPResponse)
    // handle err ...
}
```

### Logging

This module uses the classification based logging implementation in azcore. To turn on logging set `AZURE_SDK_GO_LOGGING` to `all`. If you only want to include logs for `aztables`, you must create your own logger and set the log classification as `LogCredential`.

To obtain more detailed logging, including request/response bodies and header values, make sure to leave the logger as default or enable the `LogRequest` and/or `LogResponse` classificatons. A logger that only includes credential logs can be like the following:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
// Set log to output to the console
azlog.SetListener(func(cls LogClassification, s string) {
		fmt.Println(s) // printing log out to the console
  })

// Include only azidentity credential logs
azlog.SetClassifications(azidentity.LogCredential)
```

> CAUTION: logs from credentials contain sensitive information.
> These logs must be protected to avoid compromising account security.

## Next steps

Get started with our [Table samples][tables_samples].

Several Azure Data Tables Go SDK samples are available to you in the SDK's GitHub repository. These samples provide example code for additional scenarios commonly encountered while working with Tables.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Azure.Identity` label.

## Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.


### Common Scenarios
These code samples show common scenario operations with the Azure Data tables client library.

* Create and delete tables: [sample_create_delete_table.py](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_create_delete_table.py) ([async version](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/async_samples/sample_create_delete_table_async.py))
* List and query tables: [sample_query_tables.py](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_query_tables.py) ([async version](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/async_samples/sample_query_tables_async.py))
* Insert and delete entities: [sample_insert_delete_entities.py](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_insert_delete_entities.py) ([async version](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/async_samples/sample_insert_delete_entities_async.py))
* Query and list entities: [sample_query_table.py](https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_query_table.py) ([async version](https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/async_samples/sample_query_table_async.py))
* Update, upsert, and merge entities: [sample_update_upsert_merge_entities.py](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_update_upsert_merge_entities.py) ([async version](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/async_samples/sample_update_upsert_merge_entities_async.py))
* Committing many requests in a single transaction: [sample_batching.py](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_batching.py) ([async version](https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/async_samples/sample_batching_async.py))

### Additional documentation
For more extensive documentation on Azure Data Tables, see the [Azure Data Tables documentation][Tables_product_doc] on docs.microsoft.com.

## Known Issues
A list of currently known issues relating to Cosmos DB table endpoints can be found [here](https://aka.ms/tablesknownissues).

## Contributing
This project welcomes contributions and suggestions.  Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][msft_oss_coc]. For more information see the [Code of Conduct FAQ][msft_oss_coc_faq] or contact [opencode@microsoft.com][contact_msft_oss] with any additional questions or comments.

<!-- LINKS -->
[source_code]:https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/data/aztables
[Tables_gopackage]:https://aka.ms/azsdk/go/tables
[Tables_ref_docs]:https://aka.ms/azsdk/go/tables
<!-- [Tables_ref_docs]:https://docs.microsoft.com/python/api/overview/azure/data-tables-readme-pre?view=azure-python-preview -->
[Tables_product_doc]:https://docs.microsoft.com/azure/cosmos-db/table-introduction
[Tables_samples]:https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples
[migration_guide]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/migration_guide.md

[azure_subscription]:https://azure.microsoft.com/free/
[azure_storage_account]:https://docs.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-portal
[azure_cosmos_account]:https://docs.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[pip_link]:https://pypi.org/project/pip/

[azure_create_cosmos]:https://docs.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[azure_cli_create_cosmos]:https://docs.microsoft.com/azure/cosmos-db/scripts/cli/table/create
[azure_portal_create_cosmos]:https://docs.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[azure_portal_create_account]:https://docs.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-portal
[azure_powershell_create_account]:https://docs.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-powershell
[azure_cli_create_account]: https://docs.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-cli

[azure_cli_account_url]:https://docs.microsoft.com/cli/azure/storage/account?view=azure-cli-latest#az-storage-account-show
[azure_powershell_account_url]:https://docs.microsoft.com/powershell/module/az.storage/get-azstorageaccount?view=azps-4.6.1
[azure_portal_account_url]:https://docs.microsoft.com/azure/storage/common/storage-account-overview#storage-account-endpoints

[azure_sas_token]:https://docs.microsoft.com/azure/storage/common/storage-sas-overview
[azure_shared_key]:https://docs.microsoft.com/rest/api/storageservices/authorize-with-shared-key

[azure_core_ref_docs]:https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore
[azure_core_readme]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore/README.md

[tables_error_codes]: https://docs.microsoft.com/rest/api/storageservices/table-service-error-codes

[msft_oss_coc]:https://opensource.microsoft.com/codeofconduct/
[msft_oss_coc_faq]:https://opensource.microsoft.com/codeofconduct/faq/
[contact_msft_oss]:mailto:opencode@microsoft.com

[tables_rest]: https://docs.microsoft.com/rest/api/storageservices/table-service-rest-api

<!-- [create_entity]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_insert_delete_entities.py#L51-L57
[delete_entity]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_insert_delete_entities.py#L73-L80
[update_entity]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_update_upsert_merge_entities.py#L128-L129
[query_entities]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_query_table.py#L63-L72
[get_entity]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_update_upsert_merge_entities.py#L52-L55
[upsert_entity]:https://github.com/Azure/azure-sdk-for-python/blob/main/sdk/tables/azure-data-tables/samples/sample_update_upsert_merge_entities.py#L103-L120 -->

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-python/sdk/tables/azure-data-tables/README.png)