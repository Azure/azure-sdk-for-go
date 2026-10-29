// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// liveScopes is scopes to use for live managed identity tests. Each such test must use a different
	// scope because MSAL caches managed identity tokens at the level of the process. If two tests use
	// the same scope, the second will get cached tokens and not send a request.
	liveScopes = []string{
		"https://graph.microsoft.com/.default",
		"https://management.azure.com/.default",
		"https://storage.azure.com/.default",
		"https://vault.azure.net/.default",
	}
	liveScopesMtx = &sync.Mutex{}
)

// recordMITest starts a managed identity test recording and returns ClientOptions and a scope for the test
func recordMITest(t *testing.T) (azcore.ClientOptions, string) {
	opts, stop := initRecording(t)
	t.Cleanup(stop)
	recordedScope := strings.ReplaceAll(t.Name(), "/", "_")
	if recording.GetRecordMode() == recording.PlaybackMode {
		return opts, recordedScope
	}
	liveScopesMtx.Lock()
	defer liveScopesMtx.Unlock()
	if len(liveScopes) == 0 {
		// fail instead of skipping because this would be a bug
		t.Fatal("all live managed identity test scopes have been used")
	}
	scope := liveScopes[0]
	liveScopes = liveScopes[1:]
	if recording.GetRecordMode() == recording.RecordingMode {
		actual := url.QueryEscape(strings.TrimSuffix(scope, defaultSuffix))
		err := recording.AddURISanitizer(recordedScope, actual, &recording.RecordingOptions{
			ProxyPort:    os.Getpid()%10000 + 20000, // TODO
			TestInstance: t,
			UseHTTPS:     true,
		})
		require.NoError(t, err)
	}
	return opts, scope
}

// triggerLiveTest sends a request to managed-id-test at the given URL and fails the test if the response isn't OK
func triggerLiveTest(t *testing.T, url string) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NotNil(t, res)
	if res.StatusCode != http.StatusOK {
		b, err := azruntime.Payload(res)
		require.NoError(t, err)
		t.Fatal("test application returned an error: " + string(b))
	}
}

