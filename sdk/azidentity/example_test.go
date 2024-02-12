//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Credentials that require user interaction such as [InteractiveBrowserCredential] and [DeviceCodeCredential]
// can optionally return this error instead of automatically prompting for user interaction. This allows applications
// to decide when to request user interaction. This example shows how to handle the error and authenticate a user
// interactively. It shows [InteractiveBrowserCredential] but the same pattern applies to [DeviceCodeCredential].
func ExampleAuthenticationRequiredError() {
	cred, err := azidentity.NewInteractiveBrowserCredential(
		&azidentity.InteractiveBrowserCredentialOptions{
			// This option is useful only for applications that need to control when to prompt users to
			// authenticate. If the timing of user interaction isn't important, don't set this option.
			DisableAutomaticAuthentication: true,
		},
	)
	if err != nil {
		// TODO: handle error
	}
	// this could be any client that authenticates with an azidentity credential
	client, err := newServiceClient(cred)
	if err != nil {
		// TODO: handle error
	}
	err = client.Method()
	if err != nil {
		var are *azidentity.AuthenticationRequiredError
		if errors.As(err, &are) {
			// The client requested a token and the credential requires user interaction. Whenever it's convenient
			// for the application, call Authenticate to prompt the user. Pass the error's TokenRequestOptions to
			// request a token with the parameters the client specified.
			_, err = cred.Authenticate(context.TODO(), &are.TokenRequestOptions)
			if err != nil {
				// TODO: handle error
			}
			// TODO: retry the client method; it should succeed because the credential now has the required token
		}
	}
}

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
