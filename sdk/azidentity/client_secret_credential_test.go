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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	tenantID                 = "expected-tenant"
	badTenantID              = "bad_tenant"
	clientID                 = "expected-client-id"
	secret                   = "secret"
	wrongSecret              = "wrong_secret"
	scope                    = "https://storage.azure.com/.default"
	defaultTestAuthorityHost = "login.microsoftonline.com"
)

func TestClientSecretCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewClientSecretCredential(badTenantID, clientID, secret, nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestClientSecretCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createClientSecretAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.clientSecret, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Raw().Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Raw().Body)
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
	if reqQueryParams[qpClientSecret][0] != secret {
		t.Fatalf("Unexpected secret in the client_secret header")
	}
	if reqQueryParams[qpScope][0] != scope {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Raw().URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Raw().URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientSecretCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
}

func TestClientSecretCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(accessTokenRespError)), mock.WithStatusCode(http.StatusUnauthorized))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
	var authFailed AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
	if authFailed.RawResponse() == nil {
		t.Fatalf("Expected error to include a response")
	}
}

func TestClientSecretCredential_GetTokenUnexpectedJSON(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Failed to create the credential")
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}

func TestClientSecretCredential_Live(t *testing.T) {
	if liveSP.tenantID == "" || liveSP.clientID == "" || liveSP.secret == "" {
		t.Skip("missing live service principal configuration")
	}
	cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, liveSP.secret, nil)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if tk.Token == "" {
		t.Fatalf("GetToken returned an invalid token")
	}
	if !tk.ExpiresOn.After(time.Now().UTC()) {
		t.Fatalf("GetToken returned an invalid expiration time")
	}
}
