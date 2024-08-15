# Guide to migrate from `keyvault` and to `azcertificates`

This guide is intended to assist in the migration to the `azcertificates` module from the deprecated `keyvault` module. `azcertificates` allows users to create and manage [certificates][certificates] with Azure Key Vault.

## General changes

In the past, Azure Key Vault operations were all contained in a single package. For Go, this was `github.com/Azure/azure-sdk-for-go/services/keyvault/<version>/keyvault`. 

The current strategy is to break up the Key Vault into separate modules by functionality. Now there is a specific module for keys, secrets, and certificates. This guide focuses on migrating certificate operations to use the new `azcertificates` module.

Besides, module name changes. There are a number of name differences for methods and variables. All new modules also authenticate using our [azidentity] module.

## Code examples

### `keyvault` create certificate
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
        fmt.Printf("unable to add/update certificate: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("added/updated: " + *resp.ID)
}
```

### `keyvault` create certificate
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
		// this policy is suitable for a self-signed certificate
		CertificatePolicy: &azcertificates.CertificatePolicy{
			IssuerParameters:          &azcertificates.IssuerParameters{Name: to.Ptr("self")},
			X509CertificateProperties: &azcertificates.X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
		},
	}
	// if a certificate with the same name already exists, a new version of the certificate is created
	resp, err := client.CreateCertificate(context.TODO(), "certificateName", createParams, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Created a certificate with ID:", *resp.ID)
}
```

[azidentity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[certificates]: https://learn.microsoft.com/azure/key-vault/certificates/about-certificates
