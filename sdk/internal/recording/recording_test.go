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

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type recordingTests struct {
	suite.Suite
}

func TestRecording(t *testing.T) {
	suite.Run(t, new(recordingTests))
}

func (s *recordingTests) TestInitializeRecording() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	expectedMode := Playback

	target, err := NewRecording(context, expectedMode)
	require.NoError(err)
	require.NotNil(target.RecordingFile)
	require.NotNil(target.VariablesFile)
	require.Equal(expectedMode, target.Mode)

	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TestStopDoesNotSaveVariablesWhenNoVariablesExist() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	err = target.Stop()
	require.NoError(err)

	_, err = os.ReadFile(target.VariablesFile)
	require.Equal(true, os.IsNotExist(err))
}

func (s *recordingTests) TestRecordedVariables() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { s.T().Log(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	nonExistingEnvVar := "nonExistingEnvVar"
	expectedVariableValue := "foobar"
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	// optional variables always succeed.
	require.Equal(expectedVariableValue, target.GetOptionalEnvVar(nonExistingEnvVar, expectedVariableValue, NoSanitization))

	// non existent variables return an error
	_, err = target.GetEnvVar(nonExistingEnvVar, NoSanitization)
	// mark test as succeeded
	require.Equal(envNotExistsError(nonExistingEnvVar), err.Error())

	// now create the env variable and check that it can be fetched
	os.Setenv(nonExistingEnvVar, expectedVariableValue)
	defer os.Unsetenv(nonExistingEnvVar)
	val, err := target.GetEnvVar(nonExistingEnvVar, NoSanitization)
	require.NoError(err)
	require.Equal(expectedVariableValue, val)

	err = target.Stop()
	require.NoError(err)

	// check that a variables file was created with the correct variable
	err = target.unmarshalVariablesFile(variablesMap)
	require.NoError(err)
	actualValue, ok := variablesMap[nonExistingEnvVar]
	require.Equal(true, ok)
	require.Equal(expectedVariableValue, actualValue)
}

func (s *recordingTests) TestRecordedVariablesSanitized() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	SanitizedStringVar := "sanitizedvar"
	SanitizedBase64StrigVar := "sanitizedbase64var"
	secret := "secretstring"
	secretBase64 := "asdfasdf=="
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	// call GetOptionalRecordedVariable with the Secret_String VariableType arg
	require.Equal(secret, target.GetOptionalEnvVar(SanitizedStringVar, secret, Secret_String))

	// call GetOptionalRecordedVariable with the Secret_Base64String VariableType arg
	require.Equal(secretBase64, target.GetOptionalEnvVar(SanitizedBase64StrigVar, secretBase64, Secret_Base64String))

	// Calling Stop will save the variables and apply the sanitization options
	err = target.Stop()
	require.NoError(err)

	// check that a variables file was created with the correct variables
	err = target.unmarshalVariablesFile(variablesMap)
	require.NoError(err)
	actualValue, ok := variablesMap[SanitizedStringVar]
	require.Equal(true, ok)
	// the saved value is sanitized
	require.Equal(SanitizedValue, actualValue)

	err = target.unmarshalVariablesFile(variablesMap)
	require.NoError(err)
	actualValue, ok = variablesMap[SanitizedBase64StrigVar]
	require.Equal(true, ok)
	// the saved value is sanitized
	require.Equal(SanitizedBase64Value, actualValue)
}

func (s *recordingTests) TestStopSavesVariablesIfExistAndReadsPreviousVariables() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	expectedVariableName := "someVariable"
	expectedVariableValue := "foobar"
	addedVariableName := "addedVariable"
	addedVariableValue := "fizzbuzz"
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	target.GetOptionalEnvVar(expectedVariableName, expectedVariableValue, NoSanitization)

	err = target.Stop()
	require.NoError(err)

	// check that a variables file was created with the correct variable
	err = target.unmarshalVariablesFile(variablesMap)
	require.NoError(err)
	actualValue, ok := variablesMap[expectedVariableName]
	require.True(ok)
	require.Equal(expectedVariableValue, actualValue)

	variablesMap = map[string]string{}
	target2, err := NewRecording(context, Playback)
	require.NoError(err)

	// add a new variable to the existing batch
	target2.GetOptionalEnvVar(addedVariableName, addedVariableValue, NoSanitization)

	err = target2.Stop()
	require.NoError(err)

	// check that a variables file was created with the variables loaded from the previous recording
	err = target2.unmarshalVariablesFile(variablesMap)
	require.NoError(err)
	actualValue, ok = variablesMap[addedVariableName]
	require.Truef(ok, fmt.Sprintf("Should have found %s", addedVariableName))
	require.Equal(addedVariableValue, actualValue)
	actualValue, ok = variablesMap[expectedVariableName]
	require.Truef(ok, fmt.Sprintf("Should have found %s", expectedVariableName))
	require.Equal(expectedVariableValue, actualValue)
}

