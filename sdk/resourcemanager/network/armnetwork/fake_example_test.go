// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7/fake"
)

func ExampleInterfacesServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeInterfacesServer := fake.InterfacesServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the InterfacesClient.Get() API.
		Get: func(ctx context.Context, resourceGroupName string, networkInterfaceName string, options *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, vmName, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct the response type, populating fields as required
			interfaceResp := armnetwork.InterfacesClientGetResponse{}
			interfaceResp.ID = to.Ptr("/fake/resource/id")

			// use resp to set the desired response
			resp.SetResponse(http.StatusOK, interfaceResp, nil)

			// to simulate the failure case, use errResp
			// errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armnetwork.NewInterfacesClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewInterfacesServerTransport(&fakeInterfacesServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.Get(context.TODO(), "fakeResourceGroup", "fakeInterface", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.ID)

	// APIs that haven't been faked will return an error
	_, err = client.BeginDelete(context.TODO(), "fakeResourceGroup", "fakeInterface", nil)

	fmt.Println(err.Error())

	// Output:
	// /fake/resource/id
	// fake for method BeginDelete not implemented
}
