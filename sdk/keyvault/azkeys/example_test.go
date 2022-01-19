//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

var client *azkeys.Client

func ExampleNewClient() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_CreateRSAKey() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.CreateRSAKey(context.TODO(), "new-rsa-key", &azkeys.CreateRSAKeyOptions{KeySize: to.Int32Ptr(2048)})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Key.ID)
	fmt.Println(*resp.Key.KeyType)
}

func ExampleClient_CreateECKey() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.CreateECKey(context.TODO(), "new-rsa-key", &azkeys.CreateECKeyOptions{CurveName: azkeys.JSONWebKeyCurveNameP256.ToPtr()})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Key.ID)
	fmt.Println(*resp.Key.KeyType)
}

func ExampleClient_GetKey() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetKey(context.TODO(), "key-to-retrieve", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Key.ID)
}

func ExampleClient_UpdateKeyProperties() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.UpdateKeyProperties(context.TODO(), "key-to-update", &azkeys.UpdateKeyPropertiesOptions{
		Tags: map[string]string{
			"Tag1": "val1",
		},
		KeyAttributes: &azkeys.KeyAttributes{
			RecoveryLevel: azkeys.DeletionRecoveryLevelCustomizedRecoverablePurgeable.ToPtr(),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Attributes.RecoveryLevel, resp.Tags["Tag1"])
}

func ExampleClient_BeginDeleteKey() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteKey(context.TODO(), "key-to-delete", nil)
	if err != nil {
		panic(err)
	}
	pollResp, err := resp.PollUntilDone(context.TODO(), 1*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully deleted key %s", *pollResp.Key.ID)
}

func ExampleClient_ListKeys() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.ListKeys(nil)
	for pager.NextPage(context.TODO()) {
		for _, key := range pager.PageResponse().Keys {
			fmt.Println(*key.KID)
		}
	}

	if pager.Err() != nil {
		panic(pager.Err())
	}
}
