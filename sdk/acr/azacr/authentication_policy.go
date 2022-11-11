//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

const (
	headerAuthorization = "Authorization"
	bearerHeader        = "Bearer "
)

type AuthenticationPolicyOptions struct {
}

type AuthenticationPolicy struct {
	mainResource *temporal.Resource[azcore.AccessToken, acquiringResourceState]
	cred         azcore.TokenCredential
	aadScopes    []string
	acrScope     string
	acrService   string
	authClient   *generated.AuthenticationClient
}

func NewAuthenticationPolicy(cred azcore.TokenCredential, scopes []string, authClient *generated.AuthenticationClient, opts *AuthenticationPolicyOptions) *AuthenticationPolicy {
	return &AuthenticationPolicy{
		cred:         cred,
		aadScopes:    scopes,
		authClient:   authClient,
		mainResource: temporal.NewResource(acquire),
	}
}

func (p *AuthenticationPolicy) Do(req *policy.Request) (*http.Response, error) {

	// send the original request
	resp, reqErr := req.Next()
	if reqErr != nil {
		return nil, reqErr
	}

	// do challenge process
	if resp.StatusCode == 401 {
		err := p.findServiceAndScope(resp)
		if err != nil {
			return nil, err
		}

		accessToken, err := p.getAccessToken(req)
		if err != nil {
			return nil, err
		}

		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s%s", bearerHeader, accessToken),
		)

		// send a copy of the request
		cloneReq := req.Clone(req.Raw().Context())
		resp, cloneReqErr := cloneReq.Next()
		if cloneReqErr != nil {
			return nil, cloneReqErr
		}

		return resp, nil
	}

	return resp, nil
}

func (p *AuthenticationPolicy) getAccessToken(req *policy.Request) (string, error) {
	// anonymous access
	if p.cred == nil {
		resp, err := p.authClient.ExchangeAcrRefreshTokenForAcrAccessToken(req.Raw().Context(), p.acrService, p.acrScope, "", &generated.AuthenticationClientExchangeAcrRefreshTokenForAcrAccessTokenOptions{GrantType: to.Ptr(generated.TokenGrantTypePassword)})
		if err != nil {
			return "", err
		}
		return *resp.AccessToken, nil
	}

	as := acquiringResourceState{
		p:   p,
		req: req,
	}

	refreshToken, err := p.mainResource.Get(as)
	if err != nil {
		return "", err
	}

	resp, err := p.authClient.ExchangeAcrRefreshTokenForAcrAccessToken(req.Raw().Context(), p.acrService, p.acrScope, refreshToken.Token, &generated.AuthenticationClientExchangeAcrRefreshTokenForAcrAccessTokenOptions{GrantType: to.Ptr(generated.TokenGrantTypePassword)})
	if err != nil {
		return "", err
	}
	return *resp.AccessToken, nil
}

type challengePolicyError struct {
	err error
}

func (c *challengePolicyError) Error() string {
	return c.err.Error()
}

func (*challengePolicyError) NonRetriable() {
	// marker method
}

func (c *challengePolicyError) Unwrap() error {
	return c.err
}

var _ errorinfo.NonRetriable = (*challengePolicyError)(nil)

func (p *AuthenticationPolicy) findServiceAndScope(resp *http.Response) error {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return &challengePolicyError{err: errors.New("response has no WWW-Authenticate header for challenge authentication")}
	}

	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
	parts := strings.Split(authHeader, ",")
	vals := map[string]string{}
	for _, part := range parts {
		subParts := strings.Split(part, "=")
		if len(subParts) == 2 {
			stripped := strings.ReplaceAll(subParts[1], "\"", "")
			stripped = strings.TrimSuffix(stripped, ",")
			vals[subParts[0]] = stripped
		}
	}

	if v, ok := vals["acrScope"]; ok {
		p.acrScope = v
	} else if v, ok := vals["resource"]; ok {
		p.acrScope = v
	}
	if p.acrScope == "" {
		return &challengePolicyError{err: errors.New("could not find a valid acrScope in the WWW-Authenticate header")}
	}

	if v, ok := vals["acrService"]; ok {
		p.acrService = v
	}
	if p.acrService == "" {
		return &challengePolicyError{err: errors.New("could not find a valid acrService in the WWW-Authenticate header")}
	}

	return nil
}

type acquiringResourceState struct {
	req *policy.Request
	p   *AuthenticationPolicy
}

// acquire acquires or updates the resource; only one
// thread/goroutine at a time ever calls this function
func acquire(state acquiringResourceState) (newResource azcore.AccessToken, newExpiration time.Time, err error) {
	aadToken, err := state.p.cred.GetToken(
		state.req.Raw().Context(),
		policy.TokenRequestOptions{
			Scopes: state.p.aadScopes,
		},
	)
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	refreshResp, err := state.p.authClient.ExchangeAADAccessTokenForAcrRefreshToken(state.req.Raw().Context(), generated.PostContentSchemaGrantTypeAccessToken, state.p.acrService, &generated.AuthenticationClientExchangeAADAccessTokenForAcrRefreshTokenOptions{
		AccessToken: &aadToken.Token,
	})
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	refreshToken := azcore.AccessToken{
		Token: *refreshResp.RefreshToken,
	}
	refreshToken.ExpiresOn, err = getJWTExpireTime(aadToken.Token)
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	return refreshToken, refreshToken.ExpiresOn, nil
}

func getJWTExpireTime(token string) (time.Time, error) {
	values := strings.Split(token, ".")
	if len(values) > 2 {
		value := values[1]
		padding := len(value) % 4
		if padding > 0 {
			for i := 0; i < padding; i++ {
				value += "="
			}
		}
		parsedValue, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return time.Time{}, err
		}

		var jsonValue *jwtOnlyWithExp
		err = json.Unmarshal(parsedValue, jsonValue)
		if err != nil {
			return time.Time{}, err
		}
		return jsonValue.Exp, nil
	}

	return time.Time{}, &challengePolicyError{err: errors.New("could not parse refresh token expire time")}
}

type jwtOnlyWithExp struct {
	Exp time.Time `json:"exp"`
}
