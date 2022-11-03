//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// APIVersionOptions contains options for API versions
type APIVersionOptions struct {
	Location APIVersionLocation
	Name     string
}

// APIVersionLocation indicates which part of a request identifies the service version
type APIVersionLocation int

const (
	// APIVersionLocationQueryParam indicates a query parameter
	APIVersionLocationQueryParam = 0
	// APIVersionLocationHeader indicates a header
	APIVersionLocationHeader = 1
)

// newAPIVersionPolicy constructs an APIVersionPolicy. name is the name of the query parameter or header and
// version is its value. If version is "", Do will be a no-op. If version isn't empty and name is empty,
// Do will return an error.
func newAPIVersionPolicy(version string, opts *APIVersionOptions) *apiVersionPolicy {
	if opts == nil {
		opts = &APIVersionOptions{}
	}
	return &apiVersionPolicy{location: opts.Location, name: opts.Name, version: version}
}

// apiVersionPolicy enables users to set the API version of every request a client sends.
type apiVersionPolicy struct {
	// location indicates whether "name" refers to a query parameter or header.
	location APIVersionLocation

	// name of the query param or header whose value should be overridden; provided by the client.
	name string

	// version is the value (provided by the user) that replaces the default version value.
	version string
}

// Do sets the request's API version, if the policy is configured to do so, replacing any prior value.
func (a *apiVersionPolicy) Do(req *policy.Request) (*http.Response, error) {
	if a.version != "" {
		if a.name == "" {
			// user set ClientOptions.APIVersion but the client ctor didn't set PipelineOptions.APIVersionOptions
			return nil, errors.New("this client doesn't support overriding its API version")
		}
		switch a.location {
		case APIVersionLocationHeader:
			req.Raw().Header.Set(a.name, a.version)
		case APIVersionLocationQueryParam:
			q := req.Raw().URL.Query()
			q.Set(a.name, a.version)
			req.Raw().URL.RawQuery = q.Encode()
		default:
			return nil, fmt.Errorf("unknown APIVersionLocation %d", a.location)
		}
	}
	return req.Next()
}
