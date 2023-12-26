# Examples using fakes

The examples in this module demonstrate using the `fake` subpackage for creating unit tests.
While the examples use `armcompute`, the patterns are applicable to any module with a `fake` subpackage.

## Fakes

The `fake` package found in most modules provides implementations for fake servers that can be used for testing.

To create a fake server, declare an instance of the required fake server type(s).

```go
import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5/fake"

myFakeVirtualMachinesServer := fake.VirtualMachinesServer{}
```

Next, provide func implementations for the client methods you wish to fake.
The named return variables can be used to simplify return value construction.

```go
import azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"

myFakeVirtualMachinesServer.Get = func(ctx context.Context, resourceGroupName string, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
	// TODO: resp.SetResponse(/* your fake armcompute.VirtualMachinesClientGetResponse response */)
	return
}
```

You connect the fake server to a client instance during its construction through the optional transport.

Use the `TokenCredential` type from `azcore/fake` to create a fake `azcore.TokenCredential`.

```go
import azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"

client, err := armcompute.VirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
	ClientOptions: azcore.ClientOptions{
		Transport: fake.NewVirtualMachinesServerTransport(&myFakeVirtualMachinesServer),
	},
})
```

Calling methods on the client will pass the provided values to the matching fake implementation.
The values can be arbitrary, including the zero-value for any/all parameters.

```go
resp, err := client.Get(context.TODO(), "fake-resource-group", "fake-vm", nil)
```

The values returned from client method calls are defined in the fake.
