// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/assert"
)

func TestSinglePartitionQueryWithIndexMetrics(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PopulateIndexMetrics: true}
	queryPager := container.NewQueryItemsPager("select * from docs c where c.someProp = 'some_4'", NewPartitionKeyString("1"), &opt)
	receivedIds, err := collectResultIds(t, 1, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t,
		[]string{"4"},
		receivedIds)
}

func TestSinglePartitionQuery(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PageSizeHint: 5}
	// We include an ORDER BY to ensure that ORDER BY statements are still allowed even when we're setting the cross-partition flag, as long as the query itself is still single-partition.
	queryPager := container.NewQueryItemsPager("select * from c order by c.id", NewPartitionKeyString("1"), &opt)
	receivedIds, err := collectResultIds(t, 2, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t,
		// Single partition result
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		receivedIds)
}

func TestSinglePartitionQueryInline(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PageSizeHint: 5}
	// We can specify the partition key inline in the query itself.
	queryPager := container.NewQueryItemsPager("select * from c where c.pk = '1' order by c.id", NewPartitionKey(), &opt)
	receivedIds, err := collectResultIds(t, 2, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t,
		// Single partition result
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		receivedIds)
}

func TestSinglePartitionQueryWithParameters(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{
		QueryParameters: []QueryParameter{
			{"@prop", "some_4"},
		},
	}
	// We include an ORDER BY to ensure that ORDER BY statements are still allowed even when we're setting the cross-partition flag, as long as the query itself is still single-partition.
	queryPager := container.NewQueryItemsPager("select * from c where c.someProp = @prop order by c.id", NewPartitionKeyString("1"), &opt)
	receivedIds, err := collectResultIds(t, 1, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t,
		[]string{"4"},
		receivedIds)
}

func TestSinglePartitionQueryWithProjection(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PageSizeHint: 5}

	// We include an ORDER BY to ensure that ORDER BY statements are still allowed even when we're setting the cross-partition flag, as long as the query itself is still single-partition.
	queryPager := container.NewQueryItemsPager("select value c.id from c order by c.id", NewPartitionKeyString("1"), &opt)
	receivedIds, err := collectResultIds(t, 2, queryPager, &opt, parseValueAsId)
	assert.NoError(t, err)
	assert.Equal(t,
		// Single partition result
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		receivedIds)
}

func TestCrossPartitionQuery(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c", NewPartitionKey(), &opt)
	receivedIds, err := collectResultIds(t, 5, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t,
		// Partitions should be interleaved and not re-ordered by ID.
		[]string{"0", "10", "1", "11", "2", "12", "3", "13", "4", "14", "5", "15", "6", "16", "7", "17", "8", "18", "9", "19"},
		receivedIds)
}

func TestCrossPartitionQueryFailsIfGatewayCannotSatisfyRequest(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItems(container, 2, 10))

	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c order by c.id", NewPartitionKey(), &opt)
	receivedIds, err := collectResultIds(t, 5, queryPager, &opt, parseIdProperty)
	assert.Nil(t, receivedIds)

	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Failed to query items: "))
	assert.True(t, strings.Contains(err.Error(), "BadRequest"))
	assert.True(t, strings.Contains(err.Error(), "cross partition query can not be directly served by the gateway"))
}

func TestHierarchicalPartitionQuerySinglePartition(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/parent", "/child"},
			Version: 2,
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "1",
		"parent": "parent1",
		"child":  "child1",
	}, NewPartitionKeyString("parent1").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "2",
		"parent": "parent1",
		"child":  "child2",
	}, NewPartitionKeyString("parent1").AppendString("child2")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "3",
		"parent": "parent2",
		"child":  "child1",
	}, NewPartitionKeyString("parent2").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "4",
		"parent": "parent2",
		"child":  "child2",
	}, NewPartitionKeyString("parent2").AppendString("child2")))

	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c order by c.id", NewPartitionKeyString("parent1").AppendString("child2"), &opt)
	receivedIds, err := collectResultIds(t, 1, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t, []string{"2"}, receivedIds)
}

func TestHierarchicalPartitionQueryParentPartition(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/parent", "/child"},
			Version: 2,
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "1",
		"parent": "parent1",
		"child":  "child1",
	}, NewPartitionKeyString("parent1").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "2",
		"parent": "parent1",
		"child":  "child2",
	}, NewPartitionKeyString("parent1").AppendString("child2")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "3",
		"parent": "parent2",
		"child":  "child1",
	}, NewPartitionKeyString("parent2").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "4",
		"parent": "parent2",
		"child":  "child2",
	}, NewPartitionKeyString("parent2").AppendString("child2")))

	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c where c.parent = 'parent1'", NewPartitionKey(), &opt)
	receivedIds, err := collectResultIds(t, 1, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2"}, receivedIds)
}

