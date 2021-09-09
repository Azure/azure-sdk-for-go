# Azure Storage Blob SDK for Go

## Introduction 
The Microsoft Azure Storage SDK for Go allows you to build applications that takes advantage of Azure's scalable cloud
storage. This SDK replaces the previously previewed [azblob package](azblob_track1dot5) in the
separate repository.

## Prerequisites
* Go versions 1.16 or higher. If you don't already have it, install [the Go distribution](https://golang.org/dl/)
* You must have an [Azure storage account][azure_storage_account] 

## Getting Started
* Install the Azure blob storage client library for Go with `go get`:
  ```bash
  go get github.com/Azure/azure-sdk-for-go/sdk/storage/azblob
  ```
* Import SDK in your code:
  ```golang
  import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
  ```

## Code Samples

`azblob` allows you to interact with three types of resources
* The storage accounts.
* The containers within those storage accounts.
* The blobs (block blob/ page blob/ append blob) within those containers. 

To interact with these resources, start by creating the instances of each type

### Types of credentials
The clients support different forms of authentication. Cosmos accounts can use a Shared Key Credential, Connection String, or an Shared Access Signature Token for authentication. Storage account can use the same credentials as a Cosmos account and can use the credentials in [`azidentity`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) like `azidentity.NewDefaultAzureCredential()`.

The azblob package supports any of the types that implement the `azcore.TokenCredential` interface, authorization via a Connection String, or authorization with a Shared Access Signature Token.

#### 1. Creating the client from a shared key
To use an account [shared key][azure_shared_key] (aka account key or access key), provide the key as a string. This can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or by running the following Azure CLI command:

```bash
az storage account keys list -g MyResourceGroup -n MyStorageAccount
```

Use Azure Active Directory (AAD) authentication as the credential parameter to authenticate the client:
```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
handle(err)
serviceClient, err := azblob.NewServiceClient("https://<myAccountName>.blob.core.windows.net/", cred, nil)
handle(err)
```

##### 2. Creating the client from a connection string
Depending on your use case and authorization method, you may prefer to initialize a client instance with a connection string instead of providing the account URL and credential separately. To do this, pass the
connection string to the client's `NewServiceClientFromConnectionString` method. The connection string can be found in your storage account in the [Azure Portal][azure_portal_account_url] under the "Access Keys" section or with the following Azure CLI command:

```bash
az storage account show-connection-string -g MyResourceGroup -n MyStorageAccount
```

```golang
connStr := "DefaultEndpointsProtocol=https;AccountName=<myAccountName>;AccountKey=<myAccountKey>;EndpointSuffix=core.windows.net"
serviceClient, err := azblob.NewServiceClientFromConnectionString(connStr, nil)
```

##### Creating the client from a SAS token
To use a [shared access signature (SAS) token][azure_sas_token], provide the token as a string. If your account URL includes the SAS token, omit the credential parameter. You can generate a SAS token from the Azure Portal under [Shared access signature](https://docs.microsoft.com/rest/api/storageservices/create-service-sas) or use the `ServiceClient.GetAccountSASToken` or `ContainerClient.GetContainerSASToken()` methods.

```golang
cred, err := azblob.NewSharedKeyCredential("myAccountName", "myAccountKey")
handle(err)
service, err := azblob.NewServiceClient("https://<myAccountName>.blob.core.windows.net", cred, nil)

resources := azblob.AccountSASResourceTypes{Service: true}
permission := azblob.AccountSASPermissions{Read: true}
start := time.Now()
expiry := start.AddDate(1, 0, 0)
sasUrl, err := service.GetAccountSASToken(resources, permission, start, expiry)
handle(err)

sasService, err := azblob.NewServiceClient(sasUrl, azcore.AnonymousCredential(), nil)
handle(err)
```

- For more detailed examples, please refer to zt_examples_test.go.

### Clients
Three different clients are provided to interact with the various components of the Blob Service:

1. **`ServiceClient`** 
   * Get and set account settings.
   * Query, create, and delete containers within the account.

2. **`ContainerClient`** 
   * Get and set container access settings, properties, and metadata.
   * Create, delete, and query blobs within the container.
   * `ContainerLeaseClient` to support container lease management. 

3. **`BlobClient`**
   * `AppendBlobClient`, `BlockBlobClient`, and `PageBlobClient`
   * Get and set blob properties.
   * Perform CRUD operations a given blob.
   * `BlobLeaseClient` to support blob lease management.

## Troubleshooting

### Error Handling

All I/O operations will return an `error` that can be investigated to discover more information about the error. In addition, you can investigate the raw response of any response object:
```golang
resp, err := serviceClient.CreateContainer(context.Background(), "testcontainername", nil)
if err != nil {
    err = errors.As(err, azcore.HTTPResponse)
    // handle err ...
}
```

### Logging

This module uses the classification based logging implementation in azcore. 
To turn on logging set `AZURE_SDK_GO_LOGGING` to `all`. 

If you only want to include logs for `azblob`, you must create your own logger and set the log classification as `LogCredential`.

To obtain more detailed logging, including request/response bodies and header values, make sure to leave the logger as default or enable the `LogRequest` and/or `LogResponse` classificatons. A logger that only includes credential logs can be like the following:

```golang
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
>

### Additional documentation
For more extensive documentation on Azure `azblob` SDK, see the [Azure blobs documentation][Blobs_product_doc] on docs.microsoft.com.

## License

This project is licensed under MIT.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Azure.AzBlob` label.

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
[azblob_track1dot5]:https://github.com/azure/azure-storage-blob-go
[source_code]:https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/storage/azblob
[Blobs_product_doc]:https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blobs-introduction

[azure_subscription]:https://azure.microsoft.com/free/
[azure_storage_account]:https://docs.microsoft.com/azure/storage/common/storage-account-create?tabs=azure-portal

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

[blobs_error_codes]: https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-error-codes

[msft_oss_coc]:https://opensource.microsoft.com/codeofconduct/
[msft_oss_coc_faq]:https://opensource.microsoft.com/codeofconduct/faq/
[contact_msft_oss]:mailto:opencode@microsoft.com

[blobs_rest]: https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-rest-api