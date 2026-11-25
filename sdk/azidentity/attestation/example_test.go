// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package attestation_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func Example() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	client, err := azsecrets.NewClient("https://<vault-name>.vault.azure.net/", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_, err = client.GetSecret(context.TODO(), "<secret-name>", "", nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func Example_managedIdentityOptOut() {
	cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		DisableMTLSPoP: true,
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_, err = azsecrets.NewClient("https://<vault-name>.vault.azure.net/", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}
