// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

const testURL = "http://test.contoso.com/"

func TestNewRequest(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, testURL)
	if err != nil {
		t.Fatal(err)
	}
	if m := req.Raw().Method; m != http.MethodPost {
		t.Fatalf("unexpected method %s", m)
	}
	type myValue struct{}
	var mv myValue
	if req.OperationValue(&mv) {
		t.Fatal("expected missing custom operation value")
	}
	req.SetOperationValue(myValue{})
	if !req.OperationValue(&mv) {
		t.Fatal("missing custom operation value")
	}
}

type testPolicy struct{}

func (testPolicy) Do(*Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestRequestPolicies(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, testURL)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := req.Next()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	req.policies = []Policy{}
	resp, err = req.Next()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	req.policies = []Policy{testPolicy{}}
	resp, err = req.Next()
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("unexpected nil response")
	}
}

func TestRequestBody(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, testURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := req.RewindBody(); err != nil {
		t.Fatal(err)
	}
	if err := req.Close(); err != nil {
		t.Fatal(err)
	}
	if req.Body() != nil {
		t.Fatal("expected nil body")
	}
	if req.req.GetBody != nil {
		t.Fatal("expected nil GetBody")
	}
	if err := req.SetBody(NopCloser(strings.NewReader("test")), "application/text"); err != nil {
		t.Fatal(err)
	}
	if req.Body() == nil {
		t.Fatal("unexpected nil body")
	}
	if req.req.GetBody == nil {
		t.Fatal("unexpected nil GetBody")
	}
	body, err := req.req.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "test" {
		t.Fatalf("unexpected body %s", string(b))
	}
	if err := req.RewindBody(); err != nil {
		t.Fatal(err)
	}
	if err := req.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRequestEmptyBody(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, testURL)
	require.NoError(t, err)
	require.NoError(t, req.SetBody(NopCloser(strings.NewReader("")), "application/text"))
	require.Nil(t, req.Body())
	require.NotContains(t, req.Raw().Header, shared.HeaderContentLength)
	require.Equal(t, []string{"application/text"}, req.Raw().Header[shared.HeaderContentType])

	// SetBody should treat a nil ReadSeekCloser the same as one having no content
	req, err = NewRequest(context.Background(), http.MethodPost, testURL)
	require.NoError(t, err)
	require.NoError(t, req.SetBody(nil, ""))
	require.Nil(t, req.Body())
	require.NotContains(t, req.Raw().Header, shared.HeaderContentLength)
	require.NotContains(t, req.Raw().Header, shared.HeaderContentType)

	// SetBody should allow replacing a previously set body with an empty one
	req, err = NewRequest(context.Background(), http.MethodPost, testURL)
	require.NoError(t, err)
	require.NoError(t, req.SetBody(NopCloser(strings.NewReader("content")), "application/text"))
	require.NoError(t, req.SetBody(nil, "application/json"))
	require.Nil(t, req.Body())
	require.NotContains(t, req.Raw().Header, shared.HeaderContentLength)
	require.Equal(t, []string{"application/json"}, req.Raw().Header[shared.HeaderContentType])
}

func TestRequestClone(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, testURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := req.SetBody(NopCloser(strings.NewReader("test")), "application/text"); err != nil {
		t.Fatal(err)
	}
	type ensureCloned struct {
		Count int
	}
	source := ensureCloned{Count: 12345}
	req.SetOperationValue(source)
	clone := req.Clone(context.Background())
	var cloned ensureCloned
	if !clone.OperationValue(&cloned) {
		t.Fatal("missing operation value")
	}
	if cloned.Count != source.Count {
		t.Fatal("wrong operation value")
	}
	if clone.body == nil {
		t.Fatal("missing body")
	}
}

func TestNewRequestFail(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodOptions, "://test.contoso.com/")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
	req, err = NewRequest(context.Background(), http.MethodPatch, "/missing/the/host")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
	req, err = NewRequest(context.Background(), http.MethodPatch, "mailto://nobody.contoso.com")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
}

func TestRequestWithContext(t *testing.T) {
	type ctxKey1 struct{}
	type ctxKey2 struct{}

	req1, err := NewRequest(context.WithValue(context.Background(), ctxKey1{}, 1), http.MethodPost, testURL)
	require.NoError(t, err)
	require.NotNil(t, req1.Raw().Context().Value(ctxKey1{}))

	req2 := req1.WithContext(context.WithValue(context.Background(), ctxKey2{}, 1))
	require.Nil(t, req2.Raw().Context().Value(ctxKey1{}))
	require.NotNil(t, req2.Raw().Context().Value(ctxKey2{}))

	// shallow copy, so changing req2 affects req1
	req2.Raw().Header.Add("added-req2", "value")
	require.EqualValues(t, "value", req1.Raw().Header.Get("added-req2"))
}

