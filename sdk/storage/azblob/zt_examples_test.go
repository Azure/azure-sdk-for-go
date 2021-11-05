// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func accountInfo() (string, string) {
	return os.Getenv("ACCOUNT_NAME"), os.Getenv("ACCOUNT_KEY")
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

// This example shows you how to get started using the Azure Storage Blob SDK for Go.
func Example() {
	// Use your storage account's name and key to create a credential object, used to access your account.
	// You can obtain these details from the Azure Portal.
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a service client.
	// You'll need to specify a service URL, which for blob endpoints usually makes up the syntax http(s)://<account>.blob.core.windows.net/
	service, err := NewServiceClientWithSharedKey("https://"+accountName+".blob.core.windows.net/", cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// All operations in the Azure Storage Blob SDK for Go operate on a context.Context, allowing you to control cancellation/timeout.
	ctx := context.Background() // This example has no expiry.

	// This example showcases several common operations to help you get started, such as:

	// ===== 1. Creating a container =====

	// First, branch off of the service client and create a container client.
	container := service.NewContainerClient("myContainer")
	// Then, fire off a create operation on the container client.
	// Note that, all service-side requests have an options bag attached, allowing you to specify things like metadata, public access types, etc.
	// Specifying nil omits all options.
	_, err = container.Create(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// ===== 2. Uploading/downloading a block blob =====
	// We'll specify our data up-front, rather than reading a file for simplicity's sake.
	data := "Hello world!"

	// Branch off of the container into a block blob client
	blockBlob := container.NewBlockBlobClient("HelloWorld.txt")

	// Upload data to the block blob
	_, err = blockBlob.Upload(ctx, NopCloser(strings.NewReader(data)), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Download the blob's contents and ensure that the download worked properly
	get, err := blockBlob.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Open a buffer, reader, and then download!
	downloadedData := &bytes.Buffer{}
	reader := get.Body(RetryReaderOptions{}) // RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	}
	if data != downloadedData.String() {
		log.Fatal("downloaded data doesn't match uploaded data")
	}

	// ===== 3. list blobs =====
	// The ListBlobs and ListContainers APIs return two channels, a values channel, and an errors channel.
	// You should enumerate on a range over the values channel, and then check the errors channel, as only ONE value will ever be passed to the errors channel.
	// The AutoPagerTimeout defines how long it will wait to place into the items channel before it exits & cleans itself up. A zero time will result in no timeout.
	pager := container.ListBlobsFlat(nil)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			fmt.Println(*v.Name)
		}
	}

	if err = pager.Err(); err != nil {
		log.Fatal(err)
	}

	// Delete the blob we created earlier.
	_, err = blockBlob.Delete(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the container we created earlier.
	_, err = container.Delete(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Awaiting Mohit's test PR, StorageError is barren for now.
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

	container, err := NewContainerClientWithNoCredential("https://myaccount.blob.core.windows.net/mycontainer", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = container.Create(context.Background(), nil)

	if err != nil { // an error occurred
		var stgErr StorageError

		if errors.As(err, &stgErr) { // We know this error is service-specific
			switch stgErr.ErrorCode {
			case StorageErrorCodeContainerAlreadyExists:
				// You can also look at the *http.Response that's attached to the error as well.
				if resp := stgErr.Response(); resp != nil {
					failedRequest := resp.Request
					_ = failedRequest // avoid compiler's declared but not used error
				}
			case StorageErrorCodeContainerBeingDeleted:
				// Handle this error ...
			default:
				// Handle other errors ...
			}
		}
	}
}

// This example demonstrates splitting a URL into its parts so you can examine and modify the URL in an Azure Storage fluent way.
func ExampleBlobURLParts() {
	// Let's begin with a snapshot SAS token.
	u := "https://myaccount.blob.core.windows.net/mycontainter/ReadMe.txt?" +
		"snapshot=2011-03-09T01:42:34Z&" +
		"sv=2015-02-21&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&" +
		"spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=92836758923659283652983562=="

	// Breaking the URL down into it's parts by conversion to BlobURLParts
	parts := NewBlobURLParts(u)

	// Now, we can access the parts (this example prints them.)
	fmt.Println(parts.Host, parts.ContainerName, parts.BlobName, parts.Snapshot)
	sas := parts.SAS
	fmt.Println(sas.Version(), sas.Resource(), sas.StartTime(), sas.ExpiryTime(), sas.Permissions(),
		sas.IPRange(), sas.Protocol(), sas.Identifier(), sas.Services(), sas.Signature())

	// You can also alter some of the fields and construct a new URL:
	// Note that: SAS tokens may be limited to a specific container or blob.
	// You should be careful about modifying SAS tokens, as you might take them outside of their original scope accidentally.
	parts.SAS = SASQueryParameters{}
	parts.Snapshot = ""
	parts.ContainerName = "othercontainer"

	// construct a new URL from the parts
	fmt.Print(parts.URL())
}

// This example demonstrates how to use the SAS token convenience generators.
// Though this example focuses on account SAS, these generators exist across all clients (Service, Container, Blob, and specialized Blob clients)
func ExampleServiceClient_GetSASToken() {
	// Initialize a service client
	accountName, accountKey := accountInfo()
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	serviceClient, err := NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), credential, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Provide the convenience function with relevant info (services, resource types, permissions, and duration)
	// The SAS token will be valid from this moment onwards.
	accountSAS, err := serviceClient.GetSASToken(AccountSASResourceTypes{Object: true, Service: true, Container: true},
		AccountSASPermissions{Read: true, List: true}, AccountSASServices{Blob: true}, time.Now(), time.Now().Add(48*time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	urlToSend := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, accountSAS)
	// You can hand off this URL to someone else via any mechanism you choose.

	// ******************************************

	// When someone receives the URL, they can access the resource using it in code like this, or a tool of some variety.
	serviceClient, err = NewServiceClientWithNoCredential(urlToSend, nil)
	if err != nil {
		log.Fatal(err)
	}

	// You can also break a blob URL up into it's constituent parts
	blobURLParts := NewBlobURLParts(serviceClient.URL())
	fmt.Printf("SAS expiry time = %s\n", blobURLParts.SAS.ExpiryTime())
}

// This example shows how to create and use an Azure Storage account Shared Access Signature (SAS).
func ExampleAccountSASSignatureValues_Sign() {
	accountName, accountKey := accountInfo()

	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.Sign(credential)
	if err != nil {
		log.Fatal(err)
	}

	qp := sasQueryParams.Encode()
	urlToSend := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, qp)
	// You can hand off this URL to someone else via any mechanism you choose.

	// ******************************************

	// When someone receives the URL, they can access the resource using it in code like this, or a tool of some variety.
	serviceClient, err := NewServiceClientWithNoCredential(urlToSend, nil)
	if err != nil {
		log.Fatal(err)
	}

	// You can also break a blob URL up into it's constituent parts
	blobURLParts := NewBlobURLParts(serviceClient.URL())
	fmt.Printf("SAS expiry time = %s\n", blobURLParts.SAS.ExpiryTime())
}

// This example demonstrates how to create and use a Blob service Shared Access Signature (SAS)
func ExampleBlobSASSignatureValues() {
	// Gather your account key and name from the Azure portal
	// Supplying them via environment variables is recommended.
	accountName, accountKey := accountInfo()

	// Use your storage account's name and key to form a credential object
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	containerName := "myContainer"
	blobName := "HelloWorld.txt"

	sasQueryParams, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		ContainerName: containerName,
		BlobName:      blobName,

		// To produce a container SAS, as opposed to a blob SAS, assign to permissions using ContainerSASPermissions
		// and make sure the BlobName field is ""
		Permissions: BlobSASPermissions{Add: true, Read: true, Write: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	// Create the URL of this resource you wish to access, and append the SAS query parameters.
	// Since this is a blob SAS, the URL is to the Azure Storage blob.
	qp := sasQueryParams.Encode()
	urlToSendToSomeone := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		accountName, containerName, blobName, qp)

	// At this point, you can send the URL to someone via communication method of your choice, and it will provide
	// anonymous access to the resource.

	// **************

	// When someone receives the URL, they can access the SAS-protected resource like this:
	blob, _ := NewBlobClientWithNoCredential(urlToSendToSomeone, nil)

	// if you have a SAS query parameter string, you can parse it into it's parts.
	blobURLParts := NewBlobURLParts(blob.URL())
	fmt.Printf("SAS expiry time=%v", blobURLParts.SAS.ExpiryTime())
}

// This example shows how to manipulate a container's permissions.
func ExampleContainerClient_SetAccessPolicy() {
	// Obtain your storage account's name and key from the Azure portal
	accountName, accountKey := accountInfo()

	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	uri := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	container, err := NewContainerClientWithSharedKey(uri, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Grab the background context, use no expiry
	ctx := context.Background()

	// Create the container (with no metadata and no public access)
	_, err = container.Create(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Upload a simple blob.
	blob := container.NewBlockBlobClient("HelloWorld.txt")

	_, err = blob.Upload(ctx, NopCloser(strings.NewReader("Hello World!")), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Attempt to read the blob
	get, err := http.Get(blob.URL())
	if err != nil {
		log.Fatal(err)
	}
	if get.StatusCode == http.StatusNotFound {
		_, err := container.SetAccessPolicy(ctx, &SetAccessPolicyOptions{ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: PublicAccessTypeBlob.ToPtr()}})
		if err != nil {
			log.Fatal(err)
		}

		// Now, this works:
		get, err = http.Get(blob.URL())
		if err != nil {
			log.Fatal(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(get.Body)
		var text bytes.Buffer
		_, err = text.ReadFrom(get.Body)
		if err != nil {
			return
		}
		fmt.Println(text.String())
	}
}

// This example shows how to perform operations on blob conditionally.
func ExampleBlobAccessConditions() {
	// From the Azure portal, get your Storage account's name and account key.
	accountName, accountKey := accountInfo()

	// Create a BlockBlobClient object that wraps a blob's URL and a default pipeline.
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blockBlob, err := NewBlockBlobClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/Data.txt", accountName), credential, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background() // This example uses a never-expiring context

	// This helper function displays the results of an operation; it is called frequently below.
	showResult := func(response *DownloadResponse, err error) {
		if err != nil {
			if stgErr, ok := err.(*StorageError); !ok {
				log.Fatal(err) // Network failure
			} else {
				// TODO: Port storage error
				fmt.Print("Failure: " + stgErr.Error() + "\n")
			}
		} else {
			err := response.Body(RetryReaderOptions{}).Close()
			if err != nil {
				return
			} // The client must close the response body when finished with it
			fmt.Print("Success: " + response.RawResponse.Status + "\n")
		}
	}

	showResultUpload := func(upload BlockBlobUploadResponse, err error) {
		if err != nil {
			if stgErr, ok := err.(*StorageError); !ok {
				log.Fatal(err) // Network failure
			} else {
				// TODO: Port storage error
				fmt.Print("Failure: " + stgErr.Error() + "\n")
			}
		} else {
			fmt.Print("Success: " + upload.RawResponse.Status + "\n")
		}
	}

	// Create the blob (unconditionally; succeeds)
	upload, err := blockBlob.Upload(ctx, NopCloser(strings.NewReader("Text-1")), nil)
	showResultUpload(upload, err)

	// showResult(upload, err)

	// Download blob content if the blob has been modified since we uploaded it (fails):
	showResult(blockBlob.Download(ctx, &DownloadBlobOptions{BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: upload.LastModified}}}))

	// Download blob content if the blob hasn't been modified in the last 24 hours (fails):
	showResult(blockBlob.Download(ctx, &DownloadBlobOptions{BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: to.TimePtr(time.Now().UTC().Add(time.Hour * -24))}}}))

	// Upload new content if the blob hasn't changed since the version identified by ETag (succeeds):
	upload, err = blockBlob.Upload(ctx, NopCloser(strings.NewReader("Text-2")),
		&UploadBlockBlobOptions{
			BlobAccessConditions: &BlobAccessConditions{
				ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: upload.ETag},
			},
		})
	showResultUpload(upload, err)

	// Download content if it has changed since the version identified by ETag (fails):
	showResult(blockBlob.Download(ctx,
		&DownloadBlobOptions{
			BlobAccessConditions: &BlobAccessConditions{
				ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: upload.ETag}},
		}))

	// Upload content if the blob doesn't already exist (fails):
	showResultUpload(blockBlob.Upload(ctx, NopCloser(strings.NewReader("Text-3")),
		&UploadBlockBlobOptions{
			BlobAccessConditions: &BlobAccessConditions{
				ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.StringPtr(ETagAny)},
			}}))
}

