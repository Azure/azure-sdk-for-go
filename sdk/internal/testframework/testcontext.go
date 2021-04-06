// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import ()

type TestContext interface {
	Fail(string)
	Log(string)
	Name() string
}

type testContext struct {
	fail Failer
	log  Logger
	name string
}

type Failer func(string)
type Logger func(string)
type Name func() string

func NewTestContext(f Failer, l Logger, name Name) TestContext {
	return &testContext{fail: f, log: l, name: name()}
}

func (c *testContext) Fail(msg string) {
	c.fail(msg)
	panic(msg)
}
func (c *testContext) Log(msg string) {
	c.log(msg)
}
func (c *testContext) Name() string {
	return c.name
}
