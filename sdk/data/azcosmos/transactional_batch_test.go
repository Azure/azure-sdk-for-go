// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestTransactionalBatchEncodesSupportedOperations(t *testing.T) {
	batch := NewTransactionalBatch(NewPartitionKeyString("tenant").AppendString("region").AppendNumber(7))
	etag := azcore.ETag(`"original"`)
	require.NoError(t, batch.CreateItem([]byte(`{"id":"created","n":9007199254740993}`), nil))
	require.NoError(t, batch.ReadItem("read", &TransactionalBatchReadItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.UpsertItem([]byte(`{"id":"upserted"}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.ReplaceItem("replaced", []byte(`{"id":"replaced"}`), &TransactionalBatchReplaceItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.DeleteItem("deleted", &TransactionalBatchDeleteItemOptions{IfMatchETag: &etag}))
	body, err := batch.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `[
		{"operationType":"Create","resourceBody":{"id":"created","n":9007199254740993}},
		{"operationType":"Read","id":"read","ifMatch":"\"original\""},
		{"operationType":"Upsert","resourceBody":{"id":"upserted"},"ifMatch":"\"original\""},
		{"operationType":"Replace","id":"replaced","resourceBody":{"id":"replaced"},"ifMatch":"\"original\""},
		{"operationType":"Delete","id":"deleted","ifMatch":"\"original\""}
	]`, string(body))
	require.Contains(t, string(body), "9007199254740993", "item numbers must not pass through float64")
	require.Equal(t, 3, batch.partitionKey.Len())
}

func TestTransactionalBatchOwnsAppendSnapshots(t *testing.T) {
	item := []byte(`{"id":"item","value":"original"}`)
	etag := azcore.ETag(`"original"`)
	options := &TransactionalBatchReplaceItemOptions{IfMatchETag: &etag}
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.ReplaceItem("item", item, options))
	copy(item, []byte(`{"id":"item","value":"changed!"}`))
	etag = `"changed"`
	options.IfMatchETag = nil
	body, err := batch.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `[{"operationType":"Replace","id":"item","resourceBody":{"id":"item","value":"original"},"ifMatch":"\"original\""}]`, string(body))
}

func TestTransactionalBatchReadConditionsOwnSnapshots(t *testing.T) {
	for _, test := range []struct {
		name    string
		options func(*azcore.ETag) *TransactionalBatchReadItemOptions
		body    string
	}{
		{"nil options", func(*azcore.ETag) *TransactionalBatchReadItemOptions { return nil }, `[{"operationType":"Read","id":"item"}]`},
		{"default options", func(*azcore.ETag) *TransactionalBatchReadItemOptions {
			return &TransactionalBatchReadItemOptions{}
		}, `[{"operationType":"Read","id":"item"}]`},
		{"If-Match", func(etag *azcore.ETag) *TransactionalBatchReadItemOptions {
			return &TransactionalBatchReadItemOptions{IfMatchETag: etag}
		}, `[{"operationType":"Read","id":"item","ifMatch":"\"original\""}]`},
		{"If-None-Match", func(etag *azcore.ETag) *TransactionalBatchReadItemOptions {
			return &TransactionalBatchReadItemOptions{IfNoneMatchETag: etag}
		}, `[{"operationType":"Read","id":"item","ifNoneMatch":"\"original\""}]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			etag := azcore.ETag(`"original"`)
			options := test.options(&etag)
			batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
			require.NoError(t, batch.ReadItem("item", options))
			etag = `"changed"`
			if options != nil {
				options.IfMatchETag, options.IfNoneMatchETag = nil, nil
			}
			body, err := batch.marshal()
			require.NoError(t, err)
			require.JSONEq(t, test.body, string(body))
		})
	}
}

func TestTransactionalBatchCopiesAppendIndependently(t *testing.T) {
	original := NewTransactionalBatch(NewPartitionKeyNull())
	require.NoError(t, original.ReadItem("first", nil))
	copied := original
	require.NoError(t, original.DeleteItem("original", nil))
	require.NoError(t, copied.ReadItem("copied", nil))
	first, err := original.marshal()
	require.NoError(t, err)
	second, err := copied.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `[{"operationType":"Read","id":"first"},{"operationType":"Delete","id":"original"}]`, string(first))
	require.JSONEq(t, `[{"operationType":"Read","id":"first"},{"operationType":"Read","id":"copied"}]`, string(second))
}

