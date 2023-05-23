//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"io"
	"log"
)

func Example_uploadAndDownloadBlob() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azcontainerregistry.NewBlobClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	blob := []byte("hello world")
	startRes, err := client.StartUpload(ctx, "library/hello-world", nil)
	if err != nil {
		log.Fatalf("failed to start upload blob: %v", err)
	}
	calculator := azcontainerregistry.NewBlobDigestCalculator()
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, bytes.NewReader(blob), calculator, nil)
	if err != nil {
		log.Fatalf("failed to upload blob: %v", err)
	}
	completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	if err != nil {
		log.Fatalf("failed to complete upload: %v", err)
	}
	fmt.Printf("uploaded blob digest: %s", *completeResp.DockerContentDigest)
	downloadRes, err := client.GetBlob(ctx, "library/hello-world", *completeResp.DockerContentDigest, nil)
	if err != nil {
		log.Fatalf("failed to download blob: %v", err)
	}
	reader, err := azcontainerregistry.NewDigestValidationReader(*completeResp.DockerContentDigest, downloadRes.BlobData)
	downloadBlob, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("failed to read blob: %v", err)
	}
	fmt.Printf("blob content: %s", downloadBlob)
}
