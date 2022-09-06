# Azure Key Vault Certificates client module for Go

* Certificate management (this module) - create, manage, and deploy public and private SSL/TLS certificates
* Cryptographic key management ([azkeys](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys)) - create, store, and control access to the keys used to encrypt your data
* Secrets management ([azsecrets](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets)) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets

[Source code][certificates_client_src] | [Package (pkg.go.dev)][reference_docs] |  [Product documentation][keyvault_docs] | [Samples][certificates_samples]

## Getting started

### Install the package

Install `azcertificates` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

* An [Azure subscription][azure_sub]
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* A key vault. If you need to create one, see the Key Vault documentation for instructions on doing so in the [Azure Portal][azure_keyvault_portal] or with the [Azure CLI][azure_keyvault_cli].

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. This credential type works in both local development and production environments. We recommend using a [managed identity][managed_identity] in production.

[Client][client_docs] accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a client

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", credential, nil)
}
```

## Key concepts

### Client

With a [Client][client_docs] you can get certificates from the vault, create new certificates and
new versions of existing certificates, update certificate metadata, and delete certificates. You
can also manage certificate issuers, contacts, and management policies of certificates. This is
illustrated in the [examples](#examples) below.

## Examples

This section contains code snippets covering common tasks:
* [Create a Certificate](#create-a-certificate)
* [Delete a Certificate](#delete-a-certificate)
* [List Certificates](#list-certificates)
* [Retrieve a Certificate](#retrieve-a-certificate)
* [Update Certificate Metadata](#update-certificate-metadata)

### Create a Certificate

[CreateCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.CreateCertificate)
creates a certificate to be stored in the key vault. If a certificate with the same name already exists, a new
version of the certificate is created.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	client := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", credential, nil)

	createParams := azcertificates.CreateCertificateParameters{
		// this policy is suitable for a self-signed certificate
		CertificatePolicy: &azcertificates.CertificatePolicy{
			IssuerParameters:          &azcertificates.IssuerParameters{Name: to.Ptr("self")},
			X509CertificateProperties: &azcertificates.X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
		},
	}
	resp, err := client.CreateCertificate(context.TODO(), "certificateName", createParams, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Created a certificate with ID:", *resp.ID)
}
```

### Retrieve a Certificate

[GetCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.GetCertificate)
retrieves the latest version of a certificate previously stored in the Key Vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	// passing an empty string for the version gets the latest version of the certificate
	resp, err := client.GetCertificate(context.TODO(), "certName", "", nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.ID)
}
```


### Update Certificate metadata

[UpdateCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.UpdateCertificate)
updates a certificate's metadata.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	updateParams := azcertificates.UpdateCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{Enabled: to.Ptr(false)},
	}
	// passing an empty string for the version updates the latest version of the certificate
	resp, err := client.UpdateCertificate(context.TODO(), "certName", "", updateParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.ID)
}
```

### Delete a Certificate

[DeleteCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.DeleteCertificate) requests that Key Vault delete a certificate. It returns when Key Vault has begun deleting the certificate. Deletion can take several seconds to complete, so it may be necessary to wait before performing other operations on the deleted certificate.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	// DeleteCertificate returns when Key Vault has begun deleting the certificate. That can take several
	// seconds to complete, so it may be necessary to wait before performing other operations on the
	// deleted certificate.
	resp, err := client.DeleteCertificate(context.TODO(), "certName", nil)
	if err != nil {
		// TODO: handle error
	}

	// In a soft-delete enabled vault, deleted resources can be recovered until they're purged (permanently deleted).
	fmt.Printf("Certificate will be purged at %v", *resp.ScheduledPurgeDate)
}
```

### List Certificates

[NewListCertificatesPager](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.NewListCertificatesPager) creates a pager that lists all certificates in the vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	pager := client.NewListCertificatesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, cert := range page.Value {
			fmt.Println(*cert.ID)
		}
	}
}
```

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err := client.GetCertificate(context.Background(), "certificateName", nil)
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

You can access the raw `*http.Response` returned by Key Vault using the `runtime.WithCaptureResponse` method and a context passed to any client method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

var response *http.Response
ctx := runtime.WithCaptureResponse(context.TODO(), &response)
_, err = client.GetCertificate(ctx, "certificateName", nil)
if err != nil {
    // TODO: handle error
}
// TODO: do something with response
```

###  Additional Documentation

For more extensive documentation on Azure Key Vault, see the [API reference documentation][reference_docs].

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct]. For more information, see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact opencode@microsoft.com with any additional questions or comments.

[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[azure_keyvault_cli]: https://docs.microsoft.com/azure/key-vault/general/quick-create-cli
[azure_keyvault_portal]: https://docs.microsoft.com/azure/key-vault/general/quick-create-portal
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[client_docs]: https://aka.ms/azsdk/go/azcertificates
[reference_docs]: https://aka.ms/azsdk/go/keyvault-certificates/docs
[certificates_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azcertificates
[certificates_samples]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azcertificates/example_test.go
[managed_identity]: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazcertificates%2FREADME.png)
