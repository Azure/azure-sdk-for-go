// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	certPath = "testdata/certificate.pem"
	clientID = "fake-client-id"
	tenantID = "fake-tenant"
)

var cred *azidentity.ClientCertificateCredential

func handleError(err error) {
	if err != nil {
		log.Panicf("example failed: %v", err)
	}
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
