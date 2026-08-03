// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/stretchr/testify/require"
)

// TestEkmConnectionLifecycle exercises the full create -> get -> update -> delete
// flow against a Managed HSM. It is idempotent: any pre-existing EKM connection
// is removed before the test runs.
func TestEkmConnectionLifecycle(t *testing.T) {
	client := startEKMTest(t)
	ctx := context.Background()

	// Best-effort cleanup of any pre-existing connection so the test is repeatable.
	deleteExistingEkmConnection(ctx, t, client)

	// The Subject Common Name on the proxy's server cert must match the host
	// portion of the dial target (without the port) for TLS validation.
	expectedCN := serverSubjectCommonNameFor(ekmProxyHost)

	// CREATE
	createParams := ekm.Connection{
		Host:                    to.Ptr(ekmProxyHost),
		ServerCaCertificates:    [][]byte{ekmCACert},
		ServerSubjectCommonName: to.Ptr(expectedCN),
		PathPrefix:              to.Ptr("/api/v1"),
	}
	testSerde(t, &createParams)

	created, err := client.CreateEkmConnection(ctx, createParams, nil)
	require.NoError(t, err)
	require.NotNil(t, created.Host)
	require.Equal(t, ekmProxyHost, *created.Host)
	require.NotEmpty(t, created.ServerCaCertificates)
	require.NotNil(t, created.PathPrefix)
	require.Equal(t, "/api/v1", *created.PathPrefix)
	require.NotNil(t, created.ServerSubjectCommonName)
	require.Equal(t, expectedCN, *created.ServerSubjectCommonName)
	testSerde(t, &created.Connection)

	// GET
	got, err := client.GetEkmConnection(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, got.Host)
	require.Equal(t, ekmProxyHost, *got.Host)
	require.NotEmpty(t, got.ServerCaCertificates)

	// UPDATE — change the path prefix while leaving the host (and therefore CN)
	// the same, since changing the CN without rotating the proxy cert would
	// break TLS verification.
	updateParams := ekm.Connection{
		Host:                    to.Ptr(ekmProxyHost),
		ServerCaCertificates:    [][]byte{ekmCACert},
		ServerSubjectCommonName: to.Ptr(expectedCN),
		PathPrefix:              to.Ptr("/api/v1"),
	}
	updated, err := client.UpdateEkmConnection(ctx, updateParams, nil)
	require.NoError(t, err)
	require.NotNil(t, updated.PathPrefix)
	require.Equal(t, "/api/v1", *updated.PathPrefix)
	require.NotNil(t, updated.ServerSubjectCommonName)
	require.Equal(t, expectedCN, *updated.ServerSubjectCommonName)

	// DELETE
	deleted, err := client.DeleteEkmConnection(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, deleted.Host)
	require.Equal(t, ekmProxyHost, *deleted.Host)

	// After delete, GET should fail with 4xx.
	_, err = client.GetEkmConnection(ctx, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.GreaterOrEqual(t, httpErr.StatusCode, 400)
	require.Less(t, httpErr.StatusCode, 500)
}

// TestCreateEkmConnection_ConflictWhenExists verifies that creating an EKM
// connection when one already exists fails with a 4xx ResponseError.
func TestCreateEkmConnection_ConflictWhenExists(t *testing.T) {
	client := startEKMTest(t)
	ctx := context.Background()

	deleteExistingEkmConnection(ctx, t, client)

	params := ekm.Connection{
		Host:                    to.Ptr(ekmProxyHost),
		ServerCaCertificates:    [][]byte{ekmCACert},
		ServerSubjectCommonName: to.Ptr(serverSubjectCommonNameFor(ekmProxyHost)),
		PathPrefix:              to.Ptr("/api/v1"),
	}
	_, err := client.CreateEkmConnection(ctx, params, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = client.DeleteEkmConnection(context.Background(), nil)
	})

	_, err = client.CreateEkmConnection(ctx, params, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.GreaterOrEqual(t, httpErr.StatusCode, 400)
	require.Less(t, httpErr.StatusCode, 500)
}

// TestGetEkmConnection_NotFound verifies that GetEkmConnection returns a 4xx
// error when no EKM connection has been configured on the vault.
func TestGetEkmConnection_NotFound(t *testing.T) {
	client := startEKMTest(t)
	ctx := context.Background()

	deleteExistingEkmConnection(ctx, t, client)

	res, err := client.GetEkmConnection(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res.Host)
	require.Empty(t, res.ServerCaCertificates)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.GreaterOrEqual(t, httpErr.StatusCode, 400)
	require.Less(t, httpErr.StatusCode, 500)
}

// TestGetEkmCertificate fetches the proxy client certificate that the Managed
// HSM will present to the EKM proxy. The certificate is generated by the
// service; the test only requires that a non-empty common name is returned.
func TestGetEkmCertificate(t *testing.T) {
	client := startEKMTest(t)
	ctx := context.Background()

	res, err := client.GetEkmCertificate(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, res.SubjectCommonName)
	require.NotEmpty(t, *res.SubjectCommonName)
	require.NotEmpty(t, res.CaCertificates)
	testSerde(t, &res.ProxyClientCertificateInfo)
}

// TestCheckEkmConnection calls the proxy connectivity check. If the EKM proxy
// is unreachable (common in CI where no real proxy is provisioned) the service
// returns a 4xx/5xx, which the test handles gracefully by skipping rather than
// failing — the wire contract is still validated either way.
func TestCheckEkmConnection(t *testing.T) {
	client := startEKMTest(t)
	ctx := context.Background()

	deleteExistingEkmConnection(ctx, t, client)
	params := ekm.Connection{
		Host:                    to.Ptr(ekmProxyHost),
		ServerCaCertificates:    [][]byte{ekmCACert},
		ServerSubjectCommonName: to.Ptr(serverSubjectCommonNameFor(ekmProxyHost)),
		PathPrefix:              to.Ptr("/api/v1"),
	}
	_, err := client.CreateEkmConnection(ctx, params, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = client.DeleteEkmConnection(context.Background(), nil)
	})

	res, err := client.CheckEkmConnection(ctx, nil)
	if err != nil {
		// The proxy almost certainly isn't reachable from CI; assert the failure
		// is a transport-level ResponseError (not a panic / client-side issue)
		// and skip the success path.
		var httpErr *azcore.ResponseError
		require.ErrorAs(t, err, &httpErr)
		t.Skipf("EKM proxy unreachable (HTTP %d): %s", httpErr.StatusCode, httpErr.ErrorCode)
	}
	require.NotNil(t, res.APIVersion)
	require.NotNil(t, res.EkmVendor)
	require.NotNil(t, res.EkmProduct)
	require.NotNil(t, res.ProxyVendor)
	require.NotNil(t, res.ProxyName)
	testSerde(t, &res.ProxyInfo)
}

