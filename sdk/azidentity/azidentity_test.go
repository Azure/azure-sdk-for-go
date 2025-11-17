//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/managedidentity"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
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
	accessTokenRespSuccess = []byte(fmt.Sprintf(`{"access_token": "%s","expires_in": %d,"token_type":"Bearer"}`, tokenValue, tokenExpiresIn))
	ctx                    = context.Background()
	testTRO                = policy.TokenRequestOptions{Scopes: []string{liveTestScope}}
)

// constants for this file
const (
	testHost = "https://localhost"
)

// compact removes whitespace from errors to simplify validation
func compact(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\n', '\t':
			return -1
		default:
			return r
		}
	}, s)
}

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

func TestResponseErrors(t *testing.T) {
	content := "no tokens here"
	statusCode := http.StatusTeapot
	validate := func(t *testing.T, err error) {
		require.Error(t, err)
		flatErr := compact(err.Error())
		actual := strings.Count(flatErr, compact(http.StatusText(statusCode)))
		require.Equal(t, 1, actual, "error message should include response exactly once:\n%s", err.Error())
		actual = strings.Count(flatErr, compact(content))
		require.Equal(t, 1, actual, "error message should include body exactly once:\n%s", err.Error())
	}

	for _, client := range []struct {
		name string
		ctor func(co policy.ClientOptions) (azcore.TokenCredential, error)
	}{
		{
			name: "confidential",
			ctor: func(co policy.ClientOptions) (azcore.TokenCredential, error) {
				return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, &ClientSecretCredentialOptions{ClientOptions: co})
			},
		},
		{
			name: "managed identity",
			ctor: func(co policy.ClientOptions) (azcore.TokenCredential, error) {
				return NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: co})
			},
		},
		{
			name: "public",
			ctor: func(co policy.ClientOptions) (azcore.TokenCredential, error) {
				return NewUsernamePasswordCredential(fakeTenantID, fakeClientID, "username", "password", &UsernamePasswordCredentialOptions{ClientOptions: co})
			},
		},
	} {
		t.Run(client.name, func(t *testing.T) {
			cred, err := client.ctor(policy.ClientOptions{
				Retry: policy.RetryOptions{MaxRetries: -1},
				Transport: &mockSTS{
					tokenRequestCallback: func(*http.Request) *http.Response {
						return &http.Response{
							Body:       io.NopCloser(bytes.NewBufferString(content)),
							Status:     http.StatusText(statusCode),
							StatusCode: statusCode,
						}
					},
				},
			})
			require.NoError(t, err)
			_, err = cred.GetToken(ctx, testTRO)
			validate(t, err)

			t.Run("ChainedTokenCredential", func(t *testing.T) {
				chain, err := NewChainedTokenCredential([]azcore.TokenCredential{cred}, nil)
				require.NoError(t, err)
				_, err = chain.GetToken(ctx, testTRO)
				validate(t, err)
			})
		})
	}
}

func TestTenantID(t *testing.T) {
	type tc struct {
		name           string
		ctor           func(tenant string) (azcore.TokenCredential, error)
		tenantOptional bool
	}
	for _, test := range []tc{
		{
			name: credNameAssertion,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewClientAssertionCredential(tenant, fakeClientID, func(context.Context) (string, error) { return "", nil }, nil)
			},
		},
		{
			name: credNameAzureCLI,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewAzureCLICredential(&AzureCLICredentialOptions{
					TenantID: tenant,
				})
			},
			tenantOptional: true,
		},
		{
			name: credNameAzurePowerShell,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{
					TenantID: tenant,
				})
			},
			tenantOptional: true,
		},
		{
			name: credNameAzureDeveloperCLI,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
					TenantID: tenant,
				})
			},
			tenantOptional: true,
		},
		{
			name: credNameAzurePipelines,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewAzurePipelinesCredential(tenant, fakeClientID, "service-connection", tokenValue, nil)
			},
		},
		{
			name: credNameBrowser,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
					TenantID: tenant,
				})
			},
			tenantOptional: true,
		},
		{
			name: credNameCert,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewClientCertificateCredential(tenant, fakeClientID, allCertTests[0].certs, allCertTests[0].key, nil)
			},
		},
		{
			name: credNameDeviceCode,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewDeviceCodeCredential(&DeviceCodeCredentialOptions{
					TenantID: tenant,
				})
			},
			tenantOptional: true,
		},
		{
			name: credNameOBO + "/cert",
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewOnBehalfOfCredentialWithCertificate(tenant, fakeClientID, "assertion", allCertTests[0].certs, allCertTests[0].key, nil)
			},
		},
		{
			name: credNameOBO + "/secret",
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewOnBehalfOfCredentialWithSecret(tenant, fakeClientID, "assertion", fakeSecret, nil)
			},
		},
		{
			name: credNameSecret,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewClientSecretCredential(tenant, fakeClientID, fakeSecret, nil)
			},
		},
		{
			name: credNameUserPassword,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				return NewUsernamePasswordCredential(tenant, fakeClientID, "username", "password", nil)
			},
		},
		{
			name: credNameWorkloadIdentity,
			ctor: func(tenant string) (azcore.TokenCredential, error) {
				t.Setenv(azureTenantID, tenant)
				return NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
					ClientID:      fakeClientID,
					TokenFilePath: "...",
				})
			},
		},
	} {
		t.Run(test.name+"/empty", func(t *testing.T) {
			_, err := test.ctor("")
			if test.tenantOptional {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, "tenant")
			}
		})
		t.Run(test.name+"/invalid", func(t *testing.T) {
			_, err := test.ctor(badTenantID)
			require.ErrorContains(t, err, "tenant")
		})
	}
}

