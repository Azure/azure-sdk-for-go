// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmonitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type MetricalertsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	ruleName          string
	subnetId          string
	virtualMachineId  string
	vnicId            string
	adminPassword     string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *MetricalertsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.ruleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "metricalertna", 19, false)
	testsuite.adminPassword = recording.GetEnvVariable("ADMIN_PASSWORD", "000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MetricalertsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMetricalertsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricalertsTestSuite))
}

func (testsuite *MetricalertsTestSuite) Prepare() {
	var err error
	// From step VirtualNetwork_Create
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
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "VirtualNetwork_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step NetworkInterface_Create
	template = map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"vnicId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/networkInterfaces', 'vmnic')]",
			},
		},
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
				"name":       "vmnic",
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
	deploymentExtend, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "NetworkInterface_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.vnicId = deploymentExtend.Properties.Outputs.(map[string]interface{})["vnicId"].(map[string]interface{})["value"].(string)

	// From step VirtualMachine_Create
	template = map[string]interface{}{
		"$schema":        "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"virtualMachineId": map[string]interface{}{
				"type":  "string",
				"value": "[resourceId('Microsoft.Compute/virtualMachines', 'scenariovm')]",
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
			"vnicId": map[string]interface{}{
				"type":         "string",
				"defaultValue": "$(vnicId)",
			},
		},
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "scenariovm",
				"type":       "Microsoft.Compute/virtualMachines",
				"apiVersion": "2022-03-01",
				"location":   "[parameters('location')]",
				"properties": map[string]interface{}{
					"diagnosticsProfile": map[string]interface{}{
						"bootDiagnostics": map[string]interface{}{
							"enabled": true,
						},
					},
					"hardwareProfile": map[string]interface{}{
						"vmSize": "Standard_DS1_v2",
					},
					"networkProfile": map[string]interface{}{
						"networkInterfaces": []interface{}{
							map[string]interface{}{
								"id": "[parameters('vnicId')]",
								"properties": map[string]interface{}{
									"deleteOption": "Detach",
									"primary":      true,
								},
							},
						},
					},
					"osProfile": map[string]interface{}{
						"adminPassword":            "[parameters('adminPassword')]",
						"adminUsername":            "azureuser",
						"allowExtensionOperations": true,
						"computerName":             "scenariovm",
						"secrets":                  []interface{}{},
						"windowsConfiguration": map[string]interface{}{
							"enableAutomaticUpdates": true,
							"patchSettings": map[string]interface{}{
								"assessmentMode":    "ImageDefault",
								"enableHotpatching": false,
								"patchMode":         "AutomaticByOS",
							},
							"provisionVMAgent": true,
						},
					},
					"storageProfile": map[string]interface{}{
						"dataDisks": []interface{}{},
						"imageReference": map[string]interface{}{
							"offer":     "WindowsServer",
							"publisher": "MicrosoftWindowsServer",
							"sku":       "2019-Datacenter",
							"version":   "latest",
						},
						"osDisk": map[string]interface{}{
							"name":         "[concat('scenariovm', 'osDisk')]",
							"caching":      "ReadWrite",
							"createOption": "FromImage",
							"deleteOption": "Delete",
							"diskSizeGB":   float64(127),
							"managedDisk": map[string]interface{}{
								"storageAccountType": "Premium_LRS",
							},
							"osType": "Windows",
						},
					},
				},
				"zones": []interface{}{
					"1",
				},
			},
		},
		"variables": map[string]interface{}{},
	}
	params = map[string]interface{}{
		"adminPassword": map[string]interface{}{"value": testsuite.adminPassword},
		"location":      map[string]interface{}{"value": testsuite.location},
		"vnicId":        map[string]interface{}{"value": testsuite.vnicId},
	}
	deployment = armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err = testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "VirtualMachine_Create", &deployment)
	testsuite.Require().NoError(err)
	testsuite.virtualMachineId = deploymentExtend.Properties.Outputs.(map[string]interface{})["virtualMachineId"].(map[string]interface{})["value"].(string)
}

