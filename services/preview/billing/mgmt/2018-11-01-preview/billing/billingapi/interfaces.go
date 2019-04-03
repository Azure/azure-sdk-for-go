package billingapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2018-11-01-preview/billing"
	"github.com/Azure/go-autorest/autorest"
)

// BaseClientAPI contains the set of methods on the BaseClient type.
type BaseClientAPI interface {
	UpdateAutoRenewForBillingAccount(ctx context.Context, billingAccountName string, productName string, body billing.UpdateAutoRenewRequest) (result billing.UpdateAutoRenewOperationSummary, err error)
	UpdateAutoRenewForInvoiceSection(ctx context.Context, billingAccountName string, invoiceSectionName string, productName string, body billing.UpdateAutoRenewRequest) (result billing.UpdateAutoRenewOperationSummary, err error)
}

var _ BaseClientAPI = (*billing.BaseClient)(nil)

// AccountsClientAPI contains the set of methods on the AccountsClient type.
type AccountsClientAPI interface {
	Get(ctx context.Context, billingAccountName string, expand string) (result billing.Account, err error)
	List(ctx context.Context, expand string) (result billing.AccountListResult, err error)
}

var _ AccountsClientAPI = (*billing.AccountsClient)(nil)

// AccountsWithCreateInvoiceSectionPermissionClientAPI contains the set of methods on the AccountsWithCreateInvoiceSectionPermissionClient type.
type AccountsWithCreateInvoiceSectionPermissionClientAPI interface {
	List(ctx context.Context, expand string) (result billing.AccountListResult, err error)
}

var _ AccountsWithCreateInvoiceSectionPermissionClientAPI = (*billing.AccountsWithCreateInvoiceSectionPermissionClient)(nil)

// AvailableBalanceByBillingProfileClientAPI contains the set of methods on the AvailableBalanceByBillingProfileClient type.
type AvailableBalanceByBillingProfileClientAPI interface {
	Get(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.AvailableBalance, err error)
}

var _ AvailableBalanceByBillingProfileClientAPI = (*billing.AvailableBalanceByBillingProfileClient)(nil)

// PaymentMethodsByBillingProfileClientAPI contains the set of methods on the PaymentMethodsByBillingProfileClient type.
type PaymentMethodsByBillingProfileClientAPI interface {
	List(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.PaymentMethodsListResultPage, err error)
}

var _ PaymentMethodsByBillingProfileClientAPI = (*billing.PaymentMethodsByBillingProfileClient)(nil)

// ProfilesByBillingAccountNameClientAPI contains the set of methods on the ProfilesByBillingAccountNameClient type.
type ProfilesByBillingAccountNameClientAPI interface {
	List(ctx context.Context, billingAccountName string, expand string) (result billing.ProfileListResult, err error)
}

var _ ProfilesByBillingAccountNameClientAPI = (*billing.ProfilesByBillingAccountNameClient)(nil)

// ProfilesClientAPI contains the set of methods on the ProfilesClient type.
type ProfilesClientAPI interface {
	Get(ctx context.Context, billingAccountName string, billingProfileName string, expand string) (result billing.Profile, err error)
	Update(ctx context.Context, billingAccountName string, billingProfileName string, parameters billing.Profile) (result billing.ProfilesUpdateFuture, err error)
}

var _ ProfilesClientAPI = (*billing.ProfilesClient)(nil)

// InvoiceSectionsByBillingAccountNameClientAPI contains the set of methods on the InvoiceSectionsByBillingAccountNameClient type.
type InvoiceSectionsByBillingAccountNameClientAPI interface {
	List(ctx context.Context, billingAccountName string, expand string) (result billing.InvoiceSectionListResult, err error)
}

var _ InvoiceSectionsByBillingAccountNameClientAPI = (*billing.InvoiceSectionsByBillingAccountNameClient)(nil)

