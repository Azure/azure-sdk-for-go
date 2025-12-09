// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
)

func Example_listTagsWithAnonymousAccess() {
	client, err := azcontainerregistry.NewClient("<your Container Registry's endpoint URL>", nil, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	pager := client.NewListTagsPager("library/hello-world", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Tags {
			fmt.Printf("tag: %s\n", *v.Name)
		}
	}
}
