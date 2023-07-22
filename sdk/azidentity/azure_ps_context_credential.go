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
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	credNamePSAzureContext  = "PSAzureContextCredential"
)

// used by tests to fake invoking the Powershell Context
type azurePSContextTokenProvider func(ctx context.Context, resource string, tenantID string) ([]byte, error)

// PSAzureContextCredentialOptions contains optional parameters for PSAzureContextCredentialOptions
type PSAzureContextCredentialOptions struct {
	// AdditionallyAllowedTenants specifies tenants for which the credential may acquire tokens, in addition
	// to TenantID. Add the wildcard value "*" to allow the credential to acquire tokens for any tenant the
	// logged in account can access.
	AdditionallyAllowedTenants []string
	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the Azure Context default tenant, which is typically the home tenant of the logged in user.
	TenantID string

	tokenProvider azurePSContextTokenProvider
}

// init returns an instance of PSAzureContextCredentialOptions initialized with default values.
func (o *PSAzureContextCredentialOptions) init() {
	if o.tokenProvider == nil {
		o.tokenProvider = defaultPSTokenProvider
	}
}

// PSAzureContextCredential authenticates as the identity logged
// in to Powershell via Azure Context
type PSAzureContextCredential struct {
	s             *syncer
	tokenProvider azurePSContextTokenProvider
}

// NewPSAzureContextCredential constructs an PSAzureContextCredential. Pass nil to accept default options.
func NewPSAzureContextCredential(options *PSAzureContextCredentialOptions) (*PSAzureContextCredential, error) {
	cp := PSAzureContextCredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	c := PSAzureContextCredential{tokenProvider: cp.tokenProvider}
	c.s = newSyncer(
		credNamePSAzureContext,
		cp.TenantID,
		c.requestToken,
		nil, // this credential doesn't have a silent auth method because the CLI handles caching
		syncerOptions{AdditionallyAllowedTenants: cp.AdditionallyAllowedTenants},
	)
	return &c, nil
}

// GetToken requests a token from Powershell. This credential doesn't cache tokens
// so every call invokes PS again. 
func (c *PSAzureContextCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if len(opts.Scopes) != 1 {
		return azcore.AccessToken{}, errors.New(credNameAzureCLI + ": GetToken() requires exactly one scope")
	}
	// PS expects an AAD v1 resource, not a v2 scope
	opts.Scopes = []string{strings.TrimSuffix(opts.Scopes[0], defaultSuffix)}
	return c.s.GetToken(ctx, opts)
}

func (c *PSAzureContextCredential) requestToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	b, err := c.tokenProvider(ctx, opts.Scopes[0], opts.TenantID)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	at, err := c.createAccessToken(b)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	return at, nil
}

var defaultPSTokenProvider azurePSContextTokenProvider = func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
	match, err := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf(`%s: unexpected scope "%s". Only alphanumeric characters and ".", ";", "-", and "/" are allowed`, credNamePSAzureContext, resource)
	}

	// set a default timeout for this authentication iff the application hasn't done so already
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, timeoutCLIRequest)
		defer cancel()
	}

	commandLine := "Get-AzAccessToken -ResourceTypeName " + resource
	if tenantID != "" {
		commandLine += " -TenantId " + tenantID
	}
	cliCmd := exec.CommandContext(ctx, "powershell", commandLine)
	cliCmd.Env = os.Environ()
	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		msg := stderr.String()
		var exErr *exec.ExitError
		if errors.As(err, &exErr) && exErr.ExitCode() == 127 || strings.HasPrefix(msg, "'az' is not recognized") {
			msg = "Azure CLI not found on path"
		}
		if msg == "" {
			msg = err.Error()
		}
		return nil, newCredentialUnavailableError(credNamePSAzureContext, msg)
	}

	return output, nil
}

func (c *PSAzureContextCredential) createAccessToken(tk []byte) (azcore.AccessToken, error) {
	t := struct {
		Token            string `json:"Token"`
		ExpiresOn        string `json:"ExpiresOn"`
		Type             string `json:"Type"`
		TenantId         string `json:"TenantId"`
		UserID           string `json:"UserId"`
	}{}
	err := json.Unmarshal(tk, &t)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	// the Azure CLI's "expiresOn" is local time
	exp, err := time.ParseInLocation("2006-01-02 15:04:05.999999", t.ExpiresOn, time.Local)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("Error parsing token expiration time %q: %v", t.ExpiresOn, err)
	}

	converted := azcore.AccessToken{
		Token:     t.Token,
		ExpiresOn: exp.UTC(),
	}
	return converted, nil
}

var _ azcore.TokenCredential = (*PSAzureContextCredential)(nil)
