//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/x509"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
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
	newCertTest("pemChain", "testdata/certificate-with-chain.pem", ""),
	newCertTest("pkcs12", "testdata/certificate.pfx", ""),
	newCertTest("pkcs12Encrypted", "testdata/certificate_encrypted_key.pfx", "password"),
}

func TestParseCertificates_Error(t *testing.T) {
	for _, path := range []string{
		"testdata/certificate_empty.pem",         // malformed file (no cert block)
		"testdata/certificate_encrypted_key.pfx", // requires a password we won't provide
		"testdata/certificate_nokey.pem",
		"testdata/certificate-two-keys.pem",
	} {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			_, _, err = ParseCertificates(data, nil)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestClientCertificateCredential_SendCertificateChain(t *testing.T) {
	for _, test := range allCertTests {
		t.Run(test.name, func(t *testing.T) {
			options := ClientCertificateCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: &mockSTS{tokenRequestCallback: validateX5C(t, test.certs)},
				},
				SendCertificateChain: true,
			}
			cred, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, test.key, &options)
			if err != nil {
				t.Fatal(err)
			}
			tk, err := cred.GetToken(context.Background(), testTRO)
			if err != nil {
				t.Fatal(err)
			}
			if tk.Token != tokenValue {
				t.Fatalf("unexpected token: %s", tk.Token)
			}
		})
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
	var key crypto.PrivateKey
	_, err := NewClientCertificateCredential(fakeTenantID, fakeClientID, test.certs, key, nil)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
}

func TestClientCertificateCredential_WrongKey(t *testing.T) {
	data, err := os.ReadFile("testdata/certificate-wrong-key.pem")
	if err != nil {
		t.Fatal(err)
	}
	certs, key, err := ParseCertificates(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewClientCertificateCredential("tenantID", "clientID", certs, key, nil)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientCertificateCredential_Live(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
	tests := []struct {
		name      string
		path      string
		sendChain bool
	}{
		{"PEM", liveSP.pemPath, false}, {"PKCS12", liveSP.pfxPath, false},
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
	t.Run("instance discovery disabled", func(t *testing.T) {
		if liveSP.pemPath == "" {
			t.Skip("no certificate file specified")
		}
		certData, err := os.ReadFile(liveSP.pemPath)
		if err != nil {
			t.Fatalf(`failed to read cert: %v`, err)
		}
		certs, key, err := ParseCertificates(certData, nil)
		if err != nil {
			t.Fatalf(`failed to parse cert: %v`, err)
		}
		o, stop := initRecording(t)
		defer stop()
		opts := &ClientCertificateCredentialOptions{ClientOptions: o, DisableInstanceDiscovery: true}
		cred, err := NewClientCertificateCredential(liveSP.tenantID, liveSP.clientID, certs, key, opts)
		if err != nil {
			t.Fatalf("failed to construct credential: %v", err)
		}
		testGetTokenSuccess(t, cred)
	})
}

func TestClientCertificateCredentialADFS_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveSP.clientID == "" || adfsLiveSP.certPath == "" || adfsScope == "" {
			t.Skip("set ADFS_SP_* to run this test live")
		}
	}
	certData, err := os.ReadFile(adfsLiveSP.certPath)
	if err != nil {
		t.Fatalf(`failed to read cert: %v`, err)
	}
	certs, key, err := ParseCertificates(certData, nil)
	if err != nil {
		t.Fatalf(`failed to parse cert: %v`, err)
	}
	o, stop := initRecording(t)
	defer stop()
	o.Cloud.ActiveDirectoryAuthorityHost = adfsAuthority
	opts := &ClientCertificateCredentialOptions{ClientOptions: o, DisableInstanceDiscovery: true}
	cred, err := NewClientCertificateCredential("adfs", adfsLiveSP.clientID, certs, key, opts)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
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

	tk, err := cred.GetToken(context.Background(), testTRO)
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.HasPrefix(err.Error(), credNameCert) {
		t.Fatalf("error is missing credential type prefix: %q", err.Error())
	}
}

func TestClientCertificateCredential_Regional(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		t.Skip("this test runs only in playback mode because it requires a production cert")
	}
	t.Setenv(azureRegionalAuthorityName, "westus2")
	opts, stop := initRecording(t)
	defer stop()

	f, err := os.ReadFile("testdata/certificate-with-chain.pem")
	if err != nil {
		t.Fatal(err)
	}
	cert, key, err := ParseCertificates(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	cred, err := NewClientCertificateCredential(
		liveSP.tenantID, liveSP.clientID, cert, key, &ClientCertificateCredentialOptions{SendCertificateChain: true, ClientOptions: opts},
	)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}
