# Azure OpenTelemetry Adapter Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/go%20-%20azotel%20-%20ci?branchName=main)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=6176&branchName=main)

The `azotel` module is used to connect an instance of OpenTelemetry's `TracerProvider` to an Azure SDK client.

## Getting started

**NOTE: this module requires Go 1.19 or later**

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

To add the latest version to your `go.mod` file, execute the following command.

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel
```

General documentation and examples can be found on [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel).

## Using the adapter

Once you have created an OpenTelemetry `TracerProvider`, you connect it to an Azure SDK client via its `ClientOptions`.

```go
options := azcore.ClientOptions{}
options.TracingProvider = azotel.NewTracingProvider(otelProvider, nil)
```

## Contributing
This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.
