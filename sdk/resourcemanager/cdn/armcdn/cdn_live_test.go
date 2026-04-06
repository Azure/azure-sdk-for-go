// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type CdnTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	customDomainName  string
	endpointName      string
	originGroupName   string
	originName        string
	profileName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *CdnTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.customDomainName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "customdoma", 16, false)
	testsuite.endpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "endpointna", 16, false)
	testsuite.originGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "origingrou", 16, false)
	testsuite.originName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "originname", 16, false)
	testsuite.profileName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "profilenam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *CdnTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestCdnTestSuite(t *testing.T) {
	suite.Run(t, new(CdnTestSuite))
}

func (testsuite *CdnTestSuite) Prepare() {
	var err error
	// From step Profiles_Create
	fmt.Println("Call operation: Profiles_Create")
	profilesClient, err := armcdn.NewProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	profilesClientCreateResponsePoller, err := profilesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, armcdn.Profile{
		Location: to.Ptr("global"),
		SKU: &armcdn.SKU{
			Name: to.Ptr(armcdn.SKUNameStandardMicrosoft),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, profilesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Endpoints_Create
	fmt.Println("Call operation: Endpoints_Create")
	endpointsClient, err := armcdn.NewEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	endpointsClientCreateResponsePoller, err := endpointsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.Endpoint{
		Location: to.Ptr(testsuite.location),
		Properties: &armcdn.EndpointProperties{
			ContentTypesToCompress: []*string{
				to.Ptr("text/html"),
				to.Ptr("application/octet-stream")},
			IsCompressionEnabled:       to.Ptr(true),
			IsHTTPAllowed:              to.Ptr(true),
			IsHTTPSAllowed:             to.Ptr(true),
			OriginHostHeader:           to.Ptr("www.bing.com"),
			OriginPath:                 to.Ptr("/photos"),
			QueryStringCachingBehavior: to.Ptr(armcdn.QueryStringCachingBehaviorBypassCaching),
			Origins: []*armcdn.DeepCreatedOrigin{
				{
					Name: to.Ptr("origin1"),
					Properties: &armcdn.DeepCreatedOriginProperties{
						Enabled:   to.Ptr(true),
						HostName:  to.Ptr("www.someDomain1.net"),
						HTTPPort:  to.Ptr[int32](80),
						HTTPSPort: to.Ptr[int32](443),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step OriginGroups_Create
	fmt.Println("Call operation: OriginGroups_Create")
	originGroupsClient, err := armcdn.NewOriginGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	originGroupsClientCreateResponsePoller, err := originGroupsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originGroupName, armcdn.OriginGroup{
		Properties: &armcdn.OriginGroupProperties{
			HealthProbeSettings: &armcdn.HealthProbeParameters{
				ProbeIntervalInSeconds: to.Ptr[int32](120),
				ProbePath:              to.Ptr("/health.aspx"),
				ProbeProtocol:          to.Ptr(armcdn.ProbeProtocolHTTP),
				ProbeRequestType:       to.Ptr(armcdn.HealthProbeRequestTypeGET),
			},
			Origins: []*armcdn.ResourceReference{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/endpoints/" + testsuite.endpointName + "/origins/origin1"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, originGroupsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Endpoints_Update
	fmt.Println("Call operation: Endpoints_Update")
	endpointsClientUpdateResponsePoller, err := endpointsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.EndpointUpdateParameters{
		Properties: &armcdn.EndpointPropertiesUpdateParameters{
			DefaultOriginGroup: &armcdn.ResourceReference{
				ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourcegroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/endpoints/" + testsuite.endpointName + "/originGroups/" + testsuite.originGroupName),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Origins_Create
	fmt.Println("Call operation: Origins_Create")
	originsClient, err := armcdn.NewOriginsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	originsClientCreateResponsePoller, err := originsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originName, armcdn.Origin{
		Properties: &armcdn.OriginProperties{
			Enabled:   to.Ptr(true),
			HostName:  to.Ptr("www.someDomain.net"),
			HTTPPort:  to.Ptr[int32](80),
			HTTPSPort: to.Ptr[int32](443),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, originsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles
func (testsuite *CdnTestSuite) TestProfiles() {
	var err error
	// From step Profiles_List
	fmt.Println("Call operation: Profiles_List")
	profilesClient, err := armcdn.NewProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	profilesClientNewListPager := profilesClient.NewListPager(nil)
	for profilesClientNewListPager.More() {
		_, err := profilesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Profiles_ListByResourceGroup
	fmt.Println("Call operation: Profiles_ListByResourceGroup")
	profilesClientNewListByResourceGroupPager := profilesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for profilesClientNewListByResourceGroupPager.More() {
		_, err := profilesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Profiles_Get
	fmt.Println("Call operation: Profiles_Get")
	_, err = profilesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, nil)
	testsuite.Require().NoError(err)

	// From step Profiles_Update
	fmt.Println("Call operation: Profiles_Update")
	profilesClientUpdateResponsePoller, err := profilesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, armcdn.ProfileUpdateParameters{
		Tags: map[string]*string{
			"additionalProperties": to.Ptr("Tag1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, profilesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Profiles_ListResourceUsage
	fmt.Println("Call operation: Profiles_ListResourceUsage")
	profilesClientNewListResourceUsagePager := profilesClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for profilesClientNewListResourceUsagePager.More() {
		_, err := profilesClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Profiles_ListSupportedOptimizationTypes
	fmt.Println("Call operation: Profiles_ListSupportedOptimizationTypes")
	_, err = profilesClient.ListSupportedOptimizationTypes(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/endpoints
func (testsuite *CdnTestSuite) TestEndpoints() {
	var err error
	// From step Endpoints_ListByProfile
	fmt.Println("Call operation: Endpoints_ListByProfile")
	endpointsClient, err := armcdn.NewEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	endpointsClientNewListByProfilePager := endpointsClient.NewListByProfilePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for endpointsClientNewListByProfilePager.More() {
		_, err := endpointsClientNewListByProfilePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Endpoints_Get
	fmt.Println("Call operation: Endpoints_Get")
	_, err = endpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_ListResourceUsage
	fmt.Println("Call operation: Endpoints_ListResourceUsage")
	endpointsClientNewListResourceUsagePager := endpointsClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	for endpointsClientNewListResourceUsagePager.More() {
		_, err := endpointsClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Endpoints_ValidateCustomDomain
	fmt.Println("Call operation: Endpoints_ValidateCustomDomain")
	_, err = endpointsClient.ValidateCustomDomain(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.ValidateCustomDomainInput{
		HostName: to.Ptr("www.someDomain.com"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Endpoints_PurgeContent
	fmt.Println("Call operation: Endpoints_PurgeContent")
	endpointsClientPurgeContentResponsePoller, err := endpointsClient.BeginPurgeContent(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.PurgeParameters{
		ContentPaths: []*string{
			to.Ptr("/folder1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientPurgeContentResponsePoller)
	testsuite.Require().NoError(err)

	// From step Endpoints_Stop
	fmt.Println("Call operation: Endpoints_Stop")
	endpointsClientStopResponsePoller, err := endpointsClient.BeginStop(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientStopResponsePoller)
	testsuite.Require().NoError(err)

	// From step Endpoints_Start
	fmt.Println("Call operation: Endpoints_Start")
	endpointsClientStartResponsePoller, err := endpointsClient.BeginStart(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientStartResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/endpoints/origins
func (testsuite *CdnTestSuite) TestOrigins() {
	var err error
	// From step Origins_ListByEndpoint
	fmt.Println("Call operation: Origins_ListByEndpoint")
	originsClient, err := armcdn.NewOriginsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	originsClientNewListByEndpointPager := originsClient.NewListByEndpointPager(testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	for originsClientNewListByEndpointPager.More() {
		_, err := originsClientNewListByEndpointPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Origins_Get
	fmt.Println("Call operation: Origins_Get")
	_, err = originsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originName, nil)
	testsuite.Require().NoError(err)

	// From step Origins_Update
	fmt.Println("Call operation: Origins_Update")
	originsClientUpdateResponsePoller, err := originsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originName, armcdn.OriginUpdateParameters{
		Properties: &armcdn.OriginUpdatePropertiesParameters{
			Enabled:          to.Ptr(true),
			HTTPPort:         to.Ptr[int32](42),
			HTTPSPort:        to.Ptr[int32](43),
			OriginHostHeader: to.Ptr("www.someDomain2.net"),
			Priority:         to.Ptr[int32](1),
			Weight:           to.Ptr[int32](50),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, originsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Origins_Delete
	fmt.Println("Call operation: Origins_Delete")
	originsClientDeleteResponsePoller, err := originsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, originsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/endpoints/originGroups
func (testsuite *CdnTestSuite) TestOriginGroups() {
	var err error
	// From step OriginGroups_ListByEndpoint
	fmt.Println("Call operation: OriginGroups_ListByEndpoint")
	originGroupsClient, err := armcdn.NewOriginGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	originGroupsClientNewListByEndpointPager := originGroupsClient.NewListByEndpointPager(testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	for originGroupsClientNewListByEndpointPager.More() {
		_, err := originGroupsClientNewListByEndpointPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step OriginGroups_Get
	fmt.Println("Call operation: OriginGroups_Get")
	_, err = originGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originGroupName, nil)
	testsuite.Require().NoError(err)

	// From step OriginGroups_Update
	fmt.Println("Call operation: OriginGroups_Update")
	originGroupsClientUpdateResponsePoller, err := originGroupsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.originGroupName, armcdn.OriginGroupUpdateParameters{
		Properties: &armcdn.OriginGroupUpdatePropertiesParameters{
			HealthProbeSettings: &armcdn.HealthProbeParameters{
				ProbeIntervalInSeconds: to.Ptr[int32](120),
				ProbePath:              to.Ptr("/health.aspx"),
				ProbeProtocol:          to.Ptr(armcdn.ProbeProtocolHTTP),
				ProbeRequestType:       to.Ptr(armcdn.HealthProbeRequestTypeGET),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, originGroupsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/checkNameAvailability
func (testsuite *CdnTestSuite) TestCheckNameAvailabilityWithSubscription() {
	var err error
	// From step CheckNameAvailability
	fmt.Println("Call operation: CheckNameAvailability")
	managementClient, err := armcdn.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementClient.CheckNameAvailability(testsuite.ctx, armcdn.CheckNameAvailabilityInput{
		Name: to.Ptr("sampleName"),
		Type: to.Ptr(armcdn.ResourceTypeMicrosoftCdnProfilesEndpoints),
	}, nil)
	testsuite.Require().NoError(err)

	// From step CheckNameAvailabilityWithSubscription
	fmt.Println("Call operation: CheckNameAvailabilityWithSubscription")
	_, err = managementClient.CheckNameAvailabilityWithSubscription(testsuite.ctx, armcdn.CheckNameAvailabilityInput{
		Name: to.Ptr("sampleName"),
		Type: to.Ptr(armcdn.ResourceTypeMicrosoftCdnProfilesEndpoints),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/validateProbe
func (testsuite *CdnTestSuite) TestValidateProbe() {
	var err error
	// From step ValidateProbe
	fmt.Println("Call operation: ValidateProbe")
	managementClient, err := armcdn.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementClient.ValidateProbe(testsuite.ctx, armcdn.ValidateProbeInput{
		ProbeURL: to.Ptr("https://www.bing.com/image"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/checkResourceUsage
func (testsuite *CdnTestSuite) TestResourceUsage() {
	var err error
	// From step ResourceUsage_List
	fmt.Println("Call operation: ResourceUsage_List")
	resourceUsageClient, err := armcdn.NewResourceUsageClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	resourceUsageClientNewListPager := resourceUsageClient.NewListPager(nil)
	for resourceUsageClientNewListPager.More() {
		_, err := resourceUsageClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/operations
func (testsuite *CdnTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armcdn.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/edgenodes
func (testsuite *CdnTestSuite) TestEdgeNodes() {
	var err error
	// From step EdgeNodes_List
	fmt.Println("Call operation: EdgeNodes_List")
	edgeNodesClient, err := armcdn.NewEdgeNodesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	edgeNodesClientNewListPager := edgeNodesClient.NewListPager(nil)
	for edgeNodesClientNewListPager.More() {
		_, err := edgeNodesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/cdnWebApplicationFirewallManagedRuleSets
func (testsuite *CdnTestSuite) TestManagedRuleSets() {
	var err error
	// From step ManagedRuleSets_List
	fmt.Println("Call operation: ManagedRuleSets_List")
	managedRuleSetsClient, err := armcdn.NewManagedRuleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managedRuleSetsClientNewListPager := managedRuleSetsClient.NewListPager(nil)
	for managedRuleSetsClientNewListPager.More() {
		_, err := managedRuleSetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *CdnTestSuite) Cleanup() {
	var err error
	// From step Endpoints_Delete
	fmt.Println("Call operation: Endpoints_Delete")
	endpointsClient, err := armcdn.NewEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	endpointsClientDeleteResponsePoller, err := endpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, endpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Profiles_Delete
	fmt.Println("Call operation: Endpoints_Delete")
	profilesClient, err := armcdn.NewProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	profilesClientDeleteResponsePoller, err := profilesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, profilesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
