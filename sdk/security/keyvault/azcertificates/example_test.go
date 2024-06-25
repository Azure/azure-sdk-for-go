//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

var client *azcertificates.Client

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	vaultURL := "https://<TODO: your vault name>.vault.azure.net"
	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_CreateCertificate() {
	createParams := azcertificates.CreateCertificateParameters{
		// this policy is suitable for a self-signed certificate
		CertificatePolicy: &azcertificates.CertificatePolicy{
			IssuerParameters:          &azcertificates.IssuerParameters{Name: to.Ptr("self")},
			X509CertificateProperties: &azcertificates.X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
		},
	}
	// if a certificate with the same name already exists, a new version of the certificate is created
	resp, err := client.CreateCertificate(context.TODO(), "certificateName", createParams, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Println("Created a certificate with ID:", *resp.ID)
}

func ExampleClient_GetCertificate() {
	// passing an empty string for the version gets the latest version of the certificate
	resp, err := client.GetCertificate(context.TODO(), "certName", "", nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.ID)
}

func ExampleClient_UpdateCertificate() {
	updateParams := azcertificates.UpdateCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{Enabled: to.Ptr(false)},
	}
	// passing an empty string for the version updates the latest version of the certificate
	resp, err := client.UpdateCertificate(context.TODO(), "certName", "", updateParams, nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(*resp.ID)
}

func ExampleClient_NewListCertificatePropertiesPager() {
	pager := client.NewListCertificatePropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: handle error
		}
		for _, cert := range page.Value {
			fmt.Println(*cert.ID)
		}
	}
}

func ExampleClient_DeleteCertificate() {
	// DeleteCertificate returns when Key Vault has begun deleting the certificate. That can take several
	// seconds to complete, so it may be necessary to wait before performing other operations on the
	// deleted certificate.
	resp, err := client.DeleteCertificate(context.TODO(), "certName", nil)
	if err != nil {
		// TODO: handle error
	}

	// In a soft-delete enabled vault, deleted resources can be recovered until they're purged (permanently deleted).
	fmt.Printf("Certificate will be purged at %v", *resp.ScheduledPurgeDate)
}

func ExampleClient_ImportCertificate_pfx() {
	// This example uses `ImportCertificate` to import a PFX certificate.

	// Assuming you already have a PFX containing your key pair, you can import it into Key Vault.
	// You can do this without setting a policy, but the policy is needed if you want the private key to be exportable
	// or to configure actions when a certificate is close to expiration.
	parameters := azcertificates.ImportCertificateParameters{
		Base64EncodedCertificate: to.Ptr("<TODO: pfx cert value>"),
	}

	resp, err := client.ImportCertificate(context.TODO(), "pfxCertName", parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("PFX certificate %s imported successfully.", resp.ID.Name())
}

func ExampleClient_ImportCertificate_pem() {
	// This example uses `ImportCertificate` to import a PEM certificate.

	// To import a PEM-formatted certificate, you must provide a CertificatePolicy that sets the ContentType to
	// CertificateContentType.pem or the certificate will fail to import (the default content type is PFX).
	parameters := azcertificates.ImportCertificateParameters{
		Base64EncodedCertificate: to.Ptr("<TODO: pem cert value>"),
		CertificatePolicy: &azcertificates.CertificatePolicy{

			SecretProperties: &azcertificates.SecretProperties{

				ContentType: to.Ptr("application/x-pem-file"),
			},
		},
	}

	resp, err := client.ImportCertificate(context.TODO(), "pemCertName", parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("PFX certificate %s imported successfully.", resp.ID.Name())

}