type testCache []byte

func (c *testCache) Export(_ context.Context, m cache.Marshaler, _ cache.ExportHints) (err error) {
	*c, err = m.Marshal()
	return
}

func (c *testCache) Replace(_ context.Context, u cache.Unmarshaler, _ cache.ReplaceHints) error {
	return u.Unmarshal(*c)
}

func TestUnavailableIfInDAC(t *testing.T) {
	err := newAuthenticationFailedError(credNameOBO, "it didn't work", nil)
	err = unavailableIfInDAC(err, false)
	require.ErrorAs(t, err, new(*AuthenticationFailedError))

	err = unavailableIfInDAC(err, true)
	require.ErrorAs(t, err, new(credentialUnavailable))

	err = unavailableIfInDAC(err, false)
	require.ErrorAs(t, err, new(credentialUnavailable))

	require.Equal(t, 1, strings.Count(err.Error(), credNameOBO))
}

func TestUserAuthentication(t *testing.T) {
	type authenticater interface {
		azcore.TokenCredential
		Authenticate(context.Context, *policy.TokenRequestOptions) (AuthenticationRecord, error)
	}
	for _, credential := range []struct {
		name                    string
		interactive, recordable bool
		new                     func(Cache, azcore.ClientOptions, AuthenticationRecord, bool) (authenticater, error)
	}{
		{
			name: credNameBrowser,
			new: func(c Cache, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				return NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
					AdditionallyAllowedTenants:     []string{"*"},
					AuthenticationRecord:           ar,
					Cache:                          c,
					ClientOptions:                  co,
					DisableAutomaticAuthentication: disableAutoAuth,
				})
			},
			interactive: true,
		},
		{
			name: credNameDeviceCode,
			new: func(c Cache, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				o := DeviceCodeCredentialOptions{
					AdditionallyAllowedTenants:     []string{"*"},
					AuthenticationRecord:           ar,
					Cache:                          c,
					ClientOptions:                  co,
					DisableAutomaticAuthentication: disableAutoAuth,
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
			new: func(c Cache, co azcore.ClientOptions, ar AuthenticationRecord, disableAutoAuth bool) (authenticater, error) {
				opts := UsernamePasswordCredentialOptions{
					AdditionallyAllowedTenants: []string{"*"},
					AuthenticationRecord:       ar,
					Cache:                      c,
					ClientOptions:              co,
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
				cred, err := credential.new(Cache{}, co, AuthenticationRecord{}, false)
				require.NoError(t, err)
				_, err = cred.Authenticate(context.Background(), nil)
				require.NoError(t, err)

				t.Setenv(azureAuthorityHost, cc.ActiveDirectoryAuthorityHost)
				cred, err = credential.new(Cache{}, azcore.ClientOptions{Transport: &sts}, AuthenticationRecord{}, false)
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
				if credential.name == credNameUserPassword {
					t.Skip("the test tenant requires MFA")
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

			cred, err := credential.new(Cache{}, co, AuthenticationRecord{}, false)
			require.NoError(t, err)
			ar, err := cred.Authenticate(context.Background(), &testTRO)
			require.NoError(t, err)

			// some fields of the returned AuthenticationRecord should have specific values
			require.Equal(t, developerSignOnClientID, ar.ClientID)
			require.Equal(t, supportedAuthRecordVersions[0], ar.Version)
			// all others should have nonempty values
			v := reflect.Indirect(reflect.ValueOf(&ar))
			for _, f := range reflect.VisibleFields(reflect.TypeOf(ar)) {
				s := v.FieldByIndex(f.Index).Addr().Interface().(*string)
				require.NotEmpty(t, *s)
			}
			require.Equal(t, 1, counter.count)
		})

		t.Run("PersistentCache/"+credential.name, func(t *testing.T) {
			if credential.name == credNameBrowser && !runManualTests {
				t.Skipf("set %s to run this test", azidentityRunManualTests)
			}
			tokenReqs := 0
			c := internal.NewCache(func(bool) (cache.ExportReplace, error) {
				return &testCache{}, nil
			})
			co := azcore.ClientOptions{Transport: &mockSTS{
				tokenRequestCallback: func(*http.Request) *http.Response {
					tokenReqs++
					return nil
				},
			}}

			cred, err := credential.new(c, co, AuthenticationRecord{}, false)
			require.NoError(t, err)
			record, err := cred.Authenticate(ctx, &testTRO)
			require.NoError(t, err)
			_, err = cred.GetToken(ctx, testTRO)
			require.NoError(t, err)
			require.Equal(t, 1, tokenReqs)

			// cred2 should return the token cached by cred
			cred2, err := credential.new(c, co, record, true)
			require.NoError(t, err)
			_, err = cred2.GetToken(ctx, testTRO)
			require.NoError(t, err)
			require.Equal(t, 1, tokenReqs)

			// cred should request a new token because the cached one isn't a CAE token
			caeTRO := testTRO
			caeTRO.EnableCAE = true
			_, err = cred.GetToken(ctx, caeTRO)
			require.NoError(t, err)
			require.Equal(t, 2, tokenReqs)
		})

		if credential.interactive {
			t.Run("DisableAutomaticAuthentication/"+credential.name, func(t *testing.T) {
				cred, err := credential.new(Cache{}, policy.ClientOptions{Transport: &mockSTS{}}, AuthenticationRecord{}, true)
				require.NoError(t, err)
				expected := policy.TokenRequestOptions{
					Claims:    "claims",
					EnableCAE: true,
					Scopes:    []string{"scope"},
					TenantID:  "tenant",
				}
				_, err = cred.GetToken(context.Background(), expected)
				require.Contains(t, err.Error(), credential.name)
				require.Contains(t, err.Error(), "Call Authenticate")
				var actual *AuthenticationRequiredError
				require.ErrorAs(t, err, &actual)
				require.Equal(t, expected, actual.TokenRequestOptions)

				if credential.name != credNameBrowser || runManualTests {
					_, err = cred.Authenticate(context.Background(), &testTRO)
					require.NoError(t, err)
					// silent auth should succeed this time
					_, err = cred.GetToken(context.Background(), testTRO)
					require.NoError(t, err)
				}
			})
			t.Run("DisableAutomaticAuthentication/ChainedTokenCredential/"+credential.name, func(t *testing.T) {
				cred, err := credential.new(Cache{}, policy.ClientOptions{}, AuthenticationRecord{}, true)
				require.NoError(t, err)
				expected := azcore.AccessToken{ExpiresOn: time.Now().UTC(), Token: tokenValue}
				fake := NewFakeCredential()
				fake.SetResponse(expected, nil)
				chain, err := NewChainedTokenCredential([]azcore.TokenCredential{cred, fake}, nil)
				require.NoError(t, err)
				// ChainedTokenCredential should continue iterating when a credential returns
				// AuthenticationRequiredError i.e., it should call fake.GetToken() and return the expected token
				actual, err := chain.GetToken(context.Background(), testTRO)
				require.NoError(t, err)
				require.Equal(t, expected, actual)
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

func TestGetTokenRequiresScopes(t *testing.T) {
	for _, ctor := range []func() (azcore.TokenCredential, error){
		func() (azcore.TokenCredential, error) { return NewAzureCLICredential(nil) },
		func() (azcore.TokenCredential, error) { return NewAzurePowerShellCredential(nil) },
		func() (azcore.TokenCredential, error) { return NewAzureDeveloperCLICredential(nil) },
		func() (azcore.TokenCredential, error) {
			return NewClientAssertionCredential(
				fakeTenantID, fakeClientID, func(context.Context) (string, error) { return "", nil }, nil,
			)
		},
		func() (azcore.TokenCredential, error) {
			return NewClientCertificateCredential(
				fakeTenantID, fakeClientID, allCertTests[0].certs, allCertTests[0].key, nil,
			)
		},
		func() (azcore.TokenCredential, error) {
			return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, nil)
		},
		func() (azcore.TokenCredential, error) { return NewDeviceCodeCredential(nil) },
		func() (azcore.TokenCredential, error) { return NewInteractiveBrowserCredential(nil) },
		func() (azcore.TokenCredential, error) { return NewManagedIdentityCredential(nil) },
		func() (azcore.TokenCredential, error) {
			return NewOnBehalfOfCredentialWithSecret(
				fakeTenantID, fakeClientID, "assertion", fakeSecret, nil,
			)
		},
		func() (azcore.TokenCredential, error) {
			return NewUsernamePasswordCredential(fakeTenantID, fakeClientID, fakeUsername, "password", nil)
		},
		func() (azcore.TokenCredential, error) {
			return NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
				ClientID: fakeClientID, TokenFilePath: ".", TenantID: fakeTenantID,
			})
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
	type testCase struct {
		allowed                                        []string
		desc, expected, ctorTenant, tokenRequestTenant string
		err                                            bool
	}
	cases := []testCase{
		{
			desc:               "all tenants allowed",
			allowed:            []string{"*"},
			expected:           tenantA,
			ctorTenant:         fakeTenantID,
			tokenRequestTenant: tenantA,
		},
		{
			desc:               "tenant explicitly allowed",
			allowed:            []string{tenantA, tenantB},
			ctorTenant:         fakeTenantID,
			expected:           tenantA,
			tokenRequestTenant: tenantA,
		},
		{
			desc:               "tenant explicitly allowed",
			allowed:            []string{tenantA, tenantB},
			ctorTenant:         fakeTenantID,
			expected:           tenantB,
			tokenRequestTenant: tenantB,
		},
		{
			desc:               "tenant not allowed",
			allowed:            []string{tenantA},
			ctorTenant:         fakeTenantID,
			tokenRequestTenant: tenantB,
			err:                true,
		},
	}
	for _, credential := range []struct {
		ctor           func(azcore.ClientOptions, testCase, *testing.T) (azcore.TokenCredential, error)
		env            func(testCase) map[string]string
		name           string
		tenantOptional bool
	}{
		{
			name: credNameAssertion,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := ClientAssertionCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewClientAssertionCredential(tc.ctorTenant, fakeClientID, func(context.Context) (string, error) { return "...", nil }, &o)
			},
		},
		{
			name: credNameAzureCLI,
			ctor: func(_ azcore.ClientOptions, tc testCase, t *testing.T) (azcore.TokenCredential, error) {
				o := AzureCLICredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					TenantID:                   tc.ctorTenant,
					exec: func(ctx context.Context, credName string, commandLine string) ([]byte, error) {
						require.Contains(t, commandLine, " --tenant "+tc.expected)
						return mockAzSuccess(ctx, credName, commandLine)
					},
				}
				return NewAzureCLICredential(&o)
			},
			tenantOptional: true,
		},
		{
			name: credNameAzurePowerShell,
			ctor: func(_ azcore.ClientOptions, tc testCase, t *testing.T) (azcore.TokenCredential, error) {
				o := AzurePowerShellCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					TenantID:                   tc.ctorTenant,
					exec: func(ctx context.Context, credName string, commandLine string) ([]byte, error) {
						splitCommand := strings.Split(commandLine, " ")
						encodedScript := splitCommand[len(splitCommand)-1]
						decodedScript, err := base64DecodeUTF16LE(encodedScript)
						require.NoError(t, err)
						require.Contains(t, decodedScript, fmt.Sprintf("$params['TenantId'] = '%s'", tc.expected))
						return mockAzurePowerShellSuccess(ctx, credName, commandLine)
					},
				}
				return NewAzurePowerShellCredential(&o)
			},
			tenantOptional: true,
		},
		{
			name: credNameAzureDeveloperCLI,
			ctor: func(_ azcore.ClientOptions, tc testCase, t *testing.T) (azcore.TokenCredential, error) {
				o := AzureDeveloperCLICredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					exec: func(ctx context.Context, credName string, command string) ([]byte, error) {
						require.Contains(t, command, " --tenant-id "+tc.expected)
						return mockAzdSuccess(ctx, credName, command)
					},
				}
				return NewAzureDeveloperCLICredential(&o)
			},
			tenantOptional: true,
		},
		{
			name: credNameAzurePipelines,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := AzurePipelinesCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewAzurePipelinesCredential(tc.ctorTenant, fakeClientID, "service-connection", tokenValue, &o)
			},
			env: func(testCase) map[string]string {
				return map[string]string{systemOIDCRequestURI: "https://localhost/instance"}
			},
		},
		{
			name: credNameCert,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := ClientCertificateCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewClientCertificateCredential(tc.ctorTenant, fakeClientID, allCertTests[0].certs, allCertTests[0].key, &o)
			},
		},
		{
			name: credNameDeviceCode,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := DeviceCodeCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					ClientOptions:              co,
					UserPrompt:                 func(context.Context, DeviceCodeMessage) error { return nil },
				}
				return NewDeviceCodeCredential(&o)
			},
			tenantOptional: true,
		},
		{
			name: credNameOBO,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := OnBehalfOfCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					ClientOptions:              co,
				}
				return NewOnBehalfOfCredentialWithSecret(tc.ctorTenant, fakeClientID, "assertion", fakeSecret, &o)
			},
		},
		{
			name: credNameSecret,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := ClientSecretCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewClientSecretCredential(tc.ctorTenant, fakeClientID, fakeSecret, &o)
			},
		},
		{
			name: credNameUserPassword,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := UsernamePasswordCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewUsernamePasswordCredential(tc.ctorTenant, fakeClientID, fakeUsername, "password", &o)
			},
		},
		{
			name: credNameWorkloadIdentity,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				return NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					ClientID:                   fakeClientID,
					ClientOptions:              co,
					TenantID:                   tc.ctorTenant,
					TokenFilePath:              af,
				})
			},
		},
		{
			name: "DefaultAzureCredential/" + credNameAzureCLI,
			ctor: func(_ azcore.ClientOptions, tc testCase, t *testing.T) (azcore.TokenCredential, error) {
				srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
				t.Cleanup(close)
				srv.SetResponse(mock.WithStatusCode(http.StatusBadRequest))

				c, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					ClientOptions:              policy.ClientOptions{Transport: srv},
				})
				require.NoError(t, err)
				called := false
				t.Cleanup(func() {
					if tc.err {
						require.False(t, called)
					} else {
						require.True(t, called)
					}
				})
				for _, source := range c.chain.sources {
					if c, ok := source.(*AzureCLICredential); ok {
						c.opts.exec = func(ctx context.Context, credName string, commandLine string) ([]byte, error) {
							called = true
							require.Contains(t, commandLine, " --tenant "+tc.expected)
							return mockAzSuccess(ctx, credName, commandLine)
						}
						break
					}
				}
				return c, nil
			},
			tenantOptional: true,
		},
		{
			name: "DefaultAzureCredential/" + credNameAzureDeveloperCLI,
			ctor: func(_ azcore.ClientOptions, tc testCase, t *testing.T) (azcore.TokenCredential, error) {
				srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
				t.Cleanup(close)
				srv.SetResponse(mock.WithStatusCode(http.StatusBadRequest))

				c, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
					AdditionallyAllowedTenants: tc.allowed,
					ClientOptions:              policy.ClientOptions{Transport: srv},
				})
				require.NoError(t, err)
				called := false
				t.Cleanup(func() {
					if tc.err {
						require.False(t, called)
					} else {
						require.True(t, called)
					}
				})
				for _, source := range c.chain.sources {
					switch c := source.(type) {
					case *AzureCLICredential:
						c.opts.exec = func(context.Context, string, string) ([]byte, error) {
							return nil, newCredentialUnavailableError(credNameAzureCLI, "...")
						}
					case *AzureDeveloperCLICredential:
						c.opts.exec = func(ctx context.Context, credName string, command string) ([]byte, error) {
							called = true
							require.Contains(t, command, " --tenant-id "+tc.expected)
							return mockAzdSuccess(ctx, credName, command)
						}
					}
				}
				return c, nil
			},
			tenantOptional: true,
		},
		{
			name: "DefaultAzureCredential/EnvironmentCredential",
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := DefaultAzureCredentialOptions{ClientOptions: co, TenantID: tc.tokenRequestTenant}
				return NewDefaultAzureCredential(&o)
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(tc.allowed, ";"),
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   tc.ctorTenant,
				}
			},
		},
		{
			name: "DefaultAzureCredential/EnvironmentCredential/option-overrides-env",
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := DefaultAzureCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co, TenantID: tc.tokenRequestTenant}
				return NewDefaultAzureCredential(&o)
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: "not-" + tc.tokenRequestTenant,
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   tc.ctorTenant,
				}
			},
		},
		{
			name: "DefaultAzureCredential/" + credNameWorkloadIdentity,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				o := DefaultAzureCredentialOptions{AdditionallyAllowedTenants: tc.allowed, ClientOptions: co}
				return NewDefaultAzureCredential(&o)
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(tc.allowed, ";"),
					azureAuthorityHost:              "https://login.microsoftonline.com",
					azureClientID:                   fakeClientID,
					azureFederatedTokenFile:         af,
					azureTenantID:                   tc.ctorTenant,
				}
			},
		},
		{
			name: "EnvironmentCredential/" + credNameCert,
			ctor: func(co azcore.ClientOptions, _ testCase, _ *testing.T) (azcore.TokenCredential, error) {
				return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(tc.allowed, ";"),
					azureClientCertificatePath:      "testdata/certificate.pem",
					azureClientID:                   fakeClientID,
					azureTenantID:                   tc.ctorTenant,
				}
			},
		},
		{
			name: "EnvironmentCredential/" + credNameSecret,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(tc.allowed, ";"),
					azureClientID:                   fakeClientID,
					azureClientSecret:               fakeSecret,
					azureTenantID:                   tc.ctorTenant,
				}
			},
		},
		{
			name: "EnvironmentCredential/" + credNameUserPassword,
			ctor: func(co azcore.ClientOptions, tc testCase, _ *testing.T) (azcore.TokenCredential, error) {
				return NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: co})
			},
			env: func(tc testCase) map[string]string {
				return map[string]string{
					azureAdditionallyAllowedTenants: strings.Join(tc.allowed, ";"),
					azureClientID:                   fakeClientID,
					azurePassword:                   "password",
					azureTenantID:                   tc.ctorTenant,
					azureUsername:                   fakeUsername,
				}
			},
		},
	} {
		for _, test := range cases {
			tro := policy.TokenRequestOptions{Scopes: []string{liveTestScope}, TenantID: test.tokenRequestTenant}
			t.Run(fmt.Sprintf("%s/%s", credential.name, test.desc), func(t *testing.T) {
				if credential.env != nil {
					for k, v := range credential.env(test) {
						t.Setenv(k, v)
					}
				}
				sts := mockSTS{
					tenant: test.tokenRequestTenant,
					tokenRequestCallback: func(r *http.Request) *http.Response {
						if actual := strings.Split(r.URL.Path, "/")[1]; actual != test.expected {
							t.Fatalf("expected tenant %q, got %q", test.expected, actual)
						}
						return nil
					},
				}
				c, err := credential.ctor(policy.ClientOptions{Transport: &sts}, test, t)
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

		if credential.tenantOptional {
			t.Run(credential.name+"/no tenant specified", func(t *testing.T) {
				tc := testCase{tokenRequestTenant: tenantA, expected: tenantA}
				c, err := credential.ctor(policy.ClientOptions{Transport: &mockSTS{}}, tc, t)
				require.NoError(t, err)
				_, err = c.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}, TenantID: tc.tokenRequestTenant})
				if tc.err {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		} else {
			t.Run(credential.name+"/no additionally allowed tenants", func(t *testing.T) {
				tc := testCase{err: true, ctorTenant: fakeTenantID, tokenRequestTenant: tenantA}
				if credential.env != nil {
					for k, v := range credential.env(tc) {
						t.Setenv(k, v)
					}
				}
				c, err := credential.ctor(policy.ClientOptions{Transport: &mockSTS{}}, tc, t)
				require.NoError(t, err)
				_, err = c.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}, TenantID: tc.tokenRequestTenant})
				require.Error(t, err)
			})
		}
	}

	for _, test := range cases {
		tro := policy.TokenRequestOptions{Scopes: []string{liveTestScope}, TenantID: test.tokenRequestTenant}
		t.Run(credNameBrowser, func(t *testing.T) {
			c, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
				AdditionallyAllowedTenants: test.allowed,
				// this enables testing the credential's tenant resolution without having to authenticate
				DisableAutomaticAuthentication: true,
			})
			require.NoError(t, err)
			_, err = c.GetToken(context.Background(), tro)
			if test.err {
				// the specified tenant isn't allowed, so the error should be about that
				require.ErrorContains(t, err, "AdditionallyAllowedTenants")
			} else {
				// tenant resolution should have succeeded because the specified tenant is allowed,
				// however the credential should have returned a different error because automatic
				// authentication is disabled
				var e *AuthenticationRequiredError
				require.ErrorAs(t, err, &e)
			}
		})
	}
}

