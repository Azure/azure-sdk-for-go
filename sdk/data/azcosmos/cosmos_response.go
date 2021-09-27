// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// CosmosResponse is the base response type for all responses from the Azure Cosmos DB database service.
// It contains base methods and properties that are common to all responses.
type CosmosResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// RequestCharge contains the value from the request charge header.
func (c *CosmosResponse) RequestCharge() float32 {
	requestChargeString := c.RawResponse.Header.Get(cosmosHeaderRequestCharge)
	if requestChargeString == "" {
		return 0
	}
	f, err := strconv.ParseFloat(requestChargeString, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

// ActivityId contains the value from the activity header.
func (c *CosmosResponse) ActivityId() string {
	return c.RawResponse.Header.Get(cosmosHeaderActivityId)
}

// ETag contains the value from the ETag header.
func (c *CosmosResponse) ETag() azcore.ETag {
	return azcore.ETag(c.RawResponse.Header.Get(cosmosHeaderEtag))
}
