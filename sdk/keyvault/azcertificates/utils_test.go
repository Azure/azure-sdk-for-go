//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var pathToPackage = "sdk/keyvault/azcertificates/testdata"

const fakeKvURL = "https://fakekvurl.vault.azure.net/"

var certContentNotPasswordEncoded = "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"

func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}

		err = recording.SetDefaultMatcher(nil, nil)
		if err != nil {
			panic(err)
		}

		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err = recording.AddGeneralRegexSanitizer(fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		err = recording.AddBodyKeySanitizer("$.key.kid", fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		err = recording.AddBodyKeySanitizer("$.recoveryId", fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		tenantID := os.Getenv("AZCERTIFICATES_TENANT_ID")
		err = recording.AddHeaderRegexSanitizer("WWW-Authenticate", "00000000-0000-0000-0000-000000000000", tenantID, nil)
		if err != nil {
			panic(err)
		}
	}

	// Run tests
	exitVal := m.Run()

	// 3. Reset
	if recording.GetRecordMode() != "live" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

	// 4. Error out if applicable
	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

func createClient(t *testing.T) (*Client, error) {
	vaultUrl := recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeKvURL)
	var credOptions *azidentity.ClientSecretCredentialOptions

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
		},
	}

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZCERTIFICATES_TENANT_ID")
		clientId := lookupEnvVar("AZCERTIFICATES_CLIENT_ID")
		clientSecret := lookupEnvVar("AZCERTIFICATES_CLIENT_SECRET")
		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, credOptions)
		require.NoError(t, err)
	} else {
		cred = NewFakeCredential("fake", "fake")
	}

	return NewClient(vaultUrl, cred, options)
}

func delay() time.Duration {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return time.Microsecond
	}
	return time.Second
}

func longDelay() time.Duration {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return 1 * time.Microsecond
	}
	return 3 * time.Second
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

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "faketoken",
		ExpiresOn: time.Date(2040, time.January, 1, 1, 1, 1, 1, time.UTC),
	}, nil
}

func cleanUp(t *testing.T, client *Client, certName string) {
	delResp, err := client.BeginDeleteCertificate(ctx, certName, nil)
	if err == nil {
		delPollerResp, err := delResp.PollUntilDone(ctx, delay())
		require.NoError(t, err)
		require.Contains(t, *delPollerResp.ID, certName)

		_, err = client.PurgeDeletedCertificate(ctx, certName, nil)
		require.NoError(t, err)
	}
}

func createCert(t *testing.T, client *Client, certName string) {
	resp, err := client.BeginCreateCertificate(ctx, certName, CertificatePolicy{
		IssuerParameters: &IssuerParameters{
			Name: to.StringPtr("Self"),
		},
		X509CertificateProperties: &X509CertificateProperties{
			Subject: to.StringPtr("CN=DefaultPolicy"),
		},
	}, nil)
	require.NoError(t, err)
	_, err = resp.PollUntilDone(ctx, delay())
	require.NoError(t, err)
}

func purgeCert(t *testing.T, client *Client, cert string) {
	_, err := client.PurgeDeletedCertificate(ctx, cert, nil)
	if err != nil {
		log.Printf("Was unable to purge certificate %s. %s\n", cert, err.Error())
	}
}
