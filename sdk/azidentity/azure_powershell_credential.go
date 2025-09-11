// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const credNameAzurePowerShell = "AzurePowerShellCredential"

// AzurePowerShellCredentialOptions contains optional parameters for AzurePowerShellCredential.
type AzurePowerShellCredentialOptions struct {
	// AdditionallyAllowedTenants specifies tenants to which the credential may authenticate, in addition to
	// TenantID. When TenantID is empty, this option has no effect and the credential will authenticate to
	// any requested tenant. Add the wildcard value "*" to allow the credential to authenticate to any tenant.
	AdditionallyAllowedTenants []string

	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the Azure PowerShell's default tenant, which is typically the home tenant of the logged in user.
	TenantID string

	// inDefaultChain is true when the credential is part of DefaultAzureCredential
	inDefaultChain bool

	// exec is used by tests to fake invoking Azure PowerShell
	exec executor
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

	if cp.TenantID != "" && !validTenantID(cp.TenantID) {
		return nil, errInvalidTenantID
	}

	if cp.exec == nil {
		cp.exec = shellExec
	}

	cp.AdditionallyAllowedTenants = resolveAdditionalTenants(cp.AdditionallyAllowedTenants)

	return &AzurePowerShellCredential{mu: &sync.Mutex{}, opts: cp}, nil
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

	// pass the CLI a Microsoft Entra ID v1 resource because we don't know which CLI version is installed and older ones don't support v2 scopes
	resource := strings.TrimSuffix(opts.Scopes[0], defaultSuffix)

	tenantArg := ""
	if tenant != "" {
		tenantArg = fmt.Sprintf(" -TenantId '%s'", tenant)
	}

	if opts.Claims != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(opts.Claims))
		return at, fmt.Errorf(
			"%s.GetToken(): Azure PowerShell requires multifactor authentication or additional claims. Run this command then retry the operation: Connect-AzAccount%s -ClaimsChallenge '%s'",
			credNameAzurePowerShell,
			tenantArg,
			encoded,
		)
	}

	// Inline script to handle Get-AzAccessToken differences between Az.Accounts versions with SecureString handling and minimum version requirement
	script := fmt.Sprintf(`
$ErrorActionPreference = 'Stop'
[version]$minimumVersion = '2.2.0'

$mod = Import-Module Az.Accounts -MinimumVersion $minimumVersion -PassThru -ErrorAction SilentlyContinue

if (!$mod) {
    Write-Error '%s'
	exit 1
}

# For Az.Accounts 2.17.0+ but below 5.0.0, explicitly request secure string
if ($mod.Version -ge [version]'2.17.0' -and $mod.Version -lt [version]'5.0.0') {
    $params['AsSecureString'] = $true
}

$token = Get-AzAccessToken -ResourceUrl '%s'%s

$customToken = New-Object -TypeName psobject

# If the token is a SecureString, convert to plain text using recommended pattern
if ($token.Token -is [System.Security.SecureString]) {
    $ssPtr = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($token.Token)
    try {
        $plainToken = [System.Runtime.InteropServices.Marshal]::PtrToStringBSTR($ssPtr)
    } finally {
        [System.Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ssPtr)
    }
    $customToken | Add-Member -MemberType NoteProperty -Name Token -Value $plainToken
} else {
    $customToken | Add-Member -MemberType NoteProperty -Name Token -Value $token.Token
}
$customToken | Add-Member -MemberType NoteProperty -Name ExpiresOn -Value $token.ExpiresOn.UtcDateTime.Ticks

$jsonToken = $customToken | ConvertTo-Json
return $jsonToken
`, azurePowerShellNoAzAccountModule, resource, tenantArg)

	// Windows: prefer pwsh.exe (PowerShell Core), fallback to powershell.exe (Windows PowerShell)
	// Unix: only support pwsh (PowerShell Core)
	var powershellExecutable string
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("pwsh.exe")
		if err == nil {
			powershellExecutable = "pwsh.exe"
		} else {
			powershellExecutable = "powershell.exe"
		}
	} else {
		powershellExecutable = "pwsh"
	}

	command := fmt.Sprintf("%s -NoProfile -NonInteractive -EncodedCommand %s", powershellExecutable, base64EncodeUTF16LE(script))

	c.mu.Lock()
	defer c.mu.Unlock()

	b, err := c.opts.exec(ctx, credNameAzurePowerShell, command)
	if err == nil {
		at, err = c.createAccessToken(b)
	}

	if err != nil {
		err = unavailableIfInDAC(err, c.opts.inDefaultChain)
		return at, err
	}

	msg := fmt.Sprintf("%s.GetToken() acquired a token for scope %q", credNameAzurePowerShell, strings.Join(opts.Scopes, ", "))
	log.Write(EventAuthentication, msg)
	return at, nil
}

func (c *AzurePowerShellCredential) createAccessToken(tk []byte) (azcore.AccessToken, error) {
	t := struct {
		Token     string `json:"Token"`
		ExpiresOn int64  `json:"ExpiresOn"`
	}{}
	err := json.Unmarshal(tk, &t)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	exp := ticksToUnixTime(t.ExpiresOn)

	converted := azcore.AccessToken{
		Token:     t.Token,
		ExpiresOn: exp.UTC(),
	}
	return converted, nil
}

var _ azcore.TokenCredential = (*AzurePowerShellCredential)(nil)
