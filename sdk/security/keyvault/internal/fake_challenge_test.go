// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeChallenge(t *testing.T) {
	interceptor := &FakeChallenge{}
	req, err := http.NewRequest(http.MethodGet, "https://42.vault.azure.net", nil)
	require.NoError(t, err)

	resp, err, intercepted := interceptor.Do(req)
	require.NoError(t, err)
	require.True(t, intercepted)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	headers := resp.Header
	header, ok := headers["Www-Authenticate"]
	require.True(t, ok)
	require.Equal(t, `Bearer authorization="https://fake.local/tenant" resource="https://vault.azure.net"`, header[0])

	req.Header.Set("Authorization", "fakeauthorization")
	resp, err, intercepted = interceptor.Do(resp.Request)
	require.NoError(t, err)
	require.False(t, intercepted)
	require.Nil(t, resp)
}
