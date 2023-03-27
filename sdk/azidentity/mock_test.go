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
// MSAL metadata requests. All responses are success responses. Mock access
// tokens expire in 1 hour and have the value of the "tokenValue" constant.
type mockSTS struct {
	// tenant to include in metadata responses. This value must match a test's
	// expected tenant because metadata tells MSAL where to send token requests.
	// Defaults to the "fakeTenantID" constant.
	tenant string
	// tokenRequestCallback is called for every token request
	tokenRequestCallback func(*http.Request)
}

func (m *mockSTS) Do(req *http.Request) (*http.Response, error) {
	res := http.Response{StatusCode: http.StatusOK}
	tenant := m.tenant
	if tenant == "" {
		tenant = fakeTenantID
	}
	switch s := strings.Split(req.URL.Path, "/"); s[len(s)-1] {
	case "instance":
		res.Body = io.NopCloser(bytes.NewReader(getInstanceDiscoveryResponse(tenant)))
	case "openid-configuration":
		res.Body = io.NopCloser(bytes.NewReader(getTenantDiscoveryResponse(tenant)))
	case "devicecode":
		res.Body = io.NopCloser(strings.NewReader(`{"device_code":"...","expires_in":600,"interval":60}`))
	case "token":
		if m.tokenRequestCallback != nil {
			m.tokenRequestCallback(req)
		}
		res.Body = io.NopCloser(bytes.NewReader(accessTokenRespSuccess))
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
	return &res, nil
}
