//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/golang-jwt/jwt/v4"
)

// constants used throughout this package
const (
	accessTokenRespMalformed = `{"access_token": 0, "expires_in": 3600}`
	badTenantID              = "bad_tenant"
	tokenExpiresIn           = 3600
	tokenValue               = "new_token"
)

var (
	accessTokenRespSuccess    = []byte(fmt.Sprintf(`{"access_token": "%s", "expires_in": %d}`, tokenValue, tokenExpiresIn))
	instanceDiscoveryResponse = getInstanceDiscoveryResponse(fakeTenantID)
	tenantDiscoveryResponse   = getTenantDiscoveryResponse(fakeTenantID)
)

// constants for this file
const (
	testHost = "https://localhost"
)

func getInstanceDiscoveryResponse(tenant string) []byte {
	return []byte(strings.ReplaceAll(`{
		"tenant_discovery_endpoint": "https://login.microsoftonline.com/{tenant}/v2.0/.well-known/openid-configuration",
		"api-version": "1.1",
		"metadata": [
			{
				"preferred_network": "login.microsoftonline.com",
				"preferred_cache": "login.windows.net",
				"aliases": [
					"login.microsoftonline.com",
					"login.windows.net",
					"login.microsoft.com",
					"sts.windows.net"
				]
			}
		]
	}`, "{tenant}", tenant))
}

func getTenantDiscoveryResponse(tenant string) []byte {
	return []byte(strings.ReplaceAll(`{
		"token_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token",
		"token_endpoint_auth_methods_supported": [
			"client_secret_post",
			"private_key_jwt",
			"client_secret_basic"
		],
		"jwks_uri": "https://login.microsoftonline.com/{tenant}/discovery/v2.0/keys",
		"response_modes_supported": [
			"query",
			"fragment",
			"form_post"
		],
		"subject_types_supported": [
			"pairwise"
		],
		"id_token_signing_alg_values_supported": [
			"RS256"
		],
		"response_types_supported": [
			"code",
			"id_token",
			"code id_token",
			"id_token token"
		],
		"scopes_supported": [
			"openid",
			"profile",
			"email",
			"offline_access"
		],
		"issuer": "https://login.microsoftonline.com/{tenant}/v2.0",
		"request_uri_parameter_supported": false,
		"userinfo_endpoint": "https://graph.microsoft.com/oidc/userinfo",
		"authorization_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/authorize",
		"device_authorization_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/devicecode",
		"http_logout_supported": true,
		"frontchannel_logout_supported": true,
		"end_session_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/logout",
		"claims_supported": [
			"sub",
			"iss",
			"cloud_instance_name",
			"cloud_instance_host_name",
			"cloud_graph_host_name",
			"msgraph_host",
			"aud",
			"exp",
			"iat",
			"auth_time",
			"acr",
			"nonce",
			"preferred_username",
			"name",
			"tid",
			"ver",
			"at_hash",
			"c_hash",
			"email"
		],
		"kerberos_endpoint": "https://login.microsoftonline.com/{tenant}/kerberos",
		"tenant_region_scope": "NA",
		"cloud_instance_name": "microsoftonline.com",
		"cloud_graph_host_name": "graph.windows.net",
		"msgraph_host": "graph.microsoft.com",
		"rbac_url": "https://pas.windows.net"
	}`, "{tenant}", tenant))
}

func validateX5C(t *testing.T, certs []*x509.Certificate) mock.ResponsePredicate {
	return func(req *http.Request) bool {
		err := req.ParseForm()
		if err != nil {
			t.Fatal("expected a form body")
		}
		assertion, ok := req.PostForm["client_assertion"]
		if !ok {
			t.Fatal("expected a client_assertion field")
		}
		if len(assertion) != 1 {
			t.Fatalf(`unexpected client_assertion "%v"`, assertion)
		}
		token, _ := jwt.Parse(assertion[0], nil)
		if token == nil {
			t.Fatalf("failed to parse the assertion: %s", assertion)
		}
		if v, ok := token.Header["x5c"].([]any); !ok {
			t.Fatal("missing x5c header")
		} else if actual := len(v); actual != len(certs) {
			t.Fatalf("expected %d certs, got %d", len(certs), actual)
		}
		return true
	}
}

