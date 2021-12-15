// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package recording

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
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
)

var pathToRecordings = "eng/tools/azgoperf"

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

const (
	RecordingMode     = "record"
	PlaybackMode      = "playback"
	LiveMode          = "live"
	IDHeader          = "x-recording-id"
	ModeHeader        = "x-recording-mode"
	UpstreamURIHeader = "x-recording-upstream-base-uri"
)

var rootCAs *x509.CertPool
var recordMode string
var perfTestSuite = map[string]string{}

type TransportOptions struct {
	UseHTTPS bool
	TestName string
}

func (r TransportOptions) host() string {
	if r.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

func (r TransportOptions) scheme() string {
	if r.UseHTTPS {
		return "https"
	}
	return "http"
}

func GetRecordMode() string {
	return recordMode
}

type RecordingHTTPClient struct {
	defaultClient *http.Client
	options       TransportOptions
}

// NewRecordingHTTPClient returns a type that implements `azcore.Transporter`. This will automatically route tests on the `Do` call.
func NewProxyTransport(options *TransportOptions) (*RecordingHTTPClient, error) {
	if options == nil {
		options = &TransportOptions{UseHTTPS: true}
	}
	c, err := GetHTTPClient()
	if err != nil {
		return nil, err
	}

	return &RecordingHTTPClient{
		defaultClient: c,
		options:       *options,
	}, nil
}

func GetHTTPClient() (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.RootCAs = rootCAs
	transport.TLSClientConfig.MinVersion = tls.VersionTLS12
	transport.TLSClientConfig.InsecureSkipVerify = true

	defaultHttpClient := &http.Client{
		Transport: transport,
	}
	return defaultHttpClient, nil
}

func (c RecordingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.options.ReplaceAuthority(req)
	return c.defaultClient.Do(req)
}
func (r TransportOptions) ReplaceAuthority(rawReq *http.Request) {
	originalURLHost := rawReq.URL.Host
	rawReq.URL.Scheme = r.scheme()
	rawReq.URL.Host = r.host()
	rawReq.Host = r.host()

	rawReq.Header.Set(UpstreamURIHeader, fmt.Sprintf("%v://%v", r.scheme(), originalURLHost))
	rawReq.Header.Set(ModeHeader, GetRecordMode())
	rawReq.Header.Set(IDHeader, GetRecordingId(r.TestName))
}

func GetRecordingId(s string) string {
	if v, ok := perfTestSuite[s]; ok {
		return v
	}
	return ""
}

type RecordingOptions struct{}

func (r RecordingOptions) baseURL() string {
	return "https://localhost:5001"
}

func getTestId(pathToRecordings string, t string) string {
	return path.Join(pathToRecordings, "recordings", t+".json")
}

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Start tells the test proxy to begin accepting requests for a given test
func Start(t string, options *RecordingOptions) error {
	if options == nil {
		options = &RecordingOptions{}
	}

	testId := getTestId(pathToRecordings, t)

	url := fmt.Sprintf("https://localhost:5001/%s/start", recordMode)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-recording-file", testId)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	recId := resp.Header.Get(IDHeader)
	if recId == "" {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("recording ID was not returned by the response. Response body: %s", b)
	}

	perfTestSuite[t] = recId

	return nil
}

// Stop tells the test proxy to stop accepting requests for a given test
func Stop(t string, options *RecordingOptions) error {
	if options == nil {
		options = &RecordingOptions{}
	}

	url := fmt.Sprintf("%v/%v/stop", options.baseURL(), recordMode)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	var recTest string
	var ok bool
	if recTest, ok = perfTestSuite[t]; !ok {
		return errors.New("recording ID was never set. Did you call StartRecording?")
	}
	req.Header.Set("x-recording-id", recTest)
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err == nil {
			return fmt.Errorf("proxy did not stop the recording properly: %s", string(b))
		}
		return fmt.Errorf("proxy did not stop the recording properly: %s", err.Error())
	}
	return err
}
