// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type recordingSanitizerTests struct {
	suite.Suite
}
=======
	chk "gopkg.in/check.v1"
)

type recordingSanitizerTests struct{}

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

var _ = chk.Suite(&recordingSanitizerTests{})
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type recordingSanitizerTests struct {
	suite.Suite
}
>>>>>>> c983287ea (refactor tests to testify)

const authHeader string = "Authorization"
const customHeader1 string = "Fooheader"
const customHeader2 string = "Barheader"
const nonSanitizedHeader string = "notsanitized"

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> c983287ea (refactor tests to testify)
func TestRecordingSanitizer(t *testing.T) {
	suite.Run(t, new(recordingSanitizerTests))
}

func (s *recordingSanitizerTests) TestDefaultSanitizerSanitizesAuthHeader() {
	assert := assert.New(s.T())
<<<<<<< HEAD
=======
func (s *recordingSanitizerTests) TestDefaultSanitizerSanitizesAuthHeader(c *chk.C) {
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
>>>>>>> c983287ea (refactor tests to testify)
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
<<<<<<< HEAD
<<<<<<< HEAD
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)
=======
	r, _ := recorder.NewAsMode(getTestFileName(c, false), recorder.ModeRecording, rt)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)
>>>>>>> c983287ea (refactor tests to testify)

	DefaultSanitizer(r)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(authHeader, "superSecret")

	r.RoundTrip(req)
	r.Stop()

<<<<<<< HEAD
<<<<<<< HEAD
	assert.Equal(SanitizedValue, req.Header.Get(authHeader))

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(SanitizedValue, i.Request.Headers.Get(authHeader))
	}
}

func (s *recordingSanitizerTests) TestAddSanitizedHeadersSanitizes() {
	assert := assert.New(s.T())
=======
	c.Assert(req.Header.Get(authHeader), chk.Equals, SanitizedValue)
=======
	assert.Equal(SanitizedValue, req.Header.Get(authHeader))
>>>>>>> c983287ea (refactor tests to testify)

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(SanitizedValue, i.Request.Headers.Get(authHeader))
	}
}

<<<<<<< HEAD
func (s *recordingSanitizerTests) TestAddSanitizedHeadersSanitizes(c *chk.C) {
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func (s *recordingSanitizerTests) TestAddSanitizedHeadersSanitizes() {
	assert := assert.New(s.T())
>>>>>>> c983287ea (refactor tests to testify)
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
<<<<<<< HEAD
<<<<<<< HEAD
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)
=======
	r, _ := recorder.NewAsMode(getTestFileName(c, false), recorder.ModeRecording, rt)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)
>>>>>>> c983287ea (refactor tests to testify)

	target := DefaultSanitizer(r)
	target.AddSanitizedHeaders(customHeader1, customHeader2)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(customHeader1, "superSecret")
	req.Header.Add(customHeader2, "verySecret")
	safeValue := "safeValue"
	req.Header.Add(nonSanitizedHeader, safeValue)

	r.RoundTrip(req)
	r.Stop()

<<<<<<< HEAD
<<<<<<< HEAD
	assert.Equal(SanitizedValue, req.Header.Get(customHeader1))
	assert.Equal(SanitizedValue, req.Header.Get(customHeader2))
	assert.Equal(safeValue, req.Header.Get(nonSanitizedHeader))

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(SanitizedValue, i.Request.Headers.Get(customHeader1))
		assert.Equal(SanitizedValue, i.Request.Headers.Get(customHeader2))
		assert.Equal(safeValue, i.Request.Headers.Get(nonSanitizedHeader))
	}
}

