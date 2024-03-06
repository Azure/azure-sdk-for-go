//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestPipelineErrors(t *testing.T) {
	pl := NewPipeline(nil)
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
	pl := NewPipeline(&tp)
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
