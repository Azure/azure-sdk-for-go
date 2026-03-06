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

type ApimapisTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	operationId		string
	apiId			string
	attachmentId		string
	commentId		string
	diagnosticId		string
	issueId			string
	releaseId		string
	schemaId		string
	serviceName		string
	tagDescriptionId	string
	tagId			string
	tagOperationId		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *ApimapisTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.operationId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "operationi", 16, false)
	testsuite.apiId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "apiid", 10, false)
	testsuite.attachmentId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "attachment", 16, false)
	testsuite.commentId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "commentid", 15, false)
	testsuite.diagnosticId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "diagnostic", 16, false)
	testsuite.issueId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "issueid", 13, false)
	testsuite.releaseId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "releaseid", 15, false)
	testsuite.schemaId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "schemaid", 14, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "servicenam", 16, false)
	testsuite.tagDescriptionId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tagdescrip", 16, false)
	testsuite.tagId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tagid", 11, false)
	testsuite.tagOperationId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "tagoperation", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimapisTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimapisTestSuite(t *testing.T) {
	suite.Run(t, new(ApimapisTestSuite))
}

func (testsuite *ApimapisTestSuite) Prepare() {
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

	// From step Api_CreateOrUpdate
	fmt.Println("Call operation: Api_CreateOrUpdate")
	aPIClient, err := armapimanagement.NewAPIClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aPIClientCreateOrUpdateResponsePoller, err := aPIClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.APICreateOrUpdateParameter{
		Properties: &armapimanagement.APICreateOrUpdateProperties{
			Path:	to.Ptr("petstore"),
			Format:	to.Ptr(armapimanagement.ContentFormatOpenapiLink),
			Value:	to.Ptr("https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"),
		},
	}, &armapimanagement.APIClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aPIClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApiOperation_CreateOrUpdate
	fmt.Println("Call operation: ApiOperation_CreateOrUpdate")
	aPIOperationClient, err := armapimanagement.NewAPIOperationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIOperationClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, armapimanagement.OperationContract{
		Properties: &armapimanagement.OperationContractProperties{
			TemplateParameters: []*armapimanagement.ParameterContract{
				{
					Name:		to.Ptr("uid"),
					Type:		to.Ptr("string"),
					Description:	to.Ptr("user id"),
				}},
			Method:		to.Ptr("GET"),
			DisplayName:	to.Ptr("example operation"),
			URLTemplate:	to.Ptr("/operation/customers/{uid}"),
		},
	}, &armapimanagement.APIOperationClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis
func (testsuite *ApimapisTestSuite) TestApi() {
	var err error
	// From step Api_GetEntityTag
	fmt.Println("Call operation: Api_GetEntityTag")
	aPIClient, err := armapimanagement.NewAPIClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, nil)
	testsuite.Require().NoError(err)

	// From step Api_ListByService
	fmt.Println("Call operation: Api_ListByService")
	aPIClientNewListByServicePager := aPIClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.APIClientListByServiceOptions{Filter: nil,
		Top:			nil,
		Skip:			nil,
		Tags:			nil,
		ExpandAPIVersionSet:	nil,
	})
	for aPIClientNewListByServicePager.More() {
		_, err := aPIClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Api_ListByTags
	fmt.Println("Call operation: Api_ListByTags")
	aPIClientNewListByTagsPager := aPIClient.NewListByTagsPager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.APIClientListByTagsOptions{Filter: nil,
		Top:			nil,
		Skip:			nil,
		IncludeNotTaggedApis:	nil,
	})
	for aPIClientNewListByTagsPager.More() {
		_, err := aPIClientNewListByTagsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Api_Get
	fmt.Println("Call operation: Api_Get")
	_, err = aPIClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, nil)
	testsuite.Require().NoError(err)

	// From step Api_Update
	fmt.Println("Call operation: Api_Update")
	_, err = aPIClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, "*", armapimanagement.APIUpdateContract{
		Properties: &armapimanagement.APIContractUpdateProperties{
			Path:		to.Ptr("newecho"),
			DisplayName:	to.Ptr("Echo API New"),
			ServiceURL:	to.Ptr("http://echoapi.cloudapp.net/api2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/tags
func (testsuite *ApimapisTestSuite) TestTag() {
	var err error
	// From step Tag_CreateOrUpdate
	fmt.Println("Call operation: Tag_CreateOrUpdate")
	tagClient, err := armapimanagement.NewTagClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tagClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, armapimanagement.TagCreateUpdateParameters{
		Properties: &armapimanagement.TagContractProperties{
			DisplayName: to.Ptr("tag2"),
		},
	}, &armapimanagement.TagClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Tag_AssignToApi
	fmt.Println("Call operation: Tag_AssignToApi")
	_, err = tagClient.AssignToAPI(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetEntityStateByApi
	fmt.Println("Call operation: Tag_GetEntityStateByApi")
	_, err = tagClient.GetEntityStateByAPI(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetByApi
	fmt.Println("Call operation: Tag_GetByApi")
	_, err = tagClient.GetByAPI(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_Delete
	fmt.Println("Call operation: Tag_Delete")
	_, err = tagClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/releases
func (testsuite *ApimapisTestSuite) TestApirelease() {
	var err error
	// From step ApiRelease_CreateOrUpdate
	fmt.Println("Call operation: ApiRelease_CreateOrUpdate")
	aPIReleaseClient, err := armapimanagement.NewAPIReleaseClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIReleaseClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.releaseId, armapimanagement.APIReleaseContract{
		Properties: &armapimanagement.APIReleaseContractProperties{
			APIID:	to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.ApiManagement/service/" + testsuite.serviceName + "/apis/" + testsuite.apiId),
			Notes:	to.Ptr("yahooagain"),
		},
	}, &armapimanagement.APIReleaseClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ApiRelease_GetEntityTag
	fmt.Println("Call operation: ApiRelease_GetEntityTag")
	_, err = aPIReleaseClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.releaseId, nil)
	testsuite.Require().NoError(err)

	// From step ApiRelease_ListByService
	fmt.Println("Call operation: ApiRelease_ListByService")
	aPIReleaseClientNewListByServicePager := aPIReleaseClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, &armapimanagement.APIReleaseClientListByServiceOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for aPIReleaseClientNewListByServicePager.More() {
		_, err := aPIReleaseClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiRelease_Get
	fmt.Println("Call operation: ApiRelease_Get")
	_, err = aPIReleaseClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.releaseId, nil)
	testsuite.Require().NoError(err)

	// From step ApiRelease_Update
	fmt.Println("Call operation: ApiRelease_Update")
	_, err = aPIReleaseClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.releaseId, "*", armapimanagement.APIReleaseContract{
		Properties: &armapimanagement.APIReleaseContractProperties{
			APIID:	to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.ApiManagement/service/" + testsuite.serviceName + "/apis/" + testsuite.apiId),
			Notes:	to.Ptr("yahooagain"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ApiRelease_Delete
	fmt.Println("Call operation: ApiRelease_Delete")
	_, err = aPIReleaseClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.releaseId, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/policies
func (testsuite *ApimapisTestSuite) TestApipolicy() {
	var err error
	// From step ApiPolicy_CreateOrUpdate
	fmt.Println("Call operation: ApiPolicy_CreateOrUpdate")
	aPIPolicyClient, err := armapimanagement.NewAPIPolicyClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIPolicyClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.PolicyIDNamePolicy, armapimanagement.PolicyContract{
		Properties: &armapimanagement.PolicyContractProperties{
			Value:	to.Ptr("<policies> <inbound /> <backend>    <forward-request />  </backend>  <outbound /></policies>"),
			Format:	to.Ptr(armapimanagement.PolicyContentFormatXML),
		},
	}, &armapimanagement.APIPolicyClientCreateOrUpdateOptions{IfMatch: to.Ptr("*")})
	testsuite.Require().NoError(err)

	// From step ApiPolicy_GetEntityTag
	fmt.Println("Call operation: ApiPolicy_GetEntityTag")
	_, err = aPIPolicyClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.PolicyIDNamePolicy, nil)
	testsuite.Require().NoError(err)

	// From step ApiPolicy_ListByApi
	fmt.Println("Call operation: ApiPolicy_ListByApi")
	_, err = aPIPolicyClient.ListByAPI(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, nil)
	testsuite.Require().NoError(err)

	// From step ApiPolicy_Get
	fmt.Println("Call operation: ApiPolicy_Get")
	_, err = aPIPolicyClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.PolicyIDNamePolicy, &armapimanagement.APIPolicyClientGetOptions{Format: nil})
	testsuite.Require().NoError(err)

	// From step ApiPolicy_Delete
	fmt.Println("Call operation: ApiPolicy_Delete")
	_, err = aPIPolicyClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.PolicyIDNamePolicy, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/schemas
func (testsuite *ApimapisTestSuite) TestApischema() {
	var err error
	// From step ApiSchema_CreateOrUpdate
	fmt.Println("Call operation: ApiSchema_CreateOrUpdate")
	aPISchemaClient, err := armapimanagement.NewAPISchemaClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aPISchemaClientCreateOrUpdateResponsePoller, err := aPISchemaClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.schemaId, armapimanagement.SchemaContract{
		Properties: &armapimanagement.SchemaContractProperties{
			ContentType:	to.Ptr("application/vnd.ms-azure-apim.xsd+xml"),
			Document: &armapimanagement.SchemaDocumentProperties{
				Value: to.Ptr("<s:schema elementFormDefault=\"qualified\" targetNamespace=\"http://ws.cdyne.com/WeatherWS/\" xmlns:tns=\"http://ws.cdyne.com/WeatherWS/\" xmlns:s=\"http://www.w3.org/2001/XMLSchema\" xmlns:soap12=\"http://schemas.xmlsoap.org/wsdl/soap12/\" xmlns:mime=\"http://schemas.xmlsoap.org/wsdl/mime/\" xmlns:soap=\"http://schemas.xmlsoap.org/wsdl/soap/\" xmlns:tm=\"http://microsoft.com/wsdl/mime/textMatching/\" xmlns:http=\"http://schemas.xmlsoap.org/wsdl/http/\" xmlns:soapenc=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:wsdl=\"http://schemas.xmlsoap.org/wsdl/\" xmlns:apim-wsdltns=\"http://ws.cdyne.com/WeatherWS/\">\r\n  <s:element name=\"GetWeatherInformation\">\r\n    <s:complexType />\r\n  </s:element>\r\n  <s:element name=\"GetWeatherInformationResponse\">\r\n    <s:complexType>\r\n      <s:sequence>\r\n        <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"GetWeatherInformationResult\" type=\"tns:ArrayOfWeatherDescription\" />\r\n      </s:sequence>\r\n    </s:complexType>\r\n  </s:element>\r\n  <s:complexType name=\"ArrayOfWeatherDescription\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"0\" maxOccurs=\"unbounded\" name=\"WeatherDescription\" type=\"tns:WeatherDescription\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:complexType name=\"WeatherDescription\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Description\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"PictureURL\" type=\"s:string\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:element name=\"GetCityForecastByZIP\">\r\n    <s:complexType>\r\n      <s:sequence>\r\n        <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ZIP\" type=\"s:string\" />\r\n      </s:sequence>\r\n    </s:complexType>\r\n  </s:element>\r\n  <s:element name=\"GetCityForecastByZIPResponse\">\r\n    <s:complexType>\r\n      <s:sequence>\r\n        <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"GetCityForecastByZIPResult\" type=\"tns:ForecastReturn\" />\r\n      </s:sequence>\r\n    </s:complexType>\r\n  </s:element>\r\n  <s:complexType name=\"ForecastReturn\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Success\" type=\"s:boolean\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ResponseText\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"State\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"City\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WeatherStationCity\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ForecastResult\" type=\"tns:ArrayOfForecast\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:complexType name=\"ArrayOfForecast\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"0\" maxOccurs=\"unbounded\" name=\"Forecast\" nillable=\"true\" type=\"tns:Forecast\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:complexType name=\"Forecast\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Date\" type=\"s:dateTime\" />\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Desciption\" type=\"s:string\" />\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Temperatures\" type=\"tns:temp\" />\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"ProbabilityOfPrecipiation\" type=\"tns:POP\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:complexType name=\"temp\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"MorningLow\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"DaytimeHigh\" type=\"s:string\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:complexType name=\"POP\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Nighttime\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Daytime\" type=\"s:string\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:element name=\"GetCityWeatherByZIP\">\r\n    <s:complexType>\r\n      <s:sequence>\r\n        <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ZIP\" type=\"s:string\" />\r\n      </s:sequence>\r\n    </s:complexType>\r\n  </s:element>\r\n  <s:element name=\"GetCityWeatherByZIPResponse\">\r\n    <s:complexType>\r\n      <s:sequence>\r\n        <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"GetCityWeatherByZIPResult\" type=\"tns:WeatherReturn\" />\r\n      </s:sequence>\r\n    </s:complexType>\r\n  </s:element>\r\n  <s:complexType name=\"WeatherReturn\">\r\n    <s:sequence>\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Success\" type=\"s:boolean\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ResponseText\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"State\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"City\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WeatherStationCity\" type=\"s:string\" />\r\n      <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Description\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Temperature\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"RelativeHumidity\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Wind\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Pressure\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Visibility\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WindChill\" type=\"s:string\" />\r\n      <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Remarks\" type=\"s:string\" />\r\n    </s:sequence>\r\n  </s:complexType>\r\n  <s:element name=\"ArrayOfWeatherDescription\" nillable=\"true\" type=\"tns:ArrayOfWeatherDescription\" />\r\n  <s:element name=\"ForecastReturn\" nillable=\"true\" type=\"tns:ForecastReturn\" />\r\n  <s:element name=\"WeatherReturn\" type=\"tns:WeatherReturn\" />\r\n</s:schema>"),
			},
		},
	}, &armapimanagement.APISchemaClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aPISchemaClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ApiSchema_GetEntityTag
	fmt.Println("Call operation: ApiSchema_GetEntityTag")
	_, err = aPISchemaClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.schemaId, nil)
	testsuite.Require().NoError(err)

	// From step ApiSchema_ListByApi
	fmt.Println("Call operation: ApiSchema_ListByApi")
	aPISchemaClientNewListByAPIPager := aPISchemaClient.NewListByAPIPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, &armapimanagement.APISchemaClientListByAPIOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
	})
	for aPISchemaClientNewListByAPIPager.More() {
		_, err := aPISchemaClientNewListByAPIPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiSchema_Get
	fmt.Println("Call operation: ApiSchema_Get")
	_, err = aPISchemaClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.schemaId, nil)
	testsuite.Require().NoError(err)

	// From step ApiSchema_Delete
	fmt.Println("Call operation: ApiSchema_Delete")
	_, err = aPISchemaClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.schemaId, "*", &armapimanagement.APISchemaClientDeleteOptions{Force: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/operations
func (testsuite *ApimapisTestSuite) TestApioperation() {
	var err error
	// From step ApiOperation_GetEntityTag
	fmt.Println("Call operation: ApiOperation_GetEntityTag")
	aPIOperationClient, err := armapimanagement.NewAPIOperationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIOperationClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, nil)
	testsuite.Require().NoError(err)

	// From step ApiOperation_ListByApi
	fmt.Println("Call operation: ApiOperation_ListByApi")
	aPIOperationClientNewListByAPIPager := aPIOperationClient.NewListByAPIPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, &armapimanagement.APIOperationClientListByAPIOptions{Filter: nil,
		Top:	nil,
		Skip:	nil,
		Tags:	nil,
	})
	for aPIOperationClientNewListByAPIPager.More() {
		_, err := aPIOperationClientNewListByAPIPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ApiOperation_Get
	fmt.Println("Call operation: ApiOperation_Get")
	_, err = aPIOperationClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, nil)
	testsuite.Require().NoError(err)

	// From step ApiOperation_Update
	fmt.Println("Call operation: ApiOperation_Update")
	_, err = aPIOperationClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, "*", armapimanagement.OperationUpdateContract{
		Properties: &armapimanagement.OperationUpdateContractProperties{
			Description: to.Ptr("update description"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/operations/tags
func (testsuite *ApimapisTestSuite) TestOperationstag() {
	var err error
	// From step Tag_CreateOrUpdate
	fmt.Println("Call operation: Tag_CreateOrUpdate")
	tagClient, err := armapimanagement.NewTagClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tagClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagOperationId, armapimanagement.TagCreateUpdateParameters{
		Properties: &armapimanagement.TagContractProperties{
			DisplayName: to.Ptr("tagoperation1"),
		},
	}, &armapimanagement.TagClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Tag_AssignToOperation
	fmt.Println("Call operation: Tag_AssignToOperation")
	_, err = tagClient.AssignToOperation(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, testsuite.tagOperationId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetEntityStateByOperation
	fmt.Println("Call operation: Tag_GetEntityStateByOperation")
	_, err = tagClient.GetEntityStateByOperation(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, testsuite.tagOperationId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetByOperation
	fmt.Println("Call operation: Tag_GetByOperation")
	_, err = tagClient.GetByOperation(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, testsuite.tagOperationId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_Delete
	fmt.Println("Call operation: Tag_Delete")
	_, err = tagClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagOperationId, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/apis/operations/policies
func (testsuite *ApimapisTestSuite) TestApioperationpolicy() {
	var err error
	// From step ApiOperationPolicy_CreateOrUpdate
	fmt.Println("Call operation: ApiOperationPolicy_CreateOrUpdate")
	aPIOperationPolicyClient, err := armapimanagement.NewAPIOperationPolicyClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIOperationPolicyClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, armapimanagement.PolicyIDNamePolicy, armapimanagement.PolicyContract{
		Properties: &armapimanagement.PolicyContractProperties{
			Value:	to.Ptr("<policies> <inbound /> <backend>    <forward-request />  </backend>  <outbound /></policies>"),
			Format:	to.Ptr(armapimanagement.PolicyContentFormatXML),
		},
	}, &armapimanagement.APIOperationPolicyClientCreateOrUpdateOptions{IfMatch: to.Ptr("*")})
	testsuite.Require().NoError(err)

	// From step ApiOperationPolicy_GetEntityTag
	fmt.Println("Call operation: ApiOperationPolicy_GetEntityTag")
	_, err = aPIOperationPolicyClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, armapimanagement.PolicyIDNamePolicy, nil)
	testsuite.Require().NoError(err)

	// From step ApiOperationPolicy_ListByOperation
	fmt.Println("Call operation: ApiOperationPolicy_ListByOperation")
	_, err = aPIOperationPolicyClient.ListByOperation(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, nil)
	testsuite.Require().NoError(err)

	// From step ApiOperationPolicy_Get
	fmt.Println("Call operation: ApiOperationPolicy_Get")
	_, err = aPIOperationPolicyClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, armapimanagement.PolicyIDNamePolicy, &armapimanagement.APIOperationPolicyClientGetOptions{Format: nil})
	testsuite.Require().NoError(err)

	// From step ApiOperationPolicy_Delete
	fmt.Println("Call operation: ApiOperationPolicy_Delete")
	_, err = aPIOperationPolicyClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, armapimanagement.PolicyIDNamePolicy, "*", nil)
	testsuite.Require().NoError(err)
}

func (testsuite *ApimapisTestSuite) Cleanup() {
	var err error
	// From step ApiOperation_Delete
	fmt.Println("Call operation: ApiOperation_Delete")
	aPIOperationClient, err := armapimanagement.NewAPIOperationClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = aPIOperationClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, testsuite.operationId, "*", nil)
	testsuite.Require().NoError(err)
}
