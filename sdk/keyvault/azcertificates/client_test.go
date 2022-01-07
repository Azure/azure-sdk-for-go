//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

var ctx = context.TODO()

func TestNewClient(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	client, err := NewClient("https://certvault.vault.azure.net", cred, nil)
	require.NoError(t, err)
	require.NotNil(t, client.genClient)
	require.Equal(t, "https://certvault.vault.azure.net", client.vaultURL)
}

func TestBeginCreateNewCertificate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "cert")
	require.NoError(t, err)
	resp, err := client.BeginCreateCertificate(ctx, certName, CertificatePolicy{
		IssuerParameters: &IssuerParameters{
			Name: to.StringPtr("Self"),
		},
		X509CertificateProperties: &X509CertificateProperties{
			Subject: to.StringPtr("CN=DefaultPolicy"),
		},
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.ID)

	time.Sleep(time.Second * 30)

	getResp, err := client.GetCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.ID)
}
