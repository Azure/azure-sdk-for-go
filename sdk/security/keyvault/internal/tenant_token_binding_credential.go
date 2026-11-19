// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Go regexp's \b is ASCII-only; the C# classifier uses Unicode word boundaries.
var tenantDeniedPattern = regexp.MustCompile(`(^|[^\pL\p{Mn}\p{Nd}\p{Pc}\x{200C}\x{200D}])AADSTS3921996($|[^\pL\p{Mn}\p{Nd}\p{Pc}\x{200C}\x{200D}])`)

type tenantTokenBindingCredential struct {
	credential        azcore.TokenCredential
	bearerOnlyTenants sync.Map
}

func (c *tenantTokenBindingCredential) effectiveRequestOptions(tro policy.TokenRequestOptions) policy.TokenRequestOptions {
	if _, ok := c.bearerOnlyTenants.Load(strings.ToLower(tro.TenantID)); ok {
		tro.EnableMTLSPoP = false
	}
	return tro
}

func (c *tenantTokenBindingCredential) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if err := ctx.Err(); err != nil {
		return azcore.AccessToken{}, err
	}
	tro = c.effectiveRequestOptions(tro)
	token, err := c.credential.GetToken(ctx, tro)
	if err == nil || !tro.EnableMTLSPoP || classifyAuthenticationFailure(err) != failureTenantDenied {
		return token, err
	}
	if err := ctx.Err(); err != nil {
		return azcore.AccessToken{}, err
	}
	if _, loaded := c.bearerOnlyTenants.LoadOrStore(strings.ToLower(tro.TenantID), struct{}{}); !loaded {
		log.Write("Authentication", "Key Vault: the token service denied attested token issuance for this tenant context. Requesting bearer authentication until the client is recreated.")
	}
	tro.EnableMTLSPoP = false
	return c.credential.GetToken(ctx, tro)
}

type authenticationFailure int

const (
	failureOther authenticationFailure = iota
	failureUnavailable
	failureTenantDenied
)

func classifyAuthenticationFailure(err error) authenticationFailure {
	// Accept ordinary Go wrappers, but not unclassified aggregate messages.
	for err != nil {
		if isCancellation(err) {
			return failureOther
		}
		if _, ok := err.(errorinfo.AuthenticationFailure); ok {
			break
		}
		err = errors.Unwrap(err)
	}
	if err == nil {
		return failureOther
	}
	denied, onlyUnavailable := false, true
	for err != nil {
		// errors.Is can't see a credential error's causes, so check each one visited
		if isCancellation(err) {
			return failureOther
		}
		auth, ok := err.(errorinfo.AuthenticationFailure)
		if !ok {
			if _, aggregate := err.(interface{ Unwrap() []error }); aggregate {
				return failureOther
			}
			denied = denied || tenantDeniedPattern.MatchString(err.Error())
			err = errors.Unwrap(err)
			continue
		}
		info := auth.AuthenticationFailure()
		if info.Chained {
			chainDenied := false
			for _, cause := range info.Causes {
				switch classifyAuthenticationFailure(cause) {
				case failureTenantDenied:
					chainDenied = true
				case failureUnavailable:
				default:
					return failureOther
				}
			}
			if chainDenied {
				return failureTenantDenied
			}
			if onlyUnavailable && info.Unavailable {
				return failureUnavailable
			}
			return failureOther
		}
		denied = denied || tenantDeniedPattern.MatchString(err.Error())
		onlyUnavailable = onlyUnavailable && info.Unavailable
		switch len(info.Causes) {
		case 0:
			err = nil
		case 1:
			err = info.Causes[0]
		default:
			return failureOther
		}
	}
	if denied {
		return failureTenantDenied
	}
	if onlyUnavailable {
		return failureUnavailable
	}
	return failureOther
}

func isCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
