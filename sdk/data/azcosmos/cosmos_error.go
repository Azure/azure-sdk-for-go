// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// CosmosError is used as base error for any error response from the Cosmos service.
type CosmosError struct {
	body string

	// ErrorCode is the error code returned by the service if available.
	ErrorCode string

	// StatusCode is the HTTP status code.
	StatusCode int

	// RawResponse is the underlying HTTP response.
	RawResponse *http.Response
}

func newCosmosError(response *http.Response) error {
	bytesRead, err := azruntime.Payload(response)
	if err != nil {
		return err
	}

	cError := CosmosError{
		StatusCode:  response.StatusCode,
		RawResponse: response,
		body:        string(bytesRead),
	}

	// Attempt to extract Code from body
	var cErrorResponse cosmosErrorResponse
	err = json.Unmarshal(bytesRead, &cErrorResponse)
	if err == nil {
		cError.ErrorCode = cErrorResponse.Code
	}

	return &cError
}

func newCosmosErrorWithStatusCode(statusCode int, requestCharge *float32) error {
	rawResponse := &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
	}

	if requestCharge != nil {
		rawResponse.Header.Add(cosmosHeaderRequestCharge, fmt.Sprint(*requestCharge))
	}

	return &CosmosError{
		StatusCode:  statusCode,
		RawResponse: rawResponse,
	}
}

func (e *CosmosError) Error() string {
	if e.body == "" {
		return "response contained no body"
	}
	return e.body
}

type cosmosErrorResponse struct {
	Code string `json:"Code"`
}
