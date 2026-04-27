// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

const featureNotEnabled = "FeatureNotEnabled"

type sessionCredentials struct {
	token  string
	key    string
	expiry time.Time
	// fallback indicates that session creation failed and the caller should use bearer token
	// authentication instead. This is stored as a field rather than returned as an error because
	// temporal.Resource only caches successful (non-error) results. Returning a non-error fallback
	// value allows the decision to be cached for the duration of errorExpiry, avoiding repeated
	// session creation attempts when the service indicates the feature is unavailable.
	fallback bool
}

// acquireSession is the function called by temporal.Resource to create a new session.
func acquireSession(client *generated.ContainerClient) func(context.Context) (sessionCredentials, time.Time, error) {
	return func(ctx context.Context) (creds sessionCredentials, expiry time.Time, err error) {
		resp, err := client.CreateSession(ctx, generated.CreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHMAC)}, nil)
		// Fall back to using bearer token if session is unable to be created
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				errorExpiry := time.Now().Add(5 * time.Minute)
				if respErr.StatusCode >= 500 {
					return sessionCredentials{fallback: true}, errorExpiry, nil
				}
				if respErr.StatusCode == http.StatusBadRequest && respErr.ErrorCode == featureNotEnabled {
					return sessionCredentials{fallback: true}, errorExpiry, nil
				}
				if respErr.StatusCode == http.StatusForbidden {
					return sessionCredentials{fallback: true}, errorExpiry, nil
				}
			}
			return creds, expiry, err
		}

		if resp.Expiration != nil {
			expiry = *resp.Expiration
		}
		var token, key string
		if resp.Credentials != nil {
			if resp.Credentials.SessionToken != nil {
				token = *resp.Credentials.SessionToken
			}
			if resp.Credentials.SessionKey != nil {
				key = *resp.Credentials.SessionKey
			}
		}

		return sessionCredentials{
			token:  token,
			key:    key,
			expiry: expiry,
		}, expiry, err
	}
}

func shouldRefreshSession(resource sessionCredentials, _ context.Context) bool {
	// call time.Now() instead of using Get's value so ShouldRefresh doesn't need a time.Time parameter
	return resource.expiry.Add(-30 * time.Second).Before(time.Now())
}

func getContainerClient(client *generated.ServiceClient, containerName string) *generated.ContainerClient {
	containerURL := runtime.JoinPaths(client.Endpoint(), containerName)
	return generated.NewContainerClient(containerURL, client.InternalClient())
}
