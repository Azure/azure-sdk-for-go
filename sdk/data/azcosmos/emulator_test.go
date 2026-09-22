// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore azsdk itemdb testin upserted

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// These are the end-to-end tests: they run real operations against a real service through the
// driver, so they are the only place the binding is exercised as a whole rather than in pieces.
//
// They are gated on the EMULATOR environment variable, which is what the module's ci.yml already
// sets for its Cosmos Emulator stage, so they run there without further wiring and skip everywhere
// else.
//
// A container has to exist before they run. Creating one needs container management operations,
// which are not bound yet, so the database and container ids are taken from the environment and
// default to what the emulator ships with.

const (
	// emulatorEndpoint is where the emulator listens by default.
	emulatorEndpoint = "https://localhost:8081/"
)

// emulatorClient returns a client for the emulator, or skips the test when one is not configured.
func emulatorClient(t *testing.T) (*Client, string, string) {
	t.Helper()
	return emulatorClientWithOptions(t, nil)
}

// emulatorClientWithOptions is emulatorClient with the client options under test.
func emulatorClientWithOptions(t *testing.T, options *ClientOptions) (*Client, string, string) {
	return emulatorClientConfigured(t, options, true)
}

// emulatorClientConfigured optionally initializes the client before returning it.
func emulatorClientConfigured(
	t *testing.T,
	options *ClientOptions,
	initialize bool,
) (*Client, string, string) {
	t.Helper()

	endpoint, databaseID, containerID := emulatorConfiguration(t)

	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	client, err := NewClientWithKey(endpoint, cred, options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	if initialize {
		require.NoError(t, client.Initialize(t.Context()))
	}

	return client, databaseID, containerID
}

func emulatorConfiguration(t *testing.T) (endpoint, databaseID, containerID string) {
	t.Helper()

	if os.Getenv("EMULATOR") == "" {
		t.Skip("set EMULATOR to run tests against the Cosmos DB emulator")
	}

	endpoint = os.Getenv("AZCOSMOS_ENDPOINT")
	if endpoint == "" {
		endpoint = emulatorEndpoint
	}
	databaseID = os.Getenv("AZCOSMOS_DATABASE")
	if databaseID == "" {
		databaseID = "itemdb"
	}
	containerID = os.Getenv("AZCOSMOS_CONTAINER")
	if containerID == "" {
		containerID = "items"
	}
	return
}

func setEmulatorReplication(t *testing.T, region Region, action string) {
	t.Helper()

	managementEndpoint := os.Getenv("AZCOSMOS_MANAGEMENT_ENDPOINT")
	if managementEndpoint == "" {
		t.Skip("set AZCOSMOS_MANAGEMENT_ENDPOINT to test multi-region routing")
	}
	endpoint, err := url.JoinPath(managementEndpoint, "regions", string(region), "replication", action)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusOK, response.StatusCode, string(body))
}

type recordingTokenCredential struct {
	requests chan []string
}

