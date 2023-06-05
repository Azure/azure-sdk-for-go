//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_container_NewClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientWithSharedKeyCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, cred, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientWithNoCredential() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	sharedAccessSignature, ok := os.LookupEnv("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE")
	if !ok {
		panic("AZURE_STORAGE_SHARED_ACCESS_SIGNATURE could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s?%s", accountName, containerName, sharedAccessSignature)

	containerClient, err := container.NewClientWithNoCredential(containerURL, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_NewClientFromConnectionString() {
	// Your connection string can be obtained from the Azure Portal.
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}
	containerName := "testcontainer"
	containerClient, err := container.NewClientFromConnectionString(connectionString, containerName, nil)
	handleError(err)
	fmt.Println(containerClient.URL())
}

func Example_container_ClientNewAppendBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	appendBlobClient := containerClient.NewAppendBlobClient("test_append_blob")
	handleError(err)
	fmt.Println(appendBlobClient.URL())
}

func Example_container_ClientNewBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	blobClient := containerClient.NewBlobClient("test_blob")
	handleError(err)
	fmt.Println(blobClient.URL())
}

func Example_container_ClientNewBlockBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	blockBlobClient := containerClient.NewBlockBlobClient("test_block_blob")
	handleError(err)
	fmt.Println(blockBlobClient.URL())
}

func Example_container_ClientNewPageBlobClient() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	pageBlobClient := containerClient.NewPageBlobClient("test_page_blob")
	handleError(err)
	fmt.Println(pageBlobClient.URL())
}

func Example_container_ClientCreate() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	containerCreateResponse, err := containerClient.Create(context.TODO(), &container.CreateOptions{
		Metadata: map[string]*string{"Foo": to.Ptr("Bar")},
	})
	handleError(err)
	fmt.Println(containerCreateResponse)
}

func Example_container_ClientDelete() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	containerDeleteResponse, err := containerClient.Delete(context.TODO(), nil)
	handleError(err)
	fmt.Println(containerDeleteResponse)
}

func Example_container_ClientListBlobsFlat() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
}

func Example_container_ClientListBlobsHierarchy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	maxResults := int32(5)
	pager := containerClient.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{
		Include:    container.ListBlobsInclude{Metadata: true, Tags: true},
		MaxResults: &maxResults,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, blob := range resp.ListBlobsHierarchySegmentResponse.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
}

func Example_container_ClientGetSASURL() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	permission := sas.ContainerPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	options := container.GetSASURLOptions{StartTime: &start}
	sasURL, err := containerClient.GetSASURL(permission, expiry, &options)
	handleError(err)
	_ = sasURL
}

