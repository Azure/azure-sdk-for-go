//go:build go1.16
// +build go1.16

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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

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
	cert, err := ioutil.ReadFile(localFile)
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

var recordMode string
var rootCAs *x509.CertPool

const (
	RecordingMode     = "record"
	PlaybackMode      = "playback"
	LiveMode          = "live"
	IDHeader          = "x-recording-id"
	ModeHeader        = "x-recording-mode"
	UpstreamURIHeader = "x-recording-upstream-base-uri"
)

type recordedTest struct {
	recordingId string
	liveOnly    bool
	variables   map[string]interface{}
}

var testSuite = map[string]recordedTest{}

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

func (r RecordingOptions) ReplaceAuthority(t *testing.T, rawReq *http.Request) {
	if GetRecordMode() != LiveMode && !IsLiveOnly(t) {
		originalURLHost := rawReq.URL.Host
		rawReq.URL.Scheme = r.scheme()
		rawReq.URL.Host = r.host()
		rawReq.Host = r.host()

		rawReq.Header.Set(UpstreamURIHeader, fmt.Sprintf("%v://%v", r.scheme(), originalURLHost))
		rawReq.Header.Set(ModeHeader, GetRecordMode())
		rawReq.Header.Set(IDHeader, GetRecordingId(t))
	}
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
	return path.Join(pathToRecordings, "recordings", t.Name()+".json")
}

// Start tells the test proxy to begin accepting requests for a given test
func Start(t *testing.T, pathToRecordings string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	if recordMode == LiveMode {
		return nil
	}

	if testStruct, ok := testSuite[t.Name()]; ok {
		if testStruct.liveOnly {
			// test should only be run live, don't want to generate recording
			return nil
		}
	}

	testId := getTestId(pathToRecordings, t)

	url := fmt.Sprintf("%s/%s/start", options.baseURL(), recordMode)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	marshalled, err := json.Marshal(map[string]string{"x-recording-file": testId})
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	recId := resp.Header.Get(IDHeader)
	log.Printf("test name: %s\t recording ID: %s\n", t.Name(), recId)
	if recId == "" {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("recording ID was not returned by the response. Response body: %s", b)
	}

	// Unmarshal any variables returned by the proxy
	var m map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
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

	if val, ok := testSuite[t.Name()]; ok {
		val.recordingId = recId
		val.variables = m
		testSuite[t.Name()] = val
	} else {
		testSuite[t.Name()] = recordedTest{
			recordingId: recId,
			liveOnly:    false,
			variables:   m,
		}
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

	if testStruct, ok := testSuite[t.Name()]; ok {
		if testStruct.liveOnly {
			// test should only be run live, don't want to generate recording
			return nil
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
		req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
		req.ContentLength = int64(len(marshalled))
	}

	var recTest recordedTest
	var ok bool
	if recTest, ok = testSuite[t.Name()]; !ok {
		return errors.New("recording ID was never set. Did you call StartRecording?")
	}
	req.Header.Set(IDHeader, recTest.recordingId)
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err == nil {
			return fmt.Errorf("proxy did not stop the recording properly: %s", string(b))
		}
		return fmt.Errorf("proxy did not stop the recording properly: %s", err.Error())
	}
	_ = resp.Body.Close()
	return err
}

// This looks up an environment variable and if it is not found, returns the recordedValue
func GetEnvVariable(varName string, recordedValue string) string {
	val, ok := os.LookupEnv(varName)
	if !ok || GetRecordMode() == PlaybackMode {
		return recordedValue
	}
	return val
}

func LiveOnly(t *testing.T) {
	if val, ok := testSuite[t.Name()]; ok {
		val.liveOnly = true
		testSuite[t.Name()] = val
	} else {
		testSuite[t.Name()] = recordedTest{liveOnly: true}
	}
	if GetRecordMode() == PlaybackMode {
		t.Skip("Live Test Only")
	}
}

// Function for sleeping during a test for `duration` seconds. This method will only execute when
// AZURE_RECORD_MODE = "record", if a test is running in playback this will be a noop.
func Sleep(duration time.Duration) {
	if GetRecordMode() != PlaybackMode {
		time.Sleep(duration)
	}
}

func GetRecordingId(t *testing.T) string {
	if val, ok := testSuite[t.Name()]; ok {
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
	c.options.ReplaceAuthority(c.t, req)
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
	if s, ok := testSuite[t.Name()]; ok {
		return s.liveOnly
	}
	return false
}

// GetVariables returns access to the variables stored by the test proxy for a specific test
func GetVariables(t *testing.T) map[string]interface{} {
	if s, ok := testSuite[t.Name()]; ok {
		return s.variables
	}
	return nil
}
