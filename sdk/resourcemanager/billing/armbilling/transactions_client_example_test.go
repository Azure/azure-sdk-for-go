//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armbilling_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/billing/resource-manager/Microsoft.Billing/stable/2020-05-01/examples/TransactionsListByInvoice.json
func ExampleTransactionsClient_NewListByInvoicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armbilling.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewTransactionsClient().NewListByInvoicePager("{billingAccountName}", "{invoiceName}", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.TransactionListResult = armbilling.TransactionListResult{
		// 	Value: []*armbilling.Transaction{
		// 		{
		// 			Name: to.Ptr("41000000-0000-0000-0000-000000000000"),
		// 			Type: to.Ptr("Microsoft.Billing/billingAccounts/transactions"),
		// 			ID: to.Ptr("/providers/Microsoft.Billing/BillingAccounts/{billingAccountName}/transactions/41000000-0000-0000-0000-000000000000"),
		// 			Properties: &armbilling.TransactionProperties{
		// 				AzureCreditApplied: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](2000),
		// 				},
		// 				AzurePlan: to.Ptr("Microsoft Azure Plan for DevTest"),
		// 				BillingCurrency: to.Ptr("USD"),
		// 				BillingProfileDisplayName: to.Ptr("Contoso operations billing"),
		// 				BillingProfileID: to.Ptr("/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}"),
		// 				Date: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-05-01T00:00:00.000Z"); return t}()),
		// 				Discount: to.Ptr[float32](0.1),
		// 				EffectivePrice: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](10),
		// 				},
		// 				ExchangeRate: to.Ptr[float32](1),
		// 				Invoice: to.Ptr("2344233"),
		// 				InvoiceID: to.Ptr("/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/invoices/2344233"),
		// 				InvoiceSectionDisplayName: to.Ptr("Contoso operations invoiceSection"),
		// 				InvoiceSectionID: to.Ptr("/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/22000000-0000-0000-0000-000000000000"),
		// 				Kind: to.Ptr(armbilling.TransactionTypeKindAll),
		// 				MarketPrice: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](20),
		// 				},
		// 				PricingCurrency: to.Ptr("USD"),
		// 				ProductDescription: to.Ptr("Standard D1, US West 3"),
		// 				ProductFamily: to.Ptr("Storage"),
		// 				ProductType: to.Ptr("VM Instance"),
		// 				ProductTypeID: to.Ptr("A12345"),
		// 				Quantity: to.Ptr[int32](1),
		// 				ServicePeriodEndDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-09-30T00:00:00.000Z"); return t}()),
		// 				ServicePeriodStartDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-05-01T00:00:00.000Z"); return t}()),
		// 				SubTotal: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](4500),
		// 				},
		// 				Tax: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](500),
		// 				},
		// 				TransactionAmount: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](5000),
		// 				},
		// 				TransactionType: to.Ptr(armbilling.ReservationTypePurchase),
		// 				UnitOfMeasure: to.Ptr("1 Minute"),
		// 				UnitType: to.Ptr("1 Runtime Minute"),
		// 				Units: to.Ptr[float32](11.25),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("51000000-0000-0000-0000-000000000000"),
		// 			Type: to.Ptr("Microsoft.Billing/billingAccounts/transactions"),
		// 			ID: to.Ptr("/providers/Microsoft.Billing/BillingAccounts/{billingAccountName}/transactions/51000000-0000-0000-0000-000000000000"),
		// 			Properties: &armbilling.TransactionProperties{
		// 				AzureCreditApplied: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](20),
		// 				},
		// 				AzurePlan: to.Ptr("Microsoft Azure Plan for DevTest"),
		// 				BillingCurrency: to.Ptr("USD"),
		// 				BillingProfileDisplayName: to.Ptr("Contoso operations billing"),
		// 				BillingProfileID: to.Ptr("/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}"),
		// 				Date: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-04-01T00:00:00.000Z"); return t}()),
		// 				Discount: to.Ptr[float32](0.1),
		// 				EffectivePrice: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](10),
		// 				},
		// 				ExchangeRate: to.Ptr[float32](1),
		// 				Invoice: to.Ptr("pending"),
		// 				InvoiceSectionDisplayName: to.Ptr("Contoso operations invoiceSection"),
		// 				InvoiceSectionID: to.Ptr("/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/22000000-0000-0000-0000-000000000000"),
		// 				Kind: to.Ptr(armbilling.TransactionTypeKindAll),
		// 				MarketPrice: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](20),
		// 				},
		// 				PricingCurrency: to.Ptr("USD"),
		// 				ProductDescription: to.Ptr("Standard Support"),
		// 				ProductFamily: to.Ptr("Storage"),
		// 				ProductType: to.Ptr("VM Instance"),
		// 				ProductTypeID: to.Ptr("A12345"),
		// 				Quantity: to.Ptr[int32](1),
		// 				ServicePeriodEndDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-09-30T00:00:00.000Z"); return t}()),
		// 				ServicePeriodStartDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-05-01T00:00:00.000Z"); return t}()),
		// 				SubTotal: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](45),
		// 				},
		// 				Tax: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](5),
		// 				},
		// 				TransactionAmount: &armbilling.Amount{
		// 					Currency: to.Ptr("USD"),
		// 					Value: to.Ptr[float32](50),
		// 				},
		// 				TransactionType: to.Ptr(armbilling.ReservationType("Cancel")),
		// 				UnitOfMeasure: to.Ptr("1 Minute"),
		// 				UnitType: to.Ptr("1 Runtime Minute"),
		// 				Units: to.Ptr[float32](1.25),
		// 			},
		// 	}},
		// }
	}
}
