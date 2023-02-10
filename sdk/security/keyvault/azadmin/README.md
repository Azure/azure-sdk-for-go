# Azure KeyVault Administration client library for Go

Azure Key Vault Managed HSM is a fully-managed, highly-available, single-tenant, standards-compliant cloud service that enables you to safeguard
cryptographic keys for your cloud applications using FIPS 140-2 Level 3 validated HSMs.

The Azure Key Vault administration library clients support administrative tasks such as full backup / restore and key-level role-based access control (RBAC).

## Getting started

### Install the package

Install `azadmin` and `azidentity` with `go get`:
```
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.


### Prerequisites

* An [Azure subscription][azure_sub].
* An existing Azure Key Vault. If you need to create an Azure Key Vault, you can use the [Azure CLI][azure_cli].
* Authorization to an existing Azure Key Vault using either [RBAC][rbac_guide] (recommended) or [access control][access_policy].

To create a Managed HSM resource, run the following CLI command:

```PowerShell
az keyvault create --hsm-name <your-key-vault-name> --resource-group <your-resource-group-name> --administrators <your-user-object-id> --location <your-azure-location>
```

To get `<your-user-object-id>` you can run the following CLI command:

```PowerShell
az ad user show --id <your-user-principal> --query id
```

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. This credential type works in both local development and production environments. We recommend using a [managed identity][managed_identity] in production.

Client accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a client

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

- Example AccessControlClient:
- Example BackupClient:
- Example SettingsClient:

#### Activate your managed HSM

All data plane commands are disabled until the HSM is activated. You will not be able to create keys or assign roles.
Only the designated administrators that were assigned during the create command can activate the HSM. To activate the HSM you must download the security domain.

To activate your HSM you need:

* A minimum of 3 RSA key-pairs (maximum 10)
* Specify the minimum number of keys required to decrypt the security domain (quorum)

To activate the HSM you send at least 3 (maximum 10) RSA public keys to the HSM. The HSM encrypts the security domain with these keys and sends it back.
Once this security domain is successfully downloaded, your HSM is ready to use.
You also need to specify quorum, which is the minimum number of private keys required to decrypt the security domain.

The example below shows how to use openssl to generate 3 self-signed certificates.

```PowerShell
openssl req -newkey rsa:2048 -nodes -keyout cert_0.key -x509 -days 365 -out cert_0.cer
openssl req -newkey rsa:2048 -nodes -keyout cert_1.key -x509 -days 365 -out cert_1.cer
openssl req -newkey rsa:2048 -nodes -keyout cert_2.key -x509 -days 365 -out cert_2.cer
```

Use the `az keyvault security-domain download` command to download the security domain and activate your managed HSM.
The example below uses 3 RSA key pairs (only public keys are needed for this command) and sets the quorum to 2.

```PowerShell
az keyvault security-domain download --hsm-name <your-managed-hsm-name> --sd-wrapping-keys ./certs/cert_0.cer ./certs/cert_1.cer ./certs/cert_2.cer --sd-quorum 2 --security-domain-file ContosoMHSM-SD.json
```

#### Controlling access to your managed HSM

The designated administrators assigned during creation are automatically added to the "Managed HSM Administrators" [built-in role][built_in_roles],
who are able to download a security domain and [manage roles for data plane access][access_control], among other limited permissions.

To perform other actions on keys, you need to assign principals to other roles such as "Managed HSM Crypto User", which can perform non-destructive key operations:

```PowerShell
az keyvault role assignment create --hsm-name <your-managed-hsm-name> --role "Managed HSM Crypto User" --scope / --assignee-object-id <principal-or-user-object-ID> --assignee-principal-type <principal-type>
```

Please read [best practices][best_practices] for properly securing your managed HSM.

## Key concepts

### RoleDefinition

A `RoleDefinition` is a collection of permissions. A role definition defines the operations that can be performed, such as read, write,
and delete. It can also define the operations that are excluded from allowed operations.

RoleDefinitions can be listed and specified as part of a `RoleAssignment`.

### RoleAssignment

A `RoleAssignment` is the association of a RoleDefinition to a service principal. They can be created, listed, fetched individually, and deleted.

### AccessControlClient

A `AccessControlClient` provides both synchronous and asynchronous operations allowing for management of `KeyVaultRoleDefinition` and `KeyVaultRoleAssignment` objects.

### BackupClient

A `BackupClient` provides both synchronous and asynchronous operations for performing full key backups, full key restores, and selective key restores.

### BeginFullBackup

A `BeginFullBackup` represents a long running operation for a full key backup.

### RestoreOperation

A `BeginFullRestoreOperation` represents a long running operation for both a full key and selective key restore.

## Examples

Get started with our [examples][azadmin_examples].

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err := client.GetSecret(context.Background(), "secretName", nil)
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
_, err = client.GetSecret(ctx, "secretName", nil)
if err != nil {
    // TODO: handle error
}
// TODO: do something with response
```

To learn more about other logging mechanisms see [here][logging].

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
[access_control]: https://learn.microsoft.com/azure/key-vault/managed-hsm/access-control
[access_policy]: https://learn.microsoft.com/azure/key-vault/general/assign-access-policy
[azadmin_examples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin#pkg-examples
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[rbac_guide]: https://learn.microsoft.com/azure/key-vault/general/rbac-guide
[azure_cli]: https://learn.microsoft.com/cli/azure
[best_practices]: https://learn.microsoft.com/azure/key-vault/managed-hsm/best-practices
[built_in_roles]: https://learn.microsoft.com/azure/key-vault/managed-hsm/built-in-roles
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[managed_identity]: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/

