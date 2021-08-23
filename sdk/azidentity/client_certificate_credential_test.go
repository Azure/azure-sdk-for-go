// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	certificatePath      = "testdata/certificate.pem"
	wrongCertificatePath = "wrong_certificate_path.pem"
)

func TestClientCertificateCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewClientCertificateCredential(badTenantID, clientID, certificatePath, nil)
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

func TestClientCertificateCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.cert, false, []string{scope})
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

func TestClientCertificateCredential_CreateAuthRequestSuccess_withCertificateChain(t *testing.T) {
	opts := ClientCertificateCredentialOptions{}
	opts.SendCertificateChain = true
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &opts)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.cert, cred.sendCertificateChain, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	if len(cred.cert.publicCertificates) != 1 {
		t.Fatalf("Wrong number of public certificates. Expected: %v, Received: %v", 1, len(cred.cert.publicCertificates))
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
	// create a client assertion for comparison with the one in the request
	certData, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		t.Fatalf("Failed to read certificate: %v", err)
	}
	cert, err := extractFromPEMFile(certData, "", true)
	if err != nil {
		t.Fatalf("Failed extract data from PEM file: %v", err)
	}
	assertion, err := createClientAssertionJWT(clientID, azcore.JoinPaths(AzurePublicCloud, tenantID, tokenEndpoint(oauthPath(tenantID))), cert, true)
	if err != nil {
		t.Fatalf("Failed to create client assertion: %v", err)
	}
	// Get the index that separates the header of the JWT from the payload and signature.
	// NOTE: the payload and signature cannot be used for comparison since they incorporate
	// random numbers or unique IDs when being generated.
	i := strings.Index(assertion, ".")
	if i == -1 {
		t.Fatalf("malformed JWT")
	}
	if !strings.HasPrefix(reqQueryParams[qpClientAssertion][0], assertion[:i]) {
		t.Fatalf("Client assertion failed. Expected: %v, Received: %v", assertion, reqQueryParams[qpClientAssertion][0])
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
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenSuccess_withCertificateChain(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.SendCertificateChain = true
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &options)
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
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	_, err := NewClientCertificateCredential(tenantID, clientID, wrongCertificatePath, &options)
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestClientCertificateCredential_GetTokenCheckPrivateKeyBlocks(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_formatB.pem", &options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenCheckCertificateBlocks(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_formatA.pem", &options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_GetTokenEmptyCertificate(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	_, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_empty.pem", &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_GetTokenNoPrivateKey(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	_, err := NewClientCertificateCredential(tenantID, clientID, "testdata/certificate_nokey.pem", &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestBearerPolicy_ClientCertificateCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &options)
	if err != nil {
		t.Fatalf("Did not expect an error but received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
}
