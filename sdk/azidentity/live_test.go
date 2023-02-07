//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

// configuration for live tests
var liveManagedIdentity = struct {
	clientID   string
	resourceID string
}{
	clientID:   os.Getenv("MANAGED_IDENTITY_CLIENT_ID"),
	resourceID: os.Getenv("MANAGED_IDENTITY_RESOURCE_ID"),
}

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
	fakeClientID      = "fake-client-id"
	fakeResourceID    = "/fake/resource/ID"
	fakeTenantID      = "fake-tenant"
	fakeUsername      = "fake@user"
	fakeAdfsAuthority = "fake.adfs.local"
	fakeAdfsScope     = "fake.adfs.local/fake-scope/.default"
)

var adfsLiveSP = struct {
	clientID    string
	secret      string
	certPath    string
	redirectURL string
}{
	clientID:    os.Getenv("ADFS_SP_CLIENT_ID"),
	secret:      os.Getenv("ADFS_SP_CLIENT_SECRET"),
	certPath:    os.Getenv("ADFS_SP_CERT_PATH"),
	redirectURL: os.Getenv("ADFS_SP_REDIRECT_URL"),
}

var adfsLiveUser = struct {
	clientID string
	username string
	password string
}{
	username: os.Getenv("ADFS_IDENTITY_TEST_USERNAME"),
	password: os.Getenv("ADFS_IDENTITY_TEST_PASSWORD"),
	clientID: os.Getenv("ADFS_IDENTITY_TEST_CLIENT_ID"),
}

var (
	adfsAuthority     = os.Getenv("ADFS_AUTHORITY_HOST")
	adfsScope         = os.Getenv("ADFS_SCOPE")
	liveTestScope     = "https://management.core.windows.net//.default"
	_, runManualTests = os.LookupEnv("AZIDENTITY_RUN_MANUAL_TESTS")
)

func init() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		liveManagedIdentity.clientID = fakeClientID
		liveManagedIdentity.resourceID = fakeResourceID
		liveSP.secret = "fake-secret"
		liveSP.clientID = fakeClientID
		liveSP.tenantID = fakeTenantID
		liveSP.pemPath = "testdata/certificate.pem"
		liveSP.pfxPath = "testdata/certificate.pfx"
		liveSP.sniPath = "testdata/certificate-with-chain.pem"
		liveUser.tenantID = fakeTenantID
		liveUser.username = fakeUsername
		liveUser.password = "fake-password"
		adfsLiveSP.secret = "fake-secret"
		adfsLiveSP.clientID = fakeClientID
		adfsLiveSP.certPath = "testdata/certificate.pem"
		adfsLiveUser.username = fakeUsername
		adfsLiveUser.password = "fake-password"
		adfsLiveUser.clientID = fakeClientID
		adfsScope = "https://" + fakeAdfsScope
		adfsAuthority = "https://" + fakeAdfsAuthority
	}
}

func TestMain(m *testing.M) {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		// Start from a fresh proxy
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}

		// At the end of testing we want to reset as to not interfere with other tests.
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	switch recording.GetRecordMode() {
	case recording.PlaybackMode:
		err := recording.SetBodilessMatcher(nil, nil)
		if err != nil {
			panic(err)
		}
	case recording.RecordingMode:
		// replace path variables with fake values to simplify matching (the real values aren't secret)
		pathVars := map[string]string{
			liveManagedIdentity.clientID:                    fakeClientID,
			url.QueryEscape(liveManagedIdentity.resourceID): url.QueryEscape(fakeResourceID),
			liveSP.tenantID:                                 fakeTenantID,
			liveUser.tenantID:                               fakeTenantID,
			liveUser.username:                               fakeUsername,
			strings.TrimPrefix(adfsScope, "https://"):       fakeAdfsScope,
			strings.TrimPrefix(adfsAuthority, "https://"):   fakeAdfsAuthority,
		}
		for target, replacement := range pathVars {
			if target != "" {
				err := recording.AddURISanitizer(replacement, target, nil)
				if err != nil {
					panic(err)
				}
				err = recording.AddHeaderRegexSanitizer(":path", replacement, target, nil)
				if err != nil {
					panic(err)
				}
				// replace e.g. the real tenant ID with the fake one in metadata discovery responses,
				// to ensure MSAL sends requests to the expected URI in playback
				err = recording.AddBodyRegexSanitizer(replacement, target, nil)
				if err != nil {
					panic(err)
				}
			}
		}
		// remove token request bodies (which are form encoded) because they contain
		// secrets, are irrelevant in matching, and are formed by MSAL anyway
		// (note: Cloud Shell would need an exemption from this, and that would be okay--its requests contain no secrets)
		err := recording.AddBodyRegexSanitizer("{}", `^\S+=.*`, nil)
		if err != nil {
			panic(err)
		}
		// redact secrets returned by AAD
		for _, key := range []string{"access_token", "device_code", "message", "refresh_token", "user_code"} {
			err = recording.AddBodyKeySanitizer("$."+key, "redacted", "", nil)
			if err != nil {
				panic(err)
			}
		}
	}
	val := m.Run()

	os.Exit(val)
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

// testGetTokenSuccess is a helper for happy path tests that acquires, and validates, a token from a credential
func testGetTokenSuccess(t *testing.T, cred azcore.TokenCredential, customScope ...string) {
	var scopes []string
	if customScope == nil {
		scopes = append(scopes, liveTestScope)
	} else {
		scopes = append(scopes, customScope...)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token == "" {
		t.Fatal("GetToken returned an invalid token")
	}
	if tk.ExpiresOn.Before(time.Now()) {
		t.Fatal("GetToken returned an invalid expiration time")
	}
	if tk.ExpiresOn.Location() != time.UTC {
		t.Fatal("ExpiresOn isn't UTC")
	}
	tk2, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		t.Fatal(err)
	}
	if tk2.Token != tk.Token || tk2.ExpiresOn != tk.ExpiresOn {
		t.Fatal("expected a cached token")
	}
}
