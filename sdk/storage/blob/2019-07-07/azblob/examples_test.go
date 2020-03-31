// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"
	"fmt"
	"testing"

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

func TestExample_UploadBlockBlob(t *testing.T) {
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

func TestExample_ListContainer(t *testing.T) {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	client, err := NewClient(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}
	blobClient := client.ContainerOperations()
	b, err := blobClient.ListBlobFlatSegment(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		if next := b.NextPage(context.Background()); next {
			if b.Err() != nil {
				panic(b.Err())
			}
			fmt.Println(b.PageResponse())
		}
	}
}
