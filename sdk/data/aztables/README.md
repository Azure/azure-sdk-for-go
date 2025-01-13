# Azure Tables client library for Go

Azure Tables is a NoSQL data storage service that can be accessed from anywhere in the world via authenticated calls using HTTP or HTTPS.
Tables scales as needed to support the amount of data inserted, and allows for the storing of data with non-complex accessing.
The Azure Tables client can be used to access Azure Storage or Cosmos accounts.

[Source code][source_code] | [API reference documentation][Tables_ref_docs]

## Getting started
The Azure Tables SDK can access an Azure Storage or CosmosDB account.

### Prerequisites
* Go versions 1.18 or higher
* You must have an [Azure subscription][azure_subscription] and either
    * an [Azure Storage account][azure_storage_account] or
    * an [Azure Cosmos Account][azure_cosmos_account].

#### Create account
* To create a new Storage account, you can use [Azure Portal][azure_portal_create_account], [Azure PowerShell][azure_powershell_create_account], or [Azure CLI][azure_cli_create_account]:
* To create a new Cosmos storage account, you can use the [Azure CLI][azure_cli_create_cosmos] or [Azure Portal][azure_portal_create_cosmos].

### Install the package
Install the Azure Tables client library for Go with `go get`:
```bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/aztables
```

#### Create the client
The Azure Tables library allows you to interact with two types of resources:
* the tables in your account
* the entities within those tables.
Interaction with these resources starts with an instance of a [client](#clients). To create a client object, you will need the account's table service endpoint URL and a credential that allows you to access the account. The `endpoint` can be found on the page for your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or by running the following Azure CLI command:

```bash
# Log in to Azure CLI first, this opens a browser window
az login
# Get the table service URL for the account
az storage account show -n mystorageaccount -g MyResourceGroup --query "primaryEndpoints.table"
```

Once you have the account URL, it can be used to create the service client:
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }
    client, err := aztables.NewServiceClientWithSharedKey(serviceURL, cred, nil)
    if err != nil {
        panic(err)
    }
}
```

For more information about table service URL's and how to configure custom domain names for Azure Storage check out the [official documentation][azure_portal_account_url]

#### Types of credentials

Both services (Cosmos and Storage) support the the following forms of authentication:
- Microsoft Entra ID token, using one of the collection of types from the [`azidentity`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) module, like [azidentity.DefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential). Example [here](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/aztables#example-NewServiceClient).
- Shared Key Credential
- Connection String
- Shared Access Signature Token

##### Creating the client with a Microsoft Entra ID credential
Use Microsoft Entra ID authentication as the credential parameter to authenticate the client:
```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
    if !ok {
        panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
    }
    serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)

    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        panic(err)
    }
    serviceClient, err := aztables.NewServiceClient(serviceURL, cred, nil)
    if err != nil {
        panic(err)
    }
}
```

##### Creating the client from a shared key
To use an account [shared key][azure_shared_key] (aka account key or access key), provide the key as a string. This can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or by running the following Azure CLI command:

```bash
az storage account keys list -g MyResourceGroup -n MyStorageAccount
```

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }
    serviceClient, err := aztables.NewServiceClientWithSharedKey(serviceURL, cred, nil)
    if err != nil {
        panic(err)
    }
}
```

##### Creating the client from a connection string
Depending on your use case and authorization method, you may prefer to initialize a client instance with a connection string instead of providing the account URL and credential separately. To do this, pass the
connection string to the client's `NewServiceClientFromConnectionString` method. The connection string can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or with the following Azure CLI command:

