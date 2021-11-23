//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBodilessMatcher(t *testing.T) {
	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	err = AddBodilessMatcher(t, nil)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)
}
