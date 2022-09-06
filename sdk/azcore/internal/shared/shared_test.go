//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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

func TestRetryAfter(t *testing.T) {
	if RetryAfter(nil) != 0 {
		t.Fatal("expected zero duration")
	}
	resp := &http.Response{
		Header: http.Header{},
	}
	if d := RetryAfter(resp); d > 0 {
		t.Fatalf("unexpected retry-after value %d", d)
	}
	resp.Header.Set(HeaderRetryAfter, "300")
	d := RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from seconds")
	}
	if d != 300*time.Second {
		t.Fatalf("expected 300 seconds, got %d", d/time.Second)
	}
	atDate := time.Now().Add(600 * time.Second)
	resp.Header.Set(HeaderRetryAfter, atDate.Format(time.RFC1123))
	d = RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from date")
	}
	// d will not be exactly 600 seconds but it will be close
	if s := d / time.Second; s < 598 || s > 602 {
		t.Fatalf("expected ~600 seconds, got %d", s)
	}
}

func TestTypeOfT(t *testing.T) {
	if tt := TypeOfT[bool](); tt != reflect.TypeOf(true) {
		t.Fatalf("unexpected type %s", tt)
	}
	if tt := TypeOfT[int32](); tt == reflect.TypeOf(3.14) {
		t.Fatal("didn't expect types to match")
	}
}

func TestNopClosingBytesReader(t *testing.T) {
	const val1 = "the data"
	ncbr := &NopClosingBytesReader{s: []byte(val1)}
	b, err := io.ReadAll(ncbr)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != val1 {
		t.Fatalf("got %s, want %s", string(b), val1)
	}
	const val2 = "something else"
	ncbr.Set([]byte(val2))
	b, err = io.ReadAll(ncbr)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != val2 {
		t.Fatalf("got %s, want %s", string(b), val2)
	}
	if err = ncbr.Close(); err != nil {
		t.Fatal(err)
	}
	// seek to beginning and read again
	i, err := ncbr.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if i != 0 {
		t.Fatalf("got %d, want %d", i, 0)
	}
	b, err = io.ReadAll(ncbr)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != val2 {
		t.Fatalf("got %s, want %s", string(b), val2)
	}
	// seek to middle from the end
	i, err = ncbr.Seek(-4, io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	if l := int64(len(val2)) - 4; i != l {
		t.Fatalf("got %d, want %d", l, i)
	}
	b, err = io.ReadAll(ncbr)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "else" {
		t.Fatalf("got %s, want %s", string(b), "else")
	}
	// underflow
	_, err = ncbr.Seek(-int64(len(val2)+1), io.SeekCurrent)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
}
