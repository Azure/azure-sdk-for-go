// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ThroughputOptions includes options for throughput operations.
type ThroughputOptions struct {
	// IfMatchEtag If-Match (ETag) associated with the request.
	IfMatchEtag *azcore.ETag
	// IfNoneMatchEtag If-None-Match (ETag) associated with the request.
	IfNoneMatchEtag *azcore.ETag
}

func (options *ThroughputOptions) toHeaders() *map[string]string {
	if options.IfMatchEtag == nil && options.IfNoneMatchEtag == nil {
		return nil
	}

	headers := make(map[string]string)
	if options.IfMatchEtag != nil {
		headers[headerIfMatch] = string(*options.IfMatchEtag)
	}
	if options.IfNoneMatchEtag != nil {
		headers[headerIfNoneMatch] = string(*options.IfNoneMatchEtag)
	}
	return &headers
}
