// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestItemCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_item aContainer", "read_item aContainer", "replace_item aContainer", "upsert_item aContainer", "delete_item aContainer", "patch_item aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
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

	item := map[string]interface{}{
		"id":          "1",
		"value":       "2",
		"count":       3,
		"description": "4",
	}

	container, _ := database.NewContainer("aContainer")
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

	patchItem := PatchOperations{}
	patchItem.AppendReplace("/value", "5")
	patchItem.AppendSet("/hello", "world")
	patchItem.AppendAdd("/foo", "bar")
	patchItem.AppendRemove("/description")
	patchItem.AppendIncrement("/count", 1)

	itemResponse, err = container.PatchItem(context.TODO(), pk, "1", patchItem, nil)
	if err != nil {
		t.Fatalf("Failed to patch item: %v", err)
	}

	// No content on write by default
	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}

	itemResponse, _ = container.ReadItem(context.TODO(), pk, "1", nil)

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}

	if itemResponseBody["value"] != "5" {
		t.Fatalf("Expected value to be 5, got %v", itemResponseBody["id"])
	}

	if itemResponseBody["hello"] != "world" {
		t.Fatalf("Expected hello to be world, got %v", itemResponseBody["hello"])
	}

	if itemResponseBody["foo"] != "bar" {
		t.Fatalf("Expected foo to be bar, got %v", itemResponseBody["foo"])
	}

	if itemResponseBody["count"].(float64) != float64(4) {
		t.Fatalf("Expected count to be 4, got %v", itemResponseBody["count"])
	}

	if itemResponseBody["toremove"] != nil {
		t.Fatalf("Expected toremove to be nil, got %v", itemResponseBody)
	}

	itemResponse, err = container.DeleteItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}
}

func TestItemCRUDforNullPartitionKey(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_item aContainer", "read_item aContainer", "replace_item aContainer", "upsert_item aContainer", "delete_item aContainer", "patch_item aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/partitionKey"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	item := map[string]interface{}{
		"partitionKey": nil,
		"id":           "1",
		"value":        "2",
		"count":        3,
		"description":  "4",
	}

	container, _ := database.NewContainer("aContainer")
	pk := NullPartitionKey

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

	patchItem := PatchOperations{}
	patchItem.AppendReplace("/value", "5")
	patchItem.AppendSet("/hello", "world")
	patchItem.AppendAdd("/foo", "bar")
	patchItem.AppendRemove("/description")
	patchItem.AppendIncrement("/count", 1)

	itemResponse, err = container.PatchItem(context.TODO(), pk, "1", patchItem, nil)
	if err != nil {
		t.Fatalf("Failed to patch item: %v", err)
	}

	// No content on write by default
	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}

	itemResponse, _ = container.ReadItem(context.TODO(), pk, "1", nil)

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}

	if itemResponseBody["value"] != "5" {
		t.Fatalf("Expected value to be 5, got %v", itemResponseBody["id"])
	}

	if itemResponseBody["hello"] != "world" {
		t.Fatalf("Expected hello to be world, got %v", itemResponseBody["hello"])
	}

	if itemResponseBody["foo"] != "bar" {
		t.Fatalf("Expected foo to be bar, got %v", itemResponseBody["foo"])
	}

	if itemResponseBody["count"].(float64) != float64(4) {
		t.Fatalf("Expected count to be 4, got %v", itemResponseBody["count"])
	}

	if itemResponseBody["toremove"] != nil {
		t.Fatalf("Expected toremove to be nil, got %v", itemResponseBody)
	}

	itemResponse, err = container.DeleteItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}
}

func TestItemConcurrent(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
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

	container, _ := database.NewContainer("aContainer")

	item := map[string]interface{}{
		"id":          "1",
		"value":       "2",
		"count":       3,
		"description": "4",
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	pk := NewPartitionKeyString("1")

	_, err = container.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Execute 50 concurrent operations
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = container.ReadItem(context.TODO(), pk, "1", nil)
		}()
	}
	wg.Wait()
}

