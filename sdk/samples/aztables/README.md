---
page_type: sample
languages:
  - go
products:
  - azure
  - azure-table-storage
urlFragment: tables-samples
---

# Samples for Azure Tables client library for Go

These code samples show common scenario operations with the `aztable` client library.

You can authenticate your client with a Tables API key, credential from the `azidentity` package, or a Shared Access Signature:
* See [sample_authentication.go][sample_authentication] for how to authenticate in the above cases.

These sample programs show common scenarios for the Tables client's offerings.

|**File Name**|**Description**|
|-------------|---------------|
|[sample_create_client.go][create_client]|Instantiate a table client|Authorizing a `ServiceClient` object and `Client` object |
|[sample_create_delete_table.go][create_delete_table]|Creating and deleting a table in a storage account |
|[sample_insert_delete_entities.go][insert_delete_entities]|Inserting and deleting individual entities in a table |
|[sample_query_tables.go][query_tables]|Querying tables in a storage account |
|[sample_update_upsert_merge_entities.go][update_upsert_merge]| Updating, upserting, and merging entities |
|[sample_batching.go][sample_batch]| Committing many requests in a single batch |


### Prerequisites
* Go 1.16 or later is required to use this package.
* You must have an [Azure subscription](https://azure.microsoft.com/free/) and either an
[Azure storage account](https://docs.microsoft.com/azure/storage/common/storage-account-overview) or an [Azure Cosmos Account](https://docs.microsoft.com/azure/cosmos-db/account-overview) to use this package.

## Setup

1. Install the Azure Data Tables client library for Go:
```bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/aztables
```
2. Clone or download this sample repository
3. Open the sample folder in Visual Studio Code or your IDE of choice.

## Running the samples

1. Open a terminal window and `cd` to the directory that the samples are saved in.
2. Set the environment variables specified in the sample file you wish to run.
3. Follow the usage described in the file, e.g. `go run`

## Next steps

Check out the [API reference documentation][api_reference_documentation] to learn more about
what you can do with the Azure Data Tables client library.


<!-- LINKS -->
[api_reference_documentation]: https://docs.microsoft.com/rest/api/storageservices/table-service-rest-api

[sample_authentication]:https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_authentication.go

[create_client]:https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_create_client.go

[create_delete_table]:https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_create_delete_table.go

[insert_delete_entities]: https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_insert_delete_entities.go

[query_entities]: https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_query_table.go

[query_tables]:https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_query_tables.go

[update_upsert_merge]: https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_update_entities.go

[sample_batch]:https://github.com/Azure/azure-sdk-for-go/blob/d90e7e99590c6b7b183b46e0ac69b06ced071158/sdk/samples/aztables/sample_batch.go

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go/sdk/data/aztables/README.png)