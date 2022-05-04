//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package diag

import (
	"regexp"
	"strings"
	"testing"
)

func TestCallerBasic(t *testing.T) {
	c := Caller(0)
	matched, err := regexp.MatchString(`/diag_test.go:\d+$`, c)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatalf("got %s", c)
	}
}

func TestCallerSkipFrame(t *testing.T) {
	c := Caller(1)
	matched, err := regexp.MatchString(`/testing.go:\d+$`, c)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatalf("got %s", c)
	}
}

func TestStackTraceBasic(t *testing.T) {
	trace := StackTrace(0, 1)
	trace = strings.TrimSpace(trace)
	parts := strings.Split(trace, "\n")
	const topFrame = "runtime.Callers()"
	if parts[0] != topFrame {
		t.Fatalf("got %s, expected %s", parts[0], topFrame)
	}
}

func TestStackTraceSkipFrame(t *testing.T) {
	trace := StackTrace(1, 1)
	trace = strings.TrimSpace(trace)
	parts := strings.Split(trace, "\n")
	const topFrame = "diag.StackTrace()"
	if strings.LastIndex(parts[0], topFrame) == -1 {
		t.Fatalf("%s didn't end with %s", parts[0], topFrame)
	}
}

func TestStackTraceFrameCount(t *testing.T) {
	trace := StackTrace(0, 5)
	trace = strings.TrimSpace(trace)
	parts := strings.Split(trace, "\n")
	// five stack frames, each is two lines, total 10 parts
	const length = 10
	if l := len(parts); l != length {
		t.Fatalf("expected %d frames, got %d", length, l)
	}
}
