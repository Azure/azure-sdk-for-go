# Azure KeyVault Administration client module for Go

>**Note:** The Administration module only works with [Managed HSM][managed_hsm] – functions targeting a Key Vault will fail.

* Vault administration (this module) - role-based access control (RBAC), settings, and vault-level backup and restore options
* Certificate management (azcertificates) - create, manage, and deploy public and private SSL/TLS certificates
* Cryptographic key management (azkeys) - create, store, and control access to the keys used to encrypt your data
* Secrets management (azsecrets) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets

Azure Key Vault Managed HSM is a fully-managed, highly-available, single-tenant, standards-compliant cloud service that enables you to safeguard
cryptographic keys for your cloud applications using FIPS 140-2 Level 3 validated HSMs.

The Azure Key Vault administration library clients support administrative tasks such as full backup / restore and key-level role-based access control (RBAC).

Source code | Package (pkg.go.dev)| [Product documentation][managed_hsm_docs] | Samples

## Getting started

### Install the package

Install `azadmin` and `azidentity` with `go get`:
```
go get github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication during client contruction.


### Prerequisites

* An [Azure subscription][azure_sub].
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* An existing [Key Vault Managed HSM][managed_hsm]. If you need to create one, you can do so using the Azure CLI by following the steps in [this document][create_managed_hsm].

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. This credential type works in both local development and production environments. We recommend using a [managed identity][managed_identity] in production.

The clients accept any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a client

Constructing the client also requires your Managed HSM's URL, which you can get from the Azure CLI or the Azure Portal.

## Key concepts

### RoleDefinition

A `RoleDefinition` is a collection of permissions. A role definition defines the operations that can be performed, such as read, write,
and delete. It can also define the operations that are excluded from allowed operations.

RoleDefinitions can be listed and specified as part of a `RoleAssignment`.

### RoleAssignment

A `RoleAssignment` is the association of a RoleDefinition to a service principal. They can be created, listed, fetched individually, and deleted.

### AccessControlClient

An `AccessControlClient` allows for management of `RoleDefinition` and `RoleAssignment` types.

### BackupClient

A `BackupClient` allows for performing full key backups, full key restores, and selective key restores.

## Examples

Get started with our examples.

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

settings, err := client.GetSettings(context.Background(), nil)
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
_, err = client.GetSettings(context.Background(), nil)
if err != nil {
    // TODO: handle error
}
// TODO: do something with response
```

## Contributing

This project welcomes contributions and suggestions.  Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution. For
details, visit <https://cla.microsoft.com>.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct].
For more information see the [Code of Conduct FAQ][coc_faq]
or contact opencode@microsoft.com with any
additional questions or comments.

<!-- LINKS -->
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free
[create_managed_hsm]: https://learn.microsoft.com/azure/key-vault/managed-hsm/quick-create-cli
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[managed_hsm]: https://docs.microsoft.com/azure/key-vault/managed-hsm/overview
[managed_hsm_docs]: https://learn.microsoft.com/azure/key-vault/managed-hsm/
[managed_identity]: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/

