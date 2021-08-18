// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
)

func TestCosmosClientCreateDatabase(t *testing.T) {
	emulatorTests := newEmulatorTests()
	cred, _ := NewSharedKeyCredential(emulatorTests.key)
	client, err := NewCosmosClient(emulatorTests.host, cred, nil)
	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	database := CosmosDatabaseProperties{Id: "baseDbTest"}
	resp, err := client.AddDatabase(context.TODO(), database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.RawResponse.StatusCode != 201 {
		t.Fatalf("Failed to create database: %v", resp.RawResponse)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Get(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read database: %v", err)
	}

	if resp.RawResponse.StatusCode != 200 {
		t.Fatalf("Failed to read database: %v", resp.RawResponse)
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	resp, err = resp.DatabaseProperties.Database.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}

	if resp.RawResponse.StatusCode != 204 {
		t.Fatalf("Failed to delete database: %v", resp.RawResponse)
	}
}
