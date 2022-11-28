//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

// TransportFunc is a helper to use a first-class func to satisfy the Transporter interface.
type TransportFunc func(*http.Request) (*http.Response, error)

// Do implements the Transporter interface for the TransportFunc type.
func (pf TransportFunc) Do(req *http.Request) (*http.Response, error) {
	return pf(req)
}

type nonRetriableError struct {
	error
}

func (nonRetriableError) NonRetriable() {}

func TestSyncTokenPolicy(t *testing.T) {
	stp := newSyncTokenPolicy()
	require.NotNil(t, stp)

	pl := runtime.NewPipeline("TestSyncTokenPolicy", moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{stp}}, &policy.ClientOptions{
		Transport: TransportFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
				Header: http.Header{
					"Sync-Token": []string{"jtqGc1I4=MDoyOA==;sn=28"},
				},
			}, nil
		}),
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "http://test.contoso.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	st := stp.syncTokens["jtqGc1I4"]
	require.NotZero(t, st)
	require.Equal(t, "jtqGc1I4", st.id)
	require.Equal(t, int64(28), st.seqNo)
	require.Equal(t, "MDoyOA==", st.value)
}

func TestSyncTokenPolicyError(t *testing.T) {
	stp := newSyncTokenPolicy()
	require.NotNil(t, stp)

	pl := runtime.NewPipeline("TestSyncTokenPolicy", moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{stp}}, &policy.ClientOptions{
		Transport: TransportFunc(func(req *http.Request) (*http.Response, error) {
			return nil, nonRetriableError{errors.New("failed")}
		}),
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "http://test.contoso.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.Error(t, err)
	require.Nil(t, resp)
}
