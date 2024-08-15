# Guide to migrate from `keyvault` and to `azkeys`

This guide is intended to assist in the migration to the `azkeys` module from the deprecated `keyvault` module. `azkeys` allows users to create and manage [keys][keys] with Azure Key Vault.

## General changes

In the past, Azure Key Vault operations were all contained in a single package. For Go, this was `github.com/Azure/azure-sdk-for-go/services/keyvault/<version>/keyvault`. 

The current strategy is to break up the Key Vault into separate modules by functionality. Now there is a specific module for keys, secrets, and certificates. This guide focuses on migrating keys operations to use the new `azkeys` module.

Besides, module name changes. There are a number of name differences for methods and variables. All new modules also authenticate using our [azidentity] module.

## Code examples

### `keyvault` create key
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
    authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	fmt.Println("\ncreating a key in keyvault:")
    keyParams := keyvault.KeyCreateParameters{
        Curve: &keyvault.P256,
        Kty:   &keyvault.EC,
    }
	newBundle, err := basicClient.CreateKey(context.TODO(), vaultURL, "<key name>", keyParams)
	if err != nil {
		fmt.Printf("unable to add/update key: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("added/updated: " + *newBundle.JSONWebKey.Kid)
}
```

### `azkeys` create key
```go
package main

import (
    "context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

func main() {
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azkeys.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	keyParams := azkeys.CreateKeyParameters{
		Curve: to.Ptr(azkeys.CurveNameP256K),
		Kty:   to.Ptr(azkeys.KeyTypeEC),
	}
	// if a key with the same name already exists, a new version of that key is created
	resp, err := client.CreateKey(context.TODO(), "<key name>", keyParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}
```


[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[keys]: https://learn.microsoft.com/azure/key-vault/keys/about-keys
