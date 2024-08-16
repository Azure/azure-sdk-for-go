# Guide to migrate from `keyvault` to `azsecrets`

This guide is intended to assist in the migration to the `azsecrets` module from the deprecated `keyvault` module. `azsecrets` allows users to create and manage [secrets] with Azure Key Vault.

## General changes

In the past, Azure Key Vault operations were all contained in a single package. For Go, this was `github.com/Azure/azure-sdk-for-go/services/keyvault/<version>/keyvault`. 

The new SDK divides the Key Vault API into separate modules for keys, secrets, and certificates. This guide focuses on migrating secret operations to use the new `azsecrets` module.

Besides, module name changes. There are a number of name differences for methods and variables. All new modules also authenticate using our [azidentity] module.

## Code examples

The following code example shows the difference between the old and new modules when creating a secret. The biggest differences are the client and authentication. In the `keyvault` module, users created a `keyvault.BaseClient` then added an `Authorizer` to the client to authenticate. In the `azsecrets` module, users create a credential using the [azidentity] module then use that credential to construct the client.

Another difference is that the Key Vault URL is now passed to the client once during construction, not every time a method is called.

### `keyvault` create secret

```go
import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

func main() {
    vaultURL := "https://<TODO: your vault name>.vault.azure.net"
    secretName := "mySecret"
	secretValue := "mySecretValue"

    authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	fmt.Println("\ncreating secret in keyvault:")
	var secParams keyvault.SecretSetParameters
	secParams.Value = &secretValue
	newBundle, err := basicClient.SetSecret(context.Background(), vaultURL, secretName, secParams)
	if err != nil {
		fmt.Printf("unable to add/update secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("added/updated: " + *newBundle.ID)
}
```

### `azsecrets` create secret

```go
package main

import (
    "context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func main() {
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
    secretName := "mySecret"
	secretValue := "mySecretValue"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	resp, err := client.SetSecret(context.TODO(), secretName, azsecrets.SetSecretParameters{Value: &secretValue}, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Set secret %s", resp.ID.Name())
}
```

[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[secrets]: https://learn.microsoft.com/azure/key-vault/secrets/about-secrets