// Set environment variables for the duration of a test. Restore their prior values
// after the test completes. Obviated by 1.17's T.Setenv
func setEnvironmentVariables(t *testing.T, vars map[string]string) {
	unsetSentinel := "variables having no initial value must be unset after the test"
	priorValues := make(map[string]string, len(vars))
	for k, v := range vars {
		priorValue, ok := os.LookupEnv(k)
		if ok {
			priorValues[k] = priorValue
		} else {
			priorValues[k] = unsetSentinel
		}
		err := os.Setenv(k, v)
		if err != nil {
			t.Fatalf("Unexpected error setting %s: %v", k, err)
		}
	}

	t.Cleanup(func() {
		for k, v := range priorValues {
			var err error
			if v == unsetSentinel {
				err = os.Unsetenv(k)
			} else {
				err = os.Setenv(k, v)
			}
			if err != nil {
				t.Fatalf("Unexpected error resetting %s: %v", k, err)
			}
		}
	})
}

func Test_WellKnownHosts(t *testing.T) {
	for _, cloud := range []cloud.Configuration{cloud.AzureChina, cloud.AzureGovernment, cloud.AzurePublic} {
		host, err := setAuthorityHost(cloud)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(host, "https://login.") {
			t.Fatal("unexpected ActiveDirectoryAuthorityHost: " + host)
		}
	}
}

func Test_SetEnvAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: testHost})
	authorityHost, err := setAuthorityHost(cloud.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != testHost {
		t.Fatalf(`unexpected host "%s"`, authorityHost)
	}
}

func Test_CustomAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: testHost + "/not"})
	authorityHost, err := setAuthorityHost(cloud.Configuration{ActiveDirectoryAuthorityHost: testHost})
	if err != nil {
		t.Fatal(err)
	}
	// ensure env var doesn't override explicit value
	if authorityHost != testHost {
		t.Fatalf(`expected "%s", got "%s"`, testHost, authorityHost)
	}
}

func Test_DefaultAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: ""})
	authorityHost, err := setAuthorityHost(cloud.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != cloud.AzurePublic.ActiveDirectoryAuthorityHost {
		t.Fatal("unexpected default host: " + authorityHost)
	}
}

func Test_GetTokenRequiresScopes(t *testing.T) {
	for _, ctor := range []func() (azcore.TokenCredential, error){
		func() (azcore.TokenCredential, error) { return NewAzureCLICredential(nil) },
		func() (azcore.TokenCredential, error) {
			return NewClientCertificateCredential("tenantID", "clientID", allCertTests[0].certs, allCertTests[0].key, nil)
		},
		func() (azcore.TokenCredential, error) {
			return NewClientSecretCredential("tenantID", "clientID", fakeSecret, nil)
		},
		func() (azcore.TokenCredential, error) { return NewDeviceCodeCredential(nil) },
		func() (azcore.TokenCredential, error) { return NewInteractiveBrowserCredential(nil) },
		func() (azcore.TokenCredential, error) {
			return NewUsernamePasswordCredential("tenantID", "clientID", "username", "password", nil)
		},
	} {
		cred, err := ctor()
		t.Run(fmt.Sprintf("%T", cred), func(t *testing.T) {
			if err != nil {
				t.Fatal(err)
			}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{})
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func Test_NonHTTPSAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: ""})
	authorityHost, err := setAuthorityHost(cloud.Configuration{ActiveDirectoryAuthorityHost: "http://localhost"})
	if err == nil {
		t.Fatal("Expected an error but did not receive one.")
	}
	if authorityHost != "" {
		t.Fatalf("Unexpected value in authority host string: %s", authorityHost)
	}
}

