// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/stretchr/testify/require"
)

// newTestContainerClient creates a ContainerClient backed by a mock server for testing.
func newTestContainerClient(t *testing.T, srv *mock.Server) *generated.ContainerClient {
	azClient, err := azcore.NewClient("test", "v1.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	return generated.NewContainerClient(srv.URL()+"/testcontainer", azClient)
}

// createSessionResponseXML creates a session response XML body for testing.
func createSessionResponseXML(sessionKey, sessionToken string, expiration time.Time) []byte {
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<CreateSessionResponse>
	<AuthenticationType>HMAC</AuthenticationType>
	<Id>test-session-id</Id>
	<Credentials>
		<SessionKey>` + sessionKey + `</SessionKey>
		<SessionToken>` + sessionToken + `</SessionToken>
	</Credentials>
	<Expiration>` + expiration.Format(time.RFC1123) + `</Expiration>
</CreateSessionResponse>`)
}

// createErrorResponseXML creates an error response XML body for testing.
func createErrorResponseXML(code, message string) []byte {
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<Error>
	<Code>` + code + `</Code>
	<Message>` + message + `</Message>
</Error>`)
}

func TestAcquireSession_Success(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	expiration := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	srv.AppendResponse(
		mock.WithStatusCode(http.StatusCreated),
		mock.WithBody(createSessionResponseXML("test-key", "test-token", expiration)),
	)

	client := newTestContainerClient(t, srv)

	creds, exp, err := acquireSession(client)(context.Background())
	require.NoError(t, err)
	require.Equal(t, "test-key", creds.key)
	require.Equal(t, "test-token", creds.token)
	require.Equal(t, expiration.Format(time.RFC1123), exp.Format(time.RFC1123))
}

func TestAcquireSession_FallbackToBearer_Retryable(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "ServerError500",
			statusCode: http.StatusInternalServerError,
			errorCode:  "InternalError",
		},
		{
			name:       "ServiceUnavailable503",
			statusCode: http.StatusServiceUnavailable,
			errorCode:  "ServiceUnavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			srv.SetResponse(
				mock.WithStatusCode(tt.statusCode),
				mock.WithHeader("x-ms-error-code", tt.errorCode),
				mock.WithBody(createErrorResponseXML(tt.errorCode, "error message")),
			)

			client := newTestContainerClient(t, srv)

			before := time.Now()
			creds, exp, err := acquireSession(client)(context.Background())
			require.NoError(t, err)
			require.True(t, creds.fallback)
			// expiry should be ~5 minutes from now to cache the fallback decision
			require.True(t, exp.After(before.Add(4*time.Minute)))
			require.True(t, exp.Before(before.Add(6*time.Minute)))
		})
	}
}

func TestAcquireSession_FallbackToBearer_NonRetryable(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "FeatureNotEnabled",
			statusCode: http.StatusBadRequest,
			errorCode:  featureNotEnabled,
		},
		{
			name:       "Forbidden",
			statusCode: http.StatusForbidden,
			errorCode:  "AuthorizationFailure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeFn()

			srv.AppendResponse(
				mock.WithStatusCode(tt.statusCode),
				mock.WithHeader("x-ms-error-code", tt.errorCode),
				mock.WithBody(createErrorResponseXML(tt.errorCode, "error message")),
			)

			client := newTestContainerClient(t, srv)

			before := time.Now()
			creds, exp, err := acquireSession(client)(context.Background())
			require.NoError(t, err)
			require.True(t, creds.fallback)
			// expiry should be ~5 minutes from now to cache the fallback decision
			require.True(t, exp.After(before.Add(4*time.Minute)))
			require.True(t, exp.Before(before.Add(6*time.Minute)))
		})
	}
}

func TestAcquireSession_Error(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	srv.AppendResponse(
		mock.WithStatusCode(http.StatusNotFound),
		mock.WithHeader("x-ms-error-code", "ContainerNotFound"),
		mock.WithBody(createErrorResponseXML("ContainerNotFound", "Container not found")),
	)

	client := newTestContainerClient(t, srv)

	creds, _, err := acquireSession(client)(context.Background())
	require.Error(t, err)
	// Should NOT be a fallback - this is a real error that should propagate
	require.False(t, creds.fallback)

	var respErr *azcore.ResponseError
	require.True(t, errors.As(err, &respErr))
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
}

func TestShouldRefreshSession_NotExpiringSoon(t *testing.T) {
	creds := sessionCredentials{
		expiry: time.Now().Add(5 * time.Minute),
	}
	require.False(t, shouldRefreshSession(creds, context.Background()))
}

func TestShouldRefreshSession_ExpiringSoon(t *testing.T) {
	creds := sessionCredentials{
		expiry: time.Now().Add(10 * time.Second),
	}
	require.True(t, shouldRefreshSession(creds, context.Background()))
}

func TestShouldRefreshSession_AlreadyExpired(t *testing.T) {
	creds := sessionCredentials{
		expiry: time.Now().Add(-1 * time.Minute),
	}
	require.True(t, shouldRefreshSession(creds, context.Background()))
}

func TestGetContainerClient(t *testing.T) {
	srv, closeFn := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeFn()

	serviceClient := newTestServiceClient(t, srv)
	containerClient := getContainerClient(serviceClient, "mycontainer")
	require.NotNil(t, containerClient)
	require.Contains(t, containerClient.Endpoint(), "mycontainer")
}
