// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ThroughputResponse represents the response from a throughput request.
type ThroughputResponse struct {
	// ThroughputProperties contains the unmarshalled response body in ThroughputProperties format.
	ThroughputProperties *ThroughputProperties
	cosmosResponse
}

func newThroughputResponse(resp *azcore.Response) (ThroughputResponse, error) {
	var theOffers cosmosOffersResponse
	err := resp.UnmarshalAsJSON(&theOffers)
	if err != nil {
		return ThroughputResponse{}, err
	}

	if len(theOffers.Offers) == 0 {
		return ThroughputResponse{}, newCosmosErrorWithStatusCode(http.StatusNotFound)
	}

	response := ThroughputResponse{}
	response.RawResponse = resp.Response
	response.ThroughputProperties = &theOffers.Offers[0]
	return response, nil
}