func (c *recordingTokenCredential) GetToken(
	_ context.Context,
	options policy.TokenRequestOptions,
) (azcore.AccessToken, error) {
	c.requests <- append([]string(nil), options.Scopes...)
	return azcore.AccessToken{
		Token:     "emulator-access-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

type recoveringTokenCredential struct {
	attempts atomic.Int32
}

func (c *recoveringTokenCredential) GetToken(
	_ context.Context,
	_ policy.TokenRequestOptions,
) (azcore.AccessToken, error) {
	if c.attempts.Add(1) <= 2 {
		return azcore.AccessToken{}, errors.New("temporary token acquisition failure")
	}
	return azcore.AccessToken{
		Token:     "recovered-access-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

type delayedTokenCredential struct {
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (c *delayedTokenCredential) GetToken(
	ctx context.Context,
	_ policy.TokenRequestOptions,
) (azcore.AccessToken, error) {
	c.once.Do(func() { close(c.started) })
	select {
	case <-c.release:
	case <-ctx.Done():
		return azcore.AccessToken{}, ctx.Err()
	}
	return azcore.AccessToken{
		Token:     "delayed-access-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

// emulatorContainer returns a container client for the emulator.
func emulatorContainer(t *testing.T) *ContainerClient {
	t.Helper()

	client, databaseID, containerID := emulatorClient(t)
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	return container
}

// uniqueItemID keeps tests from colliding with each other or with a previous run.
func uniqueItemID(t *testing.T) string {
	t.Helper()
	name := strings.NewReplacer("/", "-", "\\", "-", "?", "-", "#", "-").Replace(t.Name())
	return fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
}

func TestEmulatorUniqueItemIDIsResourceSafe(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {
		id := uniqueItemID(t)
		for _, invalid := range []string{"/", "\\", "?", "#"} {
			require.NotContains(t, id, invalid)
		}
	})
}

func trackEmulatorItem(t *testing.T, container *ContainerClient, pk PartitionKey, id string) {
	t.Helper()
	t.Cleanup(func() {
		_, err := container.DeleteItem(context.Background(), pk, id, nil)
		if err == nil {
			return
		}
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeNotFound, cosmosErr.Code)
	})
}

// The round trip is what proves the binding works end to end: the request reaches the service and
// the response comes back through the completion queue with its headers intact.
func TestEmulatorCreateThenReadItem(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": 42})
	require.NoError(t, err)

	created, err := container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)
	// Cosmos DB bills every request, so a charge of zero means the header never made it back.
	require.Positive(t, created.RequestCharge, "the request charge has to survive the binding")
	require.NotEmpty(t, created.ActivityID, "the activity id has to survive the binding")
	require.NotEmpty(t, created.ETag, "a create returns the new item's ETag")

	read, err := container.ReadItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Positive(t, read.RequestCharge)
	require.Equal(t, created.ETag, read.ETag, "the item has not changed since it was created")

	// A read always returns the item, whatever the content-response setting says about writes.
	var roundTripped map[string]any
	require.NoError(t, json.Unmarshal(read.Value, &roundTripped))
	require.Equal(t, id, roundTripped["id"])
	require.InDelta(t, 42.0, roundTripped["value"], 0)
}

// A missing item has to classify as CodeNotFound rather than surfacing as a bare status, which is
// the whole point of the Code type.
func TestEmulatorReadMissingItem(t *testing.T) {
	container := emulatorContainer(t)

	id := uniqueItemID(t)
	_, err := container.ReadItem(context.Background(), NewPartitionKeyString(id), id, nil)
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeNotFound, cosmosErr.Code)
	require.Equal(t, 404, cosmosErr.StatusCode)
	require.True(t, cosmosErr.FromWire, "the service produced this, not the client")
	require.NotEmpty(t, cosmosErr.ActivityID, "a wire failure still correlates with telemetry")
}

// Creating the same id twice conflicts. This also pins that a failed request still reports what it
// cost, which is the reason Error carries a request charge at all.
func TestEmulatorCreateItemConflict(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeConflict, cosmosErr.Code)
	require.Equal(t, 409, cosmosErr.StatusCode)
	require.Positive(t, cosmosErr.RequestCharge, "a rejected write is still billed")
}

// The content-response setting is per operation, and off is the default, so a create returns no
// item unless it is asked to.
func TestEmulatorCreateItemContentResponse(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name        string
		enabled     *bool
		wantContent bool
	}{
		{"default", nil, false},
		{"disabled", to(false), false},
		{"enabled", to(true), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			id := uniqueItemID(t)
			item, err := json.Marshal(map[string]any{"id": id, "pk": id})
			require.NoError(t, err)

			response, err := container.CreateItem(ctx, NewPartitionKeyString(id), id, item,
				&CreateItemOptions{Operation: OperationOptions{EnableContentResponseOnWrite: tt.enabled}})
			require.NoError(t, err)

			if tt.wantContent {
				require.NotEmpty(t, response.Value, "the item was requested back")
				return
			}
			require.Empty(t, response.Value, "no item was requested back")
		})
	}
}

// A conditional read of an unchanged item comes back with no payload, which is what makes it
// cheaper than an unconditional one. An ETag that cannot match still returns the item, which is
// what pins that the precondition is being evaluated rather than the read simply returning nothing.
func TestEmulatorReadItemIfNoneMatch(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	created, err := container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)
	require.NotEmpty(t, created.ETag)

	unchanged := created.ETag
	response, err := container.ReadItem(ctx, pk, id, &ReadItemOptions{IfNoneMatchETag: &unchanged})
	require.NoError(t, err, "an unchanged item is not an error")
	require.Empty(t, response.Value, "the item has not changed, so it is not sent")

	stale := azcore.ETag("\"00000000-0000-0000-0000-000000000000\"")
	response, err = container.ReadItem(ctx, pk, id, &ReadItemOptions{IfNoneMatchETag: &stale})
	require.NoError(t, err)
	require.NotEmpty(t, response.Value, "the ETag does not match, so the item is sent")
}

