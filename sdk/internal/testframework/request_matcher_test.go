// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/dnaeon/go-vcr/cassette"

	chk "gopkg.in/check.v1"
)

type requestMatcherTests struct{}

const matchedBody string = "Matching body."
const unMatchedBody string = "This body does not match."

// Hookup to the testing framework
var _ = chk.Suite(&requestMatcherTests{})

func (s *requestMatcherTests) TestCompareBodies(c *chk.C) {

	req := http.Request{Body: closerFromString(matchedBody)}
	recReq := cassette.Request{Body: matchedBody}

	isMatch := compareBodies(&req, recReq, c)

	c.Assert(isMatch, chk.Equals, true)

	// make the requests mis-match
	req.Body = closerFromString((unMatchedBody))

	isMatch = compareBodies(&req, recReq, c)

	c.Assert(isMatch, chk.Equals, false)
	log := c.GetTestLog()
	c.Assert(strings.HasPrefix(log, "Test recording bodies do not match"), chk.Equals, true)
}

func (s *requestMatcherTests) TestCompareHeadersIgnoresIgnoredHeaders(c *chk.C) {

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	for headerName := range ignoredHeaders {
		reqHeaders[headerName] = []string{uuid.New().String()}
		recordedHeaders[headerName] = []string{uuid.New().String()}
	}

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	isMatch := compareHeaders(&req, recReq, c)

	// All headers match
	c.Assert(isMatch, chk.Equals, true)
}

func (s *requestMatcherTests) TestCompareHeadersMatchesHeaders(c *chk.C) {

	// populate only ignored headers that do not match
	reqHeaders := make(http.Header)
	recordedHeaders := make(http.Header)
	header1 := "header1"
	headerValue := []string{"some value"}

	reqHeaders[header1] = headerValue
	recordedHeaders[header1] = headerValue

	req := http.Request{Header: reqHeaders}
	recReq := cassette.Request{Headers: recordedHeaders}

	c.Assert(
		compareHeaders(&req, recReq, c),
		chk.Equals,
		true)
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingRecHeader(c *chk.C) {

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

	c.Assert(
		compareHeaders(&req, recReq, c),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(recordingHeaderMissing+"\n", header2))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMissingReqHeader(c *chk.C) {

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

	c.Assert(
		compareHeaders(&req, recReq, c),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(requestHeaderMissing+"\n", header2))
}

func (s *requestMatcherTests) TestCompareHeadersFailsMismatchedValues(c *chk.C) {

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

	c.Assert(
		compareHeaders(&req, recReq, c),
		chk.Equals,
		false)
	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(headerValuesMismatch+"\n", header2, mismatch, headerValue))
}

func (s *requestMatcherTests) TestCompareURLs(c *chk.C) {
	scheme := "https"
	host := "foo.bar"
	req := http.Request{URL: &url.URL{Scheme: scheme, Host: host}}
	recReq := cassette.Request{URL: scheme + "://" + host}

	c.Assert(
		compareURLs(&req, recReq, c),
		chk.Equals,
		true)

	req.URL.Path = "noMatch"

	c.Assert(
		compareURLs(&req, recReq, c),
		chk.Equals,
		false)

	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(urlMismatch+"\n", req.URL.String(), recReq.URL))
}

func (s *requestMatcherTests) TestCompareMethods(c *chk.C) {
	methodPost := "POST"
	methodPatch := "PATCH"
	req := http.Request{Method: methodPost}
	recReq := cassette.Request{Method: methodPost}

	c.Assert(
		compareMethods(&req, recReq, c),
		chk.Equals,
		true)

	req.Method = methodPatch

	c.Assert(
		compareMethods(&req, recReq, c),
		chk.Equals,
		false)

	c.Assert(c.GetTestLog(), chk.Equals, fmt.Sprintf(methodMismatch+"\n", req.Method, recReq.Method))
}

func closerFromString(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}
