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
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// constants used throughout this package
const (
	accessTokenRespMalformed = `{"access_token": 0, "expires_in": 3600}`
	badTenantID              = "bad_tenant"
	tokenExpiresIn           = 3600
	tokenValue               = "new_token"
)

var (
	accessTokenRespSuccess = []byte(fmt.Sprintf(`{"access_token": "%s", "expires_in": %d}`, tokenValue, tokenExpiresIn))
	testTRO                = policy.TokenRequestOptions{Scopes: []string{liveTestScope}}
)

// constants for this file
const (
	testHost = "https://localhost"
)

func validateX5C(t *testing.T, certs []*x509.Certificate) func(*http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
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
		return nil
	}
}

// Set environment variables for the duration of a test. Restore their prior values
// after the test completes. uses t.Setenv on the key/value pairs in vars.
func setEnvironmentVariables(t *testing.T, vars map[string]string) {
	for k, v := range vars {
		t.Setenv(k, v)
	}
}

type tokenRequestCountingPolicy struct {
	count int
}

func (t *tokenRequestCountingPolicy) Do(req *policy.Request) (*http.Response, error) {
	if strings.HasSuffix(req.Raw().URL.Path, "/oauth2/v2.0/token") {
		t.count++
	}
	return req.Next()
}

func TestUserAuthentication(t *testing.T) {
	type authenticater interface {
		azcore.TokenCredential
		Authenticate(context.Context, *policy.TokenRequestOptions) (AuthenticationRecord, error)
	}
	for _, credential := range []struct {
		name                    string
		interactive, recordable bool
		new                     func(*TokenCachePersistenceOptions, azcore.ClientOptions, AuthenticationRecord, bool) (authenticater, error)
	}{
		{
			name: credNameBrowser,
			new: func(tcpo *TokenCachePersistenceOptions, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				return NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
					AuthenticationRecord:           ar,
					ClientOptions:                  co,
					DisableAutomaticAuthentication: disableAutoAuth,
					TokenCachePersistenceOptions:   tcpo,
				})
			},
			interactive: true,
		},
		{
			name: credNameDeviceCode,
			new: func(tcpo *TokenCachePersistenceOptions, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				o := DeviceCodeCredentialOptions{
					AuthenticationRecord:           ar,
					ClientOptions:                  co,
					DisableAutomaticAuthentication: disableAutoAuth,
					TokenCachePersistenceOptions:   tcpo,
				}
				if recording.GetRecordMode() == recording.PlaybackMode {
					o.UserPrompt = func(context.Context, DeviceCodeMessage) error { return nil }
				}
				return NewDeviceCodeCredential(&o)
			},
			interactive: true,
			recordable:  true,
		},
		{
			name: credNameUserPassword,
			new: func(tcpo *TokenCachePersistenceOptions, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				opts := UsernamePasswordCredentialOptions{
					AuthenticationRecord:         ar,
					ClientOptions:                co,
					TokenCachePersistenceOptions: tcpo,
				}
				return NewUsernamePasswordCredential(liveUser.tenantID, developerSignOnClientID, liveUser.username, liveUser.password, &opts)
			},
			recordable: true,
		},
	} {
		t.Run("AuthenticateDefaultScope/"+credential.name, func(t *testing.T) {
			if credential.name == credNameBrowser {
				t.Skip("the mock STS used in this test can't complete the interactive auth flow")
			}
			t.Setenv(azureAuthorityHost, "")
			customCloud := cloud.Configuration{
				ActiveDirectoryAuthorityHost: fmt.Sprintf("%s/%s", testHost, fakeTenantID),
				Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
					cloud.ResourceManager: {Audience: "https://localhost"},
				},
			}
			for _, cc := range []cloud.Configuration{cloud.AzureChina, cloud.AzureGovernment, cloud.AzurePublic, customCloud} {
				sts := mockSTS{tokenRequestCallback: func(r *http.Request) *http.Response {
					require.Contains(t, r.FormValue("scope"), cc.Services[cloud.ResourceManager].Audience+defaultSuffix)
					return nil
				}}

				co := azcore.ClientOptions{Cloud: cc, Transport: &sts}
				cred, err := credential.new(nil, co, AuthenticationRecord{}, false)
				require.NoError(t, err)
				_, err = cred.Authenticate(context.Background(), nil)
				require.NoError(t, err)

				os.Setenv(azureAuthorityHost, cc.ActiveDirectoryAuthorityHost)
				cred, err = credential.new(nil, azcore.ClientOptions{Transport: &sts}, AuthenticationRecord{}, false)
				require.NoError(t, err)
				_, err = cred.Authenticate(context.Background(), nil)
				if cc.ActiveDirectoryAuthorityHost == customCloud.ActiveDirectoryAuthorityHost {
					// Authenticate should return an error because it can't map an unknown host to a default scope
					require.ErrorIs(t, err, errScopeRequired)
				} else {
					require.NoError(t, err)
				}
			}
		})

		t.Run("Authenticate_Live_"+credential.name, func(t *testing.T) {
			switch recording.GetRecordMode() {
			case recording.LiveMode:
				if credential.interactive && !runManualTests {
					t.Skipf("set %s to run this test", azidentityRunManualTests)
				}
			case recording.PlaybackMode, recording.RecordingMode:
				if !credential.recordable {
					t.Skip("this test can't be recorded")
				}
			}
			co, stop := initRecording(t)
			defer stop()
			counter := tokenRequestCountingPolicy{}
			co.PerCallPolicies = append(co.PerCallPolicies, &counter)

			cred, err := credential.new(nil, co, AuthenticationRecord{}, false)
			require.NoError(t, err)
			ar, err := cred.Authenticate(context.Background(), &testTRO)
			require.NoError(t, err)

			// some fields of the returned AuthenticationRecord should have specific values
			require.Equal(t, ar.ClientID, developerSignOnClientID)
			require.Equal(t, ar.Version, supportedAuthRecordVersions[0])
			// all others should have nonempty values
			v := reflect.Indirect(reflect.ValueOf(&ar))
			for _, f := range reflect.VisibleFields(reflect.TypeOf(ar)) {
				s := v.FieldByIndex(f.Index).Addr().Interface().(*string)
				require.NotEmpty(t, *s)
			}
			require.Equal(t, 1, counter.count)
		})

		t.Run("PersistentCache_Live/"+credential.name, func(t *testing.T) {
			switch recording.GetRecordMode() {
			case recording.LiveMode:
				if credential.interactive && !runManualTests {
					t.Skipf("set %s to run this test", azidentityRunManualTests)
				}
			case recording.PlaybackMode, recording.RecordingMode:
				if !credential.recordable {
					t.Skip("this test can't be recorded")
				}
			}
			if runtime.GOOS != "windows" {
				t.Skip("this test runs only on Windows")
			}
			p, err := internal.CacheFilePath(t.Name())
			require.NoError(t, err)
			os.Remove(p)
			co, stop := initRecording(t)
			defer stop()
			counter := tokenRequestCountingPolicy{}
			co.PerCallPolicies = append(co.PerCallPolicies, &counter)
			tcpo := TokenCachePersistenceOptions{Name: t.Name()}

			cred, err := credential.new(&tcpo, co, AuthenticationRecord{}, true)
			require.NoError(t, err)
			record, err := cred.Authenticate(context.Background(), &testTRO)
			require.NoError(t, err)
			defer os.Remove(p)
			tk, err := cred.GetToken(context.Background(), testTRO)
			require.NoError(t, err)
			require.Equal(t, 1, counter.count)

			cred2, err := credential.new(&tcpo, co, record, true)
			require.NoError(t, err)
			tk2, err := cred2.GetToken(context.Background(), testTRO)
			require.NoError(t, err)
			require.Equal(t, tk.Token, tk2.Token)
		})

		if credential.interactive {
			t.Run("DisableAutomaticAuthentication/"+credential.name, func(t *testing.T) {
				cred, err := credential.new(nil, policy.ClientOptions{Transport: &mockSTS{}}, AuthenticationRecord{}, true)
				require.NoError(t, err)
				_, err = cred.GetToken(context.Background(), testTRO)
				require.ErrorIs(t, err, ErrAuthenticationRequired)
				if credential.name != credNameBrowser || runManualTests {
					_, err = cred.Authenticate(context.Background(), &testTRO)
					require.NoError(t, err)
					// silent auth should succeed this time
					_, err = cred.GetToken(context.Background(), testTRO)
					require.NoError(t, err)
				}
			})
		}
	}
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
						tokenProvider: func(ctx context.Context, scopes []string, tenant, subscription string) ([]byte, error) {
							if tenant != test.expected {
								t.Errorf(`unexpected tenantID "%s"`, tenant)
							}
							return mockAzTokenProviderSuccess(ctx, scopes, tenant, subscription)
						},
					}
					return NewAzureCLICredential(&o)
				},
			},
			{
				name: credNameAzureDeveloperCLI,
				ctor: func(azcore.ClientOptions) (azcore.TokenCredential, error) {
					o := AzureDeveloperCLICredentialOptions{
						AdditionallyAllowedTenants: test.allowed,
						tokenProvider: func(ctx context.Context, scopes []string, tenant string) ([]byte, error) {
							if tenant != test.expected {
								t.Errorf("unexpected tenantID %q", tenant)
							}
							return mockAzdTokenProviderSuccess(ctx, scopes, tenant)
						},
					}
					return NewAzureDeveloperCLICredential(&o)
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
					tokenRequestCallback: func(r *http.Request) *http.Response {
						if actual := strings.Split(r.URL.Path, "/")[1]; actual != test.expected {
							t.Fatalf("expected tenant %q, got %q", test.expected, actual)
						}
						return nil
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
		for _, credName := range []string{credNameAzureCLI, credNameAzureDeveloperCLI} {
			t.Run(fmt.Sprintf("DefaultAzureCredential/%s/%s", credName, test.desc), func(t *testing.T) {
				typeName := fmt.Sprintf("%T", &AzureCLICredential{})
				if credName == credNameAzureDeveloperCLI {
					typeName = fmt.Sprintf("%T", &AzureDeveloperCLICredential{})
				}
				called := false
				verifyTenant := func(tenant string) {
					called = true
					if tenant != test.expected {
						t.Fatalf("unexpected tenant %q", tenant)
					}
				}

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
				for _, source := range c.chain.sources {
					if fmt.Sprintf("%T", source) != typeName {
						continue
					}
					switch c := source.(type) {
					case *AzureCLICredential:
						c.opts.tokenProvider = func(ctx context.Context, scopes []string, tenant, subscription string) ([]byte, error) {
							verifyTenant(tenant)
							return mockAzTokenProviderSuccess(ctx, scopes, tenant, subscription)
						}
					case *AzureDeveloperCLICredential:
						c.opts.tokenProvider = func(ctx context.Context, scopes []string, tenant string) ([]byte, error) {
							verifyTenant(tenant)
							return mockAzdTokenProviderSuccess(ctx, scopes, tenant)
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
						t.Fatalf("%s wasn't invoked", credName)
					}
					break
				}
			})
		}
	}
}

func TestClaims(t *testing.T) {
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
		{
			name: credNameWorkloadIdentity,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				tokenFile := filepath.Join(t.TempDir(), "token")
				if err := os.WriteFile(tokenFile, []byte(tokenValue), os.ModePerm); err != nil {
					t.Fatalf("failed to write token file: %v", err)
				}
				o := WorkloadIdentityCredentialOptions{ClientID: fakeClientID, ClientOptions: co, TenantID: fakeTenantID, TokenFilePath: tokenFile}
				return NewWorkloadIdentityCredential(&o)
			},
		},
	} {
		for _, enableCAE := range []bool{true, false} {
			name := test.name
			if enableCAE {
				name += " CAE"
			}
			t.Run(name, func(t *testing.T) {
				reqs := 0
				sts := mockSTS{
					tokenRequestCallback: func(r *http.Request) *http.Response {
						if err := r.ParseForm(); err != nil {
							t.Error(err)
						}
						reqs++
						// Both requests should specify CP1 when CAE is enabled for the token.
						// We check only for substrings because MSAL is responsible for formatting claims.
						actual := fmt.Sprint(r.Form["claims"])
						if strings.Contains(actual, "CP1") != enableCAE {
							t.Fatalf(`unexpected claims "%v"`, actual)
						}
						if reqs == 2 {
							// the second GetToken call specifies claims we should find in the following token request
							if !strings.Contains(strings.ReplaceAll(actual, " ", ""), claim) {
								t.Fatalf(`unexpected claims "%v"`, actual)
							}
						}
						return nil
					},
				}
				o := azcore.ClientOptions{Transport: &sts}
				cred, err := test.ctor(o)
				if err != nil {
					t.Fatal(err)
				}
				tro := policy.TokenRequestOptions{EnableCAE: enableCAE, Scopes: []string{"A"}}
				if _, err = cred.GetToken(context.Background(), tro); err != nil {
					t.Fatal(err)
				}
				tro = policy.TokenRequestOptions{Claims: fmt.Sprintf("{%s}", claim), EnableCAE: enableCAE, Scopes: []string{"B"}}
				if _, err = cred.GetToken(context.Background(), tro); err != nil {
					t.Fatal(err)
				}
				if reqs != 2 {
					t.Fatalf("expected %d token requests, got %d", 2, reqs)
				}
			})
		}
	}
}