func writeArcKeyFile(t *testing.T, content string) string {
	d := ""
	switch o := runtime.GOOS; o {
	case "linux":
		d = "/var/opt/azcmagent/tokens"
	case "windows":
		pd := os.Getenv("ProgramData")
		if pd == "" {
			t.Fatal("environment variable ProgramData has no value")
		}
		d = filepath.Join(pd, "AzureConnectedMachineAgent", "Tokens")
	default:
		t.Skipf("unsupported OS %q", o)
	}
	if _, err := os.Stat(d); err != nil {
		if err = os.MkdirAll(d, 0755); err != nil {
			t.Skipf("failed to create Arc key directory: %v", err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(d) })
	}
	p := filepath.Join(d, "arc.key")
	err := os.WriteFile(p, []byte(content), 0600)
	if err != nil {
		t.Skipf("failed to write Arc key file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(p) })
	return p
}

func TestManagedIdentityCredential_AzureArc(t *testing.T) {
	expectedKey := "expected-key"
	fp := writeArcKeyFile(t, expectedKey)
	expectedPath := "/foo/token"
	expectedScope := t.Name() + "/.default"
	validateReq := func(req *http.Request) bool {
		if req.URL.Path != expectedPath {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		if p := req.URL.Query().Get("api-version"); p != azureArcAPIVersion {
			t.Fatalf("unexpected api-version: %s", p)
		}
		if p := req.URL.Query().Get("resource"); p != strings.TrimSuffix(expectedScope, defaultSuffix) {
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
	srv.AppendResponse(mock.WithHeader("WWW-Authenticate", "Basic realm="+fp), mock.WithStatusCode(401))
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
	testGetTokenSuccess(t, cred, expectedScope)
}

func newAzureArcUserAssignedTestServer(t *testing.T, id ManagedIDKind, selector, echo string) *mock.Server {
	t.Helper()
	expectedKey := "expected-key"
	keyFile := writeArcKeyFile(t, expectedKey)
	body := map[string]any{
		"access_token": tokenValue,
		"expires_in":   tokenExpiresIn,
		"token_type":   "Bearer",
	}
	if echo != "" {
		body[selector] = echo
	}
	response, err := json.Marshal(body)
	require.NoError(t, err)

	var requests atomic.Int32
	validateReq := func(req *http.Request) bool {
		assert.Equal(t, "/foo/token", req.URL.Path)
		query := req.URL.Query()
		assert.Equal(t, azureArcAPIVersion, query.Get("api-version"))
		assert.Equal(t, t.Name(), query.Get("resource"))
		assert.Equal(t, id.String(), query.Get(selector))
		for _, other := range []string{qpClientID, "object_id", "msi_res_id"} {
			if other != selector {
				assert.False(t, query.Has(other), "unexpected identity selector %s", other)
			}
		}
		assert.Equal(t, "true", req.Header.Get(headerMetadata))
		if requests.Add(1) == 1 {
			assert.Empty(t, req.Header.Get("Authorization"))
		} else {
			assert.Equal(t, "Basic "+expectedKey, req.Header.Get("Authorization"))
		}
		return true
	}
	srv, close := mock.NewServer()
	t.Cleanup(close)
	srv.AppendResponse(
		mock.WithPredicate(validateReq),
		mock.WithHeader("WWW-Authenticate", "Basic realm="+keyFile),
		mock.WithStatusCode(http.StatusUnauthorized),
	)
	srv.AppendResponse()
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(response))
	srv.AppendResponse()
	setEnvironmentVariables(t, map[string]string{
		arcIMDSEndpoint:          srv.URL(),
		identityEndpoint:         srv.URL() + "/foo/token",
		identityHeader:           "",
		identityServerThumbprint: "",
		msiEndpoint:              "",
		msiSecret:                "",
	})
	return srv
}

func TestManagedIdentityCredential_AzureArcUserAssigned(t *testing.T) {
	for _, identity := range []struct {
		name, selector string
		id             ManagedIDKind
	}{
		{name: "ClientID", selector: qpClientID, id: ClientID(fakeClientID)},
		{name: "ObjectID", selector: "object_id", id: ObjectID(fakeObjectID)},
		{name: "ResourceID", selector: "msi_res_id", id: ResourceID(fakeResourceID)},
	} {
		t.Run(identity.name, func(t *testing.T) {
			for _, response := range []struct {
				name, echo string
				success    bool
			}{
				{name: "matching identity", echo: identity.id.String(), success: true},
				{name: "missing identity"},
				{name: "mismatched identity", echo: "unexpected-" + identity.id.String()},
			} {
				t.Run(response.name, func(t *testing.T) {
					srv := newAzureArcUserAssignedTestServer(t, identity.id, identity.selector, response.echo)
					cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
						ClientOptions: azcore.ClientOptions{
							Retry:     policy.RetryOptions{MaxRetries: -1},
							Transport: srv,
						},
						ID: identity.id,
					})
					require.NoError(t, err)
					token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name() + defaultSuffix}})
					if response.success {
						require.NoError(t, err)
						require.Equal(t, tokenValue, token.Token)
					} else {
						var authErr *AuthenticationFailedError
						require.ErrorAs(t, err, &authErr)
						require.ErrorContains(t, err, "azure arc did not confirm the requested user-assigned managed identity")
						require.Equal(t, azcore.AccessToken{}, token)
					}
					require.Equal(t, 2, srv.Requests(), "both requests in the Arc challenge flow must be sent")
				})
			}
		})
	}
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
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
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
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
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
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		// TODO: check for specific error
		if err == nil {
			t.Fatal("expected an error")
		}
	})
	t.Run("key too large", func(t *testing.T) {
		size := 4097
		f := writeArcKeyFile(t, strings.Repeat(" ", size))
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.AppendResponse(
			mock.WithHeader("WWW-Authenticate", "Basic realm="+f),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.ErrorContains(t, err, fmt.Sprint(size))
	})
	t.Run("unexpected file paths", func(t *testing.T) {
		if n := runtime.GOOS; n != "linux" && n != "windows" {
			t.Skipf("unsupported OS %q", n)
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
		_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.ErrorContains(t, err, "invalid file path")

		srv.AppendResponse(
			// unexpected extension
			mock.WithHeader("WWW-Authenticate", "Basic realm="+filepath.Join(t.TempDir(), "foo")),
			mock.WithStatusCode(http.StatusUnauthorized),
		)
		cred, err = NewManagedIdentityCredential(&o)
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.ErrorContains(t, err, "invalid file")
	})
}

func TestManagedIdentityCredential_AzureContainerInstanceLive(t *testing.T) {
	// This test triggers the managed identity test app deployed to an Azure Container Instance.
	// See the bicep file and test resources scripts for details.
	ip := os.Getenv("AZIDENTITY_ACI_IP")
	if ip == "" {
		t.Skip("set AZIDENTITY_ACI_IP to run this test")
	}
	triggerLiveTest(t, "http://"+ip)
}

func TestManagedIdentityCredential_AzureFunctionsLive(t *testing.T) {
	// This test triggers the managed identity test app deployed to Azure Functions.
	// See the bicep file and test resources scripts for details.
	fn := os.Getenv("AZIDENTITY_FUNCTION_NAME")
	if fn == "" {
		t.Skip("set AZIDENTITY_FUNCTION_NAME to run this test")
	}
	url := fmt.Sprintf("https://%s.azurewebsites.net/api/HttpTrigger", fn)
	triggerLiveTest(t, url)
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
	opts, scope := recordMITest(t)
	cred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: opts})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred, scope)
}

