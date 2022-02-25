# Azure Key Vault Keys client library for Go
Azure Key Vault helps solve the following problems:
- Cryptographic key management (this library) - create, store, and control
access to the keys used to encrypt your data

[Source code][key_client_src] | [Package (pkg.go.dev)][goget_azkeys] | [API reference documentation][reference_docs] | [Product documentation][keyvault_docs]

## Getting started
### Install packages
Install [azkeys][goget_azkeys] and [azidentity][goget_azidentity] with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites
* An [Azure subscription][azure_sub]
* Go version 1.16 or higher
* A Key Vault. If you need to create one, you can use the [Azure Cloud Shell][azure_cloud_shell] to create one with these commands (replace `"my-resource-group"` and `"my-key-vault"` with your own, unique
names):

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

  > The `"vaultUri"` property is the `vault_url` used by [KeyClient][key_client_docs]

### Authenticate the client
This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate as a service principal. However, [Client][key_client_docs] accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credentials.

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

Use the output to set **AZURE_CLIENT_ID** ("appId" above), **AZURE_CLIENT_SECRET**
("password" above) and **AZURE_TENANT_ID** ("tenant" above) environment variables.
The following example shows a way to do this in Bash:
```Bash
export AZURE_CLIENT_ID="generated app id"
export AZURE_CLIENT_SECRET="random password"
export AZURE_TENANT_ID="tenant id"
```

Authorize the service principal to perform key operations in your Key Vault:
```Bash
az keyvault set-policy --name my-key-vault --spn $AZURE_CLIENT_ID --key-permissions backup delete get list create update decrypt encrypt
```
> Possible permissions:
> - Key management: backup, delete, get, list, purge, recover, restore, create, update, import
> - Cryptographic operations: decrypt, encrypt, unwrapKey, wrapKey, verify, sign

If you have enabled role-based access control (RBAC) for Key Vault instead, you can find roles like "Key Vault Crypto Officer" in our [RBAC guide][rbac_guide].
If you are managing your keys using Managed HSM, read about its [access control][access_control] that supports different built-in roles isolated from Azure Resource Manager (ARM).

#### Create a client
Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and **AZURE_TENANT_ID** environment variables are set, [DefaultAzureCredential][default_cred_ref] will be able to authenticate the [Client][key_client_docs].

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal. In the Azure Portal, this URL is the vault's "DNS Name".

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

credential, err := azidentity.NewDefaultAzureCredential(nil)

