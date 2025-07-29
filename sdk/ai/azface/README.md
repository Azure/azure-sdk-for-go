# Azure Face Client Module for Go

Azure Face service is a cognitive service that uses machine learning algorithms to analyze faces in images.

## Getting started

### Install the package

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Face module:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/ai/azface
```

### Authentication

This package uses Azure Active Directory to authenticate. You'll need to create an Azure Cognitive Services resource in the Azure portal and obtain credentials. The client can be authenticated using one of the credential types from the [azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) module.

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

cred, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    // handle error
}

client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, nil)
if err != nil {
    // handle error
}
```

## Key concepts

### Client

A `Client` provides operations for Azure Face service.

### Face Detection

The `Detect` method analyzes an image to detect human faces and returns face attributes such as emotion scores, age, gender, and face rectangles.

## Examples

### Detect faces in an image

```go
package main

import (
    "context"
    "log"

    "github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        log.Fatal(err)
    }

    client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, nil)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.TODO()
    imageURL := "https://example.com/image.jpg"

    resp, err := client.Detect(ctx, imageURL, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Detected %d faces", len(resp.Faces))
}
```

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.