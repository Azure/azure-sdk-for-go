// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

func TestParseConnectionString(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringMixedOrder(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("Id=yY;Secret=ZmZm;Endpoint=xX")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringExtraProperties(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("A=aA;Endpoint=xX;B=bB;Id=yY;C=cC;Secret=ZmZm;D=dD;")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringMissingEndoint(t *testing.T) {
	_, _, _, err := ParseConnectionString("Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Endpoint")
}

func TestParseConnectionStringMissingId(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Id")
}

func TestParseConnectionStringMissingSecret(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Secret")
}

func TestParseConnectionStringDuplicateEndoint(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Endpoint=xX;Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Endpoint")
}

func TestParseConnectionStringDuplicateId(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Id")
}

func TestParseConnectionStringDuplicateSecret(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm;Secret=zZ")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Secret")
}

func TestParseConnectionStringInvalidEncoding(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=badencoding")
	require.Error(t, err)
	require.ErrorContains(t, err, "illegal base64 data")
}

func TestNewHMACPolicy(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")

	policy := NewHMACPolicy(credential, secret)

	require.NotNil(t, policy)
	require.Equal(t, credential, policy.credential)
	require.Equal(t, secret, policy.secret)
}

func TestNewHMACPolicyWithEmptyValues(t *testing.T) {
	policy := NewHMACPolicy("", nil)

	require.NotNil(t, policy)
	require.Equal(t, "", policy.credential)
	require.Nil(t, policy.secret)
}

func TestGetContentHashBase64Empty(t *testing.T) {
	hash, err := getContentHashBase64(nil)
	require.NoError(t, err)

	expected := sha256.Sum256([]byte{})
	expectedHash := base64.StdEncoding.EncodeToString(expected[:])
	require.Equal(t, expectedHash, hash)
}

func TestGetContentHashBase64WithContent(t *testing.T) {
	content := []byte("test content")
	hash, err := getContentHashBase64(content)
	require.NoError(t, err)

	expected := sha256.Sum256(content)
	expectedHash := base64.StdEncoding.EncodeToString(expected[:])
	require.Equal(t, expectedHash, hash)
}

func TestGetHMAC(t *testing.T) {
	content := "test content"
	key := []byte("test-key")

	signature, err := getHMAC(content, key)
	require.NoError(t, err)
	require.NotEmpty(t, signature)

	signature2, err := getHMAC(content, key)
	require.NoError(t, err)
	require.Equal(t, signature, signature2)
}

func TestGetHMACWithEmptyContent(t *testing.T) {
	content := ""
	key := []byte("test-key")

	signature, err := getHMAC(content, key)
	require.NoError(t, err)
	require.NotEmpty(t, signature)
}

func TestGetHMACWithEmptyKey(t *testing.T) {
	content := "test content"
	key := []byte{}

	signature, err := getHMAC(content, key)
	require.NoError(t, err)
	require.NotEmpty(t, signature)
}

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestHMACPolicyDoGETRequest(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}

	transport := &mockTransport{response: resp}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), "GET", "https://test.azconfig.io/kv/test-key")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
	require.Contains(t, authHeader, "SignedHeaders=date;host;x-ms-content-sha256")
	require.Contains(t, authHeader, "Signature=")

	require.NotEmpty(t, req.Raw().Header.Get("Date"))
	require.NotEmpty(t, req.Raw().Header.Get("x-ms-content-sha256"))
}

func TestHMACPolicyDoPOSTRequestWithBody(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}

	transport := &mockTransport{response: resp}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	body := `{"key": "test-key", "value": "test-value"}`
	req, err := runtime.NewRequest(context.Background(), "POST", "https://test.azconfig.io/kv")
	require.NoError(t, err)

	err = req.SetBody(streaming.NopCloser(strings.NewReader(body)), "application/json")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
	require.Contains(t, authHeader, "SignedHeaders=date;host;x-ms-content-sha256")
	require.Contains(t, authHeader, "Signature=")

	require.NotEmpty(t, req.Raw().Header.Get("Date"))
	require.NotEmpty(t, req.Raw().Header.Get("x-ms-content-sha256"))
}

func TestHMACPolicyDoRequestWithQueryParams(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}

	transport := &mockTransport{response: resp}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), "GET", "https://test.azconfig.io/kv?key=test*&label=production")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
	require.NotEmpty(t, req.Raw().Header.Get("Date"))
	require.NotEmpty(t, req.Raw().Header.Get("x-ms-content-sha256"))
}

func TestHMACPolicyDoWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			credential := "test-credential"
			secret := []byte("test-secret")
			hmacPolicy := NewHMACPolicy(credential, secret)

			resp := &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}

			transport := &mockTransport{response: resp}

			pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
				PerCall: []policy.Policy{hmacPolicy},
			}, &policy.ClientOptions{
				Transport: transport,
			})

			req, err := runtime.NewRequest(context.Background(), method, "https://test.azconfig.io/kv/test-key")
			require.NoError(t, err)

			response, err := pl.Do(req)
			require.NoError(t, err)
			require.NotNil(t, response)

			authHeader := req.Raw().Header.Get("Authorization")
			require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
			require.NotEmpty(t, req.Raw().Header.Get("Date"))
			require.NotEmpty(t, req.Raw().Header.Get("x-ms-content-sha256"))
		})
	}
}

func TestHMACPolicyDoTransportError(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	transport := &mockTransport{
		response: nil,
		err:      http.ErrServerClosed,
	}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), "GET", "https://test.azconfig.io/kv/test-key")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.Error(t, err)
	require.Equal(t, http.ErrServerClosed, err)
	require.Nil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
}

func TestHMACPolicyDoWithEmptyCredentials(t *testing.T) {
	credential := ""
	secret := []byte{}
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}

	transport := &mockTransport{response: resp}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), "GET", "https://test.azconfig.io/kv/test-key")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential=")
	require.Contains(t, authHeader, "SignedHeaders=date;host;x-ms-content-sha256")
	require.Contains(t, authHeader, "Signature=")
}

func TestHMACPolicyDoWithLargeBody(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}

	transport := &mockTransport{response: resp}

	pl := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	largeBody := strings.Repeat("test data ", 1000)
	req, err := runtime.NewRequest(context.Background(), "POST", "https://test.azconfig.io/kv")
	require.NoError(t, err)

	err = req.SetBody(streaming.NopCloser(strings.NewReader(largeBody)), "text/plain")
	require.NoError(t, err)

	response, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, response)

	authHeader := req.Raw().Header.Get("Authorization")
	require.Contains(t, authHeader, "HMAC-SHA256 Credential="+credential)
	require.NotEmpty(t, req.Raw().Header.Get("x-ms-content-sha256"))
}

func TestHMACPolicySignatureConsistency(t *testing.T) {
	credential := "test-credential"
	secret := []byte("test-secret")
	hmacPolicy := NewHMACPolicy(credential, secret)

	resp1 := &http.Response{StatusCode: 200, Header: make(http.Header), Body: http.NoBody}
	resp2 := &http.Response{StatusCode: 200, Header: make(http.Header), Body: http.NoBody}

	transport1 := &mockTransport{response: resp1}
	transport2 := &mockTransport{response: resp2}

	pl1 := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport1,
	})
	pl2 := runtime.NewPipeline("azappconfig", "v2.0.0", runtime.PipelineOptions{
		PerCall: []policy.Policy{hmacPolicy},
	}, &policy.ClientOptions{
		Transport: transport2,
	})

	req1, err := runtime.NewRequest(context.Background(), "POST", "https://test.azconfig.io/kv")
	require.NoError(t, err)
	req2, err := runtime.NewRequest(context.Background(), "POST", "https://test.azconfig.io/kv")
	require.NoError(t, err)

	body := `{"key": "test", "value": "test"}`
	err = req1.SetBody(streaming.NopCloser(strings.NewReader(body)), "application/json")
	require.NoError(t, err)
	err = req2.SetBody(streaming.NopCloser(strings.NewReader(body)), "application/json")
	require.NoError(t, err)

	timestamp := "Wed, 01 Jan 2025 00:00:00 GMT"
	req1.Raw().Header.Set("Date", timestamp)
	req2.Raw().Header.Set("Date", timestamp)

	_, err = pl1.Do(req1)
	require.NoError(t, err)
	_, err = pl2.Do(req2)
	require.NoError(t, err)

	auth1 := req1.Raw().Header.Get("Authorization")
	auth2 := req2.Raw().Header.Get("Authorization")

	require.Contains(t, auth1, "Signature=")
	require.Contains(t, auth2, "Signature=")

	require.Equal(t, auth1, auth2)
}
