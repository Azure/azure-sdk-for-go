//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// sanitizeForPlayback makes sure any fields that are client-derived (ie, timestamps, etc..) are
// made consistent with what was recorded.
func sanitizeForPlayback(t *testing.T) {
	sanitizeEventFields(t)
	sanitizeTokenCreds(t)
}

// sanitizeForRecording sanitizes fields that are recorded to match what we use as  our "fake" values
// for playback.
func sanitizeForRecording(t *testing.T, egVars eventGridVars) {
	sanitizeEventFields(t)

	err := recording.AddURISanitizer(fakeVars.CE.Endpoint, regexp.QuoteMeta(egVars.CE.Endpoint), nil)
	require.NoError(t, err)

	err = recording.AddURISanitizer(fakeVars.EG.Endpoint, regexp.QuoteMeta(egVars.EG.Endpoint), nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("Aeg-Sas-Token", fakeVars.EG.SAS, "", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("Aeg-Sas-Key", fakeVars.EG.Key, "", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("Client-Request-Id", "client-request-id", "", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("client-request-id", "client-request-id", "", nil)
	require.NoError(t, err)

	sanitizeTokenCreds(t)
}

func sanitizeEventFields(t *testing.T) {
	err := recording.AddBodyRegexSanitizer("\"time\":\"2024-01-01T01:01:01.111111111Z\"", "\"time\":\"[^\"]+\"", nil)
	require.NoError(t, err)

	err = recording.AddBodyRegexSanitizer("\"eventTime\":\"2024-01-01T01:01:01.111111111Z\"", "\"eventTime\":\"[^\"]+\"", nil)
	require.NoError(t, err)

	err = recording.AddBodyRegexSanitizer("\"id\":\"11111111-1111-1111-1111-111111111111\"", "\"id\":\"[^\"]+\"", nil)
	require.NoError(t, err)
}

// sanitizeTokenCreds sanitizes all the fields we use for TokenCredential auth, like the tenant, secrets, etc..
func sanitizeTokenCreds(t *testing.T) {
	err := recording.AddBodyRegexSanitizer("client_id=fake-client-id", "client_id=.+?&", nil)
	require.NoError(t, err)

	err = recording.AddBodyRegexSanitizer("client_secret=fake-client-secret&", "client_secret=.+&", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("Client-Request-Id", "client-request-id", "", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("client-request-id", "client-request-id", "", nil)
	require.NoError(t, err)

	// sanitize the os and go fields that identity sends
	err = recording.AddHeaderRegexSanitizer("X-Client-Cpu", "amd641", ".+", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("X-Client-Os", "linux1", ".+", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("X-Client-Ver", "1.2.01", ".+", nil)
	require.NoError(t, err)
}
