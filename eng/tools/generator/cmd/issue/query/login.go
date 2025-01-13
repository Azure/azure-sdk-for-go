// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package query

import (
	"context"
	"log"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

// Login to github using the given credentials
func Login(ctx context.Context, info Info) *Client {
	var client *Client
	if info.UserInfo.IsValid() {
		client = NewClientWithUserInfo(info.UserInfo)
	} else if info.Token != "" {
		client = NewClientWithAccessToken(ctx, info.Token)
	} else {
		client = NewClient()
	}
	return client
}

func getGithubClientWithUserInfo(info UserInfo) *github.Client {
	log.Printf("Loging in with username and password")
	auth := &github.BasicAuthTransport{
		Username:  info.Username,
		Password:  info.Password,
		OTP:       info.Otp,
		Transport: nil,
	}
	return github.NewClient(auth.Client())
}

func getGithubClientWithAccessToken(ctx context.Context, token string) *github.Client {
	log.Printf("Loging in with personal access token")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// UserInfo ...
type UserInfo struct {
	Username string
	Password string
	Otp      string
}

func (u UserInfo) IsValid() bool {
	return u.Username != "" && u.Password != ""
}

// Info represents the login info
type Info struct {
	UserInfo UserInfo
	Token    string
}
