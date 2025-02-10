# Azure Key Vault Certificates client module for Go

* Certificate management (this module) - create, manage, and deploy public and private SSL/TLS certificates
* Managed HSM administration ([azadmin](https://aka.ms/azsdk/go/keyvault-admin/docs)) - role-based access control (RBAC), settings, and vault-level backup and restore options
* Cryptographic key management ([azkeys](https://aka.ms/azsdk/go/keyvault-keys/docs)) - create, store, and control access to the keys used to encrypt your data
* Secrets management ([azsecrets](https://aka.ms/azsdk/go/keyvault-secrets/docs)) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets

[Source code][certificates_client_src] | [Package (pkg.go.dev)][reference_docs] |  [Product documentation][keyvault_docs] | [Samples][certificates_samples]

## Getting started

### Install the package

Install `azcertificates` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates
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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", credential, nil)
	if err != nil {
		// TODO: handle error
	}
}
```

## Key concepts

### Client

With a [Client][client_docs], you can get certificates from the vault, create new certificates and
new versions of existing certificates, update certificate metadata, and delete certificates. You
can also manage certificate issuers, contacts, and management policies of certificates. This is
illustrated in the [examples](#examples) below.

## Examples

Get started with our [examples][certificates_samples].

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

// Includes only requests and responses in logs
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
[azure_keyvault_cli]: https://learn.microsoft.com/azure/key-vault/general/quick-create-cli
[azure_keyvault_portal]: https://learn.microsoft.com/azure/key-vault/general/quick-create-portal
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://learn.microsoft.com/azure/key-vault/
[client_docs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates#Client
[reference_docs]: https://aka.ms/azsdk/go/keyvault-certificates/docs
[certificates_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/security/keyvault/azcertificates
[certificates_samples]: https://aka.ms/azsdk/go/keyvault-certificates/docs#pkg-examples
[managed_identity]: https://learn.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview


