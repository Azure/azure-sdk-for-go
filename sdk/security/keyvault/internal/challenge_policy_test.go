//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

type credentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (cf credentialFunc) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return cf(ctx, options)
}

func TestChallengePolicy(t *testing.T) {
	accessToken := "***"
	resource := "https://vault.azure.net"
	scope := "https://vault.azure.net/.default"
	challengeResource := `Bearer authorization="https://login.microsoftonline.com/{tenant}", resource="{resource}"`
	challengeScope := `Bearer authorization="https://login.microsoftonline.com/{tenant}", scope="{resource}"`

	for _, test := range []struct {
		expectedScope, format, resource string
		disableVerify, err              bool
	}{
		// happy path: resource matches requested vault's host (vault.azure.net)
		{format: challengeResource, resource: resource, expectedScope: scope},
		{format: challengeResource, resource: resource, disableVerify: true, expectedScope: scope},
		{format: challengeScope, resource: scope, expectedScope: scope},
		{format: challengeScope, resource: scope, disableVerify: true, expectedScope: scope},
		// the policy should prefer scope to resource when a challenge specifies both
		{format: fmt.Sprintf(`%s scope="%s"`, challengeResource, scope), resource: resource, expectedScope: scope},
		{format: challengeScope + ` resource="ignore me"`, resource: scope, expectedScope: scope},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{format: challengeResource, resource: "https://vault.azure.cn", err: true},
		{format: challengeResource, resource: "https://myvault.azure.net", err: true},
		{format: challengeScope, resource: "https://vault.azure.cn/.default", err: true},
		{format: challengeScope, resource: "https://myvault.azure.net/.default", err: true},

		// the policy shouldn't return errors for the above cases when verification is disabled
		{format: challengeResource, resource: "https://vault.azure.cn", disableVerify: true, expectedScope: "https://vault.azure.cn/.default"},
		{format: challengeResource, resource: "https://myvault.azure.net", disableVerify: true, expectedScope: "https://myvault.azure.net/.default"},
		{format: challengeScope, resource: "https://vault.azure.cn/.default", disableVerify: true, expectedScope: "https://vault.azure.cn/.default"},
		{format: challengeScope, resource: "https://myvault.azure.net/.default", disableVerify: true, expectedScope: "https://myvault.azure.net/.default"},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", strings.ReplaceAll(test.format, "{resource}", test.resource)),
				mock.WithStatusCode(401),
			)
			srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
				if authz := r.Header.Values("Authorization"); len(authz) != 1 || authz[0] != "Bearer "+accessToken {
					t.Errorf(`unexpected Authorization "%s"`, authz)
				}
				return true
			}))
			srv.AppendResponse()
			authenticated := false
			cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				authenticated = true
				require.Equal(t, []string{test.expectedScope}, tro.Scopes)
				return azcore.AccessToken{Token: accessToken, ExpiresOn: time.Now().Add(time.Hour)}, nil
			})
			p := NewKeyVaultChallengePolicy(cred, &KeyVaultChallengePolicyOptions{DisableChallengeResourceVerification: test.disableVerify})
			pl := runtime.NewPipeline("", "",
				runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
				&policy.ClientOptions{Transport: srv},
			)
			req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
			require.NoError(t, err)
			_, err = pl.Do(req)
			if test.err {
				expected := fmt.Sprintf(challengeMatchError, test.resource)
				require.EqualError(t, err, expected)
				require.IsType(t, &challengePolicyError{}, err)
			} else {
				require.True(t, authenticated, "policy should have authenticated")
			}
		})
	}
}

func TestParseTenant(t *testing.T) {
	actual := parseTenant("")
	require.Empty(t, actual)

	expected := "00000000-0000-0000-0000-000000000000"
	sampleURL := "https://login.microsoftonline.com/" + expected
	actual = parseTenant(sampleURL)
	require.Equal(t, expected, actual, "tenant was not properly parsed, got %s, expected %s", actual, expected)
}