func TestItemIdEncodingRoutingGW(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer("aContainer")

	verifyEncodingScenario(t, container, "RoutingGW - PlainVanillaId", "Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithWhitespaces", "This is a test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdStartingWithWhitespaces", " Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdEndingWithWhitespace", "Test ", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdEndingWithWhitespaces", "Test  ", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithAllowedSpecialCharacters", "WithAllowedSpecial,=.:~+-@()^${}[]!_Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithBase64EncodedIdCharacters", strings.ReplaceAll("BQE1D3PdG4N4bzU9TKaCIM3qc0TVcZ2/Y3jnsRfwdHC1ombkX3F1dot/SG0/UTq9AbgdX3kOWoP6qL6lJqWeKgV3zwWWPZO/t5X0ehJzv9LGkWld07LID2rhWhGT6huBM6Q=", "/", "-"), http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdEndingWithPercentEncodedWhitespace", "IdEndingWithPercentEncodedWhitespace%20", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithPercentEncodedSpecialChar", "WithPercentEncodedSpecialChar%E9%B1%80", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithDisallowedCharQuestionMark", "Disallowed?Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithDisallowedCharForwardSlash", "Disallowed/Chars", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithDisallowedCharBackSlash", "Disallowed\\Chars", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithDisallowedCharPoundSign", "Disallowed#Chars", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithCarriageReturn", "With\rCarriageReturn", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithTab", "With\tTab", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithLineFeed", "With\nLineFeed", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "RoutingGW - IdWithUnicodeCharacters", "WithUnicode鱀", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
}

func TestItemIdEncodingComputeGW(t *testing.T) {
	emulatorTests := newEmulatorTestsWithComputeGateway(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer("aContainer")

	verifyEncodingScenario(t, container, "ComputeGW-PlainVanillaId", "Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithWhitespaces", "This is a test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdStartingWithWhitespaces", " Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdEndingWithWhitespace", "Test ", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdEndingWithWhitespaces", "Test  ", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithAllowedSpecialCharacters", "WithAllowedSpecial,=.:~+-@()^${}[]!_Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithBase64EncodedIdCharacters", strings.ReplaceAll("BQE1D3PdG4N4bzU9TKaCIM3qc0TVcZ2/Y3jnsRfwdHC1ombkX3F1dot/SG0/UTq9AbgdX3kOWoP6qL6lJqWeKgV3zwWWPZO/t5X0ehJzv9LGkWld07LID2rhWhGT6huBM6Q=", "/", "-"), http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdEndingWithPercentEncodedWhitespace", "IdEndingWithPercentEncodedWhitespace%20", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithPercentEncodedSpecialChar", "WithPercentEncodedSpecialChar%E9%B1%80", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithDisallowedCharQuestionMark", "Disallowed?Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithDisallowedCharForwardSlash", "Disallowed/Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithDisallowedCharBackSlash", "Disallowed\\Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithDisallowedCharPoundSign", "Disallowed#Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithCarriageReturn", "With\rCarriageReturn", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithTab", "With\tTab", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithLineFeed", "With\nLineFeed", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "ComputeGW-IdWithUnicodeCharacters", "WithUnicode鱀", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
}

func TestItemCRUDHierarchicalPartitionKey(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUDHierarchicalPartitionKey")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/id", "/type"},
			Kind:    PartitionKeyKindMultiHash,
			Version: 2,
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	container, err := database.NewContainer("aContainer")
	if err != nil {
		t.Fatalf("Failed to get container: %v", err)
	}

	itemAlpha := map[string]interface{}{
		"id":    "1",
		"type":  "alpha",
		"value": "0",
	}

	itemBeta := map[string]interface{}{
		"id":    "1",
		"type":  "beta",
		"value": "0",
	}

	pkAlpha := NewPartitionKey().AppendString("1").AppendString("alpha")
	pkBeta := NewPartitionKey().AppendString("1").AppendString("beta")

	marshalledAlpha, err := json.Marshal(itemAlpha)
	if err != nil {
		t.Fatal(err)
	}

	marshalledBeta, err := json.Marshal(itemBeta)
	if err != nil {
		t.Fatal(err)
	}

	item0Res, err := container.CreateItem(context.TODO(), pkAlpha, marshalledAlpha, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	if item0Res.SessionToken == nil {
		t.Fatalf("Session token is empty")
	}

	if len(item0Res.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", item0Res.Value)
	}

	item1Res, err := container.CreateItem(context.TODO(), pkBeta, marshalledBeta, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	if item1Res.SessionToken == nil {
		t.Fatalf("Session token is empty")
	}

	if len(item1Res.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", item1Res.Value)
	}

	item0Res, err = container.ReadItem(context.TODO(), pkAlpha, "1", nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if len(item0Res.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", item0Res.Value)
	}

	item1Res, err = container.ReadItem(context.TODO(), pkBeta, "1", nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if len(item1Res.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", item1Res.Value)
	}

	var item0ResBody map[string]interface{}
	err = json.Unmarshal(item0Res.Value, &item0ResBody)

	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}

	if item0ResBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", item0ResBody["id"])
	}

	if item0ResBody["type"] != "alpha" {
		t.Fatalf("Expected type to be alpha, got %v", item0ResBody["type"])
	}

	if item0ResBody["value"] != "0" {
		t.Fatalf("Expected value to be 0, got %v", item0ResBody["value"])
	}

	var item1ResBody map[string]interface{}
	err = json.Unmarshal(item1Res.Value, &item1ResBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}

	if item1ResBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", item1ResBody["id"])
	}

	if item1ResBody["type"] != "beta" {
		t.Fatalf("Expected type to be beta, got %v", item1ResBody["type"])
	}

	if item1ResBody["value"] != "0" {
		t.Fatalf("Expected value to be 0, got %v", item1ResBody["value"])
	}

	pager := container.NewQueryItemsPager("SELECT * FROM c", pkAlpha, nil)

	var alphaItems []map[string]interface{}
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to get next page: %v", err)
		}

		for _, item := range page.Items {
			var itemBody map[string]interface{}
			err = json.Unmarshal(item, &itemBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal item response: %v", err)
			}

			alphaItems = append(alphaItems, itemBody)
		}
	}

	if len(alphaItems) != 1 {
		t.Fatalf("Expected 1 item, got %v", len(alphaItems))
	}

	if !reflect.DeepEqual(alphaItems[0], item0ResBody) {
		t.Fatalf("Expected %v, got %v", item0ResBody, alphaItems[0])
	}

	pager = container.NewQueryItemsPager("SELECT * FROM c", pkBeta, nil)

	var betaItems []map[string]interface{}
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to get next page: %v", err)
		}

		for _, item := range page.Items {
			var itemBody map[string]interface{}
			err = json.Unmarshal(item, &itemBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal item response: %v", err)
			}

			betaItems = append(betaItems, itemBody)
		}
	}

	if len(betaItems) != 1 {
		t.Fatalf("Expected 1 item, got %v", len(betaItems))
	}

	if !reflect.DeepEqual(betaItems[0], item1ResBody) {
		t.Fatalf("Expected %v, got %v", item1ResBody, betaItems[0])
	}
}