```bash
az storage account show-connection-string -g MyResourceGroup -n MyStorageAccount
```

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    connStr := "DefaultEndpointsProtocol=https;AccountName=<myAccountName>;AccountKey=<myAccountKey>;EndpointSuffix=core.windows.net"
    serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
    if err != nil {
        panic(err)
    }
}
```

##### Creating the client from a SAS token
To use a [shared access signature (SAS) token][azure_sas_token], provide the token as a string. If your account URL includes the SAS token, omit the credential parameter. You can generate a SAS token from the Azure Portal under [Shared access signature](https://learn.microsoft.com/rest/api/storageservices/create-service-sas) or use the `ServiceClient.GetAccountSASToken` or `Client.GetTableSASToken()` methods.

```golang
import (
    "fmt"
    "time"

    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }
    service, err := aztables.NewServiceClientWithSharedKey("https://<myAccountName>.table.core.windows.net", cred, nil)

    resources := aztables.AccountSASResourceTypes{Service: true}
    permission := aztables.AccountSASPermissions{Read: true}
    start := time.Now()
    expiry := start.AddDate(1, 0, 0)
    sasURL, err := service.GetAccountSASToken(resources, permission, start, expiry)
    if err != nil {
        panic(err)
    }

    serviceURL := fmt.Sprintf("https://<myAccountName>.table.core.windows.net/?%s", sasURL)
    sasService, err := aztables.NewServiceClientWithNoCredential(serviceURL, nil)
    if err != nil {
        panic(err)
    }
}
```

##### Creating the client for Azurite
If you are using the [Azurite](https://github.com/Azure/Azurite) emulator you can authenticate a client with the default connection string:
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    connStr := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"
    svc, err := NewServiceClientFromConnectionString(connStr, nil)
    if err != nil {
        panic(err)
    }

    client, err := svc.CreateTable(context.TODO(), "AzuriteTable", nil)
    if err != nil {
        panic(err)
    }
}
```


## Key concepts
Common uses of the table service include:
* Storing TBs of structured data capable of serving web scale applications
* Storing datasets that do not require complex joins, foreign keys, or stored procedures and can be de-normalized for fast access
* Quickly querying data using a clustered index
* Accessing data using the OData protocol filter expressions

The following components make up the Azure Tables Service:
* The account
* A table within the account, which contains a set of entities
* An entity within a table, as a dictionary

The Azure Tables client library for Go allows you to interact with each of these components through the
use of a dedicated client object.

### Clients
Two different clients are provided to interact with the various components of the Table Service:
1. **`Client`** -
    * Interacts with a specific table (which need not exist yet).
    * Create, delete, query, and upsert entities within the specified table.
    * Create or delete the specified table itself.
2. **`ServiceClient`** -
    * Get and set account settings
    * Query, create, and delete tables within the account.
    * Get a `Client` to access a specific table using the `NewClient` method.

### Entities
Entities are similar to rows. An entity has a **`PartitionKey`**, a **`RowKey`**, and a set of properties. A property is a name value pair, similar to a column. Every entity in a table does not need to have the same properties. Entities are returned as JSON, allowing developers to use JSON marshalling and unmarshalling techniques. Additionally, you can use the `aztables.EDMEntity` to ensure proper round-trip serialization of all properties.
```golang
aztables.EDMEntity{
    Entity: aztables.Entity{
        PartitionKey: "pencils",
        RowKey: "Wooden Pencils",
    },
    Properties: map[string]any{
        "Product": "Ticonderoga Pencils",
        "Price": 5.00,
        "Count": aztables.EDMInt64(12345678901234),
        "ProductGUID": aztables.EDMGUID("some-guid-value"),
        "DateReceived": aztables.EDMDateTime(time.Date{....})
    }
}
```

## Examples

The following sections provide several code snippets covering some of the most common Table tasks, including:

