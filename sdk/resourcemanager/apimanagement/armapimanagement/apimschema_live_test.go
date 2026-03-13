// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armapimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimschemaTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	schemaId		string
	serviceName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimschemaTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.schemaId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "globalschemaid", 20, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceschema", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimschemaTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimschemaTestSuite(t *testing.T) {
	suite.Run(t, new(ApimschemaTestSuite))
}

func (testsuite *ApimschemaTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name":	to.Ptr("Contoso"),
			"Test":	to.Ptr("User"),
		},
		Location:	to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail:	to.Ptr("foo@contoso.com"),
			PublisherName:	to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:		to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity:	to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/schemas
func (testsuite *ApimschemaTestSuite) TestGlobalschema() {
	var err error
	// From step GlobalSchema_CreateOrUpdate
	fmt.Println("Call operation: GlobalSchema_CreateOrUpdate")
	globalSchemaClient, err := armapimanagement.NewGlobalSchemaClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	globalSchemaClientCreateOrUpdateResponsePoller, err := globalSchemaClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.schemaId, armapimanagement.GlobalSchemaContract{
		Properties: &armapimanagement.GlobalSchemaContractProperties{
			SchemaType:	to.Ptr(armapimanagement.SchemaTypeXML),
			Description:	to.Ptr("sample schema description"),
			Value:		"<xsd:schema xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"\r\n           xmlns:tns=\"http://tempuri.org/PurchaseOrderSchema.xsd\"\r\n           targetNamespace=\"http://tempuri.org/PurchaseOrderSchema.xsd\"\r\n           elementFormDefault=\"qualified\">\r\n <xsd:element name=\"PurchaseOrder\" type=\"tns:PurchaseOrderType\"/>\r\n <xsd:complexType name=\"PurchaseOrderType\">\r\n  <xsd:sequence>\r\n   <xsd:element name=\"ShipTo\" type=\"tns:USAddress\" maxOccurs=\"2\"/>\r\n   <xsd:element name=\"BillTo\" type=\"tns:USAddress\"/>\r\n  </xsd:sequence>\r\n  <xsd:attribute name=\"OrderDate\" type=\"xsd:date\"/>\r\n </xsd:complexType>\r\n\r\n <xsd:complexType name=\"USAddress\">\r\n  <xsd:sequence>\r\n   <xsd:element name=\"name\"   type=\"xsd:string\"/>\r\n   <xsd:element name=\"street\" type=\"xsd:string\"/>\r\n   <xsd:element name=\"city\"   type=\"xsd:string\"/>\r\n   <xsd:element name=\"state\"  type=\"xsd:string\"/>\r\n   <xsd:element name=\"zip\"    type=\"xsd:integer\"/>\r\n  </xsd:sequence>\r\n  <xsd:attribute name=\"country\" type=\"xsd:NMTOKEN\" fixed=\"US\"/>\r\n </xsd:complexType>\r\n</xsd:schema>",
		},
	}, &armapimanagement.GlobalSchemaClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, globalSchemaClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GlobalSchema_GetEntityTag
	fmt.Println("Call operation: GlobalSchema_GetEntityTag")
	_, err = globalSchemaClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.schemaId, nil)
	testsuite.Require().NoError(err)

	// From step GlobalSchema_ListByService
	fmt.Println("Call operation: GlobalSchema_ListByService")
	globalSchemaClientNewListByServicePager := globalSchemaClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.GlobalSchemaClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for globalSchemaClientNewListByServicePager.More() {
		_, err := globalSchemaClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GlobalSchema_Get
	fmt.Println("Call operation: GlobalSchema_Get")
	_, err = globalSchemaClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.schemaId, nil)
	testsuite.Require().NoError(err)

	// From step GlobalSchema_Delete
	fmt.Println("Call operation: GlobalSchema_Delete")
	_, err = globalSchemaClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.schemaId, "*", nil)
	testsuite.Require().NoError(err)
}
