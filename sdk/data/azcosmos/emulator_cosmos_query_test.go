// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
)

func TestSinglePartitionQueryWithIndexMetrics(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
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
	documentsPerPk := 1
	createSampleItems(t, container, documentsPerPk)

	receivedIds := []string{}
	queryPager := container.NewQueryItemsPager("select * from docs c where c.someProp = '2'", NewPartitionKeyString("1"), &QueryOptions{PopulateIndexMetrics: true})
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			err = json.Unmarshal(item, &itemResponseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			receivedIds = append(receivedIds, itemResponseBody["id"].(string))
		}

		if queryPager.More() && queryResponse.ContinuationToken == "" {
			t.Fatal("Query has more pages but no continuation was provided")
		}

		if queryResponse.QueryMetrics == nil {
			t.Fatal("Query metrics were not returned")
		}

		if queryResponse.IndexMetrics == nil {
			t.Fatal("Index metrics were not returned")
		}

		if queryResponse.ActivityID == "" {
			t.Fatal("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			t.Fatal("Request charge was not returned")
		}

		if len(queryResponse.Items) != 1 && len(queryResponse.Items) != 0 {
			t.Fatalf("Expected 1 items, got %d", len(queryResponse.Items))
		}
	}

	if len(receivedIds) != 1 {
		t.Fatalf("Expected 1 documents, got %d", len(receivedIds))
	}

	if receivedIds[0] != "0" {
		t.Fatalf("Expected id 0, got %s", receivedIds[0])
	}
}

func TestSinglePartitionQuery(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
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
	documentsPerPk := 10
	createSampleItems(t, container, documentsPerPk)

	numberOfPages := 0
	receivedIds := []string{}
	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c", NewPartitionKeyString("1"), &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		numberOfPages++
		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			err = json.Unmarshal(item, &itemResponseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			receivedIds = append(receivedIds, itemResponseBody["id"].(string))
		}

		if queryPager.More() && queryResponse.ContinuationToken == "" {
			t.Fatal("Query has more pages but no continuation was provided")
		}

		if queryResponse.QueryMetrics == nil {
			t.Fatal("Query metrics were not returned")
		}

		if queryResponse.IndexMetrics != nil {
			t.Fatal("Index metrics were returned")
		}

		if queryResponse.ActivityID == "" {
			t.Fatal("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			t.Fatal("Request charge was not returned")
		}

		if len(queryResponse.Items) != 5 && len(queryResponse.Items) != 0 {
			t.Fatalf("Expected 5 items, got %d", len(queryResponse.Items))
		}

		if numberOfPages == 2 && opt.ContinuationToken != "" {
			t.Fatalf("Original options should not be modified, initial continuation was empty, now it has %v", opt.ContinuationToken)
		}
	}

	if numberOfPages != 2 {
		t.Fatalf("Expected 2 pages, got %d", numberOfPages)
	}

	if len(receivedIds) != documentsPerPk {
		t.Fatalf("Expected %d documents, got %d", documentsPerPk, len(receivedIds))
	}

	for i := 0; i < documentsPerPk; i++ {
		if receivedIds[i] != strconv.Itoa(i) {
			t.Fatalf("Expected id %d, got %s", i, receivedIds[i])
		}
	}
}

func TestSinglePartitionQueryWithParameters(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
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
	documentsPerPk := 1
	createSampleItems(t, container, documentsPerPk)

	receivedIds := []string{}
	opt := QueryOptions{
		QueryParameters: []QueryParameter{
			{"@prop", "2"},
		},
	}
	queryPager := container.NewQueryItemsPager("select * from c where c.someProp = @prop", NewPartitionKeyString("1"), &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			err = json.Unmarshal(item, &itemResponseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			receivedIds = append(receivedIds, itemResponseBody["id"].(string))
		}
	}

	if len(receivedIds) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(receivedIds))
	}
}

func TestSinglePartitionQueryWithProjection(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
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
	documentsPerPk := 10
	createSampleItems(t, container, documentsPerPk)

	numberOfPages := 0
	receivedIds := []string{}
	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select value c.id from c", NewPartitionKeyString("1"), &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		numberOfPages++
		for _, item := range queryResponse.Items {
			var itemResponseBody string
			err = json.Unmarshal(item, &itemResponseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			receivedIds = append(receivedIds, itemResponseBody)
		}

		if queryPager.More() && queryResponse.ContinuationToken == "" {
			t.Fatal("Query has more pages but no continuation was provided")
		}

		if queryResponse.QueryMetrics == nil {
			t.Fatal("Query metrics were not returned")
		}

		if queryResponse.IndexMetrics != nil {
			t.Fatal("Index metrics were returned")
		}

		if queryResponse.ActivityID == "" {
			t.Fatal("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			t.Fatal("Request charge was not returned")
		}

		if len(queryResponse.Items) != 5 && len(queryResponse.Items) != 0 {
			t.Fatalf("Expected 5 items, got %d", len(queryResponse.Items))
		}

		if numberOfPages == 2 && opt.ContinuationToken != "" {
			t.Fatalf("Original options should not be modified, initial continuation was empty, now it has %v", opt.ContinuationToken)
		}
	}

	if numberOfPages != 2 {
		t.Fatalf("Expected 2 pages, got %d", numberOfPages)
	}

	if len(receivedIds) != documentsPerPk {
		t.Fatalf("Expected %d documents, got %d", documentsPerPk, len(receivedIds))
	}

	for i := 0; i < documentsPerPk; i++ {
		if receivedIds[i] != strconv.Itoa(i) {
			t.Fatalf("Expected id %d, got %s", i, receivedIds[i])
		}
	}
}

func createSampleItems(t *testing.T, container *ContainerClient, documentsPerPk int) {
	for i := 0; i < documentsPerPk; i++ {
		item := map[string]string{
			"id":       strconv.Itoa(i),
			"pk":       "1",
			"someProp": "2",
		}

		marshalled, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}

		_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		item2 := map[string]string{
			"id":       strconv.Itoa(i),
			"pk":       "2",
			"someProp": "2",
		}

		marshalled, err = json.Marshal(item2)
		if err != nil {
			t.Fatal(err)
		}

		_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("2"), marshalled, nil)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}
	}
}
