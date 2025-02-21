//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

const (
	expiresOnIntResp          = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnNonStringIntResp = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": 1560974028, "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

func TestManagedIdentityCredential_AzureArc(t *testing.T) {
	d := t.TempDir()
	before := arcKeyDirectory
	arcKeyDirectory = func() (string, error) { return d, nil }
	defer func() { arcKeyDirectory = before }()
	file, err := os.Create(filepath.Join(d, "arc.key"))
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

func TestManagedIdentityCredential_AzureArcErrors(t *testing.T) {
	for k, v := range map[string]string{
		arcIMDSEndpoint:  "https://localhost",
		identityEndpoint: "https://localhost",
	} {
		t.Setenv(k, v)
	}

	for _, test := range []struct {
		challenge, name string
		statusCode      int
	}{
		{name: "no challenge", statusCode: http.StatusUnauthorized},
		{name: "malformed challenge", challenge: "Basic realm", statusCode: http.StatusUnauthorized},
		{name: "unexpected status code", statusCode: http.StatusOK},
	} {
		t.Run(test.name, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", test.challenge),
				mock.WithStatusCode(test.statusCode),
			)
			cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
				ClientOptions: azcore.ClientOptions{Transport: srv},
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = cred.GetToken(context.Background(), testTRO)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
	t.Run("failed to get key", func(t *testing.T) {
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.SetError(fmt.Errorf("it didn't work"))
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Retry:     policy.RetryOptions{MaxRetries: -1},
				Transport: srv,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = cred.GetToken(context.Background(), testTRO)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
	t.Run("no key file", func(t *testing.T) {
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.AppendResponse(
			mock.WithHeader("WWW-Authenticate", "Basic realm="+filepath.Join(t.TempDir(), t.Name())),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
		if err != nil {
			t.Fatal(err)
		}
		_, err = cred.GetToken(context.Background(), testTRO)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
	t.Run("key too large", func(t *testing.T) {
		d := t.TempDir()
		f := filepath.Join(d, "test.key")
		err := os.WriteFile(f, bytes.Repeat([]byte("."), 4097), 0600)
		require.NoError(t, err)
		before := arcKeyDirectory
		arcKeyDirectory = func() (string, error) { return d, nil }
		defer func() { arcKeyDirectory = before }()
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.AppendResponse(
			mock.WithHeader("WWW-Authenticate", "Basic realm="+f),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, testTRO)
		require.ErrorContains(t, err, "too large")
	})
	t.Run("unexpected file paths", func(t *testing.T) {
		d, err := arcKeyDirectory()
		if err != nil {
			// test is running on an unsupported OS e.g. darwin
			t.Skip(err)
		}
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.AppendResponse(
			// unexpected directory
			mock.WithHeader("WWW-Authenticate", "Basic realm="+filepath.Join("foo", "bar.key")),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		o := ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}}
		cred, err := NewManagedIdentityCredential(&o)
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, testTRO)
		require.ErrorContains(t, err, "unexpected file path")

		srv.AppendResponse(
			// unexpected extension
			mock.WithHeader("WWW-Authenticate", "Basic realm="+filepath.Join(d, "foo")),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		cred, err = NewManagedIdentityCredential(&o)
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, testTRO)
		require.ErrorContains(t, err, "unexpected file path")
	})
	if runtime.GOOS == "windows" {
		t.Run("ProgramData not set", func(t *testing.T) {
			t.Setenv("ProgramData", "")
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", "Basic realm=foo"),
				mock.WithStatusCode(http.StatusUnauthorized),
			)
			cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
			require.NoError(t, err)
			_, err = cred.GetToken(ctx, testTRO)
			require.ErrorContains(t, err, "ProgramData")
		})
	}
}

func TestManagedIdentityCredential_AzureContainerInstanceLive(t *testing.T) {
	// This test triggers the managed identity test app deployed to an Azure Container Instance.
	// See the bicep file and test resources scripts for details.
	// It triggers the app with az because the test subscription prohibits opening ports to the internet.
	name := os.Getenv("AZIDENTITY_ACI_NAME")
	rg := os.Getenv("AZIDENTITY_RESOURCE_GROUP")
	if name == "" || rg == "" {
		t.Skip("set AZIDENTITY_ACI_NAME and AZIDENTITY_RESOURCE_GROUP to run this test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	command := fmt.Sprintf("az container exec -g %s -n %s --exec-command 'wget -qO- localhost'", rg, name)
	// using "script" as a workaround for "az container exec" requiring a tty
	// https://github.com/Azure/azure-cli/issues/17530
	cmd := exec.CommandContext(ctx, "script", "-q", "-O", "/dev/null", "-c", command)
	b, err := cmd.CombinedOutput()
	s := string(b)
	require.NoError(t, err, s)
	require.Equal(t, "test passed", s)
}

func TestManagedIdentityCredential_AzureFunctionsLive(t *testing.T) {
	// This test triggers the managed identity test app deployed to Azure Functions.
	// See the bicep file and test resources scripts for details.
	fn := os.Getenv("AZIDENTITY_FUNCTION_NAME")
	if fn == "" {
		t.Skip("set AZIDENTITY_FUNCTION_NAME to run this test")
	}
	url := fmt.Sprintf("https://%s.azurewebsites.net/api/HttpTrigger", fn)
	res, err := http.Get(url)
	require.NoError(t, err)
	if res.StatusCode != http.StatusOK {
		b, err := azruntime.Payload(res)
		require.NoError(t, err)
		t.Fatal("test application returned an error: " + string(b))
	}
}

func TestManagedIdentityCredential_AzureMLLive(t *testing.T) {
	switch recording.GetRecordMode() {
	case recording.LiveMode:
		t.Skip("this test doesn't run in live mode because it can't pass in CI")
	case recording.PlaybackMode:
		t.Setenv(defaultIdentityClientID, fakeClientID)
		t.Setenv(msiEndpoint, fakeMIEndpoint)
		t.Setenv(msiSecret, redacted)
	case recording.RecordingMode:
		missing := []string{}
		for _, v := range []string{defaultIdentityClientID, msiEndpoint, msiSecret} {
			if len(os.Getenv(v)) == 0 {
				missing = append(missing, v)
			}
		}
		if len(missing) > 0 {
			t.Skip("no value for " + strings.Join(missing, ", "))
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

func TestManagedIdentityCredential_CloudShell(t *testing.T) {
	validateReq := func(req *http.Request) *http.Response {
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
		return nil
	}
	options := ManagedIdentityCredentialOptions{}
	options.Transport = &mockSTS{tokenRequestCallback: validateReq}
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
				if q.Get(qpClientID) != "" || q.Get(miResID) != "" {
					t.Fatal("request shouldn't include a user-assigned ID")
				}
			} else {
				if q.Get(qpClientID) != "" && q.Get(miResID) != "" {
					t.Fatal("request includes two IDs")
				}
				var v string
				if _, ok := id.(ClientID); ok {
					v = q.Get(qpClientID)
				} else if _, ok := id.(ResourceID); ok {
					v = q.Get(miResID)
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
	_, err = msiCred.GetToken(context.Background(), testTRO)
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
		_, err = cred.GetToken(context.Background(), testTRO)
		if _, ok := err.(credentialUnavailable); !ok {
			t.Fatalf("expected credentialUnavailable, received %T", err)
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
	_, err = msiCred.GetToken(context.Background(), testTRO)
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
	options := ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: &mockSTS{}}}
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
	if reqQueryParams[msiResID][0] != resID {
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
	_, err = msiCred.GetToken(context.Background(), testTRO)
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
	_, err = msiCred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatalf("expected to receive an error but received none")
	}
}

func TestManagedIdentityCredential_IMDSLive(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode && !liveManagedIdentity.imds {
		t.Skip("set IDENTITY_IMDS_AVAILABLE to run this test")
	}

	t.Run("client ID", func(t *testing.T) {
		if recording.GetRecordMode() != recording.PlaybackMode && liveManagedIdentity.clientID == "" {
			t.Skip("set IDENTITY_VM_USER_ASSIGNED_MI_CLIENT_ID to run this test")
		}
		opts, stop := initRecording(t)
		defer stop()
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: opts, ID: ClientID(liveManagedIdentity.clientID)},
		)
		require.NoError(t, err)
		testGetTokenSuccess(t, cred)
	})

	t.Run("object ID", func(t *testing.T) {
		if recording.GetRecordMode() != recording.PlaybackMode && liveManagedIdentity.objectID == "" {
			t.Skip("set IDENTITY_VM_USER_ASSIGNED_MI_OBJECT_ID to run this test")
		}
		opts, stop := initRecording(t)
		defer stop()
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: opts, ID: ObjectID(liveManagedIdentity.objectID)},
		)
		require.NoError(t, err)
		testGetTokenSuccess(t, cred)
	})

	t.Run("resource ID", func(t *testing.T) {
		if recording.GetRecordMode() != recording.PlaybackMode && liveManagedIdentity.resourceID == "" {
			t.Skip("set IDENTITY_VM_USER_ASSIGNED_MI_RESOURCE_ID to run this test")
		}
		opts, stop := initRecording(t)
		defer stop()
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: opts, ID: ResourceID(liveManagedIdentity.resourceID)},
		)
		require.NoError(t, err)
		testGetTokenSuccess(t, cred)
	})

	t.Run("system assigned", func(t *testing.T) {
		opts, stop := initRecording(t)
		defer stop()
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: opts})
		require.NoError(t, err)
		testGetTokenSuccess(t, cred)
	})
}

func TestManagedIdentityCredential_IMDSRetries(t *testing.T) {
	sts := mockSTS{}
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				MaxRetries:    1,
				MaxRetryDelay: time.Nanosecond,
			},
			Transport: &sts,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cred.mic.msiType != msiTypeIMDS {
		t.SkipNow()
	}
	for _, code := range []int{404, 410, 429, 500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511} {
		reqs := 0
		sts.tokenRequestCallback = func(r *http.Request) *http.Response {
			reqs++
			return &http.Response{Body: http.NoBody, Request: r, StatusCode: code}
		}
		_, err = cred.GetToken(context.Background(), testTRO)
		if err == nil {
			t.Fatal("expected an error")
		}
		if reqs != 2 {
			t.Errorf("expected 1 retry after %d response, got %d", code, reqs-1)
		}
	}
}

func TestManagedIdentityCredential_Logs(t *testing.T) {
	logs := []string{}
	log.SetListener(func(e log.Event, msg string) {
		if e == EventAuthentication {
			logs = append(logs, msg)
		}
	})
	defer log.SetListener(nil)

	for _, id := range []ManagedIDKind{ClientID(fakeClientID), ObjectID(fakeObjectID), ResourceID(fakeResourceID), nil} {
		_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: id})
		require.NoError(t, err)
		require.Len(t, logs, 1)
		require.Contains(t, logs[0], "IMDS")
		if id != nil {
			require.Contains(t, logs[0], id.String())
			kind := ""
			switch id.(type) {
			case ClientID:
				kind = "client"
			case ObjectID:
				kind = "object"
			case ResourceID:
				kind = "resource"
			}
			require.Contains(t, logs[0], kind+" ID")
		}
		logs = nil
	}
}

