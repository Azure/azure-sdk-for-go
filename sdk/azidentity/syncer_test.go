//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestResolveTenant(t *testing.T) {
	defaultTenant := "default-tenant"
	otherTenant := "other-tenant"
	for _, test := range []struct {
		allowed          []string
		expected, tenant string
		expectError      bool
	}{
		// no alternate tenant specified -> should get default
		{expected: defaultTenant},
		{allowed: []string{""}, expected: defaultTenant},
		{allowed: []string{"*"}, expected: defaultTenant},
		{allowed: []string{otherTenant}, expected: defaultTenant},

		// alternate tenant specified and allowed -> should get that tenant
		{allowed: []string{"*"}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{otherTenant}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{"not-" + otherTenant, otherTenant}, expected: otherTenant, tenant: otherTenant},
		{allowed: []string{"not-" + otherTenant, "*"}, expected: otherTenant, tenant: otherTenant},

		// invalid or not allowed tenant -> should get an error
		{tenant: otherTenant, expectError: true},
		{allowed: []string{""}, tenant: otherTenant, expectError: true},
		{allowed: []string{defaultTenant}, tenant: otherTenant, expectError: true},
		{tenant: badTenantID, expectError: true},
		{allowed: []string{""}, tenant: badTenantID, expectError: true},
		{allowed: []string{"*", badTenantID}, tenant: badTenantID, expectError: true},
		{tenant: "invalid@tenant", expectError: true},
		{tenant: "invalid/tenant", expectError: true},
		{tenant: "invalid(tenant", expectError: true},
		{tenant: "invalid:tenant", expectError: true},
	} {
		t.Run("", func(t *testing.T) {
			s := newSyncer("", defaultTenant, test.allowed, nil, nil)
			tenant, err := s.resolveTenant(test.tenant)
			if err != nil {
				if test.expectError {
					return
				}
				t.Fatal(err)
			} else if test.expectError {
				t.Fatal("expected an error")
			}
			if tenant != test.expected {
				t.Fatalf(`expected "%s", got "%s"`, test.expected, tenant)
			}
		})
	}
}

func TestSyncer(t *testing.T) {
	silentAuths, tokenRequests := 0, 0
	s := newSyncer("", "tenant", nil,
		func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
			tokenRequests++
			return azcore.AccessToken{}, nil
		},
		func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
			var err error
			if tokenRequests == 0 {
				err = errors.New("cache miss")
			}
			silentAuths++
			return azcore.AccessToken{}, err
		},
	)
	goroutines := 50
	wg := sync.WaitGroup{}
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			_, err := s.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}
	if silentAuths != goroutines {
		t.Errorf("expected %d silent auth attempts, got %d", goroutines, silentAuths)
	}
}
