// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/test"
)

func getBuildCodeCoverage(client *azuredevops.Client, projectName string, buildId int) (*test.CodeCoverageSummary, error) {
	request, err := client.CreateRequestMessage(
		ctx,
		"GET",
		fmt.Sprintf("https://dev.azure.com/azure-sdk/%s/_apis/test/codecoverage?buildId=%d&api-version=7.0", projectName, buildId),
		"7.0",
		nil,
		"",
		"",
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := client.SendRequest(request)
	if err != nil {
		return nil, err
	}

	var codeCoverageResult test.CodeCoverageSummary
	err = client.UnmarshalBody(response, &codeCoverageResult)
	if err != nil {
		return nil, err
	}

	return &codeCoverageResult, err
}
