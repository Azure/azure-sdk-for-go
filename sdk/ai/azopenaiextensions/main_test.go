// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const RecordingDirectory = "sdk/ai/azopenaiextensions/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func printProxyLogs() {
	filePattern := "/tmp/test-proxy.log.*"

	// Find matching files
	matches, err := filepath.Glob(filePattern)
	if err != nil {
		fmt.Printf("Error finding matching files: %v\n", err)
		os.Exit(1)
	}

	// Check if any files matched
	if len(matches) == 0 {
		fmt.Printf("No files found matching the pattern: %s\n", filePattern)
		os.Exit(1)
	}

	// Print the files found (optional)
	fmt.Printf("Found %d matching files:\n", len(matches))
	for i, match := range matches {
		fmt.Printf("%d. %s\n", i+1, match)
	}

	// Read and print contents of each file
	for _, filePath := range matches {
		fmt.Printf("\n--- Contents of %s ---\n", filePath)

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			continue // Skip to the next file if this one fails
		}

		// Ensure the file is closed when we're done
		defer file.Close()

		// Read the file contents
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		// Print the file contents
		fmt.Println(string(content))
	}
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		defaultOptions := getRecordingOptions(nil)
		proxy, err := recording.StartTestProxy(RecordingDirectory, defaultOptions)
		if err != nil {
			panic(err)
		}

		err = recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
			RecordingOptions: *defaultOptions,
			ExcludedHeaders: []string{
				"X-Stainless-Arch",
				"X-Stainless-Lang",
				"X-Stainless-Os",
				"X-Stainless-Package-Version",
				"X-Stainless-Retry-Count",
				"X-Stainless-Runtime",
				"X-Stainless-Runtime-Version",
			},
		})

		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
			printProxyLogs()
		}()
	}
	os.Setenv("AOAI_OYD_ENDPOINT", os.Getenv("AOAI_ENDPOINT_USEAST"))
	os.Setenv("AOAI_OYD_MODEL", "gpt-4-0613")

	return m.Run()
}
