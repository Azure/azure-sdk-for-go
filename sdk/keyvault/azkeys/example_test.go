//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

var client azkeys.Client

func ExampleNewClient() {
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewClient(vaultURL, cred, nil)
	_ = client
}

func ExampleClient_CreateKey_rsa() {
	params := azkeys.CreateKeyParameters{
		KeySize: to.Ptr(int32(2048)),
		Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSA),
	}
	resp, err := client.CreateKey(context.TODO(), "key-name", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}

func ExampleClient_CreateKey_ec() {
	params := azkeys.CreateKeyParameters{
		Curve: to.Ptr(azkeys.JSONWebKeyCurveNameP256K),
		Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC),
	}
	resp, err := client.CreateKey(context.TODO(), "key-name", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}

func ExampleClient_DeleteKey() {
	// DeleteKey returns when Key Vault has begun deleting the key. That can take several
	// seconds to complete, so it may be necessary to wait before performing other operations
	// on the deleted key.
	resp, err := client.DeleteKey(context.TODO(), "key-name", nil)
	if err != nil {
		// TODO: handle error
	}

	// In a soft-delete enabled vault, deleted keys can be recovered until they're purged (permanently deleted).
	fmt.Printf("Key will be purged at %v", resp.ScheduledPurgeDate)
}

func ExampleClient_PurgeDeletedKey() {
	// this loop purges all the deleted keys in the vault
	pager := client.NewListDeletedKeysPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, key := range page.Value {
			_, err := client.PurgeDeletedKey(context.TODO(), key.KID.Name(), nil)
			if err != nil {
				// TODO: handle error
			}
		}
	}
}

func ExampleClient_GetKey() {
	// passing an empty string for the version parameter gets the latest version of the key
	resp, err := client.GetKey(context.TODO(), "key-name", "", nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}

func ExampleClient_UpdateKey() {
	params := azkeys.UpdateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			Expires: to.Ptr(time.Now().Add(48 * time.Hour)),
		},
		Tags: map[string]*string{"expiraton-extended": to.Ptr("true")},
	}
	// passing an empty string for the version parameter updates the latest version of the key
	updateResp, err := client.UpdateKey(context.TODO(), "key-name", "", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Enabled key %s", *updateResp.Key.KID)
}

func ExampleClient_NewListKeysPager() {
	pager := client.NewListKeysPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, key := range resp.Value {
			fmt.Println(*key.KID)
		}
	}
}
