package consumptionapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2018-10-01/consumption"
	"github.com/Azure/go-autorest/autorest"
)

// UsageDetailsClientAPI contains the set of methods on the UsageDetailsClient type.
type UsageDetailsClientAPI interface {
	List(ctx context.Context, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultPage, err error)
	ListComplete(ctx context.Context, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultIterator, err error)
	ListByBillingAccount(ctx context.Context, billingAccountID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultPage, err error)
	ListByBillingAccountComplete(ctx context.Context, billingAccountID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultIterator, err error)
	ListByBillingPeriod(ctx context.Context, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultPage, err error)
	ListByBillingPeriodComplete(ctx context.Context, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultIterator, err error)
	ListByDepartment(ctx context.Context, departmentID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultPage, err error)
	ListByDepartmentComplete(ctx context.Context, departmentID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultIterator, err error)
	ListByEnrollmentAccount(ctx context.Context, enrollmentAccountID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultPage, err error)
	ListByEnrollmentAccountComplete(ctx context.Context, enrollmentAccountID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultIterator, err error)
	ListByManagementGroup(ctx context.Context, managementGroupID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultPage, err error)
	ListByManagementGroupComplete(ctx context.Context, managementGroupID string, expand string, filter string, skiptoken string, top *int32, apply string) (result consumption.UsageDetailsListResultIterator, err error)
	ListForBillingPeriodByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultPage, err error)
	ListForBillingPeriodByBillingAccountComplete(ctx context.Context, billingAccountID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultIterator, err error)
	ListForBillingPeriodByDepartment(ctx context.Context, departmentID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultPage, err error)
	ListForBillingPeriodByDepartmentComplete(ctx context.Context, departmentID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultIterator, err error)
	ListForBillingPeriodByEnrollmentAccount(ctx context.Context, enrollmentAccountID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultPage, err error)
	ListForBillingPeriodByEnrollmentAccountComplete(ctx context.Context, enrollmentAccountID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultIterator, err error)
	ListForBillingPeriodByManagementGroup(ctx context.Context, managementGroupID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultPage, err error)
	ListForBillingPeriodByManagementGroupComplete(ctx context.Context, managementGroupID string, billingPeriodName string, expand string, filter string, apply string, skiptoken string, top *int32) (result consumption.UsageDetailsListResultIterator, err error)
}

var _ UsageDetailsClientAPI = (*consumption.UsageDetailsClient)(nil)

