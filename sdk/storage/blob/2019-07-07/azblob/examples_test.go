// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"
	"fmt"

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

func Example_CreateContainer() {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	client, err := NewClient(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}
	containerClient := client.ContainerOperations()
	c, err := containerClient.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.RawResponse)
}

func Example_UploadBlockBlob() {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	// cred, err := NewSharedKeyCredential(accountName, accountKey)
	// if err != nil {
	// 	panic(err)
	// }
	client, err := NewClient(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.BlockBlobOperations()
	cl := int64(80)
	b, err := blobClient.Upload(context.Background(), "myblockID", cl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
}

func Example_ListContainer() {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	client, err := NewClient(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	page, err := blobClient.ListBlobFlatSegment(nil)
	if err != nil {
		fmt.Println(err)
		return
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
