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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type failingModel struct{ err error }

func (m failingModel) MarshalJSON() ([]byte, error) { return nil, m.err }

type marshalModelFunc func() ([]byte, error)

func (m marshalModelFunc) MarshalJSON() ([]byte, error) { return m() }

var decodeFailure = errors.New("decode failed")

type failingDecodeModel struct{ Partial bool }

func (m *failingDecodeModel) UnmarshalJSON([]byte) error {
	m.Partial = true
	return decodeFailure
}

// These model-hook tests are serial and join every decoding goroutine before cleanup.
var modelDecodeHook func() error

type callbackModel int

func (m *callbackModel) UnmarshalJSON(body []byte) error {
	if err := modelDecodeHook(); err != nil {
		return err
	}
	var value int
	if err := json.Unmarshal(body, &value); err != nil {
		return err
	}
	*m = callbackModel(value)
	return nil
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
	body, err := marshalItem(document{Name: "example"})
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
		{json.RawMessage(" \n{\"id\":\"one\"}\t"), `{"id":"one"}`},
	} {
		body, err := marshalItem(tc.value)
		require.NoError(t, err)
		require.Equal(t, tc.want, string(body))
	}
}

func TestItemWritesRejectNonObjectBodiesBeforeSubmission(t *testing.T) {
	container := newTestContainer(t)
	container.database.client.beforeItemAcquire = func() { t.Fatal("invalid item must not reach native execution") }
	for _, value := range []any{[]byte("one"), (*struct{})(nil), nil, "text", 42, true, []int{1}} {
		for _, call := range []func() (ItemResponse, error){
			func() (ItemResponse, error) {
				return container.CreateItem(t.Context(), NewPartitionKeyString("one"), "one", value, nil)
			},
			func() (ItemResponse, error) {
				return container.ReplaceItem(t.Context(), NewPartitionKeyString("one"), "one", value, nil)
			},
			func() (ItemResponse, error) {
				return container.UpsertItem(t.Context(), NewPartitionKeyString("one"), "one", value, nil)
			},
		} {
			response, err := call()
			require.Zero(t, response)
			var failure *Error
			require.ErrorAs(t, err, &failure)
			require.Equal(t, CodeSerializationFailed, failure.Code)
			require.ErrorContains(t, err, "valid JSON object")
			require.False(t, failure.FromWire)
		}
	}
}

func TestSerializerValidatesModelHookOutput(t *testing.T) {
	for _, body := range []string{"{", "{} trailing", "{} {}", "null", "[]", `"text"`, "42", "true", " \n\t"} {
		result, err := marshalItem(marshalModelFunc(func() ([]byte, error) { return []byte(body), nil }))
		require.Nil(t, result)
		var failure *Error
		require.ErrorAs(t, err, &failure)
		require.Equal(t, CodeSerializationFailed, failure.Code)
	}
	body, err := marshalItem(marshalModelFunc(func() ([]byte, error) {
		return []byte(`{"id":"one"}`), nil
	}))
	require.NoError(t, err)
	require.Equal(t, `{"id":"one"}`, string(body))
}

