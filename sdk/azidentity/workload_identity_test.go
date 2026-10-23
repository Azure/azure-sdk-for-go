// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/customtokenproxy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func assertion(cert *x509.Certificate, key crypto.PrivateKey) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", liveSP.tenantID),
		"exp": json.Number(strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)),
		"iss": liveSP.clientID,
		"jti": uuid.New().String(),
		"nbf": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		"sub": liveSP.clientID,
	})
	x5t := sha1.Sum(cert.Raw) // nosec
	j.Header = map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
		"x5t": base64.StdEncoding.EncodeToString(x5t[:]),
	}
	return j.SignedString(key)
}

func TestWorkloadIdentityCredential_Live(t *testing.T) {
	// This test triggers the managed identity test app deployed to Azure Kubernetes Service.
	// See the bicep file and test resources scripts for details.
	// It triggers the app with kubectl because the test subscription prohibits opening ports to the internet.
	pod := os.Getenv("AZIDENTITY_POD_NAME")
	if pod == "" {
		t.Skip("set AZIDENTITY_POD_NAME to run this test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "kubectl", "exec", pod, "--", "wget", "-qO-", "localhost")
	b, err := cmd.CombinedOutput()
	s := string(b)
	require.NoError(t, err, s)
	require.Equal(t, "test passed", s)
}

func TestWorkloadIdentityCredential_Recorded(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
	// workload identity and client cert auth use the same flow. This test
	// implements cert auth with WorkloadIdentityCredential as a way to test
	// that credential in an environment that's easier to set up than AKS
	cert, err := os.ReadFile(liveSP.pemPath)
	if err != nil {
		t.Fatal(err)
	}
	certs, key, err := ParseCertificates(cert, nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := assertion(certs[0], key)
	if err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(t.TempDir(), t.Name())
	if err := os.WriteFile(f, []byte(a), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	for _, b := range []bool{true, false} {
		name := "default options"
		if b {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			co, stop := initRecording(t)
			defer stop()
			cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
				ClientID:                 liveSP.clientID,
				ClientOptions:            co,
				DisableInstanceDiscovery: b,
				EnableAzureProxy:         true,
				TenantID:                 liveSP.tenantID,
				TokenFilePath:            f,
			})
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestWorkloadIdentityCredential(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	sts := mockSTS{tenant: fakeTenantID, tokenRequestCallback: func(req *http.Request) *http.Response {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != tokenValue {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual, ok := req.PostForm["client_id"]; !ok {
			t.Error("expected a client_id")
		} else if len(actual) != 1 || actual[0] != fakeClientID {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual := strings.Split(req.URL.Path, "/")[1]; actual != fakeTenantID {
			t.Errorf(`unexpected tenant "%s"`, actual)
		}
		return nil
	}}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		ClientOptions:    policy.ClientOptions{Transport: &sts},
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    tempFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkloadIdentityCredential_Expiration(t *testing.T) {
	tokenReqs := 0
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	sts := mockSTS{tenant: fakeTenantID, tokenRequestCallback: func(req *http.Request) *http.Response {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != fmt.Sprint(tokenReqs) {
			t.Errorf(`expected assertion "%d", got "%s"`, tokenReqs, actual[0])
		}
		tokenReqs++
		return nil
	}}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		ClientOptions:    policy.ClientOptions{Transport: &sts},
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    tempFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		// tokenReqs counts requests, and its latest value is the expected client assertion and the requested scope.
		// Each iteration of this loop therefore sends a token request with a unique assertion.
		s := fmt.Sprint(tokenReqs)
		if err = os.WriteFile(tempFile, []byte(fmt.Sprint(s)), os.ModePerm); err != nil {
			t.Fatalf("failed to write token file: %v", err)
		}
		if _, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{s}}); err != nil {
			t.Fatal(err)
		}
		cred.expires = time.Now().Add(-time.Second)
	}
	if tokenReqs != 2 {
		t.Fatalf("expected 2 token requests, got %d", tokenReqs)
	}
}

func TestTestWorkloadIdentityCredential_IncompleteConfig(t *testing.T) {
	f := filepath.Join(t.TempDir(), t.Name())
	for _, env := range []map[string]string{
		{},

		{azureClientID: fakeClientID},
		{azureFederatedTokenFile: f},
		{azureTenantID: fakeTenantID},

		{azureClientID: fakeClientID, azureTenantID: fakeTenantID},
		{azureClientID: fakeClientID, azureFederatedTokenFile: f},
		{azureTenantID: fakeTenantID, azureFederatedTokenFile: f},
	} {
		t.Run("", func(t *testing.T) {
			for k, v := range env {
				t.Setenv(k, v)
			}
			if _, err := NewWorkloadIdentityCredential(nil); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestWorkloadIdentityCredential_NoFile(t *testing.T) {
	for k, v := range map[string]string{
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: filepath.Join(t.TempDir(), t.Name()),
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions:    policy.ClientOptions{Transport: &mockSTS{}},
		EnableAzureProxy: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = cred.GetToken(context.Background(), testTRO); err == nil {
		t.Fatal("expected an error")
	}
}

func TestWorkloadIdentityCredential_Options(t *testing.T) {
	clientID := "not-" + fakeClientID
	tenantID := "not-" + fakeTenantID
	wrongFile := filepath.Join(t.TempDir(), "wrong")
	rightFile := filepath.Join(t.TempDir(), "right")
	if err := os.WriteFile(rightFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	sts := mockSTS{
		tenant: tenantID,
		tokenRequestCallback: func(req *http.Request) *http.Response {
			if err := req.ParseForm(); err != nil {
				t.Error(err)
			}
			if actual, ok := req.PostForm["client_assertion"]; !ok {
				t.Error("expected a client_assertion")
			} else if len(actual) != 1 || actual[0] != tokenValue {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			if actual, ok := req.PostForm["client_id"]; !ok {
				t.Error("expected a client_id")
			} else if len(actual) != 1 || actual[0] != clientID {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			if actual := strings.Split(req.URL.Path, "/")[1]; actual != tenantID {
				t.Errorf(`unexpected tenant "%s"`, actual)
			}
			return nil
		},
	}
	// options should override environment variables
	for k, v := range map[string]string{
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: wrongFile,
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         clientID,
		ClientOptions:    policy.ClientOptions{Transport: &sts},
		EnableAzureProxy: true,
		TenantID:         tenantID,
		TokenFilePath:    rightFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token %q", tk.Token)
	}
}

// startTestTokenEndpointWithCAData starts a TLS server with custom token handler.
//
// This token server serves instance discovery, OIDC discovery document and token endpoints.
func startTestTokenEndpointWithCAData(t testing.TB, tokenHandler http.Handler) (*httptest.Server, string) {
	t.Helper()

	mux := http.NewServeMux()

	testServer := httptest.NewTLSServer(mux)

	openidDiscoveryDocumentHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		tenantID := parts[1]

		doc := tenantMetadata(r.Host, tenantID) // request.Host should be used in server side
		var m map[string]any
		if err := json.Unmarshal(doc, &m); err != nil {
			t.Fatalf("failed to unmarshal tenant metadata: %v", err)
		}
		m["token_endpoint"] = testServer.URL + "/token" // overriding the token endpoint to the proxied token endpoint

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(m); err != nil {
			t.Fatalf("failed to encode tenant metadata: %v", err)
		}
	})
	mux.Handle("/{tenantid}/v2.0/.well-known/openid-configuration", openidDiscoveryDocumentHandler)
	mux.HandleFunc("/common/discovery/instance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		host := r.Host // request.Host should be used in server side
		if ae := r.URL.Query().Get("authorization_endpoint"); ae != "" {
			u, err := url.Parse(ae)
			if err != nil {
				t.Fatalf("mockSTS failed to parse an authorization_endpoint query parameter: %v", err)
			}
			host = u.Host
		}
		_, _ = w.Write(instanceMetadata(host, "fake-tenant"))
	})
	mux.Handle("/token", tokenHandler)

	serverCert := testServer.Certificate()
	require.NotNil(t, serverCert)
	ca := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCert.Raw,
	})
	require.NotEmpty(t, ca)

	return testServer, string(ca)
}

const testClientAssertion = "test-client-assertion-abc9412787413"

// customTokenRequestPolicyFlowCheck validates the custom token request flow is
// obeying the policy token flows without skipping headers.
type customTokenRequestPolicyFlowCheck struct {
	requiredHeaderKey   string
	requiredHeaderValue string
}

func newCustomTokenRequestPolicyFlowCheck() *customTokenRequestPolicyFlowCheck {
	return &customTokenRequestPolicyFlowCheck{
		requiredHeaderKey:   "adhoc-injected-header",
		requiredHeaderValue: "adhoc-injected-header-value",
	}
}

func (c *customTokenRequestPolicyFlowCheck) Policy() policy.Policy {
	return policyFunc(func(req *policy.Request) (*http.Response, error) {
		req.Raw().Header.Set(c.requiredHeaderKey, c.requiredHeaderValue)

		return req.Next()
	})
}

func (c *customTokenRequestPolicyFlowCheck) Validate(t testing.TB, req *http.Request) {
	t.Helper()
	require.NotNil(t, req)
	require.NotNil(t, req.Header)

	// headers expecting from client even it's sending to custom token endpoint
	require.Equal(t, c.requiredHeaderValue, req.Header.Get(c.requiredHeaderKey))
}

func TestWorkloadIdentityCredential_CustomTokenEndpoint_WithCAData(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(testClientAssertion), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	policyFlowCheck := newCustomTokenRequestPolicyFlowCheck()

	customTokenEndointServerCalledTimes := new(atomic.Int32)
	customTokenEndointServer, caData := startTestTokenEndpointWithCAData(
		t,
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			customTokenEndointServerCalledTimes.Add(1)

			policyFlowCheck.Validate(t, req)

			require.NoError(t, req.ParseForm())
			require.NotEmpty(t, req.PostForm)

			require.Contains(t, req.PostForm, "client_assertion")
			require.Equal(t, req.PostForm.Get("client_assertion"), testClientAssertion)

			require.Contains(t, req.PostForm, "client_id")
			require.Equal(t, req.PostForm.Get("client_id"), fakeClientID)

			_, _ = w.Write(accessTokenRespSuccess)
		}),
	)

	t.Setenv(customtokenproxy.AzureKubernetesTokenProxy, customTokenEndointServer.URL)
	t.Setenv(customtokenproxy.AzureKubernetesCAData, caData)

	clientOptions := policy.ClientOptions{
		PerCallPolicies: []policy.Policy{
			policyFlowCheck.Policy(),
		},
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		ClientOptions:    clientOptions,
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    tempFile,
	})
	require.NoError(t, err)
	require.Nil(t, clientOptions.Transport, "constructor shouldn't mutate caller's ClientOptions")

	testGetTokenSuccess(t, cred)

	require.Equal(t, int32(1), customTokenEndointServerCalledTimes.Load())
}

func TestWorkloadIdentityCredential_CustomTokenEndpoint_InvalidSettings(t *testing.T) {
	t.Setenv(customtokenproxy.AzureKubernetesTokenProxy, "invalid-token-endpoint")
	_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    filepath.Join(t.TempDir(), "test-workload-token-file"),
	})
	require.Error(t, err)
}

func TestWorkloadIdentityCredential_CustomTokenEndpoint_WithCAFile(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(testClientAssertion), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	policyFlowCheck := newCustomTokenRequestPolicyFlowCheck()

	customTokenEndointServerCalledTimes := new(atomic.Int32)
	customTokenEndointServer, caData := startTestTokenEndpointWithCAData(
		t,
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			customTokenEndointServerCalledTimes.Add(1)

			policyFlowCheck.Validate(t, req)

			require.NoError(t, req.ParseForm())
			require.NotEmpty(t, req.PostForm)

			require.Contains(t, req.PostForm, "client_assertion")
			require.Equal(t, req.PostForm.Get("client_assertion"), testClientAssertion)

			require.Contains(t, req.PostForm, "client_id")
			require.Equal(t, req.PostForm.Get("client_id"), fakeClientID)

			_, _ = w.Write(accessTokenRespSuccess)
		}),
	)

	t.Setenv(customtokenproxy.AzureKubernetesTokenProxy, customTokenEndointServer.URL)
	d := t.TempDir()
	caFile := filepath.Join(d, "test-ca-file")
	require.NoError(t, os.WriteFile(caFile, []byte(caData), 0600))
	t.Setenv(customtokenproxy.AzureKubernetesCAFile, caFile)

	clientOptions := policy.ClientOptions{
		PerCallPolicies: []policy.Policy{
			policyFlowCheck.Policy(),
		},
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		ClientOptions:    clientOptions,
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    tempFile,
	})
	require.NoError(t, err)
	require.Nil(t, clientOptions.Transport, "constructor shouldn't mutate caller's ClientOptions")

	testGetTokenSuccess(t, cred)

	require.Equal(t, int32(1), customTokenEndointServerCalledTimes.Load())
}

