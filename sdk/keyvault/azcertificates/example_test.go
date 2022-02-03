//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

func ExampleClient_BeginCreateCertificate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginCreateCertificate(context.TODO(), "certificateName", azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{
			Name: to.StringPtr("Self"),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.StringPtr("CN=DefaultPolicy"),
		},
	}, nil)
	if err != nil {
		panic(err)
	}

	finalResponse, err := resp.PollUntilDone(context.TODO(), time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created a certificate with ID: ", *finalResponse.ID)
}

func ExampleClient_GetCertificate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetCertificate(context.TODO(), "myCertName", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.ID)

	// optionally you can get a specific version
	resp, err = client.GetCertificate(context.TODO(), "myCertName", &azcertificates.GetCertificateOptions{Version: "myCertVersion"})
	if err != nil {
		panic(err)
	}
}
func ExampleClient_UpdateCertificateProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.UpdateCertificateProperties(context.TODO(), "myCertName", &azcertificates.UpdateCertificatePropertiesOptions{
		Version: "myNewVersion",
		CertificateAttributes: &azcertificates.CertificateProperties{
			Enabled: to.BoolPtr(false),
			Expires: to.TimePtr(time.Now().Add(72 * time.Hour)),
		},
		CertificatePolicy: &azcertificates.CertificatePolicy{
			IssuerParameters: &azcertificates.IssuerParameters{
				Name: to.StringPtr("Self"),
			},
		},
		Tags: map[string]string{"Tag1": "Val1"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.ID)
	fmt.Println(*resp.KeyVaultCertificate.Properties.Enabled)
	fmt.Println(resp.Tags)
}

func ExampleClient_ListCertificates() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	poller := client.ListCertificates(nil)
	for poller.NextPage(context.TODO()) {
		for _, cert := range poller.PageResponse().Certificates {
			fmt.Println(*cert.ID)
		}
	}
	if poller.Err() != nil {
		panic(err)
	}
}

func ExampleClient_BeginDeleteCertificate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		log.Fatalf("Could not find 'AZURE_KEYVAULT_URL' in environment variables")
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	pollerResp, err := client.BeginDeleteCertificate(context.TODO(), "certToDelete", nil)
	if err != nil {
		panic(err)
	}
	finalResp, err := pollerResp.PollUntilDone(context.TODO(), time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted certificate with ID: ", *finalResp.ID)
}
