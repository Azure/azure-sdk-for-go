// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/x509"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
	cred, err := NewClientCertificateCredential(badTenantID, fakeClientID, test.certs, test.key, nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestClientCertificateCredential_GetTokenSuccess(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			cred, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, test.key, nil)
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
			cred.client = fakeConfidentialClient{}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
		})
	}
}

func TestClientCertificateCredential_GetTokenSuccess_withCertificateChain(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			options := ClientCertificateCredentialOptions{SendCertificateChain: true}
			cred, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, test.key, &options)
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
			cred.client = fakeConfidentialClient{}
			_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
			if err != nil {
				t.Fatalf("Expected an empty error but received: %s", err.Error())
			}
		})
	}
}

func TestClientCertificateCredential_GetTokenCheckPrivateKeyBlocks(t *testing.T) {
	test := allCertTests[0]
	cred, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, test.key, nil)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	cred.client = fakeConfidentialClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientCertificateCredential_NoData(t *testing.T) {
	var key crypto.PrivateKey
	_, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, []*x509.Certificate{}, key, nil)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_NoCertificate(t *testing.T) {
	test := allCertTests[0]
	_, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, []*x509.Certificate{}, test.key, nil)
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
	_, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, key, &options)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_Live(t *testing.T) {
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
			o, stop := initRecording(t)
			defer stop()
			opts := &ClientCertificateCredentialOptions{SendCertificateChain: test.sendChain, ClientOptions: o}
			cred, err := NewClientCertificateCredential(liveSP.tenantID, liveSP.clientID, certs, key, opts)
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestClientCertificateCredential_InvalidCertLive(t *testing.T) {
	test := allCertTests[0]
	o, stop := initRecording(t)
	defer stop()
	opts := &ClientCertificateCredentialOptions{ClientOptions: o}
	cred, err := NewClientCertificateCredential(liveSP.tenantID, liveSP.clientID, test.certs, test.key, opts)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}

	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if tk != nil {
		t.Fatal("GetToken returned a token")
	}
	var e AuthenticationFailedError
	if !errors.As(err, &e) {
		t.Fatal("expected AuthenticationFailedError")
	}
	if e.RawResponse() == nil {
		t.Fatal("expected RawResponse() to return a non-nil *http.Response")
	}
}
