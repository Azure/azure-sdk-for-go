// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/customtokenproxy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/managedidentity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAzureCredential_GetTokenSuccess(t *testing.T) {
	env := map[string]string{azureTenantID: fakeTenantID, azureClientID: fakeClientID, azureClientSecret: fakeSecret}
	setEnvironmentVariables(t, env)
	cred, err := NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	c := cred.chain.sources[0].(*EnvironmentCredential)
	c.cred.(*ClientSecretCredential).client.noCAE = fakeConfidentialClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}
}

func TestDefaultAzureCredential_AzureArcUserAssigned(t *testing.T) {
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTokenCredentials, credNameManagedIdentity)
	srv := newAzureArcUserAssignedTestServer(t, ClientID(fakeClientID), qpClientID, fakeClientID)
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Retry:     policy.RetryOptions{MaxRetries: -1},
			Transport: srv,
		},
	})
	require.NoError(t, err)
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name() + defaultSuffix}})
	require.NoError(t, err)
	require.Equal(t, tokenValue, token.Token)
	require.Equal(t, 2, srv.Requests(), "both requests in the Arc challenge flow must be sent")
}

func TestDefaultAzureCredential_AZURE_TOKEN_CREDENTIALS(t *testing.T) {
	if v, ok := os.LookupEnv(azureTokenCredentials); ok {
		require.NoError(t, os.Unsetenv(azureTokenCredentials))
		defer func() {
			_ = os.Setenv(azureTokenCredentials, v)
		}()
	}
	// configure EnvironmentCredential and WorkloadIdentityCredential
	// so those types appear in the chain instead of error reporters
	for k, v := range map[string]string{
		azureClientID:     fakeClientID,
		azureClientSecret: fakeSecret,
		azureTenantID:     fakeTenantID,
		// this file won't be read because the test doesn't request a token
		azureFederatedTokenFile: "/dev/null",
	} {
		t.Setenv(k, v)
	}
	fullChain := []azcore.TokenCredential{
		&EnvironmentCredential{},
		&WorkloadIdentityCredential{},
		&ManagedIdentityCredential{},
		&AzureCLICredential{},
		&AzureDeveloperCLICredential{},
		&AzurePowerShellCredential{},
	}
	firstDevIndex := 3

	t.Run("not set", func(t *testing.T) {
		actual, err := NewDefaultAzureCredential(nil)
		require.NoError(t, err)
		require.Equal(t, len(fullChain), len(actual.chain.sources))
		for i, c := range fullChain {
			require.IsType(t, c, actual.chain.sources[i])
		}
		t.Run("required", func(t *testing.T) {
			_, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{RequireAzureTokenCredentials: true})
			require.ErrorContains(t, err, azureTokenCredentials)
		})
	})

	t.Run("dev", func(t *testing.T) {
		t.Setenv(azureTokenCredentials, "dev")
		actual, err := NewDefaultAzureCredential(nil)
		require.NoError(t, err)
		require.Equal(t, len(fullChain)-firstDevIndex, len(actual.chain.sources))
		for i, c := range fullChain[firstDevIndex:] {
			require.IsType(t, c, actual.chain.sources[i])
		}
	})

	t.Run("prod", func(t *testing.T) {
		t.Setenv(azureTokenCredentials, "prod")
		actual, err := NewDefaultAzureCredential(nil)
		require.NoError(t, err)
		require.Equal(t, firstDevIndex, len(actual.chain.sources))
		for i, c := range fullChain[:firstDevIndex] {
			require.IsType(t, c, actual.chain.sources[i])
		}
	})

	for _, c := range []string{credNameEnvironment, credNameWorkloadIdentity, credNameManagedIdentity, credNameAzureCLI, credNameAzureDeveloperCLI} {
		t.Run(c, func(t *testing.T) {
			t.Setenv(azureTokenCredentials, c)
			actual, err := NewDefaultAzureCredential(nil)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual.chain.sources))
			require.Equal(t, "*azidentity."+c, fmt.Sprintf("%T", actual.chain.sources[0]))
			if c == credNameManagedIdentity {
				mic := actual.chain.sources[0].(*ManagedIdentityCredential)
				require.False(t, mic.mic.probeIMDS, "probeIMDS should be false when ManagedIdentityCredential is selected")
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(azureTokenCredentials, t.Name())
		_, err := NewDefaultAzureCredential(nil)
		require.ErrorContains(t, err, azureTokenCredentials)
	})
}

