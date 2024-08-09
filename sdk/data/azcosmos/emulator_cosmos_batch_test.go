// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"
)

func TestItemTransactionalBatch(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, spanMatcher{
		ExpectedSpans: []string{"execute_batch aContainer"},
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

	container, _ := client.NewContainer("itemCRUD", "aContainer")

	pkValue := NewPartitionKeyString("tBatch")
	batch := container.NewTransactionalBatch(pkValue)

	batch.CreateItem(emulatorTests.marshallItem("test", "tBatch"), nil)
	batch.CreateItem(emulatorTests.marshallItem("test2", "tBatch"), nil)
	batch.CreateItem(emulatorTests.marshallItem("test5", "tBatch"), nil)

	// Default behavior has no content body
	response, err := container.ExecuteTransactionalBatch(context.TODO(), batch, nil)
	if err != nil {
		t.Fatalf("Failed to execute batch: %v", err)
	}

	if len(response.OperationResults) != 3 {
		t.Fatalf("Expected 3 operation results, got %v", len(response.OperationResults))
	}

	if !response.Success {
		t.Fatalf("Expected committed to be true, got false")
	}

	for _, operationResult := range response.OperationResults {
		if operationResult.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code 201, got %v", operationResult.StatusCode)
		}

		if operationResult.RequestCharge <= 0 {
			t.Fatalf("Expected RequestCharge to be greater than 0, got %v", operationResult.RequestCharge)
		}

		if operationResult.ETag == "" {
			t.Fatalf("Expected ETag to be non-empty, got %v", operationResult.ETag)
		}

		if operationResult.ResourceBody != nil {
			t.Fatalf("Expected ResourceBody to be nil, got %v", operationResult.ResourceBody)
		}
	}

	batch2 := container.NewTransactionalBatch(pkValue)

	batch2.CreateItem(emulatorTests.marshallItem("test3", "tBatch"), nil)
	batch2.ReadItem("test2", nil)
	batch2.DeleteItem("test", nil)

	// If there is a read operation, body should be included
	response2, err := container.ExecuteTransactionalBatch(context.TODO(), batch2, nil)
	if err != nil {
		t.Fatalf("Failed to execute batch: %v", err)
	}

	if !response2.Success {
		t.Fatalf("Expected committed to be true, got false")
	}

	if len(response2.OperationResults) != 3 {
		t.Fatalf("Expected 3 operation results, got %v", len(response2.OperationResults))
	}

	for index, operationResult := range response2.OperationResults {
		if index == 0 && operationResult.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code 201, got %v", operationResult.StatusCode)
		}

		if index == 1 && operationResult.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %v", operationResult.StatusCode)
		}

		if index == 2 && operationResult.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status code 204, got %v", operationResult.StatusCode)
		}

		if operationResult.RequestCharge <= 0 {
			t.Fatalf("Expected RequestCharge to be greater than 0, got %v", operationResult.RequestCharge)
		}

		if index < 2 && operationResult.ETag == "" {
			t.Fatalf("Expected ETag to be non-empty, got %v", operationResult.ETag)
		}

		if index < 2 && operationResult.ResourceBody == nil {
			t.Fatalf("Expected ResourceBody to be not-nil, got %v", operationResult.ResourceBody)
		}

		if index == 2 && operationResult.ResourceBody != nil {
			t.Fatalf("Expected ResourceBody to be nil, got %v", operationResult.ResourceBody)
		}
	}

	// Forcing body through options
	batch3 := container.NewTransactionalBatch(pkValue)

	batch3.UpsertItem(emulatorTests.marshallItem("test4", "tBatch"), nil)
	batch3.ReplaceItem("test3", emulatorTests.marshallItem("test3", "tBatch"), nil)
	p := PatchOperations{}
	p.AppendAdd("/newField", "newValue")
	batch3.PatchItem("test5", p, nil)

	response3, err := container.ExecuteTransactionalBatch(context.TODO(), batch3, &TransactionalBatchOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to execute batch: %v", err)
	}

	if !response3.Success {
		t.Fatalf("Expected Success to be true, got false")
	}

	if len(response3.OperationResults) != 3 {
		t.Fatalf("Expected 3 operation results, got %v", len(response3.OperationResults))
	}

	for index, operationResult := range response3.OperationResults {
		if index == 0 && operationResult.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code 201, got %v", operationResult.StatusCode)
		}

		if index == 1 && operationResult.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %v", operationResult.StatusCode)
		}

		if index == 2 && operationResult.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %v", operationResult.StatusCode)
		}

		if operationResult.RequestCharge <= 0 {
			t.Fatalf("Expected RequestCharge to be greater than 0, got %v", operationResult.RequestCharge)
		}

		if operationResult.ETag == "" {
			t.Fatalf("Expected ETag to be non-empty, got %v", operationResult.ETag)
		}

		if operationResult.ResourceBody == nil {
			t.Fatalf("Expected ResourceBody not to be nil, got %v", operationResult.ResourceBody)
		}
	}
}

func TestItemTransactionalBatchError(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, spanMatcher{
		ExpectedSpans: []string{"execute_batch aContainer"},
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

	container, _ := client.NewContainer("itemCRUD", "aContainer")

	pkValue := NewPartitionKeyString("tBatch")

	_, err = container.CreateItem(context.TODO(), pkValue, emulatorTests.marshallItem("test", "tBatch"), nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	batch := container.NewTransactionalBatch(pkValue)

	batch.CreateItem(emulatorTests.marshallItem("test2", "tBatch"), nil)
	batch.CreateItem(emulatorTests.marshallItem("test", "tBatch"), nil)

	response, err := container.ExecuteTransactionalBatch(context.TODO(), batch, &TransactionalBatchOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to execute batch: %v", err)
	}

	if response.RawResponse.StatusCode != http.StatusMultiStatus {
		t.Fatalf("Expected status code 207, got %v", response.RawResponse.StatusCode)
	}

	if response.Success {
		t.Fatalf("Expected Success to be false, got true")
	}

	if len(response.OperationResults) != 2 {
		t.Fatalf("Expected 2 operation results, got %v", len(response.OperationResults))
	}

	for index, operationResult := range response.OperationResults {
		if index == 0 && operationResult.StatusCode != http.StatusFailedDependency {
			t.Fatalf("Expected status code 424, got %v", operationResult.StatusCode)
		}

		if index == 1 && operationResult.StatusCode != http.StatusConflict {
			t.Fatalf("Expected status code 409, got %v", operationResult.StatusCode)
		}
	}
}
