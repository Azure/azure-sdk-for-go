# Azure Key Vault Secrets client module for Go
Azure Key Vault helps solve the following problems:
* Secrets management (this module) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets
* Cryptographic key management ([azkeys](https://azsdk/go/keyvault-keys/docs)) - create, store, and control access to the keys used to encrypt your data
* Certificate management ([azcertificates](https://aka.ms/azsdk/go/keyvault-certificates/docs)) - create, manage, and deploy public and private SSL/TLS certificates

[Source code][module_source] | [Package (pkg.go.dev)][reference_docs] | [Product documentation][keyvault_docs] | [Samples][secrets_samples]

## Getting started

### Install packages
Install `azsecrets` and `azidentity`:
```
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets

go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.


### Prerequisites
* An [Azure subscription][azure_sub]
* Go version 1.18 or later
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

  > The `"vaultUri"` property is the `vaultURL` for [azsecrets.NewClient][secret_client_docs]

### Authenticate the client
This document demonstrates using [azidentity.DefaultAzureCredential][default_cred] to authenticate as a service principal. However, [Client][secret_client_docs] accepts any `azidentity` credential. See the [azidentity documentation][azure_identity] for information about other credentials.


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
Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and **AZURE_TENANT_ID** environment variables are set, [DefaultAzureCredential][default_cred] will be able to authenticate the Client.

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

```golang
import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}
}
```

## Key concepts
### Secret
A secret consists of a secret value and its associated metadata and management information. This library handles secret values as strings, but Azure Key Vault doesn't store them as such. For more information about secrets and how Key Vault stores and manages them, see the [Key Vault documentation](https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates).

`azseecrets.Client` can set secret values in the vault, update secret metadata, and delete secrets, as shown in the examples below.

## Examples
This section contains code snippets covering common tasks:
* [Set a Secret](#set-a-secret "Set a Secret")
* [Retrieve a Secret](#retrieve-a-secret "Retrieve a Secret")
* [Update Secret metadata](#update-secret-metadata "Update Secret metadata")
* [Delete a Secret](#delete-a-secret "Delete a Secret")
* [List Secrets](#list-secrets "List Secrets")

### Set a Secret
[SetSecret](https://aka.ms/azsdk/go/keyvault-secrets/docs#Client.SetSecret) creates new secrets and changes the values of existing secrets. If no secret with the given name exists, `SetSecret` creates a new secret with that name and the given value. If the given name is in use, `SetSecret` creates a new version of that secret, with the given value.

```golang
import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	secretName := "mySecret"
	secretValue := "mySecretValue"

	resp, err := client.SetSecret(context.TODO(), secretName, secretValue, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Set secret %s", *resp.Secret.ID)
}
```

### Retrieve a Secret
[GetSecret](https://aka.ms/azsdk/go/keyvault-secrets/docs#Client.GetSecret) retrieves a secret previously stored in the Key Vault.

```golang
import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetSecret(context.TODO(), "mySecretName", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Secret Name: %s\tSecret Value: %s", *resp.Secret.ID, *resp.Secret.Value)
}
```

### Update Secret metadata
`UpdateSecretProperties` updates a secret's metadata. It cannot change the secret's value; use [SetSecret](#set-a-secret) to set a secret's value.

```golang
import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	getResp, err := client.GetSecret(context.TODO(), "secret-to-update", nil)
	if err != nil {
		panic(err)
	}

	if getResp.Secret.Properties == nil {
		getResp.Secret.Properties = &azsecrets.Properties{}
	}
	getResp.Secret.Properties = &azsecrets.Properties{
		Enabled:     to.Ptr(true),
		ExpiresOn:   to.Ptr(time.Now().Add(48 * time.Hour)),
		NotBefore:   to.Ptr(time.Now().Add(-24 * time.Hour)),
		ContentType: to.Ptr("password"),
		Tags:        map[string]string{"Tag1": "Tag1Value"},
		// Remember to preserve the name and version
		Name:    getResp.Secret.Properties.Name,
		Version: getResp.Secret.Properties.Version,
	}
	resp, err := client.UpdateSecretProperties(context.TODO(), getResp.Secret, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated secret with ID: %s\n", *resp.Secret.ID)
}
```

### Delete a Secret
[BeginDeleteSecret](https://aka.ms/azsdk/go/keyvault-secrets/docs#Client.BeginDeleteSecret) requests Key Vault delete a secret, returning a poller which allows you to wait for the deletion to finish. Waiting is helpful when you want to purge (permanently delete) a secret as soon as possible in a vault with [soft-delete][soft_delete] enabled, as it is by default. `BeginDeleteSecret` itself is permanent in a vault with soft-delete disabled.

```golang
import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteSecret(context.TODO(), "secretToDelete", nil)
	if err != nil {
		panic(err)
	}

	// If you do not care when the secret is deleted, you do not have to
	// call resp.PollUntilDone.
	_, err = resp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		panic(err)
	}
}
```

### List secrets
[NewListPropertiesOfSecretsPager](https://aka.ms/azsdk/go/keyvault-secrets/docs#Client.NewListPropertiesOfSecretsPager) creates a `Pager` that lists the properties of all of the secrets in the client's vault, not including their secret values.

```golang
import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.NewListPropertiesOfSecretsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, v := range page.Secrets {
			fmt.Printf("Secret Name: %s\tSecret Tags: %v\n", *v.ID, v.Tags)
		}
	}
}
```

## Troubleshooting

### Error Handling

All I/O operations return a non-nil `error` when they fail. These errors include their motivating HTTP response, if any:
```golang
resp, err := client.GetSecret(context.TODO(), "mySecretName", nil)
if err != nil {
    var httpErr *azcore.ResponseError
    if errors.As(err, &httpErr) {
        // investigate httpErr.RawResponse
    }
}
```

### Logging

This module uses the classification based logging implementation in `azcore`. To turn on logging for all Azure SDK modules, set `AZURE_SDK_GO_LOGGING` to `all`. By default the logger writes to stderr. Use the `azcore/log` package to control log output. For example, logging only HTTP request and response events, and printing them to stdout:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

azlog.SetListener(func(cls azlog.Event, msg string) {
    fmt.Println(msg)
})

azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

### Accessing `http.Response`
You can access the raw `*http.Response` returned by Key Vault using the `runtime.WithCaptureResponse` method and a context passed to any client method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

func main() {
    var respFromCtx *http.Response
    ctx := runtime.WithCaptureResponse(context.TODO(), &respFromCtx)
    _, err = client.GetSecret(ctx, "mySecretName", nil)
    if err != nil {
        panic(err)
    }
    fmt.Println(respFromCtx.StatusCode)
}
```

###  Additional Documentation
See the [API reference documentation][reference_docs] for complete documentation of this module.

## Contributing
This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct]. For more information, see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact opencode@microsoft.com with any additional questions or comments.

[azure_cloud_shell]: https://shell.azure.com/bash
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[default_cred]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://aka.ms/azsdk/go/keyvault-secrets/docs
[secret_client_docs]: https://aka.ms/azsdk/go/keyvault-secrets/docs#Client
[module_source]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azsecrets
[secrets_samples]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azsecrets/example_test.go
[soft_delete]: https://docs.microsoft.com/azure/key-vault/general/soft-delete-overview

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazsecrets%2FREADME.png)