func TestDefaultAzureCredential_CLICredentialOptions(t *testing.T) {
	require := require.New(t)
	cred, err := NewDefaultAzureCredential(nil)
	require.NoError(err)
	var (
		az  *AzureCLICredential
		azd *AzureDeveloperCLICredential
	)
	for _, s := range cred.chain.sources {
		if az == nil {
			az, _ = s.(*AzureCLICredential)
		}
		if azd == nil {
			azd, _ = s.(*AzureDeveloperCLICredential)
		}
	}
	require.NotNil(az, "%T should be in the default chain", az)
	require.True(az.opts.inDefaultChain)
	require.NotNil(azd, "%T should be in the default chain", azd)
	require.True(azd.opts.inDefaultChain)
}

func TestDefaultAzureCredential_AzurePowerShellCredentialOptions(t *testing.T) {
	require := require.New(t)
	cred, err := NewDefaultAzureCredential(nil)
	require.NoError(err)
	var (
		azurePowerShell *AzurePowerShellCredential
	)
	for _, s := range cred.chain.sources {
		if azurePowerShell == nil {
			azurePowerShell, _ = s.(*AzurePowerShellCredential)
		}
	}
	require.NotNil(azurePowerShell, "%T should be in the default chain", azurePowerShell)
	require.True(azurePowerShell.opts.inDefaultChain)
}

func TestDefaultAzureCredential_ConstructorErrors(t *testing.T) {
	// ensure NewEnvironmentCredential returns an error
	t.Setenv(azureTenantID, "")

	logMsgs := []string{}
	log.SetListener(func(e log.Event, s string) {
		if e == EventAuthentication {
			logMsgs = append(logMsgs, s)
		}
	})

	cred, err := NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal(err)
	}
	// make GetToken return an error in any runtime environment
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = cred.GetToken(ctx, testTRO)
	if err == nil {
		t.Fatal("expected an error")
	}
	// these credentials' constructors returned errors because their configuration is absent;
	// those errors should be represented in the error returned by DefaultAzureCredential.GetToken()
	// and NewDefaultAzureCredential should have logged them
	for _, name := range []string{credNameEnvironment, credNameWorkloadIdentity} {
		matched, err := regexp.MatchString(name+`: .+\n`, err.Error())
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Errorf("expected an error message from %s", name)
		}
	}
	r := regexp.MustCompile(fmt.Sprintf(`(?m)NewDefaultAzureCredential failed to initialize some credentials:\n.*EnvironmentCredential:.+\n.*%s:`, credNameWorkloadIdentity))
	for _, msg := range logMsgs {
		if r.MatchString(msg) {
			return
		}
	}
	t.Fatalf("expected a log message about the constructor errors, got %s", strings.Join(logMsgs, "\n"))
}

