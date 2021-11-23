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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBodilessMatcher(t *testing.T) {
	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	err = AddBodilessMatcher(t, nil)
	require.NoError(t, err)

	// TODO: Add more to the test to confirm it actually works

	// 1. Add Body scrubber to remove entire body
	// 2. Send a request with a body in it

	err = AddBodyRegexSanitizer("*", "", nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
	require.Equal(t, data.Entries[0].RequestBody, "")
}

func TestAddBodilessMatcherNilTest(t *testing.T) {
	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	err = AddBodilessMatcher(nil, nil)
	require.NoError(t, err)

	err = AddBodyRegexSanitizer("*", "", nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", t.Name()))
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
	require.Equal(t, data.Entries[0].RequestBody, "")

	err = ResetSanitizers(nil)
	require.NoError(t, err)
}
