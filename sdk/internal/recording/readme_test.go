//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// Snippet: PolicyDefinition
// This should be a 'testdata' directory in your module. `testdata` is ignored by the go tool, making it perfect for ancillary data
var pathToPackage = "sdk/packageToTest/testdata"

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func (r recordingPolicy) Host() string {
	if r.options.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

func (r recordingPolicy) Scheme() string {
	if r.options.UseHTTPS {
		return "https"
	}
	return "http"
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

func (r *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" && !recording.IsLiveOnly(r.t) {
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = r.Scheme()
		req.Raw().URL.Host = r.Host()
		req.Raw().Host = r.Host()

		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%v://%v", r.Scheme(), originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IDHeader, recording.GetRecordingId(r.t))
	}
	return req.Next()
}

// EndSnippet

// Snippet: TestFunction
func TestSomething(t *testing.T) {
	// SnippetIgnore
	t.Skip()
	// EndSnippetIgnore
	p := NewRecordingPolicy(t, nil)
	httpClient, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &azcore.ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       httpClient,
	}

	cred := authenticate(t)

	client, err := NewClient("https://mystorageaccount.table.core.windows.net", cred, options)
	require.NoError(t, err)
	// Continue test
	// SnippetIgnore
	_ = client
	// EndSnippetIgnore
}

// EndSnippet

// Snippet: StartAndStopRecording
func TestSomethingWithStartAndStop(t *testing.T) {
	// SnippetIgnore
	t.Skip()
	// EndSnippetIgnore
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	defer recording.Stop(t, nil)

	// Continue test
}

// EndSnippet

// Snippet: UseSanitizer
func TestSomethingWithSanitizer(t *testing.T) {
	// SnippetIgnore
	t.Skip()
	// EndSnippetIgnore
	err := recording.AddURISanitizer("fakeaccountname", "my-real-account-name", nil)
	require.NoError(t, err)

	// To remove the sanitizer after this test use the following:
	defer recording.ResetSanitizers(nil)

	err = recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	defer recording.Stop(t, nil)

	// Continue test
}

// EndSnippet

// Snippet: ReadEnvVar
func TestSomethingByReadingEnvVar(t *testing.T) {
	accountName := recording.GetEnvVariable("TABLES_PRIMARY_ACCOUNT_NAME", "fakeaccountname")
	if recording.GetRecordMode() == string(recording.Record) {
		err := recording.AddURISanitizer("fakeaccountname", accountName, nil)
		require.NoError(t, err)
	}

	// Continue test
}

// EndSnippet

func authenticate(t *testing.T) azcore.TokenCredential {
	return nil
}

type SDKClient struct{}

func NewClient(endpoint string, cred azcore.TokenCredential, opt *azcore.ClientOptions) (SDKClient, error) {
	// this function is only here to make sure compile can pass
	return SDKClient{}, nil
}
