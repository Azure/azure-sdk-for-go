// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestUsernamePasswordCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(badTenantID, clientID, "username", "password", nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestUsernamePasswordCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(tenantID, clientID, "username", "password", nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createUsernamePasswordAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.username, cred.password, []string{scope})
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
	if reqQueryParams[qpResponseType][0] != "token" {
		t.Fatalf("Unexpected response type")
	}
	if reqQueryParams[qpGrantType][0] != "password" {
		t.Fatalf("Unexpected grant type")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpUsername][0] != "username" {
		t.Fatalf("Unexpected username in the username header")
	}
	if reqQueryParams[qpPassword][0] != "password" {
		t.Fatalf("Unexpected password in the password header")
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

func TestUsernamePasswordCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := UsernamePasswordCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewUsernamePasswordCredential(tenantID, clientID, "username", "password", &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestUsernamePasswordCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := UsernamePasswordCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewUsernamePasswordCredential(tenantID, clientID, "username", "wrong_password", &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestBearerPolicy_UsernamePasswordCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := UsernamePasswordCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewUsernamePasswordCredential(tenantID, clientID, "username", "password", &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}

func TestUsernamePasswordCredential_Live(t *testing.T) {
	username := os.Getenv("AZURE_IDENTITY_TEST_USERNAME")
	password := os.Getenv("AZURE_IDENTITY_TEST_PASSWORD")
	tenantID := os.Getenv("AZURE_IDENTITY_TEST_TENANTID")
	if username == "" || password == "" || tenantID == "" {
		t.Skip("no user configured")
	}
	cred, err := NewUsernamePasswordCredential(tenantID, developerSignOnClientID, username, password, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
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
