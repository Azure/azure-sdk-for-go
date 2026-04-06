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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ApimproductsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	apiId             string
	groupId           string
	productId         string
	serviceName       string
	tagId             string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *ApimproductsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.apiId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "productapiid", 18, false)
	testsuite.groupId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "productgroupid", 20, false)
	testsuite.productId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "productid", 15, false)
	testsuite.serviceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "serviceproduct", 20, false)
	testsuite.tagId, _ = recording.GenerateAlphaNumericID(testsuite.T(), "producttagid", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ApimproductsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestApimproductsTestSuite(t *testing.T) {
	suite.Run(t, new(ApimproductsTestSuite))
}

func (testsuite *ApimproductsTestSuite) Prepare() {
	var err error
	// From step ApiManagementService_CreateOrUpdate
	fmt.Println("Call operation: ApiManagementService_CreateOrUpdate")
	serviceClient, err := armapimanagement.NewServiceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	serviceClientCreateOrUpdateResponsePoller, err := serviceClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, armapimanagement.ServiceResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr(testsuite.location),
		Properties: &armapimanagement.ServiceProperties{
			PublisherEmail: to.Ptr("foo@contoso.com"),
			PublisherName:  to.Ptr("foo"),
		},
		SKU: &armapimanagement.ServiceSKUProperties{
			Name:     to.Ptr(armapimanagement.SKUTypeStandard),
			Capacity: to.Ptr[int32](1),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, serviceClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Product_CreateOrUpdate
	fmt.Println("Call operation: Product_CreateOrUpdate")
	productClient, err := armapimanagement.NewProductClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.ProductContract{
		Properties: &armapimanagement.ProductContractProperties{
			DisplayName: to.Ptr("Test Template ProductName 4"),
		},
	}, &armapimanagement.ProductClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products
func (testsuite *ApimproductsTestSuite) TestProduct() {
	var err error
	// From step Product_GetEntityTag
	fmt.Println("Call operation: Product_GetEntityTag")
	productClient, err := armapimanagement.NewProductClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, nil)
	testsuite.Require().NoError(err)

	// From step Product_ListByService
	fmt.Println("Call operation: Product_ListByService")
	productClientNewListByServicePager := productClient.NewListByServicePager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.ProductClientListByServiceOptions{Filter: nil,
		Top:          nil,
		Skip:         nil,
		ExpandGroups: nil,
		Tags:         nil,
	})
	for productClientNewListByServicePager.More() {
		_, err := productClientNewListByServicePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Product_Get
	fmt.Println("Call operation: Product_Get")
	_, err = productClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, nil)
	testsuite.Require().NoError(err)

	// From step Product_ListByTags
	fmt.Println("Call operation: Product_ListByTags")
	productClientNewListByTagsPager := productClient.NewListByTagsPager(testsuite.resourceGroupName, testsuite.serviceName, &armapimanagement.ProductClientListByTagsOptions{Filter: nil,
		Top:                      nil,
		Skip:                     nil,
		IncludeNotTaggedProducts: nil,
	})
	for productClientNewListByTagsPager.More() {
		_, err := productClientNewListByTagsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Product_Update
	fmt.Println("Call operation: Product_Update")
	_, err = productClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, "*", armapimanagement.ProductUpdateParameters{
		Properties: &armapimanagement.ProductUpdateProperties{
			DisplayName: to.Ptr("Test Template ProductName 4"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products/apis
func (testsuite *ApimproductsTestSuite) TestProductapi() {
	var err error
	// From step Api_CreateOrUpdate
	fmt.Println("Call operation: Api_CreateOrUpdate")
	aPIClient, err := armapimanagement.NewAPIClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	aPIClientCreateOrUpdateResponsePoller, err := aPIClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.apiId, armapimanagement.APICreateOrUpdateParameter{
		Properties: &armapimanagement.APICreateOrUpdateProperties{
			Path:   to.Ptr("petstore"),
			Format: to.Ptr(armapimanagement.ContentFormatOpenapiLink),
			Value:  to.Ptr("https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"),
		},
	}, &armapimanagement.APIClientBeginCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, aPIClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step ProductApi_CreateOrUpdate
	fmt.Println("Call operation: ProductApi_CreateOrUpdate")
	productAPIClient, err := armapimanagement.NewProductAPIClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productAPIClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.apiId, nil)
	testsuite.Require().NoError(err)

	// From step ProductApi_CheckEntityExists
	fmt.Println("Call operation: ProductApi_CheckEntityExists")
	_, err = productAPIClient.CheckEntityExists(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.apiId, nil)
	testsuite.Require().NoError(err)

	// From step ProductApi_ListByProduct
	fmt.Println("Call operation: ProductApi_ListByProduct")
	productAPIClientNewListByProductPager := productAPIClient.NewListByProductPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, &armapimanagement.ProductAPIClientListByProductOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for productAPIClientNewListByProductPager.More() {
		_, err := productAPIClientNewListByProductPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ProductApi_Delete
	fmt.Println("Call operation: ProductApi_Delete")
	_, err = productAPIClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.apiId, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products/groups
func (testsuite *ApimproductsTestSuite) TestProductgroup() {
	var err error
	// From step Group_CreateOrUpdate
	fmt.Println("Call operation: Group_CreateOrUpdate")
	groupClient, err := armapimanagement.NewGroupClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = groupClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.groupId, armapimanagement.GroupCreateParameters{
		Properties: &armapimanagement.GroupCreateParametersProperties{
			DisplayName: to.Ptr(testsuite.groupId),
		},
	}, &armapimanagement.GroupClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ProductGroup_CreateOrUpdate
	fmt.Println("Call operation: ProductGroup_CreateOrUpdate")
	productGroupClient, err := armapimanagement.NewProductGroupClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productGroupClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.groupId, nil)
	testsuite.Require().NoError(err)

	// From step ProductGroup_CheckEntityExists
	fmt.Println("Call operation: ProductGroup_CheckEntityExists")
	_, err = productGroupClient.CheckEntityExists(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.groupId, nil)
	testsuite.Require().NoError(err)

	// From step ProductGroup_ListByProduct
	fmt.Println("Call operation: ProductGroup_ListByProduct")
	productGroupClientNewListByProductPager := productGroupClient.NewListByProductPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, &armapimanagement.ProductGroupClientListByProductOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for productGroupClientNewListByProductPager.More() {
		_, err := productGroupClientNewListByProductPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ProductGroup_Delete
	fmt.Println("Call operation: ProductGroup_Delete")
	_, err = productGroupClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.groupId, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products/policies
func (testsuite *ApimproductsTestSuite) TestProductpolicy() {
	var err error
	// From step ProductPolicy_CreateOrUpdate
	fmt.Println("Call operation: ProductPolicy_CreateOrUpdate")
	productPolicyClient, err := armapimanagement.NewProductPolicyClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productPolicyClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.PolicyIDNamePolicy, armapimanagement.PolicyContract{
		Properties: &armapimanagement.PolicyContractProperties{
			Format: to.Ptr(armapimanagement.PolicyContentFormatXML),
		},
	}, &armapimanagement.ProductPolicyClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ProductPolicy_GetEntityTag
	fmt.Println("Call operation: ProductPolicy_GetEntityTag")
	_, err = productPolicyClient.GetEntityTag(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.PolicyIDNamePolicy, nil)
	testsuite.Require().NoError(err)

	// From step ProductPolicy_Get
	fmt.Println("Call operation: ProductPolicy_Get")
	_, err = productPolicyClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.PolicyIDNamePolicy, &armapimanagement.ProductPolicyClientGetOptions{Format: nil})
	testsuite.Require().NoError(err)

	// From step ProductPolicy_Delete
	fmt.Println("Call operation: ProductPolicy_Delete")
	_, err = productPolicyClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, armapimanagement.PolicyIDNamePolicy, "*", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products/tags
func (testsuite *ApimproductsTestSuite) TestProducttag() {
	var err error
	// From step Tag_CreateOrUpdate
	fmt.Println("Call operation: Tag_CreateOrUpdate")
	tagClient, err := armapimanagement.NewTagClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = tagClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.tagId, armapimanagement.TagCreateUpdateParameters{
		Properties: &armapimanagement.TagContractProperties{
			DisplayName: to.Ptr(testsuite.tagId),
		},
	}, &armapimanagement.TagClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step Tag_AssignToProduct
	fmt.Println("Call operation: Tag_AssignToProduct")
	_, err = tagClient.AssignToProduct(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetEntityStateByProduct
	fmt.Println("Call operation: Tag_GetEntityStateByProduct")
	_, err = tagClient.GetEntityStateByProduct(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_GetByProduct
	fmt.Println("Call operation: Tag_GetByProduct")
	_, err = tagClient.GetByProduct(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)

	// From step Tag_ListByProduct
	fmt.Println("Call operation: Tag_ListByProduct")
	tagClientNewListByProductPager := tagClient.NewListByProductPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, &armapimanagement.TagClientListByProductOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for tagClientNewListByProductPager.More() {
		_, err := tagClientNewListByProductPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Tag_DetachFromProduct
	fmt.Println("Call operation: Tag_DetachFromProduct")
	_, err = tagClient.DetachFromProduct(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, testsuite.tagId, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.ApiManagement/service/products/subscriptions
func (testsuite *ApimproductsTestSuite) TestProductsubscriptions() {
	var err error
	// From step ProductSubscriptions_List
	fmt.Println("Call operation: ProductSubscriptions_List")
	productSubscriptionsClient, err := armapimanagement.NewProductSubscriptionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	productSubscriptionsClientNewListPager := productSubscriptionsClient.NewListPager(testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, &armapimanagement.ProductSubscriptionsClientListOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
	})
	for productSubscriptionsClientNewListPager.More() {
		_, err := productSubscriptionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ApimproductsTestSuite) Cleanup() {
	var err error
	// From step Product_Delete
	fmt.Println("Call operation: Product_Delete")
	productClient, err := armapimanagement.NewProductClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = productClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.serviceName, testsuite.productId, "*", &armapimanagement.ProductClientDeleteOptions{DeleteSubscriptions: to.Ptr(true)})
	testsuite.Require().NoError(err)
}
