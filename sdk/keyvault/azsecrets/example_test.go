//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func ExampleNewClient() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}
	_ = client
}

func ExampleClient_SetSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	secret := "mySecret"
	value := "mySecretValue"

	resp, err := client.SetSecret(context.TODO(), secret, value, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Set secret %s", *resp.ID)
}

func ExampleClient_GetSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetSecret(context.TODO(), "mySecretName", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Secret Name: %s\tSecret Value: %s", *resp.ID, *resp.Value)
}

func ExampleClient_BeginDeleteSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteSecret(context.TODO(), "secretToDelete", nil)
	if err != nil {
		panic(err)
	}
	// This is optional if you don't care when the secret is deleted
	_, err = resp.PollUntilDone(context.TODO(), 250*time.Millisecond)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_ListSecrets() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.ListSecrets(nil)
	for pager.NextPage(context.TODO()) {
		for _, v := range pager.PageResponse().Secrets {
			fmt.Printf("Secret Name: %s\tSecret Tags: %v\n", *v.ID, v.Tags)
		}
	}

	if pager.Err() != nil {
		panic(pager.Err())
	}
}

func ExampleClient_RestoreSecretBackup() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BackupSecret(context.TODO(), "mySecret", nil)
	if err != nil {
		panic(err)
	}

	restoreResp, err := client.RestoreSecretBackup(context.TODO(), resp.Value, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Restored ID %s\n", *restoreResp.ID)
}

func ExampleClient_BackupSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BackupSecret(context.TODO(), "mySecret", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup secret byte value: %v", resp.Value)
}

func ExampleClient_PurgeDeletedSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteSecret(context.TODO(), "mySecretName", nil)
	if err != nil {
		panic(err)
	}
	_, err = resp.PollUntilDone(context.TODO(), 250*time.Millisecond)
	if err != nil {
		panic(err)
	}

	_, err = client.PurgeDeletedSecret(context.TODO(), "mySecretName", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_BeginRecoverDeletedSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginRecoverDeletedSecret(context.TODO(), "myDeletedSecret", nil)
	if err != nil {
		panic(err)
	}
	_, err = resp.PollUntilDone(context.TODO(), 250*time.Millisecond)
	if err != nil {
		panic(err)
	}
}
