// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPartitionKeyRangeRequestOptionsToHeaders(t *testing.T) {
	options := &partitionKeyRangeOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}
}
