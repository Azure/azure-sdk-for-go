# Guide to migrate from `azure-storage-file-go` to `azfile`

This guide is intended to assist in the migration from the `azure-storage-file-go` module to the latest releases of the `azfile` module.

## Simplified API surface area

The redesign of the `azfile` module separates clients into various sub-packages.
In previous versions, the public surface area was "flat", so all clients and supporting types were in the `azfile` package.
This made it difficult to navigate the public surface area.

## Clients

In `azure-storage-file-go` a client constructor always requires a `url.URL` and `Pipeline` parameters.

In `azfile` a client constructor always requires a `string` URL, any specified credential type, and a `*ClientOptions` for optional values.  You pass `nil` to accept default options.

```go
// new code to create service client using shared key credential
client, err := service.NewClientWithSharedKeyCredential("<my storage account URL>", sharedKeyCred, nil)
```

## Authentication

In `azure-storage-file-go` you created a `Pipeline` with the required credential type. This pipeline was then passed to the client constructor.

In `azfile`, you pass the required credential directly to the client constructor.

```go
// new code.  cred is a shared key credential created from NewSharedKeyCredential method in service package
client, err := service.NewClientWithSharedKeyCredential("<my storage account URL>", cred, nil)
```

Authentication with a shared key via `NewSharedKeyCredential` remains unchanged.

In `azure-storage-file-go` you created a `Pipeline` with `NewAnonymousCredential` to support anonymous or SAS authentication.

In `azfile` you use the constructor `NewClientWithNoCredential()` instead.

```go
// new code to create a file client using SAS
client, err := file.NewClientWithNoCredential("<public file or file with SAS URL>", nil)
```

## Listing files/directories

In `azure-storage-file-go` you explicitly created a `Marker` type that was used to page over results ([example](https://pkg.go.dev/github.com/Azure/azure-storage-file-go/azfile?utm_source=godoc#example-package)).

In `azfile`, operations that return paginated values return a `*runtime.Pager[T]`.

```go
// new code to iterate a directory
dirClient, err := directory.NewClientWithSharedKeyCredential("<directory URL>", cred, options)
// TODO: handle err

pager := dirClient.NewListFilesAndDirectoriesPager(nil)
for pager.More() {
	page, err := pager.NextPage(context.TODO())
	// process results
}
```

## Configuring the HTTP pipeline

In `azure-storage-file-go` you explicitly created a HTTP pipeline with configuration before creating a client.
This pipeline instance was then passed as an argument to the client constructor ([example](https://pkg.go.dev/github.com/Azure/azure-storage-file-go/azfile?utm_source=godoc#example-NewPipeline)).

In `azfile` a HTTP pipeline is created during client construction.  The pipeline is configured through the `azcore.ClientOptions` type.

```go
// new code
client, err := service.NewClientWithSharedKeyCredential("<my storage account URL>", cred, &service.ClientOptions{
	ClientOptions: azcore.ClientOptions{
		// configure HTTP pipeline options here
	},
})
```
