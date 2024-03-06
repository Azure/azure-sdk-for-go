# Azure KeyVault Administration client module for Go

>**Note:** The Administration module only works with [Managed HSM][managed_hsm] – functions targeting a Key Vault will fail.

* Managed HSM administration (this module) - role-based access control (RBAC), settings, and vault-level backup and restore options
* Certificate management ([azcertificates](https://aka.ms/azsdk/go/keyvault-certificates/docs)) - create, manage, and deploy public and private SSL/TLS certificates
* Cryptographic key management ([azkeys](https://aka.ms/azsdk/go/keyvault-keys/docs)) - create, store, and control access to the keys used to encrypt your data
* Secrets management ([azsecrets](https://aka.ms/azsdk/go/keyvault-secrets/docs)) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets

Azure Key Vault Managed HSM is a fully-managed, highly-available, single-tenant, standards-compliant cloud service that enables you to safeguard
cryptographic keys for your cloud applications using FIPS 140-2 Level 3 validated HSMs.

The Azure Key Vault administration library clients support administrative tasks such as full backup / restore, key-level role-based access control (RBAC), and settings management.

[Source code][azadmin_repo] | [Package (pkg.go.dev)][azadmin_pkg_go]| [Product documentation][managed_hsm_docs] | Samples ([backup][azadmin_pkg_go_samples_backup], [rbac][azadmin_pkg_go_samples_rbac], [settings][azadmin_pkg_go_samples_settings])

## Getting started

### Install the package

Install `azadmin` and `azidentity` with `go get`:
```
go get github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication. It creates a credential which is passed to the client contructor as shown in the examples below.


### Prerequisites

* An [Azure subscription][azure_sub].
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* An existing [Key Vault Managed HSM][managed_hsm]. If you need to create one, you can do so [using the Azure CLI][create_managed_hsm].

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. This credential type works in both local development and production environments. We recommend using a [managed identity][managed_identity] in production.

The clients accept any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a client

Constructing the client also requires your Managed HSM's URL, which you can get from the Azure CLI or the Azure Portal.

- [Example backup client](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup#example-NewClient)
- [Example rbac client][azadmin_pkg_go_samples_rbac]
- [Example settings client](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings#example-NewClient)

## Key concepts

### RoleDefinition

A `RoleDefinition` is a collection of permissions. A role definition defines the operations that can be performed, such as read, write, and delete. It can also define the operations that are excluded from allowed operations.

A `RoleDefinition` can be listed and specified as part of a `RoleAssignment`.

### RoleAssignment

A `RoleAssignment` is the association of a RoleDefinition to a service principal. They can be created, listed, fetched individually, and deleted.

### rbac.Client

An `rbac.Client` manages `RoleDefinition` and `RoleAssignment` types.

### backup.Client

A `backup.Client` performs full key backups, full key restores, and selective key restores.

### settings.Client

A `settings.Client` provides methods to update, get, and list settings for a Managed HSM.

## Examples

Get started with our examples:
- [backup][azadmin_pkg_go_samples_backup]  
- [rbac][azadmin_pkg_go_samples_rbac]
- [settings][azadmin_pkg_go_samples_settings]

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

// Includes only requests and responses in logs
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
[azadmin_repo]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/security/keyvault/azadmin
[azadmin_pkg_go]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin
[azadmin_pkg_go_samples_backup]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup#pkg-examples
[azadmin_pkg_go_samples_rbac]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac#pkg-examples
[azadmin_pkg_go_samples_settings]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings#pkg-examples
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free
[create_managed_hsm]: https://learn.microsoft.com/azure/key-vault/managed-hsm/quick-create-cli
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[managed_hsm]: https://docs.microsoft.com/azure/key-vault/managed-hsm/overview
[managed_hsm_docs]: https://learn.microsoft.com/azure/key-vault/managed-hsm/
[managed_identity]: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/

