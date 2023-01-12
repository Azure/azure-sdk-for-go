//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ExampleAppendBlobClient shows how to append data (in blocks) to an append blob.
// An append blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of an append blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func Example_appendblob_Client() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_append_blob.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	appendBlobClient, err := appendblob.NewClient(blobURL, cred, nil)
	handleError(err)

	_, err = appendBlobClient.Create(context.TODO(), nil)
	handleError(err)

	for i := 0; i < 5; i++ { // Append 5 blocks to the append blob
		_, err := appendBlobClient.AppendBlock(context.TODO(), streaming.NopCloser(strings.NewReader(fmt.Sprintf("Appending block #%d\n", i))), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Download the entire append blob's contents and read into a bytes.Buffer.
	get, err := appendBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	b := bytes.Buffer{}
	reader := get.Body
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(b.String())
}

// This example shows how to set an expiry time on an existing blob
// This operation is only allowed on Hierarchical Namespace enabled accounts.
func Example_appendblob_SetExpiry() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_append_blob_set_expiry.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	appendBlobClient, err := appendblob.NewClient(blobURL, cred, nil)
	handleError(err)

	// set expiry on append blob to an absolute time
	expiryTimeAbsolute := time.Now().Add(8 * time.Hour)
	_, err = appendBlobClient.SetExpiry(context.TODO(), appendblob.ExpiryTypeAbsolute(expiryTimeAbsolute), nil)
	handleError(err)

	// validate set expiry operation
	resp, err := appendBlobClient.GetProperties(context.TODO(), nil)
	handleError(err)
	if resp.ExpiresOn == nil || expiryTimeAbsolute.UTC().Format(http.TimeFormat) != (*resp.ExpiresOn).UTC().Format(http.TimeFormat) {
		return
	}
}

func Example_appendblob_Seal() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_append_blob_seal.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	appendBlobClient, err := appendblob.NewClient(blobURL, cred, nil)
	handleError(err)

	_, err = appendBlobClient.Seal(context.Background(), nil)
	handleError(err)
}