func TestNewRequestForNextLink(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		nextLink string
		expected string
	}{
		{
			name:     "empty next link joins endpoint",
			endpoint: "https://example.com",
			nextLink: "",
			expected: "https://example.com/",
		},
		{
			name:     "relative next link joins endpoint",
			endpoint: "https://example.com",
			nextLink: "page2",
			expected: "https://example.com/page2",
		},
		{
			name:     "relative next link with query joins endpoint",
			endpoint: "https://example.com",
			nextLink: "page2?skip=10",
			expected: "https://example.com/page2?skip=10",
		},
		{
			name:     "absolute next link is used as-is",
			endpoint: "https://example.com",
			nextLink: "https://other.com/page3?token=abc",
			expected: "https://other.com/page3?token=abc",
		},
		{
			name:     "absolute next link overrides endpoint",
			endpoint: "https://example.com",
			nextLink: "https://other.com/next",
			expected: "https://other.com/next",
		},
		{
			name:     "relative next link with plus in query is encoded",
			endpoint: "https://example.com",
			nextLink: "page2?filter=a+b",
			expected: "https://example.com/page2?filter=a%20b",
		},
		{
			name:     "absolute next link with plus in query is encoded",
			endpoint: "https://example.com",
			nextLink: "https://example.com/next?filter=a+b+c",
			expected: "https://example.com/next?filter=a%20b%20c",
		},
		{
			name:     "escaped plus in query is preserved",
			endpoint: "https://example.com",
			nextLink: "https://example.com/next?filter=a%2Bb",
			expected: "https://example.com/next?filter=a%2Bb",
		},
		{
			name:     "multiple query params with spaces and plus",
			endpoint: "https://example.com",
			nextLink: "https://example.com/next?a=1+2&b=hello%20world",
			expected: "https://example.com/next?a=1%202&b=hello%20world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequestForNextLink(context.Background(), http.MethodGet, tt.endpoint, tt.nextLink)
			require.NoError(t, err)
			require.Equal(t, tt.expected, req.Raw().URL.String())
		})
	}
}

func TestEncodeQueryParams(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		encodeSpaces bool
		expected     string
		expectErr    bool
	}{
		{
			name:     "no query params",
			url:      "https://example.com/foo",
			expected: "https://example.com/foo",
		},
		{
			name:     "empty query",
			url:      "https://example.com/foo?",
			expected: "https://example.com/foo?",
		},
		{
			name:     "single query param",
			url:      "https://example.com/foo?a=1",
			expected: "https://example.com/foo?a=1",
		},
		{
			name:     "params are sorted by key",
			url:      "https://example.com/foo?b=2&a=1",
			expected: "https://example.com/foo?a=1&b=2",
		},
		{
			name:         "plus encoded as space when encodeSpaces is true",
			url:          "https://example.com/foo?a=1+2",
			encodeSpaces: true,
			expected:     "https://example.com/foo?a=1%202",
		},
		{
			name:         "plus encoded as plus when encodeSpaces is false",
			url:          "https://example.com/foo?a=1+2",
			encodeSpaces: false,
			expected:     "https://example.com/foo?a=1+2",
		},
		{
			name:         "multiple plus chars with encodeSpaces",
			url:          "https://example.com/foo?a=1+2+3",
			encodeSpaces: true,
			expected:     "https://example.com/foo?a=1%202%203",
		},
		{
			name:     "escaped plus is preserved",
			url:      "https://example.com/foo?a=1%2B2",
			expected: "https://example.com/foo?a=1%2B2",
		},
		{
			name:     "semicolon is escaped",
			url:      "https://example.com/foo?a=1;2",
			expected: "https://example.com/foo?a=1%3B2",
		},
		{
			name:         "multiple params with spaces and plus",
			url:          "https://example.com/foo?a=1+2&b=hello%20world",
			encodeSpaces: true,
			expected:     "https://example.com/foo?a=1%202&b=hello%20world",
		},
		{
			name:      "invalid query escape returns error",
			url:       "https://example.com/foo?a=%zz",
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := EncodeQueryParams(tt.url, tt.encodeSpaces)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestSetBodyWithClobber(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPatch, "https://contoso.com")
	require.NoError(t, err)
	require.NotNil(t, req)
	req.req.Header.Set(shared.HeaderContentType, "clobber-me")
	require.NoError(t, SetBody(req, NopCloser(strings.NewReader(`"json-string"`)), shared.ContentTypeAppJSON, true))
	require.EqualValues(t, shared.ContentTypeAppJSON, req.req.Header.Get(shared.HeaderContentType))
}

func TestSetBodyWithNoClobber(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPatch, "https://contoso.com")
	require.NoError(t, err)
	require.NotNil(t, req)
	const mergePatch = "application/merge-patch+json"
	req.req.Header.Set(shared.HeaderContentType, mergePatch)
	require.NoError(t, SetBody(req, NopCloser(strings.NewReader(`"json-string"`)), shared.ContentTypeAppJSON, false))
	require.EqualValues(t, mergePatch, req.req.Header.Get(shared.HeaderContentType))
}
