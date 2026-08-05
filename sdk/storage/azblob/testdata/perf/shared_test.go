// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
	"math"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

type staticTokenCredential struct{}

func (staticTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func resetCredentialState(t *testing.T) {
	t.Helper()
	credentialMu.Lock()
	savedCredential := sharedCredential
	savedFactory := credentialFactory
	sharedCredential = nil
	credentialMu.Unlock()
	t.Cleanup(func() {
		credentialMu.Lock()
		sharedCredential = savedCredential
		credentialFactory = savedFactory
		credentialMu.Unlock()
	})
}

func TestNewContainerClientReusesOAuthCredential(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "https://account.example.test/")
	t.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "")
	t.Setenv("AZURE_STORAGE_CONNECTION_STRING", "")
	var calls atomic.Int32
	credentialFactory = func() (azcore.TokenCredential, error) {
		calls.Add(1)
		return staticTokenCredential{}, nil
	}

	first, err := newContainerClient("first", nil)
	require.NoError(t, err)
	second, err := newContainerClient("second", nil)
	require.NoError(t, err)

	require.Equal(t, "https://account.example.test/first", first.URL())
	require.Equal(t, "https://account.example.test/second", second.URL())
	require.Equal(t, int32(1), calls.Load())
}

func TestNewContainerClientBuildsPublicCloudURLFromAccountName(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "")
	t.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "account")
	t.Setenv("AZURE_STORAGE_CONNECTION_STRING", "")
	credentialFactory = func() (azcore.TokenCredential, error) { return staticTokenCredential{}, nil }

	client, err := newContainerClient("container", nil)

	require.NoError(t, err)
	require.Equal(t, "https://account.blob.core.windows.net/container", client.URL())
}

func TestNewContainerClientReturnsCredentialFactoryError(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "https://account.example.test")
	credentialErr := errors.New("credential unavailable")
	credentialFactory = func() (azcore.TokenCredential, error) { return nil, credentialErr }

	_, err := newContainerClient("container", nil)

	require.ErrorIs(t, err, credentialErr)
}

func TestNewContainerClientUsesConnectionStringFallback(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "")
	t.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "")
	t.Setenv("AZURE_STORAGE_CONNECTION_STRING", "DefaultEndpointsProtocol=https;AccountName=test;AccountKey=dGVzdA==;EndpointSuffix=core.windows.net")

	client, err := newContainerClient("container", nil)

	require.NoError(t, err)
	require.Equal(t, "https://test.blob.core.windows.net/container", client.URL())
}

func TestNewContainerClientRejectsInvalidAccountURL(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "://bad")
	t.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "")

	_, err := newContainerClient("container", nil)

	require.Error(t, err)
}

func TestNewContainerClientRequiresCredentials(t *testing.T) {
	resetCredentialState(t)
	t.Setenv("AZURE_STORAGE_ACCOUNT_URL", "")
	t.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "")
	t.Setenv("AZURE_STORAGE_CONNECTION_STRING", "")

	_, err := newContainerClient("container", nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "no storage credentials found")
}

func TestValidatePerfOptions(t *testing.T) {
	savedUploadMethod, savedDownloadMethod := uploadMethod, downloadMethod
	savedBlockSize, savedPageSize := commonBlockSize, listPageSize
	savedListOptions := listTestOpts
	t.Cleanup(func() {
		uploadMethod, downloadMethod = savedUploadMethod, savedDownloadMethod
		commonBlockSize, listPageSize = savedBlockSize, savedPageSize
		listTestOpts = savedListOptions
	})

	require.Error(t, validateTransferOptions("upload", -1))
	uploadMethod = "invalid"
	require.Error(t, validateTransferOptions("upload", 1))
	uploadMethod = "single"
	downloadMethod = "invalid"
	require.Error(t, validateTransferOptions("download", 1))
	downloadMethod = "stream"
	commonBlockSize = -1
	require.Error(t, validateTransferOptions("upload", 1))
	commonBlockSize = 0
	listTestOpts.count = -1
	require.Error(t, validateListOptions())
	listTestOpts.count = 1
	listTestOpts.parallelism = -1
	require.Error(t, validateListOptions())
	listTestOpts.parallelism = 0
	listPageSize = int(math.MaxInt32) + 1
	require.Error(t, validateListOptions())
	listPageSize = 0
	require.NoError(t, validateListOptions())
	require.NoError(t, validateTransferOptions("upload", 0))
}

