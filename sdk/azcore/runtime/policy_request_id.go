//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

const requestIdHeader = "x-ms-client-request-id"

type requestIdPolicy struct {
	o policy.ClientRequestOptions
}

func NewRequestIdPolicy(o *policy.ClientRequestOptions) *requestIdPolicy {
	if o == nil {
		o = &policy.ClientRequestOptions{AutoRequestId: true}
	}

	return &requestIdPolicy{o: *o}
}

func (r *requestIdPolicy) Do(req *policy.Request) (*http.Response, error) {
	if req.Raw().Header.Get(requestIdHeader) == "" {
		if r.o.RequestId != "" {
			req.Raw().Header.Set(requestIdHeader, r.o.RequestId)
		} else if r.o.AutoRequestId {
			id, err := uuid.New()
			if err != nil {
				return nil, err
			}
			req.Raw().Header.Set(requestIdHeader, id.String())
		}
	}

	return req.Next()
}
