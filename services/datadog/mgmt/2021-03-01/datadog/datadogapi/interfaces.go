// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datadog/armdatadog](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datadog/armdatadog). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package datadogapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/Azure/go-autorest/autorest"
)

// MarketplaceAgreementsClientAPI contains the set of methods on the MarketplaceAgreementsClient type.
type MarketplaceAgreementsClientAPI interface {
	CreateOrUpdate(ctx context.Context, body *datadog.AgreementResource) (result datadog.AgreementResource, err error)
	List(ctx context.Context) (result datadog.AgreementResourceListResponsePage, err error)
	ListComplete(ctx context.Context) (result datadog.AgreementResourceListResponseIterator, err error)
}

var _ MarketplaceAgreementsClientAPI = (*datadog.MarketplaceAgreementsClient)(nil)

// MonitorsClientAPI contains the set of methods on the MonitorsClient type.
type MonitorsClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, monitorName string, body *datadog.MonitorResource) (result datadog.MonitorsCreateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitorsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitorResource, err error)
	GetDefaultKey(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.APIKey, err error)
	List(ctx context.Context) (result datadog.MonitorResourceListResponsePage, err error)
	ListComplete(ctx context.Context) (result datadog.MonitorResourceListResponseIterator, err error)
	ListAPIKeys(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.APIKeyListResponsePage, err error)
	ListAPIKeysComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.APIKeyListResponseIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result datadog.MonitorResourceListResponsePage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result datadog.MonitorResourceListResponseIterator, err error)
	ListHosts(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.HostListResponsePage, err error)
	ListHostsComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.HostListResponseIterator, err error)
	ListLinkedResources(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.LinkedResourceListResponsePage, err error)
	ListLinkedResourcesComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.LinkedResourceListResponseIterator, err error)
	ListMonitoredResources(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitoredResourceListResponsePage, err error)
	ListMonitoredResourcesComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitoredResourceListResponseIterator, err error)
	RefreshSetPasswordLink(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.SetPasswordLink, err error)
	SetDefaultKey(ctx context.Context, resourceGroupName string, monitorName string, body *datadog.APIKey) (result autorest.Response, err error)
	Update(ctx context.Context, resourceGroupName string, monitorName string, body *datadog.MonitorResourceUpdateParameters) (result datadog.MonitorsUpdateFuture, err error)
}

var _ MonitorsClientAPI = (*datadog.MonitorsClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result datadog.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result datadog.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*datadog.OperationsClient)(nil)

// TagRulesClientAPI contains the set of methods on the TagRulesClient type.
type TagRulesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, monitorName string, ruleSetName string, body *datadog.MonitoringTagRules) (result datadog.MonitoringTagRules, err error)
	Get(ctx context.Context, resourceGroupName string, monitorName string, ruleSetName string) (result datadog.MonitoringTagRules, err error)
	List(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitoringTagRulesListResponsePage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.MonitoringTagRulesListResponseIterator, err error)
}

var _ TagRulesClientAPI = (*datadog.TagRulesClient)(nil)

// SingleSignOnConfigurationsClientAPI contains the set of methods on the SingleSignOnConfigurationsClient type.
type SingleSignOnConfigurationsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, monitorName string, configurationName string, body *datadog.SingleSignOnResource) (result datadog.SingleSignOnConfigurationsCreateOrUpdateFuture, err error)
	Get(ctx context.Context, resourceGroupName string, monitorName string, configurationName string) (result datadog.SingleSignOnResource, err error)
	List(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.SingleSignOnResourceListResponsePage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, monitorName string) (result datadog.SingleSignOnResourceListResponseIterator, err error)
}

var _ SingleSignOnConfigurationsClientAPI = (*datadog.SingleSignOnConfigurationsClient)(nil)
