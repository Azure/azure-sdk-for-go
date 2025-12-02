// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates_test

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/stretchr/testify/require"
)

var (
	ctx              = context.Background()
	selfSignedPolicy = azcertificates.CertificatePolicy{
		IssuerParameters:          &azcertificates.IssuerParameters{Name: to.Ptr("self")},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{Subject: to.Ptr("CN=DefaultPolicy")},
	}
)

// pollStatus calls a function until it stops returning a response error with the given status code.
// If this takes more than 2 minutes, it fails the test.
func pollStatus(t *testing.T, expectedStatus int, fn func() error) {
	var err error
	for i := 0; i < 12; i++ {
		err = fn()
		var respErr *azcore.ResponseError
		if !errors.As(err, &respErr) || respErr.StatusCode != expectedStatus {
			break
		}
		if i < 11 {
			recording.Sleep(10 * time.Second)
		}
	}
	require.NoError(t, err)
}

// pollCertOperation polls a certificate operation for up to 2 minutes, stopping when it completes.
// It fails the test if a poll fails or the operation doesn't complete successfully in the allotted time.
func pollCertOperation(t *testing.T, client *azcertificates.Client, name string) {
	var err error
	var op azcertificates.GetCertificateOperationResponse
	for i := 0; i < 24; i++ {
		op, err = client.GetCertificateOperation(ctx, name, nil)
		require.NoError(t, err)
		require.NotNil(t, op.Status)
		switch s := *op.Status; s {
		case "completed":
			return
		case "cancelled":
			t.Fatal("cert creation cancelled")
		case "inProgress":
			// sleep and continue
		default:
			t.Fatalf(`unexpected status "%s"`, s)
		}
		if i < 23 {
			recording.Sleep(5 * time.Second)
		} else {
			t.Fatal("cert creation didn't complete in time")
		}
	}
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)
}

func TestBackupRestore(t *testing.T) {
	client := startTest(t)

	certName := getName(t, "cert")
	createParams := azcertificates.CreateCertificateParameters{CertificatePolicy: &selfSignedPolicy}
	_, err := client.CreateCertificate(ctx, certName, createParams, nil)
	require.NoError(t, err)
	pollCertOperation(t, client, certName)

	backup, err := client.BackupCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.NotEmpty(t, backup.Value)
	testSerde(t, &backup.BackupCertificateResult)

	deleteResp, err := client.DeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)
	pollStatus(t, http.StatusNotFound, func() error {
		_, err = client.GetDeletedCertificate(ctx, certName, nil)
		return err
	})

	_, err = client.PurgeDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)

	var restoreResp azcertificates.RestoreCertificateResponse
	restoreParams := azcertificates.RestoreCertificateParameters{CertificateBackup: backup.Value}
	pollStatus(t, http.StatusConflict, func() error {
		restoreResp, err = client.RestoreCertificate(ctx, restoreParams, nil)
		return err
	})
	require.Equal(t, deleteResp.ID, restoreResp.ID)
	require.NotNil(t, restoreResp.Attributes)
	cleanUpCert(t, client, certName)

	// exercise otherwise unused mashalling code
	rp := azcertificates.RestoreCertificateParameters{}
	data, err := restoreParams.MarshalJSON()
	require.NoError(t, err)
	err = rp.UnmarshalJSON(data)
	require.NoError(t, err)
}

func TestContactsCRUD(t *testing.T) {
	client := startTest(t)

	contacts := azcertificates.Contacts{ContactList: []*azcertificates.Contact{
		{Email: to.Ptr("one@test.com"), Name: to.Ptr("One"), Phone: to.Ptr("1111111111")},
		{Email: to.Ptr("two@test.com"), Name: to.Ptr("Two"), Phone: to.Ptr("2222222222")},
	}}
	setResp, err := client.SetContacts(ctx, contacts, nil)
	require.NoError(t, err)
	require.Equal(t, contacts.ContactList, setResp.ContactList)

	getResp, err := client.GetContacts(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, contacts.ContactList, getResp.ContactList)

	_, err = client.DeleteContacts(ctx, nil)
	require.NoError(t, err)
}

