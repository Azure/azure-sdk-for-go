// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

const (
	// UserAgent is the string to be used in the user agent string when making requests.
	UserAgent = "azcore/" + Version

	// Version is the semantic version (see http://semver.org) of the pipeline package.
	Version = "0.1.0"
)
