// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// CosmosContainerResponse represents the response from a container request.
type CosmosContainerResponse struct {
	// ContainerProperties contains the unmarshalled response body in CosmosContainerProperties format.
	ContainerProperties *CosmosContainerProperties
	CosmosResponse
}

func newCosmosContainerResponse(resp *http.Response, container *CosmosContainer) (CosmosContainerResponse, error) {
	response := CosmosContainerResponse{
		CosmosResponse: newCosmosResponse(resp),
	}
	properties := &CosmosContainerProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.ContainerProperties = properties
	response.ContainerProperties.Container = container
	return response, nil
}
