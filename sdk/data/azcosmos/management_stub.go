// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import "context"

func (c *Client) executeManagementRequest(context.Context, managementRequest) (managementResult, error) {
	return managementResult{}, newDriverUnavailableError()
}

func (c *Client) openManagementQuery(context.Context, managementRequest, ManagementQueryOptions) (queryCursor, error) {
	return nil, newDriverUnavailableError()
}
