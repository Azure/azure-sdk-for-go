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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

var client azkeys.Client

func ExampleNewClient() {
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azkeys.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_CreateKey_rsa() {
	params := azkeys.CreateKeyParameters{
		KeySize: to.Ptr(int32(2048)),
		Kty:     to.Ptr(azkeys.KeyTypeRSA),
	}
	// if a key with the same name already exists, a new version of that key is created
	resp, err := client.CreateKey(context.TODO(), "key-name", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.Key.KID)
}

func ExampleClient_CreateKey_ec() {
	params := azkeys.CreateKeyParameters{
		Curve: to.Ptr(azkeys.CurveNameP256K),
		Kty:   to.Ptr(azkeys.KeyTypeEC),
	}
	// if a key with the same name already exists, a new version of that key is created
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
	pager := client.NewListDeletedKeyPropertiesPager(nil)
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

// UpdateKey updates the properties of a key previously stored in the key vault
func ExampleClient_UpdateKey() {
	params := azkeys.UpdateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			Expires: to.Ptr(time.Now().Add(48 * time.Hour)),
		},
		// Key Vault doesn't interpret tags. The keys and values are up to your application.
		Tags: map[string]*string{"expiration-extended": to.Ptr("true")},
	}
	// passing an empty string for the version parameter updates the latest version of the key
	updateResp, err := client.UpdateKey(context.TODO(), "key-name", "", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Enabled key %s", *updateResp.Key.KID)
}

// UpdateKeyRotationPolicy allows you to configure automatic key rotation for a key by specifying a rotation policy, and
// [Client.RotateKey] allows you to rotate a key on demand. See [Azure Key Vault documentation] for more information about key
// rotation.
//
// [Azure Key Vault documentation]: https://docs.microsoft.com/azure/key-vault/keys/how-to-configure-key-rotation
func ExampleClient_UpdateKeyRotationPolicy() {
	// this policy rotates the key every 18 months
	policy := azkeys.KeyRotationPolicy{
		LifetimeActions: []*azkeys.LifetimeAction{
			{
				Action: &azkeys.LifetimeActionType{
					Type: to.Ptr(azkeys.KeyRotationPolicyActionRotate),
				},
				Trigger: &azkeys.LifetimeActionTrigger{
					TimeAfterCreate: to.Ptr("P18M"),
				},
			},
		},
	}
	resp, err := client.UpdateKeyRotationPolicy(context.TODO(), "key-name", policy, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Updated key rotation policy at: %v", resp.Attributes.Updated)
}

func ExampleClient_NewListKeyPropertiesPager() {
	pager := client.NewListKeyPropertiesPager(nil)
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
