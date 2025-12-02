// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
)

func Example_deleteImages() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azcontainerregistry.NewClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	repositoryPager := client.NewListRepositoriesPager(nil)
	for repositoryPager.More() {
		repositoryPage, err := repositoryPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance repository page: %v", err)
		}
		for _, r := range repositoryPage.Names {
			manifestPager := client.NewListManifestsPager(*r, &azcontainerregistry.ClientListManifestsOptions{
				OrderBy: to.Ptr(azcontainerregistry.ArtifactManifestOrderByLastUpdatedOnDescending),
			})
			for manifestPager.More() {
				manifestPage, err := manifestPager.NextPage(ctx)
				if err != nil {
					log.Fatalf("failed to advance manifest page: %v", err)
				}
				imagesToKeep := 3
				for i, m := range manifestPage.Attributes {
					if i >= imagesToKeep {
						for _, t := range m.Tags {
							fmt.Printf("delete tag from image: %s", *t)
							_, err := client.DeleteTag(ctx, *r, *t, nil)
							if err != nil {
								log.Fatalf("failed to delete tag: %v", err)
							}
						}
						_, err := client.DeleteManifest(ctx, *r, *m.Digest, nil)
						if err != nil {
							log.Fatalf("failed to delete manifest: %v", err)
						}
						fmt.Printf("delete image with digest: %s", *m.Digest)
					}
				}
			}
		}
	}
}
