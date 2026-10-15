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

// CreateKey can register a reference to key material that lives in an external HSM or KMS by
// setting [KeyAttributes.ExternalKey]. Key Vault stores only the reference; cryptographic
// operations are delegated to the external system. External keys are mutually exclusive with
// [CreateKeyParameters.Kty], so leave Kty unset when supplying an ExternalKey.
func ExampleClient_CreateKey_externalKey() {
	params := azkeys.CreateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			// ID identifies the key material in the external HSM/KMS. Allowed characters are
			// [a-zA-Z0-9-] and the maximum length is 64.
			ExternalKey: &azkeys.ExternalKey{
				ID: to.Ptr("external-key-id"),
			},
		},
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

// SecureWrapKey creates a new 256-bit AES key inside a trusted execution environment (TEE) and wraps
// it with the named key encryption key, returning the wrapped key material. Securely wrapping and
// unwrapping keys requires Microsoft Azure Attestation (MAA) to attest the TEE and is supported by
// Managed HSM.
//
// The wrapping key must be created specifically for secure wrap/unwrap: it must grant the
// "secureWrapKey" and "secureUnwrapKey" key operations, must have a release policy, and must
// NOT be exportable.
// See [Azure Key Vault documentation] for more information.
//
// [Azure Key Vault documentation]: https://learn.microsoft.com/azure/key-vault/keys/about-keys
func ExampleClient_SecureWrapKey() {
	// Create an RSA-HSM wrapping key configured for secure wrap/unwrap. encodedReleasePolicy is the
	// JSON-encoded key release policy bytes (see Azure Key Vault docs for the schema).
	var encodedReleasePolicy []byte
	createParams := azkeys.CreateKeyParameters{
		Kty: to.Ptr(azkeys.KeyTypeRSAHSM),
		KeyOps: to.SliceOfPtrs(
			azkeys.KeyOperationSecureWrapKey,
			azkeys.KeyOperationSecureUnwrapKey,
		),
		ReleasePolicy: &azkeys.KeyReleasePolicy{
			EncodedPolicy: encodedReleasePolicy,
			Immutable:     to.Ptr(true),
		},
	}
	if _, err := client.CreateKey(context.TODO(), "key-name", createParams, nil); err != nil {
		// TODO: handle error
	}

	params := azkeys.SecureKeyWrapOperationParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256),
	}
	// passing an empty string for the version parameter uses the latest version of the key
	resp, err := client.SecureWrapKey(context.TODO(), "key-name", "", params, nil)
	if err != nil {
		// TODO: handle error
	}
	// resp.Value contains the wrapped key material, which SecureUnwrapKey can later reverse
	fmt.Printf("Wrapped a key with %s", *resp.Kid)
}

// SecureUnwrapKey reverses [Client.SecureWrapKey], decrypting previously wrapped key material inside a
// trusted execution environment (TEE). It requires an attestation token from Microsoft Azure
// Attestation (MAA) so the service can verify the TEE before unwrapping. The wrapping key must have
// been created with the "secureWrapKey"/"secureUnwrapKey" ops and a release policy (see
// [ExampleClient_SecureWrapKey] for the required CreateKey shape).
func ExampleClient_SecureUnwrapKey() {
	// targetAttestationToken is an attestation assertion obtained from your Microsoft Azure
	// Attestation (MAA) instance for the target TEE.
	var targetAttestationToken string
	// wrappedKey is the key material returned by a prior SecureWrapKey call.
	var wrappedKey []byte

	params := azkeys.SecureKeyUnWrapOperationParameters{
		Algorithm:              to.Ptr(azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256),
		TargetAttestationToken: to.Ptr(targetAttestationToken),
		Value:                  wrappedKey,
	}
	// passing an empty string for the version parameter uses the latest version of the key
	resp, err := client.SecureUnwrapKey(context.TODO(), "key-name", "", params, nil)
	if err != nil {
		// TODO: handle error
	}
	// resp.Value contains the unwrapped key material
	fmt.Printf("Unwrapped a key with %s", *resp.Kid)
}
