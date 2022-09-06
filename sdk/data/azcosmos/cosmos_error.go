// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func newCosmosError(response *http.Response) error {
	bytesRead, err := azruntime.Payload(response)
	if err != nil {
		return err
	}

	cError := azcore.ResponseError{
		StatusCode:  response.StatusCode,
		RawResponse: response,
	}

	// Attempt to extract Code from body
	var cErrorResponse cosmosErrorResponse
	err = json.Unmarshal(bytesRead, &cErrorResponse)
	if err == nil {
		cError.ErrorCode = cErrorResponse.Code
	}

	return &cError
}

type cosmosErrorResponse struct {
	Code string `json:"Code"`
}
