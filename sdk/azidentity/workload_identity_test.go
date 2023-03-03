//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func getAssertion(cert *x509.Certificate, key crypto.PrivateKey) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", liveSP.tenantID),
		"exp": json.Number(strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)),
		"iss": liveSP.clientID,
		"jti": uuid.New().String(),
		"nbf": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		"sub": liveSP.clientID,
	})
	x5t := sha1.Sum(cert.Raw) // nosec
	j.Header = map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
		"x5t": base64.StdEncoding.EncodeToString(x5t[:]),
	}
	return j.SignedString(key)
}

func TestWorkloadIdentityCredential_Live(t *testing.T) {
	cert, err := os.ReadFile(liveSP.pemPath)
	if err != nil {
		t.Fatal(err)
	}
	certs, key, err := ParseCertificates(cert, nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := getAssertion(certs[0], key)
	if err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(t.TempDir(), t.Name())
	if err := os.WriteFile(f, []byte(a), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	for _, b := range []bool{true, false} {
		name := "default options"
		if b {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			co, stop := initRecording(t)
			defer stop()
			o := WorkloadIdentityCredentialOptions{ClientOptions: co, DisableInstanceDiscovery: b}
			cred, err := NewWorkloadIdentityCredential(liveSP.tenantID, liveSP.clientID, f, &o)
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestWorkloadIdentityCredential(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	validateReq := func(req *http.Request) bool {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != tokenValue {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual, ok := req.PostForm["client_id"]; !ok {
			t.Error("expected a client_id")
		} else if len(actual) != 1 || actual[0] != fakeClientID {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual := strings.Split(req.URL.Path, "/")[1]; actual != fakeTenantID {
			t.Errorf(`unexpected tenant "%s"`, actual)
		}
		return true
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithBody(instanceDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(tenantDiscoveryResponse))
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()
	opts := WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: srv},
	}
	cred, err := NewWorkloadIdentityCredential(fakeTenantID, fakeClientID, tempFile, &opts)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}

func TestWorkloadIdentityCredential_Expiration(t *testing.T) {
	tokenReqs := 0
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	validateReq := func(req *http.Request) bool {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != fmt.Sprint(tokenReqs) {
			t.Errorf(`expected assertion "%d", got "%s"`, tokenReqs, actual[0])
		}
		tokenReqs++
		return true
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithBody(instanceDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(tenantDiscoveryResponse))
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()
	srv.AppendResponse(mock.WithPredicate(validateReq), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()
	opts := WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: srv},
	}
	cred, err := NewWorkloadIdentityCredential(fakeTenantID, fakeClientID, tempFile, &opts)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		// tokenReqs counts requests, and its latest value is the expected client assertion and the requested scope.
		// Each iteration of this loop therefore sends a token request with a unique assertion.
		s := fmt.Sprint(tokenReqs)
		if err = os.WriteFile(tempFile, []byte(fmt.Sprint(s)), os.ModePerm); err != nil {
			t.Fatalf("failed to write token file: %v", err)
		}
		if _, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{s}}); err != nil {
			t.Fatal(err)
		}
		cred.expires = time.Now().Add(-time.Second)
	}
	if tokenReqs != 2 {
		t.Fatalf("expected 2 token requests, got %d", tokenReqs)
	}
}