// InvoiceSectionsClientAPI contains the set of methods on the InvoiceSectionsClient type.
type InvoiceSectionsClientAPI interface {
	Create(ctx context.Context, billingAccountName string, parameters billing.InvoiceSectionProperties) (result billing.InvoiceSectionsCreateFuture, err error)
	ElevateToBillingProfile(ctx context.Context, billingAccountName string, invoiceSectionName string) (result autorest.Response, err error)
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, expand string) (result billing.InvoiceSection, err error)
	ListByBillingProfileName(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.InvoiceSectionListResult, err error)
	Update(ctx context.Context, billingAccountName string, invoiceSectionName string, parameters billing.InvoiceSection) (result billing.InvoiceSectionsUpdateFuture, err error)
}

var _ InvoiceSectionsClientAPI = (*billing.InvoiceSectionsClient)(nil)

// InvoiceSectionsWithCreateSubscriptionPermissionClientAPI contains the set of methods on the InvoiceSectionsWithCreateSubscriptionPermissionClient type.
type InvoiceSectionsWithCreateSubscriptionPermissionClientAPI interface {
	List(ctx context.Context, billingAccountName string, expand string) (result billing.InvoiceSectionListResult, err error)
}

var _ InvoiceSectionsWithCreateSubscriptionPermissionClientAPI = (*billing.InvoiceSectionsWithCreateSubscriptionPermissionClient)(nil)

// DepartmentsByBillingAccountNameClientAPI contains the set of methods on the DepartmentsByBillingAccountNameClient type.
type DepartmentsByBillingAccountNameClientAPI interface {
	List(ctx context.Context, billingAccountName string, expand string, filter string) (result billing.DepartmentListResult, err error)
}

var _ DepartmentsByBillingAccountNameClientAPI = (*billing.DepartmentsByBillingAccountNameClient)(nil)

// DepartmentsClientAPI contains the set of methods on the DepartmentsClient type.
type DepartmentsClientAPI interface {
	Get(ctx context.Context, billingAccountName string, departmentName string, expand string, filter string) (result billing.Department, err error)
}

var _ DepartmentsClientAPI = (*billing.DepartmentsClient)(nil)

// EnrollmentAccountsByBillingAccountNameClientAPI contains the set of methods on the EnrollmentAccountsByBillingAccountNameClient type.
type EnrollmentAccountsByBillingAccountNameClientAPI interface {
	List(ctx context.Context, billingAccountName string, expand string, filter string) (result billing.EnrollmentAccountListResult, err error)
}

var _ EnrollmentAccountsByBillingAccountNameClientAPI = (*billing.EnrollmentAccountsByBillingAccountNameClient)(nil)

// EnrollmentAccountsClientAPI contains the set of methods on the EnrollmentAccountsClient type.
type EnrollmentAccountsClientAPI interface {
	GetByEnrollmentAccountAccountID(ctx context.Context, billingAccountName string, enrollmentAccountName string, expand string, filter string) (result billing.EnrollmentAccount, err error)
}

var _ EnrollmentAccountsClientAPI = (*billing.EnrollmentAccountsClient)(nil)

// InvoicesByBillingAccountClientAPI contains the set of methods on the InvoicesByBillingAccountClient type.
type InvoicesByBillingAccountClientAPI interface {
	List(ctx context.Context, billingAccountName string, periodStartDate string, periodEndDate string) (result billing.InvoiceListResult, err error)
}

var _ InvoicesByBillingAccountClientAPI = (*billing.InvoicesByBillingAccountClient)(nil)

// InvoicePricesheetClientAPI contains the set of methods on the InvoicePricesheetClient type.
type InvoicePricesheetClientAPI interface {
	Download(ctx context.Context, billingAccountName string, invoiceName string) (result billing.InvoicePricesheetDownloadFuture, err error)
}

var _ InvoicePricesheetClientAPI = (*billing.InvoicePricesheetClient)(nil)

// InvoicesByBillingProfileClientAPI contains the set of methods on the InvoicesByBillingProfileClient type.
type InvoicesByBillingProfileClientAPI interface {
	List(ctx context.Context, billingAccountName string, billingProfileName string, periodStartDate string, periodEndDate string) (result billing.InvoiceListResult, err error)
}

