# Azure Batch client module for Go

Azure Batch allows users to run large-scale parallel and high-performance computing (HPC) batch jobs efficiently in Azure.

Use this module to:

- Create and manage Batch jobs and tasks
- View and perform operations on nodes in a Batch pool

## Getting started

### Install the module

Install the `azbatch` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/batch/azbatch
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

### Prerequisites

- Go, version 1.18 or higher - [Install Go](https://go.dev/doc/install)
- Azure subscription - [Create a free account](https://azure.microsoft.com/free)
- A Batch account with a linked Azure Storage account. You can create the accounts by using any of the following methods: [Azure CLI](https://learn.microsoft.com/azure/batch/quick-create-cli) | [Azure portal](https://learn.microsoft.com/azure/batch/quick-create-portal) | [Bicep](https://learn.microsoft.com/azure/batch/quick-create-bicep) | [ARM template](https://learn.microsoft.com/azure/batch/quick-create-template) | [Terraform](https://learn.microsoft.com/azure/batch/quick-create-terraform).

### Authenticate the client

Azure Batch integrates with Microsoft Entra ID for identity-based authentication of requests. You can use role-based access control (RBAC) to grant access to your Azure Batch resources to users, groups, or applications. The [Azure Identity module](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) provides types that implement Microsoft Entra ID authentication.

## Key concepts

[Azure Batch Overview](https://learn.microsoft.com/azure/batch/batch-technical-overview)

## Examples

See the [package documentation][pkgsite] for code samples.

## Troubleshooting

Please see [Troubleshooting common batch issues](https://learn.microsoft.com/troubleshoot/azure/hpc/batch/welcome-hpc-batch).

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Key Vault.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

resp, err = client.CreateJob(context.TODO(), jobContent, nil)
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
azlog.SetListener(func (_ azlog.Event, msg string) {
    fmt.Println(msg)
})

// Includes only requests and responses in logs
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

### Accessing `http.Response`

You can access the `http.Response` returned by Azure Batch to any client method using `runtime.WithCaptureResponse`:

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

var response *http.Response
ctx := runtime.WithCaptureResponse(context.TODO(), &response)
resp, err = client.CreateJob(ctx, jobContent, nil)
if err != nil {
    // TODO: handle error
}
// TODO: do something with response
```

## Contributing

This project welcomes contributions and suggestions.
Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution.
For details, visit [Contributor License Agreements](https://opensource.microsoft.com/cla/).

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment).
Simply follow the instructions provided by the bot.
You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

[pkgsite]: https://aka.ms/azsdk/go/azbatch