func TestDefaultAzureCredential_TenantID(t *testing.T) {
	before := shellExec
	t.Cleanup(func() { shellExec = before })
	expected := "expected"
	for _, override := range []bool{false, true} {
		name := "default tenant"
		if override {
			name = "TenantID set"
		}
		for _, credName := range []string{credNameAzureCLI, credNameAzureDeveloperCLI} {
			t.Run(fmt.Sprintf("%s_%s", credName, name), func(t *testing.T) {
				called := false
				shellExec = func(ctx context.Context, actualName string, command string) ([]byte, error) {
					require.Equal(t, credName, actualName)
					called = true
					tenantArg := "--tenant"
					success, err := mockAzSuccess(ctx, actualName, command)
					if credName == credNameAzureDeveloperCLI {
						tenantArg = "--tenant-id"
						success, err = mockAzdSuccess(ctx, actualName, command)
					}
					sm := regexp.MustCompile(tenantArg + ` (\w+)`).FindStringSubmatch(command)
					if override {
						require.Equal(t, 2, len(sm), "tenant not found in command line")
						require.Equal(t, expected, sm[1], "unexpected tenant in command line")
					} else {
						require.Empty(t, sm)
					}
					return success, err
				}
				t.Setenv(azureTokenCredentials, credName)
				o := DefaultAzureCredentialOptions{}
				if override {
					o.TenantID = expected
				}
				cred, err := NewDefaultAzureCredential(&o)
				if err != nil {
					t.Fatal(err)
				}
				_, err = cred.GetToken(context.Background(), testTRO)
				if err != nil {
					t.Fatal(err)
				}
				if !called {
					t.Fatalf("%s wasn't invoked", credName)
				}
			})
		}
		t.Run(fmt.Sprintf("%s_%s", credNameWorkloadIdentity, name), func(t *testing.T) {
			af := filepath.Join(t.TempDir(), "assertions")
			if err := os.WriteFile(af, []byte("assertion"), os.ModePerm); err != nil {
				t.Fatal(err)
			}
			for k, v := range map[string]string{
				azureAuthorityHost:      "https://login.microsoftonline.com",
				azureClientID:           fakeClientID,
				azureFederatedTokenFile: af,
				azureTenantID:           "un" + expected,
			} {
				t.Setenv(k, v)
			}
			o := DefaultAzureCredentialOptions{
				ClientOptions: policy.ClientOptions{
					Transport: &mockSTS{
						tenant: expected,
						tokenRequestCallback: func(r *http.Request) *http.Response {
							if actual := strings.Split(r.URL.Path, "/")[1]; actual != expected {
								t.Fatalf("expected tenant %q, got %q", expected, actual)
							}
							return nil
						},
					},
				},
			}
			if override {
				o.TenantID = expected
			}
			cred, err := NewDefaultAzureCredential(&o)
			if err != nil {
				t.Fatal(err)
			}
			_, err = cred.GetToken(context.Background(), testTRO)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDefaultAzureCredential_UserAssignedIdentity(t *testing.T) {
	t.Setenv(azureClientID, fakeClientID)
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Transport: &mockSTS{
				tokenRequestCallback: func(req *http.Request) *http.Response {
					if req.Header.Get(headerMetadata) == "" {
						return nil
					}
					for _, p := range req.URL.Query() {
						for _, v := range p {
							if strings.Contains(v, fakeClientID) {
								return nil
							}
						}
					}
					t.Fatalf("expected %q in %v", fakeClientID, req.URL.Query())
					return nil
				},
			},
		},
	})
	require.NoError(t, err)
	_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name()}})
	require.NoError(t, err)
}

func TestDefaultAzureCredential_WorkloadIdentity(t *testing.T) {
	expectedAssertion := "service account token"
	tempFile := filepath.Join(t.TempDir(), "service-account-token-file")
	if err := os.WriteFile(tempFile, []byte(expectedAssertion), os.ModePerm); err != nil {
		t.Fatalf(`failed to write temporary file "%s": %v`, tempFile, err)
	}
	sts := mockSTS{tokenRequestCallback: func(req *http.Request) *http.Response {
		if err := req.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if actual := req.PostForm["client_assertion"]; actual[0] != expectedAssertion {
			t.Fatalf(`unexpected assertion "%s"`, actual[0])
		}
		if actual := req.PostForm["client_id"]; actual[0] != fakeClientID {
			t.Fatalf(`unexpected assertion "%s"`, actual[0])
		}
		if actual := strings.Split(req.URL.Path, "/")[1]; actual != fakeTenantID {
			t.Fatalf(`unexpected tenant "%s"`, actual)
		}
		return nil
	}}
	for k, v := range map[string]string{
		azureAuthorityHost:      cloud.AzurePublic.ActiveDirectoryAuthorityHost,
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: tempFile,
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{ClientOptions: policy.ClientOptions{Transport: &sts}})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)

	t.Run("disables identity binding mode", func(t *testing.T) {
		t.Setenv(azureTokenCredentials, credNameWorkloadIdentity)
		// these values should trigger validation errors if WorkloadIdentityCredential
		// tries to configure identity binding mode...
		t.Setenv(customtokenproxy.AzureKubernetesCAData, "not a valid cert")
		t.Setenv(customtokenproxy.AzureKubernetesTokenProxy, "http://timeout.local&fail=yes#please")

		cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
		})
		require.NoError(t, err)

		// ...but ensure a timeout should it try the proxy anyway
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = cred.GetToken(ctx, testTRO)
		require.NoError(t, err)
	})
}

