// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

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

func NewTestContext(failer Failer, logger Logger, name Name) TestContext {
	return &testContext{fail: failer, log: logger, name: name()}
}

func (c *testContext) Fail(msg string) {
	c.failed = true
	c.fail(msg)
}
func (c *testContext) Log(msg string) {
	c.log(msg)
}
func (c *testContext) Name() string {
	return c.name
}
func (c *testContext) IsFailed() bool {
	return c.failed
}