func (s *recordingTests) TestUUID() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	recordedUUID1 := target.UUID()
	recordedUUID1a := target.UUID()
	require.NotEqual(recordedUUID1.String(), recordedUUID1a.String())

	err = target.Stop()
	require.NoError(err)

	target2, err := NewRecording(context, Playback)
	require.NoError(err)

	recordedUUID2 := target2.UUID()

	// The two generated UUIDs should be the same since target2 loaded the saved random seed from target
	require.Equal(recordedUUID1.String(), recordedUUID2.String())

	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TestNow() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	recordedNow1 := target.Now()

	time.Sleep(time.Millisecond * 100)

	recordedNow1a := target.Now()
	require.Equal(recordedNow1.UnixNano(), recordedNow1a.UnixNano())

	err = target.Stop()
	require.NoError(err)

	target2, err := NewRecording(context, Playback)
	require.NoError(err)

	recordedNow2 := target2.Now()

	// The two generated nows should be the same since target2 loaded the saved random seed from target
	require.Equal(recordedNow1.UnixNano(), recordedNow2.UnixNano())

	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TestGenerateAlphaNumericID() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	prefix := "myprefix"

	target, err := NewRecording(context, Playback)
	require.NoError(err)

	generated1, err := target.GenerateAlphaNumericID(prefix, 10, true)
	require.NoError(err)

	require.Equal(10, len(generated1))
	require.Equal(true, strings.HasPrefix(generated1, prefix))

	generated1a, err := target.GenerateAlphaNumericID(prefix, 10, true)
	require.NoError(err)
	require.NotEqual(generated1, generated1a)

	err = target.Stop()
	require.NoError(err)

	target2, err := NewRecording(context, Playback)
	require.NoError(err)

	generated2, err := target2.GenerateAlphaNumericID(prefix, 10, true)
	require.NoError(err)
	// The two generated Ids should be the same since target2 loaded the saved random seed from target
	require.Equal(generated2, generated1)

	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TestRecordRequestsAndDoMatching() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { require.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(context, Playback)
	require.NoError(err)
	target.recorder.SetTransport(rt)

	path, err := target.GenerateAlphaNumericID("", 5, true)
	require.NoError(err)
	reqUrl := server.URL() + "/" + path

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	_, err = target.Do(req)
	require.NoError(err)
	err = target.Stop()
	require.NoError(err)

	rec, err := cassette.Load(target.SessionName)
	require.NoError(err)

	for _, i := range rec.Interactions {
		require.Equal(reqUrl, i.Request.URL)
	}

	// re-initialize the recording
	target, err = NewRecording(context, Playback)
	require.NoError(err)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	path, err = target.GenerateAlphaNumericID("", 5, true)
	require.NoError(err)
	reqUrl = server.URL() + "/" + path
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	_, err = target.Do(req)
	require.NoError(err)
	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TestRecordRequestsAndFailMatchingForMissingRecording() {
	require := require.New(s.T())
	context := NewTestContext(func(msg string) { s.T().Log(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(context, Playback)
	require.NoError(err)
	target.recorder.SetTransport(rt)

	path, err := target.GenerateAlphaNumericID("", 5, true)
	require.NoError(err)
	reqUrl := server.URL() + "/" + path

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	_, err = target.Do(req)
	require.NoError(err)
	err = target.Stop()
	require.NoError(err)

	rec, err := cassette.Load(target.SessionName)
	require.NoError(err)

	for _, i := range rec.Interactions {
		require.Equal(reqUrl, i.Request.URL)
	}

	// re-initialize the recording
	target, err = NewRecording(context, Playback)
	require.NoError(err)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	reqUrl = server.URL() + "/" + "mismatchedRequest"
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	_, err = target.Do(req)
	require.Equal(missingRequestError(req), err.Error())
	// mark succeeded
	err = target.Stop()
	require.NoError(err)
}

func (s *recordingTests) TearDownSuite() {
	files, err := filepath.Glob("recordings/**/*.yaml")
	require.NoError(s.T(), err)
	for _, f := range files {
		err := os.Remove(f)
		require.NoError(s.T(), err)
	}
}

func TestGetEnvVariable(t *testing.T) {
	require.Equal(t, GetEnvVariable("Nonexistentevnvar", "somefakevalue"), "somefakevalue")
	temp := recordMode
	recordMode = RecordingMode
	t.Setenv("TEST_VARIABLE", "expected")
	require.Equal(t, "expected", GetEnvVariable("TEST_VARIABLE", "unexpected"))
	recordMode = temp
}

func TestRecordingOptions(t *testing.T) {
	r := RecordingOptions{
		UseHTTPS: true,
	}
	require.Equal(t, r.baseURL(), "https://localhost:5001")

	r.UseHTTPS = false
	require.Equal(t, r.baseURL(), "http://localhost:5000")
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
	byteValue, err := io.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)
	require.Equal(t, "https://azsdkengsys.azurecr.io/acr/v1/some_registry/_tags", data.Entries[0].RequestURI)
	require.Equal(t, resp.Request.URL.String(), "https://localhost:5001/acr/v1/some_registry/_tags")
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

func TestGitRootDetection(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(t, err)

	parentDir := filepath.Dir(gitRoot)
	_, err = getGitRoot(parentDir)
	require.Error(t, err)
}

func TestRecordingAssetConfigNotExist(t *testing.T) {
	absPath, relPath, err := getAssetsConfigLocation(".")
	require.NoError(t, err)
	require.Equal(t, "", absPath)
	require.Equal(t, "", relPath)
}

func TestRecordingAssetConfigOutOfBounds(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(t, err)
	parentDir := filepath.Dir(gitRoot)

	absPath, err := findAssetsConfigFile(parentDir, gitRoot)
	require.NoError(t, err)
	require.Equal(t, "", absPath)
}

func TestRecordingAssetConfig(t *testing.T) {
	cases := []struct{ expectedDirectory, searchDirectory, testFileLocation string }{
		{"sdk/internal/recording", "sdk/internal/recording", recordingAssetConfigName},
		{"sdk/internal/recording", "sdk/internal/recording/", recordingAssetConfigName},
		{"sdk/internal", "sdk/internal/recording", "../" + recordingAssetConfigName},
		{"sdk/internal", "sdk/internal/recording/", "../" + recordingAssetConfigName},
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(t, err)

	for _, c := range cases {
		_ = os.Remove(c.testFileLocation)
		o, err := os.Create(c.testFileLocation)
		require.NoError(t, err)
		o.Close()

		absPath, relPath, err := getAssetsConfigLocation(c.searchDirectory)
		// Clean up first in case of an assertion panic
		require.NoError(t, os.Remove(c.testFileLocation))
		require.NoError(t, err)

		expected := c.expectedDirectory + string(os.PathSeparator) + recordingAssetConfigName
		expected = strings.ReplaceAll(expected, "/", string(os.PathSeparator))
		require.Equal(t, expected, relPath)

		absPathExpected := filepath.Join(gitRoot, expected)
		require.Equal(t, absPathExpected, absPath)
	}
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

func TestRace(t *testing.T) {
	temp := recordMode
	recordMode = LiveMode
	t.Cleanup(func() { recordMode = temp })
	for i := 0; i < 4; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			err := Start(t, "", nil)
			require.NoError(t, err)
			GetRecordingId(t)
			GetVariables(t)
			IsLiveOnly(t)
			err = Stop(t, nil)
			require.NoError(t, err)
			LiveOnly(t)
		})
	}
}

func Test_generateAlphaNumericID(t *testing.T) {
	seed1 := int64(1234567)
	seed2 := int64(7654321)
	randomSource1 := rand.NewSource(seed1)
	randomSource2 := rand.NewSource(seed2)
	randomSource3 := rand.NewSource(seed2)
	rand1, err := generateAlphaNumericID("test", 10, false, randomSource1)
	require.NoError(t, err)
	require.Equal(t, 10, len(rand1))
	require.Equal(t, "test", rand1[0:4])
	rand2, err := generateAlphaNumericID("test", 10, false, randomSource2)
	require.NoError(t, err)
	rand3, err := generateAlphaNumericID("test", 10, false, randomSource3)
	require.NoError(t, err)
	require.Equal(t, rand2, rand3)
	require.NotEqual(t, rand1, rand2)
}

func TestGenerateAlphaNumericID(t *testing.T) {
	recordMode = RecordingMode
	err := Start(t, packagePath, nil)
	require.NoError(t, err)
	rand1, err := GenerateAlphaNumericID(t, "test", 10, false)
	require.NoError(t, err)
	rand2, err := GenerateAlphaNumericID(t, "test", 10, false)
	require.NoError(t, err)
	require.NotEqual(t, rand1, rand2)
	err = Stop(t, nil)
	require.NoError(t, err)
	recordMode = PlaybackMode
	err = Start(t, packagePath, nil)
	require.NoError(t, err)
	rand3, err := GenerateAlphaNumericID(t, "test", 10, false)
	require.NoError(t, err)
	rand4, err := GenerateAlphaNumericID(t, "test", 10, false)
	require.NoError(t, err)
	require.Equal(t, rand1, rand3)
	require.Equal(t, rand2, rand4)
	err = Stop(t, nil)
	require.NoError(t, err)
}
