// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestAAD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "aadTest")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	aadClient := emulatorTests.getAadClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_item aContainer", "read_item aContainer", "replace_item aContainer", "upsert_item aContainer", "delete_item aContainer"},
	}))

	item := map[string]string{
		"id":    "1",
		"value": "2",
	}

	container, _ := aadClient.NewContainer("aadTest", "aContainer")
	pk := NewPartitionKeyString("1")

	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	itemResponse, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	if itemResponse.SessionToken == nil {
		t.Fatalf("Session token is empty")
	}

	// No content on write by default
	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}

	itemResponse, err = container.ReadItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	var itemResponseBody map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "2" {
		t.Fatalf("Expected value to be 2, got %v", itemResponseBody["value"])
	}

	item["value"] = "3"
	marshalled, err = json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	itemResponse, err = container.ReplaceItem(context.TODO(), pk, "1", marshalled, &ItemOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	// Explicitly requesting body on write
	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "3" {
		t.Fatalf("Expected value to be 3, got %v", itemResponseBody["value"])
	}

	item["value"] = "4"
	marshalled, err = json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	itemResponse, err = container.UpsertItem(context.TODO(), pk, marshalled, &ItemOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to upsert item: %v", err)
	}

	// Explicitly requesting body on write
	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "4" {
		t.Fatalf("Expected value to be 4, got %v", itemResponseBody["value"])
	}

	itemResponse, err = container.DeleteItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}
}

type emulatorFailAccountScopeCredential struct {
	accountScope string
	delegate     azcore.TokenCredential // emulatorTokenCredential
}

func (c *emulatorFailAccountScopeCredential) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if len(tro.Scopes) == 1 && tro.Scopes[0] == c.accountScope {
		// this exact string is what your client fallback looks for
		return azcore.AccessToken{}, fmt.Errorf("AADSTS500011: simulated resource not found")
	}
	return c.delegate.GetToken(ctx, tro)
}

func TestAAD_Fallback_E2E_WithEmulator(t *testing.T) {
	t.Setenv(envCosmosScopeOverride, "") // ensure product code uses endpoint scope + fallback

	em := newEmulatorTests(t)

	// 1) Set up DB/container using key auth
	keyClient := em.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	db := em.createDatabase(t, context.TODO(), keyClient, "aadFallbackE2E")
	defer em.deleteDatabase(t, context.TODO(), db)

	props := ContainerProperties{
		ID:                     "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/id"}},
	}
	if _, err := db.CreateContainer(context.TODO(), props, nil); err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	// 2) Build a cred that FAILS on account scope to trigger fallback
	u, err := url.Parse(em.host)
	if err != nil {
		t.Fatal(err)
	}
	accountScope := fmt.Sprintf("%s://%s/.default", u.Scheme, u.Hostname())

	failingThenDelegateCred := &emulatorFailAccountScopeCredential{
		accountScope: accountScope,
		delegate:     &emulatorTokenCredential{}, // emulator-success path
	}

	// 3) Create AAD client
	aadClient, err := NewClient(em.host, failingThenDelegateCred, &ClientOptions{
		ClientOptions: azcore.ClientOptions{},
	})
	if err != nil {
		t.Fatalf("Failed to create AAD client: %v", err)
	}

	// 4) Multiple CreateItem calls
	container, err := aadClient.NewContainer("aadFallbackE2E", "aContainer")
	if err != nil {
		t.Fatalf("NewContainer: %v", err)
	}

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("id-%d", i)
		item := map[string]string{
			"id":    id, // value at /id
			"value": fmt.Sprintf("v-%d", i),
		}
		body, _ := json.Marshal(item)

		pk := NewPartitionKeyString(id)
		if _, err := container.CreateItem(context.TODO(), pk, body, nil); err != nil {
			t.Fatalf("create item %d failed (fallback should succeed): %v", i, err)
		}
	}

	readPK := NewPartitionKeyString("id-0")
	ri, err := container.ReadItem(context.TODO(), readPK, "id-0", nil)
	if err != nil || len(ri.Value) == 0 {
		t.Fatalf("read item failed: %v", err)
	}
}
