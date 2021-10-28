//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type KeyVaultChallengePolicy struct {
	Cred        azcore.TokenCredential
	Transport   policy.Transporter
	cachedToken *azcore.AccessToken
}

func (k *KeyVaultChallengePolicy) Do(req *policy.Request) (*http.Response, error) {
	// if k.cachedToken != nil {
	// 	// cached token available
	// 	cushion := time.Now().Add(-60 * time.Second)
	// 	fmt.Println(cushion)
	// 	if !k.cachedToken.ExpiresOn.After(cushion) {
	// 		decorateRequest(req, k.cachedToken)
	// 	}
	// }

	challengeReq, err := getChallengeRequest(req)
	if err != nil {
		return nil, err
	}

	challengeResp, err := k.Transport.Do(challengeReq)
	if err != nil {
		return nil, err
	}

	_, scope, err := parseAuthChallenge(challengeResp)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(scope, "/") {
		scope += ".default"
	} else {
		scope += "/.default"
	}

	token, err := k.Cred.GetToken(req.Raw().Context(), policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		return nil, err
	}
	k.cachedToken = token
	decorateRequest(req, k.cachedToken)

	return req.Next()
}

func decorateRequest(req *policy.Request, token *azcore.AccessToken) {
	req.Raw().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
}

func getChallengeRequest(orig *policy.Request) (*http.Request, error) {
	req, err := http.NewRequest(orig.Raw().Method, orig.Raw().URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = orig.Raw().Header.Clone()
	req.Header.Set("Content-Length", "0")

	return req, err
}

func parseAuthChallenge(resp *http.Response) (string, string, error) {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return "", "", errors.New("response has no WWW-Authenticate header for challenge authentication")
	}

	// Strip down to auth and resource
	// Format is "Bearer authorization=\"<site>\" resource=\"<site>\""
	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("could not understand WWW-Authenticate header: %s", authHeader)
	}

	vals := make([]string, 0)
	for _, part := range parts {
		subParts := strings.Split(part, "=")
		if len(subParts) != 2 {
			return "", "", fmt.Errorf("could not understand WWW-Authenticate header: %s", authHeader)
		}
		url := subParts[1]

		url = strings.ReplaceAll(url, "\"", "")
		vals = append(vals, url)
	}

	return vals[0], vals[1], nil
}
