//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
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
	recordingAssetConfigName    = "assets.json"
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
	r.initRandomSource()
	return generateAlphaNumericID(prefix, length, lowercaseOnly, r.src)
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
	data, err := os.ReadFile(r.VariablesFile)
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

func init() {
	recordMode = os.Getenv("AZURE_RECORD_MODE")
	if recordMode == "" {
		log.Printf("AZURE_RECORD_MODE was not set, defaulting to playback")
		recordMode = PlaybackMode
	}
	if !(recordMode == RecordingMode || recordMode == PlaybackMode || recordMode == LiveMode) {
		log.Panicf("AZURE_RECORD_MODE was not understood, options are %s, %s, or %s Received: %v.\n", RecordingMode, PlaybackMode, LiveMode, recordMode)
	}

	localFile, err := findProxyCertLocation()
	if err != nil {
		log.Println("Could not find the PROXY_CERT environment variable and was unable to locate the path in eng/common")
	}

	var certPool *x509.CertPool
	if runtime.GOOS == "windows" {
		certPool = x509.NewCertPool()
	} else {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			log.Println("could not create a system cert pool")
			log.Panicf(err.Error())
		}
	}
	cert, err := os.ReadFile(localFile)
	if err != nil {
		log.Printf("could not read file set in PROXY_CERT variable at %s.\n", localFile)
	}

	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Println("no certs appended, using system certs only")
	}

	// Set a Default matcher that ignores :path, :scheme, :authority, and :method headers
	err = SetDefaultMatcher(
		nil,
		&SetDefaultMatcherOptions{ExcludedHeaders: []string{
			":authority",
			":method",
			":path",
			":scheme",
		}},
	)
	if err != nil {
		log.Println("could not set the default matcher")
	} else {
		log.Println("default matcher was set ")
	}
}

var (
	recordMode string
	rootCAs    *x509.CertPool
)

const (
	RecordingMode           = "record"
	PlaybackMode            = "playback"
	LiveMode                = "live"
	IDHeader                = "x-recording-id"
	ModeHeader              = "x-recording-mode"
	UpstreamURIHeader       = "x-recording-upstream-base-uri"
	recordingRandSeedVarKey = "randSeed"
)

type recordedTest struct {
	recordingId      string
	liveOnly         bool
	variables        map[string]interface{}
	recordingSeed    int64
	recordingRandSrc rand.Source
}

// testMap maps test names to metadata
type testMap struct {
	m *sync.Map
}

// Load returns the named test's metadata, if it has been stored
func (t *testMap) Load(name string) (recordedTest, bool) {
	var rt recordedTest
	v, ok := t.m.Load(name)
	if ok {
		rt = v.(recordedTest)
	}
	return rt, ok
}

// Store sets metadata for the named test
func (t *testMap) Store(name string, data recordedTest) {
	t.m.Store(name, data)
}

// Remove delete metadata for the named test
func (t *testMap) Remove(name string) {
	t.m.Delete(name)
}

var testSuite = testMap{&sync.Map{}}

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type RecordingOptions struct {
	UseHTTPS        bool
	GroupForReplace string
	Variables       map[string]interface{}
	TestInstance    *testing.T
}

func defaultOptions() *RecordingOptions {
	return &RecordingOptions{
		UseHTTPS: true,
	}
}

func (r RecordingOptions) ReplaceAuthority(t *testing.T, rawReq *http.Request) *http.Request {
	if GetRecordMode() != LiveMode && !IsLiveOnly(t) {
		originalURLHost := rawReq.URL.Host

		// don't modify the original request
		cp := *rawReq
		cpURL := *cp.URL
		cp.URL = &cpURL
		cp.Header = rawReq.Header.Clone()

		cp.URL.Scheme = r.scheme()
		cp.URL.Host = r.host()
		cp.Host = r.host()

		cp.Header.Set(UpstreamURIHeader, fmt.Sprintf("%v://%v", r.scheme(), originalURLHost))
		cp.Header.Set(ModeHeader, GetRecordMode())
		cp.Header.Set(IDHeader, GetRecordingId(t))
		rawReq = &cp
	}
	return rawReq
}

func (r RecordingOptions) host() string {
	if r.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

func (r RecordingOptions) scheme() string {
	if r.UseHTTPS {
		return "https"
	}
	return "http"
}

func (r RecordingOptions) baseURL() string {
	return fmt.Sprintf("%s://%s", r.scheme(), r.host())
}

func getTestId(pathToRecordings string, t *testing.T) string {
	return filepath.Join(pathToRecordings, "recordings", t.Name()+".json")
}

func getGitRoot(fromPath string) (string, error) {
	absPath, err := filepath.Abs(fromPath)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = absPath

	root, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Unable to find git root for path '%s'", absPath)
	}

	// Wrap with Abs() to get os-specific path separators to support sub-path matching
	return filepath.Abs(strings.TrimSpace(string(root)))
}

