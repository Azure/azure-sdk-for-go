//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blobbatch_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blobbatch"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ExampleContainerBatchClient shows blob batch operations for delete and set tier using ContainerBatchClient
func Example_containerBatchClient() {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	const containerName = "testcontainer"

	// create shared key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// create container batch client
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)
	cntBatchClient, err := blobbatch.NewContainerBatchClientWithSharedKeyCredential(containerURL, cred, nil)

	// create new batch builder
	bb := blobbatch.NewBatchBuilder()

	// add operations to the batch builder
	bb.AddBlobBatchDelete("testBlob0", nil)
	bb.AddBlobBatchDelete("testBlob1", &blobbatch.DeleteOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"), // version id for deletion
	})
	bb.AddBlobBatchDelete("testBlob2", &blobbatch.DeleteOptions{
		Snapshot: to.Ptr("2023-01-03T11:57:25.6515618Z"), // snapshot for deletion
	})
	bb.AddBlobBatchDelete("testBlob3", &blobbatch.DeleteOptions{
		BlobDeleteOptions: &blob.DeleteOptions{
			DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeOnly),
			BlobDeleteType:  to.Ptr(blob.DeleteTypeNone),
		},
	})

	bb.AddBlobBatchSetTier("testBlob4", blob.AccessTierHot, nil)
	bb.AddBlobBatchSetTier("testBlob5", blob.AccessTierCool, &blobbatch.SetTierOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"),
	})

	resp, err := cntBatchClient.SubmitBatch(context.TODO(), bb)
	handleError(err)

	// print response body
	p := make([]byte, 10000)
	_, err = resp.Body.Read(p)
	handleError(err)
	fmt.Println(string(p))
}

// ExampleServiceBatchClient shows blob batch operations for delete and set tier using ServiceBatchClient
func Example_serviceBatchClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	tenantID, ok := os.LookupEnv("AZURE_TENANT_ID")
	if !ok {
		panic("AZURE_TENANT_ID could not be found")
	}
	clientID, ok := os.LookupEnv("AZURE_CLIENT_ID")
	if !ok {
		panic("AZURE_CLIENT_ID could not be found")
	}
	clientSecret, ok := os.LookupEnv("AZURE_CLIENT_SECRET")
	if !ok {
		panic("AZURE_CLIENT_SECRET could not be found")
	}

	// create client secret credential
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	handleError(err)

	// create service batch client
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	svcBatchClient, err := blobbatch.NewServiceBatchClient(serviceURL, cred, nil)

	// create new batch builder
	bb := blobbatch.NewBatchBuilder()

	// add operations to the batch builder
	bb.AddBlobBatchDelete("/container1/testBlob0", nil)
	bb.AddBlobBatchDelete("/container2/testBlob1", &blobbatch.DeleteOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"), // version id for deletion
	})
	bb.AddBlobBatchDelete("/container1/testBlob2", &blobbatch.DeleteOptions{
		Snapshot: to.Ptr("2023-01-03T11:57:25.6515618Z"), // snapshot for deletion
	})
	bb.AddBlobBatchDelete("/container2/testBlob3", &blobbatch.DeleteOptions{
		BlobDeleteOptions: &blob.DeleteOptions{
			DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeOnly),
			BlobDeleteType:  to.Ptr(blob.DeleteTypeNone),
		},
	})

	bb.AddBlobBatchSetTier("/container3/testBlob4", blob.AccessTierHot, nil)
	bb.AddBlobBatchSetTier("/container4/testBlob5", blob.AccessTierCool, &blobbatch.SetTierOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"),
	})

	resp, err := svcBatchClient.SubmitBatch(context.TODO(), bb)
	handleError(err)

	// print response body
	p := make([]byte, 10000)
	_, err = resp.Body.Read(p)
	handleError(err)
	fmt.Println(string(p))
}
