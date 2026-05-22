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

func Example_setArtifactProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client, err := azcontainerregistry.NewClient("<your Container Registry's endpoint URL>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()
	res, err := client.UpdateTagProperties(ctx, "library/hello-world", "latest", &azcontainerregistry.ClientUpdateTagPropertiesOptions{
		Value: &azcontainerregistry.TagWriteableProperties{
			CanWrite:  to.Ptr(false),
			CanDelete: to.Ptr(false),
		}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	fmt.Printf("repository library/hello-world - tag latest: 'CanWrite' property: %t, 'CanDelete' property: %t\n", *res.Tag.ChangeableAttributes.CanWrite, *res.Tag.ChangeableAttributes.CanDelete)
}
