// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcontainerregistry

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	internal       *arm.Client
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	internal, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: subscriptionID,
		internal:       internal,
	}, nil
}

// NewAgentPoolsClient creates a new instance of AgentPoolsClient.
func (c *ClientFactory) NewAgentPoolsClient() *AgentPoolsClient {
	return &AgentPoolsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewArchiveVersionsClient creates a new instance of ArchiveVersionsClient.
func (c *ClientFactory) NewArchiveVersionsClient() *ArchiveVersionsClient {
	return &ArchiveVersionsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewArchivesClient creates a new instance of ArchivesClient.
func (c *ClientFactory) NewArchivesClient() *ArchivesClient {
	return &ArchivesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewCacheRulesClient creates a new instance of CacheRulesClient.
func (c *ClientFactory) NewCacheRulesClient() *CacheRulesClient {
	return &CacheRulesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewConnectedRegistriesClient creates a new instance of ConnectedRegistriesClient.
func (c *ClientFactory) NewConnectedRegistriesClient() *ConnectedRegistriesClient {
	return &ConnectedRegistriesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewCredentialSetsClient creates a new instance of CredentialSetsClient.
func (c *ClientFactory) NewCredentialSetsClient() *CredentialSetsClient {
	return &CredentialSetsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewExportPipelinesClient creates a new instance of ExportPipelinesClient.
func (c *ClientFactory) NewExportPipelinesClient() *ExportPipelinesClient {
	return &ExportPipelinesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewImportPipelinesClient creates a new instance of ImportPipelinesClient.
func (c *ClientFactory) NewImportPipelinesClient() *ImportPipelinesClient {
	return &ImportPipelinesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	return &OperationsClient{
		internal: c.internal,
	}
}

// NewPipelineRunsClient creates a new instance of PipelineRunsClient.
func (c *ClientFactory) NewPipelineRunsClient() *PipelineRunsClient {
	return &PipelineRunsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPrivateEndpointConnectionsClient creates a new instance of PrivateEndpointConnectionsClient.
func (c *ClientFactory) NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient {
	return &PrivateEndpointConnectionsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRegistriesClient creates a new instance of RegistriesClient.
func (c *ClientFactory) NewRegistriesClient() *RegistriesClient {
	return &RegistriesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewReplicationsClient creates a new instance of ReplicationsClient.
func (c *ClientFactory) NewReplicationsClient() *ReplicationsClient {
	return &ReplicationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRunsClient creates a new instance of RunsClient.
func (c *ClientFactory) NewRunsClient() *RunsClient {
	return &RunsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewScopeMapsClient creates a new instance of ScopeMapsClient.
func (c *ClientFactory) NewScopeMapsClient() *ScopeMapsClient {
	return &ScopeMapsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewTaskRunsClient creates a new instance of TaskRunsClient.
func (c *ClientFactory) NewTaskRunsClient() *TaskRunsClient {
	return &TaskRunsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewTasksClient creates a new instance of TasksClient.
func (c *ClientFactory) NewTasksClient() *TasksClient {
	return &TasksClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewTokensClient creates a new instance of TokensClient.
func (c *ClientFactory) NewTokensClient() *TokensClient {
	return &TokensClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebhooksClient creates a new instance of WebhooksClient.
func (c *ClientFactory) NewWebhooksClient() *WebhooksClient {
	return &WebhooksClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}
