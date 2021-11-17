//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func ExampleStart() {
	err := recording.Start(&testing.T{}, "path/to/sdk/testdata", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleStop() {
	err := recording.Stop(&testing.T{}, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleResetSanitizer() {
	err := recording.ResetSanitizers(nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddBodyKeySanitizer() {
	err := recording.AddBodyKeySanitizer("$.json.path", "new-value", "regex-to-replace", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddBodyRegexSanitizer() {
	err := recording.AddBodyRegexSanitizer("my-new-value", "regex-to-replace", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddContinuationSanitizer() {
	err := recording.AddContinuationSanitizer("key", "my-new-value", true, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddGeneralRegexSanitizer() {
	err := recording.AddGeneralRegexSanitizer("my-new-value", "regex-to-scrub-secret", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddHeaderRegexSanitizer() {
	err := recording.AddHeaderRegexSanitizer("header", "my-new-value", "regex-to-scrub-secret", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddOAuthResponseSanitizer() {
	err := recording.AddOAuthResponseSanitizer(nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddRemoveHeaderSanitizer() {
	err := recording.AddRemoveHeaderSanitizer([]string{"header1", "header2"}, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddURISanitizer() {
	err := recording.AddURISanitizer("my-new-value", "my-secret", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleAddURISubscriptionIDSanitizer() {
	err := recording.AddURISubscriptionIDSanitizer("0123-4567-...", nil)
	if err != nil {
		panic(err)
	}
}
