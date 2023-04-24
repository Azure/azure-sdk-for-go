//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingestion_test

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingestion"
)

func TestUpload(t *testing.T) {
	endpoint := os.Getenv("LOGS_INGESTION_ENDPOINT")
	ruleID := os.Getenv("DATA_COLLECTION_RULE_ID")
	var streamName = "<STREAM_NAME>"

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	client, err := azingestion.NewClient(endpoint, credential, nil)
	if err != nil {
		panic(err)
	}
	letters := []any{"a", "b", "c", "d"}
	_, err = client.Upload(context.Background(), ruleID, streamName, letters, nil)
	if err != nil {
		panic(err)
	}

}
