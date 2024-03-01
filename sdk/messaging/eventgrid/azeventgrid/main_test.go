//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
)

const recordingDirectory = "sdk/messaging/eventgrid/azeventgrid/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

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

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Failed to load .env file, no integration tests will run: %s", err)
	}

	return m.Run()
}