func TestCLIArgumentValidation(t *testing.T) {
	invalidRunes := "|';&"
	for _, test := range []struct {
		ctor func() (azcore.TokenCredential, error)
		name string
	}{
		{
			ctor: func() (azcore.TokenCredential, error) {
				return NewAzureCLICredential(nil)
			},
			name: credNameAzureCLI,
		},
		{
			ctor: func() (azcore.TokenCredential, error) {
				return NewAzureDeveloperCLICredential(nil)
			},
			name: credNameAzureDeveloperCLI,
		},
	} {
		t.Run(fmt.Sprintf("%s/scope", test.name), func(t *testing.T) {
			cred, err := test.ctor()
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range invalidRunes {
				_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{
					Scopes: []string{liveTestScope + string(r)},
				})
				if err == nil {
					t.Fatalf("expected an error for a scope containing %q", r)
				}
			}
		})
		t.Run(fmt.Sprintf("%s/tenant", test.name), func(t *testing.T) {
			cred, err := test.ctor()
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range invalidRunes {
				_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{
					TenantID: fakeTenantID + string(r),
				})
				if err == nil {
					t.Fatalf("expected an error for a tenant containing %q", r)
				}
			}
		})
	}
	t.Run(credNameAzureCLI+"/subscription", func(t *testing.T) {
		for _, r := range invalidRunes {
			if _, err := NewAzureCLICredential(&AzureCLICredentialOptions{Subscription: string(r)}); err == nil {
				t.Errorf("expected an error for a subscription containing %q", r)
			}
		}
	})
}

