// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// DatabaseResponse represents the response from a database request.
type DatabaseResponse struct {
	// DatabaseProperties contains the unmarshalled response body in DatabaseProperties format.
	DatabaseProperties *DatabaseProperties
	Response
}

func newDatabaseResponse(resp *http.Response) (DatabaseResponse, error) {
	response := DatabaseResponse{
		Response: newResponse(resp),
	}
	properties := &DatabaseProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.DatabaseProperties = properties
	return response, nil
}
