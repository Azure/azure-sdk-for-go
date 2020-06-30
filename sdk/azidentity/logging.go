// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// LogCredential is the log classification that can be used for logging Azure Identity related information
const LogCredential azcore.LogClassification = "credential"

// log environment variables that can be used for credential types
func logEnvVars() {
	if !azcore.Log().Should(LogCredential) {
		return
	}
	// Log available environment variables
	envVars := []string{}
	if envCheck := os.Getenv("AZURE_TENANT_ID"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_TENANT_ID")
	}
	if envCheck := os.Getenv("AZURE_CLIENT_ID"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLIENT_ID")
	}
	if envCheck := os.Getenv("AZURE_CLIENT_SECRET"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLIENT_SECRET")
	}
	if envCheck := os.Getenv("AZURE_AUTHORITY_HOST"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_AUTHORITY_HOST")
	}
	if envCheck := os.Getenv("AZURE_CLI_PATH"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLI_PATH")
	}
	if len(envVars) > 0 {
		azcore.Log().Write(LogCredential, fmt.Sprintf("Azure Identity => Found the following environment variables: %s", strings.Join(envVars, ", ")))
	}
}

func logGetTokenSuccess(cred azcore.TokenCredential, opts azcore.TokenRequestOptions) string {
	msg := fmt.Sprintf("Azure Identity => GetToken() result for %T: SUCCESS\n", cred)
	msg += fmt.Sprintf("Azure Identity => Scopes: [%s]", strings.Join(opts.Scopes, ", "))
	return msg
}

func logGetTokenFailure(credName string) string {
	return fmt.Sprintf("Azure Identity => ERROR in GetToken() call for %s. Please check the log for the error.", credName)
}

func logCredentialError(credName string, err error) string {
	return fmt.Sprintf("Azure Identity => ERROR in %s: %s", credName, err.Error())
}

func logMSIEnv(msi msiType) string {
	switch msi {
	case 1:
		return "Azure Identity => Managed Identity environment: IMDS"
	case 2:
		return "Azure Identity => Managed Identity environment: MSI_ENDPOINT"
	case 3:
		return "Azure Identity => Managed Identity environment: MSI_ENDPOINT"
	case 4:
		return "Azure Identity => Managed Identity environment: Unavailable"
	default:
		return "Azure Identity => Managed Identity environment: Unknown"
	}
}

func addGetTokenFailureLogs(credName string, err error) {
	azcore.Log().Write(azcore.LogError, logCredentialError(credName, err))
	azcore.Log().Write(azcore.LogError, logGetTokenFailure(credName))
}
