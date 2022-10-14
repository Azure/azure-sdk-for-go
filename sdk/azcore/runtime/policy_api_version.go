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

const (
	apiVersionHeaderName = iota
	apiVersionQueryParamName
)

// APIVersionName is the name of an API version query parameter or header.
type APIVersionName interface {
	fmt.Stringer
	kind() int
}

// APIVersionHeaderName names the header a client should set with a request's API version.
type APIVersionHeaderName string

func (h APIVersionHeaderName) kind() int {
	return apiVersionHeaderName
}

func (h APIVersionHeaderName) String() string {
	return string(h)
}

// APIVersionQueryParamName names the query parameter a client should set with a request's API version.
type APIVersionQueryParamName string

func (q APIVersionQueryParamName) kind() int {
	return apiVersionQueryParamName
}

func (q APIVersionQueryParamName) String() string {
	return string(q)
}

// APIVersionPolicy enables users to set the API version of every request a client sends.
type APIVersionPolicy struct {
	// name of the query param or header to set. Should be provided by the client.
	name APIVersionName
	// version value. Should be provided by the user.
	version string
}

// NewAPIVersionPolicy constructs an APIVersionPolicy. If version is "", Do will be a no-op.
func NewAPIVersionPolicy(name APIVersionName, version string) *APIVersionPolicy {
	return &APIVersionPolicy{name, version}
}

// Do sets the request's API version, if the policy is configured to do so, replacing any prior value.
func (a *APIVersionPolicy) Do(req *policy.Request) (*http.Response, error) {
	if a.version != "" {
		if a.name == nil {
			// user set ClientOptions.APIVersion but the client ctor didn't set PipelineOptions.APIVersionName
			return nil, errors.New("this client doesn't support overriding its API version")
		}
		switch a.name.kind() {
		case apiVersionHeaderName:
			req.Raw().Header.Set(a.name.String(), a.version)
		case apiVersionQueryParamName:
			q := req.Raw().URL.Query()
			q.Set(a.name.String(), a.version)
			req.Raw().URL.RawQuery = q.Encode()
		}
	}
	return req.Next()
}
