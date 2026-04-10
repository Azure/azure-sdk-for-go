// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package querynormalization

import (
	"net/http"
	"net/url"
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

// paramEntry holds a lowercased parameter name with its value and original index
// to support stable sorting.
type paramEntry struct {
	lowercaseName string
	value         string
	index         int
}

// Do implements the policy.Policy interface on type Policy.
func (p *Policy) Do(req *policy.Request) (*http.Response, error) {
	rawURL := req.Raw().URL
	if rawURL == nil || rawURL.RawQuery == "" {
		return req.Next()
	}

	queryParams := rawURL.Query()

	// Build a flat list of paramEntry to preserve insertion order of duplicate values
	var params []paramEntry
	idx := 0
	for name, values := range queryParams {
		lowerName := strings.ToLower(name)
		for _, v := range values {
			params = append(params, paramEntry{
				lowercaseName: lowerName,
				value:         v,
				index:         idx,
			})
			idx++
		}
	}

	// Sort by lowercase name; use original index as tiebreaker to maintain
	// stable ordering for parameters with the same lowercased name
	sort.SliceStable(params, func(i, j int) bool {
		if params[i].lowercaseName != params[j].lowercaseName {
			return params[i].lowercaseName < params[j].lowercaseName
		}
		return params[i].index < params[j].index
	})

	var parts []string
	for _, entry := range params {
		parts = append(parts, url.QueryEscape(entry.lowercaseName)+"="+url.QueryEscape(entry.value))
	}
	rawURL.RawQuery = strings.Join(parts, "&")

	return req.Next()
}
