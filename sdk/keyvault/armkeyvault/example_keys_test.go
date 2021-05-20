// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleKeysClient_CreateIfNotExist() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewKeysClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.CreateIfNotExist(
		context.Background(),
		"<resource group name>",
		"<vault name>",
		"<key name>",
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Attributes: armkeyvault.Attributes{
						Enabled: to.BoolPtr(true),
					},
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			}}, nil)
	if err != nil {
		log.Fatalf("failed to create the key: %v", err)
	}
	log.Printf("key ID: %v\n", *resp.Key.ID)
}
