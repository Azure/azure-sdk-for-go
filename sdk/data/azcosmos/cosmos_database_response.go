// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// CosmosDatabaseResponse represents the response from a database request.
type CosmosDatabaseResponse struct {
	// DatabaseProperties contains the unmarshalled response body in CosmosDatabaseProperties format.
	DatabaseProperties *CosmosDatabaseProperties
	CosmosResponse
}

func newCosmosDatabaseResponse(resp *http.Response, database *CosmosDatabase) (CosmosDatabaseResponse, error) {
	response := CosmosDatabaseResponse{}
	response.RawResponse = resp
	properties := &CosmosDatabaseProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.DatabaseProperties = properties
	response.DatabaseProperties.Database = database
	return response, nil
}
