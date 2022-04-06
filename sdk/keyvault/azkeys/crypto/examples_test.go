//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto_test

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
)

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}
	_ = client // do something with client
}

func ExampleClient_Encrypt() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	encryptResponse, err := client.Encrypt(context.TODO(), crypto.EncryptionAlgRSAOAEP, []byte("plaintext"), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptResponse.Ciphertext)
}

func ExampleClient_Decrypt() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	encryptResponse, err := client.Encrypt(context.TODO(), crypto.EncryptionAlgRSAOAEP, []byte("plaintext"), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptResponse.Ciphertext)

	decryptResponse, err := client.Decrypt(context.TODO(), crypto.EncryptionAlgRSAOAEP, encryptResponse.Ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(decryptResponse.Plaintext)
}

func ExampleClient_WrapKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

	// Wrap
	wrapResp, err := client.WrapKey(context.TODO(), crypto.WrapAlgRSAOAEP, keyBytes, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(wrapResp.EncryptedKey)
}

func ExampleClient_UnwrapKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

	// Wrap
	wrapResp, err := client.WrapKey(context.TODO(), crypto.WrapAlgRSAOAEP, keyBytes, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(wrapResp.EncryptedKey)

	// Unwrap
	unwrapResp, err := client.UnwrapKey(context.TODO(), crypto.WrapAlgRSAOAEP, wrapResp.EncryptedKey, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(unwrapResp.Key)
}

func ExampleClient_Sign() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte("plaintext"))
	if err != nil {
		panic(err)
	}
	digest := hasher.Sum(nil)

	signResponse, err := client.Sign(context.TODO(), crypto.SignatureAlgRS256, digest, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(signResponse.Signature)
}

func ExampleClient_Verify() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte("plaintext"))
	if err != nil {
		panic(err)
	}
	digest := hasher.Sum(nil)

	signResponse, err := client.Sign(context.TODO(), crypto.SignatureAlgRS256, digest, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(signResponse.Signature)

	verifyResponse, err := client.Verify(context.TODO(), crypto.SignatureAlgRS256, digest, signResponse.Signature, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*verifyResponse.IsValid)
}
