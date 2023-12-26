//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fakes_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5/fake"
)

func ExampleVirtualMachinesServer_Get() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.Get() API.
		Get: func(ctx context.Context, resourceGroupName, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct the response type, populating fields as required
			vmResp := armcompute.VirtualMachinesClientGetResponse{}
			vmResp.ID = to.Ptr("/fake/resource/id")
			vmResp.Name = &vmName

			// use resp to set the desired response with the applicable HTTP status code.
			// see the doc comments of fake.VirtualMachinesServer.Get for the list of possible HTTP status codes.
			resp.SetResponse(http.StatusOK, vmResp, nil)

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.Get(context.TODO(), "fake-resource-group", "fake-vm", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.ID)
	fmt.Println(*resp.Name)

	// Output:
	// /fake/resource/id
	// fake-vm
}

func ExampleVirtualMachinesServer_Get_responseError() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.Get() API.
		Get: func(ctx context.Context, resourceGroupName, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// to simulate the failure case, use errResp.
			// the API call will return an *azcore.ResponseError with the specified values
			errResp.SetResponseError(http.StatusBadRequest, "FakeErrorCode")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	_, err = client.Get(context.TODO(), "fake-resource-group", "fake-vm", nil)

	var respErr *azcore.ResponseError
	errors.As(err, &respErr)
	fmt.Println(respErr.Error())

	// Output:
	// GET https://management.azure.com/subscriptions/fake-subscription-id/resourceGroups/fake-resource-group/providers/Microsoft.Compute/virtualMachines/fake-vm
	// --------------------------------------------------------------------------------
	// RESPONSE 400: 400 Bad Request
	// ERROR CODE: FakeErrorCode
	// --------------------------------------------------------------------------------
	// Response contained no body
	// --------------------------------------------------------------------------------
}

func ExampleVirtualMachinesServer_Get_error() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.Get() API.
		Get: func(ctx context.Context, resourceGroupName, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// to simulate the failure case, use errResp.
			// the API call will return the specified error.
			errResp.SetError(errors.New("fake-error"))

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	_, err = client.Get(context.TODO(), "fake-resource-group", "fake-vm", nil)

	fmt.Println(err.Error())

	// Output:
	// fake-error
}

func ExampleVirtualMachinesServer_Get_not_implemented() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	// in this example, no APIs have been faked.
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// APIs that haven't been faked will return an error
	_, err = client.Get(context.TODO(), "fake-resource-group", "fake-vm", nil)

	fmt.Println(err.Error())

	// Output:
	// fake for method Get not implemented
}

func ExampleVirtualMachinesServer_BeginCreateOrUpdate() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.BeginCreateOrUpdate() API.
		BeginCreateOrUpdate: func(ctx context.Context, resourceGroupName, vmName string, parameters armcompute.VirtualMachine, options *armcompute.VirtualMachinesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.VirtualMachinesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, parameters, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// responses for a long-running operation are composed of optionally one or more non-terminal
			// responses and a terminal response which can either be success or failure.

			// add two non-terminal responses to simulate the long-runninig operation being in progress.
			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.AddNonTerminalResponse(http.StatusCreated, nil)

			// create a fake item to return from the Poller[T]
			vmResp := armcompute.VirtualMachinesClientCreateOrUpdateResponse{}
			vmResp.ID = to.Ptr("/fake/resource/id")
			vmResp.Name = &vmName

			// set the terminal success response that terminates the long-running operation, returning the fake item.
			resp.SetTerminalResponse(http.StatusOK, vmResp, nil)

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	poller, err := client.BeginCreateOrUpdate(context.TODO(), "fake-resource-group", "fake-vm", armcompute.VirtualMachine{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// poll until the fake long-running operation completes.
	// we use the shortest allowed frequency so the example completes in a reasonable amount of time.
	resp, err := poller.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})

	fmt.Println(*resp.ID)
	fmt.Println(*resp.Name)

	// Output:
	// /fake/resource/id
	// fake-vm
}

