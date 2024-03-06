// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	context "context"
	"fmt"
	"time"

	gomock "github.com/golang/mock/gomock"
)

// Cancelled matches context.Context instances that are cancelled.
var Cancelled gomock.Matcher = ContextCancelledMatcher{true}

// NotCancelled matches context.Context instances that are not cancelled.
var NotCancelled gomock.Matcher = ContextCancelledMatcher{false}

// NotCancelledAndHasTimeout matches context.Context instances that are not cancelled
// AND were also created from NewContextForTest.
var NotCancelledAndHasTimeout gomock.Matcher = gomock.All(ContextCancelledMatcher{false}, ContextHasTestValueMatcher{})

// CancelledAndHasTimeout matches context.Context instances that are cancelled
// AND were also created from NewContextForTest.
var CancelledAndHasTimeout gomock.Matcher = gomock.All(ContextCancelledMatcher{true}, ContextHasTestValueMatcher{})

type ContextCancelledMatcher struct {
	// WantCancelled should be set if we expect the context should
	// be cancelled. If true, we check if Err() != nil, if false we check
	// that it's nil.
	WantCancelled bool
}

// Matches returns whether x is a match.
func (m ContextCancelledMatcher) Matches(x any) bool {
	ctx := x.(context.Context)

	if m.WantCancelled {
		return ctx.Err() != nil
	} else {
		return ctx.Err() == nil
	}
}

// String describes what the matcher matches.
func (m ContextCancelledMatcher) String() string {
	return fmt.Sprintf("want cancelled:%v", m.WantCancelled)
}

type ContextHasTestValueMatcher struct{}

func (m ContextHasTestValueMatcher) Matches(x any) bool {
	ctx := x.(context.Context)
	return ctx.Value(testContextKey(0)) == "correctContextWasUsed"
}

func (m ContextHasTestValueMatcher) String() string {
	return "has test context value"
}

type testContextKey int

// NewContextWithTimeoutForTests creates a context with a lower timeout than requested just to keep
// unit test times reasonable.
//
// It validates that the passed in timeout is the actual defaultCloseTimeout and also
// adds in a testContextKey(0) as a value, which can be used to verify that the context
// has been properly propagated.
func NewContextWithTimeoutForTests(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// (we're in the wrong package to share the value, but this is meant to match defaultCloseTimeout)
	if timeout != time.Minute {
		// panic'ing instead of require.Equal() otherwise I would need to take a 't' and not be signature
		// compatible with context.WithTimeout.
		panic(fmt.Sprintf("Incorrect close timeout: expected %s, actual %s", time.Minute, timeout))
	}

	parentWithValue := context.WithValue(parent, testContextKey(0), "correctContextWasUsed")

	// NOTE: if you're debugging then you might need to bump up this
	// value so you can single step.
	return context.WithTimeout(parentWithValue, time.Second)
}
