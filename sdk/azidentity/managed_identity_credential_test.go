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
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const (
	appServiceWindowsSuccessResp = `{"access_token": "new_token", "expires_on": "9/14/2017 00:00:00 PM +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	appServiceLinuxSuccessResp   = `{"access_token": "new_token", "expires_on": "09/14/2017 00:00:00 +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnIntResp             = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnNonStringIntResp    = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": 1560974028, "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

// TODO: replace with 1.17's T.Setenv
func clearEnvVars(envVars ...string) {
	for _, ev := range envVars {
		_ = os.Unsetenv(ev)
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

// delayPolicy adds a delay to pipeline requests. Used to test timeout behavior.
type delayPolicy struct {
	delay time.Duration
}

func (p delayPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	time.Sleep(p.delay)
	return req.Next()
}

func TestManagedIdentityCredential_AzureArc(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "arc.key"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	expectedKey := "expected-key"
	n, err := file.WriteString(expectedKey)
	if n != len(expectedKey) || err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}

	expectedPath := "/foo/token"
	validateReq := func(req *http.Request) bool {
		if req.URL.Path != expectedPath {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		if p := req.URL.Query().Get("api-version"); p != azureArcAPIVersion {
			t.Fatalf("unexpected api-version: %s", p)
		}
		if p := req.URL.Query().Get("resource"); p != strings.TrimSuffix(liveTestScope, defaultSuffix) {
			t.Fatalf("unexpected resource: %s", p)
		}
		if h := req.Header.Get("metadata"); h != "true" {
			t.Fatalf("unexpected metadata header: %s", h)
		}
		if h := req.Header.Get("Authorization"); h != "Basic "+expectedKey {
			t.Fatalf("unexpected Authorization: %s", h)
		}
		return true
	}

	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithHeader("WWW-Authenticate", "Basic realm="+file.Name()), mock.WithStatusCode(401))
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse()

	setEnvironmentVariables(t, map[string]string{
		arcIMDSEndpoint:  srv.URL(),
		identityEndpoint: srv.URL() + expectedPath,
	})
	opts := ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}}
	cred, err := NewManagedIdentityCredential(&opts)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token: %s", tk.Token)
	}
	if tk.ExpiresOn.Before(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
}

func TestManagedIdentityCredential_CloudShell(t *testing.T) {
	validateReq := func(req *http.Request) bool {
		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		if v := req.FormValue("resource"); v != strings.TrimSuffix(liveTestScope, defaultSuffix) {
			t.Fatalf("unexpected resource: %s", v)
		}
		if h := req.Header.Get("metadata"); h != "true" {
			t.Fatalf("unexpected metadata header: %s", h)
		}
		return true
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token value: %s", tk.Token)
	}
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse()
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatal("expected an error but didn't receive one")
	}
}

func TestManagedIdentityCredential_CloudShellUserAssigned(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{msiEndpoint: "http://localhost"})
	for _, id := range []ManagedIDKind{ClientID("client-id"), ResourceID("/resource/id")} {
		options := ManagedIdentityCredentialOptions{ID: id}
		msiCred, err := NewManagedIdentityCredential(&options)
		if err != nil {
			t.Fatal(err)
		}
		_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
		var authErr AuthenticationFailedError
		if !errors.As(err, &authErr) {
			t.Fatal("expected AuthenticationFailedError")
		}
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20170901Mock_windows(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	expectedSecret := "expected-secret"
	pred := func(req *http.Request) bool {
		if secret := req.Header.Get("Secret"); secret != expectedSecret {
			t.Fatalf(`unexpected Secret header "%s"`, secret)
		}
		return true
	}
	srv.AppendResponse(mock.WithPredicate(pred), mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: expectedSecret})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceLinuxSuccessResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
	t.Skip("App Service 2019-08-01 isn't supported because it's unavailable in some Functions apps.")
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: "header"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
	t.Skip("App Service 2019-08-01 isn't supported because it's unavailable in some Functions apps.")
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceLinuxSuccessResp)))
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: "header"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: "header"})
	msiCred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
	t.Skip("App Service 2019-08-01 isn't supported because it's unavailable in some Functions apps.")
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	setEnvironmentVariables(t, map[string]string{identityEndpoint: "somevalue", identityHeader: "header"})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), ClientID(fakeClientID), []string{liveTestScope})
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
	if reqQueryParams["resource"][0] != liveTestScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpClientID][0] != fakeClientID {
		t.Fatalf("Unexpected client ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAppServiceAuthRequestV20170901(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	setEnvironmentVariables(t, map[string]string{msiEndpoint: "somevalue", msiSecret: "secret"})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), ClientID(fakeClientID), []string{liveTestScope})
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
	if reqQueryParams["resource"][0] != liveTestScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams["clientid"][0] != fakeClientID {
		t.Fatalf("Unexpected client ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnStringInt(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceMockFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestManagedIdentityCredential_GetTokenIMDS400(t *testing.T) {
	res := http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(bytes.NewBufferString(""))}
	options := ManagedIdentityCredentialOptions{}
	options.Transport = newMockImds(res, res, res)
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	// cred should return credentialUnavailableError when IMDS responds 400 to a token request
	var expected credentialUnavailableError
	for i := 0; i < 3; i++ {
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
		if !errors.As(err, &expected) {
			t.Fatalf(`expected credentialUnavailableError, got %T: "%s"`, err, err.Error())
		}
	}
}

func TestManagedIdentityCredential_NewManagedIdentityCredentialFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: "https://t .com"})
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

