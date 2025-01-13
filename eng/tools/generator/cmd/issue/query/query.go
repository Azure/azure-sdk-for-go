// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package query

import (
	"context"

	"github.com/google/go-github/v62/github"
)

// Client ...
type Client struct {
	*github.Client
}

// NewClient returns a new Client without credential
func NewClient() *Client {
	return &Client{
		Client: github.NewClient(nil),
	}
}

// NewClientWithUserInfo ...
func NewClientWithUserInfo(info UserInfo) *Client {
	return &Client{
		Client: getGithubClientWithUserInfo(info),
	}
}

// NewClientWithAccessToken ...
func NewClientWithAccessToken(ctx context.Context, token string) *Client {
	return &Client{
		Client: getGithubClientWithAccessToken(ctx, token),
	}
}
