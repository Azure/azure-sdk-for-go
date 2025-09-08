//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
)

// recordingDirectory should point to the testdata directory for your package.
// When copying this template, update this path to match your package location.
// For example: "sdk/data/azappconfig/testdata" or "sdk/storage/azblob/testdata"
const recordingDirectory = "sdk/template/aztemplate/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() != recording.LiveMode {
		// NOTE: For the template package, no assets.json exists since this is a template.
		// When copying this template for your own package, you'll need to:
		// 1. Create an assets.json file in your package directory
		// 2. Set up proper test recordings using the test-proxy
		// 3. Update the recordingDirectory constant above to match your package path
		// 
		// For now, if proxy startup fails, we'll fall back to live mode gracefully.
		// In a real package with proper recordings, you would want proxy startup failures to be fatal.
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			// Template package doesn't have recordings, so proxy startup will fail.
			// This graceful fallback allows the template to work without recordings.
			// Remove this fallback when implementing a real package with proper recordings.
			return m.Run()
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	// Load environment variables from .env file for testing.
	// This is useful for local development where you can store connection strings,
	// endpoints, and other test configuration in a .env file.
	// When copying this template, create a .env file in your package root with your test variables.
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Failed to load .env file, no integration tests will run: %s", err)
	}

	return m.Run()
}