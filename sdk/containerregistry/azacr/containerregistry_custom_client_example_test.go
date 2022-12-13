//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"log"
)

func ExampleContainerRegistryClient_DeleteManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteManifest(ctx, "alpine", "3.7", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleContainerRegistryClient_GetDigest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	digest, err := client.GetDigest(ctx, "alpine", "3.7", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("digest for tag alpine:3.7 is %s", digest)
}

func ExampleContainerRegistryClient_GetManifestProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetManifestProperties(ctx, "nanoserver", "sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("manifest digest: %s\n", *res.Manifest.Digest)
	fmt.Printf("manifest size: %d\n", *res.Manifest.Size)
}

func ExampleContainerRegistryClient_UpdateManifestProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateManifestProperties(ctx, "nanoserver", "sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086", &azacr.ContainerRegistryClientUpdateManifestPropertiesOptions{Value: &azacr.ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository nanoserver - manifest sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086 - 'CanWrite' property: %t", *res.Manifest.ChangeableAttributes.CanWrite)
}
