//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcrypto_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/azcrypto"
)

var client *azcrypto.Client

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err = azcrypto.NewClient("https://<my-keyvault-url>.vault.azure.net/keys/<my-key>", cred, nil)
	if err != nil {
		panic(err)
	}
}
