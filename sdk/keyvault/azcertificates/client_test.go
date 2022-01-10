//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"fmt"
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

func TestClient_BeginCreateCertificate(t *testing.T) {
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

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	cleanUp(t, client, certName)
}

func TestClient_BeginDeleteCertificate(t *testing.T) {
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

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	delResp, err := client.BeginDeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)

	delPollerResp, err := delResp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.Contains(t, *delPollerResp.ID, certName)

	_, err = client.GetCertificate(ctx, certName, nil)
	require.Error(t, err)

	deletedResp, err := client.GetDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.Contains(t, *deletedResp.ID, certName)

	_, err = client.PurgeDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)
}

func TestClient_GetCertificateOperation(t *testing.T) {
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

	_, err = resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	resp2, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)
	require.NotNil(t, resp2.ID)

	cleanUp(t, client, certName)
}

func TestClient_BackupCertificate(t *testing.T) {
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
	_, err = resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	backup, err := client.BackupCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.Greater(t, len(backup.Value), 0)

	cleanUp(t, client, certName)
}

func TestClient_ListCertificates(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	createdCount := 0
	for i := 0; i < 4; i++ {
		name, err := createRandomName(t, fmt.Sprintf("cert%d", i))
		require.NoError(t, err)
		createCert(t, client, name)
		defer cleanUp(t, client, name)
		createdCount++
	}

	pager := client.ListCertificates(nil)
	for pager.NextPage(ctx) {
		createdCount -= len(pager.PageResponse().Value)
	}

	require.Equal(t, 0, createdCount)
	require.NoError(t, pager.Err())
}

func TestClient_ListCertificateVersions(t *testing.T) {
	t.Skip()
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	name, err := createRandomName(t, "cert")
	require.NoError(t, err)
	createCert(t, client, name)
	time.Sleep(10 * delay())
	defer cleanUp(t, client, name)

	pager := client.ListCertificateVersions(name, nil)
	count := 0
	for pager.NextPage(ctx) {
		count += len(pager.PageResponse().Value)
	}

	require.Equal(t, 1, count)
	require.NoError(t, pager.Err())

	// Add a second version
	createCert(t, client, name)
	time.Sleep(10 * delay())

	pager = client.ListCertificateVersions(name, nil)
	count = 0
	for pager.NextPage(ctx) {
		count += len(pager.PageResponse().Value)
	}

	require.Equal(t, 2, count)
	require.NoError(t, pager.Err())

	// Add a third version
	createCert(t, client, name)
	time.Sleep(10 * delay())

	pager = client.ListCertificateVersions(name, nil)
	count = 0
	for pager.NextPage(ctx) {
		count += len(pager.PageResponse().Value)
	}

	require.Equal(t, 3, count)
	require.NoError(t, pager.Err())
}

func TestClient_ImportCertificate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	importedName, err := createRandomName(t, "imported")
	require.NoError(t, err)

	importResp, err := client.ImportCertificate(ctx, importedName, certContentNotPasswordEncoded, nil)
	require.NoError(t, err)
	require.Contains(t, *importResp.ID, importedName)

	cleanUp(t, client, importedName)
}

func TestClient_CreateIssuer(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	issuerName, err := createRandomName(t, "issuer")
	require.NoError(t, err)

	resp, err := client.CreateIssuer(ctx, issuerName, "Test", &CreateIssuerOptions{
		Credentials: &IssuerCredentials{
			AccountID: to.StringPtr("keyvaultuser"),
		},
		Attributes: &IssuerAttributes{
			Enabled: to.BoolPtr(true),
		},
		OrganizationDetails: &OrganizationDetails{
			AdminDetails: []*AdministratorDetails{
				{
					FirstName: to.StringPtr("John"),
					LastName: to.StringPtr("Doe"),
					EmailAddress: to.StringPtr("admin@microsoft.com"),
					Phone: to.StringPtr("4255555555"),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, *resp.Provider, "Test")
	require.Equal(t, *resp.Credentials.AccountID, "keyvaultuser")
	require.Contains(t, *resp.ID, fmt.Sprintf("/certificates/issuers/%s", issuerName))
}
