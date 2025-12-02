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

type SimPolicyTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	dataNetworkId     string
	dataNetworkName   string
	mobileNetworkName string
	serviceId         string
	serviceName       string
	simPolicyName     string
	sliceId           string
	sliceName         string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SimPolicyTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.dataNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "datanetw", 14, false)
	testsuite.mobileNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mobilene", 14, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicen", 14, false)
	testsuite.simPolicyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "simpolic", 14, false)
	testsuite.sliceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "slicenam", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *SimPolicyTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSimPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(SimPolicyTestSuite))
}

func (testsuite *SimPolicyTestSuite) Prepare() {
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

	// From step Slices_CreateOrUpdate
	fmt.Println("Call operation: Slices_CreateOrUpdate")
	slicesClient, err := armmobilenetwork.NewSlicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	slicesClientCreateOrUpdateResponsePoller, err := slicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.sliceName, armmobilenetwork.Slice{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SlicePropertiesFormat{
			Description: to.Ptr("myFavouriteSlice"),
			Snssai: &armmobilenetwork.Snssai{
				Sd:  to.Ptr("1abcde"),
				Sst: to.Ptr[int32](1),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var slicesClientCreateOrUpdateResponse *armmobilenetwork.SlicesClientCreateOrUpdateResponse
	slicesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, slicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.sliceId = *slicesClientCreateOrUpdateResponse.ID

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
	var servicesClientCreateOrUpdateResponse *armmobilenetwork.ServicesClientCreateOrUpdateResponse
	servicesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.serviceId = *servicesClientCreateOrUpdateResponse.ID

	// From step DataNetworks_CreateOrUpdate
	fmt.Println("Call operation: DataNetworks_CreateOrUpdate")
	dataNetworksClient, err := armmobilenetwork.NewDataNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dataNetworksClientCreateOrUpdateResponsePoller, err := dataNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.dataNetworkName, armmobilenetwork.DataNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.DataNetworkPropertiesFormat{
			Description: to.Ptr("myFavouriteDataNetwork"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var dataNetworksClientCreateOrUpdateResponse *armmobilenetwork.DataNetworksClientCreateOrUpdateResponse
	dataNetworksClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, dataNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.dataNetworkId = *dataNetworksClientCreateOrUpdateResponse.ID
}

// Microsoft.MobileNetwork/mobileNetworks/{mobileNetworkName}/simPolicies/{simPolicyName}
func (testsuite *SimPolicyTestSuite) TestSimPolicies() {
	var err error
	// From step SimPolicies_CreateOrUpdate
	fmt.Println("Call operation: SimPolicies_CreateOrUpdate")
	simPoliciesClient, err := armmobilenetwork.NewSimPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	simPoliciesClientCreateOrUpdateResponsePoller, err := simPoliciesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.simPolicyName, armmobilenetwork.SimPolicy{
		Location: to.Ptr(testsuite.location),
		Properties: &armmobilenetwork.SimPolicyPropertiesFormat{
			DefaultSlice: &armmobilenetwork.SliceResourceID{
				ID: to.Ptr(testsuite.sliceId),
			},
			RegistrationTimer: to.Ptr[int32](3240),
			SliceConfigurations: []*armmobilenetwork.SliceConfiguration{
				{
					DataNetworkConfigurations: []*armmobilenetwork.DataNetworkConfiguration{
						{
							FiveQi:                              to.Ptr[int32](9),
							AdditionalAllowedSessionTypes:       []*armmobilenetwork.PduSessionType{},
							AllocationAndRetentionPriorityLevel: to.Ptr[int32](9),
							AllowedServices: []*armmobilenetwork.ServiceResourceID{
								{
									ID: to.Ptr(testsuite.serviceId),
								}},
							DataNetwork: &armmobilenetwork.DataNetworkResourceID{
								ID: to.Ptr(testsuite.dataNetworkId),
							},
							DefaultSessionType:             to.Ptr(armmobilenetwork.PduSessionTypeIPv4),
							MaximumNumberOfBufferedPackets: to.Ptr[int32](200),
							PreemptionCapability:           to.Ptr(armmobilenetwork.PreemptionCapabilityNotPreempt),
							PreemptionVulnerability:        to.Ptr(armmobilenetwork.PreemptionVulnerabilityPreemptable),
							SessionAmbr: &armmobilenetwork.Ambr{
								Downlink: to.Ptr("1 Gbps"),
								Uplink:   to.Ptr("500 Mbps"),
							},
						}},
					DefaultDataNetwork: &armmobilenetwork.DataNetworkResourceID{
						ID: to.Ptr(testsuite.dataNetworkId),
					},
					Slice: &armmobilenetwork.SliceResourceID{
						ID: to.Ptr(testsuite.sliceId),
					},
				}},
			UeAmbr: &armmobilenetwork.Ambr{
				Downlink: to.Ptr("1 Gbps"),
				Uplink:   to.Ptr("500 Mbps"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simPoliciesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SimPolicies_ListByMobileNetwork
	fmt.Println("Call operation: SimPolicies_ListByMobileNetwork")
	simPoliciesClientNewListByMobileNetworkPager := simPoliciesClient.NewListByMobileNetworkPager(testsuite.resourceGroupName, testsuite.mobileNetworkName, nil)
	for simPoliciesClientNewListByMobileNetworkPager.More() {
		_, err := simPoliciesClientNewListByMobileNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SimPolicies_Get
	fmt.Println("Call operation: SimPolicies_Get")
	_, err = simPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.simPolicyName, nil)
	testsuite.Require().NoError(err)

	// From step SimPolicies_UpdateTags
	fmt.Println("Call operation: SimPolicies_UpdateTags")
	_, err = simPoliciesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.simPolicyName, armmobilenetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SimPolicies_Delete
	fmt.Println("Call operation: SimPolicies_Delete")
	simPoliciesClientDeleteResponsePoller, err := simPoliciesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.mobileNetworkName, testsuite.simPolicyName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, simPoliciesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
