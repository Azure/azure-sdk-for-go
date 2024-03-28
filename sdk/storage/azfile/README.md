# Azure File Storage SDK for Go
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/azfile)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azfile)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/go%20-%20azdatalake%20-%20ci?branchName=main)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=4784&branchName=main)
[![Code Coverage](https://img.shields.io/azure-devops/coverage/azure-sdk/public/4784/main)](https://img.shields.io/azure-devops/coverage/azure-sdk/public/4784/main)

> Service Version: 2023-11-03

Azure File Shares offers fully managed file shares in the cloud that are accessible via the industry standard 
[Server Message Block (SMB) protocol](https://docs.microsoft.com/windows/desktop/FileIO/microsoft-smb-protocol-and-cifs-protocol-overview). 
Azure file shares can be mounted concurrently by cloud or on-premises deployments of Windows, Linux, and macOS. 
Additionally, Azure file shares can be cached on Windows Servers with Azure File Sync for fast access near where the data is being used.

[Source code][source] | [API reference documentation][docs] | [REST API documentation][rest_docs] | [Product documentation][product_docs]

## Getting started

### Install the package

Install the Azure File Storage SDK for Go with [go get][goget]:

```Powershell
go get github.com/Azure/azure-sdk-for-go/sdk/storage/azfile
```

If you plan to authenticate with Azure Active Directory (recommended), also install the [azidentity][azidentity] module.

```Powershell
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

### Prerequisites

A supported [Go][godevdl] version (the Azure SDK supports the two most recent Go releases).

You need an [Azure subscription][azure_sub] and a
[Storage Account][storage_account_docs] to use this package.

To create a new Storage Account, you can use the [Azure Portal][storage_account_create_portal],
[Azure PowerShell][storage_account_create_ps], or the [Azure CLI][storage_account_create_cli].
Here's an example using the Azure CLI:

```Powershell
az storage account create --name MyStorageAccount --resource-group MyResourceGroup --location westus --sku Standard_LRS
```

### Authenticate the client

The Azure File Storage SDK for Go allows you to interact with four types of resources: the storage
account itself, file shares, directories, and files. Interaction with these resources starts with an instance of a
client. To create a client object, you will need the storage account's file service URL and a
credential that allows you to access the storage account. The [azidentity][azidentity] module makes it easy to add 
Azure Active Directory support for authenticating Azure SDK clients with their corresponding Azure services.

```go
// create a credential for authenticating with Azure Active Directory
cred, err := azidentity.NewDefaultAzureCredential(nil)
// TODO: handle err

// create service.Client for the specified storage account that uses the above credential
client, err := service.NewClient("https://<my-storage-account-name>.file.core.windows.net/", cred, &service.ClientOptions{FileRequestIntent: to.Ptr(service.ShareTokenIntentBackup)})
// TODO: handle err
```

Learn more about enabling Azure Active Directory for authentication with Azure Storage: [Authorize access to blobs using Azure Active Directory][storage_ad]

Other options for authentication include connection strings, shared key, and shared access signatures (SAS). 
Use the appropriate client constructor function for the authentication mechanism you wish to use.

## Key concepts

Azure file shares can be used to:

- Completely replace or supplement traditional on-premises file servers or NAS devices.
- "Lift and shift" applications to the cloud that expect a file share to store file application or user data.
- Simplify new cloud development projects with shared application settings, diagnostic shares, and Dev/Test/Debug tool file shares.

### Goroutine safety
We guarantee that all client instance methods are goroutine-safe and independent of each other ([guideline](https://azure.github.io/azure-sdk/golang_introduction.html#thread-safety)). This ensures that the recommendation of reusing client instances is always safe, even across goroutines.

### Additional concepts
<!-- CLIENT COMMON BAR -->
[Client options](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/policy#ClientOptions) |
[Accessing the response](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime#WithCaptureResponse) |
[Handling failures](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#ResponseError) |
[Logging](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/log)
<!-- CLIENT COMMON BAR -->

## Examples

### Create a share and upload a file

```go
const (
shareName = "sample-share"
dirName   = "sample-dir"
fileName  = "sample-file"
)

// Get a connection string to our Azure Storage account.  You can
// obtain your connection string from the Azure Portal (click
// Access Keys under Settings in the Portal Storage account blade)
// or using the Azure CLI with:
//
//     az storage account show-connection-string --name <account_name> --resource-group <resource_group>
//
// And you can provide the connection string to your application
// using an environment variable.
connectionString := "<connection_string>"

// Path to the local file to upload
localFilePath := "<path_to_local_file>"

// Get reference to a share and create it
shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
// TODO: handle error
_, err = shareClient.Create(context.TODO(), nil)
// TODO: handle error

// Get reference to a directory and create it
dirClient := shareClient.NewDirectoryClient(dirName)
_, err = dirClient.Create(context.TODO(), nil)
// TODO: handle error

// open the file for reading
file, err := os.OpenFile(localFilePath, os.O_RDONLY, 0)
// TODO: handle error
defer file.Close()

// get the size of file
fInfo, err := file.Stat()
// TODO: handle error
fSize := fInfo.Size()

// create the file
fClient := dirClient.NewFileClient(fileName)
_, err = fClient.Create(context.TODO(), fSize, nil)
// TODO: handle error

// upload the file
err = fClient.UploadFile(context.TODO(), file, nil)
// TODO: handle error
```

### Download a file

```go
const (
shareName = "sample-share"
dirName   = "sample-dir"
fileName  = "sample-file"
)

connectionString := "<connection_string>"

// Path to the save the downloaded file
localFilePath := "<path_to_local_file>"

// Get reference to the share
shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
// TODO: handle error

// Get reference to the directory
dirClient := shareClient.NewDirectoryClient(dirName)

// Get reference to the file
fClient := dirClient.NewFileClient(fileName)

// create or open a local file where we can download the Azure File
file, err := os.Create(localFilePath)
// TODO: handle error
defer file.Close()

// Download the file
_, err = fClient.DownloadFile(context.TODO(), file, nil)
// TODO: handle error
```

### Traverse a share

```go
const shareName = "sample-share"

connectionString := "<connection_string>"

// Get reference to the share
shareClient, err := share.NewClientFromConnectionString(connectionString, shareName, nil)
// TODO: handle error

// Track the remaining directories to walk, starting from the root
var dirs []*directory.Client
dirs = append(dirs, shareClient.NewRootDirectoryClient())
for len(dirs) > 0 {
    dirClient := dirs[0]
    dirs = dirs[1:]

    // Get all the next directory's files and subdirectories
    pager := dirClient.NewListFilesAndDirectoriesPager(nil)
    for pager.More() {
        resp, err := pager.NextPage(context.TODO())
        // TODO: handle error

        for _, d := range resp.Segment.Directories {
            fmt.Println(*d.Name)
            // Keep walking down directories
            dirs = append(dirs, dirClient.NewSubdirectoryClient(*d.Name))
        }

        for _, f := range resp.Segment.Files {
            fmt.Println(*f.Name)
        }
    }
}
```

## Troubleshooting

All File service operations will return an
[*azcore.ResponseError][azcore_response_error] on failure with a
populated `ErrorCode` field. Many of these errors are recoverable.
The [fileerror][file_error] package provides the possible Storage error codes
along with various helper facilities for error handling.

```go
const (
	connectionString = "<connection_string>"
	shareName    = "sample-share"
)

// create a client with the provided connection string
client, err := service.NewClientFromConnectionString(connectionString, nil)
// TODO: handle error

// try to delete the share, avoiding any potential race conditions with an in-progress or completed deletion
_, err = client.DeleteShare(context.TODO(), shareName, nil)

if fileerror.HasCode(err, fileerror.ShareBeingDeleted, fileerror.ShareNotFound) {
    // ignore any errors if the share is being deleted or already has been deleted
} else if err != nil {
    // TODO: some other error
}
```

## Next steps

Get started with our [File samples][samples].  They contain complete examples of the above snippets and more.

## Contributing

See the [Storage CONTRIBUTING.md][storage_contrib] for details on building,
testing, and contributing to this library.

This project welcomes contributions and suggestions.  Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution. For
details, visit [cla.microsoft.com][cla].

This project has adopted the [Microsoft Open Source Code of Conduct][coc].
For more information see the [Code of Conduct FAQ][coc_faq]
or contact [opencode@microsoft.com][coc_contact] with any
additional questions or comments.

<!-- LINKS -->
[source]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/storage/azfile
[docs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azfile
[rest_docs]: https://docs.microsoft.com/rest/api/storageservices/file-service-rest-api
[product_docs]: https://docs.microsoft.com/azure/storage/files/storage-files-introduction
[godevdl]: https://go.dev/dl/
[goget]: https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them
[storage_account_docs]: https://docs.microsoft.com/azure/storage/common/storage-account-overview
[storage_account_create_ps]: https://docs.microsoft.com/azure/storage/common/storage-quickstart-create-account?tabs=azure-powershell
[storage_account_create_cli]: https://docs.microsoft.com/azure/storage/common/storage-quickstart-create-account?tabs=azure-cli
[storage_account_create_portal]: https://docs.microsoft.com/azure/storage/common/storage-quickstart-create-account?tabs=azure-portal
[azure_sub]: https://azure.microsoft.com/free/
[azcore_response_error]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#ResponseError
[file_error]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/storage/azfile/fileerror/error_codes.go
[samples]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azfile/file/examples_test.go
[storage_contrib]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com
[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[storage_ad]: https://learn.microsoft.com/azure/storage/common/storage-auth-aad
