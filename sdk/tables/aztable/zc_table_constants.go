// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "errors"

const (
	timestamp    = "Timestamp"
	partitionKey = "PartitionKey"
	rowKey       = "RowKey"
	etag         = "ETag"
	OdataType    = "@odata.type"
	ISO8601      = "2006-01-02T15:04:05.9999999Z"
)

var errPartitionKeyRowKeyError = errors.New("entity must have a PartitionKey and RowKey")
var errTooManyAccessPoliciesError = errors.New("you cannot set more than five (5) access policies at a time.")
