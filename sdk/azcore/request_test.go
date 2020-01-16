// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

type testJSON struct {
	SomeInt    int
	SomeString string
}

type testXML struct {
	SomeInt    int
	SomeString string
}

func TestRequestMarshalXML(t *testing.T) {
	u, err := url.Parse("https://contoso.com")
	if err != nil {
		panic(err)
	}
	req := NewRequest(http.MethodPost, *u)
	err = req.MarshalAsXML(testXML{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppXML {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppXML)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
}

func TestRequestEmptyPipeline(t *testing.T) {
	u, err := url.Parse("https://contoso.com")
	if err != nil {
		panic(err)
	}
	req := NewRequest(http.MethodPost, *u)
	resp, err := req.Next(context.Background())
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if !errors.Is(err, ErrNoMorePolicies) {
		t.Fatalf("expected ErrNoMorePolicies, got %v", err)
	}
}

func TestRequestMarshalJSON(t *testing.T) {
	u, err := url.Parse("https://contoso.com")
	if err != nil {
		panic(err)
	}
	req := NewRequest(http.MethodPost, *u)
	err = req.MarshalAsJSON(testJSON{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppJSON)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
}