func TestResolveTenant(t *testing.T) {
	credName := "testcred"
	defaultTenant := "default-tenant"
	otherTenant := "other-tenant"
	for _, test := range []struct {
		allowed          []string
		expected, tenant string
		expectError      bool
	}{
		// no alternate tenant specified -> should get default
		{expected: defaultTenant},
		{allowed: []string{""}, expected: defaultTenant},
		{allowed: []string{"*"}, expected: defaultTenant},
		{allowed: []string{otherTenant}, expected: defaultTenant},

		// alternate tenant specified and allowed -> should get that tenant
		{allowed: []string{"*"}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{otherTenant}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{"not-" + otherTenant, otherTenant}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{"not-" + otherTenant, "*"}, expected: otherTenant, tenant: otherTenant},

		// invalid or not allowed tenant -> should get an error
		{tenant: otherTenant, expectError: true},
		{allowed: []string{""}, tenant: otherTenant, expectError: true},
		{allowed: []string{defaultTenant}, tenant: otherTenant, expectError: true},
		{tenant: badTenantID, expectError: true},
		{allowed: []string{""}, tenant: badTenantID, expectError: true},
		{allowed: []string{"*", badTenantID}, tenant: badTenantID, expectError: true},
		{tenant: "invalid@tenant", expectError: true},
		{tenant: "invalid/tenant", expectError: true},
		{tenant: "invalid(tenant", expectError: true},
		{tenant: "invalid:tenant", expectError: true},
	} {
		t.Run("", func(t *testing.T) {
			tenant, err := resolveTenant(defaultTenant, test.tenant, credName, test.allowed)
			if err != nil {
				if test.expectError {
					if validTenantID(test.tenant) && !strings.Contains(err.Error(), credName) {
						t.Fatalf("expected error to contain %q, got %q", credName, err.Error())
					}
					return
				}
				t.Fatal(err)
			} else if test.expectError {
				t.Fatal("expected an error")
			}
			if tenant != test.expected {
				t.Fatalf(`expected "%s", got "%s"`, test.expected, tenant)
			}
		})
	}
}

