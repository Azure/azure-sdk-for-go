# Release History

## 0.3.0 (2022-02-09)

### Breaking Changes
* Updated to latest `azcore`.  Public surface area is unchanged.
* [#16978](https://github.com/Azure/azure-sdk-for-go/pull/16978): The `DownloadResponse.Body` parameter is now `*RetryReaderOptions`.

### Bugs Fixed
* Fixed Issue #16193 : `azblob.GetSASToken` wrong signed resource.
* Fixed Issue #16223 : `HttpRange` does not expose its fields.
* Fixed Issue #16254 : Issue passing reader to upload `BlockBlobClient`
* Fixed Issue #16295 : Problem with listing blobs by using of `ListBlobsHierarchy()`
* Fixed Issue #16542 : Empty `StorageError` in the Azurite environment
* Fixed Issue #16679 : Unable to access Metadata when listing blobs
* Fixed Issue #16816 : `ContainerClient.GetSASToken` doesn't allow list permission.
* Fixed Issue #16988 : Too many arguments in call to `runtime.NewResponseError`

## 0.2.0 (2021-11-03)

### Breaking Changes
* Clients now have one constructor per authentication method

## 0.1.0 (2021-09-13)

### Features Added
* This is the initial preview release of the `azblob` library
