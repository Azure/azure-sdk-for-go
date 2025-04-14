// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package benchmarks

import (
	"log"
	"os"
	"testing"

	"stress/internal/servicebus/shared"
)

func TestMain(m *testing.M) {
	if os.Getenv("ENV_FILE") == "" {
		os.Setenv("ENV_FILE", "../../../../.env")
	}

	err := shared.LoadEnvironment()

	if err != nil {
		log.Printf("Failed to load env file, benchmarks will not run: %s", err)
		return
	}

	os.Exit(m.Run())
}