var _ InvoicesByBillingProfileClientAPI = (*billing.InvoicesByBillingProfileClient)(nil)

// InvoiceClientAPI contains the set of methods on the InvoiceClient type.
type InvoiceClientAPI interface {
	Get(ctx context.Context, billingAccountName string, billingProfileName string, invoiceName string) (result billing.InvoiceSummary, err error)
}

var _ InvoiceClientAPI = (*billing.InvoiceClient)(nil)

// ProductsByBillingSubscriptionsClientAPI contains the set of methods on the ProductsByBillingSubscriptionsClient type.
type ProductsByBillingSubscriptionsClientAPI interface {
	List(ctx context.Context, billingAccountName string) (result billing.SubscriptionsListResultPage, err error)
}

var _ ProductsByBillingSubscriptionsClientAPI = (*billing.ProductsByBillingSubscriptionsClient)(nil)

// SubscriptionsByBillingProfileClientAPI contains the set of methods on the SubscriptionsByBillingProfileClient type.
type SubscriptionsByBillingProfileClientAPI interface {
	List(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.SubscriptionsListResult, err error)
}

var _ SubscriptionsByBillingProfileClientAPI = (*billing.SubscriptionsByBillingProfileClient)(nil)

// SubscriptionsByInvoiceSectionClientAPI contains the set of methods on the SubscriptionsByInvoiceSectionClient type.
type SubscriptionsByInvoiceSectionClientAPI interface {
	List(ctx context.Context, billingAccountName string, invoiceSectionName string) (result billing.SubscriptionsListResult, err error)
}

var _ SubscriptionsByInvoiceSectionClientAPI = (*billing.SubscriptionsByInvoiceSectionClient)(nil)

// SubscriptionClientAPI contains the set of methods on the SubscriptionClient type.
type SubscriptionClientAPI interface {
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, billingSubscriptionName string) (result billing.SubscriptionSummary, err error)
	Transfer(ctx context.Context, billingAccountName string, invoiceSectionName string, billingSubscriptionName string, parameters billing.TransferBillingSubscriptionRequestProperties) (result billing.SubscriptionTransferFuture, err error)
}

var _ SubscriptionClientAPI = (*billing.SubscriptionClient)(nil)

// ProductsByBillingAccountClientAPI contains the set of methods on the ProductsByBillingAccountClient type.
type ProductsByBillingAccountClientAPI interface {
	List(ctx context.Context, billingAccountName string, filter string) (result billing.ProductsListResultPage, err error)
}

var _ ProductsByBillingAccountClientAPI = (*billing.ProductsByBillingAccountClient)(nil)

// ProductsByInvoiceSectionClientAPI contains the set of methods on the ProductsByInvoiceSectionClient type.
type ProductsByInvoiceSectionClientAPI interface {
	List(ctx context.Context, billingAccountName string, invoiceSectionName string, filter string) (result billing.ProductsListResult, err error)
}

var _ ProductsByInvoiceSectionClientAPI = (*billing.ProductsByInvoiceSectionClient)(nil)

// ProductsClientAPI contains the set of methods on the ProductsClient type.
type ProductsClientAPI interface {
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, productName string) (result billing.ProductSummary, err error)
	Transfer(ctx context.Context, billingAccountName string, invoiceSectionName string, productName string, parameters billing.TransferProductRequestProperties) (result billing.ProductSummary, err error)
}

var _ ProductsClientAPI = (*billing.ProductsClient)(nil)

// TransactionsByBillingAccountClientAPI contains the set of methods on the TransactionsByBillingAccountClient type.
type TransactionsByBillingAccountClientAPI interface {
	List(ctx context.Context, billingAccountName string, startDate string, endDate string, filter string) (result billing.TransactionsListResultPage, err error)
}

var _ TransactionsByBillingAccountClientAPI = (*billing.TransactionsByBillingAccountClient)(nil)

