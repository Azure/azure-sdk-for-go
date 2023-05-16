//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"log"
)

var client *azcontainerregistry.Client

func ExampleClient_DeleteManifest() {
	resp, err := client.GetTagProperties(context.TODO(), "alpine", "3.7", nil)
	if err != nil {
		log.Fatalf("failed to get tag properties: %v", err)
	}
	_, err = client.DeleteManifest(context.TODO(), "alpine", *resp.Tag.Digest, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_DeleteRepository() {
	_, err := client.DeleteRepository(context.TODO(), "nanoserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_DeleteTag() {
	_, err := client.DeleteTag(context.TODO(), "nanoserver", "4.7.2-20180905-nanoserver-1803", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func ExampleClient_GetManifestProperties() {
	res, err := client.GetManifestProperties(context.TODO(), "nanoserver", "sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("manifest digest: %s\n", *res.Manifest.Digest)
	fmt.Printf("manifest size: %d\n", *res.Manifest.Size)
}

func ExampleClient_GetRepositoryProperties() {
	res, err := client.GetRepositoryProperties(context.TODO(), "nanoserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository name: %s\n", *res.Name)
	fmt.Printf("registry login server of the repository: %s\n", *res.RegistryLoginServer)
	fmt.Printf("repository manifest count: %d\n", *res.ManifestCount)
}

func ExampleClient_GetTagProperties() {
	res, err := client.GetTagProperties(context.TODO(), "test/bash", "latest", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("tag name: %s\n", *res.Tag.Name)
	fmt.Printf("tag digest: %s\n", *res.Tag.Digest)
}

func ExampleClient_NewListManifestsPager() {
	pager := client.NewListManifestsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Manifests.Attributes {
			fmt.Printf("manifest %d: %s\n", i+1, *v.Digest)
		}
	}
}

func ExampleClient_NewListRepositoriesPager() {
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Repositories.Names {
			fmt.Printf("repository %d: %s\n", i+1, *v)
		}
	}
}

func ExampleClient_NewListTagsPager() {
	pager := client.NewListTagsPager("nanoserver", nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for i, v := range page.Tags {
			fmt.Printf("tag %d: %s\n", i+1, *v.Name)
		}
	}
}

func ExampleClient_UpdateManifestProperties() {
	res, err := client.UpdateManifestProperties(context.TODO(), "nanoserver", "sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086", &azcontainerregistry.ClientUpdateManifestPropertiesOptions{Value: &azcontainerregistry.ManifestWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository nanoserver - manifest sha256:110d2b6c84592561338aa040b1b14b7ab81c2f9edbd564c2285dd7d70d777086 - 'CanWrite' property: %t", *res.Manifest.ChangeableAttributes.CanWrite)
}
func ExampleClient_UpdateRepositoryProperties() {
	res, err := client.UpdateRepositoryProperties(context.TODO(), "nanoserver", &azcontainerregistry.ClientUpdateRepositoryPropertiesOptions{Value: &azcontainerregistry.RepositoryWriteableProperties{
		CanWrite: to.Ptr(false),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - 'CanWrite' property: %t\n", *res.ContainerRepositoryProperties.ChangeableAttributes.CanWrite)
}

func ExampleClient_UpdateTagProperties() {
	res, err := client.UpdateTagProperties(context.TODO(), "nanoserver", "4.7.2-20180905-nanoserver-1803", &azcontainerregistry.ClientUpdateTagPropertiesOptions{
		Value: &azcontainerregistry.TagWriteableProperties{
			CanWrite: to.Ptr(false),
		}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository namoserver - tag 4.7.2-20180905-nanoserver-1803 - 'CanWrite' property: %t\n", *res.Tag.ChangeableAttributes.CanWrite)
}
