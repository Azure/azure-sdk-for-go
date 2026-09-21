// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/stretchr/testify/require"
)

type fakePolicy struct {
	t *testing.T
}

func (p *fakePolicy) Do(req *policy.Request) (*http.Response, error) {
	// Covers the PUT path. TestEntityManagerAPIVersion covers GET, so the two together cover
	// both verbs.
	require.Equal(p.t, []string{"2024-05"}, req.Raw().URL.Query()["api-version"])
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

// captureAPIVersionPolicy records the api-version the pipeline is about to send and
// answers the request itself, so it never reaches the retry, auth or transport policies.
type captureAPIVersionPolicy struct {
	apiVersions []string
}

func (p *captureAPIVersionPolicy) Do(req *policy.Request) (*http.Response, error) {
	// The whole slice rather than Get(), so a second api-version parameter is visible here
	// rather than hidden behind the first one.
	p.apiVersions = req.Raw().URL.Query()["api-version"]

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("<string>response body</string>")),
		Request:    req.Raw(),
	}, nil
}

// fakeTokenCredential is never invoked: captureAPIVersionPolicy answers the request at the
// per-call stage, before the per-retry auth policy runs. It exists because
// newEntityManagerImpl requires a non-nil provider.
type fakeTokenCredential struct{}

func (fakeTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestEntityManagerAPIVersion(t *testing.T) {
	testData := []struct {
		name     string
		override string
		expected string
	}{
		{name: "default", override: "", expected: "2024-05"},
		{name: "overridden", override: "2021-05", expected: "2021-05"},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			capture := &captureAPIVersionPolicy{}

			// PerCallPolicies run after azcore's api-version policy, so capture sees the
			// version that would go out on the wire.
			em, err := newEntityManagerImpl(
				sbauth.NewTokenProvider(fakeTokenCredential{}),
				"version",
				&policy.ClientOptions{
					APIVersion:      td.override,
					PerCallPolicies: []policy.Policy{capture},
				},
				"fake.servicebus.windows.net")
			require.NoError(t, err)

			var respBody *string
			_, err = em.Get(context.Background(), "entityPath", &respBody)
			require.NoError(t, err)

			require.Equal(t, []string{td.expected}, capture.apiVersions)
		})
	}
}
