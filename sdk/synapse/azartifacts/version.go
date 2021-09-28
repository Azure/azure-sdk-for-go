//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azartifacts

const (
	// UserAgent is the string to be used in the user agent string when making requests.
	UserAgent = "azartifacts/" + Version

	// Version is the semantic version (see http://semver.org) of this module.
	Version = "v0.1.0"
)
