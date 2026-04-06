// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armalertsmanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type AlertsManagementTestSuite struct {
	suite.Suite

	ctx                     context.Context
	cred                    azcore.TokenCredential
	options                 *arm.ClientOptions
	alertProcessingRuleName string
	location                string
	resourceGroupName       string
	subscriptionId          string
}

func (testsuite *AlertsManagementTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.alertProcessingRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "alertpro", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AlertsManagementTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAlertsManagementTestSuite(t *testing.T) {
	suite.Run(t, new(AlertsManagementTestSuite))
}

// Microsoft.AlertsManagement/operations
func (testsuite *AlertsManagementTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armalertsmanagement.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.AlertsManagement/alerts/{alertId}
func (testsuite *AlertsManagementTestSuite) TestAlerts() {
	var err error
	// From step Alerts_GetAll
	fmt.Println("Call operation: Alerts_GetAll")
	alertsClient, err := armalertsmanagement.NewAlertsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	alertsClientNewGetAllPager := alertsClient.NewGetAllPager(&armalertsmanagement.AlertsClientGetAllOptions{TargetResource: nil,
		TargetResourceType:  nil,
		TargetResourceGroup: nil,
		MonitorService:      nil,
		MonitorCondition:    nil,
		Severity:            nil,
		AlertState:          nil,
		AlertRule:           nil,
		SmartGroupID:        nil,
		IncludeContext:      nil,
		IncludeEgressConfig: nil,
		PageCount:           nil,
		SortBy:              nil,
		SortOrder:           nil,
		Select:              nil,
		TimeRange:           nil,
		CustomTimeRange:     nil,
	})
	for alertsClientNewGetAllPager.More() {
		_, err := alertsClientNewGetAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Alerts_GetSummary
	fmt.Println("Call operation: Alerts_GetSummary")
	_, err = alertsClient.GetSummary(testsuite.ctx, armalertsmanagement.AlertsSummaryGroupByFields("severity,alertState"), &armalertsmanagement.AlertsClientGetSummaryOptions{IncludeSmartGroupsCount: nil,
		TargetResource:      nil,
		TargetResourceType:  nil,
		TargetResourceGroup: nil,
		MonitorService:      nil,
		MonitorCondition:    nil,
		Severity:            nil,
		AlertState:          nil,
		AlertRule:           nil,
		TimeRange:           nil,
		CustomTimeRange:     nil,
	})
	testsuite.Require().NoError(err)
}

// Microsoft.AlertsManagement/actionRules/{alertProcessingRuleName}
func (testsuite *AlertsManagementTestSuite) TestAlertProcessingRules() {
	globalLocation := "Global"
	var err error
	// From step AlertProcessingRules_CreateOrUpdate
	fmt.Println("Call operation: AlertProcessingRules_CreateOrUpdate")
	alertProcessingRulesClient, err := armalertsmanagement.NewAlertProcessingRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = alertProcessingRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.alertProcessingRuleName, armalertsmanagement.AlertProcessingRule{
		Location: to.Ptr(globalLocation),
		Tags:     map[string]*string{},
		Properties: &armalertsmanagement.AlertProcessingRuleProperties{
			Description: to.Ptr("Remove all ActionGroups outside business hours"),
			Actions: []armalertsmanagement.ActionClassification{
				&armalertsmanagement.RemoveAllActionGroups{
					ActionType: to.Ptr(armalertsmanagement.ActionTypeRemoveAllActionGroups),
				}},
			Enabled: to.Ptr(true),
			Schedule: &armalertsmanagement.Schedule{
				Recurrences: []armalertsmanagement.RecurrenceClassification{
					&armalertsmanagement.DailyRecurrence{
						EndTime:        to.Ptr("09:00:00"),
						RecurrenceType: to.Ptr(armalertsmanagement.RecurrenceTypeDaily),
						StartTime:      to.Ptr("17:00:00"),
					},
					&armalertsmanagement.WeeklyRecurrence{
						RecurrenceType: to.Ptr(armalertsmanagement.RecurrenceTypeWeekly),
						DaysOfWeek: []*armalertsmanagement.DaysOfWeek{
							to.Ptr(armalertsmanagement.DaysOfWeekSaturday),
							to.Ptr(armalertsmanagement.DaysOfWeekSunday)},
					}},
				TimeZone: to.Ptr("Eastern Standard Time"),
			},
			Scopes: []*string{
				to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AlertProcessingRules_ListBySubscription
	fmt.Println("Call operation: AlertProcessingRules_ListBySubscription")
	alertProcessingRulesClientNewListBySubscriptionPager := alertProcessingRulesClient.NewListBySubscriptionPager(nil)
	for alertProcessingRulesClientNewListBySubscriptionPager.More() {
		_, err := alertProcessingRulesClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AlertProcessingRules_GetByName
	fmt.Println("Call operation: AlertProcessingRules_GetByName")
	_, err = alertProcessingRulesClient.GetByName(testsuite.ctx, testsuite.resourceGroupName, testsuite.alertProcessingRuleName, nil)
	testsuite.Require().NoError(err)

	// From step AlertProcessingRules_ListByResourceGroup
	fmt.Println("Call operation: AlertProcessingRules_ListByResourceGroup")
	alertProcessingRulesClientNewListByResourceGroupPager := alertProcessingRulesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for alertProcessingRulesClientNewListByResourceGroupPager.More() {
		_, err := alertProcessingRulesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AlertProcessingRules_Update
	fmt.Println("Call operation: AlertProcessingRules_Update")
	_, err = alertProcessingRulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.alertProcessingRuleName, armalertsmanagement.PatchObject{
		Properties: &armalertsmanagement.PatchProperties{
			Enabled: to.Ptr(false),
		},
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AlertProcessingRules_Delete
	fmt.Println("Call operation: AlertProcessingRules_Delete")
	_, err = alertProcessingRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.alertProcessingRuleName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.AlertsManagement/smartGroups/{smartGroupId}
func (testsuite *AlertsManagementTestSuite) TestSmartGroups() {
	var err error
	// From step SmartGroups_GetAll
	fmt.Println("Call operation: SmartGroups_GetAll")
	smartGroupsClient, err := armalertsmanagement.NewSmartGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	smartGroupsClientNewGetAllPager := smartGroupsClient.NewGetAllPager(&armalertsmanagement.SmartGroupsClientGetAllOptions{TargetResource: nil,
		TargetResourceGroup: nil,
		TargetResourceType:  nil,
		MonitorService:      nil,
		MonitorCondition:    nil,
		Severity:            nil,
		SmartGroupState:     nil,
		TimeRange:           nil,
		PageCount:           nil,
		SortBy:              nil,
		SortOrder:           nil,
	})
	for smartGroupsClientNewGetAllPager.More() {
		_, err := smartGroupsClientNewGetAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}
