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

const credNameAzurePowerShell = "AzurePowerShellCredential"

type azurePowerShellTokenProvider func(ctx context.Context, scopes []string, tenant, subscription string) ([]byte, error)

// AzurePowerShellCredentialOptions contains optional parameters for AzurePowerShellCredential.
type AzurePowerShellCredentialOptions struct {
	// AdditionallyAllowedTenants specifies tenants to which the credential may authenticate, in addition to
	// TenantID. When TenantID is empty, this option has no effect and the credential will authenticate to
	// any requested tenant. Add the wildcard value "*" to allow the credential to authenticate to any tenant.
	AdditionallyAllowedTenants []string

	// Subscription is the name or ID of a subscription. Set this to acquire tokens for an account other
	// than the Azure PowerShell's current account.
	Subscription string

	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the Azure PowerShell's default tenant, which is typically the home tenant of the logged in user.
	TenantID string

	// inDefaultChain is true when the credential is part of DefaultAzureCredential
	inDefaultChain bool
	// tokenProvider is used by tests to fake invoking az
	tokenProvider azurePowerShellTokenProvider
}

// init returns an instance of AzurePowerShellCredentialOptions initialized with default values.
func (o *AzurePowerShellCredentialOptions) init() {
	if o.tokenProvider == nil {
		o.tokenProvider = defaultAzurePowerShellTokenProvider
	}
}

// AzurePowerShellCredential authenticates as the identity logged in to Azure PowerShell.
type AzurePowerShellCredential struct {
	mu   *sync.Mutex
	opts AzurePowerShellCredentialOptions
}

// NewAzurePowerShellCredential constructs an AzurePowerShellCredential. Pass nil to accept default options.
func NewAzurePowerShellCredential(options *AzurePowerShellCredentialOptions) (*AzurePowerShellCredential, error) {
	cp := AzurePowerShellCredentialOptions{}

	if options != nil {
		cp = *options
	}

	if cp.Subscription != "" && !validSubscription(cp.Subscription) {
		return nil, fmt.Errorf(
			"%s: Subscription %q contains invalid characters. If this is the name of a subscription, use its ID instead",
			credNameAzurePowerShell,
			cp.Subscription,
		)
	}

	if cp.TenantID != "" && !validTenantID(cp.TenantID) {
		return nil, errInvalidTenantID
	}

	cp.init()
	cp.AdditionallyAllowedTenants = resolveAdditionalTenants(cp.AdditionallyAllowedTenants)

	return &AzurePowerShellCredential{mu: &sync.Mutex{}, opts: cp}, nil
}

