// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type serializerFuncs struct {
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (s *serializerFuncs) Marshal(value any) ([]byte, error) { return s.marshal(value) }
func (s *serializerFuncs) Unmarshal(data []byte, value any) error {
	return s.unmarshal(data, value)
}

type modelName string

func (n modelName) MarshalJSON() ([]byte, error) { return json.Marshal(strings.ToUpper(string(n))) }
func (n *modelName) UnmarshalJSON(body []byte) error {
	var text string
	if err := json.Unmarshal(body, &text); err != nil {
		return err
	}
	*n = modelName(strings.ToLower(text))
	return nil
}

func TestSerializerDefaultModelHooksAndRawMessage(t *testing.T) {
	type document struct {
		Name modelName `json:"customName"`
	}
	body, err := marshalItem(nil, document{Name: "example"})
	require.NoError(t, err)
	require.JSONEq(t, `{"customName":"EXAMPLE"}`, string(body))
	response := ItemResponse{Value: body}
	model, err := response.IntoModel[document]()
	require.NoError(t, err)
	require.Equal(t, document{Name: "example"}, model)
	require.JSONEq(t, `{"customName":"EXAMPLE"}`, string(response.Value))

	for _, tc := range []struct {
		value any
		want  string
	}{
		{json.RawMessage(`{"id":"one"}`), `{"id":"one"}`},
		{[]byte(`one`), `"b25l"`},
		{(*document)(nil), `null`},
	} {
		body, err := marshalItem(nil, tc.value)
		require.NoError(t, err)
		require.Equal(t, tc.want, string(body))
	}
}

func TestSerializerClientValidationAndSnapshot(t *testing.T) {
	var typedNil *serializerFuncs
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &ClientOptions{Serializer: typedNil})
	require.ErrorContains(t, err, "nil implementation")
	require.Nil(t, client)

	codec := &serializerFuncs{marshal: json.Marshal, unmarshal: json.Unmarshal}
	options := &ClientOptions{Serializer: codec}
	client, err = newClient("https://myaccount.documents.azure.com", testAccountKey, nil, options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	options.Serializer = typedNil
	require.Same(t, codec, client.options.Serializer)
	require.Nil(t, newTestContainer(t).database.client.options.Serializer)
}

func TestSerializerWriteFailureDoesNotAcquireClient(t *testing.T) {
	sentinel := errors.New("private model details")
	container := newTestContainer(t)
	container.database.client.beforeItemAcquire = func() { t.Fatal("failed serialization must not reach native execution") }
	container.database.client.options.Serializer = &serializerFuncs{
		marshal:   func(any) ([]byte, error) { return nil, sentinel },
		unmarshal: func([]byte, any) error { t.Fatal("write must not deserialize"); return nil },
	}
	for _, call := range []func() (ItemResponse, error){
		func() (ItemResponse, error) {
			return container.CreateItem(t.Context(), NewPartitionKeyString("pk"), "id", struct{}{}, nil)
		},
		func() (ItemResponse, error) {
			return container.ReplaceItem(t.Context(), NewPartitionKeyString("pk"), "id", struct{}{}, nil)
		},
		func() (ItemResponse, error) {
			return container.UpsertItem(t.Context(), NewPartitionKeyString("pk"), "id", struct{}{}, nil)
		},
	} {
		response, err := call()
		require.Zero(t, response)
		require.ErrorIs(t, err, sentinel)
		var failure *Error
		require.ErrorAs(t, err, &failure)
		require.Equal(t, CodeSerializationFailed, failure.Code)
		require.Zero(t, failure.RequestCharge)
		require.False(t, failure.FromWire)
		require.NotContains(t, err.Error(), sentinel.Error())
	}
	_, err := marshalItem(nil, make(chan int))
	var unsupported *json.UnsupportedTypeError
	require.ErrorAs(t, err, &unsupported)
}

func TestItemResponseConversionFailureMetadata(t *testing.T) {
	sentinel := errors.New("decode failed")
	response := ItemResponse{
		Response: Response{
			StatusCode: 201, SubStatus: 4, AttemptCount: 2, RequestCharge: 5.5,
			ActivityID: "synthetic-activity", Diagnostics: &Diagnostics{StatusCode: 201},
		},
		ETag: "synthetic-etag", SessionToken: "0:5", Value: []byte(`{"id":"one"}`),
		serializer: &serializerFuncs{unmarshal: func(_ []byte, value any) error {
			*value.(*map[string]any) = map[string]any{"partial": true}
			return sentinel
		}},
	}
	model, err := response.IntoModel[map[string]any]()
	require.Nil(t, model)
	require.ErrorIs(t, err, sentinel)
	require.Equal(t, &Error{
		Code: CodeSerializationFailed, Message: "decoding item response",
		StatusCode: 201, SubStatus: 4, AttemptCount: 2, RequestCharge: 5.5,
		ActivityID: response.ActivityID, Diagnostics: response.Diagnostics,
		ETag: response.ETag, SessionToken: response.SessionToken, cause: sentinel,
	}, err)
}

