//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

// Deprecated: only deprecated methods use this type.
type TestContext interface {
	Fail(string)
	Log(string)
	Name() string
	IsFailed() bool
}

// Deprecated: only deprecated methods use this type.
type Failer func(string)

// Deprecated: only deprecated methods use this type.
type Logger func(string)

// Deprecated: only deprecated methods use this type.
type Name func() string

// NewTestContext initializes a new TestContext
func NewTestContext(failer Failer, logger Logger, name Name) TestContext {
	panic(errUnsupportedAPI)
}
