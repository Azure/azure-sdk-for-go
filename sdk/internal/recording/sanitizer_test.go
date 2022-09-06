//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func reset(t *testing.T) {
	err := ResetProxy(nil)
	require.NoError(t, err)
}

type RecordingFileStruct struct {
	Entries []Entry `json:"Entries"`
}

type Entry struct {
	RequestURI      string                 `json:"RequestUri"`
	RequestMethod   string                 `json:"RequestMethod"`
	RequestHeaders  map[string]string      `json:"RequestHeaders"`
	RequestBody     string                 `json:"RequestBody"`
	StatusCode      int                    `json:"StatusCode"`
	ResponseBody    interface{}            `json:"ResponseBody"` // This should be a string, but proxy saves as an object when there is no body
	ResponseHeaders map[string]interface{} `json:"ResponseHeaders"`
}

func (e Entry) ResponseBodyByValue(k string) interface{} {
	m := e.ResponseBody.(map[string]interface{})
	return m[k]
}

func TestUriSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/"

	err = AddURISanitizer("https://replacement.com/", srvURL, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, data.Entries[0].RequestURI, "https://replacement.com/")
}

func TestHeaderRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("testproxy-header", "fakevalue")
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")
	req.Header.Set("ComplexRegex", "https://fakeaccount.table.core.windows.net")

	err = AddHeaderRegexSanitizer("testproxy-header", "Sanitized", "", nil)
	require.NoError(t, err)

	err = AddHeaderRegexSanitizer("FakeStorageLocation", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.blob\\.core\\.windows\\.net", nil)
	require.NoError(t, err)

	// This is the only failing one
	err = AddHeaderRegexSanitizer("ComplexRegex", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", &RecordingOptions{GroupForReplace: "account"})
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "Sanitized", data.Entries[0].RequestHeaders["testproxy-header"])
	require.Equal(t, "Sanitized", data.Entries[0].RequestHeaders["fakestoragelocation"])
	require.Equal(t, "https://Sanitized.table.core.windows.net", data.Entries[0].RequestHeaders["complexregex"])
}

func TestBodyKeySanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(t, err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyKeySanitizer("$.Tag", "Sanitized", "", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "Sanitized", data.Entries[0].ResponseBodyByValue("Tag"))
}

func TestBodyRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(t, err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyRegexSanitizer("Sanitized", "Value", nil)
	require.NoError(t, err)
	err = AddBodyRegexSanitizer("Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", &RecordingOptions{GroupForReplace: "account"})
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "storageaccount", data.Entries[0].ResponseBody)
	require.NotContains(t, "Value", data.Entries[0].ResponseBody)
}

func TestRemoveHeaderSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")
	req.Header.Set("ComplexRegexRemove", "https://fakeaccount.table.core.windows.net")

	err = AddRemoveHeaderSanitizer([]string{"ComplexRegexRemove", "FakeStorageLocation"}, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, []string{"ComplexRegexRemove", "FakeStorageLocation"}, data.Entries[0].ResponseHeaders)
}

func TestContinuationSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("Location", "/posts/2")

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(t, err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddContinuationSanitizer("Location", "Sanitized", true, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	req, err = http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("Location", "/posts/3")

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "Location", data.Entries[0].ResponseHeaders)
	require.NotContains(t, "Location", data.Entries[0].ResponseHeaders)
}

func TestGeneralRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddGeneralRegexSanitizer("Sanitized", "Value", nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "Value", data.Entries[0].ResponseBody)
}

func TestOAuthResponseSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddOAuthResponseSanitizer(nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
}

func TestUriSubscriptionIdSanitizer(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://management.azure.com/subscriptions/12345678-1234-1234-5678-123456789010/providers/Microsoft.ContainerRegistry/checkNameAvailability?api-version=2019-05-01")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddURISubscriptionIDSanitizer("", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "https://management.azure.com/", data.Entries[0].RequestURI)
}

func TestResetSanitizers(t *testing.T) {
	defer reset(t)

	err := ResetProxy(nil)
	require.NoError(t, err)

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")

	// Add a sanitizer
	err = AddRemoveHeaderSanitizer([]string{"FakeStorageLocation"}, &RecordingOptions{TestInstance: t})
	require.NoError(t, err)

	// Remove all sanitizers
	err = ResetProxy(&RecordingOptions{TestInstance: t})
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, data.Entries[0].RequestHeaders["fakestoragelocation"], "https://fakeaccount.blob.core.windows.net")
}

func TestSingleTestSanitizer(t *testing.T) {
	err := ResetProxy(nil)
	require.NoError(t, err)

	// The first iteration, add a sanitizer for just that test. The
	// second iteration, verify that the sanitizer was not applied.
	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("%s-%d", t.Name(), i), func(t *testing.T) {
			err = Start(t, packagePath, nil)
			require.NoError(t, err)

			if i == 0 {
				// The first time we'll set a per-test sanitizer
				// Add a sanitizer
				err = AddRemoveHeaderSanitizer([]string{"FakeStorageLocation"}, &RecordingOptions{TestInstance: t})
				require.NoError(t, err)
			}

			srvURL := "http://host.docker.internal:8080/uri-sanitizer"

			client, err := GetHTTPClient(t)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "https://localhost:5001", nil)
			require.NoError(t, err)

			req.Header.Set(UpstreamURIHeader, srvURL)
			req.Header.Set(ModeHeader, GetRecordMode())
			req.Header.Set(IDHeader, GetRecordingId(t))
			req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")

			_, err = client.Do(req)
			require.NoError(t, err)

			err = Stop(t, nil)
			require.NoError(t, err)

			// Read the file
			jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
			require.NoError(t, err)
			defer jsonFile.Close()

			var data RecordingFileStruct
			byteValue, err := io.ReadAll(jsonFile)
			require.NoError(t, err)
			err = json.Unmarshal(byteValue, &data)
			require.NoError(t, err)

			if i == 0 {
				require.NotContains(t, data.Entries[0].RequestHeaders, "fakestoragelocation")
			} else {
				require.Equal(t, data.Entries[0].RequestHeaders["fakestoragelocation"], "https://fakeaccount.blob.core.windows.net")
			}
		})
	}
}