func TestItemResponseMissingMalformedAndNull(t *testing.T) {
	for _, body := range [][]byte{nil, {}, []byte("{"), []byte("1 2"), {0x80, 1}} {
		response := ItemResponse{Value: body}
		model, err := response.IntoModel[any]()
		require.Nil(t, model)
		var failure *Error
		require.ErrorAs(t, err, &failure)
		require.Equal(t, CodeSerializationFailed, failure.Code)
	}
	response := ItemResponse{Value: []byte("null")}
	pointer, err := response.IntoModel[*int]()
	require.NoError(t, err)
	require.Nil(t, pointer)
	scalar, err := response.IntoModel[int]()
	require.NoError(t, err)
	require.Zero(t, scalar)
	var calls int
	response.serializer = &serializerFuncs{unmarshal: func([]byte, any) error { calls++; return nil }}
	response.Value = []byte{0x80, 1}
	_, err = response.IntoModel[int]()
	require.ErrorContains(t, errors.Unwrap(err), "binary JSON")
	require.Zero(t, calls)
}

func TestSerializerCustomConcurrentDecodingAfterClose(t *testing.T) {
	var calls atomic.Int32
	codec := &serializerFuncs{unmarshal: func(body []byte, target any) error {
		calls.Add(1)
		var envelope struct{ Payload int }
		if err := json.Unmarshal(body, &envelope); err != nil {
			return err
		}
		*target.(*int) = envelope.Payload
		return nil
	}}
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, &ClientOptions{Serializer: codec})
	require.NoError(t, err)
	response := ItemResponse{Value: []byte(`{"Payload":42}`), serializer: client.options.Serializer}
	require.NoError(t, client.Close())
	var workers sync.WaitGroup
	errs := make(chan error, 32)
	for range 32 {
		workers.Go(func() {
			value, err := response.IntoModel[int]()
			if err == nil && value != 42 {
				err = fmt.Errorf("decoded %d, want 42", value)
			}
			errs <- err
		})
	}
	workers.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	require.Equal(t, int32(32), calls.Load())
}

type modelQueryCursor struct {
	testQueryCursor
	page queryPage
}

func (c *modelQueryCursor) next(context.Context) (queryPage, bool, error) {
	c.calls++
	return c.page, false, nil
}

func TestTypedQueryModelsAndMetadata(t *testing.T) {
	page := queryPage{
		Response:     Response{RequestCharge: 2.5, ActivityID: "query", StatusCode: 200, AttemptCount: 2, Diagnostics: &Diagnostics{}},
		SessionToken: "0:3", IndexMetrics: "index", QueryMetrics: "metrics",
		Items: [][]byte{[]byte(`{"name":"first"}`), []byte(`null`)},
	}
	type projection struct{ Name string }
	result, err := decodeModels[*projection](t.Context(), page, nil)
	require.NoError(t, err)
	require.Equal(t, []*projection{{Name: "first"}, nil}, result.Items)
	require.Equal(t, page.Response, result.Response)
	require.Equal(t, page.SessionToken, result.SessionToken)
	require.Equal(t, page.IndexMetrics, result.IndexMetrics)
	require.Equal(t, page.QueryMetrics, result.QueryMetrics)
	page.Items = [][]byte{[]byte(`1`), []byte(`2`)}
	scalars, err := decodeModels[int](t.Context(), page, nil)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, scalars.Items)
}

