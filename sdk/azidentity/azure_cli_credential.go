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
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// used by tests to fake invoking the CLI
type azureCLITokenProvider func(ctx context.Context, resource string, tenantID string) ([]byte, error)

// AzureCLICredentialOptions contains options used to configure the AzureCLICredential
// All zero-value fields will be initialized with their default values.
type AzureCLICredentialOptions struct {
	tokenProvider azureCLITokenProvider
	// TenantID identifies the tenant the credential should authenticate in.
	// Defaults to the CLI's default tenant, which is typically the home tenant of the user logged in to the CLI.
	TenantID string
}

// init returns an instance of AzureCLICredentialOptions initialized with default values.
func (o *AzureCLICredentialOptions) init() {
	if o.tokenProvider == nil {
		o.tokenProvider = defaultTokenProvider()
	}
}

// AzureCLICredential enables authentication to Azure Active Directory using the Azure CLI command "az account get-access-token".
type AzureCLICredential struct {
	tokenProvider azureCLITokenProvider
	tenantID      string
}

// NewAzureCLICredential constructs a new AzureCLICredential with the details needed to authenticate against Azure Active Directory
// options: configure the management of the requests sent to Azure Active Directory.
func NewAzureCLICredential(options *AzureCLICredentialOptions) (*AzureCLICredential, error) {
	cp := AzureCLICredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	return &AzureCLICredential{
		tokenProvider: cp.tokenProvider,
		tenantID:      cp.TenantID,
	}, nil
}

// GetToken obtains a token from Azure Active Directory, using the Azure CLI command to authenticate.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	if len(opts.Scopes) != 1 {
		return nil, errors.New("this credential requires exactly one scope per token request")
	}
	// CLI expects an AAD v1 resource, not a v2 scope
	scope := strings.TrimSuffix(opts.Scopes[0], defaultSuffix)
	at, err := c.authenticate(ctx, scope)
	if err != nil {
		addGetTokenFailureLogs("Azure CLI Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return at, nil
}

const timeoutCLIRequest = 10 * time.Second

// authenticate creates a client secret authentication request and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// scopes: The scopes for which the token has access
func (c *AzureCLICredential) authenticate(ctx context.Context, resource string) (*azcore.AccessToken, error) {
	output, err := c.tokenProvider(ctx, resource, c.tenantID)
	if err != nil {
		return nil, err
	}

	return c.createAccessToken(output)
}

func defaultTokenProvider() func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
	return func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
		match, err := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, fmt.Errorf(`unexpected scope "%s". Only alphanumeric characters and ".", ";", "-", and "/" are allowed`, resource)
		}

		ctx, cancel := context.WithTimeout(ctx, timeoutCLIRequest)
		defer cancel()

		commandLine := "az account get-access-token -o json --resource " + resource
		if tenantID != "" {
			commandLine += " --tenant " + tenantID
		}
		var cliCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			dir := os.Getenv("SYSTEMROOT")
			if dir == "" {
				return nil, errors.New("environment variable 'SYSTEMROOT' has no value")
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
			if msg == "" {
				// if there's no output in stderr report the error message instead
				msg = err.Error()
			}
			return nil, newCredentialUnavailableError("Azure CLI Credential", msg)
		}

		return output, nil
	}
}

func (c *AzureCLICredential) createAccessToken(tk []byte) (*azcore.AccessToken, error) {
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
		return nil, err
	}

	tokenExpirationDate, err := parseExpirationDate(t.ExpiresOn)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Token Expiration Date %q: %+v", t.ExpiresOn, err)
	}

	converted := &azcore.AccessToken{
		Token:     t.AccessToken,
		ExpiresOn: *tokenExpirationDate,
	}
	return converted, nil
}

// parseExpirationDate parses either a Azure CLI or CloudShell date into a time object
func parseExpirationDate(input string) (*time.Time, error) {
	// CloudShell (and potentially the Azure CLI in future)
	expirationDate, cloudShellErr := time.Parse(time.RFC3339, input)
	if cloudShellErr != nil {
		// Azure CLI (Python) e.g. 2017-08-31 19:48:57.998857 (plus the local timezone)
		const cliFormat = "2006-01-02 15:04:05.999999"
		expirationDate, cliErr := time.ParseInLocation(cliFormat, input, time.Local)
		if cliErr != nil {
			return nil, fmt.Errorf("Error parsing expiration date %q.\n\nCloudShell Error: \n%+v\n\nCLI Error:\n%+v", input, cloudShellErr, cliErr)
		}
		return &expirationDate, nil
	}
	return &expirationDate, nil
}

var _ azcore.TokenCredential = (*AzureCLICredential)(nil)
