// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

type partitionKeyInternal struct{}

func newPartitionKeyInternal(values []interface{}) (*partitionKeyInternal, error) {
	components := make([]int, len(values))
	for _, v := range values {
		switch v.(type) {
		case nil:

		case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		default:
			return nil, fmt.Errorf("PartitionKey can only be a string, bool, or a number: '%T'", v)
		}
	}
}