func TestManagedIdentityCredential_GetTokenUnexpectedJSON(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}

func TestManagedIdentityCredential_CreateIMDSAuthRequest(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test IMDS authentication request creation.
	setEnvironmentVariables(t, map[string]string{msiEndpoint: "somevalue", msiSecret: "secret"})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createIMDSAuthRequest(context.Background(), ClientID(fakeClientID), []string{liveTestScope})
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
	if reqQueryParams["resource"][0] != liveTestScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams["client_id"][0] != fakeClientID {
		t.Fatalf("Unexpected client ID. Expected: %s, Received: %s", fakeClientID, reqQueryParams["client_id"][0])
	}
	if u := req.Raw().URL.String(); !strings.HasPrefix(u, imdsEndpoint) {
		t.Fatalf("Unexpected default authority host %s", u)
	}
	if req.Raw().URL.Scheme != "http" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestManagedIdentityCredential_GetTokenScopes(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	for _, scopes := range [][]string{nil, {}, {"a", "b"}} {
		_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: scopes})
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "scope") {
			t.Fatalf(`unexpected error "%s"`, err.Error())
		}
	}
}

func TestManagedIdentityCredential_ScopesImmutable(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
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

func TestManagedIdentityCredential_UseResourceID(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceWindowsSuccessResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	options.ID = ResourceID("sample/resource/id")
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != "new_token" {
		t.Fatalf("unexpected token returned. Expected: %s, Received: %s", "new_token", tk.Token)
	}
}

func TestManagedIdentityCredential_ResourceID_AppServiceV20190801(t *testing.T) {
	t.Skip("App Service 2019-08-01 isn't supported because it's unavailable in some Functions apps.")
	setEnvironmentVariables(t, map[string]string{identityEndpoint: "somevalue", identityHeader: "header"})
	resID := "sample/resource/id"
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(resID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), cred.id, []string{liveTestScope})
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
	if reqQueryParams["resource"][0] != liveTestScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpResID][0] != resID {
		t.Fatalf("Unexpected resource ID in resource query param")
	}
}

func TestManagedIdentityCredential_ResourceID_IMDS(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to avoid failure in the constructor
	setEnvironmentVariables(t, map[string]string{msiEndpoint: "http://localhost"})
	resID := "sample/resource/id"
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(resID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.msiType = msiTypeIMDS
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), cred.id, []string{liveTestScope})
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
	if reqQueryParams["resource"][0] != liveTestScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if reqQueryParams[qpResID][0] != resID {
		t.Fatalf("Unexpected resource ID in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnInt(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnNonStringIntResp)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

// adding an incorrect string value in expires_on
func TestManagedIdentityCredential_CreateAccessTokenExpiresOnFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "15609740s28", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL(), msiSecret: "secret"})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("expected to receive an error but received none")
	}
}