func TestSerializerWriteFailureDoesNotAcquireClient(t *testing.T) {
	sentinel := errors.New("private model details")
	container := newTestContainer(t)
	container.database.client.beforeItemAcquire = func() { t.Fatal("failed serialization must not reach native execution") }
	for _, call := range []func() (ItemResponse, error){
		func() (ItemResponse, error) {
			return container.CreateItem(t.Context(), NewPartitionKeyString("pk"), "id", failingModel{sentinel}, nil)
		},
		func() (ItemResponse, error) {
			return container.ReplaceItem(t.Context(), NewPartitionKeyString("pk"), "id", failingModel{sentinel}, nil)
		},
		func() (ItemResponse, error) {
			return container.UpsertItem(t.Context(), NewPartitionKeyString("pk"), "id", failingModel{sentinel}, nil)
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
	_, err := marshalItem(make(chan int))
	var unsupported *json.UnsupportedTypeError
	require.ErrorAs(t, err, &unsupported)
}

func TestItemResponseConversionFailureMetadata(t *testing.T) {
	response := ItemResponse{
		Response: Response{
			StatusCode: 201, SubStatus: 4, AttemptCount: 2, RequestCharge: 5.5,
			ActivityID: "synthetic-activity", Diagnostics: &Diagnostics{StatusCode: 201},
		},
		ETag: "synthetic-etag", SessionToken: "0:5", Value: []byte(`{"id":"one"}`),
		PatchTrackingID: "synthetic-patch-tracking-id",
	}
	model, err := response.IntoModel[failingDecodeModel]()
	require.Zero(t, model)
	require.ErrorIs(t, err, decodeFailure)
	require.Equal(t, &Error{
		Code: CodeSerializationFailed, Message: "decoding item response",
		StatusCode: 201, SubStatus: 4, AttemptCount: 2, RequestCharge: 5.5,
		ActivityID: response.ActivityID, Diagnostics: response.Diagnostics,
		ETag: response.ETag, SessionToken: response.SessionToken, cause: decodeFailure,
		PatchTrackingID: response.PatchTrackingID,
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
	response.Value = []byte{0x80, 1}
	_, err = response.IntoModel[int]()
	require.ErrorContains(t, errors.Unwrap(err), "binary JSON")
}

func TestSerializerCustomConcurrentDecodingAfterClose(t *testing.T) {
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, nil)
	require.NoError(t, err)
	response := ItemResponse{Value: []byte(`"EXAMPLE"`)}
	require.NoError(t, client.Close())
	var workers sync.WaitGroup
	errs := make(chan error, 32)
	for range 32 {
		workers.Go(func() {
			value, err := response.IntoModel[modelName]()
			if err == nil && value != "example" {
				err = fmt.Errorf("decoded %q, want example", value)
			}
			errs <- err
		})
	}
	workers.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}

type modelQueryCursor struct {
	testQueryCursor
	page queryPage
	end  bool
}

func (c *modelQueryCursor) next(context.Context) (queryPage, bool, error) {
	c.calls++
	return c.page, c.end, nil
}

func TestTypedQueryModelsAndMetadata(t *testing.T) {
	page := queryPage{
		Response:     Response{RequestCharge: 2.5, ActivityID: "query", StatusCode: 200, AttemptCount: 2, Diagnostics: &Diagnostics{}},
		SessionToken: "0:3", IndexMetrics: "index", QueryMetrics: "metrics",
		Items: [][]byte{[]byte(`{"name":"first"}`), []byte(`null`)},
	}
	type projection struct{ Name string }
	result, err := decodeModels[*projection](t.Context(), page)
	require.NoError(t, err)
	require.Equal(t, []*projection{{Name: "first"}, nil}, result.Items)
	require.Equal(t, page.Response, result.Response)
	require.Equal(t, page.SessionToken, result.SessionToken)
	require.Equal(t, page.IndexMetrics, result.IndexMetrics)
	require.Equal(t, page.QueryMetrics, result.QueryMetrics)
	page.Items = [][]byte{[]byte(`1`), []byte(`2`)}
	scalars, err := decodeModels[int](t.Context(), page)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, scalars.Items)
}

func TestTypedQueryFailureIsAtomicAndRetainsCheckpoint(t *testing.T) {
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
	require.Zero(t, cursor.closed)
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, "snapshot", token)
	_, err = pager.NextPage(t.Context())
	require.Error(t, err)
	require.NoError(t, pager.Close())
	_, err = pager.ContinuationToken(t.Context())
	require.Error(t, err)
	require.Equal(t, 1, cursor.calls)
	require.Equal(t, 1, cursor.closed)
}

func TestTypedQueryFinalPageFailureRetainsCheckpoint(t *testing.T) {
	cursor := &modelQueryCursor{page: queryPage{Items: [][]byte{[]byte(`"bad"`)}}, end: true}
	pager := &QueryItemsPager[int]{client: newTestContainer(t).database.client, cursor: cursor}
	defer func() { require.NoError(t, pager.Close()) }()
	page, err := pager.NextPage(t.Context())
	require.Error(t, err)
	require.Zero(t, page)
	require.False(t, pager.More())
	require.Zero(t, cursor.closed)
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, "snapshot", token)
}

func TestTypedQueryFailedPageCheckpointErrors(t *testing.T) {
	for _, unsupported := range []bool{false, true} {
		cursor := &modelQueryCursor{page: queryPage{Items: [][]byte{[]byte(`"bad"`)}}}
		if unsupported {
			cursor.checkpointErr = &Error{Code: CodeBadRequest, StatusCode: 400, SubStatus: subStatusBufferedQueryContinuationUnsupported}
		} else {
			cursor.checkpointErr = errors.New("checkpoint failed")
		}
		pager := &QueryItemsPager[int]{client: newTestContainer(t).database.client, cursor: cursor}
		_, err := pager.NextPage(t.Context())
		require.Error(t, err)
		_, err = pager.ContinuationToken(t.Context())
		require.ErrorIs(t, err, cursor.checkpointErr)
		require.False(t, pager.More())
		if unsupported {
			require.Zero(t, cursor.closed)
		} else {
			require.Equal(t, 1, cursor.closed)
		}
		require.NoError(t, pager.Close())
		require.Equal(t, 1, cursor.closed)
	}
}

