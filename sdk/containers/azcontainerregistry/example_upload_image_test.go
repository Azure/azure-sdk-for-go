// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
)

func Example_uploadImage() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azcontainerregistry.NewClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	blobClient, err := azcontainerregistry.NewBlobClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create blob client: %v", err)
	}
	ctx := context.Background()
	layer := []byte("hello world")
	startRes, err := blobClient.StartUpload(ctx, "library/hello-world", nil)
	if err != nil {
		log.Fatalf("failed to start upload layer: %v", err)
	}
	calculator := azcontainerregistry.NewBlobDigestCalculator()
	uploadResp, err := blobClient.UploadChunk(ctx, *startRes.Location, bytes.NewReader(layer), calculator, nil)
	if err != nil {
		log.Fatalf("failed to upload layer: %v", err)
	}
	completeResp, err := blobClient.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	if err != nil {
		log.Fatalf("failed to complete layer upload: %v", err)
	}
	layerDigest := *completeResp.DockerContentDigest
	config := []byte(fmt.Sprintf(`{
  architecture: "amd64",
  os: "windows",
  rootfs: {
	type: "layers",
	diff_ids: [%s],
  },
}`, layerDigest))
	startRes, err = blobClient.StartUpload(ctx, "library/hello-world", nil)
	if err != nil {
		log.Fatalf("failed to start upload config: %v", err)
	}
	calculator = azcontainerregistry.NewBlobDigestCalculator()
	uploadResp, err = blobClient.UploadChunk(ctx, *startRes.Location, bytes.NewReader(config), calculator, nil)
	if err != nil {
		log.Fatalf("failed to upload config: %v", err)
	}
	completeResp, err = blobClient.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	if err != nil {
		log.Fatalf("failed to complete config upload: %v", err)
	}
	manifest := fmt.Sprintf(`{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
	"mediaType": "application/vnd.oci.image.config.v1+json",
	"digest": "%s",
	"size": %d
  },
  "layers": [
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar",
	  "digest": "%s",
	  "size": %d,
	  "annotations": {
		"title": "artifact.txt"
	  }
    }
  ]
}`, layerDigest, len(config), *completeResp.DockerContentDigest, len(layer))
	uploadManifestRes, err := client.UploadManifest(ctx, "library/hello-world", "1.0.0", azcontainerregistry.ContentTypeApplicationVndDockerDistributionManifestV2JSON, streaming.NopCloser(bytes.NewReader([]byte(manifest))), nil)
	if err != nil {
		log.Fatalf("failed to upload manifest: %v", err)
	}
	fmt.Printf("digest of uploaded manifest: %s", *uploadManifestRes.DockerContentDigest)
}
