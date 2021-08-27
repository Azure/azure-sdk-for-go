// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
)

func TestDatabaseCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests()
	client := emulatorTests.getClient(t)

	database := CosmosDatabaseProperties{Id: "baseDbTest"}

	resp, err := client.CreateDatabase(context.TODO(), database, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read database: %v", err)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}

func TestDatabaseWithOfferCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests()
	client := emulatorTests.getClient(t)

	database := CosmosDatabaseProperties{Id: "baseDbTest"}
	tp := NewManualThroughputProperties(400)
	resp, err := client.CreateDatabase(context.TODO(), database, tp, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read database: %v", err)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}