func TestCRUD(t *testing.T) {
	client := startTest(t)

	certName := getName(t, "")
	createParams := azcertificates.CreateCertificateParameters{
		CertificatePolicy: &selfSignedPolicy,
		PreserveCertOrder: to.Ptr(true),
	}
	testSerde(t, &createParams)
	createResp, err := client.CreateCertificate(ctx, certName, createParams, nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, createResp.CSR)
	require.NotNil(t, createResp.CancellationRequested)
	require.False(t, *createResp.CancellationRequested)
	require.Nil(t, createResp.Error)
	require.NotEmpty(t, createResp.RequestID)
	require.NotEmpty(t, createResp.Status)
	require.NotEmpty(t, createResp.StatusDetails)
	require.NotEmpty(t, createResp.ID)
	require.True(t, *createResp.PreserveCertOrder)
	pollCertOperation(t, client, certName)

	getResp, err := client.GetCertificate(ctx, certName, "", nil)
	require.NoError(t, err)
	require.NotEmpty(t, getResp.ID)
	require.NotEmpty(t, getResp.KID)
	require.NotEmpty(t, getResp.SID)
	testSerde(t, &getResp.Certificate)

	updateParams := azcertificates.UpdateCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Expires: to.Ptr(time.Date(2030, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
	}
	testSerde(t, &updateParams)
	_, err = client.UpdateCertificate(ctx, certName, "", updateParams, nil)
	require.NoError(t, err)

	deleteResp, err := client.DeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)
	require.NotNil(t, deleteResp.Attributes)
	require.Equal(t, getResp.CER, deleteResp.CER)
	require.Equal(t, getResp.ContentType, deleteResp.ContentType)
	require.NotEmpty(t, deleteResp.ID)
	require.Equal(t, certName, deleteResp.ID.Name())
	require.Equal(t, getResp.ID.Version(), deleteResp.ID.Version())
	require.Equal(t, getResp.KID, deleteResp.KID)
	require.Equal(t, getResp.SID, deleteResp.SID)
	testSerde(t, &deleteResp.DeletedCertificate)

	var getDeletedResp azcertificates.GetDeletedCertificateResponse
	pollStatus(t, http.StatusNotFound, func() error {
		getDeletedResp, err = client.GetDeletedCertificate(ctx, certName, nil)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, deleteResp.ID, getDeletedResp.ID)
	require.Equal(t, deleteResp.ID.Name(), getDeletedResp.ID.Name())
	require.Equal(t, deleteResp.ID.Version(), getDeletedResp.ID.Version())
	require.Equal(t, deleteResp.DeletedCertificate, getDeletedResp.DeletedCertificate)

	_, err = client.PurgeDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)
}

