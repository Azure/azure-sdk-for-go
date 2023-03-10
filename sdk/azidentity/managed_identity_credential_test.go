//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
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
	expiresOnIntResp          = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnNonStringIntResp = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": 1560974028, "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

// TODO: replace with 1.17's T.Setenv
func clearEnvVars(envVars ...string) {
	for _, ev := range envVars {
		_ = os.Unsetenv(ev)
	}
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
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
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
	testGetTokenSuccess(t, cred)
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
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, msiCred)
}

func TestManagedIdentityCredential_AppService(t *testing.T) {
	expectedID := "expected-ID"
	expectedHeader := "header"
	for _, id := range []ManagedIDKind{ClientID(expectedID), ResourceID(expectedID), nil} {
		validateReq := func(req *http.Request) bool {
			if h := req.Header.Get("X-IDENTITY-HEADER"); h != expectedHeader {
				t.Fatalf("unexpected X-IDENTITY-HEADER: %s", h)
			}
			q := req.URL.Query()
			if v := q.Get("api-version"); v != "2019-08-01" {
				t.Fatalf(`unexpected api-version "%s"`, v)
			}
			if v := q.Get("resource"); v != strings.TrimSuffix(liveTestScope, "/.default") {
				t.Fatalf(`unexpected resource "%s"`, v)
			}
			if id == nil {
				if q.Get(qpClientID) != "" || q.Get(qpResID) != "" {
					t.Fatal("request shouldn't include a user-assigned ID")
				}
			} else {
				if q.Get(qpClientID) != "" && q.Get(qpResID) != "" {
					t.Fatal("request includes two IDs")
				}
				var v string
				if _, ok := id.(ClientID); ok {
					v = q.Get(qpClientID)
				} else if _, ok := id.(ResourceID); ok {
					v = q.Get(qpResID)
				}
				if v != id.String() {
					t.Fatalf(`unexpected id "%s"`, v)
				}
			}
			return true
		}

		t.Run(fmt.Sprintf("%T", id), func(t *testing.T) {
			srv, close := mock.NewServer()
			defer close()
			srv.AppendResponse(
				mock.WithPredicate(validateReq),
				mock.WithBody([]byte(fmt.Sprintf(
					`{"access_token": "%s", "expires_on": "%d", "resource": "https://vault.azure.net", "token_type": "Bearer", "client_id": "some-guid"}`,
					tokenValue,
					time.Now().Add(time.Hour).Unix(),
				))),
			)
			srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
			setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: expectedHeader})
			options := ManagedIdentityCredentialOptions{ID: id}
			options.Transport = srv
			cred, err := NewManagedIdentityCredential(&options)
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestManagedIdentityCredential_AppServiceError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: "secret"})
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
	if !strings.HasPrefix(err.Error(), credNameManagedIdentity) {
		t.Fatal("missing credential type prefix")
	}
}

func TestManagedIdentityCredential_GetTokenIMDS400(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusBadRequest), mock.WithBody([]byte("something went wrong")))
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	cred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	// cred should return credentialUnavailableError when IMDS responds 400 to a token request
	for i := 0; i < 3; i++ {
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
		if _, ok := err.(*credentialUnavailableError); !ok {
			t.Fatalf("expected credentialUnavailableError, received %T", err)
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
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req, err := cred.mic.createIMDSAuthRequest(context.Background(), ClientID(fakeClientID), []string{liveTestScope})
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

func TestManagedIdentityCredential_ResourceID_IMDS(t *testing.T) {
	resID := "sample/resource/id"
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(resID)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req, err := cred.mic.createAuthRequest(context.Background(), cred.mic.id, []string{liveTestScope})
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
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
}

// adding an incorrect string value in expires_on
func TestManagedIdentityCredential_CreateAccessTokenExpiresOnFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "15609740s28", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`)))
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
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
	testGetTokenSuccess(t, cred)
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
	testGetTokenSuccess(t, cred)
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
	testGetTokenSuccess(t, cred)
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
	srv.AppendResponse(mock.WithPredicate(pred), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{identityEndpoint: srv.URL(), identityHeader: expectedSecret, identityServerThumbprint: "..."})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}
