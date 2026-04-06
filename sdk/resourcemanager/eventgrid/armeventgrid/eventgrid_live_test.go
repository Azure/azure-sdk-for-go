// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armeventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type EventGridTestSuite struct {
	suite.Suite

	ctx                     context.Context
	cred                    azcore.TokenCredential
	options                 *arm.ClientOptions
	domainId                string
	domainName              string
	domainTopicName         string
	eventSubscriptionName   string
	eventhubId              string
	namespaceName           string
	partnerNamespaceName    string
	partnerRegistrationId   string
	partnerRegistrationName string
	privateEndpointName     string
	systemTopicName         string
	topicName               string
	virtualNetworksName     string
	location                string
	resourceGroupName       string
	subscriptionId          string
}

func (testsuite *EventGridTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.domainName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "domainname", 16, false)
	testsuite.domainTopicName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "domaintopi", 16, false)
	testsuite.eventSubscriptionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventsubsc", 16, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventhubna", 16, false)
	testsuite.partnerNamespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "partnernam", 16, false)
	testsuite.partnerRegistrationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "partnerreg", 16, false)
	testsuite.privateEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventgridprivateendpoint", 30, false)
	testsuite.systemTopicName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "systemtopi", 16, false)
	testsuite.topicName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "topicname", 15, false)
	testsuite.virtualNetworksName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventgridvnet", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *EventGridTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestEventGridTestSuite(t *testing.T) {
	suite.Run(t, new(EventGridTestSuite))
}

func (testsuite *EventGridTestSuite) Prepare() {
	var err error
	// From step Domains_CreateOrUpdate
	fmt.Println("Call operation: Domains_CreateOrUpdate")
	domainsClient, err := armeventgrid.NewDomainsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	domainsClientCreateOrUpdateResponsePoller, err := domainsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, armeventgrid.Domain{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armeventgrid.DomainProperties{
			InboundIPRules: []*armeventgrid.InboundIPRule{
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.30.15"),
				},
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.176.1"),
				}},
			PublicNetworkAccess: to.Ptr(armeventgrid.PublicNetworkAccessEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var domainsClientCreateOrUpdateResponse *armeventgrid.DomainsClientCreateOrUpdateResponse
	domainsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, domainsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.domainId = *domainsClientCreateOrUpdateResponse.ID

	// From step Topics_CreateOrUpdate
	fmt.Println("Call operation: Topics_CreateOrUpdate")
	topicsClient, err := armeventgrid.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicsClientCreateOrUpdateResponsePoller, err := topicsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, armeventgrid.Topic{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armeventgrid.TopicProperties{
			InboundIPRules: []*armeventgrid.InboundIPRule{
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.30.15"),
				},
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.176.1"),
				}},
			PublicNetworkAccess: to.Ptr(armeventgrid.PublicNetworkAccessEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerRegistrations_CreateOrUpdate
	fmt.Println("Call operation: PartnerRegistrations_CreateOrUpdate")
	partnerRegistrationsClient, err := armeventgrid.NewPartnerRegistrationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	partnerRegistrationsClientCreateOrUpdateResponsePoller, err := partnerRegistrationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerRegistrationName, armeventgrid.PartnerRegistration{
		Location: to.Ptr("global"),
	}, nil)
	testsuite.Require().NoError(err)
	var partnerRegistrationsClientCreateOrUpdateResponse *armeventgrid.PartnerRegistrationsClientCreateOrUpdateResponse
	partnerRegistrationsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, partnerRegistrationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.partnerRegistrationId = *partnerRegistrationsClientCreateOrUpdateResponse.ID

	// From step SystemTopics_CreateOrUpdate
	fmt.Println("Call operation: SystemTopics_CreateOrUpdate")
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	systemTopicsClientCreateOrUpdateResponsePoller, err := systemTopicsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, armeventgrid.SystemTopic{
		Location: to.Ptr("global"),
		Properties: &armeventgrid.SystemTopicProperties{
			Source:    to.Ptr("/subscriptions/" + testsuite.subscriptionId),
			TopicType: to.Ptr("Microsoft.Resources.Subscriptions"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Eventhub_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"eventhubId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.EventHub/namespaces/eventhubs', parameters('namespaceName'), 'eventhubs')]",
			},
		},
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"namespaceName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.namespaceName,
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('namespaceName')]",
				"type":       "Microsoft.EventHub/namespaces",
				"apiVersion": "2022-01-01-preview",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"disableLocalAuth":       false,
					"isAutoInflateEnabled":   false,
					"kafkaEnabled":           true,
					"maximumThroughputUnits": float64(0),
					"minimumTlsVersion":      "1.0",
					"publicNetworkAccess":    "Enabled",
					"zoneRedundant":          false,
				},
				"sku": map[string]interface{}{
					"name":     "Standard",
					"capacity": float64(1),
					"tier":     "Standard",
				},
			},
			map[string]interface{}{
				"name":       "[concat(parameters('namespaceName'), '/eventhubs')]",
				"type":       "Microsoft.EventHub/namespaces/eventhubs",
				"apiVersion": "2022-01-01-preview",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.EventHub/namespaces', parameters('namespaceName'))]",
				},
				"location": "[parameters('location')]",
				"properties": map[string]interface{}{
					"messageRetentionInDays": float64(1),
					"partitionCount":         float64(1),
					"status":                 "Active",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Eventhub_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.eventhubId = deploymentExtend.Properties.Outputs.(map[string]interface{})["eventhubId"].(map[string]interface{})["value"].(string)
}

