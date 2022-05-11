
# Azure Go Management SDK Guideline

Azure Go management SDK follows the [new Azure SDK guidelines](https://azure.github.io/azure-sdk/general_introduction.html), try to create easy-to-use APIs that are idiomatic, compatible, and dependable.

You can find the full list of management modules [here](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk).

In this guideline, we will give some instructions about the API usage pattern as well as trouble shooting method. For those are new to management Go SDK, please refer to [quickstart](./new-version-quickstart.md). For those migrate from older versions of management Go SDK, please refer to [migration guide](https://aka.ms/azsdk/go/mgmt/migration).

## Pageable Operations

### General usage

Pageable operations return final data over multiple GET requests. Each GET will receive a page of data consisting of a slice of items. You need to use New*Pager to create a pager helper for all pageable operations. With the returned *runtime.Pager[T], you could fetch pages and determine if there are more pages to fetch. For examples:

```go
pager := rgClient.NewListPager(nil)
var resourceGroups []*armresources.ResourceGroup
for pager.More() {
    nextResult, err := pager.NextPage(ctx)
    if err != nil {
        // handle error...
    }
    if nextResult.ResourceGroupListResult.Value != nil {
        resourceGroups = append(resourceGroups, nextResult.ResourceGroupListResult.Value...)
    }
}
// dealing with `resourceGroups`
```

> NOTE: No IO calls are made until the NextPage() method is invoked. The read consistency across pages is determined by the service implement.

### Item iterator

If you are not care about the underlaying detail about the pageable operation, you could use the following generic utility to create a per-item iterator for all pageable operation.

***Item iterator utility***
```go
type PageConstraint[TItem any] interface {
	Items() []*TItem
}

type Iterator[TItem any, TPage PageConstraint[TItem]] struct {
	pager *runtime.Pager[TPage]
	cur   []*TItem
	index int
}

func (iter *Iterator[TItem, TPage]) More() bool {
	return iter.pager.More() || iter.index < len(iter.cur)
}

func (iter *Iterator[TItem, TPage]) NextItem(ctx context.Context) (*TItem, error) {
	if iter.index == len(iter.cur) && !iter.pager.More() {
		return nil, errors.New("no more items")
	}
	if iter.cur == nil || iter.index == len(iter.cur) {
		// first page or page exhausted
		page, err := iter.pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		iter.cur = page.Items()
		iter.index = 0
	}
	item := iter.cur[iter.index]
	// advance item
	iter.index++
	return item, nil
}

func NewIterator[TItem any, TPage PageConstraint[TItem]](pager *runtime.Pager[TPage]) *Iterator[TItem, TPage] {
	return &Iterator[TItem, TPage]{
		pager: pager,
	}
}
```

***Usage***
```go
iter := NewIterator[*armresources.ResourceGroup](rgClient.NewListPager(nil))
for iter.More() {
    rg, err := iter.NextItem(ctx)
    if err != nil {
        // handle error...
    }
    // dealing with `rg`
}
```

### Reference

For more information, you could refer to [design guidelines of Paging](https://azure.github.io/azure-sdk/golang_introduction.html#methods-returning-collections-paging) and [API reference of pager](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime#Pager).

## Long-Running Operations

### General usage

Some operations can take a long time to complete. Azure introduces the long-running operations (LROs) to do such operations asynchronously. You need to use Begin* to start an LRO. It will return a poller that can used to keep polling for the result until LRO is done. For examples:

```go
ctx := context.TODO()
poller, err := client.BeginCreate(ctx, "resource_identifier", "additonal_parameter", nil)
if err != nil {
    // handle error...
}
resp, err = poller.PollUntilDone(ctx, 5 * time.Second)
if err != nil {
    // handle error...
}
// dealing with `resp`
```

> NOTE: You will need to pass a polling interval to ```PollUntilDone``` and tell the poller how often it should try to get the status. This number is usually small but it's best to consult the [Azure service documentation](https://docs.microsoft.com/azure/?product=featured) on best practices and recommended intervals for your specific use cases.

### Resume Tokens

Pollers provide the ability to serialize their state into a "resume token" which can be used by another process to recreate the poller. For example:

```go
ctx := context.TODO()
poller, err := client.BeginCreate(ctx, "resource_identifier", "additonal_parameter", nil)
if err != nil {
    // handle error...
}
token, err := poller.ResumeToken()
if err != nil {
    // handle error...
}

// ... 

// recreate the poller from the token
poller, err = client.BeginCreate(ctx, "", "", &armresources.ResourceGroupsClientBeginCreateOptions{
    ResumeToken: token,
})
resp, err = poller.PollUntilDone(ctx, 5 * time.Second)
if err != nil {
    // handle error...
}
// dealing with `resp`
```

> NOTE: A token can only be obtained for a poller that's in a non-terminal state. Each time you call poller.Poll(), you should get a new token to handle poller state change.

### Synchronized wrapper

If you are not care about the underlaying detail about the LRO, you could use the following generic utility to create an synchronized wrapper for all LRO.

***Synchronized wrapper utility***
```go
type OperationWaiter[TResult any] struct {
    poller *runtime.Poller[TResult]
    err    error
}

func (ow OperationWaiter[TResult]) Wait(ctx context.Context, freq time.Duration) (TResult, error) {
    if ow.err != nil {
        return *new(TResult), ow.err
    }
    return ow.poller.PollUntilDone(ctx, freq)
}

func NewOperationWaiter[TResult any](poller *runtime.Poller[TResult], err error) OperationWaiter[TResult] {
    return OperationWaiter[TResult]{poller: poller, err: err}
}
```

***Usage***
```go
resp, err := NewOperationWaiter(BeginCreate(ctx, "resource_identifier", "additonal_parameter", nil)).Wait(ctx, time.Second)
// dealing with `resp`
```
### Reference

For more information, you could refer to [design guidelines of LRO](https://azure.github.io/azure-sdk/golang_introduction.html#methods-invoking-long-running-operations) and [API reference of poller](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime#Poller).

## Client Options

### Request Retry Policy
The SDK provides a baked in retry policy for failed requests with default values that can be configured by `arm.ClientOptions.Retry`. For example:

```go
rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, credential,
    &arm.ClientOptions{
        ClientOptions: policy.ClientOptions{
            Retry: policy.RetryOptions{
                // retry for 5 times
                MaxRetries: 5,
            },
        },
    },
)
```

### Customized Policy

You can use `arm.ClientOptions.PerCallPolicies` and `arm.ClientOptions.PerRetryPolicies` option to inject customized policy to the pipeline. You can refer to `azcore` [document](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore) for further information.

### Custom HTTP Client

You can use `arm.ClientOptions.Transport` to set your own implementation of HTTP client. The HTTP client must implement the `policy.Transporter` interface. For example:

```go
// your own implementation of HTTP client
httpClient := NewYourOwnHTTPClient{}
rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, credential,
    &arm.ClientOptions{
        ClientOptions: policy.ClientOptions{
            Transport: &httpClient,
        },
    },
)
```

### Reference

More client options can be found [here](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore/policy/policy.go).


## Troubleshooting

### Logging

The SDK uses the classification-based logging implementation in `azcore`. To enable console logging for all SDK modules, set AZURE_SDK_GO_LOGGING to all. 

Use `LogOption` to configure log behavior. For example:
```go
rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, credential,
    &arm.ClientOptions{
        ClientOptions: policy.ClientOptions{
            Logging: policy.LogOptions{
                // include HTTP body for log
                IncludeBody: true,
            },
        },
    },
)
```

Use the `azcore/log` package to control log event and write log to the desired location. For example:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

// print log output to stdout
azlog.SetListener(func(event azlog.Event, s string) {
    fmt.Println(s)
})

// include only azidentity credential logs
azlog.SetEvents(azidentity.EventAuthentication)
```

### Raw HTTP response
When there is an error in the SDK request, you need to convert the error to the `azcore.ResponseError` interface to get the raw HTTP response. When the request is successful, you can get the raw HTTP response from request context.

```go
var rawResponse *http.Response
ctxWithResp := runtime.WithCaptureResponse(context.TODO(), &rawResponse)
resp, err := resourceGroupsClient.CreateOrUpdate(ctxWithResp, resourceGroupName, resourceGroupParameters, nil)
if err != nil {
    var respErr *azcore.ResponseError
    if errors.As(err, &respErr) {
        log.Fatalf("Status code: %d", respErr.RawResponse.StatusCode)
    } else {
        log.Fatalf("Other error: %+v", err)
    }
}
log.Printf("Status code: %d", rawResponse.StatusCode)
```

> NOTE: runtime is import from **github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime**


## Need help?

-   File an issue via [Github
    Issues](https://github.com/Azure/azure-sdk-for-go/issues)
-   Check [previous
    questions](https://stackoverflow.com/questions/tagged/azure+go)
    or ask new ones on StackOverflow using azure and Go tags.

## Contributing

For details on contributing to this repository, see the contributing
guide.

This project welcomes contributions and suggestions. Most contributions
require you to agree to a Contributor License Agreement (CLA) declaring
that you have the right to, and actually do, grant us the rights to use
your contribution. For details, visit <https://cla.microsoft.com>.

When you submit a pull request, a CLA-bot will automatically determine
whether you need to provide a CLA and decorate the PR appropriately
(e.g., label, comment). Simply follow the instructions provided by the
bot. You will only need to do this once across all repositories using
our CLA.

This project has adopted the Microsoft Open Source Code of Conduct. For
more information see the Code of Conduct FAQ or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any questions or comments.