func TestDefaultAzureCredential_IMDS(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{
		azureTokenCredentials: "prod", azureTenantID: "", azureClientID: "",
		azureClientSecret: "", azureFederatedTokenFile: "",
		identityEndpoint: "", identityHeader: "", identityServerThumbprint: "",
		msiEndpoint: "", msiSecret: "", arcIMDSEndpoint: "",
		"AZURE_POD_IDENTITY_AUTHORITY_HOST": "", "MSAL_MI_DISABLE_IMDS_V2": "",
	})
	require.NoError(t, os.Unsetenv(azureClientID))
	before := shellExec
	defer func() { shellExec = before }()
	shellExec = func(context.Context, string, string) ([]byte, error) {
		return nil, NewCredentialUnavailableError("CLI credentials are disabled for this test")
	}

	t.Run("probe", func(t *testing.T) {
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		// IMDSv1-only discovery avoids all native key operations, including on Windows.
		srv.AppendResponse(mock.WithStatusCode(http.StatusNotFound))
		srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
		srv.AppendResponse(mock.WithBody([]byte(`{"securityProfile":{"securityType":"Standard"}}`)))
		srv.RepeatResponse(2, mock.WithBody(accessTokenRespSuccess))
		var paths []string
		cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry: policy.RetryOptions{
					MaxRetries:  5,
					StatusCodes: []int{http.StatusInternalServerError},
				},
				Transport: srv,
				PerCallPolicies: []policy.Policy{managedIdentityTestPolicy{do: func(req *policy.Request) {
					raw := req.Raw()
					paths = append(paths, raw.URL.Path)
					if len(paths) <= 2 {
						assert.Empty(t, raw.Header.Get(headerMetadata))
						assert.NotEmpty(t, raw.Header.Get("x-ms-client-request-id"), "MSAL must construct the probe")
						deadline, ok := raw.Context().Deadline()
						assert.True(t, ok)
						assert.InDelta(t, time.Second.Seconds(), time.Until(deadline).Seconds(), 0.1)
					}
				}}},
			},
		})
		require.NoError(t, err)
		// MSAL's process-wide token cache survives go test -count repetitions.
		scope := t.Name() + uuid.NewString()
		tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
		require.NoError(t, err)
		require.Equal(t, tokenValue, tk.Token)
		_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope + "-second"}})
		require.NoError(t, err)
		assert.Equal(t, []string{
			"/metadata/identity/getplatformmetadata", "/metadata/identity/oauth2/token",
			"/metadata/instance/compute", "/metadata/identity/oauth2/token", "/metadata/identity/oauth2/token",
		}, paths, "MSAL probes v2 then v1 once; acquisitions must not re-probe")
	})

	t.Run("timeout", func(t *testing.T) {
		requests := 0
		stall := true
		chain, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{
				PerCallPolicies: []policy.Policy{policyFunc(func(req *policy.Request) (*http.Response, error) {
					requests++
					if stall {
						<-req.Raw().Context().Done()
						return nil, req.Raw().Context().Err()
					}
					return req.Next()
				})},
				Retry:     policy.RetryOptions{MaxRetries: 5},
				Transport: &mockSTS{},
			},
		})
		require.NoError(t, err)
		for i := 0; i < 2; i++ {
			start := time.Now()
			_, err = chain.GetToken(context.Background(), testTRO)
			require.ErrorContains(t, err, credNameManagedIdentity+": managed identity timed out")
			assert.InDelta(t, 1, time.Since(start).Seconds(), 0.5, "the real MSAL discovery must stop after one second")
		}
		assert.Equal(t, 2, requests, "timed-out probes must not retry or continue to v1")
		stall = false
		tk, err := chain.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.NoError(t, err)
		require.Equal(t, tokenValue, tk.Token)
	})

	t.Run("MSAL delegation and caller context", func(t *testing.T) {
		for _, timeout := range []time.Duration{5 * time.Second, 100 * time.Millisecond} {
			t.Run(timeout.String(), func(t *testing.T) {
				c, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{dac: true})
				require.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				parentDeadline, _ := ctx.Deadline()
				probes, acquisitions := 0, 0
				c.mic.msalClient = fakeManagedIdentityClient{
					capabilities: func(probe context.Context) (managedidentity.Capabilities, error) {
						probes++
						deadline, ok := probe.Deadline()
						require.True(t, ok)
						if timeout < time.Second {
							assert.Equal(t, parentDeadline, deadline)
						} else {
							assert.InDelta(t, 1, time.Until(deadline).Seconds(), 0.1)
						}
						return managedidentity.Capabilities{Source: managedidentity.DefaultToIMDS}, nil
					},
					acquireToken: func(actual context.Context, _ string, _ ...managedidentity.AcquireTokenOption) (managedidentity.AuthResult, error) {
						acquisitions++
						assert.Same(t, ctx, actual)
						assert.NoError(t, actual.Err(), "cancelling discovery must not cancel acquisition")
						return managedidentity.AuthResult{AccessToken: tokenValue}, nil
					},
				}
				for range 2 {
					_, err = c.mic.GetToken(ctx, testTRO)
					require.NoError(t, err)
				}
				assert.Equal(t, 1, probes)
				assert.Equal(t, 2, acquisitions)
			})
		}
	})

	t.Run("no detected source", func(t *testing.T) {
		c, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{dac: true})
		require.NoError(t, err)
		c.mic.msalClient = fakeManagedIdentityClient{
			capabilities: func(context.Context) (managedidentity.Capabilities, error) {
				return managedidentity.Capabilities{ErrorReason: "offline discovery failure"}, nil
			},
			acquireToken: func(context.Context, string, ...managedidentity.AcquireTokenOption) (managedidentity.AuthResult, error) {
				t.Fatal("acquisition must not run when MSAL detected no source")
				return managedidentity.AuthResult{}, nil
			},
		}
		_, err = c.GetToken(context.Background(), testTRO)
		var unavailable credentialUnavailable
		require.ErrorAs(t, err, &unavailable)
		assert.ErrorContains(t, err, "offline discovery failure")
		assert.True(t, c.mic.probeIMDS, "unsuccessful discovery must not disable later probes")
	})

	t.Run("unrelated responses aren't IMDS", func(t *testing.T) {
		srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
		defer close()
		srv.SetResponse(mock.WithStatusCode(http.StatusServiceUnavailable))
		c, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: azcore.ClientOptions{Transport: srv, Retry: policy.RetryOptions{MaxRetries: 5}},
		})
		require.NoError(t, err)
		_, err = c.GetToken(context.Background(), testTRO)
		require.ErrorContains(t, err, "no managed identity endpoint detected")
		assert.Equal(t, 2, srv.Requests(), "only MSAL's v2 and v1 probes should run, without retries")
	})
}