func TestManagedIdentityCredential_CloudShell(t *testing.T) {
	validateReq := func(req *http.Request) *http.Response {
		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		if v := req.FormValue("resource"); v != strings.TrimSuffix(t.Name(), defaultSuffix) {
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
	testGetTokenSuccess(t, msiCred, t.Name())
}

func TestManagedIdentityCredential_AppService(t *testing.T) {
	expectedID := "expected-ID"
	expectedHeader := "header"
	for _, id := range []ManagedIDKind{ClientID(expectedID), ObjectID(expectedID), ResourceID(expectedID), nil} {
		scope := fmt.Sprintf("%s/%T/.default", t.Name(), id)
		validateReq := func(t *testing.T) func(req *http.Request) bool {
			return func(req *http.Request) bool {
				if h := req.Header.Get("X-IDENTITY-HEADER"); h != expectedHeader {
					t.Fatalf("unexpected X-IDENTITY-HEADER: %s", h)
				}
				q := req.URL.Query()
				if v := q.Get("api-version"); v != "2019-08-01" {
					t.Fatalf(`unexpected api-version "%s"`, v)
				}
				if v := q.Get("resource"); v != strings.TrimSuffix(scope, "/.default") {
					t.Fatalf(`unexpected resource "%s"`, v)
				}
				clientID := q.Get(qpClientID)
				resID := q.Get(miResID)
				objectID := q.Get("object_id")
				if (clientID != "" && resID != "") || (clientID != "" && objectID != "") || (resID != "" && objectID != "") {
					t.Fatal("request includes two IDs")
				}
				if id == nil {
					if clientID != "" || resID != "" || objectID != "" {
						t.Fatal("request shouldn't include a user-assigned ID")
					}
				} else {
					actual := clientID
					switch id.(type) {
					case ObjectID:
						actual = objectID
					case ResourceID:
						actual = resID
					}
					if actual != id.String() {
						t.Errorf("expected %s, got %q", id.String(), actual)
					}
				}
				return true
			}
		}

		t.Run(fmt.Sprintf("%T", id), func(t *testing.T) {
			srv, close := mock.NewServer()
			defer close()
			srv.AppendResponse(
				mock.WithPredicate(validateReq(t)),
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
			testGetTokenSuccess(t, cred, scope)
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
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
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
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
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
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
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

func TestManagedIdentityCredential_ExpiresOnInt(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	expires := time.Now().Add(time.Hour).Unix()
	scope := t.Name()
	srv.AppendResponse(
		mock.WithBody([]byte(fmt.Sprintf(
			`{"access_token":%q,"expires_on":%d,"resource":%q,"token_type":"Bearer"}`, tokenValue, expires, scope,
		))),
	)
	setEnvironmentVariables(t, map[string]string{msiEndpoint: srv.URL()})
	options := ManagedIdentityCredentialOptions{}
	options.Transport = srv
	msiCred, err := NewManagedIdentityCredential(&options)
	require.NoError(t, err)
	tk, err := msiCred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{scope}})
	require.NoError(t, err)
	require.Equal(t, tokenValue, tk.Token)
	require.Equal(t, expires, tk.ExpiresOn.Unix())
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
	_, err = msiCred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
	if err == nil {
		t.Fatalf("expected to receive an error but received none")
	}
}

func TestManagedIdentityCredential_IMDSLive(t *testing.T) {
	// This test triggers the managed identity test app deployed to a VM.
	// See the bicep file and test resources scripts for details.
	ip := os.Getenv("AZIDENTITY_VM_IP")
	if ip == "" {
		t.Skip("set AZIDENTITY_VM_IP to run this test")
	}
	triggerLiveTest(t, "http://"+ip)
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
	if !cred.mic.imds {
		t.SkipNow()
	}
	for _, code := range []int{404, 410, 429, 500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511} {
		reqs := 0
		sts.tokenRequestCallback = func(r *http.Request) *http.Response {
			reqs++
			return &http.Response{Body: http.NoBody, Request: r, StatusCode: code}
		}
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{fmt.Sprint(code)}})
		if err == nil {
			t.Fatalf("expected an error; %d", code)
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
	for i, res := range tests {
		srv.AppendResponse(res...)

		c, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry:     policy.RetryOptions{MaxRetries: -1},
				Transport: srv,
			},
		})
		require.NoError(t, err)

		_, err = c.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name() + fmt.Sprint(i)}})
		var af *AuthenticationFailedError
		require.ErrorAs(t, err, &af, "unexpected token response from IMDS should prompt an AuthenticationFailedError")
	}
}

