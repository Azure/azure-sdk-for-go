// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore Commited NULE

package azcosmos

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func newTestContainer(t *testing.T) *ContainerClient {
	t.Helper()

	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	return container
}

func requireNotDriverUnavailable(t *testing.T, err error) {
	t.Helper()
	var cosmosErr *Error
	if errors.As(err, &cosmosErr) {
		require.NotContains(t, cosmosErr.Message, "cannot reach the Cosmos driver",
			"argument validation should run before the operation is attempted")
	}
}

func TestReadItemRejectsEmptyID(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "", nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)
}

func TestCreateItemRejectsEmptyItem(t *testing.T) {
	container := newTestContainer(t)

	for _, item := range [][]byte{nil, {}} {
		_, err := container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "item-1", item, nil)
		require.Error(t, err)
		requireNotDriverUnavailable(t, err)
	}
}

// The id addresses the item being created, so an empty one is a caller mistake that has to be
// reported as itself rather than reaching the driver as a null argument.
func TestCreateItemRejectsEmptyID(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)
}

func TestItemWriteOperationsRejectEmptyID(t *testing.T) {
	container := newTestContainer(t)
	item := []byte(`{"id":"item-1","pk":"pk"}`)
	patch := validPatchOperations(t)

	for _, tt := range []struct {
		name string
		call func() error
	}{
		{"replace", func() error {
			_, err := container.ReplaceItem(context.Background(), NewPartitionKeyString("pk"), "", item, nil)
			return err
		}},
		{"upsert", func() error {
			_, err := container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), "", item, nil)
			return err
		}},
		{"delete", func() error {
			_, err := container.DeleteItem(context.Background(), NewPartitionKeyString("pk"), "", nil)
			return err
		}},
		{"patch", func() error {
			_, err := container.PatchItem(context.Background(), NewPartitionKeyString("pk"), "", patch, nil)
			return err
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			require.ErrorContains(t, err, "item id must not be empty")
			requireNotDriverUnavailable(t, err)
		})
	}
}

func TestReplaceAndUpsertRejectEmptyItem(t *testing.T) {
	container := newTestContainer(t)

	for _, item := range [][]byte{nil, {}} {
		_, err := container.ReplaceItem(context.Background(), NewPartitionKeyString("pk"), "item-1", item, nil)
		require.ErrorContains(t, err, "item must not be empty")
		requireNotDriverUnavailable(t, err)

		_, err = container.UpsertItem(context.Background(), NewPartitionKeyString("pk"), "item-1", item, nil)
		require.ErrorContains(t, err, "item must not be empty")
		requireNotDriverUnavailable(t, err)
	}
}

func TestPatchItemRejectsInvalidOperations(t *testing.T) {
	container := newTestContainer(t)
	pk := NewPartitionKeyString("pk")

	_, err := container.PatchItem(context.Background(), pk, "item-1", PatchOperations{}, nil)
	require.ErrorContains(t, err, "at least one operation")
	requireNotDriverUnavailable(t, err)

	var operations PatchOperations
	appendErr := operations.AppendRemove("not-a-pointer")
	require.Error(t, appendErr)
	_, err = container.PatchItem(context.Background(), pk, "item-1", operations, nil)
	require.ErrorIs(t, err, appendErr)
	requireNotDriverUnavailable(t, err)
}

