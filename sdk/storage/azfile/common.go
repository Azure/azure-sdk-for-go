// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

import "net/url"

// appendToURLPath appends a string to the end of a URL's path (prefixing the string with a '/' if required)
func appendToURLPath(u string, name string) string {
	// e.g. "https://ms.com/a/b/?k1=v1&k2=v2#f"
	// When you call url.Parse() this is what you'll get:
	//     Scheme: "https"
	//     Opaque: ""
	//       User: nil
	//       Host: "ms.com"
	//       Path: "/a/b/"	This should start with a / and it might or might not have a trailing slash
	//    RawPath: ""
	// ForceQuery: false
	//   RawQuery: "k1=v1&k2=v2"
	//   Fragment: "f"
	uri, _ := url.Parse(u)

	if len(uri.Path) == 0 || uri.Path[len(uri.Path)-1] != '/' {
		uri.Path += "/" // Append "/" to end before appending name
	}
	uri.Path += name
	return uri.String()
}