func TestEmulatorReplaceAndUpsertItem(t *testing.T) {
	container := emulatorContainer(t)
	ctx := t.Context()

	t.Run("replace", func(t *testing.T) {
		id := uniqueItemID(t)
		pk := NewPartitionKeyString(id)
		trackEmulatorItem(t, container, pk, id)

		original, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": "original"})
		require.NoError(t, err)
		created, err := container.CreateItem(ctx, pk, id, original, nil)
		require.NoError(t, err)

		replacement, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": "replaced"})
		require.NoError(t, err)
		replaced, err := container.ReplaceItem(ctx, pk, id, replacement, &ReplaceItemOptions{
			Operation: OperationOptions{EnableContentResponseOnWrite: to(true)},
		})
		require.NoError(t, err)
		require.Positive(t, replaced.RequestCharge)
		require.NotEmpty(t, replaced.ActivityID)
		require.NotEmpty(t, replaced.SessionToken)
		require.NotEmpty(t, replaced.ETag)
		require.NotEqual(t, created.ETag, replaced.ETag)

		var echoed map[string]any
		require.NoError(t, json.Unmarshal(replaced.Value, &echoed))
		require.Equal(t, "replaced", echoed["value"])

		read, err := container.ReadItem(ctx, pk, id, &ReadItemOptions{
			SessionToken: replaced.SessionToken,
			Operation:    OperationOptions{ConsistencyStrategy: ReadConsistencyStrategySession},
		})
		require.NoError(t, err)
		require.Equal(t, replaced.ETag, read.ETag)
		requireItemFields(t, read.Value, map[string]any{
			"id": id, "pk": id, "value": "replaced",
		})
	})

	t.Run("upsert existing and missing", func(t *testing.T) {
		existingID := uniqueItemID(t)
		existingPK := NewPartitionKeyString(existingID)
		trackEmulatorItem(t, container, existingPK, existingID)

		original, err := json.Marshal(map[string]any{"id": existingID, "pk": existingID, "version": 1})
		require.NoError(t, err)
		_, err = container.CreateItem(ctx, existingPK, existingID, original, nil)
		require.NoError(t, err)

		updated, err := json.Marshal(map[string]any{"id": existingID, "pk": existingID, "version": 2})
		require.NoError(t, err)
		upserted, err := container.UpsertItem(ctx, existingPK, existingID, updated, &UpsertItemOptions{
			Operation: OperationOptions{EnableContentResponseOnWrite: to(true)},
		})
		require.NoError(t, err)
		require.Positive(t, upserted.RequestCharge)
		require.NotEmpty(t, upserted.ActivityID)
		require.NotEmpty(t, upserted.ETag)
		requireItemFields(t, upserted.Value, map[string]any{
			"id": existingID, "pk": existingID, "version": float64(2),
		})

		missingID := uniqueItemID(t)
		missingPK := NewPartitionKeyString(missingID)
		trackEmulatorItem(t, container, missingPK, missingID)
		missing, err := json.Marshal(map[string]any{"id": missingID, "pk": missingID, "version": 1})
		require.NoError(t, err)
		created, err := container.UpsertItem(ctx, missingPK, missingID, missing, nil)
		require.NoError(t, err)
		require.NotEmpty(t, created.ETag)

		read, err := container.ReadItem(ctx, missingPK, missingID, nil)
		require.NoError(t, err)
		requireItemFields(t, read.Value, map[string]any{
			"id": missingID, "pk": missingID, "version": float64(1),
		})
	})
}

