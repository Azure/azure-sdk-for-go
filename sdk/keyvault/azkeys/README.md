# Azure Key Vault Keys client module for Go

* Cryptographic key management (this module) - create, store, and control access to the keys used to encrypt your data
* Secrets management ([azsecrets](https://aka.ms/azsdk/go/keyvault-secrets/docs)) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets
* Certificate management ([azcertificates](https://aka.ms/azsdk/go/keyvault-certificates/docs)) - create, manage, and deploy public and private SSL/TLS certificates

[Source code][key_client_src] | [Package (pkg.go.dev)][goget_azkeys] | [Product documentation][keyvault_docs] | [Samples][keys_samples]

## Getting started

### Install packages

Install `azkeys` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

* An [Azure subscription][azure_sub]
* Go version 1.18 or higher
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

  > The `"vaultUri"` property is the ``vaultURL` argument for [NewClient][key_client_docs]

### Authenticate the client

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate as a service principal. However, [Client][key_client_docs] accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credentials.

#### Create a service principal (optional)

This [Azure Cloud Shell][azure_cloud_shell] snippet shows how to create a new service principal. Before using it, replace "my-application" with an appropriate name for your application.

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
If you are managing your keys using Managed HSM, read about its [access control][mhsm_access_control], which supports different roles isolated from Azure Resource Manager (ARM).

#### Create a client
Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and **AZURE_TENANT_ID** environment variables are set, [DefaultAzureCredential][default_cred_ref] will be able to authenticate the [Client][key_client_docs].

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
}
```

## Key concepts

### Keys

Azure Key Vault can create and store RSA and elliptic curve keys. Both can optionally be protected by hardware security modules (HSMs). Azure Key Vault can also perform cryptographic operations with them. For more information about keys and supported operations and algorithms, see the [Key Vault documentation](https://docs.microsoft.com/azure/key-vault/keys/about-keys).

[Client][key_client_docs] can create keys in the vault, get existing keys from the vault, update key metadata, and delete keys, as shown in the examples below.

## Examples

This section contains code snippets covering common tasks:
* [Configure automatic key rotation](#configure-automatic-key-rotation)
* [Create a key](#create-a-key)
* [Delete a key](#delete-a-key)
* [List keys](#list-keys)
* [Retrieve a key](#retrieve-a-key)
* [Update an existing key](#update-an-existing-key)

### Create a key

[`CreateKey`](https://aka.ms/azsdk/go/keyvault-keys/docs#Client.CreateKey) creates keys in the vault. If a key with the same name already exists, a new version of that key is created.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	// Create RSA Key
	rsaParams := azkeys.CreateKeyParameters{
		KeySize: to.Ptr(int32(2048)),
		Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSA),
	}
	resp, err := client.CreateKey(context.TODO(), "new-rsa-key", rsaParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)

	// Create EC Key
	ecParams := azkeys.CreateKeyParameters{
		Curve: to.Ptr(azkeys.JSONWebKeyCurveNameP256K),
		Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC),
	}
	resp, err = client.CreateKey(context.TODO(), "new-ec-key", ecParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}
```

### Retrieve a key

[`GetKey`](https://aka.ms/azsdk/go/keyvault-keys/docs#Client.GetKey) retrieves a key previously stored in the Vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	// passing an empty string for the version parameter gets the latest version of the key
	version := ""
	resp, err := client.GetKey(context.TODO(), "key-name", version, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}
```

### Update an existing key

[`UpdateKey`](https://aka.ms/azsdk/go/keyvault-keys/docs#Client.UpdateKey)
updates the properties of a key previously stored in the Key Vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	params := azkeys.UpdateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			Expires: to.Ptr(time.Now().Add(48 * time.Hour)),
		},
		// Key Vault doesn't interpret tags. The keys and values are up to your application.
		Tags: map[string]*string{"expiraton-extended": to.Ptr("true")},
	}
	// passing an empty string for the version parameter updates the latest version of the key
	updateResp, err := client.UpdateKey(context.TODO(), "key-name", "", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Updated key %s", *updateResp.Key.KID)
}
```

### Delete a key

[`DeleteKey`](https://aka.ms/azsdk/go/keyvault-keys/docs#Client.DeleteKey) requests that Key Vault delete a key. It returns when Key Vault has begun deleting the key. Deletion can take several seconds to complete, so it may be necessary to wait before performing other operations on the deleted key.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	// DeleteKey returns when Key Vault has begun deleting the key. That can take several
	// seconds to complete, so it may be necessary to wait before performing other operations
	// on the deleted key.
	resp, err := client.DeleteKey(context.TODO(), "key-name", nil)
	if err != nil {
		// TODO: handle error
	}

	// In a soft-delete enabled vault, deleted keys can be recovered until they're purged (permanently deleted).
	fmt.Printf("Key will be purged at %v", resp.ScheduledPurgeDate)
}
```

### Configure automatic key rotation

`UpdateKeyRotationPolicy` allows you to configure automatic key rotation for a key by specifying a rotation policy, and
`RotateKey` allows you to rotate a key on demand. See [Azure Key Vault documentation](https://docs.microsoft.com/azure/key-vault/keys/how-to-configure-key-rotation) for more information about key rotation policies.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	// this policy rotates the key every 18 months
	policy := azkeys.KeyRotationPolicy{
		LifetimeActions: []*azkeys.LifetimeActions{
			{
				Action: &azkeys.LifetimeActionsType{
					Type: to.Ptr(azkeys.ActionTypeRotate),
				},
				Trigger: &azkeys.LifetimeActionsTrigger{
					TimeAfterCreate: to.Ptr("P18M"),
				},
			},
		},
	}
	resp, err := client.UpdateKeyRotationPolicy(context.TODO(), "key-name", policy, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Updated key rotation policy at: %v", resp.Attributes.Updated)
}
```

### List keys

[`NewListKeysPager`](https://aka.ms/azsdk/go/keyvault-keys/docs#Client.NewListKeysPager) creates a pager that lists all keys in the client's vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)

	pager := client.NewListKeysPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, key := range resp.Value {
			fmt.Println(*key.KID)
		}
	}
}
```

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err := client.GetKey(context.Background(), "keyName", nil)
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
_, err = client.GetKey(ctx, "keyName", nil)
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


[mhsm_access_control]: https://docs.microsoft.com/azure/key-vault/managed-hsm/access-control
[azure_cloud_shell]: https://shell.azure.com/bash
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[default_cred_ref]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://aka.ms/azsdk/go/keyvault-keys/docs
[key_client_docs]: https://aka.ms/azsdk/go/keyvault-keys/docs#Client
[key_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azkeys/client.go
[keys_samples]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azkeys/example_test.go

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazkeys%2FREADME.png)
