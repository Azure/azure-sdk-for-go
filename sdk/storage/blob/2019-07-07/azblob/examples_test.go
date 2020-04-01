// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	endpoint = "<endpoint>"
)

const (
	tenantID     = "<tenantID>"
	clientID     = "<clientID>"
	clientSecret = "<clientSecret"
)

const (
	accountName = "<accountName>"
	accountKey  = "<accountKey>"
)

func getCredential() azcore.Credential {
	// cred, err := NewSharedKeyCredential(accountName, accountKey)
	// if err != nil {
	// 	panic(err)
	// }
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	return cred
}

func ExampleContainerOperations_Create() {
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	containerClient := client.ContainerOperations()
	c, err := containerClient.Create(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.RawResponse)
}

func ExampleBlockBlobOperations_Upload() {
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlockBlobOperations()
	block := Block{
		Name: "myblockID",
		Size: 80,
	}
	// TODO: body
	b, err := blobClient.Upload(context.Background(), block, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}

func ExampleContainerOperations_ListBlobFlatSegment() {
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
}

func ExampleAppendBlobOperations_Create() {
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
}

func ExampleBlobOperations_Delete() {
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlobOperations(nil)
	d, err := blobClient.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}

func ExampleContainerOperations_Delete() {
	client, err := NewClient(endpoint, getCredential(), nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	d, err := blobClient.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}
