// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	msiScope                     = "https://storage.azure.com"
	appServiceWindowsSuccessResp = `{"access_token": "new_token", "expires_on": "9/14/2017 00:00:00 PM +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	appServiceLinuxSuccessResp   = `{"access_token": "new_token", "expires_on": "09/14/2017 00:00:00 +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnIntResp             = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnNonStringIntResp    = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": 1560974028, "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

func clearEnvVars(envVars ...string) {
	for _, ev := range envVars {
		_ = os.Setenv(ev, "")
	}
}

// A simple fake IMDS. Similar to mock.Server but doesn't wrap httptest.Server. That's
// important because IMDS is at 169.254.169.254, not httptest.Server's default 127.0.0.1.
type mockIMDS struct {
	resp []http.Response
}

func newMockImds(responses ...http.Response) (m *mockIMDS) {
	return &mockIMDS{resp: responses}
}

func (m *mockIMDS) Do(req *http.Request) (*http.Response, error) {
	if len(m.resp) > 0 {
		resp := m.resp[0]
		m.resp = m.resp[1:]
		return &resp, nil
	}
	panic("no more responses")
}

func TestManagedIdentityCredential_GetTokenInAzureArcLive(t *testing.T) {
	if len(os.Getenv(arcIMDSEndpoint)) == 0 {
		t.Skip()
	}
	msiCred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellLive(t *testing.T) {
	if len(os.Getenv("MSI_ENDPOINT")) == 0 {
		t.Skip()
	}
	msiCred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMock(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMockFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20170901Mock_windows(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20170901 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20170901, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20170901Mock_linux(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceLinuxSuccessResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20170901 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20170901, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20190801Mock_windows(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	_ = os.Setenv("IDENTITY_ENDPOINT", srv.URL())
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20190801 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20190801, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20190801Mock_linux(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceLinuxSuccessResp)))
	_ = os.Setenv("IDENTITY_ENDPOINT", srv.URL())
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20190801 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20190801, msiCred.client.msiType)
	}
}

// Azure Functions on linux environments currently doesn't properly support the identity header,
// therefore, preference must be given to the legacy MSI_ENDPOINT variable.
func TestManagedIdentityCredential_GetTokenInAzureFunctions_linux(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	_ = os.Setenv("IDENTITY_ENDPOINT", srv.URL())
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	msiCred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20170901 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20170901, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_CreateAppServiceAuthRequestV20190801(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	_ = os.Setenv("IDENTITY_ENDPOINT", "somevalue")
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), ClientID(clientID), []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Raw().Header.Get("X-IDENTITY-HEADER") != "header" {
		t.Fatalf("Unexpected value for secret header")
	}
	reqQueryParams, err := url.ParseQuery(req.Raw().URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2019-08-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAppServiceAuthRequestV20170901(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	_ = os.Setenv("MSI_ENDPOINT", "somevalue")
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), ClientID(clientID), []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Raw().Header.Get("secret") != "secret" {
		t.Fatalf("Unexpected value for secret header")
	}
	reqQueryParams, err := url.ParseQuery(req.Raw().URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2017-09-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams["clientid"][0] != clientID {
		t.Fatalf("Unexpected client ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnStringInt(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceMockFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestManagedIdentityCredential_GetTokenIMDS400(t *testing.T) {
	resetEnvironmentVarsForTest()
	options := ManagedIdentityCredentialOptions{}
	res1 := http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}
	res2 := res1
	options.Transport = newMockImds(res1, res2)
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// cred should return CredentialUnavailableError when IMDS responds 400 to a token request.
	// Also, it shouldn't send another token request (mockIMDS will appropriately panic if it does).
	var expected CredentialUnavailableError
	for i := 0; i < 3; i++ {
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
		if !errors.As(err, &expected) {
			t.Fatalf("Expected %T, got %T", expected, err)
		}
	}
}

func TestManagedIdentityCredential_NewManagedIdentityCredentialFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", "https://t .com")
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestBearerPolicy_ManagedIdentityCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, msiScope)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}

func TestManagedIdentityCredential_GetTokenUnexpectedJSON(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}

func TestManagedIdentityCredential_CreateIMDSAuthRequest(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test IMDS authentication request creation.
	_ = os.Setenv("MSI_ENDPOINT", "somevalue")
	defer clearEnvVars("MSI_ENDPOINT")
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createIMDSAuthRequest(context.Background(), ClientID(clientID), []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Raw().Header.Get(headerMetadata) != "true" {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	reqQueryParams, err := url.ParseQuery(req.Raw().URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse IMDS query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != imdsAPIVersion {
		t.Fatalf("Unexpected IMDS API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams["client_id"][0] != clientID {
		t.Fatalf("Unexpected client ID. Expected: %s, Received: %s", clientID, reqQueryParams["client_id"][0])
	}
	if u := req.Raw().URL.String(); !strings.HasPrefix(u, imdsEndpoint) {
		t.Fatalf("Unexpected default authority host %s", u)
	}
	if req.Raw().URL.Scheme != "http" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestManagedIdentityCredential_GetTokenEnvVar(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_CLIENT_ID", "test_client_id")
	if err != nil {
		t.Fatal(err)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	at, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if at.Token != "new_token" {
		t.Fatalf("Did not receive the correct access token")
	}
}

func TestManagedIdentityCredential_GetTokenNilResource(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: nil})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	if err.Error() != "must specify a resource in order to authenticate" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagedIdentityCredential_ScopesImmutable(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	_ = os.Setenv(msiEndpoint, srv.URL())
	defer clearEnvVars(msiEndpoint)
	options := ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}}
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	scope := "https://localhost/.default"
	scopes := []string{scope}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scopes[0] != scope {
		t.Fatalf("GetToken shouldn't mutate arguments")
	}
}

func TestManagedIdentityCredential_GetTokenMultipleResources(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"resource1", "resource2"}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	if err.Error() != "can only specify one resource to authenticate with ManagedIdentityCredential" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagedIdentityCredential_UseResourceID(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	options.ID = ResourceID("sample/resource/id")
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != "new_token" {
		t.Fatalf("unexpected token returned. Expected: %s, Received: %s", "new_token", tk.Token)
	}
}

func TestManagedIdentityCredential_ResourceID_AppService(t *testing.T) {
	// setting a dummy value for IDENTITY_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	_ = os.Setenv("IDENTITY_ENDPOINT", "somevalue")
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	resID := "sample/resource/id"
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(resID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), cred.id, []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Raw().Header.Get("X-IDENTITY-HEADER") != "header" {
		t.Fatalf("Unexpected value for secret header")
	}
	reqQueryParams, err := url.ParseQuery(req.Raw().URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2019-08-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpResID][0] != resID {
		t.Fatalf("Unexpected resource ID in resource query param")
	}
}

func TestManagedIdentityCredential_ResourceID_IMDS(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to avoid failure in the constructor
	_ = os.Setenv("MSI_ENDPOINT", "http://foo.com/")
	defer clearEnvVars("MSI_ENDPOINT")
	resID := "sample/resource/id"
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(resID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.msiType = msiTypeIMDS
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), cred.id, []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	reqQueryParams, err := url.ParseQuery(req.Raw().URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2018-02-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpResID][0] != resID {
		t.Fatalf("Unexpected resource ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnInt(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnNonStringIntResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

// adding an incorrect string value in expires_on
func TestManagedIdentityCredential_CreateAccessTokenExpiresOnFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "15609740s28", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("expected to receive an error but received none")
	}
}
