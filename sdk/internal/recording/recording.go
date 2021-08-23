//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"gopkg.in/yaml.v2"
)

type Recording struct {
	SessionName              string
	RecordingFile            string
	VariablesFile            string
	Mode                     RecordMode
	variables                map[string]*string `yaml:"variables"`
	previousSessionVariables map[string]*string `yaml:"variables"`
	recorder                 *recorder.Recorder
	src                      rand.Source
	now                      *time.Time
	Sanitizer                *Sanitizer
	Matcher                  *RequestMatcher
	c                        TestContext
}

const (
	alphanumericBytes           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	alphanumericLowercaseBytes  = "abcdefghijklmnopqrstuvwxyz1234567890"
	randomSeedVariableName      = "randomSeed"
	nowVariableName             = "now"
	ModeEnvironmentVariableName = "AZURE_TEST_MODE"
)

// Inspired by https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type RecordMode string

const (
	Record   RecordMode = "record"
	Playback RecordMode = "playback"
	Live     RecordMode = "live"
)

type VariableType string

const (
	// NoSanitization indicates that the recorded value should not be sanitized.
	NoSanitization VariableType = "default"
	// Secret_String indicates that the recorded value should be replaced with a sanitized value.
	Secret_String VariableType = "secret_string"
	// Secret_Base64String indicates that the recorded value should be replaced with a sanitized valid base-64 string value.
	Secret_Base64String VariableType = "secret_base64String"
)

// NewRecording initializes a new Recording instance
func NewRecording(c TestContext, mode RecordMode) (*Recording, error) {
	// create recorder based on the test name, recordMode, variables, and sanitizers
	recPath, varPath := getFilePaths(c.Name())
	rec, err := recorder.NewAsMode(recPath, modeMap[mode], nil)
	if err != nil {
		return nil, err
	}

	// If the mode is set in the environment, let that override the requested mode
	// This is to enable support for nightly live test pipelines
	envMode := getOptionalEnv(ModeEnvironmentVariableName, string(mode))
	mode = RecordMode(*envMode)

	// initialize the Recording
	recording := &Recording{
		SessionName:              recPath,
		RecordingFile:            recPath + ".yaml",
		VariablesFile:            varPath,
		Mode:                     mode,
		variables:                make(map[string]*string),
		previousSessionVariables: make(map[string]*string),
		recorder:                 rec,
		c:                        c,
	}

	// Try loading the recording if it already exists to hydrate the variables
	err = recording.initVariables()
	if err != nil {
		return nil, err
	}

	// set the recorder Matcher
	recording.Matcher = defaultMatcher(c)
	rec.SetMatcher(recording.matchRequest)

	// wire up the sanitizer
	recording.Sanitizer = defaultSanitizer(rec)

	return recording, err
}

// GetEnvVar returns a recorded environment variable. If the variable is not found we return an error.
// variableType determines how the recorded variable will be saved.
func (r *Recording) GetEnvVar(name string, variableType VariableType) (string, error) {
	var err error
	result, ok := r.previousSessionVariables[name]
	if !ok || r.Mode == Live {

		result, err = getRequiredEnv(name)
		if err != nil {
			r.c.Fail(err.Error())
			return "", err
		}
		r.variables[name] = applyVariableOptions(result, variableType)
	}
	return *result, err
}

// GetOptionalEnvVar returns a recorded environment variable with a fallback default value.
// default Value configures the fallback value to be returned if the environment variable is not set.
// variableType determines how the recorded variable will be saved.
func (r *Recording) GetOptionalEnvVar(name string, defaultValue string, variableType VariableType) string {
	result, ok := r.previousSessionVariables[name]
	if !ok || r.Mode == Live {
		result = getOptionalEnv(name, defaultValue)
		r.variables[name] = applyVariableOptions(result, variableType)
	}
	return *result
}

// Do satisfies the azcore.Transport interface so that Recording can be used as the transport for recorded requests
func (r *Recording) Do(req *http.Request) (*http.Response, error) {
	resp, err := r.recorder.RoundTrip(req)
	if err == cassette.ErrInteractionNotFound {
		error := missingRequestError(req)
		r.c.Fail(error)
		return nil, errors.New(error)
	}
	return resp, err
}

// Stop stops the recording and saves them, including any captured variables, to disk
func (r *Recording) Stop() error {

	err := r.recorder.Stop()
	if err != nil {
		return err
	}
	if r.Mode == Live {
		return nil
	}

	if len(r.variables) > 0 {
		// Merge values from previousVariables that are not in variables to variables
		for k, v := range r.previousSessionVariables {
			if _, ok := r.variables[k]; ok {
				// skip variables that were new in the current session
				continue
			}
			r.variables[k] = v
		}

		// Marshal to YAML and save variables
		data, err := yaml.Marshal(r.variables)
		if err != nil {
			return err
		}

		f, err := r.createVariablesFileIfNotExists()
		if err != nil {
			return err
		}

		defer f.Close()

		// http://www.yaml.org/spec/1.2/spec.html#id2760395
		_, err = f.Write([]byte("---\n"))
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Recording) Now() time.Time {
	r.initNow()

	return *r.now
}

func (r *Recording) UUID() uuid.UUID {
	r.initRandomSource()
	u := uuid.UUID{}
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	// math/rand.Read() is no-fail so we omit any error checking.
	rnd := rand.New(r.src)
	rnd.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F // u.setVariant(ReservedRFC4122)

	var version byte = 4
	u[6] = (u[6] & 0xF) | (version << 4) // u.setVersion(4)
	return u
}

// GenerateAlphaNumericID will generate a recorded random alpha numeric id
// if the recording has a randomSeed already set, the value will be generated from that seed, else a new random seed will be used
func (r *Recording) GenerateAlphaNumericID(prefix string, length int, lowercaseOnly bool) (string, error) {

	if length <= len(prefix) {
		return "", errors.New("length must be greater than prefix")
	}

	r.initRandomSource()

	sb := strings.Builder{}
	sb.Grow(length)
	sb.WriteString(prefix)
	i := length - len(prefix) - 1
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for cache, remain := r.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.src.Int63(), letterIdxMax
		}
		if lowercaseOnly {
			if idx := int(cache & letterIdxMask); idx < len(alphanumericLowercaseBytes) {
				sb.WriteByte(alphanumericLowercaseBytes[idx])
				i--
			}
		} else {
			if idx := int(cache & letterIdxMask); idx < len(alphanumericBytes) {
				sb.WriteByte(alphanumericBytes[idx])
				i--
			}
		}
		cache >>= letterIdxBits
		remain--
	}
	str := sb.String()
	return str, nil
}

