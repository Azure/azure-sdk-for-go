// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/keyvault/2019-09-01/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleVaultsOperations_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
				AccessPolicies: &[]armkeyvault.AccessPolicyEntry{},
			},
		}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poll.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsOperations_BeginCreateOrUpdate_withAccessPolicies() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
				AccessPolicies: &[]armkeyvault.AccessPolicyEntry{
					{
						ObjectID: to.StringPtr("<user, service principal or security group object ID>"),
						TenantID: to.StringPtr("<tenant ID>"),
						Permissions: &armkeyvault.Permissions{
							Keys: &[]armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet,
								armkeyvault.KeyPermissionsList,
								armkeyvault.KeyPermissionsCreate,
							},
							Secrets: &[]armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet,
								armkeyvault.SecretPermissionsList,
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsOperations_BeginCreateOrUpdate_forDeployment() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
				AccessPolicies: &[]armkeyvault.AccessPolicyEntry{
					{
						ObjectID: to.StringPtr("<user, service principal or security group object ID>"),
						TenantID: to.StringPtr("<tenant ID>"),
						Permissions: &armkeyvault.Permissions{
							Keys: &[]armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet,
								armkeyvault.KeyPermissionsList,
								armkeyvault.KeyPermissionsCreate,
							},
							Secrets: &[]armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet,
								armkeyvault.SecretPermissionsList,
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsOperations_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<vault name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("vault ID: %v\n", *resp.Vault.ID)
}

func ExampleVaultsOperations_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(*resp.ResourceListResult.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range *resp.ResourceListResult.Value {
			log.Printf("vault: %v", *val.ID)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleVaultsOperations_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armkeyvault.NewVaultsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Delete(context.Background(), "<resource group name>", "<vault name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("Vault deletion status code: %v\n", resp.StatusCode)
}
