// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor/fake"
)

func ExampleMetricsServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeMetricsServer := fake.MetricsServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the MetricsClient.List() API.
		List: func(ctx context.Context, resourceURI string, options *armmonitor.MetricsClientListOptions) (resp azfake.Responder[armmonitor.MetricsClientListResponse], errResp azfake.ErrorResponder) {
			// the value of ctx, resourceURI, and options come from the API call

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// construct the response type, populating fields as required
			listResp := armmonitor.MetricsClientListResponse{}
			listResp.Namespace = to.Ptr("fake_namespace")

			// use resp to set the desired response
			resp.SetResponse(http.StatusOK, listResp, nil)

			// to simulate the failure case, use errResp
			// errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armmonitor.NewMetricsClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewMetricsServerTransport(&fakeMetricsServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.List(context.TODO(), "/fake/resource/uri", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.Namespace)

	// APIs that haven't been faked will return an error
	_, err = client.ListAtSubscriptionScope(context.TODO(), "fakeRegion", nil)

	fmt.Println(err.Error())

	// Output:
	// fake_namespace
	// fake for method ListAtSubscriptionScope not implemented
}
