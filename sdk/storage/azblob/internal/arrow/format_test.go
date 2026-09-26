// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package arrow

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

func TestUseArrow(t *testing.T) {
	for _, f := range []exported.StorageResponseFormat{exported.StorageResponseFormatAuto, exported.StorageResponseFormatXML} {
		useArrow, err := UseArrow(f)
		require.NoError(t, err)
		require.Equal(t, Supported && exported.ResolveAutoFormat(f) == exported.StorageResponseFormatArrow, useArrow)
	}

	useArrow, err := UseArrow(exported.StorageResponseFormatArrow)
	if Supported {
		require.NoError(t, err)
		require.True(t, useArrow)
	} else {
		require.ErrorIs(t, err, ErrNotSupported)
		require.False(t, useArrow)
	}
}
