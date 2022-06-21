//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
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

	certName, err := createRandomName(t, "beginCreate")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	defer cleanUp(t, client, certName)

	// want to interface with x509 std library

	cert, err := x509.ParseCertificate(pollerResp.CER)
	require.NoError(t, err)
	require.NotNil(t, cert)
}

func TestClient_BeginCreateCertificateRehydrated(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "beginCreateRehydrate")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)

	rt, err := resp.ResumeToken()
	require.NoError(t, err)

	newPoller, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), &BeginCreateCertificateOptions{ResumeToken: rt})
	require.NoError(t, err)

	pollerResp, err := newPoller.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	defer cleanUp(t, client, certName)

	// want to interface with x509 std library

	cert, err := x509.ParseCertificate(pollerResp.CER)
	require.NoError(t, err)
	require.NotNil(t, cert)
}

func TestClient_BeginDeleteCertificate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "createCert")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
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

func TestClient_BeginDeleteCertificateRehydrated(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "createCertRehydrate")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	delResp, err := client.BeginDeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)

	rt, err := delResp.ResumeToken()
	require.NoError(t, err)

	poller, err := client.BeginDeleteCertificate(ctx, certName, &BeginDeleteCertificateOptions{ResumeToken: rt})
	require.NoError(t, err)

	delPollerResp, err := poller.PollUntilDone(ctx, delay())
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

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	resp2, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)
	require.NotNil(t, resp2.ID)

	cleanUp(t, client, certName)
}
func TestClient_CancelCertificateOperation(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "cert")
	require.NoError(t, err)

	_, err = client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)

	cancelResp, err := client.CancelCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)
	require.Contains(t, *cancelResp.ID, certName)

	getResp, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)
	require.Equal(t, true, *getResp.CancellationRequested)

	_, err = client.DeleteCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)

	// Get should fail now
	_, err = client.GetCertificateOperation(ctx, certName, nil)
	require.Error(t, err)
}

func TestClient_BackupCertificate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "cert")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
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

	created := [4]string{}
	createdCount := 0
	for i := 0; i < 4; i++ {
		name, err := createRandomName(t, fmt.Sprintf("listcerts%d", i))
		created[i] = name
		require.NoError(t, err)
		createCert(t, client, name)
		defer cleanUp(t, client, name)
		createdCount++
	}

	pager := client.NewListPropertiesOfCertificatesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, cert := range page.Certificates {
			// the list of certs will contain results from other tests
			// so we just count the ones with the prefix from this test
			if strings.HasPrefix(*cert.Properties.Name, "listcerts") {
				createdCount--
			}
		}
	}

	require.Equal(t, 0, createdCount)
}

func TestClient_ListCertificateVersions(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	name, err := createRandomName(t, "cert1")
	require.NoError(t, err)
	createCert(t, client, name)
	defer cleanUp(t, client, name)

	pager := client.NewListPropertiesOfCertificateVersionsPager(name, nil)
	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		count += len(resp.Certificates)
	}

	require.Equal(t, 1, count)

	// Add a second version
	createCert(t, client, name)

	pager = client.NewListPropertiesOfCertificateVersionsPager(name, nil)
	count = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		count += len(resp.Certificates)
	}

	require.Equal(t, 2, count)

	// Add a third version
	createCert(t, client, name)

	pager = client.NewListPropertiesOfCertificateVersionsPager(name, nil)
	count = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		count += len(resp.Certificates)
	}

	require.Equal(t, 3, count)
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
	require.NotNil(t, importResp.Policy)

	cleanUp(t, client, importedName)
}

