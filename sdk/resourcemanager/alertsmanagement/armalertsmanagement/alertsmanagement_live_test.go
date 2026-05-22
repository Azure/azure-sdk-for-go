// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armalertsmanagement_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
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

func isInvalidResourceType(err error, resourceType string) bool {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return false
	}
	if !strings.EqualFold(respErr.ErrorCode, "InvalidResourceType") {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "resource type '"+strings.ToLower(resourceType)+"'")
}

// Microsoft.AlertsManagement/operations
func (testsuite *AlertsManagementTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armalertsmanagement.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	if operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		if isInvalidResourceType(err, "operations") {
			testsuite.T().Skipf("skipping operations live test due to API-version mismatch: %v", err)
		}
		testsuite.Require().NoError(err)
	}
}

// Microsoft.AlertsManagement/alerts/{alertId}
func (testsuite *AlertsManagementTestSuite) TestAlerts() {
	var err error
	// From step Alerts_GetAll
	fmt.Println("Call operation: Alerts_GetAll")
	alertsClient, err := armalertsmanagement.NewAlertsClient(fmt.Sprintf("subscriptions/%s", testsuite.subscriptionId), testsuite.cred, testsuite.options)
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
	if alertsClientNewGetAllPager.More() {
		_, err := alertsClientNewGetAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
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
	if isInvalidResourceType(err, "alertsSummary") {
		testsuite.T().Skipf("skipping alerts summary live test due to API-version mismatch: %v", err)
	}
	testsuite.Require().NoError(err)
}
