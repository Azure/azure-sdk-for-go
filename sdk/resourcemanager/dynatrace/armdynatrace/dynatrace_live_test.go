// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdynatrace_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dynatrace/armdynatrace/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type DynatraceTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	monitorName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DynatraceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.monitorName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "monitorn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DynatraceTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDynatraceTestSuite(t *testing.T) {
	suite.Run(t, new(DynatraceTestSuite))
}

func (testsuite *DynatraceTestSuite) Prepare() {
	var err error
	// From step Monitors_CreateOrUpdate
	fmt.Println("Call operation: Monitors_CreateOrUpdate")
	monitorsClient, err := armdynatrace.NewMonitorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	monitorsClientCreateOrUpdateResponsePoller, err := monitorsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, armdynatrace.MonitorResource{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"Environment": to.Ptr("Dev"),
		},
		Identity: &armdynatrace.IdentityProperties{
			Type: to.Ptr(armdynatrace.ManagedIdentityTypeSystemAssigned),
		},
		Properties: &armdynatrace.MonitorProperties{
			DynatraceEnvironmentProperties: &armdynatrace.EnvironmentProperties{
				AccountInfo:     &armdynatrace.AccountInfo{},
				EnvironmentInfo: &armdynatrace.EnvironmentInfo{},
				SingleSignOnProperties: &armdynatrace.SingleSignOnProperties{
					AADDomains: []*string{
						to.Ptr("http://www.contoso.com/")},
					SingleSignOnState: to.Ptr(armdynatrace.SingleSignOnStatesEnable),
				},
			},
			LiftrResourceCategory:         to.Ptr(armdynatrace.LiftrResourceCategoriesUnknown),
			MarketplaceSubscriptionStatus: to.Ptr(armdynatrace.MarketplaceSubscriptionStatusActive),
			MonitoringStatus:              to.Ptr(armdynatrace.MonitoringStatusEnabled),
			PlanData: &armdynatrace.PlanData{
				BillingCycle:  to.Ptr("Monthly"),
				EffectiveDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-08-30T15:14:33+02:00"); return t }()),
				PlanDetails:   to.Ptr("dynatraceapitestplan"),
				UsageType:     to.Ptr("Committed"),
			},
			ProvisioningState: to.Ptr(armdynatrace.ProvisioningStateAccepted),
			UserInfo: &armdynatrace.UserInfo{
				Country:      to.Ptr("westus2"),
				EmailAddress: to.Ptr("alice@microsoft.com"),
				FirstName:    to.Ptr("Alice"),
				LastName:     to.Ptr("Bobab"),
				PhoneNumber:  to.Ptr("123456"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, monitorsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Dynatrace.Observability/monitors/{monitorName}
func (testsuite *DynatraceTestSuite) TestMonitor() {
	var err error
	// From step Monitors_ListBySubscriptionId
	fmt.Println("Call operation: Monitors_ListBySubscriptionId")
	monitorsClient, err := armdynatrace.NewMonitorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	monitorsClientNewListBySubscriptionIDPager := monitorsClient.NewListBySubscriptionIDPager(nil)
	for monitorsClientNewListBySubscriptionIDPager.More() {
		_, err := monitorsClientNewListBySubscriptionIDPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Monitors_Get
	fmt.Println("Call operation: Monitors_Get")
	_, err = monitorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, nil)
	testsuite.Require().NoError(err)

	// From step Monitors_ListByResourceGroup
	fmt.Println("Call operation: Monitors_ListByResourceGroup")
	monitorsClientNewListByResourceGroupPager := monitorsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for monitorsClientNewListByResourceGroupPager.More() {
		_, err := monitorsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Monitors_Update
	fmt.Println("Call operation: Monitors_Update")
	_, err = monitorsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, armdynatrace.MonitorResourceUpdate{
		Tags: map[string]*string{
			"Environment": to.Ptr("Dev"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Monitors_ListMonitoredResources
	fmt.Println("Call operation: Monitors_ListMonitoredResources")
	monitorsClientNewListMonitoredResourcesPager := monitorsClient.NewListMonitoredResourcesPager(testsuite.resourceGroupName, testsuite.monitorName, nil)
	for monitorsClientNewListMonitoredResourcesPager.More() {
		_, err := monitorsClientNewListMonitoredResourcesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Monitors_ListHosts
	fmt.Println("Call operation: Monitors_ListHosts")
	monitorsClientNewListHostsPager := monitorsClient.NewListHostsPager(testsuite.resourceGroupName, testsuite.monitorName, nil)
	for monitorsClientNewListHostsPager.More() {
		_, err := monitorsClientNewListHostsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Monitors_GetVMHostPayload
	fmt.Println("Call operation: Monitors_GetVMHostPayload")
	_, err = monitorsClient.GetVMHostPayload(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, nil)
	testsuite.Require().NoError(err)

	// From step Monitors_ListAppServices
	fmt.Println("Call operation: Monitors_ListAppServices")
	monitorsClientNewListAppServicesPager := monitorsClient.NewListAppServicesPager(testsuite.resourceGroupName, testsuite.monitorName, nil)
	for monitorsClientNewListAppServicesPager.More() {
		_, err := monitorsClientNewListAppServicesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Monitors_GetMetricStatus
	fmt.Println("Call operation: Monitors_GetMetricStatus")
	_, err = monitorsClient.GetMetricStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, nil)
	testsuite.Require().NoError(err)
}

// Dynatrace.Observability/monitors/{monitorName}/tagRules/{ruleSetName}
func (testsuite *DynatraceTestSuite) TestTagRules() {
	var err error
	// From step TagRules_CreateOrUpdate
	fmt.Println("Call operation: TagRules_CreateOrUpdate")
	tagRulesClient, err := armdynatrace.NewTagRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	tagRulesClientCreateOrUpdateResponsePoller, err := tagRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, "default", armdynatrace.TagRule{
		Properties: &armdynatrace.MonitoringTagRulesProperties{
			LogRules: &armdynatrace.LogRules{
				FilteringTags: []*armdynatrace.FilteringTag{
					{
						Name:   to.Ptr("Environment"),
						Action: to.Ptr(armdynatrace.TagActionInclude),
						Value:  to.Ptr("Prod"),
					},
					{
						Name:   to.Ptr("Environment"),
						Action: to.Ptr(armdynatrace.TagActionExclude),
						Value:  to.Ptr("Dev"),
					}},
				SendAADLogs:          to.Ptr(armdynatrace.SendAADLogsStatusEnabled),
				SendActivityLogs:     to.Ptr(armdynatrace.SendActivityLogsStatusEnabled),
				SendSubscriptionLogs: to.Ptr(armdynatrace.SendSubscriptionLogsStatusEnabled),
			},
			MetricRules: &armdynatrace.MetricRules{
				FilteringTags: []*armdynatrace.FilteringTag{
					{
						Name:   to.Ptr("Environment"),
						Action: to.Ptr(armdynatrace.TagActionInclude),
						Value:  to.Ptr("Prod"),
					}},
				SendingMetrics: to.Ptr(armdynatrace.SendingMetricsStatusEnabled),
			},
			ProvisioningState: to.Ptr(armdynatrace.ProvisioningStateAccepted),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tagRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step TagRules_List
	fmt.Println("Call operation: TagRules_List")
	tagRulesClientNewListPager := tagRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.monitorName, nil)
	for tagRulesClientNewListPager.More() {
		_, err := tagRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step TagRules_Get
	fmt.Println("Call operation: TagRules_Get")
	_, err = tagRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, "default", nil)
	testsuite.Require().NoError(err)

	// From step TagRules_Delete
	fmt.Println("Call operation: TagRules_Delete")
	tagRulesClientDeleteResponsePoller, err := tagRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, "default", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, tagRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Dynatrace.Observability/operations
func (testsuite *DynatraceTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armdynatrace.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *DynatraceTestSuite) Cleanup() {
	var err error
	// From step Monitors_Delete
	fmt.Println("Call operation: Monitors_Delete")
	monitorsClient, err := armdynatrace.NewMonitorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	monitorsClientDeleteResponsePoller, err := monitorsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.monitorName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, monitorsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
