//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"time"
)

// cliTimeout is the default timeout for authentication attempts via CLI tools
const cliTimeout = 10 * time.Second

// cliTokenProvider is used by tests to fake invoking CLI authentication tools
type cliTokenProvider func(ctx context.Context, scopes []string, tenant string) ([]byte, error)

// unavailableIfInChain returns err or, if the credential was invoked by DefaultAzureCredential, a
// credentialUnavailableError having the same message. This ensures DefaultAzureCredential will try
// the next credential in its chain (another developer credential).
func unavailableIfInChain(err error, inDefaultChain bool) error {
	if err != nil && inDefaultChain {
		var unavailableErr *credentialUnavailableError
		if !errors.As(err, &unavailableErr) {
			err = newCredentialUnavailableError(credNameAzureDeveloperCLI, err.Error())
		}
	}
	return err
}

// validScope is for credentials authenticating via external tools. The authority validates scopes for all other credentials.
func validScope(scope string) bool {
	for _, r := range scope {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '.' || r == '-' || r == '_' || r == '/' || r == ':') {
			return false
		}
	}
	return true
}
