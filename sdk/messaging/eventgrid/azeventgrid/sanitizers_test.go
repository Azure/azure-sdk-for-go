//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// addSanitizers makes sure any fields that are client-derived (ie, timestamps, etc..) are
// made consistent with what was recorded.
func addSanitizers(t *testing.T) {
	sanitizeEventFields(t)

	err := recording.AddURISanitizer(fakeVars.EG.Endpoint, `^.+\?`, nil)
	require.NoError(t, err)
}

func sanitizeEventFields(t *testing.T) {
	err := recording.AddBodyRegexSanitizer("\"time\":\"2024-01-01T01:01:01.111111111Z\"", "\"time\":\"[^\"]+\"", nil)
	require.NoError(t, err)

	err = recording.AddBodyRegexSanitizer("\"eventTime\":\"2024-01-01T01:01:01.111111111Z\"", "\"eventTime\":\"[^\"]+\"", nil)
	require.NoError(t, err)

	err = recording.AddBodyRegexSanitizer("\"id\":\"11111111-1111-1111-1111-111111111111\"", "\"id\":\"[^\"]+\"", nil)
	require.NoError(t, err)
}
