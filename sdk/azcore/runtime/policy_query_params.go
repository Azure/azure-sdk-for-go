//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// queryParamsPolicy is a policy that adds custom query parameters to a request
func queryParamsPolicy(req *policy.Request) (*http.Response, error) {
	// check if any custom query params have been specified
	rawQP := req.Raw().Context().Value(shared.CtxWithQueryParametersKey{})
	if rawQP == nil {
		return req.Next()
	}

	// grab any existing query params
	var allQP url.Values
	if req.Raw().URL.RawQuery != "" {
		var err error
		allQP, err = url.ParseQuery(req.Raw().URL.RawQuery)
		if err != nil {
			return nil, err
		}
	}

	newQP := rawQP.(url.Values)
	if len(allQP) > 0 {
		// merge the existing query params with the ones in the context.
		// any overlapping values will be replaced with the context-provided values.
		for qp := range newQP {
			allQP[qp] = newQP[qp]
		}
	} else {
		allQP = newQP
	}

	req.Raw().URL.RawQuery = allQP.Encode()
	return req.Next()
}
