//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func reset(t *testing.T) {
	err := ResetProxy(&RecordingOptions{TestInstance: t})
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

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/"

	err = AddURISanitizer("https://replacement.com/", srvURL, &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Greater(t, len(data.Entries), 0)
	require.Equal(t, data.Entries[0].RequestURI, "https://replacement.com/")
}

func TestHeaderRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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

	err = AddHeaderRegexSanitizer("testproxy-header", "Sanitized", "", &RecordingOptions{TestInstance: t})
	require.NoError(t, err)

	err = AddHeaderRegexSanitizer("FakeStorageLocation", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.blob\\.core\\.windows\\.net", &RecordingOptions{TestInstance: t})
	require.NoError(t, err)

	// This is the only failing one
	err = AddHeaderRegexSanitizer("ComplexRegex", "Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", &RecordingOptions{GroupForReplace: "account", TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "Sanitized", data.Entries[0].RequestHeaders["testproxy-header"])
	require.Equal(t, "Sanitized", data.Entries[0].RequestHeaders["fakestoragelocation"])
	require.Equal(t, "https://Sanitized.table.core.windows.net", data.Entries[0].RequestHeaders["complexregex"])
}

func TestBodyKeySanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyKeySanitizer("$.Tag", "Sanitized", "", &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "Sanitized", data.Entries[0].ResponseBodyByValue("Tag"))
}

func TestBodyRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))

	err = AddBodyRegexSanitizer("Sanitized", "Value", &RecordingOptions{TestInstance: t})
	require.NoError(t, err)
	err = AddBodyRegexSanitizer("Sanitized", "https\\:\\/\\/(?<account>[a-z]+)\\.(?:table|blob|queue)\\.core\\.windows\\.net", &RecordingOptions{GroupForReplace: "account", TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "storageaccount", data.Entries[0].ResponseBody)
	require.NotContains(t, "Value", data.Entries[0].ResponseBody)
}

func TestRemoveHeaderSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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

	err = AddRemoveHeaderSanitizer([]string{"ComplexRegexRemove", "FakeStorageLocation"}, &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, []string{"ComplexRegexRemove", "FakeStorageLocation"}, data.Entries[0].ResponseHeaders)
}

func TestContinuationSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))

	err = AddContinuationSanitizer("Location", "Sanitized", true, &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "Location", data.Entries[0].ResponseHeaders)
	require.NotContains(t, "Location", data.Entries[0].ResponseHeaders)
}

func TestGeneralRegexSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddGeneralRegexSanitizer("Sanitized", "Value", &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.NotContains(t, "Value", data.Entries[0].ResponseBody)
}

func TestOAuthResponseSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	srvURL := "http://host.docker.internal:8080/uri-sanitizer"

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, srvURL)
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddOAuthResponseSanitizer(&RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
}

func TestUriSubscriptionIdSanitizer(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://management.azure.com/subscriptions/12345678-1234-1234-5678-123456789010/providers/Microsoft.ContainerRegistry/checkNameAvailability?api-version=2019-05-01")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	err = AddURISubscriptionIDSanitizer("", &RecordingOptions{TestInstance: t})
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, "https://management.azure.com/", data.Entries[0].RequestURI)
}

func TestResetSanitizers(t *testing.T) {
	defer reset(t)

	err := Start(t, packagePath, nil)
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, data.Entries[0].RequestHeaders["fakestoragelocation"], "https://fakeaccount.blob.core.windows.net")
}

func TestSingleTestSanitizer(t *testing.T) {
	// The first iteration, add a sanitizer for just that test. The
	// second iteration, verify that the sanitizer was not applied.
	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("%s-%d", t.Name(), i), func(t *testing.T) {
			err := Start(t, packagePath, nil)
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
			byteValue, err := ioutil.ReadAll(jsonFile)
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

func TestRestartProxySingleTest(t *testing.T) {
	// The first iteration, add a sanitizer for all scenarios, do not call stop.
	// The second iteration, add a sanitizer for just one scenario,
	// reset just the test instance.
	var firstTestInst *testing.T
	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("%s-%d", t.Name(), i), func(t *testing.T) {
			err := Start(t, packagePath, nil)
			require.NoError(t, err)

			if i == 0 {
				firstTestInst = t
			} else {
				err = AddGeneralRegexSanitizer("specific", "sample", &RecordingOptions{TestInstance: t})
				require.NoError(t, err)

				err = ResetProxy(&RecordingOptions{TestInstance: t})
				require.NoError(t, err)
				err = Stop(t, nil)
				require.NoError(t, err)
			}
		})
	}
	err := Stop(firstTestInst, nil)
	require.NoError(t, err)
}
