// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// DatabaseRequestOptions includes options for operations against a database.
type DatabaseRequestOptions struct {
	IfMatchEtag     *azcore.ETag
	IfNoneMatchEtag *azcore.ETag
}

func (options *DatabaseRequestOptions) toHeaders() *map[string]string {
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
