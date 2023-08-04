//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	credNameAzureDeveloperCLI = "AzureDeveloperCLICredential"
	timeoutAzdRequest         = 10 * time.Second
)

// used by tests to fake invoking the CLI
type azureDeveloperCLITokenProvider func(ctx context.Context, scope []string, tenantID string) ([]byte, error)

// AzureDeveloperCLICredentialOptions contains optional parameters for AzureDeveloperCLICredential.
type AzureDeveloperCLICredentialOptions struct {
	// AdditionallyAllowedTenants specifies tenants for which the credential may acquire tokens, in addition
	// to TenantID. Add the wildcard value "*" to allow the credential to acquire tokens for any tenant the
	// logged in account can access.
	AdditionallyAllowedTenants []string
	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the azd environment, which is the tenant where selected Azure subscription is.
	TenantID string

	tokenProvider azureDeveloperCLITokenProvider
}

// init returns an instance of AzureDeveloperCLICredentialOptions initialized with default values.
func (o *AzureDeveloperCLICredentialOptions) init() {
	if o.tokenProvider == nil {
		o.tokenProvider = defaultAzdTokenProvider
	}
}

// AzureDeveloperCLICredential authenticates as the identity logged in to the Azure Developer CLI.
type AzureDeveloperCLICredential struct {
	s             *syncer
	tokenProvider azureDeveloperCLITokenProvider
}

// NewAzureDeveloperCLICredential constructs an AzureDeveloperCLICredential. Pass nil to accept default options.
func NewAzureDeveloperCLICredential(options *AzureDeveloperCLICredentialOptions) (*AzureDeveloperCLICredential, error) {
	cp := AzureDeveloperCLICredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	c := AzureDeveloperCLICredential{tokenProvider: cp.tokenProvider}
	c.s = newSyncer(
		credNameAzureDeveloperCLI,
		cp.TenantID,
		c.requestToken,
		nil, // this credential doesn't have a silent auth method because the CLI handles caching
		syncerOptions{AdditionallyAllowedTenants: cp.AdditionallyAllowedTenants},
	)
	return &c, nil
}

// GetToken requests a token from the Azure CLI. This credential doesn't cache tokens, so every call invokes the CLI.
// This method is called automatically by Azure SDK clients.
func (c *AzureDeveloperCLICredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if len(opts.Scopes) != 1 {
		return azcore.AccessToken{}, errors.New(credNameAzureDeveloperCLI + ": GetToken() requires exactly one scope")
	}
	return c.s.GetToken(ctx, opts)
}

func (c *AzureDeveloperCLICredential) requestToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	b, err := c.tokenProvider(ctx, opts.Scopes, opts.TenantID)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	at, err := c.createAccessToken(b)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	return at, nil
}

var defaultAzdTokenProvider azureDeveloperCLITokenProvider = func(ctx context.Context, scope []string, tenantID string) ([]byte, error) {
	// set a default timeout for this authentication iff the application hasn't done so already
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, timeoutAzdRequest)
		defer cancel()
	}

	commandLine := "azd auth token -o json"
	if tenantID != "" {
		commandLine += " --tenant-id " + tenantID
	}

	for _, scopeName := range scope {
		commandLine += " --scope " + scopeName
	}

	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		dir := os.Getenv("SYSTEMROOT")
		if dir == "" {
			return nil, newCredentialUnavailableError(credNameAzureDeveloperCLI, "environment variable 'SYSTEMROOT' has no value")
		}
		cliCmd = exec.CommandContext(ctx, "cmd.exe", "/c", commandLine)
		cliCmd.Dir = dir
	} else {
		cliCmd = exec.CommandContext(ctx, "/bin/sh", "-c", commandLine)
		cliCmd.Dir = "/bin"
	}
	cliCmd.Env = os.Environ()
	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		msg := stderr.String()
		var exErr *exec.ExitError
		if errors.As(err, &exErr) && exErr.ExitCode() == 127 || strings.HasPrefix(msg, "'azd' is not recognized") {
			msg = "Azure Developer CLI not found on path"
		}
		if msg == "" {
			msg = err.Error()
		}
		return nil, newCredentialUnavailableError(credNameAzureDeveloperCLI, msg)
	}

	return output, nil
}

func (c *AzureDeveloperCLICredential) createAccessToken(tk []byte) (azcore.AccessToken, error) {
	t := struct {
		AccessToken string `json:"token"`
		ExpiresOn   string `json:"expiresOn"`
	}{}
	err := json.Unmarshal(tk, &t)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	// the Azure CLI's "expiresOn" is local time
	exp, err := time.ParseInLocation("2006-01-02T15:04:05Z", t.ExpiresOn, time.Local)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("Error parsing token expiration time %q: %v", t.ExpiresOn, err)
	}

	converted := azcore.AccessToken{
		Token:     t.AccessToken,
		ExpiresOn: exp.UTC(),
	}
	return converted, nil
}

var _ azcore.TokenCredential = (*AzureDeveloperCLICredential)(nil)
