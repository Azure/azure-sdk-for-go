// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
)

// QueryItemsResponse contains response from the query operation.
type TransactionalBatchResponse struct {
	Response
}

func newTransactionalBatchResponse(resp *http.Response) (TransactionalBatchResponse, error) {
	response := TransactionalBatchResponse{
		Response: newResponse(resp),
	}

	return response, nil
}
