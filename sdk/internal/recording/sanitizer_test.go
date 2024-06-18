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
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type sanitizerTests struct {
	suite.Suite
	proxy *TestProxyInstance
}

func TestRecordingSanitizer(t *testing.T) {
	suite.Run(t, new(sanitizerTests))
}

func (s *sanitizerTests) SetupSuite() {
	proxy, err := StartTestProxy("", nil)
	s.proxy = proxy
	require.NoError(s.T(), err)
}

func (s *sanitizerTests) TearDownSuite() {
	err1 := StopTestProxy(s.proxy)
	err2 := os.RemoveAll("testfiles")
	require.NoError(s.T(), err1)
	require.NoError(s.T(), err2)
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

func (s *sanitizerTests) TestRemoveRegisteredSanitizers() {
	require := require.New(s.T())
	defer reset(s.T())

	sanitizers := make([]string, 7)
	for i := 0; i < len(sanitizers); i++ {
		sanitizers[i] = fmt.Sprintf("AZSDK100%d", i)
	}
	require.NoError(RemoveRegisteredSanitizers(sanitizers, nil))
	require.Error(RemoveRegisteredSanitizers([]string{"unknown ID"}, nil))
}

func (s *sanitizerTests) TestUriSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/"

	err = AddURISanitizer("https://replacement.com/", srvURL, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.Equal(data.Entries[0].RequestURI, "https://replacement.com/")
}

func (s *sanitizerTests) TestHeaderRegexSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("testproxy-header", "fakevalue")
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")
	req.Header.Set("ComplexRegex", "https://fakeaccount.table.core.windows.net")

	err = AddHeaderRegexSanitizer("testproxy-header", "Sanitized", "", nil)
	require.NoError(err)

	err = AddHeaderRegexSanitizer("FakeStorageLocation", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.blob\\.core\\.windows\\.net", nil)
	require.NoError(err)

	// This is the only failing one
	opts := defaultOptions()
	opts.GroupForReplace = "account"
	err = AddHeaderRegexSanitizer("ComplexRegex", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", opts)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.Equal("Sanitized", data.Entries[0].RequestHeaders["testproxy-header"])
	require.Equal("Sanitized", data.Entries[0].RequestHeaders["fakestoragelocation"])
	require.Equal("https://Sanitized.table.core.windows.net", data.Entries[0].RequestHeaders["complexregex"])
}

func (s *sanitizerTests) TestBodyKeySanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyKeySanitizer("$.Tag", "Sanitized", "", nil)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.Equal("Sanitized", data.Entries[0].ResponseBodyByValue("Tag"))
}

func (s *sanitizerTests) TestBodyRegexSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyRegexSanitizer("Sanitized", "Value", nil)
	require.NoError(err)

	opts := defaultOptions()
	opts.GroupForReplace = "account"
	err = AddBodyRegexSanitizer("Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", opts)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.NotContains("storageaccount", data.Entries[0].ResponseBody)
	require.NotContains("Value", data.Entries[0].ResponseBody)
}

func (s *sanitizerTests) TestRemoveHeaderSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")
	req.Header.Set("ComplexRegexRemove", "https://fakeaccount.table.core.windows.net")

	err = AddRemoveHeaderSanitizer([]string{"ComplexRegexRemove", "FakeStorageLocation"}, nil)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.NotContains([]string{"ComplexRegexRemove", "FakeStorageLocation"}, data.Entries[0].ResponseHeaders)
}

func (s *sanitizerTests) TestContinuationSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("Location", "/posts/2")

	bodyValue := map[string]string{
		"key1": "value1",
	}
	marshalled, err := json.Marshal(bodyValue)
	require.NoError(err)

	req.Body = io.NopCloser(bytes.NewReader(marshalled))

	err = AddContinuationSanitizer("Location", "Sanitized", true, nil)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	req, err = http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("Location", "/posts/3")

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.NotContains("Location", data.Entries[0].ResponseHeaders)
	require.NotContains("Location", data.Entries[0].ResponseHeaders)
}

func (s *sanitizerTests) TestGeneralRegexSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	err = AddGeneralRegexSanitizer("Sanitized", "Value", nil)
	require.NoError(err)

	_, err = client.Do(req)
	require.NoError(err)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.NotContains("Value", data.Entries[0].ResponseBody)
}

func (s *sanitizerTests) TestOAuthResponseSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	err = AddOAuthResponseSanitizer(nil)
	require.NoError(err)

	_, err = client.Do(req)
	require.NoError(err)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()
	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)
}

func (s *sanitizerTests) TestUriSubscriptionIdSanitizer() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://management.azure.com/subscriptions/12345678-1234-1234-5678-123456789010/providers/Microsoft.ContainerRegistry/checkNameAvailability?api-version=2019-05-01")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	err = AddURISubscriptionIDSanitizer("", nil)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.Equal("https://management.azure.com/", data.Entries[0].RequestURI)
}

func (s *sanitizerTests) TestResetSanitizers() {
	require := require.New(s.T())
	defer reset(s.T())

	err := ResetProxy(nil)
	require.NoError(err)

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")

	opts := defaultOptions()
	opts.TestInstance = s.T()

	// Add a sanitizer
	err = AddRemoveHeaderSanitizer([]string{"FakeStorageLocation"}, opts)
	require.NoError(err)

	// Remove all sanitizers
	err = ResetProxy(opts)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)

	require.Equal(data.Entries[0].RequestHeaders["fakestoragelocation"], "https://fakeaccount.blob.core.windows.net")
}

func (s *sanitizerTests) TestSingleTestSanitizer() {
	require := require.New(s.T())
	err := ResetProxy(nil)
	require.NoError(err)

	// The first iteration, add a sanitizer for just that test. The
	// second iteration, verify that the sanitizer was not applied.
	for i := 0; i < 2; i++ {
		s.T().Run(fmt.Sprintf("%s-%d", s.T().Name(), i), func(t *testing.T) {
			err = Start(t, packagePath, nil)
			require.NoError(err)

			if i == 0 {
				// The first time we'll set a per-test sanitizer
				// Add a sanitizer
				opts := defaultOptions()
				opts.TestInstance = t
				err = AddRemoveHeaderSanitizer([]string{"FakeStorageLocation"}, opts)
				require.NoError(err)
			}

			srvURL := "http://host.docker.internal:8080/uri-sanitizer"

			client, err := GetHTTPClient(t)
			require.NoError(err)

			req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
			require.NoError(err)

			req.Header.Set(UpstreamURIHeader, srvURL)
			req.Header.Set(ModeHeader, GetRecordMode())
			req.Header.Set(IDHeader, GetRecordingId(t))
			req.Header.Set("FakeStorageLocation", "https://fakeaccount.blob.core.windows.net")

			_, err = client.Do(req)
			require.NoError(err)

			err = Stop(t, nil)
			require.NoError(err)

			// Read the file
			jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
			require.NoError(err)
			defer jsonFile.Close()

			var data RecordingFileStruct
			byteValue, err := io.ReadAll(jsonFile)
			require.NoError(err)
			err = json.Unmarshal(byteValue, &data)
			require.NoError(err)

			if i == 0 {
				require.NotContains(data.Entries[0].RequestHeaders, "fakestoragelocation")
			} else {
				require.Equal(data.Entries[0].RequestHeaders["fakestoragelocation"], "https://fakeaccount.blob.core.windows.net")
			}
		})
	}
}
