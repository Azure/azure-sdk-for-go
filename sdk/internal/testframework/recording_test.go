// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type recordingTests struct {
	suite.Suite
}

func TestRecording(t *testing.T) {
	suite.Run(t, new(recordingTests))
}

func (s *recordingTests) TestInitializeRecording() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	expectedMode := Playback

	target, err := NewRecording(context, expectedMode)
	assert.Nil(err)
	assert.NotNil(target.RecordingFile)
	assert.NotNil(target.VariablesFile)
	assert.Equal(expectedMode, target.Mode)

	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TestStopDoesNotSaveVariablesWhenNoVariablesExist() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	err = target.Stop()
	assert.Nil(err)

	_, err = ioutil.ReadFile(target.VariablesFile)
	assert.Equal(true, os.IsNotExist(err))
}

func (s *recordingTests) TestRecordedVariables() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { s.T().Log(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	nonExistingEnvVar := "nonExistingEnvVar"
	expectedVariableValue := "foobar"
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	// optional variables always succeed.
	assert.Equal(expectedVariableValue, target.GetOptionalEnvVar(nonExistingEnvVar, expectedVariableValue, NoSanitization))

	// non existent variables return an error
	val, err := target.GetEnvVar(nonExistingEnvVar, NoSanitization)
	// mark test as succeeded
	assert.Equal(envNotExistsError(nonExistingEnvVar), err.Error())

	// now create the env variable and check that it can be fetched
	os.Setenv(nonExistingEnvVar, expectedVariableValue)
	defer os.Unsetenv(nonExistingEnvVar)
	val, err = target.GetEnvVar(nonExistingEnvVar, NoSanitization)
	assert.Equal(expectedVariableValue, val)

	err = target.Stop()
	assert.Nil(err)

	// check that a variables file was created with the correct variable
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[nonExistingEnvVar]
	assert.Equal(true, ok)
	assert.Equal(expectedVariableValue, actualValue)
}

func (s *recordingTests) TestRecordedVariablesSanitized() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	SanitizedStringVar := "sanitizedvar"
	SanitizedBase64StrigVar := "sanitizedbase64var"
	secret := "secretstring"
	secretBase64 := "asdfasdf=="
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	// call GetOptionalRecordedVariable with the Secret_String VariableType arg
	assert.Equal(secret, target.GetOptionalEnvVar(SanitizedStringVar, secret, Secret_String))

	// call GetOptionalRecordedVariable with the Secret_Base64String VariableType arg
	assert.Equal(secretBase64, target.GetOptionalEnvVar(SanitizedBase64StrigVar, secretBase64, Secret_Base64String))

	// Calling Stop will save the variables and apply the sanitization options
	err = target.Stop()
	assert.Nil(err)

	// check that a variables file was created with the correct variables
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[SanitizedStringVar]
	assert.Equal(true, ok)
	// the saved value is sanitized
	assert.Equal(SanitizedValue, actualValue)

	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok = variablesMap[SanitizedBase64StrigVar]
	assert.Equal(true, ok)
	// the saved value is sanitized
	assert.Equal(SanitizedBase64Value, actualValue)
}

func (s *recordingTests) TestStopSavesVariablesIfExistAndReadsPreviousVariables() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	expectedVariableName := "someVariable"
	expectedVariableValue := "foobar"
	addedVariableName := "addedVariable"
	addedVariableValue := "fizzbuzz"
	variablesMap := map[string]string{}

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	target.GetOptionalEnvVar(expectedVariableName, expectedVariableValue, NoSanitization)

	err = target.Stop()
	assert.Nil(err)

	// check that a variables file was created with the correct variable
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[expectedVariableName]
	assert.True(ok)
	assert.Equal(expectedVariableValue, actualValue)

	variablesMap = map[string]string{}
	target2, err := NewRecording(context, Playback)
	assert.Nil(err)

	// add a new variable to the existing batch
	target2.GetOptionalEnvVar(addedVariableName, addedVariableValue, NoSanitization)

	err = target2.Stop()
	assert.Nil(err)

	// check that a variables file was created with the variables loaded from the previous recording
	target2.unmarshalVariablesFile(variablesMap)
	actualValue, ok = variablesMap[addedVariableName]
	assert.Truef(ok, fmt.Sprintf("Should have found %s", addedVariableName))
	assert.Equal(addedVariableValue, actualValue)
	actualValue, ok = variablesMap[expectedVariableName]
	assert.Truef(ok, fmt.Sprintf("Should have found %s", expectedVariableName))
	assert.Equal(expectedVariableValue, actualValue)
}

