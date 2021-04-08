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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type recordingSanitizerTests struct {
	suite.Suite
}

const authHeader string = "Authorization"
const customHeader1 string = "Fooheader"
const customHeader2 string = "Barheader"
const nonSanitizedHeader string = "notsanitized"

func TestRecordingSanitizer(t *testing.T) {
	suite.Run(t, new(recordingSanitizerTests))
}

func (s *recordingSanitizerTests) TestDefaultSanitizerSanitizesAuthHeader() {
	assert := assert.New(s.T())
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	DefaultSanitizer(r)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(authHeader, "superSecret")

	r.RoundTrip(req)
	r.Stop()

	assert.Equal(SanitizedValue, req.Header.Get(authHeader))

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(SanitizedValue, i.Request.Headers.Get(authHeader))
	}
}

func (s *recordingSanitizerTests) TestAddSanitizedHeadersSanitizes() {
	assert := assert.New(s.T())
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	target := DefaultSanitizer(r)
	target.AddSanitizedHeaders(customHeader1, customHeader2)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(customHeader1, "superSecret")
	req.Header.Add(customHeader2, "verySecret")
	safeValue := "safeValue"
	req.Header.Add(nonSanitizedHeader, safeValue)

	r.RoundTrip(req)
	r.Stop()

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
	server, cleanup := mock.NewServer()
	server.SetResponse(mock.WithStatusCode(http.StatusNoContent))
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

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

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.NotEqual(sanitizedUrl+secret, i.Request.URL)
		assert.Equal(sanitizedUrl, i.Request.URL)
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
