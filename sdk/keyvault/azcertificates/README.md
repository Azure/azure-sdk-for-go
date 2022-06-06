# Azure Key Vault Certificates client module for Go

* Certificate management (this module) - create, manage, and deploy public and private SSL/TLS certificates
* Cryptographic key management (([azkeys](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys))) - create, store, and control access to the keys used to encrypt your data
* Secrets management ([azsecrets](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets)) - securely store and control access to tokens, passwords, certificates, API keys, and other secrets

[Source code][certificates_client_src] | [Package (pkg.go.dev)][reference_docs] |  [Product documentation][keyvault_docs] | [Samples][certificates_samples]

## Getting started

### Install the package

Install `azcertificates` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

* An [Azure subscription][azure_sub]
* Go 1.18 or later
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

  > The `"vaultUri"` property is the `vaultURL` argument for [NewClient][certificate_client_docs]

### Authenticate the client

This document demonstrates using [DefaultAzureCredential][default_cred_ref] to authenticate as a service principal. However, [NewClient][certificate_client_docs]
accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credentials.

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

Authorize the service principal to perform certificate operations in your Key Vault:
```Bash
az keyvault set-policy --name my-key-vault --spn $AZURE_CLIENT_ID --certificate-permissions backup create delete get import list purge recover restore update
```
> Possible certificate permissions: backup, create, delete, deleteissuers, get, getissuers, import, list, listissuers, managecontacts, manageissuers, purge, recover, restore, setissuers, update

If you have enabled role-based access control (RBAC) for Key Vault instead, you can find roles like "Key Vault Certificates Officer" in our [RBAC guide][rbac_guide].

#### Create a client

Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and **AZURE_TENANT_ID** environment variables are set, [DefaultAzureCredential][default_cred_ref] will be able to authenticate the [Client][key_client_docs].

Constructing the client also requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", credential, nil)
	if err != nil {
		// TODO: handle error
	}
}
```

## Key concepts

### Client

With a [Client][certificate_client_docs] you can get certificates from the vault, create new certificates and
new versions of existing certificates, update certificate metadata, and delete certificates. You
can also manage certificate issuers, contacts, and management policies of certificates. This is
illustrated in the [examples](#examples) below.

## Examples

This section contains code snippets covering common tasks:
* [Create a Certificate](#create-a-certificate)
* [Delete a Certificate](#delete-a-certificate)
* [List Properties of Certificates](#list-properties-of-certificates)
* [Retrieve a Certificate](#retrieve-a-certificate)
* [Update Properties of an existing Certificate](#update-properties-of-an-existing-certificate)

### Create a Certificate

[BeginCreateCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.BeginCreateCertificate)
creates a certificate to be stored in the Azure Key Vault. If a certificate with the same name already exists, a new
version of the certificate is created. Before creating a certificate, a management policy for the certificate can be
created or our default policy will be used. This method returns a poller object that enables waiting for the operation
to complete.

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", credential, nil)
	if err != nil {
		// TODO: handle error
	}

	resp, err := client.BeginCreateCertificate(context.TODO(), "certificateName", azcertificates.Policy{
		IssuerParameters: &azcertificates.IssuerParameters{
			IssuerName: to.Ptr("Self"),
		},
		X509Properties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr("CN=DefaultPolicy"),
		},
	}, nil)
	if err != nil {
		// TODO: handle error
	}

	finalResponse, err := resp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Created a certificate with ID: ", *finalResponse.ID)
}
```
If you would like to check the status of your certificate creation, you can call `Poll(ctx context.Context)` on the poller or
[GetCertificateOperation](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.GetCertificateOperation)
with the name of the certificate.

### Retrieve a Certificate

[GetCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.GetCertificate)
retrieves the latest version of a certificate previously stored in the Key Vault.

```go
import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	resp, err := client.GetCertificate(context.TODO(), "myCertName", nil)
	if err != nil {
		// TODO: handle error
	}

	// optionally you can get a specific version
	resp, err = client.GetCertificate(context.TODO(), "myCertName", &azcertificates.GetCertificateOptions{Version: "myCertVersion"})
	if err != nil {
		// TODO: handle error
	}
}
```


### Update properties of an existing Certificate

[UpdateCertificateProperties](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.UpdateCertificateProperties)
updates a certificate previously stored in the Key Vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	resp, err := client.GetCertificate(context.TODO(), "myCertName", nil)
	if err != nil {
		// TODO: handle error
	}

	resp.Properties.Enabled = to.Ptr(false)
	updateResp, err := client.UpdateCertificateProperties(context.TODO(), *resp.Properties, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Set Enabled to %v for certificate with name %s\n", *&updateResp.Properties.Enabled, *resp.ID)
}
```

### Delete a Certificate

[BeginDeleteCertificate](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.BeginDeleteCertificate)
requests Key Vault delete a certificate, returning a poller which allows you to wait for the deletion to finish. Waiting is helpful when you want to purge (permanently delete) the certificate as soon as possible.

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	pollerResp, err := client.BeginDeleteCertificate(context.TODO(), "certToDelete", nil)
	if err != nil {
		// TODO: handle error
	}
	finalResp, err := pollerResp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Deleted certificate with ID: ", *finalResp.ID)
}
```

### List Certificates

[NewListPropertiesOfCertificatesPager](https://aka.ms/azsdk/go/keyvault-certificates/docs#Client.NewListPropertiesOfCertificatesPager) creates a pager that lists the properties of all certificates in the client's vault.

```go
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azcertificates.NewClient("https://<TODO: your vault name>.vault.azure.net", cred, nil)
	if err != nil {
		// TODO: handle error
	}

	pager := client.NewListPropertiesOfCertificatesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, cert := range page.Certificates {
			fmt.Println(*cert.ID)
		}
	}
}
```

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err := client.GetCertificate(context.Background(), "certificateName", nil)
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
_, err = client.GetCertificate(ctx, "certificateName", nil)
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

[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[azure_cloud_shell]: https://shell.azure.com/bash
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[certificate_client_docs]: https://aka.ms/azsdk/go/azcertificates
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://aka.ms/azsdk/go/keyvault-certificates/docs
[certificates_client_src]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azcertificates
[certificates_samples]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azcertificates/example_test.go

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazcertificates%2FREADME.png)
