//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNopCloser(t *testing.T) {
	nc := NopCloser(strings.NewReader("foo"))
	if err := nc.Close(); err != nil {
		t.Fatal(err)
	}
}

type testError struct {
	m string
}

func (t testError) Error() string {
	return t.m
}

func TestNewResponseError(t *testing.T) {
	err := NewResponseError(testError{m: "crash"}, &http.Response{StatusCode: http.StatusInternalServerError})
	if s := err.Error(); s != "crash" {
		t.Fatalf("unexpected error %s", s)
	}
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	re.NonRetriable() // nop
	if c := re.RawResponse().StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	var te testError
	if !errors.As(err, &te) {
		t.Fatal("unwrap failed")
	}
}

func TestDelay(t *testing.T) {
	if err := Delay(context.Background(), 5*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Delay(ctx, 5*time.Minute); err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestGetJSON(t *testing.T) {
	j, err := GetJSON(&http.Response{Body: http.NoBody})
	if !errors.Is(err, ErrNoBody) {
		t.Fatal(err)
	}
	if j != nil {
		t.Fatal("expected nil json")
	}
	j, err = GetJSON(&http.Response{Body: ioutil.NopCloser(strings.NewReader(`{ "foo": "bar" }`))})
	if err != nil {
		t.Fatal(err)
	}
	if v := j["foo"]; v != "bar" {
		t.Fatalf("unexpected value %s", v)
	}
}
