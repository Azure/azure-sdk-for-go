//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

type authFn func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

// syncer synchronizes authentication calls so that goroutines can share a credential instance
type syncer struct {
	addlTenants      []string
	mu               *sync.Mutex
	reqToken, silent authFn
	name, tenant     string
}

func newSyncer(name, tenant string, additionalTenants []string, reqToken, silentAuth authFn) *syncer {
	return &syncer{
		addlTenants: resolveAdditionalTenants(additionalTenants),
		mu:          &sync.Mutex{},
		name:        name,
		reqToken:    reqToken,
		silent:      silentAuth,
		tenant:      tenant,
	}
}

// GetToken ensures that only one goroutine authenticates at a time
func (s *syncer) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(opts.Scopes) == 0 {
		return azcore.AccessToken{}, errors.New(s.name + ".GetToken() requires at least one scope")
	}
	// we don't resolve the tenant for managed identities because they can acquire tokens only from their home tenants
	if s.name != credNameManagedIdentity {
		tenant, err := s.resolveTenant(opts.TenantID)
		if err != nil {
			return azcore.AccessToken{}, err
		}
		opts.TenantID = tenant
	}
	at, err := s.silent(ctx, opts)
	if err != nil {
		// cache miss; request a new token
		at, err = s.reqToken(ctx, opts)
	}
	if err != nil {
		// Return credentialUnavailableError directly because that type affects the behavior of credential chains.
		// Otherwise, return AuthenticationFailedError.
		var unavailableErr *credentialUnavailableError
		if !errors.As(err, &unavailableErr) {
			res := getResponseFromError(err)
			err = newAuthenticationFailedError(s.name, err.Error(), res, err)
		}
	} else if log.Should(EventAuthentication) {
		scope := strings.Join(opts.Scopes, ", ")
		msg := fmt.Sprintf(`%s.GetToken() acquired a token for scope "%s"\n`, s.name, scope)
		log.Write(EventAuthentication, msg)
	}
	return at, err
}

// resolveTenant returns the correct tenant for a token request given the credential's
// configuration, or an error when the specified tenant isn't allowed by that configuration
func (s *syncer) resolveTenant(requested string) (string, error) {
	if requested == "" || requested == s.tenant {
		return s.tenant, nil
	}
	if s.tenant == "adfs" {
		return "", errors.New("ADFS doesn't support tenants")
	}
	if !validTenantID(requested) {
		return "", errors.New(tenantIDValidationErr)
	}
	for _, t := range s.addlTenants {
		if t == "*" || t == requested {
			return requested, nil
		}
	}
	return "", fmt.Errorf(`%s isn't configured to acquire tokens for tenant %q. To enable acquiring tokens for this tenant add it to the AdditionallyAllowedTenants on the credential options, or add "*" to allow acquiring tokens for any tenant`, s.name, requested)
}

// resolveAdditionalTenants returns a copy of tenants, simplified when tenants contains a wildcard
func resolveAdditionalTenants(tenants []string) []string {
	if len(tenants) == 0 {
		return nil
	}
	for _, t := range tenants {
		// a wildcard makes all other values redundant
		if t == "*" {
			return []string{"*"}
		}
	}
	cp := make([]string, len(tenants))
	copy(cp, tenants)
	return cp
}
