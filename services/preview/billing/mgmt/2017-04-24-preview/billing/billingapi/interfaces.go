// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package billingapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/billing/mgmt/2017-04-24-preview/billing"
)

// PeriodsClientAPI contains the set of methods on the PeriodsClient type.
type PeriodsClientAPI interface {
	Get(ctx context.Context, billingPeriodName string) (result billing.Period, err error)
	List(ctx context.Context, filter string, skiptoken string, top *int32) (result billing.PeriodsListResultPage, err error)
	ListComplete(ctx context.Context, filter string, skiptoken string, top *int32) (result billing.PeriodsListResultIterator, err error)
}

var _ PeriodsClientAPI = (*billing.PeriodsClient)(nil)

// InvoicesClientAPI contains the set of methods on the InvoicesClient type.
type InvoicesClientAPI interface {
	Get(ctx context.Context, invoiceName string) (result billing.Invoice, err error)
	GetLatest(ctx context.Context) (result billing.Invoice, err error)
	List(ctx context.Context, expand string, filter string, skiptoken string, top *int32) (result billing.InvoicesListResultPage, err error)
	ListComplete(ctx context.Context, expand string, filter string, skiptoken string, top *int32) (result billing.InvoicesListResultIterator, err error)
}

var _ InvoicesClientAPI = (*billing.InvoicesClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result billing.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result billing.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*billing.OperationsClient)(nil)
