// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

type emulatorCaptureTransport struct {
	request *http.Request
}

func (tr *emulatorCaptureTransport) Do(req *http.Request) (*http.Response, error) {
	tr.request = req.Clone(req.Context())
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(apiVersionTopicXML)),
		Request:    req,
		Header:     http.Header{"Content-Type": []string{"application/atom+xml"}},
	}, nil
}

func requireEmulatorSAS(t *testing.T, req *http.Request) {
	t.Helper()
	const key = "test-key"
	header := req.Header.Get("Authorization")
	require.True(t, strings.HasPrefix(header, "SharedAccessSignature "))
	values, err := url.ParseQuery(strings.TrimPrefix(header, "SharedAccessSignature "))
	require.NoError(t, err)
	require.Equal(t, strings.ToLower(req.URL.String()), values.Get("sr"))
	require.Equal(t, "test-key-name", values.Get("skn"))
	expiry, err := strconv.ParseInt(values.Get("se"), 10, 64)
	require.NoError(t, err)
	require.Greater(t, expiry, time.Now().Unix())
	mac := hmac.New(sha256.New, []byte(key))
	_, err = mac.Write([]byte(strings.ToLower(url.QueryEscape(req.URL.String())) + "\n" + values.Get("se")))
	require.NoError(t, err)
	require.Equal(t, base64.StdEncoding.EncodeToString(mac.Sum(nil)), values.Get("sig"))
}

func TestClientEmulatorConnectionString(t *testing.T) {
	for _, test := range []struct {
		name      string
		namespace string
		flag      string
		scheme    string
	}{
		{name: "emulator", namespace: "localhost:5300", flag: ";UseDevelopmentEmulator=true", scheme: "http"},
		{name: "dynamic-port", namespace: "127.0.0.1:15300", flag: ";UseDevelopmentEmulator=true", scheme: "http"},
		{name: "IPv6", namespace: "[::1]:5300", flag: ";UseDevelopmentEmulator=true", scheme: "http"},
		{name: "container-host", namespace: "servicebus-emulator:5300", flag: ";UseDevelopmentEmulator=true", scheme: "http"},
		{name: "default-http-port", namespace: "localhost", flag: ";UseDevelopmentEmulator=true", scheme: "http"},
		{name: "boolean-case", namespace: "localhost:5300", flag: ";UseDevelopmentEmulator=TRUE", scheme: "http"},
		{name: "absent", namespace: "localhost:5300", scheme: "https"},
		{name: "false", namespace: "localhost:5300", flag: ";UseDevelopmentEmulator=false", scheme: "https"},
		{name: "azure-default", namespace: "example.servicebus.windows.net", scheme: "https"},
	} {
		t.Run(test.name, func(t *testing.T) {
			transport := &emulatorCaptureTransport{}
			client, err := NewClientFromConnectionString(
				"Endpoint=sb://"+test.namespace+"/;SharedAccessKeyName=test-key-name;SharedAccessKey=test-key"+test.flag,
				&ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}},
			)
			require.NoError(t, err)
			_, err = client.GetTopicRuntimeProperties(context.Background(), "my-topic", nil)
			require.NoError(t, err)
			require.NotNil(t, transport.request)
			require.Equal(t, test.scheme+"://"+test.namespace+"/my-topic?api-version=2024-05", transport.request.URL.String())
			requireEmulatorSAS(t, transport.request)
		})
	}
}

func TestClientEmulatorDefaultTransport(t *testing.T) {
	requests := make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requests <- req.Clone(req.Context())
		w.Header().Set("Content-Type", "application/atom+xml")
		_, _ = io.WriteString(w, apiVersionTopicXML)
	}))
	defer server.Close()
	client, err := NewClientFromConnectionString(
		"Endpoint=sb://"+strings.TrimPrefix(server.URL, "http://")+
			";SharedAccessKeyName=test-key-name;SharedAccessKey=test-key;UseDevelopmentEmulator=true",
		nil,
	)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.GetTopicRuntimeProperties(ctx, "my-topic", nil)
	require.NoError(t, err)
	require.Equal(t, "my-topic", response.TopicName)
	select {
	case req := <-requests:
		require.Equal(t, "HTTP/1.1", req.Proto)
		require.Equal(t, strings.TrimPrefix(server.URL, "http://"), req.Host)
		require.Nil(t, req.TLS)
	case <-ctx.Done():
		t.Fatal("No native HTTP request reached the server")
	}
}

func TestClientEmulatorPrecomputedSAS(t *testing.T) {
	const token = "SharedAccessSignature sr=http%3a%2f%2flocalhost%3a5300&sig=test-signature&se=4102444800&skn=test-key-name"
	transport := &emulatorCaptureTransport{}
	client, err := NewClientFromConnectionString(
		"Endpoint=sb://localhost:5300;SharedAccessSignature="+token+";UseDevelopmentEmulator=true",
		&ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}},
	)
	require.NoError(t, err)
	_, err = client.GetTopicRuntimeProperties(context.Background(), "my-topic", nil)
	require.NoError(t, err)
	require.Equal(t, "http", transport.request.URL.Scheme)
	require.Equal(t, token, transport.request.Header.Get("Authorization"))
}

type emulatorTokenCredential struct {
	scopes []string
}

func (c *emulatorTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.scopes = options.Scopes
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestClientTokenCredentialRequiresHTTPS(t *testing.T) {
	transport := &emulatorCaptureTransport{}
	credential := &emulatorTokenCredential{}
	client, err := NewClient("localhost:5300", credential,
		&ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
	require.NoError(t, err)
	_, err = client.GetTopicRuntimeProperties(context.Background(), "my-topic", nil)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:5300/my-topic?api-version=2024-05", transport.request.URL.String())
	require.Equal(t, []string{"https://servicebus.azure.net//.default"}, credential.scopes)
	require.Equal(t, "test-token", transport.request.Header.Get("Authorization"))
}

func TestClientEmulatorInvalidConnectionString(t *testing.T) {
	for _, connectionString := range []string{
		"Endpoint=sb://localhost:5300;SharedAccessKeyName=test-key-name;SharedAccessKey=test-key;UseDevelopmentEmulator=invalid",
		"Endpoint=https://localhost:5300;SharedAccessKeyName=test-key-name;SharedAccessKey=test-key;UseDevelopmentEmulator=true",
	} {
		_, err := NewClientFromConnectionString(connectionString, nil)
		require.Error(t, err)
	}
}
