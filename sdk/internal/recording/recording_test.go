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

	_, err = ioutil.ReadFile(target.VariablesFile)
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

func TestRecordingOptions(t *testing.T) {
	r := RecordingOptions{
		UseHTTPS: true,
	}
	require.Equal(t, r.HostScheme(), "https://localhost:5001")

	r.UseHTTPS = false
	require.Equal(t, r.HostScheme(), "http://localhost:5000")

	require.Equal(t, GetEnvVariable(t, "Nonexistentevnvar", "somefakevalue"), "somefakevalue")
	require.NotEqual(t, GetEnvVariable(t, "PROXY_CERT", "fake/path/to/proxycert"), "fake/path/to/proxycert")

	r.Init()
	require.Equal(t, r.Host, "localhost:5000")
	require.Equal(t, r.Scheme, "http")

	r.UseHTTPS = true
	r.Init()
	require.Equal(t, r.Host, "localhost:5001")
	require.Equal(t, r.Scheme, "https")
}

var packagePath = "sdk/internal/recording"

func TestStartStop(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := StartRecording(t, packagePath, nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamUriHeader, "https://www.bing.com/")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IdHeader, GetRecordingId(t))

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = StopRecording(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open("./recordings/TestStartStop.json")
	require.NoError(t, err)
	defer jsonFile.Close()
}

func TestUriSanitizer(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := StartRecording(t, packagePath, nil)
	require.NoError(t, err)

	err = AddUriSanitizer("replacement", "bing", nil)
	require.NoError(t, err)

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamUriHeader, "https://www.bing.com/")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IdHeader, GetRecordingId(t))

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, GetRecordingId(t))

	err = StopRecording(t, nil)
	require.NoError(t, err)

	// Make sure the file is there
	jsonFile, err := os.Open("./recordings/TestUriSanitizer.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	var data RecordingFileStruct
	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &data)
	require.NoError(t, err)

	require.Equal(t, data.Entries[0].RequestUri, "https://www.replacement.com/")
}

func TestProxyCert(t *testing.T) {
	_, err := getRootCas(t)
	require.NoError(t, err)

	tempProxyCert, ok := os.LookupEnv("PROXY_CERT")
	require.True(t, ok)
	err = os.Unsetenv("PROXY_CERT")
	require.NoError(t, err)

	_, err = getRootCas(t)
	require.NoError(t, err)

	err = os.Setenv("PROXY_CERT", "not/a/path.crt")
	require.NoError(t, err)
	_, err = GetHTTPClient(t)
	require.Error(t, err)

	os.Setenv("PROXY_CERT", tempProxyCert)
}

func TestStopRecordingNoStart(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	err := StopRecording(t, nil)
	require.Error(t, err)

	jsonFile, err := os.Open("./recordings/TestStopRecordingNoStart.json")
	require.Error(t, err)
	defer jsonFile.Close()
}

type RecordingFileStruct struct {
	Entries []Entry `json:"Entries"`
}

type Entry struct {
	RequestUri string `json:"RequestUri"`
}

func TestLiveModeOnly(t *testing.T) {
	LiveOnly(t)
	if GetRecordMode() == modePlayback {
		t.Fatalf("Test should not run in playback")
	}
}

func TestSleep(t *testing.T) {
	start := time.Now()
	Sleep(time.Second * 5)
	duration := time.Since(start)
	if GetRecordMode() == modePlayback {
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
	err := StartRecording(t, packagePath, nil)
	require.Error(t, err)

	recordMode = temp
}

func TestBackwardSlashPath(t *testing.T) {
	os.Setenv("AZURE_RECORD_MODE", "record")
	defer os.Unsetenv("AZURE_RECORD_MODE")

	packagePathBackslash := "sdk\\internal\\recordings"

	err := StartRecording(t, packagePathBackslash, nil)
	require.Error(t, err)
}
