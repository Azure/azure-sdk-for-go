//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const packagePath = "sdk/internal/recording/testdata"

type recordingTests struct {
	suite.Suite
	proxy *TestProxyInstance
}

func TestRecording(t *testing.T) {
	suite.Run(t, new(recordingTests))
}

func (s *recordingTests) SetupSuite() {
	// Ignore manual start in pipeline tests, we always want to exercise install
	os.Setenv(proxyManualStartEnv, "false")
	proxy, err := StartTestProxy("", nil)
	s.proxy = proxy
	require.NoError(s.T(), err)
}

func (s *recordingTests) TearDownSuite() {
	stopErr := StopTestProxy(s.proxy)
	require.NoError(s.T(), stopErr)

	files, err := filepath.Glob("recordings/**/*.yaml")
	require.NoError(s.T(), err)
	for _, f := range files {
		err := os.Remove(f)
		require.NoError(s.T(), err)
	}
}

func (s *recordingTests) TestGetEnvVariable() {
	require := require.New(s.T())
	require.Equal(GetEnvVariable("Nonexistentevnvar", "somefakevalue"), "somefakevalue")
	temp := recordMode
	recordMode = RecordingMode
	s.T().Setenv("TEST_VARIABLE", "expected")
	require.Equal("expected", GetEnvVariable("TEST_VARIABLE", "unexpected"))
	recordMode = temp
}

func (s *recordingTests) TestRecordingOptions() {
	require := require.New(s.T())
	r := RecordingOptions{
		UseHTTPS: true,
	}
	require.Equal(r.baseURL(), "https://localhost:5001")

	r.UseHTTPS = false
	require.Equal(r.baseURL(), "http://localhost:5000")

	r = *defaultOptions()
	require.Equal(r.baseURL(), fmt.Sprintf("https://localhost:%d", r.ProxyPort))
	// ProxyPort should be generated deterministically
	require.Equal(r.ProxyPort, defaultOptions().ProxyPort)
}

func (s *recordingTests) TestStartStop() {
	require := require.New(s.T())
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://www.bing.com/")
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
}

func (s *recordingTests) TestStartStopRecordingClient() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := NewRecordingHTTPClient(s.T(), nil)
	require.NoError(err)

	req, err := http.NewRequest("POST", "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", nil)
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
	defer func() {
		err = jsonFile.Close()
		require.NoError(err)
		err = os.Remove(jsonFile.Name())
		require.NoError(err)
	}()

	var data RecordingFileStruct
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(err)
	require.Equal(
		fmt.Sprintf("https://%s.azurecr.io/acr/v1/some_registry/_tags", SanitizedValue),
		data.Entries[0].RequestURI,
	)
	require.Equal(req.URL.String(), resp.Request.URL.String())
}

func (s *recordingTests) TestStopRecordingNoStart() {
	require := require.New(s.T())
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := Stop(s.T(), nil)
	require.Error(err)

	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.Error(err)
	defer jsonFile.Close()
}

func (s *recordingTests) TestLiveModeOnly() {
	LiveOnly(s.T())
	if GetRecordMode() == PlaybackMode {
		s.T().Fatalf("Test should not run in playback")
	}
}

func (s *recordingTests) TestSleep() {
	start := time.Now()
	Sleep(time.Millisecond * 100)
	duration := time.Since(start)
	if GetRecordMode() == PlaybackMode {
		if duration >= (time.Millisecond * 50) {
			s.T().Fatalf("Sleep took at least 50ms")
		}
	} else {
		if duration < (time.Second * 50) {
			s.T().Fatalf("Sleep took less than 50ms")
		}
	}
}

func (s *recordingTests) TestBadAzureRecordMode() {
	require := require.New(s.T())
	temp := recordMode

	recordMode = "badvalue"
	err := Start(s.T(), packagePath, nil)
	require.Error(err)

	recordMode = temp
}

func (s *recordingTests) TestBackwardSlashPath() {
	s.T().Skip("Temporarily skipping due to changes in test-proxy.")

	require := require.New(s.T())
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	packagePathBackslash := "sdk\\internal\\recording\\testdata"

	err := Start(s.T(), packagePathBackslash, nil)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)
}

func (s *recordingTests) TestLiveOnly() {
	require := require.New(s.T())
	require.Equal(IsLiveOnly(s.T()), false)
	LiveOnly(s.T())
	require.Equal(IsLiveOnly(s.T()), true)
}

