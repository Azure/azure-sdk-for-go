// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestKeyCredentialPolicy(t *testing.T) {
	const key = "foo"
	const target = "http://abc/de"
	keyPolicy := NewWebPubSubKeyCredentialPolicy(key)
	require.NotNil(t, keyPolicy)
	verifier := PolicyFunc(func(req *policy.Request) (*http.Response, error) {
		token := req.Raw().Header.Get("Authorization")

		require.True(t, strings.HasPrefix(token, "Bearer "))
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		require.NoError(t, err)
		require.Equal(t, target, claims.Claims.(jwt.MapClaims)["aud"])
		return &http.Response{}, nil
	})

	pl := newPipeline(runtime.PipelineOptions{PerCall: []policy.Policy{keyPolicy, verifier}},
		&policy.ClientOptions{})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, target)
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)
}

type PolicyFunc func(req *policy.Request) (*http.Response, error)

func (f PolicyFunc) Do(req *policy.Request) (*http.Response, error) {
	return f(req)
}

func newPipeline(plOpts runtime.PipelineOptions, options *policy.ClientOptions) runtime.Pipeline {
	return runtime.NewPipeline(ModuleName+".Client", ModuleVersion, plOpts, options)
}