func TestDeleteRecover(t *testing.T) {
	client := startTest(t)
	certName := getName(t, "")
	createParams := azcertificates.CreateCertificateParameters{CertificatePolicy: &selfSignedPolicy}
	_, err := client.CreateCertificate(ctx, certName, createParams, nil)
	require.NoError(t, err)
	pollCertOperation(t, client, certName)

	deleteResp, err := client.DeleteCertificate(ctx, certName, nil)
	require.NoError(t, err)
	pollStatus(t, http.StatusNotFound, func() error {
		_, err = client.GetDeletedCertificate(ctx, certName, nil)
		return err
	})

	recoverResp, err := client.RecoverDeletedCertificate(ctx, certName, nil)
	require.NoError(t, err)
	pollStatus(t, http.StatusNotFound, func() error {
		_, err = client.GetCertificate(ctx, certName, "", nil)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, deleteResp.Attributes, recoverResp.Attributes)
	require.Equal(t, deleteResp.ID, recoverResp.ID)
	require.Equal(t, deleteResp.ID.Name(), recoverResp.ID.Name())
	require.Equal(t, deleteResp.ID.Version(), recoverResp.ID.Version())
	require.Equal(t, deleteResp.Policy, recoverResp.Policy)
	cleanUpCert(t, client, certName)
}

func TestDisableChallengeResourceVerification(t *testing.T) {
	authResource := `"Bearer authorization="https://login.microsoftonline.com/tenant", resource="%s""`
	authScope := `"Bearer authorization="https://login.microsoftonline.com/tenant", scope="%s""`
	vaultURL := "https://fakevault.vault.azure.net"
	for _, test := range []struct {
		challenge, resource string
		disableVerify, err  bool
	}{
		// happy path: resource matches requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.net"},
		{challenge: authScope, resource: "https://vault.azure.net/.default"},
		{challenge: authResource, resource: "https://vault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.net/.default", disableVerify: true},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.cn", err: true},
		{challenge: authResource, resource: "https://myvault.azure.net", err: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", err: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", err: true},

		// the policy shouldn't return errors for the above error cases when verification is disabled
		{challenge: authResource, resource: "https://vault.azure.cn", disableVerify: true},
		{challenge: authResource, resource: "https://myvault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", disableVerify: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", disableVerify: true},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader("WWW-Authenticate", fmt.Sprintf(test.challenge, test.resource)))
			srv.AppendResponse(mock.WithStatusCode(200), mock.WithBody([]byte(`{"value":[]}`)))
			options := &azcertificates.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Transport: srv,
				},
				DisableChallengeResourceVerification: test.disableVerify,
			}
			client, err := azcertificates.NewClient(vaultURL, &azcred.Fake{}, options)
			require.NoError(t, err)
			pager := client.NewListCertificatePropertiesPager(nil)
			_, err = pager.NextPage(context.Background())
			if test.err {
				require.Error(t, err)
				require.Contains(t, err.Error(), "challenge resource")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestID(t *testing.T) {
	for _, test := range []struct{ ID, name, version string }{
		{"https://foo.vault.azure.net/certificates/name/version", "name", "version"},
		{"https://foo.vault.azure.net/certificates/name", "name", ""},
	} {
		t.Run(test.ID, func(t *testing.T) {
			ID := azcertificates.ID(test.ID)
			require.Equal(t, test.name, ID.Name())
			require.Equal(t, test.version, ID.Version())
		})
	}
}

func TestImportCertificate(t *testing.T) {
	client := startTest(t)
	certName := getName(t, "")
	importParams := azcertificates.ImportCertificateParameters{
		Base64EncodedCertificate: to.Ptr("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"),
		CertificateAttributes:    &azcertificates.CertificateAttributes{Enabled: to.Ptr(false)},
		PreserveCertOrder:        to.Ptr(true),
	}
	testSerde(t, &importParams)
	impResp, err := client.ImportCertificate(ctx, certName, importParams, nil)
	require.NoError(t, err)
	cleanUpCert(t, client, certName)
	require.Equal(t, importParams.CertificateAttributes.Enabled, impResp.Attributes.Enabled)
	require.NotEmpty(t, impResp.ID)
	require.Equal(t, certName, impResp.ID.Name())
	require.NotEmpty(t, impResp.KID)
	require.NotEmpty(t, impResp.SID)
	require.True(t, *impResp.PreserveCertOrder)
}

func TestIssuerCRUD(t *testing.T) {
	client := startTest(t)

	issuerName := getName(t, "issuer")
	setParams := azcertificates.SetIssuerParameters{
		Attributes: &azcertificates.IssuerAttributes{
			Enabled: to.Ptr(true),
		},
		Credentials: &azcertificates.IssuerCredentials{
			AccountID: to.Ptr("keyvaultuser"),
		},
		OrganizationDetails: &azcertificates.OrganizationDetails{
			AdminContacts: []*azcertificates.AdministratorContact{
				{
					FirstName: to.Ptr("First"),
					LastName:  to.Ptr("Last"),
					Email:     to.Ptr("foo@test.com"),
					Phone:     to.Ptr("42"),
				},
			},
		},
		Provider: to.Ptr("Test"),
	}
	testSerde(t, &setParams)
	setResp, err := client.SetIssuer(ctx, issuerName, setParams, &azcertificates.SetIssuerOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, setResp.ID)
	require.Equal(t, setParams.Credentials, setResp.Credentials)
	require.Equal(t, setParams.OrganizationDetails.AdminContacts[0], setResp.OrganizationDetails.AdminContacts[0])
	require.Equal(t, setParams.Provider, setResp.Provider)
	testSerde(t, &setResp.Issuer)

	getResp, err := client.GetIssuer(ctx, issuerName, nil)
	require.NoError(t, err)
	require.Equal(t, setResp.Issuer, getResp.Issuer)

	pager := client.NewListIssuerPropertiesPager(nil)
	found := false
	for pager.More() {
		page, err := pager.NextPage(ctx)
		testSerde(t, &page.IssuerPropertiesListResult)
		require.NoError(t, err)
		for _, issuer := range page.Value {
			testSerde(t, issuer)
			require.NotEmpty(t, issuer.ID)
			if *issuer.ID == *setResp.ID {
				found = true
				break
			}
		}
	}
	require.True(t, found)

	updateParams := azcertificates.UpdateIssuerParameters{
		Attributes: &azcertificates.IssuerAttributes{
			Enabled: to.Ptr(false),
		},
	}
	testSerde(t, &updateParams)
	updateResp, err := client.UpdateIssuer(ctx, issuerName, updateParams, nil)
	require.NoError(t, err)
	require.NotEqual(t, setResp.Issuer, updateResp.Issuer)

	deleteResp, err := client.DeleteIssuer(ctx, issuerName, nil)
	require.NoError(t, err)
	require.Equal(t, updateResp.Issuer, deleteResp.Issuer)
}

func TestListCertificates(t *testing.T) {
	client := startTest(t)

	tag := getName(t, "")
	count := 4
	certNames := make([]string, count)
	createParams := azcertificates.CreateCertificateParameters{
		CertificatePolicy: &selfSignedPolicy,
		Tags:              map[string]*string{tag: to.Ptr("yes")},
	}
	for i := 0; i < len(certNames); i++ {
		certNames[i] = fmt.Sprintf("%s-%d", tag, i)
		_, err := client.CreateCertificate(ctx, certNames[i], createParams, nil)
		require.NoError(t, err)
	}
	for _, name := range certNames {
		pollCertOperation(t, client, name)
	}

	listCertsPager := client.NewListCertificatePropertiesPager(&azcertificates.ListCertificatePropertiesOptions{
		IncludePending: to.Ptr(true),
	})
	for listCertsPager.More() {
		page, err := listCertsPager.NextPage(ctx)
		require.NoError(t, err)
		testSerde(t, &page.CertificatePropertiesListResult)
		for _, cert := range page.Value {
			testSerde(t, cert)
			if value, ok := cert.Tags[tag]; ok && *value == "yes" {
				require.True(t, strings.HasPrefix(cert.ID.Name(), tag))
				count--
				_, err = client.DeleteCertificate(ctx, cert.ID.Name(), nil)
				require.NoError(t, err)
			}
		}
	}
	require.Equal(t, 0, count)

	for _, name := range certNames {
		pollStatus(t, http.StatusNotFound, func() error {
			_, err := client.GetDeletedCertificate(ctx, name, nil)
			return err
		})
	}

	count = len(certNames)
	listDeletedCertsPager := client.NewListDeletedCertificatePropertiesPager(&azcertificates.ListDeletedCertificatePropertiesOptions{
		IncludePending: to.Ptr(true),
	})
	for listDeletedCertsPager.More() {
		page, err := listDeletedCertsPager.NextPage(ctx)
		require.NoError(t, err)
		testSerde(t, &page.DeletedCertificatePropertiesListResult)
		for _, cert := range page.Value {
			testSerde(t, cert)
			if value, ok := cert.Tags[tag]; ok && *value == "yes" {
				count--
				_, err = client.PurgeDeletedCertificate(ctx, cert.ID.Name(), nil)
				require.NoError(t, err)
			}
		}
	}
	require.Equal(t, 0, count)
}

func TestListCertificateVersions(t *testing.T) {
	client := startTest(t)

	name := getName(t, "")
	count := 3
	for i := 0; i < count; i++ {
		_, err := client.CreateCertificate(ctx, name, azcertificates.CreateCertificateParameters{CertificatePolicy: &selfSignedPolicy}, nil)
		require.NoError(t, err)
		pollCertOperation(t, client, name)
	}
	defer cleanUpCert(t, client, name)

	pager := client.NewListCertificatePropertiesVersionsPager(name, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		testSerde(t, &page.CertificatePropertiesListResult)
		count -= len(page.Value)
		for _, v := range page.Value {
			testSerde(t, v)
			require.Equal(t, name, v.ID.Name())
		}
	}
	require.Equal(t, count, 0)
}

// https://stackoverflow.com/questions/42643048/signing-certificate-request-with-certificate-authority
// Much of this is thanks to this response, thanks @krostar
func TestMergeCertificate(t *testing.T) {
	client := startTest(t)

	certName := getName(t, "mergeCertificate")
	policy := azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{
			Name:                    to.Ptr("Unknown"),
			CertificateTransparency: to.Ptr(false),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr("CN=MyCert"),
		},
	}
	_, err := client.CreateCertificate(ctx, certName, azcertificates.CreateCertificateParameters{CertificatePolicy: &policy}, nil)
	require.NoError(t, err)
	defer cleanUpCert(t, client, certName)

	certOpResp, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)

	data, err := os.ReadFile("testdata/ca.crt")
	require.NoError(t, err)
	block, _ := pem.Decode(data)
	require.NotNil(t, block)
	caCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	data, err = os.ReadFile("testdata/ca.key")
	require.NoError(t, err)
	pkeyBlock, _ := pem.Decode(data)
	require.NotNil(t, pkeyBlock)
	require.Equal(t, pkeyBlock.Type, "RSA PRIVATE KEY")
	pkey, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes)
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
		NotBefore:          time.Date(2030, 1, 1, 1, 1, 0, 0, time.UTC),
		NotAfter:           time.Date(2040, 1, 1, 1, 1, 0, 0, time.UTC),
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
	mergeParams := azcertificates.MergeCertificateParameters{X509Certificates: [][]byte{[]byte(certificateString)}}
	testSerde(t, &mergeParams)
	mergeResp, err := client.MergeCertificate(ctx, certName, mergeParams, nil)
	require.NoError(t, err)
	require.NotNil(t, mergeResp.Policy)
}

