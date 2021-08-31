// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// CosmosDatabaseResponse represents the response from a database request.
type CosmosDatabaseResponse struct {
	// DatabaseProperties contains the unmarshalled response body in CosmosDatabaseProperties format.
	DatabaseProperties *CosmosDatabaseProperties
	CosmosResponse
}

func newCosmosDatabaseResponse(resp *azcore.Response, database *CosmosDatabase) (CosmosDatabaseResponse, error) {
	response := CosmosDatabaseResponse{}
	response.RawResponse = resp.Response
	properties := &CosmosDatabaseProperties{}
	err := resp.UnmarshalAsJSON(properties)
	if err != nil {
		return response, err
	}
	response.DatabaseProperties = properties
	response.DatabaseProperties.Database = database
	return response, nil
}
