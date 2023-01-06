//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleBlobClient_CancelUpload() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CancelUpload(ctx, "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleBlobClient_CheckBlobExists() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CheckBlobExists(ctx, "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("blob digest: %s", *res.DockerContentDigest)
}

func ExampleBlobClient_CheckChunkExists() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CheckChunkExists(ctx, "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", "bytes=0-299", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("chunk size: %d", *res.ContentLength)
	fmt.Printf("chunk range: %s", *res.ContentRange)
}

func ExampleBlobClient_CompleteUpload() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CompleteUpload(ctx, "sha256:a05a7abc31b0caadd5058069eaa623cdd0538512dbddacc7ca5bad73c3957c5e", "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", &azcontainerregistry.BlobClientCompleteUploadOptions{
		Value: streaming.NopCloser(bytes.NewReader([]byte("usuallyEmpty"))),
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("content digest: %s", *res.DockerContentDigest)
}

func ExampleBlobClient_DeleteBlob() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteBlob(ctx, "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleBlobClient_GetBlob() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetBlob(ctx, "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// deal with the blob io
	_ = res.Body
}

func ExampleBlobClient_GetChunk() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetChunk(ctx, "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", "bytes=0-299", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// deal with the chunk io
	_ = res.Body
}

func ExampleBlobClient_GetUploadStatus() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetUploadStatus(ctx, "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload UUID: %s", *res.DockerUploadUUID)
}

func ExampleBlobClient_MountBlob() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.MountBlob(ctx, "newimage", "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	fmt.Printf("upload location: %s", *res.Location)
}

func ExampleBlobClient_StartUpload() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.StartUpload(ctx, "newimg", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload location: %s", *res.Location)
}

func ExampleBlobClient_UploadChunk() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UploadChunk(ctx, "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", streaming.NopCloser(bytes.NewReader([]byte("U29tZXRoaW5nRWxzZQ=="))), nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload location: %s", *res.Location)
}
