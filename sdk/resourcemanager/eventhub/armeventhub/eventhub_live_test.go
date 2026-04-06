// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armeventhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type EventhubTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	applicationGroupName  string
	authorizationRuleName string
	consumerGroupName     string
	eventHubName          string
	namespaceName         string
	schemaGroupName       string
	storageAccountId      string
	storageAccountName    string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *EventhubTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.applicationGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "applicatio", 16, false)
	testsuite.authorizationRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "authorizat", 16, false)
	testsuite.consumerGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "consumergr", 16, false)
	testsuite.eventHubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "eventhubna", 16, false)
	testsuite.namespaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "namespacen", 16, false)
	testsuite.schemaGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "schemagrou", 16, false)
	testsuite.storageAccountName = "storageeventhub2"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *EventhubTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestEventhubTestSuite(t *testing.T) {
	suite.Run(t, new(EventhubTestSuite))
}

func (testsuite *EventhubTestSuite) Prepare() {
	var err error
	// From step Namespaces_CreateOrUpdate
	fmt.Println("Call operation: Namespaces_CreateOrUpdate")
	namespacesClient, err := armeventhub.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientCreateOrUpdateResponsePoller, err := namespacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armeventhub.EHNamespace{
		Location: to.Ptr(testsuite.location),
		SKU: &armeventhub.SKU{
			Name: to.Ptr(armeventhub.SKUNamePremium),
			Tier: to.Ptr(armeventhub.SKUTierPremium),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step StorageAccount_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"storageAccountId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"storageAccountName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(storageAccountName)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "eastus",
				"properties": map[string]interface{}{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]interface{}{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]interface{}{
							"blob": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]interface{}{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]interface{}{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []interface{}{},
						"virtualNetworkRules": []interface{}{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]interface{}{
					"name": "Standard_RAGRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params := map[string]interface{}{
		"storageAccountName": map[string]interface{}{"value": testsuite.storageAccountName},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "StorageAccount_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step EventHubs_CreateOrUpdate
	fmt.Println("Call operation: EventHubs_CreateOrUpdate")
	eventHubsClient, err := armeventhub.NewEventHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = eventHubsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, armeventhub.Eventhub{
		Properties: &armeventhub.Properties{
			CaptureDescription: &armeventhub.CaptureDescription{
				Destination: &armeventhub.Destination{
					Name: to.Ptr("EventHubArchive.AzureBlockBlob"),
					Properties: &armeventhub.DestinationProperties{
						ArchiveNameFormat:        to.Ptr("{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}"),
						BlobContainer:            to.Ptr("container"),
						StorageAccountResourceID: to.Ptr(testsuite.storageAccountId),
					},
				},
				Enabled:           to.Ptr(true),
				Encoding:          to.Ptr(armeventhub.EncodingCaptureDescriptionAvro),
				IntervalInSeconds: to.Ptr[int32](120),
				SizeLimitInBytes:  to.Ptr[int32](10485763),
			},
			MessageRetentionInDays: to.Ptr[int64](4),
			PartitionCount:         to.Ptr[int64](4),
			Status:                 to.Ptr(armeventhub.EntityStatusActive),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/namespaces
func (testsuite *EventhubTestSuite) TTestNamespace() {
	var err error
	// From step Namespaces_CheckNameAvailability
	fmt.Println("Call operation: Namespaces_CheckNameAvailability")
	namespacesClient, err := armeventhub.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = namespacesClient.CheckNameAvailability(testsuite.ctx, armeventhub.CheckNameAvailabilityParameter{
		Name: to.Ptr("sdk-Namespace-8458"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_List
	fmt.Println("Call operation: Namespaces_List")
	namespacesClientNewListPager := namespacesClient.NewListPager(nil)
	for namespacesClientNewListPager.More() {
		_, err := namespacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_ListAuthorizationRules
	fmt.Println("Call operation: Namespaces_ListAuthorizationRules")
	namespacesClientNewListAuthorizationRulesPager := namespacesClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for namespacesClientNewListAuthorizationRulesPager.More() {
		_, err := namespacesClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_ListNetworkRuleSet
	fmt.Println("Call operation: Namespaces_ListNetworkRuleSet")
	_, err = namespacesClient.ListNetworkRuleSet(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_GetNetworkRuleSet
	fmt.Println("Call operation: Namespaces_GetNetworkRuleSet")
	_, err = namespacesClient.GetNetworkRuleSet(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListByResourceGroup
	fmt.Println("Call operation: Namespaces_ListByResourceGroup")
	namespacesClientNewListByResourceGroupPager := namespacesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for namespacesClientNewListByResourceGroupPager.More() {
		_, err := namespacesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Namespaces_Get
	fmt.Println("Call operation: Namespaces_Get")
	_, err = namespacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_Update
	fmt.Println("Call operation: Namespaces_Update")
	_, err = namespacesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, armeventhub.EHNamespace{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: Namespaces_CreateOrUpdateAuthorizationRule")
	_, err = namespacesClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armeventhub.AuthorizationRule{
		Properties: &armeventhub.AuthorizationRuleProperties{
			Rights: []*armeventhub.AccessRights{
				to.Ptr(armeventhub.AccessRightsListen),
				to.Ptr(armeventhub.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_GetAuthorizationRule
	fmt.Println("Call operation: Namespaces_GetAuthorizationRule")
	_, err = namespacesClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_ListKeys
	fmt.Println("Call operation: Namespaces_ListKeys")
	_, err = namespacesClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_RegenerateKeys
	fmt.Println("Call operation: Namespaces_RegenerateKeys")
	_, err = namespacesClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, armeventhub.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armeventhub.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_DeleteAuthorizationRule
	fmt.Println("Call operation: Namespaces_DeleteAuthorizationRule")
	_, err = namespacesClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/namespaces/eventhubs
func (testsuite *EventhubTestSuite) TTestEventhubs() {
	var err error
	// From step EventHubs_ListByNamespace
	fmt.Println("Call operation: EventHubs_ListByNamespace")
	eventHubsClient, err := armeventhub.NewEventHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	eventHubsClientNewListByNamespacePager := eventHubsClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, &armeventhub.EventHubsClientListByNamespaceOptions{Skip: nil,
		Top: nil,
	})
	for eventHubsClientNewListByNamespacePager.More() {
		_, err := eventHubsClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventHubs_ListAuthorizationRules
	fmt.Println("Call operation: EventHubs_ListAuthorizationRules")
	eventHubsClientNewListAuthorizationRulesPager := eventHubsClient.NewListAuthorizationRulesPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, nil)
	for eventHubsClientNewListAuthorizationRulesPager.More() {
		_, err := eventHubsClientNewListAuthorizationRulesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step EventHubs_Get
	fmt.Println("Call operation: EventHubs_Get")
	_, err = eventHubsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, nil)
	testsuite.Require().NoError(err)

	// From step EventHubs_CreateOrUpdateAuthorizationRule
	fmt.Println("Call operation: EventHubs_CreateOrUpdateAuthorizationRule")
	_, err = eventHubsClient.CreateOrUpdateAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.authorizationRuleName, armeventhub.AuthorizationRule{
		Properties: &armeventhub.AuthorizationRuleProperties{
			Rights: []*armeventhub.AccessRights{
				to.Ptr(armeventhub.AccessRightsListen),
				to.Ptr(armeventhub.AccessRightsSend)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step EventHubs_GetAuthorizationRule
	fmt.Println("Call operation: EventHubs_GetAuthorizationRule")
	_, err = eventHubsClient.GetAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step EventHubs_ListKeys
	fmt.Println("Call operation: EventHubs_ListKeys")
	_, err = eventHubsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)

	// From step EventHubs_RegenerateKeys
	fmt.Println("Call operation: EventHubs_RegenerateKeys")
	_, err = eventHubsClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.authorizationRuleName, armeventhub.RegenerateAccessKeyParameters{
		KeyType: to.Ptr(armeventhub.KeyTypePrimaryKey),
	}, nil)
	testsuite.Require().NoError(err)

	// From step EventHubs_DeleteAuthorizationRule
	fmt.Println("Call operation: EventHubs_DeleteAuthorizationRule")
	_, err = eventHubsClient.DeleteAuthorizationRule(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.authorizationRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/namespaces/eventhubs/consumergroups
func (testsuite *EventhubTestSuite) TTestConsumergroups() {
	var err error
	// From step ConsumerGroups_CreateOrUpdate
	fmt.Println("Call operation: ConsumerGroups_CreateOrUpdate")
	consumerGroupsClient, err := armeventhub.NewConsumerGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = consumerGroupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.consumerGroupName, armeventhub.ConsumerGroup{
		Properties: &armeventhub.ConsumerGroupProperties{
			UserMetadata: to.Ptr("New consumergroup"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ConsumerGroups_ListByEventHub
	fmt.Println("Call operation: ConsumerGroups_ListByEventHub")
	consumerGroupsClientNewListByEventHubPager := consumerGroupsClient.NewListByEventHubPager(testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, &armeventhub.ConsumerGroupsClientListByEventHubOptions{Skip: nil,
		Top: nil,
	})
	for consumerGroupsClientNewListByEventHubPager.More() {
		_, err := consumerGroupsClientNewListByEventHubPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConsumerGroups_Get
	fmt.Println("Call operation: ConsumerGroups_Get")
	_, err = consumerGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.consumerGroupName, nil)
	testsuite.Require().NoError(err)

	// From step ConsumerGroups_Delete
	fmt.Println("Call operation: ConsumerGroups_Delete")
	_, err = consumerGroupsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, testsuite.consumerGroupName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/namespaces/applicationGroups
// func (testsuite *EventhubTestSuite) TestApplicationgroup() {
// 	var err error
// 	// From step ApplicationGroup_CreateOrUpdateApplicationGroup
// 	fmt.Println("Call operation: ApplicationGroup_CreateOrUpdateApplicationGroup")
// 	applicationGroupClient, err := armeventhub.NewApplicationGroupClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
// 	testsuite.Require().NoError(err)
// 	_, err = applicationGroupClient.CreateOrUpdateApplicationGroup(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.applicationGroupName, armeventhub.ApplicationGroup{
// 		Properties: &armeventhub.ApplicationGroupProperties{
// 			ClientAppGroupIdentifier: to.Ptr("SASKeyName=KeyName"),
// 			IsEnabled:                to.Ptr(true),
// 			Policies: []armeventhub.ApplicationGroupPolicyClassification{
// 				&armeventhub.ThrottlingPolicy{
// 					Name:               to.Ptr("ThrottlingPolicy1"),
// 					Type:               to.Ptr(armeventhub.ApplicationGroupPolicyTypeThrottlingPolicy),
// 					MetricID:           to.Ptr(armeventhub.MetricIDIncomingMessages),
// 					RateLimitThreshold: to.Ptr[int64](7912),
// 				},
// 				&armeventhub.ThrottlingPolicy{
// 					Name:               to.Ptr("ThrottlingPolicy2"),
// 					Type:               to.Ptr(armeventhub.ApplicationGroupPolicyTypeThrottlingPolicy),
// 					MetricID:           to.Ptr(armeventhub.MetricIDIncomingBytes),
// 					RateLimitThreshold: to.Ptr[int64](3951729),
// 				},
// 				&armeventhub.ThrottlingPolicy{
// 					Name:               to.Ptr("ThrottlingPolicy3"),
// 					Type:               to.Ptr(armeventhub.ApplicationGroupPolicyTypeThrottlingPolicy),
// 					MetricID:           to.Ptr(armeventhub.MetricIDOutgoingBytes),
// 					RateLimitThreshold: to.Ptr[int64](245175),
// 				}},
// 		},
// 	}, nil)
// 	testsuite.Require().NoError(err)

// 	// From step ApplicationGroup_ListByNamespace
// 	fmt.Println("Call operation: ApplicationGroup_ListByNamespace")
// 	applicationGroupClientNewListByNamespacePager := applicationGroupClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
// 	for applicationGroupClientNewListByNamespacePager.More() {
// 		_, err := applicationGroupClientNewListByNamespacePager.NextPage(testsuite.ctx)
// 		testsuite.Require().NoError(err)
// 		break
// 	}

// 	// From step ApplicationGroup_Get
// 	fmt.Println("Call operation: ApplicationGroup_Get")
// 	_, err = applicationGroupClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.applicationGroupName, nil)
// 	testsuite.Require().NoError(err)

// 	// From step ApplicationGroup_Delete
// 	fmt.Println("Call operation: ApplicationGroup_Delete")
// 	_, err = applicationGroupClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.applicationGroupName, nil)
// 	testsuite.Require().NoError(err)
// }

// Microsoft.EventHub/namespaces/schemagroups
func (testsuite *EventhubTestSuite) TTestSchemaregistry() {
	var err error
	// From step SchemaRegistry_CreateOrUpdate
	fmt.Println("Call operation: SchemaRegistry_CreateOrUpdate")
	schemaRegistryClient, err := armeventhub.NewSchemaRegistryClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = schemaRegistryClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.schemaGroupName, armeventhub.SchemaGroup{
		Properties: &armeventhub.SchemaGroupProperties{
			GroupProperties:     map[string]*string{},
			SchemaCompatibility: to.Ptr(armeventhub.SchemaCompatibilityForward),
			SchemaType:          to.Ptr(armeventhub.SchemaTypeAvro),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SchemaRegistry_ListByNamespace
	fmt.Println("Call operation: SchemaRegistry_ListByNamespace")
	schemaRegistryClientNewListByNamespacePager := schemaRegistryClient.NewListByNamespacePager(testsuite.resourceGroupName, testsuite.namespaceName, &armeventhub.SchemaRegistryClientListByNamespaceOptions{Skip: nil,
		Top: nil,
	})
	for schemaRegistryClientNewListByNamespacePager.More() {
		_, err := schemaRegistryClientNewListByNamespacePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SchemaRegistry_Get
	fmt.Println("Call operation: SchemaRegistry_Get")
	_, err = schemaRegistryClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.schemaGroupName, nil)
	testsuite.Require().NoError(err)

	// From step SchemaRegistry_Delete
	fmt.Println("Call operation: SchemaRegistry_Delete")
	_, err = schemaRegistryClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.schemaGroupName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.EventHub/operations
func (testsuite *EventhubTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armeventhub.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *EventhubTestSuite) Cleanup() {
	var err error
	// From step EventHubs_Delete
	fmt.Println("Call operation: EventHubs_Delete")
	eventHubsClient, err := armeventhub.NewEventHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = eventHubsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.eventHubName, nil)
	testsuite.Require().NoError(err)

	// From step Namespaces_Delete
	fmt.Println("Call operation: Namespaces_Delete")
	namespacesClient, err := armeventhub.NewNamespacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	namespacesClientDeleteResponsePoller, err := namespacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, namespacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
