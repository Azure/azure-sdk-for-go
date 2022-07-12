//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	accessTokenRespSuccess   = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenRespMalformed = `{"access_token": 0, "expires_in": 3600}`
	badTenantID              = "bad_tenant"
	tokenValue               = "new_token"
)

// constants for this file
const (
	testHost                = "https://localhost"
	tenantDiscoveryResponse = `{
		"token_endpoint": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/oauth2/v2.0/token",
		"token_endpoint_auth_methods_supported": [
		"client_secret_post",
		"private_key_jwt",
		"client_secret_basic"
		],
		"jwks_uri": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/discovery/v2.0/keys",
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
		"issuer": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/v2.0",
		"request_uri_parameter_supported": false,
		"userinfo_endpoint": "https://graph.microsoft.com/oidc/userinfo",
		"authorization_endpoint": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/oauth2/v2.0/authorize",
		"device_authorization_endpoint": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/oauth2/v2.0/devicecode",
		"http_logout_supported": true,
		"frontchannel_logout_supported": true,
		"end_session_endpoint": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/oauth2/v2.0/logout",
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
		"kerberos_endpoint": "https://login.microsoftonline.com/3c631bb7-a9f7-4343-a5ba-a6159135f1fc/kerberos",
		"tenant_region_scope": "NA",
		"cloud_instance_name": "microsoftonline.com",
		"cloud_graph_host_name": "graph.windows.net",
		"msgraph_host": "graph.microsoft.com",
		"rbac_url": "https://pas.windows.net"
		}`
)

func validateJWTRequestContainsHeader(t *testing.T, headerName string) mock.ResponsePredicate {
	return func(req *http.Request) bool {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal("Expected a request with the JWT in the body.")
		}
		bodystr := string(body)
		kvps := strings.Split(bodystr, "&")
		assertion := strings.Split(kvps[0], "=")
		token, _ := jwt.Parse(assertion[1], nil)
		if token == nil {
			t.Fatalf("Failed to parse the JWT token: %s.", assertion[1])
		}
		if _, ok := token.Header[headerName]; !ok {
			t.Fatalf("JWT did not contain the %s header", headerName)
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
			return NewClientSecretCredential("tenantID", "clientID", "secret", nil)
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

func Test_ValidTenantIDFalse(t *testing.T) {
	if validTenantID("bad@tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad/tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad(tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad)tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad:tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
}

func Test_ValidTenantIDTrue(t *testing.T) {
	if !validTenantID("goodtenant") {
		t.Fatal("Expected to receive true, but received false")
	}
	if !validTenantID("good-tenant") {
		t.Fatal("Expected to receive true, but received false")
	}
	if !validTenantID("good.tenant") {
		t.Fatal("Expected to receive true, but received false")
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
}

func (f fakeConfidentialClient) returnResult() (confidential.AuthResult, error) {
	if f.err != nil {
		return confidential.AuthResult{}, f.err
	}
	return f.ar, nil
}

func (f fakeConfidentialClient) AcquireTokenSilent(ctx context.Context, scopes []string, options ...confidential.AcquireTokenSilentOption) (confidential.AuthResult, error) {
	if f.silentAuth {
		return f.ar, nil
	}
	return confidential.AuthResult{}, errors.New("silent authentication failed")
}

func (f fakeConfidentialClient) AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, options ...confidential.AcquireTokenByAuthCodeOption) (confidential.AuthResult, error) {
	return f.returnResult()
}

func (f fakeConfidentialClient) AcquireTokenByCredential(ctx context.Context, scopes []string) (confidential.AuthResult, error) {
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

func (f fakePublicClient) AcquireTokenSilent(ctx context.Context, scopes []string, options ...public.AcquireTokenSilentOption) (public.AuthResult, error) {
	if f.silentAuth {
		return f.ar, nil
	}
	return public.AuthResult{}, errors.New("silent authentication failed")
}

func (f fakePublicClient) AcquireTokenByUsernamePassword(ctx context.Context, scopes []string, username string, password string) (public.AuthResult, error) {
	return f.returnResult()
}

func (f fakePublicClient) AcquireTokenByDeviceCode(ctx context.Context, scopes []string) (public.DeviceCode, error) {
	if f.err != nil {
		return public.DeviceCode{}, f.err
	}
	return f.dc, nil
}

func (f fakePublicClient) AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, options ...public.AcquireTokenByAuthCodeOption) (public.AuthResult, error) {
	return f.returnResult()
}

func (f fakePublicClient) AcquireTokenInteractive(ctx context.Context, scopes []string, options ...public.InteractiveAuthOption) (public.AuthResult, error) {
	return f.returnResult()
}

var _ publicClient = (*fakePublicClient)(nil)