func TestHierarchicalPartitionQueryNoPartition(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"query_items aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "queryTests")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/parent", "/child"},
			Version: 2,
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	assert.NoError(t, err)

	container, _ := database.NewContainer("aContainer")
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "1",
		"parent": "parent1",
		"child":  "child1",
	}, NewPartitionKeyString("parent1").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "2",
		"parent": "parent1",
		"child":  "child2",
	}, NewPartitionKeyString("parent1").AppendString("child2")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "3",
		"parent": "parent2",
		"child":  "child1",
	}, NewPartitionKeyString("parent2").AppendString("child1")))
	assert.NoError(t, createSampleItem(container, map[string]interface{}{
		"id":     "4",
		"parent": "parent2",
		"child":  "child2",
	}, NewPartitionKeyString("parent2").AppendString("child2")))

	opt := QueryOptions{PageSizeHint: 5}
	queryPager := container.NewQueryItemsPager("select * from c", NewPartitionKey(), &opt)
	receivedIds, err := collectResultIds(t, 1, queryPager, &opt, parseIdProperty)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3", "4"}, receivedIds)
}

func createSampleItems(container *ContainerClient, partitions int, documentsPerPartition int) error {
	for i := 0; i < documentsPerPartition; i++ {
		// We insert documents alternating between partitions.
		// This simulates a kind of "worst-case" illustration of how cross-partition queries can interleave results since the "default" ordering is by insertion order.
		for pk := 0; pk < partitions; pk++ {
			id := strconv.Itoa(i + (pk * documentsPerPartition))
			pkStr := strconv.Itoa(pk + 1)
			err := createSampleItem(container, map[string]interface{}{
				"id":       id,
				"pk":       pkStr,
				"someProp": fmt.Sprintf("some_%s", id),
			}, NewPartitionKeyString(pkStr))
			if err != nil {
				return fmt.Errorf("Failed to create sample item: %v", err)
			}
		}
	}
	return nil
}

func createSampleItem(container *ContainerClient, item map[string]interface{}, pk PartitionKey) error {
	marshalled, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = container.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		return fmt.Errorf("Failed to create item: %v", err)
	}
	return nil
}

func parseValueAsId(item []byte) (string, error) {
	var itemResponseBody string
	err := json.Unmarshal(item, &itemResponseBody)
	if err != nil {
		return "", err
	}
	return itemResponseBody, nil
}

func parseIdProperty(item []byte) (string, error) {
	var itemResponseBody map[string]interface{}
	err := json.Unmarshal(item, &itemResponseBody)
	if err != nil {
		return "", err
	}
	return itemResponseBody["id"].(string), nil
}

func collectResultIds(t *testing.T, expectedPageCount int, queryPager *runtime.Pager[QueryItemsResponse], originalOptions *QueryOptions, idParser func([]byte) (string, error)) ([]string, error) {
	ids := []string{}
	pageCount := 0
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("Failed to query items: %v", err)
		}

		pageCount++
		for _, item := range queryResponse.Items {
			id, err := idParser(item)
			if err != nil {
				return nil, fmt.Errorf("Failed to unmarshal: %v", err)
			}
			ids = append(ids, id)
		}

		if queryPager.More() && queryResponse.ContinuationToken == nil {
			return nil, fmt.Errorf("Query has more pages but no continuation was provided")
		}

		if queryResponse.QueryMetrics == nil {
			return nil, fmt.Errorf("Query metrics were not returned")
		}

		if !originalOptions.PopulateIndexMetrics && queryResponse.IndexMetrics != nil {
			return nil, fmt.Errorf("Index metrics were returned but not requested")
		} else if originalOptions.PopulateIndexMetrics && queryResponse.IndexMetrics == nil {
			return nil, fmt.Errorf("Index metrics were requested but not returned")
		}

		if queryResponse.ActivityID == "" {
			return nil, fmt.Errorf("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			return nil, fmt.Errorf("Request charge was not returned")
		}

		if originalOptions.PageSizeHint > 0 && len(queryResponse.Items) > int(originalOptions.PageSizeHint) {
			return nil, fmt.Errorf("Expected 1-%d items, got %d", int(originalOptions.PageSizeHint), len(queryResponse.Items))
		}

		if pageCount == expectedPageCount && originalOptions.ContinuationToken != nil {
			return nil, fmt.Errorf("Original options should not be modified, initial continuation was empty, now it has %v", originalOptions.ContinuationToken)
		}
	}
	assert.Equal(t, expectedPageCount, pageCount)
	return ids, nil
}