// Argument validation runs before the context is consulted, so a caller's deterministic mistake is
// reported as itself rather than being masked by a deadline that happened to fire first.
func TestItemOperationsValidateArgumentsBeforeContext(t *testing.T) {
	container := newTestContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "", nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "item-1", nil, nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.ReplaceItem(ctx, NewPartitionKeyString("pk"), "item-1", nil, nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.UpsertItem(ctx, NewPartitionKeyString("pk"), "item-1", nil, nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.DeleteItem(ctx, NewPartitionKeyString("pk"), "", nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.PatchItem(ctx, NewPartitionKeyString("pk"), "item-1", PatchOperations{}, nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)
}

// A partition key with no components is always a caller mistake: no container has a definition
// with zero paths, and PartitionKey has a usable zero value, so "declared but never appended to"
// is an easy thing to get wrong.
func TestItemOperationsRejectEmptyPartitionKey(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.ReadItem(context.Background(), PartitionKey{}, "item-1", nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)

	_, err = container.CreateItem(context.Background(), PartitionKey{}, "item-1", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)

	_, err = container.ReplaceItem(context.Background(), PartitionKey{}, "item-1", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)

	_, err = container.UpsertItem(context.Background(), PartitionKey{}, "item-1", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)

	_, err = container.DeleteItem(context.Background(), PartitionKey{}, "item-1", nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)

	_, err = container.PatchItem(context.Background(), PartitionKey{}, "item-1", validPatchOperations(t), nil)
	require.Error(t, err)
	requireNotDriverUnavailable(t, err)
}

// An already-cancelled context must be honored rather than starting work that is bound to fail,
// and the caller must get the context's own error so errors.Is against context.Canceled works.
func TestItemOperationsHonorCancelledContext(t *testing.T) {
	container := newTestContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "item-1", nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "item-1", []byte(`{"id":"item-1"}`), nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.ReplaceItem(ctx, NewPartitionKeyString("pk"), "item-1", []byte(`{"id":"item-1"}`), nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.UpsertItem(ctx, NewPartitionKeyString("pk"), "item-1", []byte(`{"id":"item-1"}`), nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.DeleteItem(ctx, NewPartitionKeyString("pk"), "item-1", nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.PatchItem(ctx, NewPartitionKeyString("pk"), "item-1", validPatchOperations(t), nil)
	require.ErrorIs(t, err, context.Canceled)
}

// A well-formed call reaches the driver in a driver-backed build and reports that it is not
// implemented otherwise. Either way it fails as an *Error, so the documented errors.As idiom works
// whichever build a caller has.
func TestItemOperationsReportErrorsAsCosmosErrors(t *testing.T) {
	if driverAvailable {
		// The driver would try to reach the endpoint, which is what the emulator tests cover.
		t.Skip("driver-backed build: covered by the emulator tests")
	}

	container := newTestContainer(t)
	pk := NewPartitionKeyString("pk")

	_, readErr := container.ReadItem(context.Background(), pk, "item-1", nil)
	_, createErr := container.CreateItem(context.Background(), pk, "item-1", []byte(`{"id":"item-1"}`), nil)
	_, replaceErr := container.ReplaceItem(context.Background(), pk, "item-1", []byte(`{"id":"item-1"}`), nil)
	_, upsertErr := container.UpsertItem(context.Background(), pk, "item-1", []byte(`{"id":"item-1"}`), nil)
	_, deleteErr := container.DeleteItem(context.Background(), pk, "item-1", nil)
	_, patchErr := container.PatchItem(context.Background(), pk, "item-1", validPatchOperations(t), nil)

	for _, err := range []error{readErr, createErr, replaceErr, upsertErr, deleteErr, patchErr} {
		var cosmosErr *Error
		require.True(t, errors.As(err, &cosmosErr))
		require.Equal(t, CodeClientError, cosmosErr.Code)
	}
}

// The point of factoring OperationOptions out is that every operation takes the same driver-level
// settings without restating them, so a knob added there reaches all of them at once. This pins
// that both item operations carry it, which a per-type copy would not guarantee.
func TestItemOptionsShareOperationOptions(t *testing.T) {
	shared := OperationOptions{
		ConsistencyStrategy: ReadConsistencyStrategySession,
		ExcludedRegions:     []Region{RegionEastUS},
		EndToEndTimeout:     5 * time.Second,
	}

	read := ReadItemOptions{Operation: shared}
	create := CreateItemOptions{Operation: shared}
	replace := ReplaceItemOptions{Operation: shared}
	upsert := UpsertItemOptions{Operation: shared}
	deleteOptions := DeleteItemOptions{Operation: shared}
	patch := PatchItemOptions{Operation: shared}

	require.Equal(t, shared, read.Operation)
	require.Equal(t, shared, create.Operation)
	require.Equal(t, shared, replace.Operation)
	require.Equal(t, shared, upsert.Operation)
	require.Equal(t, shared, deleteOptions.Operation)
	require.Equal(t, shared, patch.Operation)
	require.Equal(t, []Region{RegionEastUS}, read.Operation.ExcludedRegions,
		"excluded regions are typed, not free strings")
}

// Close guarantees that later operations fail with CodeClientClosed. Shutdown is exactly when a
// caller's context is also likely to be cancelled, so the closed client has to win: checking the
// context first would report context.Canceled and hide the real problem.
func TestClosedClientReportedAheadOfCancelledContext(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "item-1", nil)
	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "x", []byte(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.ReplaceItem(ctx, NewPartitionKeyString("pk"), "x", []byte(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.UpsertItem(ctx, NewPartitionKeyString("pk"), "x", []byte(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.DeleteItem(ctx, NewPartitionKeyString("pk"), "x", nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	_, err = container.PatchItem(ctx, NewPartitionKeyString("pk"), "x", validPatchOperations(t), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

func TestReadItemDoesNotNestClientLifetimeLocks(t *testing.T) {
	client, err := newClient("https://myaccount.documents.azure.com", testAccountKey, nil, nil)
	require.NoError(t, err)
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)

	beforeAcquire := make(chan struct{})
	continueAcquire := make(chan struct{})
	var once sync.Once
	client.beforeItemAcquire = func() {
		once.Do(func() {
			close(beforeAcquire)
			<-continueAcquire
		})
	}

	readResult := make(chan error, 1)
	go func() {
		_, err := container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item-1", nil)
		readResult <- err
	}()
	<-beforeAcquire

	closeResult := make(chan error, 1)
	go func() { closeResult <- client.Close() }()

	select {
	case err := <-closeResult:
		require.NoError(t, err,
			"Close must not wait on a lifetime lock held before executeItem acquires it")
	case <-time.After(time.Second):
		close(continueAcquire)
		t.Fatal("Close blocked, indicating ReadItem held a nested lifetime read lock")
	}
	close(continueAcquire)

	var cosmosErr *Error
	require.ErrorAs(t, <-readResult, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

// The ABI distinguishes "inherit whatever is configured" from "use the default for this account's
// consistency level". Collapsing both onto the zero value would make one of them unreachable, and
// silently so: the read would just use the wrong strategy.
func TestReadConsistencyStrategyUnsetIsNotDefault(t *testing.T) {
	require.NotEqual(t, ReadConsistencyStrategyUnset, ReadConsistencyStrategyDefault)

	var zero ReadConsistencyStrategy
	require.Equal(t, ReadConsistencyStrategyUnset, zero, "the zero value must mean inherit")
}

func TestReadItemRejectsUnknownConsistencyStrategy(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.ReadItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		&ReadItemOptions{
			Operation: OperationOptions{
				ConsistencyStrategy: ReadConsistencyStrategy("LatestCommited"),
			},
		},
	)

	require.ErrorContains(t, err, "unknown read consistency strategy")
	requireNotDriverUnavailable(t, err)
}

func TestPatchStrategyUnsetIsNotAuto(t *testing.T) {
	require.NotEqual(t, PatchStrategyUnset, PatchStrategyAuto)

	var zero PatchStrategy
	require.Equal(t, PatchStrategyUnset, zero, "the zero value must mean inherit")
}

func TestPatchItemRejectsUnknownStrategy(t *testing.T) {
	container := newTestContainer(t)

	response, err := container.PatchItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		validPatchOperations(t),
		&PatchItemOptions{Strategy: PatchStrategy("Automatic")},
	)

	require.Equal(t, ItemResponse{}, response)
	require.ErrorContains(t, err, "unknown patch strategy")
	requireNotDriverUnavailable(t, err)
}

func TestItemOperationsRejectNULSessionToken(t *testing.T) {
	container := newTestContainer(t)
	token := SessionToken("1:2\x00:3")

	_, err := container.ReadItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		&ReadItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")

	_, err = container.CreateItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		[]byte(`{"id":"item-1","pk":"pk"}`),
		&CreateItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")

	_, err = container.ReplaceItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		[]byte(`{"id":"item-1","pk":"pk"}`),
		&ReplaceItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")

	_, err = container.UpsertItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		[]byte(`{"id":"item-1","pk":"pk"}`),
		&UpsertItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")

	_, err = container.DeleteItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		&DeleteItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")

	_, err = container.PatchItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		validPatchOperations(t),
		&PatchItemOptions{SessionToken: token},
	)
	require.ErrorContains(t, err, "session token must not contain a NUL byte")
}

func TestItemWritesRejectInvalidIfMatchETag(t *testing.T) {
	container := newTestContainer(t)
	item := []byte(`{"id":"item-1","pk":"pk"}`)
	pk := NewPartitionKeyString("pk")

	for _, etag := range []azcore.ETag{"", "\"etag\x00suffix\""} {
		t.Run(string(etag), func(t *testing.T) {
			_, replaceErr := container.ReplaceItem(context.Background(), pk, "item-1", item,
				&ReplaceItemOptions{IfMatchETag: &etag})
			_, upsertErr := container.UpsertItem(context.Background(), pk, "item-1", item,
				&UpsertItemOptions{IfMatchETag: &etag})
			_, deleteErr := container.DeleteItem(context.Background(), pk, "item-1",
				&DeleteItemOptions{IfMatchETag: &etag})
			_, patchErr := container.PatchItem(context.Background(), pk, "item-1", validPatchOperations(t),
				&PatchItemOptions{IfMatchETag: &etag})

			for _, err := range []error{replaceErr, upsertErr, deleteErr, patchErr} {
				require.Error(t, err)
				requireNotDriverUnavailable(t, err)
			}
		})
	}
}

func TestPatchClientSidePreconditionErrorIsClassified(t *testing.T) {
	original := &Error{
		Code:       CodeClientError,
		StatusCode: 412,
		Message:    "client-side patch precondition failed",
	}
	normalized := normalizeItemOperationError(operationKindPatchItem, original)

	var cosmosErr *Error
	require.ErrorAs(t, normalized, &cosmosErr)
	require.Equal(t, CodePreconditionFailed, cosmosErr.Code)
	require.False(t, cosmosErr.FromWire)
	require.NotSame(t, original, cosmosErr)
	require.Equal(t, CodeClientError, original.Code, "normalization must not mutate the driver error")

	require.Same(t, original, normalizeItemOperationError(operationKindReplaceItem, original))
	original.FromWire = true
	require.Same(t, original, normalizeItemOperationError(operationKindPatchItem, original))
}

func TestReadItemRejectsNULETag(t *testing.T) {
	container := newTestContainer(t)
	etag := azcore.ETag("\"etag\x00suffix\"")

	_, err := container.ReadItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		&ReadItemOptions{IfNoneMatchETag: &etag},
	)

	require.ErrorContains(t, err, "IfNoneMatchETag must not contain a NUL byte")
}

func TestReadItemRejectsEmptyETag(t *testing.T) {
	container := newTestContainer(t)
	etag := azcore.ETag("")

	_, err := container.ReadItem(
		context.Background(),
		NewPartitionKeyString("pk"),
		"item-1",
		&ReadItemOptions{IfNoneMatchETag: &etag},
	)

	require.ErrorContains(t, err, "IfNoneMatchETag must not be empty")
}

// The driver's budget is what guarantees an operation terminates: cancelling the context stops it
// only once the driver notices, while without a budget it is bounded by transport timeouts times a
// retry budget. Passing the caller's deadline down is what makes one number bound every layer.
func TestEndToEndTimeoutFollowsTheContextDeadline(t *testing.T) {
	t.Run("no deadline leaves the driver default", func(t *testing.T) {
		require.Zero(t, endToEndTimeout(context.Background(), 0))
	})

	t.Run("deadline becomes the budget", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		got := endToEndTimeout(ctx, 0)
		require.Positive(t, got)
		require.LessOrEqual(t, got, time.Minute)
		require.Greater(t, got, 59*time.Second, "should be what remains, not a fixed value")
	})

	t.Run("an explicit setting wins over the deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// The caller is describing how long the operation may spend, which is a different thing
		// from when they stop waiting, so it is not second-guessed.
		require.Equal(t, 5*time.Second, endToEndTimeout(ctx, 5*time.Second))
	})

	t.Run("an expired deadline stays positive", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), -time.Second)
		defer cancel()

		// Zero would read as unset at the ABI, which would remove the bound rather than tighten it.
		require.Positive(t, endToEndTimeout(ctx, 0))
	})
}

func validPatchOperations(t *testing.T) PatchOperations {
	t.Helper()
	var operations PatchOperations
	require.NoError(t, operations.AppendSet("/value", 1))
	return operations
}