// This examples shows how to create a container with metadata and then how to read & update the metadata.
func ExampleContainerClient_SetMetadata() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a containerClient object that wraps a soon-to-be-created container's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	containerClient, err := NewContainerClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context

	// Create a container with some metadata (string key/value pairs)
	// NOTE: Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Therefore, you should always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, _ := os.Executable()
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: map[string]string{"author": "Jeffrey", "app": creatingApp}})
	if err != nil {
		log.Fatal(err)
	}

	// Query the container's metadata
	get, err := containerClient.GetProperties(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Show the container's metadata
	if get.Metadata == nil {
		log.Fatal("metadata is empty!")
	}

	metadata := get.Metadata
	for k, v := range metadata {
		fmt.Print(k + "=" + v + "\n")
	}

	// Update the metadata and write it back to the container
	metadata["author"] = "Aidan" // NOTE: The keyname is in all lowercase letters
	_, err = containerClient.SetMetadata(ctx, &SetMetadataContainerOptions{Metadata: metadata})
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: The SetMetadata & SetProperties methods update the container's ETag & LastModified properties
}

// This examples shows how to create a blob with metadata and then how to read & update
// the blob's read-only properties and metadata.
func ExampleBlobClient_SetMetadata() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	BlobClient, err := NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context

	// Create a blob with metadata (string key/value pairs)
	// NOTE: Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Therefore, you should always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, _ := os.Executable()
	_, err = BlobClient.Upload(ctx, NopCloser(NopCloser(strings.NewReader("Some text"))), &UploadBlockBlobOptions{Metadata: map[string]string{"author": "Jeffrey", "app": creatingApp}})
	if err != nil {
		log.Fatal(err)
	}

	// Query the blob's properties and metadata
	get, err := BlobClient.GetProperties(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Show some of the blob's read-only properties
	fmt.Println(*get.BlobType, *get.ETag, *get.LastModified)

	// Show the blob's metadata
	if get.Metadata == nil {
		log.Fatal("No metadata returned")
	}

	metadata := get.Metadata
	for k, v := range metadata {
		fmt.Print(k + "=" + v + "\n")
	}

	// Update the blob's metadata and write it back to the blob
	metadata["editor"] = "Grant" // Add a new key/value; NOTE: The keyname is in all lowercase letters
	_, err = BlobClient.SetMetadata(ctx, metadata, nil)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: The SetMetadata method updates the blob's ETag & LastModified properties
}

