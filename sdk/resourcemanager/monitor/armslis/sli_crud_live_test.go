// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armslis_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armslis"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type SliCrudTestSuite struct {
	suite.Suite

	ctx                             context.Context
	cred                            azcore.TokenCredential
	options                         *arm.ClientOptions
	serviceGroupName                string
	sliName                         string
	amwResourceID                   string
	managedIdentityResourceID       string
	sourceAmwResourceID             string
	sourceManagedIdentityResourceID string
}

func (testsuite *SliCrudTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.serviceGroupName = getEnvOrDefault("SERVICE_GROUP_NAME", "arm-sdk-tests-sg")
	testsuite.sliName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "gosli", 12, true)
	testsuite.amwResourceID = getEnvOrDefault("AMW_RESOURCE_ID",
		"/subscriptions/6820e35f-0fe6-4af3-aad2-27414fa82621/resourceGroups/mfrei/providers/microsoft.monitor/accounts/streaming-3p-slo-am2cbn-eastus2euap-1")
	testsuite.managedIdentityResourceID = getEnvOrDefault("MANAGED_IDENTITY_RESOURCE_ID",
		"/subscriptions/6820e35f-0fe6-4af3-aad2-27414fa82621/resourceGroups/mfrei/providers/Microsoft.ManagedIdentity/userAssignedIdentities/mfrei-test-user-managed-identity")
	testsuite.sourceAmwResourceID = getEnvOrDefault("SOURCE_AMW_RESOURCE_ID", testsuite.amwResourceID)
	testsuite.sourceManagedIdentityResourceID = getEnvOrDefault("SOURCE_MANAGED_IDENTITY_RESOURCE_ID", testsuite.managedIdentityResourceID)
}

func (testsuite *SliCrudTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestSliCrudTestSuite(t *testing.T) {
	suite.Run(t, new(SliCrudTestSuite))
}

func (testsuite *SliCrudTestSuite) TestSliCrudLifecycle() {
	client, err := armslis.NewClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	// Step 1: Create SLI
	fmt.Println("Call operation: Slis_CreateOrUpdate")
	createResp, err := client.CreateOrUpdate(testsuite.ctx, testsuite.serviceGroupName, testsuite.sliName, armslis.Sli{
		Properties: &armslis.SliResource{
			Description:    to.Ptr("Live test SLI - measures latency of test API"),
			Category:       to.Ptr(armslis.CategoryLatency),
			EvaluationType: to.Ptr(armslis.EvaluationTypeWindowBased),
			EnableAlert:    to.Ptr(true),
			DestinationAmwAccounts: []*armslis.AmwAccount{
				{
					ResourceID: to.Ptr(testsuite.amwResourceID),
					Identity:   to.Ptr(testsuite.managedIdentityResourceID),
				},
			},
			BaselineProperties: &armslis.BaselineProperties{
				Baseline: &armslis.Baseline{
					Value:                     to.Ptr[float32](99),
					EvaluationPeriodDays:      to.Ptr[int32](30),
					EvaluationCalculationType: to.Ptr(armslis.EvaluationCalculationTypeCalendarDays),
				},
			},
			SliProperties: &armslis.SliProperties{
				WindowUptimeCriteria: &armslis.WindowUptimeCriteria{
					Target:     to.Ptr[float32](95),
					Comparator: to.Ptr(armslis.WindowUptimeCriteriaComparatorGreaterThanOrEqual),
				},
				Signals: &armslis.Signal{
					SignalFormula: to.Ptr("A"),
					SignalSources: []*armslis.SignalSource{
						{
							SignalSourceID:                  to.Ptr("A"),
							SourceAmwAccountManagedIdentity: to.Ptr(testsuite.sourceManagedIdentityResourceID),
							SourceAmwAccountResourceID:      to.Ptr(testsuite.sourceAmwResourceID),
							MetricNamespace:                 to.Ptr("TestMetrics"),
							MetricName:                      to.Ptr("TestLatency"),
							Filters: []*armslis.Condition{
								{
									DimensionName: to.Ptr("ApiName"),
									Operator:      to.Ptr(armslis.ConditionOperatorEqual),
									Value:         to.Ptr("TestApi"),
								},
							},
							SpatialAggregation: &armslis.SpatialAggregation{
								Type:       to.Ptr(armslis.SpatialAggregationTypeAverage),
								Dimensions: []*string{to.Ptr("Region")},
							},
							TemporalAggregation: &armslis.TemporalAggregation{
								Type:              to.Ptr(armslis.TemporalAggregationTypeAverage),
								WindowSizeMinutes: to.Ptr[int32](5),
							},
						},
					},
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(createResp.Name)
	testsuite.Equal(testsuite.sliName, *createResp.Name)

	// Step 2: Get SLI - verify it exists
	fmt.Println("Call operation: Slis_Get")
	getResp, err := client.Get(testsuite.ctx, testsuite.serviceGroupName, testsuite.sliName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(getResp.Properties)
	testsuite.Require().NotNil(getResp.Properties.Category)
	testsuite.Equal(armslis.CategoryLatency, *getResp.Properties.Category)

	// Step 3: Delete SLI
	fmt.Println("Call operation: Slis_Delete")
	_, err = client.Delete(testsuite.ctx, testsuite.serviceGroupName, testsuite.sliName, nil)
	testsuite.Require().NoError(err)

	// Step 4: Get SLI - expect 404
	fmt.Println("Call operation: Slis_Get (expect 404)")
	_, err = client.Get(testsuite.ctx, testsuite.serviceGroupName, testsuite.sliName, nil)
	testsuite.Require().Error(err)
	var respErr *azcore.ResponseError
	testsuite.Require().ErrorAs(err, &respErr)
	testsuite.Equal(404, respErr.StatusCode)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

const pathToPackage = "sdk/resourcemanager/monitor/armslis"
