// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// TransactionalBatchResponse contains response from a transactional batch operation.
type TransactionalBatchResponse struct {
	Response
	// SessionToken contains the value from the session token header to be used on session consistency.
	SessionToken string
	// OperationResults contains the individual batch operation results.
	// The order of the results is the same as the order of the operations in the batch.
	OperationResults []TransactionalBatchResult
	// Success indicates if the transaction was successfully committed.
	// If false, one of the operations in the batch failed.
	// Inspect the OperationResults, any operation with status code http.StatusFailedDependency is a dependency failure.
	// The cause of the batch failure is the first operation with status code different from http.StatusFailedDependency.
	Success bool
}

func newTransactionalBatchResponse(resp *http.Response) (TransactionalBatchResponse, error) {
	response := TransactionalBatchResponse{
		Response: newResponse(resp),
	}

	response.SessionToken = resp.Header.Get(cosmosHeaderSessionToken)

	response.Success = resp.StatusCode != http.StatusMultiStatus

	if err := runtime.UnmarshalAsJSON(resp, &response.OperationResults); err != nil {
		return TransactionalBatchResponse{}, err
	}

	return response, nil
}

// TransactionalBatchResult represents the result of a single operation in a batch.
type TransactionalBatchResult struct {
	// StatusCode contains the status code of the operation.
	StatusCode int32
	// RequestCharge contains the request charge for the operation.
	RequestCharge float32
	// ResourceBody contains the body response of the operation.
	// This property is available depending on the EnableContentResponseOnWrite option.
	ResourceBody []byte
	// ETag contains the ETag of the operation.
	ETag azcore.ETag
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (or *TransactionalBatchResult) UnmarshalJSON(b []byte) error {
	var attributes map[string]json.RawMessage
	err := json.Unmarshal(b, &attributes)
	if err != nil {
		return err
	}

	if statusCode, ok := attributes["statusCode"]; ok {
		if err := json.Unmarshal(statusCode, &or.StatusCode); err != nil {
			return err
		}
	}

	if requestCharge, ok := attributes["requestCharge"]; ok {
		if err := json.Unmarshal(requestCharge, &or.RequestCharge); err != nil {
			return err
		}
	}

	if etag, ok := attributes["eTag"]; ok {
		if err := json.Unmarshal(etag, &or.ETag); err != nil {
			return err
		}
	}

	if body, ok := attributes["resourceBody"]; ok {
		or.ResourceBody = body
	}

	return nil
}
