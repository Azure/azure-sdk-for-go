// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// constants used throughout this package
const (
	accessTokenRespError     = `{"error": "invalid_client","error_description": "Invalid client secret is provided.","error_codes": [0],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=0"}`
	accessTokenRespSuccess   = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenRespMalformed = `{"access_token": 0, "expires_in": 3600}`
	tokenValue               = "new_token"
)

// configuration for live tests
var liveSP = struct {
	tenantID string
	clientID string
	secret   string
	pemPath  string
	pfxPath  string
	sniPath  string
}{
	tenantID: os.Getenv("IDENTITY_SP_TENANT_ID"),
	clientID: os.Getenv("IDENTITY_SP_CLIENT_ID"),
	secret:   os.Getenv("IDENTITY_SP_CLIENT_SECRET"),
	pemPath:  os.Getenv("IDENTITY_SP_CERT_PEM"),
	pfxPath:  os.Getenv("IDENTITY_SP_CERT_PFX"),
	sniPath:  os.Getenv("IDENTITY_SP_CERT_SNI"),
}

var liveTestScope = "https://management.core.windows.net//.default"

func init() {
	host := AuthorityHost(os.Getenv(azureAuthorityHost))
	switch host {
	case AzureChina:
		liveTestScope = "https://management.core.chinacloudapi.cn//.default"
	case AzureGovernment:
		liveTestScope = "https://management.core.usgovcloudapi.net//.default"
	}
}

func defaultTestPipeline(srv policy.Transporter, cred azcore.TokenCredential, scope string) runtime.Pipeline {
	retryOpts := policy.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
	}
	b := runtime.NewBearerTokenPolicy(cred, []string{scope}, nil)
	return runtime.NewPipeline("azidentity-test", version, nil, []policy.Policy{b}, &azcore.ClientOptions{Retry: retryOpts, Transport: srv})
}

// constants for this file
const (
	envHostString    = "https://mock.com/"
	customHostString = "https://custommock.com/"
)

// Set environment variables for the duration of a test. Restore their prior values
// after the test completes. Obviated by 1.17's T.Setenv
func setEnvironmentVariables(t *testing.T, vars map[string]string) {
	priorValues := make(map[string]string, len(vars))
	for k, v := range vars {
		priorValues[k] = os.Getenv(k)
		err := os.Setenv(k, v)
		if err != nil {
			t.Fatalf("Unexpected error setting %s: %v", k, err)
		}
	}

	t.Cleanup(func() {
		for k, v := range priorValues {
			err := os.Setenv(k, v)
			if err != nil {
				t.Fatalf("Unexpected error resetting %s: %v", k, err)
			}
		}
	})
}

func Test_SetEnvAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: envHostString})
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != envHostString {
		t.Fatalf("Unexpected error when get host from environment variable: %v", err)
	}
}

func Test_CustomAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: envHostString})
	authorityHost, err := setAuthorityHost(customHostString)
	if err != nil {
		t.Fatal(err)
	}
	// ensure env var doesn't override explicit value
	if authorityHost != customHostString {
		t.Fatalf("Unexpected host when get host from environment variable: %v", authorityHost)
	}
}

func Test_DefaultAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: ""})
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != string(AzurePublicCloud) {
		t.Fatalf("Unexpected host when set default AuthorityHost: %v", authorityHost)
	}
}

func Test_NonHTTPSAuthorityHost(t *testing.T) {
	setEnvironmentVariables(t, map[string]string{azureAuthorityHost: ""})
	authorityHost, err := setAuthorityHost("http://foo.com")
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
