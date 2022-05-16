## Guide for migrating to `sdk/resourcemanager/**/arm**` from `services/**/mgmt/**`

This document is intended for users that are familiar with the previous version of the Azure SDK For Go for management modules (`services/**/mgmt/**`) and wish to migrate their application to the next version of Azure resource management libraries (`sdk/resourcemanager/**/arm**`)

**For users new to the Azure SDK For Go for resource management modules, please see the [README for 'sdk/azcore`](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azcore) and the README for every individual package.**

## Table of contents

* [Prerequisites](#prerequisites)
* [General Changes](#general-changes)
* [Breaking Changes](#breaking-changes)
    * [Authentication](#authentication)
    * [Error Handling](#error-handling)
    * [Long Running Operations](#long-running-operations)
    * [Pagination](#pagination)
    * [Customized Policy](#customized-policy)
    * [Custom HTTP Client](#custom-http-client)

## Prerequisites

- Go 1.18
- Latest version of resource management modules

## General Changes

The latest Azure SDK For Go for management modules is using the [Go Modules](https://github.com/golang/go/wiki/Modules) to manage the dependencies. We ship every RP as an individual module to create a more flexible user experience. Unlike previous version, the latest SDK only provides one API version in one module for each RP. The SDK stable version will always use the service's stable API version. You can search in [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go) to find the preview version and specify version num when using `go get` to retrive it (e.g., `go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault@v1.1.0-beta.1`).

## Breaking Changes

### Authentication

In the previous version (`services/**/mgmt/**`), `autorest.Authorizer` is used in authentication process.

In the latest version (`sdk/resourcemanager/**/arm**`), in order to provide a unified authentication based on Azure Identity for all Azure Go SDKs, the authentication mechanism has been re-designed and improved to offer a simpler interface.

To the show the code snippets for the change:

**Previous version (`services/**/mgmt/**`)**

```go
import "github.com/Azure/go-autorest/autorest/adal"
import "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
```

```go
authorizer, err := adal.NewServicePrincipalToken(oAuthToken, "<ClientId>", "<ClientSecret>", endpoint)
client := resources.NewGroupsClient("<SubscriptionId>")
client.Authorizer = authorizer
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
```

```go
credential, err := azidentity.NewClientSecretCredential("<TenantId>", "<ClientId>", "<ClientSecret>", nil)
client, err := armresources.NewResourceGroupsClient("<SubscriptionId>", credential, nil)
```

For detailed information on the benefits of using the new authentication types, please refer to [this page](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md)

### Error Handling

There are some minor changes in the error handling.

- When there is an error in the SDK request, in the previous version (`services/**/mgmt/**`), the return value will all be non-nil, and you can get the raw HTTP response from the response value. In the latest version (`sdk/resourcemanager/**/arm**`), the first return value will be empty and you need to convert the error to the `azcore.ResponseError` interface to get the raw HTTP response.
- In the latest version (`sdk/resourcemanager/**/arm**`), you can always get the raw HTTP response from request context regardless of request result.

**Previous version (`services/**/mgmt/**`)**

```go
resp, err := resourceGroupsClient.CreateOrUpdate(context.TODO(), resourceGroupName, resourceGroupParameters)
if err != nil {
    log.Fatalf("Status code: %d", resp.Response().StatusCode)
}
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore"
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
```

```go
var rawResponse *http.Response
ctx := context.TODO() // your context
ctxWithResp := runtime.WithCaptureResponse(ctx, &rawResponse)
resp, err := resourceGroupsClient.CreateOrUpdate(ctxWithResp, resourceGroupName, resourceGroupParameters, nil)
if err != nil {
    // with error, you can get RawResponse from context
    log.Printf("Status code: %d", rawResponse.StatusCode)
    var respErr *azcore.ResponseError
    if errors.As(err, &respErr) {
        // with error, you can also get RawResponse from error
        log.Fatalf("Status code: %d", respErr.RawResponse.StatusCode)
    } else {
        log.Fatalf("Other error: %+v", err)
    }
}
// without error, you can get RawResponse from context
log.Printf("Status code: %d", rawResponse.StatusCode)
```

### Long Running Operations

In the previous version, if a request is a long-running operation, a struct `**Future` will be returned, which is an extension of the interface `azure.FutureAPI`. You need to invoke the `future.WaitForCompletionRef` to wait until it finishes.

In the latest version, if a request is a long-running operation, the function name will start with `Begin` to indicate this function will return a poller type which contains the polling methods.

**Previous version (`services/**/mgmt/**`)**

```go
ctx := context.TODO() // your context
future, err := virtualMachinesClient.CreateOrUpdate(ctx, "<resource group name>", "<virtual machine name>", param)
if err != nil {
    log.Fatal(err)
}
if err := future.WaitForCompletionRef(ctx, virtualMachinesClient.Client); err != nil {
    log.Fatal(err)
}
vm, err := future.Result(virtualMachinesClient)
if err != nil {
    log.Fatal(err)
}
log.Printf("virtual machine ID: %v", *vm.ID)
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
poller, err := client.BeginCreateOrUpdate(context.TODO(), "<resource group name>", "<virtual machine name>", param, nil)
if err != nil {
    log.Fatal(err)
}
resp, err := poller.PollUntilDone(context.TODO(), nil)
if err != nil {
    log.Fatal(err)
}
log.Printf("virtual machine ID: %v", *resp.VirtualMachine.ID)
```

### Pagination

In the previous version, if a request is a paginated operation, a struct `**ResultPage` will be returned, which is a struct with some paging methods but no interfaces are defined regarding that.

In the latest version, if a request is a paginated operation, a struct `**Pager` will be returned that contains the paging methods.

**Previous version (`services/**/mgmt/**`)**

```go
pager, err := resourceGroupsClient.List(context.TODO(), "", nil)
if err != nil {
    log.Fatal(err)
}
for p.NotDone() {
    for _, v := range pager.Values() {
        log.Printf("resource group ID: %s\n", *rg.ID)
    }
    if err := pager.NextWithContext(context.TODO()); err != nil   {
        log.Fatal(err)
    }
}
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
ctx := context.TODO() // your context
pager := resourceGroupsClient.NewListPager(nil)
for pager.More() {
    nextResult, err := pager.NextPage(ctx)
    if err != nil {
        log.Fatalf("failed to advance page: %v", err)
    }
    for _, rg := range nextResult.Value {
        log.Printf("resource group ID: %s\n", *rg.ID)
    }
}
```

### Customized Policy

Because of adopting Azure Core which is a shared library across all Azure SDKs, there is also a minor change regarding how customized policy in configured.

In the previous version (`services/**/mgmt/**`), we use the `(autorest.Client).Sender`, `(autorest.Client).RequestInspector` and `(autorest.Client).ResponseInspector` properties in `github.com/Azure/go-autorest/autorest` module to provide customized interceptor for the HTTP traffic.

In latest version (`sdk/resourcemanager/**/arm**`), we use `arm.ClientOptions.PerCallPolicies` and `arm.ClientOptions.PerRetryPolicies` in `github.com/Azure/azure-sdk-for-go/sdk/azcore/arm` package instead to inject customized policy to the pipeline.

### Custom HTTP Client

Similar to the customized policy, there are changes regarding how the custom HTTP client is configured as well. You can now use the `arm.ClientOptions.Transport` option in `github.com/Azure/azure-sdk-for-go/sdk/azcore/arm` package to use your own implementation of HTTP client and plug in what they need into the configuration.  The HTTP client must implement the `policy.Transporter` interface.

**Previous version (`services/**/mgmt/**`)**

```go
import "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
```

```go
httpClient := NewYourOwnHTTPClient{}
client := resources.NewGroupsClient("<SubscriptionId>")
client.Sender = &httpClient
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
import "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
```

```go
httpClient := NewYourOwnHTTPClient{}
options := &arm.ClientOptions{
    ClientOptions: policy.ClientOptions{
        Transport: &httpClient,
    },
}
client, err := armresources.NewResourceGroupsClient("<SubscriptionId>", credential, options)
```

## Need help?

If you have encountered an issue during migration, please file an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues) and make sure you add the "Preview" label to the issue
