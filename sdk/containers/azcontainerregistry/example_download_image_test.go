// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
)

func Example_downloadImage() {
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

	// Get manifest
	manifestRes, err := client.GetManifest(ctx, "library/hello-world", "1.0.0", &azcontainerregistry.ClientGetManifestOptions{Accept: to.Ptr(string(azcontainerregistry.ContentTypeApplicationVndDockerDistributionManifestV2JSON))})
	if err != nil {
		log.Fatalf("failed to get manifest: %v", err)
	}
	reader, err := azcontainerregistry.NewDigestValidationReader(*manifestRes.DockerContentDigest, manifestRes.ManifestData)
	if err != nil {
		log.Fatalf("failed to create validation reader: %v", err)
	}
	manifest, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("failed to read manifest data: %v", err)
	}
	fmt.Printf("manifest: %s\n", manifest)

	// Get config
	var manifestJSON map[string]any
	err = json.Unmarshal(manifest, &manifestJSON)
	if err != nil {
		log.Fatalf("failed to unmarshal manifest: %v", err)
	}
	configDigest := manifestJSON["config"].(map[string]any)["digest"].(string)
	configRes, err := blobClient.GetBlob(ctx, "library/hello-world", configDigest, nil)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}
	reader, err = azcontainerregistry.NewDigestValidationReader(configDigest, configRes.BlobData)
	if err != nil {
		log.Fatalf("failed to create validation reader: %v", err)
	}
	config, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("failed to read config data: %v", err)
	}
	fmt.Printf("config: %s\n", config)

	// Get layers
	layers := manifestJSON["layers"].([]any)
	for _, layer := range layers {
		layerDigest := layer.(map[string]any)["digest"].(string)
		layerRes, err := blobClient.GetBlob(ctx, "library/hello-world", layerDigest, nil)
		if err != nil {
			log.Fatalf("failed to get layer: %v", err)
		}
		reader, err = azcontainerregistry.NewDigestValidationReader(layerDigest, layerRes.BlobData)
		if err != nil {
			log.Fatalf("failed to create validation reader: %v", err)
		}
		f, err := os.Create(strings.Split(layerDigest, ":")[1])
		if err != nil {
			log.Fatalf("failed to create blob file: %v", err)
		}
		_, err = io.Copy(f, reader)
		if err != nil {
			log.Fatalf("failed to write to the file: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("failed to close the file: %v", err)
		}
	}
}
