//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"io"
	"net/http"
	"net/url"
	"strings"
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

const matchedBody string = "Matching body."
const unMatchedBody string = "This body does not match."

func (s *requestMatcherTests) TestCompareBodies() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

	req := http.Request{Body: closerFromString(matchedBody)}
	recReq := cassette.Request{Body: matchedBody}

	isMatch := matcher.compareBodies(&req, recReq.Body)

	assert.Equal(true, isMatch)

	// make the requests mis-match
	req.Body = closerFromString((unMatchedBody))

	isMatch = matcher.compareBodies(&req, recReq.Body)

	assert.False(isMatch)
}

func newUUID(t *testing.T) string {
	u, err := uuid.New()
	if err != nil {
		t.Fatal(err)
	}
	return u.String()
}

func (s *requestMatcherTests) TestCompareHeadersIgnoresIgnoredHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	for headerName := range ignoredHeaders {
		reqHeaders[headerName] = []string{newUUID(s.T())}
		recordedHeaders[headerName] = []string{newUUID(s.T())}
	}

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	// All headers match
	assert.True(matcher.compareHeaders(&req, recReq))
}

func (s *requestMatcherTests) TestCompareHeadersMatchesHeaders() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	headerValue := []string{"some value"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	assert.True(matcher.compareHeaders(&req, recReq))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingRecHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

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

	assert.False(matcher.compareHeaders(&req, recReq))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingReqHeader() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

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

	assert.False(matcher.compareHeaders(&req, recReq))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMismatchedValues() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	matcher := defaultMatcher(context)

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

	assert.False(matcher.compareHeaders(&req, recReq))
}

func (s *requestMatcherTests) TestCompareURLs() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	scheme := "https"
	host := "foo.bar"
	req := http.Request{URL: &url.URL{Scheme: scheme, Host: host}}
	recReq := cassette.Request{URL: scheme + "://" + host}
	matcher := defaultMatcher(context)

	assert.True(matcher.compareURLs(&req, recReq.URL))

	req.URL.Path = "noMatch"

	assert.False(matcher.compareURLs(&req, recReq.URL))
}

func (s *requestMatcherTests) TestCompareMethods() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	methodPost := "POST"
	methodPatch := "PATCH"
	req := http.Request{Method: methodPost}
	recReq := cassette.Request{Method: methodPost}
	matcher := defaultMatcher(context)

	assert.True(matcher.compareMethods(&req, recReq.Method))

	req.Method = methodPatch

	assert.False(matcher.compareMethods(&req, recReq.Method))
}

func closerFromString(content string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(content))
}
