// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatManagementError(t *testing.T) {
	err := FormatManagementError([]byte("not-xml"), errors.New("falls back to this error instead"))
	require.EqualError(t, err, "falls back to this error instead")

	err = FormatManagementError([]byte("<Error><Code>405</Code><Detail>Some detail</Detail></Error>"), nil)
	require.EqualError(t, err, "error code: 405, Details: Some detail")
}