func ExampleVirtualMachinesServer_BeginCreateOrUpdate_responseError() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.BeginCreateOrUpdate() API.
		BeginCreateOrUpdate: func(ctx context.Context, resourceGroupName, vmName string, parameters armcompute.VirtualMachine, options *armcompute.VirtualMachinesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.VirtualMachinesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, parameters, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// setting the error response will return an *azcore.ResponseError, no polling will be performed
			errResp.SetResponseError(http.StatusBadRequest, "FakeBadRequest")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	_, err = client.BeginCreateOrUpdate(context.TODO(), "fake-resource-group", "fake-vm", armcompute.VirtualMachine{}, nil)

	var respErr *azcore.ResponseError
	errors.As(err, &respErr)
	fmt.Println(respErr.Error())

	// Output:
	// PUT https://management.azure.com/subscriptions/fake-subscription-id/resourceGroups/fake-resource-group/providers/Microsoft.Compute/virtualMachines/fake-vm
	// --------------------------------------------------------------------------------
	// RESPONSE 400: 400 Bad Request
	// ERROR CODE: FakeBadRequest
	// --------------------------------------------------------------------------------
	// Response contained no body
	// --------------------------------------------------------------------------------
}

func ExampleVirtualMachinesServer_BeginCreateOrUpdate_poller_responseError() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.BeginCreateOrUpdate() API.
		BeginCreateOrUpdate: func(ctx context.Context, resourceGroupName, vmName string, parameters armcompute.VirtualMachine, options *armcompute.VirtualMachinesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.VirtualMachinesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, parameters, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// responses for a long-running operation are composed of optionally one or more non-terminal
			// responses and a terminal response which can either be success or failure.

			// add two non-terminal responses to simulate the long-runninig operation being in progress.
			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.AddNonTerminalResponse(http.StatusCreated, nil)

			// set the terminal error response that terminates the long-running operation, returning an *azcore.ResponseError.
			resp.SetTerminalError(http.StatusConflict, "FakeStatusCreationConflict")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	poller, err := client.BeginCreateOrUpdate(context.TODO(), "fake-resource-group", "fake-vm", armcompute.VirtualMachine{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// poll until the fake long-running operation completes.
	// we use the shortest allowed frequency so the example completes in a reasonable amount of time.
	_, err = poller.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})

	var respErr *azcore.ResponseError
	errors.As(err, &respErr)
	fmt.Println(respErr.Error())

	// Output:
	// GET https://management.azure.com/subscriptions/fake-subscription-id/resourceGroups/fake-resource-group/providers/Microsoft.Compute/virtualMachines/fake-vm/get/fake/status
	// --------------------------------------------------------------------------------
	// RESPONSE 409: 409 Conflict
	// ERROR CODE: FakeStatusCreationConflict
	// --------------------------------------------------------------------------------
	// Response contained no body
	// --------------------------------------------------------------------------------
}

func ExampleVirtualMachinesServer_NewListPager() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.NewListPager() API.
		NewListPager: func(resourceGroupName string, options *armcompute.VirtualMachinesClientListOptions) (resp azfake.PagerResponder[armcompute.VirtualMachinesClientListResponse]) {
			// the values of resourceGroupName and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct one or more pages to return.
			// each page contains one or more fake items.

			page1 := armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{
							ID: to.Ptr("/fake/resource/id-1"),
						},
						{
							ID: to.Ptr("/fake/resource/id-2"),
						},
					},
				},
			}

			page2 := armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{
							ID: to.Ptr("/fake/resource/id-3"),
						},
						{
							ID: to.Ptr("/fake/resource/id-4"),
						},
					},
				},
			}

			// now add the pages to the response
			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	pager := client.NewListPager("fake-resource-group", nil)

	// iterate over the returned pages, printing the fake VM ID for each entry
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		for _, vm := range page.Value {
			fmt.Println(*vm.ID)
		}
	}

	// Output:
	// /fake/resource/id-1
	// /fake/resource/id-2
	// /fake/resource/id-3
	// /fake/resource/id-4
}