func TestAdditionallyAllowedTenants(t *testing.T) {
	af := filepath.Join(t.TempDir(), t.Name()+credNameWorkloadIdentity)
	if err := os.WriteFile(af, []byte("assertion"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	tenantA := "A"
	tenantB := "B"
	for _, test := range []struct {
		allowed                []string
		desc, expected, tenant string
		err                    bool
	}{
		{
			desc:     "all tenants allowed",
			allowed:  []string{"*"},
			expected: tenantA,
			tenant:   tenantA,
		},
		{
			desc:     "tenant explicitly allowed",
			allowed:  []string{tenantA, tenantB},
			expected: tenantA,
			tenant:   tenantA,
		},
		{
			desc:     "tenant explicitly allowed",
			allowed:  []string{tenantA, tenantB},
			expected: tenantB,
			tenant:   tenantB,
		},
		{
			desc:    "tenant not allowed",
			allowed: []string{tenantA},
			tenant:  tenantB,
			err:     true,
		},
		{
			desc:   "no additional tenants allowed",
			tenant: tenantA,
			err:    true,
		},
	} {
		tro := policy.TokenRequestOptions{Scopes: []string{liveTestScope}, TenantID: test.tenant}
		for _, subtest := range []struct {
			ctor func(azcore.ClientOptions) (azcore.TokenCredential, error)
			env  map[string]string
			name string
		}{
			{
				name: credNameAssertion,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := ClientAssertionCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co}
					return NewClientAssertionCredential(fakeTenantID, fakeClientID, func(context.Context) (string, error) { return "...", nil }, &o)
				},
			},
			{
				name: credNameAzureCLI,
				ctor: func(azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := AzureCLICredentialOptions{
						AdditionallyAllowedTenants: test.allowed,
						tokenProvider: func(ctx context.Context, resource, tenantID string) ([]byte, error) {
							if tenantID != test.expected {
								t.Errorf(`unexpected tenantID "%s"`, tenantID)
							}
							return mockCLITokenProviderSuccess(ctx, resource, tenantID)
						},
					}
					return NewAzureCLICredential(&o)
				},
			},
			{
				name: credNameCert,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := ClientCertificateCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co}
					return NewClientCertificateCredential(fakeTenantID, fakeClientID, allCertTests[0].certs, allCertTests[0].key, &o)
				},
			},
			{
				name: credNameDeviceCode,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := DeviceCodeCredentialOptions{
						AdditionallyAllowedTenants: test.allowed,
						ClientOptions:              co,
						UserPrompt:                 func(context.Context, DeviceCodeMessage) error { return nil },
					}
					return NewDeviceCodeCredential(&o)
				},
			},
			{
				name: credNameOBO,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := OnBehalfOfCredentialOptions{
						AdditionallyAllowedTenants: test.allowed,
						ClientOptions:              co,
					}
					return NewOnBehalfOfCredentialWithSecret(fakeTenantID, fakeClientID, "assertion", fakeSecret, &o)
				},
			},
			{
				name: credNameSecret,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := ClientSecretCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co}
					return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, &o)
				},
			},
			{
				name: credNameUserPassword,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := UsernamePasswordCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co}
					return NewUsernamePasswordCredential(fakeTenantID, fakeClientID, fakeUsername, "password", &o)
				},
			},
			{
				name: credNameWorkloadIdentity,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					return NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
						AdditionallyAllowedTenants: test.allowed,
						ClientID:                   fakeClientID,
						ClientOptions:              co,
						TenantID:                   fakeTenantID,
						TokenFilePath:              af,
					})
				},
			},
			{
				name: "DefaultAzureCredential/EnvironmentCredential",
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := DefaultAzureCredentialOptions{ClientOptions: co, TenantID: test.tenant}
					return NewDefaultAzureCredential(&o)
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(test.allowed, ";"),
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   fakeTenantID,
				},
			},
			{
				name: "DefaultAzureCredential/EnvironmentCredential/option-overrides-env",
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := DefaultAzureCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co, TenantID: test.tenant}
					return NewDefaultAzureCredential(&o)
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: "not-" + test.tenant,
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   fakeTenantID,
				},
			},
			{
				name: "DefaultAzureCredential/" + credNameWorkloadIdentity,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := DefaultAzureCredentialOptions{AdditionallyAllowedTenants: test.allowed, ClientOptions: co}
					return NewDefaultAzureCredential(&o)
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(test.allowed, ";"),
					azureAuthorityHost:              "https://login.microsoftonline.com",
					azureClientID:                   fakeClientID,
					azureFederatedTokenFile:         af,
					azureTenantID:                   fakeTenantID,
				},
			},
			{
				name: "EnvironmentCredential/" + credNameCert,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(test.allowed, ";"),
					azureClientCertificatePath:      "testdata/certificate.pem",
					azureClientID:                   fakeClientID,
					azureTenantID:                   fakeTenantID,
				},
			},
			{
				name: "EnvironmentCredential/" + credNameSecret,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(test.allowed, ";"),
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   fakeTenantID,
				},
			},
			{
				name: "EnvironmentCredential/" + credNameUserPassword,
				ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
					return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
				},
				env: map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(test.allowed, ";"),
					azureClientID:                   fakeClientID,
					azurePassword:                   "password",
					azureTenantID:                   fakeTenantID,
					azureUsername:                   fakeUsername,
				},
			},
		} {
			t.Run(fmt.Sprintf("%s/%s", subtest.name, test.desc), func(t *testing.T) {
				for k, v := range subtest.env {
					t.Setenv(k, v)
				}
				sts := mockSTS{
					tenant: test.tenant,
					tokenRequestCallback: func(r *http.Request) {
						if actual := strings.Split(r.URL.Path, "/")[1]; actual != test.expected {
							t.Fatalf("expected tenant %q, got %q", test.expected, actual)
						}
					},
				}
				c, err := subtest.ctor(policy.ClientOptions{Transport: &sts})
				if err != nil {
					t.Fatal(err)
				}
				tk, err := c.GetToken(context.Background(), tro)
				if err != nil {
					if test.err {
						return
					}
					t.Fatal(err)
				} else if test.err {
					t.Fatal("expected an error")
				}
				// silent authentication should succeed
				tk2, err := c.GetToken(context.Background(), tro)
				if err != nil {
					t.Fatalf(`silent authentication failed: "%v"`, err)
				}
				if tk.Token != tk2.Token {
					t.Fatalf("expected %q, got %q", tk.Token, tk2.Token)
				}
				if !tk.ExpiresOn.Equal(tk2.ExpiresOn) {
					t.Fatalf("expected %v, got %v", tk.ExpiresOn, tk2.ExpiresOn)
				}
			})
		}
		t.Run(fmt.Sprintf("DefaultAzureCredential/%s/%s", credNameAzureCLI, test.desc), func(t *testing.T) {
			// mock IMDS failure because managed identity precedes CLI in the chain
			srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.SetResponse(mock.WithStatusCode(400))
			o := DefaultAzureCredentialOptions{
				AdditionallyAllowedTenants: test.allowed,
				ClientOptions:              policy.ClientOptions{Transport: srv},
			}
			c, err := NewDefaultAzureCredential(&o)
			if err != nil {
				t.Fatal(err)
			}
			called := false
			for _, source := range c.chain.sources {
				if cli, ok := source.(*AzureCLICredential); ok {
					cli.tokenProvider = func(ctx context.Context, resource, tenantID string) ([]byte, error) {
						called = true
						if tenantID != test.expected {
							t.Fatalf(`unexpected tenantID "%s"`, tenantID)
						}
						return mockCLITokenProviderSuccess(ctx, resource, tenantID)
					}
					break
				}
			}
			if _, err := c.GetToken(context.Background(), tro); err != nil {
				if test.err {
					return
				}
				t.Fatal(err)
			} else if test.err {
				t.Fatal("expected an error")
			}
			if !called {
				t.Fatal("AzureCLICredential wasn't invoked")
			}
		})
	}
}

