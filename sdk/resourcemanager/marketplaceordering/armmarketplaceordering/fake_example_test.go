// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmarketplaceordering_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering/fake"
)

func ExampleMarketplaceAgreementsServer() {
	// first, create an instance of the fake server for the client you wish to test.
	// the type name of the server will be similar to the corresponding client, with
	// the suffix "Server" instead of "Client".
	fakeServer := fake.MarketplaceAgreementsServer{

		// next, provide implementations for the APIs you wish to fake.
		// this fake corresponds to the MarketplaceAgreementsClient.Get() API.
		Get: func(ctx context.Context, offerType armmarketplaceordering.OfferType, publisherID string, offerID string, planID string, options *armmarketplaceordering.MarketplaceAgreementsClientGetOptions) (resp azfake.Responder[armmarketplaceordering.MarketplaceAgreementsClientGetResponse], errResp azfake.ErrorResponder) {
			// the values of ctx, offerType, publisherID, offerID, planID, and options come from the API call.

			// the named return values resp and errResp are used to construct the response
			// and are meant to be mutually exclusive. if both responses have been constructed,
			// the error response is selected.

			// use resp to set the desired response
			agreementResp := armmarketplaceordering.MarketplaceAgreementsClientGetResponse{}
			agreementResp.ID = to.Ptr("/fake/resource/id")
			resp.SetResponse(http.StatusOK, agreementResp, nil)

			// to simulate the failure case, use errResp
			// errResp.SetResponseError(http.StatusBadRequest, "ThisIsASimulatedError")

			return
		},
	}

	// now create the corresponding client, connecting the fake server via the client options
	client, err := armmarketplaceordering.NewMarketplaceAgreementsClient("subscriptionID", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewMarketplaceAgreementsServerTransport(&fakeServer),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// call the API. the provided values will be passed to the fake's implementation.
	// the response or error values returned by the API call are from the fake.
	resp, err := client.Get(context.TODO(), armmarketplaceordering.OfferTypeVirtualmachine, "fakePublisherID", "fakeOfferID", "fakePlanID", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*resp.ID)

	// APIs that haven't been faked will return an error
	_, err = client.Cancel(context.TODO(), "fakePublisherID", "fakeOfferID", "fakePlanID", nil)

	fmt.Println(err.Error())

	// Output:
	// /fake/resource/id
	// fake for method Cancel not implemented
}
