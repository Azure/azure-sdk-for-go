//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package mock

import (
	"bytes"
	"testing"
)

func TestNewTrackedCloser(t *testing.T) {
	body, closed := NewTrackedCloser(bytes.NewReader([]byte{}))
	if closed() {
		t.Fatal("body wasn't closed yet")
	}
	if err := body.Close(); err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("body should be closed")
	}
}
