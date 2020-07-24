// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

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
	// check environment for shared key credential info
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	if accountName != "" && accountKey != "" {
		keyCred, err := NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			panic(err)
		}
		return keyCred
	}
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
	// remove trailing '/' if needed
	if i := len(elem[0]); elem[0][i-1] == '/' {
		elem[0] = elem[0][0 : i-1]
	}
	return strings.Join(elem, "/")
}

func ExampleContainerOperations_Create() {
	endpoint := pathJoin(getEndpoint(), containerName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	containerClient := client.ContainerOperations()
	c, err := containerClient.Create(context.Background(), nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.RawResponse.StatusCode)
	// Output: 201
}

func ExampleServiceOperations_ListContainersSegment() {
	client, err := newClient(getEndpoint(), getCredential(), nil)
	if err != nil {
		panic(err)
	}
	pager, err := client.ServiceOperations().ListContainersSegment(nil)
	if err != nil {
		panic(err)
	}
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		for _, container := range *page.EnumerationResults.ContainerItems {
			fmt.Println(*container.Name)
		}
	}
	if pager.Err() != nil {
		panic(pager.Err())
	}
	// Output: azblobsamplecontainer
}

func ExampleBlockBlobOperations_Upload() {
	endpoint := pathJoin(getEndpoint(), containerName, blockBlobName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	const blockSize = 80
	blobClient := client.BlockBlobOperations()
	body := azcore.NopCloser(bytes.NewReader(generateBlobContent(blockSize)))
	b, err := blobClient.Upload(context.Background(), blockSize, body, nil, nil, nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(b.RawResponse.StatusCode)
	// Output: 201
}

func ExampleBlobOperations_Download() {
	endpoint := pathJoin(getEndpoint(), containerName, blockBlobName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlobOperations(nil)
	b, err := blobClient.Download(context.Background(), nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(b.RawResponse.StatusCode)
	fmt.Println(string(*b.BlobType))
	// Output:
	// 200
	// BlockBlob
}

func ExampleAppendBlobOperations_Create() {
	endpoint := pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.AppendBlobOperations()
	a, err := blobClient.Create(context.Background(), 0, nil, nil, nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(a.RawResponse.StatusCode)
	// Output: 201
}

func ExampleAppendBlobOperations_AppendBlock() {
	endpoint := pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.AppendBlobOperations()
	const blockSize = 80
	body := azcore.NopCloser(bytes.NewReader(generateBlobContent(int(blockSize))))
	a, err := blobClient.AppendBlock(context.Background(), int64(blockSize), body, nil, nil, nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(a.RawResponse.StatusCode)
	// Output: 201
}

func ExampleContainerOperations_ListBlobFlatSegment() {
	endpoint := pathJoin(getEndpoint(), containerName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	pager, err := blobClient.ListBlobFlatSegment(&ContainerListBlobFlatSegmentOptions{
		Include: &[]ListBlobsIncludeItem{
			ListBlobsIncludeItemUncommittedblobs,
		},
	})
	if err != nil {
		panic(err)
	}
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		for _, blob := range *page.EnumerationResults.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
	if err = pager.Err(); err != nil {
		panic(err)
	}
	// Unordered output:
	// azblobsampleblockblob
	// azblobsampleappendblob
}

func ExampleBlobOperations_Delete() {
	endpoint := pathJoin(getEndpoint(), containerName, blockBlobName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlobOperations(nil)
	d, err := blobClient.Delete(context.Background(), nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)

	endpoint = pathJoin(getEndpoint(), containerName, appendBlobName)
	client, err = newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient = client.BlobOperations(nil)
	d, err = blobClient.Delete(context.Background(), nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)
	// Output:
	// 202
	// 202
}

func ExampleContainerOperations_Delete() {
	endpoint := pathJoin(getEndpoint(), containerName)
	client, err := newClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	d, err := blobClient.Delete(context.Background(), nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.RawResponse.StatusCode)
	// Output: 202
}
