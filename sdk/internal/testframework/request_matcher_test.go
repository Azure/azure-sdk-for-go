// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
>>>>>>> c983287ea (refactor tests to testify)
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
<<<<<<< HEAD
<<<<<<< HEAD
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type requestMatcherTests struct {
	suite.Suite
}

func TestRequestMatcher(t *testing.T) {
	suite.Run(t, new(requestMatcherTests))
}
=======
=======
	"testing"
>>>>>>> c983287ea (refactor tests to testify)

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

<<<<<<< HEAD
type requestMatcherTests struct{}
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
type requestMatcherTests struct {
	suite.Suite
}

func TestRequestMatcher(t *testing.T) {
	suite.Run(t, new(requestMatcherTests))
}
>>>>>>> c983287ea (refactor tests to testify)

const matchedBody string = "Matching body."
const unMatchedBody string = "This body does not match."

<<<<<<< HEAD
<<<<<<< HEAD
func (s *requestMatcherTests) TestCompareBodies() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
// Hookup to the testing framework
var _ = chk.Suite(&requestMatcherTests{})

func (s *requestMatcherTests) TestCompareBodies(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
func (s *requestMatcherTests) TestCompareBodies() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	req := http.Request{Body: closerFromString(matchedBody)}
	recReq := cassette.Request{Body: matchedBody}

<<<<<<< HEAD
<<<<<<< HEAD
	isMatch := compareBodies(&req, recReq, context)

	assert.Equal(true, isMatch)
<<<<<<< HEAD
=======
	isMatch := compareBodies(&req, recReq, c)
=======
	isMatch := compareBodies(&req, recReq, context)
>>>>>>> 656c2801d (eliminate dependency on go-check)

	c.Assert(isMatch, chk.Equals, true)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
>>>>>>> c983287ea (refactor tests to testify)

	// make the requests mis-match
	req.Body = closerFromString((unMatchedBody))

<<<<<<< HEAD
<<<<<<< HEAD
	isMatch = compareBodies(&req, recReq, context)

	assert.False(isMatch)
}

func (s *requestMatcherTests) TestCompareHeadersIgnoresIgnoredHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	isMatch = compareBodies(&req, recReq, c)
=======
	isMatch = compareBodies(&req, recReq, context)
>>>>>>> 656c2801d (eliminate dependency on go-check)

	assert.False(isMatch)
}

<<<<<<< HEAD
func (s *requestMatcherTests) TestCompareHeadersIgnoresIgnoredHeaders(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
func (s *requestMatcherTests) TestCompareHeadersIgnoresIgnoredHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	for headerName := range ignoredHeaders {
		reqHeaders[headerName] = []string{uuid.New().String()}
		recordedHeaders[headerName] = []string{uuid.New().String()}
	}

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	// All headers match
	assert.True(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersMatchesHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	isMatch := compareHeaders(&req, recReq, c)
=======
	isMatch := compareHeaders(&req, recReq, context)
>>>>>>> 656c2801d (eliminate dependency on go-check)

=======
>>>>>>> c983287ea (refactor tests to testify)
	// All headers match
	assert.True(compareHeaders(&req, recReq, context))
}

<<<<<<< HEAD
func (s *requestMatcherTests) TestCompareHeadersMatchesHeaders(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
func (s *requestMatcherTests) TestCompareHeadersMatchesHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	headerValue := []string{"some value"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

<<<<<<< HEAD
<<<<<<< HEAD
	assert.True(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingRecHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	c.Assert(
		compareHeaders(&req, recReq, context),
		chk.Equals,
		true)
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingRecHeader(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
	assert.True(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingRecHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	header2 := "header2"
	headerValue := []string{"some value"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	// add a new header to the just req
	reqHeaders[header2] = headerValue

<<<<<<< HEAD
<<<<<<< HEAD
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingReqHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	c.Assert(
		compareHeaders(&req, recReq, context),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(recordingHeaderMissing+"\n", header2))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingReqHeader(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingReqHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	header2 := "header2"
	headerValue := []string{"some value"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	// add a new header to just the recording
	recordedHeaders[header2] = headerValue

<<<<<<< HEAD
<<<<<<< HEAD
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMismatchedValues() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	c.Assert(
		compareHeaders(&req, recReq, context),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(requestHeaderMissing+"\n", header2))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMismatchedValues(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMismatchedValues() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	header2 := "header2"
	headerValue := []string{"some value"}
	mismatch := []string{"mismatch"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	// header names match but values are different
	recordedHeaders[header2] = headerValue
	reqHeaders[header2] = mismatch

<<<<<<< HEAD
<<<<<<< HEAD
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareURLs() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	c.Assert(
		compareHeaders(&req, recReq, context),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(headerValuesMismatch+"\n", header2, mismatch, headerValue))
}

func (s *requestMatcherTests) TestCompareURLs(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
	assert.False(compareHeaders(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareURLs() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)
	scheme := "https"
	host := "foo.bar"
	req := http.Request{URL: &url.URL{Scheme: scheme, Host: host}}
	recReq := cassette.Request{URL: scheme + "://" + host}

<<<<<<< HEAD
<<<<<<< HEAD
	assert.True(compareURLs(&req, recReq, context))

	req.URL.Path = "noMatch"

	assert.False(compareURLs(&req, recReq, context))
}

func (s *requestMatcherTests) TestCompareMethods() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
=======
	c.Assert(
		compareURLs(&req, recReq, context),
		chk.Equals,
		true)
=======
	assert.True(compareURLs(&req, recReq, context))
>>>>>>> c983287ea (refactor tests to testify)

	req.URL.Path = "noMatch"

	assert.False(compareURLs(&req, recReq, context))
}

<<<<<<< HEAD
func (s *requestMatcherTests) TestCompareMethods(c *chk.C) {
<<<<<<< HEAD
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	context := NewTestContext(func(msg string) { c.Log(msg); c.Fail() }, func(msg string) { c.Log(msg) }, func() string { return c.TestName() })
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
func (s *requestMatcherTests) TestCompareMethods() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
>>>>>>> c983287ea (refactor tests to testify)
	methodPost := "POST"
	methodPatch := "PATCH"
	req := http.Request{Method: methodPost}
	recReq := cassette.Request{Method: methodPost}

<<<<<<< HEAD
<<<<<<< HEAD
	assert.True(compareMethods(&req, recReq, context))

	req.Method = methodPatch

	assert.False(compareMethods(&req, recReq, context))
=======
	c.Assert(
		compareMethods(&req, recReq, context),
		chk.Equals,
		true)

	req.Method = methodPatch

	c.Assert(
		compareMethods(&req, recReq, context),
		chk.Equals,
		false)

	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(methodMismatch+"\n", req.Method, recReq.Method))
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	assert.True(compareMethods(&req, recReq, context))

	req.Method = methodPatch

	assert.False(compareMethods(&req, recReq, context))
>>>>>>> c983287ea (refactor tests to testify)
}

func closerFromString(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
