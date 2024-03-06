// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

type fakePolicy struct {
	t *testing.T
}

func (p *fakePolicy) Do(req *policy.Request) (*http.Response, error) {
	require.NotNil(p.t, req.Body())

	reqBytes, err := io.ReadAll(req.Body())
	require.NoError(p.t, err)
	require.Equal(p.t, "<string>hello</string>", string(reqBytes))

	// now rewind it, and try again - this is what the retry policy does.
	err = req.RewindBody()
	require.NoError(p.t, err)

	reqBytes, err = io.ReadAll(req.Body())
	require.NoError(p.t, err)
	require.Equal(p.t, "<string>hello</string>", string(reqBytes))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("<string>response body</string>")),
	}, nil
}

func TestEntityManagerRewindable(t *testing.T) {
	// prior to this I was populating the .Raw().Body field which works for first
	// requests but will fail if there is a retry since the body can't be rewound.
	pl := runtime.NewPipeline("module", "version", runtime.PipelineOptions{
		PerCall: []policy.Policy{
			&fakePolicy{t: t},
		},
	}, nil)

	em := entityManager{
		pl:   pl,
		Host: "https://localhost",
	}

	var respBody *string
	resp, err := em.Put(context.Background(), "entityPath", "hello", &respBody, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "response body", *respBody)
}
