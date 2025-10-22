//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armiothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub"
	"github.com/stretchr/testify/suite"
)

type IothubTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	certificateName   string
	iothubId          string
	resourceName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *IothubTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.certificateName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "certific", 14, false)
	testsuite.resourceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "resource", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *IothubTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestIothubTestSuite(t *testing.T) {
	suite.Run(t, new(IothubTestSuite))
}

func (testsuite *IothubTestSuite) Prepare() {
	var err error
	// From step IotHubResource_CreateOrUpdate
	fmt.Println("Call operation: IotHubResource_CreateOrUpdate")
	resourceClient, err := armiothub.NewResourceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceClientCreateOrUpdateResponsePoller, err := resourceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armiothub.Description{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Etag:     to.Ptr("AAAAAAFD6M4="),
		Properties: &armiothub.Properties{
			CloudToDevice: &armiothub.CloudToDeviceProperties{
				DefaultTTLAsISO8601: to.Ptr("PT1H"),
				Feedback: &armiothub.FeedbackProperties{
					LockDurationAsISO8601: to.Ptr("PT1M"),
					MaxDeliveryCount:      to.Ptr[int32](10),
					TTLAsISO8601:          to.Ptr("PT1H"),
				},
				MaxDeliveryCount: to.Ptr[int32](10),
			},
			EnableDataResidency:           to.Ptr(false),
			EnableFileUploadNotifications: to.Ptr(false),
			EventHubEndpoints: map[string]*armiothub.EventHubProperties{
				"events": {
					PartitionCount:      to.Ptr[int32](2),
					RetentionTimeInDays: to.Ptr[int64](1),
				},
			},
			Features:      to.Ptr(armiothub.CapabilitiesNone),
			IPFilterRules: []*armiothub.IPFilterRule{},
			MessagingEndpoints: map[string]*armiothub.MessagingEndpointProperties{
				"fileNotifications": {
					LockDurationAsISO8601: to.Ptr("PT1M"),
					MaxDeliveryCount:      to.Ptr[int32](10),
					TTLAsISO8601:          to.Ptr("PT1H"),
				},
			},
			MinTLSVersion: to.Ptr("1.2"),
			NetworkRuleSets: &armiothub.NetworkRuleSetProperties{
				ApplyToBuiltInEventHubEndpoint: to.Ptr(true),
				DefaultAction:                  to.Ptr(armiothub.DefaultActionDeny),
				IPRules: []*armiothub.NetworkRuleSetIPRule{
					{
						Action:     to.Ptr(armiothub.NetworkRuleIPActionAllow),
						FilterName: to.Ptr("rule1"),
						IPMask:     to.Ptr("131.117.159.53"),
					},
					{
						Action:     to.Ptr(armiothub.NetworkRuleIPActionAllow),
						FilterName: to.Ptr("rule2"),
						IPMask:     to.Ptr("157.55.59.128/25"),
					}},
			},
			Routing: &armiothub.RoutingProperties{
				Endpoints: &armiothub.RoutingEndpoints{
					EventHubs:         []*armiothub.RoutingEventHubProperties{},
					ServiceBusQueues:  []*armiothub.RoutingServiceBusQueueEndpointProperties{},
					ServiceBusTopics:  []*armiothub.RoutingServiceBusTopicEndpointProperties{},
					StorageContainers: []*armiothub.RoutingStorageContainerProperties{},
				},
				FallbackRoute: &armiothub.FallbackRouteProperties{
					Name:      to.Ptr("$fallback"),
					Condition: to.Ptr("true"),
					EndpointNames: []*string{
						to.Ptr("events")},
					IsEnabled: to.Ptr(true),
					Source:    to.Ptr(armiothub.RoutingSourceDeviceMessages),
				},
				Routes: []*armiothub.RouteProperties{},
			},
			StorageEndpoints: map[string]*armiothub.StorageEndpointProperties{
				"$default": {
					ConnectionString: to.Ptr(""),
					ContainerName:    to.Ptr(""),
					SasTTLAsISO8601:  to.Ptr("PT1H"),
				},
			},
		},
		SKU: &armiothub.SKUInfo{
			Name:     to.Ptr(armiothub.IotHubSKUS1),
			Capacity: to.Ptr[int64](1),
		},
	}, &armiothub.ResourceClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	var resourceClientCreateOrUpdateResponse *armiothub.ResourceClientCreateOrUpdateResponse
	resourceClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, resourceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.iothubId = *resourceClientCreateOrUpdateResponse.ID
}

