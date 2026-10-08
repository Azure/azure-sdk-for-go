// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreconditionNativeMapping(t *testing.T) {
	match, err := IfMatch(`"etag"`)
	require.NoError(t, err)
	none, err := IfNoneMatch(`"etag"`)
	require.NoError(t, err)
	for _, kind := range []operationKind{
		operationKindReadItem, operationKindCreateItem, operationKindReplaceItem,
		operationKindUpsertItem, operationKindDeleteItem,
	} {
		for _, tt := range []struct {
			condition Precondition
			wantKind  int32
			wantETag  string
		}{
			{Precondition{}, 0, ""},
			{match, 1, `"etag"`},
			{none, 2, `"etag"`},
		} {
			req := itemRequest{kind: kind, itemID: "id", partitionKey: NewPartitionKeyString("pk")}
			tt.condition.apply(&req)
			native, release := inspectNativeItemRequest(req)
			require.Equal(t, tt.wantKind, native.preconditionKind)
			require.Equal(t, tt.wantETag, native.preconditionETag)
			release()
		}
	}
}
