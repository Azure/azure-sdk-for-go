//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"fmt"
	"os"

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

	resp, err := client.CreateECKey(context.TODO(), "new-rsa-key", &azkeys.CreateECKeyOptions{CurveName: azkeys.P256.ToPtr()})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.Key.ID)
	fmt.Println(*resp.Key.KeyType)
}
