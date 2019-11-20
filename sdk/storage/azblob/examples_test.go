// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	endpoint = "https://<endpoint>.blob.core.windows.net/"
)

const (
	tenantID     = "<tenant>"
	clientID     = "<client>"
	clientSecret = "<secret>"
)

func ExampleServiceClientListContainers() {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	client, err := NewServiceClient(endpoint,
		cred,
		azcore.PipelineOptions{})
	if err != nil {
		panic(err)
	}
	iter := client.ListContainers(nil)
	for {
		p, err := iter.NextPage(context.Background())
		if errors.Is(err, azcore.IterationDone) {
			break
		} else if err != nil {
			panic(err)
		}
		for _, i := range p.ContainerItems {
			fmt.Println(i.Name)
		}
	}
}
