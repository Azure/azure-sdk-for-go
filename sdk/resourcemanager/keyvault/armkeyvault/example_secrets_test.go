//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

func ExampleSecretsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewSecretsClient("<subscription ID>", cred, nil)
	resp, err := client.CreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<vault name>",
		"<secret name>",
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Attributes: armkeyvault.Attributes{
						Enabled: to.BoolPtr(true),
					},
				},
				Value: to.StringPtr("<serect value>"),
			}}, nil)
	if err != nil {
		log.Fatalf("failed to create the key: %v", err)
	}
	log.Printf("Secret ID: %v\n", *resp.Secret.ID)
}
