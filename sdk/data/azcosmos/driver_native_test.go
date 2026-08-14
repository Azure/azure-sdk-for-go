// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// emulatorKey is the Cosmos DB emulator's well-known account key, which is published in the
// emulator documentation and grants nothing anywhere else.
const emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

// The runtime and the account reference are built without contacting the service, so this covers
// acquiring and releasing them anywhere. Creating the driver is what needs a reachable account;
// see TestNativeDriverCreationContactsTheAccount.
func TestNativeRuntimeAndAccountLifecycle(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	require.NotNil(t, d.runtime)

	require.NoError(t, d.buildAccount(driverConfig{
		endpoint:   "https://myaccount.documents.azure.com",
		accountKey: emulatorKey,
	}))
	require.NotNil(t, d.account)

	require.NoError(t, d.close())
	require.Nil(t, d.runtime, "close must not leave a dangling handle")
	require.Nil(t, d.account)
	require.NoError(t, d.close(), "close is idempotent")
}

// Repeated cycles catch a handle that close forgets to release, which a single cycle would not.
func TestNativeLifecycleIsRepeatable(t *testing.T) {
	for range 20 {
		d := &nativeDriver{}
		require.NoError(t, d.buildRuntime())
		require.NoError(t, d.buildAccount(driverConfig{
			endpoint:   "https://myaccount.documents.azure.com",
			accountKey: emulatorKey,
		}))
		require.NoError(t, d.close())
	}
}

// The endpoint is validated by the driver when the account reference is built, before any
// credential is sent. It classifies as a client error rather than a bad request: no service was
// contacted, so there is nothing to blame the service for.
func TestNativeAccountRejectsMalformedEndpoint(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	t.Cleanup(func() { _ = d.close() })

	err := d.buildAccount(driverConfig{endpoint: "not-a-url", accountKey: emulatorKey})
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientError, cosmosErr.Code)
	require.False(t, cosmosErr.FromWire)
	// The driver's own description has to survive translation, or the caller is left with a
	// status and no way to tell which of the inputs was wrong.
	require.Contains(t, cosmosErr.Message, "building the account reference")
	require.Contains(t, cosmosErr.Message, "account endpoint URL")
}

// The driver supports token credentials but the C ABI has no constructor for one, so this is the
// binding refusing rather than the driver.
func TestNativeTokenCredentialRejected(t *testing.T) {
	_, err := NewClient("https://myaccount.documents.azure.com", fakeTokenCredential{}, nil)
	require.ErrorIs(t, err, errTokenCredentialUnsupported)
}

// Creating the driver fetches the account's properties, so it does network I/O and is deferred to
// first use rather than done in the constructor. This exercises that first use.
func TestNativeDriverCreationContactsTheAccount(t *testing.T) {
	endpoint := os.Getenv("AZCOSMOS_EMULATOR_ENDPOINT")
	if endpoint == "" {
		t.Skip("set AZCOSMOS_EMULATOR_ENDPOINT to run against a live emulator")
	}

	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	client, err := NewClientWithKey(endpoint, cred, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	driver, err := client.driver.ensureDriver()
	require.NoError(t, err)
	require.NotNil(t, driver)

	// The result is cached, so a second call neither refetches nor returns a different handle.
	again, err := client.driver.ensureDriver()
	require.NoError(t, err)
	require.Equal(t, driver, again)
}

// Construction must not reach the network: an unreachable account is an operation failure, not a
// construction failure, and the constructor has no context to cancel a network call with.
func TestNativeConstructionDoesNotContactTheAccount(t *testing.T) {
	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	// Port 1 is reserved and never listening, so reaching it would fail rather than hang.
	client, err := NewClientWithKey("https://127.0.0.1:1", cred, nil)
	require.NoError(t, err, "construction should not depend on the account being reachable")
	require.NoError(t, client.Close())
}

// A closed driver has to report that rather than hand back a NULL handle, which the C ABI would
// reject with a confusing null-argument error at best and dereference at worst.
func TestNativeEnsureDriverAfterCloseIsRejected(t *testing.T) {
	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", cred, nil)
	require.NoError(t, err)

	// Captured before Close, which drops the client's reference to it.
	d := client.driver
	require.NoError(t, client.Close())

	driver, err := d.ensureDriver()
	require.Nil(t, driver, "a closed client must not yield a handle")
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

// An unreachable account has to report that the transport failed rather than that the service was
// unavailable. The driver pairs transport failures with a synthetic 503, so classifying on the
// status alone would blame the service for a local failure.
func TestNativeUnreachableAccountReportsTransportFailure(t *testing.T) {
	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)

	client, err := NewClientWithKey("https://127.0.0.1:1", cred, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	_, err = client.driver.ensureDriver()
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeTransportFailure, cosmosErr.Code)
	require.False(t, cosmosErr.FromWire, "no service responded")

	// The failure is remembered rather than retried on every operation.
	_, again := client.driver.ensureDriver()
	require.Equal(t, err, again)
}
