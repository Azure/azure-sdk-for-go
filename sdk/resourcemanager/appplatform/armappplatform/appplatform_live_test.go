//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armappplatform_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/stretchr/testify/suite"
)

type AppplatformTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	appName           string
	routeConfigName   string
	serviceName       string
	storageName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *AppplatformTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/appplatform/armappplatform/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.appName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "appname", 14, true)
	testsuite.routeConfigName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "routecon", 14, true)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicen", 14, true)
	testsuite.storageName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storagen", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AppplatformTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAppplatformTestSuite(t *testing.T) {
	suite.Run(t, new(AppplatformTestSuite))
}

func (testsuite *AppplatformTestSuite) Prepare() {
	var err error
	// From step Services_CheckNameAvailability
	fmt.Println("Call operation: Services_CheckNameAvailability")
	servicesClient, err := armappplatform.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = servicesClient.CheckNameAvailability(testsuite.ctx, testsuite.location, armappplatform.NameAvailabilityParameters{
		Name: to.Ptr(testsuite.serviceName),
		Type: to.Ptr("Microsoft.AppPlatform/Spring"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Services_CreateOrUpdate
	fmt.Println("Call operation: Services_CreateOrUpdate")
	servicesClientCreateOrUpdateResponsePoller, err := servicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ServiceResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armappplatform.ClusterResourceProperties{
			MarketplaceResource: &armappplatform.MarketplaceResource{
				Plan:      to.Ptr("tanzu-asc-ent-mtr"),
				Product:   to.Ptr("azure-spring-cloud-vmware-tanzu-2"),
				Publisher: to.Ptr("vmware-inc"),
			},
		},
		SKU: &armappplatform.SKU{
			Name: to.Ptr("E0"),
			Tier: to.Ptr("Enterprise"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Apps_CreateOrUpdate
	fmt.Println("Call operation: Apps_CreateOrUpdate")
	appsClient, err := armappplatform.NewAppsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	appsClientCreateOrUpdateResponsePoller, err := appsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, armappplatform.AppResource{
		Location: to.Ptr(testsuite.location),
		Properties: &armappplatform.AppResourceProperties{
			EnableEndToEndTLS: to.Ptr(false),
			HTTPSOnly:         to.Ptr(false),
			PersistentDisk: &armappplatform.PersistentDisk{
				MountPath: to.Ptr("/mypersistentdisk"),
				SizeInGB:  to.Ptr[int32](0),
			},
			Public: to.Ptr(false),
			TemporaryDisk: &armappplatform.TemporaryDisk{
				MountPath: to.Ptr("/mytemporarydisk"),
				SizeInGB:  to.Ptr[int32](2),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApplicationAccelerators_CreateOrUpdate
	fmt.Println("Call operation: ApplicationAccelerators_CreateOrUpdate")
	applicationAcceleratorsClient, err := armappplatform.NewApplicationAcceleratorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationAcceleratorsClientCreateOrUpdateResponsePoller, err := applicationAcceleratorsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.ApplicationAcceleratorResource{
		Properties: &armappplatform.ApplicationAcceleratorProperties{},
		SKU: &armappplatform.SKU{
			Name:     to.Ptr("E0"),
			Capacity: to.Ptr[int32](2),
			Tier:     to.Ptr("Enterprise"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationAcceleratorsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Gateways_CreateOrUpdate
	fmt.Println("Call operation: Gateways_CreateOrUpdate")
	gatewaysClient, err := armappplatform.NewGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gatewaysClientCreateOrUpdateResponsePoller, err := gatewaysClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.GatewayResource{
		Properties: &armappplatform.GatewayProperties{
			Public: to.Ptr(true),
		},
		SKU: &armappplatform.SKU{
			Name:     to.Ptr("E0"),
			Capacity: to.Ptr[int32](2),
			Tier:     to.Ptr("Enterprise"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gatewaysClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApiPortals_CreateOrUpdate
	fmt.Println("Call operation: ApiPortals_CreateOrUpdate")
	aPIPortalsClient, err := armappplatform.NewAPIPortalsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aPIPortalsClientCreateOrUpdateResponsePoller, err := aPIPortalsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.APIPortalResource{
		Properties: &armappplatform.APIPortalProperties{
			GatewayIDs: []*string{
				to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.AppPlatform/Spring/" + testsuite.serviceName + "/gateways/default")},
			Public: to.Ptr(true),
		},
		SKU: &armappplatform.SKU{
			Name:     to.Ptr("E0"),
			Capacity: to.Ptr[int32](2),
			Tier:     to.Ptr("Enterprise"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aPIPortalsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}
func (testsuite *AppplatformTestSuite) TestServices() {
	var err error
	// From step Services_ListBySubscription
	fmt.Println("Call operation: Services_ListBySubscription")
	servicesClient, err := armappplatform.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	servicesClientNewListBySubscriptionPager := servicesClient.NewListBySubscriptionPager(nil)
	for servicesClientNewListBySubscriptionPager.More() {
		_, err := servicesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Services_List
	fmt.Println("Call operation: Services_List")
	servicesClientNewListPager := servicesClient.NewListPager(testsuite.resourceGroupName, nil)
	for servicesClientNewListPager.More() {
		_, err := servicesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Services_Get
	fmt.Println("Call operation: Services_Get")
	_, err = servicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step Services_Update
	fmt.Println("Call operation: Services_Update")
	servicesClientUpdateResponsePoller, err := servicesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.ServiceResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armappplatform.ClusterResourceProperties{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Services_RegenerateTestKey
	fmt.Println("Call operation: Services_RegenerateTestKey")
	_, err = servicesClient.RegenerateTestKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.RegenerateTestKeyRequestPayload{
		KeyType: to.Ptr(armappplatform.TestKeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Services_ListTestKeys
	fmt.Println("Call operation: Services_ListTestKeys")
	_, err = servicesClient.ListTestKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step Services_DisableTestEndpoint
	fmt.Println("Call operation: Services_DisableTestEndpoint")
	_, err = servicesClient.DisableTestEndpoint(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step Services_EnableTestEndpoint
	fmt.Println("Call operation: Services_EnableTestEndpoint")
	_, err = servicesClient.EnableTestEndpoint(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}
func (testsuite *AppplatformTestSuite) TestApps() {
	var err error
	// From step Apps_List
	fmt.Println("Call operation: Apps_List")
	appsClient, err := armappplatform.NewAppsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	appsClientNewListPager := appsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for appsClientNewListPager.More() {
		_, err := appsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Apps_Get
	fmt.Println("Call operation: Apps_Get")
	_, err = appsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, &armappplatform.AppsClientGetOptions{SyncStatus: nil})
	testsuite.Require().NoError(err)

	// From step Apps_Update
	fmt.Println("Call operation: Apps_Update")
	appsClientUpdateResponsePoller, err := appsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, armappplatform.AppResource{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Apps_GetResourceUploadUrl
	fmt.Println("Call operation: Apps_GetResourceUploadUrl")
	_, err = appsClient.GetResourceUploadURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, nil)
	testsuite.Require().NoError(err)

	// From step Apps_ValidateDomain
	fmt.Println("Call operation: Apps_ValidateDomain")
	_, err = appsClient.ValidateDomain(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, armappplatform.CustomDomainValidatePayload{
		Name: to.Ptr("mydomain.io"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}
func (testsuite *AppplatformTestSuite) TestGateways() {
	var err error
	// From step Gateways_List
	fmt.Println("Call operation: Gateways_List")
	gatewaysClient, err := armappplatform.NewGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gatewaysClientNewListPager := gatewaysClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for gatewaysClientNewListPager.More() {
		_, err := gatewaysClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Gateways_Get
	fmt.Println("Call operation: Gateways_Get")
	_, err = gatewaysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step Gateways_UpdateCapacity
	// fmt.Println("Call operation: Gateways_UpdateCapacity")
	// gatewaysClientUpdateCapacityResponsePoller, err := gatewaysClient.BeginUpdateCapacity(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.SKUObject{
	// 	SKU: &armappplatform.SKU{
	// 		Name:     to.Ptr("E0"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Tier:     to.Ptr("Enterprise"),
	// 	},
	// }, nil)
	// testsuite.Require().NoError(err)
	// _, err = testutil.PollForTest(testsuite.ctx, gatewaysClientUpdateCapacityResponsePoller)
	// testsuite.Require().NoError(err)

	// From step Gateways_ValidateDomain
	fmt.Println("Call operation: Gateways_ValidateDomain")
	_, err = gatewaysClient.ValidateDomain(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.CustomDomainValidatePayload{
		Name: to.Ptr("mydomain.io"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Gateways_ListEnvSecrets
	fmt.Println("Call operation: Gateways_ListEnvSecrets")
	_, err = gatewaysClient.ListEnvSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}/routeConfigs/{routeConfigName}
func (testsuite *AppplatformTestSuite) TestGatewayRouteConfigs() {
	var err error
	// From step GatewayRouteConfigs_CreateOrUpdate
	fmt.Println("Call operation: GatewayRouteConfigs_CreateOrUpdate")
	gatewayRouteConfigsClient, err := armappplatform.NewGatewayRouteConfigsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	gatewayRouteConfigsClientCreateOrUpdateResponsePoller, err := gatewayRouteConfigsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", testsuite.routeConfigName, armappplatform.GatewayRouteConfigResource{
		Properties: &armappplatform.GatewayRouteConfigProperties{
			AppResourceID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.AppPlatform/Spring/" + testsuite.serviceName + "/apps/myApp"),
			OpenAPI: &armappplatform.GatewayRouteConfigOpenAPIProperties{
				URI: to.Ptr("https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore.json"),
			},
			Routes: []*armappplatform.GatewayAPIRoute{
				{
					Filters: []*string{
						to.Ptr("StripPrefix=2"),
						to.Ptr("RateLimit=1,1s")},
					Predicates: []*string{
						to.Ptr("Path=/api5/customer/**")},
					SsoEnabled: to.Ptr(true),
					Title:      to.Ptr("myApp route config"),
				}},
			Protocol: to.Ptr(armappplatform.GatewayRouteConfigProtocolHTTPS),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gatewayRouteConfigsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GatewayRouteConfigs_List
	fmt.Println("Call operation: GatewayRouteConfigs_List")
	gatewayRouteConfigsClientNewListPager := gatewayRouteConfigsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	for gatewayRouteConfigsClientNewListPager.More() {
		_, err := gatewayRouteConfigsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GatewayRouteConfigs_Get
	fmt.Println("Call operation: GatewayRouteConfigs_Get")
	_, err = gatewayRouteConfigsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", testsuite.routeConfigName, nil)
	testsuite.Require().NoError(err)

	// From step GatewayRouteConfigs_Delete
	fmt.Println("Call operation: GatewayRouteConfigs_Delete")
	gatewayRouteConfigsClientDeleteResponsePoller, err := gatewayRouteConfigsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", testsuite.routeConfigName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, gatewayRouteConfigsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/applicationAccelerators/{applicationAcceleratorName}
func (testsuite *AppplatformTestSuite) TestApplicationAccelerators() {
	var err error
	// From step ApplicationAccelerators_List
	fmt.Println("Call operation: ApplicationAccelerators_List")
	applicationAcceleratorsClient, err := armappplatform.NewApplicationAcceleratorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationAcceleratorsClientNewListPager := applicationAcceleratorsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for applicationAcceleratorsClientNewListPager.More() {
		_, err := applicationAcceleratorsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationAccelerators_Get
	fmt.Println("Call operation: ApplicationAccelerators_Get")
	_, err = applicationAcceleratorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/applicationAccelerators/{applicationAcceleratorName}/predefinedAccelerators/{predefinedAcceleratorName}
func (testsuite *AppplatformTestSuite) TestPredefinedAccelerators() {
	var err error
	var predefinedAcceleratorName string
	// From step PredefinedAccelerators_List
	fmt.Println("Call operation: PredefinedAccelerators_List")
	predefinedAcceleratorsClient, err := armappplatform.NewPredefinedAcceleratorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	predefinedAcceleratorsClientNewListPager := predefinedAcceleratorsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	for predefinedAcceleratorsClientNewListPager.More() {
		predefinedAcceleratorsClientListResponse, err := predefinedAcceleratorsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		predefinedAcceleratorName = *predefinedAcceleratorsClientListResponse.Value[0].Name
		break
	}

	// From step PredefinedAccelerators_Get
	fmt.Println("Call operation: PredefinedAccelerators_Get")
	_, err = predefinedAcceleratorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", predefinedAcceleratorName, nil)
	testsuite.Require().NoError(err)

	// From step PredefinedAccelerators_Disable
	fmt.Println("Call operation: PredefinedAccelerators_Disable")
	predefinedAcceleratorsClientDisableResponsePoller, err := predefinedAcceleratorsClient.BeginDisable(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", predefinedAcceleratorName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, predefinedAcceleratorsClientDisableResponsePoller)
	testsuite.Require().NoError(err)

	// From step PredefinedAccelerators_Enable
	fmt.Println("Call operation: PredefinedAccelerators_Enable")
	predefinedAcceleratorsClientEnableResponsePoller, err := predefinedAcceleratorsClient.BeginEnable(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", predefinedAcceleratorName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, predefinedAcceleratorsClientEnableResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/apiPortals/{apiPortalName}
func (testsuite *AppplatformTestSuite) TestApiPortals() {
	var err error
	// From step ApiPortals_List
	fmt.Println("Call operation: ApiPortals_List")
	aPIPortalsClient, err := armappplatform.NewAPIPortalsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aPIPortalsClientNewListPager := aPIPortalsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for aPIPortalsClientNewListPager.More() {
		_, err := aPIPortalsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiPortals_Get
	fmt.Println("Call operation: ApiPortals_Get")
	_, err = aPIPortalsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step ApiPortals_ValidateDomain
	fmt.Println("Call operation: ApiPortals_ValidateDomain")
	_, err = aPIPortalsClient.ValidateDomain(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.CustomDomainValidatePayload{
		Name: to.Ptr("mydomain.io"),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/configurationServices/{configurationServiceName}
func (testsuite *AppplatformTestSuite) TestConfigurationServices() {
	var err error
	// From step ConfigurationServices_CreateOrUpdate
	fmt.Println("Call operation: ConfigurationServices_CreateOrUpdate")
	configurationServicesClient, err := armappplatform.NewConfigurationServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	configurationServicesClientCreateOrUpdateResponsePoller, err := configurationServicesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.ConfigurationServiceResource{
		Properties: &armappplatform.ConfigurationServiceProperties{
			Settings: &armappplatform.ConfigurationServiceSettings{
				GitProperty: &armappplatform.ConfigurationServiceGitProperty{
					Repositories: []*armappplatform.ConfigurationServiceGitRepository{
						{
							Name:  to.Ptr("fake"),
							Label: to.Ptr("main"),
							Patterns: []*string{
								to.Ptr("app/dev")},
							URI: to.Ptr("https://github.com/Azure/azure-sdk-for-go"),
						}},
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationServicesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigurationServices_List
	fmt.Println("Call operation: ConfigurationServices_List")
	configurationServicesClientNewListPager := configurationServicesClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for configurationServicesClientNewListPager.More() {
		_, err := configurationServicesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConfigurationServices_Get
	fmt.Println("Call operation: ConfigurationServices_Get")
	_, err = configurationServicesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step ConfigurationServices_Validate
	fmt.Println("Call operation: ConfigurationServices_Validate")
	configurationServicesClientValidateResponsePoller, err := configurationServicesClient.BeginValidate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.ConfigurationServiceSettings{
		GitProperty: &armappplatform.ConfigurationServiceGitProperty{
			Repositories: []*armappplatform.ConfigurationServiceGitRepository{
				{
					Name:  to.Ptr("fake"),
					Label: to.Ptr("main"),
					Patterns: []*string{
						to.Ptr("app/dev")},
					URI: to.Ptr("https://github.com/Azure/azure-sdk-for-go"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationServicesClientValidateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ConfigurationServices_Delete
	fmt.Println("Call operation: ConfigurationServices_Delete")
	configurationServicesClientDeleteResponsePoller, err := configurationServicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, configurationServicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/applicationLiveViews/{applicationLiveViewName}
func (testsuite *AppplatformTestSuite) TestApplicationLiveViews() {
	var err error
	// From step ApplicationLiveViews_CreateOrUpdate
	fmt.Println("Call operation: ApplicationLiveViews_CreateOrUpdate")
	applicationLiveViewsClient, err := armappplatform.NewApplicationLiveViewsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationLiveViewsClientCreateOrUpdateResponsePoller, err := applicationLiveViewsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.ApplicationLiveViewResource{
		Properties: &armappplatform.ApplicationLiveViewProperties{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationLiveViewsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApplicationLiveViews_List
	fmt.Println("Call operation: ApplicationLiveViews_List")
	applicationLiveViewsClientNewListPager := applicationLiveViewsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for applicationLiveViewsClientNewListPager.More() {
		_, err := applicationLiveViewsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApplicationLiveViews_Get
	fmt.Println("Call operation: ApplicationLiveViews_Get")
	_, err = applicationLiveViewsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step ApplicationLiveViews_Delete
	fmt.Println("Call operation: ApplicationLiveViews_Delete")
	applicationLiveViewsClientDeleteResponsePoller, err := applicationLiveViewsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationLiveViewsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/storages/{storageName}
func (testsuite *AppplatformTestSuite) TestStorages() {
	var err error
	// From step Storages_CreateOrUpdate
	fmt.Println("Call operation: Storages_CreateOrUpdate")
	storagesClient, err := armappplatform.NewStoragesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	storagesClientCreateOrUpdateResponsePoller, err := storagesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.storageName, armappplatform.StorageResource{
		Properties: &armappplatform.StorageAccount{
			StorageType: to.Ptr(armappplatform.StorageTypeStorageAccount),
			AccountKey:  to.Ptr("account-key-of-storage-account"),
			AccountName: to.Ptr("storage-account-name"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storagesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Storages_List
	fmt.Println("Call operation: Storages_List")
	storagesClientNewListPager := storagesClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for storagesClientNewListPager.More() {
		_, err := storagesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Storages_Get
	fmt.Println("Call operation: Storages_Get")
	_, err = storagesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.storageName, nil)
	testsuite.Require().NoError(err)

	// From step Storages_Delete
	fmt.Println("Call operation: Storages_Delete")
	storagesClientDeleteResponsePoller, err := storagesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.storageName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, storagesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/serviceRegistries/{serviceRegistryName}
func (testsuite *AppplatformTestSuite) TestServiceRegistries() {
	var err error
	// From step ServiceRegistries_CreateOrUpdate
	fmt.Println("Call operation: ServiceRegistries_CreateOrUpdate")
	serviceRegistriesClient, err := armappplatform.NewServiceRegistriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceRegistriesClientCreateOrUpdateResponsePoller, err := serviceRegistriesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceRegistriesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ServiceRegistries_List
	fmt.Println("Call operation: ServiceRegistries_List")
	serviceRegistriesClientNewListPager := serviceRegistriesClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for serviceRegistriesClientNewListPager.More() {
		_, err := serviceRegistriesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ServiceRegistries_Get
	fmt.Println("Call operation: ServiceRegistries_Get")
	_, err = serviceRegistriesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step ServiceRegistries_Delete
	fmt.Println("Call operation: ServiceRegistries_Delete")
	serviceRegistriesClientDeleteResponsePoller, err := serviceRegistriesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceRegistriesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/DevToolPortals/{devToolPortalName}
func (testsuite *AppplatformTestSuite) TestDevToolPortals() {
	var err error
	// From step DevToolPortals_CreateOrUpdate
	fmt.Println("Call operation: DevToolPortals_CreateOrUpdate")
	devToolPortalsClient, err := armappplatform.NewDevToolPortalsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	devToolPortalsClientCreateOrUpdateResponsePoller, err := devToolPortalsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", armappplatform.DevToolPortalResource{
		Properties: &armappplatform.DevToolPortalProperties{
			Features: &armappplatform.DevToolPortalFeatureSettings{
				ApplicationAccelerator: &armappplatform.DevToolPortalFeatureDetail{
					State: to.Ptr(armappplatform.DevToolPortalFeatureStateEnabled),
				},
				ApplicationLiveView: &armappplatform.DevToolPortalFeatureDetail{
					State: to.Ptr(armappplatform.DevToolPortalFeatureStateEnabled),
				},
			},
			Public: to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, devToolPortalsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DevToolPortals_List
	fmt.Println("Call operation: DevToolPortals_List")
	devToolPortalsClientNewListPager := devToolPortalsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, nil)
	for devToolPortalsClientNewListPager.More() {
		_, err := devToolPortalsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DevToolPortals_Get
	fmt.Println("Call operation: DevToolPortals_Get")
	_, err = devToolPortalsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)

	// From step DevToolPortals_Delete
	fmt.Println("Call operation: DevToolPortals_Delete")
	devToolPortalsClientDeleteResponsePoller, err := devToolPortalsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, devToolPortalsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/Spring/{serviceName}/monitoringSettings/default
func (testsuite *AppplatformTestSuite) TestMonitoringSettings() {
	var err error
	// From step MonitoringSettings_UpdatePut
	fmt.Println("Call operation: MonitoringSettings_UpdatePut")
	monitoringSettingsClient, err := armappplatform.NewMonitoringSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	monitoringSettingsClientUpdatePutResponsePoller, err := monitoringSettingsClient.BeginUpdatePut(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.MonitoringSettingResource{
		Properties: &armappplatform.MonitoringSettingProperties{
			AppInsightsInstrumentationKey: to.Ptr(testsuite.subscriptionId),
			AppInsightsSamplingRate:       to.Ptr[float64](10),
			TraceEnabled:                  to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, monitoringSettingsClientUpdatePutResponsePoller)
	testsuite.Require().NoError(err)

	// From step MonitoringSettings_Get
	fmt.Println("Call operation: MonitoringSettings_Get")
	_, err = monitoringSettingsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)

	// From step MonitoringSettings_UpdatePatch
	fmt.Println("Call operation: MonitoringSettings_UpdatePatch")
	monitoringSettingsClientUpdatePatchResponsePoller, err := monitoringSettingsClient.BeginUpdatePatch(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armappplatform.MonitoringSettingResource{
		Properties: &armappplatform.MonitoringSettingProperties{
			AppInsightsInstrumentationKey: to.Ptr(testsuite.subscriptionId),
			AppInsightsSamplingRate:       to.Ptr[float64](10),
			TraceEnabled:                  to.Ptr(true),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, monitoringSettingsClientUpdatePatchResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.AppPlatform/operations
func (testsuite *AppplatformTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armappplatform.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RuntimeVersions_ListRuntimeVersions
	fmt.Println("Call operation: RuntimeVersions_ListRuntimeVersions")
	runtimeVersionsClient, err := armappplatform.NewRuntimeVersionsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = runtimeVersionsClient.ListRuntimeVersions(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step Skus_List
	fmt.Println("Call operation: Skus_List")
	sKUsClient, err := armappplatform.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(nil)
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *AppplatformTestSuite) Cleanup() {
	var err error
	// From step ApplicationAccelerators_Delete
	fmt.Println("Call operation: ApplicationAccelerators_Delete")
	applicationAcceleratorsClient, err := armappplatform.NewApplicationAcceleratorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	applicationAcceleratorsClientDeleteResponsePoller, err := applicationAcceleratorsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, applicationAcceleratorsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Apps_Delete
	fmt.Println("Call operation: Apps_Delete")
	appsClient, err := armappplatform.NewAppsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	appsClientDeleteResponsePoller, err := appsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.appName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, appsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step Services_Delete
	fmt.Println("Call operation: Services_Delete")
	servicesClient, err := armappplatform.NewServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	servicesClientDeleteResponsePoller, err := servicesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, servicesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
