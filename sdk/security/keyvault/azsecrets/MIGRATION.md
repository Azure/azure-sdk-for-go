# Guide to migrate from `keyvault` and to `azsecrets`

This guide is intended to assist in the migration to the `azsecrets` module from the deprecated `keyvault` module. `azsecrets` allows users to create and manage [secrets] with Azure Key Vault.

## General changes

In the past, Azure Key Vault operations were all contained in a single package. For Go, this was `github.com/Azure/azure-sdk-for-go/services/keyvault/<version>/keyvault`. 

The current strategy is to break up the Key Vault into separate modules by functionality. Now there is a specific module for keys, secrets, and certificates. This guide focuses on migrating secret operations to use the new `azsecrets` module.

Besides, module name changes. There are a number of name differences for methods and variables. All new modules also authenticate using our [azidentity] module.

## Code examples

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

	// If no secret with the given name exists, Key Vault creates a new secret with that name and the given value.
	// If the given name is in use, Key Vault creates a new version of that secret, with the given value.
	resp, err := client.SetSecret(context.TODO(), secretName, azsecrets.SetSecretParameters{Value: &secretValue}, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Set secret %s", resp.ID.Name())
}
```





[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[secrets]: https://learn.microsoft.com/azure/key-vault/secrets/about-secrets