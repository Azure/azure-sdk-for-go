// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestItemCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

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

	item := map[string]string{
		"id":    "1",
		"value": "2",
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

	if itemResponse.SessionToken == "" {
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

func TestItemIdEncoding(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

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

	verifyEncodingScenario(t, container, "PlainVanillaId", "Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdWithWhitespaces", "This is a test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdStartingWithWhitespaces", " Test", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdEndingWithWhitespace", "Test ", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "IdEndingWithWhitespaces", "Test  ", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "IdWithAllowedSpecialCharacters", "WithAllowedSpecial,=.:~+-@()^${}[]!_Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdWithBase64EncodedIdCharacters", strings.Replace("BQE1D3PdG4N4bzU9TKaCIM3qc0TVcZ2/Y3jnsRfwdHC1ombkX3F1dot/SG0/UTq9AbgdX3kOWoP6qL6lJqWeKgV3zwWWPZO/t5X0ehJzv9LGkWld07LID2rhWhGT6huBM6Q=", "/", "-", -1), http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdEndingWithPercentEncodedWhitespace", "IdEndingWithPercentEncodedWhitespace%20", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "IdWithPercentEncodedSpecialChar", "WithPercentEncodedSpecialChar%E9%B1%80", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "IdWithDisallowedCharQuestionMark", "Disallowed?Chars", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	verifyEncodingScenario(t, container, "IdWithDisallowedCharForwardSlash", "Disallowed/Chars", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "IdWithDisallowedCharBackSlash", "Disallowed\\Chars", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "IdWithDisallowedCharPoundSign", "Disallowed#Chars", http.StatusCreated, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized)
	verifyEncodingScenario(t, container, "IdWithCarriageReturn", "With\rCarriageReturn", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "IdWithTab", "With\tTab", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "IdWithLineFeed", "With\nLineFeed", http.StatusCreated, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest)
	verifyEncodingScenario(t, container, "IdWithUnicodeCharacters", "WithUnicode鱀", http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
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
	} else {
		if itemResponse.RawResponse.StatusCode != expectedStatus {
			t.Fatalf("[%s] Expected status code %d, got %d", name, expectedStatus, itemResponse.RawResponse.StatusCode)
		}
	}
}
