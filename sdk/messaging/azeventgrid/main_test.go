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

func TestMain(m *testing.M) {
	if err := recording.ResetProxy(nil); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file, no integration tests will run: %s", err)
	}

	os.Exit(m.Run())
}
