// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type AutoscaleTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	autoscaleSettingName string
	subnetId             string
	vmssId               string
	adminPassword        string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *AutoscaleTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.autoscaleSettingName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "autoscalesettingna", 24, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AutoscaleTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAutoscaleTestSuite(t *testing.T) {
	suite.Run(t, new(AutoscaleTestSuite))
}

func (testsuite *AutoscaleTestSuite) Prepare() {
	var err error
	// From step NetworkAndSubnet_Create
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"subnetId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'vmsssubnet')]",
			},
		},
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
			"virtualNetworksName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "vmssvnet",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []interface{}{
							"10.0.0.0/16",
						},
					},
					"subnets": []interface{}{
						map[string]interface{}{
							"name": "vmsssubnet",
							"properties": map[string]interface{}{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
				"tags": map[string]interface{}{},
			},
		},
	}
	params := map[string]interface{}{
		"location": map[string]interface{}{"value": testsuite.location},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "NetworkAndSubnet_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step NetworkInterface_Create
	template = map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
			"subnetId": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(subnetId)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "vmssnic",
				"type":       "Microsoft.Network/networkInterfaces",
				"apiVersion": "2021-08-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"ipConfigurations": []interface{}{
						map[string]interface{}{
							"name": "Ipv4config",
							"properties": map[string]interface{}{
								"subnet": map[string]interface{}{
									"id": "[parameters('subnetId')]",
								},
							},
						},
					},
				},
				"tags": map[string]interface{}{},
			},
		},
	}
	params = map[string]interface{}{
		"location": map[string]interface{}{"value": testsuite.location},
		"subnetId": map[string]interface{}{"value": testsuite.subnetId},
	}
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	_, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "NetworkInterface_Create", &deployment)
	testsuite.Require().NoError(err)

	// From step VMSS_Create
	template = map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"vmssId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Compute/virtualMachineScaleSets', parameters('virtualMachineScaleSetName'))]",
			},
		},
		"parameters": map[string]interface{}{
			"adminPassword": map[string]interface{}{
				"type":         "securestring",
				"defaultValue": "$(adminPassword)",
			},
			"location": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(location)",
			},
			"subnetId": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(subnetId)",
			},
			"virtualMachineScaleSetName": map[string]interface{}{
				"type":         "string",
				"defaultValue": "monitorvmss",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "[parameters('virtualMachineScaleSetName')]",
				"type":       "Microsoft.Compute/virtualMachineScaleSets",
				"apiVersion": "2022-03-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"scaleInPolicy": map[string]interface{}{
						"rules": []interface{}{
							"Default",
						},
					},
					"singlePlacementGroup": false,
					"upgradePolicy": map[string]interface{}{
						"mode": "Manual",
					},
					"virtualMachineProfile": map[string]interface{}{
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled": true,
							},
						},
						"networkProfile": map[string]interface{}{
							"networkInterfaceConfigurations": []interface{}{
								map[string]interface{}{
									"name": "[concat(parameters('virtualMachineScaleSetName'), 'vnet-nic01')]",
									"properties": map[string]interface{}{
										"enableAcceleratedNetworking": true,
										"enableIPForwarding":          false,
										"ipConfigurations": []interface{}{
											map[string]interface{}{
												"name": "[concat(parameters('virtualMachineScaleSetName'), 'vnet-nic01-defaultIpConfiguration')]",
												"properties": map[string]interface{}{
													"primary":                 true,
													"privateIPAddressVersion": "IPv4",
													"subnet": map[string]interface{}{
														"id": "[parameters('subnetId')]",
													},
												},
											},
										},
										"primary": true,
									},
								},
							},
						},
						"osProfile": map[string]interface{}{
							"adminPassword":            "[parameters('adminPassword')]",
							"adminUsername":            "azureuser",
							"allowExtensionOperations": true,
							"computerNamePrefix":       "vmss",
							"secrets":                  []interface{}{},
							"windowsConfiguration": map[string]interface{}{
								"enableAutomaticUpdates": true,
								"provisionVMAgent":       true,
							},
						},
						"storageProfile": map[string]interface{}{
							"imageReference": map[string]interface{}{
								"offer":     "WindowsServer",
								"publisher": "MicrosoftWindowsServer",
								"sku":       "2019-Datacenter",
								"version":   "latest",
							},
							"osDisk": map[string]interface{}{
								"caching":      "ReadWrite",
								"createOption": "FromImage",
								"diskSizeGB":   float64(127),
								"managedDisk": map[string]interface{}{
									"storageAccountType": "Premium_LRS",
								},
								"osType": "Windows",
							},
						},
					},
				},
				"sku": map[string]interface{}{
					"name":     "Standard_DS1_v2",
					"capacity": float64(2),
					"tier":     "Standard",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params = map[string]interface{}{
		"adminPassword": map[string]interface{}{"value": testsuite.adminPassword},
		"location":      map[string]interface{}{"value": testsuite.location},
		"subnetId":      map[string]interface{}{"value": testsuite.subnetId},
	}
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "VMSS_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.vmssId = deploymentExtend.Properties.Outputs.(map[string]interface{})["vmssId"].(map[string]interface{})["value"].(string)
}