// defaultAzurePowerShellTokenProvider invokes Azure PowerShell to acquire a token. It assumes
// callers have verified that all string arguments are safe to pass to Azure PowerShell.
var defaultAzurePowerShellTokenProvider azurePowerShellTokenProvider = func(ctx context.Context, scopes []string, tenantID, subscription string) ([]byte, error) {

	// pass Azure PowerShell a Microsoft Entra ID v1 resource because we don't know which Azure PowerShell version is installed and older ones don't support v2 scopes
	resource := strings.TrimSuffix(scopes[0], defaultSuffix)

	// set a default timeout for this authentication if the application hasn't done so already
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, cliTimeout)
		defer cancel()
	}

	commandLine := "Get-AzAccessToken -Resource " + resource

	if tenantID != "" {
		commandLine += " -TenantId " + tenantID
	}

	if subscription != "" {
		// subscription needs quotes because it may contain spaces
		commandLine += ` -Subscription '` + subscription + `'`
	}

	// convert PSAccessToken to JSON so it can be unmarshaled later
	commandLine += " | ConvertTo-Json"

	// Wrap the command in double quotes for PowerShell -Command
	powershellCommandArg := "\"" + commandLine + "\""

	var powershellCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		dir := os.Getenv("SYSTEMROOT")
		if dir == "" {
			return nil, newCredentialUnavailableError(credNameAzurePowerShell, "environment variable 'SYSTEMROOT' has no value")
		}

		// Prefer pwsh.exe (PowerShell Core), fallback to powershell.exe (Windows PowerShell)
		pwshPath, err := exec.LookPath("pwsh.exe")
		if err == nil {
			powershellCmd = exec.CommandContext(ctx, pwshPath, "-Command", powershellCommandArg)
		} else {
			powershellPath, err := exec.LookPath("powershell.exe")
			if err != nil {
				return nil, newCredentialUnavailableError(credNameAzurePowerShell, "Neither pwsh.exe nor powershell.exe found in PATH")
			}
			powershellCmd = exec.CommandContext(ctx, powershellPath, "-Command", powershellCommandArg)
		}
		powershellCmd.Dir = dir
	} else {
		// On Unix, only support PowerShell Core (pwsh)
		pwshPath, err := exec.LookPath("pwsh")
		if err != nil {
			return nil, newCredentialUnavailableError(credNameAzurePowerShell, "pwsh not found in PATH; PowerShell Core is required on Unix platforms")
		}
		powershellCmd = exec.CommandContext(ctx, pwshPath, "-Command", powershellCommandArg)
		powershellCmd.Dir = "/bin"
	}

	powershellCmd.Env = os.Environ()
	var stderr bytes.Buffer
	powershellCmd.Stderr = &stderr
	powershellCmd.WaitDelay = 100 * time.Millisecond

	stdout, err := powershellCmd.Output()

	if errors.Is(err, exec.ErrWaitDelay) && len(stdout) > 0 {
		// The child process wrote to stdout and exited without closing it.
		// Swallow this error and return stdout because it may contain a token.
		return stdout, nil
	}

	if err != nil {
		msg := stderr.String()
		var exErr *exec.ExitError

		if errors.As(err, &exErr) && exErr.ExitCode() == 127 || strings.HasPrefix(msg, "'Get-AzAccessToken' is not recognized") {
			msg = "Get-AzAccessToken command not found"
		}

		if msg == "" {
			msg = err.Error()
		}

		return nil, newCredentialUnavailableError(credNameAzurePowerShell, msg)
	}

	return stdout, nil
}

// GetToken requests a token from Azure PowerShell. This credential doesn't cache tokens, so every call invokes Azure PowerShell.
// This method is called automatically by Azure SDK clients.
func (c *AzurePowerShellCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	at := azcore.AccessToken{}
	if len(opts.Scopes) != 1 {
		return at, errors.New(credNameAzurePowerShell + ": GetToken() requires exactly one scope")
	}
	if !validScope(opts.Scopes[0]) {
		return at, fmt.Errorf("%s.GetToken(): invalid scope %q", credNameAzurePowerShell, opts.Scopes[0])
	}
	tenant, err := resolveTenant(c.opts.TenantID, opts.TenantID, credNameAzurePowerShell, c.opts.AdditionallyAllowedTenants)
	if err != nil {
		return at, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := c.opts.tokenProvider(ctx, opts.Scopes, tenant, c.opts.Subscription)
	if err == nil {
		at, err = c.createAccessToken(b)
	}
	if err != nil {
		err = unavailableIfInChain(err, c.opts.inDefaultChain)
		return at, err
	}
	msg := fmt.Sprintf("%s.GetToken() acquired a token for scope %q", credNameAzurePowerShell, strings.Join(opts.Scopes, ", "))
	log.Write(EventAuthentication, msg)
	return at, nil
}

func (c *AzurePowerShellCredential) createAccessToken(tk []byte) (azcore.AccessToken, error) {
	t := struct {
		Token     string `json:"Token"`
		ExpiresOn string `json:"ExpiresOn"`
	}{}
	err := json.Unmarshal(tk, &t)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	exp, err := time.Parse(time.RFC3339, t.ExpiresOn)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("%s: error parsing token expiration time %q: %v", credNameAzurePowerShell, t.ExpiresOn, err)
	}

	converted := azcore.AccessToken{
		Token:     t.Token,
		ExpiresOn: exp.UTC(),
	}
	return converted, nil
}

var _ azcore.TokenCredential = (*AzurePowerShellCredential)(nil)