func verifyEncodingScenario(t *testing.T, container *ContainerClient, name string, id string, expectedCreate int, expectedRead int, expectedReplace int, expectedDelete int) {
	item := map[string]interface{}{
		"id": id,
		"pk": id,
	}

	pk := NewPartitionKeyString(id)

	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	itemResponse, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
	verifyEncodingScenarioResponse(t, name+"Create", itemResponse, err, expectedCreate)
	itemResponse, err = container.ReadItem(context.TODO(), pk, id, nil)
	verifyEncodingScenarioResponse(t, name+"Read", itemResponse, err, expectedRead)
	itemResponse, err = container.ReplaceItem(context.TODO(), pk, id, marshalled, nil)
	verifyEncodingScenarioResponse(t, name+"Replace", itemResponse, err, expectedReplace)
	itemResponse, err = container.DeleteItem(context.TODO(), pk, id, nil)
	verifyEncodingScenarioResponse(t, name+"Delete", itemResponse, err, expectedDelete)
}

func verifyEncodingScenarioResponse(t *testing.T, name string, itemResponse ItemResponse, err error, expectedStatus int) {
	if err != nil {
		var responseErr *azcore.ResponseError
		errors.As(err, &responseErr)
		if responseErr.StatusCode != expectedStatus {
			t.Fatalf("[%s] Expected status code %d, got %d, %s", name, expectedStatus, responseErr.StatusCode, err)
		}
	} else if itemResponse.RawResponse.StatusCode != expectedStatus {
		t.Fatalf("[%s] Expected status code %d, got %d", name, expectedStatus, itemResponse.RawResponse.StatusCode)
	}
}

