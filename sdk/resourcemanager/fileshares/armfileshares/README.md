# Azure File Shares Management client library for Go

Azure File Shares Management client library for Go (`armfileshares`) provides management capabilities for Azure File Shares resources. This library follows the Azure SDK Design Guidelines for Go.

## Getting started

### Prerequisites

* Go 1.24 or later
* An Azure subscription
* Azure Storage account

### Install the package

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares
```

### Authentication

This module uses Azure Identity for authentication. Install it with:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

## Key concepts

This library provides management operations for Azure File Shares, including:
* Creating and managing file shares
* Configuring share properties
* Managing access policies

## Examples

```go
package main

import (
    "context"
    "log"
    
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
)

func main() {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        log.Fatal(err)
    }
    
    client, err := armfileshares.NewClient("<subscription-id>", cred, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the client to manage file shares
    // ...
}
```

## Troubleshooting

For troubleshooting information, please refer to the [Azure SDK for Go troubleshooting guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/TROUBLESHOOTING.md).

## Next steps

For more examples and detailed documentation, visit the [Azure SDK for Go documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go).

## Contributing

For details on contributing to this repository, see the [contributing guide][azure_sdk_for_go_contributing].

This project welcomes contributions and suggestions. Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

### Additional Helpful Links for Contributors

Many people all over the world have helped make this project better. You'll want to check out:

* [What are some good first issues for new contributors to the repo?](https://github.com/azure/azure-sdk-for-go/issues?q=is%3Aopen+is%3Aissue+label%3A%22up+for+grabs%22)
* [How to build and test your change][azure_sdk_for_go_contributing_developer_guide]
* [How you can make a change happen!][azure_sdk_for_go_contributing_pull_requests]
* Frequently Asked Questions (FAQ) and Conceptual Topics in the detailed [Azure SDK for Go wiki](https://github.com/azure/azure-sdk-for-go/wiki).

### Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter](https://www.microsoft.com/msrc/faqs-report-an-issue).

### License

Azure SDK for Go is licensed under the [MIT](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/fileshares/armfileshares/LICENSE.txt) license.

<!-- LINKS -->
[azure_sdk_for_go_contributing]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[azure_sdk_for_go_contributing_developer_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#developer-guide
[azure_sdk_for_go_contributing_pull_requests]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#pull-requests
