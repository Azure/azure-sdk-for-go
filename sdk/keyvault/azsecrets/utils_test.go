//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var pathToPackage = "sdk/keyvault/azsecrets"

const headerAuthorization = "Authorization"

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{}
	}
	p := &recordingPolicy{options: *o, t: t}
	p.options.Init()
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" {
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = "https"
		req.Raw().URL.Host = p.options.Host
		req.Raw().Host = p.options.Host

		req.Raw().Header.Set(recording.UpstreamUriHeader, fmt.Sprintf("%v://%v", p.options.Scheme, originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IdHeader, recording.GetRecordingId(p.t))
	}
	return req.Next()
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

func createClient(t *testing.T) (*Client, error) {
	vaultUrl := recording.GetEnvVariable(t, "AZURE_KEYVAULT_URL", "https://fakekvurl.vault.azure.net/")
	if recording.GetRecordMode() == "playback" {
		vaultUrl = "https://fakekvurl.vault.azure.net/"
	}

	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       client,
	}
	_ = options

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZSECRETS_TENANT_ID")
		clientId := lookupEnvVar("AZSECRETS_CLIENT_ID")
		clientSecret := lookupEnvVar("AZSECRETS_CLIENT_SECRET")
		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		require.NoError(t, err)
	} else {
		cred = NewFakeCredential("fake", "fake")
	}

	return NewClient(vaultUrl, cred, options)
}

func delay() time.Duration {
	if recording.GetRecordMode() == "playback" {
		return 1 * time.Microsecond
	}
	return 250 * time.Millisecond
}

func cleanUpSecret(t *testing.T, client *Client, secret string) {
	resp, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.PurgeDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)
}

type FakeCredential struct {
	accountName string
	accountKey  string
}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}
}

type fakeCredPolicy struct {
	cred *FakeCredential
}

func newFakeCredPolicy(cred *FakeCredential, opts runtime.AuthenticationOptions) *fakeCredPolicy {
	return &fakeCredPolicy{cred: cred}
}

func (f *fakeCredPolicy) Do(req *policy.Request) (*http.Response, error) {
	authHeader := strings.Join([]string{"Authorization ", f.cred.accountName, ":", f.cred.accountKey}, "")
	req.Raw().Header.Set(headerAuthorization, authHeader)
	return req.Next()
}

func (f *FakeCredential) NewAuthenticationPolicy(options runtime.AuthenticationOptions) policy.Policy {
	return newFakeCredPolicy(f, options)
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "faketoken",
		ExpiresOn: time.Date(2040, time.January, 1, 1, 1, 1, 1, time.UTC),
	}, nil
}