func TestEmulatorConditionalItemWrites(t *testing.T) {
	container := emulatorContainer(t)
	ctx := t.Context()
	stale := azcore.ETag("\"00000000-0000-0000-0000-000000000000\"")

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": 1})
	require.NoError(t, err)
	created, err := container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	replacement, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": 2})
	require.NoError(t, err)
	var patch PatchOperations
	require.NoError(t, patch.AppendSet("/value", 2))

	_, err = container.ReplaceItem(ctx, pk, id, replacement, &ReplaceItemOptions{IfMatchETag: &stale})
	requireWireError(t, err, CodePreconditionFailed, 412)
	_, err = container.UpsertItem(ctx, pk, id, replacement, &UpsertItemOptions{IfMatchETag: &stale})
	requireWireError(t, err, CodePreconditionFailed, 412)
	_, err = container.PatchItem(ctx, pk, id, patch, &PatchItemOptions{IfMatchETag: &stale})
	requirePatchPreconditionError(t, err)
	_, err = container.DeleteItem(ctx, pk, id, &DeleteItemOptions{IfMatchETag: &stale})
	requireWireError(t, err, CodePreconditionFailed, 412)

	replaced, err := container.ReplaceItem(ctx, pk, id, replacement, &ReplaceItemOptions{
		IfMatchETag: &created.ETag,
	})
	require.NoError(t, err)
	require.NotEmpty(t, replaced.ETag)

	current := replaced.ETag
	patched, err := container.PatchItem(ctx, pk, id, patch, &PatchItemOptions{IfMatchETag: &current})
	require.NoError(t, err)
	require.NotEmpty(t, patched.ETag)

	missingID := uniqueItemID(t)
	missingPK := NewPartitionKeyString(missingID)
	trackEmulatorItem(t, container, missingPK, missingID)
	missing, err := json.Marshal(map[string]any{"id": missingID, "pk": missingID})
	require.NoError(t, err)
	_, err = container.UpsertItem(ctx, missingPK, missingID, missing, &UpsertItemOptions{
		IfMatchETag: &stale,
	})
	require.NoError(t, err, "a missing item is created even with a stale If-Match")

	read, err := container.ReadItem(ctx, pk, id, nil)
	require.NoError(t, err)
	deleted, err := container.DeleteItem(ctx, pk, id, &DeleteItemOptions{IfMatchETag: &read.ETag})
	require.NoError(t, err)
	require.Nil(t, deleted.Value)
}

func TestEmulatorDeleteItem(t *testing.T) {
	container := emulatorContainer(t)
	ctx := t.Context()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	deleted, err := container.DeleteItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Positive(t, deleted.RequestCharge)
	require.NotEmpty(t, deleted.ActivityID)
	require.Nil(t, deleted.Value)

	_, err = container.ReadItem(ctx, pk, id, nil)
	requireWireStatus(t, err, CodeNotFound, 404)
	_, err = container.DeleteItem(ctx, pk, id, nil)
	requireWireStatus(t, err, CodeNotFound, 404)
}

func TestEmulatorPatchItem(t *testing.T) {
	container := emulatorContainer(t)
	ctx := t.Context()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	item, err := json.Marshal(map[string]any{
		"id":      id,
		"pk":      id,
		"count":   1,
		"replace": "before",
		"remove":  true,
		"source":  "move-me",
	})
	require.NoError(t, err)
	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	var operations PatchOperations
	require.NoError(t, operations.AppendAdd("/added", "new"))
	require.NoError(t, operations.AppendSet("/set", map[string]any{"nested": true}))
	require.NoError(t, operations.AppendReplace("/replace", "after"))
	require.NoError(t, operations.AppendRemove("/remove"))
	require.NoError(t, operations.AppendIncrement("/count", 2))
	require.NoError(t, operations.AppendMove("/source", "/moved"))

	patched, err := container.PatchItem(ctx, pk, id, operations, nil)
	require.NoError(t, err)
	require.Positive(t, patched.RequestCharge)
	require.NotEmpty(t, patched.ActivityID)
	require.NotEmpty(t, patched.SessionToken)
	require.NotEmpty(t, patched.ETag)
	require.NotEmpty(t, patched.Value, "Patch defaults to returning the post-image")

	var value map[string]any
	require.NoError(t, json.Unmarshal(patched.Value, &value))
	require.Equal(t, "new", value["added"])
	require.Equal(t, true, value["set"].(map[string]any)["nested"])
	require.Equal(t, "after", value["replace"])
	require.NotContains(t, value, "remove")
	require.InDelta(t, 3, value["count"], 0)
	require.Equal(t, "move-me", value["moved"])
	require.NotContains(t, value, "source")
	require.Contains(t, value, "_azsdkPatchTracking",
		"unsafe automatic PATCH records retry deduplication state on the item")
}

