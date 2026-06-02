// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package synctoken

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/exported"
)

// Policy is a pipeline policy for managing Sync-Token
// values in HTTP requests and responses.
// Don't use this type directly, use NewPolicy() instead.
type Policy struct {
	cache *Cache
}

// NewPolicy creates a new instance of Policy.
func NewPolicy(cache *Cache) *Policy {
	return &Policy{
		cache: cache,
	}
}

// Do implements the policy.Policy interface on type Policy.
func (p *Policy) Do(req *policy.Request) (*http.Response, error) {
	// add the sync token to the HTTP request
	if st := p.cache.Get(); st != "" {
		req.Raw().Header[syncTokenHeader] = []string{st}
	}

	resp, err := req.Next()
	if err != nil {
		return nil, err
	}

	// update the cache from the response if available.
	// e.g. a 404 will include a Sync-Token but a 400 will not.
	if st := resp.Header.Get(syncTokenHeader); st != "" {
		if err := p.cache.Set(exported.SyncToken(st)); err != nil {
			return nil, &nonRetriableError{err}
		}
	}

	return resp, err
}

const syncTokenHeader = "Sync-Token"

type nonRetriableError struct {
	error
}

func (*nonRetriableError) NonRetriable() {
	// marker method
}
