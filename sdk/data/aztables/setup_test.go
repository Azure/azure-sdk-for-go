//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package aztables

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func TestMain(m *testing.M) {
	// 1. Set up session level sanitizers
	if recording.GetRecordMode() == "record" {
		for _, val := range []string{"TABLES_COSMOS_ACCOUNT_NAME", "TABLES_STORAGE_ACCOUNT_NAME"} {
			account := os.Getenv(val)
			err := recording.AddUriSanitizer("fakeaccount", account, nil)
			if err != nil {
				panic(err)
			}
		}
	}

	// Run tests
	exitVal := m.Run()

	// 3. Reset
	// TODO: Add after sanitizer PR
	// err = recording.ResetSanitizers(nil)

	// 4. Error out if applicable
	os.Exit(exitVal)
}
