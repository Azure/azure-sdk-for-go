// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"errors"
	"strings"
	"testing"
)

func TestFrameError(t *testing.T) {
	err := NewFrameError(errors.New("failed"), false, 0, 0)
	trace := strings.TrimSpace(err.Error())
	parts := strings.Split(trace, "\n")
	// parts will be three:
	// failed:
	//   path/to/TestFrameError()
	//     <source code and line number>
	const length = 3
	if l := len(parts); l != length {
		t.Fatalf("expected %d frames, got %d", length, l)
	}
	if strings.LastIndex(parts[1], "TestFrameError()") == -1 {
		t.Fatalf("didn't find TestFrameError() in %s", parts[1])
	}
}

func TestFrameErrorWithStack(t *testing.T) {
	err := NewFrameError(errors.New("failed"), true, 0, 5)
	trace := strings.TrimSpace(err.Error())
	parts := strings.Split(trace, "\n")
	if l := len(parts); l < 3 {
		t.Fatalf("not enough frames, got %d", l)
	}
	// parts will be more than three but with the same top-most values as previous test
	if strings.LastIndex(parts[1], "TestFrameErrorWithStack()") == -1 {
		t.Fatalf("didn't find TestFrameErrorWithStack() in %s", parts[1])
	}
}
