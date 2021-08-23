// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosError is used as base error for any error response from the Cosmos service.
type cosmosError struct {
	body        string
	code        string
	rawResponse *http.Response
}

func (e *cosmosError) StatusCode() string {
	if e.rawResponse == nil {
		return ""
	}
	return e.rawResponse.Status
}

func (e *cosmosError) ErrorCode() string {
	return e.code
}

func (e *cosmosError) RawResponse() *http.Response {
	return e.rawResponse
}

func newCosmosError(response *azcore.Response) error {
	bytesRead, err := response.Payload()
	if err != nil {
		return err
	}

	cError := cosmosError{
		rawResponse: response.Response,
		body:        string(bytesRead),
	}

	// Attempt to extract Code from body
	var cErrorResponse cosmosErrorResponse
	err = json.Unmarshal(bytesRead, &cErrorResponse)
	if err == nil {
		cError.code = cErrorResponse.Code
	}

	return &cError
}

func (e *cosmosError) Error() string {
	if e.body == "" {
		return "response contained no body"
	}
	return e.body
}

type cosmosErrorResponse struct {
	Code string `json:"Code"`
}
