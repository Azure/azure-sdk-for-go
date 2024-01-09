// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type globalEndpointManagerPolicy struct {
	gem *globalEndpointManager
}

func (p *globalEndpointManagerPolicy) Do(req *policy.Request) (*http.Response, error) {
	shouldRefresh := p.gem.ShouldRefresh()
	if shouldRefresh {
		fmt.Println("should refresh true go routine")
		go p.gem.Update(context.Background())
	}
	fmt.Println("policy done")
	return req.Next()
}
