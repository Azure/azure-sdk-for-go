//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package querynormalization

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// Policy is a pipeline policy for normalizing query parameters in HTTP requests.
// It converts all query parameter names to lowercase and sorts them alphabetically
// while preserving the relative order of duplicate parameter names.
// Don't use this type directly, use NewPolicy() instead.
type Policy struct{}

// NewPolicy creates a new instance of Policy.
func NewPolicy() *Policy {
	return &Policy{}
}

// Do implements the policy.Policy interface on type Policy.
func (p *Policy) Do(req *policy.Request) (*http.Response, error) {
	// Normalize query parameters with error handling
	rawURL := req.Raw().URL
	if rawURL == nil || rawURL.RawQuery == "" {
		return req.Next()
	}

	// Create a slice to hold parameter entries with lowercase names
	type paramEntry struct {
		lowercaseName string
		value         string
	}

	var params []paramEntry

	// Split rawURL by &, and traverse
	for _, pair := range strings.Split(rawURL.RawQuery, "&") {
		// Split each pair by =
		parts := strings.SplitN(pair, "=", 2)

		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}

		params = append(params, paramEntry{
			lowercaseName: strings.ToLower(parts[0]),
			value:         value,
		})
	}

	// Sort by lowercase name to achieve case-insensitive alphabetical ordering
	// Go's sort.Slice is stable, so the relative order of entries with the same
	// lowercase name is preserved
	sort.SliceStable(params, func(i, j int) bool {
		return params[i].lowercaseName < params[j].lowercaseName
	})

	// Build new query parameters with lowercase names in sorted order
	var newQuery []string

	for _, entry := range params {
		if entry.lowercaseName != "" {
			newQuery = append(newQuery, fmt.Sprintf("%s=%s", entry.lowercaseName, entry.value))
		} else {
			newQuery = append(newQuery, "")
		}
	}

	// Update the request URL with normalized query parameters
	rawURL.RawQuery = strings.Join(newQuery, "&")

	return req.Next()
}