func TestDefaultAzureCredential_Attestation(t *testing.T) {
	t.Setenv(azureTokenCredentials, credNameManagedIdentity)
	for _, registered := range []bool{false, true} {
		c, err := NewDefaultAzureCredential(nil)
		require.NoError(t, err)
		require.Len(t, c.chain.sources, 1)
		mic, ok := c.chain.sources[0].(*ManagedIdentityCredential)
		require.True(t, ok)
		assert.False(t, mic.mic.attestationAvailable())
		assert.False(t, mic.mic.disableMTLSPoP, "DAC must permit automatic PoP by default")
		if registered {
			mic.mic.attestationAvailable = func() bool { return true }
		}
		binding := &tls.Certificate{}
		var want managedidentity.AcquireTokenOptions
		managedidentity.WithClaims("claims")(&want)
		if registered {
			managedidentity.WithMtlsProofOfPossession()(&want)
			managedidentity.WithAttestationSupport()(&want)
			managedidentity.WithMtlsPoPMinStrength(managedidentity.MtlsBindingStrengthKeyGuard)(&want)
		}
		mic.mic.msalClient = fakeManagedIdentityClient{
			capabilities: func(context.Context) (managedidentity.Capabilities, error) {
				require.True(t, registered)
				return managedidentity.Capabilities{
					Source: managedidentity.DefaultToIMDS, MaxSupportedBindingStrength: managedidentity.MtlsBindingStrengthKeyGuard,
				}, nil
			},
			acquireToken: func(_ context.Context, resource string, opts ...managedidentity.AcquireTokenOption) (managedidentity.AuthResult, error) {
				assert.Equal(t, "https://vault.azure.net", resource)
				var got managedidentity.AcquireTokenOptions
				for _, opt := range opts {
					opt(&got)
				}
				assert.Equal(t, want, got, "DAC must propagate PoP, attestation, and claims to MSAL")
				return managedidentity.AuthResult{AccessToken: tokenValue, BindingCertificate: binding}, nil
			},
		}
		token, err := c.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{"https://vault.azure.net/.default"}, Claims: "claims", EnableMTLSPoP: true,
		})
		require.NoError(t, err)
		assert.Same(t, binding, token.BindingCertificate, "DAC must retain MSAL's binding certificate")
	}
}

