// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"fmt"
	"testing"

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

func ExampleServiceClient_ListContainers() {
	cred := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	client, err := NewServiceClient(endpoint, cred, nil)
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

func Test_MSI_ListContainers(t *testing.T) {
	cred := azidentity.NewManagedIdentityCredential("expected_client", nil)
	client, err := NewServiceClient(endpoint, cred, nil)
	if err != nil {
		t.Fatalf("Error: " + err.Error())
	}
	iter := client.ListContainers(nil)
	for {
		p, err := iter.NextPage(context.Background())
		if errors.Is(err, azcore.IterationDone) {
			break
		} else if err != nil {
			t.Fatalf("Error: " + err.Error())
		}
		for _, i := range p.ContainerItems {
			fmt.Println(i.Name)
		}
	}
}
