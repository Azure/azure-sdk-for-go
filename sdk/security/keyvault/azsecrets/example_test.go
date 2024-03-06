//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

var client azsecrets.Client

func ExampleNewClient() {
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_SetSecret() {
	name := "mySecret"
	value := "mySecretValue"
	// If no secret with the given name exists, Key Vault creates a new secret with that name and the given value.
	// If the given name is in use, Key Vault creates a new version of that secret, with the given value.
	resp, err := client.SetSecret(context.TODO(), name, azsecrets.SetSecretParameters{Value: &value}, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Set secret %s", resp.ID.Name())
}

func ExampleClient_GetSecret() {
	// an empty string gets the latest version of the secret
	version := ""
	resp, err := client.GetSecret(context.TODO(), "mySecretName", version, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Secret Name: %s\tSecret Value: %s", resp.ID.Name(), *resp.Value)
}

func ExampleClient_DeleteSecret() {
	// DeleteSecret returns when Key Vault has begun deleting the secret. That can take several
	// seconds to complete, so it may be necessary to wait before performing other operations
	// on the deleted secret.
	resp, err := client.DeleteSecret(context.TODO(), "secretToDelete", nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("deleted secret", resp.ID.Name())
}

// List pages don't include secret values. Use [Client.GetSecret] to retrieve secret values.
func ExampleClient_NewListSecretPropertiesPager() {
	pager := client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, secret := range page.Value {
			fmt.Printf("Secret Name: %s\tSecret Tags: %v\n", secret.ID.Name(), secret.Tags)
		}
	}
}

func ExampleClient_BackupSecret() {
	backup, err := client.BackupSecret(context.TODO(), "mySecret", nil)
	if err != nil {
		// TODO: handle error
	}

	restoreResp, err := client.RestoreSecret(context.TODO(), azsecrets.RestoreSecretParameters{SecretBackup: backup.Value}, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Restored ID %s\n", *restoreResp.ID)
}

func ExampleClient_RestoreSecret() {
	backup, err := client.BackupSecret(context.TODO(), "mySecret", nil)
	if err != nil {
		// TODO: handle error
	}

	restoreResp, err := client.RestoreSecret(context.TODO(), azsecrets.RestoreSecretParameters{SecretBackup: backup.Value}, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Restored ID %s\n", *restoreResp.ID)
}

func ExampleClient_PurgeDeletedSecret() {
	// this loop purges all the deleted secrets in the vault
	pager := client.NewListDeletedSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, secret := range page.Value {
			_, err := client.PurgeDeletedSecret(context.TODO(), secret.ID.Name(), nil)
			if err != nil {
				// TODO: handle error
			}
		}
	}
}

func ExampleClient_RecoverDeletedSecret() {
	resp, err := client.RecoverDeletedSecret(context.TODO(), "myDeletedSecret", nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println("recovered deleted secret", resp.ID.Name())
}

// UpdateSecret updates a secret's metadata. It can't change the secret's value; use [Client.SetSecret] to set a secret's value.
func ExampleClient_UpdateSecretProperties() {
	updateParams := azsecrets.UpdateSecretPropertiesParameters{
		SecretAttributes: &azsecrets.SecretAttributes{
			Expires: to.Ptr(time.Now().Add(48 * time.Hour)),
		},
		// Key Vault doesn't interpret tags. The keys and values are up to your application.
		Tags: map[string]*string{"expiration-extended": to.Ptr("true")},
	}
	// an empty version updates the latest version of the secret
	version := ""
	resp, err := client.UpdateSecretProperties(context.Background(), "mySecretName", version, updateParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println("Updated secret", resp.ID.Name())
}
