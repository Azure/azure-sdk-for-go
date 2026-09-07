// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportsMultiBlock(t *testing.T) {
	tests := []struct {
		name     string
		tv       TransferValidationType
		expected bool
	}{
		{"PrecomputedCRC64", TransferValidationTypeCRC64(0), false},
		{"PrecomputedMD5", TransferValidationTypeMD5(nil), false},
		{"ComputeCRC64", TransferValidationTypeComputeCRC64(), false},
		{"ComputeStructuredMessageCRC64", TransferValidationTypeComputeStructuredMessageCRC64(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, SupportsMultiBlock(tt.tv))
		})
	}
}
