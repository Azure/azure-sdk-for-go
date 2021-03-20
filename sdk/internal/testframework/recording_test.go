// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/dnaeon/go-vcr/cassette"

	chk "gopkg.in/check.v1"
)

type recordingTests struct{}

var _ = chk.Suite(&recordingTests{})

func (r *recordingTests) Test_InitializeRecording(c *chk.C) {

	expectedMode := Playback

	target, err := NewRecording(c, expectedMode)
	c.Assert(err, chk.IsNil)
	c.Assert(target.RecordingFile, chk.NotNil)
	c.Assert(target.VariablesFile, chk.NotNil)
	c.Assert(target.Mode, chk.Equals, expectedMode)

	err = target.Stop()
	c.Assert(err, chk.IsNil)
}

func (r *recordingTests) Test_StopDoesNotSaveVariablesWhenNoVariablesExist(c *chk.C) {

	target, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	err = target.Stop()
	c.Assert(err, chk.IsNil)

	_, err = ioutil.ReadFile(target.VariablesFile)
	c.Assert(os.IsNotExist(err), chk.Equals, true)
}

func (r *recordingTests) Test_RecordedVariables(c *chk.C) {

	nonExistingEnvVar := "nonExistingEnvVar"
	expectedVariableValue := "foobar"
	variablesMap := map[string]string{}

	target, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	// optional variables always succeed.
	c.Assert(target.GetOptionalRecordedVariable(nonExistingEnvVar, expectedVariableValue), chk.Equals, expectedVariableValue)

	// non existent variables panic
	c.Assert(func() { target.GetRecordedVariable(nonExistingEnvVar) }, chk.Panics, envNotExistsError(nonExistingEnvVar))

	// now create the env variable and check that it can be fetched
	os.Setenv(nonExistingEnvVar, expectedVariableValue)
	defer os.Unsetenv(nonExistingEnvVar)
	c.Assert(target.GetRecordedVariable(nonExistingEnvVar), chk.Equals, expectedVariableValue)

	err = target.Stop()
	c.Assert(err, chk.IsNil)

	// check that a variables file was created with the correct variable
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[nonExistingEnvVar]
	c.Assert(ok, chk.Equals, true)
	c.Assert(actualValue, chk.Equals, expectedVariableValue)
}

func (r *recordingTests) Test_RecordedVariablesSanitized(c *chk.C) {

	SanitizedStringVar := "sanitizedvar"
	SanitizedBase64StrigVar := "sanitizedbase64var"
	secret := "secretstring"
	secretBase64 := "asdfasdf=="
	variablesMap := map[string]string{}

	target, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	// call GetOptionalRecordedVariable with the Secret_String VariableType arg
	c.Assert(target.GetOptionalRecordedVariable(SanitizedStringVar, secret, Secret_String), chk.Equals, secret)

	// call GetOptionalRecordedVariable with the Secret_Base64String VariableType arg
	c.Assert(target.GetOptionalRecordedVariable(SanitizedBase64StrigVar, secretBase64, Secret_Base64String), chk.Equals, secretBase64)

	// Calling Stop will save the variables and apply the sanitization options
	err = target.Stop()
	c.Assert(err, chk.IsNil)

	// check that a variables file was created with the correct variables
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[SanitizedStringVar]
	c.Assert(ok, chk.Equals, true)
	// the saved value is sanitized
	c.Assert(actualValue, chk.Equals, SanitizedValue)

	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok = variablesMap[SanitizedBase64StrigVar]
	c.Assert(ok, chk.Equals, true)
	// the saved value is sanitized
	c.Assert(actualValue, chk.Equals, SanitizedBase64Value)
}

