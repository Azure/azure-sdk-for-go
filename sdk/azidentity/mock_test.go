//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// mockSTS returns mock Azure AD responses so tests don't have to account for
// MSAL metadata requests. By default, all responses are success responses
// having a token which expires in 1 hour and whose value is the "tokenValue"
// constant. Set tokenRequestCallback to return a different *http.Response.
type mockSTS struct {
	// tenant to include in metadata responses. This value must match a test's
	// expected tenant because metadata tells MSAL where to send token requests.
	// Defaults to the "fakeTenantID" constant.
	tenant string
	// tokenRequestCallback is called for every token request. Return nil to
	// send a generic success response.
	tokenRequestCallback func(*http.Request) *http.Response
}

func (m *mockSTS) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{StatusCode: http.StatusOK}
	tenant := m.tenant
	if tenant == "" {
		tenant = fakeTenantID
	}
	switch s := strings.Split(req.URL.Path, "/"); s[len(s)-1] {
	case "instance":
		res.Body = io.NopCloser(bytes.NewReader(instanceMetadata(tenant)))
	case "openid-configuration":
		res.Body = io.NopCloser(bytes.NewReader(tenantMetadata(tenant)))
	case "devicecode":
		res.Body = io.NopCloser(strings.NewReader(`{"device_code":"...","expires_in":600,"interval":60}`))
	case "token":
		res.Body = io.NopCloser(bytes.NewReader(accessTokenRespSuccess))
		if m.tokenRequestCallback != nil {
			if r := m.tokenRequestCallback(req); r != nil {
				res = r
			}
		}
	default:
		// User realm metadata request paths look like "/common/UserRealm/user@domain".
		// Matching on the UserRealm segment avoids having to know the UPN.
		if s[len(s)-2] == "UserRealm" {
			res.Body = io.NopCloser(
				strings.NewReader(`{"account_type":"Managed","cloud_audience_urn":"urn","cloud_instance_name":"...","domain_name":"..."}`),
			)
		} else {
			panic("unexpected request " + req.URL.String())
		}
	}
	return res, nil
}

func instanceMetadata(tenant string) []byte {
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

func tenantMetadata(tenant string) []byte {
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
