// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strconv"
	"time"
	
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChangeFeedOptions defines the options for retrieving the change feed.
// Incorporate Continuation
type ChangeFeedOptions struct {
	// MaxItemCount limits the number of items returned per page.
	// Valid values are > 0. The service may return fewer items than requested.
	MaxItemCount int32 `json:"x-ms-max-item-count"`
	
	// AIM is the header that indicates this is a change feed request.
	// It should always be set to "Incremental Feed".
	// This header is required for change feed requests.
	AIM string `json:"A-IM"`
	
	// IfNoneMatch can be used with an ETag to only retrieve changes if the feed has changed.
	// If the feed hasn't changed since the ETag, a 304 Not Modified response is returned.
	IfNoneMatch *azcore.ETag `json:"If-None-Match"`
	
	// IfModifiedSince can be used with UTC time to retrieve changes if the feed has been modified after that
	// specific time
	IfModifiedSince *time.Time `json:"If-Modified-Since"`

	// PartitionKey is the logical partition key value for the request.
	// Use this to read from a specific logical partition.
	PartitionKey *PartitionKey `json:"x-ms-documentdb-partitionkey"`
	
	// PartitionKeyRangeID is the physical partition key range id for the request.
	// Use this to read from a specific physical partition instead of all partitions.
	// This is useful when you want to process changes from multiple physical partitions in parallel.
	PartitionKeyRangeID *string `json:"x-ms-documentdb-partitionkeyrangeid"`
}

func (options *ChangeFeedOptions) toHeaders() *map[string]string {
	headers := make(map[string]string)
	
	if options.MaxItemCount > 0 {
		headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(int64(options.MaxItemCount), 10)
	} else {
		headers[cosmosHeaderMaxItemCount] = int32(-1)
	}

	if options.AIM != "" {
		headers[cosmosHeaderAIM] = options.AIM
	} else {
		headers[cosmosHeaderAIM] = cosmosHeaderChangeFeedIncremental
	}
	
	if options.IfNoneMatch != nil {
		headers[headerIfNoneMatch] = string(*options.IfNoneMatch)
	}
	
	if options.IfModifiedSince != nil {
		headers[cosmosHeaderIfModifiedSince] = options.IfModifiedSince.UTC().Format(time.RFC1123)
	}
	
	if options.PartitionKey != nil {
		partitionKeyJSON, err := options.PartitionKey.toJsonString()
		if err == nil {
			headers[cosmosHeaderPartitionKey] = string(partitionKeyJSON)
		}
	}
	
	if options.PartitionKeyRangeID != nil {
		headers[cosmosHeaderPartitionKeyRangeID] = *options.PartitionKeyRangeID
	}

	if len(headers) == 0 {
		return nil
	}
	return &headers
}
