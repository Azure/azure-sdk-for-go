// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// ReadDatabaseOptions includes options ReadDatabase operation.
type ReadDatabaseOptions struct {
	IfMatchEtag     *azcore.ETag
	IfNoneMatchEtag *azcore.ETag
}

func (options *ReadDatabaseOptions) toHeaders() *map[string]string {
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

// DeleteDatabaseOptions includes options DeleteDatabase operation.
type DeleteDatabaseOptions struct {
	IfMatchEtag     *azcore.ETag
	IfNoneMatchEtag *azcore.ETag
}

func (options *DeleteDatabaseOptions) toHeaders() *map[string]string {
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

// CreateDatabaseOptions are options for the CreateDatabase operation
type CreateDatabaseOptions struct {
	// ThroughputProperties: Optional throughput configuration of the database
	ThroughputProperties *ThroughputProperties
}
