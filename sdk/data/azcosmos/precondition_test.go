// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestPreconditionConstructorsValidateETags(t *testing.T) {
	for _, constructor := range []struct {
		name string
		new  func(azcore.ETag) (Precondition, error)
		kind preconditionKind
	}{
		{"IfMatch", IfMatch, preconditionKindIfMatch},
		{"IfNoneMatch", IfNoneMatch, preconditionKindIfNoneMatch},
	} {
		t.Run(constructor.name, func(t *testing.T) {
			for _, etag := range []azcore.ETag{"", "\"etag\x00suffix\""} {
				condition, err := constructor.new(etag)
				require.ErrorContains(t, err, constructor.name)
				require.Zero(t, condition)
			}
			for _, etag := range []azcore.ETag{`"etag"`, `W/"weak-etag"`, "*"} {
				condition, err := constructor.new(etag)
				require.NoError(t, err)
				req := itemRequest{}
				condition.apply(&req)
				require.Equal(t, constructor.kind, req.preconditionKind)
				require.Equal(t, string(etag), req.preconditionETag, "ETags must reach the ABI unchanged")
			}
		})
	}
}

func TestPreconditionValueOwnershipAndReplacement(t *testing.T) {
	etag := azcore.ETag(`"first"`)
	match, err := IfMatch(etag)
	require.NoError(t, err)
	etag = `"second"`
	none, err := IfNoneMatch(etag)
	require.NoError(t, err)
	options := ReplaceItemOptions{Precondition: match}
	copy := options
	options.Precondition = none
	req := itemRequest{}
	copy.Precondition.apply(&req)
	require.Equal(t, preconditionKindIfMatch, req.preconditionKind)
	require.Equal(t, `"first"`, req.preconditionETag)
	options.Precondition.apply(&req)
	require.Equal(t, preconditionKindIfNoneMatch, req.preconditionKind)
	require.Equal(t, `"second"`, req.preconditionETag)
	(Precondition{}).apply(&req)
	require.Equal(t, preconditionKindNone, req.preconditionKind)
	require.Empty(t, req.preconditionETag)
}

func TestItemPreconditionsPreserveLifetimeOrdering(t *testing.T) {
	container := newTestContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	match, err := IfMatch(`"etag"`)
	require.NoError(t, err)
	none, err := IfNoneMatch(`"etag"`)
	require.NoError(t, err)
	pk := NewPartitionKeyString("pk")
	calls := []func(Precondition) (ItemResponse, error){
		func(p Precondition) (ItemResponse, error) {
			return container.ReadItem(ctx, pk, "id", &ReadItemOptions{Precondition: p})
		},
		func(p Precondition) (ItemResponse, error) {
			return container.CreateItem(ctx, pk, "id", []byte(`{"id":"id","pk":"pk"}`), &CreateItemOptions{Precondition: p})
		},
		func(p Precondition) (ItemResponse, error) {
			return container.ReplaceItem(ctx, pk, "id", []byte(`{"id":"id","pk":"pk"}`), &ReplaceItemOptions{Precondition: p})
		},
		func(p Precondition) (ItemResponse, error) {
			return container.UpsertItem(ctx, pk, "id", []byte(`{"id":"id","pk":"pk"}`), &UpsertItemOptions{Precondition: p})
		},
		func(p Precondition) (ItemResponse, error) {
			return container.DeleteItem(ctx, pk, "id", &DeleteItemOptions{Precondition: p})
		},
	}
	for _, condition := range []Precondition{{}, match, none} {
		for _, call := range calls {
			response, err := call(condition)
			require.Zero(t, response)
			require.ErrorIs(t, err, context.Canceled)
		}
	}
	require.NoError(t, container.database.client.Close())
	for _, call := range calls {
		response, err := call(match)
		require.Zero(t, response)
		var closed *Error
		require.ErrorAs(t, err, &closed)
		require.Equal(t, CodeClientClosed, closed.Code)
	}
}
