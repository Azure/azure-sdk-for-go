// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	deviceCode                   = "device_code"
	deviceCodeResponse           = `{"user_code":"test_code","device_code":"test_device_code","verification_uri":"https://microsoft.com/devicelogin","expires_in":900,"interval":5,"message":"To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code test_code to authenticate."}`
	deviceCodeScopes             = "user.read offline_access openid profile email"
	authorizationPendingResponse = `{"error": "authorization_pending","error_description": "Authorization pending.","error_codes": [],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=0"}`
	expiredTokenResponse         = `{"error": "expired_token","error_description": "Token has expired.","error_codes": [],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=0"}`
)

func TestDeviceCodeCredential_InvalidTenantID(t *testing.T) {
	options := DeviceCodeCredentialOptions{}
	options.TenantID = badTenantID
	cred, err := NewDeviceCodeCredential(&options)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
	var errType *CredentialUnavailableError
	if !errors.As(err, &errType) {
		t.Fatalf("Did not receive a CredentialUnavailableError. Received: %t", err)
	}
}

func TestDeviceCodeCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewDeviceCodeCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createDeviceCodeAuthRequest(context.Background(), cred.tenantID, cred.clientID, deviceCode, []string{deviceCodeScopes})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatalf("Unable to read request body")
	}
	bodyStr := string(body)
	reqQueryParams, err := url.ParseQuery(bodyStr)
	if err != nil {
		t.Fatalf("Unable to parse query params in request")
	}
	if reqQueryParams[qpGrantType][0] != deviceCodeGrantType {
		t.Fatalf("Unexpected grant type")
	}
	if reqQueryParams[qpClientID][0] != developerSignOnClientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpDeviceCode][0] != deviceCode {
		t.Fatalf("Unexpected username in the username header")
	}
	if reqQueryParams[qpScope][0] != deviceCodeScopes {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestDeviceCodeCredential_CreateAuthRequestCustomClientID(t *testing.T) {
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createDeviceCodeAuthRequest(context.Background(), cred.tenantID, cred.clientID, deviceCode, []string{deviceCodeScopes})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatalf("Unable to read request body")
	}
	bodyStr := string(body)
	reqQueryParams, err := url.ParseQuery(bodyStr)
	if err != nil {
		t.Fatalf("Unable to parse query params in request")
	}
	if reqQueryParams[qpGrantType][0] != deviceCodeGrantType {
		t.Fatalf("Unexpected grant type")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpDeviceCode][0] != deviceCode {
		t.Fatalf("Unexpected username in the username header")
	}
	if reqQueryParams[qpScope][0] != deviceCodeScopes {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
	if req.Request.URL.Path != "/organizations/oauth2/v2.0/token" {
		t.Fatalf("Did not set the right path when passing in an empty tenant ID")
	}
}

func TestDeviceCodeCredential_RequestNewDeviceCodeCustomTenantIDClientID(t *testing.T) {
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createDeviceCodeNumberRequest(context.Background(), cred.tenantID, cred.clientID, []string{deviceCodeScopes})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatalf("Unable to read request body")
	}
	bodyStr := string(body)
	reqQueryParams, err := url.ParseQuery(bodyStr)
	if err != nil {
		t.Fatalf("Unable to parse query params in request")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpScope][0] != deviceCodeScopes {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
	if req.Request.URL.Path != "/expected-tenant/oauth2/v2.0/devicecode" {
		t.Fatalf("Did not set the right path when passing in an empty tenant ID")
	}
}

func TestDeviceCodeCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := DeviceCodeCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	if tk.Token != "new_token" {
		t.Fatalf("Received an unexpected value in azcore.AccessToken.Token")
	}
}

func TestDeviceCodeCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestDeviceCodeCredential_GetTokenAuthorizationPending(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(authorizationPendingResponse)), mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithBody([]byte(authorizationPendingResponse)), mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	options.UserPrompt = func(DeviceCodeMessage) {}
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err != nil {
		t.Fatalf("Expected an empty error but received %v", err)
	}
}

func TestDeviceCodeCredential_GetTokenExpiredToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(authorizationPendingResponse)), mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithBody([]byte(expiredTokenResponse)), mock.WithStatusCode(http.StatusUnauthorized))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	options.UserPrompt = func(DeviceCodeMessage) {}
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err == nil {
		t.Fatalf("Expected an error but received none")
	}
}

func TestDeviceCodeCredential_GetTokenWithRefreshTokenFailure(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespError)), mock.WithStatusCode(http.StatusUnauthorized))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.refreshToken = "refresh_token"
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var aadErr *AADAuthenticationFailedError
	if !errors.As(err, &aadErr) {
		t.Fatalf("Did not receive an AADAuthenticationFailedError but was expecting one")
	}
}

func TestDeviceCodeCredential_GetTokenWithRefreshTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	options.UserPrompt = func(DeviceCodeMessage) {}
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.refreshToken = "refresh_token"
	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err.Error())
	}
	if tk.Token != "new_token" {
		t.Fatalf("Unexpected value for token")
	}
}

func TestBearerPolicy_DeviceCodeCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := DeviceCodeCredentialOptions{}
	options.ClientID = clientID
	options.TenantID = tenantID
	options.HTTPClient = srv
	options.AuthorityHost = srv.URL()
	options.UserPrompt = func(DeviceCodeMessage) {}
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, deviceCodeScopes)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}
