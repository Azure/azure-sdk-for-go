// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

const apiVersionTopicXML = `<entry xmlns="http://www.w3.org/2005/Atom"><title>my-topic</title><content type="application/xml"><TopicDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"><SubscriptionCount>1</SubscriptionCount><CountDetails><ScheduledMessageCount>0</ScheduledMessageCount></CountDetails><CreatedAt>2026-01-01T00:00:00Z</CreatedAt><UpdatedAt>2026-01-01T00:00:00Z</UpdatedAt><AccessedAt>2026-01-01T00:00:00Z</AccessedAt></TopicDescription></content></entry>`

// captureAPIVersionPolicy records the api-version the client is about to send and answers the
// request itself, so nothing reaches the network.
type captureAPIVersionPolicy struct {
	apiVersions []string
}

func (p *captureAPIVersionPolicy) Do(req *policy.Request) (*http.Response, error) {
	// The whole slice rather than Get(), so a second api-version parameter is visible here
	// rather than hidden behind the first one.
	p.apiVersions = req.Raw().URL.Query()["api-version"]

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(apiVersionTopicXML)),
		Request:    req.Raw(),
	}, nil
}

// fakeAPIVersionCredential is never invoked: captureAPIVersionPolicy answers the request at the
// per-call stage, before the per-retry auth policy runs. It exists because NewClient requires a
// credential.
type fakeAPIVersionCredential struct{}

func (fakeAPIVersionCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// TestClientAPIVersion covers the promise the ClientOptions doc comment makes, at the public
// boundary a customer touches. The atom-level TestEntityManagerAPIVersion covers the pipeline
// wiring; this covers the options plumbing through NewClient that carries a caller's APIVersion
// down to it.
func TestClientAPIVersion(t *testing.T) {
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

			client, err := NewClient("fake.servicebus.windows.net", fakeAPIVersionCredential{}, &ClientOptions{
				ClientOptions: azcore.ClientOptions{
					APIVersion:      td.override,
					PerCallPolicies: []policy.Policy{capture},
				},
			})
			require.NoError(t, err)

			_, err = client.GetTopicRuntimeProperties(context.Background(), "my-topic", nil)
			require.NoError(t, err)

			require.Equal(t, []string{td.expected}, capture.apiVersions)
		})
	}
}
