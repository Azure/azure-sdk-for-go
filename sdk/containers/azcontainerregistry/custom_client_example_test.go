//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"io"
	"log"
	"os"
)

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err = azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_ = client
}

func ExampleClient_GetManifest() {
	res, err := client.GetManifest(context.TODO(), "hello-world-dangling", "20190628-033033z", &azcontainerregistry.ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	manifest, err := io.ReadAll(res.ManifestData)
	if err != nil {
		log.Fatalf("failed to read manifest data: %v", err)
	}
	fmt.Printf("manifest content: %s\n", manifest)
}

func ExampleClient_UploadManifest() {
	payload, err := os.Open("example-manifest.json")
	if err != nil {
		log.Fatalf("failed to read manifest file: %v", err)
	}
	resp, err := client.UploadManifest(context.TODO(), "nanoserver", "test", "application/vnd.docker.distribution.manifest.v2+json", payload, nil)
	if err != nil {
		log.Fatalf("failed to upload manifest: %v", err)
	}
	fmt.Printf("uploaded manifest digest: %s", *resp.DockerContentDigest)
}
