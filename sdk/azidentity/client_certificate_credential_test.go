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
	certificatePath      = "testdata/certificate.pem"
	wrongCertificatePath = "wrong_certificate_path.pem"
)

func TestClientCertificateCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(cred.tenantID, cred.clientID, cred.clientCertificate, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(azcore.HeaderContentType) != azcore.HeaderURLEncoded {
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
	if reqQueryParams[qpGrantType][0] != "client_credentials" {
		t.Fatalf("Wrong grant type in request body")
	}
	if reqQueryParams[qpClientAssertionType][0] != clientAssertionType {
		t.Fatalf("Wrong client assertion type assigned to request")
	}
	if reqQueryParams[qpScope][0] != scope {
		t.Fatalf("Unexpected scope in scope header")
	}
	if len(reqQueryParams[qpClientAssertion][0]) == 0 {
		t.Fatalf("Client assertion is not present on the request")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientCertificateCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Did not expect an error but received one: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected to receive a nil error, but received: %v", err)
	}
	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

func TestClientCertificateCredential_WrongCertificatePath(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	_, err := NewClientCertificateCredential(tenantID, clientID, wrongCertificatePath, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestClientCertificateCredential_CreateAuthRequestFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	srvURL.Host = "ht @"
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create a ClientCertificateCredential")
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestClientCertificateCredential_GetTokenCheckPrivateKeyBlocks(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_formatB.pem", &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenCheckCertificateBlocks(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_formatA.pem", &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenEmptyCertificate(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_empty.pem", &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_GetTokenNoPrivateKey(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_nokey.pem", &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestBearerPolicy_ClientCertificateCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Did not expect an error but received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req := pipeline.NewRequest(http.MethodGet, srv.URL())
	_, err = req.Do(context.Background())
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
}