func TestEmulatorPatchContentResponse(t *testing.T) {
	for _, tt := range []struct {
		name          string
		clientEnabled *bool
		callEnabled   *bool
		wantContent   bool
	}{
		{"native default", nil, nil, true},
		{"operation disabled", nil, to(false), false},
		{"client disabled", to(false), nil, false},
		{"operation overrides client", to(false), to(true), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := emulatorContainerWithOptions(t, &ClientOptions{
				EnableContentResponseOnWrite: tt.clientEnabled,
			})
			id := uniqueItemID(t)
			pk := NewPartitionKeyString(id)
			trackEmulatorItem(t, container, pk, id)
			item, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": 1})
			require.NoError(t, err)
			_, err = container.CreateItem(t.Context(), pk, id, item, nil)
			require.NoError(t, err)

			var operations PatchOperations
			require.NoError(t, operations.AppendSet("/value", 2))
			response, err := container.PatchItem(t.Context(), pk, id, operations, &PatchItemOptions{
				Operation: OperationOptions{EnableContentResponseOnWrite: tt.callEnabled},
			})
			require.NoError(t, err)
			require.Equal(t, tt.wantContent, len(response.Value) > 0)
		})
	}
}

func TestEmulatorPatchMoreThanTenOperations(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), pk, id, item, nil)
	require.NoError(t, err)

	var operations PatchOperations
	for i := range 11 {
		require.NoError(t, operations.AppendSet(fmt.Sprintf("/value%d", i), i))
	}
	response, err := container.PatchItem(t.Context(), pk, id, operations, nil)
	require.NoError(t, err)

	var value map[string]any
	require.NoError(t, json.Unmarshal(response.Value, &value))
	for i := range 11 {
		require.InDelta(t, i, value[fmt.Sprintf("value%d", i)], 0)
	}
}

func TestEmulatorNewItemOperationsReportMissingItems(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item := []byte(fmt.Sprintf(`{"id":%q,"pk":%q}`, id, id))
	var operations PatchOperations
	require.NoError(t, operations.AppendSet("/value", 1))

	_, err := container.ReplaceItem(t.Context(), pk, id, item, nil)
	requireWireStatus(t, err, CodeNotFound, 404)
	_, err = container.DeleteItem(t.Context(), pk, id, nil)
	requireWireStatus(t, err, CodeNotFound, 404)
	_, err = container.PatchItem(t.Context(), pk, id, operations, nil)
	requireWireStatus(t, err, CodeNotFound, 404)
}

func TestEmulatorConcurrentItemWriteOperations(t *testing.T) {
	container := emulatorContainer(t)
	const operations = 8

	var wg sync.WaitGroup
	errs := make(chan error, operations)
	wg.Add(operations)
	for i := range operations {
		go func() {
			defer wg.Done()
			id := fmt.Sprintf("%s-%d-%d", t.Name(), time.Now().UnixNano(), i)
			pk := NewPartitionKeyString(id)
			item := []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"value":0}`, id, id))

			if _, err := container.UpsertItem(t.Context(), pk, id, item, nil); err != nil {
				errs <- err
				return
			}
			replacement := []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"value":1}`, id, id))
			if _, err := container.ReplaceItem(t.Context(), pk, id, replacement, nil); err != nil {
				errs <- err
				return
			}
			var patch PatchOperations
			if err := patch.AppendIncrement("/value", 1); err != nil {
				errs <- err
				return
			}
			if _, err := container.PatchItem(t.Context(), pk, id, patch, nil); err != nil {
				errs <- err
				return
			}
			read, err := container.ReadItem(t.Context(), pk, id, nil)
			if err != nil {
				errs <- err
				return
			}
			var value map[string]any
			if err := json.Unmarshal(read.Value, &value); err != nil {
				errs <- err
				return
			}
			if value["id"] != id || value["value"] != float64(2) {
				errs <- fmt.Errorf("item correlation failed for %s: %v", id, value)
				return
			}
			if _, err := container.DeleteItem(t.Context(), pk, id, nil); err != nil {
				errs <- err
				return
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}

func requireWireError(t *testing.T, err error, code Code, status int) {
	t.Helper()
	cosmosErr := requireWireErrorDetails(t, err, code, status)
	require.Positive(t, cosmosErr.RequestCharge)
}

