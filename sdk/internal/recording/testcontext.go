//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

type TestContext interface {
	Fail(string)
	Log(string)
	Name() string
	IsFailed() bool
}

type testContext struct {
	failed bool
	fail   Failer
	log    Logger
	name   string
}

type Failer func(string)
type Logger func(string)
type Name func() string

// NewTestContext initializes a new TestContext
func NewTestContext(failer Failer, logger Logger, name Name) TestContext {
	return &testContext{fail: failer, log: logger, name: name()}
}

// Fail calls the Failer func and makes IsFailed return true.
func (c *testContext) Fail(msg string) {
	c.failed = true
	c.fail(msg)
}

// Log calls the Logger func.
func (c *testContext) Log(msg string) {
	c.log(msg)
}

// Name calls the Name func and returns the result.
func (c *testContext) Name() string {
	return c.name
}

// IsFailed returns true if the Failer has been called.
func (c *testContext) IsFailed() bool {
	return c.failed
}
