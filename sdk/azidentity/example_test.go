//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleNewOnBehalfOfCredentialWithCertificate() {
	data, err := os.ReadFile(certPath)
	if err != nil {
		// TODO: handle error
	}

	// NewOnBehalfOfCredentialFromCertificate requires at least one *x509.Certificate, and a crypto.PrivateKey.
	// ParseCertificates returns these given certificate data in PEM or PKCS12 format. It handles common
	// scenarios but has limitations, for example it doesn't load PEM encrypted private keys.
	certs, key, err := azidentity.ParseCertificates(data, nil)
	if err != nil {
		// TODO: handle error
	}

	cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, nil)
	if err != nil {
		// TODO: handle error
	}

	// Output:
}

func ExampleNewClientCertificateCredential() {
	data, err := os.ReadFile(certPath)
	handleError(err)

	// NewClientCertificateCredential requires at least one *x509.Certificate, and a crypto.PrivateKey.
	// ParseCertificates returns these given certificate data in PEM or PKCS12 format. It handles common scenarios
	// but has limitations, for example it doesn't load PEM encrypted private keys.
	certs, key, err := azidentity.ParseCertificates(data, nil)
	handleError(err)

	cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, nil)
	handleError(err)

	// Output:
}

func ExampleNewManagedIdentityCredential_userAssigned() {
	// select a user assigned identity with its client ID...
	clientID := azidentity.ClientID("abcd1234-...")
	opts := azidentity.ManagedIdentityCredentialOptions{ID: clientID}
	cred, err = azidentity.NewManagedIdentityCredential(&opts)
	handleError(err)

	// ...or its resource ID
	resourceID := azidentity.ResourceID("/subscriptions/...")
	opts = azidentity.ManagedIdentityCredentialOptions{ID: resourceID}
	cred, err = azidentity.NewManagedIdentityCredential(&opts)
	handleError(err)
}
