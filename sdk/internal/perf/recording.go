// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
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
			panic(err)
		}
	}
	cert, err := os.ReadFile(localFile)
	if err != nil {
		log.Printf("could not read file set in PROXY_CERT variable at %s.\n", localFile)
	}

	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Println("no certs appended, using system certs only")
	}

	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
			RootCAs:            rootCAs,
		},
	}
	defaultHTTPClient = &http.Client{
		Transport: defaultTransport,
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
	playbackMode      = "playback"
	liveMode          = "live"
	idHeader          = "x-recording-id"
	modeHeader        = "x-recording-mode"
	upstreamURIHeader = "x-recording-upstream-base-uri"
)

var (
	defaultHTTPClient    *http.Client
	rootCAs              *x509.CertPool
	proxyTransportsSuite = map[string]*RecordingHTTPClient{}
)

type TransportOptions struct {
	TestName string
	proxyURL string
}

type RecordingHTTPClient struct {
	defaultClient *http.Client
	options       TransportOptions
	mode          string
	recID         string
}

// NewRecordingHTTPClient returns a type that implements `azcore.Transporter`. This will automatically route tests on the `Do` call.
func NewProxyTransport(options *TransportOptions) *RecordingHTTPClient {
	if options == nil {
		options = &TransportOptions{}
	}
	if debug {
		log.Println("Creating a new proxy transport: ", *options)
	}

	ret := &RecordingHTTPClient{
		defaultClient: defaultHTTPClient,
		options:       *options,
		mode:          "live",
	}

	proxyTransportsSuite[options.TestName] = ret

	return ret
}

func (c RecordingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.mode != liveMode {
		var err error
		if req, err = c.replaceAuthority(req); err != nil {
			return nil, err
		}
	}
	return c.defaultClient.Do(req)
}

// Change the recording mode
func (c *RecordingHTTPClient) SetMode(mode string) {
	c.mode = mode
}

func (c *RecordingHTTPClient) replaceAuthority(rawReq *http.Request) (*http.Request, error) {
	parsedProxyURL, err := url.Parse(c.options.proxyURL)
	if err != nil {
		return nil, fmt.Errorf("there was an error parsing url '%s': %s", c.options.proxyURL, err.Error())
	}
	originalURLHost := rawReq.URL.Host
	originalURLScheme := rawReq.URL.Scheme

	// don't modify the original request
	cp := *rawReq
	cpURL := *cp.URL
	cp.URL = &cpURL
	cp.Header = rawReq.Header.Clone()

	cp.URL.Scheme = parsedProxyURL.Scheme
	cp.URL.Host = parsedProxyURL.Host
	cp.Host = parsedProxyURL.Host

	cp.Header.Set(upstreamURIHeader, fmt.Sprintf("%v://%v", originalURLScheme, originalURLHost))
	cp.Header.Set(modeHeader, c.mode)
	cp.Header.Set(idHeader, c.recID)
	cp.Header.Set("x-recording-remove", "false")
	return &cp, nil
}

// start tells the test proxy to begin accepting requests for a given test
func (c *RecordingHTTPClient) start() error {
	url := fmt.Sprintf("%s/%s/start", c.options.proxyURL, c.mode)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("there was an error creating a START request: %s", err.Error())
	}

	if c.mode == playbackMode {
		req.Header.Set(idHeader, c.recID)
	}

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("there was an error communicating with the test proxy: %s", err.Error())
	}

	recID := resp.Header.Get(idHeader)
	if recID == "" {
		b, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("there was an error reading the body: %s", err.Error())
		}
		return fmt.Errorf("recording ID was not returned by the response. Response body: %s", b)
	}
	c.recID = recID

	return nil
}

// stop tells the test proxy to stop accepting requests for a given test
func (c *RecordingHTTPClient) stop() error {
	url := fmt.Sprintf("%s/%s/stop", c.options.proxyURL, c.mode)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("there was an error creating a STOP request: %s", err.Error())
	}

	if c.recID == "" {
		return errors.New("recording ID was never set. Did you call Start?")
	}

	req.Header.Set("x-recording-id", c.recID) //recTest)
	resp, err := defaultHTTPClient.Do(req)
	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err == nil {
			return fmt.Errorf("proxy did not stop the recording properly: %s", string(b))
		}
		return fmt.Errorf("proxy did not stop the recording properly: %s", err.Error())
	}
	if err != nil {
		return fmt.Errorf("there was an error communicating with the test proxy: %s", err.Error())
	}
	return nil
}