func TestOperationCRUD(t *testing.T) {
	client := startTest(t)

	certName := getName(t, "")
	policy := azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{
			Name:                    to.Ptr("Unknown"),
			CertificateTransparency: to.Ptr(false),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr("CN=MyCert"),
		},
	}
	createParams := azcertificates.CreateCertificateParameters{CertificatePolicy: &policy}
	_, err := client.CreateCertificate(ctx, certName, createParams, nil)
	require.NoError(t, err)

	params := azcertificates.UpdateCertificateOperationParameter{CancellationRequested: to.Ptr(true)}
	testSerde(t, &params)
	_, err = client.UpdateCertificateOperation(ctx, certName, params, nil)
	require.NoError(t, err)

	getResp, err := client.GetCertificateOperation(ctx, certName, nil)
	require.NoError(t, err)
	require.Equal(t, params.CancellationRequested, getResp.CancellationRequested)
	testSerde(t, &getResp.CertificateOperation)

	pollStatus(t, http.StatusConflict, func() error {
		// Key Vault returns an error when the update is slow or this delete
		// is fast such that the delete executes before the update completes
		_, err = client.DeleteCertificateOperation(ctx, certName, nil)
		return err
	})
	require.NoError(t, err)
}

func TestUpdateCertificatePolicy(t *testing.T) {
	client := startTest(t)

	certName := getName(t, "")
	policy := azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{
			CertificateTransparency: to.Ptr(false),
			Name:                    to.Ptr("Self"),
		},
		Attributes: &azcertificates.CertificateAttributes{
			Enabled:   to.Ptr(true),
			Expires:   to.Ptr(time.Date(2040, 1, 1, 1, 1, 0, 0, time.UTC)),
			NotBefore: to.Ptr(time.Date(2030, 1, 1, 1, 1, 0, 0, time.UTC)),
		},
		KeyProperties: &azcertificates.KeyProperties{
			Exportable: to.Ptr(true),
			KeySize:    to.Ptr(int32(2048)),
			KeyType:    to.Ptr(azcertificates.KeyTypeRSA),
			ReuseKey:   to.Ptr(true),
		},
		LifetimeActions: []*azcertificates.LifetimeAction{
			{
				Action: &azcertificates.LifetimeActionType{
					ActionType: to.Ptr(azcertificates.CertificatePolicyActionEmailContacts),
				},
				Trigger: &azcertificates.LifetimeActionTrigger{
					LifetimePercentage: to.Ptr(int32(98)),
				},
			},
		},
		SecretProperties: &azcertificates.SecretProperties{ContentType: to.Ptr("application/x-pkcs12")},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			EnhancedKeyUsage: []*string{to.Ptr("1.3.6.1.5.5.7.3.1"), to.Ptr("1.3.6.1.5.5.7.3.2")},
			KeyUsage:         []*azcertificates.KeyUsageType{to.Ptr(azcertificates.KeyUsageTypeDataEncipherment)},
			Subject:          to.Ptr("CN=DefaultPolicy"),
			SubjectAlternativeNames: &azcertificates.SubjectAlternativeNames{
				DNSNames: []*string{to.Ptr("localhost")},
			},
			ValidityInMonths: to.Ptr(int32(12)),
		},
	}
	_, err := client.CreateCertificate(ctx, certName, azcertificates.CreateCertificateParameters{CertificatePolicy: &policy}, nil)
	require.NoError(t, err)
	defer cleanUpCert(t, client, certName)

	getResp, err := client.GetCertificatePolicy(ctx, certName, nil)
	require.NoError(t, err)
	require.Equal(t, policy.IssuerParameters, getResp.IssuerParameters)
	require.Equal(t, policy.KeyProperties, getResp.KeyProperties)
	require.Equal(t, policy.LifetimeActions, getResp.LifetimeActions)
	require.Equal(t, policy.SecretProperties, getResp.SecretProperties)
	require.Equal(t, policy.X509CertificateProperties, getResp.X509CertificateProperties)
	updatedPolicy := azcertificates.CertificatePolicy{
		KeyProperties: &azcertificates.KeyProperties{
			Curve:      to.Ptr(azcertificates.CurveNameP256K),
			Exportable: to.Ptr(true),
			KeySize:    to.Ptr(int32(256)),
			KeyType:    to.Ptr(azcertificates.KeyTypeEC),
			ReuseKey:   to.Ptr(false),
		},
	}
	updateResp, err := client.UpdateCertificatePolicy(ctx, certName, updatedPolicy, nil)
	require.NoError(t, err)
	require.Equal(t, policy.IssuerParameters, updateResp.CertificatePolicy.IssuerParameters)
	require.Equal(t, updatedPolicy.KeyProperties, updateResp.CertificatePolicy.KeyProperties)
	require.Equal(t, policy.LifetimeActions, updateResp.CertificatePolicy.LifetimeActions)
	require.Equal(t, policy.SecretProperties, updateResp.CertificatePolicy.SecretProperties)
	require.Equal(t, policy.X509CertificateProperties, updateResp.CertificatePolicy.X509CertificateProperties)
}

