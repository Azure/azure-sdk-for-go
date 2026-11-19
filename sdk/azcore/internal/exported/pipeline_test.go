// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/sanitizer"
	"github.com/stretchr/testify/require"
)

func TestPipelineErrors(t *testing.T) {
	pl := NewPipeline(nil, nil)
	resp, err := pl.Do(nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	req, err := NewRequest(context.Background(), http.MethodGet, testURL)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	req.Raw().Header["Invalid"] = []string{string([]byte{0})}
	resp, err = pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	req, err = NewRequest(context.Background(), http.MethodGet, testURL)
	if err != nil {
		t.Fatal(err)
	}
	req.Raw().Header["Inv alid"] = []string{"value"}
	resp, err = pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

type mockTransport struct {
	succeed bool
	both    bool
}

func (m *mockTransport) Do(*http.Request) (*http.Response, error) {
	if m.both {
		return nil, nil
	}
	if m.succeed {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}
	return nil, errors.New("failed")
}

func TestPipelineDo(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodGet, testURL)
	if err != nil {
		t.Fatal(err)
	}
	tp := mockTransport{succeed: true}
	pl := NewPipeline(nil, &tp)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if sc := resp.StatusCode; sc != http.StatusOK {
		t.Fatalf("unexpected status code %d", sc)
	}
	tp.succeed = false
	resp, err = pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	tp.both = true
	resp, err = pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestPipelineSanitizer(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		require.Nil(t, PipelineSanitizer(Pipeline{}))
	})

	u, err := url.Parse("https://example.com?api-version=1&custom=visible&sig=secret")
	require.NoError(t, err)
	header := http.Header{"Content-Type": []string{"application/json"}, "Custom": []string{"visible"}}
	for _, tc := range []struct {
		name        string
		pipeline    Pipeline
		expectedURL string
		contentType string
		custom      string
	}{
		{
			name:        "nil sanitizer",
			pipeline:    NewPipeline(nil, nil),
			expectedURL: "https://example.com?api-version=1&custom=REDACTED&sig=REDACTED",
			contentType: "application/json",
			custom:      "REDACTED",
		},
		{
			name: "custom sanitizer",
			pipeline: NewPipeline(sanitizer.New(sanitizer.Options{
				AllowedHeaders:     []string{"Custom"},
				AllowedQueryParams: []string{"custom"},
			}), nil),
			expectedURL: "https://example.com?api-version=1&custom=visible&sig=REDACTED",
			contentType: "application/json",
			custom:      "visible",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := PipelineSanitizer(tc.pipeline)
			require.NotNil(t, s)
			require.Equal(t, tc.expectedURL, s.GetSanitizedURL(*u))
			require.Equal(t, "failed: "+tc.expectedURL, s.GetSanitizedURLString("failed: "+u.String(), u))
			require.Equal(t, tc.contentType, s.GetSanitizedHeader("Content-Type", header))
			require.Equal(t, tc.custom, s.GetSanitizedHeader("Custom", header))
		})
	}
}