func TestClient_IssuerCRUD(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	issuerName, err := createRandomName(t, "issuer")
	require.NoError(t, err)

	resp, err := client.CreateIssuer(ctx, issuerName, "Test", &CreateIssuerOptions{
		Credentials: &IssuerCredentials{
			AccountID: to.Ptr("keyvaultuser"),
		},
		Enabled: to.Ptr(true),
		AdministratorContacts: []*AdministratorContact{
			{
				FirstName: to.Ptr("John"),
				LastName:  to.Ptr("Doe"),
				Email:     to.Ptr("admin@microsoft.com"),
				Phone:     to.Ptr("4255555555"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, *resp.Issuer.Provider, "Test")
	require.Equal(t, *resp.Issuer.Credentials.AccountID, "keyvaultuser")
	require.Contains(t, *resp.Issuer.ID, fmt.Sprintf("/certificates/issuers/%s", issuerName))

	getResp, err := client.GetIssuer(ctx, issuerName, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Issuer.Provider, "Test")
	require.Equal(t, *getResp.Issuer.Credentials.AccountID, "keyvaultuser")
	require.Contains(t, *getResp.Issuer.ID, fmt.Sprintf("/certificates/issuers/%s", issuerName))

	issuerName2, err := createRandomName(t, "issuer2")
	require.NoError(t, err)

	createResp, err := client.CreateIssuer(ctx, issuerName2, "Test", &CreateIssuerOptions{
		Credentials: &IssuerCredentials{
			AccountID: to.Ptr("keyvaultuser2"),
		},
		Enabled: to.Ptr(true),
		AdministratorContacts: []*AdministratorContact{
			{
				FirstName: to.Ptr("John"),
				LastName:  to.Ptr("Doe"),
				Email:     to.Ptr("admin@microsoft.com"),
				Phone:     to.Ptr("4255555555"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, issuerName2, *createResp.Issuer.Name)

	// List operation
	pager := client.NewListPropertiesOfIssuersPager(nil)
	count := 0
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, issuer := range page.Issuers {
			require.Equal(t, "Test", *issuer.Provider)
			count += 1
		}
	}
	require.GreaterOrEqual(t, count, 2)

	createResp.Issuer.AdministratorContacts = []*AdministratorContact{
		{
			FirstName: to.Ptr("Jane"),
			LastName:  to.Ptr("Doey"),
			Email:     to.Ptr("admin2@microsoft.com"),
			Phone:     to.Ptr("4266666666"),
		},
	}
	// Update the certificate issuer
	updateResp, err := client.UpdateIssuer(ctx, createResp.Issuer, nil)
	require.NoError(t, err)
	require.Equal(t, 1, len(updateResp.Issuer.AdministratorContacts))
	require.Equal(t, "Jane", *updateResp.Issuer.AdministratorContacts[0].FirstName)
	require.Equal(t, "Doey", *updateResp.Issuer.AdministratorContacts[0].LastName)
	require.Equal(t, "admin2@microsoft.com", *updateResp.Issuer.AdministratorContacts[0].Email)
	require.Equal(t, "4266666666", *updateResp.Issuer.AdministratorContacts[0].Phone)

	// Delete the first issuer
	_, err = client.DeleteIssuer(ctx, issuerName, nil)
	require.NoError(t, err)

	// Get on the first issuer fails
	_, err = client.GetIssuer(ctx, issuerName, nil)
	require.Error(t, err)
}

func TestClient_ContactsCRUD(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	contacts := Contacts{ContactList: []*Contact{
		{Email: to.Ptr("admin@microsoft.com"), Name: to.Ptr("John Doe"), Phone: to.Ptr("1111111111")},
		{Email: to.Ptr("admin@contoso.com"), Name: to.Ptr("Jane Doey"), Phone: to.Ptr("2222222222")},
	}}

	resp, err := client.SetContacts(ctx, contacts.ContactList, nil)
	require.NoError(t, err)
	require.Equal(t, 2, len(resp.ContactList))

	getResp, err := client.GetContacts(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 2, len(getResp.ContactList))
	require.Equal(t, "admin@microsoft.com", *getResp.ContactList[0].Email)
	require.Equal(t, "admin@contoso.com", *getResp.ContactList[1].Email)
	require.Equal(t, "John Doe", *getResp.ContactList[0].Name)
	require.Equal(t, "Jane Doey", *getResp.ContactList[1].Name)
	require.Equal(t, "1111111111", *getResp.ContactList[0].Phone)
	require.Equal(t, "2222222222", *getResp.ContactList[1].Phone)

	deleteResp, err := client.DeleteContacts(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 2, len(deleteResp.ContactList))

	// Get should fail
	_, err = client.GetContacts(ctx, nil)
	require.Error(t, err)
}

func TestPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "policyCertificate")
	require.NoError(t, err)

	policy := Policy{
		IssuerParameters: &IssuerParameters{
			CertificateTransparency: to.Ptr(false),
			IssuerName:              to.Ptr("Self"),
		},
		Exportable: to.Ptr(true),
		KeySize:    to.Ptr(int32(2048)),
		ReuseKey:   to.Ptr(true),
		KeyType:    to.Ptr(KeyTypeRSA),
		LifetimeActions: []*LifetimeAction{
			{Action: to.Ptr(PolicyActionEmailContacts), LifetimePercentage: to.Ptr(int32(98))},
		},
		ContentType: to.Ptr(CertificateContentTypePKCS12),
		X509Properties: &X509CertificateProperties{
			EnhancedKeyUsages: []*string{to.Ptr("1.3.6.1.5.5.7.3.1"), to.Ptr("1.3.6.1.5.5.7.3.2")},
			KeyUsages:         []*KeyUsage{to.Ptr(KeyUsageDecipherOnly)},
			Subject:           to.Ptr("CN=DefaultPolicy"),
			ValidityInMonths:  to.Ptr(int32(12)),
			SubjectAlternativeNames: &SubjectAlternativeNames{
				DNSNames: []*string{to.Ptr("sdk.azure-int.net")},
			},
		},
	}

	_, err = client.BeginCreateCertificate(ctx, certName, policy, nil)
	require.NoError(t, err)

	receivedPolicy, err := client.GetCertificatePolicy(ctx, certName, nil)
	require.NoError(t, err)

	// Make sure policies are equal
	require.Equal(t, *policy.IssuerParameters.IssuerName, *receivedPolicy.Policy.IssuerParameters.IssuerName)
	require.Equal(t, *policy.Exportable, *receivedPolicy.Exportable)
	require.Equal(t, *policy.ContentType, *receivedPolicy.ContentType)

	// Update the policy
	policy.KeyType = to.Ptr(KeyTypeEC)
	policy.KeySize = to.Ptr(int32(256))
	policy.KeyCurveName = to.Ptr(KeyCurveNameP256)

	updateResp, err := client.UpdateCertificatePolicy(ctx, certName, policy, nil)
	require.NoError(t, err)

	require.Equal(t, *policy.IssuerParameters.IssuerName, *updateResp.Policy.IssuerParameters.IssuerName)
	require.Equal(t, *policy.Exportable, *updateResp.Exportable)
	require.Equal(t, *policy.ContentType, *updateResp.ContentType)
	require.Equal(t, *policy.KeyType, *updateResp.KeyType)
	require.Equal(t, *policy.KeySize, *updateResp.KeySize)
	require.Equal(t, *policy.KeyCurveName, *updateResp.KeyCurveName)

}

func TestCRUDOperations(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "cert")
	require.NoError(t, err)

	policy := Policy{
		IssuerParameters: &IssuerParameters{
			CertificateTransparency: to.Ptr(false),
			IssuerName:              to.Ptr("Self"),
		},
		Exportable: to.Ptr(true),
		KeySize:    to.Ptr(int32(2048)),
		ReuseKey:   to.Ptr(true),
		KeyType:    to.Ptr(KeyTypeRSA),
		LifetimeActions: []*LifetimeAction{
			{Action: to.Ptr(PolicyActionEmailContacts), LifetimePercentage: to.Ptr(int32(98))},
		},
		ContentType: to.Ptr(CertificateContentTypePKCS12),
		X509Properties: &X509CertificateProperties{
			EnhancedKeyUsages: []*string{to.Ptr("1.3.6.1.5.5.7.3.1"), to.Ptr("1.3.6.1.5.5.7.3.2")},
			KeyUsages:         []*KeyUsage{to.Ptr(KeyUsageDecipherOnly)},
			Subject:           to.Ptr("CN=DefaultPolicy"),
			ValidityInMonths:  to.Ptr(int32(12)),
			SubjectAlternativeNames: &SubjectAlternativeNames{
				DNSNames: []*string{to.Ptr("sdk.azure-int.net")},
			},
		},
	}

	pollerResp, err := client.BeginCreateCertificate(ctx, certName, policy, nil)
	require.NoError(t, err)
	finalResp, err := pollerResp.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)

	received, err := client.GetCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.NotNil(t, received.Policy)

	// Make sure certificates are the same
	require.Equal(t, *finalResp.ID, *received.ID)

	// Update the policy
	policy.KeyType = to.Ptr(KeyTypeEC)
	policy.KeySize = to.Ptr(int32(256))
	policy.KeyCurveName = to.Ptr(KeyCurveNameP256)

	updateResp, err := client.UpdateCertificatePolicy(ctx, certName, policy, nil)
	require.NoError(t, err)

	require.Equal(t, *policy.IssuerParameters.IssuerName, *updateResp.Policy.IssuerParameters.IssuerName)
	require.Equal(t, *policy.Exportable, *updateResp.Exportable)
	require.Equal(t, *policy.ContentType, *updateResp.ContentType)
	require.Equal(t, *policy.KeyType, *updateResp.KeyType)
	require.Equal(t, *policy.KeySize, *updateResp.KeySize)
	require.Equal(t, *policy.KeyCurveName, *updateResp.KeyCurveName)

	if received.Properties.Tags == nil {
		received.Properties.Tags = map[string]*string{}
	}
	received.Properties.Tags["tag1"] = to.Ptr("updated_values1")
	updatePropsResp, err := client.UpdateCertificateProperties(ctx, *received.Properties, nil)
	require.NoError(t, err)
	require.Equal(t, "updated_values1", *updatePropsResp.Properties.Tags["tag1"])
	require.Equal(t, *received.ID, *updatePropsResp.ID)
	require.True(t, *updatePropsResp.Properties.Enabled)

	received.Properties.Enabled = to.Ptr(false)
	resp, err := client.UpdateCertificateProperties(ctx, *received.Properties, nil)
	require.NoError(t, err)
	require.False(t, *resp.Properties.Enabled)
	require.Equal(t, "updated_values1", *resp.Properties.Tags["tag1"])
}

// https://stackoverflow.com/questions/42643048/signing-certificate-request-with-certificate-authority
// Much of this is thanks to this response, thanks @krostar
func TestMergeCertificate(t *testing.T) {
	recording.LiveOnly(t)
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "mergeCertificate")
	require.NoError(t, err)

	certPolicy := Policy{
		IssuerParameters: &IssuerParameters{
			IssuerName:              to.Ptr("Unknown"),
			CertificateTransparency: to.Ptr(false),
		},
		X509Properties: &X509CertificateProperties{
			Subject: to.Ptr("CN=MyCert"),
		},
	}

	poller, err := client.BeginCreateCertificate(ctx, certName, certPolicy, nil)
	require.NoError(t, err)
	// can't PollUntilDone for this scenario
	resp, err := poller.Poll(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" && recording.GetRecordMode() != recording.PlaybackMode {
		// sleep before moving on
		asInt, err := strconv.Atoi(retryAfter)
		require.NoError(t, err)
		time.Sleep(time.Duration(asInt) * time.Second)
	}
	defer cleanUp(t, client, certName)

	// Load public key
	data, err := ioutil.ReadFile("testdata/ca.crt")
	require.NoError(t, err)
	block, _ := pem.Decode(data)
	require.NotNil(t, block)
	caCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	data, err = ioutil.ReadFile("testdata/ca.key")
	require.NoError(t, err)
	pkeyBlock, _ := pem.Decode(data)
	require.NotNil(t, pkeyBlock)
	require.Equal(t, pkeyBlock.Type, "RSA PRIVATE KEY")
	pkey, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes)
	require.NoError(t, err)

	certOpResp, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)

	mid := base64.StdEncoding.EncodeToString(certOpResp.CSR)
	csr := fmt.Sprintf("-----BEGIN CERTIFICATE REQUEST-----\n%s\n-----END CERTIFICATE REQUEST-----", mid)

	// load certificate request
	csrblock, _ := pem.Decode([]byte(csr))
	require.NotNil(t, csrblock)
	req, err := x509.ParseCertificateRequest(csrblock.Bytes)
	require.NoError(t, err)
	require.NoError(t, req.CheckSignature())

	cert := x509.Certificate{
		SerialNumber:       big.NewInt(1),
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(24 * time.Hour),
		Issuer:             caCert.Issuer,
		Subject:            req.Subject,
		PublicKey:          req.PublicKey,
		PublicKeyAlgorithm: req.PublicKeyAlgorithm,
		SignatureAlgorithm: req.SignatureAlgorithm,
		Signature:          req.Signature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, caCert, req.PublicKey, pkey)
	require.NoError(t, err)

	// Need to strip the BEGIN/END from the certificate
	certificateString := string(certBytes)
	certificateString = strings.Replace(certificateString, "-----Begin Certificate-----", "", 1)
	certificateString = strings.Replace(certificateString, "-----End Certificate-----", "", 1)

	mergeResp, err := client.MergeCertificate(ctx, certName, [][]byte{[]byte(certificateString)}, nil)
	require.NoError(t, err)
	require.NotNil(t, mergeResp.Policy)
}