func requireWireStatus(t *testing.T, err error, code Code, status int) {
	t.Helper()
	require.NotNil(t, requireWireErrorDetails(t, err, code, status))
}

func requireWireErrorDetails(t *testing.T, err error, code Code, status int) *Error {
	t.Helper()
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, code, cosmosErr.Code)
	require.Equal(t, status, cosmosErr.StatusCode)
	require.True(t, cosmosErr.FromWire)
	require.NotEmpty(t, cosmosErr.ActivityID)
	return cosmosErr
}

func requireItemFields(t *testing.T, body []byte, expected map[string]any) {
	t.Helper()

	var item map[string]any
	require.NoError(t, json.Unmarshal(body, &item))
	for name, value := range expected {
		require.Equal(t, value, item[name], "item field %q", name)
	}
}

func requirePatchPreconditionError(t *testing.T, err error) {
	t.Helper()
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodePreconditionFailed, cosmosErr.Code)
	require.Equal(t, 412, cosmosErr.StatusCode)
	if cosmosErr.FromWire {
		require.Positive(t, cosmosErr.RequestCharge)
		require.NotEmpty(t, cosmosErr.ActivityID)
		return
	}
	require.Zero(t, cosmosErr.RequestCharge)
	require.Empty(t, cosmosErr.ActivityID)
}

// The session token a write produces has to be usable on a later read, which is how a caller reads
// their own writes under session consistency.
func TestEmulatorSessionTokenRoundTrips(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	created, err := container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)
	require.NotEmpty(t, created.SessionToken, "a write reports the session it advanced")

	read, err := container.ReadItem(ctx, pk, id, &ReadItemOptions{
		SessionToken: created.SessionToken,
		Operation:    OperationOptions{ConsistencyStrategy: ReadConsistencyStrategySession},
	})
	require.NoError(t, err)
	require.NotEmpty(t, read.Value)
}

// LatestCommitted is a quorum read independent of the account's default consistency level. A live
// round trip proves the ABI discriminant reaches the driver and service rather than only mapping
// correctly in the converter.
func TestEmulatorLatestCommittedRead(t *testing.T) {
	container := emulatorContainer(t)
	ctx := t.Context()

	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	_, err = container.CreateItem(ctx, pk, id, item, nil)
	require.NoError(t, err)

	read, err := container.ReadItem(ctx, pk, id, &ReadItemOptions{
		Operation: OperationOptions{
			ConsistencyStrategy: ReadConsistencyStrategyLatestCommitted,
		},
	})
	require.NoError(t, err)

	var value map[string]any
	require.NoError(t, json.Unmarshal(read.Value, &value))
	require.Equal(t, id, value["id"])
	require.Equal(t, id, value["pk"])
}

// A context cancelled before an operation is submitted fails through the public fast path.
func TestEmulatorPreCancelledContext(t *testing.T) {
	container := emulatorContainer(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	id := uniqueItemID(t)
	_, err := container.ReadItem(ctx, NewPartitionKeyString(id), id, nil)
	require.ErrorIs(t, err, context.Canceled)
}

// Operations run concurrently against one client, which is the case the completion queue exists to
// support: many in flight at once, each answered with its own result.
func TestEmulatorConcurrentOperations(t *testing.T) {
	container := emulatorContainer(t)
	ctx := context.Background()

	const operations = 16
	errs := make(chan error, operations)
	ids := make(chan string, operations)

	for i := range operations {
		go func() {
			id := fmt.Sprintf("%s-%d-%d", t.Name(), time.Now().UnixNano(), i)
			item, err := json.Marshal(map[string]any{"id": id, "pk": id})
			if err != nil {
				errs <- err
				return
			}
			if _, err := container.CreateItem(ctx, NewPartitionKeyString(id), id, item, nil); err != nil {
				errs <- err
				return
			}
			ids <- id
			errs <- nil
		}()
	}

	for range operations {
		require.NoError(t, <-errs)
	}

	// Every write has to be readable, which catches a completion delivered to the wrong operation.
	close(ids)
	for id := range ids {
		response, err := container.ReadItem(ctx, NewPartitionKeyString(id), id, nil)
		require.NoError(t, err, "item %s", id)

		var item map[string]any
		require.NoError(t, json.Unmarshal(response.Value, &item))
		require.Equal(t, id, item["id"], "a completion was delivered to the wrong operation")
	}
}

// A container that does not exist fails when it is resolved, before any item work is attempted.
func TestEmulatorUnknownContainer(t *testing.T) {
	client, databaseID, _ := emulatorClient(t)

	container, err := client.NewContainer(databaseID, "no-such-container")
	require.NoError(t, err, "naming a container does not contact the service")

	_, err = container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "item-1", nil)
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeNotFound, cosmosErr.Code)
}

