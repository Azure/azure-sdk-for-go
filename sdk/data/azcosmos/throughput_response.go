// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strconv"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ThroughputResponse represents the response from a throughput request.
type ThroughputResponse struct {
	// ThroughputProperties contains the unmarshalled response body in ThroughputProperties format.
	ThroughputProperties *ThroughputProperties
	Response
	// IsReplacePending returns the state of a throughput update.
	IsReplacePending bool
	// MinThroughput is minimum throughput in measurement of request units per second in the Azure Cosmos service.
	MinThroughput *int32
}

func (r *ThroughputResponse) getIsReplacePending() bool {
	isPending := r.RawResponse.Header.Get(cosmosHeaderOfferReplacePending)
	if isPending == "" {
		return false
	}

	isPendingBool, err := strconv.ParseBool(isPending)
	if err != nil {
		return false
	}

	return isPendingBool
}

func (r *ThroughputResponse) readMinThroughput() *int32 {
	minThroughput := r.RawResponse.Header.Get(cosmosHeaderOfferMinimumThroughput)
	if minThroughput == "" {
		return nil
	}

	minThroughputInt, err := strconv.ParseInt(minThroughput, 10, 32)
	if err != nil {
		return nil
	}

	minThroughputAsInt := int32(minThroughputInt)

	return &minThroughputAsInt
}

func newThroughputResponse(resp *http.Response, extraRequestCharge *float32) (ThroughputResponse, error) {
	response := ThroughputResponse{
		Response: newResponse(resp),
	}
	properties := &ThroughputProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.ThroughputProperties = properties

	if extraRequestCharge != nil {
		currentRequestCharge := response.RequestCharge + *extraRequestCharge
		response.RequestCharge = currentRequestCharge
	}

	response.IsReplacePending = response.getIsReplacePending()
	response.MinThroughput = response.readMinThroughput()
	return response, nil
}
