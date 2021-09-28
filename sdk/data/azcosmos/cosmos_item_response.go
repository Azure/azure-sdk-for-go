// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// CosmosItemResponse represents the response from an item request.
type CosmosItemResponse struct {
	// The byte content of the operation response.
	Value []byte
	CosmosResponse
}

// SessionToken contains the value from the session token header to be used on session consistency.
func (c *CosmosItemResponse) SessionToken() string {
	return c.RawResponse.Header.Get(cosmosHeaderSessionToken)
}

func newCosmosItemResponse(resp *http.Response) (CosmosItemResponse, error) {
	response := CosmosItemResponse{}
	defer resp.Body.Close()
	response.RawResponse = resp
	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}
	response.Value = body
	return response, nil
}
