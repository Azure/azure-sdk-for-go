// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Response is the base response type for all responses from the Azure Cosmos DB database service.
// It contains base methods and properties that are common to all responses.
type Response struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
	// RequestCharge contains the value from the request charge header.
	RequestCharge float32
	// ActivityID contains the value from the activity header.
	ActivityID string
	// ETag contains the value from the ETag header.
	ETag azcore.ETag
}

func newResponse(resp *http.Response) Response {
	response := Response{}
	response.RawResponse = resp
	response.RequestCharge = response.readRequestCharge()
	response.ActivityID = resp.Header.Get(cosmosHeaderActivityId)
	response.ETag = azcore.ETag(resp.Header.Get(cosmosHeaderEtag))
	return response
}

func (c *Response) readRequestCharge() float32 {
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
