// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// EUAP endpoint for Discovery API (2026-02-01-preview is only available on EUAP)
const EUAPEndpoint = "eastus2euap.management.azure.com"

// ResourceLocation is the default location for EUAP testing
const ResourceLocation = "eastus2euap"

// euapRedirectPolicy redirects requests from management.azure.com to the EUAP endpoint
type euapRedirectPolicy struct{}

func (p *euapRedirectPolicy) Do(req *policy.Request) (*http.Response, error) {
	rawReq := req.Raw()
	host := rawReq.URL.Host

	// Redirect management.azure.com to EUAP endpoint
	if host == "management.azure.com" || strings.HasSuffix(host, ".management.azure.com") {
		rawReq.URL.Host = EUAPEndpoint
		rawReq.Host = EUAPEndpoint
	}

	return req.Next()
}

// GetEUAPClientOptions returns client options configured to redirect to EUAP
func GetEUAPClientOptions() *policy.ClientOptions {
	return &policy.ClientOptions{
		PerCallPolicies: []policy.Policy{&euapRedirectPolicy{}},
	}
}
