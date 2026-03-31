// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const (
	pathToPackage = "sdk/resourcemanager/storage/armstorage/testdata"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	f := testutil.StartProxy(pathToPackage)
	defer f()

	// Storage account keys can appear as base64 in key lists and other payloads.
	if err := recording.AddBodyRegexSanitizer(`"value":"sanitized-storage-account-key"`, `"value"\s*:\s*"(?:[A-Za-z0-9+]|\\/){84,128}={0,2}"`, nil); err != nil {
		panic(err)
	}

	if err := recording.AddBodyKeySanitizer(`$..keys[*].value`, "sanitized-storage-account-key", "", nil); err != nil {
		panic(err)
	}

	regexOptions := &recording.RecordingOptions{UseHTTPS: true, GroupForReplace: "1"}
	if err := recording.AddGeneralRegexSanitizer("sanitized-storage-account-key", `"value"\s*:\s*"((?:[A-Za-z0-9+]|\\/){84,128}={0,2})"`, regexOptions); err != nil {
		panic(err)
	}

	return m.Run()
}
