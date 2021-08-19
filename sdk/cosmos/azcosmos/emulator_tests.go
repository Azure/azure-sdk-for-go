// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type emulatorTests struct {
	host string
	key  string
}

func newEmulatorTests() *emulatorTests {
	return &emulatorTests{
		host: "https://localhost:8081/",
		key:  "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
	}
}

func (e *emulatorTests) getClient(t *testing.T) *CosmosClient {
	cred, _ := NewSharedKeyCredential(e.key)
	client, err := NewCosmosClient(e.host, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func (e *emulatorTests) createDatabase(
	t *testing.T,
	ctx context.Context,
	client *CosmosClient,
	dbName string) *CosmosDatabase {
	database := CosmosDatabaseProperties{Id: dbName}
	resp, err := client.AddDatabase(ctx, database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.RawResponse.StatusCode != 201 {
		t.Fatal(e.parseErrorResponse(resp.RawResponse))
	}

	if resp.DatabaseProperties.Id != database.Id {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	return resp.DatabaseProperties.Database
}

func (e *emulatorTests) deleteDatabase(
	t *testing.T,
	ctx context.Context,
	database *CosmosDatabase) {
	resp, err := database.Delete(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}

	if resp.RawResponse.StatusCode != 204 {
		t.Fatal(e.parseErrorResponse(resp.RawResponse))
	}
}

func (e *emulatorTests) parseErrorResponse(response *http.Response) error {
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	return fmt.Errorf("Failed request with %v. \n Body %v", response, bodyString)
}
