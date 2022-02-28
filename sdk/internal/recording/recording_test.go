//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecordingOptions(t *testing.T) {
	r := RecordingOptions{
		UseHTTPS: true,
	}
	require.Equal(t, r.baseURL(), "https://localhost:5001")

	r.UseHTTPS = false
	require.Equal(t, r.baseURL(), "http://localhost:5000")

	require.Equal(t, GetEnvVariable("Nonexistentevnvar", "somefakevalue"), "somefakevalue")
	temp := recordMode
	recordMode = RecordingMode
	require.NotEqual(t, GetEnvVariable("PROXY_CERT", "fake/path/to/proxycert"), "fake/path/to/proxycert")
	recordMode = temp

	r.UseHTTPS = false
	require.Equal(t, r.baseURL(), "http://localhost:5000")

	r.UseHTTPS = true
	require.Equal(t, r.baseURL(), "https://localhost:5001")
}

var packagePath = "sdk/internal/recording/testdata"

func TestStartStop(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://www.bing.com/")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open("./testdata/recordings/TestStartStop.json")
	require.NoError(t, err)
	defer jsonFile.Close()
}

func TestStartStopRecordingClient(t *testing.T) {
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", nil)
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
	defer func() {
		err = jsonFile.Close()
		require.NoError(t, err)
		err = os.Remove(jsonFile.Name())
		require.NoError(t, err)
	}()

	var data RecordingFileStruct
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
	require.Equal(t, "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", data.Entries[0].RequestURI)
	require.Equal(t, req.URL.String(), "https://localhost:5001/acr/v1/some_registry/_tags")
}

func TestStopRecordingNoStart(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := Stop(t, nil)
	require.Error(t, err)

	jsonFile, err := os.Open("./testdata/recordings/TestStopRecordingNoStart.json")
	require.Error(t, err)
	defer jsonFile.Close()
}

func TestLiveModeOnly(t *testing.T) {
	LiveOnly(t)
	if GetRecordMode() == PlaybackMode {
		t.Fatalf("Test should not run in playback")
	}
}

func TestSleep(t *testing.T) {
	start := time.Now()
	Sleep(time.Second * 5)
	duration := time.Since(start)
	if GetRecordMode() == PlaybackMode {
		if duration > (time.Second * 1) {
			t.Fatalf("Sleep took longer than five seconds")
		}
	} else {
		if duration < (time.Second * 1) {
			t.Fatalf("Sleep took less than five seconds")
		}
	}
}

func TestBadAzureRecordMode(t *testing.T) {
	temp := recordMode

	recordMode = "badvalue"
	err := Start(t, packagePath, nil)
	require.Error(t, err)

	recordMode = temp
}

func TestBackwardSlashPath(t *testing.T) {
	t.Skip("Temporarily skipping due to changes in test-proxy.")
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	packagePathBackslash := "sdk\\internal\\recording\\testdata"

	err := Start(t, packagePathBackslash, nil)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)
}

func TestLiveOnly(t *testing.T) {
	require.Equal(t, IsLiveOnly(t), false)
	LiveOnly(t)
	require.Equal(t, IsLiveOnly(t), true)
}

func TestHostAndScheme(t *testing.T) {
	r := RecordingOptions{UseHTTPS: true}
	require.Equal(t, r.scheme(), "https")
	require.Equal(t, r.host(), "localhost:5001")

	r.UseHTTPS = false
	require.Equal(t, r.scheme(), "http")
	require.Equal(t, r.host(), "localhost:5000")
}

func TestFindProxyCertLocation(t *testing.T) {
	savedValue, ok := os.LookupEnv("PROXY_CERT")
	if ok {
		defer os.Setenv("PROXY_CERT", savedValue)
	}

	if ok {
		location, err := findProxyCertLocation()
		require.NoError(t, err)
		require.Contains(t, location, "dotnet-devcert.crt")
	}

	err := os.Unsetenv("PROXY_CERT")
	require.NoError(t, err)

	location, err := findProxyCertLocation()
	require.NoError(t, err)
	require.Contains(t, location, filepath.Join("eng", "common", "testproxy", "dotnet-devcert.crt"))
}

func TestVariables(t *testing.T) {
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	client, err := NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = Stop(t, &RecordingOptions{Variables: map[string]interface{}{"key1": "value1", "key2": "1"}})
	require.NoError(t, err)

	recordMode = PlaybackMode
	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	variables := GetVariables(t)
	require.Equal(t, variables["key1"], "value1")
	require.Equal(t, variables["key2"], "1")

	err = Stop(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer func() {
		err = jsonFile.Close()
		require.NoError(t, err)
		err = os.Remove(jsonFile.Name())
		require.NoError(t, err)
	}()
}