// MarketplacesClientAPI contains the set of methods on the MarketplacesClient type.
type MarketplacesClientAPI interface {
	List(ctx context.Context, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListComplete(ctx context.Context, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListByBillingAccount(ctx context.Context, billingAccountID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListByBillingAccountComplete(ctx context.Context, billingAccountID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListByBillingPeriod(ctx context.Context, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListByBillingPeriodComplete(ctx context.Context, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListByDepartment(ctx context.Context, departmentID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListByDepartmentComplete(ctx context.Context, departmentID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListByEnrollmentAccount(ctx context.Context, enrollmentAccountID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListByEnrollmentAccountComplete(ctx context.Context, enrollmentAccountID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListByManagementGroup(ctx context.Context, managementGroupID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListByManagementGroupComplete(ctx context.Context, managementGroupID string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListForBillingPeriodByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListForBillingPeriodByBillingAccountComplete(ctx context.Context, billingAccountID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListForBillingPeriodByDepartment(ctx context.Context, departmentID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListForBillingPeriodByDepartmentComplete(ctx context.Context, departmentID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListForBillingPeriodByEnrollmentAccount(ctx context.Context, enrollmentAccountID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListForBillingPeriodByEnrollmentAccountComplete(ctx context.Context, enrollmentAccountID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
	ListForBillingPeriodByManagementGroup(ctx context.Context, managementGroupID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultPage, err error)
	ListForBillingPeriodByManagementGroupComplete(ctx context.Context, managementGroupID string, billingPeriodName string, filter string, top *int32, skiptoken string) (result consumption.MarketplacesListResultIterator, err error)
}

var _ MarketplacesClientAPI = (*consumption.MarketplacesClient)(nil)

// BalancesClientAPI contains the set of methods on the BalancesClient type.
type BalancesClientAPI interface {
	GetByBillingAccount(ctx context.Context, billingAccountID string) (result consumption.Balance, err error)
	GetForBillingPeriodByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string) (result consumption.Balance, err error)
}

var _ BalancesClientAPI = (*consumption.BalancesClient)(nil)

// ReservationsSummariesClientAPI contains the set of methods on the ReservationsSummariesClient type.
type ReservationsSummariesClientAPI interface {
	ListByReservationOrder(ctx context.Context, reservationOrderID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultPage, err error)
	ListByReservationOrderComplete(ctx context.Context, reservationOrderID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultIterator, err error)
	ListByReservationOrderAndReservation(ctx context.Context, reservationOrderID string, reservationID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultPage, err error)
	ListByReservationOrderAndReservationComplete(ctx context.Context, reservationOrderID string, reservationID string, grain consumption.Datagrain, filter string) (result consumption.ReservationSummariesListResultIterator, err error)
}

var _ ReservationsSummariesClientAPI = (*consumption.ReservationsSummariesClient)(nil)

// ReservationsDetailsClientAPI contains the set of methods on the ReservationsDetailsClient type.
type ReservationsDetailsClientAPI interface {
	ListByReservationOrder(ctx context.Context, reservationOrderID string, filter string) (result consumption.ReservationDetailsListResultPage, err error)
	ListByReservationOrderComplete(ctx context.Context, reservationOrderID string, filter string) (result consumption.ReservationDetailsListResultIterator, err error)
	ListByReservationOrderAndReservation(ctx context.Context, reservationOrderID string, reservationID string, filter string) (result consumption.ReservationDetailsListResultPage, err error)
	ListByReservationOrderAndReservationComplete(ctx context.Context, reservationOrderID string, reservationID string, filter string) (result consumption.ReservationDetailsListResultIterator, err error)
}

var _ ReservationsDetailsClientAPI = (*consumption.ReservationsDetailsClient)(nil)

// ReservationRecommendationsClientAPI contains the set of methods on the ReservationRecommendationsClient type.
type ReservationRecommendationsClientAPI interface {
	List(ctx context.Context, filter string) (result consumption.ReservationRecommendationsListResultPage, err error)
	ListComplete(ctx context.Context, filter string) (result consumption.ReservationRecommendationsListResultIterator, err error)
}

var _ ReservationRecommendationsClientAPI = (*consumption.ReservationRecommendationsClient)(nil)

// BudgetsClientAPI contains the set of methods on the BudgetsClient type.
type BudgetsClientAPI interface {
	CreateOrUpdate(ctx context.Context, budgetName string, parameters consumption.Budget) (result consumption.Budget, err error)
	CreateOrUpdateByResourceGroupName(ctx context.Context, resourceGroupName string, budgetName string, parameters consumption.Budget) (result consumption.Budget, err error)
	Delete(ctx context.Context, budgetName string) (result autorest.Response, err error)
	DeleteByResourceGroupName(ctx context.Context, resourceGroupName string, budgetName string) (result autorest.Response, err error)
	Get(ctx context.Context, budgetName string) (result consumption.Budget, err error)
	GetByResourceGroupName(ctx context.Context, resourceGroupName string, budgetName string) (result consumption.Budget, err error)
	List(ctx context.Context) (result consumption.BudgetsListResultPage, err error)
	ListComplete(ctx context.Context) (result consumption.BudgetsListResultIterator, err error)
	ListByResourceGroupName(ctx context.Context, resourceGroupName string) (result consumption.BudgetsListResultPage, err error)
	ListByResourceGroupNameComplete(ctx context.Context, resourceGroupName string) (result consumption.BudgetsListResultIterator, err error)
}

var _ BudgetsClientAPI = (*consumption.BudgetsClient)(nil)

// PriceSheetClientAPI contains the set of methods on the PriceSheetClient type.
type PriceSheetClientAPI interface {
	Get(ctx context.Context, expand string, skiptoken string, top *int32) (result consumption.PriceSheetResult, err error)
	GetByBillingPeriod(ctx context.Context, billingPeriodName string, expand string, skiptoken string, top *int32) (result consumption.PriceSheetResult, err error)
}

var _ PriceSheetClientAPI = (*consumption.PriceSheetClient)(nil)

// TagsClientAPI contains the set of methods on the TagsClient type.
type TagsClientAPI interface {
	Get(ctx context.Context, billingAccountID string) (result consumption.TagsResult, err error)
}

var _ TagsClientAPI = (*consumption.TagsClient)(nil)

// ForecastsClientAPI contains the set of methods on the ForecastsClient type.
type ForecastsClientAPI interface {
	List(ctx context.Context, filter string) (result consumption.ForecastsListResult, err error)
}

var _ ForecastsClientAPI = (*consumption.ForecastsClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result consumption.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result consumption.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*consumption.OperationsClient)(nil)

// AggregatedCostClientAPI contains the set of methods on the AggregatedCostClient type.
type AggregatedCostClientAPI interface {
	GetByManagementGroup(ctx context.Context, managementGroupID string, filter string) (result consumption.ManagementGroupAggregatedCostResult, err error)
	GetForBillingPeriodByManagementGroup(ctx context.Context, managementGroupID string, billingPeriodName string) (result consumption.ManagementGroupAggregatedCostResult, err error)
}

var _ AggregatedCostClientAPI = (*consumption.AggregatedCostClient)(nil)

// ChargesClientAPI contains the set of methods on the ChargesClient type.
type ChargesClientAPI interface {
	ListByDepartment(ctx context.Context, billingAccountID string, departmentID string, filter string) (result consumption.ChargesListResult, err error)
	ListByEnrollmentAccount(ctx context.Context, billingAccountID string, enrollmentAccountID string, filter string) (result consumption.ChargesListResult, err error)
	ListForBillingPeriodByDepartment(ctx context.Context, billingAccountID string, departmentID string, billingPeriodName string, filter string) (result consumption.ChargeSummary, err error)
	ListForBillingPeriodByEnrollmentAccount(ctx context.Context, billingAccountID string, enrollmentAccountID string, billingPeriodName string, filter string) (result consumption.ChargeSummary, err error)
}

var _ ChargesClientAPI = (*consumption.ChargesClient)(nil)

// TenantsClientAPI contains the set of methods on the TenantsClient type.
type TenantsClientAPI interface {
	Get(ctx context.Context, billingAccountID string, billingProfileID string) (result consumption.TenantListResult, err error)
}

var _ TenantsClientAPI = (*consumption.TenantsClient)(nil)
