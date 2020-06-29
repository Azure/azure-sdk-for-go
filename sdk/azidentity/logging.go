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
	log := azcore.Log()
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
		log.Write(LogCredential, fmt.Sprintf("Azure Identity => Found the following environment variables: %s", strings.Join(envVars, ", ")))
	}
}
