# Azure Key Vault Certificates client library for Go
Azure Key Vault helps solve the following problems:
- Certificate management (this library) - create, manage, and deploy public and private SSL/TLS certificates
- Cryptographic key management
([azkeys](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azkeys)) - create, store, and control access to the keys used to encrypt your data
- Secrets management
([azsecrets](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azsecrets)) -
securely store and control access to tokens, passwords, certificates, API keys,
and other secrets

[Source code][certificates_client_src] | [pkg.go.dev][pkggodev_azcerts] | [API reference documentation][reference_docs] | [Product documentation][keyvault_docs] | [Samples][certificates_samples]

## Getting started
### Install the package
Install [azure-keyvault-certificates][pkggodev_azcerts] and [azidentity][azure_identity_goget] with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites
* An [Azure subscription][azure_sub]
* Go 1.16 or later
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

  > The `"vaultUri"` property is the `vaultURL` parameter used by the `azcertificates.NewClient` function.

### Authenticate the client
This document demonstrates using [DefaultAzureCredential][default_cred_ref] to authenticate as a service principal. However, [NewClient][certificate_client_docs]
accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credentials.

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

Authorize the service principal to perform certificate operations in your Key Vault:
```Bash
az keyvault set-policy --name my-key-vault --spn $AZURE_CLIENT_ID --certificate-permissions backup create delete get import list purge recover restore update
```
> Possible certificate permissions: backup, create, delete, deleteissuers, get, getissuers, import, list, listissuers, managecontacts, manageissuers, purge, recover, restore, setissuers, update

If you have enabled role-based access control (RBAC) for Key Vault instead, you can find roles like "Key Vault Certificates Officer" in our [RBAC guide][rbac_guide].

#### Create a client
Once the **AZURE_CLIENT_ID**, **AZURE_CLIENT_SECRET** and
**AZURE_TENANT_ID** environment variables are set,
[DefaultAzureCredential][default_cred_ref] will be able to authenticate the
[Client][certificate_client_docs].

Constructing the client also requires your vault's URL, which you can
get from the Azure CLI or the Azure Portal. In the Azure Portal, this URL is
the vault's "DNS Name".

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = azkeys.NewClient("https://my-key-vault.vault.azure.net/", credential, nil)
	if err != nil {
		panic(err)
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
* [Create a Certificate](#create-a-certificate "Create a Certificate")
* [Retrieve a Certificate](#retrieve-a-certificate "Retrieve a Certificate")
* [Update Properties of an existing Certificate](#update-properties-of-an-existing-certificate "Update Properties of an existing Certificate")
* [Delete a Certificate](#delete-a-certificate "Delete a Certificate")
* [List Properties of Certificates](#list-properties-of-certificates "List Properties of Certificates")

### Create a Certificate
[BeginCreateCertificate](https://aka.ms/azsdk/go/azcertificates)
creates a certificate to be stored in the Azure Key Vault. If a certificate with the same name already exists, a new
version of the certificate is created. Before creating a certificate, a management policy for the certificate can be
created or our default policy will be used. This method returns a long running operation poller.
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	client, err = azkeys.NewClient("https://my-key-vault.vault.azure.net/", credential, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginCreateCertificate(ctx, certName, CertificatePolicy{
		IssuerParameters: &IssuerParameters{
			Name: to.StringPtr("Self"),
		},
		X509CertificateProperties: &X509CertificateProperties{
			Subject: to.StringPtr("CN=DefaultPolicy"),
		},
	}, nil)
	if err != nil {
		panic(err)
	}

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	if err != nil {
		panic(err)
	}
	fmt.Println(*pollerResp.ID)
}
```
If you would like to check the status of your certificate creation, you can call `Poll(ctx context.Context)` on the poller or
[GetCertificateOperation](https://aka.ms/azsdk/go/azcertificates)
with the name of the certificate.

### Retrieve a Certificate
[GetCertificate](https://aka.ms/azsdk/go/azcertificates)
retrieves the latest version of a certificate previously stored in the Key Vault.
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func Example_GetCertificate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetCertificate(context.TODO(), "myCertName", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.ID)
	fmt.Println(*resp.Policy.IssuerParameters.Name)

	// optionally you can get a specific version
	resp, err = client.GetCertificate(context.TODO(), "myCertName", &azcertificates.GetCertificateOptions{Version: "myCertVersion"})
	if err != nil {
		panic(err)
	}
}
```


### Update properties of an existing Certificate
[UpdateCertificateProperties](https://aka.ms/azsdk/go/azcertificates)
updates a certificate previously stored in the Key Vault.
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.UpdateCertificateProperties(context.TODO(), "myCertName", &azcertificates.UpdateCertificatePropertiesOptions{
		Version: "myNewVersion",
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Attributes: azcertificates.Attributes{Enabled: to.BoolPtr(false)},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.ID)
	fmt.Println(*resp.Certificate.Attributes.Enabled)
}
```

### Delete a Certificate
[BeginDeleteCertificate](https://aka.ms/azsdk/go/azcertificates)
requests Key Vault delete a certificate, returning a poller which allows you to wait for the deletion to finish.
Waiting is helpful when the vault has [soft-delete][soft_delete] enabled, and you want to purge
(permanently delete) the certificate as soon as possible. When [soft-delete][soft_delete] is disabled,
`BeginDeleteCertificate` itself is permanent.

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteCertificate(context.TODO(), "myCertificateName", nil)
	if err != nil {
		panic(err)
	}

	finalResponse, err := resp.PollUntilDone(context.TODO(), time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println(*finalResponse.ID)
	fmt.Println(*finalResponse.DeletedDate)
}
```

### List  Certificates
[ListCertificates](https://aka.ms/azsdk/go/azcertificates)
lists the properties of all certificates in the specified Key Vault.
```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	poller := client.ListCertificates(nil)
	for poller.NextPage(context.TODO()) {
		for _, cert := range poller.PageResponse().Certificates {
			fmt.Println(*cert.ID)
		}
	}
	if poller.Err() != nil {
		panic(err)
	}
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

[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[azure_cloud_shell]: https://shell.azure.com/bash
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[azure_identity_goget]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[keyvault_docs]: https://docs.microsoft.com/azure/key-vault/
[pkggodev_azcerts]: https://pypi.org/project/azure-keyvault-certificates/
[certificate_client_docs]: https://aka.ms/azsdk/go/azcertificates
[rbac_guide]: https://docs.microsoft.com/azure/key-vault/general/rbac-guide
[reference_docs]: https://aka.ms/azsdk/go/azcertificates
[certificates_client_src]: https://aka.ms/azsdk/go/azcertificates
[certificates_samples]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/keyvault/azcertificates/example_test.go
[soft_delete]: https://docs.microsoft.com/azure/key-vault/general/soft-delete-overview

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazcertificates%2FREADME.png)