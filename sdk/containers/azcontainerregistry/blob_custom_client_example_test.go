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

func ExampleNewBlobClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	blobClient, err = azcontainerregistry.NewBlobClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create blob client: %v", err)
	}
	_ = blobClient
}

func ExampleBlobClient_CompleteUpload() {
	// calculator should be created when starting upload blob and passing to UploadChunk and CompleteUpload method
	calculator := azcontainerregistry.NewBlobDigestCalculator()
	res, err := blobClient.CompleteUpload(context.TODO(), "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", calculator, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("content digest: %s", *res.DockerContentDigest)
}

func ExampleBlobClient_UploadChunk() {
	// calculator should be created when starting upload blob and passing to UploadChunk and CompleteUpload method
	calculator := azcontainerregistry.NewBlobDigestCalculator()
	res, err := blobClient.UploadChunk(context.TODO(), "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", streaming.NopCloser(bytes.NewReader([]byte("U29tZXRoaW5nRWxzZQ=="))), calculator, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload location: %s", *res.Location)
}

func ExampleBlobClient_GetBlob() {
	res, err := blobClient.GetBlob(context.TODO(), "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// deal with the blob io
	_ = res.BlobData
}

func ExampleBlobClient_GetChunk() {
	// calculator should be created when starting get blob and passing to GetChunk method
	calculator := azcontainerregistry.NewBlobDigestCalculator()
	res, err := blobClient.GetChunk(context.TODO(), "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", "bytes=0-299", calculator, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// deal with the chunk io
	_ = res.ChunkData
}
