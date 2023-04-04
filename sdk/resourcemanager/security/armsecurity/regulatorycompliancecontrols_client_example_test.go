//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsecurity_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/af3f7994582c0cbd61a48b636907ad2ac95d332c/specification/security/resource-manager/Microsoft.Security/preview/2019-01-01-preview/examples/RegulatoryCompliance/getRegulatoryComplianceControlList_example.json
func ExampleRegulatoryComplianceControlsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurity.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewRegulatoryComplianceControlsClient().NewListPager("PCI-DSS-3.2", &armsecurity.RegulatoryComplianceControlsClientListOptions{Filter: nil})
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
		// page.RegulatoryComplianceControlList = armsecurity.RegulatoryComplianceControlList{
		// 	Value: []*armsecurity.RegulatoryComplianceControl{
		// 		{
		// 			Name: to.Ptr("1.1"),
		// 			Type: to.Ptr("Microsoft.Security/regulatoryComplianceControl"),
		// 			ID: to.Ptr("/subscriptions/20ff7fc3-e762-44dd-bd96-b71116dcdc23/providers/Microsoft.Security/regulatoryComplianceStandards/PCI-DSS-3.2/regulatoryComplianceControls/1.1"),
		// 			Properties: &armsecurity.RegulatoryComplianceControlProperties{
		// 				Description: to.Ptr("Common Criteria Related to Organization and Management."),
		// 				FailedAssessments: to.Ptr[int32](4),
		// 				PassedAssessments: to.Ptr[int32](7),
		// 				SkippedAssessments: to.Ptr[int32](0),
		// 				State: to.Ptr(armsecurity.StateFailed),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("2"),
		// 			Type: to.Ptr("Microsoft.Security/regulatoryComplianceControl"),
		// 			ID: to.Ptr("/subscriptions/20ff7fc3-e762-44dd-bd96-b71116dcdc23/providers/Microsoft.Security/regulatoryComplianceStandards/PCI-DSS-3.2/regulatoryComplianceControls/2"),
		// 			Properties: &armsecurity.RegulatoryComplianceControlProperties{
		// 				Description: to.Ptr("Confidential information is protected during the system design, development, testing, implementation, and change processes in accordance with confidentiality commitments and requirements."),
		// 				FailedAssessments: to.Ptr[int32](0),
		// 				PassedAssessments: to.Ptr[int32](0),
		// 				SkippedAssessments: to.Ptr[int32](10),
		// 				State: to.Ptr(armsecurity.StateSkipped),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("2.1"),
		// 			Type: to.Ptr("Microsoft.Security/regulatoryComplianceControl"),
		// 			ID: to.Ptr("/subscriptions/20ff7fc3-e762-44dd-bd96-b71116dcdc23/providers/Microsoft.Security/regulatoryComplianceStandards/PCI-DSS-3.2/regulatoryComplianceControls/2.1"),
		// 			Properties: &armsecurity.RegulatoryComplianceControlProperties{
		// 				Description: to.Ptr("Changes to confidentiality commitments and requirements are communicated to internal and external users, vendors, and other third parties whose products and services are included in the system."),
		// 				FailedAssessments: to.Ptr[int32](0),
		// 				PassedAssessments: to.Ptr[int32](0),
		// 				SkippedAssessments: to.Ptr[int32](0),
		// 				State: to.Ptr(armsecurity.StateUnsupported),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/af3f7994582c0cbd61a48b636907ad2ac95d332c/specification/security/resource-manager/Microsoft.Security/preview/2019-01-01-preview/examples/RegulatoryCompliance/getRegulatoryComplianceControl_example.json
func ExampleRegulatoryComplianceControlsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurity.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewRegulatoryComplianceControlsClient().Get(ctx, "PCI-DSS-3.2", "1.1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RegulatoryComplianceControl = armsecurity.RegulatoryComplianceControl{
	// 	Name: to.Ptr("1.1"),
	// 	Type: to.Ptr("Microsoft.Security/regulatoryComplianceControl"),
	// 	ID: to.Ptr("/subscriptions/20ff7fc3-e762-44dd-bd96-b71116dcdc23/providers/Microsoft.Security/regulatoryComplianceStandards/PCI-DSS-3.2/regulatoryComplianceControls/1.1"),
	// 	Properties: &armsecurity.RegulatoryComplianceControlProperties{
	// 		Description: to.Ptr("Common Criteria Related to Organization and Management"),
	// 		FailedAssessments: to.Ptr[int32](4),
	// 		PassedAssessments: to.Ptr[int32](7),
	// 		SkippedAssessments: to.Ptr[int32](0),
	// 		State: to.Ptr(armsecurity.StateFailed),
	// 	},
	// }
}
