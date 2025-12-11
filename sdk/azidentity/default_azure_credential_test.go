// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
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
		t.Setenv(customtokenproxy.EnvAzureKubernetesCAData, "not a valid cert")
		t.Setenv(customtokenproxy.EnvAzureKubernetesTokenProxy, "http://timeout.local&fail=yes#please")

		cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
		})
		require.NoError(t, err)

		// ...but ensure a timeout should it try the proxy anyway
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
		defer cancel()
		_, err = cred.GetToken(ctx, testTRO)
		require.NoError(t, err)
	})
}

// delayPolicy adds a delay to pipeline requests. Used to test timeout behavior.
type delayPolicy struct {
	delay time.Duration
}

func (p *delayPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if p.delay > 0 {
		select {
		case <-req.Raw().Context().Done():
			return nil, req.Raw().Context().Err()
		case <-time.After(p.delay):
			// delay has elapsed, continue on
		}
	}
	return req.Next()
}

func TestDefaultAzureCredential_IMDS(t *testing.T) {
	before := shellExec
	defer func() { shellExec = before }()
	shellExec = func(context.Context, string, string) ([]byte, error) {
		return nil, NewCredentialUnavailableError("CLI credentials are disabled for this test")
	}

	t.Run("probe", func(t *testing.T) {
		probed := false
		cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry: policy.RetryOptions{
					MaxRetries:  5,
					StatusCodes: []int{http.StatusInternalServerError},
				},
				Transport: &mockSTS{
					tokenRequestCallback: func(req *http.Request) *http.Response {
						hdr := req.Header.Get(headerMetadata)
						if probed {
							// This should be a token request. Return nil, mockSTS will respond with a token
							require.NotEmpty(t, hdr, "credential shouldn't retry probe request")
							return nil
						}
						// probe request. Respond with retriable status. The credential shouldn't retry
						probed = true
						require.Empty(t, hdr, "probe request shouldn't have Metadata header")
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
						}
					},
				},
			},
		})
		require.NoError(t, err)
		tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.NoError(t, err)
		require.True(t, probed)
		require.Equal(t, tokenValue, tk.Token)

		t.Run("Azure Container Instances", func(t *testing.T) {
			srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithBody([]byte("Required metadata header not specified or not correct")), mock.WithStatusCode(http.StatusBadRequest))
			srv.AppendResponse(mock.WithBody(accessTokenRespSuccess), mock.WithStatusCode(http.StatusOK))

			cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
				ClientOptions: policy.ClientOptions{
					Transport: srv,
				},
			})
			require.NoError(t, err)
			tk, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{t.Name()}})
			require.NoError(t, err, "DefaultAzureCredential should accept ACI's response to the probe request")
			require.Equal(t, tokenValue, tk.Token)
		})
	})

	t.Run("timeout", func(t *testing.T) {
		// shorten the timeout to speed up this test
		before := imdsProbeTimeout
		defer func() { imdsProbeTimeout = before }()
		imdsProbeTimeout = 100 * time.Millisecond

		dp := delayPolicy{2 * imdsProbeTimeout}
		chain, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
			ClientOptions: policy.ClientOptions{
				PerCallPolicies: []policy.Policy{&dp},
				Retry:           policy.RetryOptions{MaxRetries: -1},
				Transport:       &mockSTS{},
			},
		})
		require.NoError(t, err)
		for i := 0; i < 2; i++ {
			// expecting an error because managed identity times out and AzureCLICredential returns an error
			_, err = chain.GetToken(context.Background(), testTRO)
			require.ErrorContains(t, err, credNameManagedIdentity+": managed identity timed out")
		}

		// remove the delay so ManagedIdentityCredential can get a token from the fake STS
		dp.delay = 0
		tk, err := chain.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{t.Name()}})
		require.NoError(t, err)
		require.Equal(t, tokenValue, tk.Token)

		// now there should be no timeout on token requests
		dp.delay = 2 * imdsProbeTimeout
		tk, err = chain.GetToken(context.Background(), policy.TokenRequestOptions{
			// using a different scope forces a token request by bypassing the cache
			Scopes: []string{"not-" + testTRO.Scopes[0]},
		})
		require.NoError(t, err)
		require.Equal(t, tokenValue, tk.Token)
	})
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
