// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
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

func ExampleBlobClient_GetBlob() {
	const digest = "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39"
	res, err := blobClient.GetBlob(context.TODO(), "prod/bash", digest, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	reader, err := azcontainerregistry.NewDigestValidationReader(digest, res.BlobData)
	if err != nil {
		log.Fatalf("failed to create validation reader: %v", err)
	}
	f, err := os.Create("blob_file")
	if err != nil {
		log.Fatalf("failed to create blob file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close blob file: %v", err)
		}
	}()
	_, err = io.Copy(f, reader)
	if err != nil {
		log.Printf("failed to write to the file: %v", err)
	}
}

func ExampleBlobClient_GetChunk() {
	chunkSize := 1024 * 1024
	const digest = "sha256:16463e0c481e161aabb735437d30b3c9c7391c2747cc564bb927e843b73dcb39"
	current := 0
	f, err := os.Create("blob_file")
	if err != nil {
		log.Fatalf("failed to create blob file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close blob file: %v", err)
		}
	}()
	for {
		res, err := blobClient.GetChunk(context.TODO(), "prod/bash", digest, fmt.Sprintf("bytes=%d-%d", current, current+chunkSize-1), nil)
		if err != nil {
			log.Printf("failed to finish the request: %v", err)
			return
		}
		chunk, err := io.ReadAll(res.ChunkData)
		if err != nil {
			log.Printf("failed to read the chunk: %v", err)
			return
		}
		_, err = f.Write(chunk)
		if err != nil {
			log.Printf("failed to write to the file: %v", err)
			return
		}

		totalSize, _ := strconv.Atoi(strings.Split(*res.ContentRange, "/")[1])
		currentRangeEnd, _ := strconv.Atoi(strings.Split(strings.Split(*res.ContentRange, "/")[0], "-")[1])
		if totalSize == currentRangeEnd+1 {
			break
		}
		current += chunkSize
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("failed to set to the start of the file: %v", err)
		return
	}
	reader, err := azcontainerregistry.NewDigestValidationReader(digest, f)
	if err != nil {
		log.Printf("failed to create digest validation reader: %v", err)
		return
	}
	_, err = io.ReadAll(reader)
	if err != nil {
		log.Printf("failed to validate digest: %v", err)
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