func TestReadMany(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "readMany")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer("aContainer")

	item1 := map[string]interface{}{
		"id":    "1",
		"pk":    "pk1",
		"value": "first",
	}

	item2 := map[string]interface{}{
		"id":    "2",
		"pk":    "pk2",
		"value": "second",
	}

	item3 := map[string]interface{}{
		"id":    "3",
		"pk":    "pk3",
		"value": "third",
	}

	pk1 := NewPartitionKeyString("pk1")
	pk2 := NewPartitionKeyString("pk2")
	pk3 := NewPartitionKeyString("pk3")

	marshalled1, _ := json.Marshal(item1)
	marshalled2, _ := json.Marshal(item2)
	marshalled3, _ := json.Marshal(item3)

	_, err = container.CreateItem(context.TODO(), pk1, marshalled1, nil)
	if err != nil {
		t.Fatalf("Failed to create item1: %v", err)
	}

	_, err = container.CreateItem(context.TODO(), pk2, marshalled2, nil)
	if err != nil {
		t.Fatalf("Failed to create item2: %v", err)
	}

	_, err = container.CreateItem(context.TODO(), pk3, marshalled3, nil)
	if err != nil {
		t.Fatalf("Failed to create item3: %v", err)
	}

	coordinates := []ItemCoordinate{
		{ID: "1", PartitionKey: "pk1"},
		{ID: "2", PartitionKey: "pk2"},
		{ID: "3", PartitionKey: "pk3"},
	}

	response, err := container.ReadMany(context.TODO(), coordinates, nil)
	if err != nil {
		t.Fatalf("Failed to read many items: %v", err)
	}

	if len(response.Items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(response.Items))
	}

	if response.TotalRequestCharge <= 0 {
		t.Fatalf("Expected positive request charge, got %f", response.TotalRequestCharge)
	}

	var result1 map[string]interface{}
	err = json.Unmarshal(response.Items[0], &result1)
	if err != nil {
		t.Fatalf("Failed to unmarshal item 1: %v", err)
	}
	if result1["id"] != "1" || result1["value"] != "first" {
		t.Fatalf("Item 1 mismatch: %v", result1)
	}

	var result2 map[string]interface{}
	err = json.Unmarshal(response.Items[1], &result2)
	if err != nil {
		t.Fatalf("Failed to unmarshal item 2: %v", err)
	}
	if result2["id"] != "2" || result2["value"] != "second" {
		t.Fatalf("Item 2 mismatch: %v", result2)
	}

	var result3 map[string]interface{}
	err = json.Unmarshal(response.Items[2], &result3)
	if err != nil {
		t.Fatalf("Failed to unmarshal item 3: %v", err)
	}
	if result3["id"] != "3" || result3["value"] != "third" {
		t.Fatalf("Item 3 mismatch: %v", result3)
	}
}
