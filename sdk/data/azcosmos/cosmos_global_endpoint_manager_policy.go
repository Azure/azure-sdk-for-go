// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type globalEndpointManagerPolicy struct {
	gem  *globalEndpointManager
	once sync.Once
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	var err error
	p.once.Do(func() {
		err = p.gem.Update(context.Background(), true)
	})
	if p.gem.ShouldRefresh() {
		go func() {
			_ = p.gem.Update(context.Background(), false)
		}()
	}
	if err != nil {
		return nil, err
	}
	return req.Next()
}
