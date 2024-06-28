//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"encoding/hex"
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
// [Azure Key Vault documentation]: https://learn.microsoft.com/azure/key-vault/keys/how-to-configure-key-rotation
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

func toByteSlice(s string) []byte {
	if len(s)%2 == 1 {
		s = fmt.Sprintf("0%s", s)
	}
	ret, err := hex.DecodeString(s)
	if err != nil {
		// TODO: handle error
	}
	return ret
}

func ExampleClient_ImportKey() {
	// Example uploading a RSA key

	// example RSA key
	params := azkeys.ImportKeyParameters{
		Key: &azkeys.JSONWebKey{
			Kty: to.Ptr(azkeys.KeyTypeRSA),
			N:   toByteSlice("a0914d00234ac683b21b4c15d5bed887bdc959c2e57af54ae734e8f00720d775d275e455207e3784ceeb60a50a4655dd72a7a94d271e8ee8f7959a669ca6e775bf0e23badae991b4529d978528b4bd90521d32dd2656796ba82b6bbfc7668c8f5eeb5053747fd199319d29a8440d08f4412d527ff9311eda71825920b47b1c46b11ab3e91d7316407e89c7f340f7b85a34042ce51743b27d4718403d34c7b438af6181be05e4d11eb985d38253d7fe9bf53fc2f1b002d22d2d793fa79a504b6ab42d0492804d7071d727a06cf3a8893aa542b1503f832b296371b6707d4dc6e372f8fe67d8ded1c908fde45ce03bc086a71487fa75e43aa0e0679aa0d20efe35"),
			E:   toByteSlice("10001"),
			D:   toByteSlice("627c7d24668148fe2252c7fa649ea8a5a9ed44d75c766cda42b29b660e99404f0e862d4561a6c95af6a83d213e0a2244b03cd28576473215073785fb067f015da19084ade9f475e08b040a9a2c7ba00253bb8125508c9df140b75161d266be347a5e0f6900fe1d8bbf78ccc25eeb37e0c9d188d6e1fc15169ba4fe12276193d77790d2326928bd60d0d01d6ead8d6ac4861abadceec95358fd6689c50a1671a4a936d2376440a41445501da4e74bfb98f823bd19c45b94eb01d98fc0d2f284507f018ebd929b8180dbe6381fdd434bffb7800aaabdd973d55f9eaf9bb88a6ea7b28c2a80231e72de1ad244826d665582c2362761019de2e9f10cb8bcc2625649"),
			P:   toByteSlice("00d1deac8d68ddd2c1fd52d5999655b2cf1565260de5269e43fd2a85f39280e1708ffff0682166cb6106ee5ea5e9ffd9f98d0becc9ff2cda2febc97259215ad84b9051e563e14a051dce438bc6541a24ac4f014cf9732d36ebfc1e61a00d82cbe412090f7793cfbd4b7605be133dfc3991f7e1bed5786f337de5036fc1e2df4cf3"),
			Q:   toByteSlice("00c3dc66b641a9b73cd833bc439cd34fc6574465ab5b7e8a92d32595a224d56d911e74624225b48c15a670282a51c40d1dad4bc2e9a3c8dab0c76f10052dfb053bc6ed42c65288a8e8bace7a8881184323f94d7db17ea6dfba651218f931a93b8f738f3d8fd3f6ba218d35b96861a0f584b0ab88ddcf446b9815f4d287d83a3237"),
			DP:  toByteSlice("00c9a159be7265cbbabc9afcc4967eb74fe58a4c4945431902d1142da599b760e03838f8cbd26b64324fea6bdc9338503f459793636e59b5361d1e6951e08ddb089e1b507be952a81fbeaf7e76890ea4f536e25505c3f648b1e88377dfc19b4c304e738dfca07211b792286a392a704d0f444c0a802539110b7f1f121c00cff0a9"),
			DQ:  toByteSlice("00a0bd4c0a3d9f64436a082374b5caf2488bac1568696153a6a5e4cd85d186db31e2f58f024c617d29f37b4e6b54c97a1e25efec59c4d1fd3061ac33509ce8cae5c11f4cd2e83f41a8264f785e78dc0996076ee23dfdfc43d67c463afaa0180c4a718357f9a6f270d542479a0f213870e661fb950abca4a14ca290570ba7983347"),
			QI:  toByteSlice("009fe7ae42e92bc04fcd5780464bd21d0c8ac0c599f9af020fde6ab0a7e7d1d39902f5d8fb6c614184c4c1b103fb46e94cd10a6c8a40f9991a1f28269f326435b6c50276fda6493353c650a833f724d80c7d522ba16c79f0eb61f672736b68fb8be3243d10943c4ab7028d09e76cfb5892222e38bc4d35585bf35a88cd68c73b07"),
		},
	}

	// if a key with the same name already exists, a new version of that key is created
	resp, err := client.ImportKey(context.TODO(), "key-name", params, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Key %s successfully imported.", resp.Key.KID.Name())
}
