package foo

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// UserAgent returns the UserAgent string to use when sending http.Requests.
func UserAgent() string {
	return "Azure-SDK-For-Go/1.0.0 foo/2019-04-01"
}

// Version returns the semantic version (see http://semver.org) of the client.
func Version() string {
	return "1.0.0"
}
