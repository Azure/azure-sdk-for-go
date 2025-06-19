// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPartitionKeyRangeRequestOptionsToHeaders(t *testing.T) {
	options := &PartitionKeyRangeOptions{}
	if options.toHeaders() == nil {
		t.Error("Expected headers to be non-nil")
	}
}