func (s *recordingTests) TestUUID() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	recordedUUID1 := target.UUID()
	recordedUUID1a := target.UUID()
	assert.NotEqual(recordedUUID1.String(), recordedUUID1a.String())

	err = target.Stop()
	assert.Nil(err)

	target2, err := NewRecording(context, Playback)
	assert.Nil(err)

	recordedUUID2 := target2.UUID()

	// The two generated UUIDs should be the same since target2 loaded the saved random seed from target
	assert.Equal(recordedUUID1.String(), recordedUUID2.String())

	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TestNow() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	recordedNow1 := target.Now()

	time.Sleep(time.Millisecond * 100)

	recordedNow1a := target.Now()
	assert.Equal(recordedNow1.UnixNano(), recordedNow1a.UnixNano())

	err = target.Stop()
	assert.Nil(err)

	target2, err := NewRecording(context, Playback)
	assert.Nil(err)

	recordedNow2 := target2.Now()

	// The two generated nows should be the same since target2 loaded the saved random seed from target
	assert.Equal(recordedNow1.UnixNano(), recordedNow2.UnixNano())

	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TestGenerateAlphaNumericID() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })

	prefix := "myprefix"

	target, err := NewRecording(context, Playback)
	assert.Nil(err)

	generated1, err := target.GenerateAlphaNumericID(prefix, 10, true)

	assert.Equal(10, len(generated1))
	assert.Equal(true, strings.HasPrefix(generated1, prefix))

	generated1a, err := target.GenerateAlphaNumericID(prefix, 10, true)
	assert.NotEqual(generated1, generated1a)

	err = target.Stop()
	assert.Nil(err)

	target2, err := NewRecording(context, Playback)
	assert.Nil(err)

	generated2, err := target2.GenerateAlphaNumericID(prefix, 10, true)

	// The two generated Ids should be the same since target2 loaded the saved random seed from target
	assert.Equal(generated2, generated1)

	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TestRecordRequestsAndDoMatching() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(context, Playback)
	target.recorder.SetTransport(rt)

	path, err := target.GenerateAlphaNumericID("", 5, true)
	reqUrl := server.URL() + "/" + path

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	target.Do(req)
	err = target.Stop()
	assert.Nil(err)

	rec, err := cassette.Load(target.SessionName)
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(reqUrl, i.Request.URL)
	}

	// re-initialize the recording
	target, err = NewRecording(context, Playback)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	path, err = target.GenerateAlphaNumericID("", 5, true)
	reqUrl = server.URL() + "/" + path
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	target.Do(req)
	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TestRecordRequestsAndFailMatchingForMissingRecording() {
	assert := assert.New(s.T())
	context := NewTestContext(func(msg string) { s.T().Log(msg) }, func(msg string) { s.T().Log(msg) }, func() string { return s.T().Name() })
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(context, Playback)
	target.recorder.SetTransport(rt)

	path, err := target.GenerateAlphaNumericID("", 5, true)
	reqUrl := server.URL() + "/" + path

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	target.Do(req)
	err = target.Stop()
	assert.Nil(err)

	rec, err := cassette.Load(target.SessionName)
	assert.Nil(err)

	for _, i := range rec.Interactions {
		assert.Equal(reqUrl, i.Request.URL)
	}

	// re-initialize the recording
	target, err = NewRecording(context, Playback)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	reqUrl = server.URL() + "/" + "mismatchedRequest"
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	_, err = target.Do(req)
	assert.Equal(missingRequestError(req), err.Error())
	// mark succeeded
	err = target.Stop()
	assert.Nil(err)
}

func (s *recordingTests) TearDownSuite() {

	// cleanup test files
	err := os.RemoveAll("recordings")
	assert.Nil(s.T(), err)
}
