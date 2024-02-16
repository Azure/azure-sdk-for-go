//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azsystemevents_test

import (
	"azsystemevents"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testVars = struct {
	BlobURL          string
	QueueURL         string
	ConnectionString string
	QueueName        string
	SkipReason       string
}{}

func TestMain(m *testing.M) {
	var missingVars []string

	getVar := func(name string) string {
		v := os.Getenv(name)

		if v == "" {
			missingVars = append(missingVars, name)
		}

		return v
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Failed to load .env file: %s", err)
	}

	// os.Setenv("AZURE_CLIENT_ID", getVar("AZSYSTEMEVENTS_CLIENT_ID"))
	// os.Setenv("AZURE_CLIENT_SECRET", getVar("AZSYSTEMEVENTS_CLIENT_SECRET"))
	// os.Setenv("AZURE_TENANT_ID", getVar("AZSYSTEMEVENTS_TENANT_ID"))

	testVars.BlobURL = getVar("STORAGE_ACCOUNT_BLOB")
	testVars.QueueURL = getVar("STORAGE_ACCOUNT_QUEUE")
	testVars.QueueName = getVar("STORAGE_QUEUE_NAME")

	if len(missingVars) > 0 {
		testVars.SkipReason = fmt.Sprintf("WARNING: integration tests disabled, environment variables missing (%s)", strings.Join(missingVars, ","))
	}

	os.Exit(m.Run())
}

func parseManyEvents(t *testing.T, str string) []azsystemevents.Event {
	var events []azsystemevents.Event

	err := json.Unmarshal(([]byte)(str), &events)
	require.NoError(t, err)

	return events
}

func parseEvent(t *testing.T, str string) azsystemevents.Event {
	var event *azsystemevents.Event

	err := json.Unmarshal(([]byte)(str), &event)
	require.NoError(t, err)

	return *event
}

func parseManyCloudEvents(t *testing.T, str string) []messaging.CloudEvent {
	var events []messaging.CloudEvent

	err := json.Unmarshal(([]byte)(str), &events)
	require.NoError(t, err)

	return events
}

func parseCloudEvent(t *testing.T, str string) messaging.CloudEvent {
	var event *messaging.CloudEvent

	err := json.Unmarshal(([]byte)(str), &event)
	require.NoError(t, err)

	return *event
}

func deserializeSystemEvent[T any](t *testing.T, payload any) T {
	var val *T

	err := json.Unmarshal(payload.([]byte), &val)
	require.NoError(t, err)

	return *val
}

func mustParseTime(t *testing.T, str string) time.Time {
	tm, err := time.Parse(time.RFC3339Nano, str)
	require.NoError(t, err)

	return tm
}
