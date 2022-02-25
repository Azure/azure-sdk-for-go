// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Failed to load env file, NO LIVE TESTS WILL RUN: %s", err.Error())
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
