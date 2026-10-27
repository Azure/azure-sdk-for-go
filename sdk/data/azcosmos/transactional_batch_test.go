// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
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
	require.NoError(t, batch.CreateItem(json.RawMessage(`{"id":"created","n":9007199254740993}`), nil))
	require.NoError(t, batch.ReadItem("read", &TransactionalBatchReadItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.UpsertItem(json.RawMessage(`{"id":"upserted"}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag}))
	require.NoError(t, batch.ReplaceItem("replaced", json.RawMessage(`{"id":"replaced"}`), &TransactionalBatchReplaceItemOptions{IfMatchETag: &etag}))
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
	item := json.RawMessage(`{"id":"item","value":"original"}`)
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

func TestTransactionalBatchWriteSerializationMatchesPointWrites(t *testing.T) {
	for _, value := range []any{
		map[string]any{"id": "item", "value": 42},
		json.RawMessage(`{"id":"item","value":42}`),
		[]byte(`{"id":"item","value":42}`),
		(*struct{})(nil), "text", []int{1},
	} {
		expected, expectedErr := marshalItem(nil, value)
		for _, appendItem := range []func(*TransactionalBatch) error{
			func(b *TransactionalBatch) error { return b.CreateItem(value, nil) },
			func(b *TransactionalBatch) error { return b.UpsertItem(value, nil) },
			func(b *TransactionalBatch) error { return b.ReplaceItem("item", value, nil) },
		} {
			batch := NewTransactionalBatch(NewPartitionKeyString("pk"))
			err := appendItem(&batch)
			if expectedErr != nil {
				var failure *Error
				require.ErrorAs(t, err, &failure)
				require.Equal(t, CodeSerializationFailed, failure.Code)
				require.Empty(t, batch.operations)
				require.Same(t, err, batch.DeleteItem("later", nil))
				continue
			}
			require.NoError(t, err)
			require.Equal(t, json.RawMessage(expected), batch.operations[0].ResourceBody)
		}
	}
}

func TestTransactionalBatchCapturesSerializerAndOwnsSerializedModels(t *testing.T) {
	container := newTestContainer(t)
	var calls int
	var buffer []byte
	codec := &serializerFuncs{marshal: func(value any) ([]byte, error) {
		calls++
		encoded, err := json.Marshal(map[string]any{"id": "item", "custom": value})
		buffer = append(buffer[:0], encoded...)
		return buffer, err
	}}
	container.database.client.options.Serializer = codec
	batch := container.NewTransactionalBatch(NewPartitionKeyString("pk"))
	container.database.client.options.Serializer = &serializerFuncs{marshal: func(any) ([]byte, error) {
		t.Fatal("batch must retain its construction-time serializer")
		return nil, nil
	}}
	model := map[string]string{"name": "original"}
	require.NoError(t, batch.CreateItem(model, nil))
	model["name"] = "changed"
	require.NoError(t, batch.UpsertItem("second", nil))
	require.NoError(t, batch.ReplaceItem("item", []byte("third"), nil))
	clear(buffer)
	require.Equal(t, 3, calls)
	require.JSONEq(t, `{"id":"item","custom":{"name":"original"}}`, string(batch.operations[0].ResourceBody))
	require.JSONEq(t, `{"id":"item","custom":"second"}`, string(batch.operations[1].ResourceBody))
	require.JSONEq(t, `{"id":"item","custom":"dGhpcmQ="}`, string(batch.operations[2].ResourceBody))
	first, err := batch.marshal()
	require.NoError(t, err)
	second, err := batch.marshal()
	require.NoError(t, err)
	require.Equal(t, first, second)
	require.Equal(t, 3, calls, "encoding the envelope must not invoke the application serializer")
}

func TestTransactionalBatchSerializerFailurePoisonsBuilder(t *testing.T) {
	container := newTestContainer(t)
	sentinel := errors.New("custom serialization failed")
	calls := 0
	container.database.client.options.Serializer = &serializerFuncs{marshal: func(any) ([]byte, error) {
		calls++
		return nil, sentinel
	}}
	for _, appendItem := range []func(*TransactionalBatch) error{
		func(b *TransactionalBatch) error { return b.CreateItem(struct{}{}, nil) },
		func(b *TransactionalBatch) error { return b.UpsertItem(struct{}{}, nil) },
		func(b *TransactionalBatch) error { return b.ReplaceItem("item", struct{}{}, nil) },
	} {
		batch := container.NewTransactionalBatch(NewPartitionKeyString("pk"))
		err := appendItem(&batch)
		require.ErrorIs(t, err, sentinel)
		require.Same(t, err, batch.CreateItem(struct{}{}, nil))
		body, marshalErr := batch.marshal()
		require.Nil(t, body)
		require.Same(t, err, marshalErr)
	}
	require.Equal(t, 3, calls)
}

func TestTransactionalBatchInvalidAppendIsSticky(t *testing.T) {
	for _, test := range []struct {
		name   string
		append func(*TransactionalBatch) error
	}{
		{"nil body", func(b *TransactionalBatch) error { return b.CreateItem[any](nil, nil) }},
		{"empty body", func(b *TransactionalBatch) error { return b.UpsertItem(json.RawMessage{}, nil) }},
		{"invalid JSON", func(b *TransactionalBatch) error { return b.CreateItem(json.RawMessage(`{"id":`), nil) }},
		{"multiple JSON values", func(b *TransactionalBatch) error { return b.CreateItem(json.RawMessage(`{} {}`), nil) }},
		{"non-object body", func(b *TransactionalBatch) error { return b.CreateItem(json.RawMessage(`"text"`), nil) }},
		{"invalid UTF-8 body", func(b *TransactionalBatch) error { return b.CreateItem(json.RawMessage("{\"id\":\"\xff\"}"), nil) }},
		{"empty read ID", func(b *TransactionalBatch) error { return b.ReadItem("", nil) }},
		{"empty replace ID", func(b *TransactionalBatch) error { return b.ReplaceItem("", json.RawMessage(`{}`), nil) }},
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
	require.Error(t, batch.CreateItem(json.RawMessage(`{}`), nil))
	require.Error(t, batch.ReadItem("item", nil))
	require.Error(t, batch.UpsertItem(json.RawMessage(`{}`), nil))
	require.Error(t, batch.ReplaceItem("item", json.RawMessage(`{}`), nil))
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
					return b.ReplaceItem("item", json.RawMessage(`{}`), &TransactionalBatchReplaceItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.DeleteItem("item", &TransactionalBatchDeleteItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.UpsertItem(json.RawMessage(`{}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag})
				},
				func(b *TransactionalBatch) error {
					return b.UpsertItem(json.RawMessage(`{}`), &TransactionalBatchUpsertItemOptions{IfNoneMatchETag: &etag})
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
	require.Error(t, batch.UpsertItem(json.RawMessage(`{}`), &TransactionalBatchUpsertItemOptions{IfMatchETag: &etag, IfNoneMatchETag: &etag}))
	batch = NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, batch.UpsertItem(json.RawMessage(`{"id":"new"}`), &TransactionalBatchUpsertItemOptions{IfNoneMatchETag: &etag}))
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
	require.NoError(t, base.CreateItem(json.RawMessage(`{"id":"item","value":""}`), nil))
	encoded, err := base.marshal()
	require.NoError(t, err)
	for _, delta := range []int{-1, 0, 1} {
		t.Run(string(rune('b'+delta)), func(t *testing.T) {
			batch := NewTransactionalBatch(NewPartitionKeyUndefined())
			item := json.RawMessage(`{"id":"item","value":"` + strings.Repeat("x", maxTransactionalBatchBytes-len(encoded)+delta) + `"}`)
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
	require.NoError(t, batch.CreateItem(json.RawMessage(strings.Repeat(" ", maxTransactionalBatchBytes+1)+`{"id":"item"}`), nil))
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
			EndToEndTimeout:              to(time.Second),
			CustomHeaders:                map[string]string{"x-test": "original"},
			BinaryEncoding:               &BinaryEncodingOptions{Enabled: to(false)},
			ThrottlingRetry:              ThrottlingRetryOptions{MaxRetryCount: to(uint32(2))},
		},
		SessionToken: "0:1",
	}
	req, err := newTransactionalBatchRequest(batch, options)
	require.NoError(t, err)
	enabled = false
	options.Operation.ExcludedRegions[0] = RegionWestUS
	*options.Operation.EndToEndTimeout = 2 * time.Second
	options.Operation.CustomHeaders["x-test"] = "changed"
	*options.Operation.BinaryEncoding.Enabled = true
	*options.Operation.ThrottlingRetry.MaxRetryCount = 3
	options.SessionToken = "different"
	require.True(t, *req.options.EnableContentResponseOnWrite)
	require.Equal(t, []Region{RegionEastUS}, req.options.ExcludedRegions)
	require.Equal(t, SessionToken("0:1"), req.sessionToken)
	require.Equal(t, time.Second, *req.options.EndToEndTimeout)
	require.Equal(t, map[string]string{"x-test": "original"}, req.options.CustomHeaders)
	require.False(t, *req.options.BinaryEncoding.Enabled)
	require.Equal(t, uint32(2), *req.options.ThrottlingRetry.MaxRetryCount)
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
		{Operation: OperationOptions{EndToEndTimeout: to(-time.Second)}},
		{Operation: OperationOptions{CustomHeaders: map[string]string{"x-test": "invalid\nvalue"}}},
		{Operation: OperationOptions{AvailabilityStrategy: HedgingAvailability(0)}},
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

func TestTransactionalBatchRejectsNilContextAfterArgumentValidation(t *testing.T) {
	valid := NewTransactionalBatch(NewPartitionKeyString("pk"))
	require.NoError(t, valid.ReadItem("item", nil))
	for _, test := range []struct {
		name    string
		batch   TransactionalBatch
		options *TransactionalBatchOptions
		closed  bool
		want    string
	}{
		{name: "open client", batch: valid, want: "context must not be nil"},
		{name: "closed client", batch: valid, closed: true, want: "context must not be nil"},
		{name: "invalid batch", batch: NewTransactionalBatch(NewPartitionKeyString("pk")), want: "between 1 and 100"},
		{name: "invalid options", batch: valid,
			options: &TransactionalBatchOptions{Operation: OperationOptions{EndToEndTimeout: to(-time.Second)}},
			want:    "EndToEndTimeout must not be negative"},
	} {
		t.Run(test.name, func(t *testing.T) {
			container := newTestContainer(t)
			if test.closed {
				require.NoError(t, container.database.client.Close())
			}
			response, err := container.ExecuteTransactionalBatch(nil, test.batch, test.options) //nolint:staticcheck // verifies the boundary guard
			require.Zero(t, response)
			require.ErrorContains(t, err, test.want)
		})
	}
}
