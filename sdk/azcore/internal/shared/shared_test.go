//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"strings"
	"testing"
)

func TestNopCloser(t *testing.T) {
	nc := NopCloser(strings.NewReader("foo"))
	if err := nc.Close(); err != nil {
		t.Fatal(err)
	}
}