func TestClaims(t *testing.T) {
	claim := `"test":"pass"`
	for _, test := range []struct {
		ctor func(azcore.ClientOptions) (azcore.TokenCredential, error)
		env  map[string]string
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
			name: credNameAzurePipelines,
			ctor: func(co azcore.ClientOptions) (azcore.TokenCredential, error) {
				o := AzurePipelinesCredentialOptions{ClientOptions: co}
				return NewAzurePipelinesCredential(fakeTenantID, fakeClientID, "service-connection", tokenValue, &o)
			},
			env: map[string]string{systemOIDCRequestURI: "https://localhost/foo/UserRealm/foo"},
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
				for k, v := range test.env {
					t.Setenv(k, v)
				}
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
				return NewAzurePowerShellCredential(nil)
			},
			name: credNameAzurePowerShell,
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

func TestRefreshOn(t *testing.T) {
	t.Run("confidential client", func(t *testing.T) {
		expected := confidential.AuthResult{}
		expected.Metadata.RefreshOn = time.Now().Add(42 * time.Minute)
		cred, err := confidential.NewCredFromSecret("...")
		require.NoError(t, err)
		client, err := newConfidentialClient(fakeTenantID, fakeClientID, "...", cred, confidentialClientOptions{})
		require.NoError(t, err)
		client.noCAE = fakeConfidentialClient{ar: expected}
		actual, err := client.GetToken(ctx, testTRO)
		require.NoError(t, err)
		require.True(t, expected.Metadata.RefreshOn.Equal(actual.RefreshOn))
	})
	t.Run("managed identity client", func(t *testing.T) {
		expected := managedidentity.AuthResult{}
		expected.Metadata.RefreshOn = time.Now().Add(42 * time.Minute)
		cred, err := NewManagedIdentityCredential(nil)
		require.NoError(t, err)
		cred.mic.msalClient = fakeManagedIdentityClient{ar: expected}
		actual, err := cred.GetToken(ctx, testTRO)
		require.NoError(t, err)
		require.True(t, expected.Metadata.RefreshOn.Equal(actual.RefreshOn))
	})
	t.Run("public client", func(t *testing.T) {
		expected := public.AuthResult{}
		expected.Metadata.RefreshOn = time.Now().Add(42 * time.Minute)
		client, err := newPublicClient(fakeTenantID, fakeClientID, credNameBrowser, publicClientOptions{})
		require.NoError(t, err)
		client.noCAE = fakePublicClient{ar: expected}
		actual, err := client.GetToken(ctx, testTRO)
		require.NoError(t, err)
		require.True(t, expected.Metadata.RefreshOn.Equal(actual.RefreshOn))
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

	t.Run("empty default tenant", func(t *testing.T) {
		expected := "expected"
		for _, additionalTenants := range [][]string{nil, {"*"}, {expected}} {
			actual, err := resolveTenant("", expected, credName, additionalTenants)
			require.NoError(t, err, "the specified tenant should be allowed given additionalTenants %v", additionalTenants)
			require.Equal(t, expected, actual)
		}
		actual, err := resolveTenant("", expected, credName, []string{"un" + expected})
		require.Error(t, err, "the specified tenant shouldn't be allowed because additionalTenants isn't empty and doesn't contain it")
		require.Empty(t, actual)
	})
}

func TestConfidentialClientPersistentCache(t *testing.T) {
	// for WorkloadIdentityCredential
	tfp := filepath.Join(t.TempDir(), "tokenfile")
	require.NoError(t, os.WriteFile(tfp, []byte("token"), 0600))
	for _, credential := range []struct {
		name string
		new  func(azcore.ClientOptions, Cache) (azcore.TokenCredential, error)
	}{
		{
			name: credNameAssertion,
			new: func(co azcore.ClientOptions, c Cache) (azcore.TokenCredential, error) {
				o := ClientAssertionCredentialOptions{Cache: c, ClientOptions: co}
				return NewClientAssertionCredential(fakeTenantID, fakeClientID, func(context.Context) (string, error) { return "...", nil }, &o)
			},
		},
		// TODO: set SYSTEM_OIDC_REQUEST_URI, fake response
		// {
		// 	name: credNameAzurePipelines,
		// 	new: func(co azcore.ClientOptions, c Cache) (azcore.TokenCredential, error) {
		// 		o := AzurePipelinesCredentialOptions{Cache: c, ClientOptions: co}
		// 		return NewAzurePipelinesCredential(fakeTenantID, fakeClientID, "service-connection", tokenValue, &o)
		// 	},
		// },
		{
			name: credNameCert,
			new: func(co azcore.ClientOptions, c Cache) (azcore.TokenCredential, error) {
				o := ClientCertificateCredentialOptions{Cache: c, ClientOptions: co}
				return NewClientCertificateCredential(fakeTenantID, fakeClientID, allCertTests[0].certs, allCertTests[0].key, &o)
			},
		},
		{
			name: credNameSecret,
			new: func(co azcore.ClientOptions, c Cache) (azcore.TokenCredential, error) {
				o := ClientSecretCredentialOptions{Cache: c, ClientOptions: co}
				return NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, &o)
			},
		},
		{
			name: credNameWorkloadIdentity,
			new: func(co azcore.ClientOptions, c Cache) (azcore.TokenCredential, error) {
				o := WorkloadIdentityCredentialOptions{
					Cache:         c,
					ClientID:      fakeClientID,
					ClientOptions: co,
					TenantID:      fakeTenantID,
					TokenFilePath: tfp,
				}
				return NewWorkloadIdentityCredential(&o)
			},
		},
	} {
		t.Run(credential.name, func(t *testing.T) {
			tokenReqs := 0
			c := internal.NewCache(func(bool) (cache.ExportReplace, error) {
				return &testCache{}, nil
			})
			sts := mockSTS{
				tokenRequestCallback: func(*http.Request) *http.Response {
					tokenReqs++
					return nil
				},
			}
			cred, err := credential.new(policy.ClientOptions{Transport: &sts}, c)
			require.NoError(t, err)
			_, err = cred.GetToken(context.Background(), testTRO)
			require.NoError(t, err)
			_, err = cred.GetToken(ctx, testTRO)
			require.NoError(t, err)
			require.Equal(t, 1, tokenReqs)

			// cred2 should return the token cached by cred
			cred2, err := credential.new(policy.ClientOptions{Transport: &sts}, c)
			require.NoError(t, err)
			_, err = cred2.GetToken(ctx, testTRO)
			require.NoError(t, err)
			require.Equal(t, 1, tokenReqs)

			// cred should request a new token because the cached one isn't a CAE token
			caeTRO := testTRO
			caeTRO.EnableCAE = true
			_, err = cred.GetToken(ctx, caeTRO)
			require.NoError(t, err)
			require.Equal(t, 2, tokenReqs)
		})
	}
}

func TestDoForClient(t *testing.T) {
	var (
		policyHeaderName  = "PolicyHeader"
		policyHeaderValue = "policyvalue"

		reqBody  = []byte(`{"request": "azidentity"}`)
		respBody = []byte(`{"response": "golang"}`)
	)

	tests := map[string]struct {
		method  string
		path    string
		body    io.Reader
		headers http.Header
	}{
		"happy path": {
			method: http.MethodGet,
			path:   "/foo/bar",
			body:   bytes.NewBuffer(reqBody),
			headers: http.Header{
				"Header": []string{"value1", "value2"},
			},
		},
		"no body": {
			method: http.MethodGet,
			path:   "/",
			body:   http.NoBody,
		},
		"nil body": {
			method: http.MethodGet,
			path:   "/",
			body:   nil,
		},
		"headers with empty value": {
			method: http.MethodGet,
			path:   "/",
			body:   http.NoBody,
			headers: http.Header{
				"Header": nil,
			},
		},
	}

	client, err := azcore.NewClient(module, version, azruntime.PipelineOptions{
		// add PerCall policy to ensure doForClient calls .Pipeline.Do()
		PerCall: []policy.Policy{
			policyFunc(func(req *policy.Request) (*http.Response, error) {
				req.Raw().Header.Set(policyHeaderName, policyHeaderValue)
				return req.Next()
			}),
		},
	}, nil)
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tt.method, req.Method)
				assert.Equal(t, tt.path, req.URL.Path)

				rb, err := io.ReadAll(req.Body)
				assert.NoError(t, err)

				if tt.body != nil && tt.body != http.NoBody {
					assert.Equal(t, string(reqBody), string(rb))
				} else {
					assert.Empty(t, rb)
				}

				for k, v := range tt.headers {
					assert.Equal(t, v, req.Header[k])
				}

				assert.Equal(t, policyHeaderValue, req.Header.Get(policyHeaderName))

				rw.Header().Set("content-type", "application/json")
				_, err = rw.Write(respBody)
				assert.NoError(t, err)
			}))
			defer server.Close()

			req, err := http.NewRequestWithContext(context.Background(), tt.method, server.URL+tt.path, tt.body)
			require.NoError(t, err)

			for k, vs := range tt.headers {
				for _, v := range vs {
					req.Header.Add(k, v)
				}
			}

			resp, err := doForClient(client, req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, string(respBody), string(b))
		})
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

type fakeManagedIdentityClient struct {
	// set ar to have all API calls return the provided AuthResult
	ar public.AuthResult
}

func (f fakeManagedIdentityClient) AcquireToken(context.Context, string, ...managedidentity.AcquireTokenOption) (managedidentity.AuthResult, error) {
	return f.ar, nil
}

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

// ==================================================================================================================================

type policyFunc func(*policy.Request) (*http.Response, error)

// Do implements the Policy interface on policyFunc.
func (pf policyFunc) Do(req *policy.Request) (*http.Response, error) {
	return pf(req)
}
