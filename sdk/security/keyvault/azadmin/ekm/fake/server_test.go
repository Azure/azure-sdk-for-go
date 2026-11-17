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

func getServer() fake.Server {
	return fake.Server{
		CheckConnection: func(ctx context.Context, options *ekm.CheckConnectionOptions) (resp azfake.Responder[ekm.CheckConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.CheckConnectionResponse{ProxyInfo: ekm.ProxyInfo{
				APIVersion:  to.Ptr(ekmAPIVersion),
				EKMProduct:  to.Ptr(ekmProduct),
				EKMVendor:   to.Ptr(ekmVendor),
				ProxyName:   to.Ptr(proxyName),
				ProxyVendor: to.Ptr(proxyVendor),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		CreateConnection: func(ctx context.Context, ekmConnection ekm.Connection, options *ekm.CreateConnectionOptions) (resp azfake.Responder[ekm.CreateConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.CreateConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteConnection: func(ctx context.Context, options *ekm.DeleteConnectionOptions) (resp azfake.Responder[ekm.DeleteConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.DeleteConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetCertificate: func(ctx context.Context, options *ekm.GetCertificateOptions) (resp azfake.Responder[ekm.GetCertificateResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.GetCertificateResponse{ProxyClientCertificateInfo: ekm.ProxyClientCertificateInfo{
				CaCertificates:    clientCACerts,
				SubjectCommonName: to.Ptr(clientCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetConnection: func(ctx context.Context, options *ekm.GetConnectionOptions) (resp azfake.Responder[ekm.GetConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.GetConnectionResponse{Connection: ekm.Connection{
				Host:                    to.Ptr(ekmHost),
				ServerCaCertificates:    ekmCACerts,
				ServerSubjectCommonName: to.Ptr(ekmCommonName),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateConnection: func(ctx context.Context, ekmConnection ekm.Connection, options *ekm.UpdateConnectionOptions) (resp azfake.Responder[ekm.UpdateConnectionResponse], errResp azfake.ErrorResponder) {
			kvResp := ekm.UpdateConnectionResponse{Connection: ekm.Connection{
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
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	// check
	checkResp, err := client.CheckConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmAPIVersion, *checkResp.APIVersion)
	require.Equal(t, ekmProduct, *checkResp.EKMProduct)
	require.Equal(t, ekmVendor, *checkResp.EKMVendor)
	require.Equal(t, proxyName, *checkResp.ProxyName)
	require.Equal(t, proxyVendor, *checkResp.ProxyVendor)

	// create
	createResp, err := client.CreateConnection(context.TODO(), ekm.Connection{}, nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *createResp.Host)
	require.Equal(t, ekmCommonName, *createResp.ServerSubjectCommonName)
	require.Equal(t, ekmCACerts, createResp.ServerCaCertificates)

	// get
	getResp, err := client.GetConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *getResp.Host)
	require.Equal(t, ekmCommonName, *getResp.ServerSubjectCommonName)

	// update
	updateResp, err := client.UpdateConnection(context.TODO(), ekm.Connection{}, nil)
	require.NoError(t, err)
	require.Equal(t, "updated.contoso.com", *updateResp.ServerSubjectCommonName)

	// get certificate
	certResp, err := client.GetCertificate(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, clientCommonName, *certResp.SubjectCommonName)
	require.Equal(t, clientCACerts, certResp.CaCertificates)

	// delete
	deleteResp, err := client.DeleteConnection(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, ekmHost, *deleteResp.Host)
}
