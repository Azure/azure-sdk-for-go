// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newAuthenticationClient(t *testing.T) {
	client, err := NewAuthenticationClient("test", nil)
	require.NoError(t, err)
	require.NotNil(t, client)
}
