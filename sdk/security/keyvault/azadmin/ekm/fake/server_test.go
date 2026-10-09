// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

	privateEndpoint = "my-ekm-pe"
	plsAlias        = "myekmproxy-pls.00000000-0000-0000-0000-000000000000.eastus.azure.privatelinkservice"
	peJobID         = "fake-job-id"
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
		BeginCreateEkmPrivateEndpoint: func(ctx context.Context, peName string, parameters ekm.PrivateEndpointCreateParameters, options *ekm.BeginCreateEkmPrivateEndpointOptions) (resp azfake.PollerResponder[ekm.CreateEkmPrivateEndpointResponse], errResp azfake.ErrorResponder) {
			if parameters.PrivateLinkServiceID == nil || *parameters.PrivateLinkServiceID != plsAlias {
				errResp.SetResponseError(http.StatusBadRequest, "BadParameter")
				return
			}
			if peName != privateEndpoint {
				errResp.SetResponseError(http.StatusConflict, "Conflict")
				return
			}
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.SetTerminalResponse(http.StatusOK, ekm.CreateEkmPrivateEndpointResponse{PrivateEndpointOperation: ekm.PrivateEndpointOperation{
				JobID:               to.Ptr(peJobID),
				OperationType:       to.Ptr(ekm.EkmPrivateEndpointOperationTypeCreate),
				PrivateEndpointName: to.Ptr(peName),
				Status:              to.Ptr(ekm.OperationStateSucceeded),
			}}, nil)
			return
		},
		BeginDeleteEkmPrivateEndpoint: func(ctx context.Context, peName string, options *ekm.BeginDeleteEkmPrivateEndpointOptions) (resp azfake.PollerResponder[ekm.DeleteEkmPrivateEndpointResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.SetTerminalResponse(http.StatusOK, ekm.DeleteEkmPrivateEndpointResponse{PrivateEndpointOperation: ekm.PrivateEndpointOperation{
				JobID:               to.Ptr(peJobID),
				OperationType:       to.Ptr(ekm.EkmPrivateEndpointOperationTypeDelete),
				PrivateEndpointName: to.Ptr(peName),
				Status:              to.Ptr(ekm.OperationStateSucceeded),
			}}, nil)
			return
		},
		GetEkmPrivateEndpoint: func(ctx context.Context, peName string, options *ekm.GetEkmPrivateEndpointOptions) (resp azfake.Responder[ekm.GetEkmPrivateEndpointResponse], errResp azfake.ErrorResponder) {
			if peName != privateEndpoint {
				errResp.SetResponseError(http.StatusNotFound, "PrivateEndpointNotFound")
				return
			}
			resp.SetResponse(http.StatusOK, ekm.GetEkmPrivateEndpointResponse{PrivateEndpoint: ekm.PrivateEndpoint{
				Name:              to.Ptr(peName),
				ProvisioningState: to.Ptr(ekm.EkmPrivateEndpointProvisioningStateSucceeded),
				PrivateLinkServiceConnectionState: &ekm.PrivateEndpointConnectionState{
					Status: to.Ptr(ekm.EkmPrivateEndpointConnectionStatusApproved),
				},
			}}, nil)
			return
		},
		GetEkmPrivateEndpointOperationStatus: func(ctx context.Context, operationID string, options *ekm.GetEkmPrivateEndpointOperationStatusOptions) (resp azfake.Responder[ekm.GetEkmPrivateEndpointOperationStatusResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, ekm.GetEkmPrivateEndpointOperationStatusResponse{PrivateEndpointOperation: ekm.PrivateEndpointOperation{
				JobID:  to.Ptr(operationID),
				Status: to.Ptr(ekm.OperationStateRunning),
			}}, nil)
			return
		},
		ListEkmPrivateEndpoints: func(ctx context.Context, options *ekm.ListEkmPrivateEndpointsOptions) (resp azfake.Responder[ekm.ListEkmPrivateEndpointsResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, ekm.ListEkmPrivateEndpointsResponse{PrivateEndpointListResult: ekm.PrivateEndpointListResult{
				Value: []*ekm.PrivateEndpoint{{Name: to.Ptr(privateEndpoint)}},
			}}, nil)
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

func TestServerPrivateEndpoints(t *testing.T) {
	fakeServer := getServer()

	client, err := ekm.NewClient("https://fake-hsm.managedhsm.azure.net", &azfake.TokenCredential{}, &ekm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewKeyVaultServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)
	pollOpts := &runtime.PollUntilDoneOptions{Frequency: time.Millisecond}

	// create
	createPoller, err := client.BeginCreateEkmPrivateEndpoint(context.TODO(), privateEndpoint, ekm.PrivateEndpointCreateParameters{
		PrivateLinkServiceID: to.Ptr(plsAlias),
	}, nil)
	require.NoError(t, err)
	createResp, err := createPoller.PollUntilDone(context.TODO(), pollOpts)
	require.NoError(t, err)
	require.Equal(t, peJobID, *createResp.JobID)
	require.Equal(t, privateEndpoint, *createResp.PrivateEndpointName)
	require.Equal(t, ekm.EkmPrivateEndpointOperationTypeCreate, *createResp.OperationType)
	require.Equal(t, ekm.OperationStateSucceeded, *createResp.Status)

	// create failure
	_, err = client.BeginCreateEkmPrivateEndpoint(context.TODO(), "third-pe", ekm.PrivateEndpointCreateParameters{
		PrivateLinkServiceID: to.Ptr(plsAlias),
	}, nil)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusConflict, respErr.StatusCode)

	// the fake receives the request body, so a missing private link service ID is rejected
	_, err = client.BeginCreateEkmPrivateEndpoint(context.TODO(), privateEndpoint, ekm.PrivateEndpointCreateParameters{}, nil)
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)

	// get
	getResp, err := client.GetEkmPrivateEndpoint(context.TODO(), privateEndpoint, nil)
	require.NoError(t, err)
	require.Equal(t, privateEndpoint, *getResp.Name)
	require.Equal(t, ekm.EkmPrivateEndpointProvisioningStateSucceeded, *getResp.ProvisioningState)
	require.Equal(t, ekm.EkmPrivateEndpointConnectionStatusApproved, *getResp.PrivateLinkServiceConnectionState.Status)

	_, err = client.GetEkmPrivateEndpoint(context.TODO(), "missing", nil)
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)

	// list
	listResp, err := client.ListEkmPrivateEndpoints(context.TODO(), nil)
	require.NoError(t, err)
	require.Len(t, listResp.Value, 1)
	require.Equal(t, privateEndpoint, *listResp.Value[0].Name)

	// operation status
	opResp, err := client.GetEkmPrivateEndpointOperationStatus(context.TODO(), peJobID, nil)
	require.NoError(t, err)
	require.Equal(t, peJobID, *opResp.JobID)
	require.Equal(t, ekm.OperationStateRunning, *opResp.Status)

	// delete
	deletePoller, err := client.BeginDeleteEkmPrivateEndpoint(context.TODO(), privateEndpoint, nil)
	require.NoError(t, err)
	deleteResp, err := deletePoller.PollUntilDone(context.TODO(), pollOpts)
	require.NoError(t, err)
	require.Equal(t, ekm.EkmPrivateEndpointOperationTypeDelete, *deleteResp.OperationType)
	require.Equal(t, ekm.OperationStateSucceeded, *deleteResp.Status)
}
