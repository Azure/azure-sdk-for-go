// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type globalEndpointManagerPolicy struct {
	gem  *globalEndpointManager
	once sync.Once
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.once.Do(func() {
		fmt.Println("initializing")
		p.gem.Update(context.Background())
	})
	if p.gem.ShouldRefresh() {
		fmt.Println("refreshing")
		go func() {
			_ = p.gem.Update(context.Background())
		}()
	}
	return req.Next()
}
