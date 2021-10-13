//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type sanitizerTests struct {
	suite.Suite
}

const authHeader string = "Authorization"
const customHeader1 string = "Fooheader"
const customHeader2 string = "Barheader"
const nonSanitizedHeader string = "notsanitized"

func TestRecordingSanitizer(t *testing.T) {
	suite.Run(t, new(sanitizerTests))
}

func (s *sanitizerTests) TestDefaultSanitizerSanitizesAuthHeader() {
	require := require.New(s.T())
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	defaultSanitizer(r)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(authHeader, "superSecret")

	_, err := r.RoundTrip(req)
	require.NoError(err)
	err = r.Stop()
	require.NoError(err)

	require.Equal(SanitizedValue, req.Header.Get(authHeader))

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	require.NoError(err)

	for _, i := range rec.Interactions {
		require.Equal(SanitizedValue, i.Request.Headers.Get(authHeader))
	}
}

func (s *sanitizerTests) TestAddSanitizedHeadersSanitizes() {
	require := require.New(s.T())
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	target := defaultSanitizer(r)
	target.AddSanitizedHeaders(customHeader1, customHeader2)

	req, _ := http.NewRequest(http.MethodPost, server.URL(), nil)
	req.Header.Add(customHeader1, "superSecret")
	req.Header.Add(customHeader2, "verySecret")
	safeValue := "safeValue"
	req.Header.Add(nonSanitizedHeader, safeValue)

	_, err := r.RoundTrip(req)
	require.NoError(err)
	err = r.Stop()
	require.NoError(err)

	require.Equal(SanitizedValue, req.Header.Get(customHeader1))
	require.Equal(SanitizedValue, req.Header.Get(customHeader2))
	require.Equal(safeValue, req.Header.Get(nonSanitizedHeader))

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	require.NoError(err)

	for _, i := range rec.Interactions {
		require.Equal(SanitizedValue, i.Request.Headers.Get(customHeader1))
		require.Equal(SanitizedValue, i.Request.Headers.Get(customHeader2))
		require.Equal(safeValue, i.Request.Headers.Get(nonSanitizedHeader))
	}
}

func (s *sanitizerTests) TestAddUrlSanitizerSanitizes() {
	require := require.New(s.T())
	secret := "secretvalue"
	secretBody := "some body content that contains a " + secret
	server, cleanup := mock.NewServer()
	server.SetResponse(mock.WithStatusCode(http.StatusCreated), mock.WithBody([]byte(secretBody)))
	defer cleanup()
	rt := NewMockRoundTripper(server)
	r, _ := recorder.NewAsMode(getTestFileName(s.T(), false), recorder.ModeRecording, rt)

	baseUrl := server.URL() + "/"

	target := defaultSanitizer(r)
	target.AddUrlSanitizer(func(url *string) {
		*url = strings.Replace(*url, secret, SanitizedValue, -1)
	})
	target.AddBodysanitizer(func(body *string) {
		*body = strings.Replace(*body, secret, SanitizedValue, -1)
	})

	req, _ := http.NewRequest(http.MethodPost, baseUrl+secret, closerFromString(secretBody))

	_, err := r.RoundTrip(req)
	require.NoError(err)
	err = r.Stop()
	require.NoError(err)

	rec, err := cassette.Load(getTestFileName(s.T(), false))
	require.NoError(err)

	for _, i := range rec.Interactions {
		require.NotContains(i.Response.Body, secret)
		require.NotContains(i.Request.URL, secret)
		require.NotContains(i.Request.Body, secret)
		require.Contains(i.Request.URL, SanitizedValue)
		require.Contains(i.Request.Body, SanitizedValue)
		require.Contains(i.Response.Body, SanitizedValue)
	}
}

func (s *sanitizerTests) TearDownSuite() {
	require := require.New(s.T())
	// cleanup test files
	err := os.RemoveAll("testfiles")
	require.NoError(err)
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