func TestDefaultAzureCredential_AttestedProbe(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{
		azureTokenCredentials: "prod", azureTenantID: "", azureClientID: "",
		azureClientSecret: "", azureFederatedTokenFile: "",
		identityEndpoint: "", identityHeader: "", identityServerThumbprint: "",
		msiEndpoint: "", msiSecret: "", arcIMDSEndpoint: "",
	})
	require.NoError(t, os.Unsetenv(azureClientID))
	c, err := NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	var mic *ManagedIdentityCredential
	for _, source := range c.chain.sources {
		if candidate, ok := source.(*ManagedIdentityCredential); ok {
			mic = candidate
		}
	}
	require.NotNil(t, mic)
	mic.mic.attestationAvailable = func() bool { return true }
	parent, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	probes, acquisitions := 0, 0
	mic.mic.msalClient = fakeManagedIdentityClient{
		capabilities: func(probe context.Context) (managedidentity.Capabilities, error) {
			probes++
			deadline, ok := probe.Deadline()
			require.True(t, ok)
			assert.InDelta(t, 1, time.Until(deadline).Seconds(), 0.1)
			return managedidentity.Capabilities{
				Source: managedidentity.DefaultToIMDS, MaxSupportedBindingStrength: managedidentity.MtlsBindingStrengthKeyGuard,
			}, nil
		},
		acquireToken: func(actual context.Context, _ string, opts ...managedidentity.AcquireTokenOption) (managedidentity.AuthResult, error) {
			acquisitions++
			assert.NoError(t, actual.Err(), "canceling the discovery context must not cancel acquisition")
			deadline, ok := actual.Deadline()
			require.True(t, ok)
			parentDeadline, _ := parent.Deadline()
			assert.Equal(t, parentDeadline, deadline)
			var got, want managedidentity.AcquireTokenOptions
			for _, opt := range opts {
				opt(&got)
			}
			managedidentity.WithClaims("")(&want)
			managedidentity.WithMtlsProofOfPossession()(&want)
			managedidentity.WithAttestationSupport()(&want)
			managedidentity.WithMtlsPoPMinStrength(managedidentity.MtlsBindingStrengthKeyGuard)(&want)
			assert.Equal(t, want, got)
			return managedidentity.AuthResult{AccessToken: tokenValue, BindingCertificate: &tls.Certificate{}}, nil
		},
	}
	tk, err := c.GetToken(parent, policy.TokenRequestOptions{Scopes: []string{t.Name()}, EnableMTLSPoP: true})
	require.NoError(t, err)
	assert.NotNil(t, tk.BindingCertificate)
	assert.Equal(t, 1, probes, "the initial bounded probe must also supply PoP eligibility")
	assert.Equal(t, 1, acquisitions)
	assert.False(t, mic.mic.probeIMDS)
}

