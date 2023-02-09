//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const recordingRandomSeedVariableName = "recordingRandomSeed"

var (
	recordingSeed         int64
	recordingRandomSource rand.Source
)

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

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

// Host of the test proxy.
func (r *recordingPolicy) Host() string {
	if r.options.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

// Scheme of the test proxy.
func (r *recordingPolicy) Scheme() string {
	if r.options.UseHTTPS {
		return "https"
	}
	return "http"
}

// NewRecordingPolicy will create a recording policy which can be used in pipeline.
// The policy will change the destination of the request to the proxy server and add required header for the recording test.
func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

// Do with recording mode.
// When handling live request, the policy will do nothing.
// Otherwise, the policy will replace the URL of the request with the test proxy endpoint.
// After request, the policy will change back to the original URL for the request to prevent wrong polling URL for LRO.
func (r *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" && !recording.IsLiveOnly(r.t) {
		oriSchema := req.Raw().URL.Scheme
		oriHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = r.Scheme()
		req.Raw().URL.Host = r.Host()
		req.Raw().Host = r.Host()

		// replace request target to use test proxy
		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%v://%v", oriSchema, oriHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IDHeader, recording.GetRecordingId(r.t))

		resp, err = req.Next()
		// for any lro operation, need to change back to the original target to prevent
		if resp != nil {
			resp.Request.URL.Scheme = oriSchema
			resp.Request.URL.Host = oriHost
		}
		return resp, err
	} else {
		return req.Next()
	}
}

// StartRecording starts the recording with the path to store recording file.
// It will return a delegate function to stop recording.
func StartRecording(t *testing.T, pathToPackage string) func() {
	// sanitizer for any uuid string, e.g., subscriptionID
	err := recording.AddGeneralRegexSanitizer("00000000-0000-0000-0000-000000000000", `[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`, nil)
	if err != nil {
		t.Fatalf("Failed to add uuid sanitizer: %v", err)
	}
	// consolidate resource group name for recording and playback
	err = recording.AddGeneralRegexSanitizer("go-sdk-test-rg", `go-sdk-test-\d+`, nil)
	if err != nil {
		t.Fatalf("Failed to add resource group name sanitizer: %v", err)
	}
	err = recording.Start(t, pathToPackage, nil)
	if err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}
	return func() { StopRecording(t) }
}

// StopRecording stops the recording.
func StopRecording(t *testing.T) {
	err := recording.Stop(t, &recording.RecordingOptions{Variables: map[string]interface{}{recordingRandomSeedVariableName: strconv.FormatInt(recordingSeed, 10)}})
	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}
}

func initRandomSource(t *testing.T) {

	if recordingRandomSource != nil {
		return
	}

	var seed int64
	var err error

	variables := recording.GetVariables(t)
	seedString, ok := variables[recordingRandomSeedVariableName]
	if ok {
		seed, err = strconv.ParseInt(seedString.(string), 10, 64)
	}

	// We did not have a random seed already stored; create a new one
	if !ok || err != nil || recording.GetRecordMode() == "live" {
		seed = time.Now().Unix()
		val := strconv.FormatInt(seed, 10)
		variables[recordingRandomSeedVariableName] = &val
	}

	// create a Source with the seed
	recordingRandomSource = rand.NewSource(seed)
	recordingSeed = seed
}

// GenerateAlphaNumericID will generate a random alpha numeric ID.
// When handling live request, the random seed is generated.
// Otherwise, the random seed is stable and will be stored in recording file.
// The length parameter is the random part length, not include the prefix part.
func GenerateAlphaNumericID(t *testing.T, prefix string, length int) string {

	var lowercaseOnly bool = false

	initRandomSource(t)
	sb := strings.Builder{}
	sb.Grow(length)
	sb.WriteString(prefix)
	i := length - 1
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for cache, remain := recordingRandomSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = recordingRandomSource.Int63(), letterIdxMax
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
	return str
}