func ExampleVirtualMachinesServer_NewListPager_responseError() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.NewListPager() API.
		NewListPager: func(resourceGroupName string, options *armcompute.VirtualMachinesClientListOptions) (resp azfake.PagerResponder[armcompute.VirtualMachinesClientListResponse]) {
			// the values of resourceGroupName and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct one or more pages to return.
			// each page contains one or more fake items.

			page1 := armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{
							ID: to.Ptr("/fake/resource/id-1"),
						},
						{
							ID: to.Ptr("/fake/resource/id-2"),
						},
					},
				},
			}

			// now add the pages to the response
			resp.AddPage(http.StatusOK, page1, nil)

			// set an error to be returned after the first page
			resp.AddResponseError(http.StatusInternalServerError, "FakeServerError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	pager := client.NewListPager("fake-resource-group", nil)

	// iterate over the returned pages, printing the fake VM ID for each entry
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			var respErr *azcore.ResponseError
			if !errors.As(err, &respErr) {
				log.Fatal(err)
			}
			fmt.Println(respErr.Error())

			break
		}

		for _, vm := range page.Value {
			fmt.Println(*vm.ID)
		}
	}

	// Output:
	// /fake/resource/id-1
	// /fake/resource/id-2
	// GET https://management.azure.com/subscriptions/fake-subscription-id/resourceGroups/fake-resource-group/providers/Microsoft.Compute/virtualMachines/fake_page_1
	// --------------------------------------------------------------------------------
	// RESPONSE 500: 500 Internal Server Error
	// ERROR CODE: FakeServerError
	// --------------------------------------------------------------------------------
	// Response contained no body
	// --------------------------------------------------------------------------------
}

func ExampleVirtualMachinesServer_NewListPager_intermediate_responseError() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeVirtualMachinesServer := fake.VirtualMachinesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.NewListPager() API.
		NewListPager: func(resourceGroupName string, options *armcompute.VirtualMachinesClientListOptions) (resp azfake.PagerResponder[armcompute.VirtualMachinesClientListResponse]) {
			// the values of resourceGroupName and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct one or more pages to return.
			// each page contains one or more fake items.

			page1 := armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{
							ID: to.Ptr("/fake/resource/id-1"),
						},
						{
							ID: to.Ptr("/fake/resource/id-2"),
						},
					},
				},
			}

			// add the page to the response
			resp.AddPage(http.StatusOK, page1, nil)

			// set an error to be returned after the first page
			resp.AddResponseError(http.StatusRequestTimeout, "FakeServerTimeout")

			page2 := armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{
							ID: to.Ptr("/fake/resource/id-3"),
						},
						{
							ID: to.Ptr("/fake/resource/id-4"),
						},
					},
				},
			}

			// add the page to the response
			resp.AddPage(http.StatusOK, page2, nil)

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("fake-subscription-id", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	pager := client.NewListPager("fake-resource-group", nil)

	// iterate over the returned pages, printing the fake VM ID for each entry
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			var respErr *azcore.ResponseError
			if !errors.As(err, &respErr) {
				log.Fatal(err)
			}
			fmt.Println(respErr.Error())
			continue
		}

		for _, vm := range page.Value {
			fmt.Println(*vm.ID)
		}
	}

	// Output:
	// /fake/resource/id-1
	// /fake/resource/id-2
	// GET https://management.azure.com/subscriptions/fake-subscription-id/resourceGroups/fake-resource-group/providers/Microsoft.Compute/virtualMachines/fake_page_1
	// --------------------------------------------------------------------------------
	// RESPONSE 408: 408 Request Timeout
	// ERROR CODE: FakeServerTimeout
	// --------------------------------------------------------------------------------
	// Response contained no body
	// --------------------------------------------------------------------------------
	//
	// /fake/resource/id-3
	// /fake/resource/id-4
}