* [Creating a table](#creating-a-table "Creating a table")
* [Creating entities](#creating-entities "Creating entities")
* [Listing entities](#listing-entities "Listing entities")


### Creating a table
Create a table in your account and get a `Client` to perform operations on the newly created table:

```golang
import (
    "context"
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        panic(err)
    }
    accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
    if !ok {
        panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
    }
    serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)

    service, err := aztables.NewServiceClient(serviceURL, cred, nil)
    if err != nil {
        panic(err)
    }

    // Create a table
    _, err = service.CreateTable(context.TODO(), "fromServiceClient", nil)
    if err != nil {
        panic(err)
    }
}
```

### Creating entities
Create entities in the table:

```go
import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "github.com/Azure/azure-sdk-for-go/sdk/azcore"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }

    service, err := aztables.NewServiceClient("https://<myAccountName>.table.core.windows.net", cred, nil)
    if err != nil {
        panic(err)
    }

    client, err := service.NewClient("myTable")
    if err != nil {
        panic(err)
    }

    myEntity := aztables.EDMEntity{
        Entity: aztables.Entity{
            PartitionKey: "001234",
            RowKey: "RedMarker",
        },
        Properties: map[string]any{
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
    if err != nil {
        panic(err)
    }

    resp, err := client.AddEntity(context.TODO(), marshalled, nil)
    if err != nil {
        panic(err)
    }
}
```

### Listing entities
List entities in the table:

```go
import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "github.com/Azure/azure-sdk-for-go/sdk/azcore"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }
    client, err := aztables.NewClient("https://myAccountName.table.core.windows.net/myTable", cred, nil)
    if err != nil {
        panic(err)
    }

    filter := "PartitionKey eq 'markers' or RowKey eq 'Markers'"
    options := &aztables.ListEntitiesOptions{
        Filter: &filter,
        Select: to.Ptr("RowKey,Value,Product,Available"),
        Top: to.Ptr(int32(15)),
    }

    pager := client.NewListEntitiesPager(options)
    pageCount := 0
    for pager.More() {
        response, err := pager.NextPage(context.TODO())
        if err != nil {
            panic(err)
        }
        fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
        pageCount += 1

        for _, entity := range response.Entities {
            var myEntity aztables.EDMEntity
            err = json.Unmarshal(entity, &myEntity)
            if err != nil {
                panic(err)
            }

            fmt.Printf("Received: %v, %v, %v, %v\n", myEntity.RowKey, myEntity.Properties["Value"], myEntity.Properties["Product"], myEntity.Properties["Available"])
        }
    }
}
```

#### Writing Filters

##### Supported Comparison Operators
|**Operator**|**URI expression**|
|------------|------------------|
|`Equal`|`eq`|
|`GreaterThan`|`gt`|
|`GreaterThanOrEqual`|`ge`|
|`LessThan`|`lt`|
|`LessThanOrEqual`|`le`|
|`NotEqual`|`ne`|
|`And`|`and`|
|`Not`|`not`|
|`Or`|`or`|

Query strings must wrap literal values in single quotes. Literal values containing single quote characters must be escaped with a double single quote. To search for a `LastName` property of "O'Connor" use the following syntax
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("LastName eq 'O''Connor'"),
}
```

##### String Properties
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("LastName ge 'A' and LastName lt 'B'"),
}
```

##### Numeric Properties
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("Age gt 30"),
}

options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("AmountDue le 100.25"),
}
```

##### Boolean Properties
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("IsActive eq true"),
}
```