// to returns a pointer to v, for the tri-state option fields.
func to[T any](v T) *T {
	return &v
}

// emulatorContainerWithOptions returns a container client built with the options under test.
func emulatorContainerWithOptions(t *testing.T, options *ClientOptions) *ContainerClient {
	t.Helper()

	client, databaseID, containerID := emulatorClientWithOptions(t, options)
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	return container
}

// createForClientOptions creates an item and reports whether the service sent it back, which is
// what the content-response setting controls.
func createForClientOptions(t *testing.T, container *ContainerClient, operation *OperationOptions) bool {
	t.Helper()

	id := uniqueItemID(t)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	createOptions := &CreateItemOptions{}
	if operation != nil {
		createOptions.Operation = *operation
	}
	response, err := container.CreateItem(context.Background(), NewPartitionKeyString(id), id, item, createOptions)
	require.NoError(t, err)
	return len(response.Value) > 0
}

// The client-level content-response setting is only useful if an operation that does not set its
// own inherits it, so that is what this asserts.
func TestEmulatorClientContentResponseOnWrite(t *testing.T) {
	for _, tt := range []struct {
		name    string
		enabled *bool
		want    bool
	}{
		{"unset", nil, false},
		{"disabled", to(false), false},
		{"enabled", to(true), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := emulatorContainerWithOptions(t, &ClientOptions{EnableContentResponseOnWrite: tt.enabled})

			require.Equal(t, tt.want, createForClientOptions(t, container, nil),
				"an operation that sets nothing inherits the client's setting")
		})
	}
}

// An operation that sets the value overrides the client, which is what makes the client value a
// default rather than a policy.
func TestEmulatorOperationContentResponseOverridesTheClient(t *testing.T) {
	container := emulatorContainerWithOptions(t, &ClientOptions{EnableContentResponseOnWrite: to(true)})

	require.False(t, createForClientOptions(t, container, &OperationOptions{
		EnableContentResponseOnWrite: to(false),
	}), "the operation asked for no content and the client asked for content")
}

// Pausing replication makes a new item visible in East US but absent from West US. Eventual reads
// can then prove which region each routing strategy selected, rather than only proving that the
// driver accepted its preferred-region list.
func TestEmulatorRoutingStrategiesRouteReads(t *testing.T) {
	const westUSManagementName Region = "West US"
	setEmulatorReplication(t, westUSManagementName, "pause")
	t.Cleanup(func() { setEmulatorReplication(t, westUSManagementName, "resume") })

	writer := emulatorContainerWithOptions(t, &ClientOptions{
		Routing: PreferredRegions(RegionEastUS),
	})
	id := uniqueItemID(t)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = writer.CreateItem(t.Context(), NewPartitionKeyString(id), id, item, nil)
	require.NoError(t, err)

	for _, tt := range []struct {
		name       string
		routing    RoutingStrategy
		wantExists bool
	}{
		{"preferred East US", PreferredRegions(RegionEastUS), true},
		{"preferred West US", PreferredRegions(RegionWestUS), false},
		{"proximity to West US", ProximityTo(RegionWestUS), false},
		{"unknown proximity", ProximityTo("not-a-real-region"), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := emulatorContainerWithOptions(t, &ClientOptions{Routing: tt.routing})
			response, err := container.ReadItem(t.Context(), NewPartitionKeyString(id), id, &ReadItemOptions{
				Operation: OperationOptions{ConsistencyStrategy: ReadConsistencyStrategyEventual},
			})
			if tt.wantExists {
				require.NoError(t, err)
				require.NotEmpty(t, response.Value)
				return
			}

			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeNotFound, cosmosErr.Code,
				"the item exists in East US, so not found proves the read reached stale West US")
		})
	}
}