// TransactionsByBillingProfileClientAPI contains the set of methods on the TransactionsByBillingProfileClient type.
type TransactionsByBillingProfileClientAPI interface {
	List(ctx context.Context, billingAccountName string, billingProfileName string, startDate string, endDate string, filter string) (result billing.TransactionsListResult, err error)
}

var _ TransactionsByBillingProfileClientAPI = (*billing.TransactionsByBillingProfileClient)(nil)

// TransactionsByInvoiceSectionClientAPI contains the set of methods on the TransactionsByInvoiceSectionClient type.
type TransactionsByInvoiceSectionClientAPI interface {
	List(ctx context.Context, billingAccountName string, invoiceSectionName string, startDate string, endDate string, filter string) (result billing.TransactionsListResult, err error)
}

var _ TransactionsByInvoiceSectionClientAPI = (*billing.TransactionsByInvoiceSectionClient)(nil)

// PolicyClientAPI contains the set of methods on the PolicyClient type.
type PolicyClientAPI interface {
	GetByBillingProfile(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.Policy, err error)
	Update(ctx context.Context, billingAccountName string, billingProfileName string, parameters billing.Policy) (result billing.Policy, err error)
}

var _ PolicyClientAPI = (*billing.PolicyClient)(nil)

// PropertyClientAPI contains the set of methods on the PropertyClient type.
type PropertyClientAPI interface {
	Get(ctx context.Context) (result billing.Property, err error)
}

var _ PropertyClientAPI = (*billing.PropertyClient)(nil)

// TransfersClientAPI contains the set of methods on the TransfersClient type.
type TransfersClientAPI interface {
	Cancel(ctx context.Context, billingAccountName string, invoiceSectionName string, transferName string) (result billing.TransferDetails, err error)
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, transferName string) (result billing.TransferDetails, err error)
	Initiate(ctx context.Context, billingAccountName string, invoiceSectionName string, body billing.InitiateTransferRequest) (result billing.TransferDetails, err error)
	List(ctx context.Context, billingAccountName string, invoiceSectionName string) (result billing.TransferDetailsListResultPage, err error)
}

var _ TransfersClientAPI = (*billing.TransfersClient)(nil)

// RecipientTransfersClientAPI contains the set of methods on the RecipientTransfersClient type.
type RecipientTransfersClientAPI interface {
	Accept(ctx context.Context, transferName string, body billing.AcceptTransferRequest) (result billing.RecipientTransferDetails, err error)
	Decline(ctx context.Context, transferName string) (result billing.RecipientTransferDetails, err error)
	Get(ctx context.Context, transferName string) (result billing.RecipientTransferDetails, err error)
	List(ctx context.Context) (result billing.RecipientTransferDetailsListResultPage, err error)
}

var _ RecipientTransfersClientAPI = (*billing.RecipientTransfersClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result billing.OperationListResultPage, err error)
}

var _ OperationsClientAPI = (*billing.OperationsClient)(nil)

// AccountBillingPermissionsClientAPI contains the set of methods on the AccountBillingPermissionsClient type.
type AccountBillingPermissionsClientAPI interface {
	List(ctx context.Context, billingAccountName string) (result billing.PermissionsListResult, err error)
}

var _ AccountBillingPermissionsClientAPI = (*billing.AccountBillingPermissionsClient)(nil)

// InvoiceSectionsBillingPermissionsClientAPI contains the set of methods on the InvoiceSectionsBillingPermissionsClient type.
type InvoiceSectionsBillingPermissionsClientAPI interface {
	List(ctx context.Context, billingAccountName string, invoiceSectionName string) (result billing.PermissionsListResult, err error)
}

var _ InvoiceSectionsBillingPermissionsClientAPI = (*billing.InvoiceSectionsBillingPermissionsClient)(nil)

// ProfileBillingPermissionsClientAPI contains the set of methods on the ProfileBillingPermissionsClient type.
type ProfileBillingPermissionsClientAPI interface {
	List(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.PermissionsListResult, err error)
}

var _ ProfileBillingPermissionsClientAPI = (*billing.ProfileBillingPermissionsClient)(nil)

