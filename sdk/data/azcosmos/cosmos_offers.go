// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type cosmosOffers struct {
	client *Client
}

type cosmosOffersResponse struct {
	Offers []ThroughputProperties `json:"Offers"`
}

func (c cosmosOffers) ReadThroughputIfExists(
	ctx context.Context,
	targetRID string,
	requestOptions *ThroughputOptions) (ThroughputResponse, error) {
	// TODO: might want to replace with query iterator once that is in
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeOffer,
		resourceAddress: "",
	}

	path, err := generatePathForNameBased(resourceTypeOffer, "", true)
	if err != nil {
		return ThroughputResponse{}, err
	}

	azResponse, err := c.client.sendQueryRequest(
		path,
		ctx,
		fmt.Sprintf(`SELECT * FROM c WHERE c.offerResourceId = '%s'`, targetRID),
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return ThroughputResponse{}, err
	}

	var theOffers cosmosOffersResponse
	err = azruntime.UnmarshalAsJSON(azResponse, &theOffers)
	if err != nil {
		return ThroughputResponse{}, err
	}

	queryRequestCharge := newResponse(azResponse).RequestCharge
	if len(theOffers.Offers) == 0 {
		return ThroughputResponse{}, newCosmosErrorWithStatusCode(http.StatusNotFound, &queryRequestCharge)
	}

	// Now read the individual offer
	operationContext = pipelineRequestOptions{
		resourceType:    resourceTypeOffer,
		resourceAddress: theOffers.Offers[0].offerId,
		isRidBased:      true,
	}

	path, err = generatePathForNameBased(resourceTypeOffer, theOffers.Offers[0].selfLink, false)
	if err != nil {
		return ThroughputResponse{}, err
	}

	azResponse, err = c.client.sendGetRequest(
		path,
		ctx,
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return ThroughputResponse{}, err
	}

	return newThroughputResponse(azResponse, &queryRequestCharge)
}

func (c cosmosOffers) ReplaceThroughputIfExists(
	ctx context.Context,
	properties ThroughputProperties,
	targetRID string,
	requestOptions *ThroughputOptions) (ThroughputResponse, error) {

	readResponse, err := c.ReadThroughputIfExists(ctx, targetRID, requestOptions)
	if err != nil {
		return ThroughputResponse{}, err
	}

	readRequestCharge := readResponse.RequestCharge
	readResponse.ThroughputProperties.offer = properties.offer

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeOffer,
		resourceAddress: readResponse.ThroughputProperties.offerId,
		isRidBased:      true,
	}

	path, err := generatePathForNameBased(resourceTypeOffer, readResponse.ThroughputProperties.selfLink, false)
	if err != nil {
		return ThroughputResponse{}, err
	}

	azResponse, err := c.client.sendPutRequest(
		path,
		ctx,
		readResponse.ThroughputProperties,
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return ThroughputResponse{}, err
	}

	return newThroughputResponse(azResponse, &readRequestCharge)
}
