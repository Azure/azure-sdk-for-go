// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/stretchr/testify/require"
)

// These tests exercise every KeyVaultClient operation against a mock HTTP
// transport so the generated request/response code is covered even when the
// live EKM tests are skipped (the default in CI, where no EKM proxy is
// provisioned). They intentionally don't depend on EKM_PROXY_HOST or
// recordings — the goal is unit coverage of the client, not service behavior.

const mockEKMConnectionBody = `{
	"host":"ekm-proxy.example.com:443",
	"path_prefix":"/v1",
	"server_subject_common_name":"ekm-proxy.example.com",
	"server_ca_certificates":["YXp1cmUtc2RrLWZvci1nby1la20tdGVzdC1jYS1jZXJ0aWZpY2F0ZS1ieXRlcw=="]
}`

const mockEKMCertificateBody = `{
	"subject_common_name":"CN=hsm-client",
	"ca_certificates":["YXp1cmUtc2RrLWZvci1nby1la20tdGVzdC1jYS1jZXJ0aWZpY2F0ZS1ieXRlcw=="]
}`

const mockEKMProxyInfoBody = `{
	"api_version":"1.0",
	"ekm_vendor":"Contoso",
	"ekm_product":"Contoso EKM 1.0",
	"proxy_vendor":"Contoso",
	"proxy_name":"Contoso Proxy 2.0"
}`

// newMockClient builds a KeyVaultClient wired to the given mock transport.
// Retries are disabled to keep tests fast and to avoid consuming queued mock
// responses on retriable status codes (e.g. 502).
// Note: KeyVaultChallengePolicy may first send an unauthenticated request (with
// any body removed) and only send an authorized request after a 401 challenge.
	t.Helper()
	opts := &ekm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	}
	client, err := ekm.NewClient(hsmURL, &azcred.Fake{}, opts)
	require.NoError(t, err)
	return client
}

// sampleConnection returns a Connection value that round-trips cleanly through
// the generated serde so it can be used as a request body in mock tests.
func sampleConnection() ekm.Connection {
	return ekm.Connection{
		Host:                    to.Ptr("ekm-proxy.example.com:443"),
		ServerCaCertificates:    [][]byte{[]byte("ca-cert-bytes")},
		ServerSubjectCommonName: to.Ptr("ekm-proxy.example.com"),
		PathPrefix:              to.Ptr("/v1"),
	}
}

func TestCreateEkmConnection_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMConnectionBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.CreateEkmConnection(context.Background(), sampleConnection(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.Host)
	require.Equal(t, "ekm-proxy.example.com:443", *res.Host)
	require.NotNil(t, res.PathPrefix)
	require.Equal(t, "/v1", *res.PathPrefix)
	require.Len(t, res.ServerCaCertificates, 1)
}

