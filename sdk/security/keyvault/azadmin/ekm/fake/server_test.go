// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm/fake"
	"github.com/stretchr/testify/require"
)

var (
	ekmHost       = "ekm.contoso.com:443"
	ekmCommonName = "ekm.contoso.com"
	ekmCACerts    = [][]byte{[]byte("fake-ca-certificate-bytes")}

	ekmAPIVersion = "1.0"
	ekmProduct    = "fake-ekm-product"
	ekmVendor     = "fake-ekm-vendor"
	proxyName     = "fake-proxy-name"
	proxyVendor   = "fake-proxy-vendor"

	clientCommonName = "client.contoso.com"
	clientCACerts    = [][]byte{[]byte("fake-client-ca-certificate-bytes")}
)

func getServer() fake.KeyVaultServer {
	return fake.KeyVaultServer{
		CheckEkmConnection: func(ctx context.Context, options *ekm.CheckEkmConnectionOptions) (resp azfake.Responder[ekm.CheckEkmConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.CheckEkmConnectionResponse{ProxyInfo: ekm.ProxyInfo{
				APIVersion:  to.Ptr(ekmAPIVersion),
				EkmProduct:  to.Ptr(ekmProduct),
				EkmVendor:   to.Ptr(ekmVendor),
				ProxyName:   to.Ptr(proxyName),
				ProxyVendor: to.Ptr(proxyVendor),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		CreateEkmConnection: func(ctx context.Context, ekmConnection ekm.Connection, options *ekm.CreateEkmConnectionOptions) (resp azfake.Responder[ekm.CreateEkmConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.CreateEkmConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteEkmConnection: func(ctx context.Context, options *ekm.DeleteEkmConnectionOptions) (resp azfake.Responder[ekm.DeleteEkmConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.DeleteEkmConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetEkmCertificate: func(ctx context.Context, options *ekm.GetEkmCertificateOptions) (resp azfake.Responder[ekm.GetEkmCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.GetEkmCertificateResponse{ProxyClientCertificateInfo: ekm.ProxyClientCertificateInfo{
				CaCertificates:    clientCACerts,
				SubjectCommonName: to.Ptr(clientCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetEkmConnection: func(ctx context.Context, options *ekm.GetEkmConnectionOptions) (resp azfake.Responder[ekm.GetEkmConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.GetEkmConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateEkmConnection: func(ctx context.Context, ekmConnection ekm.Connection, options *ekm.UpdateEkmConnectionOptions) (resp azfake.Responder[ekm.UpdateEkmConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.UpdateEkmConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr("updated.contoso.com"),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := ekm.NewClient("https://fake-hsm.managedhsm.azure.net", &azfake.TokenCredential{}, &ekm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewKeyVaultServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	// check
	checkResp, err := client.CheckEkmConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmAPIVersion, *checkResp.APIVersion)
	require.Equal(t, ekmProduct, *checkResp.EkmProduct)
	require.Equal(t, ekmVendor, *checkResp.EkmVendor)
	require.Equal(t, proxyName, *checkResp.ProxyName)
	require.Equal(t, proxyVendor, *checkResp.ProxyVendor)

	// create
	createResp, err := client.CreateEkmConnection(context.TODO(), ekm.Connection{}, nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *createResp.Host)
	require.Equal(t, ekmCommonName, *createResp.ServerSubjectCommonName)
	require.Equal(t, ekmCACerts, createResp.ServerCaCertificates)

	// get
	getResp, err := client.GetEkmConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *getResp.Host)
	require.Equal(t, ekmCommonName, *getResp.ServerSubjectCommonName)

	// update
	updateResp, err := client.UpdateEkmConnection(context.TODO(), ekm.Connection{}, nil)
	require.NoError(t, err)
	require.Equal(t, "updated.contoso.com", *updateResp.ServerSubjectCommonName)

	// get certificate
	certResp, err := client.GetEkmCertificate(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, clientCommonName, *certResp.SubjectCommonName)
	require.Equal(t, clientCACerts, certResp.CaCertificates)

	// delete
	deleteResp, err := client.DeleteEkmConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *deleteResp.Host)
}