// This examples shows how to create a blob with HTTP Headers and then how to read & update
// the blob's HTTP headers.
func ExampleBlobHTTPHeaders() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a blob client
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/ReadMe.txt", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context

	// Create a blob with HTTP headers
	_, err = blobClient.Upload(ctx, NopCloser(NopCloser(strings.NewReader("Some text"))),
		&UploadBlockBlobOptions{HTTPHeaders: &BlobHTTPHeaders{
			BlobContentType:        to.StringPtr("text/html; charset=utf-8"),
			BlobContentDisposition: to.StringPtr("attachment"),
		}})
	if err != nil {
		log.Fatal(err)
	}

	// GetMetadata returns the blob's properties, HTTP headers, and metadata
	get, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Show some of the blob's read-only properties
	fmt.Println(*get.BlobType, *get.ETag, *get.LastModified)

	// Shows some of the blob's HTTP Headers
	httpHeaders := get.GetHTTPHeaders()
	fmt.Println(httpHeaders.BlobContentType, httpHeaders.BlobContentDisposition)

	// Update the blob's HTTP Headers and write them back to the blob
	httpHeaders.BlobContentType = to.StringPtr("text/plain")
	_, err = blobClient.SetHTTPHeaders(ctx, httpHeaders, nil)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: The SetMetadata method updates the blob's ETag & LastModified properties
}

