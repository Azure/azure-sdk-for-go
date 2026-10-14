// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

type deadlineTokenCredential struct {
	t     *testing.T
	calls int
}

func (c *deadlineTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls++
	require.Equal(c.t, []string{"test-scope"}, options.Scopes)
	deadline, ok := ctx.Deadline()
	require.True(c.t, ok)
	require.Equal(c.t, time.Minute, time.Until(deadline))
	require.NoError(c.t, ctx.Err(), "each acquisition starts with a fresh context")
	if c.calls == 1 {
		<-ctx.Done()
		return azcore.AccessToken{}, ctx.Err()
	}
	return azcore.AccessToken{Token: "refreshed", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestTokenAcquisitionTimeoutPreservesProvider(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		credential := &deadlineTokenCredential{t: t}
		state := &tokenProviderState{credential: credential, ctx: context.Background()}
		start := time.Now()
		_, err := state.getToken("test-scope")
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Equal(t, time.Minute, time.Since(start))
		require.NoError(t, state.ctx.Err(), "acquisition timeout must not cancel the shared provider")

		token, err := state.getToken("test-scope")
		require.NoError(t, err)
		require.Equal(t, "refreshed", token.Token)
		require.Equal(t, 2, credential.calls)
	})
}
