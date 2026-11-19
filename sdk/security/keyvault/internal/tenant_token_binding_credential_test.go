// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authenticationTestError struct {
	message              string
	unavailable, chained bool
	causes               []error
}

func (e *authenticationTestError) Error() string   { return e.message }
func (e *authenticationTestError) Unwrap() []error { return e.causes }
func (e *authenticationTestError) Is(target error) bool {
	return target == errorinfo.ErrAuthenticationFailed ||
		(e.unavailable && target == errorinfo.ErrCredentialUnavailable) ||
		(e.chained && target == errorinfo.ErrCredentialChain)
}

func TestTenantTokenBindingCredential(t *testing.T) {
	want := policy.TokenRequestOptions{
		Scopes: []string{"https://vault.azure.net/.default"}, TenantID: "Tenant",
		Claims: "claims", EnableCAE: true, EnableMTLSPoP: true,
	}
	cert := &tls.Certificate{}
	denied := &authenticationTestError{message: "AADSTS3921996: tenant not eligible"}
	var modes []bool
	cred := credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		modes = append(modes, tro.EnableMTLSPoP)
		require.Equal(t, want.Scopes, tro.Scopes)
		require.Equal(t, want.Claims, tro.Claims)
		require.Equal(t, want.EnableCAE, tro.EnableCAE)
		if tro.EnableMTLSPoP {
			if strings.EqualFold(tro.TenantID, want.TenantID) {
				return azcore.AccessToken{}, denied
			}
			return azcore.AccessToken{Token: "bound", BindingCertificate: cert}, nil
		}
		return azcore.AccessToken{Token: "bearer"}, nil
	})
	c := tenantTokenBindingCredential{credential: cred}
	tk, err := c.GetToken(context.Background(), want)
	require.NoError(t, err)
	require.Equal(t, "bearer", tk.Token)
	require.Nil(t, tk.BindingCertificate)
	tro := want
	tro.TenantID = "TENANT"
	_, err = c.GetToken(context.Background(), tro)
	require.NoError(t, err)
	tro.TenantID = "another tenant"
	tk, err = c.GetToken(context.Background(), tro)
	require.NoError(t, err)
	require.Same(t, cert, tk.BindingCertificate)
	otherClient := tenantTokenBindingCredential{credential: cred}
	_, err = otherClient.GetToken(context.Background(), want)
	require.NoError(t, err)
	require.Equal(t, []bool{true, false, false, true, true, false}, modes)
	require.True(t, want.EnableMTLSPoP, "the caller's options must not be modified")
}

