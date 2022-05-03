// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// QueryItemsResponse contains response from the query operation.
type TransactionalBatchResponse struct {
	Response
	// SessionToken contains the value from the session token header to be used on session consistency.
	SessionToken string
	// OperationResults contains the individual batch operation results.
	OperationResults []TransactionalBatchResponseOperationResult
	// Committed indicates if the transaction was successfully committed.
	// If false, one of the operations in the batch failed.
	// Inspect the OperationResults, any operation with status code 424 is a dependency failure.
	// The cause of the batch failure is the first operation with status code different from 424.
	Committed bool
}

func newTransactionalBatchResponse(resp *http.Response) (TransactionalBatchResponse, error) {
	response := TransactionalBatchResponse{
		Response: newResponse(resp),
	}

	response.SessionToken = resp.Header.Get(cosmosHeaderSessionToken)

	response.Committed = resp.StatusCode != http.StatusMultiStatus

	if err := runtime.UnmarshalAsJSON(resp, &response.OperationResults); err != nil {
		return TransactionalBatchResponse{}, err
	}

	return response, nil
}

type TransactionalBatchResponseOperationResult struct {
	StatusCode    int32
	RequestCharge float32
	ResourceBody  []byte
	ETag          azcore.ETag
}

func (or *TransactionalBatchResponseOperationResult) UnmarshalJSON(b []byte) error {
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
