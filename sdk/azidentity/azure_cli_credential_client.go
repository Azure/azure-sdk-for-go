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

// execManager that helps mock run Azure CLI command and getting the result from standard output and error streams.
type execManager interface {
	getAzureCLIAuthResults(resource string) (*authResults, error)
}

// execManagerStruct implements the interface execManager.
type execManagerStruct struct {
}

// azureCLICredentialClient provides the client for authenticating with Azure CLI Credential.
type azureCLICredentialClient struct {
}

// newAzureCLICredentialClient creates a new instance of the azureCLICredentialClient.
func newAzureCLICredentialClient() *azureCLICredentialClient {
	return &azureCLICredentialClient{}
}

// authenticate runs Azure CLI command for Azure CLI Credential and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// scopes: The scopes required for the token
// execManagerStruct: The struct implements the interface execManager.
func (c *azureCLICredentialClient) authenticate(ctx context.Context, scopes []string, execManagerStruct execManager) (*azcore.AccessToken, error) {
	// covert the scopes to a resource string
	resource := c.scopeToResource(scopes)

	// Validate the resource to make sure it doesn't include shell-meta characters since it gets sent as a command line argument to Azure CLI
	isResourceMatch, error := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
	if error != nil {
		return nil, error
	}
	if !isResourceMatch {
		return nil, fmt.Errorf(invalidResourceError)
	}

	// Execute Azure CLI command 'az account get-access-token --output json --resource' to return results.
	authResults, err := execManagerStruct.getAzureCLIAuthResults(resource)
	if err != nil {
		return nil, &CredentialUnavailableError{Message: authResults.errOut}
	}

	return c.createAccessToken(authResults.out)
}

// getAzureCLIAuthResults implements the execManager interface on execManagerSturct.
// It gets a token using Azure CLI for local development scenarios.
func (c *execManagerStruct) getAzureCLIAuthResults(resource string) (*authResults, error) {
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

	// Parse expiresOn string to date as required time format
	dateString := value.ExpiresOn
	timeformat := "2006-01-02 15:04:05.999999"
	expiresOnValue, err := time.ParseInLocation(timeformat, dateString, time.Local)
	if err != nil {
		return nil, err
	}

	accessToken.ExpiresOn = expiresOnValue.In(time.UTC)
	accessToken.Token = value.Token

	return accessToken, err
}

func (c *azureCLICredentialClient) scopeToResource(scopes []string) string {
	// Return the 0th scope directly if it doesn't end with suffix
	if !strings.HasSuffix(scopes[0], suffix) {
		return scopes[0]
	}

	// The following code will remove the suffix from any scopes passed into the method
	// since AzureCLICredential expect a resource string instead of a scope string
	scope := scopes[0][0:strings.Index(scopes[0], suffix)]

	return scope
}
