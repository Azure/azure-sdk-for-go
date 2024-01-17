// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type globalEndpointManagerPolicy struct {
	gem *globalEndpointManager
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	shouldRefresh := p.gem.ShouldRefresh()
	if shouldRefresh {
		go func() {
			_ = p.gem.Update(context.Background())
		}()
	}
	return req.Next()
}
