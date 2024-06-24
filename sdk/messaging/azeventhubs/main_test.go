// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	ec := testMain(m)
	os.Exit(ec)
}

func testMain(m *testing.M) int {
	if recording.GetRecordMode() == recording.LiveMode {
		// if one of our vars isn't defined just assume that they're supposed to come from an .env file
		if os.Getenv("CHECKPOINTSTORE_STORAGE_ENDPOINT") == "" {
			if err := godotenv.Load(); err != nil {
				log.Fatalf("Failed to load .env file when running live tests: %s", err.Error())
			}
		}

		// create a test storage container so our examples can run.
		defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			log.Fatalf("Failed to create DAC: %s", err)
		}

		blobClient, err := azblob.NewClient(os.Getenv("CHECKPOINTSTORE_STORAGE_ENDPOINT"), defaultAzureCred, nil)

		if err != nil {
			log.Fatalf("Failed to create blob client: %s", err)
		}

		containerName := fmt.Sprintf("container-%x", time.Now().UnixNano())

		defer func() { _, _ = blobClient.ServiceClient().DeleteContainer(context.Background(), containerName, nil) }()

		if _, err := blobClient.ServiceClient().CreateContainer(context.Background(), containerName, nil); err != nil {
			log.Fatalf("Failed to create blob container for examples: %s", err)
		}

		log.Printf("Example checkpoint container created: %s", containerName)
		os.Setenv("CHECKPOINTSTORE_STORAGE_CONTAINER_NAME", containerName)
	}

	return m.Run()
}
