## Guide for migrating to `sdk/**/arm**` from `services/**/mgmt/**`

This document is intended for users that are familiar with an older version of the Azure SDK For Go for management modules (`services/**/mgmt/**`) and wish to migrate their application to the next version of Azure resource management libraries (`sdk/**/arm**`)

**For users new to the Azure SDK For Go for resource management modules, please see the [README for 'sdk/armcore`](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/armcore) and the README for every individual package.**

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

In old version (`services/**/mgmt/**`), `autorest.Authorizer` is used in authentication process.

In new version (`sdk/**/arm**`), in order to provide a unified authentication based on Azure Identity for all Azure Go SDKs, the authentication mechanism has been re-designed and improved to offer a simpler interface.

To the show the code snippets for the change:

**In old version (`services/**/mgmt/**`)**

```go
authorizer, err := adal.NewServicePrincipalToken(oAuthToken, "<ClientId>", "<ClientSecret>", endpoint)
client := resources.NewGroupsClient("<SubscriptionId>")
client.Authorizer = authorizer
```        

**Equivalent in new version (`sdk/**/arm**`)**

```go
credential, err := azidentity.NewClientSecretCredential("<TenantId>", "<ClientId>", "<ClientSecret>", nil)
connection := armcore.NewDefaultConnection(credential, nil)
client := armresources.NewResourceGroupsClient(connection, "<SubscriptionId>")
```

For detailed information on the benefits of using the new authentication classes, please refer to [this page](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md)

### Error Handling

There are some minor changes in the error handling. 

- When there is an error in the SDK request, in the old version (`services/**/mgmt/**`), the return value will all be non-nil, and you can get the raw HTTP response from the response value. In the new version (`sdk/**/arm**`), the first return value will be empty and you need to cast the error to `HTTPResponse` interface to get the raw HTTP response. When the request is successful and there is no error returned, you will need to get the raw HTTP response in `RawResponse` property of the first return value.

**In old version (`services/**/mgmt/**`)**

```go
resp, err := resourceGroupsClient.CreateOrUpdate(context.Background(), resourceGroupName, resourceGroupParameters)
if err != nil {
	log.Fatalf("Status code: %d", resp.Response.Response.StatusCode)
}
```

**Equivalent in new version (`sdk/**/arm**`)**

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

**When there is no error in new version (`sdk/**/arm**`)**

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

In old version, if a request is a long-running operation, a struct `**Future` will be returned, which is an extension of the interface `azure.FutureAPI`. You need to invoke the `future.WaitForCompletionRef` to wait until it finishes.

In the new version, if a request is a long-running operation, the function name will start with `Begin` to indicate this function will return a poller struct which envelopes the interface `**Poller`, which is an extension of the interface `azcore.Poller`.

**In old version (`services/**/mgmt/**`)**

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

**Equivalent in new version (`sdk/**/arm**`)**

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

In old version, if a request is a paginated operation, a struct `**ResultPage` will be returned, which is a struct with some paging methods but no interfaces are defined regarding that.

In new version, if a request is a paginated operation, an interface `**Pager` will be returned, which is an extension of the interface `azcore.Pager`.

**In old version (`services/**/mgmt/**`)**

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

**Equivalent in new version (`sdk/**/arm**`)**

```go
pager := resourceGroupsClient.List(nil)
for pager.NextPage(context.Background()) {
    if err := pager.Err(); err != nil {
        log.Fatalf("failed to advance page: %v", err)
    }
    for _, rg := range pager.PageResponse().ResourceGroupListResult.Value {
        log.Printf("resource group ID: %s\n", *rg.ID)
    }
}
```

### Customized Policy

Because of adopting Azure Core which is a shared library across all Azure SDKs, there is also a minor change regarding how customized policy in configured.

In old version (`services/**/mgmt/**`), we use the `(autorest.Client).Sender`, `(autorest.Client).RequestInspector` and `(autorest.Client).ResponseInspector` properties in `github.com/Azure/go-autorest/autorest` module to provide customized interceptor for the HTTP traffic.

In new version (`sdk/**/arm**`), we use `(armcore.ConnectionOptions).PerCallPolicies` and `(armcore.ConnectionOptions).PerRetryPolicies` in `github.com/Azure/azure-sdk-for-go/sdk/armcore` module instead to inject customized policy to the pipeline.

### Custom HTTP Client

Similar to the customized policy, there are changes regarding how the custom HTTP client is configured as well. You can now use the `(armcore.ConnectionOptions).HTTPClient` option in `github.com/Azure/azure-sdk-for-go/sdk/armcore` module to use your own implementation of HTTP client and plug in what they need into the configuration.

**In old version (`services/**/mgmt/**`)** 
```go
httpClient := NewYourOwnHTTPClient{}
client := resources.NewGroupsClient("<SubscriptionId>")
client.Sender = &httpClient
```

**In new version (`sdk/**/arm**`)**

```go
httpClient := NewYourOwnHTTPClient{}
connection := armcore.NewConnection(credential, &armcore.ConnectionOptions{
    HTTPClient: &httpClient,
})
client := armresources.NewResourceGroupsClient(connection, "<SubscriptionId>")
```

## Need help?

If you have encountered an issue during migration, please file an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues) and make sure you add the "Preview" label to the issue