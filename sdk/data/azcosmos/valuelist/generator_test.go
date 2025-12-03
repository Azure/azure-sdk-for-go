package main

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	input := `package testpkg

//valueList allHeaders
const (
	foo string = "foo"
	bar string = "bar"
)
`
	expected := `// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testpkg

var allHeaders []string = []string{
	foo,
	bar,
}
`
	// Note: format.Source might adjust spacing, so we'll check for key elements
	output, err := generate("test.go", input)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	outputStr := string(output)
	if outputStr != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, outputStr)
	}
}

func TestGenerate_Multiple(t *testing.T) {
	input := `package testpkg

//valueList list1
const (
	A = "a"
)

//valueList list2
const (
	B = "b"
)
`
	expected := `// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testpkg

var list1 []string = []string{
	A,
}

var list2 []string = []string{
	B,
}
`
	output, err := generate("test.go", input)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	outputStr := string(output)
	if outputStr != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, outputStr)
	}
}

func TestGenerate_NoMarker(t *testing.T) {
	input := `package testpkg

const (
	foo = "foo"
)
`
	_, err := generate("test.go", input)
	if err == nil {
		t.Error("expected error when no marker found, got nil")
	}
}
