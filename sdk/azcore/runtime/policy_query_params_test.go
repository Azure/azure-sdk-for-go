//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestWithQueryParametersSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const (
		customQP         = "custom-qp"
		customValue      = "custom-value"
		preexistingQP    = "preexisting-qp"
		preexistingValue = "preexisting-value"
	)
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		// ensure preexisting query param wasn't removed
		qp, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return false
		}
		return qp.Get(customQP) == customValue && qp.Get(preexistingQP) == preexistingValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(policy.WithQueryParameters(context.Background(), url.Values{
		customQP: []string{customValue},
	}), http.MethodGet, srv.URL())
	require.NoError(t, err)
	qp := url.Values{}
	qp.Set(preexistingQP, preexistingValue)
	req.Raw().URL.RawQuery = qp.Encode()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}

func TestWithQueryParametersFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const (
		customQP    = "custom-qp"
		customValue = "custom-value"
	)
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		qp, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return false
		}
		return qp.Get(customQP) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
}

func TestWithQueryParametersOverwrite(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const (
		customQP    = "custom-qp"
		customValue = "custom-value"
	)
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		qp, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return false
		}
		return qp.Get(customQP) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(policy.WithQueryParameters(context.Background(), url.Values{
		customQP: []string{customValue},
	}), http.MethodGet, srv.URL()+fmt.Sprintf("?%s=overwrite-me", customQP))
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}

func TestWithQueryParametersMultipleValues(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const (
		customQP     = "custom-qp"
		customValue1 = "custom-value1"
		customValue2 = "custom-value2"
	)
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		qp, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return false
		}
		return qp[customQP][0] == customValue1 && qp[customQP][1] == customValue2
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(policy.WithQueryParameters(context.Background(), url.Values{
		customQP: []string{customValue1, customValue2},
	}), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}