func TestClaims(t *testing.T) {
	t.Skip("unskip this test after adding back CAE support")
	realCP1 := disableCP1
	t.Cleanup(func() { disableCP1 = realCP1 })
	claim := `"test":"pass"`
	for _, test := range []struct {
		ctor func(azcore.ClientOptions) (azcore.TokenCredential, error)
		name string
	}{
		{
			name: credNameAssertion,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := ClientAssertionCredentialOptions{ClientOptions: co}
				return NewClientAssertionCredential(fakeTenantID, fakeClientID, func(context.Context) (string, error) { return "...", nil }, &o)
			},
		},
		{
			name: credNameCert,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := ClientCertificateCredentialOptions{ClientOptions: co}
				return NewClientCertificateCredential(fakeTenantID, fakeClientID, allCertTests[0].certs, allCertTests[0].key, &o)
			},
		},
		{
			name: credNameDeviceCode,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := DeviceCodeCredentialOptions{
					ClientOptions: co,
					UserPrompt:    func(context.Context, DeviceCodeMessage) error { return nil },
				}
				return NewDeviceCodeCredential(&o)
			},
		},
		{
			name: credNameOBO,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := OnBehalfOfCredentialOptions{ClientOptions: co}
				return NewOnBehalfOfCredentialWithSecret(fakeTenantID, fakeClientID, "assertion", fakeSecret, &o)
			},
		},
		{
			name: credNameSecret,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := ClientSecretCredentialOptions{ClientOptions: co}
				return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, &o)
			},
		},
		{
			name: credNameUserPassword,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := UsernamePasswordCredentialOptions{ClientOptions: co}
				return NewUsernamePasswordCredential(fakeTenantID, fakeClientID, fakeUsername, "password", &o)
			},
		},
	} {
		for _, d := range []bool{true, false} {
			name := test.name
			if d {
				name += " disableCP1"
			}
			t.Run(name, func(t *testing.T) {
				disableCP1 = d
				reqs := 0
				sts := mockSTS{
					tokenRequestCallback: func(r *http.Request) {
						if err := r.ParseForm(); err != nil {
							t.Error(err)
						}
						reqs++
						// If the disableCP1 flag isn't set, both requests should specify CP1. The second
						// GetToken call specifies claims we should find in the following token request.
						// We check only for substrings because MSAL is responsible for formatting claims.
						actual := fmt.Sprint(r.Form["claims"])
						if strings.Contains(actual, "CP1") == disableCP1 {
							t.Fatalf(`unexpected claims "%v"`, actual)
						}
						if reqs == 2 {
							if !strings.Contains(strings.ReplaceAll(actual, " ", ""), claim) {
								t.Fatalf(`unexpected claims "%v"`, actual)
							}
						}
					},
				}
				o := azcore.ClientOptions{Transport: &sts}
				cred, err := test.ctor(o)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"A"}}); err != nil {
					t.Fatal(err)
				}
				// TODO: uncomment after restoring TokenRequestOptions.Claims
				// if _, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Claims: fmt.Sprintf("{%s}", claim), Scopes: []string{"B"}}); err != nil {
				// 	t.Fatal(err)
				// }
				if reqs != 2 {
					t.Fatalf("expected %d token requests, got %d", 2, reqs)
				}
			})
		}
	}
}

