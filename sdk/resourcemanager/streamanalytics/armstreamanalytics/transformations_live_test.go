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

type TransformationsTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	jobName            string
	transformationName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *TransformationsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.transformationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "transfor", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *TransformationsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestTransformationsTestSuite(t *testing.T) {
	suite.Run(t, new(TransformationsTestSuite))
}

func (testsuite *TransformationsTestSuite) Prepare() {
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

// Microsoft.StreamAnalytics/streamingjobs/{jobName}/transformations/{transformationName}
func (testsuite *TransformationsTestSuite) TestTransformations() {
	var err error
	// From step Transformations_CreateOrReplace
	fmt.Println("Call operation: Transformations_CreateOrReplace")
	transformationsClient, err := armstreamanalytics.NewTransformationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = transformationsClient.CreateOrReplace(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.transformationName, armstreamanalytics.Transformation{
		Properties: &armstreamanalytics.TransformationProperties{
			Query:          to.Ptr("Select Id, Name from inputtest"),
			StreamingUnits: to.Ptr[int32](6),
		},
	}, &armstreamanalytics.TransformationsClientCreateOrReplaceOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step Transformations_Get
	fmt.Println("Call operation: Transformations_Get")
	_, err = transformationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.transformationName, nil)
	testsuite.Require().NoError(err)

	// From step Transformations_Update
	fmt.Println("Call operation: Transformations_Update")
	_, err = transformationsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.jobName, testsuite.transformationName, armstreamanalytics.Transformation{
		Properties: &armstreamanalytics.TransformationProperties{
			Query: to.Ptr("New query"),
		},
	}, &armstreamanalytics.TransformationsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}
