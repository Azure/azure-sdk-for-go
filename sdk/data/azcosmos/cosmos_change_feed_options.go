// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strconv"
	"time"
	
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChangeFeedOptions defines the options for retrieving the change feed.
type ChangeFeedOptions struct {
	// PartitionKey is the logical partition key value for the request.
	// Use this to read from a specific logical partition.
	PartitionKey *PartitionKey

	// PartitionKeyRangeID is the physical partition key range id for the request.
	// Use this to read from a specific physical partition instead of all partitions.
	PartitionKeyRangeID *string

	// SessionToken can be set to guarantee change feed consistency with other requests.
	// More information: https://docs.microsoft.com/azure/cosmos-db/sql/consistency-levels
	SessionToken *string

	// ConsistencyLevel overrides the account defined consistency level for this operation.
	// Consistency can only be relaxed.
	ConsistencyLevel *ConsistencyLevel

	// MaxItemCount limits the number of items returned per page.
	// Valid values are > 0. The service may return fewer items than requested.
	MaxItemCount int32

	// ContinuationToken can be used to resume reading the change feed.
	// Pass the token returned from a previous ChangeFeedResponse to continue.
	ContinuationToken *string
	
	// StartTime specifies the point in time to start reading changes from.
	// Changes are returned from this time onwards.
	StartTime *time.Time
	
	// StartFromBeginning indicates that the change feed should be read from the beginning.
	// If true, StartTime is ignored.
	StartFromBeginning bool

	// IfNoneMatch can be used with an ETag to only retrieve changes if the feed has changed.
	// If the feed hasn't changed since the ETag, a 304 Not Modified response is returned.
	IfNoneMatch *azcore.ETag
}

func (options *ChangeFeedOptions) toHeaders() *map[string]string {
	headers := make(map[string]string)

	if options.PartitionKey != nil {
		partitionKeyJSON, _ := options.PartitionKey.MarshalJSON()
		headers[cosmosHeaderPartitionKey] = string(partitionKeyJSON)
	}
	
	if options.PartitionKeyRangeID != nil {
		headers[cosmosHeaderPartitionKeyRangeID] = *options.PartitionKeyRangeID
	}
	
	if options.SessionToken != nil {
		headers[cosmosHeaderSessionToken] = *options.SessionToken
	}
	
	if options.ConsistencyLevel != nil {
		headers[cosmosHeaderConsistencyLevel] = string(*options.ConsistencyLevel)
	}
	
	if options.MaxItemCount > 0 {
		headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(int64(options.MaxItemCount), 10)
	}
	
	if options.ContinuationToken != nil {
		headers[cosmosHeaderContinuationToken] = *options.ContinuationToken
	}
	
	if options.StartTime != nil {
		headers["A-IM"] = "Incremental feed"
		headers["If-Modified-Since"] = options.StartTime.UTC().Format(httpTimeFormat)
	} else if options.StartFromBeginning {
		headers["A-IM"] = "Incremental feed"
	}
	
	if options.IfNoneMatch != nil {
		headers["If-None-Match"] = string(*options.IfNoneMatch)
	}

	if len(headers) == 0 {
		return nil
	}
	return &headers
}

// httpTimeFormat is the time format used for HTTP headers (RFC1123).
const httpTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