// An application ID reaches the driver through the runtime's user-agent suffix, which the driver
// validates. The emulator does not report the user agent back, so this covers that a value the
// driver accepts survives initialization and a request.
func TestEmulatorApplicationID(t *testing.T) {
	// At the driver's 25-byte limit, so a regression in the limit shows up here too.
	container := emulatorContainerWithOptions(t, &ClientOptions{ApplicationID: "azcosmos-go-v2-e2e-testin"})

	createForClientOptions(t, container, nil)
}

// Initialize does not return until it has fetched account properties and seeded routing state. A
// later operation therefore sees the same cached driver handle.
func TestEmulatorInitializeCreatesDriver(t *testing.T) {
	client, _, _ := emulatorClient(t)

	client.driver.mu.Lock()
	created := client.driver.created
	driver := client.driver.driver
	client.driver.mu.Unlock()

	require.True(t, created)
	require.NotNil(t, driver)
	again, err := client.driver.ensureDriver(t.Context())
	require.NoError(t, err)
	require.Equal(t, driver, again)
}

// Token acquisition crosses from Rust into Go and completes asynchronously back into Rust. A full
// operation proves the callback, token lifetime and account reference all survive that round trip.
func TestEmulatorTokenCredential(t *testing.T) {
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	credential := &recordingTokenCredential{requests: make(chan []string, 1)}

	client, err := NewClient(endpoint, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.NoError(t, client.Initialize(t.Context()))

	require.Equal(t, []string{"https://cosmos.azure.com/.default"}, <-credential.requests)

	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	id := uniqueItemID(t)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	_, err = container.CreateItem(t.Context(), NewPartitionKeyString(id), id, item, nil)
	require.NoError(t, err)
	read, err := container.ReadItem(t.Context(), NewPartitionKeyString(id), id, nil)
	require.NoError(t, err)
	require.NotEmpty(t, read.Value)
}

// Rust retries token acquisition twice within one initialization. Once those transient failures
// are exhausted, a later Initialize must start a fresh driver-creation attempt.
func TestEmulatorInitializationRecoversFromTokenFailure(t *testing.T) {
	endpoint, _, _ := emulatorConfiguration(t)
	credential := &recoveringTokenCredential{}
	client, err := NewClient(endpoint, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	err = client.Initialize(t.Context())
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeAuthenticationFailed, cosmosErr.Code)
	require.Equal(t, int32(2), credential.attempts.Load())

	require.NoError(t, client.Initialize(t.Context()))
	require.Equal(t, int32(3), credential.attempts.Load())
}

// A follower with a live context must not inherit the leader's cancellation. It retries after the
// shared attempt ends and can still initialize the client.
func TestEmulatorInitializationFollowerOutlivesLeader(t *testing.T) {
	endpoint, _, _ := emulatorConfiguration(t)
	credential := &delayedTokenCredential{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	client, err := NewClient(endpoint, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	waiting := make(chan struct{}, 1)
	client.driver.beforeCreationWait = func() { waiting <- struct{}{} }

	leaderCtx, cancelLeader := context.WithCancel(t.Context())
	leaderResult := make(chan error, 1)
	go func() { leaderResult <- client.Initialize(leaderCtx) }()
	<-credential.started

	followerResult := make(chan error, 1)
	go func() { followerResult <- client.Initialize(t.Context()) }()
	<-waiting

	cancelLeader()
	require.ErrorIs(t, <-leaderResult, context.Canceled)
	close(credential.release)
	require.NoError(t, <-followerResult)
}

// Initialize is optional: the first operation performs the same work when callers skip it.
func TestEmulatorOperationInitializesLazily(t *testing.T) {
	client, databaseID, containerID := emulatorClientConfigured(t, nil, false)

	client.driver.mu.Lock()
	require.False(t, client.driver.created)
	client.driver.mu.Unlock()

	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	id := uniqueItemID(t)
	item, err := json.Marshal(map[string]any{"id": id, "pk": id})
	require.NoError(t, err)

	_, err = container.CreateItem(t.Context(), NewPartitionKeyString(id), id, item, nil)
	require.NoError(t, err)

	client.driver.mu.Lock()
	require.True(t, client.driver.created)
	require.NotNil(t, client.driver.driver)
	client.driver.mu.Unlock()
}
