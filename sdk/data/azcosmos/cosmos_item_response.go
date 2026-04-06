// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ItemResponse represents the response from an item request.
type ItemResponse struct {
	// The byte content of the operation response.
	Value []byte
	Response
	// SessionToken contains the value from the session token header to be used on session consistency.
	SessionToken *string
}

func newItemResponse(resp *http.Response) (ItemResponse, error) {
	response := ItemResponse{
		Response: newResponse(resp),
	}
	sessionToken := resp.Header.Get(cosmosHeaderSessionToken)
	if sessionToken != "" {
		response.SessionToken = &sessionToken
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}
	response.Value = body
	return response, nil
}
