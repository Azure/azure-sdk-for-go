// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcdn_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type AfdxTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	customDomainName  string
	endpointName      string
	originGroupName   string
	originName        string
	profileName       string
	routeName         string
	ruleName          string
	ruleSetName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AfdxTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.customDomainName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "afdcustomdoma", 19, false)
	testsuite.endpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "afdendpointna", 19, false)
	testsuite.originGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "afdorigingrou", 19, false)
	testsuite.originName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "afdoriginname", 19, false)
	testsuite.profileName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "afdprofilenam", 19, false)
	testsuite.routeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "routename", 15, false)
	testsuite.ruleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rulename", 14, false)
	testsuite.ruleSetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rulesetnam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AfdxTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAfdxTestSuite(t *testing.T) {
	suite.Run(t, new(AfdxTestSuite))
}

func (testsuite *AfdxTestSuite) Prepare() {
	var err error
	// From step Profiles_Create
	fmt.Println("Call operation: Profiles_Create")
	profilesClient, err := armcdn.NewProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	profilesClientCreateResponsePoller, err := profilesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, armcdn.Profile{
		Location: to.Ptr("global"),
		SKU: &armcdn.SKU{
			Name: to.Ptr(armcdn.SKUNamePremiumAzureFrontDoor),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, profilesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDCustomDomains_Create
	fmt.Println("Call operation: AFDCustomDomains_Create")
	aFDCustomDomainsClient, err := armcdn.NewAFDCustomDomainsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDCustomDomainsClientCreateResponsePoller, err := aFDCustomDomainsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.customDomainName, armcdn.AFDDomain{
		Properties: &armcdn.AFDDomainProperties{
			TLSSettings: &armcdn.AFDDomainHTTPSParameters{
				CertificateType:   to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
				MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
			},
			HostName: to.Ptr("www.someDomain.net"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDCustomDomainsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDEndpoints_Create
	fmt.Println("Call operation: AFDEndpoints_Create")
	aFDEndpointsClient, err := armcdn.NewAFDEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDEndpointsClientCreateResponsePoller, err := aFDEndpointsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.AFDEndpoint{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armcdn.AFDEndpointProperties{
			EnabledState:                      to.Ptr(armcdn.EnabledStateEnabled),
			AutoGeneratedDomainNameLabelScope: to.Ptr(armcdn.AutoGeneratedDomainNameLabelScopeTenantReuse),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDEndpointsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDOriginGroups_Create
	fmt.Println("Call operation: AFDOriginGroups_Create")
	aFDOriginGroupsClient, err := armcdn.NewAFDOriginGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginGroupsClientCreateResponsePoller, err := aFDOriginGroupsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, armcdn.AFDOriginGroup{
		Properties: &armcdn.AFDOriginGroupProperties{
			HealthProbeSettings: &armcdn.HealthProbeParameters{
				ProbeIntervalInSeconds: to.Ptr[int32](10),
				ProbePath:              to.Ptr("/path2"),
				ProbeProtocol:          to.Ptr(armcdn.ProbeProtocolNotSet),
				ProbeRequestType:       to.Ptr(armcdn.HealthProbeRequestTypeNotSet),
			},
			LoadBalancingSettings: &armcdn.LoadBalancingSettingsParameters{
				AdditionalLatencyInMilliseconds: to.Ptr[int32](1000),
				SampleSize:                      to.Ptr[int32](3),
				SuccessfulSamplesRequired:       to.Ptr[int32](3),
			},
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: to.Ptr[int32](5),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginGroupsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step RuleSets_Create
	fmt.Println("Call operation: RuleSets_Create")
	aFDOriginsClient, err := armcdn.NewAFDOriginsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginsClientCreateResponsePoller, err := aFDOriginsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, testsuite.originName, armcdn.AFDOrigin{
		Properties: &armcdn.AFDOriginProperties{
			EnabledState:     to.Ptr(armcdn.EnabledStateEnabled),
			HostName:         to.Ptr("host1.blob.core.windows.net"),
			HTTPPort:         to.Ptr[int32](80),
			HTTPSPort:        to.Ptr[int32](443),
			OriginHostHeader: to.Ptr("host1.foo.com"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step RuleSets_Create
	fmt.Println("Call operation: RuleSets_Create")
	ruleSetsClient, err := armcdn.NewRuleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = ruleSetsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/checkEndpointNameAvailability
func (testsuite *AfdxTestSuite) TestCheckEndpointNameAvailability() {
	var err error
	// From step CheckEndpointNameAvailability
	fmt.Println("Call operation: CheckEndpointNameAvailability")
	managementClient, err := armcdn.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementClient.CheckEndpointNameAvailability(testsuite.ctx, testsuite.resourceGroupName, armcdn.CheckEndpointNameAvailabilityInput{
		Name:                              to.Ptr("sampleName"),
		Type:                              to.Ptr(armcdn.ResourceTypeMicrosoftCdnProfilesAfdEndpoints),
		AutoGeneratedDomainNameLabelScope: to.Ptr(armcdn.AutoGeneratedDomainNameLabelScopeTenantReuse),
	}, nil)
	testsuite.Require().NoError(err)

	// From step AFDProfiles_CheckHostNameAvailability
	fmt.Println("Call operation: AFDProfiles_CheckHostNameAvailability")
	aFDProfilesClient, err := armcdn.NewAFDProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aFDProfilesClient.CheckHostNameAvailability(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, armcdn.CheckHostNameAvailabilityInput{
		HostName: to.Ptr("www.someDomain.net"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step AFDProfiles_ListResourceUsage
	fmt.Println("Call operation: AFDProfiles_ListResourceUsage")
	aFDProfilesClientNewListResourceUsagePager := aFDProfilesClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for aFDProfilesClientNewListResourceUsagePager.More() {
		_, err := aFDProfilesClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/profiles/customDomains
func (testsuite *AfdxTestSuite) TestAfdCustomDomains() {
	var err error
	// From step AFDCustomDomains_ListByProfile
	fmt.Println("Call operation: AFDCustomDomains_ListByProfile")
	aFDCustomDomainsClient, err := armcdn.NewAFDCustomDomainsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDCustomDomainsClientNewListByProfilePager := aFDCustomDomainsClient.NewListByProfilePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for aFDCustomDomainsClientNewListByProfilePager.More() {
		_, err := aFDCustomDomainsClientNewListByProfilePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AFDCustomDomains_Get
	fmt.Println("Call operation: AFDCustomDomains_Get")
	_, err = aFDCustomDomainsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.customDomainName, nil)
	testsuite.Require().NoError(err)

	// From step AFDCustomDomains_Update
	fmt.Println("Call operation: AFDCustomDomains_Update")
	aFDCustomDomainsClientUpdateResponsePoller, err := aFDCustomDomainsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.customDomainName, armcdn.AFDDomainUpdateParameters{
		Properties: &armcdn.AFDDomainUpdatePropertiesParameters{
			TLSSettings: &armcdn.AFDDomainHTTPSParameters{
				CertificateType:   to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
				MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDCustomDomainsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/afdEndpoints
func (testsuite *AfdxTestSuite) TestAfdEndpoints() {
	var err error
	// From step AFDEndpoints_ListByProfile
	fmt.Println("Call operation: AFDEndpoints_ListByProfile")
	aFDEndpointsClient, err := armcdn.NewAFDEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDEndpointsClientNewListByProfilePager := aFDEndpointsClient.NewListByProfilePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for aFDEndpointsClientNewListByProfilePager.More() {
		_, err := aFDEndpointsClientNewListByProfilePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AFDEndpoints_Get
	fmt.Println("Call operation: AFDEndpoints_Get")
	_, err = aFDEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)

	// From step AFDEndpoints_Update
	fmt.Println("Call operation: AFDEndpoints_Update")
	aFDEndpointsClientUpdateResponsePoller, err := aFDEndpointsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.AFDEndpointUpdateParameters{
		Properties: &armcdn.AFDEndpointPropertiesUpdateParameters{
			EnabledState: to.Ptr(armcdn.EnabledStateEnabled),
		},
		Tags: map[string]*string{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDEndpointsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDEndpoints_PurgeContent
	fmt.Println("Call operation: AFDEndpoints_PurgeContent")
	aFDEndpointsClientPurgeContentResponsePoller, err := aFDEndpointsClient.BeginPurgeContent(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.AfdPurgeParameters{
		ContentPaths: []*string{
			to.Ptr("/folder1")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDEndpointsClientPurgeContentResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDEndpoints_ListResourceUsage
	fmt.Println("Call operation: AFDEndpoints_ListResourceUsage")
	aFDEndpointsClientNewListResourceUsagePager := aFDEndpointsClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	for aFDEndpointsClientNewListResourceUsagePager.More() {
		_, err := aFDEndpointsClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AFDEndpoints_ValidateCustomDomain
	fmt.Println("Call operation: AFDEndpoints_ValidateCustomDomain")
	_, err = aFDEndpointsClient.ValidateCustomDomain(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, armcdn.ValidateCustomDomainInput{
		HostName: to.Ptr("www.someDomain.com"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/originGroups
func (testsuite *AfdxTestSuite) TestAfdOriginGroups() {
	var err error
	// From step AFDOriginGroups_ListByProfile
	fmt.Println("Call operation: AFDOriginGroups_ListByProfile")
	aFDOriginGroupsClient, err := armcdn.NewAFDOriginGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginGroupsClientNewListByProfilePager := aFDOriginGroupsClient.NewListByProfilePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for aFDOriginGroupsClientNewListByProfilePager.More() {
		_, err := aFDOriginGroupsClientNewListByProfilePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AFDOriginGroups_Get
	fmt.Println("Call operation: AFDOriginGroups_Get")
	_, err = aFDOriginGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, nil)
	testsuite.Require().NoError(err)

	// From step AFDOriginGroups_Update
	fmt.Println("Call operation: AFDOriginGroups_Update")
	aFDOriginGroupsClientUpdateResponsePoller, err := aFDOriginGroupsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, armcdn.AFDOriginGroupUpdateParameters{
		Properties: &armcdn.AFDOriginGroupUpdatePropertiesParameters{
			HealthProbeSettings: &armcdn.HealthProbeParameters{
				ProbeIntervalInSeconds: to.Ptr[int32](10),
				ProbePath:              to.Ptr("/path2"),
				ProbeProtocol:          to.Ptr(armcdn.ProbeProtocolNotSet),
				ProbeRequestType:       to.Ptr(armcdn.HealthProbeRequestTypeNotSet),
			},
			LoadBalancingSettings: &armcdn.LoadBalancingSettingsParameters{
				AdditionalLatencyInMilliseconds: to.Ptr[int32](1000),
				SampleSize:                      to.Ptr[int32](3),
				SuccessfulSamplesRequired:       to.Ptr[int32](3),
			},
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: to.Ptr[int32](5),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginGroupsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDOriginGroups_ListResourceUsage
	fmt.Println("Call operation: AFDOriginGroups_ListResourceUsage")
	aFDOriginGroupsClientNewListResourceUsagePager := aFDOriginGroupsClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, nil)
	for aFDOriginGroupsClientNewListResourceUsagePager.More() {
		_, err := aFDOriginGroupsClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/profiles/originGroups/origins
func (testsuite *AfdxTestSuite) TestAfdOrigins() {
	var err error
	// From step AFDOrigins_ListByOriginGroup
	fmt.Println("Call operation: AFDOrigins_ListByOriginGroup")
	aFDOriginsClient, err := armcdn.NewAFDOriginsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginsClientNewListByOriginGroupPager := aFDOriginsClient.NewListByOriginGroupPager(testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, nil)
	for aFDOriginsClientNewListByOriginGroupPager.More() {
		_, err := aFDOriginsClientNewListByOriginGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AFDOrigins_Get
	fmt.Println("Call operation: AFDOrigins_Get")
	_, err = aFDOriginsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, testsuite.originName, nil)
	testsuite.Require().NoError(err)

	// From step AFDOrigins_Update
	fmt.Println("Call operation: AFDOrigins_Update")
	aFDOriginsClientUpdateResponsePoller, err := aFDOriginsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, testsuite.originName, armcdn.AFDOriginUpdateParameters{
		Properties: &armcdn.AFDOriginUpdatePropertiesParameters{
			EnabledState: to.Ptr(armcdn.EnabledStateEnabled),
			HostName:     to.Ptr("host1.blob.core.windows.net"),
			HTTPPort:     to.Ptr[int32](80),
			HTTPSPort:    to.Ptr[int32](443),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/afdEndpoints/routes
func (testsuite *AfdxTestSuite) TestRoutes() {
	var err error
	// From step Routes_Create
	fmt.Println("Call operation: Routes_Create")
	routesClient, err := armcdn.NewRoutesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routesClientCreateResponsePoller, err := routesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.routeName, armcdn.Route{
		Properties: &armcdn.RouteProperties{
			CacheConfiguration: &armcdn.AfdRouteCacheConfiguration{
				CompressionSettings: &armcdn.CompressionSettings{
					ContentTypesToCompress: []*string{
						to.Ptr("text/html"),
						to.Ptr("application/octet-stream")},
					IsCompressionEnabled: to.Ptr(true),
				},
				QueryParameters:            to.Ptr("querystring=test"),
				QueryStringCachingBehavior: to.Ptr(armcdn.AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
			},
			CustomDomains: []*armcdn.ActivatedResourceReference{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/customDomains/" + testsuite.customDomainName),
				}},
			EnabledState:        to.Ptr(armcdn.EnabledStateEnabled),
			ForwardingProtocol:  to.Ptr(armcdn.ForwardingProtocolMatchRequest),
			HTTPSRedirect:       to.Ptr(armcdn.HTTPSRedirectEnabled),
			LinkToDefaultDomain: to.Ptr(armcdn.LinkToDefaultDomainEnabled),
			OriginGroup: &armcdn.ResourceReference{
				ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/originGroups/" + testsuite.originGroupName),
			},
			PatternsToMatch: []*string{
				to.Ptr("/*")},
			RuleSets: []*armcdn.ResourceReference{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/ruleSets/" + testsuite.ruleSetName),
				}},
			SupportedProtocols: []*armcdn.AFDEndpointProtocols{
				to.Ptr(armcdn.AFDEndpointProtocolsHTTPS),
				to.Ptr(armcdn.AFDEndpointProtocolsHTTP)},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Routes_ListByEndpoint
	fmt.Println("Call operation: Routes_ListByEndpoint")
	routesClientNewListByEndpointPager := routesClient.NewListByEndpointPager(testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	for routesClientNewListByEndpointPager.More() {
		_, err := routesClientNewListByEndpointPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Routes_Get
	fmt.Println("Call operation: Routes_Get")
	_, err = routesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.routeName, nil)
	testsuite.Require().NoError(err)

	// From step Routes_Update
	fmt.Println("Call operation: Routes_Update")
	routesClientUpdateResponsePoller, err := routesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.routeName, armcdn.RouteUpdateParameters{
		Properties: &armcdn.RouteUpdatePropertiesParameters{
			CacheConfiguration: &armcdn.AfdRouteCacheConfiguration{
				CompressionSettings: &armcdn.CompressionSettings{
					ContentTypesToCompress: []*string{
						to.Ptr("text/html"),
						to.Ptr("application/octet-stream")},
					IsCompressionEnabled: to.Ptr(true),
				},
				QueryStringCachingBehavior: to.Ptr(armcdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
			},
			CustomDomains: []*armcdn.ActivatedResourceReference{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/customDomains/" + testsuite.customDomainName),
				}},
			EnabledState:        to.Ptr(armcdn.EnabledStateEnabled),
			ForwardingProtocol:  to.Ptr(armcdn.ForwardingProtocolMatchRequest),
			HTTPSRedirect:       to.Ptr(armcdn.HTTPSRedirectEnabled),
			LinkToDefaultDomain: to.Ptr(armcdn.LinkToDefaultDomainEnabled),
			OriginGroup: &armcdn.ResourceReference{
				ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/originGroups/" + testsuite.originGroupName),
			},
			PatternsToMatch: []*string{
				to.Ptr("/*")},
			RuleSets: []*armcdn.ResourceReference{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Cdn/profiles/" + testsuite.profileName + "/ruleSets/" + testsuite.ruleSetName),
				}},
			SupportedProtocols: []*armcdn.AFDEndpointProtocols{
				to.Ptr(armcdn.AFDEndpointProtocolsHTTPS),
				to.Ptr(armcdn.AFDEndpointProtocolsHTTP)},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Routes_Delete
	fmt.Println("Call operation: Routes_Delete")
	routesClientDeleteResponsePoller, err := routesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, testsuite.routeName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/ruleSets
func (testsuite *AfdxTestSuite) TestRuleSets() {
	var err error
	// From step RuleSets_ListByProfile
	fmt.Println("Call operation: RuleSets_ListByProfile")
	ruleSetsClient, err := armcdn.NewRuleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ruleSetsClientNewListByProfilePager := ruleSetsClient.NewListByProfilePager(testsuite.resourceGroupName, testsuite.profileName, nil)
	for ruleSetsClientNewListByProfilePager.More() {
		_, err := ruleSetsClientNewListByProfilePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RuleSets_Get
	fmt.Println("Call operation: RuleSets_Get")
	_, err = ruleSetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, nil)
	testsuite.Require().NoError(err)

	// From step RuleSets_ListResourceUsage
	fmt.Println("Call operation: RuleSets_ListResourceUsage")
	ruleSetsClientNewListResourceUsagePager := ruleSetsClient.NewListResourceUsagePager(testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, nil)
	for ruleSetsClientNewListResourceUsagePager.More() {
		_, err := ruleSetsClientNewListResourceUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Cdn/profiles/ruleSets/rules
func (testsuite *AfdxTestSuite) TestRules() {
	var err error
	// From step Rules_Create
	fmt.Println("Call operation: Rules_Create")
	rulesClient, err := armcdn.NewRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	rulesClientCreateResponsePoller, err := rulesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, testsuite.ruleName, armcdn.Rule{
		Properties: &armcdn.RuleProperties{
			Actions: []armcdn.DeliveryRuleActionAutoGeneratedClassification{
				&armcdn.DeliveryRuleResponseHeaderAction{
					Name: to.Ptr(armcdn.DeliveryRuleActionModifyResponseHeader),
					Parameters: &armcdn.HeaderActionParameters{
						HeaderAction: to.Ptr(armcdn.HeaderActionOverwrite),
						HeaderName:   to.Ptr("X-CDN"),
						TypeName:     to.Ptr(armcdn.HeaderActionParametersTypeNameDeliveryRuleHeaderActionParameters),
						Value:        to.Ptr("MSFT"),
					},
				}},
			Conditions: []armcdn.DeliveryRuleConditionClassification{
				&armcdn.DeliveryRuleRequestMethodCondition{
					Name: to.Ptr(armcdn.MatchVariableRequestMethod),
					Parameters: &armcdn.RequestMethodMatchConditionParameters{
						MatchValues: []*armcdn.RequestMethodMatchConditionParametersMatchValuesItem{
							to.Ptr(armcdn.RequestMethodMatchConditionParametersMatchValuesItemGET)},
						NegateCondition: to.Ptr(false),
						Operator:        to.Ptr(armcdn.RequestMethodOperatorEqual),
						TypeName:        to.Ptr(armcdn.RequestMethodMatchConditionParametersTypeNameDeliveryRuleRequestMethodConditionParameters),
					},
				}},
			Order: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, rulesClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Rules_ListByRuleSet
	fmt.Println("Call operation: Rules_ListByRuleSet")
	rulesClientNewListByRuleSetPager := rulesClient.NewListByRuleSetPager(testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, nil)
	for rulesClientNewListByRuleSetPager.More() {
		_, err := rulesClientNewListByRuleSetPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Rules_Get
	fmt.Println("Call operation: Rules_Get")
	_, err = rulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)

	// From step Rules_Update
	fmt.Println("Call operation: Rules_Update")
	rulesClientUpdateResponsePoller, err := rulesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, testsuite.ruleName, armcdn.RuleUpdateParameters{
		Properties: &armcdn.RuleUpdatePropertiesParameters{
			Actions: []armcdn.DeliveryRuleActionAutoGeneratedClassification{
				&armcdn.DeliveryRuleResponseHeaderAction{
					Name: to.Ptr(armcdn.DeliveryRuleActionModifyResponseHeader),
					Parameters: &armcdn.HeaderActionParameters{
						HeaderAction: to.Ptr(armcdn.HeaderActionOverwrite),
						HeaderName:   to.Ptr("X-CDN"),
						TypeName:     to.Ptr(armcdn.HeaderActionParametersTypeNameDeliveryRuleHeaderActionParameters),
						Value:        to.Ptr("MSFT"),
					},
				}},
			Order: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, rulesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Rules_Delete
	fmt.Println("Call operation: Rules_Delete")
	rulesClientDeleteResponsePoller, err := rulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, rulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Cdn/profiles/LogAnalytics
func (testsuite *AfdxTestSuite) TestLogAnalytics() {
	var err error
	// From step LogAnalytics_GetLogAnalyticsLocations
	fmt.Println("Call operation: LogAnalytics_GetLogAnalyticsLocations")
	logAnalyticsClient, err := armcdn.NewLogAnalyticsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = logAnalyticsClient.GetLogAnalyticsLocations(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, nil)
	testsuite.Require().NoError(err)

	// From step LogAnalytics_GetLogAnalyticsRankings
	fmt.Println("Call operation: LogAnalytics_GetLogAnalyticsRankings")
	_, err = logAnalyticsClient.GetLogAnalyticsRankings(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, []armcdn.LogRanking{
		armcdn.LogRankingURL}, []armcdn.LogRankingMetric{
		armcdn.LogRankingMetricClientRequestCount}, 5, func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T06:49:27.554Z"); return t }(), func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T09:49:27.554Z"); return t }(), &armcdn.LogAnalyticsClientGetLogAnalyticsRankingsOptions{CustomDomains: []string{}})
	testsuite.Require().NoError(err)

	// From step LogAnalytics_GetLogAnalyticsResources
	fmt.Println("Call operation: LogAnalytics_GetLogAnalyticsResources")
	_, err = logAnalyticsClient.GetLogAnalyticsResources(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, nil)
	testsuite.Require().NoError(err)

	// From step LogAnalytics_GetWafLogAnalyticsMetrics
	fmt.Println("Call operation: LogAnalytics_GetWafLogAnalyticsMetrics")
	_, err = logAnalyticsClient.GetWafLogAnalyticsMetrics(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, []armcdn.WafMetric{
		armcdn.WafMetricClientRequestCount}, func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T06:49:27.554Z"); return t }(), func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T09:49:27.554Z"); return t }(), armcdn.WafGranularityPT5M, &armcdn.LogAnalyticsClientGetWafLogAnalyticsMetricsOptions{Actions: []armcdn.WafAction{
		armcdn.WafActionBlock,
		armcdn.WafActionLog},
		GroupBy:   []armcdn.WafRankingGroupBy{},
		RuleTypes: []armcdn.WafRuleType{},
	})
	testsuite.Require().NoError(err)

	// From step LogAnalytics_GetWafLogAnalyticsRankings
	fmt.Println("Call operation: LogAnalytics_GetWafLogAnalyticsRankings")
	_, err = logAnalyticsClient.GetWafLogAnalyticsRankings(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, []armcdn.WafMetric{
		armcdn.WafMetricClientRequestCount}, func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T06:49:27.554Z"); return t }(), func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-04T09:49:27.554Z"); return t }(), 5, []armcdn.WafRankingType{
		armcdn.WafRankingTypeRuleID}, &armcdn.LogAnalyticsClientGetWafLogAnalyticsRankingsOptions{Actions: []armcdn.WafAction{},
		RuleTypes: []armcdn.WafRuleType{},
	})
	testsuite.Require().NoError(err)
}

func (testsuite *AfdxTestSuite) Cleanup() {
	var err error
	// From step RuleSets_Delete
	fmt.Println("Call operation: RuleSets_Delete")
	ruleSetsClient, err := armcdn.NewRuleSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	ruleSetsClientDeleteResponsePoller, err := ruleSetsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.ruleSetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, ruleSetsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDOrigins_Delete
	fmt.Println("Call operation: AFDOrigins_Delete")
	aFDOriginsClient, err := armcdn.NewAFDOriginsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginsClientDeleteResponsePoller, err := aFDOriginsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, testsuite.originName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDOriginGroups_Delete
	fmt.Println("Call operation: AFDOriginGroups_Delete")
	aFDOriginGroupsClient, err := armcdn.NewAFDOriginGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDOriginGroupsClientDeleteResponsePoller, err := aFDOriginGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.originGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDOriginGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDEndpoints_Delete
	fmt.Println("Call operation: AFDEndpoints_Delete")
	aFDEndpointsClient, err := armcdn.NewAFDEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDEndpointsClientDeleteResponsePoller, err := aFDEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.endpointName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step AFDCustomDomains_Delete
	fmt.Println("Call operation: AFDCustomDomains_Delete")
	aFDCustomDomainsClient, err := armcdn.NewAFDCustomDomainsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aFDCustomDomainsClientDeleteResponsePoller, err := aFDCustomDomainsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.profileName, testsuite.customDomainName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aFDCustomDomainsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
