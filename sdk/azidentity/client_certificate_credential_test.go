// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type certTest struct {
	name  string
	certs []*x509.Certificate
	key   crypto.PrivateKey
}

func newCertTest(name, certPath string, password string) certTest {
	data, _ := os.ReadFile(certPath)
	certs, key, err := ParseCertificates(data, []byte(password))
	if err != nil {
		log.Panicf("failed to parse %s: %v", certPath, err)
	}
	return certTest{name: name, certs: certs, key: key}
}

var allCertTests = []certTest{
	newCertTest("pem", "testdata/certificate.pem", ""),
	newCertTest("pemB", "testdata/certificate_formatB.pem", ""),
	newCertTest("pkcs12", "testdata/certificate.pfx", ""),
	newCertTest("pkcs12Encrypted", "testdata/certificate_encrypted_key.pfx", "password"),
}

func TestClientCertificateCredential_InvalidTenantID(t *testing.T) {
	test := allCertTests[0]
	cred, err := NewClientCertificateCredential(badTenantID, clientID, test.certs, test.key, nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestClientCertificateCredential_CreateAuthRequestSuccess(t *testing.T) {
	test := allCertTests[0]
	cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.cert, false, []string{scope})
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
	if req.Raw().URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Raw().URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientCertificateCredential_CreateAuthRequestSuccess_withCertificateChain(t *testing.T) {
	test := allCertTests[0]
	opts := ClientCertificateCredentialOptions{}
	opts.SendCertificateChain = true
	cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &opts)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(context.Background(), cred.tenantID, cred.clientID, cred.cert, cred.sendCertificateChain, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Raw().Header.Get(headerContentType) != headerURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	if len(cred.cert.x5c) != 1 {
		t.Fatalf("Wrong number of public certificates. Expected: %v, Received: %v", 1, len(cred.cert.x5c))
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
	if reqQueryParams[qpGrantType][0] != "client_credentials" {
		t.Fatalf("Wrong grant type in request body")
	}
	if reqQueryParams[qpClientAssertionType][0] != clientAssertionType {
		t.Fatalf("Wrong client assertion type assigned to request")
	}
	// create a client assertion for comparison with the one in the request
	cc, err := newCertContents(test.certs, test.key.(*rsa.PrivateKey), opts.SendCertificateChain)
	if err != nil {
		t.Fatalf("%v", err)
	}
	assertion, err := createClientAssertionJWT(clientID, runtime.JoinPaths(string(AzurePublicCloud), tenantID, tokenEndpoint(oauthPath(tenantID))), cc, true)
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
	if req.Raw().URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Raw().URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientCertificateCredential_GetTokenSuccess(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
			options := ClientCertificateCredentialOptions{}
			options.AuthorityHost = AuthorityHost(srv.URL())
			options.Transport = srv
			cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &options)
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
		})
	}
}

func TestClientCertificateCredential_GetTokenSuccess_withCertificateChain(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
			options := ClientCertificateCredentialOptions{}
			options.AuthorityHost = AuthorityHost(srv.URL())
			options.SendCertificateChain = true
			options.Transport = srv
			cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &options)
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
		})
	}
}

func TestClientCertificateCredential_GetTokenInvalidCredentials(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
			options := ClientCertificateCredentialOptions{}
			options.AuthorityHost = AuthorityHost(srv.URL())
			options.Transport = srv
			cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &options)
			if err != nil {
				t.Fatalf("Did not expect an error but received one: %v", err)
			}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
			if err == nil {
				t.Fatalf("Expected to receive a nil error, but received: %v", err)
			}
			var authFailed AuthenticationFailedError
			if !errors.As(err, &authFailed) {
				t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
			}
		})
	}
}

func TestClientCertificateCredential_GetTokenCheckPrivateKeyBlocks(t *testing.T) {
	test := allCertTests[0]
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_NoData(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	var key crypto.PrivateKey
	_, err := NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{}, key, &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_NoCertificate(t *testing.T) {
	test := allCertTests[0]
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	_, err := NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{}, test.key, &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_NoPrivateKey(t *testing.T) {
	test := allCertTests[0]
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	var key crypto.PrivateKey
	_, err := NewClientCertificateCredential(tenantID, clientID, test.certs, key, &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestBearerPolicy_ClientCertificateCredential(t *testing.T) {
	test := allCertTests[0]
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientCertificateCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientCertificateCredential(tenantID, clientID, test.certs, test.key, &options)
	if err != nil {
		t.Fatalf("Did not expect an error but received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
}

func TestClientCertificateCredential_Live(t *testing.T) {
	if liveSP.clientID == "" || liveSP.tenantID == "" {
		t.Skip("missing live service principal configuration")
	}
	tests := []struct {
		name      string
		path      string
		sendChain bool
	}{
		{"PEM", liveSP.pemPath, false}, {"PKCS12", liveSP.pfxPath, false}, {"SNI", liveSP.sniPath, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.path == "" {
				t.Skip("no certificate file specified")
			}
			certData, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatalf(`failed to read cert: %v`, err)
			}
			certs, key, err := ParseCertificates(certData, nil)
			if err != nil {
				t.Fatalf(`failed to parse cert: %v`, err)
			}
			opts := &ClientCertificateCredentialOptions{SendCertificateChain: test.sendChain}
			cred, err := NewClientCertificateCredential(liveSP.tenantID, liveSP.clientID, certs, key, opts)
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
		})
	}
}