// Microsoft.Insights/autoscalesettings
func (testsuite *AutoscaleTestSuite) TestAutoscalesettings() {
	var err error
	// From step AutoscaleSettings_Create
	fmt.Println("Call operation: AutoscaleSettings_CreateOrUpdate")
	autoscaleSettingsClient, err := armmonitor.NewAutoscaleSettingsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = autoscaleSettingsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.autoscaleSettingName, armmonitor.AutoscaleSettingResource{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armmonitor.AutoscaleSetting{
			Enabled: to.Ptr(true),
			Notifications: []*armmonitor.AutoscaleNotification{
				{
					Email: &armmonitor.EmailNotification{
						CustomEmails: []*string{
							to.Ptr("gu@ms.com"),
							to.Ptr("ge@ns.net")},
						SendToSubscriptionAdministrator:    to.Ptr(true),
						SendToSubscriptionCoAdministrators: to.Ptr(true),
					},
					Operation: to.Ptr("Scale"),
				}},
			Profiles: []*armmonitor.AutoscaleProfile{
				{
					Name: to.Ptr("adios"),
					Capacity: &armmonitor.ScaleCapacity{
						Default: to.Ptr("1"),
						Maximum: to.Ptr("10"),
						Minimum: to.Ptr("1"),
					},
					FixedDate: &armmonitor.TimeWindow{
						End:      to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2015-03-05T14:30:00Z"); return t }()),
						Start:    to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2015-03-05T14:00:00Z"); return t }()),
						TimeZone: to.Ptr("UTC"),
					},
				}},
			TargetResourceURI: to.Ptr(testsuite.vmssId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AutoscaleSettings_Update
	fmt.Println("Call operation: AutoscaleSettings_Update")
	_, err = autoscaleSettingsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.autoscaleSettingName, armmonitor.AutoscaleSettingResourcePatch{
		Tags: map[string]*string{
			"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AutoscaleSettings_ListBySubscription
	fmt.Println("Call operation: AutoscaleSettings_ListBySubscription")
	autoscaleSettingsClientNewListBySubscriptionPager := autoscaleSettingsClient.NewListBySubscriptionPager(nil)
	for autoscaleSettingsClientNewListBySubscriptionPager.More() {
		_, err := autoscaleSettingsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AutoscaleSettings_ListByResourceGroup
	fmt.Println("Call operation: AutoscaleSettings_ListByResourceGroup")
	autoscaleSettingsClientNewListByResourceGroupPager := autoscaleSettingsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for autoscaleSettingsClientNewListByResourceGroupPager.More() {
		_, err := autoscaleSettingsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AutoscaleSettings_Get
	fmt.Println("Call operation: AutoscaleSettings_Get")
	_, err = autoscaleSettingsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.autoscaleSettingName, nil)
	testsuite.Require().NoError(err)

	// From step AutoscaleSettings_Delete
	fmt.Println("Call operation: AutoscaleSettings_Delete")
	_, err = autoscaleSettingsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.autoscaleSettingName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Insights/eventcategories
func (testsuite *AutoscaleTestSuite) TestEventcategories() {
	var err error
	// From step EventCategories_List
	fmt.Println("Call operation: EventCategories_List")
	eventCategoriesClient, err := armmonitor.NewEventCategoriesClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	eventCategoriesClientNewListPager := eventCategoriesClient.NewListPager(nil)
	for eventCategoriesClientNewListPager.More() {
		_, err := eventCategoriesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Insights/operations
func (testsuite *AutoscaleTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armmonitor.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
}
