// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmulatorModelWorkflow(t *testing.T) {
	type document struct {
		ID    string    `json:"id"`
		PK    string    `json:"pk"`
		Value int       `json:"value"`
		Name  modelName `json:"name"`
	}
	t.Run("model hooks", func(t *testing.T) {
		client, databaseID, containerID := emulatorClient(t)
		container, err := client.NewContainer(t.Context(), databaseID, containerID)
		require.NoError(t, err)
		id := uniqueItemID(t)
		pk := NewPartitionKeyString(id)
		item := document{ID: id, PK: id, Value: 1, Name: "example"}
		content := true
		write := OperationOptions{EnableContentResponseOnWrite: &content}
		created, err := container.CreateItem(t.Context(), pk, id, item, &CreateItemOptions{Operation: write})
		require.NoError(t, err)
		deletedItem := false
		t.Cleanup(func() {
			if !deletedItem {
				_, err := container.DeleteItem(context.Background(), pk, id, nil)
				require.NoError(t, err)
			}
		})
		require.Positive(t, created.RequestCharge)
		require.NotEmpty(t, created.ETag)
		require.NotEmpty(t, created.SessionToken)
		model, err := created.IntoModel[document]()
		require.NoError(t, err)
		require.Equal(t, item, model)

		read, err := container.ReadItem(t.Context(), pk, id, nil)
		require.NoError(t, err)
		model, err = read.IntoModel[document]()
		require.NoError(t, err)
		require.Equal(t, item, model)
		item.Value = 2
		readCondition, err := IfMatch(read.ETag)
		require.NoError(t, err)
		replaced, err := container.ReplaceItem(t.Context(), pk, id, item, &ReplaceItemOptions{
			Operation: write, Precondition: readCondition,
		})
		require.NoError(t, err)
		model, err = replaced.IntoModel[document]()
		require.NoError(t, err)
		require.Equal(t, item, model)
		item.Value = 3
		upserted, err := container.UpsertItem(t.Context(), pk, id, item, &UpsertItemOptions{Operation: write})
		require.NoError(t, err)
		model, err = upserted.IntoModel[document]()
		require.NoError(t, err)
		require.Equal(t, item, model)

		var patch PatchOperations
		require.NoError(t, patch.AppendSet("/value", 4))
		patched, err := container.PatchItem(t.Context(), pk, id, patch, &PatchItemOptions{Operation: write})
		require.NoError(t, err)
		item.Value = 4
		model, err = patched.IntoModel[document]()
		require.NoError(t, err)
		require.Equal(t, item, model)

		scope := NewFeedScopeForPartitionKey(pk)
		pager := container.NewQueryItemsPager[int](NewQuery("SELECT VALUE c.value FROM c"), scope, nil)
		t.Cleanup(func() { require.NoError(t, pager.Close()) })
		var values []int
		for pager.More() {
			page, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			values = append(values, page.Items...)
		}
		require.Equal(t, []int{4}, values)

		failed := container.NewQueryItemsPager[int](NewQuery("SELECT * FROM c"), scope, nil)
		page, err := failed.NextPage(t.Context())
		require.Zero(t, page)
		var failure *Error
		require.ErrorAs(t, err, &failure)
		require.Equal(t, CodeSerializationFailed, failure.Code)
		require.Positive(t, failure.RequestCharge)
		require.False(t, failed.More())
		require.Nil(t, failed.cursor)
		require.NoError(t, failed.Close())

		noContent := false
		suppressed, err := container.ReplaceItem(t.Context(), pk, id, item, &ReplaceItemOptions{
			Operation: OperationOptions{EnableContentResponseOnWrite: &noContent},
		})
		require.NoError(t, err)
		require.Empty(t, suppressed.Value)
		_, err = suppressed.IntoModel[document]()
		require.Error(t, err)

		deleted, err := container.DeleteItem(t.Context(), pk, id, nil)
		require.NoError(t, err)
		deletedItem = true
		require.Empty(t, deleted.Value)
		_, err = deleted.IntoModel[document]()
		require.Error(t, err)
		require.NoError(t, client.Close())
		model, err = read.IntoModel[document]()
		require.NoError(t, err, "response decoding must outlive native handles")
		require.Equal(t, 1, model.Value)
	})
}
