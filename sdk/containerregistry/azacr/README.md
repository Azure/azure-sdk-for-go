# Azure Container Registry client module for Go

Azure Container Registry allows you to store and manage container images and artifacts in a private registry for all types of container deployments.

Use the client library for Azure Container Registry to:

- List images or artifacts in a registry
- Obtain metadata for images and artifacts, repositories and tags
- Set read/write/delete properties on registry items
- Delete images and artifacts, repositories and tags

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/containerregistry/azacr) | [Package (pkg.go.dev)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr) | [REST API documentation](https://docs.microsoft.com/rest/api/containerregistry/) | [Product documentation](https://docs.microsoft.com/azure/container-registry/) | [Samples]()

## Getting started

### Install packages

Install `azacr` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

- An [Azure subscription](https://azure.microsoft.com/free/)
- A supported Go version (the Azure SDK supports the two most recent Go releases)
- A [Container Registry service instance](https://docs.microsoft.com/azure/container-registry/container-registry-intro)

To create a new Container Registry, you can use the [Azure Portal](https://docs.microsoft.com/azure/container-registry/container-registry-get-started-portal),
[Azure PowerShell](https://docs.microsoft.com/azure/container-registry/container-registry-get-started-powershell), or the [Azure CLI](https://docs.microsoft.com/azure/container-registry/container-registry-get-started-azure-cli).
Here's an example using the Azure CLI:

```Powershell
az acr create --name MyContainerRegistry --resource-group MyResourceGroup --location westus --sku Basic
```
### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential) to authenticate. 
This credential type works in both local development and production environments. 
We recommend using a [managed identity](https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview) in production.

[ContainerRegistryClient](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr#ContainerRegistryClient) and [ContainerRegistryBlobClient](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr#ContainerRegistryBlobClient) accepts any [azidentity][azure_identity] credential. 
See the [azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) documentation for more information about other credential types.

#### Create a client

Constructing the client requires your Container Registry's endpoint URL, which you can get from the Azure CLI or the Azure Portal.

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
}
```

## Key concepts

A **registry** stores Docker images and [OCI Artifacts](https://opencontainers.org/). 
An image or artifact consists of a **manifest** and **layers**. 
An image's manifest describes the layers that make up the image, and is uniquely identified by its **digest**. 
An image can also be "tagged" to give it a human-readable alias. 
An image or artifact can have zero or more **tags** associated with it, and each tag uniquely identifies the image. 
A collection of images that share the same name but have different tags, is referred to as a **repository**.

For more information please see [Container Registry Concepts](https://docs.microsoft.com/azure/container-registry/container-registry-concepts).

## Examples

Get started with our [examples](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr#pkg-examples).

### Listing repositories

```go
import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Repositories.Repositories {
			fmt.Printf("repository: %s\n", *v)
		}
	}
}
```

### List tags with anonymous access

```go
import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", nil, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	pager := client.NewListTagsPager("library/hello-world", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.TagAttributeBases {
			fmt.Printf("tag: %s\n", *v.Name)
		}
	}
}
```

### Set artifact properties

```go
import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	res, err := client.UpdateTagProperties(ctx, "library/hello-world", "latest", &azacr.ContainerRegistryClientUpdateTagPropertiesOptions{
		Value: &azacr.TagWriteableProperties{
			CanWrite: to.Ptr(false),
			CanDelete: to.Ptr(false),
		}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository library/hello-world - tag latest: 'CanWrite' property: %t, 'CanDelete' property: %t\n", *res.Tag.ChangeableAttributes.CanWrite, *res.Tag.ChangeableAttributes.CanDelete)
}
```

### Delete images

```go
import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	repositoryPager := client.NewListRepositoriesPager(nil)
	for repositoryPager.More() {
		repositoryPage, err := repositoryPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance repository page: %v", err)
		}
		for _, r := range repositoryPage.Repositories.Repositories {
			manifestPager := client.NewListManifestsPager(*r, &azacr.ContainerRegistryClientListManifestsOptions{
				OrderBy: to.Ptr(azacr.ArtifactManifestOrderByLastUpdatedOnDescending),
			})
			for manifestPager.More() {
				manifestPage, err := manifestPager.NextPage(ctx)
				if err != nil {
					log.Fatalf("failed to advance manifest page: %v", err)
				}
				imagesToKeep := 3
				for i, m := range manifestPage.Manifests.Manifests {
					if i >= imagesToKeep {
						for _, t := range m.Tags {
							fmt.Printf("delete tag from image: %s", *t)
							_, err := client.DeleteTag(ctx, *r, *t, nil)
							if err != nil {
								log.Fatalf("failed to delete tag: %v", err)
							}
						}
						_, err := client.DeleteManifest(ctx, *r, *m.Digest, nil)
						if err != nil {
							log.Fatalf("failed to delete manifest: %v", err)
						}
						fmt.Printf("delete image with digest: %s", *m.Digest)
					}
				}
			}
		}
	}
}
```

### Upload and download blob

```go
import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"io"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azacr.NewContainerRegistryBlobClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	uploadRes, err := client.UploadBlob(ctx, "library/hello-world", streaming.NopCloser(bytes.NewReader([]byte("hello world"))), nil)
	if err != nil {
		log.Fatalf("failed to upload blob: %v", err)
	}
	downloadRes, err := client.GetBlob(ctx, "library/hello-world", uploadRes.Digest, nil)
	if err != nil {
		log.Fatalf("failed to download blob: %v", err)
	}
	blob, err := io.ReadAll(downloadRes.Body)
	if err != nil {
		log.Fatalf("failed to read blob: %v", err)
	}
	fmt.Printf("blob content: %s", blob)
}
```

### Upload manifest

```go
import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azacr.NewContainerRegistryClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	blobClient, err := azacr.NewContainerRegistryBlobClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create blob client: %v", err)
	}
	ctx := context.Background()
	layer := "hello world"
	layerUploadRes, err := blobClient.UploadBlob(ctx, "library/hello-world", streaming.NopCloser(bytes.NewReader([]byte(layer))), nil)
	if err != nil {
		log.Fatalf("failed to upload layer blob: %v", err)
	}
	config := fmt.Sprintf(`{
  architecture: "amd64",
  os: "windows",
  rootfs: {
	type: "layers",
	diff_ids: [%s],
  },
}`, layerUploadRes.Digest)
	configUploadRes, err := blobClient.UploadBlob(ctx, "library/hello-world", streaming.NopCloser(bytes.NewReader([]byte(config))), nil)
	if err != nil {
		log.Fatalf("failed to upload config blob: %v", err)
	}
	manifest := fmt.Sprintf(`{
  schemaVersion: 2,
  config: {
	mediaType: "application/vnd.oci.image.config.v1+json",
	digest: %s,
	size: %d,
  },
  layers: [
	{
	  mediaType: "application/vnd.oci.image.layer.v1.tar",
	  digest: %s,
	  size: %d,
	  annotations: {
		title: "artifact.txt",
	  },
	  },
  ],
}`, configUploadRes.Digest, len(config), layerUploadRes.Digest, len(layer))
	uploadManifestRes, err := client.UploadManifest(ctx, "library/hello-world", "1.0.0", "application/vnd.oci.image.config.v1+json", streaming.NopCloser(bytes.NewReader([]byte(manifest))), nil)
	if err != nil {
		log.Fatalf("failed to upload manifest: %v", err)
	}
	fmt.Printf("digest of uploaded manifest: %s", *uploadManifestRes.DockerContentDigest)
}
```

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Container Registry.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err := client.GetRepositoryProperties(ctx, "library/hello-world", nil)
if err != nil {
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		// TODO: investigate httpErr
	} else {
		// TODO: not an HTTP error
	}
}
```

### Logging

This module uses the logging implementation in `azcore`. To turn on logging for all Azure SDK modules, set `AZURE_SDK_GO_LOGGING` to `all`. By default the logger writes to stderr. Use the `azcore/log` package to control log output. For example, logging only HTTP request and response events, and printing them to stdout:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

// Print log events to stdout
azlog.SetListener(func(cls azlog.Event, msg string) {
	fmt.Println(msg)
})

// Includes only requests and responses in credential logs
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

### Accessing `http.Response`

You can access the raw `*http.Response` returned by Container Registry using the `runtime.WithCaptureResponse` method and a context passed to any client method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

var response *http.Response
ctx := runtime.WithCaptureResponse(context.TODO(), &response)
_, err = client.GetRepositoryProperties(ctx, "library/hello-world", nil)
if err != nil {
	// TODO: handle error
}
// TODO: do something with response
```

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][https://opensource.microsoft.com/codeofconduct/]. For more information, see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact opencode@microsoft.com with any additional questions or comments.


