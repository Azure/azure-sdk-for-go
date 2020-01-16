package datafactoryapi

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
	"github.com/Azure/azure-sdk-for-go/services/preview/datafactory/mgmt/2017-09-01-preview/datafactory"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result datafactory.OperationListResponse, err error)
}

var _ OperationsClientAPI = (*datafactory.OperationsClient)(nil)

// FactoriesClientAPI contains the set of methods on the FactoriesClient type.
type FactoriesClientAPI interface {
	CancelPipelineRun(ctx context.Context, resourceGroupName string, factoryName string, runID string) (result autorest.Response, err error)
	ConfigureFactoryRepo(ctx context.Context, locationID string, factoryRepoUpdate datafactory.FactoryRepoUpdate) (result datafactory.Factory, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, factory datafactory.Factory) (result datafactory.Factory, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.Factory, err error)
	List(ctx context.Context) (result datafactory.FactoryListResponsePage, err error)
	ListComplete(ctx context.Context) (result datafactory.FactoryListResponseIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result datafactory.FactoryListResponsePage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result datafactory.FactoryListResponseIterator, err error)
	Update(ctx context.Context, resourceGroupName string, factoryName string, factoryUpdateParameters datafactory.FactoryUpdateParameters) (result datafactory.Factory, err error)
}

var _ FactoriesClientAPI = (*datafactory.FactoriesClient)(nil)

// IntegrationRuntimesClientAPI contains the set of methods on the IntegrationRuntimesClient type.
type IntegrationRuntimesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, integrationRuntime datafactory.IntegrationRuntimeResource, ifMatch string) (result datafactory.IntegrationRuntimeResource, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimeResource, err error)
	GetConnectionInfo(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimeConnectionInfo, err error)
	GetMonitoringData(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimeMonitoringData, err error)
	GetStatus(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimeStatusResponse, err error)
	ListAuthKeys(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimeAuthKeys, err error)
	ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.IntegrationRuntimeListResponsePage, err error)
	ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.IntegrationRuntimeListResponseIterator, err error)
	RegenerateAuthKey(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, regenerateKeyParameters datafactory.IntegrationRuntimeRegenerateKeyParameters) (result datafactory.IntegrationRuntimeAuthKeys, err error)
	RemoveNode(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, removeNodeParameters datafactory.IntegrationRuntimeRemoveNodeRequest) (result autorest.Response, err error)
	Start(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimesStartFuture, err error)
	Stop(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result datafactory.IntegrationRuntimesStopFuture, err error)
	SyncCredentials(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result autorest.Response, err error)
	Update(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, updateIntegrationRuntimeRequest datafactory.UpdateIntegrationRuntimeRequest) (result datafactory.IntegrationRuntimeStatusResponse, err error)
	Upgrade(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string) (result autorest.Response, err error)
}

var _ IntegrationRuntimesClientAPI = (*datafactory.IntegrationRuntimesClient)(nil)

// IntegrationRuntimeNodesClientAPI contains the set of methods on the IntegrationRuntimeNodesClient type.
type IntegrationRuntimeNodesClientAPI interface {
	Delete(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string) (result autorest.Response, err error)
	GetIPAddress(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string) (result datafactory.IntegrationRuntimeNodeIPAddress, err error)
	Update(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string, updateIntegrationRuntimeNodeRequest datafactory.UpdateIntegrationRuntimeNodeRequest) (result datafactory.SelfHostedIntegrationRuntimeNode, err error)
}

var _ IntegrationRuntimeNodesClientAPI = (*datafactory.IntegrationRuntimeNodesClient)(nil)

// LinkedServicesClientAPI contains the set of methods on the LinkedServicesClient type.
type LinkedServicesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string, linkedService datafactory.LinkedServiceResource, ifMatch string) (result datafactory.LinkedServiceResource, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string) (result datafactory.LinkedServiceResource, err error)
	ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.LinkedServiceListResponsePage, err error)
	ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.LinkedServiceListResponseIterator, err error)
}

var _ LinkedServicesClientAPI = (*datafactory.LinkedServicesClient)(nil)

// DatasetsClientAPI contains the set of methods on the DatasetsClient type.
type DatasetsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, datasetName string, dataset datafactory.DatasetResource, ifMatch string) (result datafactory.DatasetResource, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string, datasetName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string, datasetName string) (result datafactory.DatasetResource, err error)
	ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.DatasetListResponsePage, err error)
	ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.DatasetListResponseIterator, err error)
}

var _ DatasetsClientAPI = (*datafactory.DatasetsClient)(nil)

// PipelinesClientAPI contains the set of methods on the PipelinesClient type.
type PipelinesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, pipeline datafactory.PipelineResource, ifMatch string) (result datafactory.PipelineResource, err error)
	CreateRun(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string, parameters map[string]interface{}) (result datafactory.CreateRunResponse, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string, pipelineName string) (result datafactory.PipelineResource, err error)
	ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.PipelineListResponsePage, err error)
	ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.PipelineListResponseIterator, err error)
}

var _ PipelinesClientAPI = (*datafactory.PipelinesClient)(nil)

// PipelineRunsClientAPI contains the set of methods on the PipelineRunsClient type.
type PipelineRunsClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, factoryName string, runID string) (result datafactory.PipelineRun, err error)
	QueryByFactory(ctx context.Context, resourceGroupName string, factoryName string, filterParameters datafactory.PipelineRunFilterParameters) (result datafactory.PipelineRunQueryResponse, err error)
}

var _ PipelineRunsClientAPI = (*datafactory.PipelineRunsClient)(nil)

// ActivityRunsClientAPI contains the set of methods on the ActivityRunsClient type.
type ActivityRunsClientAPI interface {
	ListByPipelineRun(ctx context.Context, resourceGroupName string, factoryName string, runID string, startTime date.Time, endTime date.Time, status string, activityName string, linkedServiceName string) (result datafactory.ActivityRunsListResponsePage, err error)
	ListByPipelineRunComplete(ctx context.Context, resourceGroupName string, factoryName string, runID string, startTime date.Time, endTime date.Time, status string, activityName string, linkedServiceName string) (result datafactory.ActivityRunsListResponseIterator, err error)
}

var _ ActivityRunsClientAPI = (*datafactory.ActivityRunsClient)(nil)

// TriggersClientAPI contains the set of methods on the TriggersClient type.
type TriggersClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, triggerName string, trigger datafactory.TriggerResource, ifMatch string) (result datafactory.TriggerResource, err error)
	Delete(ctx context.Context, resourceGroupName string, factoryName string, triggerName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, factoryName string, triggerName string) (result datafactory.TriggerResource, err error)
	ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.TriggerListResponsePage, err error)
	ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result datafactory.TriggerListResponseIterator, err error)
	ListRuns(ctx context.Context, resourceGroupName string, factoryName string, triggerName string, startTime date.Time, endTime date.Time) (result datafactory.TriggerRunListResponsePage, err error)
	ListRunsComplete(ctx context.Context, resourceGroupName string, factoryName string, triggerName string, startTime date.Time, endTime date.Time) (result datafactory.TriggerRunListResponseIterator, err error)
	Start(ctx context.Context, resourceGroupName string, factoryName string, triggerName string) (result datafactory.TriggersStartFuture, err error)
	Stop(ctx context.Context, resourceGroupName string, factoryName string, triggerName string) (result datafactory.TriggersStopFuture, err error)
}

var _ TriggersClientAPI = (*datafactory.TriggersClient)(nil)
