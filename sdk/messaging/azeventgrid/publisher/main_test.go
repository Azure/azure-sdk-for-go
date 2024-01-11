//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package publisher_test

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
)

const recordingDirectory = "sdk/messaging/azeventgrid/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

var testVars = struct {
	SkipReason string
	BlobURL    string
	QueueURL   string
	QueueName  string
	eventGridVars
}{}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Failed to load .env file: %s", err)
	}

	var missingVars []string

	getVar := func(name string) string {
		v := os.Getenv(name)

		if v == "" {
			missingVars = append(missingVars, name)
		}

		return v
	}

	os.Setenv("AZURE_CLIENT_ID", getVar("AZEVENTGRID_CLIENT_ID"))
	os.Setenv("AZURE_CLIENT_SECRET", getVar("AZEVENTGRID_CLIENT_SECRET"))
	os.Setenv("AZURE_TENANT_ID", getVar("AZEVENTGRID_TENANT_ID"))

	testVars.BlobURL = getVar("STORAGE_ACCOUNT_BLOB")
	testVars.QueueURL = getVar("STORAGE_ACCOUNT_QUEUE")
	testVars.QueueName = getVar("STORAGE_QUEUE_NAME")

	testVars.eventGridVars = eventGridVars{
		EG: topicVars{Name: getVar("EVENTGRID_TOPIC_NAME"),
			Key:      getVar("EVENTGRID_TOPIC_KEY"),
			Endpoint: getVar("EVENTGRID_TOPIC_ENDPOINT"),
		},
		CE: topicVars{Name: getVar("EVENTGRID_CE_TOPIC_NAME"),
			Key:      getVar("EVENTGRID_CE_TOPIC_KEY"),
			Endpoint: getVar("EVENTGRID_CE_TOPIC_ENDPOINT"),
		},
	}

	sort.Strings(missingVars)

	if len(missingVars) > 0 {
		testVars.SkipReason = fmt.Sprintf("WARNING: integration tests disabled, environment variables missing (%s)", strings.Join(missingVars, ","))
	}

	return m.Run()
}
