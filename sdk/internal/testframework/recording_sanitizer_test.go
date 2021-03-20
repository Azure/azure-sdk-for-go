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
	chk "gopkg.in/check.v1"
)

type recordingSanitizerTests struct{}

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

var _ = chk.Suite(&recordingSanitizerTests{})

const authHeader string = "Authorization"
const customHeader1 string = "Fooheader"
const customHeader2 string = "Barheader"
const nonSanitizedHeader string = "notsanitized"

func (s *recordingSanitizerTests) TestDefaultSanitizerSanitizesAuthHeader(c *chk.C) {
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(c, false), recorder.ModeRecording, rt)

	DefaultSanitizer(r)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(authHeader, "superSecret")

	r.RoundTrip(req)
	r.Stop()

	c.Assert(req.Header.Get(authHeader), chk.Equals, SanitizedValue)

	rec, err := cassette.Load(getTestFileName(c, false))
	c.Assert(err, chk.IsNil)

	for _, i := range rec.Interactions {
		c.Assert(i.Request.Headers.Get(authHeader), chk.Equals, SanitizedValue)
	}
}

func (s *recordingSanitizerTests) TestAddSanitizedHeadersSanitizes(c *chk.C) {
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(c, false), recorder.ModeRecording, rt)

	target := DefaultSanitizer(r)
	target.AddSanitizedHeaders(customHeader1, customHeader2)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(customHeader1, "superSecret")
	req.Header.Add(customHeader2, "verySecret")
	safeValue := "safeValue"
	req.Header.Add(nonSanitizedHeader, safeValue)

	r.RoundTrip(req)
	r.Stop()

	c.Assert(req.Header.Get(customHeader1), chk.Equals, SanitizedValue)
	c.Assert(req.Header.Get(customHeader2), chk.Equals, SanitizedValue)
	c.Assert(req.Header.Get(nonSanitizedHeader), chk.Equals, safeValue)

	rec, err := cassette.Load(getTestFileName(c, false))
	c.Assert(err, chk.IsNil)

	for _, i := range rec.Interactions {
		c.Assert(i.Request.Headers.Get(customHeader1), chk.Equals, SanitizedValue)
		c.Assert(i.Request.Headers.Get(customHeader2), chk.Equals, SanitizedValue)
		c.Assert(i.Request.Headers.Get(nonSanitizedHeader), chk.Equals, safeValue)
	}
}

func (s *recordingSanitizerTests) TestAddUrlSanitizerSanitizes(c *chk.C) {
	server, cleanup := mock.NewServer()
	server.SetResponse(mock.WithStatusCode(http.StatusNoContent))
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(c, false), recorder.ModeRecording, rt)

	sanitizedUrl := server.URL()
	secret := "/secretvalue"

	target := DefaultSanitizer(r)
	target.AddUrlSanitizer(func(url *string) error {
		*url = strings.Replace(*url, secret, "", -1)
		return nil
	})

	req, _ := http.NewRequest(http.MethodPost, sanitizedUrl+secret, nil)

	r.RoundTrip(req)
	r.Stop()

	rec, err := cassette.Load(getTestFileName(c, false))
	c.Assert(err, chk.IsNil)

	for _, i := range rec.Interactions {
		c.Assert(i.Request.URL, chk.Not(chk.Equals), sanitizedUrl+secret)
		c.Assert(i.Request.URL, chk.Equals, sanitizedUrl)
	}
}

func (s *recordingSanitizerTests) TearDownSuite(c *chk.C) {
	// cleanup test files
	err := os.RemoveAll("testfiles")
	c.Log(err)
}

func getTestFileName(c *chk.C, addSuffix bool) string {
	name := "testfiles/" + c.TestName()
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
