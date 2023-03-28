//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapplicationinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestListByResourceGroup.json
func ExampleWebTestsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewWebTestsClient().NewListByResourceGroupPager("my-resource-group", nil)
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
		// page.WebTestListResult = armapplicationinsights.WebTestListResult{
		// 	Value: []*armapplicationinsights.WebTest{
		// 		{
		// 			Name: to.Ptr("my-webtest-my-component"),
		// 			Type: to.Ptr("Microsoft.Insights/webtests"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
		// 			Location: to.Ptr("southcentralus"),
		// 			Tags: map[string]*string{
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
		// 			},
		// 			Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 			Properties: &armapplicationinsights.WebTestProperties{
		// 				Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
		// 					WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
		// 				},
		// 				Description: to.Ptr(""),
		// 				Enabled: to.Ptr(false),
		// 				Frequency: to.Ptr[int32](900),
		// 				WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 				Locations: []*armapplicationinsights.WebTestGeolocation{
		// 					{
		// 						Location: to.Ptr("apac-hk-hkn-azr"),
		// 				}},
		// 				WebTestName: to.Ptr("my-webtest"),
		// 				RetryEnabled: to.Ptr(true),
		// 				SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
		// 				Timeout: to.Ptr[int32](120),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("my-webtest-my-other-component"),
		// 			Type: to.Ptr("Microsoft.Insights/webtests"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-other-component"),
		// 			Location: to.Ptr("southcentralus"),
		// 			Tags: map[string]*string{
		// 				"Test": to.Ptr("You can delete this synthetic monitor anytime"),
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-other-component": to.Ptr("Resource"),
		// 			},
		// 			Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 			Properties: &armapplicationinsights.WebTestProperties{
		// 				Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
		// 					WebTest: to.Ptr("<WebTest Name=\"342bccf4-722f-496d-b064-123456789abc\" Id=\"00a15cc1-c903-4f97-9af4-123456789abc\" Enabled=\"False\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"347e1924-9899-4c6e-ad78-123456789abc\" Version=\"1.1\" Url=\"http://my-other-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
		// 				},
		// 				Description: to.Ptr(""),
		// 				Enabled: to.Ptr(false),
		// 				Frequency: to.Ptr[int32](300),
		// 				WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 				Locations: []*armapplicationinsights.WebTestGeolocation{
		// 					{
		// 						Location: to.Ptr("us-fl-mia-edge"),
		// 				}},
		// 				WebTestName: to.Ptr("342bccf4-722f-496d-b064-123456789abc"),
		// 				RetryEnabled: to.Ptr(false),
		// 				SyntheticMonitorID: to.Ptr("my-webtest-my-other-component"),
		// 				Timeout: to.Ptr[int32](90),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestGet.json
func ExampleWebTestsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWebTestsClient().Get(ctx, "my-resource-group", "my-webtest-01-mywebservice", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.WebTest = armapplicationinsights.WebTest{
	// 	Name: to.Ptr("my-webtest-01-mywebservice"),
	// 	Type: to.Ptr("Microsoft.Insights/webtests"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/my-test-resources/providers/Microsoft.Insights/webtests/my-webtest-01-mywebservice"),
	// 	Location: to.Ptr("southcentralus"),
	// 	Tags: map[string]*string{
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-test-resources/providers/Microsoft.Insights/components/mytester": to.Ptr("Resource"),
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-test-resources/providers/Microsoft.Web/sites/mytester": to.Ptr("Resource"),
	// 	},
	// 	Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 	Properties: &armapplicationinsights.WebTestProperties{
	// 		Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
	// 			WebTest: to.Ptr("<WebTest Name=\"mytest-webtest-01\" Id=\"0317d26b-8672-4370-bd6b-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"30\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"a55ce143-4f1e-a7e6-b69e-123456789abc\" Version=\"1.1\" Url=\"http://mytester.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"30\" ParseDependentRequests=\"False\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
	// 		},
	// 		Description: to.Ptr(""),
	// 		Enabled: to.Ptr(false),
	// 		Frequency: to.Ptr[int32](900),
	// 		WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 		Locations: []*armapplicationinsights.WebTestGeolocation{
	// 			{
	// 				Location: to.Ptr("us-fl-mia-edge"),
	// 			},
	// 			{
	// 				Location: to.Ptr("apac-hk-hkn-azr"),
	// 		}},
	// 		WebTestName: to.Ptr("mytest-webtest-01"),
	// 		RetryEnabled: to.Ptr(true),
	// 		SyntheticMonitorID: to.Ptr("my-webtest-01-mywebservice"),
	// 		Timeout: to.Ptr[int32](30),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestCreate.json
func ExampleWebTestsClient_CreateOrUpdate_webTestCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWebTestsClient().CreateOrUpdate(ctx, "my-resource-group", "my-webtest-my-component", armapplicationinsights.WebTest{
		Location: to.Ptr("South Central US"),
		Kind:     to.Ptr(armapplicationinsights.WebTestKindPing),
		Properties: &armapplicationinsights.WebTestProperties{
			Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
				WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\" ><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
			},
			Description: to.Ptr("Ping web test alert for mytestwebapp"),
			Enabled:     to.Ptr(true),
			Frequency:   to.Ptr[int32](900),
			WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
			Locations: []*armapplicationinsights.WebTestGeolocation{
				{
					Location: to.Ptr("us-fl-mia-edge"),
				}},
			WebTestName:        to.Ptr("my-webtest-my-component"),
			RetryEnabled:       to.Ptr(true),
			SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
			Timeout:            to.Ptr[int32](120),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.WebTest = armapplicationinsights.WebTest{
	// 	Name: to.Ptr("my-webtest-my-component"),
	// 	Type: to.Ptr("Microsoft.Insights/webtests"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
	// 	Location: to.Ptr("southcentralus"),
	// 	Tags: map[string]*string{
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
	// 	},
	// 	Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 	Properties: &armapplicationinsights.WebTestProperties{
	// 		Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
	// 			WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\" ><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
	// 		},
	// 		Description: to.Ptr("Ping web test alert for mytestwebapp"),
	// 		Enabled: to.Ptr(true),
	// 		Frequency: to.Ptr[int32](900),
	// 		WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 		Locations: []*armapplicationinsights.WebTestGeolocation{
	// 			{
	// 				Location: to.Ptr("us-fl-mia-edge"),
	// 		}},
	// 		WebTestName: to.Ptr("my-webtest-my-component"),
	// 		RetryEnabled: to.Ptr(true),
	// 		SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
	// 		Timeout: to.Ptr[int32](120),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestUpdate.json
func ExampleWebTestsClient_CreateOrUpdate_webTestUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWebTestsClient().CreateOrUpdate(ctx, "my-resource-group", "my-webtest-my-component", armapplicationinsights.WebTest{
		Location: to.Ptr("South Central US"),
		Kind:     to.Ptr(armapplicationinsights.WebTestKindPing),
		Properties: &armapplicationinsights.WebTestProperties{
			Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
				WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"30\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\" ><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"30\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
			},
			Frequency:   to.Ptr[int32](600),
			WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
			Locations: []*armapplicationinsights.WebTestGeolocation{
				{
					Location: to.Ptr("us-fl-mia-edge"),
				},
				{
					Location: to.Ptr("apac-hk-hkn-azr"),
				}},
			WebTestName:        to.Ptr("my-webtest-my-component"),
			SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
			Timeout:            to.Ptr[int32](30),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.WebTest = armapplicationinsights.WebTest{
	// 	Name: to.Ptr("my-webtest-my-component"),
	// 	Type: to.Ptr("Microsoft.Insights/webtests"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
	// 	Location: to.Ptr("southcentralus"),
	// 	Tags: map[string]*string{
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
	// 	},
	// 	Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 	Properties: &armapplicationinsights.WebTestProperties{
	// 		Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
	// 			WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"30\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\" ><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"30\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
	// 		},
	// 		Description: to.Ptr("Ping web test alert for mytestwebapp"),
	// 		Enabled: to.Ptr(true),
	// 		Frequency: to.Ptr[int32](600),
	// 		WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 		Locations: []*armapplicationinsights.WebTestGeolocation{
	// 			{
	// 				Location: to.Ptr("us-fl-mia-edge"),
	// 			},
	// 			{
	// 				Location: to.Ptr("apac-hk-hkn-azr"),
	// 		}},
	// 		WebTestName: to.Ptr("my-webtest-my-component"),
	// 		RetryEnabled: to.Ptr(true),
	// 		SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
	// 		Timeout: to.Ptr[int32](30),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestUpdateTagsOnly.json
func ExampleWebTestsClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWebTestsClient().UpdateTags(ctx, "my-resource-group", "my-webtest-my-component", armapplicationinsights.TagsResource{
		Tags: map[string]*string{
			"Color":          to.Ptr("AzureBlue"),
			"CustomField-01": to.Ptr("This is a random value"),
			"SystemType":     to.Ptr("A08"),
			"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
			"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp":           to.Ptr("Resource"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.WebTest = armapplicationinsights.WebTest{
	// 	Name: to.Ptr("my-webtest-my-component"),
	// 	Type: to.Ptr("Microsoft.Insights/webtests"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
	// 	Location: to.Ptr("southcentralus"),
	// 	Tags: map[string]*string{
	// 		"Color": to.Ptr("AzureBlue"),
	// 		"CustomField-01": to.Ptr("This is a random value"),
	// 		"SystemType": to.Ptr("A08"),
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
	// 		"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
	// 	},
	// 	Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 	Properties: &armapplicationinsights.WebTestProperties{
	// 		Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
	// 			WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"30\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\" ><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"30\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
	// 		},
	// 		Description: to.Ptr("Ping web test alert for mytestwebapp"),
	// 		Enabled: to.Ptr(true),
	// 		Frequency: to.Ptr[int32](600),
	// 		WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
	// 		Locations: []*armapplicationinsights.WebTestGeolocation{
	// 			{
	// 				Location: to.Ptr("us-fl-mia-edge"),
	// 			},
	// 			{
	// 				Location: to.Ptr("apac-hk-hkn-azr"),
	// 		}},
	// 		WebTestName: to.Ptr("my-webtest-my-component"),
	// 		RetryEnabled: to.Ptr(true),
	// 		SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
	// 		Timeout: to.Ptr[int32](30),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestDelete.json
func ExampleWebTestsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewWebTestsClient().Delete(ctx, "my-resource-group", "my-webtest-01-mywebservice", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestList.json
func ExampleWebTestsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewWebTestsClient().NewListPager(nil)
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
		// page.WebTestListResult = armapplicationinsights.WebTestListResult{
		// 	Value: []*armapplicationinsights.WebTest{
		// 		{
		// 			Name: to.Ptr("my-webtest-my-component"),
		// 			Type: to.Ptr("Microsoft.Insights/webtests"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
		// 			Location: to.Ptr("southcentralus"),
		// 			Tags: map[string]*string{
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
		// 			},
		// 			Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 			Properties: &armapplicationinsights.WebTestProperties{
		// 				Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
		// 					WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
		// 				},
		// 				Description: to.Ptr(""),
		// 				Enabled: to.Ptr(false),
		// 				Frequency: to.Ptr[int32](900),
		// 				WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 				Locations: []*armapplicationinsights.WebTestGeolocation{
		// 				},
		// 				WebTestName: to.Ptr("my-webtest"),
		// 				RetryEnabled: to.Ptr(true),
		// 				SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
		// 				Timeout: to.Ptr[int32](120),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("my-webtest-my-other-component"),
		// 			Type: to.Ptr("Microsoft.Insights/webtests"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/my-other-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-other-component"),
		// 			Location: to.Ptr("southcentralus"),
		// 			Tags: map[string]*string{
		// 				"Test": to.Ptr("You can delete this synthetic monitor anytime"),
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-other-resource-group/providers/Microsoft.Insights/components/my-other-component": to.Ptr("Resource"),
		// 			},
		// 			Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 			Properties: &armapplicationinsights.WebTestProperties{
		// 				Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
		// 					WebTest: to.Ptr("<WebTest Name=\"342bccf4-722f-496d-b064-123456789abc\" Id=\"00a15cc1-c903-4f97-9af4-123456789abc\" Enabled=\"False\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"347e1924-9899-4c6e-ad78-123456789abc\" Version=\"1.1\" Url=\"http://my-other-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
		// 				},
		// 				Description: to.Ptr(""),
		// 				Enabled: to.Ptr(false),
		// 				Frequency: to.Ptr[int32](900),
		// 				WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 				Locations: []*armapplicationinsights.WebTestGeolocation{
		// 				},
		// 				WebTestName: to.Ptr("342bccf4-722f-496d-b064-123456789abc"),
		// 				RetryEnabled: to.Ptr(false),
		// 				SyntheticMonitorID: to.Ptr("my-webtest-my-other-component"),
		// 				Timeout: to.Ptr[int32](120),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7932c2df6c8435d6c0e5cbebbca79bce627d5f06/specification/applicationinsights/resource-manager/Microsoft.Insights/stable/2015-05-01/examples/WebTestListByComponent.json
func ExampleWebTestsClient_NewListByComponentPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapplicationinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewWebTestsClient().NewListByComponentPager("my-component", "my-resource-group", nil)
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
		// page.WebTestListResult = armapplicationinsights.WebTestListResult{
		// 	Value: []*armapplicationinsights.WebTest{
		// 		{
		// 			Name: to.Ptr("my-webtest-my-component"),
		// 			Type: to.Ptr("Microsoft.Insights/webtests"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/webtests/my-webtest-my-component"),
		// 			Location: to.Ptr("southcentralus"),
		// 			Tags: map[string]*string{
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Insights/components/my-component": to.Ptr("Resource"),
		// 				"hidden-link:/subscriptions/subid/resourceGroups/my-resource-group/providers/Microsoft.Web/sites/mytestwebapp": to.Ptr("Resource"),
		// 			},
		// 			Kind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 			Properties: &armapplicationinsights.WebTestProperties{
		// 				Configuration: &armapplicationinsights.WebTestPropertiesConfiguration{
		// 					WebTest: to.Ptr("<WebTest Name=\"my-webtest\" Id=\"678ddf96-1ab8-44c8-9274-123456789abc\" Enabled=\"True\" CssProjectStructure=\"\" CssIteration=\"\" Timeout=\"120\" WorkItemIds=\"\" xmlns=\"http://microsoft.com/schemas/VisualStudio/TeamTest/2010\" Description=\"\" CredentialUserName=\"\" CredentialPassword=\"\" PreAuthenticate=\"True\" Proxy=\"default\" StopOnError=\"False\" RecordedResultFile=\"\" ResultsLocale=\"\"><Items><Request Method=\"GET\" Guid=\"a4162485-9114-fcfc-e086-123456789abc\" Version=\"1.1\" Url=\"http://my-component.azurewebsites.net\" ThinkTime=\"0\" Timeout=\"120\" ParseDependentRequests=\"True\" FollowRedirects=\"True\" RecordResult=\"True\" Cache=\"False\" ResponseTimeGoal=\"0\" Encoding=\"utf-8\" ExpectedHttpStatusCode=\"200\" ExpectedResponseUrl=\"\" ReportingName=\"\" IgnoreHttpStatusCode=\"False\" /></Items></WebTest>"),
		// 				},
		// 				Description: to.Ptr(""),
		// 				Enabled: to.Ptr(false),
		// 				Frequency: to.Ptr[int32](900),
		// 				WebTestKind: to.Ptr(armapplicationinsights.WebTestKindPing),
		// 				Locations: []*armapplicationinsights.WebTestGeolocation{
		// 					{
		// 						Location: to.Ptr("apac-hk-hkn-azr"),
		// 				}},
		// 				WebTestName: to.Ptr("my-webtest"),
		// 				RetryEnabled: to.Ptr(true),
		// 				SyntheticMonitorID: to.Ptr("my-webtest-my-component"),
		// 				Timeout: to.Ptr[int32](120),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 			},
		// 	}},
		// }
	}
}
