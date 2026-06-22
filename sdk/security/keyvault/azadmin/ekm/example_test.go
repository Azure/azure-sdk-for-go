// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
)

var ekmClient *ekm.KeyVaultClient

func ExampleNewClient() {
	vaultURL := "https://<TODO: your managed HSM name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := ekm.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client
}

func ExampleKeyVaultClient_CreateEkmConnection() {
	// Read the DER-encoded CA certificate chain that issued the EKM proxy's
	// server certificate. Each certificate in the chain is a separate []byte.
	caCert, err := os.ReadFile("ekm-proxy-ca.der")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	conn := ekm.Connection{
		Host:                 to.Ptr("ekm-proxy.example.com:443"),
		ServerCaCertificates: [][]byte{caCert},
		// ServerSubjectCommonName must equal the host portion of Host (without
		// the port) so the Managed HSM can validate the EKM proxy's TLS cert.
		ServerSubjectCommonName: to.Ptr("ekm-proxy.example.com"),
		PathPrefix:              to.Ptr("/v1"),
	}

	created, err := ekmClient.CreateEkmConnection(context.TODO(), conn, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("Created EKM connection to host: %s", *created.Host)
}

func ExampleKeyVaultClient_GetEkmConnection() {
	got, err := ekmClient.GetEkmConnection(context.TODO(), nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("Configured EKM host: %s, path prefix: %s",
		*got.Host, *got.PathPrefix)
}

func ExampleKeyVaultClient_UpdateEkmConnection() {
	// Fetch the current configuration so we only change the fields we intend to.
	current, err := ekmClient.GetEkmConnection(context.TODO(), nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	current.PathPrefix = to.Ptr("/v2")
	updated, err := ekmClient.UpdateEkmConnection(context.TODO(), current.Connection, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("Updated EKM path prefix to: %s", *updated.PathPrefix)
}

func ExampleKeyVaultClient_DeleteEkmConnection() {
	deleted, err := ekmClient.DeleteEkmConnection(context.TODO(), nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("Deleted EKM connection that pointed at: %s", *deleted.Host)
}

func ExampleKeyVaultClient_GetEkmCertificate() {
	cert, err := ekmClient.GetEkmCertificate(context.TODO(), nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// The Managed HSM presents this client certificate to the EKM proxy.
	// Add it to the EKM proxy's trust store so the proxy accepts requests from
	// this HSM instance.
	fmt.Printf("HSM client certificate subject: %s, CA chain length: %d",
		*cert.SubjectCommonName, len(cert.CaCertificates))
}

func ExampleKeyVaultClient_CheckEkmConnection() {
	// Verifies that the Managed HSM can reach the configured EKM proxy and
	// authenticate to it. Returns information about the EKM product behind
	// the proxy on success.
	info, err := ekmClient.CheckEkmConnection(context.TODO(), nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("EKM proxy: %s %s (API %s), backed by %s %s",
		*info.ProxyVendor, *info.ProxyName, *info.APIVersion,
		*info.EkmVendor, *info.EkmProduct)
}