func TestClient_BeginRecoverDeletedCertificate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "certRecover")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
	require.NoError(t, err)
	defer cleanUp(t, client, certName)

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

	recover, err := client.BeginRecoverDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)

	recoveredResp, err := recover.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	require.Contains(t, *recoveredResp.ID, certName)
}

func TestClient_BeginRecoverDeletedCertificateRehydrated(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "certBeginRecoverRehydrated")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, NewDefaultCertificatePolicy(), nil)
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

	recover, err := client.BeginRecoverDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)

	rt, err := recover.ResumeToken()
	require.NoError(t, err)

	poller, err := client.BeginRecoverDeletedCertificate(ctx, certName, &BeginRecoverDeletedCertificateOptions{ResumeToken: rt})
	require.NoError(t, err)

	recoveredResp, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	require.NoError(t, err)
	require.Contains(t, *recoveredResp.ID, certName)

	_, err = client.GetCertificate(ctx, certName, nil)
	require.NoError(t, err)
}

func TestClient_RestoreCertificateBackup(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	certName, err := createRandomName(t, "certRestore")
	require.NoError(t, err)

	resp, err := client.BeginCreateCertificate(ctx, certName, Policy{
		IssuerParameters: &IssuerParameters{
			IssuerName: to.Ptr("Self"),
		},
		X509Properties: &X509CertificateProperties{
			Subject: to.Ptr("CN=DefaultPolicy"),
			SubjectAlternativeNames: &SubjectAlternativeNames{
				UserPrincipalNames: []*string{to.Ptr("john.doe@domain.com")},
			},
		},
	}, nil)
	require.NoError(t, err)
	defer cleanUp(t, client, certName)

	pollerResp, err := resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
	require.NotNil(t, pollerResp.ID)

	// Create a backup
	certificateBackupResp, err := client.BackupCertificate(ctx, certName, nil)
	require.NoError(t, err)

	// Delete the certificate
	deletePoller, err := client.BeginDeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)

	_, err = deletePoller.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	// Purge the cert
	_, err = client.PurgeDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)

	// Restore the cert
	// Poll until no exception
	count := 0
	for {
		resp, err := client.RestoreCertificateBackup(ctx, certificateBackupResp.Value, nil)
		if err == nil {
			require.NotNil(t, resp.Policy)
			break
		}
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			if respErr.RawResponse.StatusCode != 409 {
				require.NoError(t, err)
			}
		} else {
			require.NoError(t, err)
		}
		count += 1
		if count > 25 {
			require.NoError(t, err)
		}
		recording.Sleep(5 * time.Second)
	}
}

func TestClient_ListDeletedCertificates(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	var names []string
	createdCount := 0
	for i := 0; i < 4; i++ {
		name, err := createRandomName(t, fmt.Sprintf("delCert%d", i))
		require.NoError(t, err)
		names = append(names, name)
		createCert(t, client, name)
		createdCount++
	}
	require.Equal(t, 4, createdCount)

	for _, name := range names {
		poller, err := client.BeginDeleteCertificate(ctx, name, nil)
		require.NoError(t, err)
		_, err = poller.PollUntilDone(ctx, delay())
		require.NoError(t, err)
	}

	pager := client.NewListDeletedCertificatesPager(nil)
	deletedCount := 0
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		for _, cert := range page.DeletedCertificates {
			purgeCert(t, client, *cert.Name)
			deletedCount += 1
		}
	}
	require.Equal(t, createdCount, deletedCount)
}
