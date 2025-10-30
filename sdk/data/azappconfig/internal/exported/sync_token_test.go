// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	tokens, err := ParseSyncToken(SyncToken(""))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("  \t"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("id=val"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("id=val;"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("=;sn=1"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("id=val1;version=1"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken(";sn=1"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("sn=1;id=val"))
	require.Error(t, err)
	require.Nil(t, tokens)

	tokens, err = ParseSyncToken(SyncToken("id=val1;sn=1"))
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	require.EqualValues(t, SyncTokenValues{
		ID:      "id",
		Value:   "val1",
		Version: 1,
	}, tokens[0])

	tokens, err = ParseSyncToken(SyncToken("id1=val;sn=1,id2=val;sn=1"))
	require.NoError(t, err)
	require.Len(t, tokens, 2)
	require.EqualValues(t, []SyncTokenValues{
		{
			ID:      "id1",
			Value:   "val",
			Version: 1,
		},
		{
			ID:      "id2",
			Value:   "val",
			Version: 1,
		},
	}, tokens)
}
