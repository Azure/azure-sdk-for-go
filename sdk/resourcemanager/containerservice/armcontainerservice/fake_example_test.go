// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcontainerservice_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v8/fake"
)

func ExampleManagedClustersServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeManagedClustersServer := fake.ManagedClustersServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the VirtualMachinesClient.Get() API.
		Get: func(ctx context.Context, resourceGroupName, resourceName string, options *armcontainerservice.ManagedClustersClientGetOptions) (resp azfake.Responder[armcontainerservice.ManagedClustersClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, resourceName, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct the response type, populating fields as required
			clusterResp := armcontainerservice.ManagedClustersClientGetResponse{}
			clusterResp.ID = to.Ptr("/fake/resource/id")

			// use resp to set the desired response
			resp.SetResponse(http.StatusOK, clusterResp, nil)

			// to simulate the failure case, use errResp
			// errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	client, err := armcontainerservice.NewManagedClustersClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewManagedClustersServerTransport(&fakeManagedClustersServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.Get(context.TODO(), "fakeResourceGroup", "fakeResource", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.ID)

	// APIs that haven't been faked will return an error
	_, err = client.GetAccessProfile(context.TODO(), "fakeResourceGroup", "fakeResource", "fakeRole", nil)

	fmt.Println(err.Error())

	// Output:
	// /fake/resource/id
	// fake for method GetAccessProfile not implemented
}
