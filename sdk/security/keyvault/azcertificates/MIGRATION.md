# Guide to migrate from `keyvault` to `azcertificates`

This guide is intended to assist in the migration to the `azcertificates` module from the deprecated `keyvault` module. `azcertificates` allows users to create and manage [certificates][certificates] with Azure Key Vault.

## General changes

In the past, Azure Key Vault operations were all contained in a single package. For Go, this was `github.com/Azure/azure-sdk-for-go/services/keyvault/<version>/keyvault`. 

The new SDK divides the Key Vault API into separate modules for keys, secrets, and certificates. This guide focuses on migrating certificate operations to use the new `azcertificates` module.

There are other changes besides the module name. For example, some type and method names are different, and all new modules authenticate using our [azidentity] module.

## Code example

The following code example shows the difference between the old and new modules when creating a certificate. The biggest differences are the client and authentication. In the `keyvault` module, users created a `keyvault.BaseClient` then added an `Authorizer` to the client to authenticate. In the `azcertificates` module, users create a credential using the [azidentity] module then use that credential to construct the client.

Another difference is that the Key Vault URL is now passed to the client once during construction, not every time a method is called.

### `keyvault` create certificate
```go
import (
    "context"
    "fmt"

    "github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
    kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

func main() {
    vaultURL := "https://<TODO: your vault name>.vault.azure.net"
    authorizer, err := kvauth.NewAuthorizerFromEnvironment()
    if err != nil {
       // TODO: handle error
    }

    basicClient := keyvault.New()
    basicClient.Authorizer = authorizer

    fmt.Println("\ncreating certificate in keyvault:")
    issuerName := "self"
    subject := "CN=DefaultPolicy"
    createParams := keyvault.CertificateCreateParameters{
        CertificatePolicy: &keyvault.CertificatePolicy{
            IssuerParameters:          &keyvault.IssuerParameters{Name: &issuerName},
            X509CertificateProperties: &keyvault.X509CertificateProperties{Subject: &subject},
        }
    }
    resp, err := basicClient.CreateCertificate(context.TODO(), vaultURL, "<cert name>", createParams)
    if err != nil {
        // TODO: handle error
    }
    fmt.Println("added/updated: " + *resp.ID)
}
```

### `azcertificates` create certificate
```go
import (
    "context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

func main() {
    vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	
	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	createParams := azcertificates.CreateCertificateParameters{
		CertificatePolicy: &azcertificates.CertificatePolicy{
			IssuerParameters:          &azcertificates.IssuerParameters{Name: to.Ptr("self")},
			X509CertificateProperties: &azcertificates.X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
		},
	}
	resp, err := client.CreateCertificate(context.TODO(), "<cert name>", createParams, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Created a certificate with ID:", *resp.ID)
}
```

[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[certificates]: https://learn.microsoft.com/azure/key-vault/certificates/about-certificates