func (s *recordingTests) TestGitRootDetection() {
	require := require.New(s.T())
	cwd, err := os.Getwd()
	require.NoError(err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(err)

	parentDir := filepath.Dir(gitRoot)
	_, err = getGitRoot(parentDir)
	require.Error(err)
}

func (s *recordingTests) TestRecordingAssetConfigNotExist() {
	require := require.New(s.T())
	absPath, relPath, err := getAssetsConfigLocation(".")
	require.NoError(err)
	require.Equal("", absPath)
	require.Equal("", relPath)
}

func (s *recordingTests) TestRecordingAssetConfigOutOfBounds() {
	require := require.New(s.T())
	cwd, err := os.Getwd()
	require.NoError(err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(err)
	parentDir := filepath.Dir(gitRoot)

	absPath, err := findAssetsConfigFile(parentDir, gitRoot)
	require.NoError(err)
	require.Equal("", absPath)
}

func (s *recordingTests) TestRecordingAssetConfig() {
	require := require.New(s.T())
	cases := []struct{ expectedDirectory, searchDirectory, testFileLocation string }{
		{"sdk/internal/recording", "sdk/internal/recording", recordingAssetConfigName},
		{"sdk/internal/recording", "sdk/internal/recording/", recordingAssetConfigName},
		{"sdk/internal", "sdk/internal/recording", "../" + recordingAssetConfigName},
		{"sdk/internal", "sdk/internal/recording/", "../" + recordingAssetConfigName},
	}

	cwd, err := os.Getwd()
	require.NoError(err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(err)

	for _, c := range cases {
		_ = os.Remove(c.testFileLocation)
		o, err := os.Create(c.testFileLocation)
		require.NoError(err)
		o.Close()

		absPath, relPath, err := getAssetsConfigLocation(c.searchDirectory)
		// Clean up first in case of an assertion panic
		require.NoError(os.Remove(c.testFileLocation))
		require.NoError(err)

		expected := c.expectedDirectory + string(os.PathSeparator) + recordingAssetConfigName
		expected = strings.ReplaceAll(expected, "/", string(os.PathSeparator))
		require.Equal(expected, relPath)

		absPathExpected := filepath.Join(gitRoot, expected)
		require.Equal(absPathExpected, absPath)
	}
}

func (s *recordingTests) TestFindProxyCertLocation() {
	require := require.New(s.T())
	savedValue, ok := os.LookupEnv("PROXY_CERT")
	if ok {
		defer os.Setenv("PROXY_CERT", savedValue)
	}

	if ok {
		location, err := findProxyCertLocation()
		require.NoError(err)
		require.Contains(location, "dotnet-devcert.crt")
	}

	err := os.Unsetenv("PROXY_CERT")
	require.NoError(err)

	location, err := findProxyCertLocation()
	require.NoError(err)
	require.Contains(location, filepath.Join("eng", "common", "testproxy", "dotnet-devcert.crt"))
}

func (s *recordingTests) TestVariables() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	client, err := NewRecordingHTTPClient(s.T(), nil)
	require.NoError(err)

	req, err := http.NewRequest("POST", "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", nil)
	require.NoError(err)

	resp, err := client.Do(req)
	require.NoError(err)
	require.NotNil(resp)

	require.NotNil(GetRecordingId(s.T()))

	opts := defaultOptions()
	opts.Variables = map[string]interface{}{"key1": "value1", "key2": "1"}
	err = Stop(s.T(), opts)
	require.NoError(err)

	recordMode = PlaybackMode
	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	variables := GetVariables(s.T())
	require.Equal(variables["key1"], "value1")
	require.Equal(variables["key2"], "1")

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Make sure the file is there
	jsonFile, err := os.Open(fmt.Sprintf("./testdata/recordings/%s.json", s.T().Name()))
	require.NoError(err)
	defer func() {
		err = jsonFile.Close()
		require.NoError(err)
		err = os.Remove(jsonFile.Name())
		require.NoError(err)
	}()
}

func (s *recordingTests) TestRace() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = LiveMode
	s.T().Cleanup(func() { recordMode = temp })
	for i := 0; i < 4; i++ {
		s.T().Run("", func(t *testing.T) {
			t.Parallel()
			err := Start(t, "", nil)
			require.NoError(err)
			GetRecordingId(t)
			GetVariables(t)
			IsLiveOnly(t)
			err = Stop(t, nil)
			require.NoError(err)
			LiveOnly(t)
		})
	}
}

func (s *recordingTests) TestInnerGenerateAlphaNumericID() {
	require := require.New(s.T())
	seed1 := int64(1234567)
	seed2 := int64(7654321)
	randomSource1 := rand.NewSource(seed1)
	randomSource2 := rand.NewSource(seed2)
	randomSource3 := rand.NewSource(seed2)
	rand1, err := generateAlphaNumericID("test", 10, false, randomSource1)
	require.NoError(err)
	require.Equal(10, len(rand1))
	require.Equal("test", rand1[0:4])
	rand2, err := generateAlphaNumericID("test", 10, false, randomSource2)
	require.NoError(err)
	rand3, err := generateAlphaNumericID("test", 10, false, randomSource3)
	require.NoError(err)
	require.Equal(rand2, rand3)
	require.NotEqual(rand1, rand2)
}

func (s *recordingTests) TestGenerateAlphaNumericID() {
	require := require.New(s.T())
	recordMode = RecordingMode
	err := Start(s.T(), packagePath, nil)
	require.NoError(err)
	rand1, err := GenerateAlphaNumericID(s.T(), "test", 10, false)
	require.NoError(err)
	rand2, err := GenerateAlphaNumericID(s.T(), "test", 10, false)
	require.NoError(err)
	require.NotEqual(rand1, rand2)
	err = Stop(s.T(), nil)
	require.NoError(err)
	recordMode = PlaybackMode
	err = Start(s.T(), packagePath, nil)
	require.NoError(err)
	rand3, err := GenerateAlphaNumericID(s.T(), "test", 10, false)
	require.NoError(err)
	rand4, err := GenerateAlphaNumericID(s.T(), "test", 10, false)
	require.NoError(err)
	require.Equal(rand1, rand3)
	require.Equal(rand2, rand4)
	err = Stop(s.T(), nil)
	require.NoError(err)
}
