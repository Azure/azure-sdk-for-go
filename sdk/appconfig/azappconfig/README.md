# Azure App Configuration client library for Go
Azure App Configuration is a managed service that helps developers centralize their application and feature settings simply and securely.
It allows you to create and manage application configuration settings and retrieve their revisions from a specific point in time.

[Source code][appconfig_client_src] | [Package (pkg.go.dev)][goget_azappconfig] | [Product documentation][appconfig_docs]

## Getting started
### Install packages
Install [azappconfig][goget_azappconfig] with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig
```

### Prerequisites
* An [Azure subscription][azure_sub]
* Go version 1.16 or higher
* An App Configuration store. If you need to create one, you can use the [Azure Cloud Shell][azure_cloud_shell] to create one with these commands (replace `"my-resource-group"` and `"my-app-config"` with your own, unique
names):

  (Optional) if you want a new resource group to hold the Azure App Configuration:
  ```sh
  az group create --name my-resource-group --location westus2
  ```

  Create the Key Vault:
  ```Bash
  az appconfig create --resource-group my-resource-group --name my-app-config --location westus2
  ```

  Output:
  ```json
  {
      "creationDate": "...",
      "endpoint": "https://my-app-config.azconfig.io",
      "id": "/subscriptions/.../resourceGroups/my-resource-group/providers/Microsoft.AppConfiguration/configurationStores/my-app-config",
      "location": "westus2",
      "name": "my-app-config",
      "provisioningState": "Succeeded",
      "resourceGroup": "my-resource-group",
      "tags": {},
      "type": "Microsoft.AppConfiguration/configurationStores"
  }
  ```

  > The `"endpoint"` property is the `endpointUrl` used by [Client][appconfig_client_src]

### Authenticate the client
This document demonstrates using the connection string. However, [Client][appconfig_client_src] accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credentials.

#### Create a client
Constructing the client requires your App Configuration connection string, which you can get from the Azure Portal.
```Bash
export APPCONFIGURATION_CONNECTION_STRING="Endpoint=https://my-app-config.azconfig.io;Id=...;Secret=..."
```

```go
import (
    "os"
    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
)

client, err = azappconfig.NewClientFromConnectionString(os.Getenv("APPCONFIGURATION_CONNECTION_STRING"), nil)
```

Or, using Default Azure Credential from Azure Identity:

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

credential, err := azidentity.NewDefaultAzureCredential(nil)

client, err = azappconfig.NewClient("https://my-app-config.azconfig.io", credential, nil)
```

## Examples
This section contains code snippets covering common tasks:
* [Add a configuration setting](#add-a-configuration-setting "Add a configuration setting")
* [Get a configuration setting](#get-a-configuration-setting "Get a configuration setting")
* [Set a configuration setting](#set-a-configuration-setting "Set a configuration setting")
* [Set a configuration setting read only](#set-a-configuration-setting-read-only "Set a configuration setting read only")
* [List configuration setting revisions](#list-configuration-setting-revisions "List configuration setting revisions")
* [Delete a configuration setting](#set-a-configuration-setting "Delete a configuration setting")

### Add a configuration setting

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleAddConfigurationSetting() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    // Create configuration setting
    resp, err := client.AddSetting(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label"),
            Value: to.StringPtr("value")
        },
        nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
    fmt.Println(*resp.Value)
}
```

### Get a configuration setting

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleGetConfigurationSetting() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    // Get configuration setting
    resp, err := client.GetSetting(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label")
        },
        nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
    fmt.Println(*resp.Value)
}
```

### Set a configuration setting

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleSetConfigurationSetting() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    // Set configuration setting
    resp, err := client.SetSetting(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label"),
            Value: to.StringPtr("new_value")
        },
        nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
    fmt.Println(*resp.Value)
}
```

### Set a configuration setting read only

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleSetConfigurationSettingReadOnly() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    // Set configuration setting read only
    resp, err := client.SetReadOnly(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label")
        },
        true,
        nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
    fmt.Println(*resp.Value)
    fmt.Println(*resp.IsReadOnly)

    // Remove read only status
    resp, err := client.SetReadOnly(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label")
        },
        false,
        nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
    fmt.Println(*resp.Value)
    fmt.Println(*resp.IsReadOnly)
}
```

### List configuration setting revisions

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleListRevisions() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    revPgr := client.ListRevisions(
        azappconfig.SettingSelector{
            KeyFilter: to.StringPtr("*"),
            LabelFilter: to.StringPtr("*"),
            Fields: azappconfig.AllSettingFields()
        },
        nil)

    for revPgr.More() {
        if revResp, revErr := revPgr.NextPage(context.TODO()); revErr == nil {
            for _, setting := range revResp.Settings {
                fmt.Println(*setting.Key)
                fmt.Println(*setting.Label)
                fmt.Println(*setting.Value)
                fmt.Println(*setting.IsReadOnly)
            }
        }
    }
}
```

### Delete a configuration setting

```go
import (
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleDeleteConfigurationSetting() {
    connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
    client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
    if err != nil {
        panic(err)
    }

    // Delete configuration setting
    resp, err := client.DeleteSetting(
        context.TODO(),
        azappconfig.Setting{
            Key: to.StringPtr("key"),
            Label: to.StringPtr("label")
        },
        nil)

    if err != nil {
        panic(err)
    }
    fmt.Println(*resp.Key)
    fmt.Println(*resp.Label)
}
```

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

[azure_cloud_shell]: https://shell.azure.com/bash
[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[appconfig_docs]: https://docs.microsoft.com/azure/azure-app-configuration/
[goget_azappconfig]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig
[appconfig_client_src]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/appconfig/azappconfig/client.go

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fkeyvault%2Fazkeys%2FREADME.png)
