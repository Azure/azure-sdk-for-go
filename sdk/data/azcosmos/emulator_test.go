// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

// These are the end-to-end tests: they run real operations against a real service through the
// driver, so they are the only place the binding is exercised as a whole rather than in pieces.
//
// They are gated on the EMULATOR environment variable, which is what the module's ci.yml already
// sets for its Cosmos Emulator stage, so they run there without further wiring and skip everywhere
// else. They also need a driver-backed build, hence the build tag.
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
	t.Helper()

	if os.Getenv("EMULATOR") == "" {
		t.Skip("set EMULATOR to run tests against the Cosmos DB emulator")
	}

	endpoint := os.Getenv("AZCOSMOS_ENDPOINT")
	if endpoint == "" {
		endpoint = emulatorEndpoint
	}
	databaseID := os.Getenv("AZCOSMOS_DATABASE")
	if databaseID == "" {
		databaseID = "itemdb"
	}
	containerID := os.Getenv("AZCOSMOS_CONTAINER")
	if containerID == "" {
		containerID = "items"
	}

	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	client, err := NewClientWithKey(endpoint, cred, options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	return client, databaseID, containerID
}

// emulatorContainer returns a container client for the emulator.
func emulatorContainer(t *testing.T) *ContainerClient {
	t.Helper()

	client, databaseID, containerID := emulatorClient(t)
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	return container
}

// uniqueItemID keeps tests from colliding with each other or with a previous run, since nothing
// here deletes what it creates: item deletion is not bound yet.
func uniqueItemID(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s-%d", t.Name(), time.Now().UnixNano())
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

// Cancelling the context has to stop the operation and report the context's own error, so that
// errors.Is against context.Canceled works the way callers expect of any Go API.
func TestEmulatorCancelledContext(t *testing.T) {
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
		enabled bool
	}{
		{"disabled", false},
		{"enabled", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := emulatorContainerWithOptions(t, &ClientOptions{EnableContentResponseOnWrite: tt.enabled})

			require.Equal(t, tt.enabled, createForClientOptions(t, container, nil),
				"an operation that sets nothing inherits the client's setting")
		})
	}
}

// An operation that sets the value overrides the client, which is what makes the client value a
// default rather than a policy.
func TestEmulatorOperationContentResponseOverridesTheClient(t *testing.T) {
	container := emulatorContainerWithOptions(t, &ClientOptions{EnableContentResponseOnWrite: true})

	require.False(t, createForClientOptions(t, container, &OperationOptions{
		EnableContentResponseOnWrite: to(false),
	}), "the operation asked for no content and the client asked for content")
}

// Preferred regions name the account's only region, so this proves the driver accepted the order
// rather than that it routed anywhere in particular; a single-region emulator cannot show more.
func TestEmulatorPreferredRegions(t *testing.T) {
	container := emulatorContainerWithOptions(t, &ClientOptions{Routing: PreferredRegions(RegionEastUS)})

	createForClientOptions(t, container, nil)
}

// An application ID reaches the driver through the runtime's user-agent suffix, which the driver
// validates and rejects with a bare status. The emulator does not report the user agent back, so
// this covers that a value the driver accepts survives the whole path; rejection is covered by
// TestNewClientWithKeyRejectsUnusableOptions, which needs no service.
func TestEmulatorApplicationID(t *testing.T) {
	// At the driver's 25-byte limit, so a regression in the limit shows up here too.
	container := emulatorContainerWithOptions(t, &ClientOptions{ApplicationID: "azcosmos-go-v2-e2e-testin"})

	createForClientOptions(t, container, nil)
}

// Callers that reach first use at the same time have to end up with the one driver the client
// owns. A second creation would be a wasted account-properties fetch at best, and a handle the
// client never frees at worst.
func TestEmulatorConcurrentFirstUseCreatesOneDriver(t *testing.T) {
	client, _, _ := emulatorClient(t)

	const callers = 16
	handles := make([]uintptr, callers)
	errs := make([]error, callers)

	// Released all at once so the callers actually contend, rather than the first one finishing
	// while the rest are still being started.
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			handles[i], errs[i] = client.driver.ensureDriverHandle(context.Background())
		}(i)
	}
	close(start)
	wg.Wait()

	for i := 0; i < callers; i++ {
		require.NoError(t, errs[i])
		require.NotZero(t, handles[i])
		require.Equal(t, handles[0], handles[i], "every caller has to get the same driver")
	}
}