// Traverse up from a recording path until an asset config file is found.
// Stop searching when the root of the git repository is reached.
func findAssetsConfigFile(fromPath string, untilPath string) (string, error) {
	absPath, err := filepath.Abs(fromPath)
	if err != nil {
		return "", err
	}
	assetConfigPath := filepath.Join(absPath, recordingAssetConfigName)

	if _, err := os.Stat(assetConfigPath); err == nil {
		return assetConfigPath, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	if absPath == untilPath {
		return "", nil
	}

	parentDir := filepath.Dir(absPath)
	// This shouldn't be hit due to checks in getGitRoot, but it can't hurt to be defensive
	if parentDir == absPath || parentDir == "." {
		return "", nil
	}

	return findAssetsConfigFile(parentDir, untilPath)
}

// Returns absolute and relative paths to an asset configuration file, or an error.
func getAssetsConfigLocation(pathToRecordings string) (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	gitRoot, err := getGitRoot(cwd)
	if err != nil {
		return "", "", err
	}
	abs, err := findAssetsConfigFile(filepath.Join(gitRoot, pathToRecordings), gitRoot)
	if err != nil {
		return "", "", err
	}

	// Pass a path relative to the git root to test proxy so that paths
	// can be resolved when the repo root is mounted as a volume in a container
	rel := strings.Replace(abs, gitRoot, "", 1)
	rel = strings.TrimLeft(rel, string(os.PathSeparator))
	return abs, rel, nil
}

func requestStart(url string, testId string, assetConfigLocation string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	reqBody := map[string]string{"x-recording-file": testId}
	if assetConfigLocation != "" {
		reqBody["x-recording-assets-file"] = assetConfigLocation
	}
	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))

	return client.Do(req)
}

// Start tells the test proxy to begin accepting requests for a given test
func Start(t *testing.T, pathToRecordings string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	if recordMode == LiveMode {
		return nil
	}

	if testStruct, ok := testSuite.Load(t.Name()); ok {
		if testStruct.liveOnly {
			// test should only be run live, don't want to generate recording
			return nil
		}
	}

	testId := getTestId(pathToRecordings, t)

	absAssetLocation, relAssetLocation, err := getAssetsConfigLocation(pathToRecordings)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/start", options.baseURL(), recordMode)

	var resp *http.Response
	if absAssetLocation == "" {
		resp, err = requestStart(url, testId, "")
		if err != nil {
			return err
		}
	} else if resp, err = requestStart(url, testId, absAssetLocation); err != nil {
		return err
	} else if resp.StatusCode >= 400 {
		if resp, err = requestStart(url, testId, relAssetLocation); err != nil {
			return err
		}
	}

	recId := resp.Header.Get(IDHeader)
	if recId == "" {
		b, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("Recording ID was not returned by the response. Response body: %s", b)
	}

	// Unmarshal any variables returned by the proxy
	var m map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			return err
		}
	}

	if val, ok := testSuite.Load(t.Name()); ok {
		val.recordingId = recId
		val.variables = m
		testSuite.Store(t.Name(), val)
	} else {
		testSuite.Store(t.Name(), recordedTest{
			recordingId: recId,
			liveOnly:    false,
			variables:   m,
		})
	}
	return nil
}

// Stop tells the test proxy to stop accepting requests for a given test
func Stop(t *testing.T, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	if recordMode == LiveMode {
		return nil
	}

	if testStruct, ok := testSuite.Load(t.Name()); ok {
		if testStruct.liveOnly {
			// test should only be run live, don't want to generate recording
			return nil
		}
		if testStruct.recordingSeed != 0 {
			if options.Variables == nil {
				options.Variables = map[string]interface{}{}
			}
			options.Variables[recordingRandSeedVarKey] = strconv.FormatInt(testStruct.recordingSeed, 10)
		}
	}

	url := fmt.Sprintf("%v/%v/stop", options.baseURL(), recordMode)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if len(options.Variables) > 0 {
		req.Header.Set("Content-Type", "application/json")
		marshalled, err := json.Marshal(options.Variables)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(marshalled))
		req.ContentLength = int64(len(marshalled))
	}

	var recTest recordedTest
	var ok bool
	if recTest, ok = testSuite.Load(t.Name()); !ok {
		return errors.New("Recording ID was never set. Did you call StartRecording?")
	}
	req.Header.Set(IDHeader, recTest.recordingId)
	testSuite.Remove(t.Name())
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err == nil {
			return fmt.Errorf("proxy did not stop the recording properly: %s", string(b))
		}
		return fmt.Errorf("proxy did not stop the recording properly: %s", err.Error())
	}
	_ = resp.Body.Close()
	return err
}

