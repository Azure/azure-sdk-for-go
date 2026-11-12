// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadManySelectionOwnsInputsAndParameters(t *testing.T) {
	items := []ItemIdentity{NewItemIdentity(NewPartitionKeyString("pk"), "a")}
	selection := NewReadManyItemsSelection(items...)
	items[0] = NewItemIdentity(NewPartitionKeyString("other"), "b")
	require.Equal(t, "a", selection.identities[0].id)

	value := map[string]any{"v": 1}
	filtered := selection.WithFilter("c.v = @v")
	first, err := filtered.WithParameter("@v", value)
	require.NoError(t, err)
	value["v"] = 2
	second, err := first.WithParameter("@v", 3)
	require.NoError(t, err)
	third, err := first.WithParameter("@w", "x")
	require.NoError(t, err)

	require.Empty(t, selection.filter)
	require.Empty(t, filtered.parameters)
	require.Equal(t, []queryParameter{{Name: "@v", Value: json.RawMessage(`{"v":1}`)}}, first.parameters)
	require.Equal(t, []queryParameter{{Name: "@v", Value: json.RawMessage(`3`)}}, second.parameters)
	require.Equal(t, []queryParameter{{Name: "@v", Value: json.RawMessage(`{"v":1}`)}, {Name: "@w", Value: json.RawMessage(`"x"`)}}, third.parameters)
}

func TestReadManyParameterValidation(t *testing.T) {
	selection := NewReadManyPartitionsSelection(NewPartitionKeyString("pk")).WithFilter("c.v = @v")
	for _, name := range []string{"", "@", "v", "@with space", "@nul\x00"} {
		_, err := selection.WithParameter(name, 1)
		require.Error(t, err, name)
	}
	for _, value := range []any{math.NaN(), make(chan int)} {
		_, err := selection.WithParameter("@v", value)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
		require.NotNil(t, errors.Unwrap(err))
	}
}

func TestReadManyRequestOwnsOptions(t *testing.T) {
	options := &ReadManyOptions{
		Operation:    OperationOptions{ExcludedRegions: []Region{RegionEastUS}},
		SessionToken: "0:1",
		PageSizeHint: 5,
		MaxFanOut:    2,
	}
	req, err := newReadManyRequest(NewReadManyPartitionsSelection(NewPartitionKeyString("pk")), options)
	require.NoError(t, err)
	options.Operation.ExcludedRegions[0] = RegionWestUS
	options.PageSizeHint = 9
	require.Equal(t, []Region{RegionEastUS}, req.options.Operation.ExcludedRegions)
	require.Equal(t, int32(5), req.options.PageSizeHint)
	require.Equal(t, uint32(2), req.options.MaxFanOut)
	require.Equal(t, SessionToken("0:1"), req.options.SessionToken)
}

func TestReadManyPagerArgumentAndLifetimeOrdering(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	pk := NewPartitionKeyString("pk")
	valid := NewReadManyItemsSelection(NewItemIdentity(pk, "a"))
	withParam, err := NewReadManyItemsSelection(NewItemIdentity(pk, "a")).WithParameter("@v", 1)
	require.NoError(t, err)
	for _, tt := range []struct {
		name      string
		selection ReadManySelection
		options   *ReadManyOptions
		message   string
	}{
		{"zero selection", ReadManySelection{}, nil, "NewReadManyItemsSelection"},
		{"empty partition key", NewReadManyItemsSelection(NewItemIdentity(PartitionKey{}, "a")), nil, "partition key components"},
		{"empty partitions key", NewReadManyPartitionsSelection(PartitionKey{}), nil, "partition key components"},
		{"empty id", NewReadManyItemsSelection(NewItemIdentity(pk, "")), nil, "non-empty item id"},
		{"nul id", NewReadManyItemsSelection(NewItemIdentity(pk, "a\x00b")), nil, "NUL"},
		{"nul filter", valid.WithFilter("c.a\x00"), nil, "NUL"},
		{"parameters without filter", withParam, nil, "require a filter"},
		{"negative hint", valid, &ReadManyOptions{PageSizeHint: -1}, "page size hint"},
		{"invalid session", valid, &ReadManyOptions{SessionToken: "a\x00b"}, "session"},
		{"invalid timeout", valid, &ReadManyOptions{Operation: OperationOptions{EndToEndTimeout: to(time.Duration(-1))}}, "EndToEndTimeout"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pager := container.NewReadManyItemsPager[json.RawMessage](tt.selection, tt.options)
			page, err := pager.NextPage(ctx)
			require.ErrorContains(t, err, tt.message)
			require.Zero(t, page)
			require.NotErrorIs(t, err, context.Canceled)
			requireNotDriverUnavailable(t, err)
		})
	}
	page, err := container.NewReadManyItemsPager[json.RawMessage](valid, nil).NextPage(ctx)
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	open := newTestContainer(t)
	page, err = open.NewReadManyItemsPager[json.RawMessage](valid, nil).NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)

	page, err = open.NewReadManyItemsPager[json.RawMessage](valid, nil).NextPage(nil) //nolint:staticcheck // verifies the boundary guard
	require.ErrorContains(t, err, "context must not be nil")
	require.Zero(t, page)
}

func TestReadManyEmptySelectionReadsNothing(t *testing.T) {
	container := newTestContainer(t)
	for _, selection := range []ReadManySelection{NewReadManyItemsSelection(), NewReadManyPartitionsSelection()} {
		pager := container.NewReadManyItemsPager[json.RawMessage](selection, nil)
		require.True(t, pager.More())
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		require.Empty(t, page.Items)
		require.False(t, pager.More())
		require.NoError(t, pager.Close())
	}
}

func TestReadManyPagerCursorLifecycle(t *testing.T) {
	cursor := &testQueryCursor{}
	pager := &ReadManyItemsPager[json.RawMessage]{
		client: newTestContainer(t).database.client,
		req:    readManyRequest{selection: NewReadManyItemsSelection(NewItemIdentity(NewPartitionKeyString("pk"), "a"))},
		cursor: cursor,
	}
	empty, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Empty(t, empty.Items)
	require.True(t, pager.More(), "an empty page must not end read many")
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	_, err = pager.NextPage(t.Context())
	require.NoError(t, err)
	require.False(t, pager.More())
	require.Equal(t, 1, cursor.closed)
	require.NoError(t, pager.Close())
	require.Equal(t, 1, cursor.closed)

	failing := &testQueryCursor{nextErr: &Error{Code: CodeBadRequest}}
	pager = &ReadManyItemsPager[json.RawMessage]{
		client: newTestContainer(t).database.client,
		req:    pager.req,
		cursor: failing,
	}
	_, err = pager.NextPage(t.Context())
	require.ErrorIs(t, err, failing.nextErr)
	require.False(t, pager.More())
	require.Equal(t, 1, failing.closed)
}
