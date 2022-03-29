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

func ExampleNewClient() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}
	_ = client // do something with client
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

	resp, err := client.CreateRSAKey(context.TODO(), "new-rsa-key", &azkeys.CreateRSAKeyOptions{Size: to.Int32Ptr(2048)})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.JSONWebKey.ID)
	fmt.Println(*resp.JSONWebKey.KeyType)
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

	resp, err := client.CreateECKey(context.TODO(), "new-rsa-key", &azkeys.CreateECKeyOptions{CurveName: azkeys.CurveNameP256.ToPtr()})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.JSONWebKey.ID)
	fmt.Println(*resp.JSONWebKey.KeyType)
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
	fmt.Println(*resp.JSONWebKey.ID)
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
		Properties: &azkeys.Properties{
			RecoveryLevel: azkeys.DeletionRecoveryLevelCustomizedRecoverablePurgeable.ToPtr(),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("RecoverLevel: %s\tTag1: %s\n", *resp.Properties.RecoveryLevel, resp.Tags["Tag1"])
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

func ExampleClient_ListPropertiesOfKeys() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.ListPropertiesOfKeys(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, key := range resp.Keys {
			fmt.Println(*key.ID)
		}
	}
}

func ExampleClient_UpdateKeyRotationPolicy() {
	vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultUrl, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.UpdateKeyRotationPolicy(context.TODO(), "key-to-update", &azkeys.UpdateKeyRotationPolicyOptions{
		Attributes: &azkeys.RotationPolicyAttributes{
			ExpiryTime: to.StringPtr("P90D"),
		},
		LifetimeActions: []*azkeys.LifetimeActions{
			{
				Action: &azkeys.LifetimeActionsType{
					Type: azkeys.ActionTypeNotify.ToPtr(),
				},
				Trigger: &azkeys.LifetimeActionsTrigger{
					TimeBeforeExpiry: to.StringPtr("P30D"),
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated key rotation policy for: ", *resp.ID)

	//
	_, err = client.RotateKey(context.TODO(), "key-to-rotate", nil)
	if err != nil {
		panic(err)
	}
}
