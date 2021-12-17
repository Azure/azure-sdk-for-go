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

The last 2 released stable versions of Go are supported.

## General Changes

The latest Azure SDK For Go for management modules is using the [Go Modules](https://github.com/golang/go/wiki/Modules) to manage the dependencies. We ship every RP as an individual module to create a more flexible user experience.

## Breaking Changes

### Authentication

In the previous version (`services/**/mgmt/**`), `autorest.Authorizer` is used in authentication process.

In the latest version (`sdk/resourcemanager/**/arm**`), in order to provide a unified authentication based on Azure Identity for all Azure Go SDKs, the authentication mechanism has been re-designed and improved to offer a simpler interface.

To the show the code snippets for the change:

**Previous version (`services/**/mgmt/**`)**

```go
authorizer, err := adal.NewServicePrincipalToken(oAuthToken, "<ClientId>", "<ClientSecret>", endpoint)
client := resources.NewGroupsClient("<SubscriptionId>")
client.Authorizer = authorizer
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
credential, err := azidentity.NewClientSecretCredential("<TenantId>", "<ClientId>", "<ClientSecret>", nil)
client := armresources.NewResourceGroupsClient("<SubscriptionId>", credential, nil)
```

For detailed information on the benefits of using the new authentication types, please refer to [this page](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md)

### Error Handling

There are some minor changes in the error handling.

- When there is an error in the SDK request, in the previous version (`services/**/mgmt/**`), the return value will all be non-nil, and you can get the raw HTTP response from the response value. In the latest version (`sdk/resourcemanager/**/arm**`), the first return value will be empty and you need to convert the error to the `azcore.HTTPResponse` interface to get the raw HTTP response. When the request is successful and there is no error returned, you can get the raw HTTP response in `RawResponse` property of the first return value.

**Previous version (`services/**/mgmt/**`)**

```go
resp, err := resourceGroupsClient.CreateOrUpdate(context.Background(), resourceGroupName, resourceGroupParameters)
if err != nil {
	log.Fatalf("Status code: %d", resp.Response.Response.StatusCode)
}
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
resp, err := resourceGroupsClient.CreateOrUpdate(context.Background(), resourceGroupName, resourceGroupParameters, nil)
if err != nil {
    var respErr azcore.HTTPResponse
    if errors.As(err, &respErr) {
        log.Fatalf("Status code: %d", respErr.RawResponse().StatusCode)
    }
	log.Fatalf("Other error: %+v", err)
}
```

**When there is no error in latest version (`sdk/resourcemanager/**/arm**`)**

```go
resp, err := resourceGroupsClient.CreateOrUpdate(context.Background(), resourceGroupName, resourceGroupParameters, nil)
if err != nil {
    var respErr azcore.HTTPResponse
    if errors.As(err, &respErr) {
        log.Fatalf("Status code: %d", respErr.RawResponse().StatusCode)
    }
    log.Fatalf("Other error: %+v", err)
}
log.Printf("Status code: %d", resp.RawResponse.StatusCode)
```

### Long Running Operations

In the previous version, if a request is a long-running operation, a struct `**Future` will be returned, which is an extension of the interface `azure.FutureAPI`. You need to invoke the `future.WaitForCompletionRef` to wait until it finishes.

In the latest version, if a request is a long-running operation, the function name will start with `Begin` to indicate this function will return a poller type which contains the polling methods.

**Previous version (`services/**/mgmt/**`)**

```go
future, err := virtualMachinesClient.CreateOrUpdate(context.Background(), "<resource group name>", "<virtual machine name>", param)
if err != nil {
	log.Fatal(err)
}
if err := future.WaitForCompletionRef(context.Background(), virtualMachinesClient.Client); err != nil {
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
poller, err := client.BeginCreateOrUpdate(context.Background(), "<resource group name>", "<virtual machine name>", param, nil)
if err != nil {
	log.Fatal(err)
}
resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
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
pager, err := resourceGroupsClient.List(context.Background(), "", nil)
if err != nil {
    log.Fatal(err)
}
for p.NotDone() {
    for _, v := range pager.Values() {
        log.Printf("resource group ID: %s\n", *rg.ID)
    }
    if err := pager.NextWithContext(context.Background()); err != nil   {
        log.Fatal(err)
    }
}
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
pager := resourceGroupsClient.List(nil)
for pager.NextPage(context.Background()) {
    for _, rg := range pager.PageResponse().ResourceGroupListResult.Value {
        log.Printf("resource group ID: %s\n", *rg.ID)
    }
}
if err := pager.Err(); err != nil {
    log.Fatalf("failed to advance page: %v", err)
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
httpClient := NewYourOwnHTTPClient{}
client := resources.NewGroupsClient("<SubscriptionId>")
client.Sender = &httpClient
```

**Latest version (`sdk/resourcemanager/**/arm**`)**

```go
httpClient := NewYourOwnHTTPClient{}
options := &arm.ClientOptions{
    ClientOptions: policy.ClientOptions{
        Transport: &httpClient,
    },
}
client := armresources.NewResourceGroupsClient("<SubscriptionId>", credential, options)
```

## Need help?

If you have encountered an issue during migration, please file an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues) and make sure you add the "Preview" label to the issue