func TestAPIVersion(t *testing.T) {
	apiVersion := "7.3"
	var requireVersion = func(req *http.Request) bool {
		version := req.URL.Query().Get("api-version")
		require.Equal(t, version, apiVersion)
		return true
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(200),
		mock.WithPredicate(requireVersion),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))

	opts := &azcertificates.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := azcertificates.NewClient(vaultURL, &azcred.Fake{}, opts)
	require.NoError(t, err)

	_, err = client.GetCertificate(context.Background(), "name", "", nil)
	require.NoError(t, err)
}

func TestSubjectAlternativeNames(t *testing.T) {
	client := startTest(t)

	// Test certificate with all SAN fields including new IPAddresses and URIs
	certName := getName(t, "allsans")
	policy := azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{Name: to.Ptr("self")},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr("CN=SANTest"),
			SubjectAlternativeNames: &azcertificates.SubjectAlternativeNames{
				DNSNames:           []*string{to.Ptr("localhost"), to.Ptr("example.com")},
				Emails:             []*string{to.Ptr("admin@example.com")},
				IPAddresses:        []*string{to.Ptr("192.168.1.1"), to.Ptr("2001:0db8::1")}, // IPv4 and IPv6
				URIs:               []*string{to.Ptr("https://example.com"), to.Ptr("https://test.com/path")},
				UserPrincipalNames: []*string{to.Ptr("user@domain.com")},
			},
		},
	}

	// Create certificate and verify SANs are preserved
	createParams := azcertificates.CreateCertificateParameters{CertificatePolicy: &policy}
	testSerde(t, &createParams)
	_, err := client.CreateCertificate(ctx, certName, createParams, nil)
	require.NoError(t, err)
	pollCertOperation(t, client, certName)
	defer cleanUpCert(t, client, certName)

	// Verify all SANs are retrieved correctly
	getResp, err := client.GetCertificate(ctx, certName, "", nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.Policy)
	require.NotNil(t, getResp.Policy.X509CertificateProperties)
	require.NotNil(t, getResp.Policy.X509CertificateProperties.SubjectAlternativeNames)
	testSerde(t, getResp.Policy.X509CertificateProperties.SubjectAlternativeNames)
}
