// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelines"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/test"
)

// Since azuredevops module has some issue of fetching code coverage (https://github.com/microsoft/azure-devops-go-api/issues/124), we need to do that by ourselves.
func getBuildCodeCoverage(client *azuredevops.Client, projectName string, buildId int) (int, int, error) {
	request, err := client.CreateRequestMessage(
		context.Background(),
		"GET",
		fmt.Sprintf("https://dev.azure.com/azure-sdk/%s/_apis/test/codecoverage?buildId=%d&api-version=7.0", projectName, buildId),
		"7.0",
		nil,
		"",
		"",
		nil,
	)
	if err != nil {
		return 0, 0, err
	}

	response, err := client.SendRequest(request)
	if err != nil {
		return 0, 0, err
	}

	var codeCoverageResult test.CodeCoverageSummary
	err = client.UnmarshalBody(response, &codeCoverageResult)
	if err != nil {
		return 0, 0, err
	}

	for _, coverage := range *codeCoverageResult.CoverageData {
		if len(*coverage.CoverageStats) > 0 {
			return *(*coverage.CoverageStats)[0].Covered, *(*coverage.CoverageStats)[0].Total, nil
		}
	}

	return 0, 0, nil
}

func getCodeCoverage(pipelineClient pipelines.Client, azureDevopsClient *azuredevops.Client, info *mgmtInfo, pid int) (*int, error) {
	listRuns, err := pipelineClient.ListRuns(context.Background(), pipelines.ListRunsArgs{Project: &projectName, PipelineId: &pid})
	if err != nil && len(*listRuns) > 0 {
		return nil, err
	}

	var buildId *int
	for i := 0; i < 5 && i < len(*listRuns); i++ {
		buildId = (*listRuns)[i].Id

		// code coverage
		coveredLines, coverableLines, err := getBuildCodeCoverage(azureDevopsClient, projectName, *buildId)
		if err != nil {
			return nil, err
		}

		if coverableLines != 0 {
			info.CoveredLines = coveredLines
			info.CoverableLines = coverableLines
			break
		}
	}

	return buildId, nil
}

func getMockTestResult(testClient test.Client, info *mgmtInfo, buildId *int) error {
	buildUri := fmt.Sprintf("vstfs:///Build/Build/%d", *buildId)
	testRuns, err := testClient.GetTestRuns(context.Background(), test.GetTestRunsArgs{
		Project:  &projectName,
		BuildUri: &buildUri,
	})
	if err != nil {
		return err
	}
	for _, tr := range *testRuns {
		if strings.Contains(*tr.Name, "Test result on resourcemanager") {
			info.mockTestPass = *tr.PassedTests
			info.mockTestTotal = *tr.TotalTests
			return nil
		}
	}
	return nil
}

func getLiveTestResult(buildClient build.Client, info *mgmtInfo, buildId *int) error {
	buildLogs, err := buildClient.GetBuildLogs(context.Background(), build.GetBuildLogsArgs{
		Project: &projectName,
		BuildId: buildId,
	})
	if err != nil {
		return err
	}

	for i := 95; i < len(*buildLogs); i++ {
		logLines, err := buildClient.GetBuildLogLines(context.Background(), build.GetBuildLogLinesArgs{
			Project: &projectName,
			BuildId: buildId,
			LogId:   (*buildLogs)[i].Id,
		})
		if err != nil {
			return err
		}

		if logInfo := strings.Join(*logLines, "\n"); strings.Contains(logInfo, "Starting: Run Tests") && strings.Contains(logInfo, "Finishing: Run Tests") {
		loop:
			for _, line := range *logLines {
				if strings.Contains(line, "coverage:") {
					for _, j := range strings.Split(line, " ") {
						if strings.Contains(j, "%") {
							info.liveTestCoverage = j
							break loop
						}
					}
				}
			}
			break
		}
	}

	return nil
}