##### Datetime Properties
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("CustomerSince eq datetime'2008-07-10T00:00:00Z'"),
}
```

##### GUID Properties
```go
options := &aztables.ListEntitiesOptions{
    Filter: to.Ptr("GuidValue eq guid'a455c695-df98-5678-aaaa-81d3367e5a34'"),
}
```

#### Using Continuation Tokens
The pager exposes continuation tokens that can be used by a new pager instance to begin listing entities from a specific point. For example:
```go
import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "github.com/Azure/azure-sdk-for-go/sdk/azcore"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)
func main() {
    cred, err := aztables.NewSharedKeyCredential("<myAccountName>", "<myAccountKey>")
    if err != nil {
        panic(err)
    }
    client, err := aztables.NewClient("https://myAccountName.table.core.windows.net/myTable", cred, nil)
    if err != nil {
        panic(err)
    }

    pager := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{Top: to.Ptr(int32(10))})
    count := 0
    for pager.More() {
        response, err := pager.NextPage(context.TODO())
        if err != nil {
            panic(err)
        }

        count += len(response.Entities)

        if count > 20 {
            break
        }
    }

    newPager := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
        Top:          to.Ptr(int32(10)),
        PartitionKey: pager.NextPagePartitionKey(),
        RowKey:       pager.NextPageRowKey(),
    })

    for newPager.More() {
        // begin paging where 'pager' left off
    }
}
```

## Troubleshooting

### Error Handling

All I/O operations will return an `error` that can be investigated to discover more information about the error. In addition, you can investigate the raw response of any response object:
```golang
resp, err := client.CreateTable(context.TODO(), nil)
if err != nil {
    var respErr azcore.ResponseError
    if errors.As(err, &respErr) {
        // handle err ...
    }
}
```

### Logging

This module uses the classification based logging implementation in azcore. To turn on logging set `AZURE_SDK_GO_LOGGING` to `all`. If you only want to include logs for `aztables`, you must create your own logger and set the log classification as `LogCredential`.

To obtain more detailed logging, including request/response bodies and header values, make sure to leave the logger as default or enable the `LogRequest` and/or `LogResponse` classificatons. A logger that only includes credential logs can be like the following:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
// Set log to output to the console
log.SetListener(func(cls log.Classification, msg string) {
    fmt.Println(msg) // printing log out to the console
})

// Includes only requests and responses in credential logs
log.SetClassifications(log.Request, log.Response)
```

> CAUTION: logs from credentials contain sensitive information.
> These logs must be protected to avoid compromising account security.

## Next steps

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Azure.Tables` label.

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

### Additional documentation
For more extensive documentation on Azure Tables, see the [Azure Tables documentation][Tables_product_doc] on learn.microsoft.com.

## Known Issues
A list of currently known issues relating to Cosmos DB table endpoints can be found [here](https://aka.ms/tablesknownissues).

## Contributing
This project welcomes contributions and suggestions.  Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][msft_oss_coc]. For more information see the [Code of Conduct FAQ][msft_oss_coc_faq] or contact [opencode@microsoft.com][contact_msft_oss] with any additional questions or comments.

<!-- LINKS -->
[source_code]:https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/data/aztables
[Tables_ref_docs]:https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/aztables
[Tables_product_doc]:https://learn.microsoft.com/azure/cosmos-db/table-introduction

[azure_subscription]:https://azure.microsoft.com/free/
[azure_storage_account]:https://learn.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-portal
[azure_cosmos_account]:https://learn.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[pip_link]:https://pypi.org/project/pip/

[azure_create_cosmos]:https://learn.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[azure_cli_create_cosmos]:https://learn.microsoft.com/azure/cosmos-db/scripts/cli/table/create
[azure_portal_create_cosmos]:https://learn.microsoft.com/azure/cosmos-db/create-cosmosdb-resources-portal
[azure_portal_create_account]:https://learn.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-portal
[azure_powershell_create_account]:https://learn.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-powershell
[azure_cli_create_account]: https://learn.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-cli

[azure_cli_account_url]:https://learn.microsoft.com/cli/azure/storage/account?view=azure-cli-latest#az-storage-account-show
[azure_powershell_account_url]:https://learn.microsoft.com/powershell/module/az.storage/get-azstorageaccount?view=azps-4.6.1
[azure_portal_account_url]:https://learn.microsoft.com/azure/storage/common/storage-account-overview#storage-account-endpoints

[azure_sas_token]:https://learn.microsoft.com/azure/storage/common/storage-sas-overview
[azure_shared_key]:https://learn.microsoft.com/rest/api/storageservices/authorize-with-shared-key

[azure_core_ref_docs]:https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore
[azure_core_readme]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore/README.md

[tables_error_codes]: https://learn.microsoft.com/rest/api/storageservices/table-service-error-codes

[msft_oss_coc]:https://opensource.microsoft.com/codeofconduct/
[msft_oss_coc_faq]:https://opensource.microsoft.com/codeofconduct/faq/
[contact_msft_oss]:mailto:opencode@microsoft.com

[tables_rest]: https://learn.microsoft.com/rest/api/storageservices/table-service-rest-api

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go/sdk/data/aztables/README.png)
