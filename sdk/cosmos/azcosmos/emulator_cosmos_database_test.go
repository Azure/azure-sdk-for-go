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
	resp, err := client.AddDatabase(context.TODO(), database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Get(context.TODO(), nil)
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
