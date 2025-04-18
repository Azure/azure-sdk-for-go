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
		// Use the same context, but without the cancellation signal.
		// We DO want to preserve things like context values, but the GEM update needs to complete fully, even if the user cancels the triggering request.
		err = p.gem.Update(context.WithoutCancel(req.Raw().Context()), true)
	})
	if p.gem.ShouldRefresh() {
		go func() {
			// Use the same context, but without the cancellation signal.
			// We DO want to preserve things like context values, but the GEM update needs to complete fully, even if the user cancels the triggering request.
			_ = p.gem.Update(context.WithoutCancel(req.Raw().Context()), false)
		}()
	}
	if err != nil {
		return nil, err
	}
	return req.Next()
}