// Microsoft.Devices/operations
func (testsuite *IothubTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armiothub.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ResourceProviderCommon_GetSubscriptionQuota
	fmt.Println("Call operation: ResourceProviderCommon_GetSubscriptionQuota")
	resourceProviderCommonClient, err := armiothub.NewResourceProviderCommonClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = resourceProviderCommonClient.GetSubscriptionQuota(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Devices/IotHubs/{resourceName}
func (testsuite *IothubTestSuite) TestIotHubResource() {
	iotHubName := testsuite.resourceName
	var keyName string
	var err error
	// From step IotHubResource_CheckNameAvailability
	fmt.Println("Call operation: IotHubResource_CheckNameAvailability")
	resourceClient, err := armiothub.NewResourceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = resourceClient.CheckNameAvailability(testsuite.ctx, armiothub.OperationInputs{
		Name: to.Ptr("test-request"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IotHubResource_ListBySubscription
	fmt.Println("Call operation: IotHubResource_ListBySubscription")
	resourceClientNewListBySubscriptionPager := resourceClient.NewListBySubscriptionPager(nil)
	for resourceClientNewListBySubscriptionPager.More() {
		_, err := resourceClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_GetStats
	fmt.Println("Call operation: IotHubResource_GetStats")
	_, err = resourceClient.GetStats(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)

	// From step IotHubResource_ListEventHubConsumerGroups
	fmt.Println("Call operation: IotHubResource_ListEventHubConsumerGroups")
	resourceClientNewListEventHubConsumerGroupsPager := resourceClient.NewListEventHubConsumerGroupsPager(testsuite.resourceGroupName, testsuite.resourceName, "events", nil)
	for resourceClientNewListEventHubConsumerGroupsPager.More() {
		_, err := resourceClientNewListEventHubConsumerGroupsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_GetQuotaMetrics
	fmt.Println("Call operation: IotHubResource_GetQuotaMetrics")
	resourceClientNewGetQuotaMetricsPager := resourceClient.NewGetQuotaMetricsPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for resourceClientNewGetQuotaMetricsPager.More() {
		_, err := resourceClientNewGetQuotaMetricsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_ListByResourceGroup
	fmt.Println("Call operation: IotHubResource_ListByResourceGroup")
	resourceClientNewListByResourceGroupPager := resourceClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for resourceClientNewListByResourceGroupPager.More() {
		_, err := resourceClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_GetValidSkus
	fmt.Println("Call operation: IotHubResource_GetValidSKUs")
	resourceClientNewGetValidSKUsPager := resourceClient.NewGetValidSKUsPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for resourceClientNewGetValidSKUsPager.More() {
		_, err := resourceClientNewGetValidSKUsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_ListJobs
	fmt.Println("Call operation: IotHubResource_ListJobs")
	resourceClientNewListJobsPager := resourceClient.NewListJobsPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for resourceClientNewListJobsPager.More() {
		_, err := resourceClientNewListJobsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IotHubResource_Update
	fmt.Println("Call operation: IotHubResource_Update")
	resourceClientUpdateResponsePoller, err := resourceClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, armiothub.TagsResource{
		Tags: map[string]*string{
			"foo": to.Ptr("bar"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, resourceClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step IotHubResource_ListKeys
	fmt.Println("Call operation: IotHubResource_ListKeys")
	resourceClientNewListKeysPager := resourceClient.NewListKeysPager(testsuite.resourceGroupName, testsuite.resourceName, nil)
	for resourceClientNewListKeysPager.More() {
		nextResult, err := resourceClientNewListKeysPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		keyName = *nextResult.Value[0].KeyName
		break
	}

	// From step IotHubResource_GetKeysForKeyName
	fmt.Println("Call operation: IotHubResource_GetKeysForKeyName")
	_, err = resourceClient.GetKeysForKeyName(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, keyName, nil)
	testsuite.Require().NoError(err)

	// From step IotHub_ManualFailover
	fmt.Println("Call operation: IotHub_ManualFailover")
	client, err := armiothub.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientManualFailoverResponsePoller, err := client.BeginManualFailover(testsuite.ctx, iotHubName, testsuite.resourceGroupName, armiothub.FailoverInput{
		FailoverRegion: to.Ptr("eastus"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientManualFailoverResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *IothubTestSuite) Cleanup() {
	var err error
	// From step IotHubResource_Delete
	fmt.Println("Call operation: IotHubResource_Delete")
	resourceClient, err := armiothub.NewResourceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceClientDeleteResponsePoller, err := resourceClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.resourceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, resourceClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
