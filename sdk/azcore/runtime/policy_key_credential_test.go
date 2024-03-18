// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestKeyCredentialPolicy(t *testing.T) {
	const key = "foo"
	cred := exported.NewKeyCredential(key)

	const headerName = "fake-auth"
	policy := NewKeyCredentialPolicy(cred, headerName, nil)
	require.NotNil(t, policy)

	pl := exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		require.EqualValues(t, key, req.Header.Get(headerName))
		return &http.Response{}, nil
	}), policy)

	req, err := NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)

	policy = NewKeyCredentialPolicy(cred, headerName, &KeyCredentialPolicyOptions{
		Prefix: "Prefix: ",
	})
	require.NotNil(t, policy)

	pl = exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		require.EqualValues(t, "Prefix: "+key, req.Header.Get(headerName))
		return &http.Response{}, nil
	}), policy)

	req, err = NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)
}

func TestKeyCredentialPolicy_RequiresHTTPS(t *testing.T) {
	cred := exported.NewKeyCredential("foo")

	policy := NewKeyCredentialPolicy(cred, "fake-auth", nil)
	require.NotNil(t, policy)

	pl := exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, nil
	}), policy)

	req, err := NewRequest(context.Background(), http.MethodGet, "http://contoso.com")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.Error(t, err)
}

func TestKeyCredentialPolicy_NilCredential(t *testing.T) {
	const headerName = "fake-auth"
	policy := NewKeyCredentialPolicy(nil, headerName, nil)
	require.NotNil(t, policy)

	pl := exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		require.Zero(t, req.Header.Get(headerName))
		return &http.Response{}, nil
	}), policy)

	req, err := NewRequest(context.Background(), http.MethodGet, "http://contoso.com")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)
}

func TestKeyCredentialPolicy_InsecureAllowCredentialWithHTTP(t *testing.T) {
	cred := exported.NewKeyCredential("foo")

	policy := NewKeyCredentialPolicy(cred, "fake-auth", &KeyCredentialPolicyOptions{
		InsecureAllowCredentialWithHTTP: true,
	})
	require.NotNil(t, policy)

	pl := exported.NewPipeline(shared.TransportFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, nil
	}), policy)

	req, err := NewRequest(context.Background(), http.MethodGet, "http://contoso.com")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)
}
