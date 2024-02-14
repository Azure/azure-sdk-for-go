// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type sessionPolicy struct {
	sc interface {
		GetSessionToken(resourceAddress string) string
		SetSessionToken(resourceAddress string, containerRid string, sessionToken string)
		ClearSessionToken(resourceAddress string)
	}
}

func (p *sessionPolicy) Do(req *policy.Request) (*http.Response, error) {
	o := pipelineRequestOptions{}
	if !req.OperationValue(&o) || req.Raw().Header.Get(cosmosHeaderSessionToken) != "" {
		return req.Next()
	}

	// Reads: use session token if available
	if !o.isWriteOperation {
		sessionToken := p.sc.GetSessionToken(o.resourceAddress)
		if sessionToken != "" {
			req.Raw().Header.Set(cosmosHeaderSessionToken, sessionToken)
		}
	}

	response, err := req.Next()
	if err != nil || !(response.StatusCode >= 200 && response.StatusCode < 300) {
		// Allow potential retry attempt without session token
		p.sc.ClearSessionToken(o.resourceAddress)
		return response, err
	}

	// Successful Writes: cache session token
	if o.isWriteOperation && response.StatusCode >= 200 && response.StatusCode < 300 {
		sessionToken := response.Header.Get(cosmosHeaderSessionToken)
		containerPath := response.Header.Get(cosmosHeaderAltContentPath)
		containerRid := response.Header.Get(cosmosHeaderContentPath)

		// We currently expect the container path and RID headers to be present
		// (rather than falling back to o.resourceAddress and parsing _self)
		if sessionToken != "" && containerPath != "" && containerRid != "" {
			p.sc.SetSessionToken(containerPath, containerRid, sessionToken)
		}
	}

	return response, nil
}
