// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// ThroughputResponse represents the response from a throughput request.
type ThroughputResponse struct {
	// ThroughputProperties contains the unmarshalled response body in ThroughputProperties format.
	ThroughputProperties *ThroughputProperties
	cosmosResponse
}

func newThroughputResponse(resp *azcore.Response) (ThroughputResponse, error) {
	response := ThroughputResponse{}
	response.RawResponse = resp.Response
	properties := &ThroughputProperties{}
	err := resp.UnmarshalAsJSON(properties)
	if err != nil {
		return response, err
	}
	response.ThroughputProperties = properties
	return response, nil
}
