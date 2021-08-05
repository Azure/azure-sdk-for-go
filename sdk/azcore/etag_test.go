// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"testing"
)

func TestEtagEquals(t *testing.T) {
	e1 := NewETag("tag")
	if e1.String() != "tag" {
		t.Fatalf("ETag values are not equal")
	}

	e2 := NewETag("\"tag\"")
	if e2.String() != "\"tag\"" {
		t.Fatalf("ETag values are not equal")
	}

	e3 := NewETag("W/\"weakETag\"")
	if e3.String() != "W/\"weakETag\"" {
		t.Fatalf("ETag values are not equal")
	}
	if !e3.IsWeak() {
		t.Fatalf("ETag is expected to be weak")
	}

	strongETag := NewETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	if strongETag.String() != "\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"" {
		t.Fatalf("Etag values are not equal")
	}

	if ETagAny().IsWeak() {
		t.Fatalf("ETagAny should not be weak")
	}
}

func TestETagWeak(t *testing.T) {
	et1 := NewETag("tag")
	if !et1.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et2 := NewETag("\"tag\"")
	if !et2.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et3 := NewETag("W/\"weakETag\"")
	if !et3.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et4 := NewETag("W/\"\"")
	if !et4.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et5 := ETagAny()
	if !et5.IsWeak() {
		t.Fatalf("Expected to be weak")
	}
}

func TestEtagEquality(t *testing.T) {
	weakTag := NewETag("W/\"\"")
	weakTag1 := NewETag("W/\"1\"")
	weakTag2 := NewETag("W/\"Two\"")
	strongTag1 := NewETag("\"1\"")
	strongTag2 := NewETag("\"Two\"")
	strongTagValidChars := NewETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	weakTagValidChars := NewETag("W/\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	if weakTag.Equals(*weakTag, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag1.Equals(*weakTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag2.Equals(*weakTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTagValidChars.Equals(*weakTagValidChars, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if !strongTag1.Equals(*strongTag1, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTag2.Equals(*strongTag2, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTagValidChars.Equals(*strongTagValidChars, Strong) {
		t.Fatalf("Expected etags to be equal")
	}

	if weakTag.Equals(*weakTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTagValidChars.Equals(*strongTagValidChars, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}

	if weakTag1.Equals(*weakTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag1.Equals(*strongTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag2.Equals(*strongTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
}
