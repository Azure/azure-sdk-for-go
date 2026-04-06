// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstreamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics/v2"
	"github.com/stretchr/testify/suite"
)

type FunctionsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	functionName      string
	jobName           string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *FunctionsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.functionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "function", 14, false)
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *FunctionsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestFunctionsTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionsTestSuite))
}

func (testsuite *FunctionsTestSuite) Prepare() {
	var err error
	// From step StreamingJobs_CreateOrReplace
	fmt.Println("Call operation: StreamingJobs_CreateOrReplace")
	streamingJobsClient, err := armstreamanalytics.NewStreamingJobsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	streamingJobsClientCreateOrReplaceResponsePoller, err := streamingJobsClient.BeginCreateOrReplace(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, armstreamanalytics.StreamingJob{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1":      to.Ptr("value1"),
			"key3":      to.Ptr("value3"),
			"randomKey": to.Ptr("randomValue"),
		},
		Properties: &armstreamanalytics.StreamingJobProperties{
			CompatibilityLevel:                 to.Ptr(armstreamanalytics.CompatibilityLevelOne0),
			DataLocale:                         to.Ptr("en-US"),
			EventsLateArrivalMaxDelayInSeconds: to.Ptr[int32](16),
			EventsOutOfOrderMaxDelayInSeconds:  to.Ptr[int32](5),
			EventsOutOfOrderPolicy:             to.Ptr(armstreamanalytics.EventsOutOfOrderPolicyDrop),
			Functions:                          []*armstreamanalytics.Function{},
			Inputs:                             []*armstreamanalytics.Input{},
			OutputErrorPolicy:                  to.Ptr(armstreamanalytics.OutputErrorPolicyDrop),
			Outputs:                            []*armstreamanalytics.Output{},
			SKU: &armstreamanalytics.SKU{
				Name: to.Ptr(armstreamanalytics.SKUNameStandard),
			},
		},
	}, &armstreamanalytics.StreamingJobsClientBeginCreateOrReplaceOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingJobsClientCreateOrReplaceResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}
func (testsuite *FunctionsTestSuite) TestFunctions() {
	var err error
	// From step Functions_CreateOrReplace
	fmt.Println("Call operation: Functions_CreateOrReplace")
	functionsClient, err := armstreamanalytics.NewFunctionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = functionsClient.CreateOrReplace(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.functionName, armstreamanalytics.Function{
		Properties: &armstreamanalytics.ScalarFunctionProperties{
			Type: to.Ptr("Scalar"),
			Properties: &armstreamanalytics.FunctionConfiguration{
				Binding: &armstreamanalytics.JavaScriptFunctionBinding{
					Type: to.Ptr("Microsoft.StreamAnalytics/JavascriptUdf"),
					Properties: &armstreamanalytics.JavaScriptFunctionBindingProperties{
						Script: to.Ptr("function (x, y) { return x + y; }"),
					},
				},
				Inputs: []*armstreamanalytics.FunctionInput{
					{
						DataType: to.Ptr("Any"),
					}},
				Output: &armstreamanalytics.FunctionOutput{
					DataType: to.Ptr("Any"),
				},
			},
		},
	}, &armstreamanalytics.FunctionsClientCreateOrReplaceOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step Functions_ListByStreamingJob
	fmt.Println("Call operation: Functions_ListByStreamingJob")
	functionsClientNewListByStreamingJobPager := functionsClient.NewListByStreamingJobPager(testsuite.resourceGroupName, testsuite.jobName, &armstreamanalytics.FunctionsClientListByStreamingJobOptions{Select: nil})
	for functionsClientNewListByStreamingJobPager.More() {
		_, err := functionsClientNewListByStreamingJobPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Functions_Get
	fmt.Println("Call operation: Functions_Get")
	_, err = functionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.functionName, nil)
	testsuite.Require().NoError(err)

	// From step Functions_Update
	fmt.Println("Call operation: Functions_Update")
	_, err = functionsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.functionName, armstreamanalytics.Function{
		Properties: &armstreamanalytics.ScalarFunctionProperties{
			Type: to.Ptr("Scalar"),
			Properties: &armstreamanalytics.FunctionConfiguration{
				Binding: &armstreamanalytics.JavaScriptFunctionBinding{
					Type: to.Ptr("Microsoft.StreamAnalytics/JavascriptUdf"),
					Properties: &armstreamanalytics.JavaScriptFunctionBindingProperties{
						Script: to.Ptr("function (a, b) { return a * b; }"),
					},
				},
			},
		},
	}, &armstreamanalytics.FunctionsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Functions_Test
	fmt.Println("Call operation: Functions_Test")
	functionsClientTestResponsePoller, err := functionsClient.BeginTest(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.functionName, &armstreamanalytics.FunctionsClientBeginTestOptions{Function: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, functionsClientTestResponsePoller)
	testsuite.Require().NoError(err)

	// From step Functions_Delete
	fmt.Println("Call operation: Functions_Delete")
	_, err = functionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.functionName, nil)
	testsuite.Require().NoError(err)
}
