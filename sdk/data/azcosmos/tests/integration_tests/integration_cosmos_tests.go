// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type integrationTests struct {
	connectionString string
}

var accountTypes = []struct {
	multipleWriteRegions bool
}{
	{multipleWriteRegions: false},
	{multipleWriteRegions: true},
}

func newIntegrationTests(t *testing.T, multipleWriteRegions bool) *integrationTests {
	return newIntegrationTests(t, multipleWriteRegions)
}

func newIntegrationTestsWithEndpoint(t *testing.T, multipleWriteRegions bool) *integrationTests {
	if multipleWriteRegions {
		envCheck := os.Getenv("MULTI_WRITE_CONNECTION_STRING")
	} else {
		envCheck := os.Getenv("SINGLE_WRITE_CONNECTION_STRING")
	}
	if envCheck == "" {
		t.Skip("set MULTI_WRITE_CONNECTION_STRING or SINGLE_WRITE_CONNECTION_STRING environment variable to run this test")
	}

	return &integrationTests{
		connectionString: envCheck,
	}
}

func (e *integrationTests) getClient(o *ClientOptions, t *testing.T, tp tracing.Provider) *Client {
	cred, err := azcore.NewDefaultAzureCredential(nil)
	handle(err)

	if err != nil {
		t.Fatalf("Failed to create credential: %v", err)
	}

	client, err = azcosmos.NewClientFromConnectionString(e.connectionString, o)
	handle(err)

	return client
}


func (e *integrationTests) createDatabase(
	t *testing.T,
	ctx context.Context,
	client *Client,
	dbName string) *DatabaseClient {
	database := DatabaseProperties{ID: dbName}
	resp, err := client.CreateDatabase(ctx, database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	db, _ := client.NewDatabase(dbName)
	return db
}

func (e *integrationTests) deleteDatabase(
	t *testing.T,
	ctx context.Context,
	database *DatabaseClient) {
	_, err := database.Delete(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}
