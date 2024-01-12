//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package systemevents_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/systemevents"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testVars = struct {
	BlobURL    string
	QueueURL   string
	QueueName  string
	SkipReason string
}{}

func TestMain(m *testing.M) {
	var missingVars []string

	getVar := func(name string) string {
		v := os.Getenv(name)

		if v == "" {
			missingVars = append(missingVars, name)
		}

		return v
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Failed to load .env file: %s", err)
	}

	os.Setenv("AZURE_CLIENT_ID", getVar("AZEVENTGRID_CLIENT_ID"))
	os.Setenv("AZURE_CLIENT_SECRET", getVar("AZEVENTGRID_CLIENT_SECRET"))
	os.Setenv("AZURE_TENANT_ID", getVar("AZEVENTGRID_TENANT_ID"))

	testVars.BlobURL = getVar("STORAGE_ACCOUNT_BLOB")
	testVars.QueueURL = getVar("STORAGE_ACCOUNT_QUEUE")
	testVars.QueueName = getVar("STORAGE_QUEUE_NAME")

	if len(missingVars) > 0 {
		testVars.SkipReason = fmt.Sprintf("WARNING: integration tests disabled, environment variables missing (%s)", strings.Join(missingVars, ","))
	}

	os.Exit(m.Run())
}

func TestSystemEventForBlobs(t *testing.T) {
	if testVars.SkipReason != "" {
		t.Skipf(testVars.SkipReason)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	blobClient, err := azblob.NewClient(testVars.BlobURL, cred, nil)
	require.NoError(t, err)

	containerName := strings.ToLower(fmt.Sprintf("container%X", time.Now().UnixNano()))
	_, err = blobClient.CreateContainer(context.Background(), containerName, nil)
	require.NoError(t, err)

	// we have a system topic setup so blob uploads generate events to a queue on the same storage account.
	testBlobName := "testblob"
	_, err = blobClient.UploadBuffer(context.Background(), containerName, testBlobName, []byte{1}, nil)
	require.NoError(t, err)

	_, err = blobClient.DeleteBlob(context.Background(), containerName, testBlobName, nil)
	require.NoError(t, err)

	queueClient, err := azqueue.NewQueueClient(testVars.QueueURL+testVars.QueueName, cred, nil)
	require.NoError(t, err)

	_, err = queueClient.ClearMessages(context.Background(), nil)
	require.NoError(t, err)

	var messages []*azqueue.DequeuedMessage

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for len(messages) < 2 {
		getMessageResp, err := queueClient.DequeueMessage(ctx, nil)
		require.NoError(t, err)

		if len(getMessageResp.Messages) > 0 {
			messages = append(messages, getMessageResp.Messages...)
		}

		for _, msg := range getMessageResp.Messages {
			_, err := queueClient.DeleteMessage(ctx, *msg.MessageID, *msg.PopReceipt, nil)
			require.NoError(t, err)
		}
	}

	// there are two events that we expect since we uploaded and then deleted a blob
	for _, msg := range messages {
		decodedData, err := base64.StdEncoding.DecodeString(*msg.MessageText)
		require.NoError(t, err)

		var ce *messaging.CloudEvent

		err = json.Unmarshal(decodedData, &ce)
		require.NoError(t, err)

		require.Equal(t, fmt.Sprintf("/blobServices/default/containers/%s/blobs/%s", containerName, testBlobName), *ce.Subject)

		switch ce.Type {
		case string(systemevents.TypeStorageBlobCreatedEventData):
			var blobCreated *systemevents.StorageBlobCreatedEventData
			err = json.Unmarshal(ce.Data.([]byte), &blobCreated)
			require.NoError(t, err)

			require.Equal(t, "PutBlob", *blobCreated.API)
			require.Equal(t, "BlockBlob", *blobCreated.BlobType)

		case string(systemevents.TypeStorageBlobDeletedEventData):
			var blobDeleted systemevents.StorageBlobDeletedEventData
			err = json.Unmarshal(ce.Data.([]byte), &blobDeleted)
			require.NoError(t, err)

			require.Equal(t, "DeleteBlob", *blobDeleted.API)
			require.Equal(t, "BlockBlob", *blobDeleted.BlobType)
		}
	}
}
