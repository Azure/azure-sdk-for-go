//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_newAuthenticationClient(t *testing.T) {
	client := newAuthenticationClient("test", nil)
	require.NotNil(t, client)
}
