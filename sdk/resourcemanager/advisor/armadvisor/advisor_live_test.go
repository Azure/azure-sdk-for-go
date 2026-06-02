// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armadvisor_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/advisor/armadvisor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type AdvisorTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	name              string
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AdvisorTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.name, _ = recording.GenerateAlphaNumericID(testsuite.T(), "suppressiona", 18, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AdvisorTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAdvisorTestSuite(t *testing.T) {
	suite.Run(t, new(AdvisorTestSuite))
}

// Microsoft.Advisor/generateRecommendations
func (testsuite *AdvisorTestSuite) TestRecommendations() {
	var recommendationId string
	var resourceUri string
	var err error
	// From step Recommendations_Generate
	fmt.Println("Call operation: Recommendations_Generate")
	recommendationsClient, err := armadvisor.NewRecommendationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	recommendationsClientGenerateResponse, err := recommendationsClient.Generate(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	locationId := *recommendationsClientGenerateResponse.Location
	operationId := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}").FindAllString(locationId, -1)[1]

	// From step Recommendations_GetGenerateStatus
	fmt.Println("Call operation: Recommendations_GetGenerateStatus")
	_, err = recommendationsClient.GetGenerateStatus(testsuite.ctx, operationId, nil)
	testsuite.Require().NoError(err)

	// From step Recommendations_List
	fmt.Println("Call operation: Recommendations_List")
	recommendationsClientNewListPager := recommendationsClient.NewListPager(&armadvisor.RecommendationsClientListOptions{Filter: nil,
		Top:       to.Ptr[int32](10),
		SkipToken: nil,
	})
	for recommendationsClientNewListPager.More() {
		recommendationsClientListResponse, err := recommendationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		resourceUri, _, _ = strings.Cut(*recommendationsClientListResponse.Value[0].ID, "/providers/Microsoft.Advisor/recommendations")
		recommendationId = *recommendationsClientListResponse.Value[0].Name
		break
	}

	// From step Recommendations_Get
	fmt.Println("Call operation: Recommendations_Get")
	_, err = recommendationsClient.Get(testsuite.ctx, resourceUri, recommendationId, nil)
	testsuite.Require().NoError(err)

	// From step Suppressions_Create
	fmt.Println("Call operation: Suppressions_Create")
	suppressionsClient, err := armadvisor.NewSuppressionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = suppressionsClient.Create(testsuite.ctx, resourceUri, recommendationId, testsuite.name, armadvisor.SuppressionContract{
		Properties: &armadvisor.SuppressionProperties{
			TTL: to.Ptr("07:00:00:00"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Suppressions_List
	fmt.Println("Call operation: Suppressions_List")
	suppressionsClientNewListPager := suppressionsClient.NewListPager(&armadvisor.SuppressionsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for suppressionsClientNewListPager.More() {
		_, err := suppressionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Suppressions_Get
	fmt.Println("Call operation: Suppressions_Get")
	_, err = suppressionsClient.Get(testsuite.ctx, resourceUri, recommendationId, testsuite.name, nil)
	testsuite.Require().NoError(err)

	// From step Suppressions_Delete
	fmt.Println("Call operation: Suppressions_Delete")
	_, err = suppressionsClient.Delete(testsuite.ctx, resourceUri, recommendationId, testsuite.name, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Advisor/configurations/{configurationName}
func (testsuite *AdvisorTestSuite) TestConfigurations() {
	resourceGroup := testsuite.resourceGroupName
	var err error
	// From step Configurations_CreateInSubscription
	fmt.Println("Call operation: Configurations_CreateInSubscription")
	configurationsClient, err := armadvisor.NewConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = configurationsClient.CreateInSubscription(testsuite.ctx, armadvisor.ConfigurationNameDefault, armadvisor.ConfigData{
		Properties: &armadvisor.ConfigDataProperties{
			LowCPUThreshold: to.Ptr(armadvisor.CPUThresholdFive),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Configurations_ListBySubscription
	fmt.Println("Call operation: Configurations_ListBySubscription")
	configurationsClientNewListBySubscriptionPager := configurationsClient.NewListBySubscriptionPager(nil)
	for configurationsClientNewListBySubscriptionPager.More() {
		_, err := configurationsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Configurations_CreateInResourceGroup
	fmt.Println("Call operation: Configurations_CreateInResourceGroup")
	_, err = configurationsClient.CreateInResourceGroup(testsuite.ctx, armadvisor.ConfigurationNameDefault, resourceGroup, armadvisor.ConfigData{
		Properties: &armadvisor.ConfigDataProperties{
			Exclude: to.Ptr(false),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Configurations_ListByResourceGroup
	fmt.Println("Call operation: Configurations_ListByResourceGroup")
	configurationsClientNewListByResourceGroupPager := configurationsClient.NewListByResourceGroupPager(resourceGroup, nil)
	for configurationsClientNewListByResourceGroupPager.More() {
		_, err := configurationsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Advisor/metadata
func (testsuite *AdvisorTestSuite) TestRecommendationMetadata() {
	var err error
	var recommendationMetadataName string
	// From step RecommendationMetadata_List
	fmt.Println("Call operation: RecommendationMetadata_List")
	recommendationMetadataClient, err := armadvisor.NewRecommendationMetadataClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	recommendationMetadataClientNewListPager := recommendationMetadataClient.NewListPager(nil)
	for recommendationMetadataClientNewListPager.More() {
		recommendationMetadataClientListResponse, err := recommendationMetadataClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		recommendationMetadataName = *recommendationMetadataClientListResponse.Value[0].Name
		break
	}

	// From step RecommendationMetadata_Get
	fmt.Println("Call operation: RecommendationMetadata_Get")
	_, err = recommendationMetadataClient.Get(testsuite.ctx, recommendationMetadataName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Advisor/operations
func (testsuite *AdvisorTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armadvisor.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