func (r *recordingTests) Test_StopSavesVariablesIfExistAndReadsPreviousVariables(c *chk.C) {

	expectedVariableName := "someVariable"
	expectedVariableValue := "foobar"
	variablesMap := map[string]string{}

	target, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	target.GetOptionalRecordedVariable(expectedVariableName, expectedVariableValue)

	err = target.Stop()
	c.Assert(err, chk.IsNil)

	// check that a variables file was created with the correct variable
	target.unmarshalVariablesFile(variablesMap)
	actualValue, ok := variablesMap[expectedVariableName]
	c.Assert(ok, chk.Equals, true)
	c.Assert(actualValue, chk.Equals, expectedVariableValue)

	variablesMap = map[string]string{}
	target2, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	err = target.Stop()
	c.Assert(err, chk.IsNil)

	// check that a variables file was created with the variables loaded from the previous recording
	target2.unmarshalVariablesFile(variablesMap)
	actualValue, ok = variablesMap[expectedVariableName]
	c.Assert(ok, chk.Equals, true)
	c.Assert(actualValue, chk.Equals, expectedVariableValue)
}

func (r *recordingTests) Test_GenerateAlphaNumericId(c *chk.C) {

	prefix := "myprefix"

	target, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	generated1 := target.GenerateAlphaNumericId(prefix, 10, true)

	c.Assert(len(generated1), chk.Equals, 10)
	c.Assert(strings.HasPrefix(generated1, prefix), chk.Equals, true)

	generated1a := target.GenerateAlphaNumericId(prefix, 10, true)
	c.Assert(generated1a, chk.Not(chk.Equals), generated1)

	err = target.Stop()
	c.Assert(err, chk.IsNil)

	target2, err := NewRecording(c, Playback)
	c.Assert(err, chk.IsNil)

	generated2 := target2.GenerateAlphaNumericId(prefix, 10, true)

	// The two generated Ids should be the same since target2 loaded the saved random seed from target
	c.Assert(generated1, chk.Equals, generated2)

	err = target.Stop()
	c.Assert(err, chk.IsNil)
}

func (s *recordingTests) TestRecordRequestsAndDoMatching(c *chk.C) {
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(c, Playback)
	target.recorder.SetTransport(rt)

	reqUrl := server.URL() + "/" + target.GenerateAlphaNumericId("", 5, true)

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	target.Do(req)
	err = target.Stop()
	c.Assert(err, chk.IsNil)

	rec, err := cassette.Load(target.SessionName)
	c.Assert(err, chk.IsNil)

	for _, i := range rec.Interactions {
		c.Assert(i.Request.URL, chk.Equals, reqUrl)
	}

	// re-initialize the recording
	target, err = NewRecording(c, Playback)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	reqUrl = server.URL() + "/" + target.GenerateAlphaNumericId("", 5, true)
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	target.Do(req)
	err = target.Stop()
	c.Assert(err, chk.IsNil)
}

func (s *recordingTests) TestRecordRequestsAndFailMatchingForMissingRecording(c *chk.C) {
	server, cleanup := mock.NewServer()
	server.SetResponse()
	defer cleanup()
	rt := NewMockRoundTripper(server)

	target, err := NewRecording(c, Playback)
	target.recorder.SetTransport(rt)

	reqUrl := server.URL() + "/" + target.GenerateAlphaNumericId("", 5, true)

	req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)

	// record the request
	target.Do(req)
	err = target.Stop()
	c.Assert(err, chk.IsNil)

	rec, err := cassette.Load(target.SessionName)
	c.Assert(err, chk.IsNil)

	for _, i := range rec.Interactions {
		c.Assert(i.Request.URL, chk.Equals, reqUrl)
	}

	// re-initialize the recording
	target, err = NewRecording(c, Playback)
	target.recorder.SetTransport(rt)

	// re-create the random url using the recorded variables
	reqUrl = server.URL() + "/" + "mismatchedRequest"
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)

	// playback the request
	c.Assert(func() { target.Do(req) }, chk.Panics, missingRequestError(req))
	c.Succeed()
	err = target.Stop()
	c.Assert(err, chk.IsNil)
}

func (r *recordingTests) TearDownSuite(c *chk.C) {

	// cleanup test files
	err := os.RemoveAll("recordings")
	c.Log(err)
}