// Microsoft.Insights/metricAlerts
func (testsuite *MetricalertsTestSuite) TestMetricalert() {
	var statusName string
	var err error
	// From step MetricAlerts_Create
	fmt.Println("Call operation: MetricAlerts_CreateOrUpdate")
	metricAlertsClient, err := armmonitor.NewMetricAlertsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = metricAlertsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, armmonitor.MetricAlertResource{
		Location: to.Ptr("global"),
		Properties: &armmonitor.MetricAlertProperties{
			Description:  to.Ptr("This is the description of the rule1"),
			AutoMitigate: to.Ptr(true),
			Criteria: &armmonitor.MetricAlertMultipleResourceMultipleMetricCriteria{
				ODataType: to.Ptr(armmonitor.OdatatypeMicrosoftAzureMonitorMultipleResourceMultipleMetricCriteria),
				AllOf: []armmonitor.MultiMetricCriteriaClassification{
					&armmonitor.DynamicMetricCriteria{
						Name:             to.Ptr("High_CPU_80"),
						CriterionType:    to.Ptr(armmonitor.CriterionTypeDynamicThresholdCriterion),
						MetricName:       to.Ptr("Percentage CPU"),
						MetricNamespace:  to.Ptr("microsoft.compute/virtualmachines"),
						TimeAggregation:  to.Ptr(armmonitor.AggregationTypeEnumAverage),
						AlertSensitivity: to.Ptr(armmonitor.DynamicThresholdSensitivityMedium),
						FailingPeriods: &armmonitor.DynamicThresholdFailingPeriods{
							MinFailingPeriodsToAlert:  to.Ptr[float32](4),
							NumberOfEvaluationPeriods: to.Ptr[float32](4),
						},
						Operator: to.Ptr(armmonitor.DynamicThresholdOperatorGreaterOrLessThan),
					}},
			},
			Enabled:             to.Ptr(true),
			EvaluationFrequency: to.Ptr("PT1M"),
			Scopes: []*string{
				to.Ptr(testsuite.virtualMachineId)},
			Severity:             to.Ptr[int32](3),
			TargetResourceRegion: to.Ptr(testsuite.location),
			TargetResourceType:   to.Ptr("Microsoft.Compute/virtualMachines"),
			WindowSize:           to.Ptr("PT15M"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step MetricAlerts_ListBySubscription
	fmt.Println("Call operation: MetricAlerts_ListBySubscription")
	metricAlertsClientNewListBySubscriptionPager := metricAlertsClient.NewListBySubscriptionPager(nil)
	for metricAlertsClientNewListBySubscriptionPager.More() {
		_, err := metricAlertsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MetricAlerts_ListByResourceGroup
	fmt.Println("Call operation: MetricAlerts_ListByResourceGroup")
	metricAlertsClientNewListByResourceGroupPager := metricAlertsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for metricAlertsClientNewListByResourceGroupPager.More() {
		_, err := metricAlertsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step MetricAlerts_Get
	fmt.Println("Call operation: MetricAlerts_Get")
	_, err = metricAlertsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)

	// From step MetricAlerts_Update
	fmt.Println("Call operation: MetricAlerts_Update")
	_, err = metricAlertsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, armmonitor.MetricAlertResourcePatch{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step MetricAlertsStatus_List
	fmt.Println("Call operation: MetricAlerts_List")
	metricAlertsStatusClient, err := armmonitor.NewMetricAlertsStatusClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	metricAlertsStatusClientListResponse, err := metricAlertsStatusClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)
	statusName = *metricAlertsStatusClientListResponse.Value[0].Name

	// From step MetricAlertsStatus_ListByName
	fmt.Println("Call operation: MetricAlerts_ListByName")
	_, err = metricAlertsStatusClient.ListByName(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, statusName, nil)
	testsuite.Require().NoError(err)

	// From step MetricAlerts_Delete
	fmt.Println("Call operation: MetricAlerts_Delete")
	_, err = metricAlertsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)
}
