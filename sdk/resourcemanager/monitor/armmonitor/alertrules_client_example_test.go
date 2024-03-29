//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmonitor_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/createOrUpdateAlertRule.json
func ExampleAlertRulesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAlertRulesClient().CreateOrUpdate(ctx, "Rac46PostSwapRG", "chiricutin", armmonitor.AlertRuleResource{
		Location: to.Ptr("West US"),
		Tags:     map[string]*string{},
		Properties: &armmonitor.AlertRule{
			Name:        to.Ptr("chiricutin"),
			Description: to.Ptr("Pura Vida"),
			Actions:     []armmonitor.RuleActionClassification{},
			Condition: &armmonitor.ThresholdRuleCondition{
				DataSource: &armmonitor.RuleMetricDataSource{
					ODataType:   to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
					ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
					MetricName:  to.Ptr("Requests"),
				},
				ODataType:       to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
				Operator:        to.Ptr(armmonitor.ConditionOperatorGreaterThan),
				Threshold:       to.Ptr[float64](3),
				TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
				WindowSize:      to.Ptr("PT5M"),
			},
			IsEnabled: to.Ptr(true),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AlertRuleResource = armmonitor.AlertRuleResource{
	// 	Name: to.Ptr("chiricutin"),
	// 	Type: to.Ptr("Microsoft.Insights/alertRules"),
	// 	ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/chiricutin"),
	// 	Location: to.Ptr("West US"),
	// 	Tags: map[string]*string{
	// 		"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
	// 		"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
	// 	},
	// 	Properties: &armmonitor.AlertRule{
	// 		Name: to.Ptr("chiricutin"),
	// 		Description: to.Ptr("Pura Vida"),
	// 		Actions: []armmonitor.RuleActionClassification{
	// 		},
	// 		Condition: &armmonitor.ThresholdRuleCondition{
	// 			DataSource: &armmonitor.RuleMetricDataSource{
	// 				ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
	// 				ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
	// 				MetricName: to.Ptr("Requests"),
	// 			},
	// 			ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
	// 			Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
	// 			Threshold: to.Ptr[float64](3),
	// 			TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
	// 			WindowSize: to.Ptr("PT5M"),
	// 		},
	// 		IsEnabled: to.Ptr(true),
	// 		LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-23T21:23:52.022Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/deleteAlertRule.json
func ExampleAlertRulesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAlertRulesClient().Delete(ctx, "Rac46PostSwapRG", "chiricutin", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/getAlertRule.json
func ExampleAlertRulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAlertRulesClient().Get(ctx, "Rac46PostSwapRG", "chiricutin", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AlertRuleResource = armmonitor.AlertRuleResource{
	// 	Name: to.Ptr("chiricutin"),
	// 	Type: to.Ptr("Microsoft.Insights/alertRules"),
	// 	ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/chiricutin"),
	// 	Location: to.Ptr("West US"),
	// 	Tags: map[string]*string{
	// 		"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
	// 		"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
	// 	},
	// 	Properties: &armmonitor.AlertRule{
	// 		Name: to.Ptr("chiricutin"),
	// 		Description: to.Ptr("Pura Vida"),
	// 		Actions: []armmonitor.RuleActionClassification{
	// 		},
	// 		Condition: &armmonitor.ThresholdRuleCondition{
	// 			DataSource: &armmonitor.RuleMetricDataSource{
	// 				ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
	// 				ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
	// 				MetricName: to.Ptr("Requests"),
	// 			},
	// 			ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
	// 			Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
	// 			Threshold: to.Ptr[float64](3),
	// 			TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
	// 			WindowSize: to.Ptr("PT5M"),
	// 		},
	// 		IsEnabled: to.Ptr(true),
	// 		LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-23T21:23:52.022Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/patchAlertRule.json
func ExampleAlertRulesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAlertRulesClient().Update(ctx, "Rac46PostSwapRG", "chiricutin", armmonitor.AlertRuleResourcePatch{
		Properties: &armmonitor.AlertRule{
			Name:        to.Ptr("chiricutin"),
			Description: to.Ptr("Pura Vida"),
			Actions:     []armmonitor.RuleActionClassification{},
			Condition: &armmonitor.ThresholdRuleCondition{
				DataSource: &armmonitor.RuleMetricDataSource{
					ODataType:   to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
					ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
					MetricName:  to.Ptr("Requests"),
				},
				ODataType:       to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
				Operator:        to.Ptr(armmonitor.ConditionOperatorGreaterThan),
				Threshold:       to.Ptr[float64](3),
				TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
				WindowSize:      to.Ptr("PT5M"),
			},
			IsEnabled: to.Ptr(true),
		},
		Tags: map[string]*string{
			"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AlertRuleResource = armmonitor.AlertRuleResource{
	// 	Name: to.Ptr("chiricutin"),
	// 	Type: to.Ptr("Microsoft.Insights/alertRules"),
	// 	ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/chiricutin"),
	// 	Location: to.Ptr("West US"),
	// 	Tags: map[string]*string{
	// 		"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary"),
	// 	},
	// 	Properties: &armmonitor.AlertRule{
	// 		Name: to.Ptr("chiricutin"),
	// 		Description: to.Ptr("Pura Vida"),
	// 		Actions: []armmonitor.RuleActionClassification{
	// 		},
	// 		Condition: &armmonitor.ThresholdRuleCondition{
	// 			DataSource: &armmonitor.RuleMetricDataSource{
	// 				ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
	// 				ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
	// 				MetricName: to.Ptr("Requests"),
	// 			},
	// 			ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
	// 			Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
	// 			Threshold: to.Ptr[float64](3),
	// 			TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
	// 			WindowSize: to.Ptr("PT5M"),
	// 		},
	// 		IsEnabled: to.Ptr(true),
	// 		LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-23T21:23:52.022Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/listAlertRule.json
func ExampleAlertRulesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAlertRulesClient().NewListByResourceGroupPager("Rac46PostSwapRG", nil)
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
		// page.AlertRuleResourceCollection = armmonitor.AlertRuleResourceCollection{
		// 	Value: []*armmonitor.AlertRuleResource{
		// 		{
		// 			Name: to.Ptr("myRuleName"),
		// 			Type: to.Ptr("Microsoft.Insights/alertRules"),
		// 			ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/myRuleName"),
		// 			Location: to.Ptr("West US"),
		// 			Tags: map[string]*string{
		// 				"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
		// 				"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
		// 			},
		// 			Properties: &armmonitor.AlertRule{
		// 				Name: to.Ptr("myRuleName"),
		// 				Description: to.Ptr("Pura Vida"),
		// 				Actions: []armmonitor.RuleActionClassification{
		// 					&armmonitor.RuleEmailAction{
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleEmailAction"),
		// 						CustomEmails: []*string{
		// 							to.Ptr("gu@ms.com"),
		// 							to.Ptr("su@ms.net")},
		// 							SendToServiceOwners: to.Ptr(true),
		// 					}},
		// 					Condition: &armmonitor.ThresholdRuleCondition{
		// 						DataSource: &armmonitor.RuleMetricDataSource{
		// 							ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
		// 							ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
		// 							MetricName: to.Ptr("Requests"),
		// 						},
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
		// 						Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
		// 						Threshold: to.Ptr[float64](2),
		// 						TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
		// 						WindowSize: to.Ptr("PT5M"),
		// 					},
		// 					IsEnabled: to.Ptr(true),
		// 					LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-10T21:04:39.000Z"); return t}()),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("chiricutin0"),
		// 				Type: to.Ptr("Microsoft.Insights/alertRules"),
		// 				ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/chiricutin0"),
		// 				Location: to.Ptr("West US"),
		// 				Tags: map[string]*string{
		// 					"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
		// 					"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
		// 				},
		// 				Properties: &armmonitor.AlertRule{
		// 					Name: to.Ptr("chiricutin0"),
		// 					Description: to.Ptr("Pura Vida 0"),
		// 					Actions: []armmonitor.RuleActionClassification{
		// 					},
		// 					Condition: &armmonitor.ThresholdRuleCondition{
		// 						DataSource: &armmonitor.RuleMetricDataSource{
		// 							ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
		// 							ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
		// 							MetricName: to.Ptr("Requests"),
		// 						},
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
		// 						Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
		// 						Threshold: to.Ptr[float64](2),
		// 						TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
		// 						WindowSize: to.Ptr("PT5M"),
		// 					},
		// 					IsEnabled: to.Ptr(true),
		// 					LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-10T21:04:39.108Z"); return t}()),
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/listAlertRuleBySubscription.json
func ExampleAlertRulesClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmonitor.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAlertRulesClient().NewListBySubscriptionPager(nil)
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
		// page.AlertRuleResourceCollection = armmonitor.AlertRuleResourceCollection{
		// 	Value: []*armmonitor.AlertRuleResource{
		// 		{
		// 			Name: to.Ptr("myRuleName"),
		// 			Type: to.Ptr("Microsoft.Insights/alertRules"),
		// 			ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/myRuleName"),
		// 			Location: to.Ptr("West US"),
		// 			Tags: map[string]*string{
		// 				"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
		// 				"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
		// 			},
		// 			Properties: &armmonitor.AlertRule{
		// 				Name: to.Ptr("myRuleName"),
		// 				Description: to.Ptr("Pura Vida"),
		// 				Actions: []armmonitor.RuleActionClassification{
		// 					&armmonitor.RuleEmailAction{
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleEmailAction"),
		// 						CustomEmails: []*string{
		// 							to.Ptr("gu@ms.com"),
		// 							to.Ptr("su@ms.net")},
		// 							SendToServiceOwners: to.Ptr(true),
		// 					}},
		// 					Condition: &armmonitor.ThresholdRuleCondition{
		// 						DataSource: &armmonitor.RuleMetricDataSource{
		// 							ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
		// 							ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
		// 							MetricName: to.Ptr("Requests"),
		// 						},
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
		// 						Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
		// 						Threshold: to.Ptr[float64](2),
		// 						TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
		// 						WindowSize: to.Ptr("PT5M"),
		// 					},
		// 					IsEnabled: to.Ptr(true),
		// 					LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-10T21:04:39.000Z"); return t}()),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("chiricutin0"),
		// 				Type: to.Ptr("Microsoft.Insights/alertRules"),
		// 				ID: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/microsoft.insights/alertrules/chiricutin0"),
		// 				Location: to.Ptr("West US"),
		// 				Tags: map[string]*string{
		// 					"$type": to.Ptr("Microsoft.WindowsAzure.Management.Common.Storage.CasePreservedDictionary, Microsoft.WindowsAzure.Management.Common.Storage"),
		// 					"hidden-link:/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest": to.Ptr("Resource"),
		// 				},
		// 				Properties: &armmonitor.AlertRule{
		// 					Name: to.Ptr("chiricutin0"),
		// 					Description: to.Ptr("Pura Vida 0"),
		// 					Actions: []armmonitor.RuleActionClassification{
		// 					},
		// 					Condition: &armmonitor.ThresholdRuleCondition{
		// 						DataSource: &armmonitor.RuleMetricDataSource{
		// 							ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.RuleMetricDataSource"),
		// 							ResourceURI: to.Ptr("/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/Rac46PostSwapRG/providers/Microsoft.Web/sites/leoalerttest"),
		// 							MetricName: to.Ptr("Requests"),
		// 						},
		// 						ODataType: to.Ptr("Microsoft.Azure.Management.Insights.Models.ThresholdRuleCondition"),
		// 						Operator: to.Ptr(armmonitor.ConditionOperatorGreaterThan),
		// 						Threshold: to.Ptr[float64](2),
		// 						TimeAggregation: to.Ptr(armmonitor.TimeAggregationOperatorTotal),
		// 						WindowSize: to.Ptr("PT5M"),
		// 					},
		// 					IsEnabled: to.Ptr(true),
		// 					LastUpdatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2016-11-10T21:04:39.108Z"); return t}()),
		// 				},
		// 		}},
		// 	}
	}
}
