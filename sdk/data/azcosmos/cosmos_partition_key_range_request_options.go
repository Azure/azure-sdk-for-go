// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ReadPartitionKeyRangesOptions includes options for reading partition key ranges.
type PartitionKeyRangeOptions struct{}

// toHeaders converts the options to a map of HTTP headers.
func (options *PartitionKeyRangeOptions) toHeaders() *map[string]string {
	return nil
}