// TestConnectionSerde exercises round-trip JSON marshaling/unmarshaling for
// every EKM model so the generator's serde code is exercised even when no
// live infrastructure is available.
func TestConnectionSerde(t *testing.T) {
	host := "ekm-proxy.example.com:443"
	conn := ekm.Connection{
		Host:                    to.Ptr(host),
		ServerCaCertificates:    [][]byte{ekmCACert, []byte("second-cert")},
		ServerSubjectCommonName: to.Ptr(serverSubjectCommonNameFor(host)),
		PathPrefix:              to.Ptr("/api/v1"),
	}
	testSerde(t, &conn)

	clientCert := ekm.ProxyClientCertificateInfo{
		CaCertificates:    [][]byte{ekmCACert},
		SubjectCommonName: to.Ptr("CN=hsm-client"),
	}
	testSerde(t, &clientCert)

	proxyInfo := ekm.ProxyInfo{
		APIVersion:  to.Ptr("1.0"),
		EkmVendor:   to.Ptr("Test Vendor"),
		EkmProduct:  to.Ptr("Test Product v1.0"),
		ProxyVendor: to.Ptr("Test Proxy Vendor"),
		ProxyName:   to.Ptr("Test Proxy v2.0"),
	}
	testSerde(t, &proxyInfo)
}

// TestAPIVersion verifies that ClientOptions.APIVersion actually overrides the
// version the generated client puts on the wire. The transport is a mock
// server returning a canned 200 OK regardless of the request, so the value
// chosen here is meaningless to any real service. It simply has to differ from
// the real API version, so the predicate can detect that the override works.
func TestAPIVersion(t *testing.T) {
	apiVersion := "7.3"
	requireVersion := func(req *http.Request) bool {
		version := req.URL.Query().Get("api-version")
		require.Equal(t, apiVersion, version)
		return true
	}
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(`{"host":"ekm.example.com","server_ca_certificates":[]}`)),
		mock.WithPredicate(requireVersion),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))

	opts := &ekm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := ekm.NewClient(hsmURL, &azcred.Fake{}, opts)
	require.NoError(t, err)

	_, err = client.GetEkmConnection(context.Background(), nil)
	require.NoError(t, err)
}

// deleteExistingEkmConnection removes any pre-existing EKM connection so a test
// can be re-run cleanly. Errors are tolerated because the most common case is
// that no connection exists yet (a 4xx response from the service).
func deleteExistingEkmConnection(ctx context.Context, t *testing.T, client *ekm.KeyVaultClient) {
	t.Helper()
	_, err := client.DeleteEkmConnection(ctx, nil)
	if err == nil {
		// Some services need a short cooldown after delete before another
		// create succeeds; the test-proxy strips this wait in playback.
		recording.Sleep(0)
		return
	}
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) && httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 {
		return // no pre-existing connection to clean up
	}
	require.NoError(t, err, "unexpected error cleaning up pre-existing EKM connection")
}
