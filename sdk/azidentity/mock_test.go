//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	mockClientInfo = "eyJ1aWQiOiJjNzNjNmYyOC1hZTVmLTQxM2QtYTlhMi1lMTFlNWFmNjY4ZjgiLCJ1dGlkIjoiZTBiZDIzMjEtMDdmYS00Y2YwLTg3YjgtMDBhYTJhNzQ3MzI5In0"
	mockIDT        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Imwzc1EtNTBjQ0g0eEJWWkxIVEd3blNSNzY4MCJ9.eyJhdWQiOiIwNGIwNzc5NS04ZGRiLTQ2MWEtYmJlZS0wMmY5ZTFiZjdiNDYiLCJpc3MiOiJodHRwczovL2xvZ2luLm1pY3Jvc29mdG9ubGluZS5jb20vYzU0ZmFjODgtM2RkMy00NjFmLWE3YzQtOGEzNjhlMDM0MGIzL3YyLjAiLCJpYXQiOjE2MzcxOTEyMTIsIm5iZiI6MTYzNzE5MTIxMiwiZXhwIjoxNjM3MTk1MTEyLCJhaW8iOiJBVVFBdS84VEFBQUFQMExOZGNRUXQxNmJoSkFreXlBdjFoUGJuQVhtT0o3RXJDVHV4N0hNTjhHd2VMb2FYMWR1cDJhQ2Y0a0p5bDFzNmovSzF5R05DZmVIQlBXM21QUWlDdz09IiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTBiZDIzMjEtMDdmYS00Y2YwLTg3YjgtMDBhYTJhNzQ3MzI5LyIsIm5hbWUiOiJJZGVudGl0eSBUZXN0IFVzZXIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJpZGVudGl0eXRlc3R1c2VyQGF6dXJlc2Rrb3V0bG9vay5vbm1pY3Jvc29mdC5jb20iLCJyaCI6IjAuQVMwQWlLeFB4ZE05SDBhbnhJbzJqZ05BczVWM3NBVGJqUnBHdS00Qy1lR19lMFl0QUxFLiIsInN1YiI6ImMxYTBsY2xtbWxCYW9wc0MwVmlaLVpPMjFCT2dSUXE3SG9HRUtOOXloZnMiLCJ0aWQiOiJjNTRmYWM4OC0zZGQzLTQ2MWYtYTdjNC04YTM2OGUwMzQwYjMiLCJ1dGkiOiI5TXFOSWI5WjdrQy1QVHRtai11X0FBIiwidmVyIjoiMi4wIn0.hh5Exz9MBjTXrTuTZnz7vceiuQjcC_oRSTeBIC9tYgSO2c2sqQRpZi91qBZFQD9okayLPPKcwqXgEJD9p0-c4nUR5UQN7YSeDLmYtZUYMG79EsA7IMiQaiy94AyIe2E-oBDcLwFycGwh1iIOwwOwjbanmu2Dx3HfQx831lH9uVjagf0Aow0wTkTVCsedGSZvG-cRUceFLj-kFN-feFH3NuScuOfLR2Magf541pJda7X7oStwL_RNUFqjJFTdsiFV4e-VHK5qo--3oPU06z0rS9bosj0pFSATIVHrrS4gY7jiSvgMbG837CDBQkz5b08GUN5GlLN9jlygl1plBmbgww"
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
		if err := req.ParseForm(); err != nil {
			return nil, fmt.Errorf("mockSTS failed to parse a request body: %w", err)
		}
		if grant := req.FormValue("grant_type"); grant == "device_code" || grant == "password" {
			// include account info because we're authenticating a user
			res.Body = io.NopCloser(bytes.NewReader(
				[]byte(fmt.Sprintf(`{"access_token":"at","expires_in": 3600,"refresh_token":"rt","client_info":%q,"id_token":%q}`, mockClientInfo, mockIDT)),
			))
		} else {
			res.Body = io.NopCloser(bytes.NewReader(accessTokenRespSuccess))
		}
	default:
		// User realm metadata request paths look like "/common/UserRealm/user@domain".
		// Matching on the UserRealm segment avoids having to know the UPN.
		if s[len(s)-2] == "UserRealm" {
			res.Body = io.NopCloser(
				strings.NewReader(`{"account_type":"Managed","cloud_audience_urn":"urn","cloud_instance_name":"...","domain_name":"..."}`),
			)
		} else {
			return nil, fmt.Errorf("mockSTS received an unexpected request for %s", req.URL.String())
		}
	}
	return &res, nil
}