func TestCreateEkmConnection_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusConflict),
		mock.WithBody([]byte(`{"error":{"code":"Conflict","message":"already exists"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.CreateEkmConnection(context.Background(), sampleConnection(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusConflict, respErr.StatusCode)
}

func TestGetEkmConnection_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMConnectionBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.GetEkmConnection(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.Host)
	require.Equal(t, "ekm-proxy.example.com:443", *res.Host)
	require.NotNil(t, res.ServerSubjectCommonName)
	require.Equal(t, "ekm-proxy.example.com", *res.ServerSubjectCommonName)
}

func TestGetEkmConnection_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithBody([]byte(`{"error":{"code":"NotFound","message":"no ekm connection"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.GetEkmConnection(context.Background(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
}

func TestGetEkmConnection_MockInvalidBody(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte("not-json")),
	)

	client := newMockClient(t, srv)
	_, err := client.GetEkmConnection(context.Background(), nil)
	require.Error(t, err)
}

func TestUpdateEkmConnection_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMConnectionBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.UpdateEkmConnection(context.Background(), sampleConnection(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.PathPrefix)
	require.Equal(t, "/v1", *res.PathPrefix)
}

func TestUpdateEkmConnection_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusBadRequest),
		mock.WithBody([]byte(`{"error":{"code":"BadRequest","message":"bad request"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.UpdateEkmConnection(context.Background(), sampleConnection(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
}

func TestDeleteEkmConnection_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMConnectionBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.DeleteEkmConnection(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.Host)
	require.Equal(t, "ekm-proxy.example.com:443", *res.Host)
}

func TestDeleteEkmConnection_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithBody([]byte(`{"error":{"code":"NotFound","message":"no ekm connection"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.DeleteEkmConnection(context.Background(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
}

func TestGetEkmCertificate_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMCertificateBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.GetEkmCertificate(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.SubjectCommonName)
	require.Equal(t, "CN=hsm-client", *res.SubjectCommonName)
	require.Len(t, res.CaCertificates, 1)
}

func TestGetEkmCertificate_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusForbidden),
		mock.WithBody([]byte(`{"error":{"code":"Forbidden","message":"no permission"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.GetEkmCertificate(context.Background(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusForbidden, respErr.StatusCode)
}

func TestCheckEkmConnection_Mock(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(mockEKMProxyInfoBody)),
	)

	client := newMockClient(t, srv)
	res, err := client.CheckEkmConnection(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res.APIVersion)
	require.Equal(t, "1.0", *res.APIVersion)
	require.NotNil(t, res.EkmVendor)
	require.Equal(t, "Contoso", *res.EkmVendor)
	require.NotNil(t, res.EkmProduct)
	require.NotNil(t, res.ProxyVendor)
	require.NotNil(t, res.ProxyName)
}

func TestCheckEkmConnection_MockError(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusBadGateway),
		mock.WithBody([]byte(`{"error":{"code":"BadGateway","message":"proxy unreachable"}}`)),
	)

	client := newMockClient(t, srv)
	_, err := client.CheckEkmConnection(context.Background(), nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusBadGateway, respErr.StatusCode)
}

// TestEkmRequestRouting asserts that each operation makes the right HTTP
// method/path/api-version combination.
// It captures method/path/query via mock.WithPredicate. This is sufficient because
// KeyVaultChallengePolicy may first send an unauthenticated request (no body) before
// resending the same method/path/query after a 401 challenge.
func TestEkmRequestRouting(t *testing.T) {
	type call struct {
		method string
		path   string
	}
	for _, tc := range []struct {
		name   string
		invoke func(t *testing.T, c *ekm.KeyVaultClient)
		want   call
		body   string
	}{
		{
			name: "CheckEkmConnection",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.CheckEkmConnection(context.Background(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodPost, "/ekm/check"},
			body: mockEKMProxyInfoBody,
		},
		{
			name: "CreateEkmConnection",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.CreateEkmConnection(context.Background(), sampleConnection(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodPost, "/ekm/create"},
			body: mockEKMConnectionBody,
		},
		{
			name: "DeleteEkmConnection",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.DeleteEkmConnection(context.Background(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodDelete, "/ekm"},
			body: mockEKMConnectionBody,
		},
		{
			name: "GetEkmCertificate",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.GetEkmCertificate(context.Background(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodGet, "/ekm/certificate"},
			body: mockEKMCertificateBody,
		},
		{
			name: "GetEkmConnection",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.GetEkmConnection(context.Background(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodGet, "/ekm"},
			body: mockEKMConnectionBody,
		},
		{
			name: "UpdateEkmConnection",
			invoke: func(t *testing.T, c *ekm.KeyVaultClient) {
				_, err := c.UpdateEkmConnection(context.Background(), sampleConnection(), nil)
				require.NoError(t, err)
			},
			want: call{http.MethodPatch, "/ekm"},
			body: mockEKMConnectionBody,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeSrv()
			var got call
			var apiVersion string
			srv.AppendResponse(
				mock.WithStatusCode(http.StatusOK),
				mock.WithBody([]byte(tc.body)),
				mock.WithPredicate(func(req *http.Request) bool {
					got.method = req.Method
					got.path = req.URL.Path
					apiVersion = req.URL.Query().Get("api-version")
					return true
				}),
			)
			srv.AppendResponse(
				mock.WithStatusCode(http.StatusOK),
				mock.WithBody([]byte(tc.body)),
			)

			client := newMockClient(t, srv)
			tc.invoke(t, client)
			require.Equal(t, tc.want.method, got.method)
			require.Equal(t, tc.want.path, got.path)
			require.NotEmpty(t, apiVersion, "api-version query parameter should be set")
		})
	}
}
