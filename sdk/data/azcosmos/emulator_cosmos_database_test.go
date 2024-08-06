// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
)

func TestDatabaseCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, spanMatcher{
		ExpectedSpans: []string{"create_database baseDbTest", "read_database baseDbTest", "delete_database baseDbTest", "read_database_throughput baseDbTest"},
	}))

	database := DatabaseProperties{ID: "baseDbTest"}

	resp, err := client.CreateDatabase(context.TODO(), database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	db, _ := client.NewDatabase("baseDbTest")
	resp, err = db.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	receivedIds := []string{}
	opt := QueryDatabasesOptions{
		QueryParameters: []QueryParameter{
			{"@id", "baseDbTest"},
		},
	}
	queryPager := client.NewQueryDatabasesPager("SELECT * FROM root r WHERE r.id = @id", &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query databases: %v", err)
		}

		for _, db := range queryResponse.Databases {
			receivedIds = append(receivedIds, db.ID)
		}
	}

	if len(receivedIds) != 1 {
		t.Fatalf("Expected 1 database, got %d", len(receivedIds))
	}

	throughputResponse, err := db.ReadThroughput(context.TODO(), nil)
	if err == nil {
		t.Fatalf("Expected not finding throughput but instead got : %v", throughputResponse)
	}

	resp, err = db.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}

func TestDatabaseWithOfferCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, spanMatcher{
		ExpectedSpans: []string{"create_database baseDbTest", "read_database baseDbTest", "delete_database baseDbTest", "read_database_throughput baseDbTest", "replace_database_throughput baseDbTest"},
	}))

	database := DatabaseProperties{ID: "baseDbTest"}
	tp := NewManualThroughputProperties(400)
	resp, err := client.CreateDatabase(context.TODO(), database, &CreateDatabaseOptions{ThroughputProperties: &tp})
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	db, _ := client.NewDatabase("baseDbTest")
	resp, err = db.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	throughputResponse, err := db.ReadThroughput(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read throughput: %v", err)
	}

	mt, hasManualThroughput := throughputResponse.ThroughputProperties.ManualThroughput()
	if !hasManualThroughput {
		t.Fatalf("Expected manual throughput to be available")
	}

	if mt != 400 {
		t.Errorf("Unexpected throughput: %v", mt)
	}

	newScale := NewManualThroughputProperties(500)
	throughputResponse, err = db.ReplaceThroughput(context.TODO(), newScale, nil)
	if err != nil {
		t.Fatalf("Failed to replace throughput: %v", err)
	}

	mt, hasManualThroughput = throughputResponse.ThroughputProperties.ManualThroughput()
	if !hasManualThroughput {
		t.Fatalf("Expected manual throughput to be available")
	}

	if mt != 500 {
		t.Errorf("Unexpected throughput: %v", mt)
	}

	resp, err = db.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}
