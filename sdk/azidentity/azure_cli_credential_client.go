// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	AzureCLINotInstalled = "Azure CLI not installed"
	AzNotLogIn           = "ERROR: Please run 'az login'"
	WinAzureCLIError     = "'az' is not recognized"
	invalidResourceError = "Resource is not in expected format. Only alphanumeric characters, '.', ':', '-', and '/' are allowed"
)

type azureCliAccessToken struct {
	Token        string `json:"accessToken"`
	ExpiresOn    string `json:"expiresOn"`
	Subscription string `json:"subscription"`
	Tenant       string `json:"tenant"`
	TokenType    string `json:"tokenType"`
}

type ShellClient interface {
	getAzureCliAccessToken(command string) ([]byte, string, error)
}

// AzureCliCredentialClient provides the base for authenticating with Azure Cli Credential.
type AzureCliCredentialClient struct {
}

// NewAzureCliCredentialClient creates a new instance of the AzureCliCredentialClient.
func NewAzureCliCredentialClient() *AzureCliCredentialClient {
	return &AzureCliCredentialClient{}
}

// Authenticate runs a Azure Cli command for a Azure cli Credential and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// scopes: The scopes required for the token
func (c *AzureCliCredentialClient) authenticate(ctx context.Context, scopes []string, credentialClient ShellClient) (*azcore.AccessToken, error) {

	// Validate the resource to make sure it doesn't include shell-meta characters since it gets sent as a command line argument to Azure Cli
	scopeToResource := strings.Join(scopes, " ")
	resource := scopeToResource[0:strings.Index(scopeToResource, ".default")]
	isResourceMatch, error := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)

	if error != nil {
		return nil, error
	}
	if !isResourceMatch {
		return nil, fmt.Errorf(invalidResourceError)
	}

	// Execute Azure Cli command(az account get-access-token --output json --resource) to get a token to authenticate
	out, errout, err := credentialClient.getAzureCliAccessToken(resource)

	// Determining Azure Cli errors
	if err != nil {
		isLoginError := strings.HasPrefix(errout, AzNotLogIn)

		// Is Azure cli installed or not
		isWinError := strings.HasPrefix(errout, WinAzureCLIError)
		isOtherOsError, err := regexp.MatchString("az:(.*)not found", errout)
		if err != nil {
			return nil, err
		} else if isWinError || isOtherOsError {
			return nil, &AuthenticationFailedError{inner: errors.New(AzureCLINotInstalled)}
		}

		// Is user log in or not
		if isLoginError {
			return nil, &AuthenticationFailedError{inner: errors.New(AzNotLogIn)}
		}

		return nil, &AuthenticationFailedError{inner: errors.New(errout)}
	}

	return c.createAccessToken(out)
}

func (c *AzureCliCredentialClient) getAzureCliAccessToken(resource string) ([]byte, string, error) {
	// The command that Azure Cli would be run
	command := "az account get-access-token --output json --resource " + resource

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

	// Run Azure Cli command to return resulting Access Token or error
	cmd = exec.Command(shell, executePara, command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.Bytes(), stderr.String(), err
}

func (c *AzureCliCredentialClient) createAccessToken(output []byte) (*azcore.AccessToken, error) {
	token := azureCliAccessToken{}
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
