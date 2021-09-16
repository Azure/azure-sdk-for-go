// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// CosmosDatabaseRequestOptions includes options for operations against a database.
type CosmosDatabaseRequestOptions struct {
	IfMatchEtag     string
	IfNoneMatchEtag string
}

func (options *CosmosDatabaseRequestOptions) toHeaders() *map[string]string {
	if options.IfMatchEtag == "" && options.IfNoneMatchEtag == "" {
		return nil
	}

	headers := make(map[string]string)
	if options.IfMatchEtag != "" {
		headers[azcore.HeaderIfMatch] = options.IfMatchEtag
	}
	if options.IfNoneMatchEtag != "" {
		headers[azcore.HeaderIfNoneMatch] = options.IfNoneMatchEtag
	}
	return &headers
}
