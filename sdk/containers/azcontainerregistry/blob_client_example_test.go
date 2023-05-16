//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"log"
)

var blobClient *azcontainerregistry.BlobClient

func ExampleBlobClient_CancelUpload() {
	_, err := blobClient.CancelUpload(context.TODO(), "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleBlobClient_CheckBlobExists() {
	res, err := blobClient.CheckBlobExists(context.TODO(), "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("blob digest: %s", *res.DockerContentDigest)
}

func ExampleBlobClient_CheckChunkExists() {
	res, err := blobClient.CheckChunkExists(context.TODO(), "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", "bytes=0-299", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("chunk size: %d", *res.ContentLength)
	fmt.Printf("chunk range: %s", *res.ContentRange)
}

func ExampleBlobClient_DeleteBlob() {
	_, err := blobClient.DeleteBlob(context.TODO(), "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleBlobClient_GetUploadStatus() {
	res, err := blobClient.GetUploadStatus(context.TODO(), "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload UUID: %s", *res.DockerUploadUUID)
}

func ExampleBlobClient_MountBlob() {
	res, err := blobClient.MountBlob(context.TODO(), "newimage", "prod/bash", "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload location: %s", *res.Location)
}

func ExampleBlobClient_StartUpload() {
	res, err := blobClient.StartUpload(context.TODO(), "newimg", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("upload location: %s", *res.Location)
}
