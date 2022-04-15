//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"net/http"
	"strings"
	"testing"
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
	if err := req.SetBody(NopCloser(strings.NewReader("test")), "application/text"); err != nil {
		t.Fatal(err)
	}
	if err := req.RewindBody(); err != nil {
		t.Fatal(err)
	}
	if err := req.Close(); err != nil {
		t.Fatal(err)
	}
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
