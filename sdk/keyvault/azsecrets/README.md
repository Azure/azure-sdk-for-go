# Azure Key Vault Secrets client library for Go
Azure Key Vault helps securely store and control access to tokens, passwords, certificates, API keys, and other secrets.

[Source code][secret_client_src] | [API reference documentation][reference_docs] | [Product documentation][keyvault_docs]

## Getting started

### Install packages
Install `azsecrets` and [azure-identity][azidentity_goget]:
```
go get -u github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets
```
[azure-identity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.
```
go get -u github.com/Azure/azure-sdk-for-go/sdk/azidentity
```


### Prerequisites
* An [Azure subscription][azure_sub]
* Go version 1.16 or later
* A Key Vault. If you need to create one, you can use the [Azure Cloud Shell][azure_cloud_shell] to create one with these commands (replace `"my-resource-group"` and `"my-key-vault"` with your own, unique names):

  (Optional) if you want a new resource group to hold the Key Vault:
  ```sh
  az group create --name my-resource-group --location westus2
  ```

  Create the Key Vault:
  ```Bash
  az keyvault create --resource-group my-resource-group --name my-key-vault
  ```

  Output:
  ```json
  {
      "id": "...",
      "location": "westus2",
      "name": "my-key-vault",
      "properties": {
          "accessPolicies": [...],
          "createMode": null,
          "enablePurgeProtection": null,
          "enableSoftDelete": null,
          "enabledForDeployment": false,
          "enabledForDiskEncryption": null,
          "enabledForTemplateDeployment": null,
          "networkAcls": null,
          "provisioningState": "Succeeded",
          "sku": { "name": "standard" },
          "tenantId": "...",
          "vaultUri": "https://my-key-vault.vault.azure.net/"
      },
      "resourceGroup": "my-resource-group",
      "type": "Microsoft.KeyVault/vaults"
  }
  ```

  > The `"vaultUri"` property is the `vaultUrl` used by [azsecrets.NewClient][secret_client_docs]

### Authenticate the client
This document demonstrates using [DefaultAzureCredential][default_cred_ref] to authenticate as a service principal. However, [Client][secret_client_docs] accepts any [azure-identity][azure_identity] credential. See the [azure-identity][azure_identity] documentation for more information about other credentials.


#### Create a service principal (optional)
This [Azure Cloud Shell][azure_cloud_shell] snippet shows how to create a new service principal. Before using it, replace "your-application-name" with a more appropriate name for your service principal.

Create a service principal:
```Bash
az ad sp create-for-rbac --name http://my-application --skip-assignment
```

> Output:
> ```json
> {
>     "appId": "generated app id",
>     "displayName": "my-application",
>     "name": "http://my-application",
>     "password": "random password",
>     "tenant": "tenant id"
> }
> ```

Use the output to set **AZURE_CLIENT_ID** ("appId" above), **AZURE_CLIENT_SECRET** ("password" above) and **AZURE_TENANT_ID** ("tenant" above) environment variables. The following example shows a way to do this in Bash:
```Bash
export AZURE_CLIENT_ID="generated app id"
export AZURE_CLIENT_SECRET="random password"
export AZURE_TENANT_ID="tenant id"
```

Authorize the service principal to perform key operations in your Key Vault:
```Bash
az keyvault set-policy --name my-key-vault --spn $AZURE_CLIENT_ID --secret-permissions get set list delete backup recover restore purge
```
> Possible permissions:
> - Secret management: set, backup, delete, get, list, purge, recover, restore

If you have enabled role-based access control (RBAC) for Key Vault instead, you can find roles like "Key Vault Secrets Officer" in our [RBAC guide][rbac_guide].

#### Create a client
Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and **AZURE_TENANT_ID** environment variables are set, [DefaultAzureCredential][default_cred_ref] will be able to authenticate the Client.

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal. In the Azure Portal, this URL is the vault's "DNS Name".

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)

client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)
```

## Key concepts
### Secret
A secret consists of a secret value and its associated metadata and management information. This library handles secret values as strings, but Azure Key Vault doesn't store them as such. For more information about secrets and how Key Vault stores and manages them, see the [Key Vault documentation](https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates).

Client can set secret values in the vault, update secret metadata, and delete secrets, as shown in the [examples](#examples "examples") below.

## Examples
This section contains code snippets covering common tasks:
* [Set a Secret](#set-a-secret "Set a Secret")
* [Retrieve a Secret](#retrieve-a-secret "Retrieve a Secret")
* [Update Secret metadata](#update-secret-metadata "Update Secret metadata")
* [Delete a Secret](#delete-a-secret "Delete a Secret")
* [List Secrets](#list-secrets "List Secrets")

### Set a Secret
[SetSecret](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets@v0.1.1#Client.SetSecret) creates new secrets and changes the values of existing secrets. If no secret with the given name exists, `SetSecret` creates a new secret with that name and the given value. If the given name is in use, `SetSecret` creates a new version of that secret, with the given value.

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)

resp, err := client.SetSecret(context.Background(), "secretName", "secretValue", nil)
if err != nil {
    // handle error
}

fmt.Printf("Name: %s, Value: %s\n", *resp.ID, *resp.Value)
```