func TestTenantTokenBindingCredential_FailureClassification(t *testing.T) {
	denied := &authenticationTestError{message: "AADSTS3921996"}
	unavailable := &authenticationTestError{message: "not configured", unavailable: true}
	other := &authenticationTestError{message: "native attestation failed"}
	for _, tc := range []struct {
		name             string
		err              error
		fallback, bearer bool
	}{
		{name: "denied", err: denied, fallback: true},
		{name: "ordinary bearer denial", err: denied, bearer: true},
		{name: "wrapped denial", err: fmt.Errorf("wrapped: %w", denied), fallback: true},
		{name: "inner message", err: &authenticationTestError{causes: []error{errors.New("AADSTS3921996: denied")}}, fallback: true},
		{name: "plain text is not authentication failure", err: errors.New("AADSTS3921996")},
		{name: "wrong case", err: &authenticationTestError{message: "aadsts3921996"}},
		{name: "longer code", err: &authenticationTestError{message: "AADSTS39219960"}},
		{name: "prefix", err: &authenticationTestError{message: "XAADSTS3921996"}},
		{name: "underscore", err: &authenticationTestError{message: "AADSTS3921996_other"}},
		{name: "unicode word boundary", err: &authenticationTestError{message: "\u00e9AADSTS3921996"}},
		{name: "unicode join control", err: &authenticationTestError{message: "AADSTS3921996\u200d"}},
		{name: "unrelated terminal", err: other},
		{name: "unavailable", err: unavailable},
		{name: "aggregate", err: errors.Join(denied, other)},
		{name: "wrapped aggregate", err: &authenticationTestError{message: "AADSTS3921996", causes: []error{errors.Join(denied, other)}}},
		{name: "canceled", err: &authenticationTestError{message: "AADSTS3921996", causes: []error{context.Canceled}}},
		{name: "deadline", err: &authenticationTestError{message: "AADSTS3921996", causes: []error{context.DeadlineExceeded}}},
		{name: "chain terminal denial", err: &authenticationTestError{chained: true, causes: []error{unavailable, denied}}, fallback: true},
		{name: "chain mixed terminal", err: &authenticationTestError{message: "AADSTS3921996", chained: true, causes: []error{denied, other}}},
		{name: "chain unavailable denial", err: &authenticationTestError{chained: true, unavailable: true, causes: []error{unavailable, &authenticationTestError{message: "AADSTS3921996", unavailable: true}}}, fallback: true},
		{name: "chain canceled", err: &authenticationTestError{chained: true, causes: []error{denied, &authenticationTestError{causes: []error{context.Canceled}}}}},
		{name: "chain unknown error", err: &authenticationTestError{chained: true, causes: []error{unavailable, errors.New("AADSTS3921996")}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			c := tenantTokenBindingCredential{credential: credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				calls++
				if calls == 1 {
					require.Equal(t, !tc.bearer, tro.EnableMTLSPoP)
					return azcore.AccessToken{}, tc.err
				}
				require.False(t, tro.EnableMTLSPoP)
				return azcore.AccessToken{Token: "bearer"}, nil
			})}
			tk, err := c.GetToken(context.Background(), policy.TokenRequestOptions{EnableMTLSPoP: !tc.bearer})
			if tc.fallback {
				require.NoError(t, err)
				require.Equal(t, "bearer", tk.Token)
				require.Equal(t, 2, calls)
			} else {
				require.ErrorIs(t, err, tc.err)
				require.Equal(t, azcore.AccessToken{}, tk)
				require.Equal(t, 1, calls, "unrelated failures must not trigger bearer acquisition")
			}
		})
	}
}

func TestTenantTokenBindingCredential_BearerFailure(t *testing.T) {
	calls := 0
	bearerErr := errors.New("bearer acquisition failed")
	c := tenantTokenBindingCredential{credential: credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		calls++
		if calls == 1 {
			require.True(t, tro.EnableMTLSPoP)
			return azcore.AccessToken{}, &authenticationTestError{message: "AADSTS3921996"}
		}
		require.False(t, tro.EnableMTLSPoP, "tenant memory must survive a failed bearer acquisition")
		return azcore.AccessToken{}, bearerErr
	})}
	for range 2 {
		_, err := c.GetToken(context.Background(), policy.TokenRequestOptions{EnableMTLSPoP: true})
		require.ErrorIs(t, err, bearerErr)
	}
	require.Equal(t, 3, calls)
}

func TestTenantTokenBindingCredential_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	c := tenantTokenBindingCredential{credential: credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		calls++
		if calls == 1 {
			cancel()
		}
		require.True(t, tro.EnableMTLSPoP, "cancellation must prevent tenant memory and bearer acquisition")
		return azcore.AccessToken{}, &authenticationTestError{message: "AADSTS3921996"}
	})}
	_, err := c.GetToken(ctx, policy.TokenRequestOptions{EnableMTLSPoP: true})
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 1, calls)
	_, err = c.GetToken(ctx, policy.TokenRequestOptions{EnableMTLSPoP: true})
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 1, calls)
	require.True(t, c.effectiveRequestOptions(policy.TokenRequestOptions{EnableMTLSPoP: true}).EnableMTLSPoP)
}

func TestTenantTokenBindingCredential_Concurrent(t *testing.T) {
	c := tenantTokenBindingCredential{credential: credentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		if tro.EnableMTLSPoP {
			return azcore.AccessToken{}, &authenticationTestError{message: "AADSTS3921996"}
		}
		return azcore.AccessToken{Token: "bearer"}, nil
	})}
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Go(func() {
			tk, err := c.GetToken(context.Background(), policy.TokenRequestOptions{TenantID: fmt.Sprint(i % 2), EnableMTLSPoP: true})
			assert.NoError(t, err)
			assert.Equal(t, "bearer", tk.Token)
		})
	}
	wg.Wait()
	for _, tenant := range []string{"0", "1"} {
		require.False(t, c.effectiveRequestOptions(policy.TokenRequestOptions{TenantID: tenant, EnableMTLSPoP: true}).EnableMTLSPoP)
	}
}