// AccountBillingRoleDefinitionClientAPI contains the set of methods on the AccountBillingRoleDefinitionClient type.
type AccountBillingRoleDefinitionClientAPI interface {
	Get(ctx context.Context, billingAccountName string, billingRoleDefinitionName string) (result billing.RoleDefinition, err error)
	List(ctx context.Context, billingAccountName string) (result billing.RoleDefinitionListResult, err error)
}

var _ AccountBillingRoleDefinitionClientAPI = (*billing.AccountBillingRoleDefinitionClient)(nil)

// InvoiceSectionBillingRoleDefinitionClientAPI contains the set of methods on the InvoiceSectionBillingRoleDefinitionClient type.
type InvoiceSectionBillingRoleDefinitionClientAPI interface {
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, billingRoleDefinitionName string) (result billing.RoleDefinition, err error)
	List(ctx context.Context, billingAccountName string, invoiceSectionName string) (result billing.RoleDefinitionListResult, err error)
}

var _ InvoiceSectionBillingRoleDefinitionClientAPI = (*billing.InvoiceSectionBillingRoleDefinitionClient)(nil)

// ProfileBillingRoleDefinitionClientAPI contains the set of methods on the ProfileBillingRoleDefinitionClient type.
type ProfileBillingRoleDefinitionClientAPI interface {
	Get(ctx context.Context, billingAccountName string, billingProfileName string, billingRoleDefinitionName string) (result billing.RoleDefinition, err error)
	List(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.RoleDefinitionListResult, err error)
}

var _ ProfileBillingRoleDefinitionClientAPI = (*billing.ProfileBillingRoleDefinitionClient)(nil)

// AccountBillingRoleAssignmentClientAPI contains the set of methods on the AccountBillingRoleAssignmentClient type.
type AccountBillingRoleAssignmentClientAPI interface {
	Add(ctx context.Context, billingAccountName string, parameters billing.RoleAssignmentPayload) (result billing.RoleAssignmentListResult, err error)
	Delete(ctx context.Context, billingAccountName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	Get(ctx context.Context, billingAccountName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	List(ctx context.Context, billingAccountName string) (result billing.RoleAssignmentListResult, err error)
}

var _ AccountBillingRoleAssignmentClientAPI = (*billing.AccountBillingRoleAssignmentClient)(nil)

// InvoiceSectionBillingRoleAssignmentClientAPI contains the set of methods on the InvoiceSectionBillingRoleAssignmentClient type.
type InvoiceSectionBillingRoleAssignmentClientAPI interface {
	Add(ctx context.Context, billingAccountName string, invoiceSectionName string, parameters billing.RoleAssignmentPayload) (result billing.RoleAssignmentListResult, err error)
	Delete(ctx context.Context, billingAccountName string, invoiceSectionName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	Get(ctx context.Context, billingAccountName string, invoiceSectionName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	List(ctx context.Context, billingAccountName string, invoiceSectionName string) (result billing.RoleAssignmentListResult, err error)
}

var _ InvoiceSectionBillingRoleAssignmentClientAPI = (*billing.InvoiceSectionBillingRoleAssignmentClient)(nil)

// ProfileBillingRoleAssignmentClientAPI contains the set of methods on the ProfileBillingRoleAssignmentClient type.
type ProfileBillingRoleAssignmentClientAPI interface {
	Add(ctx context.Context, billingAccountName string, billingProfileName string, parameters billing.RoleAssignmentPayload) (result billing.RoleAssignmentListResult, err error)
	Delete(ctx context.Context, billingAccountName string, billingProfileName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	Get(ctx context.Context, billingAccountName string, billingProfileName string, billingRoleAssignmentName string) (result billing.RoleAssignment, err error)
	List(ctx context.Context, billingAccountName string, billingProfileName string) (result billing.RoleAssignmentListResult, err error)
}

var _ ProfileBillingRoleAssignmentClientAPI = (*billing.ProfileBillingRoleAssignmentClient)(nil)