### Retrieve a Secret
[GetSecret](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets@v0.1.1#Client.GetSecret) retrieves a secret previously stored in the Key Vault.

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)

resp, err := client.GetSecret(context.Background(), "mySecretName", nil)
if err != nil {
    // handle error
}

fmt.Printf("Name: %s, Value: %s\n", *resp.ID, *resp.Value)
```

### Update Secret metadata
`UpdateSecretProperties` updates a secret's metadata. It cannot change the secret's value; use [SetSecret](#set-a-secret) to set a secret's value.

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)

// Clients may specify the content type of a secret to assist in interpreting the secret data when it's retrieved
contentType := "text/plain"

// We will also disable the secret for further use

properties := azsecrets.SecretProperties{
    ContentType: &contentType,
    SecretAttributes: &azsecrets.SecretAttributes{
        Attributes: &azsecrets.Attributes{
            Enabled: to.BoolPtr(false),
        },
    },
}

resp, err := client.UpdateSecretProperties(context.Background(), "mySecretName", &properties, nil)
if err != nil {
    // handle error...
}

fmt.Printf("Updated on: %v, Content type: %v, Enabled: %v", *resp.Attributes.Updated, *resp.ContentType, *resp.Attributes.Enabled)
```

### Delete a Secret
[BeginDeleteSecret](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets@v0.1.1#Client.BeginDeleteSecret) requests Key Vault delete a secret, returning a poller which allows you to wait for the deletion to finish. Waiting is helpful when the vault has [soft-delete][soft_delete] enabled, and you want to purge (permanently delete) the secret as soon as possible. When [soft-delete][soft_delete] is disabled, `BeginDeleteSecret` itself is permanent.

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)

resp, err := client.BeginDeleteSecret(context.Background(), "secret-name", nil)
final, err := resp.PollUntilDone(context.Background(), 1 * time.Second)
```

### List secrets
[ListSecrets](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets@v0.1.1#Client.ListSecrets) lists the properties of all of the secrets in the client's vault. This list doesn't include the secret's values.

```golang
cred, err := azidentity.NewDefaultAzureCredential(nil)
client, err := azsecrets.NewClient("https://my-key-vault.vault.azure.net/", cred, nil)

pager := client.ListSecrets(nil)
for pager.NextPage(context.Background()) {
    resp := pager.PageResponse()
    for _, secret := range resp.Secrets {
        fmt.Printf("Secret ID: %s", *secret.ID)
    }
}

if pager.Err() != nil {
    // handle error
}
```

## Troubleshooting

### Error Handling

All I/O operations will return an `error` that can be investigated to discover more information about the error. In addition, you can investigate the raw response of any response object:
```golang
resp, err := client.GetSecret(context.Background(), "mySecretName", nil)
if err != nil {
    var httpErr azcore.HTTPResponse
    if errors.As(err, &httpErr) {
        // investigate httpErr.RawResponse()
    }
}
```

### Logging

This module uses the classification based logging implementation in azcore. To turn on logging set `AZURE_SDK_GO_LOGGING` to `all`. If you only want to include logs for `azsecrets`, you must create your own logger and set the log classification as `LogCredential`.

To obtain more detailed logging, including request/response bodies and header values, make sure to leave the logger as default or enable the `LogRequest` and/or `LogResponse` classificatons. A logger that only includes credential logs can be like the following:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
// Set log to output to the console
log.SetListener(func(cls log.Classification, msg string) {
		fmt.Println(msg) // printing log out to the console
})

// Includes only requests and responses in credential logs
log.SetClassifications(log.Request, log.Response)
```

> CAUTION: logs from credentials contain sensitive information.
> These logs must be protected to avoid compromising account security.

###  Additional Documentation
For more extensive documentation on Azure Key Vault, see the [API reference documentation][reference_docs].

## Contributing
This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct]. For more information, see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact opencode@microsoft.com with any additional questions or comments.

[azure_cloud_shell]: https://shell.azure.com/bash
[azure_core_exceptions]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/core/azure-core#azure-core-library-exceptions
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity/
[azidentity_goget]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[default_cred_ref]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets
[secret_client_docs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets@v0.1.1#Client
[secret_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azsecrets/client.go
[soft_delete]: https://docs.microsoft.com/azure/key-vault/general/soft-delete-overview

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazsecrets%2FREADME.png)