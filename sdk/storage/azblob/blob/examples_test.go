//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// This examples shows how to create a blob with HTTP Headers, how to read, and how to update the blob's HTTP headers.
func Example_blob_HTTPHeaders() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := blob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blob with HTTP headers
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&blockblob.UploadOptions{HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType:        to.Ptr("text/html; charset=utf-8"),
			BlobContentDisposition: to.Ptr("attachment"),
		}},
	)
	handleError(err)

	// GetMetadata returns the blob's properties, HTTP headers, and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	handleError(err)

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Shows some of the blob's HTTP Headers
	httpHeaders := blob.ParseHTTPHeaders(get)
	fmt.Println(httpHeaders.BlobContentType, httpHeaders.BlobContentDisposition)

	// Update the blob's HTTP Headers and write them back to the blob
	httpHeaders.BlobContentType = to.Ptr("text/plain")
	_, err = blobClient.SetHTTPHeaders(context.TODO(), httpHeaders, nil)
	handleError(err)
}

// This example shows how to create a blob with metadata, read blob metadata, and update a blob's read-only properties and metadata.
func Example_blob_Client_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := blob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blockblob.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blob with metadata (string key/value pairs)
	// Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	handleError(err)
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&blockblob.UploadOptions{Metadata: map[string]*string{"author": to.Ptr("Jeffrey"), "app": to.Ptr(creatingApp)}},
	)
	handleError(err)

	// Query the blob's properties and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	handleError(err)

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Show the blob's metadata
	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}

	for k, v := range get.Metadata {
		fmt.Print(k + "=" + *v + "\n")
	}

	// Update the blob's metadata and write it back to the blob
	get.Metadata["editor"] = to.Ptr("Grant")
	_, err = blobClient.SetMetadata(context.TODO(), get.Metadata, nil)
	handleError(err)
}

// This example show how to create a blob, take a snapshot of it, update the base blob,
// read from the blob snapshot, list blobs with their snapshots, and delete blob snapshots.
func Example_blob_Client_CreateSnapshot() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	credential, err := blob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	containerClient, err := container.NewClientWithSharedKeyCredential(u, credential, nil)
	handleError(err)

	// Create a blockBlobClient object to a blob in the container.
	baseBlobClient := containerClient.NewBlockBlobClient("Original.txt")

	// Create the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(streaming.NopCloser(strings.NewReader("Some text"))), nil)
	handleError(err)

	// Create a snapshot of the original blob & save its timestamp:
	createSnapshot, err := baseBlobClient.CreateSnapshot(context.TODO(), nil)
	handleError(err)
	snapshot := *createSnapshot.Snapshot

	// Modify the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("New text")), nil)
	handleError(err)

	// Download the modified blob:
	get, err := baseBlobClient.DownloadStream(context.TODO(), nil)
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

	// Show snapshot blob via original blob URI & snapshot time:
	snapshotBlobClient, _ := baseBlobClient.WithSnapshot(snapshot)
	get, err = snapshotBlobClient.DownloadStream(context.TODO(), nil)
	handleError(err)
	b.Reset()
	reader = get.Body
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(b.String())

	// FYI: You can get the base blob URL from one of its snapshot by passing "" to WithSnapshot:
	baseBlobClient, _ = snapshotBlobClient.WithSnapshot("")

	// Show all blobs in the container with their snapshots:
	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	pager := containerClient.NewListBlobsFlatPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.Segment.BlobItems {
			// Process the blobs returned
			snapTime := "N/A"
			if blob.Snapshot != nil {
				snapTime = *blob.Snapshot
			}
			fmt.Printf("Blob name: %s, Snapshot: %s\n", *blob.Name, snapTime)
		}
	}

	// Promote read-only snapshot to writable base blob:
	_, err = baseBlobClient.StartCopyFromURL(context.TODO(), snapshotBlobClient.URL(), nil)
	handleError(err)

	// When calling Delete on a base blob:
	// DeleteSnapshotsOptionOnly deletes all the base blob's snapshots but not the base blob itself
	// DeleteSnapshotsOptionInclude deletes the base blob & all its snapshots.
	// DeleteSnapshotOptionNone produces an error if the base blob has any snapshots.
	_, err = baseBlobClient.Delete(context.TODO(), &blob.DeleteOptions{DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeInclude)})
	handleError(err)
}

// This example shows how to copy a source document on the Internet to a blob.
func Example_blob_Client_StartCopyFromURL() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a containerClient object to a container where we'll create a blob and its snapshot.
	// Create a blockBlobClient object to a blob in the container.
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/CopiedBlob.bin", accountName)
	credential, err := blob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	blobClient, err := blob.NewClientWithSharedKeyCredential(blobURL, credential, nil)
	handleError(err)

	src := "https://cdn2.auth0.com/docs/media/addons/azure_blob.svg"
	startCopy, err := blobClient.StartCopyFromURL(context.TODO(), src, nil)
	handleError(err)

	copyID := *startCopy.CopyID
	copyStatus := *startCopy.CopyStatus
	for copyStatus == blob.CopyStatusTypePending {
		time.Sleep(time.Second * 2)
		getMetadata, err := blobClient.GetProperties(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		copyStatus = *getMetadata.CopyStatus
	}
	fmt.Printf("Copy from %s to %s: ID=%s, Status=%s\n", src, blobClient.URL(), copyID, copyStatus)
}

// This example demonstrates splitting a URL into its parts so you can examine and modify the URL in an Azure Storage fluent way.
func ExampleParseURL() {
	// Here is an example of a blob snapshot.
	u := "https://myaccount.blob.core.windows.net/mycontainter/ReadMe.txt?" +
		"snapshot=2011-03-09T01:42:34Z&" +
		"sv=2015-02-21&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&" +
		"spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=92836758923659283652983562=="

	// Breaking the URL down into it's parts by conversion to URLParts
	parts, _ := blob.ParseURL(u)

	// The URLParts allows access to individual portions of a Blob URL
	fmt.Printf("Host: %s\nContainerName: %s\nBlobName: %s\nSnapshot: %s\n", parts.Host, parts.ContainerName, parts.BlobName, parts.Snapshot)
	fmt.Printf("Version: %s\nResource: %s\nStartTime: %s\nExpiryTime: %s\nPermissions: %s\n", parts.SAS.Version(), parts.SAS.Resource(), parts.SAS.StartTime(), parts.SAS.ExpiryTime(), parts.SAS.Permissions())
}
