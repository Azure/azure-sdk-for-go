// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/http"
	"strconv"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ThroughputResponse represents the response from a throughput request.
type ThroughputResponse struct {
	// ThroughputProperties contains the unmarshalled response body in ThroughputProperties format.
	ThroughputProperties *ThroughputProperties
	CosmosResponse
}

// IsReplacePending returns the state of a throughput update.
func (r *ThroughputResponse) IsReplacePending() *bool {
	isPending := r.RawResponse.Header.Get(cosmosHeaderOfferReplacePending)
	if isPending == "" {
		return nil
	}

	isPendingBool, err := strconv.ParseBool(isPending)
	if err != nil {
		return nil
	}

	return &isPendingBool
}

// MinThroughput is minimum throughput in measurement of request units per second in the Azure Cosmos service.
func (r *ThroughputResponse) MinThroughput() *int {
	minThroughput := r.RawResponse.Header.Get(cosmosHeaderOfferMinimumThroughput)
	if minThroughput == "" {
		return nil
	}

	minThroughputInt, err := strconv.ParseInt(minThroughput, 10, 32)
	if err != nil {
		return nil
	}

	minThroughputAsInt := int(minThroughputInt)

	return &minThroughputAsInt
}

func newThroughputResponse(resp *http.Response, extraRequestCharge *float32) (ThroughputResponse, error) {
	response := ThroughputResponse{}
	response.RawResponse = resp
	properties := &ThroughputProperties{}
	err := azruntime.UnmarshalAsJSON(resp, properties)
	if err != nil {
		return response, err
	}
	response.ThroughputProperties = properties

	if extraRequestCharge != nil {
		currentRequestCharge := response.RequestCharge() + *extraRequestCharge
		response.RawResponse.Header.Set(cosmosHeaderRequestCharge, fmt.Sprint(currentRequestCharge))
	}

	return response, nil
}