// getRequiredEnv gets an environment variable by name and returns an error if it is not found
func getRequiredEnv(name string) (*string, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		return &env, nil
	} else {
		return nil, errors.New(envNotExistsError(name))
	}
}

// getOptionalEnv gets an environment variable by name and returns the defaultValue if not found
func getOptionalEnv(name string, defaultValue string) *string {
	env, ok := os.LookupEnv(name)
	if ok {
		return &env
	} else {
		return &defaultValue
	}
}

func (r *Recording) matchRequest(req *http.Request, rec cassette.Request) bool {
	isMatch := r.Matcher.compareMethods(req, rec.Method) &&
		r.Matcher.compareURLs(req, rec.URL) &&
		r.Matcher.compareHeaders(req, rec) &&
		r.Matcher.compareBodies(req, rec.Body)

	return isMatch
}

func missingRequestError(req *http.Request) string {
	reqUrl := req.URL.String()
	return fmt.Sprintf("\nNo matching recorded request found.\nRequest: [%s] %s\n", req.Method, reqUrl)
}

func envNotExistsError(varName string) string {
	return "Required environment variable not set: " + varName
}

// applyVariableOptions applies the VariableType transform to the value
// If variableType is not provided or Default, return result
// If variableType is Secret_String, return SanitizedValue
// If variableType isSecret_Base64String return SanitizedBase64Value
func applyVariableOptions(val *string, variableType VariableType) *string {
	var ret string

	switch variableType {
	case Secret_String:
		ret = SanitizedValue
		return &ret
	case Secret_Base64String:
		ret = SanitizedBase64Value
		return &ret
	default:
		return val
	}
}

// initRandomSource initializes the Source to be used for random value creation in this Recording
func (r *Recording) initRandomSource() {
	// if we already have a Source generated, return immediately
	if r.src != nil {
		return
	}

	var seed int64
	var err error

	// check to see if we already have a random seed stored, use that if so
	seedString, ok := r.previousSessionVariables[randomSeedVariableName]
	if ok {
		seed, err = strconv.ParseInt(*seedString, 10, 64)
	}

	// We did not have a random seed already stored; create a new one
	if !ok || err != nil || r.Mode == Live {
		seed = time.Now().Unix()
		val := strconv.FormatInt(seed, 10)
		r.variables[randomSeedVariableName] = &val
	}

	// create a Source with the seed
	r.src = rand.NewSource(seed)
}

// initNow initializes the Source to be used for random value creation in this Recording
func (r *Recording) initNow() {
	// if we already have a now generated, return immediately
	if r.now != nil {
		return
	}

	var err error
	var nowStr *string
	var newNow time.Time

	// check to see if we already have a random seed stored, use that if so
	nowStr, ok := r.previousSessionVariables[nowVariableName]
	if ok {
		newNow, err = time.Parse(time.RFC3339Nano, *nowStr)
	}

	// We did not have a random seed already stored; create a new one
	if !ok || err != nil || r.Mode == Live {
		newNow = time.Now()
		nowStr = new(string)
		*nowStr = newNow.Format(time.RFC3339Nano)
		r.variables[nowVariableName] = nowStr
	}

	// save the now value.
	r.now = &newNow
}

// getFilePaths returns (recordingFilePath, variablesFilePath)
func getFilePaths(name string) (string, string) {
	recPath := "recordings/" + name
	varPath := fmt.Sprintf("%s-variables.yaml", recPath)
	return recPath, varPath
}

// createVariablesFileIfNotExists calls os.Create on the VariablesFile and creates it if it or the path does not exist
// Callers must call Close on the result
func (r *Recording) createVariablesFileIfNotExists() (*os.File, error) {
	f, err := os.Create(r.VariablesFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Create directory for the variables if missing
		variablesDir := filepath.Dir(r.VariablesFile)
		if _, err := os.Stat(variablesDir); os.IsNotExist(err) {
			if err = os.MkdirAll(variablesDir, 0755); err != nil {
				return nil, err
			}
		}

		f, err = os.Create(r.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (r *Recording) unmarshalVariablesFile(out interface{}) error {
	data, err := ioutil.ReadFile(r.VariablesFile)
	if err != nil {
		// If the file or dir do not exist, this is not an error to report
		if os.IsNotExist(err) {
			r.c.Log(fmt.Sprintf("Did not find recording for test '%s'", r.RecordingFile))
			return nil
		} else {
			return err
		}
	} else {
		err = yaml.Unmarshal(data, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Recording) initVariables() error {
	return r.unmarshalVariablesFile(r.previousSessionVariables)
}

var modeMap = map[RecordMode]recorder.Mode{
	Record:   recorder.ModeRecording,
	Live:     recorder.ModeDisabled,
	Playback: recorder.ModeReplaying,
}
