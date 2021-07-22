---
page_type: sample
languages:
  - go
products:
  - azure
  - aztable
urlFragment: tables-samples
---

# Samples for Azure Tables client library for Go

These code samples show common scenario operations with the `aztable` client library.

You can authenticate your client with a Tables API key, credential from the `azidentity` package, or a Shared Access Signature:
* See [sample_authentication.go][sample_authentication] for how to authenticate in the above cases.

These sample programs show common scenarios for the Tables client's offerings.

|**File Name**|**Description**|
|-------------|---------------|
|[sample_create_client.go][create_client]|Instantiate a table client|Authorizing a `TableServiceClient` object and `TableClient` object|
|[sample_create_delete_table.go][create_delete_table]|Creating and deleting a table in a storage account|
|[sample_insert_delete_entities.go][insert_delete_entities]|Inserting and deleting individual entities in a table|
|[sample_query_tables.go][query_tables]|Querying tables in a storage account|
|[sample_update_upsert_merge_entities.go][update_upsert_merge]| Updating, upserting, and merging entities|
|[sample_batching.go][sample_batch]| Committing many requests in a single batch|


### Prerequisites
* Go 1.14 or later is required to use this package.
* You must have an [Azure subscription](https://azure.microsoft.com/free/) and either an
[Azure storage account](https://docs.microsoft.com/azure/storage/common/storage-account-overview) or an [Azure Cosmos Account](https://docs.microsoft.com/azure/cosmos-db/account-overview) to use this package.

## Setup

1. Install the Azure Data Tables client library for Go:
```bash
go get github.com/Azure/azure-sdk-for-go/sdk/tables/aztable
```
2. Clone or download this sample repository
3. Open the sample folder in Visual Studio Code or your IDE of choice.

## Running the samples

1. Open a terminal window and `cd` to the directory that the samples are saved in.
2. Set the environment variables specified in the sample file you wish to run.
3. Follow the usage described in the file, e.g. `go run sample_create_table.go`

## Writing Filters

### Supported Comparison Operators
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

### Example Filters

#### Filter on `PartitionKey` and `RowKey`:
```go
pk := "<my_pk>"
rk := "<my_rk>"
queryFilter := fmt.Sprintf("PartitionKey eq '%v' and RowKey eq '%v'", pk, rk)
filter :=
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter on Properties
```go
first := "<first_name>"
last := "<last_name>"
query_filter := fmt.Sprintf("FirstName eq '%v' or LastName eq '%v'", first, last)
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter with string comparison operators
```go
queryFilter := "LastName ge 'A' and LastName lt 'B'"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter with numeric properties
```go
queryFilter := "Age gt 30"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

```go
queryFilter := "AmountDue le 100.25"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter with boolean properties
```go
queryFilter := "IsActive eq true"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter with DateTime properties
```go
queryFilter := "CustomerSince eq datetime'2008-07-10T00:00:00Z'"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```

#### Filter with GUID properties
```go
queryFilter := "GuidValue eq guid'a455c695-df98-5678-aaaa-81d3367e5a34'"
pager := tableClient.Query(&QueryOptions{Filter: &queryFilter})
```


## Next steps

Check out the [API reference documentation][api_reference_documentation] to learn more about
what you can do with the Azure Data Tables client library.


<!-- LINKS -->
[api_reference_documentation]: https://docs.microsoft.com/rest/api/storageservices/table-service-rest-api

[sample_authentication]:https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_authentication.go

[create_client]:https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_create_client.go

[create_delete_table]: https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_create_delete_table.go

[insert_delete_entities]: https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_insert_delete_entities.go

[query_entities]: https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_query_table.go

[query_tables]:https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_query_tables.go

[update_upsert_merge]: https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_update_upsert_merge_entities.go

[sample_batch]: https://github.com/Azure/azure-sdk-for-python/tree/main/sdk/tables/azure-data-tables/samples/sample_batching.go

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-python/sdk/tables/azure-data-tables/README.png)