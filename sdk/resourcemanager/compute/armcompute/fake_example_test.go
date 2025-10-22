//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8/fake"
)

func ExampleVirtualMachinesServer() {
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

			// use resp to set the desired response
			resp.SetResponse(http.StatusOK, vmResp, nil)

			// to simulate the failure case, use errResp
			// errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armcompute.NewVirtualMachinesClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewVirtualMachinesServerTransport(&fakeVirtualMachinesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.Get(context.TODO(), "fakeResourceGroup", "fakeVM", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.ID)

	// APIs that haven't been faked will return an error
	_, err = client.BeginDeallocate(context.TODO(), "fakeResourceGroup", "fakeVM", nil)

	fmt.Println(err.Error())

	// Output:
	// /fake/resource/id
	// fake for method BeginDeallocate not implemented
}
