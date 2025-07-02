// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// PartitionKeyRangeOptions includes options for reading partition key ranges.
type partitionKeyRangeOptions struct{}

// toHeaders converts the options to a map of HTTP headers.
func (options *partitionKeyRangeOptions) toHeaders() *map[string]string {
	return nil
}
