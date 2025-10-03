//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package runtime

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestGetAzureStackAPIVersion(t *testing.T) {
	testCases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "empty env var",
			envValue: "",
			expected: "",
		},
		{
			name:     "wildcard returns default",
			envValue: "*",
			expected: DefaultAzureStackAPIVersion,
		},
		{
			name:     "custom version",
			envValue: "2020-01-01",
			expected: "2020-01-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue == "" {
				os.Unsetenv(AzureStackStorageAPIVersionEnvVar)
			} else {
				os.Setenv(AzureStackStorageAPIVersionEnvVar, tc.envValue)
			}
			defer os.Unsetenv(AzureStackStorageAPIVersionEnvVar)

			if got := GetAzureStackAPIVersion(); got != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestApplyAzureStackAPIVersion(t *testing.T) {
	testCases := []struct {
		name      string
		envValue  string
		options   *policy.ClientOptions
		expected  string
		shouldSet bool
	}{
		{
			name:      "nil options",
			envValue:  "2020-01-01",
			options:   nil,
			expected:  "",
			shouldSet: false,
		},
		{
			name:      "no env var",
			envValue:  "",
			options:   &policy.ClientOptions{},
			expected:  "",
			shouldSet: false,
		},
		{
			name:      "with env var",
			envValue:  "2020-01-01",
