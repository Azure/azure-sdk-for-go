// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// cosmosError is used as base error for any error response from the Cosmos service.
type cosmosError struct {
	body        string
	code        string
	rawResponse *http.Response
}

func (e *cosmosError) StatusCode() int {
	if e.rawResponse == nil {
		return 0
	}
	return e.rawResponse.StatusCode
}

func (e *cosmosError) ErrorCode() string {
	return e.code
}

func (e *cosmosError) RawResponse() *http.Response {
	return e.rawResponse
}

func newCosmosError(response *http.Response) error {
	bytesRead, err := azruntime.Payload(response)
	if err != nil {
		return err
	}

	cError := cosmosError{
		rawResponse: response,
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

func newCosmosErrorWithStatusCode(statusCode int, requestCharge *float32) error {
	rawResponse := &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
	}

	if requestCharge != nil {
		rawResponse.Header.Add(cosmosHeaderRequestCharge, fmt.Sprint(*requestCharge))
	}

	return &cosmosError{
		rawResponse: rawResponse,
	}
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
