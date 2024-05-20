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
)

// Deprecated: the local recording API that uses this type is no longer supported. Call [Start] and [Stop]
// to make recordings via the test proxy instead.
type Recording struct {
	SessionName   string
	RecordingFile string
	VariablesFile string
	Mode          RecordMode
	Sanitizer     *Sanitizer
	Matcher       *RequestMatcher
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

// Deprecated: only deprecated methods use this type. Call [Start] and [Stop] to make recordings.
type VariableType string

const (
	// NoSanitization indicates that the recorded value should not be sanitized.
	NoSanitization VariableType = "default"
	// Secret_String indicates that the recorded value should be replaced with a sanitized value.
	Secret_String VariableType = "secret_string"
	// Secret_Base64String indicates that the recorded value should be replaced with a sanitized valid base-64 string value.
	Secret_Base64String VariableType = "secret_base64String"
)

var errUnsupportedAPI = errors.New("the vcr based test recording API isn't supported. Use the test proxy instead")

// NewRecording initializes a new Recording instance
func NewRecording(c TestContext, mode RecordMode) (*Recording, error) {
	return nil, errUnsupportedAPI
}

// GetEnvVar returns a recorded environment variable. If the variable is not found we return an error.
// variableType determines how the recorded variable will be saved.
func (r *Recording) GetEnvVar(name string, variableType VariableType) (string, error) {
	return "", errUnsupportedAPI
}

// GetOptionalEnvVar returns a recorded environment variable with a fallback default value.
// default Value configures the fallback value to be returned if the environment variable is not set.
// variableType determines how the recorded variable will be saved.
func (r *Recording) GetOptionalEnvVar(name string, defaultValue string, variableType VariableType) string {
	panic(errUnsupportedAPI)
}

// Do satisfies the azcore.Transport interface so that Recording can be used as the transport for recorded requests
func (r *Recording) Do(req *http.Request) (*http.Response, error) {
	return nil, errUnsupportedAPI
}

// Stop stops the recording and saves them, including any captured variables, to disk
func (r *Recording) Stop() error {
	return errUnsupportedAPI
}

func (r *Recording) Now() time.Time {
	panic(errUnsupportedAPI)
}

func (r *Recording) UUID() uuid.UUID {
	panic(errUnsupportedAPI)
}

// GenerateAlphaNumericID will generate a recorded random alpha numeric id
// if the recording has a randomSeed already set, the value will be generated from that seed, else a new random seed will be used
func (r *Recording) GenerateAlphaNumericID(prefix string, length int, lowercaseOnly bool) (string, error) {
	return "", errUnsupportedAPI
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
	ProxyPort       int
	GroupForReplace string
	Variables       map[string]interface{}
	TestInstance    *testing.T

	// insecure allows this package's tests to configure the proxy to skip upstream TLS
	// verification so they can use a mock upstream server having a self-signed cert
	insecure bool
}

func defaultOptions() *RecordingOptions {
	return &RecordingOptions{
		UseHTTPS:  true,
		ProxyPort: os.Getpid()%10000 + 20000,
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
	if r.ProxyPort != 0 {
		return fmt.Sprintf("localhost:%d", r.ProxyPort)
	}

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
		return "", fmt.Errorf("unable to find git root for path '%s'", absPath)
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

func Start(t *testing.T, pathToRecordings string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}

	if testStruct, ok := testSuite.Load(t.Name()); ok {
		if testStruct.liveOnly {
			// test should only be run live, don't want to generate recording
			return nil
		}
	}
	if options == nil {
		options = defaultOptions()
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
	if err != nil {
		return err
	}
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
	origScheme := req.URL.Scheme
	origHost := req.URL.Host
	req = c.options.ReplaceAuthority(c.t, req)
	resp, err := c.defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// if the request succeeds, restore the scheme/host with their original values.
	// this is imporant for things like LROs that might use the originating URL to
	// poll for status and/or fetch the final result.
	resp.Request.URL.Scheme = origScheme
	resp.Request.URL.Host = origHost
	return resp, nil
}

// NewRecordingHTTPClient returns a type that implements `azcore.Transporter`. This will automatically route tests on the `Do` call.
func NewRecordingHTTPClient(t *testing.T, options *RecordingOptions) (*RecordingHTTPClient, error) {
	if options == nil {
		options = defaultOptions()
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
