//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/tls"
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
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
		ClientID:      fakeClientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      fakeTenantID,
		TokenFilePath: tempFile,
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
		ClientID:      fakeClientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      fakeTenantID,
		TokenFilePath: tempFile,
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

func TestWorkloadIdentityCredential_SNIPolicy(t *testing.T) {
	called := false
	expected := ""
	newServer := func() ([]byte, *url.URL) {
		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintln(w, string(accessTokenRespSuccess))
		}))
		t.Cleanup(ts.Close)
		ts.TLS = &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				called = true
				if expected == "" {
					t.Error("test bug: expected server name not set; should match test server's DNS name")
				} else if actual := info.ServerName; actual != expected {
					t.Errorf("expected %q, got %q", expected, actual)
				}
				return nil, nil
			},
		}
		ts.StartTLS()
		cert := ts.Certificate()
		expected = cert.DNSNames[0]
		pemData := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)
		return pemData, u
	}

	pemData, u := newServer()
	caFile := filepath.Join(t.TempDir(), t.Name())
	require.NoError(t, os.WriteFile(caFile, pemData, 0600))

	f := filepath.Join(t.TempDir(), t.Name())
	require.NoError(t, os.WriteFile(f, []byte(tokenValue), 0600))

	for k, v := range map[string]string{
		aksSNIName:              expected,
		aksTokenEndpoint:        u.Host,
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: f,
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}

	for _, test := range []struct {
		name string
		vars map[string]string
	}{
		{"no cert specified", nil},
		{"two certs specified", map[string]string{aksCAData: "...", aksCAFile: "..."}},
	} {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.vars {
				t.Setenv(k, v)
			}
			_, err := NewWorkloadIdentityCredential(nil)
			require.ErrorContains(t, err, aksCAData)
			require.ErrorContains(t, err, aksCAFile)
		})
	}
	o := WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Transport: &mockSTS{
				tokenRequestCallback: func(*http.Request) *http.Response {
					t.Fatal("credential should have sent token request to endpoint specified in " + aksTokenEndpoint)
					return nil
				},
			},
		},
	}
	for k, v := range map[string]string{
		aksCAData: string(pemData),
		aksCAFile: caFile,
	} {
		called = false
		t.Run(k, func(t *testing.T) {
			t.Setenv(k, v)
			cred, err := NewWorkloadIdentityCredential(&o)
			require.NoError(t, err)

			tk, err := cred.GetToken(ctx, testTRO)
			require.NoError(t, err)
			require.Equal(t, tokenValue, tk.Token)
			require.True(t, called, "test bug: test server's GetCertificate function wasn't called")

			t.Run("race", func(t *testing.T) {
				cred, err := NewWorkloadIdentityCredential(&o)
				require.NoError(t, err)
				wg := sync.WaitGroup{}
				ch := make(chan error, 1)
				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						if _, err := cred.GetToken(ctx, testTRO); err != nil {
							select {
							case ch <- err:
							default:
							}
						}
					}()
				}
				wg.Wait()
				select {
				case err := <-ch:
					t.Fatal(err)
				default:
				}
			})
		})
	}

	t.Run("file", func(t *testing.T) {
		t.Setenv(aksCAFile, caFile)
		t.Run("updated", func(t *testing.T) {
			p, err := newAKSTokenRequestPolicy()
			require.NoError(t, err)
			pl := runtime.NewPipeline("", "", runtime.PipelineOptions{}, &policy.ClientOptions{
				PerRetryPolicies: []policy.Policy{p},
				Transport: &mockSTS{
					tokenRequestCallback: func(*http.Request) *http.Response {
						t.Fatal("policy should have sent this request to the AKS endpoint")
						return nil
					},
				},
			})

			called = false
			r, err := runtime.NewRequest(ctx, http.MethodGet, u.String()+"/tenant/token")
			require.NoError(t, err)
			_, err = pl.Do(r)
			require.NoError(t, err)
			require.True(t, called, "test bug: test server's GetCertificate function wasn't called")

			// need a new server because a started one's TLS cert is immutable. Unfortunately, a new
			// server listens on a different port, so we need to update the policy's host. This is
			// why this test exercises the policy directly rather than through a credential instance
			pemData, u := newServer()
			p.host = u.Host
			require.NoError(t, os.WriteFile(caFile, pemData, 0600))

			called = false
			r, err = runtime.NewRequest(ctx, http.MethodGet, u.String()+"/tenant/token")
			require.NoError(t, err)
			_, err = pl.Do(r)
			require.NoError(t, err)
			require.True(t, called, "test bug: test server's GetCertificate function wasn't called")
		})
		t.Run("invalid", func(t *testing.T) {
			require.NoError(t, os.WriteFile(caFile, []byte("not a cert"), 0600))
			_, err := NewWorkloadIdentityCredential(nil)
			require.ErrorContains(t, err, "couldn't parse")
			require.ErrorContains(t, err, aksCAFile)
		})
		t.Run("empty", func(t *testing.T) {
			require.NoError(t, os.Truncate(caFile, 0))
			_, err := NewWorkloadIdentityCredential(nil)
			require.ErrorContains(t, err, "empty file")
		})
		t.Run("not found", func(t *testing.T) {
			require.NoError(t, os.Remove(caFile))
			_, err := NewWorkloadIdentityCredential(nil)
			require.ErrorContains(t, err, "no such file")
		})
	})
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
		ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
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
		ClientID:      clientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      tenantID,
		TokenFilePath: rightFile,
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
