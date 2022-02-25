// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ContainerResponse represents the response from a container request.
type ContainerResponse struct {
	// ContainerProperties contains the unmarshalled response body in ContainerProperties format.
	ContainerProperties *ContainerProperties
	Response
}

func newContainerResponse(resp *http.Response) (ContainerResponse, error) {
	response := ContainerResponse{
		Response: newResponse(resp),
	}
	properties := &ContainerProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.ContainerProperties = properties
	return response, nil
}
