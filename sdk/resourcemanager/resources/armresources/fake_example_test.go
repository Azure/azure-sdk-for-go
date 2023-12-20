//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/fake"
)

func ExampleServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeServer := fake.Server{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the Client.CheckExistence() API.
		CheckExistence: func(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, apiVersion string, options *armresources.ClientCheckExistenceOptions) (resp azfake.Responder[armresources.ClientCheckExistenceResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, apiVersion, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// use resp to set the desired response
			resp.SetResponse(http.StatusNoContent, armresources.ClientCheckExistenceResponse{Success: true}, nil)

			// to simulate the failure case, use errResp
			//errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armresources.NewClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.CheckExistence(context.TODO(), "fakeResourceGroup", "fakeProvider", "fakeParent", "fakeType", "fakeResource", "fakeAPIVersion", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Success)

	// APIs that haven't been faked will return an error
	_, err = client.Get(context.TODO(), "fakeResourceGroup", "fakeProvider", "fakeParent", "fakeType", "fakeResource", "fakeAPIVersion", nil)

	fmt.Println(err.Error())

	// Output:
	// true
	// fake for method Get not implemented
}
