// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	// the name of the sample container to create
	containerName = "azblobsamplecontainer"

	// the name of the sample block blob
	blockBlobName = "azblobsampleblockblob"

	// the name of the sample append blob
	appendBlobName = "azblobsampleappendblob"
)

// returns a credential that can be used to authenticate with blob storage
func getCredential() azcore.Credential {
	// NewEnvironmentCredential() will read various environment vars
	// to obtain a credential.  see the documentation for more info.
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		panic(err)
	}
	return cred
}

// comes from environment var AZURE_STORAGE_ENDPOINT
// e.g. https://<your_storage_account>.blob.core.windows.net
func getEndpoint() string {
	storageEndpoint := os.Getenv("AZURE_STORAGE_ENDPOINT")
	if storageEndpoint == "" {
		panic("missing environment variable AZURE_STORAGE_ENDPOINT")
	}
	return storageEndpoint
}

// returns a byte slice for dummy blob content
func generateBlobContent(size int) []byte {
	content := make([]byte, size, size)
	for i := 0; i < size; i++ {
		content[i] = byte(i % 256)
	}
	return content
}

// concatenates all elements together with a '/' between each element
func pathJoin(elem ...string) string {
	// add a trailing '/' if it doesn't already end with one
	if i := len(elem[0]); elem[0][i-1] != '/' {
		elem[0] = elem[0] + "/"
	}
	// join all the parts together
	for i := 1; i < len(elem); i++ {
		elem[0] = elem[0] + elem[i] + "/"
	}
	return elem[0]
}

func ExampleContainerOperations_Create() {
	endpoint := pathJoin(getEndpoint(), containerName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	containerClient := client.ContainerOperations()
	c, err := containerClient.Create(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.RawResponse.StatusCode)
	// Output: 201
}

func ExampleBlockBlobOperations_Upload() {
	endpoint := pathJoin(getEndpoint(), containerName, blockBlobName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blockSize := int32(80)
	blockID := "myblockID"
	blobClient := client.BlockBlobOperations()
	block := Block{
		Name: &blockID,
		Size: &blockSize,
	}
	body := azcore.NopCloser(bytes.NewReader(generateBlobContent(int(blockSize))))
	b, err := blobClient.Upload(context.Background(), block, body, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(b.RawResponse.StatusCode)
	// Output: 201
}

/*func ExampleContainerOperations_ListBlobFlatSegment() {
	endpoint := pathJoin(storageEndpoint, containerName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	page, err := blobClient.ListBlobFlatSegment(nil)
	if err != nil {
		panic(err)
	}
	for page.NextPage(context.Background()) {
		resp := page.PageResponse()
		fmt.Println(*resp.EnumerationResults)
	}
	if err = page.Err(); err != nil {
		panic(err)
	}
	if page.PageResponse() == nil {
		panic("unexpected nil payload")
	}
}*/

func ExampleAppendBlobOperations_Create() {
	endpoint := pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.AppendBlobOperations()
	s := "text/plain"
	a, err := blobClient.Create(context.Background(), int64(10), &AppendBlobCreateOptions{BlobContentType: &s})
	if err != nil {
		panic(err)
	}
	fmt.Println(a.RawResponse.StatusCode)
	// Output: 201
}

func ExampleAppendBlobOperations_AppendBlock() {
	endpoint := pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.AppendBlobOperations()
	const blockSize = 80
	body := azcore.NopCloser(bytes.NewReader(generateBlobContent(int(blockSize))))
	a, err := blobClient.AppendBlock(context.Background(), int64(blockSize), body, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(a.RawResponse.StatusCode)
	// Output: 201
}

func ExampleBlobOperations_Delete_BlockBlob() {
	endpoint := pathJoin(getEndpoint(), containerName, blockBlobName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlobOperations(nil)
	d, err := blobClient.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)
	// Output: 202
}

func ExampleBlobOperations_Delete_AppendBlob() {
	endpoint := pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlobOperations(nil)
	d, err := blobClient.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)
	// Output: 202
}

func ExampleContainerOperations_Delete() {
	endpoint := pathJoin(getEndpoint(), containerName)
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	d, err := blobClient.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)
	// Output: 202
}