func TestTokenCachePersistenceOptions(t *testing.T) {
	af := filepath.Join(t.TempDir(), t.Name()+credNameWorkloadIdentity)
	if err := os.WriteFile(af, []byte("assertion"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	before := internal.NewCache
	t.Cleanup(func() { internal.NewCache = before })
	for _, test := range []struct {
		desc    string
		options *TokenCachePersistenceOptions
		err     error
	}{
		{
			desc: "nil options",
		},
		{
			desc:    "default options",
			options: &TokenCachePersistenceOptions{},
		},
		{
			desc:    "all options set",
			options: &TokenCachePersistenceOptions{AllowUnencryptedStorage: true, Name: "name"},
		},
	} {
		internal.NewCache = func(o *internal.TokenCachePersistenceOptions, _ bool) (cache.ExportReplace, error) {
			if (test.options == nil) != (o == nil) {
				t.Fatalf("expected %v, got %v", test.options, o)
			}
			if test.options != nil {
				if test.options.AllowUnencryptedStorage != o.AllowUnencryptedStorage {
					t.Fatalf("expected AllowUnencryptedStorage %v, got %v", test.options.AllowUnencryptedStorage, o.AllowUnencryptedStorage)
				}
				if test.options.Name != o.Name {
					t.Fatalf("expected Name %q, got %q", test.options.Name, o.Name)
				}
			}
			return nil, nil
		}
		for _, subtest := range []struct {
			ctor func(azcore.ClientOptions, *TokenCachePersistenceOptions) (azcore.TokenCredential, error)
			env  map[string]string
			name string
		}{
			{
				name: credNameAssertion,
				ctor: func(co azcore.ClientOptions, tco *TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
					o := ClientAssertionCredentialOptions{ClientOptions: co, TokenCachePersistenceOptions: tco}
					return NewClientAssertionCredential(fakeTenantID, fakeClientID, func(context.Context) (string, error) { return "...", nil }, &o)
				},
			},
			{
				name: credNameCert,
				ctor: func(co azcore.ClientOptions, tco *TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
					o := ClientCertificateCredentialOptions{ClientOptions: co, TokenCachePersistenceOptions: tco}
					return NewClientCertificateCredential(fakeTenantID, fakeClientID, allCertTests[0].certs, allCertTests[0].key, &o)
				},
			},
			{
				name: credNameDeviceCode,
				ctor: func(co azcore.ClientOptions, tco *TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
					o := DeviceCodeCredentialOptions{
						ClientOptions:                co,
						TokenCachePersistenceOptions: tco,
						UserPrompt:                   func(context.Context, DeviceCodeMessage) error { return nil },
					}
					return NewDeviceCodeCredential(&o)
				},
			},
			{
				name: credNameSecret,
				ctor: func(co azcore.ClientOptions, tco *TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
					o := ClientSecretCredentialOptions{ClientOptions: co, TokenCachePersistenceOptions: tco}
					return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, &o)
				},
			},
			{
				name: credNameUserPassword,
				ctor: func(co azcore.ClientOptions, tco *TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
					o := UsernamePasswordCredentialOptions{ClientOptions: co, TokenCachePersistenceOptions: tco}
					return NewUsernamePasswordCredential(fakeTenantID, fakeClientID, fakeUsername, "password", &o)
				},
			},
		} {
			t.Run(fmt.Sprintf("%s/%s", subtest.name, test.desc), func(t *testing.T) {
				for k, v := range subtest.env {
					t.Setenv(k, v)
				}
				c, err := subtest.ctor(policy.ClientOptions{Transport: &mockSTS{}}, test.options)
				if err != nil {
					t.Fatal(err)
				}
				_, err = c.GetToken(context.Background(), testTRO)
				if err != nil {
					if !errors.Is(err, test.err) {
						t.Fatalf("expected %v, got %v", test.err, err)
					}
				} else if test.err != nil {
					t.Fatal("expected an error")
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

var _ msalConfidentialClient = (*fakeConfidentialClient)(nil)

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

var _ msalPublicClient = (*fakePublicClient)(nil)
