// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

// configuration for live tests
var liveSP = struct {
	tenantID string
	clientID string
	secret   string
	pemPath  string
	pfxPath  string
	sniPath  string
}{
	tenantID: os.Getenv("IDENTITY_SP_TENANT_ID"),
	clientID: os.Getenv("IDENTITY_SP_CLIENT_ID"),
	secret:   os.Getenv("IDENTITY_SP_CLIENT_SECRET"),
	pemPath:  os.Getenv("IDENTITY_SP_CERT_PEM"),
	pfxPath:  os.Getenv("IDENTITY_SP_CERT_PFX"),
	sniPath:  os.Getenv("IDENTITY_SP_CERT_SNI"),
}

var liveUser = struct {
	tenantID string
	username string
	password string
}{
	username: os.Getenv("AZURE_IDENTITY_TEST_USERNAME"),
	password: os.Getenv("AZURE_IDENTITY_TEST_PASSWORD"),
	tenantID: os.Getenv("AZURE_IDENTITY_TEST_TENANTID"),
}

const (
	fakeClientID   = "fake-client-id"
	fakeResourceID = "/fake/resource/ID"
	fakeTenantID   = "fake-tenant"
)

var liveTestScope = "https://management.core.windows.net//.default"

func init() {
	host := AuthorityHost(os.Getenv(azureAuthorityHost))
	switch host {
	case AzureChina:
		liveTestScope = "https://management.core.chinacloudapi.cn//.default"
	case AzureGovernment:
		liveTestScope = "https://management.core.usgovcloudapi.net//.default"
	}

	if recording.GetRecordMode() == recording.PlaybackMode {
		liveSP.secret = "fake-secret"
		liveSP.clientID = fakeClientID
		liveSP.tenantID = fakeTenantID
		liveSP.pemPath = "testdata/certificate.pem"
		liveSP.pfxPath = "testdata/certificate.pfx"
		liveSP.sniPath = "testdata/certificate-with-chain.pem"
		liveUser.tenantID = fakeTenantID
		liveUser.username = "fake-user"
		liveUser.password = "fake-password"
	}
}

func TestMain(m *testing.M) {
	if recording.GetRecordMode() == recording.RecordingMode {
		// remove default sanitizers such as the OAuth response sanitizer
		err := recording.ResetSanitizers(nil)
		if err != nil {
			panic(err)
		}
		// replace path variables with fake values to simplify matching (these IDs aren't secret)
		if id, ok := os.LookupEnv("MANAGED_IDENTITY_CLIENT_ID"); ok {
			err = recording.AddURISanitizer(fakeClientID, id, nil)
			if err != nil {
				panic(err)
			}
			err = recording.AddHeaderRegexSanitizer(":path", fakeClientID, id, nil)
			if err != nil {
				panic(err)
			}
		}
		if id, ok := os.LookupEnv("MANAGED_IDENTITY_RESOURCE_ID"); ok {
			replacement := url.QueryEscape(fakeResourceID)
			target := url.QueryEscape(id)
			err = recording.AddURISanitizer(replacement, target, nil)
			if err != nil {
				panic(err)
			}
			err = recording.AddHeaderRegexSanitizer(":path", replacement, target, nil)
			if err != nil {
				panic(err)
			}
		}
		for _, tenantID := range []string{liveSP.tenantID, liveUser.tenantID} {
			if tenantID != "" {
				err = recording.AddURISanitizer(fakeTenantID, tenantID, nil)
				if err != nil {
					panic(err)
				}
				err = recording.AddHeaderRegexSanitizer(":path", fakeTenantID, tenantID, nil)
				if err != nil {
					panic(err)
				}
			}
		}
		// remove token request bodies (which are form encoded) because they contain
		// secrets, are irrelevant in matching, and are formed by MSAL anyway
		// (note: Cloud Shell would need an exemption from this, and that would be okay--its requests contain no secrets)
		err = recording.AddBodyRegexSanitizer("{}", `^\S+=\w+`, nil)
		if err != nil {
			panic(err)
		}
		for _, key := range []string{"access_token", "refresh_token"} {
			err = recording.AddBodyKeySanitizer("$."+key, "redacted", "", nil)
			if err != nil {
				panic(err)
			}
		}
		defer func() {
			err := recording.ResetSanitizers(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	os.Exit(m.Run())
}

func initRecording(t *testing.T) (policy.ClientOptions, func()) {
	err := recording.Start(t, "sdk/azidentity/testdata", nil)
	if err != nil {
		t.Fatal(err)
	}
	transport, err := recording.GetHTTPClient(t)
	if err != nil {
		t.Fatal(err)
	}
	clientOpts := policy.ClientOptions{Transport: transport, PerCallPolicies: []policy.Policy{newRecordingPolicy(t)}}
	return clientOpts, func() {
		err := recording.Stop(t, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type recordingPolicy struct {
	t *testing.T
}

func newRecordingPolicy(t *testing.T) policy.Policy {
	return &recordingPolicy{t: t}
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	mode := recording.GetRecordMode()
	if mode != recording.LiveMode && !recording.IsLiveOnly(p.t) {
		r := req.Raw()
		originalURL := r.URL
		r.Header.Set(recording.IDHeader, recording.GetRecordingId(p.t))
		r.Header.Set(recording.ModeHeader, mode)
		r.Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%s://%s", originalURL.Scheme, originalURL.Host))
		r.Host = "localhost:5001"
		r.URL.Host = r.Host
		r.URL.Scheme = "https"
	}
	return req.Next()
}
