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
	testAuthCode    = "12345"
	testRedirectURI = "http://localhost"
)

func TestAuthorizationCodeCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewAuthorizationCodeCredential(badTenantID, clientID, testAuthCode, testRedirectURI, nil)
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

func TestAuthorizationCodeCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewAuthorizationCodeCredential(tenantID, clientID, testAuthCode, testRedirectURI, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := cred.client.createAuthorizationCodeAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.authCode, cred.clientSecret, "", cred.redirectURI, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(headerContentType) != headerURLEncoded {
		t.Fatal("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatal("Unable to read request body")
	}
	bodyStr := string(body)
	reqQueryParams, err := url.ParseQuery(bodyStr)
	if err != nil {
		t.Fatal("Unable to parse query params in request")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatal("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpScope][0] != scope {
		t.Fatal("Unexpected scope in scope header")
	}
	if reqQueryParams[qpCode][0] != testAuthCode {
		t.Fatal("Unexpected authorization code")
	}
	if reqQueryParams[qpRedirectURI][0] != testRedirectURI {
		t.Fatal("Unexpected redirectURI")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatal("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatal("Wrong request scheme")
	}
}

func TestAuthorizationCodeCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := AuthorizationCodeCredentialOptions{}
	options.ClientSecret = secret
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewAuthorizationCodeCredential(tenantID, clientID, testAuthCode, testRedirectURI, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
}

func TestAuthorizationCodeCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(accessTokenRespError)), mock.WithStatusCode(http.StatusUnauthorized))
	options := AuthorizationCodeCredentialOptions{}
	options.ClientSecret = secret
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewAuthorizationCodeCredential(tenantID, clientID, testAuthCode, testRedirectURI, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	} else {
		var respError *AADAuthenticationFailedError
		if !errors.As(authFailed.Unwrap(), &respError) {
			t.Fatalf("Expected: AADAuthenticationFailedError, Received: %T", err)
		} else {
			if len(respError.Message) == 0 {
				t.Fatalf("Did not receive an error message")
			}
			if len(respError.Description) == 0 {
				t.Fatalf("Did not receive an error description")
			}
			if len(respError.Timestamp) == 0 {
				t.Fatalf("Did not receive a timestamp")
			}
			if len(respError.TraceID) == 0 {
				t.Fatalf("Did not receive a TraceID")
			}
			if len(respError.CorrelationID) == 0 {
				t.Fatalf("Did not receive a CorrelationID")
			}
			if len(respError.URL) == 0 {
				t.Fatalf("Did not receive an error URL")
			}
			if respError.Response == nil {
				t.Fatalf("Did not receive an error response")
			}
		}
	}
}

func TestAuthorizationCodeCredential_GetTokenUnexpectedJSON(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	options := AuthorizationCodeCredentialOptions{}
	options.ClientSecret = secret
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewAuthorizationCodeCredential(tenantID, clientID, testAuthCode, testRedirectURI, &options)
	if err != nil {
		t.Fatalf("Failed to create the credential")
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}