func TestManagedIdentityCredential_ServiceFabric(t *testing.T) {
	scope := t.Name()
	expectedSecret := "expected-secret"
	validateReq := func(req *http.Request) bool {
		if secret := req.Header.Get("Secret"); secret != expectedSecret {
			t.Fatalf("unexpected Secret header %q", secret)
		}
		if p := req.URL.Query().Get("api-version"); p != serviceFabricAPIVersion {
			t.Fatalf("unexpected api-version: %s", p)
		}
		if p := req.URL.Query().Get("resource"); p != scope {
			t.Fatalf("unexpected resource: %s", p)
		}
		return true
	}

	// Serve TLS with a certificate whose thumbprint the test knows, so azidentity's owned transport
	// exercises MSAL's certificate pinning end to end. No Transport override, so azidentity owns the
	// transport and can pin.
	cert := allCertTests[0].certs[0]
	tlsCert := tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: allCertTests[0].key, Leaf: cert}
	srv, close := mock.NewTLSServer(mock.WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{tlsCert}}))
	defer close()
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()

	// Service Fabric publishes the SHA-1 thumbprint of its endpoint certificate.
	thumbprint := sha1.Sum(cert.Raw)
	setEnvironmentVariables(t, map[string]string{
		identityEndpoint:         srv.URL(),
		identityHeader:           expectedSecret,
		identityServerThumbprint: hex.EncodeToString(thumbprint[:]),
	})
	cred, err := NewManagedIdentityCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred, scope)
}

type nonHTTPClientTransporter struct{}

func (nonHTTPClientTransporter) Do(*http.Request) (*http.Response, error) {
	return nil, nil
}

// TestManagedIdentityCredential_ServiceFabricCallerTransport verifies a caller-supplied *http.Client
// Transport becomes the base MSAL pins, while a non-*http.Client transporter is rejected.
func TestManagedIdentityCredential_ServiceFabricCallerTransport(t *testing.T) {
	cert := allCertTests[0].certs[0]
	thumbprint := sha1.Sum(cert.Raw)

	t.Run("custom *http.Client is used as the pinning base", func(t *testing.T) {
		scope := t.Name()
		tlsCert := tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: allCertTests[0].key, Leaf: cert}
		srv, close := mock.NewTLSServer(mock.WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{tlsCert}}))
		defer close()
		srv.AppendResponse(mock.WithBody(accessTokenRespSuccess))
		srv.AppendResponse()

		// A Proxy hook on the caller's transport records that its transport is the one making requests.
		var mu sync.Mutex
		proxyCalled := false
		callerTransport := &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				mu.Lock()
				proxyCalled = true
				mu.Unlock()
				return nil, nil
			},
		}
		setEnvironmentVariables(t, map[string]string{
			identityEndpoint:         srv.URL(),
			identityHeader:           "expected-secret",
			identityServerThumbprint: hex.EncodeToString(thumbprint[:]),
		})
		opts := &ManagedIdentityCredentialOptions{}
		opts.Transport = &http.Client{Transport: callerTransport}
		cred, err := NewManagedIdentityCredential(opts)
		require.NoError(t, err)
		testGetTokenSuccess(t, cred, scope)
		mu.Lock()
		defer mu.Unlock()
		require.True(t, proxyCalled, "expected the caller's transport to be used as the pinning base")
	})

	t.Run("non-*http.Client Transport is rejected", func(t *testing.T) {
		setEnvironmentVariables(t, map[string]string{
			identityEndpoint:         fakeMIEndpoint,
			identityHeader:           "expected-secret",
			identityServerThumbprint: hex.EncodeToString(thumbprint[:]),
		})
		opts := &ManagedIdentityCredentialOptions{}
		opts.Transport = nonHTTPClientTransporter{}
		_, err := NewManagedIdentityCredential(opts)
		require.Error(t, err)
		require.Contains(t, err.Error(), "*http.Client")
	})
}

