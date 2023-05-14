// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// PartitionKeyRangesResponse represents the response from a partition key ranges request.
type PartitionKeyRangesResponse struct {
	//  PartitionKeyRangesProperties contains the unmarshalled response body in PartitionKeyRangesResponse format.
	PartitionKeyRangesProperties *PartitionKeyRangesProperties
	Response
}

func newPartitionKeyRangeResponse(resp *http.Response) (PartitionKeyRangesResponse, error) {
	response := PartitionKeyRangesResponse{
		Response: newResponse(resp),
	}
	properties := &PartitionKeyRangesProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.PartitionKeyRangesProperties = properties

	return response, nil
}
