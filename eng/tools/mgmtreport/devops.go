// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelines"
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

func getBuildId(pipelineClient pipelines.Client, azureDevopsClient *azuredevops.Client, pid int) (int, error) {
	listRuns, err := pipelineClient.ListRuns(context.Background(), pipelines.ListRunsArgs{Project: &projectName, PipelineId: &pid})
	if err != nil {
		return 0, err
	}

	var buildId *int
	for i := 0; i < 5 && i < len(*listRuns); i++ {
		buildId = (*listRuns)[i].Id

		_, coverableLines, err := getBuildCodeCoverage(azureDevopsClient, projectName, *buildId)
		if err != nil {
			return 0, err
		}

		if coverableLines != 0 {
			break
		}
	}

	return *buildId, nil
}

func getLogID(buildClient build.Client, buildId int) (int, int, error) {
	var mockTestLogId int
	var liveTestLogId int

	result, err := buildClient.GetBuildTimeline(context.Background(), build.GetBuildTimelineArgs{
		Project: &projectName,
		BuildId: &buildId,
	})
	if err != nil {
		return mockTestLogId, liveTestLogId, err
	}

	for _, record := range *result.Records {
		if mockTestLogId != 0 && liveTestLogId != 0 {
			break
		} else if *record.State == build.TimelineRecordStateValues.Completed && *record.Result == build.TaskResultValues.Succeeded {
			if *record.Name == "Mock Test" && mockTestLogId == 0 {
				mockTestLogId = *record.Log.Id
			} else if *record.Name == "Run Tests" && liveTestLogId == 0 {
				liveTestLogId = *record.Log.Id
			}
		}
	}

	return mockTestLogId, liveTestLogId, nil
}

func getTestResult(buildClient build.Client, buildId, logId int) (int, int, string, error) {
	logLines, err := buildClient.GetBuildLogLines(context.Background(), build.GetBuildLogLinesArgs{
		Project: &projectName,
		BuildId: &buildId,
		LogId:   &logId,
	})
	if err != nil {
		return 0, 0, "", err
	}

	logResult := strings.Join(*logLines, "\n")
	totalTests := regexp.MustCompile("=== RUN.*/").FindAllString(logResult, -1)

	passedTests := regexp.MustCompile("--- PASS:.*/").FindAllString(logResult, -1)

	coverages := regexp.MustCompile("coverage:.*").FindAllString(logResult, -1)

	var coverage string
	if len(coverages) > 0 {
		for _, s := range strings.Split(coverages[0], " ") {
			if strings.Contains(s, "%") {
				coverage = s
			}
		}
	}

	return len(totalTests), len(passedTests), coverage, nil
}

func getOperation(buildClient build.Client, buildId, mockTestLogId int) (map[string]struct{}, int, error) {
	logLines, err := buildClient.GetBuildLogLines(context.Background(), build.GetBuildLogLinesArgs{
		Project: &projectName,
		BuildId: &buildId,
		LogId:   &mockTestLogId,
	})
	if err != nil {
		return nil, 0, err
	}

	logResult := strings.Join(*logLines, "\n")
	skipOperations := regexp.MustCompile("--- SKIP:.*/").FindAllString(logResult, -1)

	allOperation := make(map[string]struct{})
	var v struct{}
	all := regexp.MustCompile("=== RUN.*/.*").FindAllString(logResult, -1)
	for _, o := range all {
		_, after, b := strings.Cut(o, "/")
		if b {
			operation := strings.TrimPrefix(after, "Test")
			if _, ok := allOperation[operation]; !ok {
				allOperation[operation] = v
			}
		}
	}

	return allOperation, len(skipOperations), nil
}

func getLiveTestOperationCalls(buildClient build.Client, allOperation map[string]struct{}, buildId, liveTestLogId int) (int, error) {
	logLines, err := buildClient.GetBuildLogLines(context.Background(), build.GetBuildLogLinesArgs{
		Project: &projectName,
		BuildId: &buildId,
		LogId:   &liveTestLogId,
	})
	if err != nil {
		return 0, err
	}

	logResult := strings.Join(*logLines, "\n")
	result := regexp.MustCompile("Call operation:.*").FindAllString(logResult, -1)

	callOperation := make(map[string]struct{})
	var v struct{}
	for _, co := range result {
		_, after, b := strings.Cut(co, "Call operation: ")
		if b {
			if _, ok := allOperation[after]; ok {
				if _, ok = callOperation[after]; !ok {
					callOperation[after] = v
				}
			}
		}
	}

	return len(callOperation), nil
}