func TestManagedIdentityCredential_ServiceFabricThumbprintMismatch(t *testing.T) {
	for _, customTransport := range []bool{false, true} {
		name := "default transport"
		if customTransport {
			name = "custom *http.Client"
		}
		t.Run(name, func(t *testing.T) {
			cert := allCertTests[0].certs[0]
			tlsCert := tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: allCertTests[0].key, Leaf: cert}
			srv, close := mock.NewTLSServer(mock.WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{tlsCert}}))
			defer close()
			var requests atomic.Int32
			srv.AppendResponse(mock.WithPredicate(func(*http.Request) bool {
				requests.Add(1)
				return true
			}), mock.WithBody(accessTokenRespSuccess))
			srv.AppendResponse()

			thumbprint := sha1.Sum(cert.Raw)
			thumbprint[0] ^= 0xff
			setEnvironmentVariables(t, map[string]string{
				identityEndpoint:         srv.URL(),
				identityHeader:           "expected-secret",
				identityServerThumbprint: hex.EncodeToString(thumbprint[:]),
			})
			opts := &ManagedIdentityCredentialOptions{
				ClientOptions: azcore.ClientOptions{Retry: policy.RetryOptions{MaxRetries: -1}},
			}
			if customTransport {
				opts.Transport = &http.Client{Transport: &http.Transport{}}
			}
			cred, err := NewManagedIdentityCredential(opts)
			require.NoError(t, err)
			token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
			var authErr *AuthenticationFailedError
			require.ErrorAs(t, err, &authErr)
			require.Contains(t, err.Error(), "TLS certificate thumbprint from Service Fabric did not match IDENTITY_SERVER_THUMBPRINT")
			require.Empty(t, token.Token)
			require.Zero(t, requests.Load(), "certificate verification must fail before sending an HTTP request")
		})
	}
}

func TestManagedIdentityCredential_ServiceFabricRedirect(t *testing.T) {
	for _, customTransport := range []bool{false, true} {
		name := "default transport"
		if customTransport {
			name = "custom *http.Client"
		}
		t.Run(name, func(t *testing.T) {
			for _, status := range []int{
				http.StatusMovedPermanently,
				http.StatusFound,
				http.StatusSeeOther,
				http.StatusTemporaryRedirect,
				http.StatusPermanentRedirect,
			} {
				t.Run(fmt.Sprint(status), func(t *testing.T) {
					cert := allCertTests[0].certs[0]
					tlsCert := tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: allCertTests[0].key, Leaf: cert}
					tlsConfig := &tls.Config{Certificates: []tls.Certificate{tlsCert}}
					target, closeTarget := mock.NewTLSServer(mock.WithTLSConfig(tlsConfig))
					defer closeTarget()
					var targetRequests atomic.Int32
					target.AppendResponse(mock.WithPredicate(func(*http.Request) bool {
						targetRequests.Add(1)
						return true
					}), mock.WithBody(accessTokenRespSuccess))
					target.AppendResponse()

					srv, close := mock.NewTLSServer(mock.WithTLSConfig(tlsConfig))
					defer close()
					var requests atomic.Int32
					srv.AppendResponse(mock.WithPredicate(func(*http.Request) bool {
						requests.Add(1)
						return true
					}), mock.WithStatusCode(status), mock.WithHeader("Location", target.URL()))
					srv.AppendResponse()

					thumbprint := sha1.Sum(cert.Raw)
					setEnvironmentVariables(t, map[string]string{
						identityEndpoint:         srv.URL(),
						identityHeader:           "expected-secret",
						identityServerThumbprint: hex.EncodeToString(thumbprint[:]),
					})
					opts := &ManagedIdentityCredentialOptions{
						ClientOptions: azcore.ClientOptions{Retry: policy.RetryOptions{MaxRetries: -1}},
					}
					if customTransport {
						opts.Transport = &http.Client{Transport: &http.Transport{}}
					}
					cred, err := NewManagedIdentityCredential(opts)
					require.NoError(t, err)
					token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
					var authErr *AuthenticationFailedError
					require.ErrorAs(t, err, &authErr)
					require.Contains(t, err.Error(), "redirects are not permitted for managed identity on Service Fabric")
					require.Empty(t, token.Token)
					require.EqualValues(t, 1, requests.Load(), "the pinned endpoint must receive the request")
					require.Zero(t, targetRequests.Load(), "the redirect target must not receive a request")
				})
			}
		})
	}
}

func TestManagedIdentityCredential_UnsupportedID(t *testing.T) {
	t.Run("Azure Arc", func(t *testing.T) {
		t.Setenv(identityEndpoint, fakeMIEndpoint)
		t.Setenv(arcIMDSEndpoint, fakeMIEndpoint)
		for _, id := range []ManagedIDKind{ClientID(fakeClientID), ObjectID(fakeObjectID), ResourceID(fakeResourceID)} {
			_, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ID: id})
			require.NoError(t, err)
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