client, err = azkeys.NewClient("https://my-key-vault.vault.azure.net/", credential, nil)
```

## Key concepts
### Keys
Azure Key Vault can create and store RSA and elliptic curve keys. Both can optionally be protected by hardware security modules (HSMs). Azure Key Vault can also perform cryptographic operations with them. For more information about keys and supported operations and algorithms, see the [Key Vault documentation](https://docs.microsoft.com/azure/key-vault/keys/about-keys).

[Client][key_client_docs] can create keys in the vault, get existing keys from the vault, update key metadata, and delete keys, as shown in the [examples](#examples "examples") below.

## Examples
This section contains code snippets covering common tasks:
* [Create a key](#create-a-key "Create a key")
* [Retrieve a key](#retrieve-a-key "Retrieve a key")
* [Update an existing key](#update-an-existing-key "Update an existing key")
* [Delete a key](#delete-a-key "Delete a key")
<!-- * [Configure automatic key rotation](#configure-automatic-key-rotation "Configure automatic key rotation") -->
* [List keys](#list-keys "List keys")
<!-- * [Perform cryptographic operations](#cryptographic-operations) -->

### Create a key
[`CreateRSAKey`](https://aka.ms/azsdk/go/keyvault-keys) and
[`CreateECKey`](https://aka.ms/azsdk/go/keyvault-keys) create RSA and elliptic curve keys in the vault, respectively. If a key with the same name already exists, a new version of that key is created.

```go
import (
    "fmt"

    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleCreateKeys() {
    vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        panic(err)
    }

    client, err := azkeys.NewClient(vaultUrl, cred, nil)
    if err != nil {
        panic(err)
    }

    // Create RSA Key
    resp, err := client.CreateRSAKey(context.TODO(), "new-rsa-key", &azkeys.CreateRSAKeyOptions{KeySize: to.Int32Ptr(2048)})
    if err != nil {
        panic(err)
    }
    fmt.Println(*resp.Key.ID)
    fmt.Println(*resp.Key.KeyType)

    // Create EC Key
    resp, err := client.CreateECKey(context.TODO(), "new-rsa-key", &azkeys.CreateECKeyOptions{CurveName: azkeys.P256.ToPtr()})
    if err != nil {
        panic(err)
    }
    fmt.Println(*resp.Key.ID)
    fmt.Println(*resp.Key.KeyType)
}
```

### Retrieve a key
[`GetKey`](https://aka.ms/azsdk/go/keyvault-keys) retrieves a key previously stored in the Vault.
```go
import (
    "fmt"

    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleRetrieveKey() {
    credential, err := azidentity.NewDefaultAzureCredential(nil)

    client, err = azkeys.NewClient("https://my-key-vault.vault.azure.net/", credential, nil)
    resp, err := client.GetKey(context.TODO(), "key-to-retrieve", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Key.ID)
}
```

### Update an existing key
[`UpdateKeyProperties`](https://aka.ms/azsdk/go/keyvault-keys)
updates the properties of a key previously stored in the Key Vault.
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleClient_UpdateKeyProperties() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.UpdateKeyProperties(context.TODO(), "key-to-update", &azkeys.UpdateKeyPropertiesOptions{
		Tags: map[string]string{
			"Tag1": "val1",
		},
		KeyAttributes: &azkeys.KeyAttributes{
			RecoveryLevel: azkeys.CustomizedRecoverablePurgeable.ToPtr(),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Attributes.RecoveryLevel, *resp.Tags["Tag1"])
}
```

### Delete a key
[`BeginDeleteKey`](https://aka.ms/azsdk/go/keyvault-keys) requests Key Vault delete a key, returning a poller which allows you to wait for the deletion to finish. Waiting is helpful when the vault has [soft-delete][soft_delete] enabled, and you want to purge (permanently delete) the key as soon as possible. When [soft-delete][soft_delete] is disabled, `BeginDeleteKey` itself is permanent.

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleClient_BeginDeleteKey() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	pollResp, err := resp.PollUntilDone(context.TODO(), 1 * time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully deleted key %s", *pollResp.Key.ID)
}
```

### Configure automatic key rotation
`update_key_rotation_policy` allows you to configure automatic key rotation for a key by specifying a rotation policy.
In addition, `rotate_key` allows you to rotate a key on-demand by creating a new version of the given key.

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

credential, err := azidentity.NewDefaultAzureCredential(nil)
handle(err)

client, err = azkeys.NewClient("https://my-key-vault.vault.azure.net/", credential, nil)
handle(err)

// Set the key's automated rotation policy to rotate the key 30 days before the key expires
resp, err = client.UpdateKeyRotationPolicy(ctx, "key-name", &UpdateKeyRotationPolicyOptions{
    Attributes: &KeyRotationPolicyAttributes{
        ExpiryTime: to.StringPtr("P90D"),
    },
    LifetimeActions: []*LifetimeActions{
        {
            Action: &LifetimeActionsType{
                Type: ActionTypeNotify.ToPtr(),
            },
            Trigger: &LifetimeActionsTrigger{
                TimeBeforeExpiry: to.StringPtr("P30D"),
            },
        },
    },
})
handle(err)

currentPolicyResp, err := client.GetKeyRotationPolicy(context.TODO(), "key-name", nil)
handle(err)

// Finally, you can rotate a key on-demand by creating a new version of the key
rotatedResp, err := client.RotateKey.rotate_key(context.TODO(), "key-name", nil)
handle(err)
```

### List keys
[`ListKeys`](https://aka.ms/azsdk/go/keyvault-keys) lists the properties of all of the keys in the client's vault.

```go
import (
    "fmt"

    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleClient_ListKeys() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.ListKeys(nil)
	for pager.NextPage(context.TODO()) {
		for _, key := range pager.PageResponse().Keys {
			fmt.Println(*key.KID)
		}
	}

	if pager.Err() != nil {
		panic(pager.Err())
	}
}
```

### Cryptographic operations
[CryptographyClient](https://aka.ms/azsdk/go/keyvault-keys)
enables cryptographic operations (encrypt/decrypt, wrap/unwrap, sign/verify) using a particular key.

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

credential, err := azidentity.NewDefaultAzureCredential(nil)

client, err = crypto.NewClient("https://my-key-vault.vault.azure.net/keys/<my-key>/<key-version>", credential, nil)

encryptResponse, err := cryptoClient.Encrypt(ctx, AlgorithmRSAOAEP, []byte("plaintext"), nil)
```

See the [package documentation][crypto_client_docs] for more details of the cryptography API.

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

###  Additional Documentation
For more extensive documentation on Azure Key Vault, see the
[API reference documentation][reference_docs].

## Contributing
This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct].
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact opencode@microsoft.com with any additional questions or comments.

[access_control]: https://docs.microsoft.com/azure/key-vault/managed-hsm/access-control
[azure_cloud_shell]: https://shell.azure.com/bash
[azure_core_exceptions]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/core/azure-core#azure-core-library-exceptions
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[goget_azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[default_cred_ref]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[goget_azkeys]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://aka.ms/azsdk/go/keyvault-keys
[key_client_docs]: https://aka.ms/azsdk/go/keyvault-keys#Client
[crypto_client_docs]: https://aka.ms/azsdk/go/keyvault-keys/crypto/docs
[key_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/fd86ba6a0ece5a0658dd16f8d3d564493369a8a2/sdk/keyvault/azkeys/client.go
[key_samples]: https://github.com/Azure/azure-sdk-for-go/tree/fd86ba6a0ece5a0658dd16f8d3d564493369a8a2/sdk/keyvault/azkeys/example_test.go
[soft_delete]: https://docs.microsoft.com/azure/key-vault/general/soft-delete-overview

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazkeys%2FREADME.png)
