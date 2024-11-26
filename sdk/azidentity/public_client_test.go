// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/stretchr/testify/require"
)

func TestPublicClientLogging(t *testing.T) {
	logMsgs := []string{}
	log.SetListener(func(e log.Event, msg string) {
		if e == EventAuthentication {
			logMsgs = append(logMsgs, msg)
		}
	})
	defer log.SetListener(nil)

	pc, err := newPublicClient(fakeTenantID, fakeClientID, credNameUserPassword, publicClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{},
		},
	})
	require.NoError(t, err)

	// client should log token scopes when acquiring a token from the cache or authority
	expected := fmt.Sprintf(scopeLogFmt, credNameUserPassword, strings.Join(testTRO.Scopes, ", "))
	for i := 0; i < 2; i++ {
		logMsgs = []string{}
		_, err = pc.GetToken(ctx, testTRO)
		require.NoError(t, err)

		scopesLogged := false
		for _, msg := range logMsgs {
			require.Contains(t, msg, credNameUserPassword)
			if strings.Contains(msg, testTRO.Scopes[0]) {
				scopesLogged = true
				require.Equal(t, expected, msg)
			}
		}
		require.True(t, scopesLogged)
	}
}
