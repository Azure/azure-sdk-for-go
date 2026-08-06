// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"io"
	"net/http"
	"net/url"
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

func TestEncodeQueryParams(t *testing.T) {
	for _, test := range []struct {
		name    string
		url     string
		options *EncodeQueryParamsOptions
		want    string
	}{
		{
			name: "no query params",
			url:  "https://contoso.com/path",
			want: "https://contoso.com/path",
		},
		{
			name: "empty query string",
			url:  "https://contoso.com/path?",
			want: "https://contoso.com/path",
		},
		{
			name: "existing params are re-encoded",
			url:  "https://contoso.com/path?b=2&a=1",
			want: "https://contoso.com/path?a=1&b=2",
		},
		{
			name: "semicolon in existing value is escaped",
			url:  "https://contoso.com/path?a=1;2",
			want: "https://contoso.com/path?a=1%3B2",
		},
		{
			name: "space in existing value uses %20 not +",
			url:  "https://contoso.com/path?a=hello world",
			want: "https://contoso.com/path?a=hello%20world",
		},
		{
			name: "plus in existing value is treated as a space and encoded as %20",
			url:  "https://contoso.com/path?a=1+2",
			want: "https://contoso.com/path?a=1%202",
		},
		{
			name:    "literal plus in option value is escaped as %2B",
			url:     "https://contoso.com/path",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"1+2"}}},
			want:    "https://contoso.com/path?a=1%2B2",
		},
		{
			name: "reserved characters in existing value are escaped",
			url:  "https://contoso.com/path?a=a b&c=d/e",
			want: "https://contoso.com/path?a=a%20b&c=d%2Fe",
		},
		{
			name:    "nil options is a no-op",
			url:     "https://contoso.com/path?a=1",
			options: nil,
			want:    "https://contoso.com/path?a=1",
		},
		{
			name:    "options params added to URL with no query string",
			url:     "https://contoso.com/path",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"1"}}},
			want:    "https://contoso.com/path?a=1",
		},
		{
			name:    "options params added to URL with trailing ?",
			url:     "https://contoso.com/path?",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"1"}}},
			want:    "https://contoso.com/path?a=1",
		},
		{
			name:    "options params merged with existing params",
			url:     "https://contoso.com/path?a=1",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"b": {"2"}}},
			want:    "https://contoso.com/path?a=1&b=2",
		},
		{
			name:    "options overwrite an existing single-valued key",
			url:     "https://contoso.com/path?a=1",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"2"}}},
			want:    "https://contoso.com/path?a=2",
		},
		{
			name:    "options overwrite an existing multi-valued key",
			url:     "https://contoso.com/path?a=1&a=2",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"3"}}},
			want:    "https://contoso.com/path?a=3",
		},
		{
			name:    "options overwrite an existing key with multiple new values",
			url:     "https://contoso.com/path?a=1",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"2", "3"}}},
			want:    "https://contoso.com/path?a=2&a=3",
		},
		{
			name:    "multiple option values for a key",
			url:     "https://contoso.com/path",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"1", "2"}}},
			want:    "https://contoso.com/path?a=1&a=2",
		},
		{
			name:    "option values with special characters are escaped",
			url:     "https://contoso.com/path",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{"a": {"x;y z"}}},
			want:    "https://contoso.com/path?a=x%3By%20z",
		},
		{
			name:    "empty options parameters is a no-op",
			url:     "https://contoso.com/path?a=1",
			options: &EncodeQueryParamsOptions{Parameters: url.Values{}},
			want:    "https://contoso.com/path?a=1",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := EncodeQueryParams(test.url, test.options)
			require.NoError(t, err)
			require.EqualValues(t, test.want, got)
		})
	}
}

func TestEncodeQueryParamsInvalid(t *testing.T) {
	// an invalid percent-encoding in the existing query string surfaces as an error
	_, err := EncodeQueryParams("https://contoso.com/path?a=%zz", nil)
	require.Error(t, err)
}
