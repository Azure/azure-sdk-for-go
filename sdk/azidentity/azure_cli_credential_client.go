// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	suffix               = "/.default"
	getTokenCommand      = "az account get-access-token -o json "
	resourceArgumentName = "--resource "
	invalidResourceError = "Resource is not in expected format. Only alphanumeric characters, '.', ':', '-', and '/' are allowed"
)

type azureCLIAccessToken struct {
	Token        string `json:"accessToken"`
	ExpiresOn    string `json:"expiresOn"`
	Subscription string `json:"subscription"`
	Tenant       string `json:"tenant"`
	TokenType    string `json:"tokenType"`
}

// execManager that helps mock invoking a Azure CLI command and getting the result from standard output and error streams.
type execManager interface {
	getAzureCLIAccessToken(resource string) ([]byte, string, error)
}

// azureCLICredentialClient provides the client for authenticating with Azure CLI Credential.
type azureCLICredentialClient struct {
}

// execManage invokes Azure CLI command 'account get-access-token' to get a token.
type execManage struct {
}

// newAzureCLICredentialClient creates a new instance of the azureCLICredentialClient.
func newAzureCLICredentialClient() *azureCLICredentialClient {
	return &azureCLICredentialClient{}
}

// authenticate runs Azure CLI command for Azure CLI Credential and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// scopes: The scopes required for the token
func (c *azureCLICredentialClient) authenticate(ctx context.Context, scopes []string, execManage execManager) (*azcore.AccessToken, error) {
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

	// Execute Azure CLI command 'az account get-access-token --output json --resource' to get a token to authenticate
	out, errout, err := execManage.getAzureCLIAccessToken(resource)
	if err != nil {
		return nil, &CredentialUnavailableError{Message: errout}
	}

	return c.createAccessToken(out)
}

func (c *execManage) getAzureCLIAccessToken(resource string) ([]byte, string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	var cmd *exec.Cmd
	var shell string
	var executePara string

	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
		executePara = "/c"
	} else {
		shell = "/bin/sh"
		executePara = "-c"
	}

	// Run Azure CLI command to return resulting Access Token or error
	command := getTokenCommand + resourceArgumentName + resource
	cmd = exec.Command(shell, executePara, command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.Bytes(), stderr.String(), err
}

func (c *azureCLICredentialClient) createAccessToken(output []byte) (*azcore.AccessToken, error) {
	token := azureCLIAccessToken{}
	accessToken := &azcore.AccessToken{}

	err := json.Unmarshal(output, &token)
	if err != nil {
		return nil, err
	}

	// Parse expiresOn string to date as required time format
	dateString := token.ExpiresOn
	timeformat := "2006-01-02 15:04:05.999999"
	expiresOnValue, err := time.ParseInLocation(timeformat, dateString, time.Local)
	if err != nil {
		return nil, err
	}

	accessToken.ExpiresOn = expiresOnValue.In(time.UTC)
	accessToken.Token = token.Token

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