func TestWorkloadIdentityCredential_CustomTokenEndpoint_AKSSetup(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(testClientAssertion), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	policyFlowCheck := newCustomTokenRequestPolicyFlowCheck()
	sniName := "test-sni.example.com"

	customTokenEndointServerCalledTimes := new(atomic.Int32)
	customTokenEndointServer, caData := startTestTokenEndpointWithCAData(
		t,
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			customTokenEndointServerCalledTimes.Add(1)

			policyFlowCheck.Validate(t, req)

			require.NotNil(t, req.TLS)
			require.Equal(t, req.TLS.ServerName, sniName, "when SNI is set, request should set SNI")

			require.NoError(t, req.ParseForm())
			require.NotEmpty(t, req.PostForm)

			require.Contains(t, req.PostForm, "client_assertion")
			require.Equal(t, req.PostForm.Get("client_assertion"), testClientAssertion)

			require.Contains(t, req.PostForm, "client_id")
			require.Equal(t, req.PostForm.Get("client_id"), fakeClientID)

			_, _ = w.Write(accessTokenRespSuccess)
		}),
	)

	t.Setenv(customtokenproxy.AzureKubernetesTokenProxy, customTokenEndointServer.URL)
	t.Setenv(customtokenproxy.AzureKubernetesSNIName, sniName)

	d := t.TempDir()
	caFile := filepath.Join(d, "test-ca-file")
	require.NoError(t, os.WriteFile(caFile, []byte(caData), 0600))
	t.Setenv(customtokenproxy.AzureKubernetesCAFile, caFile)

	clientOptions := policy.ClientOptions{
		PerCallPolicies: []policy.Policy{
			policyFlowCheck.Policy(),
		},
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:         fakeClientID,
		ClientOptions:    clientOptions,
		EnableAzureProxy: true,
		TenantID:         fakeTenantID,
		TokenFilePath:    tempFile,
	})
	require.NoError(t, err)
	require.Nil(t, clientOptions.Transport, "constructor shouldn't mutate caller's ClientOptions")

	testGetTokenSuccess(t, cred)

	require.Equal(t, int32(1), customTokenEndointServerCalledTimes.Load())
}
