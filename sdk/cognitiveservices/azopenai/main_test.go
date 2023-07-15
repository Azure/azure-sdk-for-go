// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("No .env file - can't run examples or live tests\n")
		}
	}

	os.Exit(m.Run())
}
