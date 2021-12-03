//go:build go1.16
// +build go1.16

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

var client *crypto.Client

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_Encrypt() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	encryptResponse, err := client.Encrypt(context.TODO(), crypto.AlgorithmRSAOAEP, []byte("plaintext"), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptResponse.Result)
}

func ExampleClient_Decrypt() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	encryptResponse, err := client.Encrypt(context.TODO(), crypto.AlgorithmRSAOAEP, []byte("plaintext"), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptResponse.Result)

	decryptResponse, err := client.Decrypt(context.TODO(), crypto.AlgorithmRSAOAEP, encryptResponse.Result, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(decryptResponse.Result)
}

func ExampleClient_WrapKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

	// Wrap
	wrapResp, err := client.WrapKey(context.TODO(), crypto.RSAOAEP, keyBytes, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(wrapResp.Result)
}

func ExampleClient_UnwrapKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

	// Wrap
	wrapResp, err := client.WrapKey(context.TODO(), crypto.RSAOAEP, keyBytes, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(wrapResp.Result)

	// Unwrap
	unwrapResp, err := client.UnwrapKey(context.TODO(), crypto.RSAOAEP, wrapResp.Result, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(unwrapResp.Result)
}

func ExampleClient_Sign() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte("plaintext"))
	if err != nil {
		panic(err)
	}
	digest := hasher.Sum(nil)

	signResponse, err := client.Sign(context.TODO(), crypto.RS256, digest, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(signResponse.Result)
}

func ExampleClient_Verify() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = crypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte("plaintext"))
	if err != nil {
		panic(err)
	}
	digest := hasher.Sum(nil)

	signResponse, err := client.Sign(context.TODO(), crypto.RS256, digest, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(signResponse.Result)

	verifyResponse, err := client.Verify(context.TODO(), crypto.RS256, digest, signResponse.Result, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*verifyResponse.IsValid)
}