func TestTransactionalBatchEscapesEnvelopeStrings(t *testing.T) {
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	id := "quoted\"\\\n\u2603"
	etag := azcore.ETag("etag\"\\\n\u2603")
	require.NoError(t, batch.ReadItem(id, &TransactionalBatchReadItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.ReadItem(id, &TransactionalBatchReadItemOptions{IfNoneMatchETag: &etag}))
	body, err := batch.marshal()
	require.NoError(t, err)
	var operations []transactionalBatchOperation
	require.NoError(t, json.Unmarshal(body, &operations))
	require.Equal(t, id, operations[0].ID)
	require.Equal(t, string(etag), operations[0].IfMatch)
	require.Equal(t, id, operations[1].ID)
	require.Equal(t, string(etag), operations[1].IfNoneMatch)
}

func TestTransactionalBatchInvalidAppendIsSticky(t *testing.T) {
	for _, test := range []struct {
		name   string
		append func(*TransactionalBatch) error
	}{
		{"nil body", func(b *TransactionalBatch) error { return b.CreateItem(nil, nil) }},
		{"empty body", func(b *TransactionalBatch) error { return b.UpsertItem([]byte{}, nil) }},
		{"invalid JSON", func(b *TransactionalBatch) error { return b.CreateItem([]byte(`{"id":`), nil) }},
		{"multiple JSON values", func(b *TransactionalBatch) error { return b.CreateItem([]byte(`{} {}`), nil) }},
		{"invalid UTF-8 body", func(b *TransactionalBatch) error { return b.CreateItem([]byte{'"', 0xff, '"'}, nil) }},
		{"empty read ID", func(b *TransactionalBatch) error { return b.ReadItem("", nil) }},
		{"empty replace ID", func(b *TransactionalBatch) error { return b.ReplaceItem("", []byte(`{}`), nil) }},
		{"empty delete ID", func(b *TransactionalBatch) error { return b.DeleteItem("", nil) }},
		{"NUL ID", func(b *TransactionalBatch) error { return b.ReadItem("a\x00b", nil) }},
		{"invalid UTF-8 ID", func(b *TransactionalBatch) error { return b.ReadItem(string([]byte{0xff}), nil) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
			require.NoError(t, batch.ReadItem("first", nil))
			err := test.append(&batch)
			require.Error(t, err)
			require.Equal(t, err, batch.DeleteItem("later", nil))
			body, executeErr := batch.marshal()
			require.Nil(t, body)
			require.Equal(t, err, executeErr)
			require.Len(t, batch.operations, 1, "an ignored append failure must not produce a partial transaction")
		})
	}
	var batch *TransactionalBatch
	require.Error(t, batch.CreateItem([]byte(`{}`), nil))
	require.Error(t, batch.ReadItem("item", nil))
	require.Error(t, batch.UpsertItem([]byte(`{}`), nil))
	require.Error(t, batch.ReplaceItem("item", []byte(`{}`), nil))
	require.Error(t, batch.DeleteItem("item", nil))
}

func TestTransactionalBatchValidatesConditions(t *testing.T) {
	for _, value := range []string{"", "a\x00b", string([]byte{0xff})} {
		t.Run(value, func(t *testing.T) {
			etag := azcore.ETag(value)
			constructors := []func(*TransactionalBatch) error{
				func(b *TransactionalBatch) error {
					return b.ReadItem("item", &TransactionalBatchReadItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.ReadItem("item", &TransactionalBatchReadItemOptions{IfNoneMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.ReplaceItem("item", []byte(`{}`), &TransactionalBatchReplaceItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.DeleteItem("item", &TransactionalBatchDeleteItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.UpsertItem([]byte(`{}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.UpsertItem([]byte(`{}`), &TransactionalBatchUpsertItemOptions{IfNoneMatchETag: &etag})
				},
			}
			for _, appendOperation := range constructors {
				batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
				require.Error(t, appendOperation(&batch))
			}
		})
	}
	etag := azcore.ETag("*")
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.ReadItem("first", nil))
	err := batch.ReadItem("item", &TransactionalBatchReadItemOptions{IfMatchETag: &etag, IfNoneMatchETag: &etag})
	require.ErrorContains(t, err, "must not both be set")
	body, marshalErr := batch.marshal()
	require.Nil(t, body)
	require.Equal(t, err, marshalErr, "an invalid conditional read must not submit the preceding operation")
	batch = NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.Error(t, batch.UpsertItem([]byte(`{}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag, IfNoneMatchETag: &etag}))
	batch = NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.UpsertItem([]byte(`{"id":"new"}`), &TransactionalBatchUpsertItemOptions{IfNoneMatchETag: &etag}))
	body, err = batch.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `[{"operationType":"Upsert","resourceBody":{"id":"new"},"ifNoneMatch":"*"}]`, string(body))
}

func TestTransactionalBatchOperationCountLimits(t *testing.T) {
	batch := NewTransactionalBatch(NewPartitionKeyBool(false))
	_, err := batch.marshal()
	require.ErrorContains(t, err, "between 1 and 100")
	for range maxTransactionalBatchOperations {
		require.NoError(t, batch.ReadItem("item", nil))
	}
	body, err := batch.marshal()
	require.NoError(t, err)
	var operations []transactionalBatchOperation
	require.NoError(t, json.Unmarshal(body, &operations))
	require.Len(t, operations, maxTransactionalBatchOperations)
	require.ErrorContains(t, batch.ReadItem("too-many", nil), "must not exceed 100")
	_, err = batch.marshal()
	require.ErrorContains(t, err, "must not exceed 100")
}

func TestTransactionalBatchEncodedPayloadLimits(t *testing.T) {
	base := NewTransactionalBatch(NewPartitionKeyUndefined())
	require.NoError(t, base.CreateItem([]byte(`{"id":"item","value":""}`), nil))
	encoded, err := base.marshal()
	require.NoError(t, err)
	for _, delta := range []int{-1, 0, 1} {
		t.Run(string(rune('b'+delta)), func(t *testing.T) {
			batch := NewTransactionalBatch(NewPartitionKeyUndefined())
			item := []byte(`{"id":"item","value":"` + strings.Repeat("x", maxTransactionalBatchBytes-len(encoded)+delta) + `"}`)
			require.NoError(t, batch.CreateItem(item, nil))
			body, err := batch.marshal()
			if delta > 0 {
				require.ErrorContains(t, err, "2097152 bytes")
				require.Nil(t, body)
			} else {
				require.NoError(t, err)
				require.Len(t, body, maxTransactionalBatchBytes+delta)
			}
		})
	}
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.CreateItem([]byte(strings.Repeat(" ", maxTransactionalBatchBytes+1)+`{"id":"item"}`), nil))
	body, err := batch.marshal()
	require.NoError(t, err, "the limit applies to the encoded envelope, not the unencoded input length")
	require.JSONEq(t, `[{"operationType":"Create","resourceBody":{"id":"item"}}]`, string(body))
}

func TestTransactionalBatchValidatesPartitionKeyValues(t *testing.T) {
	for _, key := range []PartitionKey{
		{}, NewPartitionKeyNumber(math.NaN()), NewPartitionKeyNumber(math.Inf(1)),
		NewPartitionKeyNumber(math.Inf(-1)), NewPartitionKeyString(string([]byte{0xff})),
	} {
		batch := NewTransactionalBatch(key)
		require.NoError(t, batch.ReadItem("item", nil))
		_, err := batch.marshal()
		require.Error(t, err)
	}
	for _, key := range []PartitionKey{
		NewPartitionKeyNull(), NewPartitionKeyUndefined(), NewPartitionKeyBool(false),
		NewPartitionKeyNumber(0), NewPartitionKeyString("a\x00b"),
		NewPartitionKeyString("tenant").AppendNull().AppendUndefined(),
	} {
		batch := NewTransactionalBatch(key)
		require.NoError(t, batch.ReadItem("item", nil))
		_, err := batch.marshal()
		require.NoError(t, err)
	}
}

func TestTransactionalBatchRequestOwnsOptions(t *testing.T) {
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.ReadItem("item", nil))
	enabled := true
	options := &TransactionalBatchOptions{
		Operation: OperationOptions{
			EnableContentResponseOnWrite: &enabled,
			ExcludedRegions:              []Region{RegionEastUS},
			EndToEndTimeout:              time.Second,
		},
		SessionToken: "0:1",
	}
	req, err := newTransactionalBatchRequest(batch, options)
	require.NoError(t, err)
	enabled = false
	options.Operation.ExcludedRegions[0] = RegionWestUS
	options.SessionToken = "different"
	require.True(t, *req.options.EnableContentResponseOnWrite)
	require.Equal(t, []Region{RegionEastUS}, req.options.ExcludedRegions)
	require.Equal(t, SessionToken("0:1"), req.sessionToken)
	require.Equal(t, time.Second, req.options.EndToEndTimeout)
	require.Empty(t, req.itemID)
	require.Equal(t, operationKindBatch, req.kind)
	require.NoError(t, batch.DeleteItem("later", nil))
	require.JSONEq(t, `[{"operationType":"Read","id":"item"}]`, string(req.body))
}

func TestTransactionalBatchArgumentLifetimeAndContextOrdering(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
	response, err := container.ExecuteTransactionalBatch(ctx, batch, nil)
	require.Zero(t, response)
	require.ErrorContains(t, err, "between 1 and 100")
	require.NotErrorIs(t, err, context.Canceled)
	require.NoError(t, batch.ReadItem("item", nil))
	response, err = container.ExecuteTransactionalBatch(ctx, batch, nil)
	require.Zero(t, response)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	container = newTestContainer(t)
	response, err = container.ExecuteTransactionalBatch(ctx, batch, nil)
	require.Zero(t, response)
	require.ErrorIs(t, err, context.Canceled)
	for _, options := range []*TransactionalBatchOptions{
		{SessionToken: "a\x00b"},
		{Operation: OperationOptions{ConsistencyStrategy: ReadConsistencyStrategy("invalid")}},
	} {
		response, err = container.ExecuteTransactionalBatch(ctx, batch, options)
		require.Zero(t, response)
		require.Error(t, err)
		require.NotErrorIs(t, err, context.Canceled)
	}
	if !driverAvailable {
		response, err = container.ExecuteTransactionalBatch(t.Context(), batch, nil)
		require.Zero(t, response)
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeClientError, cosmosErr.Code)
		require.Contains(t, cosmosErr.Message, "this build cannot reach the Cosmos driver")
	}
}
