// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v4"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	mobileNetworkName string
	serviceName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ServiceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicen", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ServiceTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (testsuite *ServiceTestSuite) Prepare() {
	var err error
	// From step MobileNetworks_CreateOrUpdate
	fmt.Println("Call operation: MobileNetworks_CreateOrUpdate")
	mobileNetworksClient, err := armmobilenetwork.NewMobileNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	mobileNetworksClientCreateOrUpdateResponsePoller, err := mobileNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, armmobilenetwork.MobileNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.PropertiesFormat{
			PublicLandMobileNetworkIdentifier: &armmobilenetwork.PlmnID{
				Mcc: to.Ptr("001"),
				Mnc: to.Ptr("01"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, mobileNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.MobileNetwork/mobileNetworks/{mobileNetworkName}/services/{serviceName}
func (testsuite *ServiceTestSuite) TestServices() {
	var err error
	// From step Services_CreateOrUpdate
	fmt.Println("Call operation: Services_CreateOrUpdate")
	servicesClient, err := armmobilenetwork.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	servicesClientCreateOrUpdateResponsePoller, err := servicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.serviceName, armmobilenetwork.Service{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.ServicePropertiesFormat{
			PccRules: []*armmobilenetwork.PccRuleConfiguration{
				{
					RuleName:       to.Ptr("default-rule"),
					RulePrecedence: to.Ptr[int32](255),
					RuleQosPolicy: &armmobilenetwork.PccRuleQosPolicy{
						FiveQi:                              to.Ptr[int32](9),
						AllocationAndRetentionPriorityLevel: to.Ptr[int32](9),
						MaximumBitRate: &armmobilenetwork.Ambr{
							Downlink: to.Ptr("1 Gbps"),
							Uplink:   to.Ptr("500 Mbps"),
						},
						PreemptionCapability:    to.Ptr(armmobilenetwork.PreemptionCapabilityNotPreempt),
						PreemptionVulnerability: to.Ptr(armmobilenetwork.PreemptionVulnerabilityPreemptable),
					},
					ServiceDataFlowTemplates: []*armmobilenetwork.ServiceDataFlowTemplate{
						{
							Direction: to.Ptr(armmobilenetwork.SdfDirectionUplink),
							Ports:     []*string{},
							RemoteIPList: []*string{
								to.Ptr("10.3.4.0/24")},
							TemplateName: to.Ptr("IP-to-server"),
							Protocol: []*string{
								to.Ptr("ip")},
						}},
					TrafficControl: to.Ptr(armmobilenetwork.TrafficControlPermissionEnabled),
				}},
			ServicePrecedence: to.Ptr[int32](255),
			ServiceQosPolicy: &armmobilenetwork.QosPolicy{
				FiveQi:                              to.Ptr[int32](9),
				AllocationAndRetentionPriorityLevel: to.Ptr[int32](9),
				MaximumBitRate: &armmobilenetwork.Ambr{
					Downlink: to.Ptr("1 Gbps"),
					Uplink:   to.Ptr("500 Mbps"),
				},
				PreemptionCapability:    to.Ptr(armmobilenetwork.PreemptionCapabilityNotPreempt),
				PreemptionVulnerability: to.Ptr(armmobilenetwork.PreemptionVulnerabilityPreemptable),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Services_ListByMobileNetwork
	fmt.Println("Call operation: Services_ListByMobileNetwork")
	servicesClientNewListByMobileNetworkPager := servicesClient.NewListByMobileNetworkPager(testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	for servicesClientNewListByMobileNetworkPager.More() {
		_, err := servicesClientNewListByMobileNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Services_Get
	fmt.Println("Call operation: Services_Get")
	_, err = servicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step Services_UpdateTags
	fmt.Println("Call operation: Services_UpdateTags")
	_, err = servicesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.serviceName, armmobilenetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Services_Delete
	fmt.Println("Call operation: Services_Delete")
	servicesClientDeleteResponsePoller, err := servicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
