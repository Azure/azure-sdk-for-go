// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdatalakestore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datalake-store/armdatalakestore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	accountName            string
	firewallRuleName       string
	trustedIdProviderName  string
	virtualNetworkRuleName string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *AccountTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 8+6, true)
	testsuite.firewallRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "firewall", 8+6, false)
	testsuite.trustedIdProviderName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "trustedi", 8+6, false)
	testsuite.virtualNetworkRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualn", 8+6, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "centralus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *AccountTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func (testsuite *AccountTestSuite) Prepare() {
	var err error
	// From step Accounts_CheckNameAvailability
	fmt.Println("Call operation: Accounts_CheckNameAvailability")
	accountsClient, err := armdatalakestore.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountsClient.CheckNameAvailability(testsuite.ctx, testsuite.location, armdatalakestore.CheckNameAvailabilityParameters{
		Name: to.Ptr(testsuite.accountName),
		Type: to.Ptr(armdatalakestore.CheckNameAvailabilityParametersTypeMicrosoftDataLakeStoreAccounts),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_Create
	fmt.Println("Call operation: Accounts_Create")
	accountsClientCreateResponsePoller, err := accountsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armdatalakestore.CreateDataLakeStoreAccountParameters{
		Identity: &armdatalakestore.EncryptionIdentity{
			Type: to.Ptr("SystemAssigned"),
		},
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"test_key": to.Ptr("test_value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, accountsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DataLakeStore/accounts/{accountName}
func (testsuite *AccountTestSuite) TestAccounts() {
	var err error
	// From step Accounts_List
	fmt.Println("Call operation: Accounts_List")
	accountsClient, err := armdatalakestore.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountsClientNewListPager := accountsClient.NewListPager(&armdatalakestore.AccountsClientListOptions{Filter: nil,
		Top:     nil,
		Skip:    nil,
		Select:  nil,
		Orderby: nil,
		Count:   nil,
	})
	for accountsClientNewListPager.More() {
		_, err := accountsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Accounts_Get
	fmt.Println("Call operation: Accounts_Get")
	_, err = accountsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_ListByResourceGroup
	fmt.Println("Call operation: Accounts_ListByResourceGroup")
	accountsClientNewListByResourceGroupPager := accountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdatalakestore.AccountsClientListByResourceGroupOptions{Filter: nil,
		Top:     nil,
		Skip:    nil,
		Select:  nil,
		Orderby: nil,
		Count:   nil,
	})
	for accountsClientNewListByResourceGroupPager.More() {
		_, err := accountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Accounts_Update
	fmt.Println("Call operation: Accounts_Update")
	accountsClientUpdateResponsePoller, err := accountsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armdatalakestore.UpdateDataLakeStoreAccountParameters{
		Tags: map[string]*string{
			"test_key": to.Ptr("test_value"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, accountsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.DataLakeStore/accounts/{accountName}/virtualNetworkRules/{virtualNetworkRuleName}
func (testsuite *AccountTestSuite) TestVirtualNetworkRules() {
	var subnetId string
	var err error
	// From step Create_VirtualNetworkAndSubnet
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworksName'), 'default')]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworksName": map[string]any{
				"type":         "string",
				"defaultValue": "datalakestorevnet",
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworksName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2021-05-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
								"serviceEndpoints": []any{
									map[string]any{
										"service": "Microsoft.AzureActiveDirectory",
										"locations": []any{
											"[parameters('location')]",
										},
									},
								},
							},
						},
					},
				},
				"tags": map[string]any{},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetworkAndSubnet", &deployment)
	testsuite.Require().NoError(err)
	subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)

	// From step VirtualNetworkRules_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkRules_CreateOrUpdate")
	virtualNetworkRulesClient, err := armdatalakestore.NewVirtualNetworkRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = virtualNetworkRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.virtualNetworkRuleName, armdatalakestore.CreateOrUpdateVirtualNetworkRuleParameters{
		Properties: &armdatalakestore.CreateOrUpdateVirtualNetworkRuleProperties{
			SubnetID: to.Ptr(subnetId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_ListByAccount
	fmt.Println("Call operation: VirtualNetworkRules_ListByAccount")
	virtualNetworkRulesClientNewListByAccountPager := virtualNetworkRulesClient.NewListByAccountPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for virtualNetworkRulesClientNewListByAccountPager.More() {
		_, err := virtualNetworkRulesClientNewListByAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkRules_Get
	fmt.Println("Call operation: VirtualNetworkRules_Get")
	_, err = virtualNetworkRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.virtualNetworkRuleName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_Update
	fmt.Println("Call operation: VirtualNetworkRules_Update")
	_, err = virtualNetworkRulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.virtualNetworkRuleName, &armdatalakestore.VirtualNetworkRulesClientUpdateOptions{
		Parameters: &armdatalakestore.UpdateVirtualNetworkRuleParameters{
			Properties: &armdatalakestore.UpdateVirtualNetworkRuleProperties{
				SubnetID: to.Ptr(subnetId),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step VirtualNetworkRules_Delete
	fmt.Println("Call operation: VirtualNetworkRules_Delete")
	_, err = virtualNetworkRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.virtualNetworkRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataLakeStore/accounts/{accountName}/firewallRules/{firewallRuleName}
func (testsuite *AccountTestSuite) TestFirewallRules() {
	var err error
	// From step FirewallRules_CreateOrUpdate
	fmt.Println("Call operation: FirewallRules_CreateOrUpdate")
	firewallRulesClient, err := armdatalakestore.NewFirewallRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = firewallRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.firewallRuleName, armdatalakestore.CreateOrUpdateFirewallRuleParameters{
		Properties: &armdatalakestore.CreateOrUpdateFirewallRuleProperties{
			EndIPAddress:   to.Ptr("2.2.2.2"),
			StartIPAddress: to.Ptr("1.1.1.1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_ListByAccount
	fmt.Println("Call operation: FirewallRules_ListByAccount")
	firewallRulesClientNewListByAccountPager := firewallRulesClient.NewListByAccountPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for firewallRulesClientNewListByAccountPager.More() {
		_, err := firewallRulesClientNewListByAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallRules_Get
	fmt.Println("Call operation: FirewallRules_Get")
	_, err = firewallRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)

	// From step FirewallRules_Update
	fmt.Println("Call operation: FirewallRules_Update")
	_, err = firewallRulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.firewallRuleName, &armdatalakestore.FirewallRulesClientUpdateOptions{
		Parameters: &armdatalakestore.UpdateFirewallRuleParameters{
			Properties: &armdatalakestore.UpdateFirewallRuleProperties{
				EndIPAddress:   to.Ptr("2.2.2.2"),
				StartIPAddress: to.Ptr("1.1.1.1"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step FirewallRules_Delete
	fmt.Println("Call operation: FirewallRules_Delete")
	_, err = firewallRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.firewallRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders/{trustedIdProviderName}
func (testsuite *AccountTestSuite) TestTrustedIdProviders() {
	var err error
	// From step TrustedIdProviders_CreateOrUpdate
	fmt.Println("Call operation: TrustedIdProviders_CreateOrUpdate")
	trustedIDProvidersClient, err := armdatalakestore.NewTrustedIDProvidersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = trustedIDProvidersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.trustedIdProviderName, armdatalakestore.CreateOrUpdateTrustedIDProviderParameters{
		Properties: &armdatalakestore.CreateOrUpdateTrustedIDProviderProperties{
			IDProvider: to.Ptr("https://sts.windows.net/ea9ec534-a3e3-4e45-ad36-3afc5bb291c1"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step TrustedIdProviders_ListByAccount
	fmt.Println("Call operation: TrustedIdProviders_ListByAccount")
	trustedIDProvidersClientNewListByAccountPager := trustedIDProvidersClient.NewListByAccountPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for trustedIDProvidersClientNewListByAccountPager.More() {
		_, err := trustedIDProvidersClientNewListByAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TrustedIdProviders_Get
	fmt.Println("Call operation: TrustedIdProviders_Get")
	_, err = trustedIDProvidersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.trustedIdProviderName, nil)
	testsuite.Require().NoError(err)

	// From step TrustedIdProviders_Update
	fmt.Println("Call operation: TrustedIdProviders_Update")
	_, err = trustedIDProvidersClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.trustedIdProviderName, &armdatalakestore.TrustedIDProvidersClientUpdateOptions{
		Parameters: &armdatalakestore.UpdateTrustedIDProviderParameters{
			Properties: &armdatalakestore.UpdateTrustedIDProviderProperties{
				IDProvider: to.Ptr("https://sts.windows.net/ea9ec534-a3e3-4e45-ad36-3afc5bb291c1"),
			},
		},
	})
	testsuite.Require().NoError(err)

	// From step TrustedIdProviders_Delete
	fmt.Println("Call operation: TrustedIdProviders_Delete")
	_, err = trustedIDProvidersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.trustedIdProviderName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.DataLakeStore/operations
func (testsuite *AccountTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armdatalakestore.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = operationsClient.List(testsuite.ctx, nil)
	testsuite.Require().NoError(err)

	// From step Locations_GetCapability
	fmt.Println("Call operation: Locations_GetCapability")
	locationsClient, err := armdatalakestore.NewLocationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = locationsClient.GetCapability(testsuite.ctx, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// From step Locations_GetUsage
	fmt.Println("Call operation: Locations_GetUsage")
	locationsClientNewGetUsagePager := locationsClient.NewGetUsagePager(testsuite.location, nil)
	for locationsClientNewGetUsagePager.More() {
		_, err := locationsClientNewGetUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *AccountTestSuite) Cleanup() {
	var err error
	// From step Accounts_Delete
	fmt.Println("Call operation: Accounts_Delete")
	accountsClient, err := armdatalakestore.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountsClientDeleteResponsePoller, err := accountsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, accountsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