func getRandomSource(t *testing.T) rand.Source {
	if testStruct, ok := testSuite.Load(t.Name()); ok {
		if testStruct.recordingRandSrc != nil {
			return testStruct.recordingRandSrc
		}
	}

	var seed int64
	var err error

	variables := GetVariables(t)
	seedString, ok := variables[recordingRandSeedVarKey]
	if ok {
		seed, err = strconv.ParseInt(seedString.(string), 10, 64)
	}

	// We did not have a random seed already stored; create a new one
	if !ok || err != nil || GetRecordMode() == "live" {
		seed = time.Now().Unix()
	}

	source := rand.NewSource(seed)
	if testStruct, ok := testSuite.Load(t.Name()); ok {
		testStruct.recordingSeed = seed
		testStruct.recordingRandSrc = source
		testSuite.Store(t.Name(), testStruct)
	}

	return source
}

// GenerateAlphaNumericID will generate a recorded random alpha numeric id.
// When live mode or the recording has a randomSeed already set, the value will be generated from that seed, else a new random seed will be used.
func GenerateAlphaNumericID(t *testing.T, prefix string, length int, lowercaseOnly bool) (string, error) {
	return generateAlphaNumericID(prefix, length, lowercaseOnly, getRandomSource(t))
}

func generateAlphaNumericID(prefix string, length int, lowercaseOnly bool, randomSource rand.Source) (string, error) {
	if length <= len(prefix) {
		return "", errors.New("length must be greater than prefix")
	}

	sb := strings.Builder{}
	sb.Grow(length)
	sb.WriteString(prefix)
	i := length - len(prefix) - 1
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for cache, remain := randomSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSource.Int63(), letterIdxMax
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

// GetEnvVariable looks up an environment variable and if it is not found, returns the recordedValue
func GetEnvVariable(varName string, recordedValue string) string {
	val, ok := os.LookupEnv(varName)
	if !ok || GetRecordMode() == PlaybackMode {
		return recordedValue
	}
	return val
}

func LiveOnly(t *testing.T) {
	if val, ok := testSuite.Load(t.Name()); ok {
		val.liveOnly = true
		testSuite.Store(t.Name(), val)
	} else {
		testSuite.Store(t.Name(), recordedTest{liveOnly: true})
	}
	if GetRecordMode() == PlaybackMode {
		t.Skip("Live Test Only")
	}
}

// Sleep during a test for `duration` seconds. This method will only execute when
// AZURE_RECORD_MODE = "record", if a test is running in playback this will be a noop.
func Sleep(duration time.Duration) {
	if GetRecordMode() != PlaybackMode {
		time.Sleep(duration)
	}
}

func GetRecordingId(t *testing.T) string {
	if val, ok := testSuite.Load(t.Name()); ok {
		return val.recordingId
	} else {
		return ""
	}
}

func GetRecordMode() string {
	return recordMode
}

func findProxyCertLocation() (string, error) {
	fileLocation, ok := os.LookupEnv("PROXY_CERT")
	if ok {
		return fileLocation, nil
	}

	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Print("Could not find PROXY_CERT environment variable or toplevel of git repository, please set PROXY_CERT to location of certificate found in eng/common/testproxy/dotnet-devcert.crt")
		return "", err
	}
	topLevel := bytes.NewBuffer(out).String()
	return filepath.Join(topLevel, "eng", "common", "testproxy", "dotnet-devcert.crt"), nil
}

type RecordingHTTPClient struct {
	defaultClient *http.Client
	options       RecordingOptions
	t             *testing.T
}

func (c RecordingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req = c.options.ReplaceAuthority(c.t, req)
	return c.defaultClient.Do(req)
}

// NewRecordingHTTPClient returns a type that implements `azcore.Transporter`. This will automatically route tests on the `Do` call.
func NewRecordingHTTPClient(t *testing.T, options *RecordingOptions) (*RecordingHTTPClient, error) {
	if options == nil {
		options = &RecordingOptions{UseHTTPS: true}
	}
	c, err := GetHTTPClient(t)
	if err != nil {
		return nil, err
	}

	return &RecordingHTTPClient{
		defaultClient: c,
		options:       *options,
		t:             t,
	}, nil
}

func GetHTTPClient(t *testing.T) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.RootCAs = rootCAs
	transport.TLSClientConfig.MinVersion = tls.VersionTLS12
	transport.TLSClientConfig.InsecureSkipVerify = true

	defaultHttpClient := &http.Client{
		Transport: transport,
	}
	return defaultHttpClient, nil
}

func IsLiveOnly(t *testing.T) bool {
	if s, ok := testSuite.Load(t.Name()); ok {
		return s.liveOnly
	}
	return false
}

// GetVariables returns access to the variables stored by the test proxy for a specific test
func GetVariables(t *testing.T) map[string]interface{} {
	if s, ok := testSuite.Load(t.Name()); ok {
		return s.variables
	}
	return nil
}
