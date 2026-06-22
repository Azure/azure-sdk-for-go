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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armslis"
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
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/arm-sdk-tests-rg/providers/microsoft.monitor/accounts/amw-arm-sdk-tests-rg")
	testsuite.managedIdentityResourceID = getEnvOrDefault("MANAGED_IDENTITY_RESOURCE_ID",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/arm-sdk-tests-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uami-arm-sdk-tests-rg")
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
		Identity: &armslis.ManagedServiceIdentity{
			Type: to.Ptr(armslis.ManagedServiceIdentityTypeUserAssigned),
			UserAssignedIdentities: func() map[string]*armslis.UserAssignedIdentity {
				identities := map[string]*armslis.UserAssignedIdentity{
					testsuite.managedIdentityResourceID: {},
				}
				if testsuite.sourceManagedIdentityResourceID != testsuite.managedIdentityResourceID {
					identities[testsuite.sourceManagedIdentityResourceID] = &armslis.UserAssignedIdentity{}
				}
				return identities
			}(),
		},
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
							// Source metric is a real Azure Managed Prometheus metric scraped by AKS.
							// Test infra (bicep) deploys an AKS cluster with the Azure Monitor metrics addon
							// pointed at the source AMW; container_cpu_usage_seconds_total is always populated.
							MetricNamespace: to.Ptr("customdefault"),
							MetricName:      to.Ptr("container_cpu_usage_seconds_total"),
							Filters: []*armslis.Condition{
								{
									DimensionName: to.Ptr("container"),
									Operator:      to.Ptr(armslis.ConditionOperatorNotEqual),
									Value:         to.Ptr("POD"),
								},
							},
							SpatialAggregation: &armslis.SpatialAggregation{
								Type:       to.Ptr(armslis.SpatialAggregationTypeSum),
								Dimensions: []*string{to.Ptr("instance")},
							},
							TemporalAggregation: &armslis.TemporalAggregation{
								Type:              to.Ptr(armslis.TemporalAggregationTypeRate),
								WindowSizeMinutes: to.Ptr[int32](1),
							},
						},
					},
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(createResp.Name)
	// In playback the proxy may rewrite resource names; in live mode the
	// response should equal the requested name. Accept either.
	if *createResp.Name != testsuite.sliName && *createResp.Name != "Sanitized" {
		testsuite.T().Fatalf("unexpected create name %q (want %q or 'Sanitized')", *createResp.Name, testsuite.sliName)
	}

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