// Microsoft.EventGrid/domains
func (testsuite *EventGridTestSuite) TestDomains() {
	var err error
	// From step Domains_ListBySubscription
	fmt.Println("Call operation: Domains_ListBySubscription")
	domainsClient, err := armeventgrid.NewDomainsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	domainsClientNewListBySubscriptionPager := domainsClient.NewListBySubscriptionPager(&armeventgrid.DomainsClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for domainsClientNewListBySubscriptionPager.More() {
		_, err := domainsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Domains_ListByResourceGroup
	fmt.Println("Call operation: Domains_ListByResourceGroup")
	domainsClientNewListByResourceGroupPager := domainsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.DomainsClientListByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for domainsClientNewListByResourceGroupPager.More() {
		_, err := domainsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Domains_Get
	fmt.Println("Call operation: Domains_Get")
	_, err = domainsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, nil)
	testsuite.Require().NoError(err)

	// From step Domains_Update
	fmt.Println("Call operation: Domains_Update")
	domainsClientUpdateResponsePoller, err := domainsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, armeventgrid.DomainUpdateParameters{
		Properties: &armeventgrid.DomainUpdateParameterProperties{
			InboundIPRules: []*armeventgrid.InboundIPRule{
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.30.15"),
				},
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.176.1"),
				}},
			PublicNetworkAccess: to.Ptr(armeventgrid.PublicNetworkAccessEnabled),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Domains_ListSharedAccessKeys
	fmt.Println("Call operation: Domains_ListSharedAccessKeys")
	_, err = domainsClient.ListSharedAccessKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, nil)
	testsuite.Require().NoError(err)

	// From step Domains_RegenerateKey
	fmt.Println("Call operation: Domains_RegenerateKey")
	_, err = domainsClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, armeventgrid.DomainRegenerateKeyRequest{
		KeyName: to.Ptr("key1"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/domains/eventSubscriptions
func (testsuite *EventGridTestSuite) TestDomainEventSubscriptions() {
	var err error
	// From step DomainEventSubscriptions_CreateOrUpdate
	fmt.Println("Call operation: DomainEventSubscriptions_CreateOrUpdate")
	domainEventSubscriptionsClient, err := armeventgrid.NewDomainEventSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	domainEventSubscriptionsClientCreateOrUpdateResponsePoller, err := domainEventSubscriptionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.eventSubscriptionName, armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.EventHubEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeEventHub),
				Properties: &armeventgrid.EventHubEventSubscriptionDestinationProperties{
					ResourceID: to.Ptr(testsuite.eventhubId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainEventSubscriptionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DomainEventSubscriptions_List
	fmt.Println("Call operation: DomainEventSubscriptions_List")
	domainEventSubscriptionsClientNewListPager := domainEventSubscriptionsClient.NewListPager(testsuite.resourceGroupName, testsuite.domainName, &armeventgrid.DomainEventSubscriptionsClientListOptions{Filter: nil,
		Top: nil,
	})
	for domainEventSubscriptionsClientNewListPager.More() {
		_, err := domainEventSubscriptionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DomainEventSubscriptions_Get
	fmt.Println("Call operation: DomainEventSubscriptions_Get")
	_, err = domainEventSubscriptionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step DomainEventSubscriptions_Update
	fmt.Println("Call operation: DomainEventSubscriptions_Update")
	domainEventSubscriptionsClientUpdateResponsePoller, err := domainEventSubscriptionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.eventSubscriptionName, armeventgrid.EventSubscriptionUpdateParameters{
		Labels: []*string{
			to.Ptr("label1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainEventSubscriptionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DomainEventSubscriptions_GetDeliveryAttributes
	fmt.Println("Call operation: DomainEventSubscriptions_GetDeliveryAttributes")
	_, err = domainEventSubscriptionsClient.GetDeliveryAttributes(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step DomainEventSubscriptions_Delete
	fmt.Println("Call operation: DomainEventSubscriptions_Delete")
	domainEventSubscriptionsClientDeleteResponsePoller, err := domainEventSubscriptionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainEventSubscriptionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/domains/topics
func (testsuite *EventGridTestSuite) TestDomainTopics() {
	var err error
	// From step DomainTopics_CreateOrUpdate
	fmt.Println("Call operation: DomainTopics_CreateOrUpdate")
	domainTopicsClient, err := armeventgrid.NewDomainTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	domainTopicsClientCreateOrUpdateResponsePoller, err := domainTopicsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.domainTopicName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainTopicsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DomainTopics_ListByDomain
	fmt.Println("Call operation: DomainTopics_ListByDomain")
	domainTopicsClientNewListByDomainPager := domainTopicsClient.NewListByDomainPager(testsuite.resourceGroupName, testsuite.domainName, &armeventgrid.DomainTopicsClientListByDomainOptions{Filter: nil,
		Top: nil,
	})
	for domainTopicsClientNewListByDomainPager.More() {
		_, err := domainTopicsClientNewListByDomainPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DomainTopics_Get
	fmt.Println("Call operation: DomainTopics_Get")
	_, err = domainTopicsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.domainTopicName, nil)
	testsuite.Require().NoError(err)

	// From step DomainTopicEventSubscriptions_CreateOrUpdate
	fmt.Println("Call operation: DomainTopicEventSubscriptions_CreateOrUpdate")
	domainTopicEventSubscriptionsClient, err := armeventgrid.NewDomainTopicEventSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	domainTopicEventSubscriptionsClientCreateOrUpdateResponsePoller, err := domainTopicEventSubscriptionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, testsuite.eventSubscriptionName, armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.EventHubEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeEventHub),
				Properties: &armeventgrid.EventHubEventSubscriptionDestinationProperties{
					ResourceID: to.Ptr(testsuite.eventhubId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainTopicEventSubscriptionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DomainTopicEventSubscriptions_List
	fmt.Println("Call operation: DomainTopicEventSubscriptions_List")
	domainTopicEventSubscriptionsClientNewListPager := domainTopicEventSubscriptionsClient.NewListPager(testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, &armeventgrid.DomainTopicEventSubscriptionsClientListOptions{Filter: nil,
		Top: nil,
	})
	for domainTopicEventSubscriptionsClientNewListPager.More() {
		_, err := domainTopicEventSubscriptionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DomainTopicEventSubscriptions_Get
	fmt.Println("Call operation: DomainTopicEventSubscriptions_Get")
	_, err = domainTopicEventSubscriptionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step DomainTopicEventSubscriptions_GetDeliveryAttributes
	fmt.Println("Call operation: DomainTopicEventSubscriptions_GetDeliveryAttributes")
	_, err = domainTopicEventSubscriptionsClient.GetDeliveryAttributes(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step DomainTopicEventSubscriptions_Delete
	fmt.Println("Call operation: DomainTopicEventSubscriptions_Delete")
	domainTopicEventSubscriptionsClientDeleteResponsePoller, err := domainTopicEventSubscriptionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainTopicEventSubscriptionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step DomainTopics_Delete
	fmt.Println("Call operation: DomainTopics_Delete")
	domainTopicsClientDeleteResponsePoller, err := domainTopicsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.domainName, testsuite.domainTopicName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, domainTopicsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/topics
func (testsuite *EventGridTestSuite) TestTopics() {
	var err error
	// From step Topics_ListBySubscription
	fmt.Println("Call operation: Topics_ListBySubscription")
	topicsClient, err := armeventgrid.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicsClientNewListBySubscriptionPager := topicsClient.NewListBySubscriptionPager(&armeventgrid.TopicsClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for topicsClientNewListBySubscriptionPager.More() {
		_, err := topicsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Topics_ListByResourceGroup
	fmt.Println("Call operation: Topics_ListByResourceGroup")
	topicsClientNewListByResourceGroupPager := topicsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.TopicsClientListByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for topicsClientNewListByResourceGroupPager.More() {
		_, err := topicsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Topics_Get
	fmt.Println("Call operation: Topics_Get")
	_, err = topicsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, nil)
	testsuite.Require().NoError(err)

	// From step Topics_Update
	fmt.Println("Call operation: Topics_Update")
	topicsClientUpdateResponsePoller, err := topicsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, armeventgrid.TopicUpdateParameters{
		Properties: &armeventgrid.TopicUpdateParameterProperties{
			InboundIPRules: []*armeventgrid.InboundIPRule{
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.30.15"),
				},
				{
					Action: to.Ptr(armeventgrid.IPActionTypeAllow),
					IPMask: to.Ptr("12.18.176.1"),
				}},
			PublicNetworkAccess: to.Ptr(armeventgrid.PublicNetworkAccessEnabled),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Topics_ListSharedAccessKeys
	fmt.Println("Call operation: Topics_ListSharedAccessKeys")
	_, err = topicsClient.ListSharedAccessKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, nil)
	testsuite.Require().NoError(err)

	// From step Topics_RegenerateKey
	fmt.Println("Call operation: Topics_RegenerateKey")
	topicsClientRegenerateKeyResponsePoller, err := topicsClient.BeginRegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, armeventgrid.TopicRegenerateKeyRequest{
		KeyName: to.Ptr("key1"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicsClientRegenerateKeyResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/topics/eventSubscriptions
func (testsuite *EventGridTestSuite) TestTopicEventSubscriptions() {
	var err error
	// From step TopicEventSubscriptions_CreateOrUpdate
	fmt.Println("Call operation: TopicEventSubscriptions_CreateOrUpdate")
	topicEventSubscriptionsClient, err := armeventgrid.NewTopicEventSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicEventSubscriptionsClientCreateOrUpdateResponsePoller, err := topicEventSubscriptionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, testsuite.eventSubscriptionName, armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.EventHubEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeEventHub),
				Properties: &armeventgrid.EventHubEventSubscriptionDestinationProperties{
					ResourceID: to.Ptr(testsuite.eventhubId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicEventSubscriptionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TopicEventSubscriptions_List
	fmt.Println("Call operation: TopicEventSubscriptions_List")
	topicEventSubscriptionsClientNewListPager := topicEventSubscriptionsClient.NewListPager(testsuite.resourceGroupName, testsuite.topicName, &armeventgrid.TopicEventSubscriptionsClientListOptions{Filter: nil,
		Top: nil,
	})
	for topicEventSubscriptionsClientNewListPager.More() {
		_, err := topicEventSubscriptionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TopicEventSubscriptions_Get
	fmt.Println("Call operation: TopicEventSubscriptions_Get")
	_, err = topicEventSubscriptionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step TopicEventSubscriptions_Update
	fmt.Println("Call operation: TopicEventSubscriptions_Update")
	topicEventSubscriptionsClientUpdateResponsePoller, err := topicEventSubscriptionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, testsuite.eventSubscriptionName, armeventgrid.EventSubscriptionUpdateParameters{
		Labels: []*string{
			to.Ptr("label1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicEventSubscriptionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TopicEventSubscriptions_GetDeliveryAttributes
	fmt.Println("Call operation: TopicEventSubscriptions_GetDeliveryAttributes")
	_, err = topicEventSubscriptionsClient.GetDeliveryAttributes(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step TopicEventSubscriptions_Delete
	fmt.Println("Call operation: TopicEventSubscriptions_Delete")
	topicEventSubscriptionsClientDeleteResponsePoller, err := topicEventSubscriptionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicEventSubscriptionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/partnerRegistrations
func (testsuite *EventGridTestSuite) TestPartnerRegistrations() {
	var err error
	// From step PartnerRegistrations_ListBySubscription
	fmt.Println("Call operation: PartnerRegistrations_ListBySubscription")
	partnerRegistrationsClient, err := armeventgrid.NewPartnerRegistrationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	partnerRegistrationsClientNewListBySubscriptionPager := partnerRegistrationsClient.NewListBySubscriptionPager(&armeventgrid.PartnerRegistrationsClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for partnerRegistrationsClientNewListBySubscriptionPager.More() {
		_, err := partnerRegistrationsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerRegistrations_Get
	fmt.Println("Call operation: PartnerRegistrations_Get")
	_, err = partnerRegistrationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerRegistrationName, nil)
	testsuite.Require().NoError(err)

	// From step PartnerRegistrations_ListByResourceGroup
	fmt.Println("Call operation: PartnerRegistrations_ListByResourceGroup")
	partnerRegistrationsClientNewListByResourceGroupPager := partnerRegistrationsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.PartnerRegistrationsClientListByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for partnerRegistrationsClientNewListByResourceGroupPager.More() {
		_, err := partnerRegistrationsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerRegistrations_Update
	fmt.Println("Call operation: PartnerRegistrations_Update")
	partnerRegistrationsClientUpdateResponsePoller, err := partnerRegistrationsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerRegistrationName, armeventgrid.PartnerRegistrationUpdateParameters{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerRegistrationsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/partnerNamespaces
func (testsuite *EventGridTestSuite) TestPartnerNamespaces() {
	var err error
	// From step PartnerNamespaces_CreateOrUpdate
	fmt.Println("Call operation: PartnerNamespaces_CreateOrUpdate")
	partnerNamespacesClient, err := armeventgrid.NewPartnerNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	partnerNamespacesClientCreateOrUpdateResponsePoller, err := partnerNamespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, armeventgrid.PartnerNamespace{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
		Properties: &armeventgrid.PartnerNamespaceProperties{
			PartnerRegistrationFullyQualifiedID: to.Ptr(testsuite.partnerRegistrationId),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerNamespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerNamespaces_ListBySubscription
	fmt.Println("Call operation: PartnerNamespaces_ListBySubscription")
	partnerNamespacesClientNewListBySubscriptionPager := partnerNamespacesClient.NewListBySubscriptionPager(&armeventgrid.PartnerNamespacesClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for partnerNamespacesClientNewListBySubscriptionPager.More() {
		_, err := partnerNamespacesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerNamespaces_ListByResourceGroup
	fmt.Println("Call operation: PartnerNamespaces_ListByResourceGroup")
	partnerNamespacesClientNewListByResourceGroupPager := partnerNamespacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.PartnerNamespacesClientListByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for partnerNamespacesClientNewListByResourceGroupPager.More() {
		_, err := partnerNamespacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerNamespaces_Get
	fmt.Println("Call operation: PartnerNamespaces_Get")
	_, err = partnerNamespacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, nil)
	testsuite.Require().NoError(err)

	// From step PartnerNamespaces_Update
	fmt.Println("Call operation: PartnerNamespaces_Update")
	partnerNamespacesClientUpdateResponsePoller, err := partnerNamespacesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, armeventgrid.PartnerNamespaceUpdateParameters{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerNamespacesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerNamespaces_ListSharedAccessKeys
	fmt.Println("Call operation: PartnerNamespaces_ListSharedAccessKeys")
	_, err = partnerNamespacesClient.ListSharedAccessKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, nil)
	testsuite.Require().NoError(err)

	// From step PartnerNamespaces_RegenerateKey
	fmt.Println("Call operation: PartnerNamespaces_RegenerateKey")
	_, err = partnerNamespacesClient.RegenerateKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, armeventgrid.PartnerNamespaceRegenerateKeyRequest{
		KeyName: to.Ptr("key1"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step PartnerNamespaces_Delete
	fmt.Println("Call operation: PartnerNamespaces_Delete")
	partnerNamespacesClientDeleteResponsePoller, err := partnerNamespacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerNamespaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerNamespacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/partnerConfigurations/default
func (testsuite *EventGridTestSuite) TestPartnerConfigurations() {
	var err error
	// From step PartnerConfigurations_CreateOrUpdate
	fmt.Println("Call operation: PartnerConfigurations_CreateOrUpdate")
	partnerConfigurationsClient, err := armeventgrid.NewPartnerConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	partnerConfigurationsClientCreateOrUpdateResponsePoller, err := partnerConfigurationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.PartnerConfiguration{
		Location: to.Ptr("global"),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerConfigurationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerConfigurations_ListBySubscription
	fmt.Println("Call operation: PartnerConfigurations_ListBySubscription")
	partnerConfigurationsClientNewListBySubscriptionPager := partnerConfigurationsClient.NewListBySubscriptionPager(&armeventgrid.PartnerConfigurationsClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for partnerConfigurationsClientNewListBySubscriptionPager.More() {
		_, err := partnerConfigurationsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerConfigurations_Get
	fmt.Println("Call operation: PartnerConfigurations_Get")
	_, err = partnerConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)

	// From step PartnerConfigurations_ListByResourceGroup
	fmt.Println("Call operation: PartnerConfigurations_ListByResourceGroup")
	partnerConfigurationsClientNewListByResourceGroupPager := partnerConfigurationsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for partnerConfigurationsClientNewListByResourceGroupPager.More() {
		_, err := partnerConfigurationsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step PartnerConfigurations_Update
	fmt.Println("Call operation: PartnerConfigurations_Update")
	partnerConfigurationsClientUpdateResponsePoller, err := partnerConfigurationsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.PartnerConfigurationUpdateParameters{
		Properties: &armeventgrid.PartnerConfigurationUpdateParameterProperties{
			DefaultMaximumExpirationTimeInDays: to.Ptr[int32](100),
		},
		Tags: map[string]*string{
			"tag1": to.Ptr("value11"),
			"tag2": to.Ptr("value22"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerConfigurationsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerConfigurations_AuthorizePartner
	fmt.Println("Call operation: PartnerConfigurations_AuthorizePartner")
	_, err = partnerConfigurationsClient.AuthorizePartner(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.Partner{
		PartnerName: to.Ptr("Auth0"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step PartnerConfigurations_UnauthorizePartner
	fmt.Println("Call operation: PartnerConfigurations_UnauthorizePartner")
	_, err = partnerConfigurationsClient.AuthorizePartner(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.Partner{
		PartnerName: to.Ptr("Auth0"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step PartnerConfigurations_Delete
	fmt.Println("Call operation: PartnerConfigurations_Delete")
	partnerConfigurationsClientDeleteResponsePoller, err := partnerConfigurationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerConfigurationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/systemTopics
func (testsuite *EventGridTestSuite) TestSystemTopics() {
	var err error
	// From step SystemTopics_ListBySubscription
	fmt.Println("Call operation: SystemTopics_ListBySubscription")
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	systemTopicsClientNewListBySubscriptionPager := systemTopicsClient.NewListBySubscriptionPager(&armeventgrid.SystemTopicsClientListBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for systemTopicsClientNewListBySubscriptionPager.More() {
		_, err := systemTopicsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SystemTopics_ListByResourceGroup
	fmt.Println("Call operation: SystemTopics_ListByResourceGroup")
	systemTopicsClientNewListByResourceGroupPager := systemTopicsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.SystemTopicsClientListByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for systemTopicsClientNewListByResourceGroupPager.More() {
		_, err := systemTopicsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SystemTopics_Get
	fmt.Println("Call operation: SystemTopics_Get")
	_, err = systemTopicsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, nil)
	testsuite.Require().NoError(err)

	// From step SystemTopics_Update
	fmt.Println("Call operation: SystemTopics_Update")
	systemTopicsClientUpdateResponsePoller, err := systemTopicsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, armeventgrid.SystemTopicUpdateParameters{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/systemTopics/eventSubscriptions
func (testsuite *EventGridTestSuite) TestSystemTopicEventSubscriptions() {
	var err error
	// From step SystemTopicEventSubscriptions_CreateOrUpdate
	fmt.Println("Call operation: SystemTopicEventSubscriptions_CreateOrUpdate")
	systemTopicEventSubscriptionsClient, err := armeventgrid.NewSystemTopicEventSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	systemTopicEventSubscriptionsClientCreateOrUpdateResponsePoller, err := systemTopicEventSubscriptionsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, testsuite.eventSubscriptionName, armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.EventHubEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeEventHub),
				Properties: &armeventgrid.EventHubEventSubscriptionDestinationProperties{
					ResourceID: to.Ptr(testsuite.eventhubId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicEventSubscriptionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SystemTopicEventSubscriptions_ListBySystemTopic
	fmt.Println("Call operation: SystemTopicEventSubscriptions_ListBySystemTopic")
	systemTopicEventSubscriptionsClientNewListBySystemTopicPager := systemTopicEventSubscriptionsClient.NewListBySystemTopicPager(testsuite.resourceGroupName, testsuite.systemTopicName, &armeventgrid.SystemTopicEventSubscriptionsClientListBySystemTopicOptions{Filter: nil,
		Top: nil,
	})
	for systemTopicEventSubscriptionsClientNewListBySystemTopicPager.More() {
		_, err := systemTopicEventSubscriptionsClientNewListBySystemTopicPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SystemTopicEventSubscriptions_Get
	fmt.Println("Call operation: SystemTopicEventSubscriptions_Get")
	_, err = systemTopicEventSubscriptionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step SystemTopicEventSubscriptions_Update
	fmt.Println("Call operation: SystemTopicEventSubscriptions_Update")
	systemTopicEventSubscriptionsClientUpdateResponsePoller, err := systemTopicEventSubscriptionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, testsuite.eventSubscriptionName, armeventgrid.EventSubscriptionUpdateParameters{
		Labels: []*string{
			to.Ptr("label1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicEventSubscriptionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SystemTopicEventSubscriptions_GetDeliveryAttributes
	fmt.Println("Call operation: SystemTopicEventSubscriptions_GetDeliveryAttributes")
	_, err = systemTopicEventSubscriptionsClient.GetDeliveryAttributes(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step SystemTopicEventSubscriptions_Delete
	fmt.Println("Call operation: SystemTopicEventSubscriptions_Delete")
	systemTopicEventSubscriptionsClientDeleteResponsePoller, err := systemTopicEventSubscriptionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicEventSubscriptionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/eventSubscriptions
func (testsuite *EventGridTestSuite) TestEventSubscriptions() {
	var err error
	// From step EventSubscriptions_CreateOrUpdate
	fmt.Println("Call operation: EventSubscriptions_CreateOrUpdate")
	eventSubscriptionsClient, err := armeventgrid.NewEventSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	eventSubscriptionsClientCreateOrUpdateResponsePoller, err := eventSubscriptionsClient.BeginCreateOrUpdate(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, testsuite.eventSubscriptionName, armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.EventHubEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeEventHub),
				Properties: &armeventgrid.EventHubEventSubscriptionDestinationProperties{
					ResourceID: to.Ptr(testsuite.eventhubId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, eventSubscriptionsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step EventSubscriptions_ListGlobalBySubscription
	fmt.Println("Call operation: EventSubscriptions_ListGlobalBySubscription")
	eventSubscriptionsClientNewListGlobalBySubscriptionPager := eventSubscriptionsClient.NewListGlobalBySubscriptionPager(&armeventgrid.EventSubscriptionsClientListGlobalBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListGlobalBySubscriptionPager.More() {
		_, err := eventSubscriptionsClientNewListGlobalBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListGlobalBySubscriptionForTopicType
	fmt.Println("Call operation: EventSubscriptions_ListGlobalBySubscriptionForTopicType")
	eventSubscriptionsClientNewListGlobalBySubscriptionForTopicTypePager := eventSubscriptionsClient.NewListGlobalBySubscriptionForTopicTypePager("Microsoft.Resources.Subscriptions", &armeventgrid.EventSubscriptionsClientListGlobalBySubscriptionForTopicTypeOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListGlobalBySubscriptionForTopicTypePager.More() {
		_, err := eventSubscriptionsClientNewListGlobalBySubscriptionForTopicTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListGlobalByResourceGroup
	fmt.Println("Call operation: EventSubscriptions_ListGlobalByResourceGroup")
	eventSubscriptionsClientNewListGlobalByResourceGroupPager := eventSubscriptionsClient.NewListGlobalByResourceGroupPager(testsuite.resourceGroupName, &armeventgrid.EventSubscriptionsClientListGlobalByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListGlobalByResourceGroupPager.More() {
		_, err := eventSubscriptionsClientNewListGlobalByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListRegionalByResourceGroup
	fmt.Println("Call operation: EventSubscriptions_ListRegionalByResourceGroup")
	eventSubscriptionsClientNewListRegionalByResourceGroupPager := eventSubscriptionsClient.NewListRegionalByResourceGroupPager(testsuite.resourceGroupName, testsuite.location, &armeventgrid.EventSubscriptionsClientListRegionalByResourceGroupOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListRegionalByResourceGroupPager.More() {
		_, err := eventSubscriptionsClientNewListRegionalByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListByDomainTopic
	fmt.Println("Call operation: EventSubscriptions_ListByDomainTopic")
	eventSubscriptionsClientNewListByDomainTopicPager := eventSubscriptionsClient.NewListByDomainTopicPager(testsuite.resourceGroupName, testsuite.domainName, testsuite.topicName, &armeventgrid.EventSubscriptionsClientListByDomainTopicOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListByDomainTopicPager.More() {
		_, err := eventSubscriptionsClientNewListByDomainTopicPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_Get
	fmt.Println("Call operation: EventSubscriptions_Get")
	_, err = eventSubscriptionsClient.Get(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step EventSubscriptions_ListRegionalBySubscription
	fmt.Println("Call operation: EventSubscriptions_ListRegionalBySubscription")
	eventSubscriptionsClientNewListRegionalBySubscriptionPager := eventSubscriptionsClient.NewListRegionalBySubscriptionPager(testsuite.location, &armeventgrid.EventSubscriptionsClientListRegionalBySubscriptionOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListRegionalBySubscriptionPager.More() {
		_, err := eventSubscriptionsClientNewListRegionalBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListRegionalBySubscriptionForTopicType
	fmt.Println("Call operation: EventSubscriptions_ListRegionalBySubscriptionForTopicType")
	eventSubscriptionsClientNewListRegionalBySubscriptionForTopicTypePager := eventSubscriptionsClient.NewListRegionalBySubscriptionForTopicTypePager(testsuite.location, "Microsoft.EventHub.namespaces", &armeventgrid.EventSubscriptionsClientListRegionalBySubscriptionForTopicTypeOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListRegionalBySubscriptionForTopicTypePager.More() {
		_, err := eventSubscriptionsClientNewListRegionalBySubscriptionForTopicTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListRegionalByResourceGroupForTopicType
	fmt.Println("Call operation: EventSubscriptions_ListRegionalByResourceGroupForTopicType")
	eventSubscriptionsClientNewListRegionalByResourceGroupForTopicTypePager := eventSubscriptionsClient.NewListRegionalByResourceGroupForTopicTypePager(testsuite.resourceGroupName, testsuite.location, "Microsoft.EventHub.namespaces", &armeventgrid.EventSubscriptionsClientListRegionalByResourceGroupForTopicTypeOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListRegionalByResourceGroupForTopicTypePager.More() {
		_, err := eventSubscriptionsClientNewListRegionalByResourceGroupForTopicTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListByResource
	fmt.Println("Call operation: EventSubscriptions_ListByResource")
	eventSubscriptionsClientNewListByResourcePager := eventSubscriptionsClient.NewListByResourcePager(testsuite.resourceGroupName, "Microsoft.EventGrid", "topics", testsuite.topicName, &armeventgrid.EventSubscriptionsClientListByResourceOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListByResourcePager.More() {
		_, err := eventSubscriptionsClientNewListByResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_ListGlobalByResourceGroupForTopicType
	fmt.Println("Call operation: EventSubscriptions_ListGlobalByResourceGroupForTopicType")
	eventSubscriptionsClientNewListGlobalByResourceGroupForTopicTypePager := eventSubscriptionsClient.NewListGlobalByResourceGroupForTopicTypePager(testsuite.resourceGroupName, "Microsoft.Resources.ResourceGroups", &armeventgrid.EventSubscriptionsClientListGlobalByResourceGroupForTopicTypeOptions{Filter: nil,
		Top: nil,
	})
	for eventSubscriptionsClientNewListGlobalByResourceGroupForTopicTypePager.More() {
		_, err := eventSubscriptionsClientNewListGlobalByResourceGroupForTopicTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventSubscriptions_Update
	fmt.Println("Call operation: EventSubscriptions_Update")
	eventSubscriptionsClientUpdateResponsePoller, err := eventSubscriptionsClient.BeginUpdate(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, testsuite.eventSubscriptionName, armeventgrid.EventSubscriptionUpdateParameters{
		Labels: []*string{
			to.Ptr("label1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, eventSubscriptionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step EventSubscriptions_GetDeliveryAttributes
	fmt.Println("Call operation: EventSubscriptions_GetDeliveryAttributes")
	_, err = eventSubscriptionsClient.GetDeliveryAttributes(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)

	// From step EventSubscriptions_Delete
	fmt.Println("Call operation: EventSubscriptions_Delete")
	eventSubscriptionsClientDeleteResponsePoller, err := eventSubscriptionsClient.BeginDelete(testsuite.ctx, "subscriptions/"+testsuite.subscriptionId, testsuite.eventSubscriptionName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, eventSubscriptionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/operations
func (testsuite *EventGridTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armeventgrid.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.EventGrid/topicTypes
func (testsuite *EventGridTestSuite) TestTopicTypes() {
	var err error
	// From step TopicTypes_List
	fmt.Println("Call operation: TopicTypes_List")
	topicTypesClient, err := armeventgrid.NewTopicTypesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicTypesClientNewListPager := topicTypesClient.NewListPager(nil)
	for topicTypesClientNewListPager.More() {
		_, err := topicTypesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TopicTypes_Get
	fmt.Println("Call operation: TopicTypes_Get")
	_, err = topicTypesClient.Get(testsuite.ctx, "Microsoft.Storage.StorageAccounts", nil)
	testsuite.Require().NoError(err)

	// From step TopicTypes_ListEventTypes
	fmt.Println("Call operation: TopicTypes_ListEventTypes")
	topicTypesClientNewListEventTypesPager := topicTypesClient.NewListEventTypesPager("Microsoft.Storage.StorageAccounts", nil)
	for topicTypesClientNewListEventTypesPager.More() {
		_, err := topicTypesClientNewListEventTypesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.EventGrid/verifiedPartners
func (testsuite *EventGridTestSuite) TestVerifiedPartners() {
	var err error
	// From step VerifiedPartners_List
	fmt.Println("Call operation: VerifiedPartners_List")
	verifiedPartnersClient, err := armeventgrid.NewVerifiedPartnersClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	verifiedPartnersClientNewListPager := verifiedPartnersClient.NewListPager(&armeventgrid.VerifiedPartnersClientListOptions{Filter: nil,
		Top: nil,
	})
	for verifiedPartnersClientNewListPager.More() {
		_, err := verifiedPartnersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VerifiedPartners_Get
	fmt.Println("Call operation: VerifiedPartners_Get")
	_, err = verifiedPartnersClient.Get(testsuite.ctx, "Auth0", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventGrid/{parentType}/{parentName}/privateEndpointConnections
func (testsuite *EventGridTestSuite) TestPrivateEndpointConnections() {
	parentName := testsuite.domainName
	parentType := "domains"
	var privateEndpointConnectionName string
	var err error
	// From step PrivateEndpoint_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]interface{}{
			"domainId": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.domainId,
			},
			"privateEndpointName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.privateEndpointName,
			},
			"virtualNetworksName": map[string]interface{}{
				"type":         "string",
				"defaultValue": testsuite.virtualNetworksName,
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2020-11-01",
				"location":   "eastus",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []interface{}{
							"10.0.0.0/16",
						},
					},
					"enableDdosProtection": false,
					"subnets": []interface{}{
						map[string]interface{}{
							"name": "default",
							"properties": map[string]interface{}{
								"addressPrefix":                     "10.0.0.0/24",
								"delegations":                       []interface{}{},
								"privateEndpointNetworkPolicies":    "Disabled",
								"privateLinkServiceNetworkPolicies": "Enabled",
							},
						},
					},
					"virtualNetworkPeerings": []interface{}{},
				},
			},
			map[string]interface{}{
				"name":       "[concat(parameters('privateEndpointName'), '-nic')]",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "eastus",
				"properties": map[string]interface{}{
					"dnsSettings": map[string]interface{}{
						"dnsServers": []interface{}{},
					},
					"enableIPForwarding": false,
					"ipConfigurations": []interface{}{
						map[string]interface{}{
							"name": "privateEndpointIpConfig.ab24488f-044e-43f0-b9d1-af1f04071719",
							"properties": map[string]interface{}{
								"primary":                   true,
								"privateIPAddress":          "10.0.0.4",
								"privateIPAddressVersion":   "IPv4",
								"privateIPAllocationMethod": "Dynamic",
								"subnet": map[string]interface{}{
									"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
								},
							},
						},
					},
				},
			},
			map[string]interface{}{
				"name":       "[concat(parameters('virtualNetworksName'), '/default')]",
				"type":       "Microsoft.Network/virtualNetworks/subnets",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworksName'))]",
				},
				"properties": map[string]interface{}{
					"addressPrefix":                     "10.0.0.0/24",
					"delegations":                       []interface{}{},
					"privateEndpointNetworkPolicies":    "Disabled",
					"privateLinkServiceNetworkPolicies": "Enabled",
				},
			},
			map[string]interface{}{
				"name":       "[parameters('privateEndpointName')]",
				"type":       "Microsoft.Network/privateEndpoints",
				"apiVersion": "2020-11-01",
				"dependsOn": []interface{}{
					"[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
				},
				"location": "eastus",
				"properties": map[string]interface{}{
					"customDnsConfigs":                    []interface{}{},
					"manualPrivateLinkServiceConnections": []interface{}{},
					"privateLinkServiceConnections": []interface{}{
						map[string]interface{}{
							"name": "[parameters('privateEndpointName')]",
							"properties": map[string]interface{}{
								"groupIds": []interface{}{
									"domain",
								},
								"privateLinkServiceConnectionState": map[string]interface{}{
									"description":     "Auto-Approved",
									"actionsRequired": "None",
									"status":          "Approved",
								},
								"privateLinkServiceId": "[parameters('domainId')]",
							},
						},
					},
					"subnet": map[string]interface{}{
						"id": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
					},
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "PrivateEndpoint_Create", &deployment)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_ListByResource
	fmt.Println("Call operation: PrivateEndpointConnections_ListByResource")
	privateEndpointConnectionsClient, err := armeventgrid.NewPrivateEndpointConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateEndpointConnectionsClientNewListByResourcePager := privateEndpointConnectionsClient.NewListByResourcePager(testsuite.resourceGroupName, armeventgrid.PrivateEndpointConnectionsParentType(parentType), parentName, &armeventgrid.PrivateEndpointConnectionsClientListByResourceOptions{Filter: nil,
		Top: nil,
	})
	for privateEndpointConnectionsClientNewListByResourcePager.More() {
		nextResult, err := privateEndpointConnectionsClientNewListByResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		privateEndpointConnectionName = *nextResult.Value[0].Name
		break
	}

	// From step PrivateEndpointConnections_Update
	fmt.Println("Call operation: PrivateEndpointConnections_Update")
	privateEndpointConnectionsClientUpdateResponsePoller, err := privateEndpointConnectionsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.PrivateEndpointConnectionsParentType(parentType), parentName, privateEndpointConnectionName, armeventgrid.PrivateEndpointConnection{
		Properties: &armeventgrid.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armeventgrid.ConnectionState{
				Description:     to.Ptr("approving connection"),
				ActionsRequired: to.Ptr("None"),
				Status:          to.Ptr(armeventgrid.PersistedConnectionStatusRejected),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, privateEndpointConnectionsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step PrivateEndpointConnections_Get
	fmt.Println("Call operation: PrivateEndpointConnections_Get")
	_, err = privateEndpointConnectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, armeventgrid.PrivateEndpointConnectionsParentType(parentType), parentName, privateEndpointConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step PrivateLinkResources_ListByResource
	fmt.Println("Call operation: PrivateLinkResources_ListByResource")
	privateLinkResourcesClient, err := armeventgrid.NewPrivateLinkResourcesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	privateLinkResourcesClientNewListByResourcePager := privateLinkResourcesClient.NewListByResourcePager(testsuite.resourceGroupName, parentType, parentName, &armeventgrid.PrivateLinkResourcesClientListByResourceOptions{Filter: nil,
		Top: nil,
	})
	for privateLinkResourcesClientNewListByResourcePager.More() {
		_, err := privateLinkResourcesClientNewListByResourcePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *EventGridTestSuite) Cleanup() {
	var err error
	// From step SystemTopics_Delete
	fmt.Println("Call operation: SystemTopics_Delete")
	systemTopicsClient, err := armeventgrid.NewSystemTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	systemTopicsClientDeleteResponsePoller, err := systemTopicsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.systemTopicName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, systemTopicsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step PartnerRegistrations_Delete
	fmt.Println("Call operation: PartnerRegistrations_Delete")
	partnerRegistrationsClient, err := armeventgrid.NewPartnerRegistrationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	partnerRegistrationsClientDeleteResponsePoller, err := partnerRegistrationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.partnerRegistrationName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, partnerRegistrationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Topics_Delete
	fmt.Println("Call operation: Topics_Delete")
	topicsClient, err := armeventgrid.NewTopicsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	topicsClientDeleteResponsePoller, err := topicsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.topicName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, topicsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