// ExampleBlockBlobClient shows how to upload a lot of data (in blocks) to a blob.
// A block blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// Therefore, the maximum size of a block blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func ExampleBlockBlobClient() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a ContainerClient object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyBlockBlob.txt", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	BlobClient, err := NewBlockBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context

	// These helper functions convert a binary block ID to a base-64 string and vice versa
	// NOTE: The blockID must be <= 64 bytes and ALL blockIDs for the block must be the same length
	blockIDBinaryToBase64 := func(blockID []byte) string { return base64.StdEncoding.EncodeToString(blockID) }
	blockIDBase64ToBinary := func(blockID string) []byte { _binary, _ := base64.StdEncoding.DecodeString(blockID); return _binary }

	// These helper functions convert an int block ID to a base-64 string and vice versa
	blockIDIntToBase64 := func(blockID int) string {
		binaryBlockID := (&[4]byte{})[:] // All block IDs are 4 bytes long
		binary.LittleEndian.PutUint32(binaryBlockID, uint32(blockID))
		return blockIDBinaryToBase64(binaryBlockID)
	}
	blockIDBase64ToInt := func(blockID string) int {
		blockIDBase64ToBinary(blockID)
		return int(binary.LittleEndian.Uint32(blockIDBase64ToBinary(blockID)))
	}

	// Upload 4 blocks to the blob (these blocks are tiny; they can be up to 100MB each)
	words := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(words)) // The collection of block IDs (base 64 strings)

	// Upload each block sequentially (one after the other); for better performance, you want to upload multiple blocks in parallel)
	for index, word := range words {
		// This example uses the index as the block ID; convert the index/ID into a base-64 encoded string as required by the service.
		// NOTE: Over the lifetime of a blob, all block IDs (before base 64 encoding) must be the same length (this example uses 4 byte block IDs).
		base64BlockIDs[index] = blockIDIntToBase64(index) // Some people use UUIDs for block IDs

		// Upload a block to this blob specifying the Block ID and its content (up to 100MB); this block is uncommitted.
		_, err := BlobClient.StageBlock(ctx, base64BlockIDs[index], NopCloser(strings.NewReader(word)), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// After all the blocks are uploaded, atomically commit them to the blob.
	_, err = BlobClient.CommitBlockList(ctx, base64BlockIDs, nil)
	if err != nil {
		log.Fatal(err)
	}

	// For the blob, show each block (ID and size) that is a committed part of it.
	getBlock, err := BlobClient.GetBlockList(ctx, BlockListTypeAll, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range getBlock.BlockList.CommittedBlocks {
		fmt.Printf("Block ID=%d, Size=%d\n", blockIDBase64ToInt(*block.Name), block.Size)
	}

	// Download the blob in its entirety; download operations do not take blocks into account.
	// NOTE: For really large blobs, downloading them like allocates a lot of memory.
	get, err := BlobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	blobData := &bytes.Buffer{}
	reader := get.Body(RetryReaderOptions{})
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	fmt.Println(blobData)
}

// ExampleAppendBlobClient shows how to append data (in blocks) to an append blob.
// An append blob can have a maximum of 50,000 blocks; each block can have a maximum of 100MB.
// Therefore, the maximum size of an append blob is slightly more than 4.75 TB (100 MB X 50,000 blocks).
func ExampleAppendBlobClient() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a ContainerClient object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyAppendBlob.txt", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	appendBlobClient, err := NewAppendBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context
	_, err = appendBlobClient.Create(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ { // Append 5 blocks to the append blob
		_, err := appendBlobClient.AppendBlock(ctx, NopCloser(strings.NewReader(fmt.Sprintf("Appending block #%d\n", i))), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Download the entire append blob's contents and show it.
	get, err := appendBlobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Buffer{}
	reader := get.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	fmt.Println(b.String())
}

// ExamplePageBlobClient shows how to manipulate a page blob with PageBlobClient.
// A page blob is a collection of 512-byte pages optimized for random read and write operations.
// The maximum size for a page blob is 8 TB.
func ExamplePageBlobClient() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a ContainerClient object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/MyPageBlob.txt", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := NewPageBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context
	_, err = blobClient.Create(ctx, PageBlobPageBytes*4, nil)
	if err != nil {
		log.Fatal(err)
	}

	page := [PageBlobPageBytes]byte{}
	copy(page[:], "Page 0")
	_, err = blobClient.UploadPages(ctx, NopCloser(bytes.NewReader(page[:])), nil)
	if err != nil {
		log.Fatal(err)
	}

	copy(page[:], "Page 1")
	_, err = blobClient.UploadPages(ctx, NopCloser(bytes.NewReader(page[:])), &UploadPagesOptions{PageRange: &HttpRange{0, 2 * PageBlobPageBytes}})
	if err != nil {
		log.Fatal(err)
	}

	getPages, err := blobClient.GetPageRanges(ctx, HttpRange{0 * PageBlobPageBytes, 10 * PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, pr := range getPages.PageList.PageRange {
		fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
	}

	_, err = blobClient.ClearPages(ctx, HttpRange{0 * PageBlobPageBytes, 1 * PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}

	getPages, err = blobClient.GetPageRanges(ctx, HttpRange{0 * PageBlobPageBytes, 10 * PageBlobPageBytes}, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, pr := range getPages.PageList.PageRange {
		fmt.Printf("Start=%d, End=%d\n", pr.Start, pr.End)
	}

	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	blobData := &bytes.Buffer{}
	reader := get.Body(RetryReaderOptions{})
	_, err = blobData.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	fmt.Printf("%#v", blobData.Bytes())
}

// This example show how to create a blob, take a snapshot of it, update the base blob,
// read from the blob snapshot, list blobs with their snapshots, and hot to delete blob snapshots.
func Example_blobSnapshots() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a ContainerClient object to a container where we'll create a blob and its snapshot.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	containerClient, err := NewContainerClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a BlockBlobClient object to a blob in the container.
	baseBlobClient := containerClient.NewBlockBlobClient("Original.txt")

	ctx := context.Background() // This example uses a never-expiring context

	// Create the original blob:
	_, err = baseBlobClient.Upload(ctx, NopCloser(NopCloser(strings.NewReader("Some text"))), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a snapshot of the original blob & save its timestamp:
	createSnapshot, err := baseBlobClient.CreateSnapshot(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	snapshot := *createSnapshot.Snapshot

	// Modify the original blob & show it:
	_, err = baseBlobClient.Upload(ctx, NopCloser(strings.NewReader("New text")), nil)
	if err != nil {
		log.Fatal(err)
	}

	get, err := baseBlobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	b := bytes.Buffer{}
	reader := get.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	fmt.Println(b.String())

	// Show snapshot blob via original blob URI & snapshot time:
	snapshotBlobClient := baseBlobClient.WithSnapshot(snapshot)
	get, err = snapshotBlobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	b.Reset()
	reader = get.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	if err != nil {
		return
	}
	err = reader.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	fmt.Println(b.String())

	// FYI: You can get the base blob URL from one of its snapshot by passing "" to WithSnapshot:
	baseBlobClient = snapshotBlobClient.WithSnapshot("")

	// Show all blobs in the container with their snapshots:
	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	pager := containerClient.ListBlobsFlat(nil)

	for pager.NextPage(ctx) {
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
	_, err = baseBlobClient.StartCopyFromURL(ctx, snapshotBlobClient.URL(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// When calling Delete on a base blob:
	// DeleteSnapshotsOptionOnly deletes all the base blob's snapshots but not the base blob itself
	// DeleteSnapshotsOptionInclude deletes the base blob & all its snapshots.
	// DeleteSnapshotOptionNone produces an error if the base blob has any snapshots.
	_, err = baseBlobClient.Delete(ctx, &DeleteBlobOptions{DeleteSnapshots: DeleteSnapshotsOptionTypeInclude.ToPtr()})
	if err != nil {
		log.Fatal(err)
	}
}

func Example_progressUploadDownload() {
	// Create a credentials object with your Azure Storage Account name and key.
	accountName, accountKey := accountInfo()
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// From the Azure portal, get your Storage account blob service URL endpoint.
	cURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)

	// Create an serviceClient object that wraps the service URL and a request pipeline to making requests.
	containerClient, err := NewContainerClientWithSharedKey(cURL, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context
	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	blobClient := containerClient.NewBlockBlobClient("Data.bin")

	// requestBody is the stream of data to write
	requestBody := NopCloser(strings.NewReader("Some text to write"))

	// Wrap the request body in a RequestBodyProgress and pass a callback function for progress reporting.
	_, err = blobClient.Upload(ctx, streaming.NewRequestProgress(NopCloser(requestBody), func(bytesTransferred int64) {
		fmt.Printf("Wrote %d of %d bytes.", bytesTransferred, requestBody)
	}), &UploadBlockBlobOptions{HTTPHeaders: &BlobHTTPHeaders{
		BlobContentType:        to.StringPtr("text/html; charset=utf-8"),
		BlobContentDisposition: to.StringPtr("attachment"),
	}})
	if err != nil {
		log.Fatal(err)
	}

	// Here's how to read the blob's data with progress reporting:
	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Wrap the response body in a ResponseBodyProgress and pass a callback function for progress reporting.
	responseBody := streaming.NewResponseProgress(get.Body(RetryReaderOptions{}),
		func(bytesTransferred int64) {
			fmt.Printf("Read %d of %d bytes.", bytesTransferred, *get.ContentLength)
		})

	downloadedData := &bytes.Buffer{}
	_, err = downloadedData.ReadFrom(responseBody)
	if err != nil {
		return
	}
	err = responseBody.Close()
	if err != nil {
		return
	} // The client must close the response body when finished with it
	// The downloaded blob data is in downloadData's buffer
}

// This example shows how to copy a source document on the Internet to a blob.
func ExampleBlobClient_startCopy() {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := accountInfo()

	// Create a ContainerClient object to a container where we'll create a blob and its snapshot.
	// Create a BlockBlobClient object to a blob in the container.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/CopiedBlob.bin", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := NewBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background() // This example uses a never-expiring context

	src := "https://cdn2.auth0.com/docs/media/addons/azure_blob.svg"
	startCopy, err := blobClient.StartCopyFromURL(ctx, src, nil)
	if err != nil {
		log.Fatal(err)
	}

	copyID := *startCopy.CopyID
	copyStatus := *startCopy.CopyStatus
	for copyStatus == CopyStatusTypePending {
		time.Sleep(time.Second * 2)
		getMetadata, err := blobClient.GetProperties(ctx, nil)
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
//	accountName, accountKey := accountInfo()
//
//	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
//	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
//	credential, err := NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	blockBlobURL, err := NewBlockBlobClient(u, credential, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background() // This example uses a never-expiring context
//
//	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
//	response, err := UploadFileToBlockBlob(ctx, file, blockBlobURL,
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
//	err = DownloadBlobToFile(context.Background(), blockBlobURL.BlobClient, 0, CountToEnd, destFile,
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
	accountName, accountKey := accountInfo()

	// Create a BlobClient object to a blob in the container (we assume the container & blob already exist).
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlob.bin", accountName)
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobClient, err := NewBlobClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	contentLength := int64(0) // Used for progress reporting to report the total number of bytes being downloaded.

	// Download returns an intelligent retryable stream around a blob; it returns an io.ReadCloser.
	dr, err := blobClient.Download(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	rs := dr.Body(RetryReaderOptions{})

	// NewResponseBodyProgress wraps the GetRetryStream with progress reporting; it returns an io.ReadCloser.
	stream := streaming.NewResponseProgress(rs,
		func(bytesTransferred int64) {
			fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
		})
	defer func(stream io.ReadCloser) {
		err := stream.Close()
		if err != nil {

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
	_ = written // Avoid compiler's "declared and not used" error
}

//func ExampleUploadStreamToBlockBlob() {
//	// From the Azure portal, get your Storage account blob service URL endpoint.
//	accountName, accountKey := accountInfo()
//
//	// Create a BlockBlobURL object to a blob in the container (we assume the container already exists).
//	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlockBlob.bin", accountName)
//	credential, err := NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	blockBlobURL, err := NewBlockBlobClient(u, credential, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	ctx := context.Background() // This example uses a never-expiring context
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
//	_, err = UploadStreamToBlockBlob(ctx, bytes.NewReader(data), blockBlobURL,
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
	accountName, accountKey := accountInfo()

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create an containerClient object that wraps the container's URL and a default pipeline.
	u := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName)
	containerClient, err := NewContainerClientWithSharedKey(u, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	generatedUuid, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	leaseID := to.StringPtr(generatedUuid.String())
	containerLeaseClient, err := containerClient.NewContainerLeaseClient(leaseID)
	if err != nil {
		log.Fatal(err)
	}

	// All operations allow you to specify a timeout via a Go context.Context object.
	ctx := context.Background() // This example uses a never-expiring context

	// Now acquire a lease on the container.
	// You can choose to pass an empty string for proposed ID so that the service automatically assigns one for you.
	duration := int32(60)
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: &duration})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The container is leased for delete operations with lease ID", *acquireLeaseResponse.LeaseID)

	// The container cannot be deleted without providing the lease ID.
	_, err = containerLeaseClient.Delete(ctx, nil)
	if err == nil {
		log.Fatal("delete should have failed")
	}
	fmt.Println("The container cannot be deleted while there is an active lease")

	// We can release the lease now and the container can be deleted.
	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease on the container is now released")

	// Acquire a lease again to perform other operations.
	// Duration is still 60
	acquireLeaseResponse, err = containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: &duration})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The container is leased again with lease ID", *acquireLeaseResponse.LeaseID)

	// We can change the ID of an existing lease.
	// A lease ID can be any valid GUID string format.
	newLeaseID, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}

	newLeaseID[0] = 1
	changeLeaseResponse, err := containerLeaseClient.ChangeLease(ctx,
		&ChangeLeaseContainerOptions{ProposedLeaseID: to.StringPtr(newLeaseID.String())})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease ID was changed to", *changeLeaseResponse.LeaseID)

	// The lease can be renewed.
	renewLeaseResponse, err := containerLeaseClient.RenewLease(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease was renewed with the same ID", *renewLeaseResponse.LeaseID)

	// Finally, the lease can be broken and we could prevent others from acquiring a lease for a period of time
	duration = 60
	_, err = containerLeaseClient.BreakLease(ctx, &BreakLeaseContainerOptions{BreakPeriod: &duration})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The lease was broken, and nobody can acquire a lease for 60 seconds")
}
