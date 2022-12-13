//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/containerregistry/azacr"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleContainerRegistryClient_CheckDockerV2Support() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CheckDockerV2Support(ctx, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleContainerRegistryClient_GetManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetManifest(ctx, "hello-world-dangling", "20190628-033033z", &azacr.ContainerRegistryClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	manifest, err := io.ReadAll(res.Body)
	fmt.Printf("manifest content: %s\n", manifest)
}

func ExampleContainerRegistryClient_NewListRepositoriesPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Repositories.Repositories {
			fmt.Printf("repository %d: %s\n", i+1, *v)
		}
	}
}

func ExampleContainerRegistryClient_GetRepositoryProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetRepositoryProperties(ctx, "nanoserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository name: %s\n", *res.Name)
	fmt.Printf("registry login server of the repository: %s\n", *res.RegistryLoginServer)
	fmt.Printf("repository manifest count: %d\n", *res.ManifestCount)
}

func ExampleContainerRegistryClient_DeleteRepository() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteRepository(ctx, "nanoserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleContainerRegistryClient_UpdateRepositoryProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateRepositoryProperties(ctx, "nanoserver", &azacr.ContainerRegistryClientUpdateRepositoryPropertiesOptions{Value: &azacr.RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - 'CanWrite' property: %t\n", *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite)
}

func ExampleContainerRegistryClient_NewListTagsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListTagsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.TagAttributeBases {
			fmt.Printf("tag %d: %s\n", i+1, *v.Name)
		}
	}
}

func ExampleContainerRegistryClient_GetTagProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetTagProperties(ctx, "test/bash", "latest", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("tag name: %s\n", *res.Tag.Name)
	fmt.Printf("tag digest: %s\n", *res.Tag.Digest)
}

func ExampleContainerRegistryClient_UpdateTagProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateTagProperties(ctx, "nanoserver", "4.7.2-20180905-nanoserver-1803", &azacr.ContainerRegistryClientUpdateTagPropertiesOptions{
		Value: &azacr.TagWriteableProperties{
			CanWrite: to.Ptr(false),
		}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - tag 4.7.2-20180905-nanoserver-1803 - 'CanWrite' property: %t\n", *res.Tag.ChangeableAttributes.CanWrite)
}

func ExampleContainerRegistryClient_DeleteTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteTag(ctx, "nanoserver", "4.7.2-20180905-nanoserver-1803", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleContainerRegistryClient_NewListManifestsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListManifestsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Manifests.Manifests {
			fmt.Printf("manifest %d: %s\n", i+1, *v.Digest)
		}
	}
}

func ExampleContainerRegistryClient_UploadManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azacr.NewContainerRegistryClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	payload, err := os.Open("example-manifest.json")
	if err != nil {
		log.Fatalf("failed to read manifest file: %v", err)
	}
	resp, err := client.UploadManifest(ctx, "nanoserver", "test", "application/vnd.docker.distribution.manifest.v2+json", payload, nil)
	fmt.Printf("uploaded manifest digest: %s", *resp.DockerContentDigest)
}