func (s *recordingSanitizerTests) TestAddUrlSanitizerSanitizes() {
	assert := assert.New(s.T())
	secret := "secretvalue"
	secretBody := "some body content that contains a " + secret
<<<<<<< HEAD
	server, cleanup := mock.NewServer()
	server.SetResponse(mock.WithStatusCode(http.StatusCreated), mock.WithBody([]byte(secretBody)))
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	baseUrl := server.URL() + "/"

	target := DefaultSanitizer(r)
	target.AddUrlSanitizer(func(url *string) {
		*url = strings.Replace(*url, secret, SanitizedValue, -1)
	})
	target.AddBodysanitizer(func(body *string) {
		*body = strings.Replace(*body, secret, SanitizedValue, -1)
	})

	req, _ := http.NewRequest(http.MethodPost, baseUrl+secret, closerFromString(secretBody))
=======
	c.Assert(req.Header.Get(customHeader1), chk.Equals, SanitizedValue)
	c.Assert(req.Header.Get(customHeader2), chk.Equals, SanitizedValue)
	c.Assert(req.Header.Get(nonSanitizedHeader), chk.Equals, safeValue)
=======
	assert.Equal(SanitizedValue, req.Header.Get(customHeader1))
	assert.Equal(SanitizedValue, req.Header.Get(customHeader2))
	assert.Equal(safeValue, req.Header.Get(nonSanitizedHeader))
>>>>>>> c983287ea (refactor tests to testify)

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(SanitizedValue, i.Request.Headers.Get(customHeader1))
		assert.Equal(SanitizedValue, i.Request.Headers.Get(customHeader2))
		assert.Equal(safeValue, i.Request.Headers.Get(nonSanitizedHeader))
	}
}

func (s *recordingSanitizerTests) TestAddUrlSanitizerSanitizes() {
	assert := assert.New(s.T())
=======
>>>>>>> a450961e3 (fb)
	server, cleanup := mock.NewServer()
	server.SetResponse(mock.WithStatusCode(http.StatusCreated), mock.WithBody([]byte(secretBody)))
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	baseUrl := server.URL() + "/"

	target := DefaultSanitizer(r)
	target.AddUrlSanitizer(func(url *string) {
		*url = strings.Replace(*url, secret, SanitizedValue, -1)
	})
	target.AddBodysanitizer(func(body *string) {
		*body = strings.Replace(*body, secret, SanitizedValue, -1)
	})

<<<<<<< HEAD
	req, _ := http.NewRequest(http.MethodPost, sanitizedUrl+secret, nil)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	req, _ := http.NewRequest(http.MethodPost, baseUrl+secret, closerFromString(secretBody))
>>>>>>> a450961e3 (fb)

	r.RoundTrip(req)
	r.Stop()

<<<<<<< HEAD
<<<<<<< HEAD
	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.NotContains(i.Response.Body, secret)
		assert.NotContains(i.Request.URL, secret)
		assert.NotContains(i.Request.Body, secret)
		assert.Contains(i.Request.URL, SanitizedValue)
		assert.Contains(i.Request.Body, SanitizedValue)
		assert.Contains(i.Response.Body, SanitizedValue)
	}
}

func (s *recordingSanitizerTests) TearDownSuite() {
	assert := assert.New(s.T())
	// cleanup test files
	err := os.RemoveAll("testfiles")
	assert.Nil(err)
}

func getTestFileName(t *testing.T, addSuffix bool) string {
	name := "testfiles/" + t.Name()
=======
	rec, err := cassette.Load(getTestFileName(c, false))
	c.Assert(err, chk.IsNil)
=======
	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)
>>>>>>> c983287ea (refactor tests to testify)

	for _, i := range rec.Interactions {
		assert.NotContains(i.Response.Body, secret)
		assert.NotContains(i.Request.URL, secret)
		assert.NotContains(i.Request.Body, secret)
		assert.Contains(i.Request.URL, SanitizedValue)
		assert.Contains(i.Request.Body, SanitizedValue)
		assert.Contains(i.Response.Body, SanitizedValue)
	}
}

func (s *recordingSanitizerTests) TearDownSuite() {
	assert := assert.New(s.T())
	// cleanup test files
	err := os.RemoveAll("testfiles")
	assert.Nil(err)
}

<<<<<<< HEAD
func getTestFileName(c *chk.C, addSuffix bool) string {
	name := "testfiles/" + c.TestName()
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func getTestFileName(t *testing.T, addSuffix bool) string {
	name := "testfiles/" + t.Name()
>>>>>>> c983287ea (refactor tests to testify)
	if addSuffix {
		name = name + ".yaml"
	}
	return name
}

type mockRoundTripper struct {
	server *mock.Server
}

func NewMockRoundTripper(server *mock.Server) *mockRoundTripper {
	return &mockRoundTripper{server: server}
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.server.Do(req)
}