func TestTypedQueryFailureIsAtomicAndCloses(t *testing.T) {
	cursor := &modelQueryCursor{page: queryPage{
		Response:     Response{RequestCharge: 7, ActivityID: "page", StatusCode: 200},
		SessionToken: "0:3", Items: [][]byte{[]byte("1"), []byte(`"not an integer"`), []byte("3")},
	}}
	pager := &QueryItemsPager[int]{client: newTestContainer(t).database.client, cursor: cursor}
	page, err := pager.NextPage(t.Context())
	require.Zero(t, page)
	var failure *Error
	require.ErrorAs(t, err, &failure)
	require.Equal(t, CodeSerializationFailed, failure.Code)
	require.Equal(t, float64(7), failure.RequestCharge)
	require.Equal(t, SessionToken("0:3"), failure.SessionToken)
	require.Contains(t, failure.Message, "index 1")
	var mismatch *json.UnmarshalTypeError
	require.ErrorAs(t, err, &mismatch)
	require.False(t, pager.More())
	require.Equal(t, 1, cursor.closed)
	token, err := pager.ContinuationToken(t.Context())
	require.Error(t, err)
	require.Empty(t, token)
	require.NoError(t, pager.Close())
	require.Equal(t, 1, cursor.calls)
	require.Equal(t, 1, cursor.closed)
}

func TestTypedQueryCancellationDuringConversion(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	client := newTestContainer(t).database.client
	client.options.Serializer = &serializerFuncs{unmarshal: func(body []byte, target any) error {
		cancel()
		return json.Unmarshal(body, target)
	}}
	cursor := &modelQueryCursor{page: queryPage{
		Response: Response{RequestCharge: 3}, Items: [][]byte{[]byte("1")},
	}}
	pager := &QueryItemsPager[int]{client: client, cursor: cursor}
	page, err := pager.NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
	var failure *Error
	require.ErrorAs(t, err, &failure)
	require.Equal(t, CodeOperationCancelled, failure.Code)
	require.Equal(t, float64(3), failure.RequestCharge)
	require.Equal(t, 1, cursor.closed)
}

func TestTypedQueryCloseWaitsForConversion(t *testing.T) {
	client := newTestContainer(t).database.client
	entered, release := make(chan struct{}), make(chan struct{})
	unblock := sync.OnceFunc(func() { close(release) })
	t.Cleanup(unblock)
	client.options.Serializer = &serializerFuncs{unmarshal: func(body []byte, target any) error {
		close(entered)
		<-release
		return json.Unmarshal(body, target)
	}}
	cursor := &modelQueryCursor{page: queryPage{Items: [][]byte{[]byte("1")}}}
	pager := &QueryItemsPager[int]{client: client, cursor: cursor}
	nextDone := make(chan error, 1)
	go func() { _, err := pager.NextPage(t.Context()); nextDone <- err }()
	<-entered
	locked := pager.mu.TryLock()
	if locked {
		pager.mu.Unlock()
	}
	require.False(t, locked, "conversion must hold the checkpoint/close lock")
	closeDone := make(chan error, 1)
	go func() { closeDone <- pager.Close() }()
	clientDone := make(chan error, 1)
	go func() { clientDone <- client.Close() }()
	unblock()
	require.NoError(t, <-nextDone)
	require.NoError(t, <-closeDone)
	require.NoError(t, <-clientDone)
	require.Equal(t, 1, cursor.closed)
}

func BenchmarkTypedModel(b *testing.B) {
	for _, size := range []int{1024, 16384} {
		body := []byte(`{"value":"` + strings.Repeat("x", size) + `"}`)
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			b.ReportAllocs()
			response := ItemResponse{Value: body}
			for b.Loop() {
				_, err := response.IntoModel[struct{ Value string }]()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkManualModel(b *testing.B) {
	for _, size := range []int{1024, 16384} {
		body := []byte(`{"value":"` + strings.Repeat("x", size) + `"}`)
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var model struct{ Value string }
				if err := json.Unmarshal(body, &model); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchmarkQueryModels(b *testing.B, typed bool) {
	for _, count := range []int{10, 100} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			page := queryPage{Items: make([][]byte, count)}
			for i := range page.Items {
				page.Items[i] = []byte(`{"value":"` + strings.Repeat("x", 1024) + `"}`)
			}
			b.ReportAllocs()
			for b.Loop() {
				if typed {
					if _, err := decodeModels[struct{ Value string }](b.Context(), page, nil); err != nil {
						b.Fatal(err)
					}
				} else {
					items := make([]struct{ Value string }, len(page.Items))
					for i, body := range page.Items {
						var model struct{ Value string }
						if err := json.Unmarshal(body, &model); err != nil {
							b.Fatal(err)
						}
						items[i] = model
					}
				}
			}
		})
	}
}

func BenchmarkTypedQuery(b *testing.B)  { benchmarkQueryModels(b, true) }
func BenchmarkManualQuery(b *testing.B) { benchmarkQueryModels(b, false) }