func TestManagedIdentityCredential_UnexpectedIMDSResponse(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	tests := [][]mock.ResponseOption{
		{mock.WithBody([]byte("not json")), mock.WithStatusCode(http.StatusOK)},
	}
	// credential should return AuthenticationFailedError when a token request ends with a retriable response
	ro := policy.RetryOptions{}
	setIMDSRetryOptionDefaults(&ro)
	for _, c := range ro.StatusCodes {
		tests = append(tests, []mock.ResponseOption{mock.WithStatusCode(c)})
	}
	for _, res := range tests {
		srv.AppendResponse(res...)

		c, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry:     policy.RetryOptions{MaxRetries: -1},
				Transport: srv,
			},
		})
		require.NoError(t, err)

		_, err = c.GetToken(ctx, testTRO)
		var af *AuthenticationFailedError
		require.ErrorAs(t, err, &af, "unexpected token response from IMDS should prompt an AuthenticationFailedError")
	}
}

func TestManagedIdentityCredential_ServiceFabric(t *testing.T) {
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

func TestManagedIdentityCredential_UnsupportedID(t *testing.T) {
	t.Run("Azure Arc", func(t *testing.T) {
		t.Setenv(identityEndpoint, fakeMIEndpoint)
		t.Setenv(arcIMDSEndpoint, fakeMIEndpoint)
		for _, id := range []ManagedIDKind{ClientID(fakeClientID), ObjectID(fakeObjectID), ResourceID(fakeResourceID)} {
			_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: id})
			require.Errorf(t, err, "expected an error for %T", id)
		}
	})
	t.Run("Azure ML", func(t *testing.T) {
		t.Setenv(msiEndpoint, fakeMIEndpoint)
		t.Setenv(msiSecret, "...")
		_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ResourceID(fakeResourceID)})
		require.Error(t, err)
		_, err = NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ObjectID(fakeObjectID)})
		require.Error(t, err)
		_, err = NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: ClientID(fakeClientID)})
		require.NoError(t, err)
	})
	t.Run("Cloud Shell", func(t *testing.T) {
		t.Setenv(msiEndpoint, fakeMIEndpoint)
		for _, id := range []ManagedIDKind{ClientID(fakeClientID), ObjectID(fakeObjectID), ResourceID(fakeResourceID)} {
			_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: id})
			require.Errorf(t, err, "expected an error for %T", id)
		}
	})
	t.Run("Service Fabric", func(t *testing.T) {
		t.Setenv(identityEndpoint, fakeMIEndpoint)
		t.Setenv(identityHeader, "...")
		t.Setenv(identityServerThumbprint, "...")
		for _, id := range []ManagedIDKind{ClientID(fakeClientID), ObjectID(fakeObjectID), ResourceID(fakeResourceID)} {
			_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: id})
			require.Errorf(t, err, "expected an error for %T", id)
		}
	})
}
