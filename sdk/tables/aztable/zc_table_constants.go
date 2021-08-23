// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "errors"

const (
	timestamp    = "Timestamp"
	partitionKey = "PartitionKey"
	rowKey       = "RowKey"
	etag         = "ETag"
)

var errPartitionKeyRowKeyError = errors.New("entity must have a PartitionKey and RowKey")
var errTooManyAccessPoliciesError = errors.New("you cannot set more than five (5) access policies at a time")
