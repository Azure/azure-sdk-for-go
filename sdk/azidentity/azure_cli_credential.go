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
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const credNameAzureCLI = "AzureCLICredential"

// AzureCLICredentialOptions contains optional parameters for AzureCLICredential.
type AzureCLICredentialOptions struct {
	// AdditionallyAllowedTenants specifies tenants for which the credential may acquire tokens, in addition
	// to TenantID. Add the wildcard value "*" to allow the credential to acquire tokens for any tenant the
	// logged in account can access.
	AdditionallyAllowedTenants []string

	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the CLI's default tenant, which is typically the home tenant of the logged in user.
	TenantID string

	// inDefaultChain is true when the credential is part of DefaultAzureCredential
	inDefaultChain bool
	// tokenProvider is used by tests to fake invoking az
	tokenProvider cliTokenProvider
}

// init returns an instance of AzureCLICredentialOptions initialized with default values.
func (o *AzureCLICredentialOptions) init() {
	if o.tokenProvider == nil {
		o.tokenProvider = defaultAzTokenProvider
	}
}

// AzureCLICredential authenticates as the identity logged in to the Azure CLI.
type AzureCLICredential struct {
	mu   *sync.Mutex
	opts AzureCLICredentialOptions
}

// NewAzureCLICredential constructs an AzureCLICredential. Pass nil to accept default options.
func NewAzureCLICredential(options *AzureCLICredentialOptions) (*AzureCLICredential, error) {
	cp := AzureCLICredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	cp.AdditionallyAllowedTenants = resolveAdditionalTenants(cp.AdditionallyAllowedTenants)
	return &AzureCLICredential{mu: &sync.Mutex{}, opts: cp}, nil
}

// GetToken requests a token from the Azure CLI. This credential doesn't cache tokens, so every call invokes the CLI.
// This method is called automatically by Azure SDK clients.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	at := azcore.AccessToken{}
	if len(opts.Scopes) != 1 {
		return at, errors.New(credNameAzureCLI + ": GetToken() requires exactly one scope")
	}
	tenant, err := resolveTenant(c.opts.TenantID, opts.TenantID, credNameAzureCLI, c.opts.AdditionallyAllowedTenants)
	if err != nil {
		return at, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := c.opts.tokenProvider(ctx, opts.Scopes, tenant)
	if err == nil {
		at, err = c.createAccessToken(b)
	}
	if err != nil {
		err = unavailableIfInChain(err, c.opts.inDefaultChain)
		return at, err
	}
	msg := fmt.Sprintf("%s.GetToken() acquired a token for scope %q", credNameAzureCLI, strings.Join(opts.Scopes, ", "))
	log.Write(EventAuthentication, msg)
	return at, nil
}

var defaultAzTokenProvider cliTokenProvider = func(ctx context.Context, scopes []string, tenantID string) ([]byte, error) {
	if !validScope(scopes[0]) {
		return nil, fmt.Errorf("%s.GetToken(): invalid scope %q", credNameAzureCLI, scopes[0])
	}
	// pass the CLI a Microsoft Entra ID v1 resource because we don't know which CLI version is installed and older ones don't support v2 scopes
	resource := strings.TrimSuffix(scopes[0], defaultSuffix)
	// set a default timeout for this authentication iff the application hasn't done so already
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, cliTimeout)
		defer cancel()
	}
	commandLine := "az account get-access-token -o json --resource " + resource
	if tenantID != "" {
		commandLine += " --tenant " + tenantID
	}
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		dir := os.Getenv("SYSTEMROOT")
		if dir == "" {
			return nil, newCredentialUnavailableError(credNameAzureCLI, "environment variable 'SYSTEMROOT' has no value")
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
		if errors.As(err, &exErr) && exErr.ExitCode() == 127 || strings.HasPrefix(msg, "'az' is not recognized") {
			msg = "Azure CLI not found on path"
		}
		if msg == "" {
			msg = err.Error()
		}
		return nil, newCredentialUnavailableError(credNameAzureCLI, msg)
	}

	return output, nil
}

func (c *AzureCLICredential) createAccessToken(tk []byte) (azcore.AccessToken, error) {
	t := struct {
		AccessToken      string `json:"accessToken"`
		Authority        string `json:"_authority"`
		ClientID         string `json:"_clientId"`
		ExpiresOn        string `json:"expiresOn"`
		IdentityProvider string `json:"identityProvider"`
		IsMRRT           bool   `json:"isMRRT"`
		RefreshToken     string `json:"refreshToken"`
		Resource         string `json:"resource"`
		TokenType        string `json:"tokenType"`
		UserID           string `json:"userId"`
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
		Token:     t.AccessToken,
		ExpiresOn: exp.UTC(),
	}
	return converted, nil
}

var _ azcore.TokenCredential = (*AzureCLICredential)(nil)
