// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatManagementError(t *testing.T) {
	err := FormatManagementError([]byte("not-xml"))
	require.EqualError(t, err, "body:not-xml error:EOF")

	err = FormatManagementError([]byte("<Error><Code>405</Code><Detail>Some detail</Detail></Error>"))
	require.EqualError(t, err, "error code: 405, Details: Some detail")
}
