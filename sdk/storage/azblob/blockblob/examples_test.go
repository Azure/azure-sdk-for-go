//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ExampleBlockBlobClient shows how to upload data (in blocks) to a blob.
// A block blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of a block blob is slightly more than 190 TiB (4000 MiB X 50,000 blocks).
func Example_blockblob_Client() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_block_blob.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	blockBlobClient, err := blockblob.NewClient(blobURL, cred, nil)
	handleError(err)

	// NOTE: The blockID must be <= 64 bytes and ALL blockIDs for the block must be the same length
	blockIDBinaryToBase64 := func(blockID []byte) string { return base64.StdEncoding.EncodeToString(blockID) }
	blockIDBase64ToBinary := func(blockID string) []byte { _binary, _ := base64.StdEncoding.DecodeString(blockID); return _binary }

	// These helper functions convert an int block ID to a base-64 string and vice versa
	blockIDIntToBase64 := func(blockID int) string {
		binaryBlockID := &[4]byte{} // All block IDs are 4 bytes long
		binary.LittleEndian.PutUint32(binaryBlockID[:], uint32(blockID))
		return blockIDBinaryToBase64(binaryBlockID[:])
	}
	blockIDBase64ToInt := func(blockID string) int {
		blockIDBase64ToBinary(blockID)
		return int(binary.LittleEndian.Uint32(blockIDBase64ToBinary(blockID)))
	}

	// Upload 4 blocks to the blob (these blocks are tiny; they can be up to 100MB each)
	words := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(words)) // The collection of block IDs (base 64 strings)

	// Upload each block sequentially (one after the other)
	for index, word := range words {
		// This example uses the index as the block ID; convert the index/ID into a base-64 encoded string as required by the service.
		// NOTE: Over the lifetime of a blob, all block IDs (before base 64 encoding) must be the same length (this example uses 4 byte block IDs).
		base64BlockIDs[index] = blockIDIntToBase64(index)

		// Upload a block to this blob specifying the Block ID and its content (up to 100MB); this block is uncommitted.
		_, err := blockBlobClient.StageBlock(context.TODO(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// After all the blocks are uploaded, atomically commit them to the blob.
	_, err = blockBlobClient.CommitBlockList(context.TODO(), base64BlockIDs, nil)
	handleError(err)

	// For the blob, show each block (ID and size) that is a committed part of it.
	getBlock, err := blockBlobClient.GetBlockList(context.TODO(), blockblob.BlockListTypeAll, nil)
	handleError(err)
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fmt.Printf("Block ID=%d, Size=%d\n", blockIDBase64ToInt(*block.Name), block.Size)
	}

	// Download the blob in its entirety; download operations do not take blocks into account.
	blobDownloadResponse, err := blockBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)

	blobData := &bytes.Buffer{}
	reader := blobDownloadResponse.Body
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(blobData)
}

// This example shows how to copy a large stream in blocks (chunks) to a block blob.
func Example_blockblob_Client_UploadFile() {
	file, err := os.Open("BigFile.bin") // Open the file we want to upload
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	fileSize, err := file.Stat() // Get the size of the file (stream)
	if err != nil {
		log.Fatal(err)
	}

	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blockBlobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := blockBlobClient.UploadFile(context.TODO(), file,
		&blockblob.UploadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Printf("Uploaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	_ = response // Avoid compiler's "declared and not used" error

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(destFile *os.File) {
		_ = destFile.Close()

	}(destFile)

	// Perform download
	_, err = blockBlobClient.DownloadFile(context.TODO(), destFile,
		&blob.DownloadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			}})

	if err != nil {
		log.Fatal(err)
	}
}

// This example shows how to set an expiry time on an existing blob
// This operation is only allowed on Hierarchical Namespace enabled accounts.
func Example_blockblob_SetExpiry() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_block_blob_set_expiry.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	blockBlobClient, err := blockblob.NewClient(blobURL, cred, nil)
	handleError(err)

	// set expiry on block blob 4 hours relative to now
	_, err = blockBlobClient.SetExpiry(context.TODO(), blockblob.ExpiryTypeRelativeToNow(4*time.Hour), nil)
	handleError(err)

	// validate set expiry operation
	resp, err := blockBlobClient.GetProperties(context.TODO(), nil)
	handleError(err)
	if resp.ExpiresOn == nil {
		return
	}
}

// This example shows how to set up log callback to dump SDK events
func Example_blockblob_uploadLogs() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	blobName := "test_block_blob_set_expiry.txt"
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/testcontainer/%s", accountName, blobName)

	azlog.SetEvents(azblob.EventUpload, azlog.EventRequest, azlog.EventResponse)
	azlog.SetListener(func(cls azlog.Event, msg string) {
		if cls == azblob.EventUpload {
			fmt.Println(msg)
		}
	})

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	blockBlobClient, err := blockblob.NewClient(blobURL, cred, nil)
	handleError(err)

	// set expiry on block blob 4 hours relative to now
	_, err = blockBlobClient.SetExpiry(context.TODO(), blockblob.ExpiryTypeRelativeToNow(4*time.Hour), nil)
	handleError(err)

	// validate set expiry operation
	resp, err := blockBlobClient.GetProperties(context.TODO(), nil)
	handleError(err)
	if resp.ExpiresOn == nil {
		return
	}
}