func TestDefaultAzureCredential_UnexpectedIMDSResponse(t *testing.T) {
	before := shellExec
	defer func() { shellExec = before }()
	shellExec = mockAzSuccess

	const dockerDesktopPrefix = "connecting to 169.254.169.254:80: connecting to 169.254.169.254:80: dial tcp 169.254.169.254:80: connectex: A socket operation was attempted to an unreachable "
	for _, test := range []struct {
		desc string
		res  [][]mock.ResponseOption
	}{
		{
			"Docker Desktop",
			[][]mock.ResponseOption{
				{
					mock.WithBody([]byte(dockerDesktopPrefix + "host.")),
					mock.WithStatusCode(http.StatusForbidden),
				},
				{
					mock.WithBody([]byte(dockerDesktopPrefix + "host.")),
					mock.WithStatusCode(http.StatusForbidden),
				},
			},
		},
		{
			"Docker Desktop",
			[][]mock.ResponseOption{
				{
					mock.WithBody([]byte(dockerDesktopPrefix + "network.")),
					mock.WithStatusCode(http.StatusForbidden),
				},
				{
					mock.WithBody([]byte(dockerDesktopPrefix + "network.")),
					mock.WithStatusCode(http.StatusForbidden),
				},
			},
		},
		{
			"IMDS: no identity assigned",
			[][]mock.ResponseOption{
				{mock.WithStatusCode(http.StatusBadRequest)},
				{
					mock.WithBody([]byte(`{"error":"invalid_request","error_description":"Identity not found"}`)),
					mock.WithStatusCode(http.StatusBadRequest),
				},
			},
		},
		{
			"no token in response",
			[][]mock.ResponseOption{
				{mock.WithStatusCode(http.StatusOK)},
				{mock.WithBody([]byte(`{"error": "no token here"}`)), mock.WithStatusCode(http.StatusOK)},
			},
		},
		{
			"non-JSON token response",
			[][]mock.ResponseOption{
				{mock.WithStatusCode(http.StatusOK)},
				{mock.WithBody([]byte("not json")), mock.WithStatusCode(http.StatusOK)},
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			for _, res := range test.res {
				srv.AppendResponse(res...)
			}
			c, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
				ClientOptions: policy.ClientOptions{Transport: srv},
			})
			require.NoError(t, err)
			tk, err := c.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{strings.ReplaceAll(t.Name(), "#", "")}})
			require.NoError(t, err, "expected a token from AzureCLICredential")
			require.Equal(t, tokenValue, tk.Token, "expected a token from AzureCLICredential")
		})
	}
}

func TestDefaultAzureCredential_UnsupportedMIClientID(t *testing.T) {
	fail := true
	before := shellExec
	defer func() { shellExec = before }()
	shellExec = func(ctx context.Context, credName string, commandLine string) ([]byte, error) {
		if fail {
			return nil, errors.New("fail")
		}
		return mockAzSuccess(ctx, credName, commandLine)
	}
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(msiEndpoint, fakeMIEndpoint)

	cred, err := NewDefaultAzureCredential(nil)
	require.NoError(t, err, "an unsupported client ID isn't a constructor error")

	_, err = cred.GetToken(ctx, testTRO)
	require.ErrorContains(t, err, "Cloud Shell", "error should mention the unsupported ID")

	fail = false
	_, err = cred.GetToken(ctx, testTRO)
	require.NoError(t, err, "expected a token from AzureCLICredential")
}
