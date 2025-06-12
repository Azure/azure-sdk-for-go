// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"fmt"

	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates/fake"
	"github.com/stretchr/testify/require"
)

var (
	version     = "123"
	certName    = "certName"
	issuerName  = "issuerName"
	expiredTime = time.Date(2030, 1, 1, 1, 1, 1, 0, time.Local)
	contact1    = azcertificates.Contact{Email: to.Ptr("one@localhost"), Name: to.Ptr("One"), Phone: to.Ptr("1111111111")}
	contact2    = azcertificates.Contact{Email: to.Ptr("two@localhost"), Name: to.Ptr("Two"), Phone: to.Ptr("2222222222")}
)

func getServer() fake.Server {
	return fake.Server{
		BackupCertificate: func(ctx context.Context, name string, options *azcertificates.BackupCertificateOptions) (resp azfake.Responder[azcertificates.BackupCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.BackupCertificateResponse{
				BackupCertificateResult: azcertificates.BackupCertificateResult{
					Value: []byte("testing"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		CreateCertificate: func(ctx context.Context, name string, parameters azcertificates.CreateCertificateParameters, options *azcertificates.CreateCertificateOptions) (resp azfake.Responder[azcertificates.CreateCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.CreateCertificateResponse{
				CertificateOperation: azcertificates.CertificateOperation{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, "pending"))),
				},
			}
			resp.SetResponse(http.StatusAccepted, kvResp, nil)
			return
		},
		DeleteCertificate: func(ctx context.Context, name string, options *azcertificates.DeleteCertificateOptions) (resp azfake.Responder[azcertificates.DeleteCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.DeleteCertificateResponse{
				DeletedCertificate: azcertificates.DeletedCertificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteCertificateOperation: func(ctx context.Context, name string, options *azcertificates.DeleteCertificateOperationOptions) (resp azfake.Responder[azcertificates.DeleteCertificateOperationResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.DeleteCertificateOperationResponse{
				CertificateOperation: azcertificates.CertificateOperation{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteContacts: func(ctx context.Context, options *azcertificates.DeleteContactsOptions) (resp azfake.Responder[azcertificates.DeleteContactsResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.DeleteContactsResponse{
				Contacts: azcertificates.Contacts{
					ContactList: []*azcertificates.Contact{
						{Email: to.Ptr("one@localhost"), Name: to.Ptr("One"), Phone: to.Ptr("1111111111")},
						{Email: to.Ptr("two@localhost"), Name: to.Ptr("Two"), Phone: to.Ptr("2222222222")},
					},
					ID: to.Ptr("https://fake.vault.azure.net/certificates/contacts"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteIssuer: func(ctx context.Context, name string, options *azcertificates.DeleteIssuerOptions) (resp azfake.Responder[azcertificates.DeleteIssuerResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.DeleteIssuerResponse{
				Issuer: azcertificates.Issuer{
					ID: to.Ptr(fmt.Sprintf("https://fake.vault.azure.net/certificates/issuers/%s", name)),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetCertificate: func(ctx context.Context, name string, version string, options *azcertificates.GetCertificateOptions) (resp azfake.Responder[azcertificates.GetCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetCertificateOperation: func(ctx context.Context, name string, options *azcertificates.GetCertificateOperationOptions) (resp azfake.Responder[azcertificates.GetCertificateOperationResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetCertificateOperationResponse{
				CertificateOperation: azcertificates.CertificateOperation{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, "pending"))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetCertificatePolicy: func(ctx context.Context, name string, options *azcertificates.GetCertificatePolicyOptions) (resp azfake.Responder[azcertificates.GetCertificatePolicyResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetCertificatePolicyResponse{
				CertificatePolicy: azcertificates.CertificatePolicy{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetContacts: func(ctx context.Context, options *azcertificates.GetContactsOptions) (resp azfake.Responder[azcertificates.GetContactsResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetContactsResponse{
				Contacts: azcertificates.Contacts{
					ContactList: []*azcertificates.Contact{
						&contact1,
						&contact2,
					},
					ID: to.Ptr("https://fake.vault.azure.net/certificates/contacts"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetDeletedCertificate: func(ctx context.Context, name string, options *azcertificates.GetDeletedCertificateOptions) (resp azfake.Responder[azcertificates.GetDeletedCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetDeletedCertificateResponse{
				DeletedCertificate: azcertificates.DeletedCertificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetIssuer: func(ctx context.Context, name string, options *azcertificates.GetIssuerOptions) (resp azfake.Responder[azcertificates.GetIssuerResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.GetIssuerResponse{
				Issuer: azcertificates.Issuer{
					ID: to.Ptr(fmt.Sprintf("https://fake.vault.azure.net/certificates/issuers/%s", name)),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		ImportCertificate: func(ctx context.Context, name string, parameters azcertificates.ImportCertificateParameters, options *azcertificates.ImportCertificateOptions) (resp azfake.Responder[azcertificates.ImportCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.ImportCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
					Attributes: &azcertificates.CertificateAttributes{
						Expires: parameters.CertificateAttributes.Expires,
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		NewListCertificatePropertiesPager: func(options *azcertificates.ListCertificatePropertiesOptions) (resp azfake.PagerResponder[azcertificates.ListCertificatePropertiesResponse]) {
			page1 := azcertificates.ListCertificatePropertiesResponse{
				CertificatePropertiesListResult: azcertificates.CertificatePropertiesListResult{
					Value: []*azcertificates.CertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert1/123")),
						},
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert2/123")),
						},
					},
				},
			}
			page2 := azcertificates.ListCertificatePropertiesResponse{
				CertificatePropertiesListResult: azcertificates.CertificatePropertiesListResult{
					Value: []*azcertificates.CertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert3/123")),
						},
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert3/123")),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListCertificatePropertiesVersionsPager: func(name string, options *azcertificates.ListCertificatePropertiesVersionsOptions) (resp azfake.PagerResponder[azcertificates.ListCertificatePropertiesVersionsResponse]) {
			page1 := azcertificates.ListCertificatePropertiesVersionsResponse{
				CertificatePropertiesListResult: azcertificates.CertificatePropertiesListResult{
					Value: []*azcertificates.CertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%d", name, 1))),
						},
						{
							ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%d", name, 2))),
						},
					},
				},
			}
			page2 := azcertificates.ListCertificatePropertiesVersionsResponse{
				CertificatePropertiesListResult: azcertificates.CertificatePropertiesListResult{
					Value: []*azcertificates.CertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%d", name, 3))),
						},
						{
							ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%d", name, 4))),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListDeletedCertificatePropertiesPager: func(options *azcertificates.ListDeletedCertificatePropertiesOptions) (resp azfake.PagerResponder[azcertificates.ListDeletedCertificatePropertiesResponse]) {
			page1 := azcertificates.ListDeletedCertificatePropertiesResponse{
				DeletedCertificatePropertiesListResult: azcertificates.DeletedCertificatePropertiesListResult{
					Value: []*azcertificates.DeletedCertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert1/123")),
						},
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert2/123")),
						},
					},
				},
			}
			page2 := azcertificates.ListDeletedCertificatePropertiesResponse{
				DeletedCertificatePropertiesListResult: azcertificates.DeletedCertificatePropertiesListResult{
					Value: []*azcertificates.DeletedCertificateProperties{
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert3/123")),
						},
						{
							ID: to.Ptr(azcertificates.ID("https://fake-vault.vault.azure.net/certificates/cert4/123")),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListIssuerPropertiesPager: func(options *azcertificates.ListIssuerPropertiesOptions) (resp azfake.PagerResponder[azcertificates.ListIssuerPropertiesResponse]) {
			page1 := azcertificates.ListIssuerPropertiesResponse{
				IssuerPropertiesListResult: azcertificates.IssuerPropertiesListResult{
					Value: []*azcertificates.IssuerProperties{
						{
							ID: to.Ptr("name1"),
						},
						{
							ID: to.Ptr("name2"),
						},
					},
				},
			}
			page2 := azcertificates.ListIssuerPropertiesResponse{
				IssuerPropertiesListResult: azcertificates.IssuerPropertiesListResult{
					Value: []*azcertificates.IssuerProperties{
						{
							ID: to.Ptr("name3"),
						},
						{
							ID: to.Ptr("name4"),
						},
					},
				},
			}
			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		MergeCertificate: func(ctx context.Context, name string, parameters azcertificates.MergeCertificateParameters, options *azcertificates.MergeCertificateOptions) (resp azfake.Responder[azcertificates.MergeCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.MergeCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusCreated, kvResp, nil)
			return
		},
		PurgeDeletedCertificate: func(ctx context.Context, name string, options *azcertificates.PurgeDeletedCertificateOptions) (resp azfake.Responder[azcertificates.PurgeDeletedCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.PurgeDeletedCertificateResponse{}
			resp.SetResponse(http.StatusNoContent, kvResp, nil)
			return
		},
		RecoverDeletedCertificate: func(ctx context.Context, name string, options *azcertificates.RecoverDeletedCertificateOptions) (resp azfake.Responder[azcertificates.RecoverDeletedCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.RecoverDeletedCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		RestoreCertificate: func(ctx context.Context, parameters azcertificates.RestoreCertificateParameters, options *azcertificates.RestoreCertificateOptions) (resp azfake.Responder[azcertificates.RestoreCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.RestoreCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", certName, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		SetContacts: func(ctx context.Context, contacts azcertificates.Contacts, options *azcertificates.SetContactsOptions) (resp azfake.Responder[azcertificates.SetContactsResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.SetContactsResponse{
				Contacts: contacts,
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		SetIssuer: func(ctx context.Context, name string, parameter azcertificates.SetIssuerParameters, options *azcertificates.SetIssuerOptions) (resp azfake.Responder[azcertificates.SetIssuerResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.SetIssuerResponse{
				Issuer: azcertificates.Issuer{
					ID: to.Ptr(fmt.Sprintf("https://fake.vault.azure.net/certificates/issuers/%s", name)),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateCertificate: func(ctx context.Context, name string, version string, parameters azcertificates.UpdateCertificateParameters, options *azcertificates.UpdateCertificateOptions) (resp azfake.Responder[azcertificates.UpdateCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.UpdateCertificateResponse{
				Certificate: azcertificates.Certificate{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
					Attributes: &azcertificates.CertificateAttributes{
						Expires: parameters.CertificateAttributes.Expires,
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateCertificateOperation: func(ctx context.Context, name string, certificateOperation azcertificates.UpdateCertificateOperationParameter, options *azcertificates.UpdateCertificateOperationOptions) (resp azfake.Responder[azcertificates.UpdateCertificateOperationResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.UpdateCertificateOperationResponse{
				CertificateOperation: azcertificates.CertificateOperation{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, "pending"))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateCertificatePolicy: func(ctx context.Context, name string, certificatePolicy azcertificates.CertificatePolicy, options *azcertificates.UpdateCertificatePolicyOptions) (resp azfake.Responder[azcertificates.UpdateCertificatePolicyResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.UpdateCertificatePolicyResponse{
				CertificatePolicy: azcertificates.CertificatePolicy{
					ID: to.Ptr(azcertificates.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/certificates/%s/%s", name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateIssuer: func(ctx context.Context, name string, parameter azcertificates.UpdateIssuerParameters, options *azcertificates.UpdateIssuerOptions) (resp azfake.Responder[azcertificates.UpdateIssuerResponse], errResp azfake.ErrorResponder) {
			kvResp := azcertificates.UpdateIssuerResponse{
				Issuer: azcertificates.Issuer{
					ID: to.Ptr(fmt.Sprintf("https://fake.vault.azure.net/certificates/issuers/%s", name)),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := azcertificates.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &azcertificates.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	// backup certificate
	backupResp, err := client.BackupCertificate(context.Background(), certName, nil)
	require.NoError(t, err)
	require.NotNil(t, backupResp.Value)

	// create certificate
	createResp, err := client.CreateCertificate(context.Background(), certName, azcertificates.CreateCertificateParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, certName, createResp.ID.Name())
	require.Equal(t, "pending", createResp.ID.Version())

	// delete certificate
	deleteResp, err := client.DeleteCertificate(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, deleteResp.ID.Name())

	// delete certificate operation
	deleteOperationResp, err := client.DeleteCertificateOperation(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, deleteOperationResp.ID.Name())

	// delete contacts
	deleteContactsResp, err := client.DeleteContacts(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, deleteContactsResp.ContactList, 2)
	require.NotNil(t, deleteContactsResp.ID)

	// delete issuer
	deleteIssuerResp, err := client.DeleteIssuer(context.Background(), issuerName, nil)
	require.NoError(t, err)
	require.Contains(t, *deleteIssuerResp.ID, issuerName)

	// get certificate operation
	getOpResp, err := client.GetCertificateOperation(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, getOpResp.ID.Name())
	require.Equal(t, "pending", getOpResp.ID.Version())

	// get certificate
	getResp, err := client.GetCertificate(context.Background(), certName, "", nil)
	require.NoError(t, err)
	require.Equal(t, certName, getResp.ID.Name())
	require.Empty(t, getResp.ID.Version())

	// get certificate policy
	getCertPolicyResp, err := client.GetCertificatePolicy(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, getCertPolicyResp.ID.Name())

	// get contacts
	getContactsResp, err := client.GetContacts(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, getContactsResp.ContactList, 2)
	require.NotNil(t, getContactsResp.ID)

	// get deleted certificate
	getDeletedCertResp, err := client.GetDeletedCertificate(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, getDeletedCertResp.ID.Name())

	// get issuer
	getIssuerResp, err := client.GetIssuer(context.Background(), issuerName, nil)
	require.NoError(t, err)
	require.Contains(t, *getIssuerResp.ID, issuerName)

	// import certificate
	importParams := azcertificates.ImportCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Expires: to.Ptr(expiredTime),
		},
	}
	importResp, err := client.ImportCertificate(context.Background(), certName, importParams, nil)
	require.NoError(t, err)
	require.Equal(t, certName, importResp.ID.Name())
	require.Equal(t, version, importResp.ID.Version())
	require.True(t, expiredTime.Equal(*importResp.Attributes.Expires))

	// list cert properties
	certPropertiesPager := client.NewListCertificatePropertiesPager(nil)
	for certPropertiesPager.More() {
		page, err := certPropertiesPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, cert := range page.Value {
			require.NotNil(t, cert.ID)
			require.Contains(t, cert.ID.Name(), "cert")
		}
	}

	// list cert properties versions
	certPropertiesVersionsPager := client.NewListCertificatePropertiesVersionsPager(certName, nil)
	for certPropertiesVersionsPager.More() {
		page, err := certPropertiesVersionsPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, cert := range page.Value {
			require.NotNil(t, cert.ID)
			require.Equal(t, cert.ID.Name(), "certName")
		}
	}

	// list deleted cert properties
	deletedCertPropertiesPager := client.NewListDeletedCertificatePropertiesPager(nil)
	for certPropertiesPager.More() {
		page, err := deletedCertPropertiesPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, cert := range page.Value {
			require.NotNil(t, cert.ID)
			require.Contains(t, cert.ID.Name(), "cert")
		}
	}

	issuerPager := client.NewListIssuerPropertiesPager(nil)
	for certPropertiesPager.More() {
		page, err := issuerPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, issuer := range page.Value {
			require.Contains(t, issuer.ID, "name")
		}
	}

	// merge certificate
	mergeResp, err := client.MergeCertificate(context.Background(), certName, azcertificates.MergeCertificateParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, certName, mergeResp.ID.Name())
	require.Equal(t, version, mergeResp.ID.Version())

	// purge deleted certificate
	purgeResp, err := client.PurgeDeletedCertificate(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Empty(t, purgeResp)

	recoverResp, err := client.RecoverDeletedCertificate(context.Background(), certName, nil)
	require.NoError(t, err)
	require.Equal(t, certName, recoverResp.ID.Name())

	restoreResp, err := client.RestoreCertificate(context.Background(), azcertificates.RestoreCertificateParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, certName, restoreResp.ID.Name())

	contacts := azcertificates.Contacts{
		ContactList: []*azcertificates.Contact{&contact1, &contact2},
		ID:          to.Ptr("https://fake.vault.azure.net/certificates/contacts"),
	}
	setContactsResp, err := client.SetContacts(context.Background(), contacts, nil)
	require.NoError(t, err)
	require.Len(t, setContactsResp.ContactList, 2)
	require.NotNil(t, setContactsResp.ID)

	setIssuerResp, err := client.SetIssuer(context.Background(), issuerName, azcertificates.SetIssuerParameters{}, nil)
	require.NoError(t, err)
	require.Contains(t, *setIssuerResp.ID, issuerName)

	// update certificate
	updateParams := azcertificates.UpdateCertificateParameters{
		CertificateAttributes: &azcertificates.CertificateAttributes{
			Expires: to.Ptr(time.Date(2030, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
	}
	updateResp, err := client.UpdateCertificate(context.Background(), certName, version, updateParams, nil)
	require.NoError(t, err)
	require.Equal(t, certName, updateResp.ID.Name())
	require.Equal(t, version, updateResp.ID.Version())

	updateOpResp, err := client.UpdateCertificateOperation(context.Background(), certName, azcertificates.UpdateCertificateOperationParameter{}, nil)
	require.NoError(t, err)
	require.Equal(t, certName, updateOpResp.ID.Name())
	require.Equal(t, "pending", updateOpResp.ID.Version())

	updateCertPolicyResp, err := client.UpdateCertificatePolicy(context.Background(), certName, azcertificates.CertificatePolicy{}, nil)
	require.NoError(t, err)
	require.Equal(t, certName, updateCertPolicyResp.ID.Name())

	updateIssuerResp, err := client.UpdateIssuer(context.Background(), issuerName, azcertificates.UpdateIssuerParameters{}, nil)
	require.NoError(t, err)
	require.Contains(t, *updateIssuerResp.ID, issuerName)
}
