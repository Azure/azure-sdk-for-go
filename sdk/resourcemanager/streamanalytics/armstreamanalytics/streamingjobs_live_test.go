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

type StreamingjobsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	jobName           string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *StreamingjobsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *StreamingjobsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestStreamingjobsTestSuite(t *testing.T) {
	suite.Run(t, new(StreamingjobsTestSuite))
}

// Microsoft.StreamAnalytics/operations
func (testsuite *StreamingjobsTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armstreamanalytics.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.StreamAnalytics/streamingjobs/{jobName}
func (testsuite *StreamingjobsTestSuite) TestStreamingJobs() {
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

	// From step StreamingJobs_List
	fmt.Println("Call operation: StreamingJobs_List")
	streamingJobsClientNewListPager := streamingJobsClient.NewListPager(&armstreamanalytics.StreamingJobsClientListOptions{Expand: nil})
	for streamingJobsClientNewListPager.More() {
		_, err := streamingJobsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StreamingJobs_ListByResourceGroup
	fmt.Println("Call operation: StreamingJobs_ListByResourceGroup")
	streamingJobsClientNewListByResourceGroupPager := streamingJobsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armstreamanalytics.StreamingJobsClientListByResourceGroupOptions{Expand: nil})
	for streamingJobsClientNewListByResourceGroupPager.More() {
		_, err := streamingJobsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StreamingJobs_Get
	fmt.Println("Call operation: StreamingJobs_Get")
	_, err = streamingJobsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, &armstreamanalytics.StreamingJobsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step StreamingJobs_Update
	fmt.Println("Call operation: StreamingJobs_Update")
	_, err = streamingJobsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, armstreamanalytics.StreamingJob{
		Properties: &armstreamanalytics.StreamingJobProperties{
			EventsLateArrivalMaxDelayInSeconds: to.Ptr[int32](13),
			EventsOutOfOrderMaxDelayInSeconds:  to.Ptr[int32](21),
		},
	}, &armstreamanalytics.StreamingJobsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step StreamingJobs_Delete
	fmt.Println("Call operation: StreamingJobs_Delete")
	streamingJobsClientDeleteResponsePoller, err := streamingJobsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, streamingJobsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
