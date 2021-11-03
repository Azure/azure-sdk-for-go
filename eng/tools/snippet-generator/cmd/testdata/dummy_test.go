//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testdata

import "testing"

// Snippet: Test1
func TestSomething(t *testing.T) {
	t.Skip()
}

// EndSnippet

func TestSomething2(t *testing.T) {
	t.Skip()
	// Snippet: Test2
	x := 1
	y := x + 1
	// EndSnippet
	t.Log(x + y)
}

func TestExtra(t *testing.T) {
	// Snippet: Test3
	t.Log("testing")
	// EndSnippet
}
