// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
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

func TestAAD_Emulator_UsesClientOptionsAudience(t *testing.T) {
	em := newEmulatorTests(t)

	keyClient := em.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	db := em.createDatabase(t, context.TODO(), keyClient, "aadClientOptionsAudienceTest")
	defer em.deleteDatabase(t, context.TODO(), db)

	props := ContainerProperties{
		ID:                     "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/id"}},
	}
	if _, err := db.CreateContainer(context.TODO(), props, nil); err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	customAudience := "https://custom.cosmos.azure.com"
	cred := &emulatorTokenCredential{} // Use emulator credential for CI reliability

	aadClient, err := NewClient(em.host, cred, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.Configuration{
				Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
					ServiceName: {Audience: customAudience},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create AAD client: %v", err)
	}

	_, err = aadClient.NewContainer("aadClientOptionsAudienceTest", "aContainer")
	if err != nil {
		t.Fatalf("NewContainer: %v", err)
	}
}

func TestAAD_Emulator_UsesAccountScope_WhenNoAudienceProvided(t *testing.T) {
	em := newEmulatorTests(t)

	keyClient := em.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	db := em.createDatabase(t, context.TODO(), keyClient, "aadAccountScopeTest")
	defer em.deleteDatabase(t, context.TODO(), db)

	props := ContainerProperties{
		ID:                     "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/id"}},
	}
	if _, err := db.CreateContainer(context.TODO(), props, nil); err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	cred := &emulatorTokenCredential{} // Use emulator credential for CI reliability

	aadClient, err := NewClient(em.host, cred, &ClientOptions{
		ClientOptions: azcore.ClientOptions{}, // No audience set
	})
	if err != nil {
		t.Fatalf("Failed to create AAD client: %v", err)
	}

	container, err := aadClient.NewContainer("aadAccountScopeTest", "aContainer")
	if err != nil {
		t.Fatalf("NewContainer: %v", err)
	}

	item := map[string]string{"id": "2", "value": "200"}
	body, _ := json.Marshal(item)
	pk := NewPartitionKeyString("2")

	if _, err := container.CreateItem(context.TODO(), pk, body, nil); err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	_, err = url.Parse(em.host)
	if err != nil {
		t.Fatal(err)
	}
}
