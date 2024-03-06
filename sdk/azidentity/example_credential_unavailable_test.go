//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// timeoutWrapper signals ChainedTokenCredential to try another credential when managed identity times out
type timeoutWrapper struct {
	cred    *azidentity.ManagedIdentityCredential
	timeout time.Duration
}

// GetToken implements the azcore.TokenCredential interface
func (w *timeoutWrapper) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var tk azcore.AccessToken
	var err error
	if w.timeout > 0 {
		c, cancel := context.WithTimeout(ctx, w.timeout)
		defer cancel()
		tk, err = w.cred.GetToken(c, opts)
		if ce := c.Err(); errors.Is(ce, context.DeadlineExceeded) {
			// The Context reached its deadline, probably because no managed identity is available.
			// A credential unavailable error signals the chain to try its next credential, if any.
			err = azidentity.NewCredentialUnavailableError("managed identity timed out")
		} else {
			// some managed identity implementation is available, so don't apply the timeout to future calls
			w.timeout = 0
		}
	} else {
		tk, err = w.cred.GetToken(ctx, opts)
	}
	return tk, err
}

// This example demonstrates a small wrapper that sets a deadline for authentication and signals
// [ChainedTokenCredential] to try another credential when managed identity authentication times
// out, as it would for example in a local development environment.
func ExampleNewChainedTokenCredential_managedIdentityTimeout() {
	mic, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	azCLI, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		// TODO: handle error
	}
	creds := []azcore.TokenCredential{
		&timeoutWrapper{mic, time.Second},
		azCLI,
	}
	chain, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		// TODO: handle error
	}
	// TODO: construct a client with the credential chain
	_ = chain
}
