// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
)

// NewUniqueRequestIDPolicy creates a policy object that sets the request's x-ms-client-request-id header if it doesn't already exist.
func NewUniqueRequestIDPolicy() Policy {
	return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
		const xMsClientRequestID = "x-ms-client-request-id"
		id := req.Request.Header.Get(xMsClientRequestID)
		if id == "" {
			// Add a unique request ID if the caller didn't specify one already
			req.Request.Header.Set(xMsClientRequestID, "TODO" /*newUUID().String()*/)
		}
		return req.Do(ctx)
	})
}
