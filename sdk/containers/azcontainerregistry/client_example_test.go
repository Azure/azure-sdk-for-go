//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleClient_DeleteManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	resp, err := client.GetTagProperties(ctx, "alpine", "3.7", nil)
	if err != nil {
		log.Fatalf("failed to get tag properties: %v", err)
	}
	_, err = client.DeleteManifest(ctx, "alpine", *resp.Tag.Digest, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_DeleteRepository() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteRepository(ctx, "nanoserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_DeleteTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.DeleteTag(ctx, "nanoserver", "4.7.2-20180905-nanoserver-1803", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_GetManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetManifest(ctx, "hello-world-dangling", "20190628-033033z", &azcontainerregistry.ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	manifest, err := io.ReadAll(res.ManifestData)
	fmt.Printf("manifest content: %s\n", manifest)
}

func ExampleClient_GetManifestProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
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

func ExampleClient_GetRepositoryProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
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

func ExampleClient_GetTagProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
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

func ExampleClient_NewListManifestsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListManifestsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.ACRManifests.Manifests {
			fmt.Printf("manifest %d: %s\n", i+1, *v.Digest)
		}
	}
}

func ExampleClient_NewListRepositoriesPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Repositories.Names {
			fmt.Printf("repository %d: %s\n", i+1, *v)
		}
	}
}

func ExampleClient_NewListTagsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListTagsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Tags {
			fmt.Printf("tag %d: %s\n", i+1, *v.Name)
		}
	}
}

func ExampleClient_UpdateManifestProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateManifestProperties(ctx, "nanoserver", "sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086", &azcontainerregistry.ClientUpdateManifestPropertiesOptions{Value: &azcontainerregistry.ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository nanoserver - manifest sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086 - 'CanWrite' property: %t", *res.Manifest.ChangeableAttributes.CanWrite)
}
func ExampleClient_UpdateRepositoryProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateRepositoryProperties(ctx, "nanoserver", &azcontainerregistry.ClientUpdateRepositoryPropertiesOptions{Value: &azcontainerregistry.RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - 'CanWrite' property: %t\n", *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite)
}

func ExampleClient_UpdateTagProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateTagProperties(ctx, "nanoserver", "4.7.2-20180905-nanoserver-1803", &azcontainerregistry.ClientUpdateTagPropertiesOptions{
		Value: &azcontainerregistry.TagWriteableProperties{
			CanWrite: to.Ptr(false),
		}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - tag 4.7.2-20180905-nanoserver-1803 - 'CanWrite' property: %t\n", *res.Tag.ChangeableAttributes.CanWrite)
}

func ExampleClient_UploadManifest() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := azcontainerregistry.NewClient("https://example.azurecr.io", cred, nil)
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