func TestUint16Flag(t *testing.T) {
	var value uint16Flag
	require.Equal(t, "0", value.String())
	require.Equal(t, "0", (*uint16Flag)(nil).String())
	require.NoError(t, value.Set("65535"))
	require.Equal(t, "65535", value.String())
	require.Error(t, value.Set("65536"))
	require.Error(t, value.Set("-1"))
	require.Error(t, value.Set("not-a-number"))
}

func TestGenerateRandomBytes(t *testing.T) {
	_, err := generateRandomBytes(-1)
	require.Error(t, err)
	zero, err := generateRandomBytes(0)
	require.NoError(t, err)
	require.Empty(t, zero)
	data, err := generateRandomBytes(32)
	require.NoError(t, err)
	require.Len(t, data, 32)
}

func TestCleanupContainerOnError(t *testing.T) {
	transport := &failingDeleteTransport{}
	client, err := container.NewClientWithNoCredential("https://example.test/container", &container.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
	})
	require.NoError(t, err)

	var noErr error
	cleanupContainerOnError(&noErr, client)
	require.Zero(t, transport.calls.Load())
	original := errors.New("setup failed")
	cleanupContainerOnError(&original, nil)
	require.ErrorContains(t, original, "setup failed")
	cleanupContainerOnError(&original, client)
	require.ErrorContains(t, original, "setup failed")
	require.ErrorIs(t, original, errDeleteContainer)
}

var errDeleteContainer = errors.New("delete failed")

type failingDeleteTransport struct{ calls atomic.Int32 }

func (t *failingDeleteTransport) Do(*http.Request) (*http.Response, error) {
	t.calls.Add(1)
	return nil, errDeleteContainer
}

type failingSetupTransport struct {
	deletes atomic.Int32
}

var errSetupUpload = errors.New("setup upload failed")

func (t *failingSetupTransport) Do(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodDelete {
		t.deletes.Add(1)
		return setupResponse(req, http.StatusAccepted), nil
	}
	if req.URL.Query().Get("restype") == "container" {
		return setupResponse(req, http.StatusCreated), nil
	}
	return nil, errSetupUpload
}

func setupResponse(req *http.Request, status int) *http.Response {
	return &http.Response{Request: req, StatusCode: status, Status: http.StatusText(status), Header: http.Header{}, Body: http.NoBody}
}

func TestListSetupDeletesContainerAfterSeedFailure(t *testing.T) {
	transport := &failingSetupTransport{}
	savedFactory := containerClientFactory
	savedOptions := listTestOpts
	containerClientFactory = func(containerName string, _ *container.ClientOptions) (*container.Client, error) {
		return container.NewClientWithNoCredential("https://example.test/"+containerName, &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
	}
	listTestOpts = listTestOptions{count: 1, parallelism: 1}
	t.Cleanup(func() {
		containerClientFactory = savedFactory
		listTestOpts = savedOptions
	})

	_, err := NewListTest(context.Background(), perf.PerfTestOptions{})

	require.ErrorIs(t, err, errSetupUpload)
	require.Equal(t, int32(1), transport.deletes.Load())
}

func TestDownloadSetupDeletesContainerAfterSeedFailure(t *testing.T) {
	transport := &failingSetupTransport{}
	savedFactory := containerClientFactory
	savedOptions := downloadTestOpts
	savedMethod := downloadMethod
	containerClientFactory = func(containerName string, _ *container.ClientOptions) (*container.Client, error) {
		return container.NewClientWithNoCredential("https://example.test/"+containerName, &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
	}
	downloadTestOpts.size = 4
	downloadMethod = "stream"
	t.Cleanup(func() {
		containerClientFactory = savedFactory
		downloadTestOpts = savedOptions
		downloadMethod = savedMethod
	})

	_, err := NewDownloadTest(context.Background(), perf.PerfTestOptions{})

	require.ErrorIs(t, err, errSetupUpload)
	require.Equal(t, int32(1), transport.deletes.Load())
}
