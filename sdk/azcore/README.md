# Azure Core Client Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/azcore)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/go%20-%20azcore%20-%20ci?branchName=main)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=1843&branchName=main)
[![Code Coverage](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1843/main)](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1843/main)

The `azcore` module provides a set of common interfaces and types for Go SDK client modules.
These modules follow the [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html).

## Getting started

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Typically, you will not need to explicitly install `azcore` as it will be installed as a client module dependency.
To add the latest version to your `go.mod` file, execute the following command.

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/azcore
```

General documentation and examples can be found on [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore).

## Key concepts

The main shared concepts of `azcore` include:

- Accessing HTTP response details for the returned model of any SDK client operation.
- Polling long-running operations (LROs), via `Poller`.
- Pagenation for pageable apis, via `Pager`.
- Abstractions for Azure SDK credentials (`TokenCredential`).
- HTTP pipeline and HTTP policies such as retry and logging, which are configurable via service client specific options.

### Long Running Operations

Some operations take a long time to complete and require polling for their status. Methods starting long-running operations return `Poller[T]` types.

You can intermittently poll whether the operation has finished by using the `Poll()` method on the returned `Poller[T]`.  Alternatively, if you just want to wait until the operation completes, you can use `PollUntilDone()`.


#### PollUntilDone Example
```go
client, err := exampleservice.NewClient()

poller, err := client.SomeLongRunningAPI()
if err != nil {
    panic(err)
}

pollResp, err := poller.PollUntilDone(context.TODO())
if err != nil {
    panic(err)
}
// Do something with the result
//  i.e. fmt.Printf("Retrieved value %s", *pollResp.Value.Property)

```

#### Poll Example
```go
client, err := exampleservice.NewClient()

poller, err := client.SomeLongRunningAPI()
if err != nil {
    panic(err)
}

for {
    _, err := poller.Poll()
    if err != nil {
        panic(err)
    }
    if poller.Done(){
        result, err := poller.Result(context.TODO())
        if err != nil {
            panic(err)
        }
        //Do Something with the result
        break
    }
}
```
### Consuming Service Methods Returning `Pager`

If a service call returns multiple values in pages, it would return a `runtime.Pager` as a result.  You can iterate over pages using the following coding pattern:

```go
client, err := exampleservice.NewClient()

pager := client.SomePageableAPI()
for pager.More() {
    resp, err := pager.NextPage(context.TODO())
    if err != nil {
        panic(err)
    }
    for _, value := range resp.Values {
    }
}
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