func TestManagedIdentityCredential_IMDSLive(t *testing.T) {
	switch recording.GetRecordMode() {
	case recording.LiveMode:
		t.Skip("this test doesn't run in live mode because it can't pass in CI")
	case recording.RecordingMode:
		// record iff either managed identity environment variable is set, because
		// otherwise there's no reason to believe the test is running on a VM
		if len(liveManagedIdentity.clientID)+len(liveManagedIdentity.resourceID) == 0 {
			t.Skip("neither MANAGED_IDENTITY_CLIENT_ID nor MANAGED_IDENTITY_RESOURCE_ID is set")
		}
	}
	opts, stop := initRecording(t)
	defer stop()
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: opts})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token == "" {
		t.Fatal("GetToken returned an invalid token")
	}
	if tk.ExpiresOn.Before(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
}

func TestManagedIdentityCredential_IMDSClientIDLive(t *testing.T) {
	switch recording.GetRecordMode() {
	case recording.LiveMode:
		t.Skip("this test doesn't run in live mode because it can't pass in CI")
	case recording.RecordingMode:
		if liveManagedIdentity.clientID == "" {
			t.Skip("MANAGED_IDENTITY_CLIENT_ID isn't set")
		}
	}
	opts, stop := initRecording(t)
	defer stop()
	o := ManagedIdentityCredentialOptions{ClientOptions: opts, ID: ClientID(liveManagedIdentity.clientID)}
	cred, err := NewManagedIdentityCredential(&o)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token == "" {
		t.Fatal("GetToken returned an invalid token")
	}
	if tk.ExpiresOn.Before(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
}

func TestManagedIdentityCredential_IMDSResourceIDLive(t *testing.T) {
	switch recording.GetRecordMode() {
	case recording.LiveMode:
		t.Skip("this test doesn't run in live mode because it can't pass in CI")
	case recording.RecordingMode:
		if liveManagedIdentity.resourceID == "" {
			t.Skip("MANAGED_IDENTITY_RESOURCE_ID isn't set")
		}
	}
	opts, stop := initRecording(t)
	defer stop()
	o := ManagedIdentityCredentialOptions{ClientOptions: opts, ID: ResourceID(liveManagedIdentity.resourceID)}
	cred, err := NewManagedIdentityCredential(&o)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token == "" {
		t.Fatal("GetToken returned an invalid token")
	}
	if tk.ExpiresOn.Before(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
}

func TestManagedIdentityCredential_IMDSTimeoutExceeded(t *testing.T) {
	resetEnvironmentVarsForTest()
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{delayPolicy{delay: time.Microsecond}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	cred.client.imdsTimeout = time.Nanosecond
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	var expected credentialUnavailableError
	if !errors.As(err, &expected) {
		t.Fatalf(`expected credentialUnavailableError, got %T: "%v"`, err, err)
	}
	if tk != nil {
		t.Fatal("GetToken returned a token and an error")
	}
}

func TestManagedIdentityCredential_IMDSTimeoutSuccess(t *testing.T) {
	resetEnvironmentVarsForTest()
	res := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(accessTokenRespSuccess))}
	options := ManagedIdentityCredentialOptions{}
	options.Transport = newMockImds(res, res)
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	cred.client.imdsTimeout = time.Minute
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf(`got unexpected token "%s"`, tk.Token)
	}
	if !tk.ExpiresOn.After(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
	if cred.client.imdsTimeout > 0 {
		t.Fatal("credential didn't remove IMDS timeout after receiving a response")
	}
}

func TestManagedIdentityCredential_ServiceFabric(t *testing.T) {
	resetEnvironmentVarsForTest()
	expectedSecret := "expected-secret"
	pred := func(req *http.Request) bool {
		if secret := req.Header.Get("Secret"); secret != expectedSecret {
			t.Fatalf(`unexpected Secret header "%s"`, secret)
		}
		if p := req.URL.Query().Get("api-version"); p != serviceFabricAPIVersion {
			t.Fatalf("unexpected api-version: %s", p)
		}
		if p := req.URL.Query().Get("resource"); p != strings.TrimSuffix(liveTestScope, defaultSuffix) {
			t.Fatalf("unexpected resource: %s", p)
		}
		return true
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithPredicate(pred), mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: expectedSecret, identityServerThumbprint: "..."})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf(`got unexpected token "%s"`, tk.Token)
	}
	if !tk.ExpiresOn.After(time.Now().UTC()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
}
