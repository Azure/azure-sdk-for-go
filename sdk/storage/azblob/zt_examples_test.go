//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// This example shows you how to get started using the Azure Blob Storage SDK for Go.
func Example() {
	// Your account name and key can be obtained from the Azure Portal.
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	service, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// ===== 1. Create a container =====

	// First, create a container client, and use the Create method to create a new container in your account
	container, _ := service.NewContainerClient("mycontainer")
	// All functions that make service requests have an options struct as the final parameter.
	// The options struct allows you to specify optional parameters such as metadata, public access types, etc.
	// If you want to use the default options, pass in nil.
	_, err = container.Create(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// ===== 2. Upload and Download a block blob =====
	data := "Hello world!"

	// Create a new BlockBlobClient from the ContainerClient
	blockBlob, _ := container.NewBlockBlobClient("HelloWorld.txt")

	// Upload data to the block blob
	_, err = blockBlob.Upload(context.TODO(), streaming.NopCloser(strings.NewReader(data)), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Download the blob's contents and ensure that the download worked properly
	get, err := blockBlob.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Use the bytes.Buffer object to read the downloaded data.
	downloadedData := &bytes.Buffer{}
	reader := get.Body(nil) // RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return
	}

	err = reader.Close()
	if err != nil {
		return
	}
	if data != downloadedData.String() {
		log.Fatal("downloaded data does not match uploaded data")
	} else {
		fmt.Printf("Downloaded data: %s\n", downloadedData.String())
	}

	// ===== 3. List blobs =====
	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
	// PageResponse() can be used to iterate over the results of the specific page.
	// Always check the Err() method after paging to see if an error was returned by the pager. A pager will return either an error or the page of results.
	pager := container.ListBlobsFlat(nil)
	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()

		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			fmt.Println(*v.Name)
		}
	}

	if err = pager.Err(); err != nil {
		log.Fatal(err)
	}

	// Delete the blob.
	_, err = blockBlob.Delete(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the container.
	_, err = container.Delete(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleStorageError() {
	/* This example demonstrates how to handle errors returned from the various Client methods. All these methods return an
	   object implementing the azcore.Response interface and an object implementing Go's error interface.
	   The error result is nil if the request was successful; your code can safely use the Response interface object.
	   If the error is non-nil, the error could be due to:

	1. An invalid argument passed to the method. You should not write code to handle these errors;
	   instead, fix these errors as they appear during development/testing.

	2. A network request didn't reach an Azure Storage Service. This usually happens due to a bad URL or
	   faulty networking infrastructure (like a router issue). In this case, an object implementing the
	   net.Error interface will be returned. The net.Error interface offers Timeout and Temporary methods
	   which return true if the network error is determined to be a timeout or temporary condition. If
	   your pipeline uses the retry policy factory, then this policy looks for Timeout/Temporary and
	   automatically retries based on the retry options you've configured. Because of the retry policy,
	   your code will usually not call the Timeout/Temporary methods explicitly other than possibly logging
	   the network failure.

	3. A network request did reach the Azure Storage Service but the service failed to perform the
	   requested operation. In this case, an object implementing the StorageError interface is returned.
	   The StorageError interface also implements the net.Error interface and, if you use the retry policy,
	   you would most likely ignore the Timeout/Temporary methods. However, the StorageError interface exposes
	   richer information such as a service error code, an error description, details data, and the
	   service-returned http.Response. And, from the http.Response, you can get the initiating http.Request.
	*/

	container, err := azblob.NewContainerClientWithNoCredential("https://myaccount.blob.core.windows.net/mycontainer", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = container.Create(context.TODO(), nil)

	if err != nil {
		var stgErr azblob.StorageError

		if errors.As(err, &stgErr) { // We know this error is service-specific
			switch stgErr.ErrorCode {
			case azblob.StorageErrorCodeContainerAlreadyExists:
				// You can also look at the *http.Response that's attached to the error as well.
				if resp := stgErr.Response(); resp != nil {
					fmt.Println(resp.Request)
				}
			case azblob.StorageErrorCodeContainerBeingDeleted:
				// Handle this error ...
			default:
				// Handle other errors ...
			}
		}
	}
}

// This example demonstrates splitting a URL into its parts so you can examine and modify the URL in an Azure Storage fluent way.
func ExampleBlobURLParts() {
	// Here is an example of a blob snapshot.
	u := "https://myaccount.blob.core.windows.net/mycontainter/ReadMe.txt?" +
		"snapshot=2011-03-09T01:42:34Z&" +
		"sv=2015-02-21&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&" +
		"spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=92836758923659283652983562=="

	// Breaking the URL down into it's parts by conversion to BlobURLParts
	parts, _ := azblob.NewBlobURLParts(u)

	// The BlobURLParts allows access to individual portions of a Blob URL
	fmt.Printf("Host: %s\nContainerName: %s\nBlobName: %s\nSnapshot: %s\n", parts.Host, parts.ContainerName, parts.BlobName, parts.Snapshot)
	fmt.Printf("Version: %s\nResource: %s\nStartTime: %s\nExpiryTime: %s\nPermissions: %s\n", parts.SAS.Version(), parts.SAS.Resource(), parts.SAS.StartTime(), parts.SAS.ExpiryTime(), parts.SAS.Permissions())

	// You can alter fields to construct a new URL:
	// Note: SAS tokens may be limited to a specific container or blob, be careful modifying SAS tokens, you might take them outside of their original scope accidentally.
	parts.SAS = azblob.SASQueryParameters{}
	parts.Snapshot = ""
	parts.ContainerName = "othercontainer"

	// construct a new URL from the parts
	fmt.Print(parts.URL())
}

// This example demonstrates how to use the SAS token generators, which exist across all clients (Service, Container, Blob, and specialized Blob clients).
// This example specifically generates an account SAS token.
func ExampleServiceClient_GetSASToken() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	serviceClient, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Provide the relevant info (services, resource types, permissions, and duration)
	// The SAS token will be valid from this moment onwards.
	accountSAS, err := serviceClient.GetSASToken(
		azblob.AccountSASResourceTypes{Object: true, Service: true, Container: true},
		azblob.AccountSASPermissions{Read: true, List: true},
		azblob.AccountSASServices{Blob: true},
		time.Now(),
		time.Now().Add(48*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	// This URL can be used to authenticate requests now
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, accountSAS)
	serviceClientUsingSAS, err := azblob.NewServiceClientWithNoCredential(sasURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// For example, you can now list all the containers
	pager := serviceClientUsingSAS.ListContainers(nil)
	for pager.NextPage(context.TODO()) {
		for _, container := range pager.PageResponse().ContainerItems {
			fmt.Println(*container.Name)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	// You can also break a blob URL up into it's constituent parts
	blobURLParts, _ := azblob.NewBlobURLParts(serviceClientUsingSAS.URL())
	fmt.Printf("SAS expiry time = %s\n", blobURLParts.SAS.ExpiryTime())
}

// This example shows how to create and use an Azure Storage account Shared Access Signature (SAS).
func ExampleAccountSASSignatureValues_Sign() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	sasQueryParams, err := azblob.AccountSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   azblob.AccountSASPermissions{Read: true, List: true}.String(),
		Services:      azblob.AccountSASServices{Blob: true}.String(),
		ResourceTypes: azblob.AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.Sign(credential)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := sasQueryParams.Encode()
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, queryParams)

	// This URL can be used to authenticate requests now
	serviceClient, err := azblob.NewServiceClientWithNoCredential(sasURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// You can also break a blob URL up into it's constituent parts
	blobURLParts, _ := azblob.NewBlobURLParts(serviceClient.URL())
	fmt.Printf("SAS expiry time = %s\n", blobURLParts.SAS.ExpiryTime())
}

// This example demonstrates how to create and use a Blob service Shared Access Signature (SAS)
func ExampleBlobSASSignatureValues() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	containerName := "mycontainer"
	blobName := "HelloWorld.txt"

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		ContainerName: containerName,
		BlobName:      blobName,
		Permissions:   azblob.BlobSASPermissions{Add: true, Read: true, Write: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	// Create the SAS URL for the resource you wish to access, and append the SAS query parameters.
	qp := sasQueryParams.Encode()
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, containerName, blobName, qp)

	// Access the SAS-protected resource
	blob, err := azblob.NewBlobClientWithNoCredential(sasURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// if you have a SAS query parameter string, you can parse it into it's parts.
	blobURLParts, _ := azblob.NewBlobURLParts(blob.URL())
	fmt.Printf("SAS expiry time=%v", blobURLParts.SAS.ExpiryTime())
}

// This example shows how to manipulate a container's permissions.
func ExampleContainerClient_SetAccessPolicy() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	uri := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	container, err := azblob.NewContainerClientWithSharedKey(uri, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create the container
	_, err = container.Create(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Upload a simple blob.
	blob, _ := container.NewBlockBlobClient("HelloWorld.txt")

	_, err = blob.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("Hello World!")), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Attempt to read the blob
	get, err := http.Get(blob.URL())
	if err != nil {
		log.Fatal(err)
	}
	if get.StatusCode == http.StatusNotFound {
		// Change the blob to be a public access blob
		_, err := container.SetAccessPolicy(
			context.TODO(),
			&azblob.SetAccessPolicyOptions{
				ContainerSetAccessPolicyOptions: azblob.ContainerSetAccessPolicyOptions{
					Access: azblob.PublicAccessTypeBlob.ToPtr(),
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		// Now, this works:
		get, err = http.Get(blob.URL())
		if err != nil {
			log.Fatal(err)
		}
		var text bytes.Buffer
		_, err = text.ReadFrom(get.Body)
		if err != nil {
			return
		}
		defer get.Body.Close()

		fmt.Println("Public access blob data: ", text.String())
	}
}

// This example shows how to perform operations on blob conditionally.
func ExampleBlobAccessConditions() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blockBlob, err := azblob.NewBlockBlobClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/Data.txt", accountName), credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// This function displays the results of an operation
	showResult := func(response *azblob.DownloadResponse, err error) {
		if err != nil {
			var stgErr *azblob.StorageError
			if errors.As(err, &stgErr) {
				log.Fatalf("Failure: %s\n", stgErr.Error())
			} else {
				log.Fatal(err) // Network failure
			}
		} else {
			err := response.Body(nil).Close()
			if err != nil {
				log.Fatal(err)
			}
			// The client must close the response body when finished with it
			fmt.Printf("Success: %s\n", response.RawResponse.Status)
		}

		// Close the response
		if err != nil {
			return
		}
		fmt.Printf("Success: %s\n", response.RawResponse.Status)
	}

	showResultUpload := func(upload azblob.BlockBlobUploadResponse, err error) {
		if err != nil {
			var stgErr *azblob.StorageError
			if errors.As(err, &stgErr) {
				log.Fatalf("Failure: " + stgErr.Error() + "\n")
			} else {
				log.Fatal(err) // Network failure
			}
		}
		fmt.Print("Success: " + upload.RawResponse.Status + "\n")
	}

	// Create the blob
	upload, err := blockBlob.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("Text-1")), nil)
	showResultUpload(upload, err)

	// Download blob content if the blob has been modified since we uploaded it (fails):
	downloadResp, err := blockBlob.Download(
		context.TODO(),
		&azblob.DownloadBlobOptions{
			BlobAccessConditions: &azblob.BlobAccessConditions{
				ModifiedAccessConditions: &azblob.ModifiedAccessConditions{
					IfModifiedSince: upload.LastModified,
				},
			},
		},
	)
	showResult(&downloadResp, err)

	// Download blob content if the blob hasn't been modified in the last 24 hours (fails):
	downloadResp, err = blockBlob.Download(
		context.TODO(),
		&azblob.DownloadBlobOptions{
			BlobAccessConditions: &azblob.BlobAccessConditions{
				ModifiedAccessConditions: &azblob.ModifiedAccessConditions{
					IfUnmodifiedSince: to.Ptr(time.Now().UTC().Add(time.Hour * -24))},
			},
		},
	)
	showResult(&downloadResp, err)

	// Upload new content if the blob hasn't changed since the version identified by ETag (succeeds):
	showResultUpload(blockBlob.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Text-2")),
		&azblob.UploadBlockBlobOptions{
			BlobAccessConditions: &azblob.BlobAccessConditions{
				ModifiedAccessConditions: &azblob.ModifiedAccessConditions{IfMatch: upload.ETag},
			},
		},
	))

	// Download content if it has changed since the version identified by ETag (fails):
	downloadResp, err = blockBlob.Download(
		context.TODO(),
		&azblob.DownloadBlobOptions{
			BlobAccessConditions: &azblob.BlobAccessConditions{
				ModifiedAccessConditions: &azblob.ModifiedAccessConditions{IfNoneMatch: upload.ETag}},
		})
	showResult(&downloadResp, err)

	// Upload content if the blob doesn't already exist (fails):
	showResultUpload(blockBlob.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Text-3")),
		&azblob.UploadBlockBlobOptions{
			BlobAccessConditions: &azblob.BlobAccessConditions{
				ModifiedAccessConditions: &azblob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(string(azcore.ETagAny))},
			},
		}))
}

func ExampleContainerClient_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	containerClient, err := azblob.NewContainerClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName), credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a container with some metadata, key names are converted to lowercase before being sent to the service.
	// You should always use lowercase letters, especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = containerClient.Create(context.TODO(), &azblob.CreateContainerOptions{Metadata: map[string]string{"author": "Jeffrey", "app": creatingApp}})
	if err != nil {
		log.Fatal(err)
	}

	// Query the container's metadata
	get, err := containerClient.GetProperties(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	if get.Metadata == nil {
		log.Fatal("metadata is empty!")
	}

	for k, v := range get.Metadata {
		fmt.Printf("%s=%s\n", k, v)
	}

	// Update the metadata and write it back to the container
	get.Metadata["author"] = "Aidan"
	_, err = containerClient.SetMetadata(context.TODO(), &azblob.SetMetadataContainerOptions{Metadata: get.Metadata})
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: SetMetadata & SetProperties methods update the container's ETag & LastModified properties
}

// This example shows how to create a blob with metadata, read blob metadata, and update a blob's read-only properties and metadata.
func ExampleBlobClient_SetMetadata() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a blob with metadata (string key/value pairs)
	// Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&azblob.UploadBlockBlobOptions{Metadata: map[string]string{"author": "Jeffrey", "app": creatingApp}},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Query the blob's properties and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Show the blob's metadata
	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}

	for k, v := range get.Metadata {
		fmt.Print(k + "=" + v + "\n")
	}

	// Update the blob's metadata and write it back to the blob
	get.Metadata["editor"] = "Grant"
	_, err = blobClient.SetMetadata(context.TODO(), get.Metadata, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// This examples shows how to create a blob with HTTP Headers, how to read, and how to update the blob's HTTP headers.
func ExampleBlobHTTPHeaders() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a blob with HTTP headers
	_, err = blobClient.Upload(
		context.TODO(),
		streaming.NopCloser(strings.NewReader("Some text")),
		&azblob.UploadBlockBlobOptions{HTTPHeaders: &azblob.BlobHTTPHeaders{
			BlobContentType:        to.Ptr("text/html; charset=utf-8"),
			BlobContentDisposition: to.Ptr("attachment"),
		}},
	)
	if err != nil {
		log.Fatal(err)
	}

	// GetMetadata returns the blob's properties, HTTP headers, and metadata
	get, err := blobClient.GetProperties(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Show some of the blob's read-only properties
	fmt.Printf("BlobType: %s\nETag: %s\nLastModified: %s\n", *get.BlobType, *get.ETag, *get.LastModified)

	// Shows some of the blob's HTTP Headers
	httpHeaders := get.GetHTTPHeaders()
	fmt.Println(httpHeaders.BlobContentType, httpHeaders.BlobContentDisposition)

	// Update the blob's HTTP Headers and write them back to the blob
	httpHeaders.BlobContentType = to.Ptr("text/plain")
	_, err = blobClient.SetHTTPHeaders(context.TODO(), httpHeaders, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleBlockBlobClient shows how to upload data (in blocks) to a blob.
// A block blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of a block blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func ExampleBlockBlobClient() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a ContainerClient object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyBlockBlob.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: The blockID must be <= 64 bytes and ALL blockIDs for the block must be the same length
	blockIDBinaryToBase64 := func(blockID []byte) string { return base64.StdEncoding.EncodeToString(blockID) }
	blockIDBase64ToBinary := func(blockID string) []byte { _binary, _ := base64.StdEncoding.DecodeString(blockID); return _binary }

	// These helper functions convert an int block ID to a base-64 string and vice versa
	blockIDIntToBase64 := func(blockID int) string {
		binaryBlockID := (&[4]byte{}) // All block IDs are 4 bytes long
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
		_, err := blobClient.StageBlock(context.TODO(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// After all the blocks are uploaded, atomically commit them to the blob.
	_, err = blobClient.CommitBlockList(context.TODO(), base64BlockIDs, nil)
	if err != nil {
		log.Fatal(err)
	}

	// For the blob, show each block (ID and size) that is a committed part of it.
	getBlock, err := blobClient.GetBlockList(context.TODO(), azblob.BlockListTypeAll, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fmt.Printf("Block ID=%d, Size=%d\n", blockIDBase64ToInt(*block.Name), block.Size)
	}

	// Download the blob in its entirety; download operations do not take blocks into account.
	get, err := blobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	blobData := &bytes.Buffer{}
	reader := get.Body(nil)
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

// ExampleAppendBlobClient shows how to append data (in blocks) to an append blob.
// An append blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// The maximum size of an append blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func ExampleAppendBlobClient() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyAppendBlob.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	appendBlobClient, err := azblob.NewAppendBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = appendBlobClient.Create(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ { // Append 5 blocks to the append blob
		_, err := appendBlobClient.AppendBlock(context.TODO(), streaming.NopCloser(strings.NewReader(fmt.Sprintf("Appending block #%d\n", i))), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Download the entire append blob's contents and read into a bytes.Buffer.
	get, err := appendBlobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Buffer{}
	reader := get.Body(nil)
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

// ExamplePageBlobClient shows how to manipulate a page blob with PageBlobClient.
// A page blob is a collection of 512-byte pages optimized for random read and write operations.
// The maximum size for a page blob is 8 TB.
func ExamplePageBlobClient() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a ContainerClient object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyPageBlob.txt", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewPageBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = blobClient.Create(context.TODO(), azblob.PageBlobPageBytes*4, nil)
	if err != nil {
		log.Fatal(err)
	}

	page := make([]byte, azblob.PageBlobPageBytes)
	copy(page, "Page 0")
	_, err = blobClient.UploadPages(context.TODO(), streaming.NopCloser(bytes.NewReader(page)), nil)
	if err != nil {
		log.Fatal(err)
	}

	copy(page, "Page 1")
	_, err = blobClient.UploadPages(
		context.TODO(),
		streaming.NopCloser(bytes.NewReader(page)),
		&azblob.UploadPagesOptions{PageRange: &azblob.HttpRange{Offset: 0, Count: 2 * azblob.PageBlobPageBytes}},
	)
	if err != nil {
		log.Fatal(err)
	}

	getPages, err := blobClient.GetPageRanges(context.TODO(), azblob.HttpRange{Offset: 0, Count: 10 * azblob.PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, pr := range getPages.PageList.PageRange {
		fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
	}

	_, err = blobClient.ClearPages(context.TODO(), azblob.HttpRange{Offset: 0, Count: 1 * azblob.PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}

	getPages, err = blobClient.GetPageRanges(context.TODO(), azblob.HttpRange{Offset: 0, Count: 10 * azblob.PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, pr := range getPages.PageList.PageRange {
		fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
	}

	get, err := blobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	blobData := &bytes.Buffer{}
	reader := get.Body(nil)
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	fmt.Println(blobData.String())
}

// This example show how to create a blob, take a snapshot of it, update the base blob,
// read from the blob snapshot, list blobs with their snapshots, and delete blob snapshots.
func Example_blobSnapshots() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	containerClient, err := azblob.NewContainerClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a BlockBlobClient object to a blob in the container.
	baseBlobClient, _ := containerClient.NewBlockBlobClient("Original.txt")

	// Create the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(streaming.NopCloser(strings.NewReader("Some text"))), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a snapshot of the original blob & save its timestamp:
	createSnapshot, err := baseBlobClient.CreateSnapshot(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	snapshot := *createSnapshot.Snapshot

	// Modify the original blob:
	_, err = baseBlobClient.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("New text")), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Download the modified blob:
	get, err := baseBlobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Buffer{}
	reader := get.Body(nil)
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
	get, err = snapshotBlobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	b.Reset()
	reader = get.Body(nil)
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
	pager := containerClient.ListBlobsFlat(nil)

	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()
		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			// Process the blobs returned
			snapTime := "N/A"
			if blob.Snapshot != nil {
				snapTime = *blob.Snapshot
			}
			fmt.Printf("Blob name: %s, Snapshot: %s\n", *blob.Name, snapTime)
		}
	}

	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	// Promote read-only snapshot to writable base blob:
	_, err = baseBlobClient.StartCopyFromURL(context.TODO(), snapshotBlobClient.URL(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// When calling Delete on a base blob:
	// DeleteSnapshotsOptionOnly deletes all the base blob's snapshots but not the base blob itself
	// DeleteSnapshotsOptionInclude deletes the base blob & all its snapshots.
	// DeleteSnapshotOptionNone produces an error if the base blob has any snapshots.
	_, err = baseBlobClient.Delete(context.TODO(), &azblob.DeleteBlobOptions{DeleteSnapshots: azblob.DeleteSnapshotsOptionTypeInclude.ToPtr()})
	if err != nil {
		log.Fatal(err)
	}
}

func Example_progressUploadDownload() {
	// Create a credentials object with your Azure Storage Account name and key.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// From the Azure portal, get your Storage account blob service URL endpoint.
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)

	// Create an serviceClient object that wraps the service URL and a request pipeline to making requests.
	containerClient, err := azblob.NewContainerClientWithSharedKey(containerURL, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	blobClient, _ := containerClient.NewBlockBlobClient("Data.bin")

	// requestBody is the stream of data to write
	requestBody := streaming.NopCloser(strings.NewReader("Some text to write"))

	// Wrap the request body in a RequestBodyProgress and pass a callback function for progress reporting.
	requestProgress := streaming.NewRequestProgress(streaming.NopCloser(requestBody), func(bytesTransferred int64) {
		fmt.Printf("Wrote %d of %d bytes.", bytesTransferred, requestBody)
	})
	_, err = blobClient.Upload(context.TODO(), requestProgress, &azblob.UploadBlockBlobOptions{
		HTTPHeaders: &azblob.BlobHTTPHeaders{
			BlobContentType:        to.Ptr("text/html; charset=utf-8"),
			BlobContentDisposition: to.Ptr("attachment"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Here's how to read the blob's data with progress reporting:
	get, err := blobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Wrap the response body in a ResponseBodyProgress and pass a callback function for progress reporting.
	responseBody := streaming.NewResponseProgress(
		get.Body(nil),
		func(bytesTransferred int64) {
			fmt.Printf("Read %d of %d bytes.", bytesTransferred, *get.ContentLength)
		},
	)

	downloadedData := &bytes.Buffer{}
	_, err = downloadedData.ReadFrom(responseBody)
	if err != nil {
		return
	}
	err = responseBody.Close()
	if err != nil {
		return
	}
	fmt.Printf("Downloaded data: %s\n", downloadedData.String())
}

// This example shows how to copy a source document on the Internet to a blob.
func ExampleBlobClient_startCopy() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a ContainerClient object to a container where we'll create a blob and its snapshot.
	// Create a BlockBlobClient object to a blob in the container.
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/CopiedBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewBlobClientWithSharedKey(blobURL, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	src := "https://cdn2.auth0.com/docs/media/addons/azure_blob.svg"
	startCopy, err := blobClient.StartCopyFromURL(context.TODO(), src, nil)
	if err != nil {
		log.Fatal(err)
	}

	copyID := *startCopy.CopyID
	copyStatus := *startCopy.CopyStatus
	for copyStatus == azblob.CopyStatusTypePending {
		time.Sleep(time.Second * 2)
		getMetadata, err := blobClient.GetProperties(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		copyStatus = *getMetadata.CopyStatus
	}
	fmt.Printf("Copy from %s to %s: ID=%s, Status=%s\n", src, blobClient.URL(), copyID, copyStatus)
}

// // This example shows how to copy a large stream in blocks (chunks) to a block blob.
//func ExampleUploadFileToBlockBlob() {
//	file, err := os.Open("BigFile.bin") // Open the file we want to upload
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//		}
//	}(file)
//	fileSize, err := file.Stat() // Get the size of the file (stream)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// From the Azure portal, get your Storage account blob service URL endpoint.
//	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
//
//	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
//	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
//	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	blockBlobURL, err := azblob.NewBlockBlobClient(u, credential, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
//	response, err := UploadFileToBlockBlob(context.TODO(), file, blockBlobURL,
//		HighLevelUploadToBlockBlobOption{
//			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
//			Progress: func(bytesTransferred int64) {
//				fmt.Printf("Uploaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
//			},
//		})
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = response // Avoid compiler's "declared and not used" error
//
//	// Set up file to download the blob to
//	destFileName := "BigFile-downloaded.bin"
//	destFile, err := os.Create(destFileName)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func(destFile *os.File) {
//		_ = destFile.Close()
//
//	}(destFile)
//
//	// Perform download
//	err = DownloadBlobToFile(context.TODO(), blockBlobURL.BlobClient, 0, CountToEnd, destFile,
//		HighLevelDownloadFromBlobOptions{
//			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
//			Progress: func(bytesTransferred int64) {
//				fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
//			}})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}

// This example shows how to download a large stream with intelligent retries. Specifically, if
// the connection fails while reading, continuing to read from this stream initiates a new
// GetBlob call passing a range that starts from the last byte successfully read before the failure.
func ExampleBlobClient_Download() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Create a BlobClient object to a blob in the container (we assume the container & blob already exist).
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlob.bin", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := azblob.NewBlobClientWithSharedKey(blobURL, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	contentLength := int64(0) // Used for progress reporting to report the total number of bytes being downloaded.

	// Download returns an intelligent retryable stream around a blob; it returns an io.ReadCloser.
	dr, err := blobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	rs := dr.Body(nil)

	// NewResponseBodyProgress wraps the GetRetryStream with progress reporting; it returns an io.ReadCloser.
	stream := streaming.NewResponseProgress(
		rs,
		func(bytesTransferred int64) {
			fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
		},
	)
	defer func(stream io.ReadCloser) {
		err := stream.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stream) // The client must close the response body when finished with it

	file, err := os.Create("BigFile.bin") // Create the file to hold the downloaded blob contents.
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	written, err := io.Copy(file, stream) // Write to the file by reading from the blob (with intelligent retries).
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote %d bytes.\n", written)
}

//func ExampleUploadStreamToBlockBlob() {
//	// From the Azure portal, get your Storage account blob service URL endpoint.
//	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
//
//	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
//	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
//	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	blockBlobURL, err := azblob.NewBlockBlobClient(u, credential, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create some data to test the upload stream
//	blobSize := 8 * 1024 * 1024
//	data := make([]byte, blobSize)
//	_, err = rand.Read(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Perform UploadStreamToBlockBlob
//	bufferSize := 2 * 1024 * 1024 // Configure the size of the rotating buffers that are used when uploading
//	maxBuffers := 3               // Configure the number of rotating buffers that are used when uploading
//	_, err = UploadStreamToBlockBlob(context.TODO(), bytes.NewReader(data), blockBlobURL,
//		UploadStreamToBlockBlobOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})
//
//	// Verify that upload was successful
//	if err != nil {
//		log.Fatal(err)
//	}
//}

// This example shows how to perform various lease operations on a container.
// The same lease operations can be performed on individual blobs as well.
// A lease on a container prevents it from being deleted by others, while a lease on a blob
// protects it from both modifications and deletions.
func ExampleContainerLeaseClient() {
	// From the Azure portal, get your Storage account's name and account key.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create an containerClient object that wraps the container's URL and a default pipeline.
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	containerClient, err := azblob.NewContainerClientWithSharedKey(containerURL, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a unique ID for the lease
	// A lease ID can be any valid GUID string format. To generate UUIDs, consider the github.com/google/uuid package
	leaseID := "36b1a876-cf98-4eb2-a5c3-6d68489658ff"
	containerLeaseClient, err := containerClient.NewContainerLeaseClient(to.Ptr(leaseID))
	if err != nil {
		log.Fatal(err)
	}

	// Now acquire a lease on the container.
	// You can choose to pass an empty string for proposed ID so that the service automatically assigns one for you.
	duration := int32(60)
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(context.TODO(), &azblob.AcquireLeaseContainerOptions{Duration: &duration})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The container is leased for delete operations with lease ID", *acquireLeaseResponse.LeaseID)

	// The container cannot be deleted without providing the lease ID.
	_, err = containerLeaseClient.Delete(context.TODO(), nil)
	if err == nil {
		log.Fatal("delete should have failed")
	}
	fmt.Println("The container cannot be deleted while there is an active lease")

	// We can release the lease now and the container can be deleted.
	_, err = containerLeaseClient.ReleaseLease(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease on the container is now released")

	// Acquire a lease again to perform other operations.
	// Duration is still 60
	acquireLeaseResponse, err = containerLeaseClient.AcquireLease(context.TODO(), &azblob.AcquireLeaseContainerOptions{Duration: &duration})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The container is leased again with lease ID", *acquireLeaseResponse.LeaseID)

	// We can change the ID of an existing lease.
	newLeaseID := "6b3e65e5-e1bb-4a3f-8b72-13e9bc9cd3bf"
	changeLeaseResponse, err := containerLeaseClient.ChangeLease(context.TODO(),
		&azblob.ChangeLeaseContainerOptions{ProposedLeaseID: to.Ptr(newLeaseID)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease ID was changed to", *changeLeaseResponse.LeaseID)

	// The lease can be renewed.
	renewLeaseResponse, err := containerLeaseClient.RenewLease(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease was renewed with the same ID", *renewLeaseResponse.LeaseID)

	// Finally, the lease can be broken and we could prevent others from acquiring a lease for a period of time
	_, err = containerLeaseClient.BreakLease(context.TODO(), &azblob.BreakLeaseContainerOptions{BreakPeriod: to.Ptr[int32](60)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease was broken, and nobody can acquire a lease for 60 seconds")
}
