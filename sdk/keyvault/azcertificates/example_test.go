//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates_test

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

	resp, err := client.BeginCreateCertificate(context.TODO(), "certificateName", azcertificates.NewDefaultCertificatePolicy(), nil)
	if err != nil {
		panic(err)
	}

	finalResponse, err := resp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
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

	getResp, err := client.GetCertificate(context.TODO(), "myCertName", nil)
	if err != nil {
		panic(err)
	}
	getResp.Properties.Enabled = to.Ptr(false)
	getResp.Properties.Tags["Tag1"] = to.Ptr("Val1")

	resp, err := client.UpdateCertificateProperties(context.TODO(), "myCertName", *getResp.Properties, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp.ID)
	fmt.Println(*resp.Properties.Enabled)
	fmt.Println(resp.Properties.Tags)
}

func ExampleClient_NewListPropertiesOfCertificatesPager() {
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

	pager := client.NewListPropertiesOfCertificatesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, cert := range page.Certificates {
			fmt.Println(*cert.ID)
		}
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
	finalResp, err := pollerResp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted certificate with ID: ", *finalResp.ID)
}

func ExampleClient_MergeCertificate() {
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

	certPolicy := azcertificates.NewDefaultCertificatePolicy()
	certPolicy.IssuerParameters.CertificateTransparency = to.Ptr(false)

	// Load public key
	data, err := ioutil.ReadFile("path/to/cert/ca.crt")
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(data)
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	data, err = ioutil.ReadFile("path/to/certkey/ca.key")
	if err != nil {
		panic(err)
	}
	pkeyBlock, _ := pem.Decode(data)
	pkey, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginCreateCertificate(context.TODO(), "myCertName", certPolicy, nil)
	if err != nil {
		panic(err)
	}
	_, err = resp.PollUntilDone(context.TODO(), &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		panic(err)
	}

	certOpResp, err := client.GetCertificateOperation(context.TODO(), "myCertName", nil)
	if err != nil {
		panic(err)
	}

	mid := base64.StdEncoding.EncodeToString(certOpResp.CSR)
	csr := fmt.Sprintf("-----BEGIN CERTIFICATE REQUEST-----\n%s\n-----END CERTIFICATE REQUEST-----", mid)

	// load certificate request
	csrblock, _ := pem.Decode([]byte(csr))
	req, err := x509.ParseCertificateRequest(csrblock.Bytes)
	if err != nil {
		panic(err)
	}

	cert := x509.Certificate{
		SerialNumber:       big.NewInt(1),
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(24 * time.Hour),
		Issuer:             caCert.Issuer,
		Subject:            req.Subject,
		PublicKey:          req.PublicKey,
		PublicKeyAlgorithm: req.PublicKeyAlgorithm,
		Signature:          req.Signature,
		SignatureAlgorithm: req.SignatureAlgorithm,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, caCert, req.PublicKey, pkey)
	if err != nil {
		panic(err)
	}

	// Need to strip the BEGIN/END from the certificate
	certificateString := string(certBytes)
	certificateString = strings.Replace(certificateString, "-----Begin Certificate-----", "", 1)
	certificateString = strings.Replace(certificateString, "-----End Certificate-----", "", 1)

	mergeResp, err := client.MergeCertificate(context.TODO(), "myCertName", [][]byte{[]byte(certificateString)}, nil)
	if err != nil {
		panic(err)
	}
	_ = mergeResp
}
