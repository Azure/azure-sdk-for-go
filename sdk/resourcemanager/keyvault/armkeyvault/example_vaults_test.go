//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

func ExampleVaultsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("<subscription ID>", cred, nil)
	poll, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<vault name>",
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr("<Azure location>"),
			Properties: &armkeyvault.VaultProperties{
				TenantID: to.StringPtr("<tenant ID>"),
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{},
			},
		}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poll.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create a vault: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsClient_BeginCreateOrUpdate_withAccessPolicies() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("<subscription ID>", cred, nil)
	poll, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<vault name>",
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr("<Azure location>"),
			Properties: &armkeyvault.VaultProperties{
				TenantID: to.StringPtr("<tenant ID>"),
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						ObjectID: to.StringPtr("<user, service principal or security group object ID>"),
						TenantID: to.StringPtr("<tenant ID>"),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
						},
					},
				},
			},
		}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poll.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create a vault: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsClient_BeginCreateOrUpdate_forDeployment() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("<subscription ID>", cred, nil)
	poll, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<vault name>",
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr("<Azure location>"),
			Properties: &armkeyvault.VaultProperties{
				TenantID:                     to.StringPtr("<tenant ID>"),
				EnabledForDeployment:         to.BoolPtr(true),
				EnabledForTemplateDeployment: to.BoolPtr(true),
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						ObjectID: to.StringPtr("<user, service principal or security group object ID>"),
						TenantID: to.StringPtr("<tenant ID>"),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
						},
					},
				},
			},
		}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poll.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create a vault: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<vault name>", nil)
	if err != nil {
		log.Fatalf("failed to get the vault: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("<subscription ID>", cred, nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceListResult.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range resp.ResourceListResult.Value {
			log.Printf("vault: %v", *val.ID)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleVaultsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient("<subscription ID>", cred, nil)
	_, err = client.Delete(context.Background(), "<resource group name>", "<vault name>", nil)
	if err != nil {
		log.Fatalf("failed to delete the vault: %v", err)
	}
}