func TestTypedQueryCancellationDuringConversion(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	client := newTestContainer(t).database.client
	modelDecodeHook = func() error { cancel(); return nil }
	t.Cleanup(func() { modelDecodeHook = nil })
	cursor := &modelQueryCursor{page: queryPage{
		Response: Response{RequestCharge: 3}, Items: [][]byte{[]byte("1")},
	}}
	pager := &QueryItemsPager[callbackModel]{client: client, cursor: cursor}
	page, err := pager.NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
	var failure *Error
	require.ErrorAs(t, err, &failure)
	require.Equal(t, CodeOperationCancelled, failure.Code)
	require.Equal(t, float64(3), failure.RequestCharge)
	require.Equal(t, 1, cursor.closed)
}

func TestTypedQueryCloseDoesNotWaitForConversion(t *testing.T) {
	client := newTestContainer(t).database.client
	entered, release := make(chan struct{}), make(chan struct{})
	unblock := sync.OnceFunc(func() { close(release) })
	t.Cleanup(unblock)
	modelDecodeHook = func() error {
		close(entered)
		<-release
		return nil
	}
	t.Cleanup(func() { modelDecodeHook = nil })
	cursor := &modelQueryCursor{page: queryPage{Items: [][]byte{[]byte("1")}}}
	pager := &QueryItemsPager[callbackModel]{client: client, cursor: cursor}
	nextDone := make(chan error, 1)
	go func() { _, err := pager.NextPage(t.Context()); nextDone <- err }()
	<-entered
	locked := pager.mu.TryLock()
	if locked {
		pager.mu.Unlock()
	}
	require.True(t, locked, "conversion must not hold the checkpoint/close lock")
	closeDone := make(chan error, 1)
	go func() { closeDone <- pager.Close() }()
	clientDone := make(chan error, 1)
	go func() { clientDone <- client.Close() }()
	for _, done := range []chan error{closeDone, clientDone} {
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("shutdown blocked on model conversion")
		}
	}
	unblock()
	require.NoError(t, <-nextDone)
	require.Equal(t, 1, cursor.closed)
}

func TestTypedQueryModelHookReentry(t *testing.T) {
	for _, action := range []string{"more", "next", "checkpoint", "pager close", "client close", "close then fail"} {
		t.Run(action, func(t *testing.T) {
			client := newTestContainer(t).database.client
			cursor := &modelQueryCursor{page: queryPage{Items: [][]byte{[]byte("1")}}}
			pager := &QueryItemsPager[callbackModel]{client: client, cursor: cursor}
			defer func() { require.NoError(t, pager.Close()) }()
			sentinel := errors.New("conversion failed")
			modelDecodeHook = func() error {
				switch action {
				case "more":
					require.True(t, pager.More())
				case "next":
					page, err := pager.NextPage(t.Context())
					require.Zero(t, page)
					require.ErrorContains(t, err, "conversion is in progress")
				case "checkpoint":
					token, err := pager.ContinuationToken(t.Context())
					require.Empty(t, token)
					require.ErrorContains(t, err, "conversion is in progress")
				case "pager close", "close then fail":
					require.NoError(t, pager.Close())
				case "client close":
					require.NoError(t, client.Close())
				}
				if action == "close then fail" {
					return sentinel
				}
				return nil
			}
			t.Cleanup(func() { modelDecodeHook = nil })
			done := make(chan struct{})
			go func() {
				defer close(done)
				page, err := pager.NextPage(t.Context())
				if action == "close then fail" {
					require.ErrorIs(t, err, sentinel)
					require.Zero(t, page)
				} else {
					require.NoError(t, err)
					require.Equal(t, []callbackModel{1}, page.Items)
				}
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("model hook deadlocked")
			}
			require.Equal(t, 1, cursor.calls)
			if action == "pager close" || action == "close then fail" {
				require.False(t, pager.More())
				require.Equal(t, 1, cursor.closed)
				_, err := pager.ContinuationToken(t.Context())
				require.Error(t, err)
			}
		})
	}
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
					if _, err := decodeModels[struct{ Value string }](b.Context(), page); err != nil {
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
