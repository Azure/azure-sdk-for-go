// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

type SessionCredentials = generated.SessionCredentials

type Provider interface {
	GetSessionCredentials(ctx context.Context, containerURL string) (SessionCredentials, error)
	RefreshSessionCredentials(ctx context.Context, containerURL string) (SessionCredentials, error)
}

type SessionPolicy struct {
	bearerTokenPolicy policy.Policy
	opts              SessionOptions
	provider          Provider
}

func NewSessionPolicy(opts SessionOptions, bearerTokenPolicy policy.Policy, oauthServiceClient *generated.ServiceClient) (policy.Policy, error) {
	var provider Provider
	switch opts.Mode {
	case SessionModeSingleContainer:
		if opts.AccountName == "" {
			return nil, errors.New("account name is required for singlecontainer mode")
		}
		if opts.ContainerName == "" {
			return nil, errors.New("container name is required for singlecontainer mode")
		}
		provider = NewSingleContainerProvider(oauthServiceClient, opts.ContainerName)
	default:
		return nil, fmt.Errorf("unsupported session mode %v", opts.Mode)
	}

	return &SessionPolicy{
		bearerTokenPolicy: bearerTokenPolicy,
		opts:              opts,
		provider:          provider,
	}, nil
}

// Do implements the policy.SessionPolicy interface.
func (p *SessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Look at request URL - if its a blob API, try to get a session token for the container and apply it to the request.
	if containerName, ok := parseBlobURL(req.Raw().URL.String()); ok {
		// Get session token for the container and apply it to the request
		sessionCreds, err := p.provider.GetSessionCredentials(req.Raw().Context(), containerName)
		if err == nil {
			var key, token string
			if sessionCreds.SessionKey != nil {
				key = *sessionCreds.SessionKey
			}
			if sessionCreds.SessionToken != nil {
				token = *sessionCreds.SessionToken
			}
			cred, err := NewSharedKeyCredential(p.opts.AccountName, key)
			if err != nil {
				return nil, err
			}

			if d := getHeader(shared.HeaderXmsDate, req.Raw().Header); d == "" {
				req.Raw().Header.Set(shared.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
			}
			stringToSign, err := cred.buildStringToSign(req.Raw())
			if err != nil {
				return nil, err
			}
			signature, err := cred.computeHMACSHA256(stringToSign)
			if err != nil {
				return nil, err
			}
			authHeader := strings.Join([]string{"Session ", token, ":", signature}, "")
			req.Raw().Header.Set(shared.HeaderAuthorization, authHeader)

			response, err := req.Next()
			if err != nil && response != nil && response.StatusCode == http.StatusForbidden {
				// Service failed to authenticate request, log it
				log.Write(azlog.EventResponse, "===== HTTP Forbidden status, String-to-Sign:\n"+stringToSign+"\n===============================\n")
			}

		} else if !errors.Is(err, ErrFallbackToBearer) {
			return nil, err
		}
	}

	// Fall back to bearer token policy
	return p.bearerTokenPolicy.Do(req)
}

// parseBlobURL parses a blob URL and returns the container name if it's a valid blob URL.
// Returns the container name and true if the URL is a blob URL, empty string and false otherwise.
func parseBlobURL(rawURL string) (containerName string, isBlobURL bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}

	// Path format: /<container>/<blob> or /<container>
	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		return "", false
	}

	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", false
	}

	containerName = parts[0]
	// It's a blob URL if there's a blob path after the container
	isBlobURL = len(parts) > 1 && parts[1] != ""

	return containerName, isBlobURL
}
