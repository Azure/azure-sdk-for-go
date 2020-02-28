// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	suffix               = "/.default"
	cmd                  = "cmd.exe"
	bash                 = "/bin/sh"
	windir               = "windir"
	getTokenCommand      = "az account get-access-token -o json"
	resourceArgumentName = "--resource"
	invalidResourceError = "Resource is not in expected format. Only alphanumeric characters, '.', ':', '-', and '/' are allowed"
)

type authResults struct {
	out    []byte
	errOut string
}

// azureCLIAccessTokenProvider provides an AccessToken, either by Azure CLI or by mocking.
type azureCLIAccessTokenProvider interface {
	getAzureCLIAuthResults(resource string) (*authResults, error)
}

// azureCLIAccessTokenProviderStruct implements the interface azureCLIAccessTokenProvider, to run Azure CLI command.
type azureCLIAccessTokenProviderStruct struct {
}

// azureCLICredentialClient provides the client for authenticating with Azure CLI Credential.
type azureCLICredentialClient struct {
	azAccessTokenProvider azureCLIAccessTokenProvider
}

// newAzureCLICredentialClient creates a new instance of the azureCLICredentialClient.
func newAzureCLICredentialClient(azAccessTokenProvider azureCLIAccessTokenProvider) *azureCLICredentialClient {
	if azAccessTokenProvider == nil {
		azAccessTokenProvider = &azureCLIAccessTokenProviderStruct{}
	}

	return &azureCLICredentialClient{azAccessTokenProvider: azAccessTokenProvider}
}

// authenticate runs Azure CLI command for Azure CLI Credential and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// scopes: The scopes required for the token
func (c *azureCLICredentialClient) authenticate(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	// convert the scopes to a resource string
	resource := c.scopeToResource(scopes)

	// Validate the resource to make sure it doesn't include shell-meta characters.
	isResourceMatch, error := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
	if error != nil {
		return nil, error
	}
	if !isResourceMatch {
		return nil, fmt.Errorf(invalidResourceError)
	}

	authResults, err := c.azAccessTokenProvider.getAzureCLIAuthResults(resource)
	if err != nil {
		return nil, &CredentialUnavailableError{Message: authResults.errOut}
	}

	return c.createAccessToken(authResults.out)
}

// getAzureCLIAuthResults implements the azureCLIAccessTokenProvider interface on azureCLIAccessTokenProviderStruct.
// Execute Azure CLI command 'az account get-access-token --output json --resource' to return results.
func (c *azureCLIAccessTokenProviderStruct) getAzureCLIAuthResults(resource string) (*authResults, error) {
	// Developer can set the path what the install path for Azure CLI is.
	azureCLIPath := "AzureCLIPath"

	// The default install path are used to find Azure CLI for windows. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
	azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	azureCLIDefaultPath := "/bin:/sbin:/usr/bin:/usr/local/bin"

	results := &authResults{}
	var stderr bytes.Buffer
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv(windir)))
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
		cmd.Args = append(cmd.Args, "/c")
	} else {
		cmd = exec.Command(bash)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
		cmd.Args = append(cmd.Args, "-c")
	}

	cmd.Args = append(cmd.Args, getTokenCommand, resourceArgumentName, resource)
	cmd.Stderr = &stderr
	output, err := cmd.Output()

	results.out = output
	results.errOut = stderr.String()

	return results, err
}

func (c *azureCLICredentialClient) createAccessToken(output []byte) (*azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string `json:"accessToken"`
		ExpiresOn    string `json:"expiresOn"`
		Subscription string `json:"subscription"`
		Tenant       string `json:"tenant"`
		TokenType    string `json:"tokenType"`
	}{}
	accessToken := &azcore.AccessToken{}
	err := json.Unmarshal(output, &value)
	if err != nil {
		return nil, err
	}

	dateString := value.ExpiresOn
	timeformat := "2006-01-02 15:04:05.999999"

	// The expiresOnValue return from  the Azure CLI is local time. So, get local time first, then parse to UTC.
	expiresOnValue, err := time.ParseInLocation(timeformat, dateString, time.Local)
	if err != nil {
		return nil, err
	}

	accessToken.ExpiresOn = expiresOnValue.In(time.UTC)
	accessToken.Token = value.Token

	return accessToken, err
}

func (c *azureCLICredentialClient) scopeToResource(scopes []string) string {
	// Return the first scope directly if it doesn't end with suffix
	if !strings.HasSuffix(scopes[0], suffix) {
		return scopes[0]
	}

	// Remove suffix from first scope since Azure CLI command parameter "--resource" don't need suffix.
	scope := scopes[0][0:strings.Index(scopes[0], suffix)]

	return scope
}
