// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
)

type cosmosOffers struct {
	connection *cosmosClientConnection
}

type cosmosOffersResponse struct {
	Offers []ThroughputProperties `json:"Offers"`
}

func (c cosmosOffers) ReadThroughputIfExists(
	ctx context.Context,
	targetRID string,
	requestOptions *ThroughputRequestOptions) (ThroughputResponse, error) {
	// TODO: might want to replace with query once that is in
	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeOffer,
		resourceAddress: "",
	}

	path, err := generatePathForNameBased(resourceTypeOffer, "", true)
	if err != nil {
		return ThroughputResponse{}, err
	}

	azResponse, err := c.connection.sendQueryRequest(
		path,
		ctx,
		fmt.Sprintf(`SELECT * FROM c WHERE c.offerResourceId = '%s'`, targetRID),
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return ThroughputResponse{}, err
	}

	return newThroughputResponse(azResponse)
}
