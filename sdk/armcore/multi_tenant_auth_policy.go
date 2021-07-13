// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	bearerTokenPrefix            = "Bearer "
	headerAuthorizationAuxiliary = "x-ms-authorization-auxiliary"
)

type multiTenantPolicy struct {
	bearerPolicy azcore.Policy
}

func newMultiTenantPolicy(cred azcore.TokenCredential, opts authenticationPolicyOptions) azcore.Policy {
	return &multiTenantPolicy{
		bearerPolicy: azcore.NewTokenRefreshPolicy(
			cred,
			&multiTenantProcessor{
				completeAuxTenants: opts.AuxiliaryTenants,
			},
			azcore.AuthenticationPolicyOptions{
				Options: opts.Options,
			},
		),
	}
}

func (p *multiTenantPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	return p.bearerPolicy.Do(req)
}

type authenticationPolicyOptions struct {
	Options          azcore.TokenRequestOptions
	AuxiliaryTenants []string
}

type multiTenantProcessor struct {
	// completeAuxTenants holds the complete original list of auxiliary tenants
	completeAuxTenants []string
	// requestTenants holds only the tenants that need their token to be refreshed
	requestTenants []string
	// hold the list of expires_on times associated with the tenants
	expiresOn map[string]time.Time
}

func (p *multiTenantProcessor) IsZeroOrExpired(eo map[string]time.Time) bool {
	tenants := []string{}
	// the case where the map of expires_on times is empty or nil, get a token for all of the
	// auxiliary tenants
	if len(eo) == 0 {
		p.requestTenants = append(tenants, p.completeAuxTenants...)
		return true
	}
	getToken := false
	// if there is a zero value time in the list of expires_on times or a token has expired, add
	// that tenant to the list for refreshing
	for tenant, t := range eo {
		if t.IsZero() || t.Before(time.Now()) {
			tenants = append(tenants, tenant)
			getToken = true
		}
	}
	p.requestTenants = tenants
	return getToken
}

func (p *multiTenantProcessor) ShouldRefresh(eo map[string]time.Time) bool {
	const window = 2 * time.Minute
	now, getToken, tenants := time.Now(), false, []string{}
	for tenant, t := range eo {
		if t.Add(-window).Before(now) {
			tenants = append(tenants, tenant)
			getToken = true
		}
	}
	p.requestTenants = tenants
	return getToken
}

func (p *multiTenantProcessor) Fetch(ctx context.Context, cred azcore.TokenCredential, opts azcore.TokenRequestOptions) (string, error) {
	auxH := []string{}
	if p.expiresOn == nil {
		p.expiresOn = map[string]time.Time{}
	}
	for _, t := range p.requestTenants {
		opts.TenantID = t
		tk, err := cred.GetToken(ctx, opts)
		if err != nil {
			return "", err
		}
		bearerTk := bearerTokenPrefix + tk.Token
		auxH = append(auxH, bearerTk)
		p.expiresOn[t] = tk.ExpiresOn
	}
	return strings.Join(auxH, ", "), nil
}

func (p *multiTenantProcessor) Header() string {
	return headerAuthorizationAuxiliary
}