// This example shows how to manipulate a container's permissions.
func Example_container_ClientSetAccessPolicy() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	// Create the container
	_, err = containerClient.Create(context.TODO(), nil)
	handleError(err)

	// Upload a simple blob.
	blockBlobClient := containerClient.NewBlockBlobClient("HelloWorld.txt")
	handleError(err)

	_, err = blockBlobClient.Upload(context.TODO(), streaming.NopCloser(strings.NewReader("Hello World!")), nil)
	handleError(err)

	// Attempt to read the blob
	get, err := http.Get(blockBlobClient.URL())
	handleError(err)
	if get.StatusCode == http.StatusNotFound {
		// ChangeLease the blob to be public access blob
		_, err := containerClient.SetAccessPolicy(
			context.TODO(),
			&container.SetAccessPolicyOptions{
				Access: to.Ptr(container.PublicAccessTypeBlob),
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		// Now, this works
		get, err = http.Get(blockBlobClient.URL())
		if err != nil {
			log.Fatal(err)
		}
		var text bytes.Buffer
		_, err = text.ReadFrom(get.Body)
		if err != nil {
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(get.Body)

		fmt.Println("Public access blob data: ", text.String())
	}
}

func Example_container_ClientSetMetadata() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	containerName := "testcontainer"
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	// Create a container with some metadata, key names are converted to lowercase before being sent to the service.
	// You should always use lowercase letters, especially when querying a map for a metadata key.
	creatingApp, err := os.Executable()
	handleError(err)
	_, err = containerClient.Create(context.TODO(), &container.CreateOptions{Metadata: map[string]*string{"author": to.Ptr("azblob"), "app": to.Ptr(creatingApp)}})
	handleError(err)

	// Query the container's metadata
	containerGetPropertiesResponse, err := containerClient.GetProperties(context.TODO(), nil)
	handleError(err)

	if containerGetPropertiesResponse.Metadata == nil {
		log.Fatal("metadata is empty!")
	}

	for k, v := range containerGetPropertiesResponse.Metadata {
		fmt.Printf("%s=%s\n", k, *v)
	}

	// Update the metadata and write it back to the container
	containerGetPropertiesResponse.Metadata["author"] = to.Ptr("Mohit")
	_, err = containerClient.SetMetadata(context.TODO(), &container.SetMetadataOptions{Metadata: containerGetPropertiesResponse.Metadata})
	handleError(err)

	// NOTE: SetMetadata & SetProperties methods update the container's ETag & LastModified properties
}

// ExampleContainerBatchDelete shows blob batch operations for delete and set tier.
func Example_container_BatchDelete() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	const containerName = "testcontainer"

	// create shared key credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// create container batch client
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)
	cntBatchClient, err := container.NewClientWithSharedKeyCredential(containerURL, cred, nil)
	handleError(err)

	// create new batch builder
	bb, err := cntBatchClient.NewBatchBuilder()
	handleError(err)

	// add operations to the batch builder
	err = bb.Delete("testBlob0", nil)
	handleError(err)

	err = bb.Delete("testBlob1", &container.BatchDeleteOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"), // version id for deletion
	})
	handleError(err)

	err = bb.Delete("testBlob2", &container.BatchDeleteOptions{
		Snapshot: to.Ptr("2023-01-03T11:57:25.6515618Z"), // snapshot for deletion
	})
	handleError(err)

	err = bb.Delete("testBlob3", &container.BatchDeleteOptions{
		DeleteOptions: blob.DeleteOptions{
			DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeOnly),
			BlobDeleteType:  to.Ptr(blob.DeleteTypeNone),
		},
	})
	handleError(err)

	resp, err := cntBatchClient.SubmitBatch(context.TODO(), bb, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get response for individual sub-requests
	for _, resp := range resp.Responses {
		if resp.ContainerName != nil && resp.BlobName != nil {
			fmt.Println("Container: " + *resp.ContainerName)
			fmt.Println("Blob: " + *resp.BlobName)
		}
		if resp.Error == nil {
			fmt.Println("Successful sub-request")
		} else {
			fmt.Println("Error: " + resp.Error.Error())
		}
	}
}

// ExampleContainerBatchSetTier shows blob batch operations for delete and set tier.
func Example_container_BatchSetTier() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	tenantID, ok := os.LookupEnv("AZURE_STORAGE_TENANT_ID")
	if !ok {
		panic("AZURE_STORAGE_TENANT_ID could not be found")
	}
	clientID, ok := os.LookupEnv("AZURE_STORAGE_CLIENT_ID")
	if !ok {
		panic("AZURE_STORAGE_CLIENT_ID could not be found")
	}
	clientSecret, ok := os.LookupEnv("AZURE_STORAGE_CLIENT_SECRET")
	if !ok {
		panic("AZURE_STORAGE_CLIENT_SECRET could not be found")
	}

	const containerName = "testcontainer"

	// create client secret credential
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	handleError(err)

	// create container batch client
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)
	cntBatchClient, err := container.NewClient(containerURL, cred, nil)
	handleError(err)

	// create new batch builder
	bb, err := cntBatchClient.NewBatchBuilder()
	handleError(err)

	// add operations to the batch builder
	err = bb.SetTier("testBlob1", blob.AccessTierHot, nil)
	handleError(err)

	err = bb.SetTier("testBlob2", blob.AccessTierCool, &container.BatchSetTierOptions{
		VersionID: to.Ptr("2023-01-03T11:57:25.4067017Z"),
	})
	handleError(err)

	err = bb.SetTier("testBlob3", blob.AccessTierCool, &container.BatchSetTierOptions{
		Snapshot: to.Ptr("2023-01-03T11:57:25.6515618Z"),
	})
	handleError(err)

	err = bb.SetTier("testBlob4", blob.AccessTierCool, &container.BatchSetTierOptions{
		SetTierOptions: blob.SetTierOptions{
			RehydratePriority: to.Ptr(blob.RehydratePriorityStandard),
		},
	})
	handleError(err)

	resp, err := cntBatchClient.SubmitBatch(context.TODO(), bb, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get response for individual sub-requests
	for _, resp := range resp.Responses {
		if resp.ContainerName != nil && resp.BlobName != nil {
			fmt.Println("Container: " + *resp.ContainerName)
			fmt.Println("Blob: " + *resp.BlobName)
		}
		if resp.Error == nil {
			fmt.Println("Successful sub-request")
		} else {
			fmt.Println("Error: " + resp.Error.Error())
		}
	}
}