// ==================================================================================================================================

type fakeConfidentialClient struct {
	// set ar to have all API calls return the provided AuthResult
	ar confidential.AuthResult

	// set err to have all API calls return the provided error
	err error

	// set true to have silent auth succeed
	silentAuth bool

	// optional callbacks for validating MSAL call args
	oboCallback func(context.Context, string, []string)
}

func (f fakeConfidentialClient) returnResult() (confidential.AuthResult, error) {
	if f.err != nil {
		return confidential.AuthResult{}, f.err
	}
	return f.ar, nil
}

func (f fakeConfidentialClient) AcquireTokenSilent(ctx context.Context, scopes []string, options ...confidential.AcquireSilentOption) (confidential.AuthResult, error) {
	if f.silentAuth {
		return f.ar, nil
	}
	return confidential.AuthResult{}, errors.New("silent authentication failed")
}

func (f fakeConfidentialClient) AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, options ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error) {
	return f.returnResult()
}

func (f fakeConfidentialClient) AcquireTokenByCredential(ctx context.Context, scopes []string, options ...confidential.AcquireByCredentialOption) (confidential.AuthResult, error) {
	return f.returnResult()
}

func (f fakeConfidentialClient) AcquireTokenOnBehalfOf(ctx context.Context, userAssertion string, scopes []string, options ...confidential.AcquireOnBehalfOfOption) (confidential.AuthResult, error) {
	if f.oboCallback != nil {
		f.oboCallback(ctx, userAssertion, scopes)
	}
	return f.returnResult()
}

var _ confidentialClient = (*fakeConfidentialClient)(nil)

// ==================================================================================================================================

type fakePublicClient struct {
	// set ar to have all API calls return the provided AuthResult
	ar public.AuthResult

	// similar to ar but for device code APIs
	dc public.DeviceCode

	// set err to have all API calls return the provided error
	err error

	// set true to have silent auth succeed
	silentAuth bool
}

func (f fakePublicClient) returnResult() (public.AuthResult, error) {
	if f.err != nil {
		return public.AuthResult{}, f.err
	}
	return f.ar, nil
}

func (f fakePublicClient) AcquireTokenSilent(ctx context.Context, scopes []string, options ...public.AcquireSilentOption) (public.AuthResult, error) {
	if f.silentAuth {
		return f.ar, nil
	}
	return public.AuthResult{}, errors.New("silent authentication failed")
}

func (f fakePublicClient) AcquireTokenByUsernamePassword(ctx context.Context, scopes []string, username string, password string, options ...public.AcquireByUsernamePasswordOption) (public.AuthResult, error) {
	return f.returnResult()
}

func (f fakePublicClient) AcquireTokenByDeviceCode(ctx context.Context, scopes []string, options ...public.AcquireByDeviceCodeOption) (public.DeviceCode, error) {
	if f.err != nil {
		return public.DeviceCode{}, f.err
	}
	return f.dc, nil
}

func (f fakePublicClient) AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, options ...public.AcquireByAuthCodeOption) (public.AuthResult, error) {
	return f.returnResult()
}

func (f fakePublicClient) AcquireTokenInteractive(ctx context.Context, scopes []string, options ...public.AcquireInteractiveOption) (public.AuthResult, error) {
	return f.returnResult()
}

var _ publicClient = (*fakePublicClient)(nil)
