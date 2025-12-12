// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
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
	location := "v2/blobland/blobs/uploads/2b28c60d-d296-44b7-b2b4-1f01c63195c6?_nouploadcache=false&_state=VYABvUSCNW2yY5e5VabLHppXqwU0K7cvT0YUdq57KBt7Ik5hbWUiOiJibG9ibGFuZCIsIlVVSUQiOiIyYjI4YzYwZC1kMjk2LTQ0YjctYjJiNC0xZjAxYzYzMTk1YzYiLCJPZmZzZXQiOjAsIlN0YXJ0ZWRBdCI6IjIwMTktMDgtMjdUMjM6NTI6NDcuMDUzNjU2Mjg1WiJ9"
	f, err := os.Open("blob-file")
	if err != nil {
		log.Fatalf("failed to read blob file: %v", err)
	}
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatalf("failed to calculate blob size: %v", err)
	}
	chunkSize := int64(5)
	current := int64(0)
	for {
		end := current + chunkSize
		if end > size {
			end = size
		}
		chunkReader := io.NewSectionReader(f, current, end-current)
		uploadResp, err := blobClient.UploadChunk(context.TODO(), location, chunkReader, calculator, &azcontainerregistry.BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(current)), RangeEnd: to.Ptr(int32(end - 1))})
		if err != nil {
			log.Fatalf("failed to upload chunk: %v", err)
		}
		location = *uploadResp.Location
		current = end
		if current >= size {
			break
		}
	}
	fmt.Printf("upload location: %s", location)
}
