// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sasurl

import "strings"

// Append appends the encoded SAS query string to resourceURL. resourceURL may already carry a
// query string (e.g. "snapshot", "versionid", "sharesnapshot", or custom parameters), in which
// case the SAS is merged into it with "&" instead of introducing a second "?".
func Append(resourceURL, sasQuery string) string {
	if sasQuery == "" {
		return resourceURL
	}
	separator := "?"
	if strings.Contains(resourceURL, "?") {
		separator = "&"
	}
	return resourceURL + separator + sasQuery
}

// AppendToAccountURL is like Append, but first adds a trailing slash to the account URL's path
// (to be consistent with the portal). Any existing query string is split off beforehand so the
// slash lands on the path and not on a query value.
func AppendToAccountURL(accountURL, sasQuery string) string {
	path, rawQuery, hasQuery := strings.Cut(accountURL, "?")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if hasQuery {
		path += "?" + rawQuery
	}
	return Append(path, sasQuery)
}
